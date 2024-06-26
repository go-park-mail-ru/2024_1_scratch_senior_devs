package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/loadtls"

	metricsmw "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/metrics"

	grpcNote "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc"
	generatedNote "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"

	noteRepo "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/repo"
	noteUsecase "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/usecase"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/olivere/elastic/v7"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() (err error) {
	logFile, err := os.OpenFile(os.Getenv("NOTE_LOG_FILE"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("error opening log file: " + err.Error())
		return
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(logFile, os.Stdout), &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.LoadConfig(os.Getenv("CONFIG_FILE"), logger)

	db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("error connecting to postgres: " + err.Error())
		return
	}
	defer db.Close()

	elasticClient, err := elastic.NewClient(elastic.SetURL(os.Getenv("ELASTIC_URL")))
	if err != nil {
		logger.Error("error connecting to elasticsearch: " + err.Error())
		return
	}

	tlsCredentials, err := loadtls.LoadTLSCredentials(cfg.Grpc.NoteIP)
	if err != nil {
		logger.Error("fail to load TLS credentials" + err.Error())
		return
	}

	postgresMetrics, err := metrics.NewDatabaseMetrics("postgres", "note")
	if err != nil {
		logger.Error("can`t create metrics (note postgres): " + err.Error())
	}

	elasticMetrics, err := metrics.NewDatabaseMetrics("elastic", "note")
	if err != nil {
		logger.Error("can`t create metrics (note elastic): " + err.Error())
	}

	grpcMetrics, err := metrics.NewGrpcMetrics("note")
	if err != nil {
		logger.Error("can`t create metrics (note grpc): " + err.Error())
	}

	NoteBaseRepo := noteRepo.CreateNotePostgres(db, &postgresMetrics)
	NoteSearchRepo := noteRepo.CreateNoteElastic(elasticClient, cfg.Elastic, &elasticMetrics)

	NoteUsecase := noteUsecase.CreateNoteUsecase(NoteBaseRepo, NoteSearchRepo, cfg.Elastic, cfg.Constraints, &sync.WaitGroup{})
	NoteDelivery := grpcNote.NewGrpcNoteHandler(NoteUsecase)

	MetricsMiddleware := metricsmw.NewGrpcMw(grpcMetrics)
	LogMiddleware := log.NewGrpcLogMw(logger)

	gRPCServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.ChainUnaryInterceptor(MetricsMiddleware.ServerMetricsInterceptor, LogMiddleware.ServerLogsInterceptor),
	)
	generatedNote.RegisterNoteServer(gRPCServer, NoteDelivery)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.PathPrefix("/metrics").Handler(promhttp.Handler())
	http.Handle("/", r)
	httpSrv := http.Server{Handler: r, Addr: fmt.Sprintf(":%s", cfg.Grpc.NoteMetricsPort)}
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			logger.Error("fail httpSrv.ListenAndServe")
		}
	}()

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Grpc.NotePort))
		if err != nil {
			logger.Error(err.Error())
		}
		if err := gRPCServer.Serve(listener); err != nil {
			logger.Error(err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	gRPCServer.GracefulStop()
	return nil
}

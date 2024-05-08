package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/loadtls"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpcAuth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc"
	generatedAuth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc/gen"
	authRepo "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/repo"
	authUsecase "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/usecase"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/log"
	metricsmw "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/metrics"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/jackc/pgx/v4/pgxpool"
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
	logFile, err := os.OpenFile(os.Getenv("AUTH_LOG_FILE"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

	tlsCredentials, err := loadtls.LoadTLSCredentials(cfg.Grpc.AuthIP)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	postgresMetrics, err := metrics.NewDatabaseMetrics("auth postgres")
	if err != nil {
		logger.Error("cant create metrics")
	}

	AuthRepo := authRepo.CreateAuthRepo(db, &postgresMetrics)
	AuthUsecase := authUsecase.CreateAuthUsecase(AuthRepo, cfg.AuthUsecase, cfg.Validation)
	AuthDelivery := grpcAuth.NewGrpcAuthHandler(AuthUsecase)

	grpcMetrics, err := metrics.NewGrpcMetrics("auth")
	if err != nil {
		logger.Error("cant create metrics")
	}
	metricsMw := metricsmw.NewGrpcMw(*grpcMetrics)
	logMw := log.NewGrpcLogMw(logger)
	gRPCServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.ChainUnaryInterceptor(metricsMw.ServerMetricsInterceptor, logMw.ServerLogsInterceptor),
	)
	generatedAuth.RegisterAuthServer(gRPCServer, AuthDelivery)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.PathPrefix("/metrics").Handler(promhttp.Handler())
	http.Handle("/", r)
	httpSrv := http.Server{Handler: r, Addr: fmt.Sprintf(":%d", 7071)}
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			logger.Error("fail httpSrv.ListenAndServe")
		}
	}()
	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Grpc.AuthPort))
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

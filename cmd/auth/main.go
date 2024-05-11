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
	"syscall"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	grpcAuth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc"
	generatedAuth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc/gen"

	authRepo "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/repo"
	authUsecase "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/usecase"

	metricsmw "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/metrics"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/loadtls"
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

	redisOpts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		logger.Error("error connecting to redis: " + err.Error())
		return
	}
	redisDB := redis.NewClient(redisOpts)

	tlsCredentials, err := loadtls.LoadTLSCredentials(cfg.Grpc.AuthIP)
	if err != nil {
		logger.Error("fail to load TLS credentials" + err.Error())
		return
	}

	postgresMetrics, err := metrics.NewDatabaseMetrics("postgres", "auth")
	if err != nil {
		logger.Error("can`t create metrics (auth postgres): " + err.Error())
	}

	redisMetrics, err := metrics.NewDatabaseMetrics("redis", "main")
	if err != nil {
		logger.Error("can`t create metrics (main redis): " + err.Error())
	}

	grpcMetrics, err := metrics.NewGrpcMetrics("auth")
	if err != nil {
		logger.Error("can`t create metrics (auth grpc): " + err.Error())
	}

	BlockerRepo := authRepo.CreateBlockerRepo(*redisDB, cfg.Blocker, &redisMetrics)
	BlockerUsecase := authUsecase.CreateBlockerUsecase(BlockerRepo, cfg.Blocker)

	AuthRepo := authRepo.CreateAuthRepo(db, &postgresMetrics)
	AuthUsecase := authUsecase.CreateAuthUsecase(AuthRepo, cfg.AuthUsecase, cfg.Validation)
	AuthDelivery := grpcAuth.NewGrpcAuthHandler(AuthUsecase, BlockerUsecase)

	MetricsMiddleware := metricsmw.NewGrpcMw(grpcMetrics)
	LogMiddleware := log.NewGrpcLogMw(logger)

	gRPCServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.ChainUnaryInterceptor(MetricsMiddleware.ServerMetricsInterceptor, LogMiddleware.ServerLogsInterceptor),
	)
	generatedAuth.RegisterAuthServer(gRPCServer, AuthDelivery)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.PathPrefix("/metrics").Handler(promhttp.Handler())
	http.Handle("/", r)
	httpSrv := http.Server{Handler: r, Addr: fmt.Sprintf(":%s", cfg.Grpc.AuthMetricsPort)}
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

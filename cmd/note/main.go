package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	grpcNote "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc"
	generatedNote "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"
	noteRepo "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/repo"
	noteUsecase "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/usecase"

	"github.com/joho/godotenv"
	"github.com/olivere/elastic/v7"

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
	logFile, err := os.OpenFile(os.Getenv("MAIN_LOG_FILE"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

	NoteBaseRepo := noteRepo.CreateNotePostgres(db)
	NoteSearchRepo := noteRepo.CreateNoteElastic(elasticClient, cfg.Elastic)

	NoteUsecase := noteUsecase.CreateNoteUsecase(NoteBaseRepo, NoteSearchRepo, cfg.Elastic, &sync.WaitGroup{})
	NoteDelivery := grpcNote.NewGrpcNoteHandler(NoteUsecase)

	gRPCServer := grpc.NewServer()
	generatedNote.RegisterNoteServer(gRPCServer, NoteDelivery)

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

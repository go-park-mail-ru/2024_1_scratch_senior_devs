package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/olivere/elastic/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/redis/go-redis/v9"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/path"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/protection"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/recover"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	grpcAuth "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc/gen"

	authDelivery "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/http"
	authRepo "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/repo"
	authUsecase "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/usecase"
	_ "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/docs"

	noteDelivery "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/http"
	noteRepo "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/repo"
	noteUsecase "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/usecase"

	attachDelivery "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/attach/delivery/http"
	attachRepo "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/attach/repo"
	attachUsecase "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/attach/usecase"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
}

// @title 			YouNote API
// @version 		1.0
// @description 	API for YouNote service
// @host 			you-note.ru
func main() {
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

	redisOpts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		logger.Error("error connecting to redis: " + err.Error())
		return
	}
	redisDB := redis.NewClient(redisOpts)

	elasticClient, err := elastic.NewClient(elastic.SetURL(os.Getenv("ELASTIC_URL")))
	if err != nil {
		logger.Error("error connecting to elasticsearch: " + err.Error())
		return
	}

	authConn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Grpc.AuthIP, cfg.Grpc.AuthPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("fail grpc.Dial auth: " + err.Error())
		return
	}
	defer authConn.Close()

	JwtMiddleware := protection.CreateJwtMiddleware(cfg.AuthHandler.Jwt)
	CsrfMiddleware := protection.CreateCsrfMiddleware(cfg.AuthHandler.Csrf)

	logMW := log.CreateLogMiddleware(logger)

	NoteBaseRepo := noteRepo.CreateNotePostgres(db)
	NoteSearchRepo := noteRepo.CreateNoteElastic(elasticClient, cfg.Elastic)
	NoteUsecase := noteUsecase.CreateNoteUsecase(NoteBaseRepo, NoteSearchRepo, cfg.Elastic, &sync.WaitGroup{})
	NoteDelivery := noteDelivery.CreateNotesHandler(NoteUsecase)

	BlockerRepo := authRepo.CreateBlockerRepo(*redisDB, cfg.Blocker)
	BlockerUsecase := authUsecase.CreateBlockerUsecase(BlockerRepo, cfg.Blocker)

	AuthClient := grpcAuth.NewAuthClient(authConn)
	AuthDelivery := authDelivery.CreateAuthHandler(AuthClient, BlockerUsecase, NoteUsecase, cfg.AuthHandler, cfg.Validation)

	AttachRepo := attachRepo.CreateAttachRepo(db)
	AttachUsecase := attachUsecase.CreateAttachUsecase(AttachRepo, NoteBaseRepo)
	AttachDelivery := attachDelivery.CreateAttachHandler(AttachUsecase, cfg.Attach)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	r.Use(
		logMW,
		protection.CorsMiddleware,
		recover.RecoverMiddleware,
	)

	r.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet, http.MethodOptions)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.Handle("/signup", http.HandlerFunc(AuthDelivery.SignUp)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/login", http.HandlerFunc(AuthDelivery.SignIn)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/logout", JwtMiddleware(CsrfMiddleware(http.HandlerFunc(AuthDelivery.LogOut)))).Methods(http.MethodDelete, http.MethodOptions)
		auth.Handle("/check_user", JwtMiddleware(http.HandlerFunc(AuthDelivery.CheckUser))).Methods(http.MethodGet, http.MethodOptions)
		auth.Handle("/get_qr", JwtMiddleware(http.HandlerFunc(AuthDelivery.GetQRCode))).Methods(http.MethodGet, http.MethodOptions)
		auth.Handle("/disable_2fa", JwtMiddleware(CsrfMiddleware(http.HandlerFunc(AuthDelivery.DisableSecondFactor)))).Methods(http.MethodDelete, http.MethodOptions)
	}

	note := r.PathPrefix("/note").Subrouter()
	note.Use(protection.ReadAndCloseBody, JwtMiddleware, CsrfMiddleware)
	{
		note.Handle("/get_all", http.HandlerFunc(NoteDelivery.GetAllNotes)).Methods(http.MethodGet, http.MethodOptions)
		note.Handle("/{id}", http.HandlerFunc(NoteDelivery.GetNote)).Methods(http.MethodGet, http.MethodOptions)
		note.Handle("/add", http.HandlerFunc(NoteDelivery.AddNote)).Methods(http.MethodPost, http.MethodOptions)
		note.Handle("/{id}/edit", http.HandlerFunc(NoteDelivery.UpdateNote)).Methods(http.MethodPost, http.MethodOptions)
		note.Handle("/{id}/delete", http.HandlerFunc(NoteDelivery.DeleteNote)).Methods(http.MethodDelete, http.MethodOptions)
		note.Handle("/{id}/add_attach", http.HandlerFunc(AttachDelivery.AddAttach)).Methods(http.MethodPost, http.MethodOptions)
	}

	profile := r.PathPrefix("/profile").Subrouter()
	profile.Use(protection.ReadAndCloseBody, JwtMiddleware, CsrfMiddleware)
	{
		profile.Handle("/get", http.HandlerFunc(AuthDelivery.GetProfile)).Methods(http.MethodGet, http.MethodOptions)
		profile.Handle("/update", http.HandlerFunc(AuthDelivery.UpdateProfile)).Methods(http.MethodPost, http.MethodOptions)
		profile.Handle("/update_avatar", http.HandlerFunc(AuthDelivery.UpdateProfileAvatar)).Methods(http.MethodPost, http.MethodOptions)
	}

	attach := r.PathPrefix("/attach").Subrouter()
	attach.Use(JwtMiddleware, CsrfMiddleware)
	{
		attach.Handle("/{id}", http.HandlerFunc(AttachDelivery.GetAttach)).Methods(http.MethodGet, http.MethodOptions)
		attach.Handle("/{id}/delete", http.HandlerFunc(AttachDelivery.DeleteAttach)).Methods(http.MethodDelete, http.MethodOptions)
	}

	http.Handle("/", r)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	server := http.Server{
		Handler:           path.PathMiddleware(r),
		Addr:              ":8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Info("Server stopped")
		}
	}()
	logger.Info("Server started")

	sig := <-signalCh
	logger.Info("Received signal: " + sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed: " + err.Error())
	}
}

package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	err := godotenv.Load()
	if err != nil {
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
		logger.Info("error connecting to postgres: " + err.Error())
		return
	}
	defer db.Close()

	redisOpts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		logger.Info("error connecting to redis: " + err.Error())
		return
	}
	redisDB := redis.NewClient(redisOpts)

	JwtMiddleware := protection.CreateJwtMiddleware(logger, cfg.AuthHandler.Jwt)
	CsrfMiddleware := protection.CreateCsrfMiddleware(logger, cfg.AuthHandler.Csrf)
	RecoverMiddleware := recover.CreateRecoverMiddleware(logger)

	BlockerRepo := authRepo.CreateBlockerRepo(*redisDB, logger, cfg.Blocker)
	BlockerUsecase := authUsecase.CreateBlockerUsecase(BlockerRepo, logger, cfg.Blocker)

	AuthRepo := authRepo.CreateAuthRepo(db, logger)
	AuthUsecase := authUsecase.CreateAuthUsecase(AuthRepo, logger, cfg.AuthUsecase, cfg.UserValidation)
	AuthDelivery := authDelivery.CreateAuthHandler(AuthUsecase, BlockerUsecase, logger, cfg.AuthHandler, cfg.UserValidation)

	NoteRepo := noteRepo.CreateNoteRepo(db, logger)
	NoteUsecase := noteUsecase.CreateNoteUsecase(NoteRepo, logger)
	NoteDelivery := noteDelivery.CreateNotesHandler(NoteUsecase, logger)

	AttachRepo := attachRepo.CreateAttachRepo(db, logger)
	AttachUsecase := attachUsecase.CreateAttachUsecase(AttachRepo, NoteRepo, logger)
	AttachDelivery := attachDelivery.CreateAttachHandler(AttachUsecase, logger, cfg.Attach)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	r.Use(
		protection.CorsMiddleware,
		log.LogMiddleware,
		RecoverMiddleware,
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

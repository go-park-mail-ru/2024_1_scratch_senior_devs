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

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/cors"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/jwt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/path"
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
// @host 			you-note.ru:8080
func main() {
	db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	logFile, err := os.OpenFile(os.Getenv("MAIN_LOG_FILE"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(logFile, os.Stdout), &slog.HandlerOptions{Level: slog.LevelInfo}))

	JwtMiddleware := jwt.CreateJwtMiddleware(logger)
	RecoverMiddleware := recover.CreateRecoverMiddleware(logger)

	AuthRepo := authRepo.CreateAuthRepo(db, logger)
	AuthUsecase := authUsecase.CreateAuthUsecase(AuthRepo, logger)
	AuthDelivery := authDelivery.CreateAuthHandler(AuthUsecase, logger)

	NoteRepo := noteRepo.CreateNoteRepo(db, logger)
	NoteUsecase := noteUsecase.CreateNoteUsecase(NoteRepo, logger)
	NoteDelivery := noteDelivery.CreateNotesHandler(NoteUsecase, logger)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	r.Use(cors.CorsMiddleware, log.LogMiddleware, RecoverMiddleware)

	r.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet, http.MethodOptions)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.Handle("/signup", http.HandlerFunc(AuthDelivery.SignUp)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/login", http.HandlerFunc(AuthDelivery.SignIn)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/logout", JwtMiddleware(http.HandlerFunc(AuthDelivery.LogOut))).Methods(http.MethodDelete, http.MethodOptions)
		auth.Handle("/check_user", JwtMiddleware(http.HandlerFunc(AuthDelivery.CheckUser))).Methods(http.MethodGet, http.MethodOptions)
		auth.Handle("/get_qr", JwtMiddleware(http.HandlerFunc(AuthDelivery.GetQRCode))).Methods(http.MethodGet, http.MethodOptions)
	}

	note := r.PathPrefix("/note").Subrouter()
	note.Use(JwtMiddleware)
	{
		note.Handle("/get_all", http.HandlerFunc(NoteDelivery.GetAllNotes)).Methods(http.MethodGet, http.MethodOptions)
		note.Handle("/{id}", http.HandlerFunc(NoteDelivery.GetNote)).Methods(http.MethodGet, http.MethodOptions)
		note.Handle("/add", http.HandlerFunc(NoteDelivery.AddNote)).Methods(http.MethodPost, http.MethodOptions)
		note.Handle("/{id}/edit", http.HandlerFunc(NoteDelivery.UpdateNote)).Methods(http.MethodPost, http.MethodOptions)
		note.Handle("/{id}/delete", http.HandlerFunc(NoteDelivery.DeleteNote)).Methods(http.MethodDelete, http.MethodOptions)
	}

	profile := r.PathPrefix("/profile").Subrouter()
	profile.Use(jwt.ReadAndCloseBody, JwtMiddleware)
	{
		profile.Handle("/get", http.HandlerFunc(AuthDelivery.GetProfile)).Methods(http.MethodGet, http.MethodOptions)
		profile.Handle("/update", http.HandlerFunc(AuthDelivery.UpdateProfile)).Methods(http.MethodPost, http.MethodOptions)
		profile.Handle("/update_avatar", http.HandlerFunc(AuthDelivery.UpdateProfileAvatar)).Methods(http.MethodPost, http.MethodOptions)
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

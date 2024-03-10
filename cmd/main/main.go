package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/docs"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware"

	authDelivery "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/http"
	authRepo "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/repo"
	authUsecase "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/usecase"

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

	logger := slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))

	AuthRepo := authRepo.CreateAuthRepo(db)
	AuthUsecase := authUsecase.CreateAuthUsecase(AuthRepo)
	AuthDelivery := authDelivery.CreateAuthHandler(AuthUsecase, logger)

	NoteRepo := noteRepo.CreateNoteRepo(db, logger)
	NoteUsecase := noteUsecase.CreateNoteUsecase(NoteRepo)
	NoteDelivery := noteDelivery.CreateNotesHandler(NoteUsecase, logger)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	swagger := r.PathPrefix("/swagger").Subrouter()
	swagger.Use(middleware.CorsMiddlewareForSwagger)

	r.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet, http.MethodOptions)

	auth := r.PathPrefix("/auth").Subrouter()
	auth.Use(middleware.CorsMiddleware)
	{
		auth.Handle("/signup", http.HandlerFunc(AuthDelivery.SignUp)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/login", http.HandlerFunc(AuthDelivery.SignIn)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/logout", middleware.JwtMiddleware(http.HandlerFunc(AuthDelivery.LogOut))).Methods(http.MethodDelete, http.MethodOptions) // POST ?
		auth.Handle("/check_user", middleware.JwtMiddleware(http.HandlerFunc(AuthDelivery.CheckUser))).Methods(http.MethodGet, http.MethodOptions)
	}

	note := r.PathPrefix("/note").Subrouter()
	note.Use(middleware.CorsMiddleware, middleware.JwtMiddleware)
	{
		note.Handle("/get_all", http.HandlerFunc(NoteDelivery.GetAllNotes)).Methods(http.MethodGet, http.MethodOptions)
		note.Handle("/{id}", http.HandlerFunc(NoteDelivery.GetNote)).Methods(http.MethodGet, http.MethodOptions)
		note.Handle("/add", http.HandlerFunc(NoteDelivery.AddNote)).Methods(http.MethodPost, http.MethodOptions)
	}

	http.Handle("/", r)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	srv := http.Server{
		Handler:           r,
		Addr:              ":8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("Server stopped")
		}
	}()
	fmt.Println("Server started")

	sig := <-signalCh
	fmt.Printf("Received signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown failed: %v\n", err)
	}
}

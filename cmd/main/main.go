package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/authmw"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/authmw"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	authDelivery "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/http"
	authRepo "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/repo"
	authUsecase "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/usecase"

	noteDelivery "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/http"
	noteRepo "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/repo"
	noteUsecase "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/usecase"
)

func main() {
	db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	AuthRepo := authRepo.CreateAuthRepo(db)
	AuthUsecase := authUsecase.CreateAuthUsecase(*AuthRepo)
	AuthDelivery := authDelivery.CreateAuthHandler(AuthUsecase)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.Handle("/signup", http.HandlerFunc(AuthDelivery.SignUp)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/login", http.HandlerFunc(AuthDelivery.SignIn)).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/logout", authmw.JwtMiddleware(http.HandlerFunc(AuthDelivery.LogOut))).Methods(http.MethodPost, http.MethodOptions)
		auth.Handle("/check_user", authmw.JwtMiddleware(http.HandlerFunc(AuthDelivery.CheckUser))).Methods(http.MethodPost, http.MethodOptions)
	}
	NoteRepo := noteRepo.CreateNotesRepo(db)
	NoteUsecase := noteUsecase.CreateNotesUsecase(*NoteRepo)
	NoteDelivery := noteDelivery.CreateNotesHandler(*NoteUsecase)

	note := r.PathPrefix("/note").Subrouter()
	note.Use(authmw.JwtMiddleware)
	{
		note.Handle("/get_all", http.HandlerFunc(NoteDelivery.GetAllNotes)).Methods(http.MethodGet, http.MethodOptions)

	}

	http.Handle("/", r)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	srv := http.Server{
		Handler:           r,
		Addr:              "127.0.0.1:8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println("error starting server")
		}
	}()

	// Ловим сигнал
	sig := <-signalCh
	log.Printf("Received signal: %v\n", sig)

	// Дальше начинаем обратный отсчёт
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Гасим сервер с обязательством закончить за 30 секунд всю работу
	// time.Sleep(time.Second * 5)
	log.Println("i am here")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}
}

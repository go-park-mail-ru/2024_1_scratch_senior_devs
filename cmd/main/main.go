package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/authmw"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	authDelivery "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/http"
	authRepo "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/repo"
	authUsecase "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/usecase"
)

func main() {
	db, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		"", // database username
		"", // postgres password
		"", // postgres host
		"", // postgres port
		"", // database name
	))
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

	http.Handle("/", r)

	srv := http.Server{
		Handler:           r,
		Addr:              "127.0.0.1:8080",
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	srv.ListenAndServe()
}

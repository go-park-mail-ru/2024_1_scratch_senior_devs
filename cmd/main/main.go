package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	})

	http.Handle("/", r)

	srv := http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
	}

	srv.ListenAndServe()
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to create database pool: %v", err)
	}
	defer pool.Close()

	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok"}`))
		})

		r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if err := pool.Ping(r.Context()); err != nil {
				log.Printf("database ping failed: %v", err)
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"status":"not_ready"}`))
				return
			}

			w.Write([]byte(`{"status":"ok"}`))
		})
	})

	port := ":8080"
	log.Printf("Server starting on port %s...", port)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

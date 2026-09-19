package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	})

	port := ":8080"
	log.Printf("Server starting on port %s...", port)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

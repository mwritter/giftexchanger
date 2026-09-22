package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mwritter/giftswap/services/api/internal/auth"
	"github.com/mwritter/giftswap/services/api/internal/server"
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

	authService := auth.NewService(pool, auth.LogMailer{}, auth.Config{
		BaseURL:      envOr("APP_BASE_URL", "http://localhost:3000"),
		EntryURL:     envOr("APP_ENTRY_URL", "http://localhost:3000/dashboard"),
		ErrorURL:     envOr("APP_ERROR_URL", "http://localhost:3000/auth/error"),
		MagicLinkTTL: durationOr("MAGIC_LINK_TTL", 15*time.Minute),
		SessionTTL:   durationOr("SESSION_TTL", 30*24*time.Hour),
		CookieSecure: boolOr("COOKIE_SECURE", false),
	})

	handler := server.New(pool, authService)

	port := envOr("PORT", "8080")
	if port[0] != ':' {
		port = ":" + port
	}
	log.Printf("Server starting on port %s...", port)

	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

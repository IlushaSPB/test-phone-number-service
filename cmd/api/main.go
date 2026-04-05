package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/IlushaSPB/test-phone-number-service/internal/config"
	"github.com/IlushaSPB/test-phone-number-service/internal/handler"
	"github.com/IlushaSPB/test-phone-number-service/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer pool.Close()

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}
	log.Println("database connection established")

	h := handler.New(pool)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/api/numbers/import", h.Import)
	mux.HandleFunc("/api/numbers/search", h.Search)

	chain := middleware.Logging(mux)

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, chain); err != nil {
		log.Fatal(err)
	}
}

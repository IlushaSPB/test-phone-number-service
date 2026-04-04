package handler

import (
	"context"

	"github.com/IlushaSPB/test-phone-number-service/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Handler {
	return &Handler{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (h *Handler) Ping(ctx context.Context) error {
	return h.pool.Ping(ctx)
}

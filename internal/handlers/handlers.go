package handlers

import (
	"github.com/johan253/idme/internal/config"
	"github.com/johan253/idme/internal/db"
)

type Handler struct {
	cfg *config.Config
	q   *db.Queries
}

func New(cfg *config.Config, q *db.Queries) *Handler {
	return &Handler{cfg: cfg, q: q}
}

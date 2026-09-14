package handlers

import (
	"github.com/johan253/idme/internal/config"
	"github.com/johan253/idme/internal/db"
	"github.com/johan253/idme/internal/keys"
)

type Handler struct {
	cfg *config.Config
	q   *db.Queries
	m   *keys.Manager
}

func New(cfg *config.Config, q *db.Queries, m *keys.Manager) *Handler {
	return &Handler{cfg: cfg, q: q, m: m}
}

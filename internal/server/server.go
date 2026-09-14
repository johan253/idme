package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/johan253/idme/internal/config"
	"github.com/johan253/idme/internal/db"
	"github.com/johan253/idme/internal/keys"
)

type Server struct {
	cfg *config.Config
	q   *db.Queries
	m   *keys.Manager
}

func NewServer(cfg *config.Config, q *db.Queries, m *keys.Manager) *http.Server {
	NewServer := &Server{
		cfg: cfg,
		q:   q,
		m:   m,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

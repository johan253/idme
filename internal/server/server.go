package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/johan253/idme/internal/config"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *http.Server {
	NewServer := &Server{
		cfg: cfg,
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

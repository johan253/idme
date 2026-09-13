package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddle "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/johan253/idme/internal/handlers"
	"github.com/johan253/idme/internal/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(chimiddle.Logger)
	r.Use(chimiddle.RequestID)
	r.Use(chimiddle.Recoverer)
	r.Use(chimiddle.RedirectSlashes)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	h := handlers.New(s.cfg, s.q)

	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/logout", h.Logout)

	r.With(middleware.Authorization(s.cfg)).Route("/user", func(r chi.Router) {
		r.Get("/", h.GetUser)
		r.Delete("/", h.DeleteUser)
	})

	return r
}

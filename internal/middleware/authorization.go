package middleware

import (
	"log"
	"net/http"

	"github.com/johan253/idme/internal/auth"
	"github.com/johan253/idme/internal/config"
)

func Authorization(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.ExtractToken(r)
			if err != nil {
				log.Printf("authorization: extract token: %v", err)
				http.Error(w, auth.ErrUnauthorized.Error(), http.StatusUnauthorized)
				return
			}
			claims, err := auth.Parse(cfg, token)
			if err != nil {
				log.Printf("authorization: parse token: %v", err)
				http.Error(w, auth.ErrUnauthorized.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(auth.WithClaims(r.Context(), claims)))
		})
	}
}

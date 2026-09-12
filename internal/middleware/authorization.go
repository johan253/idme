package middleware

import (
	"errors"
	"net/http"
)

var UnauthorizedError = errors.New("unauthorized")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, UnauthorizedError.Error(), http.StatusUnauthorized)
		}
	})
}

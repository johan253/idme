package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/johan253/idme/internal/config"
)

// WithClaims returns a copy of ctx with the given claims attached.
func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// ExtractToken pulls a JWT from either the access_token cookie or a Bearer
// Authorization header.
func ExtractToken(r *http.Request) (string, error) {
	return extractToken(r)
}

type ctxKey string

const claimsKey ctxKey = "claims"

var ErrUnauthorized = errors.New("unauthorized")

type User struct {
	Id       string   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

type Claims struct {
	User
	jwt.RegisteredClaims
}

func Sign(cfg *config.Config, user *User) (string, error) {
	iat := time.Now()
	exp := iat.Add(time.Duration(cfg.JwtTtlSeconds) * time.Second)

	claims := &Claims{
		User: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(iat),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JwtSecret))
}

func Parse(cfg *config.Config, token string) (*Claims, error) {
	claims := &Claims{}
	tok, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnauthorized
		}
		return []byte(cfg.JwtSecret), nil
	})
	if err != nil || !tok.Valid {
		return nil, ErrUnauthorized
	}
	return claims, nil
}

func ClaimsFrom(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey).(*Claims)
	return claims, ok
}

func extractToken(r *http.Request) (string, error) {
	if cookie, err := r.Cookie("access_token"); err == nil {
		return cookie.Value, nil
	}
	header := r.Header.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer "), nil
	}
	return "", ErrUnauthorized
}

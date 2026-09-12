package handlers

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/johan253/idme/internal/auth"
	"github.com/johan253/idme/internal/db"
	"github.com/johan253/idme/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type registerResponse struct {
	AccessToken string    `json:"access_token"`
	User        auth.User `json:"user"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
}

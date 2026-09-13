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
	Username string `json:"username" validate:"required,min=3,max=30,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=72,password_secure"`
	Email    string `json:"email" validate:"required,email"`
}

type registerResponse struct {
	AccessToken string    `json:"access_token"`
	User        auth.User `json:"user"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	req := &registerRequest{}
	if err := utils.ReadJSONFromBody(r, req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := utils.Validator.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, utils.FriendlyValidationError(err))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	user, err := h.q.CreateUser(r.Context(), db.CreateUserParams{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Email:        req.Email,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			writeError(w, http.StatusConflict, "Username or email already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	authUser := &auth.User{
		Id:       user.ID.String(),
		Username: user.Username,
		Roles:    user.Roles,
	}
	accessToken, err := auth.Sign(h.cfg, authUser)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}
	utils.WriteJSONToBody(w, http.StatusCreated, &registerResponse{
		User:        *authUser,
		AccessToken: accessToken,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "Login endpoint not implemented yet")
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "Logout endpoint not implemented yet")
}

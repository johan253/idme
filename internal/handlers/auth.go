package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/johan253/idme/internal/auth"
	"github.com/johan253/idme/internal/db"
	"github.com/johan253/idme/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=72,password_secure"`
}

type registerRequest struct {
	loginRequest
	Email string `json:"email" validate:"required,email"`
}

type authResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	User         auth.User `json:"user"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	req := &registerRequest{}
	if err := utils.ReadJSONFromBody(r, req); err != nil {
		log.Printf("register: decode request body: %v", err)
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := utils.Validator.Struct(req); err != nil {
		log.Printf("register: validation failed for username=%q email=%q: %v", req.Username, req.Email, err)
		writeError(w, http.StatusBadRequest, utils.FriendlyValidationError(err))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("register: bcrypt hash: %v", err)
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	user, err := h.q.CreateUser(r.Context(), db.CreateUserParams{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Email:        req.Email,
		Roles:        []string{"user"},
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			log.Printf("register: unique violation for username=%q email=%q: %v", req.Username, req.Email, err)
			writeError(w, http.StatusConflict, "Username or email already exists")
			return
		}
		log.Printf("register: create user (username=%q email=%q): %v", req.Username, req.Email, err)
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
		log.Printf("register: sign token for user id=%s: %v", authUser.Id, err)
		writeError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}
	utils.WriteJSONToBody(w, http.StatusCreated, &authResponse{
		User:        *authUser,
		AccessToken: accessToken,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	req := &loginRequest{}
	if err := utils.ReadJSONFromBody(r, req); err != nil {
		log.Printf("login: decode request body: %v", err)
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	// Validate the request before enumerating users
	if err := utils.Validator.Struct(req); err != nil {
		log.Printf("login: validation failed for username=%q: %v", req.Username, err)
		writeError(w, http.StatusBadRequest, "Invalid username or password")
		return
	}
	user, err := h.q.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		log.Printf("login: get user by username=%q: %v", req.Username, err)
		writeError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Printf("login: password mismatch for username=%q: %v", req.Username, err)
		writeError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}
	authUser := &auth.User{
		Id:       user.ID.String(),
		Username: user.Username,
		Roles:    user.Roles,
	}
	token, err := auth.Sign(h.cfg, authUser)
	if err != nil {
		log.Printf("login: sign token for user id=%s: %v", user.ID.String(), err)
		writeError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}
	utils.WriteJSONToBody(w, http.StatusOK, &authResponse{
		User:        *authUser,
		AccessToken: token,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "Logout endpoint not implemented yet")
}

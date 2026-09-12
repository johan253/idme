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
	Password string `json:"password" validate:"required,min=8,max=100,password_secure"`
	Email    string `json:"email" validate:"required,email"`
}

type registerResponse struct {
	AccessToken string    `json:"access_token,omitempty"`
	User        auth.User `json:"user,omitempty"`
	Error       string    `json:"error,omitempty"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	req := &registerRequest{}
	// Read and decode request body
	if err := utils.ReadJSONFromBody(r, req); err != nil {
		utils.WriteJSONToBody(w, http.StatusBadRequest, &registerResponse{
			Error: "Invalid request body",
		})
		return
	}
	// Validate required fields
	validate := utils.NewValidator()
	if err := validate.Struct(req); err != nil {
		utils.WriteJSONToBody(w, http.StatusBadRequest, &registerResponse{
			Error: err.Error(),
		})
		return
	}
	// Hash and salt the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteJSONToBody(w, http.StatusInternalServerError, &registerResponse{
			Error: "Internal server error",
		})
		return
	}
	// Create the user in the database
	user, err := h.q.CreateUser(r.Context(), db.CreateUserParams{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Email:        req.Email,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		// Check if the error is a unique constraint violation (duplicate username or email)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			utils.WriteJSONToBody(w, http.StatusConflict, &registerResponse{
				Error: "Username or email already exists",
			})
			return
		}
		utils.WriteJSONToBody(w, http.StatusInternalServerError, &registerResponse{
			Error: "Failed to create user",
		})
		return
	}
	// Generate a JWT token for the newly registered user
	authUser := &auth.User{
		Id:       user.ID.String(),
		Username: user.Username,
		Roles:    user.Roles,
	}
	accessToken, err := auth.Sign(h.cfg, authUser)
	if err != nil {
		utils.WriteJSONToBody(w, http.StatusInternalServerError, &registerResponse{
			Error: "Failed to generate access token",
		})
		return
	}
	// Return the user and access token in the response
	utils.WriteJSONToBody(w, http.StatusCreated, &registerResponse{
		User:        *authUser,
		AccessToken: accessToken,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONToBody(w, http.StatusNotImplemented, map[string]string{
		"error": "Login endpoint not implemented yet",
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONToBody(w, http.StatusNotImplemented, map[string]string{
		"error": "Logout endpoint not implemented yet",
	})
}

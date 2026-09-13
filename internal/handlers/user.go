package handlers

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/johan253/idme/internal/auth"
	"github.com/johan253/idme/internal/utils"
)

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFrom(r.Context())
	if !ok {
		log.Printf("get user: no claims in context")
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var userId pgtype.UUID
	if err := userId.Scan(claims.Id); err != nil {
		log.Printf("get user: scan claims.Id=%v: %v", claims.Id, err)
		writeError(w, http.StatusInternalServerError, "Failed to parse user ID")
		return
	}
	user, err := h.q.GetUserById(r.Context(), userId)
	if err != nil {
		log.Printf("get user: get user by ID=%d: %v", claims.Id, err)
		writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}
	authUser := &auth.User{
		Id:       user.ID.String(),
		Username: user.Username,
		Roles:    user.Roles,
	}
	utils.WriteJSONToBody(w, http.StatusOK, authUser)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFrom(r.Context())
	if !ok {
		log.Printf("delete user: no claims in context")
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var userId pgtype.UUID
	if err := userId.Scan(claims.Id); err != nil {
		log.Printf("delete user: scan claims.Id=%v: %v", claims.Id, err)
		writeError(w, http.StatusInternalServerError, "Failed to parse user ID")
		return
	}
	err := h.q.DeleteUser(r.Context(), userId)
	if err != nil {
		log.Printf("delete user: delete user by ID=%d: %v", claims.Id, err)
		writeError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	utils.WriteJSONToBody(w, http.StatusOK, &successResponse{Message: "User deleted successfully"})
}

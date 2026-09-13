package handlers

import (
	"net/http"

	"github.com/johan253/idme/internal/utils"
)

type errorResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Message string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, msg string) {
	utils.WriteJSONToBody(w, status, errorResponse{Error: msg})
}

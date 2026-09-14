package handlers

import (
	"net/http"

	"github.com/johan253/idme/internal/utils"
)

func (h *Handler) JWKS(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONToBody(w, http.StatusOK, h.m.PublishSet())
}

package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSONToBody(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

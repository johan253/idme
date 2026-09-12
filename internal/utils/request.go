package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeJSONBody(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Disallow unknown fields
	if err := decoder.Decode(&dst); err != nil {
		return fmt.Errorf("failed to decode JSON body: %w", err)
	}
	return nil
}

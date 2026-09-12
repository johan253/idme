package utils

import (
	"encoding/json"
	"net/http"
)

func ReadJSONFromBody(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

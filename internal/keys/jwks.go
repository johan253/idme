package keys

import (
	"encoding/base64"
	"math/big"
)

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

// PublishSet returns the JWKS document for all keys currently in the cache.
func (m *Manager) PublishSet() JWKS {
	all := m.All()
	out := JWKS{Keys: make([]JWK, 0, len(all))}
	for _, k := range all {
		out.Keys = append(out.Keys, JWK{
			Kty: "RSA",
			Kid: k.Kid,
			Use: "sig",
			Alg: "RS256",
			N:   base64.RawURLEncoding.EncodeToString(k.Public.N.Bytes()),
			E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(k.Public.E)).Bytes()),
		})
	}
	return out
}

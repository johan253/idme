package keys

import (
	"encoding/base64"
)

type JWK struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	X   string `json:"x"`
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
			Kty: "OKP",
			Crv: "Ed25519",
			Kid: k.Kid,
			Use: "sig",
			Alg: "EdDSA",
			X:   base64.RawURLEncoding.EncodeToString(k.Public),
		})
	}
	return out
}

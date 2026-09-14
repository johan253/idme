package keys

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// KeyPair is a signing key held in memory. Private is nil for verify-only
// entries loaded from other pods' rotations before the active flip.
type KeyPair struct {
	Kid     string
	Public  ed25519.PublicKey
	Private ed25519.PrivateKey
}

// GenerateEd25519 creates a fresh Ed25519 keypair with a random kid.
func GenerateEd25519() (*KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	kid, err := randomKid()
	if err != nil {
		return nil, err
	}
	return &KeyPair{Kid: kid, Public: pub, Private: priv}, nil
}

func randomKid() (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func encodePublicPEM(pub ed25519.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})), nil
}

func decodePublicPEM(s string) (ed25519.PublicKey, error) {
	block, _ := pem.Decode([]byte(s))
	if block == nil {
		return nil, errors.New("invalid public pem")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ed25519Pub, ok := pub.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("expected ed25519 public key, got %T", pub)
	}
	return ed25519Pub, nil
}

func encodePrivatePEM(priv ed25519.PrivateKey) ([]byte, error) {
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), nil
}

func decodePrivatePEM(b []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("invalid private pem")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ed25519Priv, ok := priv.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("expected ed25519 private key, got %T", priv)
	}
	return ed25519Priv, nil
}

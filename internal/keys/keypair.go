package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

const rsaKeyBits = 2048

// KeyPair is a signing key held in memory. Private is nil for verify-only
// entries loaded from other pods' rotations before the active flip.
type KeyPair struct {
	Kid     string
	Public  *rsa.PublicKey
	Private *rsa.PrivateKey
}

// GenerateRSA creates a fresh RSA keypair with a random kid.
func GenerateRSA() (*KeyPair, error) {
	priv, err := rsa.GenerateKey(rand.Reader, rsaKeyBits)
	if err != nil {
		return nil, err
	}
	kid, err := randomKid()
	if err != nil {
		return nil, err
	}
	return &KeyPair{Kid: kid, Public: &priv.PublicKey, Private: priv}, nil
}

func randomKid() (string, error) {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func encodePublicPEM(pub *rsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})), nil
}

func decodePublicPEM(s string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(s))
	if block == nil {
		return nil, errors.New("invalid public pem")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("expected rsa public key, got %T", pub)
	}
	return rsaPub, nil
}

func encodePrivatePEM(priv *rsa.PrivateKey) []byte {
	der := x509.MarshalPKCS1PrivateKey(priv)
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: der})
}

func decodePrivatePEM(b []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, errors.New("invalid private pem")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

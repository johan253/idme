package keys

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
)

// Cipher seals and opens private key material at rest. The default
// implementation uses AES-256-GCM with a KEK held in memory; swap for a
// KMS-backed implementation without touching callers.
type Cipher interface {
	Seal(plaintext []byte) ([]byte, error)
	Open(ciphertext []byte) ([]byte, error)
}

type aesGCMCipher struct {
	aead cipher.AEAD
}

// NewAESGCMCipher builds a Cipher from a 32-byte KEK.
func NewAESGCMCipher(kek []byte) (Cipher, error) {
	if len(kek) != 32 {
		return nil, fmt.Errorf("kek must be 32 bytes, got %d", len(kek))
	}
	block, err := aes.NewCipher(kek)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &aesGCMCipher{aead: aead}, nil
}

func (c *aesGCMCipher) Seal(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return c.aead.Seal(nonce, nonce, plaintext, nil), nil
}

func (c *aesGCMCipher) Open(ciphertext []byte) ([]byte, error) {
	ns := c.aead.NonceSize()
	if len(ciphertext) < ns {
		return nil, errors.New("ciphertext too short")
	}
	nonce, body := ciphertext[:ns], ciphertext[ns:]
	return c.aead.Open(nil, nonce, body, nil)
}

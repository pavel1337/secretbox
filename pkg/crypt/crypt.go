package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

type Crypter interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

type AESGCM struct {
	key []byte
}

var _ Crypter = &AESGCM{}

func NewAESGCM(encryptionKey string) (*AESGCM, error) {
	key := []byte(encryptionKey)
	if len(key) == 0 {
		return nil, fmt.Errorf("AESGCM: encryption key cannot be empty")
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("AESGCM: encryption key is not 32 characters: %v", encryptionKey)
	}
	return &AESGCM{
		key: key,
	}, nil
}

func (c *AESGCM) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (c *AESGCM) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

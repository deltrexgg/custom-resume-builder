package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)


type APIKeyRequest struct {
	APIKey string `json:"api_key"`
}

type APIKeyResponse struct {
	EncryptedKey string `json:"encrypted_key"`
}

func getEncryptionKey() ([]byte, error) {
	key := os.Getenv("API_TOKEN_ENCRYPTION")

	if key == "" {
		return nil, fmt.Errorf("API_TOKEN_ENCRYPTION is not set")
	}

	keyBytes := []byte(key)

	if len(keyBytes) != 32 {
		return nil, fmt.Errorf(
			"API_TOKEN_ENCRYPTION must be exactly 32 bytes",
		)
	}

	return keyBytes, nil
}

func EncryptAPIKey(apiKey string) (string, error) {
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(
		nonce,
		nonce,
		[]byte(apiKey),
		nil,
	)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptAPIKey(encryptedKey string) (string, error) {
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return "", fmt.Errorf("invalid encrypted key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()

	if len(data) < nonceSize {
		return "", fmt.Errorf("invalid encrypted key")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	plaintext, err := gcm.Open(
		nil,
		nonce,
		ciphertext,
		nil,
	)

	if err != nil {
		return "", fmt.Errorf("decrypt API key: %w", err)
	}

	return string(plaintext), nil
}

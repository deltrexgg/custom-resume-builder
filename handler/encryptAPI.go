package handler

import (
	"custom-resume-builder/utils"
	"encoding/json"
	"net/http"
)

func APIEncrypt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"Only POST requests are allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var body utils.APIKeyRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(
			w,
			"Invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if body.APIKey == "" {
		http.Error(
			w,
			"api_key is required",
			http.StatusBadRequest,
		)
		return
	}

	encryptedKey, err := utils.EncryptAPIKey(body.APIKey)
	if err != nil {
		http.Error(
			w,
			"Failed to encrypt API key",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(utils.APIKeyResponse{
		EncryptedKey: encryptedKey,
	}); err != nil {
		http.Error(
			w,
			"Failed to encode response",
			http.StatusInternalServerError,
		)
		return
	}
}
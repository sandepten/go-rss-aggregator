package auth

import (
	"errors"
	"net/http"
)

// GetAPIKey gets the API key from the request headers
// Example:
// X-API-KEY: 123456
func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("X-API-KEY")
	if apiKey == "" {
		return "", errors.New("API key not found")
	}
	return apiKey, nil
}

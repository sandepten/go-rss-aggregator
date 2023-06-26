package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey gets the API key from the request headers
// Example:
// Authorization: ApiKey {API_KEY}
func GetAPIKey(headers http.Header) (string, error) {
	head := headers.Get("Authorization")
	if head == "" {
		return "", errors.New("api key not found")
	}
	val := strings.Split(head, " ")
	if len(val) != 2 {
		return "", errors.New("invalid API key format")
	}
	if val[0] != "ApiKey" {
		return "", errors.New("invalid API key format")
	}
	return val[1], nil
}

package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	possibleKey := headers.Get("Authorization")

	if len(possibleKey) == 0 {
		return "", errors.New("No possible apikey string authorized")
	}
	if !strings.HasPrefix(possibleKey, "ApiKey ") {
		return "", errors.New("Not prefixed by Apikey")
	}
	possibleKey = strings.TrimPrefix(possibleKey, "ApiKey ")

	if len(possibleKey) == 0 {
		return "", errors.New("No apikey after trimming")
	}
	return possibleKey, nil
}

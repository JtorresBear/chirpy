package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if len(tokenString) == 0 {
		return "", errors.New("No token string authorized")
	}
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return "", errors.New("Not prefixed by Bearer")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer")
	tokenString = strings.TrimSpace(tokenString)

	if len(tokenString) == 0 {
		return "", errors.New("No token after Bearer")
	}

	return tokenString, nil
}

package tokenextractor

import (
	"errors"
	"net/http"
	"strings"
)

func ExtractTokenFromHttpRequest(r *http.Request) (string, error) {
	authHeaderArray, ok := r.Header["Authorization"]
	if !ok {
		return "", errors.New("authorization header is missing")
	}

	if len(authHeaderArray) == 0 {
		return "", errors.New("authorization header is empty")
	}
	authHeader := authHeaderArray[0]
	words := strings.Split(authHeader, " ")
	if len(words) < 2 || words[0] != "Bearer" {
		return "", errors.New("authorization is not in format 'Bearer ...'")
	}
	return words[1], nil
}

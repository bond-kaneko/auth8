package auth

import (
	"errors"
	"net/http"
)

func GetAccessToken(r *http.Request) (string, error) {
	if _, ok := r.Header["Authorization"]; !ok {
		return "", errors.New("no access token")
	}

	return "", nil
}

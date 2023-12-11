package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("no API Key found in Headers")
	}

	vals := strings.Split(auth, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed Auth Header")
	}

	if vals[0] != "APIKey" {
		return "", errors.New("malformed Auth Header")
	}

	return vals[1], nil
}
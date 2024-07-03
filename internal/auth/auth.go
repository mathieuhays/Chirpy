package auth

import (
	"errors"
	"net/http"
	"slices"
	"strings"
)

const (
	TypeBearer = "Bearer"
	TypeApiKey = "ApiKey"
)

var (
	errMalformed   = errors.New("malformed authorization header")
	errInvalidType = errors.New("invalid authorization type")
)

var validTypes = []string{
	TypeApiKey,
	TypeBearer,
}

type Authorization struct {
	Name  string
	Value string
}

func GetAuthorization(header http.Header) (Authorization, error) {
	parts := strings.Split(header.Get("Authorization"), " ")
	if len(parts) != 2 {
		return Authorization{}, errMalformed
	}

	if !slices.Contains(validTypes, parts[0]) {
		return Authorization{}, errInvalidType
	}

	return Authorization{
		Name:  parts[0],
		Value: parts[1],
	}, nil
}

package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type chirpyClaims struct {
	jwt.RegisteredClaims
}

func generateRefreshToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func generateAccessToken(userId int, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour).UTC()),
		Subject:   strconv.Itoa(userId),
	})
	return token.SignedString([]byte(jwtSecret))
}

func verifyAccessToken(accessToken string, jwtSecret string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &chirpyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*chirpyClaims)
	if !ok {
		return 0, errors.New("could not type cast claims")
	}

	return strconv.Atoi(claims.Subject)
}

func getTokenFromRequest(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}

package main

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (a *apiConfig) handlePostLogin(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email            string
		Password         string
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		writeError(w, errors.New("could not decode payload"), http.StatusBadRequest)
		return
	}

	if len(payload.Email) == 0 || len(payload.Password) == 0 {
		writeError(w, errors.New("payload incomplete"), http.StatusBadRequest)
		return
	}

	if !isValidEmail(payload.Email) {
		writeError(w, errors.New("invalid email"), http.StatusBadRequest)
		return
	}

	dbUser, exists := a.database.GetUserByEmail(payload.Email)
	if !exists {
		writeError(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(payload.Password)); err != nil {
		writeError(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	expiresIn := time.Duration(payload.ExpiresInSeconds) * time.Second
	day := time.Hour * time.Duration(24)
	if expiresIn == 0 || expiresIn > day {
		expiresIn = day
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
		Subject:   strconv.Itoa(dbUser.Id),
	})
	signedString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		log.Println(err)
		writeError(w, errors.New("something went wrong"), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		Id:    dbUser.Id,
		Email: dbUser.Email,
		Token: signedString,
	})
}

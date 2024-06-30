package main

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (a *apiConfig) handlerPutUsers(w http.ResponseWriter, r *http.Request) {
	authorization := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	type ChirpyClaims struct {
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(authorization, &ChirpyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.jwtSecret), nil
	})

	if err != nil {
		log.Println(err)
		writeError(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*ChirpyClaims)
	if !ok {
		writeError(w, errSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		writeError(w, errSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	var payload struct {
		Email    string
		Password string
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
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

	dbUser, exists := a.database.GetUser(userId)
	if !exists {
		writeError(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	password, err := validatePassword(payload.Password)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	newDbUser, err := a.database.UpdateUser(dbUser.Id, payload.Email, password)
	if err != nil {
		writeError(w, errSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, user{
		Id:    newDbUser.Id,
		Email: newDbUser.Email,
	})
}

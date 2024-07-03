package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

func (a *apiConfig) handlePostLogin(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string
		Password string
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
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(payload.Password)); err != nil {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	signedString, err := generateAccessToken(dbUser.Id, a.jwtSecret)
	if err != nil {
		log.Println(err)
		writeError(w, errSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	newToken, err := generateRefreshToken()
	if err != nil {
		log.Println(err)
		writeError(w, errSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	session, err := a.database.CreateSession(dbUser.Id, newToken, time.Now().Add(time.Hour*24*6))
	if err != nil {
		log.Println(err)
		writeError(w, errSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, struct {
		Id           int    `json:"id"`
		Email        string `json:"email"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		Id:           dbUser.Id,
		Email:        dbUser.Email,
		IsChirpyRed:  dbUser.IsChirpyRed,
		Token:        signedString,
		RefreshToken: session.Token,
	})
}

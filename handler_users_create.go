package main

import (
	"encoding/json"
	"errors"
	"github.com/mathieuhays/Chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (a *apiConfig) handlerPostUsers(w http.ResponseWriter, r *http.Request) {
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

	password, err := validatePassword(payload.Password)
	if err != nil {
		writeError(w, err, http.StatusBadRequest)
		return
	}

	newUser, err := a.database.CreateUser(payload.Email, password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			writeError(w, errors.New("password is too long"), http.StatusBadRequest)
			return
		}

		if errors.Is(err, database.ErrEmailAlreadyExists) {
			writeError(w, err, http.StatusBadRequest)
			return
		}

		writeError(w, errSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, user{
		Id:          newUser.Id,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed,
	})
}

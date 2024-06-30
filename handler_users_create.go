package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (a *apiConfig) handlerPostUsers(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email string
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		writeError(w, errors.New("could not decode payload"), http.StatusBadRequest)
		return
	}

	if len(payload.Email) == 0 {
		writeError(w, errors.New("payload incomplete"), http.StatusBadRequest)
		return
	}

	if !isValidEmail(payload.Email) {
		writeError(w, errors.New("invalid email"), http.StatusBadRequest)
		return
	}

	newUser, err := a.database.CreateUser(payload.Email)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, user{
		Id:    newUser.Id,
		Email: newUser.Email,
	})
}

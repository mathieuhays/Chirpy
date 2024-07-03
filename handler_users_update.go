package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (a *apiConfig) handlerPutUsers(w http.ResponseWriter, r *http.Request) {
	userId, err := verifyAccessToken(getTokenFromRequest(r), a.jwtSecret)
	if err != nil {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
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
		writeError(w, errUnauthorized, http.StatusUnauthorized)
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

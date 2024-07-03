package main

import (
	"encoding/json"
	"errors"
	"github.com/mathieuhays/Chirpy/internal/auth"
	"net/http"
)

func (a *apiConfig) handlerPostPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetAuthorization(r.Header)
	if err != nil || token.Name != auth.TypeApiKey {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	if token.Value != a.polkaApiKey {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	var payload struct {
		Event string
		Data  struct {
			UserId *int `json:"user_id"`
		}
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		writeError(w, errors.New("could not decode payload"), http.StatusBadRequest)
		return
	}

	if payload.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if payload.Data.UserId == nil {
		writeError(w, errors.New("invalid parameters"), http.StatusBadRequest)
		return
	}

	dbUser, exists := a.database.GetUser(*payload.Data.UserId)
	if !exists {
		writeError(w, errors.New("not found"), http.StatusNotFound)
		return
	}

	dbUser.IsChirpyRed = true
	if _, err = a.database.UpdateUser(dbUser); err != nil {
		writeError(w, errSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

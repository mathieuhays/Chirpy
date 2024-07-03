package main

import (
	"encoding/json"
	"errors"
	"github.com/mathieuhays/Chirpy/internal/auth"
	"net/http"
)

func (a *apiConfig) handlerPostChirps(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetAuthorization(r.Header)
	if err != nil || token.Name != auth.TypeBearer {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	userId, err := verifyAccessToken(token.Value, a.jwtSecret)
	if err != nil {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	dbUser, exists := a.database.GetUser(userId)
	if !exists {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	var payload struct {
		Body string
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		writeError(w, errors.New("could not decode payload"), http.StatusBadRequest)
		return
	}

	if len(payload.Body) == 0 {
		writeError(w, errors.New("payload incomplete"), http.StatusBadRequest)
		return
	}

	if len(payload.Body) > chirpMaxLength {
		writeError(w, errors.New("chirp is too long"), http.StatusBadRequest)
		return
	}

	cleanBody := censorProfanities(payload.Body, map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	})

	newChirp, err := a.database.CreateChirp(cleanBody, dbUser.Id)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, chirp{
		Id:       newChirp.Id,
		Body:     newChirp.Body,
		AuthorId: newChirp.AuthorId,
	})
}

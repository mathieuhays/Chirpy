package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (a *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := a.database.GetChirps()
	if err != nil {
		writeError(w, err, 500)
		return
	}

	writeJSON(w, 200, chirps)
}

func (a *apiConfig) handlerPostChirps(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Body string
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		writeError(w, errors.New("could not decode payload"), 400)
		return
	}

	if len(payload.Body) == 0 {
		writeError(w, errors.New("payload incomplete"), 400)
		return
	}

	newChirp, err := a.database.CreateChirp(payload.Body)
	if err != nil {
		writeError(w, err, 400)
		return
	}

	writeJSON(w, 201, newChirp)
}

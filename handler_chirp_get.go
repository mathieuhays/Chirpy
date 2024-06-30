package main

import (
	"errors"
	"net/http"
	"strconv"
)

func (a *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		writeError(w, errors.New("Invalid parameter"), http.StatusBadRequest)
		return
	}

	chirp, exists := a.database.GetChirp(chirpID)
	if !exists {
		writeError(w, errors.New("not found"), http.StatusNotFound)
		return
	}

	writeJSON(w, 200, chirp)
}

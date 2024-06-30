package main

import (
	"errors"
	"net/http"
	"strconv"
)

func (a *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := a.database.GetChirps()
	if err != nil {
		writeError(w, err, 500)
		return
	}

	var output = make([]chirp, len(chirps))
	for i, c := range chirps {
		output[i] = chirp{
			Id:   c.Id,
			Body: c.Body,
		}
	}

	writeJSON(w, 200, chirps)
}

func (a *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		writeError(w, errors.New("Invalid parameter"), http.StatusBadRequest)
		return
	}

	dbChirp, exists := a.database.GetChirp(chirpID)
	if !exists {
		writeError(w, errors.New("not found"), http.StatusNotFound)
		return
	}

	writeJSON(w, 200, chirp{
		Id:   dbChirp.Id,
		Body: dbChirp.Body,
	})
}

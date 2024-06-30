package main

import (
	"net/http"
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

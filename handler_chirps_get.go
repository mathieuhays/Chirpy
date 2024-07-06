package main

import (
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func (a *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := a.database.GetChirps()
	if err != nil {
		writeError(w, err, 500)
		return
	}

	authorIdRaw := r.URL.Query().Get("author_id")
	sortMethod := strings.ToLower(r.URL.Query().Get("sort"))
	var authorId *int

	if authorIdRaw != "" {
		id, err := strconv.Atoi(authorIdRaw)
		if err != nil {
			writeError(w, errors.New("invalid parameter"), http.StatusBadRequest)
			return
		}

		authorId = &id
	}

	var output = make([]chirp, 0)
	for _, c := range chirps {
		if authorId != nil && c.AuthorId != *authorId {
			continue
		}

		output = append(output, chirp{
			Id:       c.Id,
			Body:     c.Body,
			AuthorId: c.AuthorId,
		})
	}

	if sortMethod == "desc" {
		sort.Slice(output, func(i, j int) bool {
			return output[i].Id > output[j].Id
		})
	} else {
		sort.Slice(output, func(i, j int) bool {
			return output[i].Id < output[j].Id
		})
	}

	writeJSON(w, 200, output)
}

func (a *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		writeError(w, errors.New("invalid parameter"), http.StatusBadRequest)
		return
	}

	dbChirp, exists := a.database.GetChirp(chirpID)
	if !exists {
		writeError(w, errors.New("not found"), http.StatusNotFound)
		return
	}

	writeJSON(w, 200, chirp{
		Id:       dbChirp.Id,
		Body:     dbChirp.Body,
		AuthorId: dbChirp.AuthorId,
	})
}

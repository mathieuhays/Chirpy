package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	var content struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&content)
	if err != nil {
		writeError(w, errors.New("request malformed"), http.StatusBadRequest)
		return
	}

	if len(content.Body) == 0 {
		writeError(w, errors.New("Chirp is empty"), http.StatusBadRequest)
		return
	}

	validationErr := validateChirp(content.Body)
	if validationErr != nil {
		writeError(w, validationErr, http.StatusBadRequest)
		return
	}

	profanities := []string{"kerfuffle", "sharbert", "fornax"}
	cleanedBody := censorProfanities(content.Body, profanities)

	writeJSON(w, http.StatusOK, struct {
		CleanedBody string `json:"cleaned_body"`
	}{CleanedBody: cleanedBody})
}

func writeError(w http.ResponseWriter, err error, statusCode int) {
	writeJSON(w, statusCode, errorResponse{Error: err.Error()})
}

func writeJSON(w http.ResponseWriter, statusCode int, object interface{}) {
	data, mErr := json.Marshal(object)
	if mErr != nil {
		log.Printf("error marshalling error object: %v", mErr.Error())
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

package main

import (
	"encoding/json"
	"errors"
	"io"
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

	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Printf("body reader error: %v", err.Error())
		writeError(w, errors.New("cannot read request"), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &content)
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

	respData, err := json.Marshal(struct {
		Valid bool `json:"valid"`
	}{Valid: true})
	if err != nil {
		writeError(w, errors.New("Something went wrong"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respData)
}

func writeError(w http.ResponseWriter, err error, statusCode int) {
	data, mErr := json.Marshal(errorResponse{Error: err.Error()})
	if mErr != nil {
		log.Printf("error marshalling error object: %v", mErr.Error())
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

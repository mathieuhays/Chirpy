package main

import (
	"errors"
	"strings"
)

const chirpMaxLength = 140

type chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

func validateChirp(body string) error {
	if len(body) > chirpMaxLength {
		return errors.New("chirp is too long")
	}

	return nil
}

func censorProfanities(content string, profanities map[string]struct{}) string {
	words := strings.Split(content, " ")

	for i, word := range words {
		lower := strings.ToLower(word)

		if _, ok := profanities[lower]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

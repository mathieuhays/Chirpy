package main

import (
	"errors"
	"strings"
)

const chirpMaxLength = 140

func validateChirp(content string) error {
	if len(content) > 140 {
		return errors.New("Chirp is too long")
	}

	return nil
}

func censorProfanities(content string, profanities []string) string {
	words := strings.Split(content, " ")

	for i, word := range words {
		for _, profanity := range profanities {
			if profanity == strings.ToLower(word) {
				words[i] = "****"
			}
		}
	}

	return strings.Join(words, " ")
}

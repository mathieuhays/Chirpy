package chirp

import (
	"errors"
	"strings"
)

const maxLength = 140

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func validateChirp(body string) error {
	if len(body) > maxLength {
		return errors.New("chirp is too long")
	}

	return nil
}

func NewChirp(body string) (Chirp, error) {
	if err := validateChirp(body); err != nil {
		return Chirp{}, err
	}

	cleanBody := censorProfanities(body, []string{"kerfuffle", "sharbert", "fornax"})

	return Chirp{Body: cleanBody}, nil
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

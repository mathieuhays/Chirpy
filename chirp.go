package main

import "errors"

func validateChirp(content string) error {
	if len(content) > 140 {
		return errors.New("Chirp is too long")
	}

	return nil
}

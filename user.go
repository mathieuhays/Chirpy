package main

import (
	"errors"
	"net/mail"
	"strings"
)

var (
	errPasswordTooShort = errors.New("password must be at least 6 characters long")
)

type user struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func isValidEmail(email string) bool {
	// no domain (or too many)
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	// missing tld
	if !strings.Contains(parts[1], ".") {
		return false
	}

	// failsafe in case it catches something we don't
	_, err := mail.ParseAddress(email)
	return err == nil
}

func validatePassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errPasswordTooShort
	}

	// @TODO add other requirements like a mix of lowercase, uppercase, digits and special characters

	return password, nil
}

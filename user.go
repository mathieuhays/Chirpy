package main

import (
	"net/mail"
	"strings"
)

type user struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
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

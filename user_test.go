package main

import (
	"errors"
	"fmt"
	"testing"
)

func Test_isValidEmail(t *testing.T) {
	tests := []struct {
		email string
		want  bool
	}{
		{"obviously.invalid", false},
		{"another-bad-one", false},
		{"email@example.com", true},
		{"missing@tld", false},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("email test #%d", i), func(t *testing.T) {
			if got := isValidEmail(tt.email); got != tt.want {
				t.Errorf("isValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidPassword(t *testing.T) {
	tests := []struct {
		password string
		want     error
	}{
		{"test", errPasswordTooShort},
		{"weak?", errPasswordTooShort},
		{"iamastrongpassword", nil},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("password test #%d", i), func(t *testing.T) {
			if _, got := validatePassword(tt.password); !errors.Is(got, tt.want) {
				t.Errorf("validatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

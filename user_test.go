package main

import (
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

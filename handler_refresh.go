package main

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

func (a *apiConfig) handlePostRefresh(w http.ResponseWriter, r *http.Request) {
	authorization := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	if len(authorization) == 0 {
		writeError(w, errors.New("not authorized"), http.StatusUnauthorized)
		return
	}

	session, exists := a.database.GetSession(authorization)
	if !exists {
		writeError(w, errors.New("not authorized"), http.StatusUnauthorized)
		return
	}

	if session.Expiration.Before(time.Now()) {
		_ = a.database.RevokeSession(session.Token)
		writeError(w, errors.New("not authorized"), http.StatusUnauthorized)
		return
	}

	newToken, err := generateAccessToken(session.User, a.jwtSecret)
	if err != nil {
		writeError(w, errSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: newToken,
	})
}

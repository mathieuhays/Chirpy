package main

import (
	"net/http"
	"time"
)

func (a *apiConfig) handlePostRefresh(w http.ResponseWriter, r *http.Request) {
	authorization := getTokenFromRequest(r)
	if len(authorization) == 0 {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	session, exists := a.database.GetSession(authorization)
	if !exists {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	if session.Expiration.Before(time.Now()) {
		_ = a.database.RevokeSession(session.Token)
		writeError(w, errUnauthorized, http.StatusUnauthorized)
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

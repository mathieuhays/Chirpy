package main

import (
	"errors"
	"net/http"
	"strings"
)

func (a *apiConfig) handlePostRevoke(w http.ResponseWriter, r *http.Request) {
	authorization := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	if len(authorization) == 0 {
		writeError(w, errors.New("not authorized"), http.StatusUnauthorized)
		return
	}

	if _, exists := a.database.GetSession(authorization); !exists {
		writeError(w, errors.New("not authorized"), http.StatusUnauthorized)
		return
	}

	err := a.database.RevokeSession(authorization)
	if err != nil {
		writeError(w, errSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

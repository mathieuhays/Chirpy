package main

import (
	"net/http"
)

func (a *apiConfig) handlePostRevoke(w http.ResponseWriter, r *http.Request) {
	authorization := getTokenFromRequest(r)
	if len(authorization) == 0 {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	if _, exists := a.database.GetSession(authorization); !exists {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	err := a.database.RevokeSession(authorization)
	if err != nil {
		writeError(w, errSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

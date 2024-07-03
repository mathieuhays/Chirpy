package main

import (
	"github.com/mathieuhays/Chirpy/internal/auth"
	"net/http"
)

func (a *apiConfig) handlePostRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetAuthorization(r.Header)
	if err != nil || token.Name != auth.TypeBearer {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	if _, exists := a.database.GetSession(token.Value); !exists {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	err = a.database.RevokeSession(token.Value)
	if err != nil {
		writeError(w, errSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package main

import (
	"errors"
	"github.com/mathieuhays/Chirpy/internal/auth"
	"net/http"
	"strconv"
)

func (a *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetAuthorization(r.Header)
	if err != nil || token.Name != auth.TypeBearer {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	userId, err := verifyAccessToken(token.Value, a.jwtSecret)
	if err != nil {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	dbUser, exists := a.database.GetUser(userId)
	if !exists {
		writeError(w, errUnauthorized, http.StatusUnauthorized)
		return
	}

	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		writeError(w, errors.New("invalid parameter"), http.StatusBadRequest)
		return
	}

	dbChirp, exists := a.database.GetChirp(chirpID)
	if !exists {
		writeError(w, errors.New("not found"), http.StatusNotFound)
		return
	}

	if dbChirp.AuthorId != dbUser.Id {
		writeError(w, errors.New("forbidden"), http.StatusForbidden)
		return
	}

	err = a.database.DeleteChirp(dbChirp.Id)
	if err != nil {
		writeError(w, errSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

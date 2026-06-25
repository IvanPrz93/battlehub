package main

import (
	"net/http"

	"github.com/IvanPrz93/battlehub/internal/auth"
)

func (cfg *apiConfig) handlerMyMatchHistory(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	matches, err := cfg.db.GetMatchHistory(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Cannot find user", err)
		return
	}

	var respBody []Match

	for _, match := range matches {
		respBody = append(respBody, Match{
			ID:        		match.ID,
			Player1ID:		match.Player1ID,
			Player2ID:		match.Player2ID,
			WinnerID:		match.WinnerID.UUID,
			CreatedAt: 		match.CreatedAt,
			CompletedAt:	match.CompletedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, respBody)

}

func (cfg *apiConfig) handlerGetMatchHistory(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	user, err := cfg.db.GetUserByUsername(r.Context(), username)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Cannot find user", err)
		return
	}

	matches, err := cfg.db.GetMatchHistory(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Cannot find user", err)
		return
	}

	var respBody []Match

	for _, match := range matches {
		respBody = append(respBody, Match{
			ID:        		match.ID,
			Player1ID:		match.Player1ID,
			Player2ID:		match.Player2ID,
			WinnerID:		match.WinnerID.UUID,
			CreatedAt: 		match.CreatedAt,
			CompletedAt:	match.CompletedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, respBody)

}
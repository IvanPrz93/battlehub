package main

import (
	"net/http"

	"github.com/IvanPrz93/battlehub/internal/auth"
)

func (cfg *apiConfig) handlerMyProfile(w http.ResponseWriter, r *http.Request) {

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

	user, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
			ID:        		user.ID,
			Username:		user.Username,
			Email:			user.Email,
			CreatedAt: 		user.CreatedAt,
		})
}

func (cfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	user, err := cfg.db.GetUserByUsername(r.Context(), username)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
			ID:        		user.ID,
			Username:		user.Username,
			Email:			user.Email,
			CreatedAt: 		user.CreatedAt,
		})

}
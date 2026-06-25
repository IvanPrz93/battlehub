package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/IvanPrz93/battlehub/internal/database"
	"github.com/IvanPrz93/battlehub/internal/auth"
	"github.com/google/uuid"
)

type User struct {
	ID        		uuid.UUID 	`json:"id"`
	Username       	string		`json: "username"`
	Email     		string    	`json:"email"`
	HashedPassword 	string 		`json: "hashed_password"`
	CreatedAt 		time.Time 	`json:"created_at"`
	Wins           	int32		`json:"wins"`
	Loses          	int32		`json:"loses"`
	Rating         	int32		`json:"rating"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username	string	`json: "username"`
		Email 		string 	`json: "email"`
		Password	string	`json: "password"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	arg := database.CreateUserParams{
		Username: params.Username,
		Email: params.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.db.CreateUser(r.Context(), arg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        		user.ID,
			Username:		user.Username,
			Email:			user.Email,
			CreatedAt: 		user.CreatedAt,
			Wins:           user.Wins,
			Loses:          user.Loses,
			Rating:         user.Rating,
		},
	})
}

package main

import (
	"encoding/json"
	"net/http"
	"time"
	"database/sql"
	"math"

	"github.com/IvanPrz93/battlehub/internal/database"
	"github.com/google/uuid"
)

type Match struct {
	ID        		uuid.UUID 		`json:"id"`
	Player1ID 		uuid.UUID 		`json:"player1_id"`
	Player2ID 		uuid.UUID 		`json:"player2_id"`
	WinnerID       	uuid.UUID		`json: "winner_id"`
	CreatedAt 		time.Time 		`json:"created_at"`
	CompletedAt 	sql.NullTime	`json:"completed_at"`
}

func (cfg *apiConfig) createNewMatch(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Player1ID 		uuid.UUID 	`json:"player1_id"`
		Player2ID 		uuid.UUID 	`json:"player2_id"`
	}

	type response struct {
		Match
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	match, err := cfg.db.CreateMatch(r.Context(), database.CreateMatchParams{
		Player1ID: params.Player1ID,
		Player2ID: params.Player2ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create match", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Match: Match{
			ID:        		match.ID,
			Player1ID:		match.Player1ID,
			Player2ID:		match.Player2ID,
			CreatedAt: 		match.CreatedAt,
		},
	})
}


func (cfg *apiConfig) endMatch(w http.ResponseWriter, r *http.Request) {
	matchID := r.PathValue("id")
	match_uuid, err := uuid.Parse(matchID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid match Id: ", err)
		return
	}
	
	type parameters struct {
		WinnerID 	uuid.UUID 	`json:"winner_id"`
	}

	type response struct {
		Match
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	tx, err := cfg.dbConn.BeginTx(r.Context(), nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't start transaction", err)
		return 
	}

	defer tx.Rollback()

	txQueries := cfg.db.WithTx(tx)

	match, err := txQueries.UpdateMatchResults(r.Context(), database.UpdateMatchResultsParams{
		ID: 		match_uuid,
		WinnerID: 	uuid.NullUUID{
			UUID:  params.WinnerID,
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update match results", err)
		return
	}

	if !(match.Player1ID == match.WinnerID.UUID || match.Player2ID == match.WinnerID.UUID) {
		respondWithError(w, http.StatusBadRequest, "Winner was not a player in this match", err)
		return
	}

	var loserID uuid.UUID
	if match.Player1ID == match.WinnerID.UUID {
		loserID = match.Player2ID
	} else {
		loserID = match.Player1ID
	}

	winner, err := txQueries.GetUserByID(r.Context(), match.WinnerID.UUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user", err)
		return
	}

	loser, err := txQueries.GetUserByID(r.Context(), loserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user", err)
		return
	}

	const K = 32
	expectedScore := 1/(1 + math.Pow(10, (float64(loser.Rating) - float64(winner.Rating))/400))
	ratingChange := K * (1 - expectedScore)

	winner.Rating += int32(ratingChange)
	loser.Rating -= int32(ratingChange)

	winner.Wins += 1
	loser.Loses += 1

	_, err = txQueries.UpdateUserResults(r.Context(), database.UpdateUserResultsParams{
		ID:     winner.ID,
		Wins:   winner.Wins,
		Loses:  winner.Loses,
		Rating: winner.Rating,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update winner", err)
		return
	}

	_, err = txQueries.UpdateUserResults(r.Context(), database.UpdateUserResultsParams{
		ID:     loser.ID,
		Wins:   loser.Wins,
		Loses:  loser.Loses,
		Rating: loser.Rating,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update loser", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't commit transaction", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		Match: Match{
			ID:        		match.ID,
			Player1ID:		match.Player1ID,
			Player2ID:		match.Player2ID,
			WinnerID:		match.WinnerID.UUID,
			CreatedAt: 		match.CreatedAt,
			CompletedAt:	match.CompletedAt,
		},
	})
}

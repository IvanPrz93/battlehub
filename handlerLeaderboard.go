package main

import (
	"net/http"
	"strconv"
)

type LeaderboardRow struct {
	Rank     	int64	`json:"rank"`
	Username 	string	`json:"username"`
	Wins     	int32	`json:"wins"`
	Loses    	int32	`json:"loses"`
	Rating   	int32	`json:"rating"`
}

func (cfg *apiConfig) handlerLeaderboard(w http.ResponseWriter, r *http.Request) {
	var limitNum int64
	var err error
	limit := r.URL.Query().Get("limit")
	
	if limit != ""{
		limitNum, err = strconv.ParseInt(limit, 10, 32)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unrecognized limit", err)
			return
		}
	} else {
		limitNum = 20
	}

	leaderboard, err := cfg.db.GetLeaderboard(r.Context(), int32(limitNum))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Cannot list leaderboard", err)
		return
	}

	var respBody []LeaderboardRow

	for _, pos := range leaderboard {
		respBody = append(respBody, LeaderboardRow{
			Rank:		pos.Rank,
			Username:	pos.Username,
			Wins:		pos.Wins,
			Loses:		pos.Loses,
			Rating:		pos.Rating,
		})
	}

	respondWithJSON(w, http.StatusOK, respBody)
	
}
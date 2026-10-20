package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/JtorresBear/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) chirpHandler(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	type reqBody struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(respW, 500, "Could not read request")
		return
	}

	chirpReq := reqBody{}
	err = json.Unmarshal(data, &chirpReq)
	if err != nil {
		respondWithError(respW, 400, "Something wrong with data")
		return
	}

	if len(chirpReq.Body) > 140 {
		respondWithError(respW, 400, "Chirp is too long")
		return
	}
	chirpSplit := strings.Split(chirpReq.Body, " ")

	for index, word := range chirpSplit {
		for _, bWord := range bannedWords {
			if strings.ToLower(word) == bWord {
				chirpSplit[index] = "****"
			}
		}
	}
	chirpReq.Body = strings.Join(chirpSplit, " ")
	params := database.CreateChirpParams{}
	params.Body = chirpReq.Body
	params.UserID = chirpReq.UserID
	sqlChirp, err := cfg.dbQueries.CreateChirp(req.Context(), params)
	if err != nil {
		respondWithError(respW, 500, "Couldn't create chirp")
		return
	}
	resChirp := Chirp{
		ID:        sqlChirp.ID,
		CreatedAt: sqlChirp.CreatedAt,
		UpdatedAt: sqlChirp.UpdatedAt,
		Body:      sqlChirp.Body,
		UserID:    sqlChirp.UserID,
	}

	respondWithJson(respW, 201, resChirp)
}

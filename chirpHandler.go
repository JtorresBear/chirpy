package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/JtorresBear/chirpy/internal/auth"
	"github.com/JtorresBear/chirpy/internal/database"
)

func (cfg *apiConfig) chirpHandler(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(respW, 401, "Invalid Token")
		return
	}
	authUser, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(respW, 401, "Invalid Token")
		return
	}
	type reqBody struct {
		Body string `json:"body"`
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
	params.UserID = authUser
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

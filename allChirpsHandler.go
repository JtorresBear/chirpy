package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) GetAllChirps(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	sqlChirps, err := cfg.dbQueries.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(respW, 500, "Couldn't load chirps from database")
		return
	}

	resChirps := make([]Chirp, 0)

	for _, chirp := range sqlChirps {
		resChirp := Chirp{}
		resChirp.ID = chirp.ID
		resChirp.CreatedAt = chirp.CreatedAt
		resChirp.UpdatedAt = chirp.UpdatedAt
		resChirp.Body = chirp.Body
		resChirp.UserID = chirp.UserID
		resChirps = append(resChirps, resChirp)
	}
	respondWithJson(respW, 200, resChirps)
}

func (cfg *apiConfig) getChirp(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	possibleUUID := req.PathValue("chirpID")

	parsedUUID, err := uuid.Parse(possibleUUID)
	if err != nil {
		respondWithError(respW, 500, "Couldn't get the uuid")
		return
	}

	sqlChirp, err := cfg.dbQueries.GetChirp(req.Context(), parsedUUID)
	if err != nil {
		respondWithError(respW, 404, "no such chirp")
		return
	}
	chirp := Chirp{
		ID:        sqlChirp.ID,
		CreatedAt: sqlChirp.CreatedAt,
		UpdatedAt: sqlChirp.UpdatedAt,
		Body:      sqlChirp.Body,
		UserID:    sqlChirp.UserID,
	}

	respondWithJson(respW, 200, chirp)
}

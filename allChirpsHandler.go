package main

import (
	"net/http"
	"sort"

	"github.com/JtorresBear/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) GetAllChirps(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	authID := req.URL.Query().Get("author_id")
	parsedAuthID := uuid.UUID{}
	sqlChirps := []database.Chirp{}
	var err error
	if authID == "" {
		sqlChirps, err = cfg.dbQueries.GetAllChirps(req.Context())
		if err != nil {
			respondWithError(respW, 500, "Couldn't load chirps from database")
			return
		}
	} else {
		var err2 error
		parsedAuthID, err2 = uuid.Parse(authID)
		if err2 != nil {
			respondWithError(respW, 500, err2.Error())
			return
		}
		sqlChirps, err2 = cfg.dbQueries.GetAllChirpsFromUser(req.Context(), parsedAuthID)
		if err2 != nil {
			respondWithError(respW, 500, err2.Error())
		}
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

	sortOrder := req.URL.Query().Get("sort")
	if sortOrder == "desc" {
		sort.Slice(resChirps, func(i int, j int) bool {
			return resChirps[i].CreatedAt.After(resChirps[j].CreatedAt)
		})
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

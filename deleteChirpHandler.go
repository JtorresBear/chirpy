package main

import (
	"net/http"

	"github.com/JtorresBear/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) deleteChirpHandler(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	possibleUUID := req.PathValue("chirpID")

	parsedUUID, err := uuid.Parse(possibleUUID)
	if err != nil {
		respondWithError(respW, 404, "Chirp not found")
		return
	}
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(respW, 401, "Invalid Token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(respW, 401, err.Error())
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(req.Context(), parsedUUID)

	if err != nil {
		respondWithError(respW, 401, err.Error())
		return
	}

	if userID != chirp.UserID {
		respondWithError(respW, 403, "Not author of chirp")
		return
	}

	err = cfg.dbQueries.DeleteChirp(req.Context(), chirp.ID)
	if err != nil {
		respondWithError(respW, 404, err.Error())
		return
	}
	respW.WriteHeader(204)
}

package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/JtorresBear/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) upgradeToRedHandler(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	possibleKey, err := auth.GetAPIKey(req.Header)
	if possibleKey != cfg.polkaKey {
		respondWithError(respW, 401, "You don't have access")
		return
	}
	type reqBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `josn:"data"`
	}
	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(respW, 500, err.Error())
		return
	}
	reqEvent := reqBody{}
	err = json.Unmarshal(data, &reqEvent)
	if err != nil {
		respondWithError(respW, 500, err.Error())
		return
	}

	if reqEvent.Event != "user.upgraded" {
		respW.WriteHeader(204)
		return
	}
	if reqEvent.Event == "user.upgraded" {
		_, err = cfg.dbQueries.GetUserById(req.Context(), reqEvent.Data.UserID)
		if err != nil {
			respondWithError(respW, 404, err.Error())
			return
		}
		err = cfg.dbQueries.UpgradToRed(req.Context(), reqEvent.Data.UserID)
		if err != nil {
			respondWithError(respW, 404, err.Error())
			return
		}
	}
	respW.WriteHeader(204)
}

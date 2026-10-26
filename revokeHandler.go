package main

import (
	"net/http"

	"github.com/JtorresBear/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	refToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(respW, 401, err.Error())
		return
	}
	err = cfg.dbQueries.UpdateRevoked_at(req.Context(), refToken)
	if err != nil {
		respondWithError(respW, 401, err.Error())
		return
	}

	respW.WriteHeader(204)
}

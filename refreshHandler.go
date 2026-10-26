package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/JtorresBear/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	refToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(respW, 401, err.Error())
		return
	}
	tokenRow, err := cfg.dbQueries.GetTokenRow(req.Context(), refToken)

	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(respW, 401, "Token Doesn't Exist")
		return
	}
	if tokenRow.RevokedAt.Valid {
		respondWithError(respW, 401, "Token Is Revoked")
		return
	}
	if !time.Now().UTC().Before(tokenRow.ExpiresAt) {
		respondWithError(respW, 401, "Token Is Revoked")
		return
	}
	accessToken, err := auth.MakeJWT(tokenRow.UserID, cfg.secret)

	respondWithJson(respW, 200, struct {
		Token string `json:"token"`
	}{Token: accessToken})
}

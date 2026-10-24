package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/JtorresBear/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type authcred struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		Expiration int    `json:"expires_in_seconds"`
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(respW, 500, err.Error())
		return
	}
	cred := authcred{}

	err = json.Unmarshal(data, &cred)
	if err != nil {
		respondWithError(respW, 400, err.Error())
		return
	}
	if cred.Expiration < 1 || cred.Expiration > 3600 {
		cred.Expiration = 3600
	}
	actualCred, err := cfg.dbQueries.GetUserPassword(req.Context(), cred.Email)
	if err != nil {
		respondWithError(respW, 401, "incorrect email or password")
		return
	}
	match, err := auth.CheckPasswordHash(cred.Password, actualCred.HashedPassword)
	if err != nil {
		respondWithError(respW, 401, "incorrect email or password")
		return
	}
	if match != true {
		respondWithError(respW, 401, "incorrect email or password")
		return
	}
	token, err := auth.MakeJWT(actualCred.ID, cfg.secret, time.Duration(cred.Expiration)*time.Second)

	if err != nil {
		respondWithError(respW, 500, err.Error())
		return
	}
	user := User{
		ID:        actualCred.ID,
		CreatedAt: actualCred.CreatedAt,
		UpdatedAt: actualCred.UpdatedAt,
		Email:     actualCred.Email,
		Token:     token,
	}
	respondWithJson(respW, 200, user)
}

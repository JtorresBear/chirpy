package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/JtorresBear/chirpy/internal/auth"
	"github.com/JtorresBear/chirpy/internal/database"
)

func (cfg *apiConfig) loginHandler(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type authcred struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
	token, err := auth.MakeJWT(actualCred.ID, cfg.secret)

	if err != nil {
		respondWithError(respW, 500, err.Error())
		return
	}
	refreshToken := auth.MakeRefreshToken()

	refTokParams := database.CreateRefreshTokenParams{}
	refTokParams.Token = refreshToken
	refTokParams.UserID = actualCred.ID

	_, err = cfg.dbQueries.CreateRefreshToken(req.Context(), refTokParams)
	if err != nil {
		respondWithError(respW, 400, err.Error())
		return
	}

	user := User{
		ID:           actualCred.ID,
		CreatedAt:    actualCred.CreatedAt,
		UpdatedAt:    actualCred.UpdatedAt,
		Email:        actualCred.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}
	respondWithJson(respW, 200, user)
}

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/JtorresBear/chirpy/internal/auth"
	"github.com/JtorresBear/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) updatePasswordHandler(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(respW, 401, err.Error())
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(respW, 401, err.Error())
		return
	}
	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(respW, 500, "Could not read request")
		return
	}
	email_password := reqBody{}

	err = json.Unmarshal(data, &email_password)
	if err != nil {
		respondWithError(respW, 400, "Something wrong with data")
		return
	}
	params := database.UpdatePasswordParams{}
	pw, err := auth.HashPassword(email_password.Password)
	if err != nil {
		respondWithError(respW, 400, err.Error())
		return
	}
	params.HashedPassword = pw
	params.Email = email_password.Email
	params.ID = userID
	userCred, err := cfg.dbQueries.UpdatePassword(req.Context(), params)

	returnUser := struct {
		ID        uuid.UUID `json:"id"`
		EMAIL     string    `json:"email"`
		UPDATEDAT time.Time `json:"updated_at"`
		CREATEDAT time.Time `json:"created_at"`
	}{
		ID:        userCred.ID,
		EMAIL:     userCred.Email,
		UPDATEDAT: userCred.UpdatedAt,
		CREATEDAT: userCred.CreatedAt,
	}

	respondWithJson(respW, 200, returnUser)
}

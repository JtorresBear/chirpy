package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/JtorresBear/chirpy/internal/auth"
	"github.com/JtorresBear/chirpy/internal/database"
)

func (cfg *apiConfig) postUsersHandler(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type reqBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	//response body will be the user struct

	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(respW, 500, "Could not read request")
		return
	}
	emailPws := reqBody{}
	user := User{}

	err = json.Unmarshal(data, &emailPws)
	if err != nil {
		respondWithError(respW, 400, "Something wrong with the data")
		return
	}
	emailPws.Password, err = auth.HashPassword(emailPws.Password)
	if err != nil {
		respondWithError(respW, 500, err.Error())
	}
	dbUser, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          emailPws.Email,
		HashedPassword: emailPws.Password,
	})
	if err != nil {
		respondWithError(respW, 500, "Something went wrong with creating a user")
		return
	}
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	user.Email = dbUser.Email

	respondWithJson(respW, 201, user)
}

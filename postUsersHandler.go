package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func (cfg *apiConfig) postUsersHandler(respW http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	type reqBody struct {
		Email string `json:"email"`
	}
	//response body will be the user struct

	data, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithError(respW, 500, "Could not read request")
		return

	}
	email := reqBody{}
	user := User{}

	err = json.Unmarshal(data, &email)
	if err != nil {
		respondWithError(respW, 400, "Something wrong with the data")
		return
	}
	dbUser, err := cfg.dbQueries.CreateUser(req.Context(), email.Email)
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

package main

import "net/http"

func (cfg *apiConfig) resetHandler(respW http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(respW, 403, "FORBIDDEN")
		return
	}
	err := cfg.dbQueries.DeleteAllUsers(req.Context())
	if err != nil {
		respondWithError(respW, 500, "Couldn't delete all users")
	}
	cfg.fileserverHits.Store(0)
	respW.Header().Set("Content-Type", "text/plain; charset=utf-8")
	respW.WriteHeader(200)
	respW.Write([]byte("OK"))
}

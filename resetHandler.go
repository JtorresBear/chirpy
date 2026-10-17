package main

import "net/http"

func (cfg *apiConfig) resetHandler(respW http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	respW.Header().Set("Content-Type", "text/plain; charset=utf-8")
	respW.WriteHeader(200)
	respW.Write([]byte("OK"))
}

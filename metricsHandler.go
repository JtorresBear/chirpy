package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) metricsHandler(respW http.ResponseWriter, req *http.Request) {
	respW.Header().Set("Content-Type", "text/html; charset=utf-8")
	respW.WriteHeader(200)
	metrics := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %v times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	respW.Write([]byte(metrics))
}

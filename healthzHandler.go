package main

import "net/http"

func healthZHandler(respW http.ResponseWriter, req *http.Request) {
	respW.Header().Set("Content-Type", "text/plain; charset=utf-8")
	respW.WriteHeader(200)
	respW.Write([]byte("OK"))
}

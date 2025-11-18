package main

import "net/http"

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /", getHome())
	return mux
}

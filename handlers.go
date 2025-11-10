package main

import (
	"net/http"
	"time"
)

func getHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
		w.Write([]byte("home\n"))
	})
}

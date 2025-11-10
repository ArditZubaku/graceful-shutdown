package main

import (
	"fmt"
	"net/http"
)

func main() {
	s := http.Server{
		Addr:    ":8000",
		Handler: routes(),
	}

	fmt.Println("Listening on :8000")
	if err := s.ListenAndServe(); err != nil {
		fmt.Printf("Stopped listening: %v\n", err)
	}
}

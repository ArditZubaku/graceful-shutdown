package main

import (
	"fmt"
	"net/http"
	"time"
)

func getHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// go backgroundWork()
		// Increase the waitgroup/register the goroutine outside to avoid racing conditions
		wg.Add(1)
		go func() {
			defer wg.Done()
			backgroundWork()
		}()
		w.Write([]byte("home\n"))
	})
}

func backgroundWork() {
	fmt.Println("Background work started.")
	time.Sleep(10 * time.Second)
	fmt.Println("Background work ended.")
}

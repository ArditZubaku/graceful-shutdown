package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s := http.Server{
		Addr:    ":8000",
		Handler: routes(),
	}

	go func() {
		fmt.Println("Listening on :8000")
		if err := s.ListenAndServe(); err != nil {
			fmt.Printf("Stopped listening: %v\n", err)
		}
	}()

	shutdown, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer stop()

	// We are blocking until we receive a shutdown signal (CTRL+C for example)
	<-shutdown.Done()

	fmt.Println("Shutting server down...")
	// We give the server X seconds to shutdown gracefully or we do it forcefully
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// This starts the shutdown process, waits until all the handlers finish and then kills the server
	// Stops accepting new connections when in shutdown mode.
	if err := s.Shutdown(ctx); err != nil {
		fmt.Printf("Shut down with error: %v", err)
	} else {
		fmt.Println("Successfully shut down the server!")
	}
}

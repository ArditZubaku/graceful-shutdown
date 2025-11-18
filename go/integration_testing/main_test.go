package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"syscall"
	"testing"
	"time"
)

func TestServerGracefulShutdown(t *testing.T) {
	completed := "completed"

	server := &http.Server{
		Addr: ":54332",
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Second)
				_, err := w.Write([]byte(completed))
				if err != nil {
					log.Fatal("Write error:", err)
				}
			},
		),
	}

	ctx := t.Context()

	serverErrorCh := make(chan error, 1)
	serverStartedCh := make(chan struct{})

	go func() {
		err := runServer(ctx, server, 5*time.Second, serverStartedCh)
		if !errors.Is(err, http.ErrServerClosed) {
			serverErrorCh <- err
		}
		close(serverErrorCh)
	}()

	// start long request asynchronously
	resCh := make(chan *http.Response, 1)
	resErrCh := make(chan error, 1)
	go func() {
		<-serverStartedCh
		res, err := http.Get("http://localhost" + server.Addr)
		if err != nil {
			resErrCh <- err
			return
		}
		resCh <- res
	}()

	// give the handler a moment to enter its sleep
	time.Sleep(100 * time.Millisecond)

	// send SIGINT while handler still sleeping
	sysErr := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if sysErr != nil {
		t.Fatalf("Failed to send SIGINT: %v", sysErr)
	}

	select {
	case err := <-serverErrorCh:
		if err != nil {
			t.Fatalf("Expected no server error, got: %v", err)
		}

	case <-time.After(10 * time.Second):
		t.Fatal("Server did not shut down in time")

	case err := <-resErrCh:
		//cancel() // Signal server goroutine to exit
		t.Fatalf("HTTP GET request failed: %v", err)

	case res := <-resCh:
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Fatalf("Failed to close response body: %v", err)
			}
		}(res.Body)

		if res.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 OK, got %v", res.Status)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if string(body) != completed {
			t.Fatalf("Expected response body 'completed', got '%s'", string(body))
		}
	}

}

func TestServerTimeoutDuringShutdown(t *testing.T) {
	completed := "completed"

	server := &http.Server{
		Addr: ":54331",
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(10 * time.Second)
				_, err := w.Write([]byte(completed))
				if err != nil {
					log.Fatal("Write error:", err)
				}
			},
		),
	}

	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- runServer(
			context.Background(),
			server,
			5*time.Millisecond,
			make(chan struct{}),
		)
	}()

	reqErrCh := make(chan error, 1)
	go func() {
		_, err := http.Get("http://localhost" + server.Addr)
		reqErrCh <- err
	}()

	time.Sleep(1 * time.Second)

	sysErr := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if sysErr != nil {
		t.Fatalf("Failed to send SIGINT: %v", sysErr)
	}

	if <-reqErrCh == nil {
		t.Errorf("Expected request to fail due to server shutdown, but it succeeded")
	}

	serverErr := <-serverErrCh
	if !errors.Is(serverErr, context.DeadlineExceeded) {
		t.Errorf("Expected server error to be context.DeadlineExceeded, got: %v", serverErr)
	}
}

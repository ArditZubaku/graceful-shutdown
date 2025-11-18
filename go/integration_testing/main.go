package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server := createServer()

	if err := runServer(
		context.Background(),
		server,
		10*time.Second,
		make(chan struct{}),
	); err != nil {
		slog.Error(
			"Server error",
			slog.String("reason", err.Error()),
		)
	}
}

func createServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/slow",
		func(w http.ResponseWriter, r *http.Request) {
			slog.Info("Slow request started...")
			time.Sleep(8 * time.Second)
			n, err := fmt.Fprintf(
				w,
				"Slow request completed at: %v\n",
				time.Now(),
			)
			if err != nil || n == 0 {
				slog.Error(
					"Failed to write response",
					slog.String("reason", err.Error()),
				)
			}
		},
	)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func runServer(
	ctx context.Context,
	server *http.Server,
	shutdownTimeout time.Duration,
	serverStarted chan<- struct{},
) error {
	serverErrChan := make(chan error, 1)

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}
	close(serverStarted)

	go func() {
		slog.Info("Starting server...")
		if err := server.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
			serverErrChan <- err
		}
		close(serverErrChan)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-serverErrChan:
		return err
	case <-stop:
		slog.InfoContext(ctx, "Shutdown signal received")
	case <-ctx.Done():
		slog.InfoContext(ctx, "Operation canceled", slog.String("reason", ctx.Err().Error()))
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		if closeErr := server.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}

	slog.Info("Server exited gracefully!")
	return nil
}

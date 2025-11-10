# Graceful Shutdown

A simple Go HTTP server demonstrating graceful shutdown patterns with background goroutines.

## What it does

This application shows how to properly shut down a Go HTTP server while ensuring all background work completes before the application exits.

### Key Features

- **HTTP Server**: Listens on port 8000 with a single `/` endpoint
- **Background Work**: Each request spawns a 10-second background goroutine
- **Graceful Shutdown**: Handles SIGTERM/SIGINT signals (Ctrl+C) properly
- **Goroutine Tracking**: Uses `sync.WaitGroup` to wait for all background work to finish

## How it works

1. **Server Start**: HTTP server starts on `:8000`
2. **Request Handling**: GET requests to `/` trigger background work that takes 10 seconds
3. **Signal Handling**: Application listens for shutdown signals (Ctrl+C)
4. **Graceful Shutdown**:
   - Stops accepting new connections
   - Gives existing handlers up to 15 seconds to complete
   - Waits for all background goroutines to finish
   - Exits cleanly

## Running the application and testing the shutdown

1. Start the server: `go run .`
2. Make a request: `curl http://localhost:8000`
3. Immediately press Ctrl+C
4. Watch the server wait for background work to complete before exiting

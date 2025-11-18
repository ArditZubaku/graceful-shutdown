# Graceful Shutdown - Node.js Version

A Node.js HTTP server demonstrating graceful shutdown patterns with background work, mirroring the Go implementation.

## What it does

This application shows how to properly shut down a Node.js HTTP server while ensuring all background work completes before the application exits.

### Key Features

- **HTTP Server**: Listens on port 8000 with a single `/` endpoint
- **Background Work**: Each request starts background work that takes 10 seconds
- **Graceful Shutdown**: Handles SIGTERM/SIGINT signals (Ctrl+C) properly
- **Task Tracking**: Tracks active background tasks and waits for completion

## How it works

1. **Server Start**: HTTP server starts on `:8000`
2. **Request Handling**:
   - GET requests to `/` trigger background work that takes 10 seconds
   - Response is sent immediately while work continues in background
3. **Signal Handling**: Application listens for shutdown signals (Ctrl+C)
4. **Graceful Shutdown**:
   - Stops accepting new connections
   - Waits for all background tasks to complete
   - Times out after 15 seconds for forceful shutdown
   - Exits cleanly

## Running the application

```bash
# Install dependencies (none required for this basic version)
npm install

# Start the server
npm start
# or
node server.js

# For development with auto-restart
npm run dev
```

Visit `http://localhost:8000` to trigger background work, then press Ctrl+C to see the graceful shutdown in action.

## Key Differences from Go Version

| Aspect              | Go                                               | Node.js                              |
| ------------------- | ------------------------------------------------ | ------------------------------------ |
| **Concurrency**     | Goroutines with `sync.WaitGroup`                 | `setTimeout` with counter tracking   |
| **HTTP Server**     | `net/http` with automatic goroutines per request | `http.createServer` with event loop  |
| **Shutdown**        | `server.Shutdown()` with context timeout         | `server.close()` with manual timeout |
| **Background Work** | `time.Sleep(10 * time.Second)`                   | `setTimeout(() => {}, 10000)`        |

## Testing the shutdown

1. Start the server: `npm start`
2. Make a request: `curl http://localhost:8000`
3. Immediately press Ctrl+C
4. Watch the server wait for background work to complete before exiting

Both implementations demonstrate the same graceful shutdown pattern adapted to their respective runtime environments.

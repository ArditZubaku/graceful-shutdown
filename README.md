# Graceful Shutdown Implementations

This repository demonstrates graceful shutdown patterns in different programming languages and runtimes. The goal is to show how to properly handle application shutdown while ensuring all background work completes before the process exits.

## What is Graceful Shutdown?

Graceful shutdown is a pattern where an application:

1. **Stops accepting new requests** when a shutdown signal is received
2. **Completes all ongoing work** (background tasks, database operations, etc.)
3. **Exits cleanly** without data loss or corruption
4. **Has a timeout mechanism** to prevent hanging indefinitely

This is crucial for production applications to ensure data integrity and proper resource cleanup.

## Implementations

### Go Implementation

- **Features**: Uses goroutines and `sync.WaitGroup` for concurrent background work
- **Shutdown**: Leverages `http.Server.Shutdown()` with context timeout
- **Documentation**: [Go README](go/README.md)

### Node.js Implementation

- **Features**: Uses event loop and `setTimeout` for asynchronous background work
- **Shutdown**: Manual tracking with `server.close()` and timeout handling
- **Documentation**: [Node.js README](nodejs/README.md)

## Common Pattern

Both implementations follow the same high-level pattern:

1. **HTTP Server** listening on port 8000 with a `/` endpoint
2. **Background Work** triggered by each request (10-second simulation)
3. **Signal Handling** for SIGTERM/SIGINT (Ctrl+C)
4. **Graceful Shutdown** with 15-second timeout for forceful exit

## Key Differences

| Aspect                | Go                            | Node.js                      |
| --------------------- | ----------------------------- | ---------------------------- |
| **Concurrency Model** | Goroutines (true parallelism) | Event Loop (single-threaded) |
| **Background Work**   | `time.Sleep()` in goroutine   | `setTimeout()` callback      |
| **Work Tracking**     | `sync.WaitGroup`              | Manual counter               |
| **Server Shutdown**   | Built-in `Shutdown()` method  | Manual `close()` + tracking  |

## Testing Both Implementations

```bash
# Test Go version
cd go/
go run .
# In another terminal: curl http://localhost:8000
# Press Ctrl+C to see graceful shutdown

# Test Node.js version
cd nodejs/
node server.js
# In another terminal: curl http://localhost:8000
# Press Ctrl+C to see graceful shutdown
```

## Use Cases

This pattern is essential for:

- **Web servers** handling background processing
- **Microservices** with database connections
- **Queue workers** processing jobs
- **Any application** that needs clean shutdown for data integrity

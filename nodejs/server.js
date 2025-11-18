const http = require("http");
const url = require("url");

// Track active background work
let activeBackgroundTasks = 0;
let isShuttingDown = false;

// Create HTTP server
const server = http.createServer((req, res) => {
  const parsedUrl = url.parse(req.url, true);

  if (req.method === "GET" && parsedUrl.pathname === "/") {
    handleHome(req, res);
  } else {
    res.writeHead(404, { "Content-Type": "text/plain" });
    res.end("Not Found\n");
  }
});

function handleHome(_req, res) {
  // Start background work (similar to Go's goroutine)
  startBackgroundWork();

  // Respond immediately (like Go handler)
  res.writeHead(200, { "Content-Type": "text/plain" });
  res.end("home\n");
}

function startBackgroundWork() {
  activeBackgroundTasks++;
  console.log(
    `Background work started. Active tasks: ${activeBackgroundTasks}`,
  );

  // Simulate 10-second background work (like Go's time.Sleep)
  setTimeout(() => {
    activeBackgroundTasks--;
    console.log(
      `Background work ended. Active tasks: ${activeBackgroundTasks}`,
    );

    // If we're shutting down and this was the last task, exit
    if (isShuttingDown && activeBackgroundTasks === 0) {
      console.log("All background tasks finished. Exiting...");
      process.exit(0);
    }
  }, 10000);
}

// Start server
const PORT = 8000;
server.listen(PORT, () => {
  console.log(`Listening on :${PORT}`);
});

// Graceful shutdown handling
process.on("SIGTERM", gracefulShutdown);
process.on("SIGINT", gracefulShutdown);

function gracefulShutdown(signal) {
  console.log(`\nReceived ${signal}. Shutting server down...`);
  isShuttingDown = true;

  // Stop accepting new connections (like Go's s.Shutdown())
  server.close(() => {
    console.log("Server stopped accepting new connections");

    if (activeBackgroundTasks === 0) {
      console.log("No background tasks running. Exiting...");
      process.exit(0);
    } else {
      console.log(
        `Waiting for ${activeBackgroundTasks} background tasks to finish...`,
      );

      // Set a timeout for forceful shutdown (like Go's context.WithTimeout)
      setTimeout(() => {
        console.log("Forceful shutdown after timeout");
        process.exit(1);
      }, 15000);
    }
  });
}


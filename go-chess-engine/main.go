package main

import (
	"io"
	"log"
	"os"

	"go-chess-engine/uci"
)

func main() {
	// 1. Set up the logger first and foremost.
	logger := setupLogger("uci.log")
	logger.Println("--- Engine Started ---")

	// 2. Pass the logger to the handler.
	handler := uci.NewHandler(logger)
	handler.Loop()

	logger.Println("--- Engine Quit ---")
}

// setupLogger configures and returns a new logger instance.
func setupLogger(filename string) *log.Logger {
	// Open the log file. O_TRUNC will clear the file on each start.
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		// If we can't create a log file, something is very wrong.
		// Fallback to standard error.
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create a multi-writer to log to both file and console (optional, but great for debugging)
	multiWriter := io.MultiWriter(os.Stderr, logFile)

	// Create a new logger. LstdFlags adds date and time to each log entry.
	// You can remove os.Stderr from multiWriter if you only want to log to the file.
	return log.New(multiWriter, "", log.LstdFlags)
}

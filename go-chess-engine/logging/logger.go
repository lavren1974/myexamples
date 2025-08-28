package logging

import (
	"log"
	"os"
)

// Log is the global logger instance for the entire application.
var Log *log.Logger

// init() is a special Go function that runs before main().
// This is the perfect place to set up our logger, guaranteeing
// it's ready before any other part of the application runs.
func init() {
	// Open the log file. O_TRUNC will clear the file on each start.
	logFile, err := os.OpenFile("uci.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		// If we can't create a log file, we have to panic.
		log.Fatalf("FATAL: Failed to open log file: %v", err)
	}

	// Initialize the global logger.
	Log = log.New(logFile, "", log.Ltime)
}

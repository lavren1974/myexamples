package main

import (
	// Import our new logging package. Its init() function will run automatically.
	"go-chess-engine/logging"
	"go-chess-engine/uci"
)

func main() {
	// The logger is already set up by the logging package's init() function.
	logging.Log.Println("--- Engine Started ---")

	// Create and run the UCI handler. We no longer need to pass the logger.
	handler := uci.NewHandler()
	handler.Loop()

	logging.Log.Println("--- Engine Quit ---")
}

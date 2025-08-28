package uci

import (
	"bufio"
	"fmt"
	"go-chess-engine/chess"
	"go-chess-engine/engine"
	"go-chess-engine/logging" // Import the logging package
	"os"
	"strings"
)

// Handler no longer needs a logger field.
type Handler struct {
	state  *chess.State
	engine *engine.Engine
}

// NewHandler is now simpler.
func NewHandler() *Handler {
	return &Handler{
		state:  chess.New(),
		engine: engine.New(),
	}
}

// Loop now uses the global logger.
func (h *Handler) Loop() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		logging.Log.Printf("Received: %s", command) // Use global logger

		fields := strings.Fields(command)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "uci":
			h.handleUci()
		case "isready":
			h.handleIsReady()
		case "ucinewgame":
			h.handleUciNewGame()
		case "position":
			h.handlePosition(fields)
		case "go":
			h.handleGo()
		case "quit":
			return
		}
	}
}

func (h *Handler) handleUci() {
	h.sendResponse("id name GoNativeRefactored")
	h.sendResponse("id author Go Developer")
	h.sendResponse("uciok")
}

func (h *Handler) handleIsReady() {
	h.sendResponse("readyok")
}

// ... (handleUciNewGame and handlePosition are unchanged) ...
func (h *Handler) handleUciNewGame() {
	h.state = chess.New()
}

func (h *Handler) handlePosition(fields []string) {
	var movesIndex = -1
	if len(fields) > 1 && fields[1] == "startpos" {
		h.state = chess.New()
		movesIndex = 2
	} else if len(fields) > 2 && fields[1] == "fen" {
		fenStr := ""
		for i := 2; i < len(fields); i++ {
			if fields[i] == "moves" {
				movesIndex = i
				break
			}
			fenStr += fields[i] + " "
		}
		h.state = chess.MustParseFEN(strings.TrimSpace(fenStr))
	}

	if movesIndex != -1 && movesIndex+1 < len(fields) {
		for i := movesIndex + 1; i < len(fields); i++ {
			move := chess.ParseMove(fields[i])
			h.state.ApplyMove(move)
		}
	}
}

func (h *Handler) handleGo() {
	bestMove := h.engine.FindBestMove(h.state)
	h.sendResponse(fmt.Sprintf("bestmove %s", chess.FormatMove(bestMove)))
}

// sendResponse now uses the global logger.
func (h *Handler) sendResponse(msg string) {
	logging.Log.Printf("Sending: %s", msg) // Use global logger
	fmt.Println(msg)
}

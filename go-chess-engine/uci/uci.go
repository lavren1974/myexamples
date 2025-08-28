package uci

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"go-chess-engine/chess"
	"go-chess-engine/engine"
)

// Handler manages the UCI communication session.
type Handler struct {
	state  *chess.State
	engine *engine.Engine
	logger *log.Logger // Add a logger field
}

// NewHandler now accepts a logger.
func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		state:  chess.New(),
		engine: engine.New(),
		logger: logger, // Store the logger
	}
}

// Loop is the main UCI command loop.
func (h *Handler) Loop() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		h.logger.Printf("Received: %s", command) // Use the handler's logger

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

func (h *Handler) handleUciNewGame() {
	h.state = chess.New()
}

func (h *Handler) handlePosition(fields []string) {
	// ... (this function's code does not need to change)
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

// sendResponse now belongs to the handler so it can access the logger.
func (h *Handler) sendResponse(msg string) {
	h.logger.Printf("Sending: %s", msg)
	fmt.Println(msg)
}

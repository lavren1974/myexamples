package uci

import (
	"bufio"
	"fmt"
	"go-chess-engine/chess"
	"go-chess-engine/engine"
	"go-chess-engine/logging"
	"os"
	"strings"
)

type Handler struct {
	board  chess.Board // This is now the INTERFACE, not a concrete type
	engine *engine.Engine
}

func NewHandler() *Handler {
	return &Handler{
		// Use the factory to create the board from the starting position
		board:  chess.NewBoardFromConfig(chess.StartFEN),
		engine: engine.New(),
	}
}

func (h *Handler) Loop() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		logging.Log.Printf("Received: %s", command)

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
	// Re-create the board from the starting position
	h.board = chess.NewBoardFromConfig(chess.StartFEN)
}

func (h *Handler) handlePosition(fields []string) {
	var movesIndex = -1
	var fen = chess.StartFEN

	if len(fields) > 1 && fields[1] == "startpos" {
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
		fen = strings.TrimSpace(fenStr)
	}

	// Create the new board from the specified FEN
	h.board = chess.NewBoardFromConfig(fen)

	// Apply moves if they are provided
	if movesIndex != -1 && movesIndex+1 < len(fields) {
		for i := movesIndex + 1; i < len(fields); i++ {
			move := chess.ParseMove(fields[i])
			h.board.ApplyMove(move)
		}
	}
}

func (h *Handler) handleGo() {
	// The engine needs to receive the board interface
	bestMove := h.engine.FindBestMove(h.board)
	h.sendResponse(fmt.Sprintf("bestmove %s", chess.FormatMove(bestMove)))
}

// This is the corrected function signature.
// It now correctly has the (h *Handler) receiver.
func (h *Handler) sendResponse(msg string) {
	logging.Log.Printf("Sending: %s", msg)
	fmt.Println(msg)
}

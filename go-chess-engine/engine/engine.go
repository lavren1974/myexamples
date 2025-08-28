package engine

import (
	"go-chess-engine/chess"
	"math/rand"
	"time"
)

type Engine struct {
	// Future fields: transposition tables, search settings, etc.
}

func New() *Engine {
	return &Engine{}
}

// FindBestMove now uses the legal move generator.
func (e *Engine) FindBestMove(s *chess.State) chess.Move {
	// Call the new legal move generator
	moves := s.GenerateLegalMoves()

	if len(moves) == 0 {
		// If there are no legal moves, the game is over (checkmate or stalemate).
		// Return a null move.
		return chess.Move{}
	}

	rand.Seed(time.Now().UnixNano())
	return moves[rand.Intn(len(moves))]
}

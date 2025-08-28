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

// FindBestMove is the main thinking function.
// For now, it just picks a random move.
func (e *Engine) FindBestMove(s *chess.State) chess.Move {
	moves := s.GenerateMoves()
	if len(moves) == 0 {
		return chess.Move{} // Return a null move
	}

	rand.Seed(time.Now().UnixNano())
	return moves[rand.Intn(len(moves))]
}

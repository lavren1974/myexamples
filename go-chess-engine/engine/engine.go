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
// FindBestMove now accepts the Board INTERFACE
func (e *Engine) FindBestMove(b chess.Board) chess.Move {
	moves := b.GenerateLegalMoves() // This call works on both ArrayBoard and Bitboard!
	if len(moves) == 0 {
		return chess.Move{}
	}
	rand.Seed(time.Now().UnixNano())
	return moves[rand.Intn(len(moves))]
}

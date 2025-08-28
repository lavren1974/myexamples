package chess

import (
	"go-chess-engine/config"
	"log"
)

// NewBoardFromConfig is a factory function that creates a board
// based on the global application configuration.
func NewBoardFromConfig(fen string) Board {
	rep := config.AppConfig.BoardRepresentation
	log.Printf("INFO: Creating new board with representation: '%s'", rep)

	switch rep {
	case "bitboard":
		return NewBitboard(fen)
	case "array":
		return NewArrayBoard(fen)
	default:
		log.Printf("WARN: Unknown board representation '%s'. Defaulting to 'array'.", rep)
		return NewArrayBoard(fen)
	}
}

package chess

import "strings"

// State holds all information about the current state of the game.
type State struct {
	Board      [64]Piece
	SideToMove Color
	// Castling rights, en passant target, etc., would go here.
}

// New creates a new game state from the standard starting position.
func New() *State {
	return MustParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
}

// ApplyMove updates the board state by making a move.
func (s *State) ApplyMove(m Move) {
	s.Board[m.To] = s.Board[m.From]
	s.Board[m.From] = Empty
	if s.SideToMove == White {
		s.SideToMove = Black
	} else {
		s.SideToMove = White
	}
}

// MustParseFEN is a helper that panics if the FEN is invalid.
func MustParseFEN(fen string) *State {
	s := &State{}
	fields := strings.Fields(fen)

	// Piece placement
	rank, file := 7, 0
	for _, char := range fields[0] {
		if char == '/' {
			rank--
			file = 0
		} else if char >= '1' && char <= '8' {
			file += int(char - '0')
		} else {
			square := rank*8 + file
			s.Board[square] = pieceFromChar(char)
			file++
		}
	}

	// Active color
	if fields[1] == "w" {
		s.SideToMove = White
	} else {
		s.SideToMove = Black
	}
	return s
}

func pieceFromChar(c rune) Piece {
	switch c {
	case 'p':
		return BlackPawn
	case 'n':
		return BlackKnight
	case 'b':
		return BlackBishop
	case 'r':
		return BlackRook
	case 'q':
		return BlackQueen
	case 'k':
		return BlackKing
	case 'P':
		return WhitePawn
	case 'N':
		return WhiteKnight
	case 'B':
		return WhiteBishop
	case 'R':
		return WhiteRook
	case 'Q':
		return WhiteQueen
	case 'K':
		return WhiteKing
	}
	return Empty
}

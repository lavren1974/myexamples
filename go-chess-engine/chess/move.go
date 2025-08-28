package chess

import "strconv"

type Color int8
type Piece int8
type Move struct {
	From, To int
}

const (
	White Color = 0
	Black Color = 1
)

const (
	Empty Piece = iota
	WhitePawn
	WhiteKnight
	WhiteBishop
	WhiteRook
	WhiteQueen
	WhiteKing
	BlackPawn
	BlackKnight
	BlackBishop
	BlackRook
	BlackQueen
	BlackKing
)

// FormatMove converts a Move struct to a UCI-compliant string (e.g., "e2e4").
func FormatMove(m Move) string {
	return indexToSquare(m.From) + indexToSquare(m.To)
}

// ParseMove converts a UCI string to a Move struct.
func ParseMove(s string) Move {
	return Move{
		From: squareToIndex(s[0:2]),
		To:   squareToIndex(s[2:4]),
	}
}

// --- Private utility functions ---

func squareToIndex(s string) int {
	file := int(s[0] - 'a')
	rank, _ := strconv.Atoi(string(s[1]))
	return (rank-1)*8 + file
}

func indexToSquare(i int) string {
	file := 'a' + rune(i%8)
	rank := '1' + rune(i/8)
	return string(file) + string(rank)
}

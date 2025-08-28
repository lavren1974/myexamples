package chess

import "strconv"

type Color int8
type Piece int8
type Move struct {
	From, To int
}

const (
	White   Color = 0
	Black   Color = 1
	NoColor Color = 2 // Useful for error checking
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

// PieceColor returns the color of a piece. Returns NoColor for Empty.
func (p Piece) Color() Color {
	if p >= WhitePawn && p <= WhiteKing {
		return White
	}
	if p >= BlackPawn && p <= BlackKing {
		return Black
	}
	return NoColor
}

// ... (rest of the file remains the same)
// FormatMove, ParseMove, squareToIndex, indexToSquare
func FormatMove(m Move) string {
	return indexToSquare(m.From) + indexToSquare(m.To)
}

func ParseMove(s string) Move {
	return Move{
		From: squareToIndex(s[0:2]),
		To:   squareToIndex(s[2:4]),
	}
}

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

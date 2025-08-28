package chess

import "strconv"

type Color int8
type Piece int8
type Move struct {
	From      int
	To        int
	Promotion Piece
}

const (
	White   Color = 0
	Black   Color = 1
	NoColor Color = 2
)

// This is the corrected constant block.
const (
	Empty       Piece = iota // 0
	WhitePawn                // 1
	WhiteKnight              // 2
	WhiteBishop              // 3
	WhiteRook                // 4
	WhiteQueen               // 5
	WhiteKing                // 6
	BlackPawn                // 7
	BlackKnight              // 8
	BlackBishop              // 9
	BlackRook                // 10
	BlackQueen               // 11
	BlackKing                // 12
)

func (p Piece) Color() Color {
	if p >= WhitePawn && p <= WhiteKing {
		return White
	}
	if p >= BlackPawn && p <= BlackKing {
		return Black
	}
	return NoColor
}

func FormatMove(m Move) string {
	baseMove := indexToSquare(m.From) + indexToSquare(m.To)
	if m.Promotion != Empty {
		switch m.Promotion {
		case WhiteQueen, BlackQueen:
			baseMove += "q"
		case WhiteRook, BlackRook:
			baseMove += "r"
		case WhiteBishop, BlackBishop:
			baseMove += "b"
		case WhiteKnight, BlackKnight:
			baseMove += "n"
		}
	}
	return baseMove
}

func ParseMove(s string) Move {
	move := Move{
		From:      squareToIndex(s[0:2]),
		To:        squareToIndex(s[2:4]),
		Promotion: Empty,
	}
	if len(s) == 5 {
		var pieceColor Color
		if s[3] == '8' {
			pieceColor = White
		} else {
			pieceColor = Black
		}

		promoChar := s[4]
		switch promoChar {
		case 'q':
			if pieceColor == White {
				move.Promotion = WhiteQueen
			} else {
				move.Promotion = BlackQueen
			}
		case 'r':
			if pieceColor == White {
				move.Promotion = WhiteRook
			} else {
				move.Promotion = BlackRook
			}
		case 'b':
			if pieceColor == White {
				move.Promotion = WhiteBishop
			} else {
				move.Promotion = BlackBishop
			}
		case 'n':
			if pieceColor == White {
				move.Promotion = WhiteKnight
			} else {
				move.Promotion = BlackKnight
			}
		}
	}
	return move
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

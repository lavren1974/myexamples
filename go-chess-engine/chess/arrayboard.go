package chess

import "strings"

// ArrayBoard is our original implementation using a simple array.
type ArrayBoard struct {
	Board           [64]Piece
	sideToMove      Color
	whiteKingSquare int
	blackKingSquare int
	// New fields for castling rights
	whiteKingsideCastle  bool
	whiteQueensideCastle bool
	blackKingsideCastle  bool
	blackQueensideCastle bool
}

// NewArrayBoard creates a new board from a FEN string.
func NewArrayBoard(fen string) *ArrayBoard {
	b := &ArrayBoard{}
	fields := strings.Fields(fen)

	// 1. Piece placement (and find kings)
	rank, file := 7, 0
	for _, char := range fields[0] {
		if char == '/' {
			rank--
			file = 0
		} else if char >= '1' && char <= '8' {
			file += int(char - '0')
		} else {
			square := rank*8 + file
			piece := pieceFromChar(char)
			b.Board[square] = piece
			if piece == WhiteKing {
				b.whiteKingSquare = square
			} else if piece == BlackKing {
				b.blackKingSquare = square
			}
			file++
		}
	}

	// 2. Side to move
	if fields[1] == "w" {
		b.sideToMove = White
	} else {
		b.sideToMove = Black
	}

	// 3. Parse castling rights
	if len(fields) > 2 {
		for _, char := range fields[2] {
			switch char {
			case 'K':
				b.whiteKingsideCastle = true
			case 'Q':
				b.whiteQueensideCastle = true
			case 'k':
				b.blackKingsideCastle = true
			case 'q':
				b.blackQueensideCastle = true
			}
		}
	}
	// Note: en passant and move counters are still ignored for now.
	return b
}

// --- Methods to satisfy the Board interface ---

// Replace your existing ApplyMove with this new version.
func (b *ArrayBoard) ApplyMove(m Move) {
	piece := b.Board[m.From]

	// --- Handle the actual move ---
	// If it's a castling move (king moving two squares), we must also move the rook.
	if (piece == WhiteKing || piece == BlackKing) && dist(m.From%8, m.To%8) == 2 {
		// Kingside castle
		if m.To%8 > m.From%8 {
			rookFrom, rookTo := m.To+1, m.To-1
			b.Board[rookTo] = b.Board[rookFrom]
			b.Board[rookFrom] = Empty
		} else { // Queenside castle
			rookFrom, rookTo := m.To-2, m.To+1
			b.Board[rookTo] = b.Board[rookFrom]
			b.Board[rookFrom] = Empty
		}
	}

	// Standard piece placement (including promotion)
	if m.Promotion != Empty {
		b.Board[m.To] = m.Promotion
	} else {
		b.Board[m.To] = piece
	}
	b.Board[m.From] = Empty

	// --- Update state after the move ---

	// 1. Update king's position if it moved
	if piece == WhiteKing {
		b.whiteKingSquare = m.To
	} else if piece == BlackKing {
		b.blackKingSquare = m.To
	}

	// 2. Revoke castling rights if a king or rook moves for the first time
	if piece == WhiteKing {
		b.whiteKingsideCastle = false
		b.whiteQueensideCastle = false
	} else if piece == BlackKing {
		b.blackKingsideCastle = false
		b.blackQueensideCastle = false
	} else if m.From == 0 { // a1 rook
		b.whiteQueensideCastle = false
	} else if m.From == 7 { // h1 rook
		b.whiteKingsideCastle = false
	} else if m.From == 56 { // a8 rook
		b.blackQueensideCastle = false
	} else if m.From == 63 { // h8 rook
		b.blackKingsideCastle = false
	}

	// 3. Switch side to move
	if b.sideToMove == White {
		b.sideToMove = Black
	} else {
		b.sideToMove = White
	}
}

func (b *ArrayBoard) SideToMove() Color {
	return b.sideToMove
}

func (b *ArrayBoard) IsCheckmate() bool {
	if len(b.GenerateLegalMoves()) == 0 {
		kingSq := b.blackKingSquare
		if b.sideToMove == White {
			kingSq = b.whiteKingSquare
		}
		return b.isSquareAttacked(kingSq, oppositeColor(b.sideToMove))
	}
	return false
}

func (b *ArrayBoard) IsStalemate() bool {
	if len(b.GenerateLegalMoves()) == 0 {
		kingSq := b.blackKingSquare
		if b.sideToMove == White {
			kingSq = b.whiteKingSquare
		}
		return !b.isSquareAttacked(kingSq, oppositeColor(b.sideToMove))
	}
	return false
}

// --- Move Generation Logic (Moved from movegen.go) ---

var (
	rookDirections   = []int{-8, -1, 1, 8}
	bishopDirections = []int{-9, -7, 7, 9}
	knightOffsets    = []int{-17, -15, -10, -6, 6, 10, 15, 17}
	kingOffsets      = []int{-9, -8, -7, -1, 1, 7, 8, 9}
)

func (b *ArrayBoard) GenerateLegalMoves() []Move {
	var legalMoves []Move
	pseudoLegalMoves := b.generatePseudoLegalMoves()
	for _, move := range pseudoLegalMoves {
		tempState := *b
		tempState.ApplyMove(move)
		var kingSquare int
		if b.sideToMove == White {
			kingSquare = tempState.whiteKingSquare
		} else {
			kingSquare = tempState.blackKingSquare
		}
		if !tempState.isSquareAttacked(kingSquare, tempState.sideToMove) {
			legalMoves = append(legalMoves, move)
		}
	}
	return legalMoves
}

func (b *ArrayBoard) isSquareAttacked(sq int, byColor Color) bool {
	if byColor == White {
		if sq/8 > 0 {
			if sq%8 > 0 {
				if b.Board[sq-9] == WhitePawn {
					return true
				}
			}
			if sq%8 < 7 {
				if b.Board[sq-7] == WhitePawn {
					return true
				}
			}
		}
	} else {
		if sq/8 < 7 {
			if sq%8 > 0 {
				if b.Board[sq+7] == BlackPawn {
					return true
				}
			}
			if sq%8 < 7 {
				if b.Board[sq+9] == BlackPawn {
					return true
				}
			}
		}
	}
	for _, offset := range knightOffsets {
		targetSq := sq + offset
		if targetSq >= 0 && targetSq < 64 && dist(sq%8, targetSq%8) <= 2 {
			piece := b.Board[targetSq]
			if piece.Color() == byColor && (piece == WhiteKnight || piece == BlackKnight) {
				return true
			}
		}
	}
	for _, offset := range kingOffsets {
		targetSq := sq + offset
		if targetSq >= 0 && targetSq < 64 && dist(sq%8, targetSq%8) <= 1 {
			piece := b.Board[targetSq]
			if piece.Color() == byColor && (piece == WhiteKing || piece == BlackKing) {
				return true
			}
		}
	}
	for _, dir := range rookDirections {
		for targetSq := sq + dir; ; targetSq += dir {
			if targetSq < 0 || targetSq >= 64 || dist((targetSq-dir)%8, targetSq%8) > 1 {
				break
			}
			piece := b.Board[targetSq]
			if piece != Empty {
				if piece.Color() == byColor && (piece == WhiteRook || piece == BlackRook || piece == WhiteQueen || piece == BlackQueen) {
					return true
				}
				break
			}
		}
	}
	for _, dir := range bishopDirections {
		for targetSq := sq + dir; ; targetSq += dir {
			if targetSq < 0 || targetSq >= 64 || dist((targetSq-dir)%8, targetSq%8) != 1 {
				break
			}
			piece := b.Board[targetSq]
			if piece != Empty {
				if piece.Color() == byColor && (piece == WhiteBishop || piece == BlackBishop || piece == WhiteQueen || piece == BlackQueen) {
					return true
				}
				break
			}
		}
	}
	return false
}

func (b *ArrayBoard) generatePseudoLegalMoves() []Move {
	var moves []Move
	for from := 0; from < 64; from++ {
		piece := b.Board[from]
		if piece.Color() != b.sideToMove {
			continue
		}
		switch piece {
		case WhitePawn, BlackPawn:
			moves = append(moves, b.generatePawnMoves(from)...)
		case WhiteKnight, BlackKnight:
			moves = append(moves, b.generateKnightMoves(from)...)
		case WhiteBishop, BlackBishop:
			moves = append(moves, b.generateSlidingMoves(from, bishopDirections)...)
		case WhiteRook, BlackRook:
			moves = append(moves, b.generateSlidingMoves(from, rookDirections)...)
		case WhiteQueen, BlackQueen:
			moves = append(moves, b.generateSlidingMoves(from, bishopDirections)...)
			moves = append(moves, b.generateSlidingMoves(from, rookDirections)...)
		case WhiteKing, BlackKing:
			moves = append(moves, b.generateKingMoves(from)...)
		}
	}
	return moves
}

func (b *ArrayBoard) generateSlidingMoves(from int, directions []int) []Move {
	var moves []Move
	isRook := b.Board[from] == WhiteRook || b.Board[from] == BlackRook
	for _, dir := range directions {
		prevSquare := from
		for {
			to := prevSquare + dir
			if to < 0 || to >= 64 {
				break
			}
			rankDist, fileDist := dist(to/8, prevSquare/8), dist(to%8, prevSquare%8)
			if isRook {
				if rankDist > 1 || fileDist > 1 || (rankDist == 1 && fileDist == 1) {
					break
				}
			} else {
				if rankDist != 1 || fileDist != 1 {
					break
				}
			}
			targetPiece := b.Board[to]
			if targetPiece == Empty {
				moves = append(moves, Move{From: from, To: to})
			} else {
				if targetPiece.Color() != b.sideToMove {
					moves = append(moves, Move{From: from, To: to})
				}
				break
			}
			prevSquare = to
		}
	}
	return moves
}

func (b *ArrayBoard) generatePawnMoves(from int) []Move {
	var moves []Move
	var pushDir, startRank, promotionRank int
	var captureDirs []int
	if b.sideToMove == White {
		pushDir, startRank, promotionRank = 8, 1, 7
		captureDirs = []int{7, 9}
	} else {
		pushDir, startRank, promotionRank = -8, 6, 0
		captureDirs = []int{-7, -9}
	}
	oneStep := from + pushDir
	if oneStep >= 0 && oneStep < 64 && b.Board[oneStep] == Empty {
		isPromotion := (oneStep / 8) == promotionRank
		b.addPawnMove(&moves, from, oneStep, isPromotion)
		if from/8 == startRank {
			twoSteps := from + 2*pushDir
			if b.Board[twoSteps] == Empty {
				b.addPawnMove(&moves, from, twoSteps, false)
			}
		}
	}
	for _, capDir := range captureDirs {
		to := from + capDir
		if to < 0 || to >= 64 || dist(from%8, to%8) != 1 {
			continue
		}
		targetPiece := b.Board[to]
		if targetPiece != Empty && targetPiece.Color() != b.sideToMove {
			isPromotion := (to / 8) == promotionRank
			b.addPawnMove(&moves, from, to, isPromotion)
		}
	}
	return moves
}

func (b *ArrayBoard) addPawnMove(moves *[]Move, from, to int, isPromotion bool) {
	if isPromotion {
		if b.sideToMove == White {
			*moves = append(*moves, Move{From: from, To: to, Promotion: WhiteQueen})
			*moves = append(*moves, Move{From: from, To: to, Promotion: WhiteRook})
			*moves = append(*moves, Move{From: from, To: to, Promotion: WhiteBishop})
			*moves = append(*moves, Move{From: from, To: to, Promotion: WhiteKnight})
		} else {
			*moves = append(*moves, Move{From: from, To: to, Promotion: BlackQueen})
			*moves = append(*moves, Move{From: from, To: to, Promotion: BlackRook})
			*moves = append(*moves, Move{From: from, To: to, Promotion: BlackBishop})
			*moves = append(*moves, Move{From: from, To: to, Promotion: BlackKnight})
		}
	} else {
		*moves = append(*moves, Move{From: from, To: to, Promotion: Empty})
	}
}

func (b *ArrayBoard) generateKnightMoves(from int) []Move {
	var moves []Move
	for _, offset := range knightOffsets {
		to := from + offset
		if to < 0 || to >= 64 {
			continue
		}
		if dist(from%8, to%8) > 2 || dist(from/8, to/8) > 2 {
			continue
		}
		targetPiece := b.Board[to]
		if targetPiece == Empty || targetPiece.Color() != b.sideToMove {
			moves = append(moves, Move{From: from, To: to})
		}
	}
	return moves
}

// Replace your existing generateKingMoves with this new version.
func (b *ArrayBoard) generateKingMoves(from int) []Move {
	var moves []Move

	// 1. Standard King Moves
	for _, offset := range kingOffsets {
		to := from + offset
		if to < 0 || to >= 64 || dist(from%8, to%8) > 1 {
			continue
		}
		targetPiece := b.Board[to]
		if targetPiece == Empty || targetPiece.Color() != b.sideToMove {
			moves = append(moves, Move{From: from, To: to})
		}
	}

	// 2. Castling Moves
	opponentColor := oppositeColor(b.sideToMove)
	// Don't generate castling moves if the king is currently in check
	if b.isSquareAttacked(from, opponentColor) {
		return moves
	}

	if b.sideToMove == White {
		// Kingside (O-O)
		if b.whiteKingsideCastle && b.Board[5] == Empty && b.Board[6] == Empty {
			if !b.isSquareAttacked(5, opponentColor) && !b.isSquareAttacked(6, opponentColor) {
				moves = append(moves, Move{From: from, To: 6})
			}
		}
		// Queenside (O-O-O)
		if b.whiteQueensideCastle && b.Board[1] == Empty && b.Board[2] == Empty && b.Board[3] == Empty {
			if !b.isSquareAttacked(2, opponentColor) && !b.isSquareAttacked(3, opponentColor) {
				moves = append(moves, Move{From: from, To: 2})
			}
		}
	} else { // Black's turn
		// Kingside (O-O)
		if b.blackKingsideCastle && b.Board[61] == Empty && b.Board[62] == Empty {
			if !b.isSquareAttacked(61, opponentColor) && !b.isSquareAttacked(62, opponentColor) {
				moves = append(moves, Move{From: from, To: 62})
			}
		}
		// Queenside (O-O-O)
		if b.blackQueensideCastle && b.Board[57] == Empty && b.Board[58] == Empty && b.Board[59] == Empty {
			if !b.isSquareAttacked(58, opponentColor) && !b.isSquareAttacked(59, opponentColor) {
				moves = append(moves, Move{From: from, To: 58})
			}
		}
	}
	return moves
}

// --- Helper functions ---

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

func oppositeColor(c Color) Color {
	if c == White {
		return Black
	}
	return White
}
func dist(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

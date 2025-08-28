package chess

import (
	"strings"
)

// Bitboard is the new implementation using bitboards.
type Bitboard struct {
	byPiece         [13]bitboard
	byColor         [2]bitboard
	sideToMove      Color
	whiteKingSquare int
	blackKingSquare int
}

// NewBitboard creates a bitboard representation from a FEN string.
func NewBitboard(fen string) *Bitboard {
	b := &Bitboard{}
	fields := strings.Fields(fen)
	rank, file := 7, 0
	for _, char := range fields[0] {
		if char == '/' {
			rank--
			file = 0
		} else if char >= '1' && char <= '8' {
			file += int(char - '0')
		} else {
			sq := rank*8 + file
			piece := pieceFromChar(char)
			b.byPiece[piece].setBit(sq)
			if piece.Color() != NoColor {
				b.byColor[piece.Color()].setBit(sq)
			}
			if piece == WhiteKing {
				b.whiteKingSquare = sq
			}
			if piece == BlackKing {
				b.blackKingSquare = sq
			}
			file++
		}
	}
	if fields[1] == "w" {
		b.sideToMove = White
	} else {
		b.sideToMove = Black
	}
	return b
}

// --- Methods to satisfy the Board interface ---

func (b *Bitboard) ApplyMove(m Move) {
	movingPiece, _ := b.pieceAt(m.From)
	moveMask := bitboard((1 << m.From) | (1 << m.To))
	isCapture := b.byColor[oppositeColor(b.sideToMove)].getBit(m.To)
	var capturedPiece Piece
	if isCapture {
		capturedPiece, _ = b.pieceAt(m.To)
	}
	b.byPiece[movingPiece] ^= moveMask
	b.byColor[b.sideToMove] ^= moveMask
	if isCapture {
		b.byPiece[capturedPiece].clearBit(m.To)
		b.byColor[oppositeColor(b.sideToMove)].clearBit(m.To)
	}
	if m.Promotion != Empty {
		b.byPiece[movingPiece].clearBit(m.To)
		b.byPiece[m.Promotion].setBit(m.To)
	}
	if movingPiece == WhiteKing {
		b.whiteKingSquare = m.To
	}
	if movingPiece == BlackKing {
		b.blackKingSquare = m.To
	}
	b.sideToMove = oppositeColor(b.sideToMove)
}

func (b *Bitboard) GenerateLegalMoves() []Move {
	var legalMoves []Move
	pseudoLegalMoves := b.generatePseudoLegalMoves()
	for _, move := range pseudoLegalMoves {
		tempBoard := *b
		tempBoard.ApplyMove(move)
		var kingSq int
		if b.sideToMove == White {
			kingSq = tempBoard.whiteKingSquare
		} else {
			kingSq = tempBoard.blackKingSquare
		}
		if !tempBoard.isSquareAttacked(kingSq, tempBoard.sideToMove) {
			legalMoves = append(legalMoves, move)
		}
	}
	return legalMoves
}

func (b *Bitboard) SideToMove() Color { return b.sideToMove }
func (b *Bitboard) IsCheckmate() bool { return len(b.GenerateLegalMoves()) == 0 && b.isKingInCheck() }
func (b *Bitboard) IsStalemate() bool { return len(b.GenerateLegalMoves()) == 0 && !b.isKingInCheck() }
func (b *Bitboard) isKingInCheck() bool {
	kingSq := b.blackKingSquare
	if b.sideToMove == White {
		kingSq = b.whiteKingSquare
	}
	return b.isSquareAttacked(kingSq, oppositeColor(b.sideToMove))
}

// --- Bitboard-Specific Logic ---

// THIS IS THE NEW, FULLY CORRECTED isSquareAttacked FUNCTION
func (b *Bitboard) isSquareAttacked(sq int, byColor Color) bool {
	// Pawn attacks
	if byColor == White {
		if sq/8 > 0 {
			if sq%8 > 0 && b.byPiece[WhitePawn].getBit(sq-9) {
				return true
			}
			if sq%8 < 7 && b.byPiece[WhitePawn].getBit(sq-7) {
				return true
			}
		}
	} else {
		if sq/8 < 7 {
			if sq%8 > 0 && b.byPiece[BlackPawn].getBit(sq+7) {
				return true
			}
			if sq%8 < 7 && b.byPiece[BlackPawn].getBit(sq+9) {
				return true
			}
		}
	}

	// Knight attacks
	var knightPiece Piece = WhiteKnight
	if byColor == Black {
		knightPiece = BlackKnight
	}
	for _, offset := range knightOffsets {
		to := sq + offset
		if to >= 0 && to < 64 && dist(sq%8, to%8) <= 2 {
			if b.byPiece[knightPiece].getBit(to) {
				return true
			}
		}
	}

	// King attacks
	var kingPiece Piece = WhiteKing
	if byColor == Black {
		kingPiece = BlackKing
	}
	for _, offset := range kingOffsets {
		to := sq + offset
		if to >= 0 && to < 64 && dist(sq%8, to%8) <= 1 {
			if b.byPiece[kingPiece].getBit(to) {
				return true
			}
		}
	}

	occupied := b.byColor[White] | b.byColor[Black]

	// Rook and Queen attacks (along ranks and files)
	var rook, queen Piece = WhiteRook, WhiteQueen
	if byColor == Black {
		rook, queen = BlackRook, BlackQueen
	}
	for _, dir := range rookDirections {
		for to := sq + dir; ; to += dir {
			// This is the robust wrap-around check
			if to < 0 || to >= 64 {
				break
			}
			if dir == 1 || dir == -1 { // Horizontal
				if to/8 != sq/8 {
					break
				}
			}

			if b.byPiece[rook].getBit(to) || b.byPiece[queen].getBit(to) {
				return true
			}
			if occupied.getBit(to) {
				break
			} // Blocked by another piece
		}
	}

	// Bishop and Queen attacks (along diagonals)
	var bishop Piece = WhiteBishop
	if byColor == Black {
		bishop = BlackBishop
	}
	for _, dir := range bishopDirections {
		for to := sq + dir; ; to += dir {
			// Robust wrap-around check for diagonals
			if to < 0 || to >= 64 || dist(to%8, (to-dir)%8) != 1 {
				break
			}

			if b.byPiece[bishop].getBit(to) || b.byPiece[queen].getBit(to) {
				return true
			}
			if occupied.getBit(to) {
				break
			} // Blocked
		}
	}

	return false
}

func (b *Bitboard) generatePseudoLegalMoves() []Move {
	var moves []Move
	occupied := b.byColor[White] | b.byColor[Black]
	myPieces := b.byColor[b.sideToMove]
	enemyPieces := b.byColor[oppositeColor(b.sideToMove)]
	b.generatePawnMoves(&moves, ^occupied, enemyPieces)
	b.generateKnightMoves(&moves, myPieces)
	b.generateSlidingMoves(&moves, occupied, myPieces)
	b.generateKingMoves(&moves, myPieces)
	return moves
}

func (b *Bitboard) generateSlidingMoves(moves *[]Move, occupied, myPieces bitboard) {
	var rooks, bishops, queens bitboard
	if b.sideToMove == White {
		rooks, bishops, queens = b.byPiece[WhiteRook], b.byPiece[WhiteBishop], b.byPiece[WhiteQueen]
	} else {
		rooks, bishops, queens = b.byPiece[BlackRook], b.byPiece[BlackBishop], b.byPiece[BlackQueen]
	}
	piecesToMove := rooks | queens
	for piecesToMove != 0 {
		from := piecesToMove.lsb()
		for _, dir := range rookDirections {
			for i := 1; ; i++ {
				to := from + dir*i
				if to < 0 || to >= 64 {
					break
				}
				if dir == 1 || dir == -1 {
					if to/8 != from/8 {
						break
					}
				}
				if myPieces.getBit(to) {
					break
				}
				*moves = append(*moves, Move{From: from, To: to})
				if occupied.getBit(to) {
					break
				}
			}
		}
		piecesToMove.clearBit(from)
	}
	piecesToMove = bishops | queens
	for piecesToMove != 0 {
		from := piecesToMove.lsb()
		for _, dir := range bishopDirections {
			for i := 1; ; i++ {
				to := from + dir*i
				if to < 0 || to >= 64 || dist(to%8, (to-dir*i)%8) != i {
					break
				}
				if myPieces.getBit(to) {
					break
				}
				*moves = append(*moves, Move{From: from, To: to})
				if occupied.getBit(to) {
					break
				}
			}
		}
		piecesToMove.clearBit(from)
	}
}

func (b *Bitboard) generateKingMoves(moves *[]Move, myPieces bitboard) {
	from := b.whiteKingSquare
	if b.sideToMove == Black {
		from = b.blackKingSquare
	}
	for _, offset := range kingOffsets {
		to := from + offset
		if to >= 0 && to < 64 && dist(from%8, to%8) <= 1 {
			if !myPieces.getBit(to) {
				*moves = append(*moves, Move{From: from, To: to})
			}
		}
	}
}

func (b *Bitboard) generatePawnMoves(moves *[]Move, empty, enemy bitboard) {
	// ... (unchanged)
	var pawns, singlePush, doublePush bitboard
	if b.sideToMove == White {
		pawns = b.byPiece[WhitePawn]
		singlePush = (pawns << 8) & empty
		doublePush = ((singlePush & Rank3) << 8) & empty
		for singlePush != 0 {
			to := singlePush.lsb()
			*moves = append(*moves, Move{From: to - 8, To: to})
			singlePush.clearBit(to)
		}
		for doublePush != 0 {
			to := doublePush.lsb()
			*moves = append(*moves, Move{From: to - 16, To: to})
			doublePush.clearBit(to)
		}
		capturesWest := (pawns << 7) & enemy & ^FileH
		capturesEast := (pawns << 9) & enemy & ^FileA
		for capturesWest != 0 {
			to := capturesWest.lsb()
			*moves = append(*moves, Move{From: to - 7, To: to})
			capturesWest.clearBit(to)
		}
		for capturesEast != 0 {
			to := capturesEast.lsb()
			*moves = append(*moves, Move{From: to - 9, To: to})
			capturesEast.clearBit(to)
		}
	} else {
		pawns = b.byPiece[BlackPawn]
		singlePush = (pawns >> 8) & empty
		doublePush = ((singlePush & Rank6) >> 8) & empty
		for singlePush != 0 {
			to := singlePush.lsb()
			*moves = append(*moves, Move{From: to + 8, To: to})
			singlePush.clearBit(to)
		}
		for doublePush != 0 {
			to := doublePush.lsb()
			*moves = append(*moves, Move{From: to + 16, To: to})
			doublePush.clearBit(to)
		}
		capturesWest := (pawns >> 9) & enemy & ^FileH
		capturesEast := (pawns >> 7) & enemy & ^FileA
		for capturesWest != 0 {
			to := capturesWest.lsb()
			*moves = append(*moves, Move{From: to + 9, To: to})
			capturesWest.clearBit(to)
		}
		for capturesEast != 0 {
			to := capturesEast.lsb()
			*moves = append(*moves, Move{From: to + 7, To: to})
			capturesEast.clearBit(to)
		}
	}
}

func (b *Bitboard) generateKnightMoves(moves *[]Move, myPieces bitboard) {
	// ... (unchanged)
	knights := b.byPiece[WhiteKnight]
	if b.sideToMove == Black {
		knights = b.byPiece[BlackKnight]
	}
	for knights != 0 {
		from := knights.lsb()
		var attacks bitboard
		if from < 48 {
			if from%8 > 0 {
				attacks.setBit(from + 15)
			}
			if from%8 < 7 {
				attacks.setBit(from + 17)
			}
		}
		if from > 15 {
			if from%8 > 0 {
				attacks.setBit(from - 17)
			}
			if from%8 < 7 {
				attacks.setBit(from - 15)
			}
		}
		if from%8 > 1 {
			if from < 56 {
				attacks.setBit(from + 6)
			}
			if from > 7 {
				attacks.setBit(from - 10)
			}
		}
		if from%8 < 6 {
			if from < 56 {
				attacks.setBit(from + 10)
			}
			if from > 7 {
				attacks.setBit(from - 6)
			}
		}
		validMoves := attacks & ^myPieces
		for validMoves != 0 {
			to := validMoves.lsb()
			*moves = append(*moves, Move{From: from, To: to})
			validMoves.clearBit(to)
		}
		knights.clearBit(from)
	}
}

func (b *Bitboard) pieceAt(sq int) (Piece, bool) {
	// ... (unchanged)
	for i := WhitePawn; i <= BlackKing; i++ {
		if b.byPiece[i].getBit(sq) {
			return i, true
		}
	}
	return Empty, false
}

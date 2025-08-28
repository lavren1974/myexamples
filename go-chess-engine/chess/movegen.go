package chess

var (
	rookDirections   = []int{-8, -1, 1, 8}
	bishopDirections = []int{-9, -7, 7, 9}
	knightOffsets    = []int{-17, -15, -10, -6, 6, 10, 15, 17}
	kingOffsets      = []int{-9, -8, -7, -1, 1, 7, 8, 9}
)

func (s *State) GenerateLegalMoves() []Move {
	var legalMoves []Move
	pseudoLegalMoves := s.generatePseudoLegalMoves()
	for _, move := range pseudoLegalMoves {
		tempState := *s
		tempState.ApplyMove(move)
		var kingSquare int
		if s.SideToMove == White {
			kingSquare = tempState.whiteKingSquare
		} else {
			kingSquare = tempState.blackKingSquare
		}
		if !tempState.isSquareAttacked(kingSquare, tempState.SideToMove) {
			legalMoves = append(legalMoves, move)
		}
	}
	return legalMoves
}

func (s *State) isSquareAttacked(sq int, byColor Color) bool {
	if byColor == White {
		if sq/8 > 0 {
			if sq%8 > 0 {
				if s.Board[sq-9] == WhitePawn {
					return true
				}
			}
			if sq%8 < 7 {
				if s.Board[sq-7] == WhitePawn {
					return true
				}
			}
		}
	} else {
		if sq/8 < 7 {
			if sq%8 > 0 {
				if s.Board[sq+7] == BlackPawn {
					return true
				}
			}
			if sq%8 < 7 {
				if s.Board[sq+9] == BlackPawn {
					return true
				}
			}
		}
	}
	for _, offset := range knightOffsets {
		targetSq := sq + offset
		if targetSq >= 0 && targetSq < 64 && dist(sq%8, targetSq%8) <= 2 {
			piece := s.Board[targetSq]
			if piece.Color() == byColor && (piece == WhiteKnight || piece == BlackKnight) {
				return true
			}
		}
	}
	for _, offset := range kingOffsets {
		targetSq := sq + offset
		if targetSq >= 0 && targetSq < 64 && dist(sq%8, targetSq%8) <= 1 {
			piece := s.Board[targetSq]
			if piece.Color() == byColor && (piece == WhiteKing || piece == BlackKing) {
				return true
			}
		}
	}
	for _, dir := range rookDirections {
		for targetSq := sq + dir; ; targetSq += dir {
			if targetSq < 0 || targetSq >= 64 || dist(targetSq%8, (targetSq-dir)%8) > 1 {
				break
			}
			piece := s.Board[targetSq]
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
			if targetSq < 0 || targetSq >= 64 || dist(targetSq%8, (targetSq-dir)%8) != 1 {
				break
			}
			piece := s.Board[targetSq]
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

func (s *State) generatePseudoLegalMoves() []Move {
	var moves []Move
	for from := 0; from < 64; from++ {
		piece := s.Board[from]
		if piece.Color() != s.SideToMove {
			continue
		}
		switch piece {
		case WhitePawn, BlackPawn:
			moves = append(moves, s.generatePawnMoves(from)...)
		case WhiteKnight, BlackKnight:
			moves = append(moves, s.generateKnightMoves(from)...)
		case WhiteBishop, BlackBishop:
			moves = append(moves, s.generateSlidingMoves(from, bishopDirections)...)
		case WhiteRook, BlackRook:
			moves = append(moves, s.generateSlidingMoves(from, rookDirections)...)
		case WhiteQueen, BlackQueen:
			moves = append(moves, s.generateSlidingMoves(from, bishopDirections)...)
			moves = append(moves, s.generateSlidingMoves(from, rookDirections)...)
		case WhiteKing, BlackKing:
			moves = append(moves, s.generateKingMoves(from)...)
		}
	}
	return moves
}

// THIS IS THE NEW, ROBUST SLIDING MOVE GENERATOR
func (s *State) generateSlidingMoves(from int, directions []int) []Move {
	var moves []Move
	isRook := s.Board[from] == WhiteRook || s.Board[from] == BlackRook

	for _, dir := range directions {
		prevSquare := from
		for {
			to := prevSquare + dir

			// 1. Off-board check
			if to < 0 || to >= 64 {
				break
			}

			// 2. Wrap-around check
			rankDist := dist(to/8, prevSquare/8)
			fileDist := dist(to%8, prevSquare%8)

			if isRook {
				if rankDist > 1 || fileDist > 1 || (rankDist == 1 && fileDist == 1) {
					break
				}
			} else { // Is Bishop or Queen
				if rankDist != 1 || fileDist != 1 {
					break
				}
			}

			// 3. Check target square content
			targetPiece := s.Board[to]
			if targetPiece == Empty {
				moves = append(moves, Move{From: from, To: to})
			} else {
				if targetPiece.Color() != s.SideToMove {
					moves = append(moves, Move{From: from, To: to})
				}
				break // Blocked by a piece
			}

			prevSquare = to // Advance to the next square in this direction
		}
	}
	return moves
}

// --- (The rest of the file is unchanged from our previous fix) ---

func (s *State) generatePawnMoves(from int) []Move {
	var moves []Move
	var pushDir, startRank, promotionRank int
	var captureDirs []int
	if s.SideToMove == White {
		pushDir, startRank, promotionRank = 8, 1, 7
		captureDirs = []int{7, 9}
	} else {
		pushDir, startRank, promotionRank = -8, 6, 0
		captureDirs = []int{-7, -9}
	}
	oneStep := from + pushDir
	if oneStep >= 0 && oneStep < 64 && s.Board[oneStep] == Empty {
		isPromotion := (oneStep / 8) == promotionRank
		addPawnMove(s, &moves, from, oneStep, isPromotion)
		if from/8 == startRank {
			twoSteps := from + 2*pushDir
			if s.Board[twoSteps] == Empty {
				addPawnMove(s, &moves, from, twoSteps, false)
			}
		}
	}
	for _, capDir := range captureDirs {
		to := from + capDir
		if to < 0 || to >= 64 || dist(from%8, to%8) != 1 {
			continue
		}
		targetPiece := s.Board[to]
		if targetPiece != Empty && targetPiece.Color() != s.SideToMove {
			isPromotion := (to / 8) == promotionRank
			addPawnMove(s, &moves, from, to, isPromotion)
		}
	}
	return moves
}

func addPawnMove(s *State, moves *[]Move, from, to int, isPromotion bool) {
	if isPromotion {
		if s.SideToMove == White {
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

func (s *State) generateKnightMoves(from int) []Move {
	var moves []Move
	for _, offset := range knightOffsets {
		to := from + offset
		if to < 0 || to >= 64 {
			continue
		}
		if dist(from%8, to%8) > 2 || dist(from/8, to/8) > 2 {
			continue
		}
		targetPiece := s.Board[to]
		if targetPiece == Empty || targetPiece.Color() != s.SideToMove {
			moves = append(moves, Move{From: from, To: to})
		}
	}
	return moves
}

func (s *State) generateKingMoves(from int) []Move {
	var moves []Move
	for _, offset := range kingOffsets {
		to := from + offset
		if to < 0 || to >= 64 {
			continue
		}
		if dist(from%8, to%8) > 1 || dist(from/8, to/8) > 1 {
			continue
		}
		targetPiece := s.Board[to]
		if targetPiece == Empty || targetPiece.Color() != s.SideToMove {
			moves = append(moves, Move{From: from, To: to})
		}
	}
	return moves
}

func dist(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

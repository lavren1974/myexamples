package chess

// Direction offsets for sliding pieces and kings
var (
	rookDirections   = []int{-8, -1, 1, 8}
	bishopDirections = []int{-9, -7, 7, 9}
	knightOffsets    = []int{-17, -15, -10, -6, 6, 10, 15, 17}
	kingOffsets      = []int{-9, -8, -7, -1, 1, 7, 8, 9}
)

// GenerateMoves creates a list of pseudo-legal moves for the current position.
func (s *State) GenerateMoves() []Move {
	var moves []Move
	for from := 0; from < 64; from++ {
		piece := s.Board[from]
		if piece.Color() != s.SideToMove {
			continue // Skip pieces that don't belong to the current player
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
			// A queen is just a rook and a bishop combined
			moves = append(moves, s.generateSlidingMoves(from, bishopDirections)...)
			moves = append(moves, s.generateSlidingMoves(from, rookDirections)...)
		case WhiteKing, BlackKing:
			moves = append(moves, s.generateKingMoves(from)...)
		}
	}
	return moves
}

// generateSlidingMoves handles rooks, bishops, and queens.
func (s *State) generateSlidingMoves(from int, directions []int) []Move {
	var moves []Move
	for _, dir := range directions {
		for to := from + dir; ; to += dir {
			// Check for off-board
			if to < 0 || to >= 64 || distToEdge(to, -dir) == 0 {
				break
			}

			targetPiece := s.Board[to]

			if targetPiece == Empty {
				// Quiet move to an empty square
				moves = append(moves, Move{From: from, To: to})
			} else {
				// Capture move
				if targetPiece.Color() != s.SideToMove {
					moves = append(moves, Move{From: from, To: to})
				}
				// Stop searching in this direction (blocked by own or enemy piece)
				break
			}
		}
	}
	return moves
}

// generateKnightMoves handles knight movement.
func (s *State) generateKnightMoves(from int) []Move {
	var moves []Move
	for _, offset := range knightOffsets {
		to := from + offset
		if to < 0 || to >= 64 {
			continue
		}
		// Check for board wrap-around
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

// generateKingMoves handles king movement.
func (s *State) generateKingMoves(from int) []Move {
	var moves []Move
	for _, offset := range kingOffsets {
		to := from + offset
		if to < 0 || to >= 64 {
			continue
		}
		// Check for board wrap-around
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

// generatePawnMoves handles all pawn moves: pushes and captures.
func (s *State) generatePawnMoves(from int) []Move {
	var moves []Move

	// Determine direction and starting rank
	var pushDir, startRank, secondRank int
	var captureDirs []int
	if s.SideToMove == White {
		pushDir = 8
		startRank, secondRank = 1, 2 // Ranks 2 and 3
		captureDirs = []int{7, 9}
	} else {
		pushDir = -8
		startRank, secondRank = 6, 5 // Ranks 7 and 6
		captureDirs = []int{-7, -9}
	}

	// 1. Single Push
	oneStep := from + pushDir
	if oneStep >= 0 && oneStep < 64 && s.Board[oneStep] == Empty {
		moves = append(moves, Move{From: from, To: oneStep})

		// 2. Double Push from starting rank
		if from/8 == startRank {
			twoSteps := from + 2*pushDir
			if twoSteps/8 == secondRank+(startRank-1) && s.Board[twoSteps] == Empty {
				moves = append(moves, Move{From: from, To: twoSteps})
			}
		}
	}

	// 3. Captures
	for _, capDir := range captureDirs {
		to := from + capDir
		if to < 0 || to >= 64 {
			continue
		}
		// Prevent wrap-around captures
		if dist(from%8, to%8) != 1 {
			continue
		}

		targetPiece := s.Board[to]
		if targetPiece != Empty && targetPiece.Color() != s.SideToMove {
			moves = append(moves, Move{From: from, To: to})
		}
	}

	// Note: En-passant and promotions are not yet implemented.
	return moves
}

// --- Utility functions for move generation ---

// dist calculates the distance between two numbers (for rank/file checks).
func dist(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

// distToEdge calculates moves until a piece hits the edge of the board.
// Used to prevent wrap-around for sliding pieces.
func distToEdge(sq int, dir int) int {
	file, rank := sq%8, sq/8
	switch dir {
	case 1:
		return 7 - file // East
	case -1:
		return file // West
	case 8:
		return 7 - rank // North
	case -8:
		return rank // South
	case 7:
		return min(file, rank) // South-West
	case -7:
		return 7 - max(file, rank) // North-East
	case 9:
		return min(7-file, 7-rank) // North-East
	case -9:
		return min(file, 7-rank) // South-West
	}
	return 0 // Should not happen
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

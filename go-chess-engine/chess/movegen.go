package chess

// GenerateMoves creates a list of pseudo-legal moves.
// NOTE: This is still the highly simplified version.
func (s *State) GenerateMoves() []Move {
	var moves []Move
	for from := 0; from < 64; from++ {
		piece := s.Board[from]
		if piece == Empty {
			continue
		}

		isWhitePiece := piece >= WhitePawn && piece <= WhiteKing
		if (s.SideToMove == White) != isWhitePiece {
			continue
		}

		switch piece {
		case WhiteKnight, BlackKnight:
			knightOffsets := []int{-17, -15, -10, -6, 6, 10, 15, 17}
			for _, offset := range knightOffsets {
				to := from + offset
				if to >= 0 && to < 64 {
					fromFile, toFile := from%8, to%8
					if abs(fromFile-toFile) <= 2 {
						moves = append(moves, Move{From: from, To: to})
					}
				}
			}
		case WhitePawn:
			to := from + 8
			if to < 64 && s.Board[to] == Empty {
				moves = append(moves, Move{From: from, To: to})
			}
		case BlackPawn:
			to := from - 8
			if to >= 0 && s.Board[to] == Empty {
				moves = append(moves, Move{From: from, To: to})
			}
		}
	}
	return moves
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

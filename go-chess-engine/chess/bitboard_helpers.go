package chess

import "math/bits"

// A bitboard is a 64-bit integer representing squares on the board.
type bitboard uint64

// --- Bitboard Constants (Masks) ---
const (
	FileA bitboard = 0x0101010101010101
	FileB bitboard = FileA << 1
	FileC bitboard = FileA << 2
	FileD bitboard = FileA << 3
	FileE bitboard = FileA << 4
	FileF bitboard = FileA << 5
	FileG bitboard = FileA << 6
	FileH bitboard = FileA << 7

	Rank1 bitboard = 0x00000000000000FF
	Rank2 bitboard = Rank1 << 8
	Rank3 bitboard = Rank1 << 16
	Rank4 bitboard = Rank1 << 24
	Rank5 bitboard = Rank1 << 32
	Rank6 bitboard = Rank1 << 40
	Rank7 bitboard = Rank1 << 48
	Rank8 bitboard = Rank1 << 56
)

// --- Bitwise Helper Functions ---

// setBit sets a bit at a given square index.
func (b *bitboard) setBit(sq int) {
	*b |= (1 << sq)
}

// clearBit clears a bit at a given square index.
func (b *bitboard) clearBit(sq int) {
	*b &= ^(1 << sq)
}

// getBit checks if a bit is set at a given square index.
func (b bitboard) getBit(sq int) bool {
	return (b>>sq)&1 == 1
}

// popcount (population count) counts the number of set bits.
func (b bitboard) popcount() int {
	return bits.OnesCount64(uint64(b))
}

// lsb (least significant bit) finds the index of the first set bit.
// Also known as bitScanForward. Returns 0 if the bitboard is empty.
func (b bitboard) lsb() int {
	return bits.TrailingZeros64(uint64(b))
}

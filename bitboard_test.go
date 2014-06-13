package bitboard

import (
	"testing"
)

func TestTest(t *testing.T) {
	bitboard := New()
	bitboard.PrettyPrint()
}

func TestMovePiece(t *testing.T) {
	bitboard := New()
	bitboard.MovePieceFileRank("a", 2, "a", 3)
	if bitboard.GetPieceRowFile("a", 3) == -1 && bitboard.GetPieceRowFile("a", 2) != -1 {
		t.Error("There shoule be something here. Was %v but expected -1", bitboard.GetPieceRowFile("a", 3))
	}
}

func BenchmarkSetPiece(b *testing.B) {
	bitboard := New()
	for i := 0; i < b.N; i++ {
		bitboard.SetPiece(1, 2, 3)
	}
}

func BenchmarkMovePiece(b *testing.B) {
	bitboard := New()
	for i := 0; i < b.N; i++ {
		bitboard.MovePiece(0, 0, 0, 1)
		bitboard.MovePiece(0, 1, 0, 0)
	}
}

func BenchmarkMovePieceFileRank(b *testing.B) {
	bitboard := New()
	for i := 0; i < b.N; i++ {
		bitboard.MovePieceFileRank("a", 1, "a", 2)
		bitboard.MovePieceFileRank("a", 2, "a", 1)
	}
}

func BenchmarkNewBoard(b *testing.B) {
	bitboards := make([]BitBoard, b.N)
	for i := 0; i < b.N-1; i++ {
		bitboards[i] = New()
	}
}

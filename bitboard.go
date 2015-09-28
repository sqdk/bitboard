package bitboard

import (
	"fmt"
	"github.com/sqdk/bitops"
)

const (
	OUT_OF_BOUNDS = -2
	EMPTY_CELL    = -1
	WPawnId       = 0
	WRooksId      = 1
	WKnightsId    = 2
	WBishopsId    = 3
	WQueenId      = 4
	WKingId       = 5
	BPawnId       = 6
	BRooksId      = 7
	BKnightsId    = 8
	BBishopsId    = 9
	BQueenId      = 10
	BKingId       = 11
)

/*
	Each uint64 represents the board for a single piece type.
	Every bit indicates weather the given piece is at the specific location.

	Pieces can be moved by 0-based XY coordinate or by rank and file.

	The LookupTable is used to indicate if any piece is at the given location on the board
	This will shorten the time needed for a lot of operations and only increase memory overhead
 	by 1/13 of the total memory usage.

 	uint64 to board mapping:

	  1  2  3  4  5  6  7  8   <- Rank (coordinate y)
	a 0  1  2  3  4  5  6  7
	b 8  9  10 11 12 13 14 15
	c 16 17 18 19 20 21 22 23
	d 24 25 26 27 28 29 30 31
	e 32 33 34 35 36 37 38 39
	f 40 41 42 43 44 45 46 47
	g 48 49 50 51 52 53 54 55
	h 56 57 58 59 60 61 62 63
	^
	File (coordinate x)
*/
type BitBoard struct {
	Board       [12]uint64
	LookupTable uint64
}

/*
	Creates a new chess board by initializing the Board array
	with the appropriate values.
*/
func New() BitBoard {
	var board BitBoard
	board.Board[WPawnId] = uint64(4629771061636907072)
	board.Board[BPawnId] = uint64(144680345676153346)
	board.Board[WRooksId] = uint64(9223372036854775936)
	board.Board[BRooksId] = uint64(72057594037927937)
	board.Board[BKnightsId] = uint64(281474976710912)
	board.Board[WKnightsId] = uint64(36028797018996736)
	board.Board[WBishopsId] = uint64(140737496743936)
	board.Board[BBishopsId] = uint64(1099511693312)
	board.Board[WQueenId] = uint64(2147483648)
	board.Board[WKingId] = uint64(549755813888)
	board.Board[BKingId] = uint64(4294967296)
	board.Board[BQueenId] = uint64(16777216)

	board.LookupTable = uint64(0xFFFF00000000FFFF)
	return board
}

//Used to minimize memory overhead during test
func (board *BitBoard) ResetBoard() {
	board.Board[WPawnId] = uint64(4629771061636907072)
	board.Board[BPawnId] = uint64(144680345676153346)
	board.Board[WRooksId] = uint64(9223372036854775936)
	board.Board[BRooksId] = uint64(72057594037927937)
	board.Board[BKnightsId] = uint64(281474976710912)
	board.Board[WKnightsId] = uint64(36028797018996736)
	board.Board[WBishopsId] = uint64(140737496743936)
	board.Board[BBishopsId] = uint64(1099511693312)
	board.Board[WQueenId] = uint64(2147483648)
	board.Board[WKingId] = uint64(549755813888)
	board.Board[BKingId] = uint64(4294967296)
	board.Board[BQueenId] = uint64(16777216)

	board.LookupTable = uint64(0xFFFF00000000FFFF)
}

func Clone(b BitBoard) BitBoard {
	var board BitBoard
	for i := 0; i < len(b.Board); i++ {
		board.Board[i] = b.Board[i]
	}
	board.LookupTable = b.LookupTable
	return board
}

/*
	Moves piece by rank and file instead of XY. This is a lot slower because of the
	expensive operations needed to decipher the string input to XY coordinates.
	Intended use is for easy input of moves with standard chess notation.
*/
func (b *BitBoard) MovePieceFileRank(fileStart string, rankStart int, fileEnd string, rankEnd int) {
	xstart := RankToX(rankStart)
	xend := RankToX(rankEnd)
	ystart := FileToY(fileStart)
	yend := FileToY(fileEnd)
	b.MovePiece(xstart, ystart, xend, yend)
}

/*
	Returns the piece in the given file and rank as an integer
	indicating the type and color of the piece.
*/
func (b *BitBoard) GetPieceRowFile(file string, rank int) int {
	return b.GetPiece(RankToX(rank), FileToY(file))
}

func (b *BitBoard) MovePiece(xstart, ystart, xend, yend int) {
	piece := b.GetPiece(xstart, ystart)
	destinationPiece := b.GetPiece(xend, yend)
	if piece == EMPTY_CELL {
		return
	}
	if destinationPiece != EMPTY_CELL {
		b.RemovePieceFast(xend, yend, destinationPiece)
	}
	b.SetPiece(xend, yend, piece)
	b.RemovePieceFast(xstart, ystart, piece)
}

func (b *BitBoard) MovePieceFast(xstart, ystart, xend, yend, piece int) {
	b.SetPiece(xend, yend, piece)
	b.RemovePieceFast(xstart, ystart, piece)
}

func (b *BitBoard) GetPiece(x, y int) int {
	if x < 0 || y < 0 || x > 7 || y > 7 {
		return OUT_OF_BOUNDS
	}

	//Query lookup table to check if there is a piece in the specific position
	if !bitops.QueryBit(&b.LookupTable, xyToIndex(y, x)) {
		return EMPTY_CELL
	}

	for pieceIndex := 0; pieceIndex < 12; pieceIndex++ {
		if bitops.QueryBit(&b.Board[pieceIndex], xyToIndex(x, y)) {
			return pieceIndex
		}
	}
	return EMPTY_CELL
}

func (b *BitBoard) SetPiece(x, y, piece int) {
	bitops.SetBit(&b.Board[piece], xyToIndex(x, y), true)
	bitops.SetBit(&b.LookupTable, xyToIndex(y, x), true)
}

func (b *BitBoard) RemovePiece(x, y int) {
	piece := b.GetPiece(x, y)
	if piece == EMPTY_CELL {
		return
	}

	bitops.SetBit(&b.Board[piece], xyToIndex(x, y), false)
	bitops.SetBit(&b.LookupTable, xyToIndex(y, x), false)
}

func (b *BitBoard) RemovePieceFast(x, y, piece int) {
	bitops.SetBit(&b.Board[piece], xyToIndex(x, y), false)
	bitops.SetBit(&b.LookupTable, xyToIndex(y, x), false)
}

func xyToIndex(x, y int) int {
	return x + (y * 8)
}

func FileToY(file string) int {
	return int(file[0]) - 97
	/*
		switch file {
		case "a":
			return 0
		case "b":
			return 1
		case "c":
			return 2
		case "d":
			return 3
		case "e":
			return 4
		case "f":
			return 5
		case "g":
			return 6
		case "h":
			return 7
		}
		return -1*/
}

func YToFile(y int) string {
	return string(y + 97)
	/*
		if y <= 7 {
			return YToFileLookup[y]
		}
		return "?"*/
}

func RankToX(rank int) int {
	return 8 - rank
}

func XToRank(x int) int {
	return x + 1
}

func (b *BitBoard) PrettyPrint() {
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			switch b.GetPiece(x, y) {
			case WPawnId:
				fmt.Print("P")
			case BPawnId:
				fmt.Print("p")
			case WRooksId:
				fmt.Print("R")
			case BRooksId:
				fmt.Print("r")
			case WKnightsId:
				fmt.Print("N")
			case BKnightsId:
				fmt.Print("n")
			case WBishopsId:
				fmt.Print("B")
			case BBishopsId:
				fmt.Print("b")
			case WKingId:
				fmt.Print("K")
			case BKingId:
				fmt.Print("k")
			case WQueenId:
				fmt.Print("Q")
			case BQueenId:
				fmt.Print("q")
			default:
				fmt.Print(".")
			}

		}
		fmt.Println()
	}
}

func (b *BitBoard) PrettyPrintMark(coords []XYPair) {
	fmt.Println(b.PrettyPrintMarkToString(coords))
}

func (b *BitBoard) PrettyPrintMarkToString(coords []XYPair) string {
	s := make([]byte, 0)
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			piece := b.GetPiece(x, y)

			mark := false
			for i := 0; i < len(coords); i++ {
				if x == coords[i].X && y == coords[i].Y {
					if piece != -1 {
						s = append(s, '+')
					} else {
						s = append(s, '*')
					}
					mark = true
					break
				}
			}
			if mark {
				continue
			}

			switch piece {
			case WPawnId:
				s = append(s, 'P')
			case BPawnId:
				s = append(s, 'p')
			case WRooksId:
				s = append(s, 'R')
			case BRooksId:
				s = append(s, 'r')
			case WKnightsId:
				s = append(s, 'N')
			case BKnightsId:
				s = append(s, 'n')
			case WBishopsId:
				s = append(s, 'B')
			case BBishopsId:
				s = append(s, 'b')
			case WKingId:
				s = append(s, 'K')
			case BKingId:
				s = append(s, 'k')
			case WQueenId:
				s = append(s, 'Q')
			case BQueenId:
				s = append(s, 'q')
			default:
				s = append(s, '.')
			}
		}
		if x != 7 {
			s = append(s, '\n')
		}
	}
	return string(s)
}

type XYPair struct {
	X int
	Y int
}

func IsPieceWhite(piece int) bool {
	if piece >= WPawnId && piece < BPawnId {
		return true
	}
	return false
}

package model

import (
	"fmt"
)

var (
	whitePiecesRunes                = []ChessPieceType{'\u2654', '\u2655', '\u2656', '\u2657', '\u2658'}
	whitePawnRune    ChessPieceType = '\u2659'
	blackPiecesRunes                = []ChessPieceType{'\u265A', '\u265B', '\u265C', '\u265D', '\u265E'}
	blackPawnRune    ChessPieceType = '\u265F'
)

func NewGame(config *GameConfig) *Game {
	field, whitePieces, whitePawns, blackPieces, blackPawns := InitCells(config.Size)
	first := initPlayer(config.FirstPlayerName, PlayerColorWhite, whitePieces, whitePawns)
	second := initPlayer(config.FirstPlayerName, PlayerColorBlack, blackPieces, blackPawns)

	game := &Game{
		first:  first,
		second: second,
		board:  &Board{cells: field},
	}
	return game
}

func initPlayer(pname string, color Color, pieces []Cell, pawns []Cell) *Player {
	chessPieces := make([]*ChessPiece, 0, len(pieces)*2)
	for _, c := range pawns {
		chessPieces = append(chessPieces, c.piece)
	}
	for _, c := range pieces {
		chessPieces = append(chessPieces, c.piece)
	}
	return &Player{
		name:   pname,
		color:  color,
		pieces: chessPieces,
	}
}

func InitCells(size int) (f []Cell, firstPlayerPieces []Cell, firstPlayerPawns []Cell, secondPlayerPieces []Cell, secondPlayerPawns []Cell) {
	field := make([]Cell, 0, size*size)
	for i := range size {
		names := make([]string, 0, size)
		for h := range size {
			names = append(names, getColumnChar(h+1))
		}

		for j := range size {
			c := Cell{}
			var val rune
			if (i+j)%2 == 0 {
				val = '#'
			} else {
				val = ' '
			}
			c.line = j
			c.value = &val
			c.column = names[j]
			field = append(field, c)
		}
	}

	whitePieces := fillCellsWithPieces(field[0:size], whitePiecesRunes)
	whitePawns := fillWithPawn(field[size:size*2], whitePawnRune)
	blackPawns := fillWithPawn(field[len(field)-size*2:len(field)-size], blackPawnRune)
	blackPieces := fillCellsWithPieces(field[len(field)-size:], blackPiecesRunes)
	return field, whitePieces, whitePawns, blackPieces, blackPawns
}

func fillWithPawn(cell []Cell, pawn ChessPieceType) []Cell {
	for i, c := range cell {
		if isEmptyCell(*c.value) {
			pawn := &ChessPiece{
				value:    pawn,
				currentX: c.line,
				currentY: c.column,
			}
			cell[i].piece = pawn
			cell[i].hasPiece = true
			cell[i].value = (*rune)(&pawn.value)
		}
	}
	return cell
}

func fillCellsWithPieces(cell []Cell, pieces []ChessPieceType) []Cell {
	for i, c := range cell {
		if isEmptyCell(*c.value) {
			chess := &ChessPiece{
				value:    pieces[i%len(pieces)],
				currentX: c.line,
				currentY: c.column,
			}
			cell[i].piece = chess
			cell[i].hasPiece = true
			cell[i].value = (*rune)(&chess.value)
		}
	}
	return cell
}

func isEmptyCell(cell rune) bool {
	return cell == ' ' || cell == '#'
}

// DebugPrintCells prints cells as a square board with size cells in each row.
func DebugPrintCells(cells []Cell, size int) {
	if size <= 0 {
		return
	}

	for i, cell := range cells {
		if i%size == 0 {
			fmt.Printf("%3d", i/size+1)
		}

		if cell.hasPiece {
			fmt.Printf("%3c", cell.piece.value)
		} else {
			fmt.Printf("%3c", *cell.value)
		}
		if (i+1)%size == 0 {
			fmt.Println()
		}
	}

	fmt.Printf("%3s", "")
	for i := 0; i < size && i < len(cells); i++ {
		fmt.Printf("%3s", cells[i].column)
	}
	fmt.Println()
}

func getColumnChar(i int) string {
	result := ""
	for i > 0 {
		i--
		result = string(rune('A'+i%26)) + result
		i /= 26
	}

	return result
}

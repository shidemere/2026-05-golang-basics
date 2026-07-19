// Package board let us work with game board.
package board

import (
	"fmt"
)

var (
	whitePieces = []rune{'\u2654', '\u2655', '\u2656', '\u2657', '\u2658'}
	whitePawn   = '\u2659'
	blackPieces = []rune{'\u265A', '\u265B', '\u265C', '\u265D', '\u265E'}
	blackPawn   = '\u265F'
)

func CreateBoard(size int) ([][]rune, error) {
	matrix := make([][]rune, size)
	for i := range size {
		matrix[i] = make([]rune, size)
		for j := range size {
			if (i+j)%2 == 0 {
				matrix[i][j] = '#'
			} else {
				matrix[i][j] = ' '
			}
		}
	}

	fillBoardWithChess(matrix)
	return matrix, nil
}

func fillBoardWithChess(matrix [][]rune) {
	fillWithPieces(matrix[0], whitePieces)
	fillWithPawn(matrix[1], whitePawn)
	fillWithPawn(matrix[len(matrix)-2], blackPawn)
	fillWithPieces(matrix[len(matrix)-1], blackPieces)
}

func fillWithPieces(line, pieces []rune) {
	for i, cell := range line {
		if isEmptyCell(cell) {
			line[i] = pieces[i%len(pieces)]
		}
	}
}

func fillWithPawn(line []rune, pawn rune) {
	for i, cell := range line {
		if isEmptyCell(cell) {
			line[i] = pawn
		}
	}
}

func isEmptyCell(cell rune) bool {
	return cell == ' ' || cell == '#'
}

func PrintBoard(matrix [][]rune, size int) {
	for lineidx, row := range matrix {
		// print numbers
		if size > lineidx {
			fmt.Printf("%3d", lineidx+1)
		}
		for _, r := range row {
			// %c prints the actual character representation of the rune
			fmt.Printf("%3c", r)
		}
		fmt.Println()
	}

	rowNumberWidth := max(3, len(fmt.Sprint(size)))
	fmt.Printf("%*s", rowNumberWidth, "")

	for i := range size {
		ch := getColumnChar(i + 1)
		fmt.Printf("%3s", ch)
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

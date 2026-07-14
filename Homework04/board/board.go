// Package board let us work with game board.
package board

import (
	"fmt"
	"slices"
)

var (
	whitePieces = []rune{'\u2654', '\u2655', '\u2656', '\u2657', '\u2658'}
	whitePawn   = '\u2659'
	blackPieces = []rune{'\u265A', '\u265B', '\u265C', '\u265D', '\u265E'}
	blackPawn   = '\u265F'
)

func CreateBoard(size int) ([][]rune, error) {
	matrix := make([][]rune, 0, size)
	var row []rune
	for i := range size {
		// fill left column with numbers
		row = make([]rune, 0, size)
		for j := range size {
			if (i+j)%2 == 0 {
				row = append(row, '#')
			} else {
				row = append(row, ' ')
			}
		}
		matrix = append(matrix, row)
	}
	fillBoardWithChess(matrix)
	return matrix, nil
}

func fillBoardWithChess(matrix [][]rune) {
	fillSecondLineWithPieces(matrix[0])
	fillWithPawn(matrix[1], whitePawn)
	fillWithPawn(matrix[len(matrix)-2], blackPawn)
	fillBottomLineWithBlackPieces(matrix[len(matrix)-1])
}

func fillBottomLineWithBlackPieces(line []rune) {
	counter := 0
	for i, v := range line {
		if counter == len(whitePieces) {
			counter = 0
		}
		if v == ' ' || v == '#' {
			line[i] = blackPieces[counter]
		}
		counter++
	}
}

func fillWithPawn(line []rune, r rune) {
	for i, v := range line {
		if v == ' ' || v == '#' {
			line[i] = r
		}
	}
}

func fillSecondLineWithPieces(firstLine []rune) {
	counter := 0
	for i, v := range firstLine {
		if counter == len(whitePieces) {
			counter = 0
		}
		if v == ' ' || v == '#' {
			firstLine[i] = whitePieces[counter]
		}
		counter++
	}
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

	switch {
	case size < 10:
		fmt.Printf("  ")
	case size > 10 && size < 100:
		fmt.Print("   ")
	case size > 100 && size < 1000:
		fmt.Print("    ")
	}

	for i := range size {
		ch := getColumnChar(i + 1)
		fmt.Printf("%3s", ch)
	}
	fmt.Println()
}

func getColumnChar(i int) string {
	if i <= 26 {
		return string(rune(64 + i))
	}

	res := make([]rune, 0)
	for i > 0 {
		i--
		remainder := i % 26 // find a letter
		letter := rune('A' + remainder)

		res = append(res, letter)
		i = i / 26
	}

	slices.Reverse(res)
	return string(res)
}

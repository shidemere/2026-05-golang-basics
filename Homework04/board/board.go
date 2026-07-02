// Package board let us work with game board.
package board

import (
	"fmt"
)

var (
	chars       = []rune{'H', 'G', 'F', 'E', 'D', 'C', 'B', 'A'}
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
		row = append(row, rune(i)+'1') // because rune(i) will show char with code point of i

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
	lastLine := fillLastLine(size)
	matrix = append(matrix, lastLine)
	return matrix, nil
}

func fillBoardWithChess(matrix [][]rune) {
	fillSecondLineWithPieces(matrix[0])
	fillWithPawn(matrix[1], whitePawn)
	fillWithPawn(matrix[len(matrix)-2], blackPawn)
	fillBottomLineWithBlackPieces(matrix[len(matrix)-1])
}

func fillLastLine(size int) []rune {
	lastLineWithLetters := make([]rune, 0, size+1) // + 1 for whitespace in beginning
	lastLineWithLetters = append(lastLineWithLetters, ' ')
	lastLineWithLetters = fillLastLineWIthLettersBySizeOfBoard(lastLineWithLetters, size)
	return lastLineWithLetters
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

func fillLastLineWIthLettersBySizeOfBoard(lastLineWithLetters []rune, size int) []rune {
	counter := 0
	for range size {
		if counter == len(chars) {
			counter = 0
		}
		lastLineWithLetters = append(lastLineWithLetters, chars[counter])
		counter++
	}
	return lastLineWithLetters
}

func PrintBoard(matrix [][]rune) {
	for _, row := range matrix {
		for _, r := range row {
			// %c prints the actual character representation of the rune
			fmt.Printf("%c ", r)
		}
		fmt.Println()
	}
}

// Package board let us work with game board.
package board

import "fmt"

func CreateBoard(size int) ([][]rune, error) {
	matrix := make([][]rune, 0, size)
	var row []rune
	for i := range size {
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
	return matrix, nil
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

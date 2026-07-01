// Package board let us work with game board.
package board

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

package main

import (
	"fmt"

	"github.com/shidemere/Homework03/board"
)

func main() {
	matrix, _ := board.CreateBoard(8)
	for _, row := range matrix {
		for _, r := range row {
			// %c prints the actual character representation of the rune
			fmt.Printf("%c ", r)
		}
		fmt.Println()
	}
}

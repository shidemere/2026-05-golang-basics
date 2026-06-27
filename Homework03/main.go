package main

import (
	"fmt"
	"os"

	"github.com/shidemere/Homework03/board"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("you need to run program with size of board.\nExample: go run main.go 8")
		os.Exit(1)
	}
	matrix, _ := board.CreateBoard(8)

	for _, row := range matrix {
		for _, r := range row {
			// %c prints the actual character representation of the rune
			fmt.Printf("%c ", r)
		}
		fmt.Println()
	}
}

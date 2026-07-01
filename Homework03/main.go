package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/shidemere/Homework03/board"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("you need to run program with size of board.\nExample: go run main.go 8")
		os.Exit(1)
	}
	arg := os.Args[1]
	val, err := strconv.Atoi(arg)
	if err != nil {
		log.Fatalf("incorrect input argument: %v", err)
	}
	matrix, _ := board.CreateBoard(val)

	for _, row := range matrix {
		for _, r := range row {
			// %c prints the actual character representation of the rune
			fmt.Printf("%c ", r)
		}
		fmt.Println()
	}
}

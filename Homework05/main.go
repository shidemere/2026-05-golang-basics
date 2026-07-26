package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/shidemere/2026-05-golang-basics/Homework05/model"
)

func main() {
	val := handleBoardSizeInput()
	first, second, err := handlePlayersNames()
	if err != nil {
		log.Fatalf("can't start game: %v", err)
	}
	cfg := &model.GameConfig{
		FirstPlayerName:  first,
		SecondPlayerName: second,
		Size:             val,
	}
	game := model.NewGame(cfg)
	model.DebugPrintCells(game.GetBoard().GetCells(), val)
}

func handleBoardSizeInput() int {
	if len(os.Args) < 2 {
		fmt.Println("you need to run program with size of board.\nExample: go run main.go 8")
		os.Exit(1)
	}
	arg := os.Args[1]
	val, err := strconv.Atoi(arg)
	if err != nil {
		log.Fatalf("incorrect input argument: %v", err)
	}

	return val
}

func handlePlayersNames() (first, second string, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	var counter int
	for {
		fmt.Printf("Пожалуйста, введите имя %d го игрока: ", counter+1)
		if scanner.Scan() {
			if scanner.Err() != nil {
				return "", "", fmt.Errorf("failed to handle player name: %w", err)
			}
			input := scanner.Text()
			err := validatePlayerName(input)
			if err != nil {
				return "", "", fmt.Errorf("failed to handle player name: %v", err)
			}

			if first == "" {
				first = input
				counter++
				continue
			}
			if first != "" && second == "" {
				second = input
				counter++
			}

			if first != "" && second != "" {
				break
			}
		}
	}
	return first, second, nil
}

func validatePlayerName(input string) any {
	if len(input) <= 3 || len(input) > 64 {
		return errors.New("name too long or too short, should be betwenn 3 and 64 character")
	}

	return nil
}

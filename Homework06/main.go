package main

import (
	"log"

	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/console"
	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/model"
	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/service"
)

func main() {
	size := console.HandleBoardSizeInput()
	first, second, err := console.HandlePlayersNames()
	if err != nil {
		log.Fatalf("can't start game: %v", err)
	}
	cfg := &model.GameConfig{
		FirstPlayerName:  first,
		SecondPlayerName: second,
		Size:             size,
	}
	game := model.NewGame(cfg)
	console.PrintBoard(game.GetBoard().GetCells(), size)

	err = service.HandleGameProcess(game, size)
	if err != nil {
		log.Fatalf("can't handle game process: %v", err)
	}
}

package main

import (
	"log"
	"time"

	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/console"
	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/model"
	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/repository"
	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/service"
)

const randomEntitySaveInterval = 5 * time.Second

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

	repo := repository.NewRepository()
	randomService := service.NewRandomService(repo)
	ticker := time.NewTicker(randomEntitySaveInterval)
	defer ticker.Stop()

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := randomService.Save(); err != nil {
					log.Printf("can't save random entity: %v", err)
				}
			case <-done:
				return
			}
		}
	}()

	err = service.HandleGameProcess(game, size)
	if err != nil {
		log.Fatalf("can't handle game process: %v", err)
	}
}

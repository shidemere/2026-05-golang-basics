package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/shidemere/2026-05-golang-basics/Homework07/internal/console"
	"github.com/shidemere/2026-05-golang-basics/Homework07/internal/model"
	"github.com/shidemere/2026-05-golang-basics/Homework07/internal/repository"
	"github.com/shidemere/2026-05-golang-basics/Homework07/internal/service"
)

const randomEntitySaveInterval = 5 * time.Second

func main() {
	var wg sync.WaitGroup
	size, countOfBoards := console.HandleBoardSizeInputAndCount()
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

	triples := make(chan model.Triple, countOfBoards)
	for i := range countOfBoards {
		wg.Add(1)
		go func() {
			wg.Done()
			err = service.HandleGameProcess(game, size, triples, i)
			if err != nil {
				log.Fatalf("can't handle game process: %v", err)
			}
		}()
	}

	// You need to separate console drawing in separate gorutine
	// For this you need to find a way send your actual board state in separate gorutine AND in separate gorutine you also need to listen other gorutines stop signals
	// You can use channels. Create new struct with name like "Pair" where will be slice of cells and size. And then give this channel to game-gorutine for writing
	// 			and give this channel to draw-gorutine for reading.
	// When channel will be closed -> you will understand the signals

	go func() {
		wg.Add(1)
		defer wg.Done()
		for t := range triples {
			time.Sleep(1 * time.Second)
			fmt.Printf("%s BOARD #%d %s\n", strings.Repeat("_", 25), t.NumberOfBoard, strings.Repeat("_", 25))

			console.CleanBoard()
			console.PrintBoard(t.Cells, t.Size)
			fmt.Printf("%s\n", strings.Repeat("_", 50))
		}
	}()
	wg.Wait()
}

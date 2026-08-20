package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
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
	size, countOfBoards := console.HandleBoardSizeInputAndCount()
	scanner := bufio.NewScanner(os.Stdin)

	first, second, err := console.HandlePlayersNames(scanner)
	if err != nil {
		log.Fatalf("can't start game: %v", err)
	}

	cfg := &model.GameConfig{
		FirstPlayerName:  first,
		SecondPlayerName: second,
		Size:             size,
	}

	if countOfBoards == 1 {
		console.PrintInstructions()
	} else {
		console.PrintConcurrentInstructions()
	}
	fmt.Println("Нажмите Enter, чтобы продолжить...")
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Fatalf("can't start game: %v", err)
		}
		return
	}

	stopRandomService := startRandomServiceWork()
	defer stopRandomService()

	if countOfBoards == 1 {
		if err := service.HandleGameProcess(model.NewGame(cfg), size, scanner); err != nil {
			log.Fatalf("can't handle game process: %v", err)
		}
		return
	}

	if err := runConcurrentGames(cfg, size, countOfBoards, scanner); err != nil {
		log.Fatalf("can't handle concurrent games: %v", err)
	}
}

func runConcurrentGames(cfg *model.GameConfig, size, countOfBoards int, scanner *bufio.Scanner) error {
	games := make([]*model.Game, countOfBoards)
	commands := make([]chan model.GameMove, countOfBoards)
	for i := range countOfBoards {
		games[i] = model.NewGame(cfg)
		commands[i] = make(chan model.GameMove)
	}

	snapshots := make(chan model.StateSnapshot)
	results := make(chan model.CommandResult, countOfBoards)
	initialStateReady := make(chan struct{})

	var gameWG sync.WaitGroup
	var drawerWG sync.WaitGroup
	var outputMu sync.Mutex

	drawerWG.Go(func() {
		drawSnapshots(snapshots, size, countOfBoards, initialStateReady, &outputMu)
	})

	for i := range countOfBoards {
		gameWG.Go(func() {
			service.RunConcurrentGame(games[i], i, commands[i], snapshots, results)
		})
	}

	<-initialStateReady

	activeBoards := make(map[int]bool, countOfBoards)
	for i := range countOfBoards {
		activeBoards[i] = true
	}

	for len(activeBoards) > 0 {
		commandsSent := 0
		for boardNumber := range countOfBoards {
			if !activeBoards[boardNumber] {
				continue
			}

			command, err := readConcurrentCommand(scanner, boardNumber, &outputMu)
			if err != nil {
				if errors.Is(err, io.EOF) {
					closeActiveCommandChannels(commands, activeBoards)
					gameWG.Wait()
					close(snapshots)
					drawerWG.Wait()
					return nil
				}
				return err
			}

			commands[boardNumber] <- command
			commandsSent++
		}

		for range commandsSent {
			result := <-results
			if result.Err != nil {
				outputMu.Lock()
				fmt.Printf("Доска #%d: команда завершилась с ошибкой: %v\n", result.NumberOfBoard+1, result.Err)
				outputMu.Unlock()
			}
			if result.Finished {
				delete(activeBoards, result.NumberOfBoard)
			}
		}
	}

	gameWG.Wait()
	close(snapshots)
	drawerWG.Wait()
	return nil
}

func readConcurrentCommand(scanner *bufio.Scanner, boardNumber int, outputMu *sync.Mutex) (model.GameMove, error) {
	for {
		outputMu.Lock()
		fmt.Printf("Команда для доски #%d (Автоход N или Сдаться): ", boardNumber+1)
		outputMu.Unlock()

		move, _, err := console.ReadAndConvertPlayerInput(scanner)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return model.GameMove{}, io.EOF
			}
			outputMu.Lock()
			fmt.Printf("Не удалось обработать команду: %v\n", err)
			outputMu.Unlock()
			continue
		}

		if move.Type == model.Move {
			outputMu.Lock()
			fmt.Println("Ручные ходы доступны только в режиме одной доски")
			outputMu.Unlock()
			continue
		}

		return *move, nil
	}
}

func drawSnapshots(
	snapshots <-chan model.StateSnapshot,
	size int,
	initialStateCount int,
	initialStateReady chan<- struct{},
	outputMu *sync.Mutex,
) {
	initialStatesDrawn := 0
	for snapshot := range snapshots {
		outputMu.Lock()
		fmt.Printf("\n%s ДОСКА #%d %s\n", strings.Repeat("_", 20), snapshot.NumberOfBoard+1, strings.Repeat("_", 20))
		switch {
		case snapshot.Finished:
			fmt.Printf("Игрок %s сдался. Игра завершена.\n", snapshot.LastMovePlayerName)
		case snapshot.LastMovePlayerName != "":
			fmt.Printf(
				"Игрок %s выполнил ход за %s. Осталось автоходов: %d.\n",
				snapshot.LastMovePlayerName,
				snapshot.LastMoveDuration.Round(time.Millisecond),
				snapshot.RemainingAutoMoves,
			)
			fmt.Printf("Следующий игрок: %s.\n", snapshot.CurrentPlayerName)
		default:
			fmt.Printf("Текущий игрок: %s. Ожидание команды.\n", snapshot.CurrentPlayerName)
		}
		console.PrintBoard(snapshot.Board, size)
		fmt.Printf("%s\n", strings.Repeat("_", 50))
		outputMu.Unlock()

		if initialStatesDrawn < initialStateCount {
			initialStatesDrawn++
			if initialStatesDrawn == initialStateCount {
				close(initialStateReady)
			}
		}
	}
}

func closeActiveCommandChannels(commands []chan model.GameMove, activeBoards map[int]bool) {
	for boardNumber := range activeBoards {
		close(commands[boardNumber])
	}
}

func startRandomServiceWork() func() {
	repo := repository.NewRepository()
	randomService := service.NewRandomService(repo)
	done := make(chan struct{})
	var once sync.Once

	go func() {
		ticker := time.NewTicker(randomEntitySaveInterval)
		defer ticker.Stop()
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

	return func() {
		once.Do(func() {
			close(done)
		})
	}
}

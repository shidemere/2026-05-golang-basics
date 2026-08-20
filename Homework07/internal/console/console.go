// Package console need for work with console output and input
package console

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/shidemere/2026-05-golang-basics/Homework07/internal/model"
)

// PrintBoard prints cells as a square board with size cells in each row.
func PrintBoard(cells []model.Cell, size int) {
	if size <= 0 {
		return
	}

	for i, cell := range cells {
		if i%size == 0 {
			fmt.Printf("%3d", i/size+1)
		}

		if cell.HasPiece() {
			fmt.Printf("%3c", cell.GetPiece().GetValue())
		} else {
			fmt.Printf("%3c", *cell.GetValue())
		}
		if (i+1)%size == 0 {
			fmt.Println()
		}
	}

	fmt.Printf("%3s", "")
	for i := 0; i < size && i < len(cells); i++ {
		fmt.Printf("%3s", cells[i].GetColumn())
	}
	fmt.Println()
}

func CleanBoard() {
	fmt.Print("\033[H\033[2J")
}

func PrintInstructions() {
	fmt.Println("\n---------------------")
	fmt.Println(`
		Необходимо ввести команду. 
		В текущей реализации доступны 3 команды:
		1. Сдаться - прерывает игру
		2. Автоход {количество} - делает определенное количество ходов
		3. Ход {старая позиция} {новая позиция} - переводит одну из фигур из одной позиции в другую.
		`)
}

func PrintConcurrentInstructions() {
	fmt.Println("\n---------------------")
	fmt.Println(`
		Команды вводятся последовательно для каждой активной доски.
		В режиме нескольких досок доступны 2 команды:
		1. Сдаться - завершает игру на выбранной доске.
		2. Автоход {количество} - запускает указанное количество автоматических ходов.
		Ручные ходы в этом режиме недоступны.
		`)
}

func HandleBoardSizeInputAndCount() (int, int) {
	if len(os.Args) < 3 {
		fmt.Println("you need to run program with size and count of boards.\nExample: go run . 8 2")
		os.Exit(1)
	}
	arg := os.Args[1]
	val, err := strconv.Atoi(arg)
	if err != nil {
		log.Fatalf("incorrect input argument: %v", err)
	}
	if val < 2 {
		log.Fatalf("board size must be at least 2, got %d", val)
	}

	sec := os.Args[2]
	secV, err := strconv.Atoi(sec)
	if err != nil {
		log.Fatalf("incorrect second input argument: %v", err)
	}
	if secV <= 0 {
		log.Fatalf("count of boards must be positive, got %d", secV)
	}
	return val, secV
}

func HandlePlayersNames(scanner *bufio.Scanner) (first, second string, err error) {
	if scanner == nil {
		return "", "", errors.New("scanner is nil")
	}

	var counter int
	for {
		fmt.Printf("Пожалуйста, введите имя %d го игрока: ", counter+1)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", "", fmt.Errorf("failed to handle player name: %w", err)
			}
			return "", "", io.EOF
		}

		input := scanner.Text()
		if err := validatePlayerName(input); err != nil {
			return "", "", fmt.Errorf("failed to handle player name: %v", err)
		}

		if first == "" {
			first = input
			counter++
			continue
		}
		if second == "" {
			second = input
			counter++
		}

		if first != "" && second != "" {
			break
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

func ReadAndConvertPlayerInput(scanner *bufio.Scanner) (*model.GameMove, string, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, "", fmt.Errorf("can't process game move: %w", err)
		}
		return nil, "", io.EOF
	}

	input := scanner.Text()
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return nil, "", errors.New("команда не задана")
	}

	switch fields[0] {
	case "Сдаться":
		if len(fields) != 1 {
			return nil, "", errors.New("команда сдачи задается в формате: Сдаться")
		}
		return &model.GameMove{Type: model.GiveUP}, input, nil
	case "Ход":
		return &model.GameMove{Type: model.Move}, input, nil
	case "Автоход":
		if len(fields) != 2 {
			return nil, "", errors.New("неправильно задан автоход, необходимо задать в формате: Автоход {количество}")
		}

		count, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, "", fmt.Errorf("количество автоходов %q не является числом", fields[1])
		}
		if count <= 0 {
			return nil, "", fmt.Errorf("количество автоходов должно быть больше нуля, получено %d", count)
		}

		return &model.GameMove{Type: model.Auto, AutoMoveCount: count}, input, nil
	default:
		return nil, "", fmt.Errorf("неизвестная команда %q", fields[0])
	}
}

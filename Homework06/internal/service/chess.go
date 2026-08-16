// Package service provides logic
package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/console"
	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/model"
)

func HandleGameProcess(g *model.Game, size int) error {
	console.PrintInstructions()
	buffer := bufio.NewScanner(os.Stdin)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		currentPlayer := g.GetCurrentPlayer()
		fmt.Printf("Игрок %s, ожидание ввода: \n", currentPlayer.GetPlayerName())

		move, err := HandlePlayerInput(g.GetBoard(), currentPlayer, buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			fmt.Printf("Не удалось обработать команду: %v\n", err)
			continue
		}

		switch move.Type {
		case model.GiveUP:
			fmt.Printf("%s сдался. Игра окончена\n", currentPlayer.GetPlayerName())
			defineWinner(chooseInactivePlayer(g))
			return nil
		case model.Auto:
			if err := makeAutoMoves(g, move.AutoMoveCount, size, random, time.Sleep); err != nil {
				fmt.Printf("Не удалось выполнить автоход: %v\n", err)
			}
		case model.Move:
			completeMove(g, move, size)
		default:
			fmt.Printf("Неизвестный тип хода: %d\n", move.Type)
		}
	}
}

func HandlePlayerInput(b *model.Board, player *model.Player, scanner *bufio.Scanner) (*model.GameMove, error) {
	move, input, err := console.ReadAndConvertPlayerInput(scanner)
	if err != nil {
		return nil, fmt.Errorf("can't read or convert player input: %w", err)
	}

	switch move.Type {
	case model.GiveUP:
		return move, nil
	case model.Move:
		fields := strings.Fields(input)
		if len(fields) != 3 {
			return nil, errors.New("неправильно задан ход, необходимо задать в формате: ход {старая позиция} {новая позиция}")
		}
		oldM, err := getCellByString(fields[1], b)
		if err != nil {
			return nil, fmt.Errorf("cant handler player input: %v", err)
		}
		newM, err := getCellByString(fields[2], b)
		if err != nil {
			return nil, fmt.Errorf("cant handler player input: %v", err)
		}
		return makeMove(oldM, newM, player)
	case model.Auto:
		return move, nil
	}

	return nil, fmt.Errorf("неподдерживаемый тип хода: %d", move.Type)
}

func makeAutoMoves(g *model.Game, count, size int, random *rand.Rand, sleep func(time.Duration)) error {
	for range count {
		player := g.GetCurrentPlayer()
		delay := time.Duration(random.Intn(3)+2) * time.Second
		sleep(delay)

		move, err := makeRandomMove(g.GetBoard(), player, random)
		if err != nil {
			return fmt.Errorf("игрок %s: %w", player.GetPlayerName(), err)
		}

		completeMove(g, move, size)
	}

	return nil
}

func makeRandomMove(b *model.Board, player *model.Player, random *rand.Rand) (*model.GameMove, error) {
	cells := b.GetCells()
	sourceCells := make([]*model.Cell, 0)
	destinationCells := make([]*model.Cell, 0)

	for i := range cells {
		cell := &cells[i]
		switch {
		case cell.HasPiece() && cell.GetPiece() != nil && cell.GetPiece().GetColor() == player.GetColor():
			sourceCells = append(sourceCells, cell)
		case !cell.HasPiece():
			destinationCells = append(destinationCells, cell)
		}
	}

	if len(sourceCells) == 0 {
		return nil, errors.New("нет фигур для перемещения")
	}
	if len(destinationCells) == 0 {
		return nil, errors.New("нет свободных клеток")
	}

	oldCell := sourceCells[random.Intn(len(sourceCells))]
	newCell := destinationCells[random.Intn(len(destinationCells))]
	return makeMove(oldCell, newCell, player)
}

func getCellByString(input string, b *model.Board) (*model.Cell, error) {
	cell, err := connectMoveStringToCell(input, b)
	if err != nil {
		return nil, fmt.Errorf("can't parse move from input because cant connect string to cell: %v", err)
	}

	return cell, nil
}

func makeMove(oldM *model.Cell, newM *model.Cell, player *model.Player) (*model.GameMove, error) {
	if !oldM.HasPiece() || oldM.GetPiece() == nil {
		return nil, model.MoveChessPieceNotExistError{
			Position: oldM.GetColumn() + strconv.Itoa(oldM.GetLine()),
		}
	}

	if oldM.GetPiece().GetColor() != player.GetColor() {
		return nil, &model.ColorMismatchError{
			PlayerName:  player.GetPlayerName(),
			PlayerColor: player.GetColor(),
			ChessColor:  oldM.GetPiece().GetColor(),
		}
	}
	piece := oldM.GetPiece()
	newM.SetHasPiece(true)
	newM.SetPiece(piece)
	oldM.SetHasPiece(false)
	oldM.SetPiece(nil)
	piece.SetCurrentX(newM.GetLine())
	piece.SetCurrentY(newM.GetColumn())

	result := &model.GameMove{
		Type:          model.Move,
		OldPosition:   oldM,
		NewPosition:   newM,
		CurrentPlayer: player,
		Piece:         newM.GetPiece(),
	}
	return result, nil
}

func completeMove(g *model.Game, move *model.GameMove, size int) {
	g.AddMove(*move)
	g.ChangeCurrentPlayer()
	console.CleanBoard()
	console.PrintBoard(g.GetBoard().GetCells(), size)
}

func connectMoveStringToCell(s string, b *model.Board) (*model.Cell, error) {
	splited := strings.Split(s, "")
	if len(splited) < 2 {
		return nil, errors.New("при указывании позиции для хода необходимо следовать формату {колонка}{строка}")
	}

	var column strings.Builder
	clEnd := 0
	for i, v := range splited {
		if unicode.IsUpper([]rune(v)[0]) {
			column.WriteString(v)
		} else {
			if clEnd == 0 {
				clEnd = i
			}
		}
	}

	ln := strings.Join(splited[clEnd:], "")

	line, err := strconv.Atoi(ln)
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать номер строки: %v", err)
	}

	for i, v := range b.GetCells() {
		if v.GetColumn() == column.String() && v.GetLine() == line {
			return &b.GetCells()[i], nil
		}
	}

	return nil, errors.New("для такой позиции не найдена ячейка")
}

func chooseInactivePlayer(g *model.Game) *model.Player {
	if g.GetCurrentPlayer() == g.GetSecondPlayer() {
		return g.GetFirstPlayer()
	}
	return g.GetSecondPlayer()
}

func defineWinner(p *model.Player) {
	fmt.Printf("Поздравляем игрока %s с победой!\n", p.GetPlayerName())
}

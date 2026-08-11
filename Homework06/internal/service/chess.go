// Package service provides logic
package service

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/console"
	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/model"
)

var playerCount int = 0

func HandleGameProcess(g *model.Game, size int) error {
	console.PrintInstructions()
	buffer := bufio.NewScanner(os.Stdin)
	for {
		currentPlayer := g.GetCurrentPlayer()
		fmt.Printf("Игрок %s, ожидание ввода: \n", currentPlayer.GetPlayerName())
		move, err := HandlePlayerInput(g.GetBoard(), currentPlayer, buffer)
		if err != nil {
			return fmt.Errorf("cant handle game process: %v", err)
		}
		if move.Type == model.GiveUP {
			fmt.Printf("%s сдался. Игра окончена\n", currentPlayer.GetPlayerName())
			defineWinner(chooseInactivePlayer(g))
			break
		}

		g.AddMove(*move)
		g.ChangeCurrentPlayer()
		console.CleanBoard()
		console.PrintBoard(g.GetBoard().GetCells(), size)
	}
	return nil
}

func HandlePlayerInput(b *model.Board, player *model.Player, scanner *bufio.Scanner) (*model.GameMove, error) {
	move, input, err := console.ReadAndConverPlayerInput(b, player, scanner)
	if err != nil {
		return nil, fmt.Errorf("can't read or convert plyaer input: %v", err)
	}
	switch move.Type {
	case model.GiveUP:
		return move, nil
	case model.Bishop:
		splited := strings.Split(input, " ")
		if len(splited) != 3 {
			return nil, errors.New("неправильно задан ход, необходимо задать в формате: ход {старая позиция} {новая позиция}")
		}

		// need to take second and third word, and push it in the separate method which will return to me Cell
		// then i need to build GameMove with it
		oldM, err := connectMoveStringToCell(splited[1], b)
		if err != nil {
			return nil, fmt.Errorf("can't parse move from input because old cell: %v", err)
		}

		newM, err := connectMoveStringToCell(splited[2], b)
		if err != nil {
			return nil, fmt.Errorf("can't parse move from input because new cell: %v", err)
		}

		// need to check if piece exist for moving
		if !oldM.HasPiece() || oldM.GetPiece() == nil {
			return nil, &model.MoveChessPieceNotExistError{splited[1]}
		}

		piece := oldM.GetPiece()
		// move piece and build GameMove
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
	case model.Auto:
		splited := strings.Split(input, " ")
		if len(splited) != 2 {
			return nil, errors.New("неправильно задан автоход, необходимо задать в формате: Автоход {количество}")
		}

		// TODO I don't clearly understand how i need to do this shit
		// actually i need to think about some algorithm for generating movies
		// Howewer, I also need to keep count of auto moves outside, because i can do only one single move in one time

	}

	return nil, errors.New("nothing to read")
}

// TODO it will break with input >10 because we have 3 chars, not two. And split for three is not work well now
func connectMoveStringToCell(s string, b *model.Board) (*model.Cell, error) {
	splited := strings.Split(s, "")
	if len(splited) != 2 {
		return nil, errors.New("при указывании позиции для хода необходимо следовать формату {колонка}{строка}")
	}
	column := splited[0]
	line, err := strconv.Atoi(splited[1])
	if err != nil {
		return nil, fmt.Errorf("не удалось преобразовать номер строки: %v", err)
	}

	for i, v := range b.GetCells() {
		if v.GetColumn() == column && v.GetLine() == line {
			return &b.GetCells()[i], nil
		}
	}

	return nil, errors.New("для такой позиции не найдена ячейка")
}

func chooseInactivePlayer(g *model.Game) *model.Player {
	if playerCount%2 == 0 {
		return g.GetSecondPlayer()
	} else {
		return g.GetFirstPlayer()
	}
}

func chooseCurrentPlayer(g *model.Game) *model.Player {
	if playerCount%2 == 0 {
		return g.GetFirstPlayer()
	} else {
		return g.GetSecondPlayer()
	}
}

func defineWinner(p *model.Player) {
	fmt.Printf("Поздравляем игрока %s с победой!\n", p.GetPlayerName())
}

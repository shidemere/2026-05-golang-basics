// Package service provides logic
package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/shidemere/2026-05-golang-basics/Homework07/internal/console"
	"github.com/shidemere/2026-05-golang-basics/Homework07/internal/model"
)

func HandleGameProcess(
	g *model.Game,
	size int,
	buffer *bufio.Scanner,
) error {
	if g == nil {
		return errors.New("game is nil")
	}
	if buffer == nil {
		return errors.New("scanner is nil")
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	console.PrintBoard(g.GetBoard().GetCells(), size)

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
			if err := makeAutoMoves(g, move.AutoMoveCount, random, time.Sleep, func(player string, duration time.Duration, _ int) {
				fmt.Printf("Игрок %s выполнил ход за %s\n", player, duration.Round(time.Millisecond))
				console.PrintBoard(g.GetBoard().GetCells(), size)
			}); err != nil {
				fmt.Printf("Не удалось выполнить автоход: %v\n", err)
			}
		case model.Move:
			completeMove(g, move)
			console.PrintBoard(g.GetBoard().GetCells(), size)
		default:
			fmt.Printf("Неизвестный тип хода: %d\n", move.Type)
		}
	}
}

func RunConcurrentGame(
	g *model.Game,
	numberOfBoard int,
	commands <-chan model.GameMove,
	drawer chan<- model.StateSnapshot,
	results chan<- model.CommandResult,
) {
	if g == nil {
		results <- model.CommandResult{NumberOfBoard: numberOfBoard, Finished: true, Err: errors.New("game is nil")}
		return
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	sendStateToDrawer(g, drawer, numberOfBoard, "", 0, 0, false)

	for command := range commands {
		switch command.Type {
		case model.GiveUP:
			player := g.GetCurrentPlayer()
			sendStateToDrawer(g, drawer, numberOfBoard, player.GetPlayerName(), 0, 0, true)
			results <- model.CommandResult{NumberOfBoard: numberOfBoard, Finished: true}
			return
		case model.Auto:
			err := makeAutoMoves(g, command.AutoMoveCount, random, time.Sleep, func(player string, duration time.Duration, remaining int) {
				sendStateToDrawer(g, drawer, numberOfBoard, player, duration, remaining, false)
			})
			results <- model.CommandResult{NumberOfBoard: numberOfBoard, Err: err}
		default:
			results <- model.CommandResult{
				NumberOfBoard: numberOfBoard,
				Err:           fmt.Errorf("unsupported command type %d", command.Type),
			}
		}
	}
}

func sendStateToDrawer(
	g *model.Game,
	drawer chan<- model.StateSnapshot,
	number int,
	lastMovePlayer string,
	lastMoveDuration time.Duration,
	remaining int,
	finished bool,
) {
	cells := g.GetBoard().GetCells()
	boardCopy := make([]model.Cell, len(cells))
	copy(boardCopy, cells)

	snapshot := model.StateSnapshot{
		NumberOfBoard:      number,
		Board:              boardCopy,
		CurrentPlayerName:  g.GetCurrentPlayer().GetPlayerName(),
		LastMovePlayerName: lastMovePlayer,
		LastMoveDuration:   lastMoveDuration,
		RemainingAutoMoves: remaining,
		Finished:           finished,
	}
	drawer <- snapshot
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

func makeAutoMoves(
	g *model.Game,
	count int,
	random *rand.Rand,
	sleep func(time.Duration),
	afterMove func(player string, duration time.Duration, remaining int),
) error {
	for moveNumber := range count {
		player := g.GetCurrentPlayer()
		delay := time.Duration(random.Intn(3)+2) * time.Second
		startedAt := time.Now()
		sleep(delay)

		move, err := makeRandomMove(g.GetBoard(), player, random)
		if err != nil {
			return fmt.Errorf("игрок %s: %w", player.GetPlayerName(), err)
		}

		completeMove(g, move)
		if afterMove != nil {
			afterMove(player.GetPlayerName(), time.Since(startedAt), count-moveNumber-1)
		}
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

func completeMove(g *model.Game, move *model.GameMove) {
	g.AddMove(*move)
	g.ChangeCurrentPlayer()
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

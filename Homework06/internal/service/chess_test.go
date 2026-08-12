package service

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/model"
)

func TestMakeAutoMovesAlternatesPlayers(t *testing.T) {
	game := model.NewGame(&model.GameConfig{
		FirstPlayerName:  "first",
		SecondPlayerName: "second",
		Size:             8,
	})

	whiteBefore := occupiedPositions(game.GetBoard(), model.PlayerColorWhite)
	blackBefore := occupiedPositions(game.GetBoard(), model.PlayerColorBlack)
	delays := make([]time.Duration, 0, 2)

	err := makeAutoMoves(
		game,
		2,
		8,
		rand.New(rand.NewSource(1)),
		func(delay time.Duration) {
			delays = append(delays, delay)
		},
	)
	if err != nil {
		t.Fatalf("makeAutoMoves returned an error: %v", err)
	}

	if game.GetCurrentPlayer() != game.GetFirstPlayer() {
		t.Fatal("after two auto moves the turn must return to the first player")
	}
	if positionsEqual(whiteBefore, occupiedPositions(game.GetBoard(), model.PlayerColorWhite)) {
		t.Error("the first player's pieces did not move")
	}
	if positionsEqual(blackBefore, occupiedPositions(game.GetBoard(), model.PlayerColorBlack)) {
		t.Error("the second player's pieces did not move")
	}
	if got := occupiedCount(game.GetBoard()); got != 32 {
		t.Fatalf("auto moves changed the number of pieces: got %d, want 32", got)
	}
	if len(delays) != 2 {
		t.Fatalf("got %d delays, want 2", len(delays))
	}
	for _, delay := range delays {
		if delay < 2*time.Second || delay > 4*time.Second {
			t.Errorf("delay %s is outside the 2-4 second range", delay)
		}
	}
}

func occupiedPositions(board *model.Board, color model.Color) map[string]struct{} {
	positions := make(map[string]struct{})
	for _, cell := range board.GetCells() {
		if cell.HasPiece() && cell.GetPiece() != nil && cell.GetPiece().GetColor() == color {
			positions[cell.GetColumn()+strconv.Itoa(cell.GetLine())] = struct{}{}
		}
	}
	return positions
}

func positionsEqual(first, second map[string]struct{}) bool {
	if len(first) != len(second) {
		return false
	}
	for position := range first {
		if _, ok := second[position]; !ok {
			return false
		}
	}
	return true
}

func occupiedCount(board *model.Board) int {
	count := 0
	for _, cell := range board.GetCells() {
		if cell.HasPiece() {
			count++
		}
	}
	return count
}

package service

import (
	"math/rand"
	"testing"
	"time"

	"github.com/shidemere/2026-05-golang-basics/Homework07/internal/model"
)

func TestMakeAutoMovesReportsEveryMove(t *testing.T) {
	game := model.NewGame(&model.GameConfig{
		FirstPlayerName:  "first",
		SecondPlayerName: "second",
		Size:             8,
	})
	random := rand.New(rand.NewSource(1))
	remainingMoves := make([]int, 0, 3)

	err := makeAutoMoves(
		game,
		3,
		random,
		func(time.Duration) {},
		func(_ string, _ time.Duration, remaining int) {
			remainingMoves = append(remainingMoves, remaining)
		},
	)
	if err != nil {
		t.Fatalf("makeAutoMoves returned an error: %v", err)
	}

	want := []int{2, 1, 0}
	if len(remainingMoves) != len(want) {
		t.Fatalf("got %d updates, want %d", len(remainingMoves), len(want))
	}
	for i := range want {
		if remainingMoves[i] != want[i] {
			t.Errorf("update %d has %d remaining moves, want %d", i, remainingMoves[i], want[i])
		}
	}
}

func TestRunConcurrentGameReportsInitialAndFinalStateOnGiveUp(t *testing.T) {
	game := model.NewGame(&model.GameConfig{
		FirstPlayerName:  "first",
		SecondPlayerName: "second",
		Size:             8,
	})
	commands := make(chan model.GameMove, 1)
	snapshots := make(chan model.StateSnapshot, 2)
	results := make(chan model.CommandResult, 1)
	commands <- model.GameMove{Type: model.GiveUP}
	close(commands)

	RunConcurrentGame(game, 2, commands, snapshots, results)

	initial := <-snapshots
	if initial.NumberOfBoard != 2 || initial.Finished {
		t.Fatalf("unexpected initial snapshot: %+v", initial)
	}
	final := <-snapshots
	if final.NumberOfBoard != 2 || !final.Finished {
		t.Fatalf("unexpected final snapshot: %+v", final)
	}
	result := <-results
	if result.NumberOfBoard != 2 || !result.Finished || result.Err != nil {
		t.Fatalf("unexpected result: %+v", result)
	}
}

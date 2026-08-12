package service

import (
	"math/rand"
	"testing"

	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/repository"
)

func TestRandomServiceCreatesValidEntities(t *testing.T) {
	service := &RandomService{random: rand.New(rand.NewSource(1))}

	player := service.CreateRandomPlayer()
	if player == nil {
		t.Fatal("CreateRandomPlayer returned nil")
	}
	if player.GetPlayerName() == "" {
		t.Error("random player has an empty name")
	}
	if len(player.GetChessPieces()) == 0 {
		t.Error("random player has no chess pieces")
	}

	piece := service.CreateRandomChessPiece()
	if piece == nil {
		t.Fatal("CreateRandomChessPiece returned nil")
	}
	if piece.GetCurrentX() <= 0 || piece.GetCurrentY() == "" {
		t.Error("random chess piece has an invalid position")
	}

	move := service.CreateRandomMove()
	if move == nil {
		t.Fatal("CreateRandomMove returned nil")
	}
	if move.CurrentPlayer == nil || move.Piece == nil || move.OldPosition == nil || move.NewPosition == nil {
		t.Error("random move is incomplete")
	}
}

func TestRandomServiceSavesOneEntityPerCall(t *testing.T) {
	repo := repository.NewRepository()
	service := &RandomService{
		Repository: repo,
		random:     rand.New(rand.NewSource(1)),
	}

	for i := 1; i <= 10; i++ {
		if err := service.Save(); err != nil {
			t.Fatalf("Save returned an error: %v", err)
		}

		total := len(repo.Players()) + len(repo.Pieces()) + len(repo.Moves())
		if total != i {
			t.Fatalf("after %d calls repository contains %d entities", i, total)
		}
	}
}

package repository

import (
	"testing"

	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/model"
)

func TestRepositoryDistributesEntitiesByType(t *testing.T) {
	repo := NewRepository()
	game := model.NewGame(&model.GameConfig{
		FirstPlayerName:  "first",
		SecondPlayerName: "second",
		Size:             8,
	})
	player := game.GetFirstPlayer()
	piece := player.GetChessPieces()[0]
	move := &model.GameMove{Type: model.Move, CurrentPlayer: player, Piece: piece}

	for _, entity := range []model.Entity{player, piece, move} {
		if err := repo.Save(entity); err != nil {
			t.Fatalf("Save(%T) returned an error: %v", entity, err)
		}
	}

	if got := len(repo.Players()); got != 1 {
		t.Errorf("got %d players, want 1", got)
	}
	if got := len(repo.Pieces()); got != 1 {
		t.Errorf("got %d pieces, want 1", got)
	}
	if got := len(repo.Moves()); got != 1 {
		t.Errorf("got %d moves, want 1", got)
	}
}

func TestRepositoryRejectsNilEntities(t *testing.T) {
	repo := NewRepository()

	if err := repo.Save(nil); err == nil {
		t.Error("Save(nil) returned no error")
	}

	var player *model.Player
	if err := repo.Save(player); err == nil {
		t.Error("Save((*model.Player)(nil)) returned no error")
	}
}

func TestRepositoryReturnsSliceCopies(t *testing.T) {
	repo := NewRepository()
	game := model.NewGame(&model.GameConfig{
		FirstPlayerName:  "first",
		SecondPlayerName: "second",
		Size:             8,
	})
	if err := repo.Save(game.GetFirstPlayer()); err != nil {
		t.Fatalf("Save returned an error: %v", err)
	}

	players := repo.Players()
	players[0] = nil

	if repo.Players()[0] == nil {
		t.Error("Players exposed the repository's internal slice")
	}
}

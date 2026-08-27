// Package repository provides work with data
package repository

import (
	"errors"
	"fmt"
	"sync"

	"github.com/shidemere/2026-05-golang-basics/Homework08/internal/model"
)

type Repository struct {
	mu      sync.RWMutex
	players []*model.Player
	moves   []*model.GameMove
	pieces  []*model.ChessPiece
}

func NewRepository() *Repository {
	return &Repository{
		players: make([]*model.Player, 0),
		moves:   make([]*model.GameMove, 0),
		pieces:  make([]*model.ChessPiece, 0),
	}
}

func (r *Repository) Save(e model.Entity) error {
	if r == nil {
		return errors.New("repository is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	switch v := e.(type) {
	case *model.Player:
		if v == nil {
			return errors.New("player is nil")
		}
		r.players = append(r.players, v)
	case *model.ChessPiece:
		if v == nil {
			return errors.New("chess piece is nil")
		}
		r.pieces = append(r.pieces, v)
	case *model.GameMove:
		if v == nil {
			return errors.New("game move is nil")
		}
		r.moves = append(r.moves, v)
	default:
		return fmt.Errorf("unsupported entity type %T", e)
	}

	return nil
}

func (r *Repository) Players() []*model.Player {
	if r == nil {
		return nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.Player, len(r.players))
	copy(result, r.players)
	return result
}

func (r *Repository) Moves() []*model.GameMove {
	if r == nil {
		return nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.GameMove, len(r.moves))
	copy(result, r.moves)
	return result
}

func (r *Repository) Pieces() []*model.ChessPiece {
	if r == nil {
		return nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.ChessPiece, len(r.pieces))
	copy(result, r.pieces)
	return result
}

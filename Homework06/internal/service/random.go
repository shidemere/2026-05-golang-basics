package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/model"
	"github.com/shidemere/2026-05-golang-basics/Homework06/internal/repository"
)

type RandomService struct {
	Repository *repository.Repository
	random     *rand.Rand
}

func NewRandomService(repo *repository.Repository) *RandomService {
	return &RandomService{
		Repository: repo,
		random:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *RandomService) CreateRandomMove() *model.GameMove {
	game := s.createRandomGame()
	player := game.GetFirstPlayer()
	if s.randomSource().Intn(2) == 1 {
		player = game.GetSecondPlayer()
	}

	move, err := makeRandomMove(game.GetBoard(), player, s.randomSource())
	if err != nil {
		return nil
	}

	return move
}

func (s *RandomService) CreateRandomPlayer() *model.Player {
	game := s.createRandomGame()
	if s.randomSource().Intn(2) == 0 {
		return game.GetFirstPlayer()
	}

	return game.GetSecondPlayer()
}

func (s *RandomService) CreateRandomChessPiece() *model.ChessPiece {
	player := s.CreateRandomPlayer()
	pieces := player.GetChessPieces()
	return pieces[s.randomSource().Intn(len(pieces))]
}

func (s *RandomService) Save() error {
	if s.Repository == nil {
		return errors.New("repository is nil")
	}

	var entity model.Entity
	switch s.randomSource().Intn(3) {
	case 0:
		entity = s.CreateRandomPlayer()
	case 1:
		entity = s.CreateRandomChessPiece()
	case 2:
		entity = s.CreateRandomMove()
	}

	if entity == nil {
		return errors.New("can't create random entity")
	}

	return s.Repository.Save(entity)
}

func (s *RandomService) createRandomGame() *model.Game {
	return model.NewGame(&model.GameConfig{
		FirstPlayerName:  s.randomPlayerName(),
		SecondPlayerName: s.randomPlayerName(),
		Size:             8,
	})
}

func (s *RandomService) randomPlayerName() string {
	return fmt.Sprintf("Player-%06d", s.randomSource().Intn(1_000_000))
}

func (s *RandomService) randomSource() *rand.Rand {
	if s.random == nil {
		s.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	return s.random
}

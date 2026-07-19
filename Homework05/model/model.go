// Package model provides structs for playing in Chess
package model

type Color int

const (
	PlayerColorWhite = iota
	PlayerColorBlack
)

type ChessPieceType rune

const (
	King = iota
	Queen
	Rook
	Bishop
	Knight
	Pawn
)

type Game struct {
	first  *Player
	second *Player
	board  *Board
	moves  []GameMove
}

func (g *Game) GetBoard() *Board {
	return g.board
}

type GameMove struct {
	NewPosition   *Cell
	OldPosition   *Cell
	CurrentPlayer *Player
	Piece         *ChessPiece
}

type GameConfig struct {
	FirstPlayerName  string
	SecondPlayerName string
	Size             int
}

type Player struct {
	name   string
	color  Color
	pieces []*ChessPiece
}

type ChessPiece struct {
	value    ChessPieceType
	color    Color
	currentX int
	currentY string
}

type Board struct {
	cells []Cell
}

func (b *Board) GetCells() []Cell {
	return b.cells
}

type Cell struct {
	line     int
	column   string
	value    *rune
	hasPiece bool
	piece    *ChessPiece
}

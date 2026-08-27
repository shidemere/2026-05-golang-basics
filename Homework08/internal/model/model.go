// Package model provides structs for playing in Chess
package model

import (
	"fmt"
	"time"
)

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

type MoveType int

const (
	GiveUP = iota
	Move
	Auto
)

type MoveChessPieceNotExistError struct {
	Position string
}

func (m MoveChessPieceNotExistError) Error() string {
	return fmt.Sprintf("on position %s piece does not exist, nothing to move", m.Position)
}

type ColorMismatchError struct {
	PlayerName  string
	PlayerColor Color
	ChessColor  Color
}

func (c ColorMismatchError) Error() string {
	return fmt.Sprintf("Player with name %s and color %v can't move chess piece with color %v", c.PlayerName, c.PlayerColor, c.ChessColor)
}

type Game struct {
	first         *Player
	second        *Player
	currentPlayer *Player
	board         *Board
	moves         []GameMove
}

func (g *Game) GetBoard() *Board {
	return g.board
}

func (g *Game) GetFirstPlayer() *Player {
	return g.first
}

func (g *Game) GetSecondPlayer() *Player {
	return g.second
}

func (g *Game) GetCurrentPlayer() *Player {
	return g.currentPlayer
}

func (g *Game) ChangeCurrentPlayer() {
	switch g.currentPlayer {
	case g.first:
		g.currentPlayer = g.second
	case g.second:
		g.currentPlayer = g.first
	}
}

type Entity interface {
	isEntity() bool
}
type GameMove struct {
	Type          MoveType
	NewPosition   *Cell
	OldPosition   *Cell
	CurrentPlayer *Player
	Piece         *ChessPiece
	AutoMoveCount int
}

func (m *GameMove) isEntity() bool {
	return true
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

func (p *Player) isEntity() bool {
	return true
}

func (p *Player) GetChessPieces() []*ChessPiece {
	return p.pieces
}

func (p *Player) GetPlayerName() string {
	return p.name
}

func (p *Player) GetColor() Color {
	return p.color
}

type ChessPiece struct {
	value    ChessPieceType
	color    Color
	currentX int
	currentY string
}

func (c *ChessPiece) isEntity() bool {
	return true
}

func (p *ChessPiece) GetValue() ChessPieceType {
	return p.value
}

func (p *ChessPiece) GetCurrentX() int {
	return p.currentX
}

func (p *ChessPiece) GetCurrentY() string {
	return p.currentY
}

func (p *ChessPiece) SetCurrentX(i int) {
	p.currentX = i
}

func (p *ChessPiece) SetCurrentY(s string) {
	p.currentY = s
}

func (p *ChessPiece) GetColor() Color {
	return p.color
}

type Board struct {
	cells []Cell
}

func (b *Board) GetCells() []Cell {
	return b.cells
}

func (g *Game) AddMove(move GameMove) {
	g.moves = append(g.moves, move)
}

type Cell struct {
	line     int
	column   string
	value    *rune
	hasPiece bool
	piece    *ChessPiece
}

func (c *Cell) GetLine() int {
	return c.line
}

func (c *Cell) SetLine(i int) {
	c.line = i
}

func (c *Cell) GetColumn() string {
	return c.column
}

func (c *Cell) SetColumn(s string) {
	c.column = s
}

func (c *Cell) GetValue() *rune {
	return c.value
}

func (c *Cell) SetValue(r *rune) {
	c.value = r
}

func (c *Cell) HasPiece() bool {
	return c.hasPiece
}

func (c *Cell) SetHasPiece(b bool) {
	c.hasPiece = b
}

func (c *Cell) GetPiece() *ChessPiece {
	return c.piece
}

func (c *Cell) SetPiece(p *ChessPiece) {
	c.piece = p
}

type StateSnapshot struct {
	NumberOfBoard      int
	Board              []Cell
	CurrentPlayerName  string
	LastMovePlayerName string
	LastMoveDuration   time.Duration
	RemainingAutoMoves int
	Finished           bool
}

type CommandResult struct {
	NumberOfBoard int
	Finished      bool
	Err           error
}

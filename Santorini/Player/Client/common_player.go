package client

import (
	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	common "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	result "github.com/CS4500-F18/dare-rebr/Santorini/Common/Tournament"
	strategy "github.com/CS4500-F18/dare-rebr/Santorini/Player/Strategy"
)

// Player data required
type player struct {
	// The player's name
	name string

	// This player's opponent
	opponent string

	// The strategy this player is using
	strategy strategy.IStrategy
}

// Get the name of this Player
func (p player) Name() string {
	return p.name
}

func (p player) SetName(newName string) {
	p.name = newName
	p.strategy.SetName(newName)
}

// Set this player's opponent
func (p player) SetOpponent(name string) {
	p.opponent = name
	p.strategy.SetOpponent(name)
}

// Get the location to place your next worker
func (p player) PlaceWorker(b board.IBoard) board.Pos {
	return p.strategy.WorkerPlacement(b)
}

// Get the next turn, including which worker to act on
func (p player) NextTurn(b board.IBoard) common.Turn {
	return p.strategy.WorkerTurn(b)
}

func (p player) ReceiveTournamentResults(result []result.MatchResult) {
	//There you go, enjoy.
}

func (p player) Opponent() string {
	return p.opponent
}

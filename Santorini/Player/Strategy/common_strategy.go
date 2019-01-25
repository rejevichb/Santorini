package strategy

import (
	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	iplayer "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
)

//Represents a strategy object that can take a turn given the state of the board
type IStrategy interface {
	//Set the name of the player we are creating moves for
	SetName(name string)

	//Set the name of the player this strategy is playing against
	SetOpponent(name string)

	//Returns the Turn to take for the given Board state, including which Worker to act on
	//If the move wins, the build will be bogus
	WorkerTurn(board.IBoard) iplayer.Turn

	//Returns the location to place a Worker given the Board state
	WorkerPlacement(board.IBoard) board.Pos
}

// Helper types defining functions from board and player name to position/turn
type placeStrategy func(b board.IBoard, player, opponent string) (board.Pos, error)
type turnStrategy func(b board.IBoard, player, opponent string) (iplayer.Turn, error)

// A structure to hold data relevant to a Strategy
type basicStrategy struct {
	placeIdea placeStrategy
	turnIdea  turnStrategy
	player    string
	opponent  string
}

func (b *basicStrategy) SetName(name string) {
	b.player = name
}

func (b *basicStrategy) SetOpponent(name string) {
	b.opponent = name
}

// Execute the placement strategy
func (b *basicStrategy) WorkerPlacement(board board.IBoard) board.Pos {
	p, _ := b.placeIdea(board, b.player, b.opponent)
	// Do something with error
	return p
}

// Execute the turn generation strategy
func (b *basicStrategy) WorkerTurn(board board.IBoard) iplayer.Turn {
	t, _ := b.turnIdea(board, b.player, b.opponent)
	// Do something with error
	return t
}

//Creates a new strategy that can determine turns and worker positions
func NewStrategy(player string, placeIdea placeStrategy, turnIdea turnStrategy) IStrategy {
	return &basicStrategy{
		placeIdea: placeIdea,
		turnIdea:  turnIdea,
		player:    player,
	}
}

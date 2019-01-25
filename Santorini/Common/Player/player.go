package player

import (
	"github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	result "github.com/CS4500-F18/dare-rebr/Santorini/Common/Tournament"
)

//The purpose of the player package is to encapsulate the work that a player needs to be responsible for. This includes
//going through the setup of a game, taking turns (including deciding on a turn based on the current state of the game),
//and leaving the game once finished.

//Represents a player of the game Santorini
type IPlayer interface {
	//Returns the name of the string identifying this player.
	Name() string

	//Assign a new Name to this Player
	SetName(newName string)

	//inform the player of who they are playing against
	SetOpponent(opponent string)

	//Returns the next location on the given IBoard that this Player
	//would like to place a Worker at
	PlaceWorker(b board.IBoard) board.Pos

	//Returns the next move and build, and the worker they wish to act on,
	//from the state of the given IBoard
	NextTurn(b board.IBoard) Turn

	//Receives the results of the tournament for this player.
	ReceiveTournamentResults(result []result.MatchResult)

	//Returns the string name of the opponent player
	Opponent() string
}

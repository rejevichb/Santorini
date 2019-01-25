package sandbox

import (
	"time"

	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	iplayer "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	result "github.com/CS4500-F18/dare-rebr/Santorini/Common/Tournament"
)

const (
	TIMEOUT_UNIT    = time.Millisecond
	TIMEOUT_DEFAULT = 10000 // 10000 milliseconds = 10 seconds
)

//WrappedPlayer represents a Player interface wrapped in an option type,
//such that the Player's response calls could fail, which would result in an error
//returned from the WrappedPlayer
type WrappedPlayer interface {
	// Method signatures remain the same as in the Player interface (Common/Player),
	// with the added functionality of an error if communication to the Player fails
	Name() (string, error)
	SetName(newName string) error
	SetOpponent(opponent string) error
	PlaceWorker(b board.IBoard) (board.Pos, error)
	NextTurn(b board.IBoard) (iplayer.Turn, error)
	ReceiveTournamentResult(result result.TournamentResult) error
}

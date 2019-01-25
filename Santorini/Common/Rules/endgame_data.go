package rules

import "errors"

//Constants for error messages for internal debugging
const (
	RULE_BROKEN_MSG  = "rules violation"
	CANNOT_MOVE_MSG  = "player lost the game due to no valid moves possible"
	WINNING_MOVE_MSG = "player won the game via a valid move"
)

var (
	RULE_BROKEN_ERR  = errors.New(RULE_BROKEN_MSG)
	CANNOT_MOVE_ERR  = errors.New(CANNOT_MOVE_MSG)
	WINNING_MOVE_ERR = errors.New(WINNING_MOVE_MSG)
)

//Represents the end of a Game
type GameResult struct {
	Winner     string
	Loser      string
	Reason     string
	BrokenRule bool
}

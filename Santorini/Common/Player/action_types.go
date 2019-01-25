package player

import (
	"errors"

	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	rules "github.com/CS4500-F18/dare-rebr/Santorini/Common/Rules"
)

//Represents the specification of a whole turn by a strategy.
type Turn struct {
	//The id of the worker selected for this turn.
	WID int

	//The position on the board the selected worker should move to.
	MoveTo board.Pos

	//The position on the board the selected worker should build on,
	//once it has moved.
	BuildAt board.Pos
}

const (
	INVALID_MOVE_MSG  = "Move is invalid for the given board state"
	INVALID_BUILD_MSG = "Build is invalid for the given board state"
)

var (
	INVALID_MOVE_ERR  = errors.New(INVALID_MOVE_MSG)
	INVALID_BUILD_ERR = errors.New(INVALID_BUILD_MSG)
)

func (t Turn) ValidMove() bool {
	validWID := board.ValidWID(t.WID)
	moveInBounds := t.MoveTo.InBounds()

	return validWID && moveInBounds
}

func (t Turn) ValidBuild() bool {
	validWID := board.ValidWID(t.WID)
	buildInBounds := t.BuildAt.InBounds()

	return validWID && buildInBounds
}

func (t Turn) Move(player string, b board.IBoard) (board.IBoard, board.IWorker, error) {
	if !t.ValidMove() {
		return b, nil, INVALID_MOVE_ERR
	}

	worker, err := b.FindWorker(player, t.WID)
	if err != nil {
		return b, nil, err
	}

	if !rules.CheckMove(b, worker.Pos(), t.MoveTo) {
		return b, worker, rules.RULE_BROKEN_ERR
	}

	b, err = b.Move(player, t.WID, t.MoveTo)
	worker = b.WorkersFor(player)[t.WID]
	return b, worker, err
}

func (t Turn) Build(player string, b board.IBoard) (board.IBoard, board.IWorker, error) {
	if !t.ValidBuild() {
		return b, nil, INVALID_BUILD_ERR
	}

	worker, err := b.FindWorker(player, t.WID)
	if err != nil {
		return b, nil, err
	}

	if !rules.CheckBuild(b, worker.Pos(), t.BuildAt) {
		return b, worker, rules.RULE_BROKEN_ERR
	}

	b, err = b.AddFloor(t.BuildAt)
	return b, worker, err
}

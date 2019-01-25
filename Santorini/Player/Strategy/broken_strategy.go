package strategy

import (
	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	iplayer "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
)

func BrokenStrategy(player string) IStrategy {
	return NewStrategy(player, FarPlacement, BrokenTurn)
}

func BrokenTurn(b board.IBoard, player, opponent string) (iplayer.Turn, error) {
	return iplayer.Turn{0, board.Pos{0, 0}, board.Pos{0, 0}}, nil
}

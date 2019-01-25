package strategy

import (
	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	iplayer "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
)

func InfiniteTurnStrategy(player string) IStrategy {
	return NewStrategy(player, FarPlacement, InfiniteTurn)
}

func InfinitePlaceStrategy(player string) IStrategy {
	return NewStrategy(player, InfinitePlacement, StayAliveTurn)
}

func InfinitePlacement(b board.IBoard, player, opponent string) (board.Pos, error) {
	for {
	}
	return FarPlacement(b, player, opponent)
}

func InfiniteTurn(b board.IBoard, player, opponent string) (iplayer.Turn, error) {
	for {
	}
	return WinningTurn(b, player, opponent, 3)
}

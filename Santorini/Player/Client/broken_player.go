package client

import (
	common "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	strategy "github.com/CS4500-F18/dare-rebr/Santorini/Player/Strategy"
)

//Creates a new player that doesn't abide by the rules
func BrokenPlayer(name string) common.IPlayer {
	return player{
		name:     name,
		strategy: strategy.BrokenStrategy(name),
	}
}

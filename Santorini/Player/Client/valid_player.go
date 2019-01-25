package client

import (
	common "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	strategy "github.com/CS4500-F18/dare-rebr/Santorini/Player/Strategy"
)

//Creates a new player object that abides by rules
func ValidPlayer(name string) common.IPlayer {
	return player{
		name:     name,
		strategy: strategy.ValidStrategy(name),
	}
}

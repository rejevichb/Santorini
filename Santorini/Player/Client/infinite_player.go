package client

import (
	common "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	strategy "github.com/CS4500-F18/dare-rebr/Santorini/Player/Strategy"
)

//Creates a new player that loops forever on turns
func InfiniteTurnPlayer(name string) common.IPlayer {
	return player{
		name:     name,
		strategy: strategy.InfiniteTurnStrategy(name),
	}
}

//Creates a new player that loops forever on placement
func InfinitePlacementPlayer(name string) common.IPlayer {
	return player{
		name:     name,
		strategy: strategy.InfinitePlaceStrategy(name),
	}
}

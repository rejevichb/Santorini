package main

import (
	player "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	client "github.com/CS4500-F18/dare-rebr/Santorini/Player/Client"
)

func Player(name string) player.IPlayer {
	return client.InfiniteTurnPlayer(name)
}

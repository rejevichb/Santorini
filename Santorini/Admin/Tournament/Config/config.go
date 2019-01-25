package config

import (
	sandbox "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Sandbox"
	iobs "github.com/CS4500-F18/dare-rebr/Santorini/Observer"
)

// TourneyConfiguration for a Tournament
type TournamentConfig interface {
	GenerateComponents() ([]sandbox.WrappedPlayer, []iobs.IObserver)
}

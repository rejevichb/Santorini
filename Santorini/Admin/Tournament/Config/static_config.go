package config

import (
	"encoding/json"
	"fmt"
	"os"

	sandbox "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Sandbox"
	iplayer "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	obs "github.com/CS4500-F18/dare-rebr/Santorini/Observer"
	client "github.com/CS4500-F18/dare-rebr/Santorini/Player/Client"
	remote "github.com/CS4500-F18/dare-rebr/Santorini/Remote/Relay"
)

// Constants for KIND
const (
	VALID    = "good"
	BROKEN   = "breaker"
	INFINITE = "infinite"
)

// TourneyConfiguration for a Tournament
type StaticConfig struct {
	Players   []StaticPlayer   `json:"players"`
	Observers []StaticObserver `json:"observers"`
	IP        string           `json:"ip"`
	Port      int              `json:"port"`
}

// Create Tournament-usable pieces from a TourneyConfiguration
func (c StaticConfig) GenerateComponents() ([]sandbox.WrappedPlayer, []obs.IObserver) {
	players := make([]sandbox.WrappedPlayer, 0)
	for _, p := range c.Players {
		players = append(players, sandbox.NewTimeoutPlayer(sandbox.TIMEOUT_DEFAULT, c.playerFromSpec(p)))

	}

	observers := make([]obs.IObserver, 0)
	for _, o := range c.Observers {
		observers = append(observers, c.observerFromSpec(o))
	}

	return players, observers
}

// ClientRelays returns each Player in this config in its own remote relay,
// alongside specified observers
func (c StaticConfig) ClientRelays() ([]remote.IRelay, []obs.IObserver) {
	relays := make([]remote.IRelay, 0)
	for _, p := range c.Players {
		player := c.playerFromSpec(p)
		relays = append(relays, remote.NewPlayerRelay(player))
	}

	observers := make([]obs.IObserver, 0)
	for _, o := range c.Observers {
		observers = append(observers, c.observerFromSpec(o))
	}

	return relays, observers
}

// Return a Player from the given Player JSON
func (c StaticConfig) playerFromSpec(p StaticPlayer) iplayer.IPlayer {
	switch p.Kind {
	case VALID:
		return client.ValidPlayer(p.Name)

	case BROKEN:
		return client.BrokenPlayer(p.Name)

	case INFINITE:
		return client.InfinitePlacementPlayer(p.Name)
	}

	return nil
}

// Return an Observer from the given Observer JSON
func (c StaticConfig) observerFromSpec(o StaticObserver) obs.IObserver {
	return obs.NewJSONObserver(o.Name, os.Stdout)
}

// Player JSON
type StaticPlayer struct {
	Kind     string
	Name     string
	Location string
}

//Decode JSON bytes into struct
func (p *StaticPlayer) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&p.Kind, &p.Name, &p.Location}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in StaticPlayer: %d != %d", g, e)
	}
	return nil
}

// Observer JSON
type StaticObserver struct {
	Name     string
	Location string
}

//Decode Observer JSON bytes into struct
func (o *StaticObserver) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&o.Name, &o.Location}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in StaticObserver: %d != %d", g, e)
	}
	return nil
}

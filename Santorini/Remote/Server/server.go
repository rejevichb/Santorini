package remote

import (
	"errors"
	"fmt"

	sandbox "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Sandbox"
	tournament "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Tournament"
	config "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Tournament/Config"
	result "github.com/CS4500-F18/dare-rebr/Santorini/Common/Tournament"
)

// An integer (0 or 1) representing a boolean
type IntBool bool

// Convert a JSON integer into the boolean true or false
func (bit *IntBool) UnmarshalJSON(data []byte) error {
	asString := string(data)
	if asString == "1" || asString == "true" {
		*bit = true
	} else if asString == "0" || asString == "false" {
		*bit = false
	} else {
		return errors.New(fmt.Sprintf("Boolean unmarshal error: invalid input %s", asString))
	}
	return nil
}

// Configuration JSON for setting up a Server
type ServerConfig struct {
	//the minimum amount of players before a server can begin a game
	MinPlayers int `json:"min players"`

	//how long the server should wait until starting (in seconds)
	WaitingFor int `json:"waiting for"`

	//port to host the server on
	Port int `json:"port"`

	//true if repeat, false if no repeat
	Repeat int `json:"repeat"`
}

// An empty structure representing a Server
type server struct{}

// Create a new generic server
func NewServer() server {
	return server{}
}

// Starts a new server from the given configuration, then returns a slice of
// tournament results from the tournaments run
func (serv server) Start(cfg ServerConfig) []result.TournamentResult {
	remoteConfig := config.NewRemoteConfig(cfg.MinPlayers, cfg.Port, cfg.WaitingFor, sandbox.TIMEOUT_DEFAULT)
	results := make([]result.TournamentResult, 0)

	if cfg.Repeat == 1 {
		for {
			manager := tournament.NewManager(3)
			manager.RunWithConfig(remoteConfig)
		}
	} else {
		manager := tournament.NewManager(3)
		result := manager.RunWithConfig(remoteConfig)
		results = append(results, result)
	}

	return results
}

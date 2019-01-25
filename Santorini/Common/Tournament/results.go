package tournament

import (
	"encoding/json"

	rules "github.com/CS4500-F18/dare-rebr/Santorini/Common/Rules"
)

// A MatchResult holds who won and who lost within a Game
type MatchResult struct {
	//The winning user's name
	Winner string `json:"winner"`

	//The losing user's name
	Loser string `json:"loser"`

	//Whether this match result is due to a player breaking a rule
	RuleBroken bool `json:"-"`

	//Each of the individual games' results
	GameResults []rules.GameResult `json:"-"`
}

func NewMatchResult(winner string, loser string, ruleBroken bool, games []rules.GameResult) MatchResult {
	return MatchResult{
		Winner:      winner,
		Loser:       loser,
		RuleBroken:  ruleBroken,
		GameResults: games,
	}
}

type TournamentResult struct {
	// The returned results from each game
	Games []MatchResult

	// The users kicked for breaking rules
	Kicked []string
}

func (t TournamentResult) MarshalJSON() ([]byte, error) {
	results := make([]interface{}, 0)

	for _, g := range t.Games {
		if g.RuleBroken {
			results = append(results, []string{g.Winner, g.Loser, "irregular"})
		} else {
			results = append(results, []string{g.Winner, g.Loser})
		}
	}

	return json.Marshal(results)
}

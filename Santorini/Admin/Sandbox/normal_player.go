package sandbox

import (
	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	iplayer "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	result "github.com/CS4500-F18/dare-rebr/Santorini/Common/Tournament"
)

type NormalPlayer struct {
	player iplayer.IPlayer
}

func NewNormalPlayer(p iplayer.IPlayer) NormalPlayer {
	return NormalPlayer{player: p}
}

//Get the name of this Player
func (t NormalPlayer) Name() (string, error) {
	return t.player.Name(), nil
}

//Set the player's name with a timeout
func (t NormalPlayer) SetName(newName string) error {
	t.player.SetName(newName)
	return nil
}

//Get the location to place your next worker
func (t NormalPlayer) PlaceWorker(b board.IBoard) (board.Pos, error) {
	return t.player.PlaceWorker(b), nil
}

//Get the next turn, including which worker to act on
func (t NormalPlayer) NextTurn(b board.IBoard) (iplayer.Turn, error) {
	return t.player.NextTurn(b), nil
}

//sets the opponent of this player
func (t NormalPlayer) SetOpponent(name string) error {
	t.player.SetOpponent(name)
	return nil
}
//receivs tournament results
func (t NormalPlayer) ReceiveTournamentResults(results []result.MatchResult) error {
	return nil
}

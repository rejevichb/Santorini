package sandbox

import (
	"fmt"
	"time"

	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	iplayer "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	result "github.com/CS4500-F18/dare-rebr/Santorini/Common/Tournament"
)

var (
	TIMEOUT_ERROR = func(method string) error {
		return fmt.Errorf("Player timed out on call to %s", method)
	}
)

// TimeoutPlayer wraps player calls in a timeout
type TimeoutPlayer struct {
	timeout int
	player  iplayer.IPlayer
}

// Timeout, in TIMEOUT_UNITs, alongside the player being wrapped
func NewTimeoutPlayer(timeout int, p iplayer.IPlayer) TimeoutPlayer {
	return TimeoutPlayer{timeout: timeout, player: p}
}

//Get the name of this Player
func (t TimeoutPlayer) Name() (string, error) {
	timeout := time.Duration(t.timeout) * TIMEOUT_UNIT
	ch := make(chan string)
	quit := make(chan bool, 1)

	go nameTimeout(t.player, ch, quit)

	select {
	case res := <-ch:
		return res, nil
	case <-time.After(timeout):
		quit <- true
		return "", TIMEOUT_ERROR("Name()")
	}
}

//Set the player's name with a timeout
func (t TimeoutPlayer) SetName(newName string) error {
	go t.player.SetName(newName)
	return nil
}

//Get the location to place your next worker
func (t TimeoutPlayer) PlaceWorker(b board.IBoard) (board.Pos, error) {
	timeout := time.Duration(t.timeout) * time.Millisecond
	ch := make(chan board.Pos)
	quit := make(chan bool, 1)

	go placeTimeout(t.player, b, ch, quit)

	select {
	case res := <-ch:
		return res, nil
	case <-time.After(timeout):
		quit <- true
		return board.Pos{X: -1, Y: -1}, TIMEOUT_ERROR("PlaceWorker()")
	}
}

//Get the next turn, including which worker to act on
func (t TimeoutPlayer) NextTurn(b board.IBoard) (iplayer.Turn, error) {
	timeout := time.Duration(t.timeout) * TIMEOUT_UNIT
	ch := make(chan iplayer.Turn)
	quit := make(chan bool, 1)

	go turnTimeout(t.player, b, ch, quit)

	select {
	case res := <-ch:
		return res, nil
	case <-time.After(timeout):
		quit <- true
		return iplayer.Turn{
			WID:     -1,
			MoveTo:  board.Pos{X: -1, Y: -1},
			BuildAt: board.Pos{X: -1, Y: -1},
		}, TIMEOUT_ERROR("NextTurn()")
	}
}

// Receive an opponent we are playing against
func (t TimeoutPlayer) SetOpponent(name string) error {
	go t.player.SetOpponent(name)
	return nil
}

// Receive the results of a finished Tournament
func (t TimeoutPlayer) ReceiveTournamentResult(result result.TournamentResult) error {
	return nil
}

// ######### Go sucks ########## //
// Each of these does the same thing:
//            - Request some data from a player, given a Board
//            - Receive data, or end the request after a timeout

func nameTimeout(p iplayer.IPlayer, ch chan string, q chan bool) {
	for {
		select {
		case <-q:
			return
		default:
			ch <- p.Name()
			return
		}
	}
}

func placeTimeout(p iplayer.IPlayer, b board.IBoard, ch chan board.Pos, q chan bool) {
	for {
		select {
		case <-q:
			return
		default:
			ch <- p.PlaceWorker(b)
			return
		}
	}
}

func turnTimeout(p iplayer.IPlayer, b board.IBoard, ch chan iplayer.Turn, q chan bool) {
	for {
		select {
		case <-q:
			return
		default:
			ch <- p.NextTurn(b)
			return
		}
	}
}

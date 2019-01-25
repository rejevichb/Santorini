package referee

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/CS4500-F18/dare-rebr/Santorini/Common/Rules"
	lib "github.com/CS4500-F18/dare-rebr/Santorini/Lib"

	sandbox "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Sandbox"
	"github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	iplayer "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	obs "github.com/CS4500-F18/dare-rebr/Santorini/Observer"
	"github.com/CS4500-F18/dare-rebr/Santorini/Player/Client"
)

const PLAYER_1 = "uno"
const PLAYER_2 = "dos"

// For testing purposes
func SetupBoard(p_1, p_2, p_3, p_4 board.Pos) board.IBoard {
	testBoard := board.IBoard(board.BaseBoard())
	testBoard, _ = testBoard.PlaceWorker(p_1, PLAYER_1)
	testBoard, _ = testBoard.PlaceWorker(p_2, PLAYER_1)
	testBoard, _ = testBoard.PlaceWorker(p_3, PLAYER_2)
	testBoard, _ = testBoard.PlaceWorker(p_4, PLAYER_2)
	return testBoard
}

// Create a new referee with a 3 second timeout
func newRef(name1 string, p1 iplayer.IPlayer, name2 string, p2 iplayer.IPlayer) *referee {
	return &referee{
		names: [board.PlayerCount]string{name1, name2},
		players: [board.PlayerCount]sandbox.WrappedPlayer{
			sandbox.NewTimeoutPlayer(3000, p1),
			sandbox.NewTimeoutPlayer(3000, p2),
		},
		observers: []obs.IObserver{},
	}
}

//Helper to create a referee with valid Players 1 and 2
func getReferee() *referee {
	p1 := client.ValidPlayer(PLAYER_1)
	p2 := client.ValidPlayer(PLAYER_2)
	r := newRef(PLAYER_1, p1, PLAYER_2, p2)
	return r
}

//Helper to create a referee with valid player 1 and broken player 2
func getRiggedReferee() *referee {
	p1 := client.ValidPlayer(PLAYER_1)
	p2 := client.BrokenPlayer(PLAYER_2)
	r := newRef(PLAYER_1, p1, PLAYER_2, p2)
	return r
}

//Test that attaching an observer
func TestReferee_AttachObserver(t *testing.T) {
	ref := getReferee()
	observer := obs.NewJSONObserver("observer1", os.Stdout)

	ref.AttachObserver(observer)

	if len(ref.observers) != 1 {
		t.Fail()
	}
}

//Test that removing an observer should modify the referee's observers
func TestReferee_DetachObserver(t *testing.T) {
	ref := getReferee()
	observer := obs.NewJSONObserver("observer1", os.Stdout)

	ref.AttachObserver(observer)

	if len(ref.observers) != 1 {
		t.Fail()
	}

	ref.DetachObserver(observer)

	if len(ref.observers) != 0 {
		t.Fail()
	}
}

//Test that notifying an observer writes to that observer's target
//NOTE this test relies on the Observer's behavior, should be stubbed
func TestNotifyObserver(t *testing.T) {
	var buf strings.Builder

	ref := getReferee()
	observer := obs.NewJSONObserver("observer1", &buf)

	ref.AttachObserver(observer)

	if len(ref.observers) != 1 {
		t.Error("Referee observers not incremented")
	}

	b := SetupBoard(board.Pos{X: 3, Y: 3}, board.Pos{X: 1, Y: 1}, board.Pos{X: 5, Y: 5}, board.Pos{X: 4, Y: 4})

	ref.NotifyAll(b)

	boardBytes, e := json.Marshal(b)
	if e != nil {
		t.Error("Failed to marshal board into JSON")
	}

	if string(boardBytes) != lib.StripSpaces(buf.String()) {
		t.Error("Board bytes not equal to expected string")
	}
}

//Test that the result creator accesses from the referee's names
func TestReferee_result(t *testing.T) {
	ref := getReferee()

	gRes := ref.result(0, 1, "test", true)

	if gRes.BrokenRule != true {
		t.Fail()
	}

	if gRes.Loser != PLAYER_2 || gRes.Winner != PLAYER_1 {
		t.Fail()
	}

	if gRes.Reason != "test" {
		t.Fail()
	}
}

// Test a single game between a valid and broken player ends in
// a broken rule from the broken player
func TestReferee_playSingleGame_ValidAndBrokenWorker(t *testing.T) {
	r := getRiggedReferee()

	gResult := r.playSingleGame(board.BaseBoard())

	if gResult.Winner != PLAYER_1 || gResult.Loser != PLAYER_2 {
		t.Fail()
	}

	if !gResult.BrokenRule {
		t.Fail()
	}

	if gResult.Reason != rules.RULE_BROKEN_MSG {
		t.Fail()
	}
}

//A game between a broken and a valid player should
//have one game that has a broken rule
func TestReferee_Play(t *testing.T) {
	r := getRiggedReferee()

	results := r.Play()

	if len(results) != 1 {
		t.Fail()
	}
	gResult := results[0]

	if gResult.Winner != PLAYER_1 || gResult.Loser != PLAYER_2 {
		t.Fail()
	}

	if !gResult.BrokenRule {
		t.Fail()
	}

	if gResult.Reason != rules.RULE_BROKEN_MSG {
		t.Fail()
	}
}

// Multiple games between a valid player and a broken player should return one
// game, with the broken player as the loser, and having a rule broken
func TestReferee_BestOf_OneRuleBreaker(t *testing.T) {
	ref := getRiggedReferee()

	results := ref.BestOf(10)

	for _, r := range results {
		if r.Reason != rules.RULE_BROKEN_MSG || !r.BrokenRule {
			t.Fail()
		}
		if r.Loser != PLAYER_2 || r.Winner != PLAYER_1 {
			t.Fail()
		}
	}
}

// A game between all valid players should not have a broken rule
func TestReferee_BestOf_AllValidPlayers(t *testing.T) {
	ref := getReferee()

	results := ref.BestOf(3)

	for _, r := range results {
		if r.BrokenRule {
			t.Fail()
		}
	}
}

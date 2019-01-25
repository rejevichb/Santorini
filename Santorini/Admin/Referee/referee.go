package referee

import (
	sandbox "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Sandbox"
	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	output "github.com/CS4500-F18/dare-rebr/Santorini/Common/JSON"
	rules "github.com/CS4500-F18/dare-rebr/Santorini/Common/Rules"
	obs "github.com/CS4500-F18/dare-rebr/Santorini/Observer"
)

// The referee carries a game to its completion given that the strategies are
// instantiated previously, and that each strategy has a player name associated
// with it
type IReferee interface {
	Play() []rules.GameResult
	BestOf(count int) []rules.GameResult
	AttachObserver(obs obs.IObserver)
	DetachObserver(obs obs.IObserver)
}

//The referee maintains the names of its players, the players it is
//accessing, and the observers on its games
type referee struct {
	//The names for players 1 and 2
	names [board.PlayerCount]string

	//players 1 and 2 (names above)
	players [board.PlayerCount]sandbox.WrappedPlayer

	//Observers on this game
	observers []obs.IObserver
}

// instantiates a new referee that can run a game, with the given names and players
func NewReferee(name1 string, p1 sandbox.WrappedPlayer, name2 string, p2 sandbox.WrappedPlayer) IReferee {
	return &referee{
		names:     [board.PlayerCount]string{name1, name2},
		players:   [board.PlayerCount]sandbox.WrappedPlayer{p1, p2},
		observers: []obs.IObserver{},
	}
}

// Runs a game and returns the map of player names to wins, and the result of
// the last game in the series
// The Referee will end the game and declare a winner if any call to a Strategy
// object takes longer than a given timeout time
func (r *referee) Play() []rules.GameResult {
	return r.BestOf(1)
}

// Runs count games and returns the map of player names to wins, and the result
// of the final game in the series
// The Referee will end the game and award all wins to the other player
// if any call to a strategy object takes longer than a given timeout time.
// NOTE: `games` is a natural odd number
func (r *referee) BestOf(games int) []rules.GameResult {
	results := make([]rules.GameResult, 0)

	wins := make(map[string]int)
	wins[r.names[0]] = 0
	wins[r.names[1]] = 0

	for turnPlayer, player := range r.players {
		otherPlayer := opponent(turnPlayer)
		err := player.SetOpponent(r.names[otherPlayer])
		if err != nil {
			return []rules.GameResult{{
				Winner:     r.names[otherPlayer],
				Loser:      r.names[turnPlayer],
				Reason:     rules.RULE_BROKEN_MSG,
				BrokenRule: true,
			}}
		}
	}

	for i := 0; i < games; i++ {
		board := board.BaseBoard()
		result := r.playSingleGame(board)

		results = append(results, result)
		wins[result.Winner]++
		if result.BrokenRule {
			break
		}
		if wins[result.Winner] > games/2 { // strictly greater catches the `1` case
			break
		}
	}

	return results
}

// Runs a game and returns the winner.
// The Referee will end the game and declare a winner if any call to a Strategy
// object takes longer than a given timeout time
func (r *referee) playSingleGame(b board.IBoard) rules.GameResult {
	r.NotifyAll(b)

	// Phase 1 (Placing Workers):
	b, endGame := r.startWorkerPlacement(b)
	if endGame.Winner != "" {
		r.NotifyAll(b)
		r.NotifyAll(endGame)
		return endGame
	}

	// Phase 2: (Moving and Building)
	b, result := r.handleGameTurns(b)
	r.NotifyAll(b)
	r.NotifyAll(result)
	return result
}

// Runs the loop to receive worker placements from each player
func (r *referee) startWorkerPlacement(b board.IBoard) (board.IBoard, rules.GameResult) {
	for wIdx := 0; wIdx < board.WorkerCount; wIdx++ {
		for turnPlayer := 0; turnPlayer < board.PlayerCount; turnPlayer++ {
			otherPlayer := (turnPlayer + 1) % 2

			breakResult := rules.GameResult{
				Winner:     r.names[otherPlayer],
				Loser:      r.names[turnPlayer],
				Reason:     rules.RULE_BROKEN_MSG,
				BrokenRule: true,
			}

			workerLocation, err := r.players[turnPlayer].PlaceWorker(b)
			if err != nil {
				return nil, breakResult
			}

			if !rules.CheckPlaceWorker(b, workerLocation) {
				return nil, breakResult
			} else {
				if b, err = b.PlaceWorker(workerLocation, r.names[turnPlayer]); err != nil {
					return nil, breakResult
				} else {
					r.NotifyAll(b)
				}
			}
		}
	}

	return b, rules.GameResult{}
}

// Play through a game, taking turns from each player alternating until one player
// has won the game/lost the game, then return the result.
func (r *referee) handleGameTurns(b board.IBoard) (board.IBoard, rules.GameResult) {
	turnPlayer := 0
	otherPlayer := (turnPlayer + 1) % 2

	//each iteration of this for loop represents a "turn" and this loop runs until the game has been won,
	//at which point a GameResult is returned.
	for {
		//check to see if the player's whose turn it is has already lost the game.
		if rules.CheckLossPreMove(b, r.names[turnPlayer]) {
			return b, r.result(otherPlayer, turnPlayer, rules.CANNOT_MOVE_MSG, false)
		}

		//returns a full turn, (build, move) for a the turnPlayer. If the time limit is exceeded,
		//the turn is skipped
		turn, err := r.players[turnPlayer].NextTurn(b)
		if err != nil {
			break
		}

		worker, err := b.FindWorker(r.names[turnPlayer], turn.WID)
		if err != nil {
			break
		}
		moveDir := output.DirectionFrom2Pos(worker.Pos(), turn.MoveTo)

		b, worker, err = turn.Move(r.names[turnPlayer], b)
		if err != nil {
			break
		}

		//if the game has been won, return from this method with the board and GameResult
		if won := rules.CheckWinPostMove(b, r.names[turnPlayer]); won {
			message := output.MoveJSON{worker.Name(), moveDir}
			r.NotifyAll(message)
			return b, r.result(turnPlayer, otherPlayer, rules.WINNING_MOVE_MSG, false)
		}

		buildDir := output.DirectionFrom2Pos(worker.Pos(), turn.BuildAt)
		b, worker, err = turn.Build(r.names[turnPlayer], b)
		if err != nil {
			break
		}

		obsTurn := output.MoveBuildJSON{worker.Name(), moveDir, buildDir}
		r.NotifyAll(obsTurn)
		r.NotifyAll(b)

		//switch the active player (player whose turn it is) and the waiting (non-turn) player
		turnPlayer = (turnPlayer + 1) % 2
		otherPlayer = (turnPlayer + 1) % 2
	}

	return b, r.result(otherPlayer, turnPlayer, rules.RULE_BROKEN_MSG, true)
}

// Create a game result from a winner and loser's idx, a game end reason, and
// whether or not a rule was broken
func (r referee) result(winner, loser int, reason string, broken bool) rules.GameResult {
	return rules.GameResult{
		Winner:     r.names[winner],
		Loser:      r.names[loser],
		Reason:     reason,
		BrokenRule: broken,
	}
}

//Get the opponent from a given player's index
func opponent(pIdx int) int {
	return (pIdx + 1) % board.PlayerCount
}

/*########## OBSERVER HANDLING ##########*/

// Attach an observer to the game(s) being refereed
func (r *referee) AttachObserver(o obs.IObserver) {
	r.observers = append(r.observers, o)
}

// Detach an observer to the game(s) being refereed
func (r *referee) DetachObserver(obs obs.IObserver) {
	for idx, existing := range r.observers {
		if existing == obs {
			r.observers = append(r.observers[:idx], r.observers[idx+1:]...)
		}
	}
}

// Inform all observers to the game(s) of a board change or player turn
func (r referee) NotifyAll(update interface{}) {
	for _, obs := range r.observers {
		NotifyObserver(obs, update)
	}
}

// Inform an observers of a board change or player turn
func NotifyObserver(observer obs.IObserver, update interface{}) {
	switch parsed := update.(type) {
	case board.IBoard:
		observer.ReceiveBoard(parsed)

	case output.MoveBuildJSON:
		observer.ReceiveTurn(parsed)

	case output.MoveJSON:
		observer.ReceiveWinningMove(parsed)

	case rules.GameResult:
		observer.ReceiveEndgame(parsed)
	}
}

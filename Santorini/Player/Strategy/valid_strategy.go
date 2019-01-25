package strategy

import (
	"errors"

	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	iplayer "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	rules "github.com/CS4500-F18/dare-rebr/Santorini/Common/Rules"
	// "fmt"
)

const (
	CANNOT_SURVIVE = "Unable to stay alive from given board state and depth"
	CANNOT_WIN     = "Unable to win from given board state and depth"
)

func ValidStrategy(player string) IStrategy {
	return NewStrategy(player, FarPlacement, StayAliveTurn)
}

func DiagonalPlacement(b board.IBoard, player string) (board.Pos, error) {
	for index := 0; index < board.NormalBoardSize; index++ {
		diagPos := board.Pos{X: index, Y: index}
		worker := b.WorkerAt(diagPos)
		if worker == nil {
			return diagPos, nil
		}
	}
	return board.Pos{-1, -1}, errors.New("No free diagonal spaces")
}

func FarPlacement(b board.IBoard, player, opponent string) (board.Pos, error) {
	enemyWorkers := b.WorkersFor(opponent)

	var bestPos board.Pos
	bestDist := 0
	width, height := b.Dimensions()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pos := board.Pos{X: x, Y: y}
			tempDist := 0
			for _, enemyWorker := range enemyWorkers {
				if enemyWorker != nil {
					tempDist += board.PosDistance(pos, enemyWorker.Pos())
				}
			}
			worker := b.WorkerAt(pos)
			if worker == nil && tempDist > bestDist {
				bestDist = tempDist
				bestPos = pos
			}
		}
	}
	return bestPos, nil
}

func validTurns(b board.IBoard, player string) []iplayer.Turn {
	var turns []iplayer.Turn

	workers := b.WorkersFor(player)
	for _, worker := range workers {
		workerPos := worker.Pos()

		for _, moveTarget := range workerPos.Neighbors() {
			if !rules.CheckMove(b, workerPos, moveTarget) {
				continue
			}

			boardPostMove, err := b.Move(player, worker.ID(), moveTarget)
			if err != nil {
				continue
			}

			workerPostMove := worker.Move(moveTarget)
			postMovePos := workerPostMove.Pos()

			for _, buildTarget := range moveTarget.Neighbors() {
				if !rules.CheckBuild(boardPostMove, postMovePos, buildTarget) {
					continue
				}

				newTurn := iplayer.Turn{
					WID:     worker.ID(),
					MoveTo:  moveTarget,
					BuildAt: buildTarget,
				}

				turns = append(turns, newTurn)
			}
		}
	}
	return turns
}

func StayAliveTurn(b board.IBoard, player, opponent string) (iplayer.Turn, error) {
	return SurvivingTurn(b, player, opponent, 3)
}

func WinnableTurn(b board.IBoard, player, opponent string) (iplayer.Turn, error) {
	return WinningTurn(b, player, opponent, 3)
}

//Return a Turn that will keep you alive for depth turns
func SurvivingTurn(b board.IBoard, player, opponent string, depth int) (iplayer.Turn, error) {
	lastTurn := iplayer.Turn{WID: -1, MoveTo: board.Pos{X: -1, Y: -1}, BuildAt: board.Pos{X: -1, Y: -1}}
	workers := b.WorkersFor(player)

	for _, turn := range validTurns(b, player) {
		lastTurn = turn

		if depth > 0 {
			b, _ = b.Move(player, workers[turn.WID].ID(), turn.MoveTo)
			b, _ = b.AddFloor(turn.BuildAt)

			_, losableTurnErr := WinningTurn(b, opponent, player, depth-1)
			if losableTurnErr == nil {
				continue
			}
		}
		return turn, nil
	}
	return lastTurn, errors.New(CANNOT_SURVIVE)
}

//Return a Turn that wins you the Game, or an error if no such Turn exists
func WinningTurn(b board.IBoard, player, opponent string, depth int) (iplayer.Turn, error) {
	lastTurn := iplayer.Turn{WID: -1, MoveTo: board.Pos{X: -1, Y: -1}, BuildAt: board.Pos{X: -1, Y: -1}}
	workers := b.WorkersFor(player)

	for _, turn := range validTurns(b, player) {
		lastTurn = turn

		if depth > 0 {
			b, _ = b.Move(player, workers[turn.WID].ID(), turn.MoveTo)
			b, _ = b.AddFloor(turn.BuildAt)

			win := rules.CheckWinPostMove(b, player)

			if win {
				return turn, nil
			}

			_, otherSurvivableErr := SurvivingTurn(b, opponent, player, depth-1)
			if otherSurvivableErr == nil {
				continue
			}
		}
		return turn, nil
	}
	return lastTurn, errors.New(CANNOT_WIN)
}

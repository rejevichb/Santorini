package rules

import (
	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
)

//The purpose of the rules package is to orchestrate checks for all logical (non-
//physical [I.E. workers cannot go on rounded 4th floor building rooftops])
//rules about the game. This will include managing which player's turn it is,
//managing what moves during a turn are valid or not,
//and relaying this information to the player, and determining any ending conditions
//for the game.

const WinningWorkerHeight = 3

//A MoveRule is a predicate style function that takes ITiles and returns true
//if the given move would be valid, and false if not.
//Simply checks ahead of time if the turn would be valid.
type moveRule func(b board.IBoard, moveFrom, moveTo board.ITile) bool

// Check that a Move location is on a Board
func moveInBounds(b board.IBoard, moveFrom, moveTo board.ITile) bool {
	return moveTo.Pos().InBounds()
}

// Checks that we can only move to a neighboring ITile
func adjacencyMoveRule(b board.IBoard, moveFrom, moveTo board.ITile) bool {
	return moveFrom.IsNeighbor(moveTo)
}

//Checks that we only go up one floor during our move
func onlyUp1FloorRule(b board.IBoard, moveFrom, moveTo board.ITile) bool {
	return (moveFrom.FloorCount() + 1) >= moveTo.FloorCount()
}

// Returns whether the ITile that is being moved to is vacant
func vacantTileMoveRule(b board.IBoard, moveFrom, moveTo board.ITile) bool {
	return b.WorkerAt(moveTo.Pos()) == nil
}

//A BuildRule is a predicate style function that takes a turn and returns true if the given build would be valid, and false if not.
//DOES NOT ALTER THE STATE OF THE BOARD, WORKER, OR PLAYER.
//Simply checks ahead of time if the turn would be valid.
type buildRule func(b board.IBoard, workerTile, buildAt board.ITile) bool

// Check that a Build location is on a Board
func buildInBounds(b board.IBoard, workerTile, buildAt board.ITile) bool {
	return buildAt.Pos().InBounds()
}

// Checks that we can only build on a neighboring ITile
func adjacencyBuildRule(b board.IBoard, workerTile, buildAt board.ITile) bool {
	return workerTile.IsNeighbor(buildAt)
}

func heightBuildRule(b board.IBoard, workerTile, buildAt board.ITile) bool {
	return buildAt.FloorCount() < board.MaxBuildingHeight
}

//Returns whether the ITile that is being built on is vacant
func vacantTileBuildRule(b board.IBoard, workerTile, buildAt board.ITile) bool {
	return b.WorkerAt(buildAt.Pos()) == nil
}

//A PlaceRule is a predicate style function that takes a PlaceWorker and returns true if it would be valid, and false if not.
//DOES NOT ALTER THE STATE OF THE BOARD, WORKER, OR PLAYER
//Simply checks ahead of time if the place worker would be valid.
type placeRule func(b board.IBoard, placeAt board.ITile) bool

// Check that a Build location is on a Board
func placeInBounds(b board.IBoard, placeAt board.ITile) bool {
	return placeAt.Pos().InBounds()
}

//Returns whether the ITile that is being placed on is vacant
func vacantTilePlaceRule(b board.IBoard, placeAt board.ITile) bool {
	return b.WorkerAt(placeAt.Pos()) == nil
}

// Win conditions check board state from a specific player's perspective
type gameOverCondition func(b board.IBoard, playerName string) bool

// Returns whether the given player has a worker on the goal floor
// Returns false if the player is not playing on the given board

func goalFloorReachedCondition(b board.IBoard, player string) bool {
	workers := b.WorkersFor(player)
	for _, w := range workers {
		if w == nil {
			continue
		}

		tile, err := b.TileAt(w.Pos())
		if err != nil {
			continue
		}
		if tile.FloorCount() == WinningWorkerHeight {
			return true
		}
	}
	return false
}

// Returns whether the given player has a move available
// Returns false if the player is not playing on the given board
func furtherMoveImpossibleCondition(b board.IBoard, player string) bool {
	workers := b.WorkersFor(player)

	for _, w := range workers {
		if w == nil {
			continue
		}

		workerPos := w.Pos()
		workerTile, err := b.TileAt(workerPos)
		if err != nil {
			continue
		}

		for _, neighborPos := range workerPos.Neighbors() {
			targetTile, err := b.TileAt(neighborPos)
			if err != nil {
				return true
			}

			validMove := true
			for _, rule := range moveRules {
				validMove = validMove && rule(b, workerTile, targetTile)
			}

			if validMove {
				return false
			}
		}
	}
	return true
}

// Rules for vanilla Santorini moves
var moveRules = []moveRule{moveInBounds, adjacencyMoveRule, onlyUp1FloorRule, vacantTileMoveRule}

// Rules for vanilla Santorini builds
var buildRules = []buildRule{buildInBounds, adjacencyBuildRule, heightBuildRule, vacantTileBuildRule}

// Rules for vanilla Santorini player places
var placeRules = []placeRule{placeInBounds, vacantTilePlaceRule}

// The conditions that would result in a game loss before a player can move
var lossConditions = []gameOverCondition{furtherMoveImpossibleCondition}

// The conditions that would result in a game win after a player has moved
var winConditions = []gameOverCondition{goalFloorReachedCondition}

// Checks that all of the rules for worker placement pass
func CheckPlaceWorker(b board.IBoard, targetPos board.Pos) bool {
	targetTile, err := b.TileAt(targetPos)
	if err != nil {
		return false
	}

	validPlace := true
	for _, rule := range placeRules {
		validPlace = validPlace && rule(b, targetTile)
	}
	return validPlace
}

// Checks that all of the rules for building on cells pass
func CheckBuild(b board.IBoard, workerPos board.Pos, targetPos board.Pos) bool {
	workerTile, err := b.TileAt(workerPos)
	if err != nil {
		return false
	}

	targetTile, err := b.TileAt(targetPos)
	if err != nil {
		return false
	}

	validBuild := true
	for _, rule := range buildRules {
		// if !rule(build) {
		// 	fmt.Printf("BUILD RULE BROKEN: %v\n", idx)
		// }
		validBuild = validBuild && rule(b, workerTile, targetTile)
	}
	return validBuild
}

// Checks that all of the rules for movement of a worker pass
func CheckMove(b board.IBoard, workerPos, targetPos board.Pos) bool {
	workerTile, err := b.TileAt(workerPos)
	if err != nil {
		return false
	}

	targetTile, err := b.TileAt(targetPos)
	if err != nil {
		return false
	}

	validMove := true
	for _, rule := range moveRules {
		// if !rule(build) {
		// 	fmt.Printf("BUILD RULE BROKEN: %v\n", idx)
		// }
		validMove = validMove && rule(b, workerTile, targetTile)
	}
	return validMove
}

// Checks if a player has lost the game before a move
// NOTE Returns false if the player is not playing on the given board
func CheckLossPreMove(b board.IBoard, player string) bool {
	for _, condition := range lossConditions {
		if loss := condition(b, player); loss {
			// fmt.Println("Player won PRE move", player, "won because of check:", idx)
			return true
		}
	}
	return false
}

// Checks if a player has won the game after a move
// NOTE Returns false if the player is not playing on the given board
func CheckWinPostMove(b board.IBoard, player string) bool {
	for _, condition := range winConditions {
		if win := condition(b, player); win {
			// fmt.Println("Player won POST move", player, "won because of check:", ind)
			return true
		}
	}
	return false
}

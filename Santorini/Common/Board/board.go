package board

import (
	"fmt"
	"sort"

	"github.com/CS4500-F18/dare-rebr/Santorini/Lib"
)

//Purpose: To represent the physical game board of the game Santorini.
//Will just contain the logic regarding the setup and usage of pieces
//on the game board, not the logic regarding how the game is played.

const (
	// MaxBuildingHeight is the maximum number of floors that are allowed on any given ITile
	MaxBuildingHeight = 4

	// NormalBoardSize is the width and height of a normal board in Classic Santorini.
	NormalBoardSize = 6

	// PlayerCount is how many players in a game
	PlayerCount = 2

	// WorkerCount is how many workers per Player in a game
	WorkerCount = 2

	// TotalWorkers is how many workers should be in a typical game
	TotalWorkers = WorkerCount * PlayerCount
)

// Unified Error messages
const (
	PLAYER_NOT_FOUND  = "Player %s not found on given Board"
	WORKER_ID_INVALID = "Worker ID %v is not valid"
	WORKER_NOT_FOUND  = "Worker not found at position (%v, %v)"
	POS_NOT_FOUND     = "No Tile found for Position (%v, %v)"
	MAX_WORKERS       = "Board already reached limit for Worker count"
	MAX_PLAYERS       = "Board already has maximum Players"
	WORKER_NOT_PLACED = "Worker with ID %v has not been placed yet"
	WORKER_NOT_EXIST  = "Worker %s%v not found on the board"
)

//A Board is a struct which has the following receivers
type IBoard interface {
	//Returns the ITile in this Board at the coordinate pair represented by the given target Pos.
	//If the target Pos is outside of this Board, returns an error instead.
	TileAt(target Pos) (ITile, error)
	//Returns all of the workers currently on this Board for the given player.
	//Returns an error if there are no workers for the given player.
	WorkersFor(player string) []IWorker
	Workers() []IWorker
	//Attempts to move the given worker to the given target position. If this move is impossible, returns an error.
	Move(playerName string, wID int, target Pos) (IBoard, error)
	//Attempts to add a floor to the given target position. If this move is impossible (the Worker's AddFloorTo() returns an error), returns an error.
	AddFloor(target Pos) (IBoard, error)
	//Confirms that AcceptWorkerArrival() returns no errors (returning an error if so) for the ITile corresponding to the given Pos,
	//creates a new Worker, places the Worker on that ITile and returns a pointer to the Worker struct.
	//The Worker is added the the Board's collection of Workers.
	PlaceWorker(p Pos, owner string) (IBoard, error)
	//Returns a slice of strings with the names of all of the players using this board
	Players() []string
	//Returns the width and height of the Board
	Dimensions() (int, int)
	//Find the given player's worker, or return an error if it doesn't exist
	FindWorker(player string, workerID int) (IWorker, error)
	//Return the Worker at the given Pos, or an error if there is none
	WorkerAt(pos Pos) IWorker
}

type board struct {
	//A Board has a TileMap, representing the (x, y) set of Tiles
	// (access to Tiles will be in the form of board.tiles[x][y], once you
	// know that x and y are valid accessors)
	// (x, y) represents (column, row)
	tiles TileMap

	//Workers on the Board
	workers WorkerSet
}

func (b board) Dimensions() (int, int) {
	return NormalBoardSize, NormalBoardSize
}

// Return the Tile at the given Pos, or an error if there is none
func (b board) TileAt(target Pos) (ITile, error) {
	if !target.InBounds() {
		return nil, fmt.Errorf(POS_NOT_FOUND, target.X, target.Y)
	}
	return b.tiles[target.X][target.Y], nil
}

//get the workers for the given player
func (b board) WorkersFor(playerName string) []IWorker {
	targetWorkers := make(WorkerSet, 0)

	for _, worker := range b.workers {
		if worker != nil && worker.Owner() == playerName {
			targetWorkers = append(targetWorkers, worker)
		}
	}

	sort.Sort(targetWorkers)
	return targetWorkers
}

//get all workers on this Board
func (b board) Workers() []IWorker {
	targetWorkers := make(WorkerSet, 0)
	for _, worker := range b.workers {
		if worker != nil {
			targetWorkers = append(targetWorkers, worker)
		}
	}

	sort.Sort(targetWorkers)
	return targetWorkers
}

func (b board) Move(playerName string, workerID int, target Pos) (IBoard, error) {
	worker, err := b.FindWorker(playerName, workerID)
	if err != nil {
		return b, err
	}

	// Knowing the Move should work, perform it on the board w/ mutation
	newBoard, err := b.performMove(playerName, worker, target)
	if err != nil {
		return b, err
	}

	return newBoard, nil
}

// Return the Board after adding a floor to the given Pos
// Returns an error if any step fails
func (b board) AddFloor(target Pos) (IBoard, error) {
	targetTile, err := b.TileAt(target)
	if err != nil {
		return b, err
	}

	targetTile, err = targetTile.AddFloor()
	if err != nil {
		return b, err
	}

	newBoard := b.deepCopy()
	newBoard.tiles[targetTile.Pos().X][targetTile.Pos().Y] = targetTile

	return newBoard, nil
}

//Places a worker at the given target position, for the given owner
//returns an error if you are trying to add a player past the max player count
//returns an error if you there is are already the max number of workers
//returns an error if the position is out of bounds
func (b board) PlaceWorker(targetPos Pos, owner string) (IBoard, error) {
	if len(b.Players()) >= PlayerCount && !lib.StringPresent(b.Players(), owner) {
		return b, fmt.Errorf(MAX_PLAYERS)
	}

	newBoard := b.deepCopy()

	err := newBoard.addWorker(targetPos, owner)
	if err != nil {
		return b, err
	}

	return newBoard, nil
}

//Return the Players whose workers are on this Board
func (b board) Players() []string {
	players := make([]string, 0)

	for _, worker := range b.workers {
		if worker != nil && !lib.StringPresent(players, worker.Owner()) {
			players = append(players, worker.Owner())
		}
	}

	return players
}

// Find the given Player's worker with the given ID
func (b board) FindWorker(player string, workerID int) (IWorker, error) {
	if !ValidWID(workerID) {
		return nil, fmt.Errorf(WORKER_ID_INVALID, workerID)
	}

	workers := b.WorkersFor(player)

	for _, worker := range workers {
		if worker != nil && worker.ID() == workerID {
			return worker, nil
		}
	}

	return nil, fmt.Errorf(WORKER_NOT_EXIST, player, workerID)
}

// Find the Worker at the given Pos, or nil if there is none
func (b board) WorkerAt(pos Pos) IWorker {
	for _, w := range b.workers {
		if w != nil && w.Pos() == pos {
			return w
		}
	}
	return nil
}

//Copies the contents of this Board into a new Board
func (b board) deepCopy() board {
	newMap := make(TileMap)
	for x, set := range b.tiles {
		newSet := make(TileSet)
		for y, tile := range set {
			newSet[y] = tile
		}
		newMap[x] = newSet
	}

	newWorkers := make(WorkerSet, len(b.workers))
	for idx, worker := range b.workers {
		newWorkers[idx] = worker
	}

	return board{
		tiles:   newMap,
		workers: newWorkers,
	}
}

//Mutate the board to move the Worker
func (b board) performMove(playerName string, existing IWorker, target Pos) (board, error) {
	// New Board
	newBoard := b.deepCopy()

	// Move worker on New Board
	movedWorker := existing.Move(target)
	newBoard.replaceWorker(playerName, existing.ID(), movedWorker)

	return newBoard, nil
}

//Mutate the board and replace the targetted worker with the given worker
func (b board) replaceWorker(playerName string, workerID int, newWorker IWorker) {
	for idx, w := range b.workers {
		if w != nil && w.Owner() == playerName && w.ID() == workerID {
			b.workers[idx] = newWorker
			return
		}
	}
}

//Add a worker at the given position owned by the given owner
//Returns an error if the Board's worker array is full
func (b board) addWorker(pos Pos, owner string) error {
	existing := b.WorkersFor(owner)
	if len(existing) >= WorkerCount {
		return fmt.Errorf(MAX_WORKERS)
	}

	for idx, worker := range b.workers {
		if worker == nil {
			newWorker := NewWorker(pos, owner, len(existing))
			b.workers[idx] = newWorker
			return nil
		}
	}

	return fmt.Errorf(MAX_WORKERS)
}

func (b board) workerCount() int {
	count := 0
	for _, worker := range b.workers {
		if worker != nil {
			count++
		}
	}
	return count
}

//Constructs a new Board struct with all tiles set to 0,
//and no workers
func BaseBoard() board {
	return newBoard(emptyTileMap(), emptyWorkerSet())
}

//Constructs a new board with the given tile/worker set
func newBoard(tiles TileMap, workers WorkerSet) board {
	return board{
		tiles:   tiles,
		workers: workers,
	}
}

//Constructs a board with a select set of tiles,
//with all else being height 0
func BoardWithTiles(tiles []ITile) board {
	baseMap := emptyTileMap()
	for _, tile := range tiles {
		pos := tile.Pos()
		baseMap[pos.X][pos.Y] = tile
	}

	return board{
		tiles:   baseMap,
		workers: emptyWorkerSet(),
	}
}

//Constructs a new Board with the given workers,
//and all tiles set to height 0
func BoardWithWorkers(workers []IWorker) board {
	return newBoard(emptyTileMap(), workers)
}

//returns an empty worker set
func emptyWorkerSet() WorkerSet {
	return make(WorkerSet, TotalWorkers)
}

//returns a map of tiles with posns for a 6x6 grid, all with heights of 0 and marked unoccupied
func emptyTileMap() TileMap {
	tiles := make(TileMap)
	for x := 0; x < NormalBoardSize; x++ {
		tiles[x] = make(TileSet)
		for y := 0; y < NormalBoardSize; y++ {
			p := Pos{X: x, Y: y}
			tiles[x][y] = NewTile(p)
		}
	}
	return tiles
}

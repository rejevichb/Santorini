package json

import (
	"encoding/json"
	"fmt"

	common "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	iplayer "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
)

var EMPTY = []int{}

// Move and Build Turn with Direction structs
type MoveBuildJSON struct {
	WorkerName string
	MoveDir    Direction
	BuildDir   Direction
}

// Just a Move Turn
type MoveJSON struct {
	WorkerName string
	MoveDir    Direction
}

//Represents a Direction according to the specification system given in our homeworks
type Direction struct {
	EastWest   string
	NorthSouth string
}

// Direction strings
const (
	EAST  = "EAST"
	WEST  = "WEST"
	NORTH = "NORTH"
	SOUTH = "SOUTH"
	PUT   = "PUT"
)

//Convert a Direction to and from JSON bytes
func (d Direction) MarshalJSON() ([]byte, error) {
	result := [2]string{d.EastWest, d.NorthSouth}
	return json.Marshal(result)
}
func (d *Direction) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&d.EastWest, &d.NorthSouth}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Direction: %d != %d", g, e)
	}
	return nil
}

//Get a position from an existing position and a direction
func PosFromDirection(currentPos common.Pos, d Direction) common.Pos {
	newPos := common.Pos{}
	switch d.EastWest {
	case "EAST":
		newPos.X = currentPos.X + 1
	case "WEST":
		newPos.X = currentPos.X - 1
	case "PUT":
		newPos.X = currentPos.X
	}
	switch d.NorthSouth {
	case "NORTH":
		newPos.Y = currentPos.Y - 1
	case "SOUTH":
		newPos.Y = currentPos.Y + 1
	case "PUT":
		newPos.Y = currentPos.Y
	}
	return newPos
}

//Get a Direction from an origin and target position
func DirectionFrom2Pos(p1, p2 common.Pos) Direction {
	d := Direction{}

	colDiff := (p1.X - p2.X)
	rowDiff := (p1.Y - p2.Y)

	switch colDiff {
	case -1:
		d.EastWest = EAST
	case 0:
		d.EastWest = PUT
	case 1:
		d.EastWest = WEST
	}

	switch rowDiff {
	case -1:
		d.NorthSouth = SOUTH
	case 0:
		d.NorthSouth = PUT
	case 1:
		d.NorthSouth = NORTH
	}

	return d
}

// A move and build from a player, with direction strings
type MoveBuildTurn struct {
	//Player name + Worker ID
	WorkerName string

	// Move directions
	MoveEW string
	MoveNS string

	// Build directions
	BuildEW string
	BuildNS string
}

// Convert a move and build to an array of its fields
func (m MoveBuildTurn) MarshalJSON() ([]byte, error) {
	tmp := []interface{}{m.WorkerName, m.MoveEW, m.MoveNS, m.BuildEW, m.BuildNS}
	return json.Marshal(tmp)
}

// Convert JSON bytes to a MoveBuild struct with strings
func (m *MoveBuildTurn) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&m.WorkerName, &m.MoveEW, &m.MoveNS, &m.BuildEW, &m.BuildNS}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("Wrong number of fields in MoveBuildTurn: %d != %d", g, e)
	}
	return nil
}

// Move/build with strings to move/build with Direction structs
func (mb MoveBuildTurn) ToDirection() MoveBuildJSON {
	return MoveBuildJSON{
		mb.WorkerName,
		Direction{mb.MoveEW, mb.MoveNS},
		Direction{mb.BuildEW, mb.BuildNS},
	}
}

// Move build with Direction structs to a Turn, given a Board
// Returns an error if: - the worker's name is invalid
//                      - there is no worker with the given player & ID
func (mb MoveBuildTurn) ToTurn(b common.IBoard) (iplayer.Turn, error) {
	turn := iplayer.Turn{}

	withDir := mb.ToDirection()

	playerName, workerID, err := common.ParseWorkerName(mb.WorkerName)
	if err != nil {
		return turn, err
	}

	turn.WID = workerID

	worker, err := b.FindWorker(playerName, workerID)
	if err != nil {
		return turn, err
	}

	movePos := PosFromDirection(worker.Pos(), withDir.MoveDir)
	turn.MoveTo = movePos

	buildPos := PosFromDirection(movePos, withDir.BuildDir)
	turn.BuildAt = buildPos

	return turn, nil
}

// Convert a Turn, with a Board and the Player whose turn it is, to a move/build
// with string directions
func MoveBuildFromTurn(player string, b common.IBoard, turn iplayer.Turn) MoveBuildTurn {
	worker, _ := b.FindWorker(player, turn.WID)

	moveDir := DirectionFrom2Pos(worker.Pos(), turn.MoveTo)
	buildDir := DirectionFrom2Pos(turn.MoveTo, turn.BuildAt)

	return MoveBuildTurn{
		WorkerName: worker.Name(),
		MoveEW:     moveDir.EastWest,
		MoveNS:     moveDir.NorthSouth,
		BuildEW:    buildDir.EastWest,
		BuildNS:    buildDir.NorthSouth,
	}
}

// A move turn, with direction strings and a worker name (player name + worker ID)
type MoveTurn struct {
	// Player name + Worker ID
	WorkerName string

	// Move directions
	MoveEW string
	MoveNS string
}

// Create a MoveTurn from JSON bytes
func (m *MoveTurn) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&m.WorkerName, &m.MoveEW, &m.MoveNS}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("Wrong number of fields in MoveTurn: %d != %d", g, e)
	}
	return nil
}

// Change the Move's string directions to a Direction struct
func (m MoveTurn) ToDirection() MoveJSON {
	return MoveJSON{
		m.WorkerName,
		Direction{m.MoveEW, m.MoveNS},
	}
}

// Convert a Move to a Turn, given a board
// Returns an error if the worker name is invalid, or the worker doesn't exist
func (mb MoveTurn) ToTurn(b common.IBoard) (iplayer.Turn, error) {
	turn := iplayer.Turn{}

	withDir := mb.ToDirection()

	playerName, workerID, err := common.ParseWorkerName(mb.WorkerName)
	if err != nil {
		return turn, err
	}

	turn.WID = workerID

	worker, err := b.FindWorker(playerName, workerID)
	if err != nil {
		return turn, err
	}

	movePos := PosFromDirection(worker.Pos(), withDir.MoveDir)
	turn.MoveTo = movePos

	turn.BuildAt = common.Pos{X: -1, Y: -1}

	return turn, nil
}

// Create a Move from a Turn, the Board that Turn is acting on, and the Player
// who is enacting that Turn
func MoveFromTurn(player string, b common.IBoard, turn iplayer.Turn) MoveTurn {
	worker, _ := b.FindWorker(player, turn.WID)

	moveDir := DirectionFrom2Pos(worker.Pos(), turn.MoveTo)

	return MoveTurn{
		WorkerName: worker.Name(),
		MoveEW:     moveDir.EastWest,
		MoveNS:     moveDir.NorthSouth,
	}
}

// Renaming a player
type Rename struct {
	Name string
}

// Convert JSON bytes to a struct wrapping a string
func (n Rename) UnmarshalJSON(buf []byte) error {
	var throwaway string
	tmp := []interface{}{&throwaway, &n.Name}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Rename: %d != %d", g, e)
	}
	if tmp[0] != "playing-as" {
		// should have "playing-as", but we ignore it
		// once we confirm presence
		return fmt.Errorf("Invalid first argument to Rename: %s", tmp[0])
	}
	return nil
}

// Convert a renaming to a JSON array, including "playing-as" as the first element
func (n *Rename) MarshalJSON() ([]byte, error) {
	tmp := []string{"playing-as", n.Name}
	return json.Marshal(tmp)
}

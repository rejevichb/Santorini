package json

import (
	"encoding/json"
	"fmt"

	common "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	rules "github.com/CS4500-F18/dare-rebr/Santorini/Common/Rules"
)

const YES = "yes"
const NO = "no"

//Represents any Command following a Board specification
type Command struct {
	Type      string
	Worker    string
	Direction Direction
}

//Unmarshals the JSON Array into an ACTUAL JSON object
func (c *Command) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&c.Type, &c.Worker, &c.Direction}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Command: %d != %d", g, e)
	}
	return nil
}

// Execute the given Command on a Board
func (command Command) Execute(b common.IBoard) (common.IBoard, interface{}) {
	playerName, workerID, err := common.ParseWorkerName(command.Worker)
	if err != nil {
		panic(err)
	}

	worker, err := b.FindWorker(playerName, workerID)
	if err != nil {
		panic(err)
	}
	workerPos := worker.Pos()
	targetPos := PosFromDirection(workerPos, command.Direction)

	switch command.Type {
	case "move":
		if !rules.CheckMove(b, workerPos, targetPos) {
			return b, NO
		}
		post, err := b.Move(playerName, workerID, targetPos)
		if err != nil {
			return post, NO
		}
		return post, EMPTY

	case "build", "+build":
		if !rules.CheckBuild(b, workerPos, targetPos) {
			return b, NO
		}
		post, err := b.AddFloor(targetPos)
		if err != nil {
			return post, NO
		}
		return post, EMPTY

	case "neighbors":
		workerTile, _ := b.TileAt(workerPos)
		targetTile, err := b.TileAt(targetPos)
		if err != nil {
			return b, NO
		}
		if workerTile.IsNeighbor(targetTile) {
			return b, YES
		} else {
			return b, NO
		}

	case "occupied?":
		if worker := b.WorkerAt(targetPos); worker != nil {
			return b, YES
		} else {
			return b, NO
		}

	case "height":
		tile, _ := b.TileAt(targetPos)
		return b, tile.FloorCount()
	}
	panic(fmt.Sprintf("Invalid Command: %v", command))
}

//Represents any Command following a Board specification
type ShortCommand struct {
	Type      string
	Direction Direction
}

//Unmarshals the JSON Array into an ACTUAL JSON object
func (s *ShortCommand) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&s.Type, &s.Direction}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Command: %d != %d", g, e)
	}
	return nil
}

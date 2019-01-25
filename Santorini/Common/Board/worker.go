package board

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// A WorkerSet is a list of IWorkers
type WorkerSet []IWorker

//Helpers for sorting
func (set WorkerSet) Len() int {
	return len(set)
}
func (set WorkerSet) Less(i, j int) bool {
	firstNil := set[i] == nil
	nonNil := (set[i] != nil && set[j] != nil)
	return (firstNil) || (nonNil && set[i].ID() < set[j].ID())
}
func (set WorkerSet) Swap(i, j int) {
	temp := set[i]
	set[i] = set[j]
	set[j] = temp
}

func (set *WorkerSet) UnmarshalJSON(buf []byte) error {
	var typed WorkerSet
	var impl []Worker
	err := json.Unmarshal(buf, &impl)

	for _, worker := range impl {
		typed = append(typed, IWorker(worker))
	}

	set = &typed

	return err
}

//A Worker is a struct which has the following receivers
type IWorker interface {
	//Returns the Pos this Worker is on
	Pos() Pos
	//Returns a string representing the player that owns this Worker.
	Owner() string
	//Returns an integer ID differentiating a player's workers
	ID() int
	//Returns a new Worker having moved the existing to the other Pos
	Move(target Pos) IWorker
	//Worker's Name
	Name() string
}

type Worker struct {
	currentPos Pos
	owner      string
	id         int
}

func (w Worker) Pos() Pos {
	return w.currentPos
}

func (w Worker) Owner() string {
	return w.owner
}

func (w Worker) ID() int {
	return w.id
}

func (w Worker) Move(target Pos) IWorker {
	return Worker{
		currentPos: target,
		owner:      w.owner,
		id:         w.id,
	}
}

//Return this worker's name
func (w Worker) Name() string {
	return w.owner + (strconv.Itoa(w.id + 1))
}

//Constructs a new Worker struct with the given arguments.
//Returns the newly created struct.
func NewWorker(p Pos, owner string, id int) IWorker {
	return Worker{
		currentPos: p,
		owner:      owner,
		id:         id,
	}
}

// Return whether the given Worker ID is valid
func ValidWID(n int) bool {
	return n >= 0 && n < WorkerCount
}

//Break out a Worker's name into its owner's name, and its ID
//returns an error if there is an issue parsing
func ParseWorkerName(name string) (string, int, error) {
	if len(name) < 2 {
		return "", -1, errors.New("Worker name not long enough")
	}

	playerName := name[:len(name)-1]
	workerID, err := strconv.Atoi(name[len(name)-1:])
	if err != nil {
		return "", -1, err
	}
	return playerName, (workerID - 1), nil
}

func (w Worker) MarshalJSON() ([]byte, error) {
	arr := []interface{}{w.Name(), w.Pos().X, w.Pos().Y}

	return json.Marshal(arr)
}

//Unmarshals the JSON Array into an ACTUAL JSON object
func (w *Worker) UnmarshalJSON(buf []byte) error {
	var workerName string
	tmp := []interface{}{&workerName, &w.currentPos.X, &w.currentPos.Y}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("Wrong number of fields in Worker: %d != %d", g, e)
	}

	if owner, id, err := ParseWorkerName(workerName); err != nil {
		return fmt.Errorf("Failed to parse Worker name: %s", workerName)
	} else {
		w.owner = owner
		w.id = id
	}

	return nil
}

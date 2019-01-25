package board

import (
	"bytes"
	"testing"

	"github.com/CS4500-F18/dare-rebr/Santorini/Lib"
)

const (
	p1 = "uno"
	p2 = "dos"
)

func TestJSONStreams_Encoder(t *testing.T) {
	var b []byte
	rw := bytes.NewBuffer(b)
	encoder, decoder := lib.JSONStreams(rw)

	workers := []IWorker{
		NewWorker(Pos{X: 0, Y: 0}, p1, 0),
		NewWorker(Pos{X: 0, Y: 1}, p1, 1),
		NewWorker(Pos{X: 1, Y: 0}, p2, 0),
		NewWorker(Pos{X: 1, Y: 1}, p2, 1),
	}
	inBoard := BoardWithWorkers(workers)

	err := encoder.Encode(inBoard)
	if err != nil {
		t.Fatalf("Failed to encode board: %v", err)
	}

	boardVal := BaseBoard()
	err = decoder.Decode(&boardVal)
	if err != nil {
		t.Fatalf("Failure to decode board: %v", err)
	}
	outBoard := IBoard(boardVal)

	if !sliceEqualityString(inBoard.Players(), outBoard.Players()) {
		t.Fatalf("Differing Players")
	}

	if !sliceEqualityWorker(inBoard.WorkersFor(p1), outBoard.WorkersFor(p1)) {
		t.Fatalf("Differing Workers")
	}
}

func sliceEqualityString(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sliceEqualityWorker(first, second []IWorker) bool {
	if len(first) != len(second) {
		return false
	}

	for i := range first {

		a := first[i]
		b := second[i]

		if (a == nil) != (b == nil) {
			return false
		}

		if a.Name() != b.Name() || a.Owner() != b.Owner() {
			return false
		}

		if a.ID() != b.ID() {
			return false
		}

		if a.Pos().X != b.Pos().X || a.Pos().Y != b.Pos().Y {
			return false
		}
	}
	return true
}

func placeWorkersOnBoard(b IBoard, names [2]string, locs [4]Pos) (IBoard, error) {
	b, e1 := b.PlaceWorker(locs[0], names[0])
	if e1 != nil {
		return b, e1
	}

	b, e2 := b.PlaceWorker(locs[1], names[0])
	if e2 != nil {
		return b, e2
	}
	b, e3 := b.PlaceWorker(locs[2], names[1])
	if e3 != nil {
		return b, e3
	}
	b, e4 := b.PlaceWorker(locs[3], names[1])
	if e4 != nil {
		return b, e4
	}
	return b, nil
}

package board

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//corner
func TestNeighbors(t *testing.T) {
	inPos := Pos{0, 0}
	expectedOut := []Pos{
		{0, 1}, {1, 0}, {1, 1},
	}

	out := inPos.Neighbors()

	if !cmp.Equal(out, expectedOut) {
		fmt.Printf("Expected %v", expectedOut)
		fmt.Printf("Actual %v", out)
		t.Fail()
	}
}

//wall
func TestGetNeighbors2(t *testing.T) {
	inPos := Pos{5, 3}
	expectedOut := []Pos{
		{4, 2}, {4, 3}, {4, 4}, {5, 2}, {5, 4},
	}

	out := inPos.Neighbors()

	if !cmp.Equal(out, expectedOut) {
		fmt.Printf("Expected %v", expectedOut)
		fmt.Printf("Actual %v", out)
		t.Fail()
	}
}

//middle
func TestGetNeighbors3(t *testing.T) {
	inPos := Pos{3, 3}
	expectedOut := []Pos{
		{2, 2}, {2, 3}, {2, 4}, {3, 2},
		{3, 4}, {4, 2}, {4, 3}, {4, 4},
	}

	out := inPos.Neighbors()

	if !cmp.Equal(out, expectedOut) {
		fmt.Printf("Expected %v", expectedOut)
		fmt.Printf("Actual %v", out)
		t.Fail()
	}

	tile := NewTile(inPos)

	for _, p := range expectedOut {
		newTile := NewTile(p)
		if !tile.IsNeighbor(newTile) {
			t.Fail()
		}
	}
}

func TestPosDistance2(t *testing.T) {
	centerPos := Pos{3, 3}

	if PosDistance(centerPos, Pos{4, 4}) != 1 {
		t.Fail()
	}
	if PosDistance(centerPos, Pos{2, 2}) != 1 {
		t.Fail()
	}
	if PosDistance(centerPos, Pos{2, 4}) != 1 {
		t.Fail()
	}
	if PosDistance(centerPos, Pos{4, 2}) != 1 {
		t.Fail()
	}

}

func TestPosDistance1(t *testing.T) {
	centerPos := Pos{3, 3}

	if PosDistance(centerPos, Pos{4, 3}) != 1 {
		t.Fail()
	}
	if PosDistance(centerPos, Pos{3, 4}) != 1 {
		t.Fail()
	}
	if PosDistance(centerPos, Pos{3, 2}) != 1 {
		t.Fail()
	}
	if PosDistance(centerPos, Pos{2, 3}) != 1 {
		t.Fail()
	}

}

//does find worker work after initial worker placement
func TestFindWorker(t *testing.T) {
	workers := []IWorker{
		NewWorker(Pos{0, 0}, PLAYER_1, 0),
		NewWorker(Pos{1, 1}, PLAYER_1, 1),
		NewWorker(Pos{0, 1}, PLAYER_2, 0),
		NewWorker(Pos{1, 0}, PLAYER_2, 1),
	}
	b := BoardWithWorkers(workers)

	p1w1, eA := b.FindWorker(PLAYER_1, 0)
	p1w2, eB := b.FindWorker(PLAYER_1, 1)
	p2w1, eC := b.FindWorker(PLAYER_2, 0)
	p2w2, eD := b.FindWorker(PLAYER_2, 1)

	posITileA := p1w1.Pos()
	posITileB := p1w2.Pos()
	posITileC := p2w1.Pos()
	posITileD := p2w2.Pos()

	if posITileA.X != 0 || posITileA.Y != 0 || eA != nil {
		t.Fail()
	}

	if posITileB.X != 1 || posITileB.Y != 1 || eB != nil {
		t.Fail()
	}

	if posITileC.X != 0 || posITileC.Y != 1 || eC != nil {
		t.Fail()
	}

	if posITileD.X != 1 || posITileD.Y != 0 || eD != nil {
		t.Fail()
	}
}

//does find worker work after a move
func TestFindWorker_PostMove(t *testing.T) {
	board := SetupBoard(Pos{0, 0}, Pos{1, 1}, Pos{5, 5}, Pos{4, 4})

	//move workers - no validation check in move, should be able to move anywhere on the board
	board, eM1 := board.Move(PLAYER_1, 0, Pos{0, 1})
	board, eM2 := board.Move(PLAYER_1, 1, Pos{0, 0})
	board, eM3 := board.Move(PLAYER_2, 0, Pos{4, 5})
	board, eM4 := board.Move(PLAYER_2, 1, Pos{5, 4})

	if eM1 != nil || eM2 != nil || eM3 != nil || eM4 != nil {
		t.Fail()
	}

	p1w1, eA := board.FindWorker(PLAYER_1, 0)
	p1w2, eB := board.FindWorker(PLAYER_1, 1)
	p2w1, eC := board.FindWorker(PLAYER_2, 0)
	p2w2, eD := board.FindWorker(PLAYER_2, 1)

	posITileA := p1w1.Pos()
	posITileB := p1w2.Pos()
	posITileC := p2w1.Pos()
	posITileD := p2w2.Pos()

	if posITileA.X != 0 || posITileA.Y != 1 || eA != nil {
		t.Fail()
	}

	if posITileB.X != 0 || posITileB.Y != 0 || eB != nil {
		t.Fail()
	}

	if posITileC.X != 4 || posITileC.Y != 5 || eC != nil {
		t.Fail()
	}

	if posITileD.X != 5 || posITileD.Y != 4 || eD != nil {
		t.Fail()
	}
}

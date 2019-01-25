package board

import (
	"testing"

	lib "github.com/CS4500-F18/dare-rebr/Santorini/Lib"
)

//test that New tile has the intended Pos
func TestTile_NewTile_Pos(t *testing.T) {
	tile := NewTile(Pos{0, 0})
	p := Pos{0, 0}
	if tile.Pos() != p {
		t.Fail()
	}
}

//test that the new tile has the intended pos
func TestTile_NewTile_Pos2(t *testing.T) {
	tile := NewTile(Pos{3, 5})
	p := Pos{3, 5}
	if tile.Pos() != p {
		t.Fail()
	}
}

//test that the new tile has proper default values
func TestTile_NewTile_DefaultValues(t *testing.T) {
	tile := NewTile(Pos{4, 4})
	if tile.FloorCount() != 0 {
		t.Fail()
	}
}

//test that valid neighbors are recognized as such
func Testtile_IsNeighbor_ValidNeighbor(t *testing.T) {
	t0 := NewTile(Pos{0, 0})

	t1 := NewTile(Pos{0, 1})
	t2 := NewTile(Pos{1, 1})
	t3 := NewTile(Pos{1, 0})

	n1 := t0.IsNeighbor(t1)
	n2 := t0.IsNeighbor(t2)
	n3 := t0.IsNeighbor(t3)

	if !(n1 && n2 && n3) {
		t.Fail()
	}
}

func TestTile_AddFloor_Valid(t *testing.T) {
	tile := tile{pos: Pos{0, 0}, floorCount: 0}

	var newTile ITile = tile
	for i := 0; i < MaxBuildingHeight; i++ {
		newTile, _ = newTile.AddFloor()
		if height := newTile.FloorCount(); height != i+1 {
			t.Errorf("Floor count not incremented (expected %v, got %v)", i+1, height)
		}
	}
}

func TestTile_AddFloor_Error(t *testing.T) {
	tile := tile{pos: Pos{0, 0}, floorCount: 4}

	//expecting an error, attempting to build 5th floor
	_, err := tile.AddFloor()
	if err == nil {
		t.Fail()
	}
}

func TestIntPresent(t *testing.T) {
	tArr := []int{2, 3, 4, 5, 6, 7, 8}

	for _, v := range tArr {
		if !lib.IntPresent(tArr, v) {
			t.Fail()
		}
	}
}

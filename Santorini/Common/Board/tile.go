package board

import (
	"fmt"
)

type TileSet map[int]ITile
type TileMap map[int]TileSet

//A ITile is a struct which has the following receivers
type ITile interface {
	//Returns a Pos struct representing the coordinate pair position of this ITile
	Pos() Pos
	//Returns the number of floors on this ITile. 0 Implies the ITile is empty,
	//1-3 implies there is a building of that height on the ITile, and
	//4 implies there is a completed building on the ITile.
	FloorCount() int
	//Adds a floor to the given ITile and returns a new ITile.
	//Returns nil for an error if this operation was successful,
	//and an error if this operation is impossible.
	//Examples of impossible AddFloors include:
	// * add when this ITile's GetFloorCount() would return a value greater or equal to 4
	AddFloor() (ITile, error)
	//Returns true if the given ITile is a neighbor of the current ITile. Returns false otherwise.
	IsNeighbor(ITile) bool
}

type tile struct {
	pos        Pos
	floorCount int
}

func (t tile) Pos() Pos {
	return t.pos
}

func (t tile) FloorCount() int {
	return t.floorCount
}

func (t tile) AddFloor() (ITile, error) {
	if t.floorCount >= MaxBuildingHeight {
		return t, fmt.Errorf("The maximum floor height of %v has been reached", MaxBuildingHeight)
	}
	return tile{pos: t.pos, floorCount: t.floorCount + 1}, nil
}

func (t tile) IsNeighbor(other ITile) bool {
	myPos := t.Pos()
	otherPos := other.Pos()

	same := otherPos == myPos
	boardDist := PosDistance(myPos, otherPos)

	return !same && boardDist == 1
}

//Constructs a new tile struct with no floors at the given position.
//Returns the newly created struct.
func NewTile(p Pos) ITile {
	return tile{pos: p}
}

// Helper to create an ITile at the given height
func CustomTile(p Pos, height int) ITile {
	tile := NewTile(p)
	for i := 0; i < height; i++ {
		tile, _ = tile.AddFloor()
	}
	return tile
}

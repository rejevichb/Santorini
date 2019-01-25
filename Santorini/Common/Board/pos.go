package board

import (
	"encoding/json"
	"fmt"
	"math"
)

//A Pos is a struct containing two integers
type Pos struct {
	//X is an integer representing the X coordinate of a coordinate pair
	X int
	//Y is an integer representing the Y coordinate of a coordinate pair
	Y int
}

// Return whether a given Position is in bounds on any Board
func (p Pos) InBounds() bool {
	return inBounds(p.X) && inBounds(p.Y)
}

// Return whether a given integer is within the width/height boundary of any Board
func inBounds(n int) bool {
	return n >= 0 && n < NormalBoardSize
}

//Returns the manhattan distance between the given two positions
func PosDistance(a Pos, b Pos) int {
	return int(math.Max(math.Abs(float64(a.X-b.X)), math.Abs(float64(a.Y-b.Y))))
}

// Returns a list of positions that represent the x y coordinates of all of the neighbors of the given position,
// only if they exist
func (existing Pos) Neighbors() []Pos {
	x := existing.X
	y := existing.Y

	moves := make([]Pos, 0)
	for xShift := -1; xShift <= 1; xShift++ {
		for yShift := -1; yShift <= 1; yShift++ {
			targetPos := Pos{x + xShift, y + yShift}
			if targetPos != existing && targetPos.InBounds() {
				moves = append(moves, targetPos)
			}
		}
	}

	return moves
}

func (p Pos) MarshalJSON() ([]byte, error) {
	arr := []interface{}{p.X, p.Y}
	return json.Marshal(arr)
}

func (p *Pos) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&p.X, &p.Y}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("Wrong number of fields in Pos: %d != %d", g, e)
	}
	return nil
}

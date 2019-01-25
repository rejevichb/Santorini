package rules

import (
	"testing"

	"github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	common "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
)

const PLAYER_1 = "uno"
const PLAYER_2 = "dos"

var (
	empty = board.BaseBoard()
)

/*
#########################################################
##################### Move Rules ########################
#########################################################
*/

// For testing purposes
func SetupBoard(p_1, p_2, p_3, p_4 board.Pos) board.IBoard {
	og_board := board.BaseBoard()
	testBoard, _ := og_board.PlaceWorker(p_1, PLAYER_1)
	testBoard, _ = testBoard.PlaceWorker(p_2, PLAYER_1)
	testBoard, _ = testBoard.PlaceWorker(p_3, PLAYER_2)
	testBoard, _ = testBoard.PlaceWorker(p_4, PLAYER_2)
	return testBoard
}

//################### moveInBounds(m Move) bool ###################

//test that an in-bound move returns true
func TestRules_moveInBounds(t *testing.T) {
	workerPos := common.Pos{X: 1, Y: 1}
	targetPos := common.Pos{X: 0, Y: 0}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	if !moveInBounds(empty, workerTile, targetTile) {
		t.Fail()
	}
}

//test an out-of bounds move returns false
func TestRules_moveInBounds_OutOfBounds(t *testing.T) {
	workerPos := common.Pos{X: 1, Y: 1}
	targetPos := common.Pos{X: 5, Y: 6}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	if moveInBounds(empty, workerTile, targetTile) {
		t.Fail()
	}
}

/*
################### adjacencyMoveRule(m Move) bool ###################
*/
func TestRules_adjacencyBuildRule_Adjacent(t *testing.T) {
	workerPos := common.Pos{X: 0, Y: 0}
	targetPos := common.Pos{X: 1, Y: 1}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	valid := adjacencyMoveRule(empty, workerTile, targetTile)
	if !valid {
		t.Fail()
	}
}

func TestRules_adjacencyBuildRule_Adjacent2(t *testing.T) {
	workerPos := common.Pos{X: 4, Y: 4}
	targetPos := common.Pos{X: 5, Y: 4}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	valid := adjacencyMoveRule(empty, workerTile, targetTile)

	if !valid {
		t.Fail()
	}
}

func TestRules_adjacencyMoveRule_NotAdjacent(t *testing.T) {
	workerPos := common.Pos{X: 1, Y: 1}
	targetPos := common.Pos{X: 0, Y: 3}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	valid := adjacencyMoveRule(empty, workerTile, targetTile)
	if valid {
		t.Fail()
	}
}

//tests a worker trying to move to the cell he is on
func TestRules_adjacencyMoveRule_SameCell(t *testing.T) {
	workerPos := common.Pos{X: 4, Y: 4}
	targetPos := common.Pos{X: 4, Y: 4}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	valid := adjacencyMoveRule(empty, workerTile, targetTile)

	if valid {
		t.Fail()
	}
}

/*
################### onlyUp1FloorRule(m Move) bool ###################
*/

//testing a worker on a 0 floor tile to a 1 floor tile
func TestRules_onlyUp1FloorRule_SameFloor(t *testing.T) {
	workerPos := common.Pos{X: 4, Y: 4}
	targetPos := common.Pos{X: 4, Y: 5}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	valid := onlyUp1FloorRule(empty, workerTile, targetTile)

	if !valid {
		t.Fail()
	}
}

//testing a worker moving from a 1 floor tile to a 2 floor tile
func TestRules_onlyUp1FloorRule_1Up(t *testing.T) {
	workerPos := common.Pos{X: 4, Y: 4}
	targetPos := common.Pos{X: 5, Y: 4}

	workerTile := common.CustomTile(workerPos, 1)
	targetTile := common.CustomTile(targetPos, 2)

	valid := onlyUp1FloorRule(empty, workerTile, targetTile)

	if !valid {
		t.Fail()
	}
}

//test a worker moving from a 1 floor tile to a 3 floor tile
func TestRules_onlyUp1FloorRule_2Up(t *testing.T) {
	workerPos := common.Pos{X: 5, Y: 5}
	targetPos := common.Pos{X: 4, Y: 4}

	workerTile := common.CustomTile(workerPos, 1)
	targetTile := common.CustomTile(targetPos, 3)

	valid := onlyUp1FloorRule(empty, workerTile, targetTile)

	if valid {
		t.Fail()
	}
}

/*
################### vacantTileMoveRule(m Move) bool ###################
*/
//tests a move to a vacant tile
func TestRules_vacantTileMoveRule_vacantTile(t *testing.T) {
	workerPos := common.Pos{X: 5, Y: 5}
	targetPos := common.Pos{X: 4, Y: 4}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	valid := vacantTileMoveRule(empty, workerTile, targetTile)

	if !valid {
		t.Fail()
	}
}

//tests an attempted move to an occupied tile
func TestRules_vacantTileMoveRule_OccupiedTile(t *testing.T) {
	workerPos := common.Pos{X: 5, Y: 5}
	targetPos := common.Pos{X: 4, Y: 4}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	b := board.BoardWithWorkers([]board.IWorker{board.NewWorker(targetPos, PLAYER_1, 0)})

	valid := vacantTileMoveRule(b, workerTile, targetTile)

	if valid {
		t.Fail()
	}
}

/*
#########################################################
#########################################################
##################### Build Rules #######################
#########################################################
#########   type BuildRule func(Build) bool   ###########
*/

/*
################### buildInBounds(b Build) bool ###################
*/
//test a build on a cell that is in bounds
func TestRules_buildInBounds_InBounds(t *testing.T) {
	workerPos := common.Pos{X: 1, Y: 1}
	targetPos := common.Pos{X: 0, Y: 0}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	if !buildInBounds(empty, workerTile, targetTile) {
		t.Fail()
	}
}

//test a build on a cell that is not in bounds
func TestRules_buildInBounds_NotInBounds(t *testing.T) {
	workerPos := common.Pos{X: 1, Y: 1}
	targetPos := common.Pos{X: -1, Y: 0}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	if buildInBounds(empty, workerTile, targetTile) {
		t.Fail()
	}

	targetPos = common.Pos{X: 4, Y: 6}
	workerTile = common.CustomTile(workerPos, 0)
	targetTile = common.CustomTile(targetPos, 0)

	if buildInBounds(empty, workerTile, targetTile) {
		t.Fail()
	}
}

/*
################### adjacencyBuildRule(b Build) bool ###################
*/
// tests a build on a cell that is adjacent to the worker making the build
func TestRules_adjacencyBuildRule_IsAdjacent(t *testing.T) {
	workerPos := common.Pos{X: 1, Y: 1}
	targetPos := common.Pos{X: 1, Y: 0}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	valid := adjacencyBuildRule(empty, workerTile, targetTile)

	if !valid {
		t.Fail()
	}
}

//tests a build on a cell not adjacent to the worker who is building
func TestRules_adjacencyBuildRule_NotAdjacent(t *testing.T) {
	workerPos := common.Pos{X: 1, Y: 1}
	targetPos := common.Pos{X: 1, Y: 3}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	valid := adjacencyBuildRule(empty, workerTile, targetTile)

	if valid {
		t.Fail()
	}
}

/*
################### heightBuildRule(b Build) bool ###################
*/
//testing valid height rule when the worker and buildOn heights are 0
func TestRules_heightBuildRule_ValidHeight(t *testing.T) {
	workerPos := common.Pos{X: 1, Y: 1}
	targetPos := common.Pos{X: 1, Y: 0}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	valid := heightBuildRule(empty, workerTile, targetTile)

	if !valid {
		t.Fail()
	}
}

//test a build on a tile with height 3
func TestRules_heightBuildRule_ValidHeight2(t *testing.T) {
	workerPos := common.Pos{X: 1, Y: 1}
	targetPos := common.Pos{X: 1, Y: 0}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 3)

	valid := heightBuildRule(empty, workerTile, targetTile)

	if !valid {
		t.Fail()
	}
}

//test a build on a tile with height 4
func TestRules_heightBuildRule_MaxHeight(t *testing.T) {

	workerPos := common.Pos{X: 1, Y: 1}
	targetPos := common.Pos{X: 1, Y: 0}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 4)

	valid := heightBuildRule(empty, workerTile, targetTile)

	if valid {
		t.Fail()
	}
}

/*
################### vacantTileBuildRule(b Build) bool ###################
*/
//test that a build on a vacant cell is valid
func TestRules_vacantTileBuildRule_Vacant(t *testing.T) {
	workerPos := common.Pos{X: 1, Y: 1}
	targetPos := common.Pos{X: 1, Y: 0}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	valid := vacantTileBuildRule(empty, workerTile, targetTile)

	if !valid {
		t.Fail()
	}
}

//test that a build on an occupied cell is not valid
func TestRules_vacantTileBuildRule_NotVacant(t *testing.T) {
	workerPos := common.Pos{X: 1, Y: 1}
	targetPos := common.Pos{X: 1, Y: 1}

	workerTile := common.CustomTile(workerPos, 0)
	targetTile := common.CustomTile(targetPos, 0)

	b := board.BoardWithWorkers([]board.IWorker{board.NewWorker(targetPos, PLAYER_1, 0)})

	valid := vacantTileBuildRule(b, workerTile, targetTile)

	if valid {
		t.Fail()
	}
}

/*
#########################################################
#########################################################
##################### Place Rules #######################
#########################################################
######   type PlaceRule func(PlaceWorker) bool   ########
*/

/*
################### placeInBounds(p PlaceWorker) bool ###################
*/
func TestRules_placeInBounds_InBounds(t *testing.T) {
	placementTile := common.CustomTile(common.Pos{X: 1, Y: 0}, 0)

	if !placeInBounds(empty, placementTile) {
		t.Fail()
	}
}

func TestRules_placeInBounds_NotInBounds(t *testing.T) {

	if placeInBounds(empty, common.NewTile(common.Pos{X: -1, Y: 0})) {
		t.Fail()
	}

	if placeInBounds(empty, common.NewTile(common.Pos{X: 4, Y: -1})) {
		t.Fail()
	}
}

/*
################### vacantTilePlaceRule(p PlaceWorker) bool ###################
*/
func Test_vacantTilePlaceRule_Vacant(t *testing.T) {
	placeAtPos := common.Pos{X: 3, Y: 1}
	placeAtTile := common.CustomTile(placeAtPos, 1)

	if !vacantTilePlaceRule(empty, placeAtTile) {
		t.Fail()
	}
}

func Test_vacantTilePlaceRule_Occupied(t *testing.T) {
	placeAtPos := common.Pos{X: 3, Y: 1} //testing placing where a worker already is
	placeAtTile := common.CustomTile(placeAtPos, 2)

	b := board.BoardWithWorkers([]board.IWorker{board.NewWorker(placeAtPos, PLAYER_1, 0)})

	if vacantTilePlaceRule(b, placeAtTile) { // should be false - already a worker on 3,1
		t.Fail()
	}
}

func get2TilesFrom2Pos(pFromWorker, pTarget common.Pos, b common.IBoard, t *testing.T) (fromWorkerTile, toTile common.ITile) {
	fromWorkerTile, err := b.TileAt(pFromWorker)
	if err != nil {
		t.Error("error on IBoard.TileAt(p Pos)")
	}
	toTile, err = b.TileAt(pTarget)
	if err != nil {
		t.Error("error on IBoard.TileAt(p Pos)")
	}
	return
}

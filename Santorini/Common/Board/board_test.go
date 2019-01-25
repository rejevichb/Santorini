package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWorkers(t *testing.T) {
	testBoard := SetupBoard(Pos{0, 0}, Pos{1, 1}, Pos{0, 1}, Pos{1, 0})
	pl1w := testBoard.WorkersFor(PLAYER_1)
	assert.Equal(t, len(pl1w), 2, "Should find workers")

	assert.Equal(t, PLAYER_1, pl1w[1].Owner(), "worker should be owned by who placed it")
	assert.Equal(t, Pos{0, 0}, pl1w[0].Pos(), "worker should be where it was placed")
	assert.Equal(t, Pos{1, 1}, pl1w[1].Pos(), "worker should be where it was placed")
	noneWorkers := testBoard.WorkersFor("gibberish")
	assert.True(t, len(noneWorkers) == 0, "Should find not workers")
}

func TestGetITile(t *testing.T) {
	testBoard := BaseBoard()

	_, err := testBoard.TileAt(Pos{0, -1})
	assert.NotNil(t, err, "Should not be able to fetch out of bounds ITile")

	ITile, err := testBoard.TileAt(Pos{1, 1})
	assert.Nil(t, err, "Should be able to fetch in bounds ITile")
	assert.Equal(t, 0, ITile.FloorCount(), "Should have 0 floors when not built on")
}

//The method SetupBoard(...) is used to create a new board and place workers on the board
//The player names are PLAYER_1 and PLAYER_2 respectively.
//The 4 Pos inputted represent the intended worker placement for each player:
//       where:
//              p_1: Player 1 worker 1
//              p_2: Player 1 worker 2
//              p_3: Player 2 worker 1
//              p_4: Player 2 worker 2
// This function returns an IBoard with the worker placements already made.
//see common_test.go for impl of SetupBoard(...)

/*
################### IBoard.TileAt(Pos) (ITile, error) ###################
*/

const PLAYER_1 = "uno"
const PLAYER_2 = "dos"

// For testing purposes
func SetupBoard(positions ...Pos) IBoard {
	workers := []IWorker{
		NewWorker(positions[0], PLAYER_1, 0),
		NewWorker(positions[1], PLAYER_1, 1),
		NewWorker(positions[2], PLAYER_2, 0),
		NewWorker(positions[3], PLAYER_2, 1),
	}
	return BoardWithWorkers(workers)
}

//SANITY CHECK - make sure that initialization happens properly for the entire board
func TestBoard_TileAt_AfterInitialization_BeforeWorkerPlacement(t *testing.T) {
	testBoard := BaseBoard()

	for x := 0; x < NormalBoardSize; x++ {
		for y := 0; y < NormalBoardSize; y++ {
			tile, er := testBoard.TileAt(Pos{x, y})
			if tile.FloorCount() != 0 || er != nil {
				t.Fail()
			}
		}
	}
}

func TestBoard_InitialWorkers(t *testing.T) {
	testBoard := BaseBoard()
	workers := testBoard.Workers()

	assert.Equal(t, len(workers), 0, "Board should have 0 IWorkers")

	for _, w := range workers {
		assert.Equal(t, w, nil, "Board should have all nil workers")
	}
}

//Test that a tile with a worker actually has
func TestBoard_TileAt_AfterWorkerPlacement(t *testing.T) {
	board := SetupBoard(Pos{0, 0}, Pos{1, 1}, Pos{5, 5}, Pos{4, 4})

	//get a worker-occupied tile and verify worker occupation & other data
	tile, e := board.TileAt(Pos{0, 0})

	if e != nil || tile.Pos().X != 0 || tile.Pos().Y != 0 {
		t.Fail()
	}
}

func TestBoard_TileAt_AfterWorkerMove(t *testing.T) {
	board := SetupBoard(Pos{3, 3}, Pos{1, 1}, Pos{5, 5}, Pos{4, 4})

	//move the worker
	board, _ = board.Move(PLAYER_1, 0, Pos{2, 2}) //not testing move here

	//get a worker-occupied tile AFTER a worker move and verify worker occupation & other data
	tile, e := board.TileAt(Pos{2, 2})

	if e != nil || tile.Pos().X != 2 || tile.Pos().Y != 2 {
		t.Fail()
	}
}

func TestBoard_TileAt_AfterFloorAdded(t *testing.T) {
	board := SetupBoard(Pos{3, 3}, Pos{1, 1}, Pos{5, 5}, Pos{4, 4})

	board, _ = board.AddFloor(Pos{2, 2}) //not testing add floor
	board, _ = board.AddFloor(Pos{2, 2})

	ITile, e := board.TileAt(Pos{2, 2})

	if ITile.FloorCount() != 2 || e != nil || ITile.Pos().X != 2 || ITile.Pos().Y != 2 {
		t.Fail()
	}
}

func TestBoard_TileAt_AfterFloorAdded2(t *testing.T) {
	board := SetupBoard(Pos{3, 3}, Pos{1, 1}, Pos{5, 5}, Pos{4, 4})

	board, _ = board.AddFloor(Pos{0, 5}) //not testing add floor
	board, _ = board.AddFloor(Pos{0, 5})

	tile, e := board.TileAt(Pos{0, 5})

	if tile.FloorCount() != 2 || e != nil || tile.Pos().X != 0 || tile.Pos().Y != 5 {
		t.Fail()
	}
}

/*
################### IBoard.WorkersFor() ([WorkerCount]Worker, error) ###################
*/
//Helper method - creates a board with workers at preset locations
// returns each worker for the specified player, and the board that they come from
func WorkersHelper(t *testing.T, p string) (IWorker, IWorker, IBoard) {
	b := SetupBoard(Pos{0, 0}, Pos{1, 1}, Pos{5, 5}, Pos{4, 4})

	workers := b.WorkersFor(p)

	w0 := workers[0]
	w1 := workers[1]

	return w0, w1, b
}

//make sure Workers returns the correct pos after placement
func TestBoard_Workers_Pos_AfterPlacement(t *testing.T) {
	w0, w1, _ := WorkersHelper(t, PLAYER_1)

	if w0.Pos().X != 0 || w0.Pos().Y != 0 {
		t.Error("Worker 0 not at expected location:", w0.Pos())
	}

	if w1.Pos().X != 1 || w1.Pos().Y != 1 {
		t.Error("Worker 1 not at expected location:", w1.Pos())
	}
}

//test that the board returns the correct worker location after a worker move
func TestBoard_Workers_Pos_AfterMove(t *testing.T) {
	_, _, b := WorkersHelper(t, PLAYER_1)

	b, e := b.Move(PLAYER_1, 0, Pos{0, 1})
	if e != nil {
		t.Fail()
	}

	newWorkers := b.WorkersFor(PLAYER_1)

	if newWorkers[0].Pos().X != 0 || newWorkers[0].Pos().Y != 1 {
		t.Fail()
	}
}

//make sure the workers returned in Workers have the correct owner
func TestBoard_Workers_WorkerOwner(t *testing.T) {
	w0, w1, _ := WorkersHelper(t, PLAYER_1)
	w2, w3, _ := WorkersHelper(t, PLAYER_2)

	if w0.Owner() != PLAYER_1 || w1.Owner() != PLAYER_1 {
		t.Fail()
	}
	if w2.Owner() != PLAYER_2 || w3.Owner() != PLAYER_2 {
		t.Fail()
	}
}

/*
################### IBoard.Move(w Worker, target Pos) (IBoard, error) ###################
*/
func TestBoard_Move1(t *testing.T) {
	b := SetupBoard(Pos{3, 3}, Pos{1, 1}, Pos{5, 5}, Pos{4, 4})

	//get a worker we want to move from the board
	p1Workers := b.WorkersFor(PLAYER_1) //not testing get workers here, ignore error
	w0 := p1Workers[0]
	if w0.Pos().X != 3 || w0.Pos().Y != 3 {
		t.Fail()
	}

	b, e := b.Move(PLAYER_1, 0, Pos{2, 2})

	p1Workers = b.WorkersFor(PLAYER_1)
	w0 = p1Workers[0]

	if e != nil {
		t.Fail()
	}

	if w0.Pos().X != 2 || w0.Pos().Y != 2 {
		t.Fail()
	}
}

func TestBoard_Move2(t *testing.T) {
	b := SetupBoard(Pos{3, 3}, Pos{1, 1}, Pos{5, 5}, Pos{4, 4})

	//get a worker we want to move from the board
	p2Workers := b.WorkersFor(PLAYER_2) //not testing get workers here, ignore error
	worker2 := p2Workers[1]

	if worker2.Pos().X != 4 || worker2.Pos().Y != 4 {
		t.Errorf("Player 2's first worker not at expected location pre-move, Expected %v, Actual %v", Pos{4, 4}, worker2.Pos())
	}

	newB, err := b.Move(PLAYER_2, 1, Pos{4, 3})
	if err != nil {
		t.Error("Error moving Player 2's first worker:", err)
	}

	p2Workers = newB.WorkersFor(PLAYER_2)
	worker2 = p2Workers[1]

	if worker2.Pos().X != 4 || worker2.Pos().Y != 3 {
		t.Errorf("Player 2's first worker not at expected location post-move, Expected %v, Actual %v", Pos{4, 3}, worker2.Pos())
	}
}

func TestBoard_Move3(t *testing.T) {
	b := SetupBoard(Pos{2, 2}, Pos{1, 1}, Pos{3, 3}, Pos{4, 4})

	p1Workers := b.WorkersFor(PLAYER_1) //not testing get workers
	w1 := p1Workers[1]

	if w1.Pos().X != 1 || w1.Pos().Y != 1 {
		t.Fail()
	}

	b, e := b.Move(PLAYER_1, 1, Pos{1, 2})

	if e != nil {
		t.Fail()
	}

	p1Workers = b.WorkersFor(PLAYER_1)
	w1 = p1Workers[1]

	if w1.Pos().X != 1 || w1.Pos().Y != 2 {
		t.Fail()
	}
}

/*
################### IBoard.AddFloor(target Pos) (IBoard, error) ###################
*/

func TestBoard_AddFloor(t *testing.T) {
	b := SetupBoard(Pos{2, 2}, Pos{1, 1}, Pos{3, 3}, Pos{4, 4})

	b, e := b.AddFloor(Pos{2, 5})

	if e != nil {
		t.Fail()
	}

	ITile, _ := b.TileAt(Pos{2, 5})

	if ITile.FloorCount() != 1 {
		t.Fail()
	}
}

func TestBoard_AddFloor2(t *testing.T) {
	b := SetupBoard(Pos{2, 2}, Pos{1, 1}, Pos{3, 3}, Pos{4, 4})

	b, e := b.AddFloor(Pos{3, 4})
	b, e2 := b.AddFloor(Pos{3, 4})

	if e != nil || e2 != nil {
		t.Fail()
	}

	ITile, _ := b.TileAt(Pos{3, 4})

	if ITile.FloorCount() != 2 {
		t.Fail()
	}
}

func TestBoard_AddFloor3(t *testing.T) {
	b := SetupBoard(Pos{2, 2}, Pos{1, 1}, Pos{3, 3}, Pos{4, 4})

	b, e := b.AddFloor(Pos{5, 5})
	b, e1 := b.AddFloor(Pos{5, 5})
	b, e2 := b.AddFloor(Pos{5, 5})
	b, e3 := b.AddFloor(Pos{5, 5})

	if e != nil || e1 != nil || e2 != nil || e3 != nil {
		t.Fail()
	}

	ITile, _ := b.TileAt(Pos{5, 5})

	if ITile.FloorCount() != 4 {
		t.Fail()
	}
}

func TestBoard_AddFloor_Max(t *testing.T) {
	b := SetupBoard(Pos{2, 2}, Pos{1, 1}, Pos{3, 3}, Pos{4, 4})

	b, e := b.AddFloor(Pos{5, 5})
	b, e1 := b.AddFloor(Pos{5, 5})
	b, e2 := b.AddFloor(Pos{5, 5})
	b, e3 := b.AddFloor(Pos{5, 5})

	if e != nil || e1 != nil || e2 != nil || e3 != nil {
		t.Fail()
	}

	b, e4 := b.AddFloor(Pos{5, 5}) //try to add 5th floor

	if e4 == nil {
		t.Fail()
	}

	ITile, _ := b.TileAt(Pos{5, 5})

	//make sure the floor count is 4 not 5
	if ITile.FloorCount() != 4 {
		t.Fail()
	}

}

/*
################### IBoard.PlaceWorker(p Pos, owner string) (IBoard, error) ###################
*/
func TestBoard_PlaceWorker(t *testing.T) {
	og_b := BaseBoard()

	b, e := og_b.PlaceWorker(Pos{5, 5}, PLAYER_1)
	b, e2 := b.PlaceWorker(Pos{4, 4}, PLAYER_1)

	workers := b.WorkersFor(PLAYER_1) //not testing get workers
	w0 := workers[0]
	w1 := workers[1]

	//fail on error
	if e != nil || e2 != nil {
		t.Errorf("Failures in PlaceWorker")
	}

	//fail on bad worker data
	if w0.Pos().X != 5 || w0.Pos().Y != 5 || w0.Owner() != PLAYER_1 {
		t.Errorf("Player 1's first worker not at expected location, found at (%v, %v)", w0.Pos().X, w0.Pos().Y)
	}

	//fail on bad worker data
	if w1.Pos().X != 4 || w1.Pos().Y != 4 || w1.Owner() != PLAYER_1 {
		t.Errorf("Player w's second worker not at expected location, found at (%v, %v)", w1.Pos().X, w1.Pos().Y)
	}
}

func TestBoard_PlaceWorker_Player2(t *testing.T) {
	og_b := BaseBoard()

	b, e := og_b.PlaceWorker(Pos{5, 5}, PLAYER_1)
	b, e2 := b.PlaceWorker(Pos{0, 0}, PLAYER_2)

	workers1 := b.WorkersFor(PLAYER_1)
	workers2 := b.WorkersFor(PLAYER_2)
	assert.NotNil(t, workers1[0], "Player 1 should have a non-nil first worker")
	assert.NotNil(t, workers2[0], "Player 2 should have a non-nil first worker")

	assert.Nil(t, e, "No error placing first worker")
	assert.Nil(t, e2, "No error placing second worker")
}

/*
################### IBoard.Players() []string ###################
*/
func TestBoard_Players(t *testing.T) {
	workers := []IWorker{
		NewWorker(Pos{0, 0}, PLAYER_1, 0),
		NewWorker(Pos{1, 1}, PLAYER_2, 1),
	}
	b := BoardWithWorkers(workers)

	players := b.Players()
	if players[0] != PLAYER_1 || players[1] != PLAYER_2 {
		t.Fail()
	}
}

/*
################### IBoard.Dimensions() (int, int) ###################
*/
func TestBoard_Dimensions(t *testing.T) {
	board := BaseBoard()

	xDem, yDem := board.Dimensions()

	if xDem != NormalBoardSize || yDem != NormalBoardSize {
		t.Fail()
	}
}

package design

//Purpose: To represent the physical game board of the game Santorini.
//Will just contain the logic regarding the setup and usage of pieces
//on the game board, not the logic regarding how the game is played.

//A Pos is a struct containing two integers
type Pos struct {
	//X is an integer representing the X coordinate of a coordinate pair
	X int
	//Y is an integer representing the Y coordinate of a coordinate pair
	Y int
}

//A ITile is a struct which has the following receivers
type ITile interface {
	//Returns a Pos struct representing the coordinate pair position of this ITile
	GetPos() Pos
	//Returns the number of floors on this ITile. 0 Implies the ITile is empty,
	//1-3 implies there is a building of that height on the ITile, and
	//4 implies there is a completed building on the ITile.
	GetFloorCount() int
	//Adds a floor to the given ITile. Returns nil if this operation was successful,
	//and an error if this operation is impossible.
	//Examples of impossible AddFloors include:
	// * add when this ITile's GetFloorCount() would return a value greater or equal to 4
	AddFloor() error
	//Accepts a Worker to be on this ITile. If this is impossible, returns an error.
	//Examples of impossible AcceptWorkerArrivals include:
	// * accept when this ITile's GetFloorCount() would return a value greater or equal to 4
	// * accept when there already is a Worker on this ITile
	AcceptWorkerArrival() error
	//Accepts a Worker to leave this ITile. If this is impossible, returns an error.
	//Examples of impossible AcceptWorkerDeparture include:
	// * accept when there is no Worker on this ITile
	AcceptWorkerDeparture() error
}

//A Worker is a struct which has the following receivers
type Worker interface {
	//Returns a pointer to a ITile struct representing the ITile this Worker is currently on.
	GetITile() (*ITile, error)
	//Returns an integer representing the player that owns this Worker.
	GetOwner() int
	//Moves the worker to the given ITile by changing the ITile the Worker is currently on.
	//Returns an error in the case that this Move is impossible, and doesn't change
	//the state of this Worker.
	//Examples of impossible Moves include:
	// * a move to a ITile whose AcceptWorkerArrival returns an error.
	// * a move to a non-adjacent ITile (the x and y coordinates of the target ITile are within plus or minus one of the coordinates this Worker is on
	// * a move to a ITile whose GetFloodCount() returns a value higher more than one more than returned by the ITile the Worker currently is on
	Move(other ITile) error
	//Adds a floor to the given ITile by changing the ITile the Worker is given.
	//Returns an error in the case that this AddFloorTo is impossible, and doesn't change
	//the state of this ITile.
	//Examples of impossible AddFloorTos include:
	// * add to a ITile whose AddFloor returns an error.
	// * add to a non-adjacent ITile (the x and y coordinates of the target ITile are within plus or minus one of the coordinates this Worker is on
	AddFloorTo(other ITile) error
}

//A Board is a struct which has the following receivers
type Board interface {
	//Returns the ITile in this Board at the coordinate pair represented by the given target Pos.
	//If the target Pos is outside of this Board, returns an error instead.
	GetITileAt(target Pos) (ITile, error)
	//Returns all of the workers currently on this Board
	GetWorkers() []Worker
	//Attempts to move the given worker to the given target position. If this move is impossible (the Worker's Move() returns an error), returns an error.
	Move(w Worker, target Pos) error
	//Attempts to have the given worker add a floor to the given target position. If this move is impossible (the Worker's AddFloorTo() returns an error), returns an error.
	AddFloor(w Worker, target Pos) error
	//Confirms that AcceptWorkerArrival() returns no errors (returning an error if so) for the ITile corresponding to the given Pos,
	//creates a new Worker, places the Worker on that ITile and returns a pointer to the Worker struct.
	//The Worker is added the the Board's collection of Workers.
	PlaceWorker(p Pos, owner int) (*Worker, error)
	//Returns the height and width (as number of ITiles) for the board
	GetDimensions() (width int, height int)
}

package design

import common "github.com/CS4500-F18/dare-rebr/Santorini/Common"

//the purpose of a strategy object is to encapsulate the decision making logic for a
//player of Santorini, whether this strategy be a human player or an autonomous logic program

//Represents a strategy object that can take a turn given the state of the read only board at any time
type Strategy interface {
	//Returns the turn that is desired to be taken (as a combination of a move and a build.
	//This function will be called by the admin when it is desired that the strategy take a turn
	//And the strategy will use either internal logic or input from a human player to create a
	//turn object. This object will be sent to the matching player object for this stratgy,
	//which will make move and build requests to the admin on behalf of the strategy,
	//to prevent the strategy from having access to any admin only instances.
	GetTurn() Turn
}

//Represents the specification of a whole turn by a strategy.
type Turn struct {
	//The id of the worker selected for this turn.
	W int
	//The position on the board the selected worker should move to.
	MoveTo common.Pos
	//The position on the board the selected worker should build on,
	//once it has moved.
	AddTo common.Pos
}

//Creates a new strategy object that will be able to take turns.
//The given readonlyboard is the board that the game of Santorini
//will be taking place on, and the given string is the name of
//the player that this strategy is using.
func NewStrategy(b ReadOnlyBoard, player string) Strategy {
	return nil //To be completed
}

//Represents a read only version of the information contained by a cell on the board
type Cell struct {
	//an integer representing the height of this cell (0 being flat, and each integer above 0 representing one additional floor)
	height int
	//a string representing a worker occupying this cell if there is one occupying this cell, and the empty string if otherwise
	//a worker string will always be a string representing the name of the player that owns it, followed by a single digit which
	//represents its id.
	worker string
}

//Represents a read only version of a board object from the common package. Does not allow access to the board itself,
//but allows the update of the read only copy at any time. Unless this update is called, it is stateless.
type ReadOnlyBoard interface {
	//Returns a map of Positions to Cell objects, which represent the state of that position on the board
	View() map[common.Pos]Cell
	//Updates the representation of the board using the current state of the internally stored board
	Update()
}

//Takes in a board and returns a new ReadOnlyBoard which returns the state of the board
//in a non mutable fashion
func NewReadOnlyBoard(b common.IBoard) ReadOnlyBoard {
	return nil //To be completed
}

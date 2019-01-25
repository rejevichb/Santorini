package design

// The referee carries a game to its completion given that the strategies and players are instantiated by a
// main function or connection code

type Referee interface {
	// Runs a game and returns the winner.
	// The Referee will end the game and declare a winner if any call to a strategy object takes longer than a given timeout time
	// The referee initializes the game board at the start of this function, using the given board dimensions.
	// It then initializes players based off of the given names. It stores these for future use
	// Phase 1 (placing):
	//  - until workerCount amount of workers have been placed for each player:
	//  - the referee calls Strategy.GetPlace()
	//  - pipes result through player.PlaceWorkerRequest() for the matching player
	//  - pipes result through CheckPlaceWorker(). if it returns false, then the game ends and the other player wins.
	//  - if the result is true, it does board.PlaceWorker() with the corresponding information and also adds the worker to player's available workers
	// Phase 2:
	// - Referee calls Strategy.GetTurn() and saves the result
	// - pipes the Move part of the Turn through Player.MakeMoveReguest() for the matching player
	// - pipes the result through CheckMove(). if it returns false, the game ends and the other player wins
	// - if the result is true, it does board.Move() with the corresponding information
	// - it checks if there is a winner with CheckPostMove() and returns a winner if present
	// - it marks the selected worker on the current player
	// - these steps are then repeated for building, using Build specific, instead of Move specific, moves
	// - the current player is then switched and return to the beginning of Phase 2
	Play() string
}

// instantiates a new referee that can run a game
func NewReferee(stratA Strategy, nameA string, stratB Strategy, nameB string, workerCount int, boardHeight int) Referee {
	return nil
}

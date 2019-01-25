package design

// //The purpose of the player package is to encapsulate the work that a player needs to be responsible for. This includes
// //going through the setup of a game, taking turns (including deciding on a turn based on the current state of the game),
// //and leaving the game once finished.

// //Represents a player of the game Santorini
// type Player interface {
// 	//Returns the name of the string identifying this player.
// 	GetName() string
// 	//Shows a view of the board so that the Player can make decisions based on this information.
// 	ViewRequest()
// 	//Returns a function that, when passed a Board object (given by the admin) creates a complete PlaceWorker to be attempted.
// 	PlaceWorkerRequest(Pos) func(Board) PlaceWorker
// 	//Returns a function that, when passed a Board object (given by the admin) creates a complete Turn to be attempted.
// 	TakeTurnRequest(worker string, moveTo Pos, addTo Pos) func(Board) Turn
// }

// //Creates a new player object given the string name of this player.
// func NewPlayer(name string) Player {
// 	return nil //To be Player
// }

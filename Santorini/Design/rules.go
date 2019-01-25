package design

// //The purpose of the rules package is to orchestrate checks for all logical (non physical [I.E. workers cannot go on rounded 4th floor building rooftops])
// //rules about the game. This will include managing which player's turn it is, managing what moves during a turn are valid or not,
// //and relaying this information to the player, and determining any ending conditions for the game.

// //A Rule is a predicate style function that takes a turn and returns true if it would be valid, and false if not.
// // DOES NOT ALTER THE STATE OF THE BOARD, WORKER, OR PLAYER.
// //Simply checks ahead of time if the turn would be valid.
// //type Rule func(Turn) bool

// //Ex.
// /*
// func AdjacencyRule(t Turn) {
// 	w := t.GetWorker()
// 	b := b.GetBoard()
// 	mITile,err := b.GetITileAt(GetMoveTo())
// 	aITile,err := b.GetITileAt(GetAddTo())
// 	return w.GetITile().IsNeighbors(mITile) && mITile.isNeighbors(aITile)
// }

// func OnlyUp1FloorRule(t Turn) {
// 	//Checks that we only go up one floor during our move
// } */

// //A PlaceRule is a predicate style function that takes a PlaceWorker and returns true if it would be valid, and false if not.
// //DOES NOT ALTER THE STATE OF THE BOARD, WORKER, OR PLAYER
// //Simply checks ahead of time if the place worker would be valid.
// //type PlaceRule func(PlaceWorker) bool

// //Represents the set of all rules used to check a Turn for validity
// //type RuleSet []Rule

// //Represents the set of all placerules used to check a PlaceWorker for validity.
// //type PlaceRuleSet []PlaceRule

// //Ex.
// //const BasicRuleSet = []Rule{AdjacencyRule, OnlyUp1FloorRule, ... }
// //const BasicPlaceRuleSet = []PlaceRuleSet{NonOverlappingRule, ... }

// //A turn represents a full set of moves that make up a turn. This includes specifying a worker,
// //moving the worker, and having the worker execute a build.

// //Takes in a player,a board, a worker, and two Positions. The first position
// //represents the position of the ITile the worker should move to, and the second represents the position
// //of the ITile the worker should add to. The player is the player taking the action, the board is the current state of
// //the board, and the worker is the worker chosen to execute the turn.
// type Turn struct {
// 	P      Player
// 	B      Board
// 	W      Worker
// 	MoveTo Pos
// 	AddTo  Pos
// }

// //Takes in a player, a board, and a position. The position represents the position that the
// //given player would like to create a new worker on the given board.
// type PlaceWorker struct {
// 	P       Player
// 	B       Board
// 	PlaceAt Pos
// }

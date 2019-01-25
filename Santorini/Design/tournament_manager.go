package design

/*
  A Tournament Manager should:
    - accept new Users (name, strategy)
    - Start a new round-robin tournament (between Users already accepted)
      - Start new games between two Users (delegate to Referee)
      - record the winner of a game (a user)
      - record the loser of a game (a user)
    - Kick out Users (for breaking rules, give all cheater's wins to their opponent)
    - Report the results of a Game (run by a Referee)
    - Report the results of a Tournament
    - Connect an outside Observer to a running Game

  Beginning, middle, end:
  A Tournament, once initialized, will start accepting Users until it reaches
  the maximum number of Users acceptable or the timeout (in seconds) lapses.

  At this point, it will start pairing Users and creating new Referees to administer
  games between those users until each user has played with each other user in the
  tournament (if a User cheats in a game, they are kicked out and are no longer
  included in the set of potential Users to play against). If a User isn't done
  with a game, the manager (if there are no other unmatched opponents available)
  will wait for that game to end before the manager starts a game between them
  and another User.

  Once a winner has been determined from the Tournament (with tiebreaker games being
  run if, after each User has played each other, there is more than one player tied
  for the most wins), that User is reported to whoever started the Tournament (top-
  level server?), and the Tournament will end.
*/
type IManager interface {
	// Start listening for New Users, and once time has elapsed or the requisite
	// number of Users have entered, calls `runTournament` and returns the winner.
	// (Interface)
	Start(maxUsers, timeout int) User

	// Enter a new user into the next Tournament to be run, or deny that User if their
	// name is the empty string, or is already in use.
	// NOTE Will not make any changes if the Tournament has started
	// (Interface)
	NewUser(name string, strategy Strategy)

	// Kick a User for cheating within a game, and transfer their wins to
	// their opponent.
	// NOTE mutate tournament Users state by removing the cheater and updating
	//       the opponent's recorded wins
	// NOTE mutate tournament by disconnecting Observers watching the cheater
	// (implementation)
	kickUser(cheater, opponent User)

	// Run a Tournament between all registered users, and then report the winner
	// NOTE: prevent other Users from joining the Tournament
	// NOTE: The winner of the Tournament is whoever has the most wins at the end
	//       (tiebreaker games will be performed if there is more than User tied
	//       for the most wins)
	// (implementation)
	runTournament() User

	// Start a single game between two Users in a thread, and record the winner of
	// the game through `recordResults` when it is finished.
	// NOTE mutate tournament Users state via `KickUser`
	// NOTE mutate tournament Users state (update played opponents for each User)
	// NOTE mutate tournament by disconnecting Observers watching the cheater
	// (implementation)
	runGame(user1, user2 User) User

	// Mutate the tournament's Users state, and update the winner and loser
	// win/loss counts
	// (implementation)
	recordResults(winner, loser User)

	// Accept an observer on the given User, and attach it to the given user's
	// Game if they are playing any
	// (Interface)
	AcceptObserver(o IObserver, u User)

	// Disconnect an Observer, and detach it from the Referee (if any) they are
	// reading from.
	// (Interface)
	DisconnectObserver(o IObserver)
}

// A user represents the top-level player of multiple games
// A user has a name and a strategy, which is handed to the Referee to
// enact the playing of one or multiple games (based on the referee)
// Likewise, a user has a game history attached, representing their win/loss count
// and previously-played games
// NOTE A User can only play one game at a time
type User struct {
	// The User's name (unique within the Tournament)
	Name string

	// The Strategy a User is using (networked, AI, etc.)
	Strategy Strategy

	// A User's game history (wins, losses, matchups)
	History UserRecord

	// The index of the referee currently playing a game with this User
	CurrentGame int
}

// A Record is the stateful representation of a User's previous wins, losses,
// and previous Users faced (stored as list of names)
type UserRecord struct {
	// How many games the user has won
	Wins int

	// How many games the user has lost
	Losses int

	// Previous game results
	Opponents []GameResult
}

// A GameResult holds descriptive data about a Game's endstate:
// - who won
// - who lost
// - why the game ended
type GameResult struct {
	// The winning user's name
	Winner string

	// The losing user's name
	Loser string

	// A unique code indicating why a game ended
	// One of:
	// - Normal/Valid End (one user had a worker reach height 3)
	// - Rule broken (someone broke the rules and should be kicked out)
	EndCode int
}

// The manager maintains User state, Observers on different Users, and the
// different ongoing Games being Refereed
type Manager struct {
	// Whether the Tournament is running, or waiting for more Players
	Running bool

	// The Users within a running Tournament
	Users []User

	// A map of User names to the observers watching that User's games (as they occur)
	Observers map[string][]IObserver

	// The Referees administering games within this Tournament
	RunningGames []Referee
}

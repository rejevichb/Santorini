package tournament

import (
	"strings"

	ref "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Referee"
	sandbox "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Sandbox"
	cfg "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Tournament/Config"
	result "github.com/CS4500-F18/dare-rebr/Santorini/Common/Tournament"
	lib "github.com/CS4500-F18/dare-rebr/Santorini/Lib"
	obs "github.com/CS4500-F18/dare-rebr/Santorini/Observer"
)

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
*/
type IManager interface {
	//Runs the Tournament and reports both cheaters and the game results
	RunWithConfig(t cfg.TournamentConfig) result.TournamentResult
}

// A potential game to be run between two users
type userPair struct {
	player1, player2 *user
}

//The manager maintains User state, Observers on different Users, and the
//different ongoing Matches being Refereed
type manager struct {
	//How many games to play between each pair of opponents
	gamesPerRound int

	//The Map of taken names
	existingNames map[string]bool

	//The Users within a running Tournament
	Users []user

	//Completed Matches within this Tournament
	Matches []result.MatchResult

	//List of Observers watching games
	Observers []obs.IObserver

	//Users who have misbehaved during a game
	Excluded []user
}

//Return a new Tournament Manager with the given configuration values
func NewManager(games int) *manager {
	return &manager{
		gamesPerRound: games,
		Users:         make([]user, 0),
		existingNames: make(map[string]bool),
		Matches:       make([]result.MatchResult, 0),
		Observers:     make([]obs.IObserver, 0),
		Excluded:      make([]user, 0),
	}
}

//Load players and observers from a configuration
func (m *manager) RunWithConfig(c cfg.TournamentConfig) result.TournamentResult {
	players, observers := c.GenerateComponents()

	for _, player := range players {
		m.acceptPlayer(player)
	}

	for _, observer := range observers {
		m.AttachObserver(observer)
	}

	return m.run()
}

//Run a Tournament between all Users added to the Tournament
func (m *manager) run() result.TournamentResult {
	potentialGames := generateTuples(m.Users)

	for _, pair := range potentialGames {
		if m.runnablePair(pair.UserA, pair.UserB) {
			m.runSeries(pair.UserA, pair.UserB)
		}
	}

	result := result.TournamentResult{
		Games:  m.Matches,
		Kicked: userNames(m.Excluded),
	}

	for _, user := range append(m.Users, m.Excluded...) {
		user.Conn.ReceiveTournamentResult(result)
	}

	return result
}

//A Pair of two users representing a potential game matchup
type UserPair struct {
	UserA, UserB user
}

//Generate all possible game pairs from the users in the tournament
func generateTuples(arr []user) []UserPair {
	pairs := make([]UserPair, 0)
	for idxA, uA := range arr {
		for idxB, uB := range arr {
			if idxB > idxA {
				pairs = append(pairs, UserPair{uA, uB})
			}
		}
	}

	return pairs
}

func (m *manager) runnablePair(a, b user) bool {
	return !m.userCheated(a) && !m.userCheated(b)
}

//Run a series of games between users A and B, knowing they both have not cheated
func (m *manager) runSeries(a, b user) {
	referee := ref.NewReferee(a.Name, a.Conn, b.Name, b.Conn)
	m.AttachObservers(referee)
	gameSet := referee.BestOf(m.gamesPerRound)
	lastGame := gameSet[len(gameSet)-1]

	if lastGame.BrokenRule {
		if lastGame.Loser == a.Name {
			m.handleCheater(a)
		} else {
			m.handleCheater(b)
		}
	}
	matchResult := result.NewMatchResult(lastGame.Winner, lastGame.Loser, lastGame.BrokenRule, gameSet)
	m.Matches = append(m.Matches, matchResult)

	m.DetachObservers(referee)
}

//Add an Observer to the tournament
func (m *manager) AttachObserver(o obs.IObserver) {
	m.Observers = append(m.Observers, o)
}

//Remove an Observer from the tournament
func (m *manager) DetachObserver(target obs.IObserver) {
	for idx, o := range m.Observers {
		if o.Name() == target.Name() {
			m.Observers = append(m.Observers[:idx], m.Observers[idx+1:]...)
		}
	}
}

//Connect all Observers to a Referee
func (m *manager) AttachObservers(r ref.IReferee) {
	for _, observer := range m.Observers {
		r.AttachObserver(observer)
	}
}

//Remove all Observers from a Referee
func (m *manager) DetachObservers(r ref.IReferee) {
	for _, observer := range m.Observers {
		r.DetachObserver(observer)
	}
}

//Add users to the tournament, resolving invalid names or name conflicts
//by assigning the faulty-named player a new name of "abc...xyzabc..."
func (m *manager) acceptPlayer(player sandbox.WrappedPlayer) {
	originalName, _ := player.Name()
	lowercase := strings.ToLower(originalName)

	if !m.validName(lowercase) || !m.addUnique(lowercase, player) {
		newName := ""

		for i := 0; true; i++ {
			newName = newName + string(lib.ALPHA[i%len(lib.ALPHA)]) // repeat a-z
			if m.addUnique(newName, player) {
				player.SetName(newName)
				return
			}
		}
	}
}

//Is the given player name valid?
// A name is valid if it's composed of exclusively lowercase characterss
// and is non-empty
func (m *manager) validName(name string) bool {
	return name != "" && lib.IsLowercase(name)
}

//Add the given player with the given name if the player
//has a unique name, or return false if not unique
func (m *manager) addUnique(name string, player sandbox.WrappedPlayer) bool {
	if _, ok := m.existingNames[name]; !ok {
		m.Users = append(m.Users, NewUser(name, player))
		m.existingNames[name] = true
		return true
	}
	return false
}

//Remove a rule-violating player from the manager's users and match history
func (m *manager) handleCheater(target user) {
	m.removeCheater(target)
	m.removeFromHistory(target)
}

//Remove the target from the list of Users within the Tournament
func (m *manager) removeCheater(target user) {
	for idx, u := range m.Users {
		if target.eq(u) {
			m.Users = append(m.Users[:idx], m.Users[idx+1:]...)
			delete(m.existingNames, target.Name)
			m.Excluded = append(m.Excluded, target)
		}
	}
}

//Remove this user from the history of completed matches
func (m *manager) removeFromHistory(target user) {
	for idx, match := range m.Matches {
		if match.Winner == target.Name {
			if match.RuleBroken {
				m.Matches = append(m.Matches[:idx], m.Matches[idx+1:]...)
			} else {
				m.Matches[idx].RuleBroken = true
				m.Matches[idx].Winner = m.Matches[idx].Loser
				m.Matches[idx].Loser = target.Name
			}
			break
		} else if match.Loser == target.Name {
			m.Matches[idx].RuleBroken = true
			break
		}
		// if you're not the winner or loser, don't change anything
	}
}

//Did the given user cheat?
//(is the given user in the list of kicked players)
func (m *manager) userCheated(u user) bool {
	for _, user := range m.Excluded {
		if user.eq(u) {
			return true
		}
	}
	return false
}

/********* USER METHODS *********/

//A user is a named player
type user struct {
	//The User's name (unique within the Tournament)
	Name string

	//The Strategy a User is using (networked, AI, etc.)
	Conn sandbox.WrappedPlayer
}

//create a new user from a name and wrapped player
func NewUser(name string, player sandbox.WrappedPlayer) user {
	return user{
		Name: name,
		Conn: player,
	}
}

//return whether two users are equal
//(two users are equal if their names are equal)
func (u user) eq(other user) bool {
	return u.Name == other.Name
}

//fold over users and return their names
func userNames(users []user) []string {
	names := make([]string, 0)
	for _, u := range users {
		names = append(names, u.Name)
	}
	return names
}

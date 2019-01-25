package remote

import (
	"encoding/json"
	"net"
	"time"

	sandbox "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Sandbox"
	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	data "github.com/CS4500-F18/dare-rebr/Santorini/Common/JSON"
	iplayer "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	result "github.com/CS4500-F18/dare-rebr/Santorini/Common/Tournament"
	lib "github.com/CS4500-F18/dare-rebr/Santorini/Lib"
)

//Component representing a Player's remote connection
//that will send JSON data to a player
//NOTE implements IPlayer interface
type ProxyPlayer struct {
	//The TCP connection to the player
	conn net.Conn

	// The timeout for TCP communication
	timeout int
}

// Create a proxy Player over the given connection with the
// given timeout in milliseconds
func NewProxyPlayer(conn net.Conn, timeout int) ProxyPlayer {
	return ProxyPlayer{conn: conn, timeout: timeout}
}

// Get the name of this Player
func (p ProxyPlayer) Name() (string, error) {
	_, decoder := lib.JSONStreams(p.conn)

	var name string
	err := p.getWithTimeout(p.timeout, &name, decoder)
	if err != nil {
		return "", err
	}
	return name, nil
}

//SetName tells a Player what their assigned Name is, given
//that it may be non-unique
func (p ProxyPlayer) SetName(newName string) error {
	encoder, _ := lib.JSONStreams(p.conn)
	rename := data.Rename{Name: newName}
	return encoder.Encode(rename)
}

//PlaceWorker gets the location to place your next worker
func (p ProxyPlayer) PlaceWorker(b board.IBoard) (board.Pos, error) {
	encoder, decoder := lib.JSONStreams(p.conn)

	// Send Workers
	workers := b.Workers()
	if err := encoder.Encode(workers); err != nil {
		return board.Pos{X: -1, Y: -1}, err
	} else {
		// Listen for response
		var pos board.Pos
		err := p.getWithTimeout(p.timeout, &pos, decoder)
		return pos, err
	}
}

//NextTurn gets the next turn, including which worker ID to act on
func (p ProxyPlayer) NextTurn(b board.IBoard) (iplayer.Turn, error) {
	encoder, decoder := lib.JSONStreams(p.conn)

	invalid := iplayer.Turn{-1, board.Pos{X: -1, Y: -1}, board.Pos{X: -1, Y: -1}}

	if err := encoder.Encode(b); err != nil {
		return iplayer.Turn{}, err
	} else {
		// Listen for response
		var iface interface{}
		err := p.getWithTimeout(p.timeout, &iface, decoder)
		if err != nil {
			return invalid, err
		}
		buf, _ := json.Marshal(iface)

		if turn, err := p.TryGiveUp(buf); err == nil {
			return turn, nil
		} else if turn, err := p.TryMoveTurn(b, buf); err == nil {
			return turn, nil
		} else if turn, err := p.TryMoveBuildTurn(b, buf); err == nil {
			return turn, nil
		} else {
			return invalid, err
		}
	}
}

// Receive an attempted give-up from a Player
func (p ProxyPlayer) TryGiveUp(buf []byte) (iplayer.Turn, error) {
	turn := iplayer.Turn{-1, board.Pos{-1, -1}, board.Pos{-1, -1}}

	var str string
	err := json.Unmarshal(buf, &str)
	if err != nil {
		return turn, err
	}
	return turn, nil
}

// Attempt to decode into a move/build turn
func (p ProxyPlayer) TryMoveBuildTurn(b board.IBoard, buf []byte) (iplayer.Turn, error) {
	turn := iplayer.Turn{-1, board.Pos{-1, -1}, board.Pos{-1, -1}}

	var mbt data.MoveBuildTurn
	err := json.Unmarshal(buf, &mbt)
	if err != nil {
		return turn, err
	}

	turn, err = mbt.ToTurn(b)
	return turn, err
}

// Attempt to decode into a solely-move turn
func (p ProxyPlayer) TryMoveTurn(b board.IBoard, buf []byte) (iplayer.Turn, error) {
	turn := iplayer.Turn{-1, board.Pos{X: -1, Y: -1}, board.Pos{X: -1, Y: -1}}

	var mt data.MoveTurn
	err := json.Unmarshal(buf, &mt)
	if err != nil {
		return turn, err
	}

	turn, err = mt.ToTurn(b)
	return turn, err
}

//Opponent informs the Player of the opponent they are playing
func (p ProxyPlayer) SetOpponent(name string) error {
	encoder, _ := lib.JSONStreams(p.conn)
	return encoder.Encode(name)
}

// Get the results of a tournament
func (p ProxyPlayer) ReceiveTournamentResult(result result.TournamentResult) error {
	defer p.conn.Close()
	encoder, _ := lib.JSONStreams(p.conn)
	err := encoder.Encode(result)
	return err
}

// Generic accessor with a timeout on the TCP connection
func (p ProxyPlayer) getWithTimeout(timeout int, target interface{}, decoder *json.Decoder) error {
	duration := time.Duration(timeout) * sandbox.TIMEOUT_UNIT
	p.conn.SetDeadline(time.Now().Add(duration)) // Stop I/O after duration
	defer p.conn.SetDeadline(time.Time{})        // Zero value -- no timeout

	err := decoder.Decode(target)
	return err
}

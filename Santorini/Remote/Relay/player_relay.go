package remote

import (
	"encoding/json"
	"net"
	"strconv"

	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	data "github.com/CS4500-F18/dare-rebr/Santorini/Common/JSON"
	iplayer "github.com/CS4500-F18/dare-rebr/Santorini/Common/Player"
	lib "github.com/CS4500-F18/dare-rebr/Santorini/Lib"
)

/*
  The Player Relay is the component that acts as the TCP communication portion
  on the client-side
*/

// will need to:
// - connect initixally
// - listen for new requests
// - respond to requests in JSON
// - wrap a Player

type PlayerRelay struct {
	player iplayer.IPlayer
}

type IRelay interface {
	Connect(host string, port int, err chan bool) error
	Register(*json.Encoder) error
	SendPlacement(*json.Encoder, board.IBoard) error
	SendTurn(*json.Encoder, board.IBoard) error
	ListenAndRespond(net.Conn, chan bool)
}

//create an unconnected PlayerRelay given a player.
func NewPlayerRelay(p iplayer.IPlayer) PlayerRelay {
	return PlayerRelay{player: p}
}

//attempts to connect to the IP and port, returns an error if the connection cannot be established.
func (r PlayerRelay) Connect(host string, port int, done chan bool) error {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	} else {
		go r.ListenAndRespond(conn, done)
		return nil
	}
}

//register this remote player with the server
func (r PlayerRelay) Register(encoder *json.Encoder) error {
	if err := encoder.Encode(r.player.Name()); err != nil {
		return err
	} else {
		return nil
	}
}

//sends a placement action using a the strategy of the playerRelay's wrapped IPayer
func (r PlayerRelay) SendPlacement(encoder *json.Encoder, b board.IBoard) error {
	placement := r.player.PlaceWorker(b)
	if err := encoder.Encode(placement); err != nil {
		return err
	} else {
		return nil
	}
}

// sends a turn action over TCP using the strategy of the playerRelay's wrapped IPayer
func (r PlayerRelay) SendTurn(encoder *json.Encoder, b board.IBoard) error {
	turn := r.player.NextTurn(b)

	turnJSON := data.MoveBuildFromTurn(r.player.Name(), b, turn)
	if err := encoder.Encode(turnJSON); err != nil {
		return err
	} else {
		return nil
	}
}

//switch over datatype based on what response we get.
/*Can get:
  - String  (this player's new name)
  - Placements (Worker placements on the board to use in placement)
  - Board (Board to enact a Turn on)
  - Tournament results (Game's over, we're done here folks)
*/
func (r PlayerRelay) ListenAndRespond(conn net.Conn, done chan bool) {
	encoder, decoder := lib.JSONStreams(conn)
	defer conn.Close()
	r.Register(encoder)

	for {
		if !decoder.More() {
			continue
		} else { // We have something to work with
			buf, _ := r.GetData(decoder)

			if err := r.TryRename(buf); err == nil {
				continue
			} else if err := r.TryOpponent(buf); err == nil {
				continue
			} else if err := r.TryPlacement(encoder, buf); err == nil {
				continue
			} else if err := r.TryBoard(encoder, buf); err == nil {
				continue
			} else if err := r.TryResult(buf); err == nil {
				done <- true
				return
			}
		}
	}
}

func (r PlayerRelay) GetData(decoder *json.Decoder) ([]byte, error) {
	// actually get data from the connection
	var obj interface{}
	decoder.Decode(&obj)
	return json.Marshal(obj)
}

//tries to marshal the given bytes data into a name request,
//if successful, set name for this player.
func (r PlayerRelay) TryRename(buf []byte) error {
	// New name for this player
	var rename data.Rename
	err := json.Unmarshal(buf, &rename)
	if err == nil {
		r.player.SetName(rename.Name)
	}
	return err
}

//tries to marshal the given bytes data into a set opponent request,
//if successful, set name for this player.
func (r PlayerRelay) TryOpponent(buf []byte) error {
	// Opponent's name
	var s string
	err := json.Unmarshal(buf, &s)
	if err == nil {
		r.player.SetOpponent(s)
	}
	return err
}

//tries to marshal the given bytes data into a list of worker placements,
//if successful, send the next placement for this player.
func (r PlayerRelay) TryPlacement(encoder *json.Encoder, buf []byte) error {
	// Worker placements
	var workers []board.Worker
	err := json.Unmarshal(buf, &workers)
	if err == nil {
		interfaces := make([]board.IWorker, 0)
		for _, worker := range workers {
			interfaces = append(interfaces, board.IWorker(worker))
		}
		withWorkers := board.BoardWithWorkers(interfaces)
		err = r.SendPlacement(encoder, withWorkers)
		if err != nil {
			return err
		}
	}
	return err
}

//tries to marshal the given bytes data into a board data strcture,
//if that works, send the next player turn.
func (r PlayerRelay) TryBoard(encoder *json.Encoder, buf []byte) error {
	// Board to enact a Turn on
	b := board.BaseBoard()
	err := json.Unmarshal(buf, &b)
	if err == nil {
		err = r.SendTurn(encoder, b)
		if err != nil {
			return err
		}
	}
	return err
}

func (r PlayerRelay) TryResult(buf []byte) error {
	var results []interface{}
	err := json.Unmarshal(buf, &results)
	return err
}

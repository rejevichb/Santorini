package remote

import (
	"bytes"
	"net"
	"testing"

	"github.com/CS4500-F18/dare-rebr/Santorini/Lib"
	"github.com/CS4500-F18/dare-rebr/Santorini/Player/Client"
)

//placement json

//turn json

//set name

//set opponent

//test name change, opponent name change
func TestConn_TryPlacement(t *testing.T) {
	go func(t *testing.T) {
		// start a server on a port
		l, err := net.Listen("tcp", ":8080")
		if err != nil {
			t.Fatal(err)
		}
		// when you get a connection
		conn, err := l.Accept()
		if err != nil {
			return
		}

		defer conn.Close()

		// encode an empty array into bytes
		enc, _ := lib.JSONStreams(conn)
		enc.Encode([]interface{}{})
	}(t)

	// create a relay with a player
	var b_buf []byte
	rw := bytes.NewBuffer(b_buf)
	enc, _ := lib.JSONStreams(rw)

	//make player and player relay
	player := client.ValidPlayer("player1")
	relay := NewPlayerRelay(player)

	address := net.JoinHostPort("localhost", "8080")
	conn, err := net.Dial("tcp", address)
	// err better not be there
	if err != nil {
		t.Errorf("error on client net.Dial")
	}
	_, dec := lib.JSONStreams(conn)

	// call GetData and make sure you got stuff
	data, err := relay.GetData(dec)

	// call TryPlacement with an encoder you can read from
	e := relay.TryPlacement(enc, data)
	if e != nil {
		t.Errorf("try placement failed")
	}

}

// func TestConn_TryOpponent(t *testing.T) {
// 	go func(t *testing.T) {
// 		// start a server on a port
// 		l, err := net.Listen("tcp", ":8080")
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		// when you get a connection
// 		conn, err := l.Accept()
// 		if err != nil {
// 			return
// 		}
//
// 		defer conn.Close()
//
// 		// encode an empty array into bytes
// 		enc, _ := lib.JSONStreams(conn)
// 		enc.Encode([]interface{}{"oop"})
// 	}(t)
//
// 	//make player and player relay
// 	player := client.ValidPlayer("player1")
// 	relay := NewPlayerRelay(player)
//
// 	address := net.JoinHostPort("localhost", "8080")
// 	conn, err := net.Dial("tcp", address)
// 	if err != nil {
// 		t.Errorf("error on client net.Dial")
// 	}
//
// 	_, dec := lib.JSONStreams(conn)
//
// 	// call GetData and make sure you got stuff
// 	data, err := relay.GetData(dec)
//
// 	// call TryPlacement with an encoder you can read from
// 	e := relay.TryOpponent(data)
// 	if e != nil {
// 		t.Errorf("try placement failed")
// 	}
//
// 	if relay.player.Opponent() != "oop" {
// 		t.Errorf("player relay object's opponent's name is not correct")
// 	}
// }

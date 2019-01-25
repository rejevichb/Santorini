package config

import (
	"net"
	"strconv"
	"time"

	sandbox "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Sandbox"
	obs "github.com/CS4500-F18/dare-rebr/Santorini/Observer"
	remote "github.com/CS4500-F18/dare-rebr/Santorini/Remote/Player"
)

// TourneyConfiguration for a Tournament
type RemoteConfig struct {
	//min players
	playerCount int

	//port to listen on
	port int

	//time limit to accept new players (milliseconds)
	acceptLimit int

	//timeout in milliseconds for underlying players
	timeout int
}

func NewRemoteConfig(players, port, limit, timeout int) RemoteConfig {
	return RemoteConfig{players, port, limit, timeout}
}

// Create Tournament-usable pieces from a TourneyConfiguration
// Wait for Players til you hit the time limit, re-run if below min
// NOTE: On re-run, keep previous players until you hit the minimum limit.
func (c RemoteConfig) GenerateComponents() ([]sandbox.WrappedPlayer, []obs.IObserver) {
	serv, err := net.Listen("tcp", ":"+strconv.Itoa(c.port))
	if err != nil {
		panic(err)
	}

	connections := make(chan net.Conn)

	go acceptConnections(serv, connections)

	proxies := make([]sandbox.WrappedPlayer, 0)
	timer := time.After(time.Duration(c.acceptLimit) * time.Second)
	observers := make([]obs.IObserver, 0)

	for {
		select {
		case <-timer:
			if len(proxies) >= c.playerCount {
				// After the timer has lapsed, if you have the minimal player count:
				// Return everyone so far
				return proxies, observers
			} else {
				timer = time.After(time.Duration(c.acceptLimit) * time.Second)
			}

		case conn := <-connections:
			player := remote.NewProxyPlayer(conn, c.timeout)
			proxies = append(proxies, player)
		}
	}
}

func acceptConnections(l net.Listener, c chan net.Conn) {
	for {
		conn, err := l.Accept()
		if err == nil {
			c <- conn
			time.Sleep(1 * time.Second)
		}
	}
}

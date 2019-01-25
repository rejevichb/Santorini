package reader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	static "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Tournament/Config"
)

func RunTest(r io.Reader) {
	decoder := json.NewDecoder(r)

	var config static.StaticConfig
	decoder.Decode(&config)

	relays, _ := config.ClientRelays()

	done := make(chan bool)
	started := 0
	for _, relay := range relays {
		err := relay.Connect(config.IP, config.Port, done)
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to %s:%v", config.IP, config.Port))
		} else {
			started++
		}
	}

	if started == 0 {
		return
	}

	for {
		select {
		case <-done:
			os.Exit(0) // This is a hack
			started--
			if started <= 0 {
				return
			}
		}
	}
}

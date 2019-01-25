package main

import (
	"fmt"
	"os"

	"github.com/CS4500-F18/dare-rebr/4/Aux/sheetclient"
)

func main() {
	// Server args
	if args := os.Args; len(args) < 3 {
		fmt.Println("Error: must provide a server address as an argument.\nEx: `client google.com`")
		return
	} else {
		serverAddr := args[1]
		name := args[2]
		// Delegate to client
		fmt.Println("Setting up client pointed at:", serverAddr)
		sheetclient.Run(serverAddr, name)
	}

}

// Package main provides a single client to encode a data stream
// into a specified format, and output the stream into reverse
// chronological order to two defined IO streams.
package main

import (
	"os"
	"spread_client"
)

func main() {
	spread_client.Read_Stream(os.Stdin, os.Stdout)
}

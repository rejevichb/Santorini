package main

import (
	"os"

	reader "github.com/CS4500-F18/dare-rebr/Santorini/13_Server/Aux/reader"
)

func main() {
	reader.RunTest(os.Stdin, os.Stdout)
}

package main

import (
	"os"

	"github.com/CS4500-F18/dare-rebr/Santorini/7_Harness/Aux/reader"
)

func main() {
	r := os.Stdin
	w := os.Stdout

	reader.RunTest(r, w)
}

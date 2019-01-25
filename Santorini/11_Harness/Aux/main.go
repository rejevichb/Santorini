package main

import (
	"os"

	reader "github.com/CS4500-F18/dare-rebr/Santorini/11_Harness/aux/reader"
)

func main() {
	reader.RunTest(os.Stdin, os.Stdout)
}

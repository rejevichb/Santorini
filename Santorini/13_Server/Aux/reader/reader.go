package reader

import (
	"encoding/json"
	"io"

	server "github.com/CS4500-F18/dare-rebr/Santorini/Remote/Server"
)

func RunTest(r io.Reader, w io.Writer) {
	decoder := json.NewDecoder(r)
	encoder := json.NewEncoder(w)

	var configuration server.ServerConfig
	err := decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}

	serv := server.NewServer()

	results := serv.Start(configuration)
	if len(results) == 1 {
		encoder.Encode(results[0])
	}
}

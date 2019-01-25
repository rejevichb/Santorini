package lib

import (
	"encoding/json"
	"io"
)

func JSONStreams(c io.ReadWriter) (*json.Encoder, *json.Decoder) {
	encoder := json.NewEncoder(c)
	decoder := json.NewDecoder(c)
	return encoder, decoder
}

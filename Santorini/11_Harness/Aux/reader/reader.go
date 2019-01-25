package reader

import (
	"encoding/json"
	"io"

	tourny "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Tournament"
	cfg "github.com/CS4500-F18/dare-rebr/Santorini/Admin/Tournament/Config"
)

func RunTest(r io.Reader, w io.Writer) {
	decoder := json.NewDecoder(r)
	encoder := json.NewEncoder(w)

	var configuration cfg.StaticConfig
	decoder.Decode(&configuration)

	tournament := tourny.NewManager(3)
	result := tournament.RunWithConfig(configuration)

	encoder.Encode(result.Kicked)
	encoder.Encode(result.Games)
}

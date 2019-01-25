package reader

import (
	"encoding/json"
	"io"

	common "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	spec "github.com/CS4500-F18/dare-rebr/Santorini/Common/JSON"
)

func RunTest(reader io.Reader, writer io.Writer) {
	decoder := json.NewDecoder(reader)
	encoder := json.NewEncoder(writer)

	bStruct := common.BaseBoard()
	err := decoder.Decode(&bStruct)
	if err != nil {
		panic(err)
	}
	b := common.IBoard(bStruct)

	var output interface{}
	var command spec.Command

	for decoder.More() {
		if err = decoder.Decode(&command); err != nil {
			panic(err)
		} else {
			b, output = command.Execute(b)
			encoder.Encode(output)
		}
	}

	return
}

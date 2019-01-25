package reader

//Responsible for handling reading in custom JSON values from test harness

import (
	"encoding/json"
	"io"

	common "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	spec "github.com/CS4500-F18/dare-rebr/Santorini/Common/JSON"
)

const BoardSize = 6

var selectedWorker = ""
var afterMovePos common.Pos

func RunTest(reader io.Reader, writer io.Writer) {
	decoder := json.NewDecoder(reader)
	encoder := json.NewEncoder(writer)

	bStruct := common.BaseBoard()
	err := decoder.Decode(&bStruct)
	if err != nil {
		panic(err)
	}
	b := common.IBoard(bStruct)

	//Look for move command
	var command spec.Command
	decoder.Decode(&command)

	var output interface{}
	b, output = command.Execute(b)
	if output == spec.NO {
		encoder.Encode(spec.NO)
		return
	}

	//Look for +build command
	if decoder.More() {
		var shortComm spec.ShortCommand
		decoder.Decode(&shortComm)
		command = spec.Command{Type: shortComm.Type, Worker: command.Worker, Direction: shortComm.Direction}

		b, output = command.Execute(b)
		if output == spec.NO {
			encoder.Encode(spec.NO)
			return
		}
	}

	encoder.Encode(spec.YES)
}

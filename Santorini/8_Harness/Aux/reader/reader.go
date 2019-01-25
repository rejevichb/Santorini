package reader

//Responsible for handling reading in custom JSON values from test harness

import (
	"encoding/json"
	"io"

	common "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	spec "github.com/CS4500-F18/dare-rebr/Santorini/Common/JSON"
	strategy "github.com/CS4500-F18/dare-rebr/Santorini/Player/Strategy"
)

func RunTest(reader io.Reader, writer io.Writer) {
	decoder := json.NewDecoder(reader)
	encoder := json.NewEncoder(writer)

	var stratPlayer string
	decoder.Decode(&stratPlayer)

	bStruct := common.BaseBoard()
	err := decoder.Decode(&bStruct)
	if err != nil {
		panic(err)
	}
	b := common.IBoard(bStruct)

	var roundsAhead float64
	decoder.Decode(&roundsAhead)

	var command spec.Command
	var short spec.ShortCommand
	if decoder.More() {
		decoder.Decode(&command)
		if decoder.More() {
			decoder.Decode(&short)
		}
	}

	// existing move
	if (command != spec.Command{}) {
		b, _ = command.Execute(b)
	}

	if (short != spec.ShortCommand{}) {
		command = spec.Command{Type: short.Type, Worker: command.Worker, Direction: short.Direction}
		b, _ = command.Execute(b)
	}

	var otherPlayer string
	for _, name := range b.Players() {
		if name != stratPlayer {
			otherPlayer = name
		}

	}

	_, err = strategy.SurvivingTurn(b, stratPlayer, otherPlayer, int(roundsAhead))

	if err != nil {
		encoder.Encode(spec.NO)
	} else {
		encoder.Encode(spec.YES)
	}

	return
}

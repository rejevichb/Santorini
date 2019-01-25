package observer

import (
	"encoding/json"
	"fmt"
	"io"

	board "github.com/CS4500-F18/dare-rebr/Santorini/Common/Board"
	output "github.com/CS4500-F18/dare-rebr/Santorini/Common/JSON"
	rules "github.com/CS4500-F18/dare-rebr/Santorini/Common/Rules"
)

type IObserver interface {
	// This Observer's name
	Name() string

	//Receive an updated board
	ReceiveBoard(b board.IBoard)

	//Receive the final Move that wins a Game
	ReceiveWinningMove(t output.MoveJSON)

	//Receive a full Turn performed
	ReceiveTurn(t output.MoveBuildJSON)

	//Receive an endgame state
	ReceiveEndgame(end rules.GameResult)
}

type JsonObserver struct {
	name   string
	output io.Writer
}

func NewJSONObserver(name string, stream io.Writer) IObserver {
	return JsonObserver{name, stream}
}

func (o JsonObserver) Name() string {
	return o.name
}

//Print the board state
//Receive an updated board
func (o JsonObserver) ReceiveBoard(b board.IBoard) {
	marshalled, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	str := string(marshalled)

	for idx := 0; idx < len(str)-1; idx++ {
		if string(str[idx])+string(str[idx+1]) == "]," {
			newSet := "],\n"
			str = str[:idx] + newSet + str[idx+2:]
			idx += len(newSet)
		}
	}
	o.output.Write([]byte(str + "\n"))
}

//Receive the final Move that wins a Game
func (o JsonObserver) ReceiveWinningMove(move output.MoveJSON) {
	jsonMove, _ := json.Marshal(move.MoveDir)
	jsonWorker, _ := json.Marshal(move.WorkerName)

	str := fmt.Sprintf("[%s, %s]\n", jsonWorker, string(jsonMove))
	o.output.Write([]byte(str))
}

//Receive a full Turn performed
func (o JsonObserver) ReceiveTurn(turn output.MoveBuildJSON) {
	jsonMove, err := json.Marshal(turn.MoveDir)
	if err != nil {
		panic(err)
	}
	jsonBuild, _ := json.Marshal(turn.BuildDir)
	jsonWorker, _ := json.Marshal(turn.WorkerName)

	str := fmt.Sprintf("[%s, %s, %s]\n", jsonWorker, string(jsonMove), string(jsonBuild))
	o.output.Write([]byte(str))
}

//Receive an endgame state
func (o JsonObserver) ReceiveEndgame(end rules.GameResult) {
	if end.Reason == rules.RULE_BROKEN_MSG {
		o.output.Write([]byte("\"" + end.Loser + " Lost: " + end.Reason + "\"\n"))
	} else {
		o.output.Write([]byte("\"" + end.Winner + " Won\"\n"))
	}
}

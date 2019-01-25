// Package spread_client provides a function to
// decode input and encode output from two streams.
package spread_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"spreadsheet_adap"
)

func Read_Stream(rdr io.Reader, wrtr io.Writer) error {
	spreadsheets := make(map[string]spreadsheet_adap.Spreadsheet)

	dec := json.NewDecoder(rdr)

	for {
		var currJSON interface{}

		err := dec.Decode(&currJSON)

		if err != nil {
			if err != io.EOF {
				return errors.New("Error while attempting to read stream.")
			}

			break
		}

		execute(&wrtr, &spreadsheets, currJSON)
	}
	return nil
}

func execute(out *io.Writer, spreads *map[string]spreadsheet_adap.Spreadsheet, json interface{}) {
	commands := initCommands()

	jArr, ok := json.([]interface{})

	if !ok || len(jArr) <= 0 {
		return
	}

	commandName, ok := jArr[0].(string)

	if !ok {
		return
	}

	if cmd := commands[commandName]; cmd != nil {
		cmd(out, spreads, jArr)
	}

}

func initCommands() map[string]func(*io.Writer, *map[string]spreadsheet_adap.Spreadsheet, []interface{}) {
	cmds := make(map[string]func(*io.Writer, *map[string]spreadsheet_adap.Spreadsheet, []interface{}))
	cmds["set"] = setCommand
	cmds["at"] = atCommand
	cmds["sheet"] = sheetCommand

	return cmds
}

func getSafeInt(n interface{}) (int, bool) {
	x, ok := n.(float64)

	if ok {
		return int(x), x == float64(int64(x))
	}

	return -1, false
}
func setCommand(out *io.Writer, spreads *map[string]spreadsheet_adap.Spreadsheet, json []interface{}) {

	if len(json) != 5 {
		return
	}

	sheatName, strConv := json[1].(string)
	xCoord, getX := getSafeInt(json[2])
	yCoord, getY := getSafeInt(json[3])
	jf, err := translateJF(json[4])

	if err != nil {
		return
	}

	if strConv && getX && getY && err == nil {
		sheat, gotSheet := (*spreads)[sheatName]
		if gotSheet {
			sheat.Write(jf, yCoord, xCoord)
		}
	}

	return
}

func atCommand(out *io.Writer, spreads *map[string]spreadsheet_adap.Spreadsheet, json []interface{}) {
	if len(json) != 4 {
		return
	}

	shName, trans := json[1].(string)
	spread, found := (*spreads)[shName]
	xCoord, getX := json[2].(float64)
	yCoord, getY := json[3].(float64)

	if trans && found && getX && getY && xCoord == float64(int64(xCoord)) && yCoord == float64(int64(yCoord)) {
		val, err := spread.Value(int(yCoord), int(xCoord))
		if err == nil {
			fmt.Fprintf(*out, "%f\n", val)
		} else {
			if err.Error() == "Circular Reference" {
				fmt.Fprint(*out, "false\n")
			}
		}
	}

}

func sheetCommand(out *io.Writer, spreads *map[string]spreadsheet_adap.Spreadsheet, json []interface{}) {
	if len(json) != 3 {
		return
	}

	name, nameTrans := json[1].(string)

	linearBlock, succesfulLinear := json[2].([]interface{})
	rows := len(linearBlock)
	rowHolder := make([][]interface{}, rows)
	if rows == 0 {
		return
	}

	if !succesfulLinear {
		return
	}
	first, err := linearBlock[0].([]interface{})
	cols := len(first)

	if !err {
		return
	}

	for row := 0; row < rows; row++ {
		madeRow, yes := linearBlock[row].([]interface{})
		if !yes || len(madeRow) != cols {
			return
		}
		rowHolder[row] = madeRow
	}

	if nameTrans {
		transBlock, errBlock := convertSheet(rowHolder)
		if errBlock == nil {
			(*spreads)[name] = transBlock
		}
	}
}

func convertSheet(sheetJson [][]interface{}) (spreadsheet_adap.Spreadsheet, error) {
	rows := len(sheetJson)
	spreadSlice := make([][]spreadsheet_adap.Equation, rows)
	if rows <= 0 {
		return spreadsheet_adap.MakeSpreadsheet(spreadSlice), nil
	}

	cols := len(sheetJson[0])
	for y := 0; y < rows; y++ {
		if len(sheetJson[y]) != cols {
			return spreadsheet_adap.MakeBlankSpreadsheet(), errors.New("Non-rectangular input")
		}

		newRow := make([]spreadsheet_adap.Equation, cols)

		for x := 0; x < cols; x++ {
			trans, transErr := translateJF(sheetJson[y][x])

			if transErr != nil {
				return spreadsheet_adap.MakeBlankSpreadsheet(), transErr
			}

			newRow[x] = trans
		}

		spreadSlice[y] = newRow
	}

	return spreadsheet_adap.MakeSpreadsheet(spreadSlice), nil
}

func translateJF(json interface{}) (spreadsheet_adap.Equation, error) {
	switch v := json.(type) {
	case float64:
		return spreadsheet_adap.MakeConstant(v), nil
	case int:
		return spreadsheet_adap.MakeConstant(float64(v)), nil
	case []interface{}:
		return parseArr(v)
	default:
		return spreadsheet_adap.MakeBlankEq(), errors.New("Does not parse to valid JF.")
	}
}

func parseArr(jsonArr []interface{}) (spreadsheet_adap.Equation, error) {
	if len(jsonArr) != 3 {
		return spreadsheet_adap.MakeBlankEq(), errors.New("Does not parse to valid JF.")
	}

	if v, ok := jsonArr[0].(string); ok {
		fx, okX := jsonArr[1].(float64)
		fy, okY := jsonArr[2].(float64)
		if v == ">" && okX && okY && fx == float64(int64(fx)) && fy == float64(int64(fy)) {
			return spreadsheet_adap.MakeRef(int(fy), int(fx)), nil
		}
	}

	v, ok := jsonArr[1].(string)
	if ok {
		if v == "+" {
			eq1, ok1 := translateJF(jsonArr[0])
			eq2, ok2 := translateJF(jsonArr[2])
			if ok1 == nil && ok2 == nil {
				return spreadsheet_adap.MakeAdd(eq1, eq2), nil
			}
		}

		if v == "*" {
			eq1, ok1 := translateJF(jsonArr[0])
			eq2, ok2 := translateJF(jsonArr[2])

			if ok1 == nil && ok2 == nil {
				return spreadsheet_adap.MakeMult(eq1, eq2), nil
			}
		}
	}
	return spreadsheet_adap.MakeBlankEq(), errors.New("Does not parse to valid JF.")
}

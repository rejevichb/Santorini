package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"

	sheet "github.com/CS4500-F18/dare-rebr/3"
)

func MALFORMED_ERROR(i interface{}) error {
	return fmt.Errorf("JSON input not well-defined: %v of type %v", i, reflect.TypeOf(i))
}

func main() {
	sheets := make(map[string]sheet.Spreadsheet)
	decoder := json.NewDecoder(os.Stdin)

	var inputJSON interface{}

	for {
		err := decoder.Decode(&inputJSON)
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}

		switch inputJSON.(type) {
		case []interface{}:

			inputArray := inputJSON.([]interface{})
			if len(inputArray) < 3 {
				continue
			}

			command, ok := inputArray[0].(string)
			if !ok {
				continue
			}
			name, ok := inputArray[1].(string)
			if !ok {
				continue
			}

			funcArgs := inputArray[2:]

			if output, err := executeCommand(sheets, command, name, funcArgs); err != nil {
				fmt.Println("Execution error: ", err)
				continue
			} else if output != nil {
				fmt.Println("Output: ", output)
			}

		default:
			// Discard value
		}
	}
}

func executeCommand(sheets map[string]sheet.Spreadsheet, command, name string, args []interface{}) (*float64, error) {
	cellSheet, existing := sheets[name]

	switch command {
	case "sheet":
		return nil, executeCreate(sheets, name, args[0])

	case "set":
		if !existing {
			return nil, fmt.Errorf("Named sheet '%s' doesn't exist", name)
		}

		return nil, executeSet(cellSheet, args)

	case "at":
		if !existing {
			return nil, fmt.Errorf("Named sheet '%s' doesn't exist", name)
		}

		value, err := executeAt(cellSheet, args)

		return &value, err
	}

	return nil, fmt.Errorf("Invalid sheet command: '%s'", command)
}

/******* Execute each command with the args it wants *******/

func executeCreate(sheets map[string]sheet.Spreadsheet, name string, unparsedCells interface{}) error {
	fmt.Printf("Making sheet with name '%s' and cells %v\n", name, unparsedCells)

	_, ok := sheets[name]
	if ok {
		return fmt.Errorf("Spreadsheet already exists with name %s", name)
	}

	cells, validArray := unparsedCells.([]interface{})
	if !validArray {
		return MALFORMED_ERROR(cells)
	}

	cellSheet := make([][]sheet.Cell, 0)

	for x, col := range cells {
		cellSheet = append(cellSheet, make([]sheet.Cell, 0))
		colArr := col.([]interface{})

		for _, cell := range colArr {
			if parsedEq, err := parseEquation(cell); err != nil {
				return err
			} else {
				cellSheet[x] = append(cellSheet[x], sheet.CreateCell(parsedEq))
			}
		}
	}

	newSheet := createSheet(cellSheet)
	sheets[name] = newSheet

	return nil
}

func executeSet(cellSheet sheet.Spreadsheet, args []interface{}) error {
	fmt.Printf("Set with args %v on sheet %v\n", args, cellSheet)

	x, ok := toInt(args[0])
	if !ok {
		return MALFORMED_ERROR(args[0])
	}
	y, ok := toInt(args[1])
	if !ok {
		return MALFORMED_ERROR(args[1])
	}

	row := y
	col := colString(x)

	if equation, err := parseEquation(args[2]); err != nil {
		return err
	} else {
		setSheetCell(cellSheet, row, col, equation)
		return nil
	}
}

func executeAt(cellSheet sheet.Spreadsheet, args []interface{}) (float64, error) {
	fmt.Printf("At with args %v on sheet %v\n", args, cellSheet)

	x, ok := toInt(args[0])
	if !ok {
		return 0, MALFORMED_ERROR(args[0])
	}
	y, ok := toInt(args[1])
	if !ok {
		return 0, MALFORMED_ERROR(args[1])
	}

	row := y
	col := colString(x)

	return getSheetCell(cellSheet, row, col)
}

///////////////////////////////////////////
/******* Package-specific adapters *******/
///////////////////////////////////////////

// Using sheet package, create new spreadsheet
func createSheet(cells [][]sheet.Cell) sheet.Spreadsheet {
	// Implementation-specific
	return sheet.CreateSpreadSheet(cells)
}

// Using sheet package, set cell equation
func setSheetCell(s sheet.Spreadsheet, row int, col string, cellVal sheet.Equation) {
	newCell := sheet.CreateCell(cellVal)

	s.ChangeCell(newCell, row, col)
}

// Using sheet package, get cell value
func getSheetCell(s sheet.Spreadsheet, row int, col string) (float64, error) {
	return s.Evaluate(row, col)
}

/////////////////////////
/******* Parsing *******/
/////////////////////////

// NOTE: This adapter connects generic arrays to sheet package structs
func parseEquation(jArr interface{}) (sheet.Equation, error) {
	// fmt.Printf("Parsing %v of type %s\n", jArr, reflect.TypeOf(jArr))

	switch jArr.(type) {
	case float64:
		return sheet.CreateConstant(jArr.(float64)), nil

	case []interface{}:
		arr := jArr.([]interface{})
		if len(arr) != 3 {
			return nil, fmt.Errorf("Cannot parse JF Array of length %v", len(arr))
		}

		left, center, right := arr[0], arr[1], arr[2]
		if left == ">" {
			x, ok := toInt(center)
			if !ok {
				return nil, MALFORMED_ERROR(left)
			}
			y, ok := toInt(right)
			if !ok {
				return nil, MALFORMED_ERROR(right)
			}

			row := y
			col := colString(x)

			return sheet.CreateReference(row, col), nil
		} else if center == "*" {
			return parseFunction(left, right, sheet.Mult)
		} else if center == "+" {
			return parseFunction(left, right, sheet.Add)
		}
	}

	return nil, fmt.Errorf("Invalid JF value: %v", jArr)
}

func parseReference(xVal, yVal interface{}) (ref sheet.Reference, err error) {
	x, ok := toInt(xVal)
	if !ok {
		return sheet.Reference{}, MALFORMED_ERROR(xVal)
	}
	y, ok := toInt(yVal)
	if !ok {
		return sheet.Reference{}, MALFORMED_ERROR(yVal)
	}

	row := y
	col := colString(x)

	return sheet.CreateReference(row, col), nil
}

func parseFunction(left, right interface{}, op func(float64, float64) float64) (sheet.Function, error) {
	l, el := parseEquation(left)
	r, er := parseEquation(right)
	if el != nil {
		return sheet.Function{}, el
	} else if er != nil {
		return sheet.Function{}, er
	} else {
		return sheet.CreateFunction(op, l, r), nil
	}
}

func toInt(f interface{}) (int, bool) {
	if fl, ok := f.(float64); ok {
		return int(fl), true
	}
	return 0, false
}

func colString(xVal int) string {

	// 0 = A, 1 = B, ... 25 = Z,
	permCount := sheet.GetPerm(xVal)
	str := make([]rune, 0)

	i := 0.0

	// Start at i=0, see if X > 26^i, add that "char" to string, i++
	for {
		upper := math.Pow(26, i)
		rem := math.Remainder(float64(xVal), upper)

		char := rune(int(rem) + 65)
		str = append(str[:0], char)

		if float64(permCount) > math.Pow(26, i) {
			i++
		} else {
			break
		}
	}

	return string(str)
}

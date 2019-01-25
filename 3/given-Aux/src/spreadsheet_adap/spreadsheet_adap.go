package spreadsheet_adap

import (
	"errors"
	"math"
	"spreadsheet"
)

type Spreadsheet interface {
	Write(eq spreadsheet.Equation, row int, col int) error
	Value(row int, col int) (float64, error)
}

type Coord struct {
	x int
	y int
}

type mySpread struct {
	delegate *spreadsheet.Spreadsheet
}

func (spread mySpread) Write(eq spreadsheet.Equation, row int, col int) error {
	nRow, nCol, err := convertCoords(row, col)

	if err != nil {
		return err
	}

	spread.delegate.ChangeCell(spreadsheet.CreateCell(eq), nRow, nCol)
	return nil
}

func (spread mySpread) Value(row int, col int) (float64, error) {
	nRow, nCol, err := convertCoords(row, col)

	if err != nil {
		return 0, err
	}

	return spread.delegate.Evaluate(nRow, nCol)
}

func MakeSpreadsheet(eq [][]spreadsheet.Equation) Spreadsheet {
	total := make([][]spreadsheet.Cell, len(eq))
	for row := 0; row < len(eq); row++ {
		crRow := make([]spreadsheet.Cell, len(eq[row]))
		for col := 0; col < len(eq[row]); col++ {
			crRow[col] = spreadsheet.CreateCell(eq[row][col])
		}
		total[row] = crRow
	}

	del := spreadsheet.CreateSpreadSheet(total)
	return mySpread{&del}
}

func MakeBlankSpreadsheet() Spreadsheet {
	del := spreadsheet.CreateSpreadSheet(make([][]spreadsheet.Cell, 0))
	return mySpread{&del}
}

func convertCoords(row int, col int) (int, string, error) {
	nRow := row + 1
	nCol, err := convertCol(col)

	return nRow, nCol, err
}

func convertCol(col int) (string, error) {
	if col < 0 {
		return "", errors.New("No negative conversion.")
	}

	lengthOfPerm := 0
	for spreadsheet.GetPerm(lengthOfPerm+1) <= col {
		lengthOfPerm += 1
	}

	num := col - spreadsheet.GetPerm(lengthOfPerm)

	strTotal := ""
	for i := lengthOfPerm; i >= 0; i-- {
		place := int(math.Pow(26, float64(i)))
		for mult := 25; mult >= 0; mult-- {
			if mult*place <= num {
				strTotal = strTotal + getCharOf(mult)
				num = num - mult*place
				break
			}
		}
	}
	return strTotal, nil
}

func getCharOf(mult int) string {
	ascii := mult + 65
	return string(ascii)
}

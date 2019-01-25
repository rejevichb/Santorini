// Package spreadsheet provides the interface spreadsheet functionality
// defined in specification.
package spreadsheet

import (
	"errors"
	"fmt"
	"math"
)

type Spreadsheet struct {
	units [][]Cell
}

// Cell is a struct the holds an equation.
type Cell struct {
	eq Equation
}

// Reference is a struct that represents another cell in the spreadsheet
type Reference struct {
	row int
	col string
}

// Equation is an interface that declares the ability to be evaluated.
type Equation interface {
	evaluate(spread *Spreadsheet, seen map[coord]bool) (float64, error)
}

// coord is a simple pair structure for storing coordinates
type coord struct {
	x int
	y int
}

// Function is a struct representing the application of the operator upon two equations
type Function struct {
	operator func(float64, float64) float64
	eq1      Equation
	eq2      Equation
}

// Constant is an Equation representing a floating point number.
type Constant struct {
	value float64
}

func CreateSpreadSheet(units [][]Cell) Spreadsheet {
	return Spreadsheet{units}
}

// GetPerm generates the permutation of letters A-Z up to the given n.
// Nessecary for calculating columns.
func GetPerm(n int) int {
	total := 0

	for i := n; i > 0; i-- {
		total += int(math.Pow(26, float64(i)))
	}

	return total
}

// The changeCell function manipulates the Spreadsheet by altering the
// spreadsheet pointed to by overwriting the old cell at location with
// a given cell.
func (spread *Spreadsheet) ChangeCell(c Cell, row int, col string) error {
	trueCol, err := convertToInt(col)

	if err != nil {
		return errors.New("Improper column format")
	}

	trueRow := row - 1

	if !spread.withinBounds(trueRow, trueCol) {
		return errors.New("Cell not within bounds. ")
	}

	spread.units[trueRow][trueCol] = c

	return nil
}

// CreateConstant is a factory function for a Constant
func CreateConstant(value float64) Constant {
	return Constant{value}
}

// CreateReference is a factory function that creates a Reference represented by the row and col
func CreateReference(row int, col string) Reference {
	return Reference{row, col}
}

// The IncreaseRow function provides a non-error functionality
// to increase the row count of the given Spreadsheet.
func (spread *Spreadsheet) IncreaseRow() {
	var colSz int
	if len(spread.units) <= 0 {
		colSz = 1
	} else {
		colSz = len(spread.units[0])
	}

	nRow := make([]Cell, colSz)

	for i := 0; i < len(nRow); i++ {
		nRow[i] = Cell{Constant{0}}
	}

	spread.units = append(spread.units, nRow)
}

// The IncreaseRow function provides a non-error functionality
// to increase the column count of the given Spreadsheet.
func (spread *Spreadsheet) IncreaseColumn() {
	for i := 0; i < len(spread.units); i++ {
		spread.units[i] = append(spread.units[i], Cell{Constant{0}})
	}
}

// The Evaluate function computes and returns a returned float value of the
// cell at the location in the Spreadsheet. Error is returned if a value
// cannot be computed.
func (spread *Spreadsheet) Evaluate(row int, col string) (float64, error) {
	iCol, err := convertToInt(col)

	if err != nil {
		return -1, errors.New("Improper column format")
	}

	fixedRow := row - 1

	if !spread.withinBounds(fixedRow, iCol) {
		rows, cols := spread.size()
		fmt.Println(fixedRow, iCol, rows, cols)
		return -1, errors.New("Cell not within bounds.")
	}

	return spread.units[fixedRow][iCol].evaluate(spread, make(map[coord]bool))
}

// CreateCell is a factory function for Cells taking the boxed equation.
func CreateCell(eq Equation) Cell {
	return Cell{eq}
}

// Add adds a and b
func Add(a float64, b float64) float64 {
	return a + b
}

// Mult multiplies a and b
func Mult(a float64, b float64) float64 {
	return a * b
}

// CreateFunction is a factory method for creating Function objects
func CreateFunction(operator func(float64, float64) float64, eq1 Equation, eq2 Equation) Function {
	return Function{operator, eq1, eq2}
}

// Function size returns the number of rows and columns that can
// store data in the given Spreadsheet.
func (spread *Spreadsheet) size() (int, int) {
	if len(spread.units) <= 0 {
		return 0, 0
	} else {
		return len(spread.units), len(spread.units[0])
	}
}

// The convertToInt function provides the full functionality of
// converting an alphabetical numbering scheme into a numeric indexing scheme.
func convertToInt(col string) (int, error) {
	total := 0
	for i := 0; i < len(col); i++ {
		charInd := len(col) - 1 - i
		power := int(math.Pow(26, float64(i)))
		actChar := col[charInd]
		convChar := int(actChar) - 65

		if convChar < 0 || convChar > 25 {
			return -1, errors.New("Invalid string format.")
		}
		total += power * convChar
	}

	return total + GetPerm(len(col)-1), nil
}

// The withinBounds function returns if the given parameters
// represent a valid location in the Spreadsheet pointed to.
func (spread *Spreadsheet) withinBounds(trueRow int, trueCol int) bool {
	return trueRow >= 0 && trueRow < len(spread.units) && trueCol >= 0 && trueCol < len(spread.units[trueRow])
}

// evaluate evaluates a Cell in the context of the given spread. seen is an accumulator for detecting circular references.
func (c Cell) evaluate(spread *Spreadsheet, seen map[coord]bool) (float64, error) {
	return c.eq.evaluate(spread, seen)
}

// evaulate of a constant is simply the value of the constant.
func (c Constant) evaluate(spread *Spreadsheet, seen map[coord]bool) (float64, error) {
	return c.value, nil
}

// evaluates the Function by applying the operator to the evaluation of the two equations.
func (f Function) evaluate(spread *Spreadsheet, seen map[coord]bool) (float64, error) {
	a1, err1 := f.eq1.evaluate(spread, seen)
	a2, err2 := f.eq2.evaluate(spread, seen)

	errF := err1
	if err2 != nil {
		errF = err2
	}
	return f.operator(a1, a2), errF
}

// evaluate of a Reference evaluates to the evaluation of the cell represented by the Reference
func (ref Reference) evaluate(spread *Spreadsheet, seen map[coord]bool) (float64, error) {
	row := ref.row - 1

	col, err := convertToInt(ref.col)

	if err != nil {
		return 0, err
	}

	if seen[coord{row, col}] {
		return -1, errors.New("Circular Reference")
	}

	if !spread.withinBounds(row, col) {
		return -1, fmt.Errorf("Out of bounds reference row: %d col: %d", row, col)
	}

	seen[coord{row, col}] = true

	return spread.units[row][col].evaluate(spread, seen)
}

package sheet

import (
	"fmt"
)

/////////////
//* TYPES *//
/////////////

type Operation string

const (
	Constant       Operation = "constant"
	Addition       Operation = "addition"
	Multiplication Operation = "operation"
)

type Number struct {
	Val int
}

type Formula interface {
	Calculate(s Sheet) (float64, error)
	validFormula(A, B Formula, op Operation) bool
}

type Equation struct {
	A  Formula
	B  Formula
	Op Operation
}

type Sheet struct {
	Cells [][]Formula
}

type CellRef struct {
	X, Y int
}

/////////////////////////////
//* Methods and Functions *//
/////////////////////////////

const INVALID_STRING = "Invalid Formula: A: %s, B: %s, Op: %s"

func (s *Sheet) Set(X, Y int, f Formula) error {
	if X >= len(s.Cells) {
		return fmt.Errorf("invalid CellRef X-value: %v", X)
	}
	if Y >= len(s.Cells[0]) {
		return fmt.Errorf("invalid CellRef Y-value: %v", Y)
	}

	s.Cells[X][Y] = f

	return nil
}

func newFormula(A, B Formula, op Operation) (Formula, error) {
	if isValid := validFormula(A, B, op); !isValid {
		return nil, fmt.Errorf(INVALID_STRING, A, B, op)
	}
	switch op {
	case Constant:
		return &Equation{A, B, Constant}, nil
	case Multiplication:
		return &Equation{A, B, Multiplication}, nil
	case Addition:
		return &Equation{A, B, Addition}, nil
	}

	return nil, fmt.Errorf("Invalid operation: %s", op)
}

// Calculate the Formula at (X, Y)
func (s *Sheet) Calculate(X, Y int) (float64, error) {
	cell := s.Cells[X][Y]

	return cell.Calculate(*s)
}

/////////////////////////////////////////
//* Calculation functions on Formulae *//
/////////////////////////////////////////

// Perform the equation's operation on the calculated result of a and b
// Return a float64 val
func (e *Equation) Calculate(s Sheet) (float64, error) {

	if isValid := validFormula(e.A, e.B, e.Op); !isValid {
		return 0, fmt.Errorf(INVALID_STRING, e.A, e.B, e.Op)
	}

	switch e.Op {
	case Constant:
		return e.A.Calculate(s)

	case Multiplication:
		v1, err := e.A.Calculate(s)
		if err != nil {
			return 0, err
		}
		v2, err := e.B.Calculate(s)
		if err != nil {
			return 0, err
		}
		return v1 * v2, nil

	case Addition:
		v1, err := e.A.Calculate(s)
		if err != nil {
			return 0, err
		}
		v2, err := e.B.Calculate(s)
		if err != nil {
			return 0, err
		}
		return v1 + v2, nil
	}

	return 0, fmt.Errorf("Invalid Equation: %s", e)
}

// Return the Number
func (n *Number) Calculate(s Sheet) (float64, error) {
	return float64(n.Val), nil
}

func (c *CellRef) Calculate(s Sheet) (float64, error) {
	return s.Calculate(c.X, c.Y)
}

//////////////////////////////////////////////
//* formula validation and error handling  *//
//////////////////////////////////////////////

func validFormula(A, B Formula, op Operation) bool {
	switch op {
	case Constant:
		if A == nil {
			return false
		}
		if B != nil && A != nil {
			return false
		}
	case Multiplication, Addition:
		if A == nil || B == nil {
			return false
		}
	}
	return true //, nil
}

func (e *Equation) validFormula(A, B Formula, op Operation) bool {
	return validFormula(e.A, e.B, e.Op)
}

func (n *Number) validFormula(A, B Formula, op Operation) bool {
	return true
}

func (c *CellRef) validFormula(A, B Formula, op Operation) bool {
	return true
}

func (s *Sheet) validFormula(A, B Formula, op Operation) bool {
	return true
}

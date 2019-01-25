package sheet

import (
	"testing"
)

func createSheet(rows, cols int) Sheet {
	cells := make([][]Formula, 0)

	for x := 0; x < rows; x++ {
		cells = append(cells, make([]Formula, cols, cols))
	}

	return Sheet{cells}
}

func TestNumberAddition(t *testing.T) {
	sh := createSheet(2, 2)
	e := Equation{&Number{1}, &Number{2}, Addition}

	if v, err := e.Calculate(sh); v != 3 || err != nil {
		t.Fail()
	}
}

func TestNumberMult(t *testing.T) {
	sh := createSheet(2, 2)
	e := Equation{&Number{1}, &Number{2}, Multiplication}

	if v, err := e.Calculate(sh); v != 2 || err != nil {
		t.Fail()
	}
}

func TestNumberMultNestedFormula(t *testing.T) {
	sh := createSheet(4, 4)
	e := Equation{&Number{1}, &Number{2}, Multiplication}
	e1 := Equation{&e, &Number{2}, Multiplication}

	if v, err := e1.Calculate(sh); v != 4 || err != nil {
		t.Fail()
	}
}

func TestNumberMultNestedFormula2(t *testing.T) {
	sh := createSheet(3, 1)
	e := Equation{&Number{1}, &Number{2}, Multiplication}
	e1 := Equation{&Number{1}, &Number{1}, Addition}
	e2 := Equation{&e, &e1, Multiplication}

	if v, err := e2.Calculate(sh); v != 4 || err != nil {
		t.Fail()
	}
}

func TestMalformedEquation(t *testing.T) {
	sh := createSheet(5, 10)
	e := Equation{nil, nil, Constant}
	_, err := e.Calculate(sh)

	if err == nil {
		t.Fail()
	}
}

func TestCellCalculate(t *testing.T) {
	sh := createSheet(2, 2)
	sh.Set(0, 0, &Number{7})

	cell := CellRef{0, 0}
	if v, err := cell.Calculate(sh); v != 7 || err != nil {
		t.Fail()
	}
}

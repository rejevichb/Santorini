package spreadsheet

import (
	"fmt"
	"strconv"
	"testing"
	"testing_utils"
)

func TestConvertToInt(t *testing.T) {
	tables := []struct {
		in  string
		out int
	}{
		{"A", 0},
		{"C", 2},
		{"Z", 25},
		{"AA", 26},
		{"AB", 27},
		{"ZZ", 26*25 + 25 + 26},
		{"CBA", 2080},
	}

	for _, table := range tables {
		ans, _ := convertToInt(table.in)

		if ans != table.out {
			t.Errorf("For input: \"%s\" expected: %d got %d", table.in, table.out, ans)
		}
	}
}

func TestGetPerm(t *testing.T) {
	if GetPerm(2) != 26*26+26 {
		t.Error("Get perm got " + strconv.Itoa(GetPerm(2)))
	}
}

func TestSize(t *testing.T) {
	s := Spreadsheet{}

	row, col := s.size()
	testing_utils.AssertEqual(0, row, t)
	testing_utils.AssertEqual(0, col, t)

	s.IncreaseRow()

	row, col = s.size()
	testing_utils.AssertEqual(1, row, t)
	testing_utils.AssertEqual(1, col, t)

	s.IncreaseColumn()

	row, col = s.size()

	testing_utils.AssertEqual(1, row, t)
	testing_utils.AssertEqual(2, col, t)

	s.IncreaseRow()

	row, col = s.size()
	testing_utils.AssertEqual(2, row, t)
	testing_utils.AssertEqual(2, col, t)

	val, _ := s.Evaluate(1, "A")
	testing_utils.AssertEqualFloat(0, val, t)

	s.ChangeCell(Cell{Constant{2}}, 1, "B")
	s.ChangeCell(Cell{Constant{3}}, 2, "B")
	s.ChangeCell(Cell{Function{Add, Reference{1, "B"}, Reference{2, "B"}}}, 1, "A")

	val2, err := s.Evaluate(1, "A")

	row, col = s.size()
	testing_utils.AssertEqual(2, row, t)
	testing_utils.AssertEqual(2, col, t)

	testing_utils.AssertTrue(s.withinBounds(1, 1), t)

	if err != nil {
		t.Errorf("Evaluate raised error %s", err)
	}

	testing_utils.AssertEqualFloat(5, val2, t)

	s.ChangeCell(Cell{Constant{-2}}, 1, "B")

	val, err = s.Evaluate(1, "A")
	testing_utils.AssertEqualFloat(val, 1, t)

}

func TestIntegrationOfAdapter(t *testing.T) {
	sheet := make([][]Cell, 2)
	row0 := make([]Cell, 2)
	row1 := make([]Cell, 2)

	row0[0] = CreateCell(CreateConstant(1))
	row0[1] = CreateCell(CreateReference(0, "A"))

	row1[0] = CreateCell(CreateConstant(2))
	row1[1] = CreateCell(CreateConstant(3))

	sheet[0] = row0
	sheet[1] = row1

	rlSheet := CreateSpreadSheet(sheet)

	val, _ := rlSheet.Evaluate(1, "A")
	testing_utils.AssertEqualFloat(val, 1, t)

	ref1 := CreateReference(2, "A")
	ref2 := CreateReference(2, "B")
	mulRef := CreateCell(CreateFunction(Mult, ref1, ref2))
	rlSheet.ChangeCell(mulRef, 1, "A")
	rlSheet.ChangeCell(CreateCell(CreateReference(1, "A")), 1, "B")

	val2, err := rlSheet.Evaluate(1, "B")
	fmt.Println(err)
	testing_utils.AssertEqualFloat(6, val2, t)

	val, _ = rlSheet.Evaluate(1, "A")
	testing_utils.AssertEqualFloat(6, val, t)

	val, _ = rlSheet.Evaluate(2, "A")
	testing_utils.AssertEqualFloat(2, val, t)

	val, _ = rlSheet.Evaluate(2, "B")
	testing_utils.AssertEqualFloat(3, val, t)
}

package spreadsheet_adap

import "testing"

func TestConvertCol(t *testing.T) {
	tables := []struct {
		in  int
		out string
	}{
		{0, "A"},
		{2, "C"},
		{25, "Z"},
		{26, "AA"},
		{27, "AB"},
		{26*25 + 25 + 26, "ZZ"},
		{2080, "CBA"},
	}

	for _, table := range tables {
		ans, _ := convertCol(table.in)

		if ans != table.out {
			t.Errorf("For input: %d expected: \"%s\" got \"%s\"", table.in, table.out, ans)
		}
	}
}

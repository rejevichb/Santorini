package spreadsheet_adap

import "spreadsheet"

type Equation = spreadsheet.Equation

func MakeConstant(n float64) spreadsheet.Equation {
	return spreadsheet.CreateConstant(n)
}

func MakeMult(e1 spreadsheet.Equation, e2 spreadsheet.Equation) spreadsheet.Equation {
	return spreadsheet.CreateFunction(spreadsheet.Mult, e1, e2)
}

func MakeAdd(e1 spreadsheet.Equation, e2 spreadsheet.Equation) spreadsheet.Equation {
	return spreadsheet.CreateFunction(spreadsheet.Add, e1, e2)
}

func MakeRef(row int, col int) spreadsheet.Equation {
	nRow, nCol, _ := convertCoords(row, col)

	return spreadsheet.CreateReference(nRow, nCol)
}

func MakeBlankEq() spreadsheet.Equation {
	return MakeConstant(0)
}

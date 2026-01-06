package frame

import (
	"fmt"
)

type Frame struct {
	Cols       []*Column
	NameColMap map[string]int
}

func (f *Frame) Select(c ...ColExpr) (*Frame, error) {
	// return new frame, deep copy columns
	newCols := make([]*Column, len(c))
	newNameColMap := make(map[string]int)

	for i, col := range c {
		colIdx, ok := f.NameColMap[col.Name]
		// column not in Frame
		if !ok {
			return nil, fmt.Errorf("Column '%s' not recognized", col.Name)
		}

		vec := f.Cols[colIdx].Vec.DeepCopy()
		selectedCol := Column{
			Name:  col.Name,
			DType: f.Cols[colIdx].DType,
			Vec:   vec,
		}
		newCols[i] = &selectedCol
		newNameColMap[selectedCol.Name] = i
	}
	return &Frame{Cols: newCols, NameColMap: newNameColMap}, nil
}

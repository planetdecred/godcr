package decredmaterial

import (
	"gioui.org/layout"
)

type GridLayout struct {
	List      *layout.List
	Alignment layout.Alignment
	Direction layout.Direction
	RowCount  int
}

func (g GridLayout) Layout(gtx layout.Context, num int, el GridElement) layout.Dimensions {
	rows := make([]layout.Widget, 0)

	currentRow := make([]layout.FlexChild, 0)
	for i := 0; i < num; i++ {
		index := i
		currentRow = append(currentRow, layout.Rigid(func(gtx C) D { return el(gtx, index) }))

		if len(currentRow) >= g.RowCount {
			rowCopy := currentRow
			rows = append(rows, func(gtx C) D {
				return layout.Flex{Alignment: g.Alignment}.Layout(gtx, rowCopy...)
			})
			currentRow = make([]layout.FlexChild, 0)
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, func(gtx C) D {
			return layout.Flex{}.Layout(gtx, currentRow...)
		})
	}

	return g.List.Layout(gtx, len(rows), func(gtx C, index int) D {
		return g.Direction.Layout(gtx, func(gtx C) D {
			return rows[index](gtx)
		})
	})
}

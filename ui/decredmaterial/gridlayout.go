package decredmaterial

import (
	"gioui.org/layout"
)

type GridLayout struct {
	List              *layout.List
	HorizontalSpacing layout.Spacing
	Alignment         layout.Alignment
	Direction         layout.Direction
	RowCount          int
}

func (g GridLayout) Layout(gtx layout.Context, num int, el GridElement) layout.Dimensions {
	rows := make([]layout.Widget, 0)

	currentRow := make([]layout.FlexChild, 0)

	appendRow := func(row []layout.FlexChild) {
		rows = append(rows, func(gtx C) D {
			return layout.Flex{Alignment: g.Alignment, Spacing: g.HorizontalSpacing}.Layout(gtx, row...)
		})
		currentRow = make([]layout.FlexChild, 0)
	}

	for i := 0; i < num; i++ {
		index := i
		currentRow = append(currentRow, layout.Rigid(func(gtx C) D { return el(gtx, index) }))

		if len(currentRow) >= g.RowCount {
			appendRow(currentRow)
		}
	}

	if len(currentRow) > 0 {
		appendRow(currentRow)
	}

	return g.List.Layout(gtx, len(rows), func(gtx C, index int) D {
		return g.Direction.Layout(gtx, func(gtx C) D {
			return rows[index](gtx)
		})
	})
}

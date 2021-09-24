package renderers

import (
	//"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
)

type cellAlign int

const (
	cellAlignLeft cellAlign = iota
	cellAlignRight
	cellAlignCenter
	cellAlignCopyHeader
)

type cell struct {
	content       string
	alignment     cellAlign
	contentLength float64
}

type row struct {
	isHeader bool
	cells    []cell
}

type table struct {
	theme *decredmaterial.Theme
	rows  []row
}

func newTable(theme *decredmaterial.Theme) *table {
	return &table{
		theme: theme,
	}
}

func (t *table) startNextRow() {
	t.rows = append(t.rows, row{})
}

func (t *table) addCell(content string, alignment cellAlign, isHeader bool) {
	if len(t.rows) == 0 {
		return
	}

	cell := cell{
		content:       content,
		contentLength: float64(len(content)),
		alignment:     alignment,
	}

	rowIndex := len(t.rows) - 1
	t.rows[rowIndex].isHeader = isHeader
	t.rows[rowIndex].cells = append(t.rows[rowIndex].cells, cell)
}

// normalize ensure that the table has the same number of cells
// in each rows, header or not.
func (t *table) normalize() {

}

func (t *table) setAlignment() {
	if len(t.rows) == 0 {
		return
	}

	for i := range t.rows {
		if i == 0 {
			continue
		}

		for cellIndex := range t.rows[i].cells {
			t.rows[i].cells[cellIndex].alignment = t.rows[0].cells[cellIndex].alignment
		}
	}
}

func (t *table) layoutCellLabel(gtx C, c cell, isHeader bool) D {
	var w layout.Direction
	switch c.alignment {
	case cellAlignLeft:
		w = layout.W
	case cellAlignRight:
		w = layout.E
	default:
		w = layout.Center
	}

	lbl := t.theme.Body2(c.content)
	if isHeader {
		lbl.Font.Weight = text.Bold
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.UniformInset(unit.Dp(7)).Layout(gtx, func(gtx C) D {
		return w.Layout(gtx, lbl.Layout)
	})
}

func (t *table) layoutRow(gtx C, r row) D {
	var maxHeight int
	width := float64(gtx.Constraints.Max.X) / float64(len(r.cells))
	line := t.theme.SeparatorVertical(0, 1)
	line.Color = t.theme.Color.Gray1

	return (&layout.List{Axis: layout.Horizontal}).Layout(gtx, len(r.cells), func(gtx C, i int) D {
		gtx.Constraints.Min.X = int(width)
		gtx.Constraints.Max.X = int(width)
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				dims := t.layoutCellLabel(gtx, r.cells[i], r.isHeader)
				if maxHeight < dims.Size.Y {
					maxHeight = dims.Size.Y
				}
				return dims
			}),
			layout.Expanded(func(gtx C) D {
				line.Height = maxHeight
				return line.Layout(gtx)
			}),
		)
	})
}

func (t *table) render() layout.Widget {
	return func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return (&layout.List{Axis: layout.Vertical}).Layout(gtx, len(t.rows), func(gtx C, i int) D {
			var bgCol color.NRGBA
			if i == 0 || i%2 != 0 {
				bgCol = t.theme.Color.Surface
			} else {
				bgCol = t.theme.Color.Background
			}

			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					border := widget.Border{Color: t.theme.Color.Gray1, CornerRadius: unit.Dp(0), Width: unit.Dp(1)}
					return border.Layout(gtx, func(gtx C) D {
						return decredmaterial.Fill(gtx, bgCol)
					})
				}),
				layout.Stacked(func(gtx C) D {
					return t.layoutRow(gtx, t.rows[i])
				}),
			)
		})
	}
}

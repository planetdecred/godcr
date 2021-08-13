package renderers

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

func getLabel(lbl decredmaterial.Label) decredmaterial.Label {
	return lbl
}

func setStyle(lbl *decredmaterial.Label, style string) {
	var s text.Style

	switch style {
	case "italic":
		s = text.Italic
	case "regular":
		s = text.Regular
	}

	lbl.Font.Style = s
}

func setWeight(lbl *decredmaterial.Label, weight string) {
	var w text.Weight

	switch weight {
	case "normal":
		w = text.Normal
	case "medium":
		w = text.Medium
	case "bold", "strong":
		w = text.Bold
	default:
		w = lbl.Font.Weight
	}

	lbl.Font.Weight = w
}

func getHeading(txt string, level int, theme *decredmaterial.Theme) decredmaterial.Label {
	var lblWdgt func(string) decredmaterial.Label

	switch level {
	case 5:
		lblWdgt = theme.H5
	case 6:
		lblWdgt = theme.H6
	default:
		lblWdgt = theme.H4
	}

	lbl := lblWdgt(txt)
	lbl.Font.Weight = text.Bold

	return lbl
}

func renderStrike(lbl decredmaterial.Label, theme *decredmaterial.Theme) layout.Widget {
	return func(gtx C) D {
		var dims D
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				dims = lbl.Layout(gtx)
				return dims
			}),
			layout.Expanded(func(gtx C) D {
				return layout.Inset{
					Top: unit.Dp((float32(dims.Size.Y) / float32(2))),
				}.Layout(gtx, func(gtx C) D {
					l := theme.Separator()
					l.Color = lbl.Color
					l.Width = dims.Size.X
					return l.Layout(gtx)
				})
			}),
		)
	}
}

func renderBlockQuote(lbl decredmaterial.Label, theme *decredmaterial.Theme) layout.Widget {
	words := strings.Fields(lbl.Text)

	return func(gtx C) D {
		var dims D

		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				l := theme.SeparatorVertical(dims.Size.Y, 10)
				l.Color = theme.Color.Gray
				return l.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				dims = layout.Inset{
					Left: unit.Dp(4),
				}.Layout(gtx, func(gtx C) D {
					return decredmaterial.GridWrap{
						Axis:      layout.Horizontal,
						Alignment: layout.Start,
					}.Layout(gtx, len(words), func(gtx C, i int) D {
						lbl.Text = words[i] + " "
						return lbl.Layout(gtx)
					})
				})

				return dims
			}),
		)
	}
}

func renderHorizontalLine(theme *decredmaterial.Theme) layout.Widget {
	l := theme.Separator()
	l.Width = 1
	return l.Layout
}

func renderEmptyLine(theme *decredmaterial.Theme, isList bool) layout.Widget {
	var padding = -5

	if isList {
		padding = -10
	}

	return func(gtx C) D {
		dims := theme.Body2("").Layout(gtx)
		dims.Size.Y = dims.Size.Y + padding
		return dims
	}
}

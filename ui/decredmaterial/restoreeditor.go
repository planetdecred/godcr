// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/values"
)

type RestoreEditor struct {
	t          *Theme
	Edit       Editor
	TitleLabel Label
	LineColor  color.NRGBA
	m10        unit.Value
	height     int
}

func (t *Theme) RestoreEditor(editor *widget.Editor, hint string, title string) RestoreEditor {

	e := t.Editor(editor, hint)
	e.Bordered = false
	return RestoreEditor{
		t:          t,
		Edit:       e,
		TitleLabel: t.Body2(title),
		LineColor:  t.Color.Gray1,
		m10:        values.MarginPadding10,
		height:     30,
	}
}

func (re RestoreEditor) Layout(gtx layout.Context) layout.Dimensions {
	border := widget.Border{Color: re.LineColor, CornerRadius: values.MarginPadding8, Width: values.MarginPadding2}
	return border.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Left:  re.m10,
					Right: re.m10,
				}.Layout(gtx, func(gtx C) D {
					return re.TitleLabel.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				l := re.t.VLine(re.height, 2)
				l.Color = re.t.Color.Gray1
				return layout.Inset{}.Layout(gtx, func(gtx C) D {
					return l.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				edit := re.Edit.Layout(gtx)
				re.height = edit.Size.Y
				return edit
			}),
		)
	})
}

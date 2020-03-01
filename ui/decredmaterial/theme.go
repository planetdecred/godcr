// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	lightBlue = color.RGBA{41, 112, 255, 255}
	orange    = color.RGBA{237, 109, 71, 255}
	green     = color.RGBA{46, 214, 161, 255}

	keyBlue   = color.RGBA{0x29, 0x70, 0xFF, 255}
	turquiose = color.RGBA{0x2E, 0xD6, 0xA1, 255}
	darkBlue  = color.RGBA{0x09, 0x14, 0x40, 255}
)

type Theme struct {
	Shaper text.Shaper
	Color  struct {
		Primary    color.RGBA
		Secondary  color.RGBA
		Text       color.RGBA
		Hint       color.RGBA
		InvText    color.RGBA
		Success    color.RGBA
		Danger     color.RGBA
		Background color.RGBA
		Surface    color.RGBA
	}
	Icon struct {
		ContentCreate *Icon
		ContentAdd    *Icon
	}
	TextSize              unit.Value
	checkBoxCheckedIcon   *Icon
	checkBoxUncheckedIcon *Icon
	radioCheckedIcon      *Icon
	radioUncheckedIcon    *Icon
}

func NewTheme() *Theme {
	t := &Theme{
		Shaper: font.Default(),
	}
	t.Color.Primary = rgb(0x3f51b5)
	t.Color.Text = rgb(0x000000)
	t.Color.Hint = rgb(0xbbbbbb)
	t.Color.InvText = rgb(0xffffff)
	t.TextSize = unit.Sp(16)

	t.checkBoxCheckedIcon = mustIcon(NewIcon(icons.ToggleCheckBox))
	t.checkBoxUncheckedIcon = mustIcon(NewIcon(icons.ToggleCheckBoxOutlineBlank))
	t.radioCheckedIcon = mustIcon(NewIcon(icons.ToggleRadioButtonChecked))
	t.radioUncheckedIcon = mustIcon(NewIcon(icons.ToggleRadioButtonUnchecked))

	return t
}

func (t *Theme) Background(gtx *layout.Context, w layout.Widget) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			fillWithColor(gtx, ARGB(0x22444444))
		}),
		layout.Stacked(w),
	)
}

func (t *Theme) Surface(gtx *layout.Context, w layout.Widget) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			fillWithColor(gtx, RGB(0xffffff))
		}),
		layout.Stacked(w),
	)
}

func mustIcon(ic *Icon, err error) *Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func rgb(c uint32) color.RGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}

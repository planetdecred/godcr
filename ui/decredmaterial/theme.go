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
	// decred primary colors

	keyblue   = rgb(0x2970ff)
	turquiose = rgb(0x2ed6a1)
	darkblue  = rgb(0x091440)

	// decred complemetary colors

	lightblue = rgb(0x70cbff)
	orange    = rgb(0xed6d47)
	green     = rgb(0x41bf53)
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
	t.Color.Primary = keyblue
	t.Color.Text = darkblue
	t.Color.Hint = rgb(0xbbbbbb)
	t.Color.InvText = rgb(0xffffff)
	t.Color.Background = argb(0x22444444)
	t.Color.Surface = rgb(0xffffff)
	t.Color.Success = green
	t.Color.Danger = rgb(0xff0000)
	t.TextSize = unit.Sp(16)

	t.checkBoxCheckedIcon = mustIcon(NewIcon(icons.ToggleCheckBox))
	t.checkBoxUncheckedIcon = mustIcon(NewIcon(icons.ToggleCheckBoxOutlineBlank))
	t.radioCheckedIcon = mustIcon(NewIcon(icons.ToggleRadioButtonChecked))
	t.radioUncheckedIcon = mustIcon(NewIcon(icons.ToggleRadioButtonUnchecked))

	return t
}

func (t *Theme) Background(gtx *layout.Context, w layout.Widget) {
	layout.Stack{
		Alignment: layout.Center,
	}.Layout(gtx,
		layout.Expanded(func() {
			fill(gtx, t.Color.Background)
		}),
		layout.Stacked(w),
	)
}

func (t *Theme) Surface(gtx *layout.Context, w layout.Widget) {
	layout.Stack{
		Alignment: layout.Center,
	}.Layout(gtx,
		layout.Expanded(func() {
			fillMax(gtx, t.Color.Surface)
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

func fillMax(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Max, Y: cs.Height.Max}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
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

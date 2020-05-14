// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	// decred primary colors

	keyblue = rgb(0x2970ff)
	//turquiose = rgb(0x2ed6a1)
	darkblue = rgb(0x091440)

	// decred complemetary colors

	//lightblue = rgb(0x70cbff)
	//orange    = rgb(0xed6d47)
	green = rgb(0x41bf53)
)

type Theme struct {
	Shaper text.Shaper
	Color  struct {
		Primary    color.RGBA
		Secondary  color.RGBA
		Text       color.RGBA
		Hint       color.RGBA
		Overlay    color.RGBA
		InvText    color.RGBA
		Success    color.RGBA
		Danger     color.RGBA
		Background color.RGBA
		Surface    color.RGBA
		Gray       color.RGBA
		Black      color.RGBA
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
	chevronUpIcon         *Icon
	chevronDownIcon       *Icon
}

func NewTheme() *Theme {
	t := &Theme{
		Shaper: font.Default(),
	}
	t.Color.Primary = keyblue
	t.Color.Text = darkblue
	t.Color.Hint = rgb(0xbbbbbb)
	t.Color.InvText = rgb(0xffffff)
	t.Color.Overlay = rgb(0x000000)
	t.Color.Background = argb(0x22444444)
	t.Color.Surface = rgb(0xffffff)
	t.Color.Success = green
	t.Color.Danger = rgb(0xff0000)
	t.Color.Gray = rgb(0x808080)
	t.Color.Black = rgb(0x000000)
	t.TextSize = unit.Sp(16)

	t.checkBoxCheckedIcon = mustIcon(NewIcon(icons.ToggleCheckBox))
	t.checkBoxUncheckedIcon = mustIcon(NewIcon(icons.ToggleCheckBoxOutlineBlank))
	t.radioCheckedIcon = mustIcon(NewIcon(icons.ToggleRadioButtonChecked))
	t.radioUncheckedIcon = mustIcon(NewIcon(icons.ToggleRadioButtonUnchecked))
	t.chevronUpIcon = mustIcon(NewIcon(icons.NavigationExpandLess))
	t.chevronDownIcon = mustIcon(NewIcon(icons.NavigationExpandMore))

	return t
}

func (t *Theme) Modal(gtx *layout.Context, w layout.Widget) {
	overlayColor := t.Color.Black
	overlayColor.A = 200

	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			fillMax(gtx, overlayColor)
		}),
		layout.Stacked(func() {
			inset := layout.Inset{
				Top: unit.Dp(50),
			}
			inset.Layout(gtx, func() {
				fillMax(gtx, t.Color.Surface)
				inset := layout.Inset{
					Top:   unit.Dp(7),
					Left:  unit.Dp(25),
					Right: unit.Dp(25),
				}
				inset.Layout(gtx, w)
			})
		}),
	)
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
			fill(gtx, t.Color.Surface)
		}),
		layout.Stacked(w),
	)
}

func (t *Theme) alert(gtx *layout.Context, txt string, bgColor color.RGBA) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(2)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Width.Min),
					Y: float32(gtx.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, bgColor)
		}),
		layout.Stacked(func() {
			layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
				lbl := t.Body2(txt)
				lbl.Color = t.Color.Surface
				lbl.Layout(gtx)
			})
		}),
	)
}

func (t *Theme) ErrorAlert(gtx *layout.Context, txt string) {
	t.alert(gtx, txt, t.Color.Danger)
}

func (t *Theme) SuccessAlert(gtx *layout.Context, txt string) {
	t.alert(gtx, txt, t.Color.Success)
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

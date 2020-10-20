// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"gioui.org/widget/material"

	"golang.org/x/image/draw"

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

type (
	C = layout.Context
	D = layout.Dimensions

	ReadClipboard struct{}
)

type Theme struct {
	Shaper text.Shaper
	Base   *material.Theme
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
		ContentCreate *widget.Icon
		ContentAdd    *widget.Icon
	}
	TextSize              unit.Value
	checkBoxCheckedIcon   *widget.Icon
	checkBoxUncheckedIcon *widget.Icon
	radioCheckedIcon      *widget.Icon
	radioUncheckedIcon    *widget.Icon
	chevronUpIcon         *widget.Icon
	chevronDownIcon       *widget.Icon

	Clipboard     chan string
	ReadClipboard chan interface{}
}

func NewTheme(fontCollection []text.FontFace) *Theme {
	t := &Theme{
		Shaper: text.NewCache(fontCollection),
		Base:   material.NewTheme(fontCollection),
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

	t.checkBoxCheckedIcon = mustIcon(widget.NewIcon(icons.ToggleCheckBox))
	t.checkBoxUncheckedIcon = mustIcon(widget.NewIcon(icons.ToggleCheckBoxOutlineBlank))
	t.radioCheckedIcon = mustIcon(widget.NewIcon(icons.ToggleRadioButtonChecked))
	t.radioUncheckedIcon = mustIcon(widget.NewIcon(icons.ToggleRadioButtonUnchecked))
	t.chevronUpIcon = mustIcon(widget.NewIcon(icons.NavigationExpandLess))
	t.chevronDownIcon = mustIcon(widget.NewIcon(icons.NavigationExpandMore))
	t.Clipboard = make(chan string)
	return t
}

func (t *Theme) Background(gtx layout.Context, w layout.Widget) {
	layout.Stack{
		Alignment: layout.N,
	}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return fill(gtx, t.Color.Background)
		}),
		layout.Stacked(w),
	)
}

func (t *Theme) Surface(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return layout.Stack{
		Alignment: layout.Center,
	}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return fill(gtx, t.Color.Surface)
		}),
		layout.Stacked(w),
	)
}

func (t *Theme) ImageIcon(gtx layout.Context, icon image.Image, size int) layout.Dimensions {
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: size, Y: size}})
	draw.ApproxBiLinear.Scale(img, img.Bounds(), icon, icon.Bounds(), draw.Src, nil)
	iconOp := paint.NewImageOp(img)

	i := Image{Src: iconOp}
	i.Scale = float32(size) / float32(gtx.Px(unit.Dp(float32(size))))
	return i.Layout(gtx)
}

func (t *Theme) alert(gtx layout.Context, txt string, bgColor color.RGBA) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			clip.RRect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Min.X),
					Y: float32(gtx.Constraints.Min.Y),
				}},
			}.Add(gtx.Ops)
			return fill(gtx, bgColor)
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X / 3
			gtx.Constraints.Max.X = gtx.Constraints.Max.X / 3
			gtx.Constraints.Min.Y = 80
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx C) D {

					lbl := t.Body2(txt)
					lbl.Alignment = text.Start
					lbl.Color = t.Color.Surface
					return lbl.Layout(gtx)
				})
			})
		}),
	)
}

func (t *Theme) ErrorAlert(gtx layout.Context, txt string) layout.Dimensions {
	return t.alert(gtx, txt, t.Color.Danger)
}

func (t *Theme) SuccessAlert(gtx layout.Context, txt string) layout.Dimensions {
	return t.alert(gtx, txt, t.Color.Success)
}

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
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

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func fillMax(gtx layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Max.X, Y: cs.Max.Y}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
}

func fill(gtx layout.Context, col color.RGBA) layout.Dimensions {
	cs := gtx.Constraints
	d := image.Point{X: cs.Min.X, Y: cs.Min.Y}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	return layout.Dimensions{Size: d}
}

// mulAlpha scales all color components by alpha/255.
func mulAlpha(c color.RGBA, alpha uint8) color.RGBA {
	a := uint16(alpha)
	return color.RGBA{
		A: uint8(uint16(c.A) * a / 255),
		R: uint8(uint16(c.R) * a / 255),
		G: uint8(uint16(c.G) * a / 255),
		B: uint8(uint16(c.B) * a / 255),
	}
}

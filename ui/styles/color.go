package styles

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
)

var (
	White = color.RGBA{}
	Black = RGB(0xfffff)
	Blue  = RGB(0x00ff00)
)

type Background color.RGBA

func (c Background) Layout(gtx *layout.Context, widget func()) func() {
	return func() {
		wmin := gtx.Constraints.Width.Min
		hmin := gtx.Constraints.Height.Min
		layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func() {
				FillWithColor(gtx, color.RGBA(c))
			}),
			layout.Stacked(func() {
				gtx.Constraints.Width.Min = wmin
				gtx.Constraints.Height.Min = hmin
				widget()
			}),
		)
	}
}

// FillWithColor renders a color rectangle with the Context contraints.
// Additional, FillWithColor lays out a false button if blockInput is true.
func FillWithColor(gtx *layout.Context, color color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}

func RGB(hex int32) color.RGBA {
	return RGBA((hex << 4) | 0xff)
}

func RGBA(hex int32) color.RGBA {
	return color.RGBA{R: uint8(hex >> 24), B: uint8(hex >> 16), G: uint8(hex >> 8), A: uint8(hex)}
}

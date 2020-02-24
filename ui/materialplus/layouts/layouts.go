package layouts

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
)

// FillWithColor paints a color on the Context's max contraints.
// Restores the dimensions after painting.
func FillWithColor(gtx *layout.Context, color color.RGBA) {
	cs := gtx.Constraints
	dmin := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	rect := f32.Rectangle{
		Max: f32.Point{X: float32(cs.Width.Max), Y: float32(cs.Height.Max)},
	}
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{Rect: rect}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: dmin}
	op.InvalidateOp{}.Add(gtx.Ops)
}

// Faded halfs the color alpha
func Faded(c color.RGBA) color.RGBA {
	c.A /= 2
	return c
}

// RGB converts a 24 bit int color to color.RGBA with bits 0x__RRGGBB
func RGB(c uint32) color.RGBA {
	return ARGB(0xff000000 | c)
}

// ARGB converts a 32 bit int color to color.RGBA with bits 0xAARRGGBB
func ARGB(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

package materialplus

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
)

// fillWithColor paints a colored rectangle with gtx min contraints
func fillWithColor(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}

// Faded halfs the color alpha
func Faded(c color.RGBA) color.RGBA {
	c.A /= 2
	return c
}

// RGB converts a 24 bit color to color.RGBA with bits 0x__RRGGBB
func RGB(c uint32) color.RGBA {
	return ARGB(0xff000000 | c)
}

// ARGB converts a 32 bit color to color.RGBA with bits 0xAARRGGBB
func ARGB(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func toMax(gtx *layout.Context) {
	gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
	gtx.Constraints.Height.Min = gtx.Constraints.Height.Max
}

package styles

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
)

var (
	White = color.RGBA{}
	Black = rgba(0xfffffff)
	Blue  = rgba(0x00ff00ff)
)

var (
	BackgroundBlack = Background{
		Color: Black,
	}
	BackgroundBlue = Background{
		Color: Blue,
	}
)

type Background struct {
	Color      color.RGBA
	BlockInput bool
}

func (c Background) Layout(gtx *layout.Context, widget func()) func() {
	return func() {
		FillWithColor(gtx, c.Color, c.BlockInput)
		widget()
	}
}

// FillWithColor renders a color rectangle with the Context contraints.
// Additional, FillWithColor lays out a false button if blockInput is true.
func FillWithColor(gtx *layout.Context, color color.RGBA, blockInput bool) {
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{
		Max: f32.Point{
			X: float32(gtx.Constraints.Width.Max),
			Y: float32(gtx.Constraints.Height.Max),
		},
	}}.Add(gtx.Ops)
	if blockInput {
		new(widget.Button).Layout(gtx)
	}
}

func rgba(hex int32) color.RGBA {
	return color.RGBA{R: uint8(hex >> 24), B: uint8(hex >> 16), G: uint8(hex >> 8), A: uint8(hex)}
}

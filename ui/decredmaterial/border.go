package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// Border lays out a widget and draws a border inside it.
type Border struct {
	Color  color.NRGBA
	Radius CornerRadius
	Width  unit.Dp
}

func (b Border) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := w(gtx)
	sz := dims.Size

	width := gtx.Dp(b.Width)
	sz.X -= width
	sz.Y -= width

	r := image.Rectangle{Max: sz}
	r = r.Add(image.Point{X: width / 2, Y: width / 2})

	tr := gtx.Dp(unit.Dp(b.Radius.TopRight))
	tl := gtx.Dp(unit.Dp(b.Radius.TopLeft))
	br := gtx.Dp(unit.Dp(b.Radius.BottomRight))
	bl := gtx.Dp(unit.Dp(b.Radius.BottomLeft))
	radius := clip.RRect{
		Rect: r,
		NW:   tl, NE: tr, SE: br, SW: bl,
	}

	paint.FillShape(gtx.Ops,
		b.Color,
		clip.Stroke{
			Path:  radius.Path(gtx.Ops),
			Width: float32(width),
		}.Op(),
	)

	return dims
}

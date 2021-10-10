package decredmaterial

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// Border lays out a widget and draws a border inside it.
type Border struct {
	Color  color.NRGBA
	Radius CornerRadius
	Width  unit.Value
}

func (b Border) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	dims := w(gtx)
	sz := layout.FPt(dims.Size)

	width := float32(gtx.Px(b.Width))
	sz.X -= width
	sz.Y -= width

	r := f32.Rectangle{Max: sz}
	r = r.Add(f32.Point{X: width * 0.5, Y: width * 0.5})

	tr := float32(gtx.Px(unit.Dp(b.Radius.TopRight)))
	tl := float32(gtx.Px(unit.Dp(b.Radius.TopLeft)))
	br := float32(gtx.Px(unit.Dp(b.Radius.BottomRight)))
	bl := float32(gtx.Px(unit.Dp(b.Radius.BottomLeft)))
	radius := clip.RRect{
		Rect: r,
		NW:   tl, NE: tr, SE: br, SW: bl,
	}

	paint.FillShape(gtx.Ops,
		b.Color,
		clip.Stroke{
			Path:  radius.Path(gtx.Ops),
			Width: width,
		}.Op(),
	)

	return dims
}

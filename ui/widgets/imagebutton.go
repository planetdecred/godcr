package widgets

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	gioText "gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ImageButton struct {
	Text         string
	Font         gioText.Font
	Color        color.RGBA
	Background   color.RGBA
	Size         unit.Value
	Padding      unit.Value
	CornerRadius unit.Value
	shaper       gioText.Shaper

	alignment layout.Alignment
	Axis      layout.Axis
	Src       *material.Image
}

func NewImageButton(img *material.Image, text string) ImageButton {
	return ImageButton{
		Text: text,
		Font: gioText.Font{
			Size: unit.Sp(16).Scale(14.0 / 16.0),
		},
		Color:      rgb(0xffffff),
		Background: rgb(0x3f51b5),
		shaper:     font.Default(),
		alignment:  layout.Middle,
		Axis:       layout.Horizontal,

		Src:     img,
		Size:    unit.Dp(56),
		Padding: unit.Dp(16),
	}
}

func (b ImageButton) Layout(gtx *layout.Context, button *widget.Button, buttonTextSpace float32) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(8)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Width.Min),
					Y: float32(gtx.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)

			fill(gtx, b.Background)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),

		layout.Stacked(func() {
			iconAndLabel := layout.Flex{Axis: b.Axis, Alignment: layout.Middle}

			icon := layout.Rigid(func() {
				layout.Inset{Top: b.Padding, Left: b.Padding, Bottom: b.Padding, Right: unit.Dp(buttonTextSpace / 2)}.Layout(gtx, func() {
					b.Src.Scale = 0.2
					b.Src.Layout(gtx)
				})
			})

			label := layout.Rigid(func() {
				layout.Inset{Top: b.Padding, Right: b.Padding, Bottom: b.Padding, Left: unit.Dp(buttonTextSpace / 2)}.Layout(gtx, func() {
					paint.ColorOp{Color: b.Color}.Add(gtx.Ops)
					widget.Label{Alignment: gioText.Start}.Layout(gtx, b.shaper, b.Font, b.Text)
				})

			})

			if b.Src != nil {
				iconAndLabel.Layout(gtx, icon, label)
			} else {
				iconAndLabel.Layout(gtx, label)
			}

			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)
}

func drawInk(gtx *layout.Context, c widget.Click) {
	d := gtx.Now().Sub(c.Time)
	t := float32(d.Seconds())
	const duration = 0.5
	if t > duration {
		return
	}

	t = t / duration
	var stack op.StackOp
	stack.Push(gtx.Ops)
	size := float32(gtx.Px(unit.Dp(700))) * t
	rr := size * .5
	col := byte(0xaa * (1 - t*t))
	ink := paint.ColorOp{Color: color.RGBA{A: col, R: col, G: col, B: col}}
	ink.Add(gtx.Ops)

	op.TransformOp{}.Offset(c.Position).Offset(f32.Point{
		X: -rr,
		Y: -rr,
	}).Add(gtx.Ops)

	clip.Rect{
		Rect: f32.Rectangle{Max: f32.Point{
			X: size,
			Y: size,
		}},
		NE: rr, NW: rr, SE: rr, SW: rr,
	}.Op(gtx.Ops).Add(gtx.Ops)

	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: size, Y: size}}}.Add(gtx.Ops)
	stack.Pop()
	op.InvalidateOp{}.Add(gtx.Ops)
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

func rgb(c uint32) color.RGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

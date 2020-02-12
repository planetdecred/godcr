package widgets

import (
	"image"
	"image/color"

	"github.com/raedahgroup/godcr-gio/ui/helper"

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
	Text       string
	Font       gioText.Font
	Color      color.RGBA
	Background color.RGBA
	HPadding   unit.Value
	VPadding   unit.Value
	shaper     gioText.Shaper

	alignment layout.Alignment
	Axis      layout.Axis
	Src       *material.Image
}

func NewImageButton(img *material.Image, text string) ImageButton {
	img.Scale = 0.2

	return ImageButton{
		Text: text,
		Font: gioText.Font{
			Size: unit.Sp(16).Scale(14.0 / 16.0),
		},
		Color:      helper.RGB(0xffffff),
		Background: helper.RGB(0x3f51b5),
		shaper:     font.Default(),
		alignment:  layout.Middle,
		Axis:       layout.Horizontal,

		Src:      img,
		VPadding: unit.Dp(16),
		HPadding: unit.Dp(16),
	}
}

func (b ImageButton) Layout(gtx *layout.Context, button *widget.Button, buttonTextSpace float32) {
	hmin := gtx.Constraints.Width.Min
	vmin := gtx.Constraints.Height.Min

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

			helper.Fill(gtx, b.Background)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),

		layout.Stacked(func() {
			iconAndLabel := layout.Flex{Axis: b.Axis, Alignment: layout.Middle}
			gtx.Constraints.Width.Min = hmin
			gtx.Constraints.Height.Min = vmin

			icon := layout.Rigid(func() {
				layout.Inset{Top: b.VPadding, Left: b.HPadding, Bottom: b.VPadding, Right: unit.Dp(buttonTextSpace / 2)}.Layout(gtx, func() {
					b.Src.Layout(gtx)
				})
			})

			label := layout.Rigid(func() {
				layout.Inset{Top: b.VPadding, Right: b.HPadding, Bottom: b.VPadding, Left: unit.Dp(buttonTextSpace / 2)}.Layout(gtx, func() {
					paint.ColorOp{Color: b.Color}.Add(gtx.Ops)
					widget.Label{}.Layout(gtx, b.shaper, b.Font, b.Text)
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

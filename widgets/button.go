package widgets

import (
	"image"
	"image/color"
	"image/draw"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/raedahgroup/godcr-gio/helper"
	"golang.org/x/exp/shiny/iconvg"
)

type (
	Button struct {
		button    *widget.Button
		icon      *Icon
		padding   unit.Value
		isRounded bool
		size      unit.Value

		text string

		alignment   Alignment
		Color       color.RGBA
		Background  color.RGBA
		borderColor color.RGBA
	}

	Icon struct {
		imgSize int
		img     image.Image
		src     []byte
		color   color.RGBA
		op      paint.ImageOp
		material.Icon
	}
)

const (
	defaultButtonPadding = 10
)

// NewIcon returns a new Icon from IconVG data.
func NewIcon(data []byte) (*Icon, error) {
	_, err := iconvg.DecodeMetadata(data)
	if err != nil {
		return nil, err
	}
	return &Icon{
		src:   data,
		color: helper.WhiteColor,
	}, nil
}

func NewButton(txt string, icon *Icon) *Button {
	theme := helper.GetTheme()

	btn := &Button{
		button:      new(widget.Button),
		icon:        icon,
		text:        txt,
		padding:     unit.Dp(defaultButtonPadding),
		Background:  helper.DecredDarkBlueColor,
		Color:       helper.WhiteColor,
		alignment:   AlignLeft,
		borderColor: helper.BackgroundColor,
	}

	if icon != nil {
		btn.padding = unit.Dp(13)
		btn.size = unit.Dp(46)
	}

	btn.Background = theme.Color.Primary
	btn.Color = helper.WhiteColor

	return btn
}

func (b *Button) SetText(txt string) *Button {
	b.text = txt
	return b
}

func (b *Button) SetPadding(padding int) *Button {
	b.padding = unit.Dp(float32(padding))
	return b
}

func (b *Button) SetSize(size int) *Button {
	b.size = unit.Dp(float32(size))
	return b
}

func (b *Button) SetBackgroundColor(color color.RGBA) *Button {
	b.Background = color
	return b
}

func (b *Button) SetBorderColor(color color.RGBA) *Button {
	b.borderColor = color
	return b
}

func (b *Button) SetColor(color color.RGBA) *Button {
	b.Color = color
	if b.icon != nil {
		b.icon.color = color
	}

	return b
}

func (b *Button) MakeRound() *Button {
	b.isRounded = true
	return b
}

func (b *Button) SetAlignment(alignment Alignment) *Button {
	b.alignment = alignment
	return b
}

func (b *Button) Draw(ctx *layout.Context, onClick func()) {
	theme := helper.GetTheme()
	for b.button.Clicked(ctx) {
		onClick()
	}

	if b.icon != nil {
		b.drawIconButton(ctx, theme)
		return
	}

	col := b.Color
	bgcol := b.Background
	hmin := ctx.Constraints.Width.Min
	vmin := ctx.Constraints.Height.Min
	layout.Stack{Alignment: layout.Center}.Layout(ctx,
		layout.Expanded(func() {
			minWidth := ctx.Constraints.Width.Min
			minHeight := ctx.Constraints.Height.Min
			rr := float32(ctx.Px(unit.Dp(4)))

			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(minWidth),
					Y: float32(minHeight),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(ctx.Ops).Add(ctx.Ops)
			Fill(ctx, b.borderColor)

			layout.Align(layout.Center).Layout(ctx, func() {
				ctx.Constraints.Width.Min = minWidth - 2
				ctx.Constraints.Height.Min = minHeight - 2

				clip.Rect{
					Rect: f32.Rectangle{Max: f32.Point{
						X: float32(ctx.Constraints.Width.Min),
						Y: float32(ctx.Constraints.Height.Min),
					}},
					NE: rr, NW: rr, SE: rr, SW: rr,
				}.Op(ctx.Ops).Add(ctx.Ops)
				Fill(ctx, bgcol)
				for _, c := range b.button.History() {
					drawInk(ctx, c)
				}
			})
		}),
		layout.Stacked(func() {
			ctx.Constraints.Width.Min = hmin
			ctx.Constraints.Height.Min = vmin
			layout.Align(layout.Center).Layout(ctx, func() {
				layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(ctx, func() {
					paint.ColorOp{Color: col}.Add(ctx.Ops)
					widget.Label{}.Layout(ctx, theme.Shaper, theme.Fonts.Bold, b.text)
				})
			})
			pointer.Rect(image.Rectangle{Max: ctx.Dimensions.Size}).Add(ctx.Ops)
			b.button.Layout(ctx)
		}),
	)
}

func (b *Button) drawIconButton(ctx *layout.Context, theme *helper.Theme) {
	iconSize := ctx.Px(b.size) - 2*ctx.Px(b.padding)

	layout.Stack{}.Layout(ctx,
		layout.Expanded(func() {
			size := float32(ctx.Constraints.Width.Min)
			rr := float32(4)
			if b.isRounded {
				rr = size * .5
			}
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{X: size, Y: size}},
				NE:   rr, NW: rr, SE: rr, SW: rr,
			}.Op(ctx.Ops).Add(ctx.Ops)
			Fill(ctx, b.Background)

			rect := image.Rectangle{
				Max: ctx.Dimensions.Size,
			}
			pointer.Rect(rect).Add(ctx.Ops)
			b.button.Layout(ctx)

			for _, c := range b.button.History() {
				drawInk(ctx, c)
			}
		}),
		layout.Stacked(func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(ctx,
				layout.Rigid(func() {
					layout.UniformInset(b.padding).Layout(ctx, func() {
						ico := b.icon.image(iconSize)
						ico.Add(ctx.Ops)
						paint.PaintOp{
							Rect: f32.Rectangle{Max: toPointF(ico.Size())},
						}.Add(ctx.Ops)
					})
				}),
				layout.Rigid(func() {
					if b.text != "" {
						layout.UniformInset(b.padding).Layout(ctx, func() {
							paint.ColorOp{Color: b.Color}.Add(ctx.Ops)
							widget.Label{}.Layout(ctx, theme.Shaper, theme.Fonts.Bold, b.text)
						})
					}
				}),
			)
		}),
	)
}

func (ic *Icon) SetColor(col color.RGBA) *Icon {
	ic.color = col
	return ic
}

func (ic *Icon) Draw(ctx *layout.Context, size int) {
	ico := ic.image(size)
	ico.Add(ctx.Ops)
	paint.PaintOp{
		Rect: f32.Rectangle{Max: toPointF(ico.Size())},
	}.Add(ctx.Ops)
}

func (ic *Icon) image(sz int) paint.ImageOp {
	if sz == ic.imgSize {
		return ic.op
	}
	m, _ := iconvg.DecodeMetadata(ic.src)
	dx, dy := m.ViewBox.AspectRatio()
	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: sz, Y: int(float32(sz) * dy / dx)}})
	var ico iconvg.Rasterizer
	ico.SetDstImage(img, img.Bounds(), draw.Src)

	m.Palette[0] = ic.color
	//color.RGBA{A: 0xff, R: 0xff, G: 0xff, B: 0xff}
	iconvg.Decode(&ico, ic.src, &iconvg.DecodeOptions{
		Palette: &m.Palette,
	})
	ic.op = paint.NewImageOp(img)
	ic.imgSize = sz
	return ic.op
}

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func drawInk(ctx *layout.Context, c widget.Click) {
	d := ctx.Now().Sub(c.Time)
	t := float32(d.Seconds())
	const duration = 0.5
	if t > duration {
		return
	}
	t = t / duration
	var stack op.StackOp
	stack.Push(ctx.Ops)
	size := float32(ctx.Px(unit.Dp(700))) * t
	rr := size * .5
	col := byte(0xaa * (1 - t*t))
	ink := paint.ColorOp{Color: color.RGBA{A: col, R: col, G: col, B: col}}
	ink.Add(ctx.Ops)
	op.TransformOp{}.Offset(c.Position).Offset(f32.Point{
		X: -rr,
		Y: -rr,
	}).Add(ctx.Ops)
	clip.Rect{
		Rect: f32.Rectangle{Max: f32.Point{
			X: float32(size),
			Y: float32(size),
		}},
		NE: rr, NW: rr, SE: rr, SW: rr,
	}.Op(ctx.Ops).Add(ctx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(size), Y: float32(size)}}}.Add(ctx.Ops)
	stack.Pop()
	op.InvalidateOp{}.Add(ctx.Ops)
}

func Fill(ctx *layout.Context, col color.RGBA) {
	cs := ctx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(ctx.Ops)
	paint.PaintOp{Rect: dr}.Add(ctx.Ops)
	ctx.Dimensions = layout.Dimensions{Size: d}
}

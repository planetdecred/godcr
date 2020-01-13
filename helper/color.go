package helper

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

var (
	WhiteColor = color.RGBA{255, 255, 255, 255}
	BlackColor = color.RGBA{0, 0, 0, 255}
	GrayColor  = color.RGBA{128, 128, 128, 255}

	DangerColor  = color.RGBA{215, 58, 73, 255}
	SuccessColor = color.RGBA{227, 98, 9, 255}

	DecredDarkBlueColor  = color.RGBA{9, 20, 64, 255}
	DecredLightBlueColor = color.RGBA{41, 112, 255, 255}

	DecredOrangeColor = color.RGBA{237, 109, 71, 255}
	DecredGreenColor  = color.RGBA{46, 214, 161, 255} //color.RGBA{65, 191, 83, 255}

	BackgroundColor = color.RGBA{248, 249, 250, 255}
)

func PaintArea(ctx *layout.Context, color color.RGBA, x int, y int) {
	borderRadius := float32(6)
	borderWidth := 1
	if y < 21 {
		borderRadius = float32(4)
		borderWidth = 0
	}

	clip.Rect{
		Rect: f32.Rectangle{
			Max: f32.Point{
				X: float32(x),
				Y: float32(y),
			},
		},
		NE: borderRadius,
		NW: borderRadius,
		SE: borderRadius,
		SW: borderRadius,
	}.Op(ctx.Ops).Add(ctx.Ops)
	Fill(ctx, GrayColor, x, y)

	innerWidth := x - borderWidth
	innerHeight := y - borderWidth

	clip.Rect{
		Rect: f32.Rectangle{
			Max: f32.Point{
				X: float32(innerWidth),
				Y: float32(innerHeight),
			},
			Min: f32.Point{
				X: float32(borderWidth),
				Y: float32(borderWidth),
			},
		},
		NE: borderRadius,
		NW: borderRadius,
		SE: borderRadius,
		SW: borderRadius,
	}.Op(ctx.Ops).Add(ctx.Ops)
	Fill(ctx, color, innerWidth, innerHeight)
}

func PaintCircle(ctx *layout.Context, color color.RGBA, size float32) {
	borderRadius := size * .5

	clip.Rect{
		Rect: f32.Rectangle{
			Max: f32.Point{
				X: float32(size),
				Y: float32(size),
			},
		},
		NE: borderRadius,
		NW: borderRadius,
		SE: borderRadius,
		SW: borderRadius,
	}.Op(ctx.Ops).Add(ctx.Ops)
	Fill(ctx, color, int(size), int(size))
}

func PaintFooter(ctx *layout.Context, color color.RGBA, x int, y int) {
	paint.ColorOp{
		Color: GrayColor,
	}.Add(ctx.Ops)

	paint.PaintOp{
		Rect: f32.Rectangle{
			Max: f32.Point{
				X: float32(x),
				Y: float32(y),
			},
		},
	}.Add(ctx.Ops)

	paint.ColorOp{
		Color: color,
	}.Add(ctx.Ops)

	paint.PaintOp{
		Rect: f32.Rectangle{
			Max: f32.Point{
				X: float32(x),
				Y: float32(y),
			},
			Min: f32.Point{
				Y: 1,
			},
		},
	}.Add(ctx.Ops)
}

func Fill(ctx *layout.Context, col color.RGBA, x, y int) {
	//cs := ctx.Constraints
	d := image.Point{X: x, Y: y}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(ctx.Ops)
	paint.PaintOp{Rect: dr}.Add(ctx.Ops)
	ctx.Dimensions = layout.Dimensions{Size: d}
}

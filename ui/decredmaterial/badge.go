package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Badge struct {
	Background color.NRGBA
}

func (t *Theme) Badge() *Badge {
	return &Badge{
		Background: t.Color.Primary,
	}
}

func (b *Badge) Layout(gtx C, label Label) D {
	min := image.Point{
		X: 23,
		Y: 23,
	}
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			radius := 0.5 * float32(min.X)

			clip.UniformRRect(f32.Rectangle{Max: f32.Point{
				X: float32(min.X),
				Y: float32(min.Y),
			}}, radius).Add(gtx.Ops)
			paint.Fill(gtx.Ops, b.Background)
			return layout.Dimensions{Size: min}
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints.Min = min
			return layout.Center.Layout(gtx, func(gtx C) D {
				return label.Layout(gtx)
			})
		}),
	)
}

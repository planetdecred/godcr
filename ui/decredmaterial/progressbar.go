// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type ProgressBarStyle struct {
	Radius    CornerRadius
	Height    unit.Value
	Width     unit.Value
	Direction layout.Direction
	material.ProgressBarStyle
}

func (t *Theme) ProgressBar(progress int) ProgressBarStyle {
	return ProgressBarStyle{ProgressBarStyle: material.ProgressBar(t.Base, float32(progress)/100)}
}

// This achieves a progress bar using linear layouts.
func (p ProgressBarStyle) Layout2(gtx C) D {
	if p.Width.V <= 0 {
		p.Width = unit.Px(float32(gtx.Constraints.Max.X))
	}

	return p.Direction.Layout(gtx, func(gtx C) D {
		return LinearLayout{
			Width:      gtx.Px(p.Width),
			Height:     gtx.Px(p.Height),
			Background: p.TrackColor,
			Border:     Border{Radius: p.Radius},
		}.Layout2(gtx, func(gtx C) D {

			return LinearLayout{
				Width:      int(p.Width.V * clamp1(p.Progress)),
				Height:     gtx.Px(p.Height),
				Background: p.Color,
				Border:     Border{Radius: p.Radius},
			}.Layout(gtx)
		})
	})
}

func (p ProgressBarStyle) Layout(gtx layout.Context) layout.Dimensions {
	shader := func(width float32, color color.NRGBA) layout.Dimensions {
		maxHeight := p.Height
		if p.Height.V <= 0 {
			maxHeight = unit.Dp(4)
		}

		d := image.Point{X: int(width), Y: gtx.Px(maxHeight)}
		height := float32(gtx.Px(maxHeight))

		tr := float32(gtx.Px(unit.Dp(p.Radius.TopRight)))
		tl := float32(gtx.Px(unit.Dp(p.Radius.TopLeft)))
		br := float32(gtx.Px(unit.Dp(p.Radius.BottomRight)))
		bl := float32(gtx.Px(unit.Dp(p.Radius.BottomLeft)))

		defer clip.RRect{
			Rect: f32.Rectangle{Max: f32.Pt(width, height)},
			NW:   tl, NE: tr, SE: br, SW: bl,
		}.Push(gtx.Ops).Pop()

		paint.ColorOp{Color: color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		return layout.Dimensions{Size: d}
	}

	if p.Width.V <= 0 {
		p.Width = unit.Px(float32(gtx.Constraints.Max.X))
	}

	progressBarWidth := p.Width.V
	return layout.Stack{Alignment: layout.W}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return shader(progressBarWidth, p.TrackColor)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			fillWidth := progressBarWidth * clamp1(p.Progress)
			fillColor := p.Color
			if gtx.Queue == nil {
				fillColor = Disabled(fillColor)
			}
			return shader(fillWidth, fillColor)
		}),
	)
}

// clamp1 limits v to range [0..1].
func clamp1(v float32) float32 {
	if v >= 1 {
		return 1
	} else if v <= 0 {
		return 0
	} else {
		return v
	}
}

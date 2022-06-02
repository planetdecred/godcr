// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/ui/values"
)

type ProgressBarStyle struct {
	Radius    CornerRadius
	Height    unit.Dp
	Width     unit.Dp
	Direction layout.Direction
	material.ProgressBarStyle
}

type ProgressBarItem struct {
	Value   int
	Color   color.NRGBA
	SubText string
}

// MultiLayerProgressBar shows the percentage of the mutiple progress layer
// against the total/expected progress.
type MultiLayerProgressBar struct {
	t *Theme

	items  []ProgressBarItem
	Radius CornerRadius
	Height unit.Dp
	Width  int
	total  int
}

func (t *Theme) ProgressBar(progress int) ProgressBarStyle {
	return ProgressBarStyle{ProgressBarStyle: material.ProgressBar(t.Base, float32(progress)/100)}
}

func (t *Theme) MultiLayerProgressBar(total int, items []ProgressBarItem) *MultiLayerProgressBar {
	mp := &MultiLayerProgressBar{
		t: t,

		total:  total,
		Height: values.MarginPadding8,
		items:  items,
	}

	return mp
}

// This achieves a progress bar using linear layouts.
func (p ProgressBarStyle) Layout2(gtx C) D {
	if p.Width <= 0 {
		p.Width = unit.Dp(gtx.Constraints.Max.X)
	}

	return p.Direction.Layout(gtx, func(gtx C) D {
		return LinearLayout{
			Width:      gtx.Dp(p.Width),
			Height:     gtx.Dp(p.Height),
			Background: p.TrackColor,
			Border:     Border{Radius: p.Radius},
		}.Layout2(gtx, func(gtx C) D {

			return LinearLayout{
				Width:      int(float32(p.Width) * clamp1(p.Progress)),
				Height:     gtx.Dp(p.Height),
				Background: p.Color,
				Border:     Border{Radius: p.Radius},
			}.Layout(gtx)
		})
	})
}

func (p ProgressBarStyle) Layout(gtx layout.Context) layout.Dimensions {
	shader := func(width int, color color.NRGBA) layout.Dimensions {
		maxHeight := p.Height
		if p.Height <= 0 {
			maxHeight = unit.Dp(4)
		}

		d := image.Point{X: int(width), Y: gtx.Dp(maxHeight)}
		height := gtx.Dp(maxHeight)

		tr := gtx.Dp(unit.Dp(p.Radius.TopRight))
		tl := gtx.Dp(unit.Dp(p.Radius.TopLeft))
		br := gtx.Dp(unit.Dp(p.Radius.BottomRight))
		bl := gtx.Dp(unit.Dp(p.Radius.BottomLeft))

		defer clip.RRect{
			Rect: image.Rectangle{Max: image.Pt(width, height)},
			NW:   tl, NE: tr, SE: br, SW: bl,
		}.Push(gtx.Ops).Pop()

		paint.ColorOp{Color: color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		return layout.Dimensions{Size: d}
	}

	if p.Width <= 0 {
		p.Width = unit.Dp(gtx.Constraints.Max.X)
	}

	progressBarWidth := int(p.Width)
	return layout.Stack{Alignment: layout.W}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return shader(progressBarWidth, p.TrackColor)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			fillWidth := int(float32(progressBarWidth) * clamp1(p.Progress))
			fillColor := p.Color
			if gtx.Queue == nil {
				fillColor = Disabled(fillColor)
			}
			return shader(fillWidth, fillColor)
		}),
	)
}

// TODO: Allow more than just 2 layers and make it dynamic
func (mp *MultiLayerProgressBar) progressBarLayout(gtx C) D {
	r := gtx.Dp(values.MarginPadding0)
	if mp.Width <= 0 {
		mp.Width = gtx.Constraints.Max.X
	}

	// progressScale represent the different progress bar layers
	progressScale := func(width int, color color.NRGBA) layout.Dimensions {
		d := image.Point{X: int(width), Y: gtx.Dp(mp.Height)}

		defer clip.RRect{
			Rect: image.Rectangle{Max: image.Point{X: width, Y: gtx.Dp(mp.Height)}},
			NE:   r, NW: r, SE: r, SW: r,
		}.Push(gtx.Ops).Pop()

		paint.ColorOp{Color: color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		return layout.Dimensions{
			Size: d,
		}
	}

	calProgressWidth := func(progress int) int {
		val := (progress / mp.total) * 100
		return (mp.Width / 100) * val
	}

	// This takes only 2 layers
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			width := calProgressWidth(mp.items[0].Value)
			if width == 0 {
				return D{}
			}
			return progressScale(width, mp.items[0].Color)
		}),
		layout.Rigid(func(gtx C) D {
			width := calProgressWidth(mp.items[1].Value)
			if width == 0 {
				return D{}
			}
			return progressScale(width, mp.items[1].Color)
		}),
	)
}

func (mp *MultiLayerProgressBar) Layout(gtx C, labelWdg layout.Widget) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(labelWdg),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, mp.progressBarLayout)
		}),
	)
}

// clamp1 limits mp.to range [0..1].
func clamp1(v float32) float32 {
	if v >= 1 {
		return 1
	} else if v <= 0 {
		return 0
	} else {
		return v
	}
}

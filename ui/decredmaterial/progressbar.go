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
	Height unit.Value
	material.ProgressBarStyle
}

func (t *Theme) ProgressBar(progress int) ProgressBarStyle {
	return ProgressBarStyle{ProgressBarStyle: material.ProgressBar(t.Base, float32(progress)/100)}
}

func (p ProgressBarStyle) Layout(gtx layout.Context) layout.Dimensions {
	shader := func(width float32, color color.NRGBA) layout.Dimensions {
		maxHeight := unit.Dp(4)
		if p.Height.V > 0 {
			maxHeight = p.Height
		}
		rr := float32(gtx.Px(unit.Dp(2)))

		d := image.Point{X: int(width), Y: gtx.Px(maxHeight)}

		height := float32(gtx.Px(maxHeight))
		clip.UniformRRect(f32.Rectangle{Max: f32.Pt(width, height)}, rr).Add(gtx.Ops)
		paint.ColorOp{Color: color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		return layout.Dimensions{Size: d}
	}

	progressBarWidth := float32(gtx.Constraints.Max.X)
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

// Disabled blends color towards the luminance and multiplies alpha.
// Blending towards luminance will desaturate the color.
// Multiplying alpha blends the color together more with the background.
func Disabled(c color.NRGBA) (d color.NRGBA) {
	const r = 80 // blend ratio
	lum := approxLuminance(c)
	return color.NRGBA{
		R: byte((int(c.R)*r + int(lum)*(256-r)) / 256),
		G: byte((int(c.G)*r + int(lum)*(256-r)) / 256),
		B: byte((int(c.B)*r + int(lum)*(256-r)) / 256),
		A: byte(int(c.A) * (128 + 32) / 256),
	}
}

// approxLuminance is a fast approximate version of RGBA.Luminance.
func approxLuminance(c color.NRGBA) byte {
	const (
		r = 13933 // 0.2126 * 256 * 256
		g = 46871 // 0.7152 * 256 * 256
		b = 4732  // 0.0722 * 256 * 256
		t = r + g + b
	)
	return byte((r*int(c.R) + g*int(c.G) + b*int(c.B)) / t)
}

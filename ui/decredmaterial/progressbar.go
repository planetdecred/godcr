// SPDX-License-Identifier: Unlicense OR MIT

package decredmaterial

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	// "gioui.org/widget"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/ui/values"
)

type ProgressBarStyle struct {
	Radius    CornerRadius
	Height    unit.Value
	Width     unit.Value
	Direction layout.Direction
	material.ProgressBarStyle
}

type ProgressBarItem struct {
	Value   float32
	Color   color.NRGBA
	SubText string
}

// VoteBar shows the range/percentage of the yes votes and no votes against the total required.
type MultiLayerProgressBar struct {
	t *Theme

	items  []ProgressBarItem
	Radius CornerRadius
	Height unit.Value
	Width  float32
	total  float32
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

func (t *Theme) MultiLayerProgressBar(total float32, items []ProgressBarItem) *MultiLayerProgressBar {
	mp := &MultiLayerProgressBar{
		t: t,

		total:  total,
		Height: values.MarginPadding8,
		items:  items,
	}

	return mp
}

func (mp *MultiLayerProgressBar) progressBarLayout(gtx C) D {
	r := float32(gtx.Px(values.MarginPadding0))
	mp.Width = float32(gtx.Constraints.Max.X)

	// progressScale represent the different progress bar layers
	progressScale := func(width float32, color color.NRGBA) layout.Dimensions {
		d := image.Point{X: int(width), Y: gtx.Px(mp.Height)}

		defer clip.RRect{
			Rect: f32.Rectangle{Max: f32.Point{X: width, Y: float32(gtx.Px(mp.Height))}},
			NE:   r, NW: r, SE: r, SW: r,
		}.Push(gtx.Ops).Pop()

		paint.ColorOp{Color: color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		return layout.Dimensions{
			Size: d,
		}
	}

	return layout.Stack{Alignment: layout.W}.Layout(gtx,
		// layout.Stacked(func(gtx layout.Context) layout.Dimensions {
		// 	return progressScale(mp.Width, mp.t.Color.Gray2)
		// }),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			// progressLayer := make([]layout.Widget, 0)
			var w []layout.Widget

			for _, item := range mp.items {
				val := (item.Value / mp.total) * 100
				width := (mp.Width / 100) * val
				fmt.Println(val)
				fmt.Println(width)
				fmt.Println(item)
				fmt.Println(mp.Width)

				w = append(w, func(gtx C) D {
					if width == 0 {
						return D{}
					}

					// return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					// layout.Rigid(mp.t.Label(values.TextSize14, item.SubText).Layout),
					// layout.Rigid(func(gtx C) D {
					return progressScale(width, item.Color)
					// })
					// )
				})

			}

			// fmt.Println(len(progressLayer))
			// return layout.Flex{}.Layout(gtx, progressLayer...)
			list := &layout.List{Axis: layout.Horizontal}
			return list.Layout(gtx, len(w), func(gtx C, i int) D {
				return w[i](gtx)
			})

		}),
	)
}

func (mp *MultiLayerProgressBar) Layout(gtx C) D {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding5, Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					// layout.Rigid(func(gtx C) D {
					// 	return layout.Flex{}.Layout(gtx,
					// 		layout.Rigid(func(gtx C) D {
					// 			yesLabel := mp.Theme.Body1("Yes: ")
					// 			return mp.layoutIconAndText(gtx, yesLabel, mp.yesVotes, mp.yesColor)
					// 		}),
					// 		layout.Rigid(func(gtx C) D {
					// 			noLabel := mp.Theme.Body1("No: ")
					// 			return mp.layoutIconAndText(gtx, noLabel, mp.noVotes, mp.noColor)
					// 		}),
					// 		layout.Flexed(1, func(gtx C) D {
					// 			return layout.E.Layout(gtx, func(gtx C) D {
					// 				return mp.layoutInfo(gtx)
					// 			})
					// 		}),
					// 	)
					// }),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, mp.progressBarLayout)
					}),
				)
			})
		}),
	)
}

// func (v *VoteBar) layoutIconAndText(gtx C, lbl decredmaterial.Label, count float32, col color.NRGBA) D {
// 	return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
// 		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
// 			layout.Rigid(func(gtx C) D {
// 				return layout.Inset{Right: values.MarginPadding5, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
// 					mp.legendIcon.Color = col
// 					return mp.legendIcon.Layout(gtx, values.MarginPadding10)
// 				})
// 			}),
// 			layout.Rigid(func(gtx C) D {
// 				lbl.Font.Weight = text.SemiBold
// 				return lbl.Layout(gtx)
// 			}),
// 			layout.Rigid(func(gtx C) D {
// 				percentage := (count / mp.totalVotes) * 100
// 				if percentage != percentage {
// 					percentage = 0
// 				}
// 				percentageStr := strconv.FormatFloat(float64(percentage), 'f', 1, 64) + "%"
// 				countStr := strconv.FormatFloat(float64(count), 'f', 0, 64)

// 				return mp.Theme.Body1(fmt.Sprintf("%s (%s)", countStr, percentageStr)).Layout(gtx)
// 			}),
// 		)
// 	})
// }

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

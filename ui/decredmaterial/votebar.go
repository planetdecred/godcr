package decredmaterial

import (
	"image/color"
	"strconv"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

type VoteBar struct {
	yesVotes    float32
	noVotes     float32
	yesColor    color.RGBA
	noColor     color.RGBA
	legendLabel Label
}

func (t *Theme) VoteBar(yesVotes, noVotes float32) *VoteBar {
	return &VoteBar{
		yesVotes:    yesVotes,
		noVotes:     noVotes,
		yesColor:    t.Color.Success,
		noColor:     t.Color.Danger,
		legendLabel: t.Body2(""),
	}
}

func (v *VoteBar) Layout(gtx layout.Context) layout.Dimensions {
	totalVotes := v.yesVotes + v.noVotes
	yesFlex := v.yesVotes / totalVotes
	noFlex := v.noVotes / totalVotes

	gtx.Constraints.Max.Y = 7
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(yesFlex, func(gtx C) D {
			return fillMax(gtx, v.yesColor)
		}),
		layout.Flexed(noFlex, func(gtx C) D {
			return fillMax(gtx, v.noColor)
		}),
	)
}

func (v *VoteBar) LayoutWithLegend(gtx layout.Context, icon *widget.Icon) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						icon.Color = v.yesColor
						return v.layoutIconAndText(gtx, v.yesVotes, icon)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Left: unit.Dp(10),
						}.Layout(gtx, func(gtx C) D {
							icon.Color = v.noColor
							return v.layoutIconAndText(gtx, v.noVotes, icon)
						})
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return v.Layout(gtx)
		}),
	)
}

func (v *VoteBar) layoutIconAndText(gtx layout.Context, count float32, icon *widget.Icon) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Right: unit.Dp(5), Top: unit.Dp(5)}.Layout(gtx, func(gtx C) D {
				return icon.Layout(gtx, unit.Dp(10))
			})
		}),
		layout.Rigid(func(gtx C) D {
			v.legendLabel.Text = strconv.FormatFloat(float64(count), 'f', 0, 64)
			return v.legendLabel.Layout(gtx)
		}),
	)
}

package decredmaterial

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/values"
)

type VoteBar struct {
	yesVotes           float32
	noVotes            float32
	eligibleVotes      float32
	totalVotes         float32
	requiredPercentage float32
	passPercentage     float32
	yesColor           color.NRGBA
	noColor            color.NRGBA
	bgColor            color.NRGBA

	yesLabel                      Label
	noLabel                       Label
	legendLabel                   Label
	totalVotesLabel               Label
	requirementLabel              Label
	passTooltipLabel              Label
	totalVotesTooltipLabel        Label
	quorumRequirementTooltipLabel Label
	totalVotesCountLabel          Label
	quorumRequirementCountLabel   Label

	passTooltip   *Tooltip
	quorumTooltip *Tooltip
	infoIcon      *widget.Icon
	legendIcon    *widget.Icon
}

const (
	voteBarHeight      = 8
	voteBarRadius      = 5
	voteBarThumbWidth  = 4
	voteBarThumbHeight = 15
)

func (t *Theme) VoteBar(infoIcon, legendIcon *widget.Icon) VoteBar {
	return VoteBar{
		yesColor:                      t.Color.Success,
		noColor:                       t.Color.Danger,
		yesLabel:                      t.Body2("Yes:"),
		noLabel:                       t.Body2("No:"),
		legendLabel:                   t.Body2(""),
		requirementLabel:              t.Body2(""),
		totalVotesLabel:               t.Body2(""),
		passTooltip:                   t.Tooltip(),
		quorumTooltip:                 t.Tooltip(),
		infoIcon:                      infoIcon,
		legendIcon:                    legendIcon,
		bgColor:                       t.Color.Gray,
		passTooltipLabel:              t.Caption(""),
		totalVotesTooltipLabel:        t.Caption("Total votes"),
		quorumRequirementTooltipLabel: t.Caption("Quorum requirement"),
		totalVotesCountLabel:          t.Caption(""),
		quorumRequirementCountLabel:   t.Caption(""),
	}
}

func (v *VoteBar) SetParams(yesVotes, noVotes, eligibleVotes, requiredPercentage, passPercentage float32) *VoteBar {
	totalVotes := yesVotes + noVotes

	v.yesVotes = yesVotes
	v.noVotes = noVotes
	v.eligibleVotes = eligibleVotes
	v.passPercentage = passPercentage
	v.totalVotes = totalVotes
	v.requiredPercentage = requiredPercentage
	v.totalVotesLabel.Text = fmt.Sprintf("%d", int(totalVotes))
	v.passTooltipLabel.Text = fmt.Sprintf("%d %% Yes votes required for approval", int(v.passPercentage))
	v.totalVotesCountLabel.Text = strconv.FormatFloat(float64(totalVotes), 'f', 0, 64)

	return v
}

func (v *VoteBar) Layout(gtx C) D {
	yesRatio := (v.yesVotes / v.totalVotes) * float32(gtx.Constraints.Max.X)

	// draw yes voteBar
	st := op.Save(gtx.Ops)
	rrect := f32.Rectangle{
		Max: f32.Point{
			X: yesRatio,
			Y: voteBarHeight,
		},
	}
	clip.RRect{
		Rect: rrect,
		NW:   voteBarRadius,
		SW:   voteBarRadius,
	}.Add(gtx.Ops)
	paint.Fill(gtx.Ops, v.yesColor)
	st.Load()

	// draw no voteBar
	st = op.Save(gtx.Ops)
	rrect = f32.Rectangle{
		Min: f32.Point{
			X: yesRatio,
		},
		Max: f32.Point{
			X: float32(gtx.Constraints.Max.X),
			Y: voteBarHeight,
		},
	}
	clip.RRect{
		Rect: rrect,
		NE:   voteBarRadius,
		SE:   voteBarRadius,
	}.Add(gtx.Ops)
	paint.Fill(gtx.Ops, v.noColor)
	st.Load()

	st = op.Save(gtx.Ops)
	thumbLeftPos := (v.passPercentage / 100) * float32(gtx.Constraints.Max.X)
	rect := image.Rectangle{
		Min: image.Point{
			X: int(thumbLeftPos - float32(voteBarThumbWidth/2)),
			Y: -voteBarThumbHeight + 7,
		},
		Max: image.Point{
			X: int(thumbLeftPos + voteBarThumbWidth),
			Y: voteBarThumbHeight,
		},
	}
	clip.Rect(rect).Add(gtx.Ops)
	paint.Fill(gtx.Ops, v.bgColor)
	inset := layout.Inset{Left: unit.Dp(110)}
	v.passTooltip.Layout(gtx, rect, inset, func(gtx C) D {
		return v.passTooltipLabel.Layout(gtx)
	})
	st.Load()

	return D{
		Size: image.Point{
			X: gtx.Constraints.Max.X,
			Y: voteBarHeight,
		},
	}
}

func (v *VoteBar) LayoutWithLegend(gtx C) D {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding5, Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								v.legendIcon.Color = v.yesColor
								return v.layoutIconAndText(gtx, v.yesLabel, v.yesVotes)
							}),
							layout.Rigid(func(gtx C) D {
								v.legendIcon.Color = v.noColor
								return v.layoutIconAndText(gtx, v.noLabel, v.noVotes)
							}),
							layout.Flexed(1, func(gtx C) D {
								return layout.E.Layout(gtx, func(gtx C) D {
									return v.layoutInfo(gtx)
								})
							}),
						)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, v.Layout)
					}),
				)
			})
		}),
	)
}

func (v *VoteBar) layoutIconAndText(gtx C, lbl Label, count float32) D {
	return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: values.MarginPadding5, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return v.legendIcon.Layout(gtx, values.MarginPadding10)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return lbl.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				percentage := (count / v.totalVotes) * 100
				percentageStr := strconv.FormatFloat(float64(percentage), 'f', 1, 64) + "%"
				countStr := strconv.FormatFloat(float64(count), 'f', 0, 64)

				v.legendLabel.Text = fmt.Sprintf("%s (%s)", countStr, percentageStr)
				return v.legendLabel.Layout(gtx)
			}),
		)
	})
}

func (v *VoteBar) layoutInfo(gtx C) D {
	quorumRequirement := (v.requiredPercentage / 100) * v.eligibleVotes
	v.requirementLabel.Text = fmt.Sprintf("/%d votes", int(quorumRequirement))
	v.quorumRequirementCountLabel.Text = strconv.FormatFloat(float64(quorumRequirement), 'f', 0, 64)

	dims := layout.Flex{}.Layout(gtx,
		layout.Rigid(v.totalVotesLabel.Layout),
		layout.Rigid(v.requirementLabel.Layout),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				rect := image.Rectangle{
					Min: gtx.Constraints.Min,
					Max: gtx.Constraints.Max,
				}
				rect.Max.Y = 20
				v.layoutInfoTooltip(gtx, rect)
				return v.infoIcon.Layout(gtx, unit.Dp(20))
			})
		}),
	)

	return dims
}

func (v *VoteBar) layoutInfoTooltip(gtx C, rect image.Rectangle) {
	inset := layout.Inset{Left: unit.Dp(-165)}
	v.quorumTooltip.Layout(gtx, rect, inset, func(gtx C) D {
		gtx.Constraints.Min.X = 150
		gtx.Constraints.Max.X = 150
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(v.totalVotesTooltipLabel.Layout),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, v.totalVotesCountLabel.Layout)
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(v.quorumRequirementTooltipLabel.Layout),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, v.quorumRequirementCountLabel.Layout)
					}),
				)
			}),
		)
	})
}

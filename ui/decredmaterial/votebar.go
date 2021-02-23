package decredmaterial

import (
	"fmt"
	"image/color"
	"strconv"

	"gioui.org/layout"
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

	yesLabel         Label
	noLabel          Label
	legendLabel      Label
	totalVotesLabel  Label
	requirementLabel Label

	passTooltip   *Tooltip
	quorumTooltip *Tooltip
	infoIcon      *widget.Icon
	legendIcon    *widget.Icon
}

const (
	voteBarHeight = 8
)

func (t *Theme) VoteBar(infoIcon, legendIcon *widget.Icon) VoteBar {
	return VoteBar{
		yesColor:         t.Color.Success,
		noColor:          t.Color.Danger,
		yesLabel:         t.Body2("Yes:"),
		noLabel:          t.Body2("No:"),
		legendLabel:      t.Body2(""),
		requirementLabel: t.Body2(""),
		totalVotesLabel:  t.Body2(""),
		passTooltip:      t.Tooltip("", Left),
		quorumTooltip:    t.Tooltip("", Left),
		infoIcon:         infoIcon,
		legendIcon:       legendIcon,
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

	return v
}

func (v *VoteBar) Layout(gtx C) D {
	yesFlex := float32(v.yesVotes / v.totalVotes)
	noFlex := float32(v.noVotes / v.totalVotes)

	yesRadius := CornerRadius{
		SW: 5,
		NW: 5,
	}

	noRadius := CornerRadius{
		SE: 5,
		NE: 5,
	}

	gtx.Constraints.Max.Y = voteBarHeight
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(yesFlex, func(gtx C) D {
					return fillMax(gtx, v.yesColor, yesRadius)
				}),
				layout.Flexed(noFlex, func(gtx C) D {
					return fillMax(gtx, v.noColor, noRadius)
				}),
			)
		}),
	)
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
		layout.Expanded(func(gtx C) D {
			leftPos := (v.passPercentage / 100) * float32(gtx.Constraints.Max.X)
			return layout.Inset{
				Left: unit.Dp(leftPos),
				Top: values.MarginPadding20,
			}.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = 3
				gtx.Constraints.Min.Y = 28

				v.passTooltip.SetText(fmt.Sprintf("%d %% Yes votes required for approval", int(v.passPercentage)))
				return v.passTooltip.Layout(gtx, func(gtx C) D {
					return fill(gtx, v.yesColor)
				})
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

	v.quorumTooltip.SetText(fmt.Sprintf("%d votes cast, quorum requirement is %d", int(v.totalVotes), int(quorumRequirement)))
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(v.totalVotesLabel.Layout),
		layout.Rigid(v.requirementLabel.Layout),
		layout.Rigid(func(gtx C) D {
			return v.quorumTooltip.Layout(gtx, func(gtx C) D {
				return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return v.infoIcon.Layout(gtx, unit.Dp(20))
				})
			})
		}),
	)
}

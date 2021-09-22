package proposal

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

// VoteBar widget implements voting stat for proposals.
// VoteBar shows the range/percentage of the yes votes and no votes against the total required.
type VoteBar struct {
	*load.Load

	yesVotes           float32
	noVotes            float32
	eligibleVotes      float32
	totalVotes         float32
	requiredPercentage float32
	passPercentage     float32

	yesColor color.NRGBA
	noColor  color.NRGBA

	passTooltip   *decredmaterial.Tooltip
	quorumTooltip *decredmaterial.Tooltip

	infoIcon   *widget.Icon
	legendIcon *widget.Icon
}

var voteBarThumbWidth = 2

func NewVoteBar(l *load.Load) *VoteBar {
	vb := &VoteBar{
		Load: l,

		yesColor:      l.Theme.Color.Success,
		noColor:       l.Theme.Color.Danger,
		passTooltip:   l.Theme.Tooltip(),
		quorumTooltip: l.Theme.Tooltip(),
	}

	vb.infoIcon = l.Icons.ActionInfo
	vb.infoIcon.Color = l.Theme.Color.Gray

	vb.legendIcon = l.Icons.ImageBrightness1
	vb.legendIcon.Color = l.Theme.Color.InactiveGray

	return vb
}

func (v *VoteBar) SetYesNoVoteParams(yesVotes, noVotes float32) *VoteBar {
	v.yesVotes = yesVotes
	v.noVotes = noVotes

	v.totalVotes = yesVotes + noVotes

	return v
}

func (v *VoteBar) SetVoteValidityParams(eligibleVotes, requiredPercentage, passPercentage float32) *VoteBar {
	v.eligibleVotes = eligibleVotes
	v.passPercentage = passPercentage
	v.requiredPercentage = requiredPercentage

	return v
}

func (v *VoteBar) Layout(gtx C) D {
	var rW, rE float32
	r := float32(gtx.Px(unit.Dp(4)))
	progressBarWidth := float32(gtx.Constraints.Max.X)
	quorumRequirement := (v.requiredPercentage / 100) * v.eligibleVotes

	yesVotes := (v.yesVotes / quorumRequirement) * 100
	noVotes := (v.noVotes / quorumRequirement) * 100
	yesWidth := (progressBarWidth / 100) * yesVotes
	noWidth := (progressBarWidth / 100) * noVotes

	// progressScale represent the different progress bar layers
	progressScale := func(width float32, color color.NRGBA, layer int) layout.Dimensions {
		maxHeight := unit.Dp(8)
		rW, rE = 0, 0
		if layer == 2 {
			if width >= progressBarWidth {
				rE = r
			}
			rW = r
		} else if layer == 3 {
			if v.yesVotes == 0 {
				rW = r
			}
			rE = r
		} else {
			rE, rW = r, r
		}
		d := image.Point{X: int(width), Y: gtx.Px(maxHeight)}

		defer clip.RRect{
			Rect: f32.Rectangle{Max: f32.Point{X: width, Y: float32(gtx.Px(maxHeight))}},
			NE:   rE, NW: rW, SE: rE, SW: rW,
		}.Push(gtx.Ops).Pop()

		paint.ColorOp{Color: color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		return layout.Dimensions{
			Size: d,
		}
	}

	if yesWidth > progressBarWidth || noWidth > progressBarWidth || (yesWidth+noWidth) > progressBarWidth {
		yes := (v.yesVotes / v.totalVotes) * 100
		no := (v.noVotes / v.totalVotes) * 100
		noWidth = (progressBarWidth / 100) * no
		yesWidth = (progressBarWidth / 100) * yes
		rE = r
	} else if yesWidth < 0 {
		yesWidth, noWidth = 0, 0
	}

	return layout.Stack{Alignment: layout.W}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return progressScale(progressBarWidth, v.Theme.Color.Gray1, 1)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if yesWidth == 0 {
						return D{}
					}
					return progressScale(yesWidth, v.yesColor, 2)
				}),
				layout.Rigid(func(gtx C) D {
					if noWidth == 0 {
						return D{}
					}
					return progressScale(noWidth, v.noColor, 3)
				}),
			)
		}),
		layout.Stacked(v.requiredYesVotesIndicator),
	)
}

func (v *VoteBar) votesIndicatorTooltip(gtx C, r image.Rectangle, tipPos float32) {
	insetLeft := tipPos - float32(voteBarThumbWidth/2) - 205
	inset := layout.Inset{Left: unit.Dp(insetLeft), Top: unit.Dp(25)}
	v.passTooltip.Layout(gtx, r, inset, func(gtx C) D {
		txt := fmt.Sprintf("%d %% Yes votes required for approval", int(v.passPercentage))
		return v.Theme.Caption(txt).Layout(gtx)
	})
}

func (v *VoteBar) requiredYesVotesIndicator(gtx C) D {
	thumbLeftPos := (v.passPercentage / 100) * float32(gtx.Constraints.Max.X)
	rect := image.Rectangle{
		Min: image.Point{
			X: int(thumbLeftPos - float32(voteBarThumbWidth/2)),
			Y: -1,
		},
		Max: image.Point{
			X: int(int(thumbLeftPos) + voteBarThumbWidth),
			Y: 45,
		},
	}
	clip.Rect(rect).Add(gtx.Ops)
	paint.Fill(gtx.Ops, v.Theme.Color.InactiveGray)
	v.votesIndicatorTooltip(gtx, rect, thumbLeftPos)
	return D{
		Size: rect.Max,
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
								yesLabel := v.Theme.Body1("Yes: ")
								return v.layoutIconAndText(gtx, yesLabel, v.yesVotes)
							}),
							layout.Rigid(func(gtx C) D {
								v.legendIcon.Color = v.noColor
								noLabel := v.Theme.Body1("No: ")
								return v.layoutIconAndText(gtx, noLabel, v.noVotes)
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

func (v *VoteBar) layoutIconAndText(gtx C, lbl decredmaterial.Label, count float32) D {
	return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: values.MarginPadding5, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Px(unit.Dp(10))
					return v.legendIcon.Layout(gtx, clr)
				})
			}),
			layout.Rigid(func(gtx C) D {
				lbl.Font.Weight = text.Bold
				return lbl.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				percentage := (count / v.totalVotes) * 100
				if percentage != percentage {
					percentage = 0
				}
				percentageStr := strconv.FormatFloat(float64(percentage), 'f', 1, 64) + "%"
				countStr := strconv.FormatFloat(float64(count), 'f', 0, 64)

				return v.Theme.Body1(fmt.Sprintf("%s (%s)", countStr, percentageStr)).Layout(gtx)
			}),
		)
	})
}

func (v *VoteBar) layoutInfo(gtx C) D {
	quorumRequirement := (v.requiredPercentage / 100) * v.eligibleVotes

	requirement := v.Theme.Body2(fmt.Sprintf("/%d votes", int(quorumRequirement)))
	requirement.Color = v.Theme.Color.Gray

	dims := layout.Flex{}.Layout(gtx,
		layout.Rigid(v.Theme.Body2(fmt.Sprintf("%d", int(v.totalVotes))).Layout),
		layout.Rigid(requirement.Layout),
		layout.Rigid(func(gtx C) D {
			rect := image.Rectangle{
				Min: gtx.Constraints.Min,
				Max: gtx.Constraints.Max,
			}
			rect.Max.Y = 20
			v.layoutInfoTooltip(gtx, rect)
			return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Px(unit.Dp(20))
				return v.infoIcon.Layout(gtx, v.requirementLabel.Color)
			})
		}),
	)

	return dims
}

func (v *VoteBar) layoutInfoTooltip(gtx C, rect image.Rectangle) {
	inset := layout.Inset{Top: unit.Dp(20), Left: unit.Dp(-180)}

	col := v.Theme.Color.Gray4
	totalVotesTooltipLabel := v.Theme.Caption("Total votes")
	totalVotesTooltipLabel.Color = col

	totalVotesCountLabel := v.Theme.Caption(strconv.FormatFloat(float64(v.totalVotes), 'f', 0, 64))
	totalVotesCountLabel.Color = col

	quorumRequirementTooltip := v.Theme.Caption("Quorum requirement")
	quorumRequirementTooltip.Color = col

	txt := strconv.FormatFloat(float64((v.requiredPercentage/100)*v.eligibleVotes), 'f', 0, 64)
	quorumRequirementCount := v.Theme.Caption(txt)
	quorumRequirementCount.Color = col

	v.quorumTooltip.Layout(gtx, rect, inset, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Px(unit.Dp(180))
		gtx.Constraints.Max.X = gtx.Px(unit.Dp(180))
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(totalVotesTooltipLabel.Layout),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, totalVotesCountLabel.Layout)
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(quorumRequirementTooltip.Layout),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, quorumRequirementCount.Layout)
					}),
				)
			}),
		)
	})
}

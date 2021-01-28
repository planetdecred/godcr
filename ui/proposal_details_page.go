package ui

import (
	//"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"
	"runtime"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/utils"
	"github.com/planetdecred/godcr/ui/values"
)

const PageProposalDetails = "proposaldetails"

type ProposalPage struct {
	theme            *decredmaterial.Theme
	proposal         **dcrlibwallet.Proposal
	line             *decredmaterial.Line
	clickables       map[string]*widget.Clickable
	backButton       decredmaterial.IconButton
	legendIcon       *widget.Icon
	container        *layout.List
	renderedMarkdown []layout.Widget
}

func (win *Window) ProposalPage(common pageCommon) layout.Widget {
	pg := ProposalPage{
		theme:      common.theme,
		proposal:   &win.selectedProposal,
		backButton: common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack),
		legendIcon: common.icons.imageBrightness1,
		line:       common.theme.Line(),
		container:  &layout.List{Axis: layout.Vertical},
	}
	pg.backButton.Color = common.theme.Color.Hint
	pg.backButton.Size = values.MarginPadding30
	pg.line.Color = pg.theme.Color.Hint

	return func(gtx C) D {
		pg.Handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *ProposalPage) Handle(common pageCommon) {
	for pg.backButton.Button.Clicked() {
		*common.page = PageProposals
	}

	for to, c := range pg.clickables {
		//fmt.Println(c)
		for c.Clicked() {
			pg.goToURL(to)
		}
	}
}

func (pg *ProposalPage) goToURL(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Error(err)
	}
}

func (pg *ProposalPage) Layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	return c.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.W.Layout(gtx, func(gtx C) D {
						return pg.backButton.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return pg.layoutProposalDescription(gtx)
				}),
			)
		})
	})
}

func (pg *ProposalPage) layoutProposalDescription(gtx layout.Context) layout.Dimensions {
	proposal := *pg.proposal
	w := []layout.Widget{
		func(gtx C) D {
			return pg.layoutProposalHeader(gtx)
		},
		func(gtx C) D {
			return pg.layoutProposalDetailsSubHeader(gtx)
		},
		func(gtx C) D {
			category := proposal.Category
			if category == dcrlibwallet.ProposalCategoryApproved || category == dcrlibwallet.ProposalCategoryActive || category == dcrlibwallet.ProposalCategoryRejected {
				return layout.Inset{
					Top:    unit.Dp(8),
					Bottom: unit.Dp(8),
				}.Layout(gtx, func(gtx C) D {
					return pg.theme.VoteBar(float32(proposal.YesVotes), float32(proposal.NoVotes)).LayoutWithLegend(gtx, pg.legendIcon)
				})
			}
			return layout.Dimensions{}
		},
		func(gtx C) D {
			return layout.Inset{
				Top:    unit.Dp(12),
				Bottom: unit.Dp(12),
			}.Layout(gtx, func(gtx C) D {
				pg.line.Width = gtx.Constraints.Max.X
				return pg.line.Layout(gtx)
			})
		},
	}

	// add this so that the markdown document is only parsed once
	// this fixes the issue of links not triggering when clicked
	if pg.renderedMarkdown == nil {
		r := utils.RenderMarkdown(gtx, pg.theme, pg.getProposalText())
		pg.renderedMarkdown, pg.clickables = r.Layout()
	}

	w = append(w, pg.renderedMarkdown...)
	return pg.container.Layout(gtx, len(w), func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(0)).Layout(gtx, w[i])
	})
}

func (pg *ProposalPage) layoutProposalHeader(gtx layout.Context) layout.Dimensions {
	proposal := *pg.proposal

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(0.55, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return getTitleLabel(pg.theme, proposal.Name).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return getSubtitleLabel(pg.theme, proposal.Token).Layout(gtx)
				}),
			)
		}),
		layout.Flexed(0.45, func(gtx C) D {
			if proposal.Category == dcrlibwallet.ProposalCategoryPre || proposal.Category == dcrlibwallet.ProposalCategoryAbandoned {
				return layout.E.Layout(gtx, func(gtx C) D {
					return getSubtitleLabel(pg.theme, fmt.Sprintf("Last updated %s", timeAgo(proposal.Timestamp))).Layout(gtx)
				})
			}
			return layout.Dimensions{}
		}),
	)
}

func (pg *ProposalPage) layoutProposalDetailsSubHeader(gtx layout.Context) layout.Dimensions {
	proposal := *pg.proposal

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.layoutProposalDetailsSubHeaderRow(gtx, "Created by:", proposal.Username)
		}),
		layout.Rigid(func(gtx C) D {
			return pg.layoutProposalDetailsSubHeaderRow(gtx, "Version:", proposal.Version)
		}),
		layout.Rigid(func(gtx C) D {
			return pg.layoutProposalDetailsSubHeaderRow(gtx, "Last updated:", timeAgo(proposal.Timestamp))
		}),
	)
}

func (pg *ProposalPage) layoutProposalDetailsSubHeaderRow(gtx layout.Context, leftText, rightText string) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(0.03, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return getSubtitleLabel(pg.theme, leftText).Layout(gtx)
			})
		}),
		layout.Flexed(0.2, func(gtx C) D {
			return layout.Inset{
				Left: unit.Dp(4),
			}.Layout(gtx, func(gtx C) D {
				return getTitleLabel(pg.theme, rightText).Layout(gtx)
			})
		}),
	)
}

func (pg *ProposalPage) getProposalText() []byte {
	proposal := *pg.proposal
	desc, _ := base64.StdEncoding.DecodeString(proposal.IndexFile)

	return desc
}

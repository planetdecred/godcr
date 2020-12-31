package ui

import (
	//"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"
	//"regexp"
	"runtime"
	//"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	//"gioui.org/widget/material"

	//"github.com/PuerkitoBio/goquery"
	//"gitlab.com/golang-commonmark/markdown"
	//"github.com/gomarkdown/markdown"
	//"github.com/gomarkdown/markdown/parser"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/utils"
	"github.com/planetdecred/godcr/ui/values"
)

const PageProposalDetails = "proposaldetails"

type renderFunc func(C, []string) layout.Widget

type ProposalPage struct {
	theme      *decredmaterial.Theme
	proposal   **dcrlibwallet.Proposal
	line       *decredmaterial.Line
	pageLinks  map[string]layout.Dimensions
	clickables map[string]*widget.Clickable
	backButton decredmaterial.IconButton
	legendIcon *widget.Icon
	container  *layout.List
}

const (
	//markdownLinkPlaceholder = "[[link]]"
	//markdownHeadingPlaceholder = "[[heading--"

	spacer     = "@@@@"
	linkTag    = "[[link"
	headingTag = "[[heading"
)

var (
	tags = []string{"li", "ul", "td", "tr", "th", "table", "strong", "p", "h1", "h2", "h3", "h4", "h5", "h6"}
)

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
		return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx C) D {
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
			return pg.layoutProposalHeader(gtx, false)
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
					yes, no := calculateVotes(proposal.VoteSummary.OptionsResult)
					return pg.theme.VoteBar(yes, no).LayoutWithLegend(gtx, pg.legendIcon)
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

	r := utils.RenderMarkdown(gtx, pg.theme, pg.getProposalText())
	w = append(w, r.Layout()...)
	return pg.container.Layout(gtx, len(w), func(gtx C, i int) D {
		return layout.UniformInset(unit.Dp(0)).Layout(gtx, w[i])
	})
}

func (pg *ProposalPage) layoutProposalHeader(gtx layout.Context, truncateTitle bool) layout.Dimensions {
	proposal := *pg.proposal

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(0.55, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return getTitleLabel(pg.theme, proposal.Name).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return getSubtitleLabel(pg.theme, proposal.CensorshipRecord.Token).Layout(gtx)
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

var has = false

func (pg *ProposalPage) getProposalText() []byte {
	proposal := *pg.proposal

	var desc []byte
	for i := range proposal.Files {
		if proposal.Files[i].Name == "index.md" {
			desc, _ = base64.StdEncoding.DecodeString(proposal.Files[i].Payload)
			break
		}
	}

	return desc
}

package ui

import (
	"encoding/base64"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PageProposalDetails = "proposaldetails"

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

var (
	markdownRegex           = regexp.MustCompile(`\[[^][]+]\((https?://[^()]+)\)`)
	markdownLinkPlaceholder = "[[link]]"
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

	ws := pg.getProposalDescriptionTextParts(gtx)
	w = append(w, ws...)

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

func (pg *ProposalPage) getProposalDescriptionTextParts(gtx layout.Context) []layout.Widget {
	proposal := *pg.proposal

	var desc []byte
	for i := range proposal.Files {
		if proposal.Files[i].Name == "index.md" {
			desc, _ = base64.StdEncoding.DecodeString(proposal.Files[i].Payload)
			break
		}
	}

	return pg.PrepareText(gtx, string(desc))
}

func (pg *ProposalPage) PrepareText(gtx layout.Context, text string) []layout.Widget {
	// first extract all links and replace them with a placeholder
	ls := markdownRegex.FindAllStringSubmatch(text, -1)
	links := make([]string, len(ls))
	for i := range ls {
		index := i
		links[index] = ls[index][0]
	}

	// replace all links with placeholder
	txt := markdownRegex.ReplaceAllString(text, markdownLinkPlaceholder)

	prevRune := ""
	paragraphs := strings.FieldsFunc(txt, func(r rune) bool {
		match := string(r) == "\n" && prevRune == "\n"
		prevRune = string(r)

		return match
	})

	w := []layout.Widget{}
	for i := range paragraphs {
		index := i
		words := strings.Split(paragraphs[index], " ")
		w = append(w, func(gtx C) D {
			dims := decredmaterial.GridWrap{
				Axis:      layout.Horizontal,
				Alignment: layout.End,
			}.Layout(gtx, len(words), func(gtx C, i int) D {
				word := words[i]
				if words[i] == markdownLinkPlaceholder || strings.HasPrefix(words[i], markdownLinkPlaceholder) {
					link := links[0]
					links = links[1:]
					word := strings.Replace(words[i], markdownLinkPlaceholder, link, -1)
					return pg.layoutLinkWord(gtx, word)
				}
				return pg.theme.Body2(strings.TrimSpace(word) + " ").Layout(gtx)
			})
			dims.Size.Y += 20
			return dims
		})
	}
	return w
}

func (pg *ProposalPage) layoutLinkWord(gtx layout.Context, link string) layout.Dimensions {
	text, _ := getStringInBetweenTwoString(link, "[", "]")
	linkRef, _ := getStringInBetweenTwoString(link, "(", ")")

	if pg.clickables == nil {
		pg.clickables = map[string]*widget.Clickable{}
	}

	if _, ok := pg.clickables[linkRef]; !ok {
		pg.clickables[linkRef] = new(widget.Clickable)
	}

	return material.Clickable(gtx, pg.clickables[linkRef], func(gtx C) D {
		lbl := pg.theme.Body2(text + " ")
		lbl.Color = pg.theme.Color.Primary
		return lbl.Layout(gtx)
	})
}

func getStringInBetweenTwoString(str string, startS string, endS string) (result string, found bool) {
	s := strings.Index(str, startS)
	if s == -1 {
		return result, false
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return result, false
	}
	result = newS[:e]
	return result, true
}

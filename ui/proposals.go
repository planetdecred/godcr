package ui

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/ararog/timeago"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const (
	PageProposals          = "proposals"
	proposalSyncPaneHeight = 100
)

type ProposalsPage struct {
	theme                       *decredmaterial.Theme
	wallet                      *wallet.Wallet
	pageListContainer           *layout.List
	proposalListContainer       *layout.List
	tabTitles                   []string
	tabContainer                *decredmaterial.Tabs
	outline                     decredmaterial.Outline
	isShowingProposalDetails    bool
	proposals                   *map[int32][]dcrlibwallet.Proposal
	selectedProposal            **dcrlibwallet.Proposal
	clickables                  []*gesture.Click
	syncListener                *proposalSyncListener
	notSyncedIcon               *widget.Icon
	syncingIcon                 image.Image
	notSyncedStatusLabel        decredmaterial.Label
	fetchingTokenInventoryLabel decredmaterial.Label
	syncButton                  decredmaterial.Button
	cancelSyncButton            decredmaterial.Button
	syncingLabel                decredmaterial.Label
}

func (win *Window) ProposalsPage(common pageCommon) layout.Widget {
	pg := &ProposalsPage{
		theme:                       common.theme,
		wallet:                      win.wallet,
		proposalListContainer:       &layout.List{Axis: layout.Vertical},
		pageListContainer:           &layout.List{Axis: layout.Vertical},
		tabContainer:                decredmaterial.NewTabs(common.theme),
		tabTitles:                   []string{"In Discussion", "Voting", "Approved", "Rejected", "Abandoned"},
		isShowingProposalDetails:    false,
		proposals:                   &win.proposals,
		outline:                     common.theme.Outline(),
		selectedProposal:            &win.proposal,
		syncListener:                &proposalSyncListener{},
		notSyncedStatusLabel:        common.theme.H6("Not Synced"),
		fetchingTokenInventoryLabel: common.theme.H6("Fetching token inventory..."),
		syncingLabel:                common.theme.H6("Syncing..."),
	}
	pg.tabContainer.Position = decredmaterial.Top
	pg.notSyncedIcon = common.icons.navigationCancel
	pg.notSyncedIcon.Color = common.theme.Color.Danger

	pg.syncingIcon = common.icons.syncingIcon

	pg.syncButton = common.theme.Button(new(widget.Clickable), "Connect")
	pg.syncButton.TextSize = values.TextSize10
	pg.syncButton.Background = color.RGBA{}
	pg.syncButton.Color = common.theme.Color.Text

	pg.cancelSyncButton = common.theme.Button(new(widget.Clickable), "Cancel")
	pg.cancelSyncButton.TextSize = values.TextSize10
	pg.cancelSyncButton.Background = color.RGBA{}
	pg.cancelSyncButton.Color = common.theme.Color.Text

	tabItems := make([]decredmaterial.TabItem, len(pg.tabTitles))
	for i := range pg.tabTitles {
		tabItems[i] = decredmaterial.NewTabItem(pg.tabTitles[i], nil)
	}
	pg.tabContainer.SetTabs(tabItems)

	return func(gtx C) D {
		pg.Handler(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *ProposalsPage) SetClickables() {
	proposals := pg.getProposalsForCurrentTab()
	pg.clickables = make([]*gesture.Click, len(proposals))
	for i := range proposals {
		pg.clickables[i] = &gesture.Click{}
	}
}

func (pg *ProposalsPage) Handler(c pageCommon) {
	if pg.clickables == nil {
		pg.SetClickables()
	}

	for pg.tabContainer.ChangeEvent() {
		pg.SetClickables()
	}

	for pg.syncButton.Button.Clicked() {
		if !pg.syncListener.isSyncing {
			pg.wallet.StartProposalsSync(pg.syncListener)
		}
	}
}

func (pg *ProposalsPage) showProposalDetails(index int, c pageCommon) {
	category := pg.getSelectedProposalsCategory()
	proposals := *pg.proposals
	currentProposals := proposals[category]

	for i := range currentProposals {
		if currentProposals[i].Category == category && i == index {
			*pg.selectedProposal = &currentProposals[i]
			*c.page = PageProposalDetails
			break
		}
	}
}

func (pg *ProposalsPage) getSelectedProposalsCategory() int32 {
	switch pg.tabContainer.Selected {
	case 0:
		return dcrlibwallet.ProposalCategoryPre
	case 1:
		return dcrlibwallet.ProposalCategoryActive
	case 2:
		return dcrlibwallet.ProposalCategoryApproved
	case 3:
		return dcrlibwallet.ProposalCategoryRejected
	case 4:
		return dcrlibwallet.ProposalCategoryAbandoned
	default:
		return dcrlibwallet.ProposalCategoryAll
	}
}

func (pg *ProposalsPage) getProposalsForCurrentTab() []dcrlibwallet.Proposal {
	proposals := *pg.proposals
	return proposals[pg.getSelectedProposalsCategory()]
}

func (pg *ProposalsPage) Layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	for index, click := range pg.clickables {
		for _, e := range click.Events(gtx) {
			if e.Type == gesture.TypeClick {
				pg.showProposalDetails(index, c)
			}
		}
	}

	proposalListContainerHeight := gtx.Constraints.Max.Y - proposalSyncPaneHeight
	pageContent := []func(gtx C) D{
		func(gtx C) D {
			gtx.Constraints.Max.Y = proposalListContainerHeight
			return pg.layoutProposalsList(gtx)
		},
		func(gtx C) D {
			return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx C) D {
				return decredmaterial.Card{Color: pg.theme.Color.Surface}.Layout(gtx, func(gtx C) D {
					if pg.syncListener.isSyncing {
						return pg.layoutIsSyncingStatus(gtx)
					}
					return pg.layoutSyncStartSection(gtx)
				})
			})
		},
	}

	return c.Layout(gtx, func(gtx C) D {
		return pg.pageListContainer.Layout(gtx, len(pageContent), func(gtx C, i int) D {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, pageContent[i])
		})
	})
}

func (pg *ProposalsPage) layoutSyncStartSection(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: unit.Dp(10)}.Layout(gtx, func(gtx C) D {
					return pg.notSyncedIcon.Layout(gtx, unit.Dp(20))
				})
			}),
			layout.Rigid(func(gtx C) D {
				return pg.notSyncedStatusLabel.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx C) D {
					border := widget.Border{Color: pg.theme.Color.Hint, CornerRadius: values.MarginPadding5, Width: values.MarginPadding1}
					return border.Layout(gtx, func(gtx C) D {
						return pg.syncButton.Layout(gtx)
					})
				})
			}),
		)
	})
}

func (pg *ProposalsPage) layoutIsSyncingStatus(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.layoutSyncStatus(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.progressBarRow(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.layoutSyncProgressCounter(gtx)
			}),
		)
	})
}

func (pg *ProposalsPage) layoutSyncStatus(gtx layout.Context) layout.Dimensions {
	syncStatusLabel := pg.notSyncedStatusLabel
	btn := pg.syncButton
	if pg.syncListener.isSyncing {
		syncStatusLabel = pg.syncingLabel
		btn = pg.cancelSyncButton
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Right: unit.Dp(10)}.Layout(gtx, func(gtx C) D {
				if pg.syncListener.isSyncing {
					return pg.theme.ImageIcon(gtx, pg.syncingIcon, 20)
				}

				return pg.notSyncedIcon.Layout(gtx, unit.Dp(20))
			})
		}),
		layout.Rigid(func(gtx C) D {
			return syncStatusLabel.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				border := widget.Border{Color: pg.theme.Color.Hint, CornerRadius: values.MarginPadding5, Width: values.MarginPadding1}
				return border.Layout(gtx, func(gtx C) D {
					if pg.syncListener.isSyncing {
						return btn.Layout(gtx)
					}
					return btn.Layout(gtx)
				})
			})
		}),
	)
}

// syncBoxTitleRow lays out the progress bar.
func (pg *ProposalsPage) progressBarRow(gtx layout.Context) layout.Dimensions {
	percentageProgress := 0
	if pg.syncListener.progress != nil {
		percentageProgress = int(pg.syncListener.progress.ProposalsFetchProgress)
	}

	p := pg.theme.ProgressBar(int(percentageProgress))
	p.Color = pg.theme.Color.Success
	return p.Layout(gtx)
}

func (pg *ProposalsPage) layoutSyncProgressCounter(gtx layout.Context) layout.Dimensions {
	percentageProgress := 0
	if pg.syncListener.progress != nil {
		percentageProgress = int(pg.syncListener.progress.ProposalsFetchProgress)
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.theme.Body1(fmt.Sprintf("%d%%", percentageProgress)).Layout(gtx)
		}),
	)
}

func (pg *ProposalsPage) layoutProposalsList(gtx layout.Context) layout.Dimensions {
	return pg.tabContainer.Layout(gtx, func(gtx C) D {
		proposals := pg.getProposalsForCurrentTab()

		return pg.proposalListContainer.Layout(gtx, len(proposals), func(gtx C, i int) D {
			if len(pg.clickables) > 0 && len(pg.clickables) == len(proposals) {
				click := pg.clickables[i]
				pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
				click.Add(gtx.Ops)
			}

			return layout.UniformInset(unit.Dp(3)).Layout(gtx, func(gtx C) D {
				return decredmaterial.Card{Color: pg.theme.Color.Surface}.Layout(gtx, func(gtx C) D {
					return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
						return pg.layoutProposalHeader(gtx, proposals[i])
					})
				})
			})
		})
	})
}

func (pg *ProposalsPage) layoutProposalHeader(gtx layout.Context, proposal dcrlibwallet.Proposal) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(0.55, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return getTitleLabel(pg.theme, truncateString(proposal.Name, 60)).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return getSubtitleLabel(pg.theme, truncateString(proposal.CensorshipRecord.Token, 35)).Layout(gtx)
				}),
			)
		}),
		layout.Flexed(0.45, func(gtx C) D {
			if proposal.Category == dcrlibwallet.ProposalCategoryPre || proposal.Category == dcrlibwallet.ProposalCategoryAbandoned {
				return layout.E.Layout(gtx, func(gtx C) D {
					return getSubtitleLabel(pg.theme, fmt.Sprintf("Last updated %s", timeAgo(proposal.Timestamp))).Layout(gtx)
				})
			} else {
				yes, no := calculateVotes(proposal.VoteSummary.OptionsResult)
				return pg.theme.VoteBar(yes, no).Layout(gtx)
			}
		}),
	)
}

type proposalSyncListener struct {
	syncStage int
	isSyncing bool
	progress  *dcrlibwallet.ProposalsFetchProgressReport
}

func (p *proposalSyncListener) OnSyncStarted() {
	p.isSyncing = true
}

func (p *proposalSyncListener) OnProposalsDiscovery() {
	p.syncStage = dcrlibwallet.ProposalsDiscoverySyncState
}

func (p *proposalSyncListener) OnProposalsFetched(proposalsFetchProgress *dcrlibwallet.ProposalsFetchProgressReport) {
	if p.syncStage == dcrlibwallet.ProposalsDiscoverySyncState {
		p.syncStage = dcrlibwallet.ProposalsFetchSyncState
	}
	p.progress = proposalsFetchProgress
}

func (p *proposalSyncListener) OnSyncCompleted() {}

func (p *proposalSyncListener) OnSyncCanceled() {}

func (p *proposalSyncListener) OnSyncEndedWithError(err error) {}

func calculateVotes(options []dcrlibwallet.ProposalVoteOptionResult) (float32, float32) {
	var yes, no float32

	for i := range options {
		if options[i].Option.ID == "yes" {
			yes = float32(options[i].VotesReceived)
		} else {
			no = float32(options[i].VotesReceived)
		}
	}

	return yes, no
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

func timeAgo(timestamp int64) string {
	timeAgo, _ := timeago.TimeAgoWithTime(time.Now(), time.Unix(timestamp, 0))
	return timeAgo
}

func getTitleLabel(theme *decredmaterial.Theme, txt string) decredmaterial.Label {
	lbl := theme.Body1(txt)
	lbl.Color = theme.Color.Text
	return lbl
}

func getSubtitleLabel(theme *decredmaterial.Theme, txt string) decredmaterial.Label {
	lbl := theme.Caption(txt)
	lbl.Color = theme.Color.Hint
	return lbl
}

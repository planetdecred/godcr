package ui

import (
	"fmt"
	"image"
	"image/color"
	"time"

	//"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/ararog/timeago"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const (
	PageProposals = "Proposals"
)

type proposalNotificationListeners struct {
	wallet *wallet.Wallet
}

func (p proposalNotificationListeners) OnNewProposal(proposalID int, censorshipToken string) {
	p.wallet.GetProposalUpdate(censorshipToken, 1)
}

func (p proposalNotificationListeners) OnProposalVoteStarted(proposalID int, censorshipToken string) {
	p.wallet.GetProposalUpdate(censorshipToken, 2)
}

func (p proposalNotificationListeners) OnProposalVoteFinished(proposalID int, censorshipToken string) {
	p.wallet.GetProposalUpdate(censorshipToken, 3)
}

type proposalItem struct {
	proposal dcrlibwallet.Proposal
	button   *widget.Clickable
}

type ProposalsPage struct {
	theme                          *decredmaterial.Theme
	wallet                         *wallet.Wallet
	pageListContainer              *layout.List
	proposalListContainer          *layout.List
	tabTitles                      []string
	tabContainer                   *decredmaterial.Tabs
	outline                        decredmaterial.Outline
	isSyncing                      bool
	hasFetchedSavedProposals       bool
	hasRegisteredProposalListeners bool
	notSyncingIcon                 *widget.Icon
	syncingIcon                    image.Image
	notSyncingStatusLabel          decredmaterial.Label
	syncButton                     decredmaterial.Button
	cancelSyncButton               decredmaterial.Button
	syncingLabel                   decredmaterial.Label
	proposals                      map[int32][]proposalItem
	latestProposals                *[]dcrlibwallet.Proposal
	updatedProposal                **wallet.UpdatedProposal
	selectedProposal               **dcrlibwallet.Proposal
}

func (win *Window) ProposalsPage(common pageCommon) layout.Widget {
	pg := &ProposalsPage{
		theme:                 common.theme,
		wallet:                win.wallet,
		proposalListContainer: &layout.List{Axis: layout.Vertical},
		pageListContainer:     &layout.List{Axis: layout.Vertical},
		tabContainer:          decredmaterial.NewTabs(common.theme),
		tabTitles:             []string{"In Discussion", "Voting", "Approved", "Rejected", "Abandoned"},
		proposals:             make(map[int32][]proposalItem),
		outline:               common.theme.Outline(),
		isSyncing:             false,
		notSyncingStatusLabel: common.theme.H6("Not Syncing"),
		syncingLabel:          common.theme.H6("Syncing..."),
		latestProposals:       &win.latestProposals,
		updatedProposal:       &win.updatedProposal,
		selectedProposal:      &win.selectedProposal,
	}

	pg.tabContainer.Position = decredmaterial.Top
	pg.notSyncingIcon = common.icons.navigationCancel
	pg.notSyncingIcon.Color = common.theme.Color.Danger

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

func (pg *ProposalsPage) Handler(c pageCommon) {
	for proposalGroupIndex := range pg.proposals {
		pgIndex := proposalGroupIndex
		for proposalItemIndex := range pg.proposals[pgIndex] {
			piIndex := proposalItemIndex
			for pg.proposals[pgIndex][piIndex].button.Clicked() {
				*pg.selectedProposal = &pg.proposals[pgIndex][piIndex].proposal //&proposalItem.proposal
				*c.page = PageProposalDetails
			}
		}
	}

	for pg.syncButton.Button.Clicked() {
		if !pg.isSyncing {
			pg.isSyncing = true
			pg.wallet.StartProposalsSync()
		}
	}

	for pg.cancelSyncButton.Button.Clicked() {
		if pg.isSyncing {
			pg.wallet.CancelProposalsSync()
			pg.isSyncing = false
		}
	}

	if *pg.latestProposals != nil {
		pg.addLatestProposals()
	}

	if *pg.updatedProposal != nil {
		pg.addUpdatedProposal()
	}
}

func (pg *ProposalsPage) addProposal(proposal dcrlibwallet.Proposal) {
	proposalItem := proposalItem{
		proposal: proposal,
		button:   new(widget.Clickable),
	}
	pg.proposals[proposal.Category] = append(pg.proposals[proposal.Category], proposalItem)
}

func (pg *ProposalsPage) addLatestProposals() {
	latestProposals := *pg.latestProposals
	for _, v := range latestProposals {
		pg.addProposal(v)
	}
	*pg.latestProposals = nil
}

func (pg *ProposalsPage) addUpdatedProposal() {
	updatedProposal := *pg.updatedProposal

	if updatedProposal.UpdateType != 1 {
		for proposalGroupIndex, proposalGroup := range pg.proposals {
			for proposalItemIndex, proposalItem := range proposalGroup {
				if proposalItem.proposal.CensorshipRecord.Token == updatedProposal.Proposal.CensorshipRecord.Token {
					pg.proposals[proposalGroupIndex] = append(pg.proposals[proposalGroupIndex][:proposalItemIndex], pg.proposals[proposalGroupIndex][proposalItemIndex+1:]...)
				}
			}
		}
	}
	pg.addProposal(*updatedProposal.Proposal)
	*pg.updatedProposal = nil
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

func (pg *ProposalsPage) getProposalsForCurrentTab() []proposalItem {
	return pg.proposals[pg.getSelectedProposalsCategory()]
}

func (pg *ProposalsPage) Layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	if !pg.hasFetchedSavedProposals {
		pg.wallet.GetProposals()
		pg.hasFetchedSavedProposals = true
	}

	if !pg.hasRegisteredProposalListeners {
		pg.wallet.AddProposalNotificationListener(proposalNotificationListeners{pg.wallet})
		pg.hasRegisteredProposalListeners = true
	}

	return c.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return pg.layoutProposalsList(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx C) D {
					return decredmaterial.Card{Color: pg.theme.Color.Surface}.Layout(gtx, func(gtx C) D {
						if pg.isSyncing {
							return pg.layoutIsSyncingSection(gtx)
						}
						return pg.layoutSyncStartSection(gtx)
					})
				})
			}),
		)
	})
}

func (pg *ProposalsPage) layoutSyncStartSection(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: unit.Dp(10)}.Layout(gtx, func(gtx C) D {
					return pg.notSyncingIcon.Layout(gtx, unit.Dp(20))
				})
			}),
			layout.Rigid(func(gtx C) D {
				return pg.notSyncingStatusLabel.Layout(gtx)
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

func (pg *ProposalsPage) layoutIsSyncingSection(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Right: unit.Dp(10)}.Layout(gtx, func(gtx C) D {
								return pg.theme.ImageIcon(gtx, pg.syncingIcon, 20)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return pg.syncingLabel.Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							return layout.E.Layout(gtx, func(gtx C) D {
								border := widget.Border{Color: pg.theme.Color.Hint, CornerRadius: values.MarginPadding5, Width: values.MarginPadding1}
								return border.Layout(gtx, func(gtx C) D {
									return pg.cancelSyncButton.Layout(gtx)
								})
							})
						}),
					)
				})
			}),
		)
	})
}

func (pg *ProposalsPage) layoutProposalsList(gtx layout.Context) layout.Dimensions {
	return pg.tabContainer.Layout(gtx, func(gtx C) D {
		proposals := pg.getProposalsForCurrentTab()

		return pg.proposalListContainer.Layout(gtx, len(proposals), func(gtx C, i int) D {
			return layout.UniformInset(unit.Dp(3)).Layout(gtx, func(gtx C) D {
				return decredmaterial.Card{Color: pg.theme.Color.Surface}.Layout(gtx, func(gtx C) D {
					return material.Clickable(gtx, proposals[i].button, func(gtx C) D {
						return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
							return pg.layoutProposalHeader(gtx, proposals[i])
						})
					})
				})
			})
		})
	})
}

func (pg *ProposalsPage) layoutProposalHeader(gtx layout.Context, proposalItem proposalItem) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(0.55, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return getTitleLabel(pg.theme, truncateString(proposalItem.proposal.Name, 60)).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return getSubtitleLabel(pg.theme, truncateString(proposalItem.proposal.CensorshipRecord.Token, 35)).Layout(gtx)
				}),
			)
		}),
		layout.Flexed(0.45, func(gtx C) D {
			if proposalItem.proposal.Category == dcrlibwallet.ProposalCategoryPre || proposalItem.proposal.Category == dcrlibwallet.ProposalCategoryAbandoned {
				return layout.E.Layout(gtx, func(gtx C) D {
					return getSubtitleLabel(pg.theme, fmt.Sprintf("Last updated %s", timeAgo(proposalItem.proposal.Timestamp))).Layout(gtx)
				})
			}
			yes, no := calculateVotes(proposalItem.proposal.VoteSummary.OptionsResult)
			return pg.theme.VoteBar(yes, no).Layout(gtx)
		}),
	)
}

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

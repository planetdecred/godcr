package proposal

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"strings"
	"sync"
	"time"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

const ProposalsPageID = "Proposals"

type proposalItem struct {
	proposal     dcrlibwallet.Proposal
	tooltip      *decredmaterial.Tooltip
	tooltipLabel decredmaterial.Label
	voteBar      *VoteBar
}

type ProposalsPage struct {
	*load.Load

	ctx        context.Context // page context
	ctxCancel  context.CancelFunc
	proposalMu sync.Mutex

	multiWallet *dcrlibwallet.MultiWallet

	//categoryList to be removed with new update to UI.
	categoryList  *decredmaterial.ClickableList
	proposalsList *decredmaterial.ClickableList

	tabCard        decredmaterial.Card
	itemCard       decredmaterial.Card
	syncCard       decredmaterial.Card
	updatedLabel   decredmaterial.Label
	lastSyncedTime string

	proposalItems         []proposalItem
	proposalCount         []int
	selectedCategoryIndex int

	legendIcon    *decredmaterial.Icon
	infoIcon      *decredmaterial.Icon
	updatedIcon   *decredmaterial.Icon
	syncButton    *widget.Clickable
	startSyncIcon *decredmaterial.Image
	timerIcon     *decredmaterial.Image

	showSyncedCompleted bool
	isSyncing           bool
}

var (
	proposalCategoryTitles = []string{"In discussion", "Voting", "Approved", "Rejected", "Abandoned"}
	proposalCategories     = []int32{
		dcrlibwallet.ProposalCategoryPre,
		dcrlibwallet.ProposalCategoryActive,
		dcrlibwallet.ProposalCategoryApproved,
		dcrlibwallet.ProposalCategoryRejected,
		dcrlibwallet.ProposalCategoryAbandoned,
	}
)

func NewProposalsPage(l *load.Load) *ProposalsPage {
	pg := &ProposalsPage{
		Load:                  l,
		multiWallet:           l.WL.MultiWallet,
		selectedCategoryIndex: -1,
	}

	pg.initLayoutWidgets()

	return pg
}

func (pg *ProposalsPage) ID() string {
	return ProposalsPageID
}

func (pg *ProposalsPage) OnResume() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())

	pg.listenForSyncNotifications()

	pg.proposalMu.Lock()
	selectedCategory := pg.selectedCategoryIndex
	pg.proposalMu.Unlock()
	if selectedCategory == -1 {
		pg.countProposals()
		pg.loadProposals(0)
	}

	pg.isSyncing = pg.multiWallet.Politeia.IsSyncing()
}

func (pg *ProposalsPage) countProposals() {
	proposalCount := make([]int, len(proposalCategories))
	for i, category := range proposalCategories {
		count, err := pg.multiWallet.Politeia.Count(category)
		if err == nil {
			proposalCount[i] = int(count)
		}
	}

	pg.proposalMu.Lock()
	pg.proposalCount = proposalCount
	pg.proposalMu.Unlock()
}

func (pg *ProposalsPage) loadProposals(category int) {
	proposals, err := pg.multiWallet.Politeia.GetProposalsRaw(proposalCategories[category], 0, 0, true)
	if err != nil {
		pg.proposalMu.Lock()
		pg.proposalItems = make([]proposalItem, 0)
		pg.proposalMu.Unlock()
	} else {
		proposalItems := make([]proposalItem, len(proposals))
		for i := 0; i < len(proposals); i++ {
			proposal := proposals[i]
			item := proposalItem{
				proposal: proposals[i],
				voteBar:  NewVoteBar(pg.Load),
			}

			if proposal.Category == dcrlibwallet.ProposalCategoryPre {
				tooltipLabel := pg.Theme.Caption("")
				tooltipLabel.Color = pg.Theme.Color.Gray
				if proposal.VoteStatus == 1 {
					tooltipLabel.Text = "Waiting for author to authorize voting"
				} else if proposal.VoteStatus == 2 {
					tooltipLabel.Text = "Waiting for admin to trigger the start of voting"
				}

				item.tooltip = pg.Theme.Tooltip()
				item.tooltipLabel = tooltipLabel
			}

			proposalItems[i] = item
		}
		pg.proposalMu.Lock()
		pg.selectedCategoryIndex = category
		pg.proposalItems = proposalItems
		pg.proposalMu.Unlock()
	}
}

func (pg *ProposalsPage) Handle() {
	//categoryList to be removed with new update to UI.
	if clicked, selectedItem := pg.categoryList.ItemClicked(); clicked {
		go pg.loadProposals(selectedItem)
	}

	if clicked, selectedItem := pg.proposalsList.ItemClicked(); clicked {
		pg.proposalMu.Lock()
		selectedProposal := pg.proposalItems[selectedItem].proposal
		pg.proposalMu.Unlock()

		pg.ChangeFragment(newProposalDetailsPage(pg.Load, &selectedProposal))
	}

	for pg.syncButton.Clicked() {
		pg.isSyncing = true
		go pg.multiWallet.Politeia.Sync()
	}

	if pg.showSyncedCompleted {
		time.AfterFunc(time.Second*3, func() {
			pg.showSyncedCompleted = false
		})
	}
}

func (pg *ProposalsPage) listenForSyncNotifications() {
	go func() {
		for {
			var notification interface{}

			select {
			case notification = <-pg.Receiver.NotificationsUpdate:
			case <-pg.ctx.Done():
				return
			}

			switch n := notification.(type) {
			case wallet.Proposal:
				if n.ProposalStatus == wallet.Synced {
					pg.isSyncing = false
					pg.showSyncedCompleted = true

					pg.proposalMu.Lock()
					selectedCategory := pg.selectedCategoryIndex
					pg.proposalMu.Unlock()
					if selectedCategory != -1 {
						pg.countProposals()
						pg.loadProposals(selectedCategory)
					}
				}
			}
		}
	}()
}

func (pg *ProposalsPage) OnClose() {
	pg.ctxCancel()
}

// - Layout

func (pg *ProposalsPage) initLayoutWidgets() {
	//categoryList to be removed with new update to UI.
	pg.categoryList = pg.Theme.NewClickableList(layout.Horizontal)
	pg.itemCard = pg.Theme.Card()
	pg.syncButton = new(widget.Clickable)

	pg.infoIcon = decredmaterial.NewIcon(pg.Icons.ActionInfo)
	pg.infoIcon.Color = pg.Theme.Color.Gray

	pg.legendIcon = decredmaterial.NewIcon(pg.Icons.ImageBrightness1)
	pg.legendIcon.Color = pg.Theme.Color.InactiveGray

	pg.updatedIcon = decredmaterial.NewIcon(pg.Icons.NavigationCheck)
	pg.updatedIcon.Color = pg.Theme.Color.Success

	pg.updatedLabel = pg.Theme.Body2("Updated")
	pg.updatedLabel.Color = pg.Theme.Color.Success

	radius := decredmaterial.Radius(0)
	pg.tabCard = pg.Theme.Card()
	pg.tabCard.Radius = radius

	pg.syncCard = pg.Theme.Card()
	pg.syncCard.Radius = radius

	pg.proposalsList = pg.Theme.NewClickableList(layout.Vertical)
	pg.proposalsList.DividerHeight = values.MarginPadding8

	pg.timerIcon = pg.Icons.TimerIcon

	pg.startSyncIcon = pg.Icons.Restore
}

func (pg *ProposalsPage) Layout(gtx C) D {
	border := widget.Border{Color: pg.Theme.Color.Gray1, CornerRadius: values.MarginPadding0, Width: values.MarginPadding1}
	borderLayout := func(gtx layout.Context, body layout.Widget) layout.Dimensions {
		return border.Layout(gtx, body)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return borderLayout(gtx, pg.layoutTabs)
				}),
				layout.Rigid(func(gtx C) D {
					return borderLayout(gtx, func(gtx C) D {
						return pg.syncCard.Layout(gtx, func(gtx C) D {
							m := values.MarginPadding12
							if pg.showSyncedCompleted || pg.isSyncing {
								m = values.MarginPadding14
							}
							return layout.UniformInset(m).Layout(gtx, func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
											return layout.Center.Layout(gtx, pg.layoutSyncSection)
										})
									}),
									layout.Rigid(func(gtx C) D {
										if pg.showSyncedCompleted || pg.isSyncing {
											return D{}
										}
										lastUpdatedInfo := pg.Theme.Body2(components.TimeAgo(pg.multiWallet.Politeia.GetLastSyncedTimeStamp()))
										lastUpdatedInfo.Color = pg.Theme.Color.Text
										return layout.Inset{Top: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
											return lastUpdatedInfo.Layout(gtx)
										})
									}),
								)
							})
						})
					})
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding24}.Layout(gtx, pg.layoutContent)
		}),
	)
}

func (pg *ProposalsPage) layoutContent(gtx C) D {
	pg.proposalMu.Lock()
	proposalItems := pg.proposalItems
	pg.proposalMu.Unlock()
	if len(proposalItems) == 0 {
		return pg.layoutNoProposalsFound(gtx)
	}
	return pg.layoutProposalsList(gtx)
}

func (pg *ProposalsPage) layoutProposalsList(gtx C) D {
	pg.proposalMu.Lock()
	proposalItems := pg.proposalItems
	pg.proposalMu.Unlock()
	return pg.proposalsList.Layout(gtx, len(proposalItems), func(gtx C, i int) D {
		return components.UniformHorizontalPadding(gtx, func(gtx C) D {
			return layout.Inset{
				Top:    values.MarginPadding2,
				Bottom: values.MarginPadding2,
			}.Layout(gtx, func(gtx C) D {
				return pg.itemCard.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
						item := proposalItems[i]
						proposal := item.proposal
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return pg.layoutAuthorAndDate(gtx, item)
							}),
							layout.Rigid(func(gtx C) D {
								return pg.layoutTitle(gtx, proposal)
							}),
							layout.Rigid(func(gtx C) D {
								if proposal.Category == dcrlibwallet.ProposalCategoryActive ||
									proposal.Category == dcrlibwallet.ProposalCategoryApproved ||
									proposal.Category == dcrlibwallet.ProposalCategoryRejected {
									return pg.layoutProposalVoteBar(gtx, item)
								}
								return D{}
							}),
						)
					})
				})
			})
		})
	})
}

func (pg *ProposalsPage) layoutNoProposalsFound(gtx C) D {
	pg.proposalMu.Lock()
	selectedCategory := pg.selectedCategoryIndex
	pg.proposalMu.Unlock()
	str := fmt.Sprintf("No %s proposals", strings.ToLower(proposalCategoryTitles[selectedCategory]))

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Center.Layout(gtx, pg.Theme.Body1(str).Layout)
}

func (pg *ProposalsPage) layoutSyncSection(gtx C) D {
	if pg.showSyncedCompleted {
		return pg.layoutIsSyncedSection(gtx)
	} else if pg.multiWallet.Politeia.IsSyncing() {
		return pg.layoutIsSyncingSection(gtx)
	}
	return pg.layoutStartSyncSection(gtx)
}

func (pg *ProposalsPage) layoutTabs(gtx C) D {
	width := float32(gtx.Constraints.Max.X-20) / float32(len(proposalCategoryTitles))
	pg.proposalMu.Lock()
	selectedCategory := pg.selectedCategoryIndex
	pg.proposalMu.Unlock()

	return pg.tabCard.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Left:  values.MarginPadding12,
			Right: values.MarginPadding12,
		}.Layout(gtx, func(gtx C) D {
			// categoryList to be removed with new update to UI.
			return pg.categoryList.Layout(gtx, len(proposalCategoryTitles), func(gtx C, i int) D {
				gtx.Constraints.Min.X = int(width)
				return layout.Stack{Alignment: layout.S}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						return layout.UniformInset(values.MarginPadding14).Layout(gtx, func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										lbl := pg.Theme.Body1(proposalCategoryTitles[i])
										lbl.Color = pg.Theme.Color.Gray
										if selectedCategory == i {
											lbl.Color = pg.Theme.Color.Primary
										}
										return lbl.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Left: values.MarginPadding4, Top: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
											c := pg.Theme.Card()
											c.Color = pg.Theme.Color.LightGray
											c.Radius = decredmaterial.Radius(8.5)
											lbl := pg.Theme.Body2(strconv.Itoa(pg.proposalCount[i]))
											lbl.Color = pg.Theme.Color.Gray
											if selectedCategory == i {
												c.Color = pg.Theme.Color.Primary
												lbl.Color = pg.Theme.Color.Surface
											}
											return c.Layout(gtx, func(gtx C) D {
												return layout.Inset{
													Left:  values.MarginPadding5,
													Right: values.MarginPadding5,
												}.Layout(gtx, lbl.Layout)
											})
										})
									}),
								)
							})
						})
					}),
					layout.Stacked(func(gtx C) D {
						if selectedCategory != i {
							return D{}
						}
						tabHeight := gtx.Px(unit.Dp(2))
						tabRect := image.Rect(0, 0, int(width), tabHeight)
						paint.FillShape(gtx.Ops, pg.Theme.Color.Primary, clip.Rect(tabRect).Op())
						return layout.Dimensions{
							Size: image.Point{X: int(width), Y: tabHeight},
						}
					}),
				)
			})
		})
	})
}

func (pg *ProposalsPage) layoutAuthorAndDate(gtx C, item proposalItem) D {
	proposal := item.proposal
	grayCol := pg.Theme.Color.Gray

	nameLabel := pg.Theme.Body2(proposal.Username)
	nameLabel.Color = grayCol

	dotLabel := pg.Theme.H4(" . ")
	dotLabel.Color = grayCol

	versionLabel := pg.Theme.Body2("Version " + proposal.Version)
	versionLabel.Color = grayCol

	stateLabel := pg.Theme.Body2(fmt.Sprintf("%v /2", proposal.VoteStatus))
	stateLabel.Color = grayCol

	timeAgoLabel := pg.Theme.Body2(components.TimeAgo(proposal.Timestamp))
	timeAgoLabel.Color = grayCol

	var categoryLabel decredmaterial.Label
	var categoryLabelColor color.NRGBA
	switch proposal.Category {
	case dcrlibwallet.ProposalCategoryApproved:
		categoryLabel = pg.Theme.Body2("Approved")
		categoryLabelColor = pg.Theme.Color.Success
	case dcrlibwallet.ProposalCategoryActive:
		categoryLabel = pg.Theme.Body2("Voting")
		categoryLabelColor = pg.Theme.Color.Primary
	case dcrlibwallet.ProposalCategoryRejected:
		categoryLabel = pg.Theme.Body2("Rejected")
		categoryLabelColor = pg.Theme.Color.Danger
	case dcrlibwallet.ProposalCategoryAbandoned:
		categoryLabel = pg.Theme.Body2("Abandoned")
		categoryLabelColor = grayCol
	case dcrlibwallet.ProposalCategoryPre:
		categoryLabel = pg.Theme.Body2("In discussion")
		categoryLabelColor = grayCol
	}
	categoryLabel.Color = categoryLabelColor

	pg.proposalMu.Lock()
	selectedCategory := pg.selectedCategoryIndex
	pg.proposalMu.Unlock()

	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(nameLabel.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPaddingMinus22}.Layout(gtx, dotLabel.Layout)
				}),
				layout.Rigid(versionLabel.Layout),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(categoryLabel.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPaddingMinus22}.Layout(gtx, dotLabel.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					if proposalCategories[selectedCategory] == dcrlibwallet.ProposalCategoryPre {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(stateLabel.Layout),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									rect := image.Rectangle{
										Min: gtx.Constraints.Min,
										Max: gtx.Constraints.Max,
									}
									rect.Max.Y = 20
									pg.layoutInfoTooltip(gtx, rect, item)
									return pg.infoIcon.Layout(gtx, values.MarginPadding20)
								})
							}),
						)
					}

					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							if proposalCategories[selectedCategory] == dcrlibwallet.ProposalCategoryActive {
								return layout.Inset{
									Right: values.MarginPadding4,
									Top:   values.MarginPadding3,
								}.Layout(gtx, pg.timerIcon.Layout12dp)
							}
							return D{}
						}),
						layout.Rigid(timeAgoLabel.Layout),
					)
				}),
			)
		}),
	)
}

func (pg *ProposalsPage) layoutInfoTooltip(gtx C, rect image.Rectangle, item proposalItem) {
	inset := layout.Inset{Top: values.MarginPadding20, Left: values.MarginPaddingMinus195}
	item.tooltip.Layout(gtx, rect, inset, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Px(values.MarginPadding195)
		gtx.Constraints.Max.X = gtx.Px(values.MarginPadding195)
		return item.tooltipLabel.Layout(gtx)
	})
}

func (pg *ProposalsPage) layoutTitle(gtx C, proposal dcrlibwallet.Proposal) D {
	lbl := pg.Theme.H6(proposal.Name)
	lbl.Font.Weight = text.Bold
	return layout.Inset{Top: values.MarginPadding4}.Layout(gtx, lbl.Layout)
}

func (pg *ProposalsPage) layoutProposalVoteBar(gtx C, item proposalItem) D {
	proposal := item.proposal
	yes := float32(proposal.YesVotes)
	no := float32(proposal.NoVotes)
	quorumPercent := float32(proposal.QuorumPercentage)
	passPercentage := float32(proposal.PassPercentage)
	eligibleTickets := float32(proposal.EligibleTickets)

	return item.voteBar.
		SetYesNoVoteParams(yes, no).
		SetVoteValidityParams(eligibleTickets, quorumPercent, passPercentage).
		SetProposalDetails(proposal.NumComments, proposal.PublishedAt, proposal.Token).
		Layout(gtx)
}

func (pg *ProposalsPage) layoutIsSyncedSection(gtx C) D {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.updatedIcon.Layout(gtx, values.MarginPadding20)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, pg.updatedLabel.Layout)
		}),
	)
}

func (pg *ProposalsPage) layoutIsSyncingSection(gtx C) D {
	th := material.NewTheme(gofont.Collection())
	gtx.Constraints.Min.X = gtx.Px(unit.Dp(20))
	loader := material.Loader(th)
	loader.Color = pg.Theme.Color.Gray
	return loader.Layout(gtx)
}

func (pg *ProposalsPage) layoutStartSyncSection(gtx C) D {
	return material.Clickable(gtx, pg.syncButton, func(gtx C) D {
		return pg.startSyncIcon.Layout24dp(gtx)
	})
}

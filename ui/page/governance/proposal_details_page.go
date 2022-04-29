package governance

import (
	"context"
	"fmt"
	"time"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/renderers"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const ProposalDetailsPageID = "proposal_details"

type proposalItemWidgets struct {
	widgets    []layout.Widget
	clickables map[string]*widget.Clickable
}

type ProposalDetails struct {
	*load.Load
	*listeners.ProposalNotificationListener //not needed.

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	descriptionList *layout.List

	proposal      *dcrlibwallet.Proposal
	proposalItems map[string]proposalItemWidgets

	scrollbarList *widget.List
	rejectedIcon  *widget.Icon
	successIcon   *widget.Icon

	redirectIcon *decredmaterial.Image
	copyIcon     *decredmaterial.Image

	viewInPoliteiaBtn *decredmaterial.Clickable
	copyRedirectURL   *decredmaterial.Clickable

	descriptionCard decredmaterial.Card
	vote            decredmaterial.Button
	backButton      decredmaterial.IconButton

	voteBar            *components.VoteBar
	loadingDescription bool
}

func NewProposalDetailsPage(l *load.Load, proposal *dcrlibwallet.Proposal) *ProposalDetails {
	pg := &ProposalDetails{
		Load: l,

		loadingDescription: false,
		proposal:           proposal,
		descriptionCard:    l.Theme.Card(),
		descriptionList:    &layout.List{Axis: layout.Vertical},
		scrollbarList: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		redirectIcon:      l.Theme.Icons.RedirectIcon,
		proposalItems:     make(map[string]proposalItemWidgets),
		rejectedIcon:      l.Theme.Icons.NavigationCancel,
		successIcon:       l.Theme.Icons.ActionCheckCircle,
		viewInPoliteiaBtn: l.Theme.NewClickable(true),
		copyRedirectURL:   l.Theme.NewClickable(false),
		voteBar:           components.NewVoteBar(l.Theme),
	}

	pg.backButton, _ = components.SubpageHeaderButtons(l.Theme)

	pg.vote = l.Theme.Button("Vote")
	pg.vote.TextSize = values.TextSize14
	pg.vote.Background = l.Theme.Color.Primary
	pg.vote.Color = l.Theme.Color.Surface
	pg.vote.Inset = layout.Inset{
		Top:    values.MarginPadding8,
		Bottom: values.MarginPadding8,
		Left:   values.MarginPadding12,
		Right:  values.MarginPadding12,
	}

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *ProposalDetails) ID() string {
	return ProposalDetailsPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *ProposalDetails) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	pg.listenForSyncNotifications()
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *ProposalDetails) HandleUserInteractions() {
	for token := range pg.proposalItems {
		for location, clickable := range pg.proposalItems[token].clickables {
			if clickable.Clicked() {
				components.GoToURL(location)
			}
		}
	}

	if pg.vote.Clicked() {
		newVoteModal(pg.Load, pg.proposal).Show()
	}

	for pg.viewInPoliteiaBtn.Clicked() {
		host := "https://proposals.decred.org/record/" + pg.proposal.Token
		if pg.WL.MultiWallet.NetType() == dcrlibwallet.Testnet3 {
			host = "https://test-proposals.decred.org/record/" + pg.proposal.Token
		}

		info := modal.NewInfoModal(pg.Load).
			Title("View on Politeia").
			Body("Copy and paste the link below in your browser, to view proposal on Politeia dashboard.").
			SetCancelable(true).
			UseCustomWidget(func(gtx C) D {
				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						border := widget.Border{Color: pg.Theme.Color.Gray4, CornerRadius: values.MarginPadding10, Width: values.MarginPadding2}
						wrapper := pg.Theme.Card()
						wrapper.Color = pg.Theme.Color.Gray4
						return border.Layout(gtx, func(gtx C) D {
							return wrapper.Layout(gtx, func(gtx C) D {
								return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
									return layout.Flex{}.Layout(gtx,
										layout.Flexed(0.9, pg.Theme.Body1(host).Layout),
										layout.Flexed(0.1, func(gtx C) D {
											return layout.E.Layout(gtx, func(gtx C) D {
												return layout.Inset{Top: values.MarginPadding7}.Layout(gtx, func(gtx C) D {
													if pg.copyRedirectURL.Clicked() {
														clipboard.WriteOp{Text: host}.Add(gtx.Ops)
														pg.Toast.Notify("URL copied")
													}
													return pg.copyRedirectURL.Layout(gtx, pg.Theme.Icons.CopyIcon.Layout24dp)
												})
											})
										}),
									)
								})
							})
						})
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:  values.MarginPaddingMinus10,
							Left: values.MarginPadding10,
						}.Layout(gtx, func(gtx C) D {
							label := pg.Theme.Body2("Web URL")
							label.Color = pg.Theme.Color.GrayText2
							return label.Layout(gtx)
						})
					}),
				)
			}).
			PositiveButton("Got it", func(isChecked bool) {})
		pg.ShowModal(info)
	}
}

func (pg *ProposalDetails) listenForSyncNotifications() {
	if pg.ProposalNotificationListener == nil {
		return
	}
	pg.ProposalNotificationListener = listeners.NewProposalNotificationListener()
	err := pg.WL.MultiWallet.Politeia.AddNotificationListener(pg.ProposalNotificationListener, ProposalDetailsPageID)
	if err != nil {
		log.Errorf("Error adding politeia notification listener: %v", err)
		return
	}

	go func() {
		for {
			select {
			case notification := <-pg.ProposalNotifChan:
				if notification.ProposalStatus == wallet.Synced {
					proposal, err := pg.WL.MultiWallet.Politeia.GetProposalRaw(pg.proposal.Token)
					if err == nil {
						pg.proposal = proposal
						pg.RefreshWindow()
					}
				}
			// is this really needed since listener has been set up on main.go
			case <-pg.ctx.Done():
				pg.WL.MultiWallet.Politeia.RemoveNotificationListener(ProposalDetailsPageID)
				close(pg.ProposalNotifChan)
				pg.ProposalNotificationListener = nil

				return
			}
		}
	}()
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *ProposalDetails) OnNavigatedFrom() {
	pg.ctxCancel()
}

// - Layout

func (pg *ProposalDetails) layoutProposalVoteBar(gtx C) D {
	proposal := pg.proposal

	yes := float32(proposal.YesVotes)
	no := float32(proposal.NoVotes)
	quorumPercent := float32(proposal.QuorumPercentage)
	passPercentage := float32(proposal.PassPercentage)
	eligibleTickets := float32(proposal.EligibleTickets)

	return pg.voteBar.
		SetYesNoVoteParams(yes, no).
		SetVoteValidityParams(eligibleTickets, quorumPercent, passPercentage).
		SetProposalDetails(proposal.NumComments, proposal.PublishedAt, proposal.Token).
		Layout(gtx)
}

func (pg *ProposalDetails) layoutProposalVoteAction(gtx C) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return pg.vote.Layout(gtx)
}

func (pg *ProposalDetails) layoutInDiscussionState(gtx C) D {
	stateText1 := "Waiting for author to authorize voting"
	stateText2 := "Waiting for admin to trigger the start of voting"

	proposal := pg.proposal

	c := func(gtx layout.Context, val int32, info string) layout.Dimensions {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				if proposal.VoteStatus == val || proposal.VoteStatus < val {
					c := pg.Theme.Card()
					c.Color = pg.Theme.Color.Primary
					c.Radius = decredmaterial.Radius(9.5)
					lbl := pg.Theme.Body1(fmt.Sprint(val))
					lbl.Color = pg.Theme.Color.Surface
					if proposal.VoteStatus < val {
						c.Color = pg.Theme.Color.Gray4
						lbl.Color = pg.Theme.Color.GrayText3
					}
					return c.Layout(gtx, func(gtx C) D {
						m := values.MarginPadding6
						return layout.Inset{Left: m, Right: m}.Layout(gtx, lbl.Layout)
					})
				}
				icon := decredmaterial.NewIcon(pg.successIcon)
				icon.Color = pg.Theme.Color.Primary
				return layout.Inset{
					Left:   values.MarginPaddingMinus2,
					Right:  values.MarginPaddingMinus2,
					Bottom: values.MarginPaddingMinus2,
				}.Layout(gtx, func(gtx C) D {
					return icon.Layout(gtx, values.MarginPadding24)
				})
			}),
			layout.Rigid(func(gtx C) D {
				col := pg.Theme.Color.Primary
				txt := info + "..."
				if proposal.VoteStatus != val {
					txt = info
					col = pg.Theme.Color.GrayText3
					if proposal.VoteStatus > 1 {
						col = pg.Theme.Color.Text
					}
				}
				lbl := pg.Theme.Body1(txt)
				lbl.Color = col
				return layout.Inset{Left: values.MarginPadding16}.Layout(gtx, lbl.Layout)
			}),
		)
	}

	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return c(gtx, 1, stateText1)
		}),
		layout.Rigid(func(gtx C) D {
			height, width := gtx.Px(values.MarginPadding26), gtx.Px(values.MarginPadding4)
			line := pg.Theme.Line(height, width)
			if proposal.VoteStatus > 1 {
				line.Color = pg.Theme.Color.Primary
			} else {
				line.Color = pg.Theme.Color.Gray2
			}
			return layout.Inset{Left: values.MarginPadding8}.Layout(gtx, line.Layout)
		}),
		layout.Rigid(func(gtx C) D {
			return c(gtx, 2, stateText2)
		}),
	)
}

func (pg *ProposalDetails) layoutNormalTitle(gtx C) D {
	var label decredmaterial.Label
	var icon *decredmaterial.Icon
	proposal := pg.proposal
	switch proposal.Category {
	case dcrlibwallet.ProposalCategoryApproved:
		label = pg.Theme.Body2("Approved")
		icon = decredmaterial.NewIcon(pg.successIcon)
		icon.Color = pg.Theme.Color.Success
	case dcrlibwallet.ProposalCategoryRejected:
		label = pg.Theme.Body2("Rejected")
		icon = decredmaterial.NewIcon(pg.rejectedIcon)
		icon.Color = pg.Theme.Color.Danger
	case dcrlibwallet.ProposalCategoryAbandoned:
		label = pg.Theme.Body2("Abandoned")
	case dcrlibwallet.ProposalCategoryActive:
		label = pg.Theme.Body2("Voting in progress...")
	}
	timeagoLabel := pg.Theme.Body2(components.TimeAgo(proposal.Timestamp))

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if icon == nil {
						return D{}
					}
					return icon.Layout(gtx, values.MarginPadding20)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, label.Layout)
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								if proposal.Category == dcrlibwallet.ProposalCategoryActive {
									ic := pg.Theme.Icons.TimerIcon
									if pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.DarkModeConfigKey, false) {
										ic = pg.Theme.Icons.TimerDarkMode
									}
									return layout.Inset{
										Right: values.MarginPadding4,
										Top:   values.MarginPadding3,
									}.Layout(gtx, ic.Layout12dp)
								}
								return D{}
							}),
							layout.Rigid(timeagoLabel.Layout),
						)
					})
				}),
			)
		}),
		layout.Rigid(pg.lineSeparator(layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding10})),
		layout.Rigid(pg.layoutProposalVoteBar),
		layout.Rigid(func(gtx C) D {
			if proposal.Category != dcrlibwallet.ProposalCategoryActive {
				return D{}
			}
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.lineSeparator(layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding10})),
				layout.Rigid(pg.layoutProposalVoteAction),
			)
		}),
	)
}

func (pg *ProposalDetails) layoutTitle(gtx C) D {
	proposal := pg.proposal

	return pg.descriptionCard.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			if proposal.Category == dcrlibwallet.ProposalCategoryPre {
				return pg.layoutInDiscussionState(gtx)
			}
			return pg.layoutNormalTitle(gtx)
		})
	})
}

func (pg *ProposalDetails) layoutDescription(gtx C) D {
	grayCol := pg.Theme.Color.GrayText2
	proposal := pg.proposal

	dotLabel := pg.Theme.H4(" . ")
	dotLabel.Color = grayCol

	userLabel := pg.Theme.Body2(proposal.Username)
	userLabel.Color = grayCol

	versionLabel := pg.Theme.Body2("Version " + proposal.Version)
	versionLabel.Color = grayCol

	publishedLabel := pg.Theme.Body2("Published " + components.TimeAgo(proposal.PublishedAt))
	publishedLabel.Color = grayCol

	updatedLabel := pg.Theme.Body2("Updated " + components.TimeAgo(proposal.Timestamp))
	updatedLabel.Color = grayCol

	w := []layout.Widget{
		func(gtx C) D {
			lbl := pg.Theme.H5(proposal.Name)
			lbl.Font.Weight = text.SemiBold
			return lbl.Layout(gtx)
		},
		pg.lineSeparator(layout.Inset{Top: values.MarginPadding16, Bottom: values.MarginPadding16}),
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(userLabel.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPaddingMinus22}.Layout(gtx, dotLabel.Layout)
				}),
				layout.Rigid(publishedLabel.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPaddingMinus22}.Layout(gtx, dotLabel.Layout)
				}),
				layout.Rigid(versionLabel.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPaddingMinus22}.Layout(gtx, dotLabel.Layout)
				}),
				layout.Rigid(updatedLabel.Layout),
			)
		},
		pg.lineSeparator(layout.Inset{Top: values.MarginPadding16, Bottom: values.MarginPadding16}),
	}

	_, ok := pg.proposalItems[proposal.Token]
	if ok {
		w = append(w, pg.proposalItems[proposal.Token].widgets...)
	} else {
		loading := func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, layout.Flexed(1, func(gtx C) D {
				return layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8)}.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, material.Loader(pg.Theme.Base).Layout)
				})
			}))
		}

		w = append(w, loading)
	}

	w = append(w, pg.layoutRedirect("View on Politeia", pg.redirectIcon, pg.viewInPoliteiaBtn))

	return pg.descriptionCard.Layout(gtx, func(gtx C) D {
		return pg.Theme.List(pg.scrollbarList).Layout(gtx, 1, func(gtx C, i int) D {
			return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
				return pg.descriptionList.Layout(gtx, len(w), func(gtx C, i int) D {
					return layout.UniformInset(unit.Dp(0)).Layout(gtx, w[i])
				})
			})
		})
	})
}

func (pg *ProposalDetails) layoutRedirect(text string, icon *decredmaterial.Image, btn *decredmaterial.Clickable) layout.Widget {
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.lineSeparator(layout.Inset{Top: values.MarginPadding12, Bottom: values.MarginPadding12})),
			layout.Rigid(func(gtx C) D {
				return btn.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return pg.Theme.Body1(text).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{}.Layout(gtx, func(gtx C) D {
								return layout.E.Layout(gtx, icon.Layout24dp)
							})
						}),
					)
				})
			}),
		)
	}
}

func (pg *ProposalDetails) lineSeparator(inset layout.Inset) layout.Widget {
	return func(gtx C) D {
		return inset.Layout(gtx, pg.Theme.Separator().Layout)
	}
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *ProposalDetails) Layout(gtx C) D {
	proposal := pg.proposal
	_, ok := pg.proposalItems[proposal.Token]
	if !ok && !pg.loadingDescription {
		pg.loadingDescription = true
		go func() {
			var proposalDescription string
			if proposal.IndexFile != "" && proposal.IndexFileVersion == proposal.Version {
				proposalDescription = proposal.IndexFile
			} else {
				var err error
				proposalDescription, err = pg.WL.MultiWallet.Politeia.FetchProposalDescription(proposal.Token)
				if err != nil {
					fmt.Printf("Error loading proposal description: %v", err)
					time.Sleep(7 * time.Second)
					pg.loadingDescription = false
					return
				}
			}

			r := renderers.RenderMarkdown(gtx, pg.Theme, proposalDescription)
			proposalWidgets, proposalClickables := r.Layout()
			pg.proposalItems[proposal.Token] = proposalItemWidgets{
				widgets:    proposalWidgets,
				clickables: proposalClickables,
			}
			pg.loadingDescription = false
		}()
	}

	body := func(gtx C) D {
		page := components.SubPage{
			// App: pg.App,
			Title:      components.TruncateString(proposal.Name, 40),
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, pg.layoutTitle)
					}),
					layout.Rigid(pg.layoutDescription),
				)
			},
			ExtraItem: pg.viewInPoliteiaBtn,
			Extra: func(gtx C) D {
				return layout.Inset{}.Layout(gtx, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Top: values.MarginPadding5,
								}.Layout(gtx, pg.Theme.Caption("View in politeia").Layout)
							}),
							layout.Rigid(pg.redirectIcon.Layout24dp),
						)
					})
				})
			},
		}
		return page.Layout(gtx)
	}
	return components.UniformPadding(gtx, body)
}

package governance

import (
	"context"
	"fmt"
	"sync"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	// "github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

// const ModalInputVote = "input_vote_modal"

type agendaVoteModal struct {
	*load.Load
	modal decredmaterial.Modal

	detailsMu      sync.Mutex
	detailsCancel  context.CancelFunc
	// voteDetails    *dcrlibwallet.Agenda.Choices
	voteDetailsErr error

	agenda *dcrlibwallet.Agenda
	isVoting bool

	walletSelector *WalletSelector
	materialLoader material.LoaderStyle
	yesVote        decredmaterial.CheckBoxStyle
	// noVote         *inputVoteOptionsWidgets
	voteBtn        decredmaterial.Button
	cancelBtn      decredmaterial.Button
}

func newAgendaVoteModal(l *load.Load) *agendaVoteModal {
	avm := &agendaVoteModal{
		Load:           l,
		modal:          *l.Theme.ModalFloatTitle(),
		// agenda:       agenda,
		materialLoader: material.Loader(material.NewTheme(gofont.Collection())),
		voteBtn:        l.Theme.Button("Vote"),
		cancelBtn:      l.Theme.OutlineButton("Cancel"),
	}

	avm.voteBtn.Background = l.Theme.Color.Gray3
	avm.voteBtn.Color = l.Theme.Color.Surface

	// avm.yesVote = newInputVoteOptions(avm.Load, "Yes")
	// avm.noVote = newInputVoteOptions(avm.Load, "No")
	// avm.noVote.activeBg = l.Theme.Color.Orange2
	// avm.noVote.dotColor = l.Theme.Color.Danger

	avm.walletSelector = NewWalletSelector(l).
		Title("Voting wallet").
		WalletSelected(func(w *dcrlibwallet.Wallet) {

			avm.detailsMu.Lock()
			// avm.yesVote.reset()
			// avm.noVote.reset()
			// cancel current loading thread if any.
			if avm.detailsCancel != nil {
				avm.detailsCancel()
			}

			// ctx, cancel := context.WithCancel(context.Background())
			// avm.detailsCancel = cancel

			// avm.voteDetails = nil
			avm.voteDetailsErr = nil

			avm.detailsMu.Unlock()

			avm.RefreshWindow()

			// go func() {
			// 	// voteDetails, err := avm.WL.MultiWallet.Politeia.ProposalVoteDetailsRaw(w.ID, avm.proposal.Token)
			// 	avm.detailsMu.Lock()
			// 	if !components.ContextDone(ctx) {
			// 		// avm.voteDetails = voteDetails
			// 		avm.voteDetailsErr = err
			// 	}
			// 	avm.detailsMu.Unlock()
			// }()
		}).
		WalletValidator(func(w *dcrlibwallet.Wallet) bool {
			return !w.IsWatchingOnlyWallet()
		})
	return avm
}

func (avm *agendaVoteModal) ModalID() string {
	return ModalInputVote
}

func (avm *agendaVoteModal) OnResume() {
	avm.walletSelector.SelectFirstValidWallet()
}

func (avm *agendaVoteModal) OnDismiss() {

}

func (avm *agendaVoteModal) Show() {
	avm.ShowModal(avm)
}

func (avm *agendaVoteModal) Dismiss() {
	avm.DismissModal(avm)
}

// func (avm *agendaVoteModal) eligibleVotes() int {
// 	avm.detailsMu.Lock()
// 	// voteDetails := avm.voteDetails
// 	avm.detailsMu.Unlock()

// 	if voteDetails == nil {
// 		return 0
// 	}

// 	return len(voteDetails.EligibleTickets)
// }

// func (avm *agendaVoteModal) remainingVotes() int {
// 	avm.detailsMu.Lock()
// 	voteDetails := avm.voteDetails
// 	avm.detailsMu.Unlock()

// 	if voteDetails == nil {
// 		return 0
// 	}

// 	return len(voteDetails.EligibleTickets) - (avm.yesVote.voteCount() + avm.noVote.voteCount())
// }

func (avm *agendaVoteModal) sendVotes() {
	avm.detailsMu.Lock()
	// tickets := avm.voteDetails.EligibleTickets
	avm.detailsMu.Unlock()

	// votes := make([]*dcrlibwallet.ProposalVote, 0)
	// addVotes := func(bit string, count int) {
	// 	for i := 0; i < count; i++ {

	// 		// get and pop
	// 		var eligibleTicket *dcrlibwallet.EligibleTicket
	// 		eligibleTicket, tickets = tickets[0], tickets[1:]

	// 		vote := &dcrlibwallet.ProposalVote{
	// 			Ticket: eligibleTicket,
	// 			Bit:    bit,
	// 		}

	// 		votes = append(votes, vote)
	// 	}
	// }

	// addVotes(dcrlibwallet.VoteBitYes, 5)
	// addVotes(dcrlibwallet.VoteBitNo, 4)

	modal.NewPasswordModal(avm.Load).
		Title("Confirm to vote").
		NegativeButton("Cancel", func() {
			avm.isVoting = false
		}).
		PositiveButton("Confirm", func(password string, pm *modal.PasswordModal) bool {
			// go func() {
			// 	err := avm.WL.MultiWallet.Politeia.CastVotes(avm.walletSelector.selectedWallet.ID, votes, avm.proposal.Token, password)
			// 	if err != nil {
			// 		pm.SetError(err.Error())
			// 		pm.SetLoading(false)
			// 		return
			// 	}
			// 	pm.Dismiss()
			// 	avm.Toast.Notify("Vote sent successfully, refreshing proposals!")
			// 	go avm.WL.MultiWallet.Politeia.Sync()
			// 	avm.Dismiss()
			// }()

			return false
		}).Show()
}

func (avm *agendaVoteModal) Handle() {
	for avm.cancelBtn.Clicked() {
		if avm.isVoting {
			continue
		}
		avm.Dismiss()
	}

	// avm.handleVoteCountButtons(avm.yesVote)
	// avm.handleVoteCountButtons(avm.noVote)

	// totalVotes := avm.yesVote.voteCount() + avm.noVote.voteCount()
	// validToVote := totalVotes > 0 && totalVotes <= avm.eligibleVotes()
	// avm.voteBtn.SetEnabled(validToVote)

	for avm.voteBtn.Clicked() {
		if avm.isVoting {
			break
		}

		// if !validToVote {
		// 	break
		// }

		avm.isVoting = true
		// avm.sendVotes()
	}
}

// - Layout

func (avm *agendaVoteModal) Layout(gtx layout.Context) D {
	avm.detailsMu.Lock()
	// voteDetails := vm.voteDetails
	voteDetailsErr := avm.voteDetailsErr
	avm.detailsMu.Unlock()
	w := []layout.Widget{
		func(gtx C) D {
			t := avm.Theme.H6("Vote")
			t.Font.Weight = text.SemiBold
			return t.Layout(gtx)
		},
		func(gtx C) D {
			return avm.walletSelector.Layout(gtx)
		},
		func(gtx C) D {
			// if voteDetails != nil {
			// 	return D{}
			// }

			if voteDetailsErr != nil {
				return avm.Theme.Label(values.TextSize16, voteDetailsErr.Error()).Layout(gtx)
			}

			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding24)
			return avm.materialLoader.Layout(gtx)
		},
		func(gtx C) D {
			// if voteDetails == nil {
			// 	return D{}
			// }

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtc C) D {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								// if voteDetails.YesVotes == 0 {
								// 	return layout.Dimensions{}
								// }

								wrap := avm.Theme.Card()
								wrap.Color = avm.Theme.Color.Green50
								wrap.Radius = decredmaterial.Radius(8)
								// if voteDetails.NoVotes > 0 {
								// 	wrap.Radius.TopRight = 0
								// 	wrap.Radius.BottomRight = 0
								// }
								return wrap.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									inset := layout.Inset{
										Left:   values.MarginPadding12,
										Top:    values.MarginPadding8,
										Right:  values.MarginPadding12,
										Bottom: values.MarginPadding8,
									}
									return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												card := avm.Theme.Card()
												card.Color = avm.Theme.Color.Green500
												card.Radius = decredmaterial.Radius(4)
												return card.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													gtx.Constraints.Min.X += gtx.Px(values.MarginPadding8)
													gtx.Constraints.Min.Y += gtx.Px(values.MarginPadding8)
													return layout.Dimensions{Size: gtx.Constraints.Min}
												})
											}),
											layout.Rigid(func(gtx C) D {
												return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
													label := avm.Theme.Body2(fmt.Sprintf("Yes: %d", 3))
													return label.Layout(gtx)
												})
											}),
										)
									})
								})
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								// if voteDetails.NoVotes == 0 {
								// 	return layout.Dimensions{}
								// }

								wrap := avm.Theme.Card()
								wrap.Color = avm.Theme.Color.Orange2
								wrap.Radius = decredmaterial.Radius(8)
								// if voteDetails.YesVotes > 0 {
								// 	wrap.Radius.TopLeft = 0
								// 	wrap.Radius.BottomLeft = 0
								// }
								return wrap.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									inset := layout.Inset{
										Left:   values.MarginPadding12,
										Top:    values.MarginPadding8,
										Right:  values.MarginPadding12,
										Bottom: values.MarginPadding8,
									}
									return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												card := avm.Theme.Card()
												card.Color = avm.Theme.Color.Danger
												card.Radius = decredmaterial.Radius(4)
												return card.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
													gtx.Constraints.Min.X += gtx.Px(values.MarginPadding8)
													gtx.Constraints.Min.Y += gtx.Px(values.MarginPadding8)
													return layout.Dimensions{Size: gtx.Constraints.Min}
												})
											}),
											layout.Rigid(func(gtx C) D {
												return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
													label := avm.Theme.Body2(fmt.Sprintf("No: %d", 3))
													return label.Layout(gtx)
												})
											}),
										)
									})
								})
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					// if voteDetails == nil {
					// 	return D{}
					// }

					text := fmt.Sprintf("You have %d votes", 7)
					return avm.Theme.Label(values.TextSize16, text).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return avm.yesVote.Layout(gtx)
					// return avm.inputOptions(gtx, avm.yesVote)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top: values.MarginPadding8,
					}.Layout(gtx, func(gtx C) D {
						return avm.yesVote.Layout(gtx)
						// return avm.inputOptions(gtx, avm.noVote)
					})
				}),
			)
		},
		func(gtx C) D {
			// if voteDetails != nil && vm.yesVote.voteCount()+vm.noVote.voteCount() > len(voteDetails.EligibleTickets) {
			// 	label := vm.Theme.Label(values.TextSize14, "You don’t have enough votes")
			// 	label.Color = vm.Theme.Color.Danger
			// 	return label.Layout(gtx)
			// }

			return D{}
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, avm.cancelBtn.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						if avm.isVoting {
							return avm.materialLoader.Layout(gtx)
						}
						return avm.voteBtn.Layout(gtx)
					}),
				)
			})
		},
	}

	return avm.modal.Layout(gtx, w)
}

func (avm *agendaVoteModal) inputOptions(gtx layout.Context, wdg *inputVoteOptionsWidgets) D {
	wrap := avm.Theme.Card()
	wrap.Color = avm.Theme.Color.Gray4
	dotColor := avm.Theme.Color.Gray3
	if wdg.voteCount() > 0 {
		wrap.Color = wdg.activeBg
		dotColor = wdg.dotColor
	}
	return wrap.Layout(gtx, func(gtx C) D {
		inset := layout.Inset{
			Top:    values.MarginPadding8,
			Bottom: values.MarginPadding8,
			Left:   values.MarginPadding16,
			Right:  values.MarginPadding8,
		}
		return inset.Layout(gtx, func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(.4, func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							card := avm.Theme.Card()
							card.Color = dotColor
							card.Radius = decredmaterial.Radius(4)
							return card.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Min.X += gtx.Px(values.MarginPadding8)
								gtx.Constraints.Min.Y += gtx.Px(values.MarginPadding8)
								return layout.Dimensions{Size: gtx.Constraints.Min}
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
								return avm.Theme.Body2(wdg.label).Layout(gtx)
							})
						}),
					)
				}),
				layout.Flexed(.6, func(gtx C) D {
					border := widget.Border{
						Color:        avm.Theme.Color.Gray2,
						CornerRadius: values.MarginPadding8,
						Width:        values.MarginPadding2,
					}

					return border.Layout(gtx, func(gtx C) D {
						card := avm.Theme.Card()
						card.Color = avm.Theme.Color.Surface
						return card.Layout(gtx, func(gtx C) D {
							var height int
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
								layout.Flexed(1, func(gtx C) D {
									dims := layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return wdg.decrement.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											gtx.Constraints.Min.X, gtx.Constraints.Max.X = 30, 30
											return wdg.input.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return wdg.increment.Layout(gtx)
										}),
									)
									height = dims.Size.Y
									return dims
								}),
								layout.Flexed(0.02, func(gtx C) D {
									line := avm.Theme.Line(height, gtx.Px(values.MarginPadding2))
									line.Color = avm.Theme.Color.Gray2
									return line.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									return wdg.max.Layout(gtx)
								}),
							)
						})
					})
				}),
			)
		})
	})
}

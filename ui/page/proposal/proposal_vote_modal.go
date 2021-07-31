package proposal

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
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const ModalInputVote = "input_vote_modal"

type voteModal struct {
	*load.Load
	modal decredmaterial.Modal

	detailsMu      sync.Mutex
	detailsCancel  context.CancelFunc
	voteDetails    *dcrlibwallet.ProposalVoteDetails
	voteDetailsErr error

	proposal *dcrlibwallet.Proposal
	isVoting bool

	walletSelector *WalletSelector
	materialLoader material.LoaderStyle
	yesVote        *inputVoteOptionsWidgets
	noVote         *inputVoteOptionsWidgets
	voteBtn        decredmaterial.Button
	cancelBtn      decredmaterial.Button
}

func newVoteModal(l *load.Load, proposal *dcrlibwallet.Proposal) *voteModal {
	vm := &voteModal{
		Load:           l,
		modal:          *l.Theme.ModalFloatTitle(),
		proposal:       proposal,
		materialLoader: material.Loader(material.NewTheme(gofont.Collection())),
		voteBtn:        l.Theme.Button(new(widget.Clickable), "Vote"),
		cancelBtn:      l.Theme.Button(new(widget.Clickable), "Cancel"),
	}

	vm.cancelBtn.Background = vm.Theme.Color.Surface
	vm.cancelBtn.Color = vm.Theme.Color.Primary

	vm.voteBtn.TextSize, vm.cancelBtn.TextSize = values.TextSize16, values.TextSize16
	vm.voteBtn.Font.Weight, vm.cancelBtn.Font.Weight = text.Bold, text.Bold
	vm.voteBtn.Background = l.Theme.Color.Gray1
	vm.voteBtn.Color = l.Theme.Color.Surface

	vm.yesVote = newInputVoteOptions(vm.Load, "Yes")
	vm.noVote = newInputVoteOptions(vm.Load, "No")
	vm.noVote.activeBg = l.Theme.Color.Orange2
	vm.noVote.dotColor = l.Theme.Color.Danger

	vm.walletSelector = NewWalletSelector(l).
		Title("Voting wallet").
		WalletSelected(func(w *dcrlibwallet.Wallet) {

			vm.detailsMu.Lock()
			vm.yesVote.reset()
			vm.noVote.reset()
			// cancel current loading thread if any.
			if vm.detailsCancel != nil {
				vm.detailsCancel()
			}

			ctx, cancel := context.WithCancel(context.Background())
			vm.detailsCancel = cancel

			vm.voteDetails = nil
			vm.voteDetailsErr = nil

			vm.detailsMu.Unlock()

			vm.RefreshWindow()

			go func() {

				voteDetails, err := vm.WL.MultiWallet.Politeia.ProposalVoteDetailsRaw(w.ID, vm.proposal.Token)
				vm.detailsMu.Lock()
				if !components.ContextDone(ctx) {
					vm.voteDetails = voteDetails
					vm.voteDetailsErr = err
				}
				vm.detailsMu.Unlock()
			}()
		}).
		WalletValidator(func(w *dcrlibwallet.Wallet) bool {
			return !w.IsWatchingOnlyWallet()
		})
	return vm
}

func (cm *voteModal) ModalID() string {
	return ModalInputVote
}

func (vm *voteModal) OnResume() {
	vm.walletSelector.SelectFirstValidWallet()
}

func (cm *voteModal) OnDismiss() {

}

func (cm *voteModal) Show() {
	cm.ShowModal(cm)
}

func (cm *voteModal) Dismiss() {
	cm.DismissModal(cm)
}

func (vm *voteModal) eligibleVotes() int {
	vm.detailsMu.Lock()
	voteDetails := vm.voteDetails
	vm.detailsMu.Unlock()

	if voteDetails == nil {
		return 0
	}

	return len(voteDetails.EligibleTickets)
}

func (vm *voteModal) remainingVotes() int {
	vm.detailsMu.Lock()
	voteDetails := vm.voteDetails
	vm.detailsMu.Unlock()

	if voteDetails == nil {
		return 0
	}

	return len(voteDetails.EligibleTickets) - (vm.yesVote.voteCount() + vm.noVote.voteCount())
}

func (vm *voteModal) sendVotes() {
	vm.detailsMu.Lock()
	tickets := vm.voteDetails.EligibleTickets
	vm.detailsMu.Unlock()

	votes := make([]*dcrlibwallet.ProposalVote, 0)
	addVotes := func(bit string, count int) {
		for i := 0; i < count; i++ {

			// get and pop
			var eligibleTicket *dcrlibwallet.EligibleTicket
			eligibleTicket, tickets = tickets[0], tickets[1:]

			vote := &dcrlibwallet.ProposalVote{
				Ticket: eligibleTicket,
				Bit:    bit,
			}

			votes = append(votes, vote)
		}
	}

	addVotes(dcrlibwallet.VoteBitYes, vm.yesVote.voteCount())
	addVotes(dcrlibwallet.VoteBitNo, vm.noVote.voteCount())

	modal.NewPasswordModal(vm.Load).
		Title("Confirm to vote").
		NegativeButton("Cancel", func() {
			vm.isVoting = false
		}).
		PositiveButton("Confirm", func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := vm.WL.MultiWallet.Politeia.CastVotes(vm.walletSelector.selectedWallet.ID, votes, vm.proposal.Token, password)
				if err != nil {
					pm.SetError(err.Error())
					pm.SetLoading(false)
					return
				}
				pm.Dismiss()
				vm.CreateToast("Vote sent successfully, refreshing proposals!", true)
				go vm.WL.MultiWallet.Politeia.Sync()
				vm.Dismiss()
			}()

			return false
		}).Show()
}

func (vm *voteModal) Handle() {
	for vm.cancelBtn.Button.Clicked() {
		if vm.isVoting {
			continue
		}
		vm.Dismiss()
	}

	vm.handleVoteCountButtons(vm.yesVote)
	vm.handleVoteCountButtons(vm.noVote)

	totalVotes := vm.yesVote.voteCount() + vm.noVote.voteCount()
	validToVote := totalVotes > 0 && totalVotes <= vm.eligibleVotes()
	vm.voteBtn.SetEnabled(validToVote)

	for vm.voteBtn.Clicked() {
		if vm.isVoting {
			break
		}

		if !validToVote {
			break
		}

		vm.isVoting = true
		vm.sendVotes()
	}
}

// - Layout

func (cm *voteModal) Layout(gtx layout.Context) D {
	cm.detailsMu.Lock()
	voteDetails := cm.voteDetails
	voteDetailsErr := cm.voteDetailsErr
	cm.detailsMu.Unlock()
	w := []layout.Widget{
		func(gtx C) D {
			t := cm.Theme.H6("Vote")
			t.Font.Weight = text.Bold
			return t.Layout(gtx)
		},
		func(gtx C) D {
			if voteDetails == nil {
				return D{}
			}

			text := fmt.Sprintf("You have %d votes", len(voteDetails.EligibleTickets))
			return cm.Theme.Label(values.TextSize16, text).Layout(gtx)
		},
		func(gtx C) D {
			return cm.walletSelector.Layout(gtx)
		},
		func(gtx C) D {
			if voteDetails != nil {
				return D{}
			}

			if voteDetailsErr != nil {
				return cm.Theme.Label(values.TextSize16, voteDetailsErr.Error()).Layout(gtx)
			}

			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding24)
			return cm.materialLoader.Layout(gtx)
		},
		func(gtx C) D {
			if voteDetails == nil {
				return D{}
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						// Top: values.MarginPadding16,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return cm.inputOptions(gtx, cm.yesVote)
					})

				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top: values.MarginPadding8,
					}.Layout(gtx, func(gtx C) D {
						return cm.inputOptions(gtx, cm.noVote)
					})
				}),
			)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, cm.cancelBtn.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						if cm.isVoting {
							return cm.materialLoader.Layout(gtx)
						}
						return cm.voteBtn.Layout(gtx)
					}),
				)
			})
		},
	}

	return cm.modal.Layout(gtx, w, 850)
}

func (cm *voteModal) inputOptions(gtx layout.Context, wdg *inputVoteOptionsWidgets) D {
	wrap := cm.Theme.Card()
	wrap.Color = cm.Theme.Color.LightGray
	dotColor := cm.Theme.Color.InactiveGray
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
							card := cm.Theme.Card()
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
								label := cm.Theme.Body2(wdg.label)
								label.Color = cm.Theme.Color.DeepBlue
								return label.Layout(gtx)
							})
						}),
					)
				}),
				layout.Flexed(.6, func(gtx C) D {
					border := widget.Border{
						Color:        cm.Theme.Color.Gray1,
						CornerRadius: values.MarginPadding8,
						Width:        values.MarginPadding2,
					}

					return border.Layout(gtx, func(gtx C) D {
						card := cm.Theme.Card()
						card.Color = cm.Theme.Color.Surface
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
											gtx.Constraints.Min.X, gtx.Constraints.Max.X = 100, 100
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
									line := cm.Theme.Line(height, gtx.Px(values.MarginPadding2))
									line.Color = cm.Theme.Color.Gray1
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

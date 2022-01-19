package governance

import (
	"context"
	// "fmt"
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

	detailsMu     sync.Mutex
	detailsCancel context.CancelFunc
	// voteDetails    *dcrlibwallet.Agenda.Choices
	voteDetailsErr error

	agenda   *dcrlibwallet.Agenda
	isVoting bool

	consensusPage *ConsensusPage

	walletSelector *WalletSelector
	materialLoader material.LoaderStyle
	abstainVote    decredmaterial.CheckBoxStyle
	yesVote        decredmaterial.CheckBoxStyle
	noVote         decredmaterial.CheckBoxStyle
	// agendaVoteOptions         []decredmaterial.CheckBoxStyle
	items             map[string]string //[key]str-key
	itemKeys          []string
	defaultValue      string // str-key
	initialValue      string
	currentValue      string
	optionsRadioGroup *widget.Enum
	voteBtn           decredmaterial.Button
	cancelBtn         decredmaterial.Button
}

func newAgendaVoteModal(l *load.Load, agenda *dcrlibwallet.Agenda) *agendaVoteModal {
	avm := &agendaVoteModal{
		Load:          l,
		modal:         *l.Theme.ModalFloatTitle(),
		agenda:        agenda,
		consensusPage: NewConsensusPage(l),
		// defaultValue: agenda.VotingPreference,
		materialLoader:    material.Loader(material.NewTheme(gofont.Collection())),
		abstainVote:       l.Theme.CheckBox(new(widget.Bool), "Abstain"),
		yesVote:           l.Theme.CheckBox(&widget.Bool{}, "Yes"),
		noVote:            l.Theme.CheckBox(&widget.Bool{}, "No"),
		optionsRadioGroup: new(widget.Enum),
		voteBtn:           l.Theme.Button("Vote"),
		cancelBtn:         l.Theme.OutlineButton("Cancel"),
	}

	avm.voteBtn.Background = l.Theme.Color.Gray3
	avm.voteBtn.Color = l.Theme.Color.Surface

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

	ArrVoteOptions := make(map[string]string)
	for i := range agenda.Choices {
		ArrVoteOptions[agenda.Choices[i].Id] = agenda.Choices[i].Id
	}

	// sort keys to keep order when refreshed
	sortedKeys := make([]string, 0)
	for k := range ArrVoteOptions {
		sortedKeys = append(sortedKeys, k)
	}
	avm.itemKeys = sortedKeys
	avm.items = ArrVoteOptions
	return avm
}

func (avm *agendaVoteModal) ModalID() string {
	return ModalInputVote
}

func (avm *agendaVoteModal) OnResume() {
	avm.walletSelector.SelectFirstValidWallet()

	initialValue := avm.agenda.VotingPreference
	if initialValue == "" {
		initialValue = avm.defaultValue
	}

	avm.initialValue = initialValue
	avm.currentValue = initialValue

	avm.optionsRadioGroup.Value = avm.currentValue
}

func (avm *agendaVoteModal) OnDismiss() {

}

func (avm *agendaVoteModal) Show() {
	avm.ShowModal(avm)
}

func (avm *agendaVoteModal) Dismiss() {
	avm.DismissModal(avm)
}

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

	modal.NewPasswordModal(avm.Load).
		Title("Confirm to vote").
		NegativeButton("Cancel", func() {
			avm.isVoting = false
		}).
		PositiveButton("Confirm", func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := avm.WL.MultiWallet.Consensus.SetVoteChoice(avm.walletSelector.selectedWallet.ID, "", avm.agenda.AgendaID, avm.optionsRadioGroup.Value, "", password)
				if err != nil {
					pm.SetError(err.Error())
					pm.SetLoading(false)
					return
				}
				pm.Dismiss()
				avm.Toast.Notify("Vote updated successfully, refreshing agendas!")
				go avm.WL.MultiWallet.Consensus.GetAllAgendasForWallet(avm.walletSelector.selectedWallet.ID, false)
				avm.Dismiss()
			}()

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

	for avm.optionsRadioGroup.Changed() {
		avm.currentValue = avm.optionsRadioGroup.Value
		// avm.wallet.SaveConfigValueForKey(avm.preferenceKey, avm.optionsRadioGroup.Value)
		// avm.updateButtonClicked()
	}

	validToVote := avm.optionsRadioGroup.Value != "" && avm.optionsRadioGroup.Value != avm.initialValue
	avm.voteBtn.SetEnabled(validToVote)
	if avm.voteBtn.Enabled() {
		avm.voteBtn.Background = avm.Theme.Color.Primary
	}

	for avm.voteBtn.Clicked() {
		if avm.isVoting {
			break
		}

		if !validToVote {
			break
		}

		avm.isVoting = true
		avm.sendVotes()
	}

	if avm.modal.BackdropClicked(true) {
		avm.Dismiss()
	}
}

// - Layout

func (avm *agendaVoteModal) Layout(gtx layout.Context) D {
	avm.detailsMu.Lock()
	avm.detailsMu.Unlock()
	w := []layout.Widget{
		func(gtx C) D {
			t := avm.Theme.H6("Change Vote")
			t.Font.Weight = text.SemiBold
			return t.Layout(gtx)
		},
		func(gtx C) D {
			t := avm.Theme.Body1("Select one of the options below to vote")
			return t.Layout(gtx)
		},
		func(gtx C) D {
			return avm.walletSelector.Layout(gtx)
		},
		func(gtx C) D {
			// if voteDetails == nil {
			// 	return D{}
			// }

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx, avm.layoutItems()...)
				}),
			)
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

func (avm *agendaVoteModal) layoutItems() []layout.FlexChild {

	items := make([]layout.FlexChild, 0)
	for _, k := range avm.itemKeys {
		radioItem := layout.Rigid(avm.Load.Theme.RadioButton(avm.optionsRadioGroup, k, avm.items[k], avm.Load.Theme.Color.DeepBlue, avm.Load.Theme.Color.Primary).Layout)

		items = append(items, radioItem)
	}

	return items
}

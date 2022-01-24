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
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/page/staking"
	"github.com/planetdecred/godcr/ui/values"
)

// const ModalInputVote = "input_vote_modal"

type agendaVoteModal struct {
	*load.Load
	modal          decredmaterial.Modal
	voteSuccessful func()

	detailsMu     sync.Mutex
	detailsCancel context.CancelFunc
	LiveTickets   []*dcrlibwallet.Transaction

	agenda               *dcrlibwallet.Agenda
	vspIsFetched         bool
	liveTicketsIsFetched bool
	isVoting             bool
	loadCount            int

	consensusPage *ConsensusPage

	walletSelector    *WalletSelector
	vspSelector       *staking.VSPSelector
	ticketSelector    *ticketSelector
	materialLoader    material.LoaderStyle
	items             map[string]string //[key]str-key
	itemKeys          []string
	defaultValue      string // str-key
	initialValue      string
	currentValue      string
	optionsRadioGroup *widget.Enum
	voteBtn           decredmaterial.Button
	cancelBtn         decredmaterial.Button
}

func newAgendaVoteModal(l *load.Load, agenda *dcrlibwallet.Agenda, consensusPage *ConsensusPage) *agendaVoteModal {
	avm := &agendaVoteModal{
		Load:          l,
		modal:         *l.Theme.ModalFloatTitle(),
		agenda:        agenda,
		consensusPage: consensusPage,
		// voteSuccessful: onVoteSuccessful,
		materialLoader:    material.Loader(material.NewTheme(gofont.Collection())),
		optionsRadioGroup: new(widget.Enum),
		voteBtn:           l.Theme.Button("Vote"),
		cancelBtn:         l.Theme.OutlineButton("Cancel"),
	}

	avm.voteBtn.Background = l.Theme.Color.Gray3
	avm.voteBtn.Color = l.Theme.Color.Surface

	avm.walletSelector = NewWalletSelector(l).
		Title("Voting wallet").
		WalletSelected(func(w *dcrlibwallet.Wallet) {
			avm.loadCount = 0
			avm.detailsMu.Lock()
			// cancel current loading thread if any.
			if avm.detailsCancel != nil {
				avm.detailsCancel()
			}

			avm.detailsMu.Unlock()
			avm.FetchLiveTickets(w.ID)
			avm.RefreshWindow()
			avm.loadCount++
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

	avm.vspIsFetched = len((*l.WL.VspInfo).List) > 0

	return avm
}

func (avm *agendaVoteModal) FetchLiveTickets(walletID int) {
	go func() {
		avm.liveTicketsIsFetched = false

		wallet := avm.WL.MultiWallet.WalletWithID(walletID)
		tickets, err := staking.WalletLiveTickets(wallet)
		if err != nil {
			avm.Toast.NotifyError(err.Error())
			return
		}

		liveTickets := make([]*dcrlibwallet.Transaction, 0)
		txItems, err := staking.StakeToTransactionItems(avm.Load, tickets, true, func(filter int32) bool {
			switch filter {
			case dcrlibwallet.TxFilterUnmined:
				fallthrough
			case dcrlibwallet.TxFilterImmature:
				fallthrough
			case dcrlibwallet.TxFilterLive:
				return true
			}

			return false
		})
		if err != nil {
			avm.Toast.NotifyError(err.Error())
			return
		}

		for _, liveTicket := range txItems {
			liveTickets = append(liveTickets, liveTicket.Transaction)
		}

		avm.LiveTickets = liveTickets
		avm.liveTicketsIsFetched = true
		avm.RefreshWindow()
	}()
}

func (avm *agendaVoteModal) ModalID() string {
	return ModalInputVote
}

func (avm *agendaVoteModal) OnResume() {
	avm.walletSelector.SelectFirstValidWallet()

	avm.vspSelector = staking.NewVSPSelector(avm.Load).Title("Select a vsp")
	if avm.vspIsFetched && components.StringNotEmpty(avm.WL.GetRememberVSP()) {
		avm.vspSelector.SelectVSP(avm.WL.GetRememberVSP())
	}

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
	avm.detailsMu.Unlock()

	modal.NewPasswordModal(avm.Load).
		Title("Confirm to vote").
		NegativeButton("Cancel", func() {
			avm.isVoting = false
		}).
		PositiveButton("Confirm", func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := avm.WL.MultiWallet.Consensus.SetVoteChoice(avm.walletSelector.selectedWallet.ID, avm.vspSelector.SelectedVSP().Info.PubKey, avm.vspSelector.SelectedVSP().Host, avm.agenda.AgendaID, avm.optionsRadioGroup.Value, "", password)
				// err := avm.WL.MultiWallet.Consensus.SetVoteChoice(avm.walletSelector.selectedWallet.ID, avm.vspSelector.SelectedVSP().Info.PubKey, avm.vspSelector.SelectedVSP().Host, avm.agenda.AgendaID, avm.optionsRadioGroup.Value, avm.ticketSelector.SelectedTicket().Hash, password)
				if err != nil {
					pm.SetError(err.Error())
					pm.SetLoading(false)
					return
				}
				pm.Dismiss()
				avm.Toast.Notify("Vote updated successfully, refreshing agendas!")

				avm.Dismiss()
				go avm.consensusPage.FetchAgendas()
			}()

			return false
		}).Show()
}

func (avm *agendaVoteModal) Handle() {
	if avm.vspSelector.Changed() {
		avm.WL.RememberVSP(avm.vspSelector.SelectedVSP().Host)
	}

	// reselect vsp if there's a delay in fetching the VSP List
	if !avm.vspIsFetched && len((*avm.WL.VspInfo).List) > 0 {
		if avm.WL.GetRememberVSP() != "" {
			avm.vspSelector.SelectVSP(avm.WL.GetRememberVSP())
			avm.vspIsFetched = true
		}
	}

	for avm.cancelBtn.Clicked() {
		if avm.isVoting {
			continue
		}
		avm.Dismiss()
	}

	for avm.optionsRadioGroup.Changed() {
		avm.currentValue = avm.optionsRadioGroup.Value
	}

	if avm.liveTicketsIsFetched {
		if avm.loadCount == 1 {
			avm.loadCount++
			avm.ticketSelector = newTicketSelector(avm.Load, avm.LiveTickets).Title("Select a ticket")
		}

	}

	validToVote := avm.optionsRadioGroup.Value != "" && avm.optionsRadioGroup.Value != avm.initialValue && avm.vspSelector.SelectedVSP() != nil && avm.ticketSelector.SelectedTicket() != nil
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
			if len(avm.LiveTickets) < 1 {
				return D {}
			}
			return avm.vspSelector.Layout(gtx)
		},
		func(gtx C) D {
			if !avm.liveTicketsIsFetched {
				gtx.Constraints.Min.X = gtx.Px(values.MarginPadding24)
				return avm.materialLoader.Layout(gtx)
			}
			if len(avm.LiveTickets) < 1 {
				return D {}
			}
			return avm.ticketSelector.Layout(gtx)
		},
		func(gtx C) D {
			if !avm.liveTicketsIsFetched {
				gtx.Constraints.Min.X = gtx.Px(values.MarginPadding24)
				return avm.materialLoader.Layout(gtx)
			}
			text := fmt.Sprintf("You have %d live tickets for the selected wallet [%s]", len(avm.LiveTickets), avm.walletSelector.SelectedWallet().Name)
			return avm.Theme.Label(values.TextSize16, text).Layout(gtx)
		},
		func(gtx C) D {
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

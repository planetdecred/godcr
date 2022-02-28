package governance

import (
	"context"
	"fmt"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

type agendaVoteModal struct {
	*load.Load
	modal decredmaterial.Modal

	LiveTickets []*dcrlibwallet.Transaction

	agenda    *dcrlibwallet.Agenda
	isVoting  bool
	loadCount int

	consensusPage *ConsensusPage

	walletSelector    *WalletSelector
	vspSelector       *components.VSPSelector
	ticketSelector    *ticketSelector
	spendingPassword  decredmaterial.Editor
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
		Load:              l,
		modal:             *l.Theme.ModalFloatTitle(),
		agenda:            agenda,
		consensusPage:     consensusPage,
		materialLoader:    material.Loader(material.NewTheme(gofont.Collection())),
		optionsRadioGroup: new(widget.Enum),
		spendingPassword:  l.Theme.EditorPassword(new(widget.Editor), "Spending password"),
		voteBtn:           l.Theme.Button("Vote"),
		cancelBtn:         l.Theme.OutlineButton("Cancel"),
	}

	avm.voteBtn.Background = l.Theme.Color.Gray3
	avm.voteBtn.Color = l.Theme.Color.Surface

	avm.walletSelector = NewWalletSelector(l).
		Title("Voting wallet").
		WalletSelected(func(w *dcrlibwallet.Wallet) {
			avm.loadCount = 0
			avm.FetchLiveTickets(w.ID)
			avm.RefreshWindow()
			avm.loadCount++

			// update agenda options prefrence to that of the selected wallet
			consensusItems := components.LoadAgendas(avm.Load, w, false)
			for _, consensusItem := range consensusItems {
				if consensusItem.Agenda.AgendaID == agenda.AgendaID {
					ArrVoteOptions := make(map[string]string)
					for i := range consensusItem.Agenda.Choices {
						ArrVoteOptions[agenda.Choices[i].Id] = consensusItem.Agenda.Choices[i].Id
					}

					// sort keys to keep order when refreshed
					sortedKeys := make([]string, 0)
					for k := range ArrVoteOptions {
						sortedKeys = append(sortedKeys, k)
					}
					avm.itemKeys = sortedKeys
					avm.items = ArrVoteOptions

					initialValue := consensusItem.Agenda.VotingPreference
					if initialValue == "" {
						initialValue = avm.defaultValue
					}

					avm.initialValue = initialValue
					avm.currentValue = initialValue

					avm.optionsRadioGroup.Value = avm.currentValue
				}
			}
		}).
		WalletValidator(func(w *dcrlibwallet.Wallet) bool {
			return !w.IsWatchingOnlyWallet()
		})

	return avm
}

func (avm *agendaVoteModal) FetchLiveTickets(walletID int) {
	go func() {
		wallet := avm.WL.MultiWallet.WalletWithID(walletID)
		tickets, err := components.WalletLiveTickets(wallet)
		if err != nil {
			avm.Toast.NotifyError(err.Error())
			return
		}

		liveTickets := make([]*dcrlibwallet.Transaction, 0)
		txItems, err := components.StakeToTransactionItems(avm.Load, tickets, true, func(filter int32) bool {
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
		avm.RefreshWindow()
	}()
}

func (avm *agendaVoteModal) ModalID() string {
	return ModalInputVote
}

func (avm *agendaVoteModal) OnResume() {
	avm.walletSelector.SelectFirstValidWallet()

	avm.vspSelector = components.NewVSPSelector(avm.Load).Title("Select a vsp")
	lastUsedVSP := avm.WL.MultiWallet.LastUsedVSP()
	if len(avm.WL.MultiWallet.KnownVSPs()) == 0 {
		// TODO: Does this modal need this list?
		go avm.WL.MultiWallet.ReloadVSPList(context.TODO())
	} else if components.StringNotEmpty(lastUsedVSP) {
		avm.vspSelector.SelectVSP(lastUsedVSP)
	}
	initialValue := avm.agenda.VotingPreference
	if initialValue == "" {
		initialValue = avm.defaultValue
	}

	avm.initialValue = initialValue
	avm.currentValue = initialValue

	avm.optionsRadioGroup.Value = avm.currentValue
}

func (avm *agendaVoteModal) OnDismiss() {}

func (avm *agendaVoteModal) Show() {
	avm.ShowModal(avm)
}

func (avm *agendaVoteModal) Dismiss() {
	avm.DismissModal(avm)
}

func (avm *agendaVoteModal) Handle() {
	if avm.vspSelector.Changed() {
		avm.WL.MultiWallet.SaveLastUsedVSP(avm.vspSelector.SelectedVSP().Host)
	}

	// reselect vsp if there's a delay in fetching the VSP List
	lastUsedVSP := avm.WL.MultiWallet.LastUsedVSP()
	if len(avm.WL.MultiWallet.KnownVSPs()) > 0 && lastUsedVSP != "" {
		avm.vspSelector.SelectVSP(lastUsedVSP)
	}

	for avm.cancelBtn.Clicked() {
		if avm.isVoting {
			continue
		}
		avm.Dismiss()
	}

	_, isChanged := decredmaterial.HandleEditorEvents(avm.spendingPassword.Editor)
	if isChanged {
		avm.spendingPassword.SetError("")
	}

	for avm.optionsRadioGroup.Changed() {
		avm.currentValue = avm.optionsRadioGroup.Value
	}

	if len(avm.LiveTickets) != 0 {
		if avm.loadCount == 1 {
			avm.loadCount++
			avm.ticketSelector = newTicketSelector(avm.Load, avm.LiveTickets).Title("Select a ticket")
		}

	}

	validToVote := avm.optionsRadioGroup.Value != "" && avm.optionsRadioGroup.Value != avm.initialValue && avm.vspSelector.SelectedVSP() != nil && avm.spendingPassword.Editor.Text() != ""
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
	w := []layout.Widget{
		func(gtx C) D {
			t := avm.Theme.H6("Change Vote")
			t.Font.Weight = text.SemiBold
			return t.Layout(gtx)
		},
		func(gtx C) D {
			return avm.Theme.Body1("Select one of the options below to vote").Layout(gtx)
		},
		func(gtx C) D {
			return avm.walletSelector.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if len(avm.LiveTickets) == 0 {
						return D{}
					}
					return avm.vspSelector.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					// to we need this loader? it seems redundant.
					// if len(avm.LiveTickets) == 0 {
					// 	gtx.Constraints.Min.X = gtx.Px(values.MarginPadding24)
					// 	return avm.materialLoader.Layout(gtx)
					// }
					var ticketCountLabel decredmaterial.Label
					text := fmt.Sprintf("You have %d live tickets for the selected wallet [%s]", len(avm.LiveTickets), avm.walletSelector.SelectedWallet().Name)
					ticketCountLabel = avm.Theme.Label(values.MarginPadding14, text)
					if len(avm.LiveTickets) < 1 {
						ticketCountLabel.Color = avm.Theme.Color.Danger
					}
					return ticketCountLabel.Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if len(avm.LiveTickets) == 0 {
						return D{}
					}
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, avm.layoutItems()...)
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if len(avm.LiveTickets) < 1 {
						return D{}
					}
					return avm.spendingPassword.Layout(gtx)
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
		radioBtn := avm.Load.Theme.RadioButton(avm.optionsRadioGroup, k, avm.items[k], avm.Load.Theme.Color.DeepBlue, avm.Load.Theme.Color.Primary)
		radioItem := layout.Rigid(radioBtn.Layout)
		items = append(items, radioItem)
	}

	return items
}

func (avm *agendaVoteModal) sendVotes() {
	go func() {
		password := []byte(avm.spendingPassword.Editor.Text())

		defer func() {
			avm.isVoting = false
		}()

		vspHost := avm.vspSelector.SelectedVSP().Host
		pubKey := avm.vspSelector.SelectedVSP().PubKey
		radioBtnGrp := avm.optionsRadioGroup.Value
		err := avm.walletSelector.selectedWallet.SetVoteChoice(vspHost, pubKey, avm.agenda.AgendaID, radioBtnGrp, "", password)
		if err != nil {
			if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
				avm.spendingPassword.SetError("Invalid password")
			} else {
				avm.Toast.NotifyError(err.Error())
			}
			return
		}
		avm.Dismiss()
		avm.Toast.Notify("Vote updated successfully, refreshing agendas!")

		avm.Dismiss()
		go avm.consensusPage.FetchAgendas()
	}()
}

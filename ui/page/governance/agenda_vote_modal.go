package governance

import (
	"sort"

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

	votableTickets []*dcrlibwallet.Transaction // tickets that have not been spent by a vote or revocation (unspent) and that have not expired (unexpired).

	agenda           *dcrlibwallet.Agenda
	isVoting         bool
	modalUpdateCount int // this keeps track of the number of times the modal has been updated.

	onPreferenceUpdated func()

	walletSelector    *WalletSelector
	ticketSelector    *ticketSelector
	spendingPassword  decredmaterial.Editor
	materialLoader    material.LoaderStyle
	voteChoices       []string
	initialValue      string
	optionsRadioGroup *widget.Enum
	voteBtn           decredmaterial.Button
	cancelBtn         decredmaterial.Button
}

func newAgendaVoteModal(l *load.Load, agenda *dcrlibwallet.Agenda, onPreferenceUpdated func()) *agendaVoteModal {
	avm := &agendaVoteModal{
		Load:                l,
		modal:               *l.Theme.ModalFloatTitle(),
		agenda:              agenda,
		onPreferenceUpdated: onPreferenceUpdated,
		materialLoader:      material.Loader(material.NewTheme(gofont.Collection())),
		optionsRadioGroup:   new(widget.Enum),
		spendingPassword:    l.Theme.EditorPassword(new(widget.Editor), "Spending password"),
		voteBtn:             l.Theme.Button("Update Preference"),
		cancelBtn:           l.Theme.OutlineButton("Cancel"),
	}

	avm.voteBtn.Background = l.Theme.Color.Gray3
	avm.voteBtn.Color = l.Theme.Color.Surface

	avm.walletSelector = NewWalletSelector(l).
		Title("Select wallet").
		WalletSelected(func(w *dcrlibwallet.Wallet) {
			avm.modalUpdateCount = 0 // modal just opened.

			avm.FetchUnspentUnexpiredTickets(w.ID)
			avm.modalUpdateCount++

			// update agenda options prefrence to that of the selected wallet
			consensusItems := components.LoadAgendas(avm.Load, w, false)
			for _, consensusItem := range consensusItems {
				if consensusItem.Agenda.AgendaID == agenda.AgendaID {
					voteChoices := make([]string, len(consensusItem.Agenda.Choices))
					for i := range consensusItem.Agenda.Choices {
						voteChoices[i] = consensusItem.Agenda.Choices[i].Id
					}

					avm.voteChoices = voteChoices
					avm.initialValue = consensusItem.Agenda.VotingPreference
					avm.optionsRadioGroup.Value = avm.initialValue
				}
			}
		}).
		WalletValidator(func(w *dcrlibwallet.Wallet) bool {
			return !w.IsWatchingOnlyWallet()
		})

	return avm
}

func (avm *agendaVoteModal) FetchUnspentUnexpiredTickets(walletID int) {
	go func() {
		wallet := avm.WL.MultiWallet.WalletWithID(walletID)
		tickets, err := wallet.UnspentUnexpiredTickets()
		if err != nil {
			avm.Toast.NotifyError(err.Error())
			return
		}

		// sort by newest first
		sort.Slice(tickets[:], func(i, j int) bool {
			var timeStampI = tickets[i].Timestamp
			var timeStampJ = tickets[j].Timestamp
			return timeStampI > timeStampJ
		})
		avm.votableTickets = make([]*dcrlibwallet.Transaction, len(tickets))
		for i := range tickets {
			avm.votableTickets[i] = &tickets[i]
		}
		avm.RefreshWindow()
	}()
}

func (avm *agendaVoteModal) ModalID() string {
	return ModalInputVote
}

func (avm *agendaVoteModal) OnResume() {
	avm.walletSelector.SelectFirstValidWallet()

	avm.initialValue = avm.agenda.VotingPreference
	avm.optionsRadioGroup.Value = avm.initialValue
}

func (avm *agendaVoteModal) OnDismiss() {}

func (avm *agendaVoteModal) Show() {
	avm.ShowModal(avm)
}

func (avm *agendaVoteModal) Dismiss() {
	avm.DismissModal(avm)
}

func (avm *agendaVoteModal) Handle() {
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

	if len(avm.votableTickets) != 0 {
		if avm.modalUpdateCount == 1 { // modal window has been updated once.
			avm.modalUpdateCount++
			avm.ticketSelector = newTicketSelector(avm.Load, avm.votableTickets).Title("Select a ticket")
		}
	}

	validToVote := avm.optionsRadioGroup.Value != "" && avm.optionsRadioGroup.Value != avm.initialValue && avm.spendingPassword.Editor.Text() != ""
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
			t := avm.Theme.H6("Update Voting Preference")
			t.Font.Weight = text.SemiBold
			return t.Layout(gtx)
		},
		avm.Theme.Body1("Select one of the options below to vote").Layout,
		avm.walletSelector.Layout,
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, avm.layoutItems()...)
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(avm.spendingPassword.Layout),
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
	for _, voteChoice := range avm.voteChoices {
		radioBtn := avm.Load.Theme.RadioButton(avm.optionsRadioGroup, voteChoice, voteChoice, avm.Load.Theme.Color.DeepBlue, avm.Load.Theme.Color.Primary)
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

		choiceID := avm.optionsRadioGroup.Value
		err := avm.walletSelector.selectedWallet.SetVoteChoice(avm.agenda.AgendaID, choiceID, "", password)
		if err != nil {
			if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
				avm.spendingPassword.SetError("Invalid password")
			} else {
				avm.Toast.NotifyError(err.Error())
			}
			return
		}
		avm.Toast.Notify("Vote preference updated successfully")

		avm.Dismiss()
		avm.onPreferenceUpdated()
	}()
}

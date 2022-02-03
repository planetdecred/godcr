package staking

import (
	"context"
	"fmt"
	"strconv"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const ticketBuyerModalID = "staking_modal"

type ticketBuyerModal struct {
	*load.Load

	settingsSaved func()
	onCancel      func()

	modal           decredmaterial.Modal
	cancel          decredmaterial.Button
	saveSettingsBtn decredmaterial.Button

	balToMaintainEditor decredmaterial.Editor

	accountSelector *components.AccountSelector
	vspSelector     *vspSelector
}

func newTicketBuyerModal(l *load.Load) *ticketBuyerModal {
	tb := &ticketBuyerModal{
		Load: l,

		cancel:          l.Theme.OutlineButton("Cancel"),
		saveSettingsBtn: l.Theme.Button("Save"),
		modal:           *l.Theme.ModalFloatTitle(),
		vspSelector:     newVSPSelector(l).title("Select a vsp"),
	}

	tb.balToMaintainEditor = l.Theme.Editor(new(widget.Editor), "Balance to maintain (DCR)")
	tb.balToMaintainEditor.Editor.SingleLine = true

	tb.saveSettingsBtn.SetEnabled(false)

	return tb
}

func (tb *ticketBuyerModal) OnSettingsSaved(settingsSaved func()) *ticketBuyerModal {
	tb.settingsSaved = settingsSaved
	return tb
}

func (tb *ticketBuyerModal) OnCancel(cancel func()) *ticketBuyerModal {
	tb.onCancel = cancel
	return tb
}

func (tb *ticketBuyerModal) selectFirstWalletAccount() {
	err := tb.accountSelector.SelectFirstWalletValidAccount()
	if err != nil {
		tb.Toast.NotifyError(err.Error())
	}
}

func (tb *ticketBuyerModal) OnResume() {
	tb.initializeAccountSelector()
	if len(tb.WL.MultiWallet.KnownVSPs()) == 0 {
		// TODO: Does this modal need this list?
		go tb.WL.MultiWallet.ReloadVSPList(context.TODO())
	}

	// loop through all available wallets and select the one with ticket buyer config.
	// if non, set the selected wallet to the first.
	// temporary work around for only one wallet.
	// TODO: extend functionality to allow for multiwallet config
	for _, wal := range tb.WL.SortedWalletList() {
		if wal.TicketBuyerConfigIsSet() {
			tbConfig := wal.AutoTicketsBuyerConfig()
			acct, err := wal.GetAccount(tbConfig.PurchaseAccount)
			if err != nil {
				tb.Toast.NotifyError(err.Error())
			}

			tb.vspSelector.selectVSP(tbConfig.VspHost)
			tb.balToMaintainEditor.Editor.SetText(strconv.FormatFloat(dcrlibwallet.AmountCoin(tbConfig.BalanceToMaintain), 'f', 0, 64))
			tb.accountSelector.SetSelectedAccount(acct)
			break
		}
	}

	if tb.accountSelector.SelectedAccount() == nil {
		err := tb.accountSelector.SelectFirstWalletValidAccount()
		if err != nil {
			tb.Toast.NotifyError(err.Error())
		}
	}
}

func (tb *ticketBuyerModal) Layout(gtx layout.Context) layout.Dimensions {
	l := []layout.Widget{
		func(gtx C) D {
			t := tb.Theme.H6("Auto ticket purchase")
			t.Font.Weight = text.SemiBold
			return t.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding8,
						Bottom: values.MarginPadding16,
					}.Layout(gtx, tb.accountSelector.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return tb.balToMaintainEditor.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding16,
						Bottom: values.MarginPadding16,
					}.Layout(gtx, func(gtx C) D {
						return tb.vspSelector.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Right: values.MarginPadding4,
						}.Layout(gtx, tb.cancel.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						return tb.saveSettingsBtn.Layout(gtx)
					}),
				)
			})
		},
	}

	return tb.modal.Layout(gtx, l)
}

func (tb *ticketBuyerModal) canSave() bool {
	if tb.vspSelector.selectedVSP == nil {
		return false
	}

	if tb.balToMaintainEditor.Editor.Text() == "" {
		return false
	}

	return true
}

func (tb *ticketBuyerModal) ModalID() string {
	return ticketBuyerModalID
}

func (tb *ticketBuyerModal) Show() {
	tb.ShowModal(tb)
}

func (tb *ticketBuyerModal) Dismiss() {
	tb.DismissModal(tb)
}

func (tb *ticketBuyerModal) initializeAccountSelector() {
	tb.accountSelector = components.NewAccountSelector(tb.Load).
		Title("Purchasing account").
		AccountSelected(func(selectedAccount *dcrlibwallet.Account) {}).
		AccountValidator(func(account *dcrlibwallet.Account) bool {
			wal := tb.WL.MultiWallet.WalletWithID(account.WalletID)

			// Imported and watch only wallet accounts are invalid for sending
			accountIsValid := account.Number != dcrlibwallet.ImportedAccountNumber && !wal.IsWatchingOnlyWallet()

			if wal.ReadBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, false) {
				// privacy is enabled for selected wallet
				accountIsValid = account.Number == wal.MixedAccountNumber()
			}
			return accountIsValid
		})
}

func (tb *ticketBuyerModal) OnDismiss() {}

func (tb *ticketBuyerModal) Handle() {
	tb.saveSettingsBtn.SetEnabled(tb.canSave())

	if tb.cancel.Clicked() || tb.modal.BackdropClicked(true) {
		tb.onCancel()
		tb.Dismiss()
	}

	if tb.saveSettingsBtn.Clicked() {
		vspHost := tb.vspSelector.selectedVSP.Host
		amount, err := strconv.ParseFloat(tb.balToMaintainEditor.Editor.Text(), 64)
		if err != nil {
			tb.Toast.NotifyError(err.Error())
			return
		}

		balToMaintain := dcrlibwallet.AmountAtom(amount)
		account := tb.accountSelector.SelectedAccount()
		wal := tb.WL.MultiWallet.WalletWithID(account.WalletID)
		if wal == nil {
			tb.Toast.NotifyError(fmt.Sprintf("wallet with ID: %v does not exist", wal.ID))
			return
		}

		// clear all wallets config when a new wallet is selected.
		// temporary work around for only one wallet.
		// TODO: extend functionality to allow for multiwallet config
		for _, wal := range tb.WL.SortedWalletList() {
			tb.WL.MultiWallet.ClearTicketBuyerConfig(wal.ID)
		}

		wal.SetAutoTicketsBuyerConfig(vspHost, account.Number, balToMaintain)
		tb.settingsSaved()
		tb.Dismiss()
	}
}

package privacy

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/values"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type sharedModalConf struct {
	*load.Load
	wallet   *dcrlibwallet.Wallet
	checkBox decredmaterial.CheckBoxStyle
}

func showInfoModal(conf *sharedModalConf, title, body, btnText string, isError, alignCenter bool) {
	icon := decredmaterial.NewIcon(decredmaterial.MustIcon(widget.NewIcon(icons.AlertError)))
	icon.Color = conf.Theme.Color.DeepBlue
	if !isError {
		icon = decredmaterial.NewIcon(conf.Theme.Icons.ActionCheckCircle)
		icon.Color = conf.Theme.Color.Success
	}

	info := modal.NewInfoModal(conf.Load).
		Icon(icon).
		Title(title).
		Body(body).
		PositiveButton(btnText, func(isChecked bool) {})

	if alignCenter {
		align := layout.Center
		info.SetContentAlignment(align, align)
	}

	conf.ShowModal(info)
}

func showModalSetupMixerInfo(conf *sharedModalConf) {
	info := modal.NewInfoModal(conf.Load).
		Title("Set up mixer by creating two needed accounts").
		SetupWithTemplate(modal.SetupMixerInfoTemplate).
		CheckBox(conf.checkBox, false).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton("Begin setup", func(movefundsChecked bool) {
			showModalSetupMixerAcct(conf, movefundsChecked)
		})
	conf.ShowModal(info)
}

func showModalSetupMixerAcct(conf *sharedModalConf, movefundsChecked bool) {
	accounts, _ := conf.wallet.GetAccountsRaw()
	txt := "There are existing accounts named mixed or unmixed. Please change the name to something else for now. You can change them back after the setup."
	for _, acct := range accounts.Acc {
		if acct.Name == "mixed" || acct.Name == "unmixed" {
			alert := decredmaterial.NewIcon(decredmaterial.MustIcon(widget.NewIcon(icons.AlertError)))
			alert.Color = conf.Theme.Color.DeepBlue
			info := modal.NewInfoModal(conf.Load).
				Icon(alert).
				Title("Account name is taken").
				Body(txt).
				PositiveButton("Go back & rename", func(movefundsChecked bool) {
					conf.PopFragment()
				})
			conf.ShowModal(info)
			return
		}
	}

	modal.NewPasswordModal(conf.Load.Theme, nil).
		Title("Confirm to create needed accounts").
		NegativeButton("Cancel", func() {}).
		PositiveButton("Confirm", func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := conf.wallet.CreateMixerAccounts("mixed", "unmixed", password)
				if err != nil {
					pm.SetError(err.Error())
					pm.SetLoading(false)
					return
				}
				conf.WL.MultiWallet.SetBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, true)

				if movefundsChecked {
					err := moveFundsFromDefaultToUnmixed(conf, password)
					if err != nil {
						log.Error(err)
						txt := fmt.Sprintf("Error moving funds: %s.\n%s", err.Error(), "Auto funds transfer has been skipped. Move funds to unmixed account manually from the send page.")
						showInfoModal(conf, "Move funds to unmixed account", txt, "Got it", true, false)
					}
				}

				pm.Dismiss()

				conf.ChangeFragment(NewAccountMixerPage(conf.Load, conf.wallet))
			}()

			return false
		}).Show()
}

// moveFundsFromDefaultToUnmixed moves funds from the default wallet account to the
// newly created unmixed account
func moveFundsFromDefaultToUnmixed(conf *sharedModalConf, password string) error {
	acc, err := conf.wallet.GetAccountsRaw()
	if err != nil {
		return err
	}

	// get the first account in the wallet as this is the default
	sourceAccount := acc.Acc[0]
	destinationAccount := conf.wallet.UnmixedAccountNumber()

	destinationAddress, err := conf.wallet.CurrentAddress(destinationAccount)
	if err != nil {
		return err
	}

	unsignedTx, err := conf.WL.MultiWallet.NewUnsignedTx(sourceAccount.WalletID, sourceAccount.Number)
	if err != nil {
		return err
	}

	// get tx fees
	feeAndSize, err := unsignedTx.EstimateFeeAndSize()
	if err != nil {
		return err
	}

	// calculate max amount to be sent
	amountAtom := sourceAccount.Balance.Spendable - feeAndSize.Fee.AtomValue
	err = unsignedTx.AddSendDestination(destinationAddress, amountAtom, true)
	if err != nil {
		return err
	}

	// send fund
	_, err = unsignedTx.Broadcast([]byte(password))
	if err != nil {
		return err
	}

	showInfoModal(conf, "Transaction sent!", "", "Got it", false, true)

	return err
}

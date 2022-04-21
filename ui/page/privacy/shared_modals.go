package privacy

import (
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
	wallet *dcrlibwallet.Wallet
}

func showModalSetupMixerInfo(conf *sharedModalConf) {
	info := modal.NewInfoModal(conf.Load).
		Title("Set up mixer by creating two needed accounts").
		SetupWithTemplate(modal.SetupMixerInfoTemplate).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton("Begin setup", func() {
			showModalSetupMixerAcct(conf)
		})
	conf.ShowModal(info)
}

func showModalSetupMixerAcct(conf *sharedModalConf) {
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
				PositiveButton("Go back & rename", func() {
					conf.PopFragment()
				})
			conf.ShowModal(info)
			return
		}
	}

	modal.NewPasswordModal(conf.Load).
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
				pm.Dismiss()

				conf.ChangeFragment(NewAccountMixerPage(conf.Load, conf.wallet))
			}()

			return false
		}).Show()
}

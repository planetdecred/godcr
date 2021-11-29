package dexclient

import (
	"fmt"
	"strconv"

	"decred.org/dcrdex/client/asset/btc"
	"decred.org/dcrdex/client/asset/dcr"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

type dexCreateWalletWidget struct {
	*load.Load
	sourceAccountSelector *components.AccountSelector
	backButton            decredmaterial.IconButton
	createNewWallet       decredmaterial.Button
	walletPassword        decredmaterial.Editor
	appPassword           decredmaterial.Editor
	walletInfoWdg         *walletInfoWidget
	isSending             bool
	walletCreated         func()
}

type walletInfoWidget struct {
	image    *decredmaterial.Image
	coinName string
	coinID   uint32
}

func newDexCreateWalletWidget(l *load.Load) *dexCreateWalletWidget {
	dcw := &dexCreateWalletWidget{
		Load:            l,
		walletPassword:  l.Theme.EditorPassword(new(widget.Editor), "Wallet Password"),
		appPassword:     l.Theme.EditorPassword(new(widget.Editor), "App Password"),
		createNewWallet: l.Theme.Button("Add"),
	}

	dcw.createNewWallet.TextSize = values.TextSize12
	dcw.createNewWallet.Background = l.Theme.Color.Primary
	dcw.appPassword.Editor.SingleLine = true
	dcw.appPassword.Editor.SetText("")

	dcw.sourceAccountSelector = components.NewAccountSelector(dcw.Load).
		Title("Select DCR account to use with DEX").
		AccountSelected(func(selectedAccount *dcrlibwallet.Account) {}).
		AccountValidator(func(account *dcrlibwallet.Account) bool {
			// Filter out imported account and mixed.
			wal := dcw.WL.MultiWallet.WalletWithID(account.WalletID)
			if account.Number == load.MaxInt32 ||
				account.Number == wal.MixedAccountNumber() {
				return false
			}
			return true
		})
	err := dcw.sourceAccountSelector.SelectFirstWalletValidAccount()
	if err != nil {
		dcw.Toast.NotifyError(err.Error())
	}

	dcw.backButton, _ = components.SubpageHeaderButtons(l)

	return dcw
}

func (dcw *dexCreateWalletWidget) handle() {
	if dcw.createNewWallet.Button.Clicked() {
		if dcw.appPassword.Editor.Text() == "" || dcw.isSending {
			return
		}

		dcw.isSending = true
		go func() {
			defer func() {
				dcw.isSending = false
			}()

			coinID := dcw.walletInfoWdg.coinID
			coinName := dcw.walletInfoWdg.coinName
			if dcw.Dexc().HasWallet(int32(coinID)) {
				dcw.Toast.NotifyError(fmt.Sprintf("already connected a %s wallet", coinName))
				return
			}

			settings := make(map[string]string)
			var walletType string
			appPass := []byte(dcw.appPassword.Editor.Text())
			walletPass := []byte(dcw.walletPassword.Editor.Text())

			switch coinID {
			case dcr.BipID:
				selectedAccount := dcw.sourceAccountSelector.SelectedAccount()
				settings[dcrlibwallet.DexDcrWalletIDConfigKey] = strconv.Itoa(selectedAccount.WalletID)
				settings["account"] = selectedAccount.Name
				settings["password"] = dcw.walletPassword.Editor.Text()
				walletType = dcrlibwallet.CustomDexDcrWalletType
			case btc.BipID:
				walletType = "SPV" // decred.org/dcrdex/client/asset/btc.walletTypeSPV
				walletPass = nil   // Core doesn't accept wallet passwords for dex-managed spv wallets.
			}

			err := dcw.Dexc().AddWallet(coinID, walletType, settings, appPass, walletPass)
			if err != nil {
				dcw.Toast.NotifyError(err.Error())
				return
			}

			dcw.walletCreated()
		}()
	}
}

func (dcw *dexCreateWalletWidget) layout(gtx layout.Context) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return dcw.Load.Theme.Label(values.TextSize20, "Add a").Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Left: values.MarginPadding8, Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
							ic := dcw.walletInfoWdg.image
							ic.Scale = 0.2
							return dcw.walletInfoWdg.image.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return dcw.Load.Theme.Label(values.TextSize20, fmt.Sprintf("%s Wallet", dcw.walletInfoWdg.coinName)).Layout(gtx)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return dcw.Load.Theme.Label(values.TextSize14, "Your wallet is required to pay registration fees.").Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			if dcw.walletInfoWdg.coinID == dcr.BipID {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
							return dcw.sourceAccountSelector.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
							return dcw.walletPassword.Layout(gtx)
						})
					}),
				)
			}
			return D{}
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
				return dcw.appPassword.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
				if dcw.isSending {
					dcw.createNewWallet.Background = dcw.Theme.Color.Hint
				} else {
					dcw.createNewWallet.Background = dcw.Theme.Color.Primary
				}
				return dcw.createNewWallet.Layout(gtx)
			})
		}),
	)
}

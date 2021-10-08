package dexclient

import (
	"fmt"
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/dexc"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const dexUnlockWalletModalID = "dex_unlock_wallet_modal"

// walletInfoWidget the data will be show up in the unlock or create wallet modal.
type walletInfoWidget struct {
	image    *decredmaterial.Image
	coin     string
	coinName string
	coinID   uint32
}

type unlockWalletModal struct {
	*load.Load
	modal            *decredmaterial.Modal
	unlockWallet     decredmaterial.Button
	appPassword      decredmaterial.Editor
	isSending        bool
	walletInfoWidget *walletInfoWidget
	unlocked         func([]byte)
}

func newUnlockWalletModal(l *load.Load) *unlockWalletModal {
	md := &unlockWalletModal{
		Load:         l,
		modal:        l.Theme.ModalFloatTitle(),
		appPassword:  l.Theme.EditorPassword(new(widget.Editor), "Password"),
		unlockWallet: l.Theme.Button("Unlock"),
		walletInfoWidget: &walletInfoWidget{
			image:    coinImageBySymbol(&l.Icons, dexc.DefaultAsset),
			coin:     dexc.DefaultAsset,
			coinName: "Decred",
			coinID:   dexc.DefaultAssetID,
		},
	}

	md.unlockWallet.TextSize = values.TextSize12
	md.unlockWallet.Background = l.Theme.Color.Primary
	md.appPassword.Editor.SingleLine = true

	return md
}

func (md *unlockWalletModal) ModalID() string {
	return dexUnlockWalletModalID
}

func (md *unlockWalletModal) Show() {
	md.ShowModal(md)
}

func (md *unlockWalletModal) Dismiss() {
	md.DismissModal(md)
}

func (md *unlockWalletModal) OnDismiss() {
	md.appPassword.Editor.SetText("")
}

func (md *unlockWalletModal) OnResume() {
	md.appPassword.Editor.Focus()
}

func (md *unlockWalletModal) Handle() {
	if md.unlockWallet.Button.Clicked() {
		if strings.Trim(md.appPassword.Editor.Text(), " ") == "" || md.isSending {
			return
		}

		md.isSending = true
		go func() {
			status := md.DL.WalletState(md.walletInfoWidget.coinID)
			if status == nil {
				md.isSending = false
				md.Toast.NotifyError(fmt.Sprintf("No wallet for %d", md.walletInfoWidget.coinID))
				return
			}

			err := md.DL.OpenWallet(md.walletInfoWidget.coinID, []byte(md.appPassword.Editor.Text()))
			if err != nil {
				md.isSending = false
				md.Toast.NotifyError(err.Error())
				return
			}
			md.isSending = false
			md.unlocked([]byte(md.appPassword.Editor.Text()))
			md.Dismiss()
		}()
	}
}

func (md *unlockWalletModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return md.Load.Theme.Label(values.TextSize20, "Unlock").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding8, Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						ic := md.walletInfoWidget.image
						ic.Scale = 0.2
						return md.walletInfoWidget.image.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return md.Load.Theme.Label(values.TextSize20, fmt.Sprintf("%s Wallet", md.walletInfoWidget.coinName)).Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return md.Load.Theme.Label(values.TextSize14, `App Password
Your app password is always required when performing sensitive wallet operations.`).Layout(gtx)
		},
		func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
				return md.appPassword.Layout(gtx)
			})
		},
		func(gtx C) D {
			return md.unlockWallet.Layout(gtx)
		},
	}

	return md.modal.Layout(gtx, w, 900)
}

// coinImageBySymbol returns image widget for supported asset coins.
func coinImageBySymbol(icons *load.Icons, coinName string) *decredmaterial.Image {
	m := map[string]*decredmaterial.Image{
		"btc": icons.BTC,
		"dcr": icons.DCR,
		"bch": icons.BCH,
		"ltc": icons.LTC,
	}
	coin, ok := m[strings.ToLower(coinName)]

	if !ok {
		return nil
	}

	return coin
}

package dexclient

import (
	"fmt"
	"strings"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/dexc"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const dexCreateWalletModalID = "dex_create_wallet_modal"

type createWalletModal struct {
	*load.Load
	modal            *decredmaterial.Modal
	createNewWallet  decredmaterial.Button
	walletPassword   decredmaterial.Editor
	appPassword      decredmaterial.Editor
	accountName      decredmaterial.Editor
	walletInfoWidget *walletInfoWidget
	isSending        bool
	walletCreated    func()
}

func newCreateWalletModal(l *load.Load, wallInfo *walletInfoWidget) *createWalletModal {
	md := &createWalletModal{
		Load:             l,
		modal:            l.Theme.ModalFloatTitle(),
		accountName:      l.Theme.Editor(new(widget.Editor), "Account Name"),
		walletPassword:   l.Theme.EditorPassword(new(widget.Editor), "Wallet Password"),
		appPassword:      l.Theme.EditorPassword(new(widget.Editor), "App Password"),
		createNewWallet:  l.Theme.Button(new(widget.Clickable), "Add"),
		walletInfoWidget: wallInfo,
	}

	md.createNewWallet.TextSize = values.TextSize12
	md.createNewWallet.Background = l.Theme.Color.Primary
	md.appPassword.Editor.SingleLine = true
	md.appPassword.Editor.SetText("")

	return md
}

func (md *createWalletModal) ModalID() string {
	return dexCreateWalletModalID
}

func (md *createWalletModal) Show() {
	md.ShowModal(md)
}

func (md *createWalletModal) Dismiss() {
	md.DismissModal(md)
}

func (md *createWalletModal) OnDismiss() {
	md.accountName.Editor.SetText("")
}

func (md *createWalletModal) OnResume() {
	md.accountName.Editor.Focus()
}

func (md *createWalletModal) Handle() {
	if md.createNewWallet.Button.Clicked() {
		if strings.Trim(md.appPassword.Editor.Text(), " ") == "" || md.isSending {
			return
		}

		md.isSending = true
		go func() {
			coinID := md.walletInfoWidget.coinID
			config, err := md.DL.AutoWalletConfig(coinID)

			if err != nil {
				md.isSending = false
				md.Toast.NotifyError(err.Error())
				return
			}

			for assetID, supportedAsset := range md.DL.User().Assets {
				if assetID == coinID {
					for _, cfgOpt := range supportedAsset.Info.ConfigOpts {
						if key := cfgOpt.Key; key == "fallbackfee" ||
							key == "feeratelimit" ||
							key == "redeemconftarget" ||
							key == "rpcbind" ||
							key == "rpcport" ||
							key == "txsplit" {
							config[key] = fmt.Sprintf("%v", cfgOpt.DefaultValue)
						}
					}

					break
				}
			}

			// Bitcoin
			config["walletname"] = md.accountName.Editor.Text()
			config["rpcport"] = "18332"

			// Decred
			config["account"] = md.accountName.Editor.Text()
			config["password"] = md.walletPassword.Editor.Text()

			form := &dexc.NewWalletForm{
				AssetID: coinID,
				Config:  config,
				Pass:    []byte(md.walletPassword.Editor.Text()),
				AppPW:   []byte(md.appPassword.Editor.Text()),
			}

			has := md.DL.WalletState(form.AssetID) != nil
			if has {
				md.Toast.NotifyError(fmt.Sprintf("already have a wallet for %d", form.AssetID))
				return
			}

			// Wallet does not exist yet. Try to create it.
			err = md.DL.CreateWallet(form.AppPW, form.Pass, &core.WalletForm{
				AssetID: form.AssetID,
				Config:  form.Config,
			})
			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
			}

			md.walletCreated()
			md.Dismiss()
		}()
	}
}

func (md *createWalletModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return md.Load.Theme.Label(values.TextSize20, "Add a").Layout(gtx)
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
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if md.walletInfoWidget.coinID == dexc.DefaultAssetID {
						return md.Load.Theme.Label(values.TextSize14, "Your Decred wallet is required to pay registration fees.").Layout(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return md.accountName.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return md.walletPassword.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return md.appPassword.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return md.createNewWallet.Layout(gtx)
		},
	}

	return md.modal.Layout(gtx, w, 900)
}

package dexclient

import (
	"fmt"
	"strconv"

	"decred.org/dcrdex/client/asset/btc"
	"decred.org/dcrdex/client/asset/dcr"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const dexCreateWalletModalID = "dex_create_wallet_modal"

type createWalletModal struct {
	*load.Load
	sourceAccountSelector *components.AccountSelector
	modal                 *decredmaterial.Modal
	submit                decredmaterial.Button
	cancel                decredmaterial.Button
	walletPassword        decredmaterial.Editor
	appPassword           decredmaterial.Editor
	walletInfoWidget      *walletInfoWidget
	materialLoader        material.LoaderStyle
	isSending             bool
	walletCreated         func()
}

type walletInfoWidget struct {
	image    *decredmaterial.Image
	coinName string
	coinID   uint32
}

func newCreateWalletModal(l *load.Load, wallInfo *walletInfoWidget, walletCreated func()) *createWalletModal {
	md := &createWalletModal{
		Load:             l,
		modal:            l.Theme.ModalFloatTitle(),
		walletPassword:   l.Theme.EditorPassword(new(widget.Editor), "Wallet Password"),
		appPassword:      l.Theme.EditorPassword(new(widget.Editor), "App Password"),
		submit:           l.Theme.Button("Add"),
		cancel:           l.Theme.OutlineButton("Cancel"),
		materialLoader:   material.Loader(material.NewTheme(gofont.Collection())),
		walletInfoWidget: wallInfo,
		walletCreated:    walletCreated,
	}

	md.appPassword.Editor.SingleLine = true
	md.appPassword.Editor.SetText("")

	md.sourceAccountSelector = components.NewAccountSelector(md.Load).
		Title("Select DCR account to use with DEX").
		AccountSelected(func(selectedAccount *dcrlibwallet.Account) {}).
		AccountValidator(func(account *dcrlibwallet.Account) bool {
			// Filter out imported account and mixed.
			wal := md.WL.MultiWallet.WalletWithID(account.WalletID)
			if account.Number == load.MaxInt32 ||
				account.Number == wal.MixedAccountNumber() {
				return false
			}
			return true
		})
	err := md.sourceAccountSelector.SelectFirstWalletValidAccount()
	if err != nil {
		md.Toast.NotifyError(err.Error())
	}

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
}

func (md *createWalletModal) OnResume() {
}

func (md *createWalletModal) Handle() {
	if md.cancel.Button.Clicked() && !md.isSending {
		md.Dismiss()
	}

	if md.submit.Button.Clicked() {
		if md.appPassword.Editor.Text() == "" || md.isSending {
			return
		}

		md.isSending = true
		go func() {
			defer func() {
				md.isSending = false
			}()

			coinID := md.walletInfoWidget.coinID
			coinName := md.walletInfoWidget.coinName
			if md.Dexc().HasWallet(int32(coinID)) {
				md.Toast.NotifyError(fmt.Sprintf("already connected a %s wallet", coinName))
				return
			}

			settings := make(map[string]string)
			var walletType string
			appPass := []byte(md.appPassword.Editor.Text())
			walletPass := []byte(md.walletPassword.Editor.Text())

			switch coinID {
			case dcr.BipID:
				selectedAccount := md.sourceAccountSelector.SelectedAccount()
				settings[dcrlibwallet.DexDcrWalletIDConfigKey] = strconv.Itoa(selectedAccount.WalletID)
				settings["account"] = selectedAccount.Name
				settings["password"] = md.walletPassword.Editor.Text()
				walletType = dcrlibwallet.CustomDexDcrWalletType
			case btc.BipID:
				walletType = "SPV" // decred.org/dcrdex/client/asset/btc.walletTypeSPV
				walletPass = nil   // Core doesn't accept wallet passwords for dex-managed spv wallets.
			}

			err := md.Dexc().AddWallet(coinID, walletType, settings, appPass, walletPass)
			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
			}

			md.Dismiss()
			md.walletCreated()
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
					return md.Load.Theme.Label(values.TextSize14, "Your wallet is required to pay registration fees.").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if md.walletInfoWidget.coinID == dcr.BipID {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
									return md.sourceAccountSelector.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
									return md.walletPassword.Layout(gtx)
								})
							}),
						)
					}
					return D{}
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return md.appPassword.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if md.isSending {
							return D{}
						}
						return layout.Inset{
							Right:  values.MarginPadding4,
							Bottom: values.MarginPadding15,
						}.Layout(gtx, md.cancel.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						if md.isSending {
							return layout.Inset{
								Top:    values.MarginPadding10,
								Bottom: values.MarginPadding15,
							}.Layout(gtx, md.materialLoader.Layout)
						}
						return md.submit.Layout(gtx)
					}),
				)
			})
		},
	}

	return md.modal.Layout(gtx, w)
}

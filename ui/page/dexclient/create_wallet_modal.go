package dexclient

import (
	"context"
	"fmt"
	"strconv"

	"decred.org/dcrdex/client/asset/btc"
	"decred.org/dcrdex/client/asset/dcr"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

type createWalletModal struct {
	*load.Load
	*decredmaterial.Modal

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	sourceAccountSelector *components.AccountSelector
	submit                decredmaterial.Button
	cancel                decredmaterial.Button
	walletPassword        decredmaterial.Editor
	appPassword           decredmaterial.Editor
	walletInfoWidget      *walletInfoWidget
	materialLoader        material.LoaderStyle
	isSending             bool
	dexClientPassword     string
	isRegisterAction      bool
	walletCreated         func()
	cancelClicked         func()
}

type walletInfoWidget struct {
	image    *decredmaterial.Image
	coinName string
	coinID   uint32
}

func newCreateWalletModal(l *load.Load, wallInfo *walletInfoWidget, appPass string) *createWalletModal {
	md := &createWalletModal{
		Load:             l,
		Modal:            l.Theme.ModalFloatTitle("dex_create_wallet_modal"),
		walletPassword:   l.Theme.EditorPassword(new(widget.Editor), "Wallet Password"),
		appPassword:      l.Theme.EditorPassword(new(widget.Editor), "App Password"),
		submit:           l.Theme.Button("Add"),
		cancel:           l.Theme.OutlineButton("Cancel"),
		materialLoader:   material.Loader(l.Theme.Base),
		walletInfoWidget: wallInfo,
		walletCreated:    walletCreated,
	}
	md.submitBtn.SetEnabled(false)
	md.sourceAccountSelector = components.NewAccountSelector(md.Load, nil).
		Title(strSelectAccountForDex).
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

	return md
}

func (md *createWalletModal) OnDismiss() {
	md.ctxCancel()
}

func (md *createWalletModal) OnResume() {
	md.ctx, md.ctxCancel = context.WithCancel(context.TODO())
	md.sourceAccountSelector.ListenForTxNotifications(md.ctx, md.ParentWindow())

	err := md.sourceAccountSelector.SelectFirstWalletValidAccount(nil)
	if err != nil {
		md.Toast.NotifyError(err.Error())
	}
}

func (md *createWalletModal) SetRegisterAction(registerAction bool) *createWalletModal {
	md.isRegisterAction = registerAction
	return md
}

func (md *createWalletModal) CancelClicked(clicked func()) *createWalletModal {
	md.cancelClicked = clicked
	return md
}

func (md *createWalletModal) WalletCreated(walletCreated func()) *createWalletModal {
	md.walletCreated = walletCreated
	return md
}

func (md *createWalletModal) validateInputs(isRequiredWalletPassword bool) (bool, string, string) {
	appPass := md.dexClientPassword
	if appPass == "" {
		appPass = md.appPassword.Editor.Text()
	}

	if appPass == "" {
		md.submitBtn.SetEnabled(false)
		return false, "", ""
	}

	wallPassword := md.walletPassword.Editor.Text()
	if isRequiredWalletPassword && wallPassword == "" {
		md.submitBtn.SetEnabled(false)
		return false, "", ""
	}

	md.submitBtn.SetEnabled(true)
	return true, appPass, wallPassword
}

func (md *createWalletModal) Handle() {
	canSubmit, appPass, walletPass := md.validateInputs(md.walletInfoWidget.coinID == dcr.BipID)

	if isWalletPasswordSubmit, _ := decredmaterial.HandleEditorEvents(md.walletPassword.Editor); isWalletPasswordSubmit {
		if canSubmit {
			md.doCreateWallet([]byte(appPass), []byte(walletPass))
		} else if md.dexClientPassword == "" {
			md.appPassword.Editor.Focus()
		}
	}

	isSubmit, _ := decredmaterial.HandleEditorEvents(md.appPassword.Editor)
	if canSubmit && (md.submitBtn.Clicked() || isSubmit) {
		md.doCreateWallet([]byte(appPass), []byte(walletPass))
	}

	if md.cancelBtn.Clicked() && !md.isSending {
		md.Dismiss()
		md.cancelClicked()
	}
}

func (md *createWalletModal) doCreateWallet(appPass, walletPass []byte) {
	if md.isSending {
		return
	}

	md.isSending = true
	md.modal.SetDisabled(true)
	go func() {
		defer func() {
			md.isSending = false
			md.modal.SetDisabled(false)
		}()

		coinID := md.walletInfoWidget.coinID
		coinName := md.walletInfoWidget.coinName
		if md.Dexc().HasWallet(int32(coinID)) {
			md.Toast.NotifyError(fmt.Sprintf(nStrAlreadyConnectWallet, coinName))
			return
		}

		md.isSending = true
		md.Modal.SetDisabled(true)

		go func() {
			defer func() {
				md.isSending = false
				md.Modal.SetDisabled(false)
			}()

		md.Dismiss()
		md.walletCreated()
	}()
}

func (md *createWalletModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return md.Load.Theme.Label(values.TextSize20, strAddA).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding8, Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						ic := md.walletInfoWidget.image
						ic.Scale = 0.2
						return md.walletInfoWidget.image.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return md.Load.Theme.Label(values.TextSize20, fmt.Sprintf(nStrNameWallet, md.walletInfoWidget.coinName)).Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if !md.isRegisterAction {
						return D{}
					}
					return md.Load.Theme.Label(values.TextSize14, strRequireWalletPayFee).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if md.walletInfoWidget.coinID == dcr.BipID {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
									return md.sourceAccountSelector.Layout(md.ParentWindow(), gtx)
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
					if md.dexClientPassword != "" {
						return D{}
					}
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
						}.Layout(gtx, md.cancelBtn.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						if md.isSending {
							return layout.Inset{
								Top:    values.MarginPadding10,
								Bottom: values.MarginPadding15,
							}.Layout(gtx, md.materialLoader.Layout)
						}
						return md.submitBtn.Layout(gtx)
					}),
				)
			})
		},
	}

	return md.Modal.Layout(gtx, w)
}

package dexclient

import (
	"context"
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

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	sourceAccountSelector *components.AccountSelector
	modal                 *decredmaterial.Modal
	submitBtn             decredmaterial.Button
	cancelBtn             decredmaterial.Button
	walletPassword        decredmaterial.Editor
	appPassword           decredmaterial.Editor
	walletInfoWidget      *walletInfoWidget
	materialLoader        material.LoaderStyle
	isSending             bool
	dexClientPassword     string
	isRegisterAction      bool
	walletCreated         func(md *createWalletModal)
}

type walletInfoWidget struct {
	image    *decredmaterial.Image
	coinName string
	coinID   uint32
}

func newCreateWalletModal(l *load.Load, wallInfo *walletInfoWidget, appPass string, walletCreated func(md *createWalletModal)) *createWalletModal {
	md := &createWalletModal{
		Load:              l,
		modal:             l.Theme.ModalFloatTitle(),
		walletPassword:    l.Theme.EditorPassword(&widget.Editor{Submit: true}, strWalletPassword),
		appPassword:       l.Theme.EditorPassword(&widget.Editor{Submit: true}, strAppPassword),
		submitBtn:         l.Theme.Button(strSubmit),
		cancelBtn:         l.Theme.OutlineButton(values.String(values.StrCancel)),
		materialLoader:    material.Loader(material.NewTheme(gofont.Collection())),
		walletInfoWidget:  wallInfo,
		walletCreated:     walletCreated,
		dexClientPassword: appPass,
	}
	md.submitBtn.SetEnabled(false)
	md.sourceAccountSelector = components.NewAccountSelector(md.Load, nil).
		Title(strSellectAccountForDex).
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
	md.ctxCancel()
}

func (md *createWalletModal) OnResume() {
	md.ctx, md.ctxCancel = context.WithCancel(context.TODO())
	md.sourceAccountSelector.ListenForTxNotifications(md.ctx)

	err := md.sourceAccountSelector.SelectFirstWalletValidAccount(nil)
	if err != nil {
		md.Toast.NotifyError(err.Error())
	}
}

func (md *createWalletModal) SetRegisterAction(registerAction bool) *createWalletModal {
	md.isRegisterAction = registerAction
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
	isRequiredWalletPassword := md.walletInfoWidget.coinID == dcr.BipID
	canSubmit, appPass, walletPass := md.validateInputs(isRequiredWalletPassword)

	if isWalletPasswordSubmit, _ := decredmaterial.HandleEditorEvents(md.walletPassword.Editor); isWalletPasswordSubmit {
		if md.dexClientPassword != "" && canSubmit {
			if isRequiredWalletPassword {
				md.doCreateWallet([]byte(appPass), []byte(walletPass))
			} else {
				md.doCreateWallet([]byte(appPass), nil)
			}
		} else {
			md.appPassword.Editor.Focus()
		}
	}

	isSubmit, _ := decredmaterial.HandleEditorEvents(md.appPassword.Editor)
	if canSubmit && (md.submitBtn.Button.Clicked() || isSubmit) {
		if isRequiredWalletPassword {
			md.doCreateWallet([]byte(appPass), []byte(walletPass))
		} else {
			md.doCreateWallet([]byte(appPass), nil)
		}
	}

	if md.cancelBtn.Button.Clicked() && !md.isSending {
		md.Dismiss()
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

		settings := make(map[string]string)
		var walletType string
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
		md.walletCreated(md)
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

	return md.modal.Layout(gtx, w)
}

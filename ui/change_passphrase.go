package ui

import (
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

const PageWalletPassphrase = "walletPassphrase"

type walletPassphrasePage struct {
	container                     layout.List
	newPassW, confPassW, oldPassW widget.Editor
	newPass, confPass, oldPass    decredmaterial.Editor
	passwordBar                   *decredmaterial.ProgressBar
	backButtonW, savePasswordW    widget.Button
	backButton                    decredmaterial.IconButton
	savePassword                  decredmaterial.Button
	errorLabel                    decredmaterial.Label
}

func (win *Window) WalletPassphrasePage(common pageCommon) layout.Widget {
	page := &walletPassphrasePage{
		container: layout.List{
			Axis: layout.Vertical,
		},
		oldPass:      common.theme.Editor("Enter old password"),
		newPass:      common.theme.Editor("Enter new password"),
		confPass:     common.theme.Editor("Confirm new password"),
		savePassword: common.theme.Button("Change"),
		errorLabel:   common.theme.Caption(""),
		backButton:   common.theme.PlainIconButton(common.icons.navigationArrowBack),
	}
	page.oldPass.IsRequired = true
	page.oldPassW.SingleLine = true
	page.newPass.IsRequired = true
	page.newPassW.SingleLine = true
	page.confPass.IsRequired = true
	page.confPassW.SingleLine = true
	page.savePassword.TextSize = unit.Dp(11)
	page.passwordBar = common.theme.ProgressBar(0)
	page.errorLabel.Color = common.theme.Color.Danger

	page.backButton.Color = common.theme.Color.Hint
	page.backButton.Size = unit.Dp(32)

	return func() {
		page.Layout(common)
		page.handle(common)
	}
}

// Layout lays out the widgets for the change wallet passphrase page.
func (page *walletPassphrasePage) Layout(common pageCommon) {
	gtx := common.gtx
	wdgs := []func(){
		func() {
			layout.W.Layout(common.gtx, func() {
				page.backButton.Layout(common.gtx, &page.backButtonW)
			})
			layout.Inset{Left: unit.Dp(44)}.Layout(common.gtx, func() {
				common.theme.H5("Change Wallet Password").Layout(gtx)
			})
		},
		func() {
			layout.Flex{}.Layout(gtx,
				layout.Rigid(func() {
					layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func() {
						common.theme.Body1("Are about changing the passphrase for").Layout(gtx)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Left: unit.Dp(5)}.Layout(gtx, func() {
						txt := common.theme.H5(common.info.Wallets[*common.selectedWallet].Name)
						txt.Color = common.theme.Color.Danger
						txt.Layout(gtx)
					})
				}),
			)
		},
		func() {
			page.errorLabel.Layout(gtx)
		},
		func() {
			page.oldPass.Layout(gtx, &page.oldPassW)
		},
		func() {
			page.passwordStrength(common)
		},
		func() {
			page.newPass.Layout(gtx, &page.newPassW)
		},
		func() {
			page.confPass.Layout(gtx, &page.confPassW)
		},
		func() {
			layout.Inset{Top: unit.Dp(20)}.Layout(gtx, func() {
				page.savePassword.Layout(gtx, &page.savePasswordW)
			})
		},
	}

	common.Layout(gtx, func() {
		layout.UniformInset(unit.Dp(20)).Layout(gtx, func() {
			page.container.Layout(gtx, len(wdgs), func(i int) {
				layout.UniformInset(unit.Dp(3)).Layout(gtx, wdgs[i])
			})
		})
	})
}

func (page *walletPassphrasePage) passwordStrength(common pageCommon) {
	layout.Inset{Top: unit.Dp(10)}.Layout(common.gtx, func() {
		common.gtx.Constraints.Height.Max = 20
		page.passwordBar.Layout(common.gtx)
	})
}

// Handle handles all widget inputs on the main wallets page.
func (page *walletPassphrasePage) handle(common pageCommon) {
	gtx := common.gtx

	page.handlerEditorEvents(common, &page.confPassW)
	page.handlerEditorEvents(common, &page.newPassW)
	page.handlerEditorEvents(common, &page.oldPassW)

	if page.savePasswordW.Clicked(gtx) && page.validateInputs(common) {
		op := page.oldPassW.Text()
		np := page.newPassW.Text()

		err := common.wallet.ChangeWalletPassphrase(common.info.Wallets[*common.selectedWallet].ID, op, np)
		if err != nil {
			log.Debug("Error changing wallet password " + err.Error())
			if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
				page.errorLabel.Text = "Passphrase is incorrect"
			} else {
				page.errorLabel.Text = err.Error()
			}
			return
		}

		page.errorLabel.Text = "Password changed successfully"
		page.errorLabel.Color = common.theme.Color.Success
		page.savePassword.Text = "Changed"
		page.savePassword.Background = common.theme.Color.Success
		page.resetFields()
	}

	if strings.Trim(page.newPassW.Text(), " ") != "" {
		strength := dcrlibwallet.ShannonEntropy(page.newPassW.Text()) / 4.0
		switch {
		case strength > 0.6:
			page.passwordBar.ProgressColor = common.theme.Color.Success
		case strength > 0.3 && strength <= 0.6:
			page.passwordBar.ProgressColor = color.RGBA{248, 231, 27, 255}
		default:
			page.passwordBar.ProgressColor = common.theme.Color.Danger
		}

		page.passwordBar.Progress = strength * 100
	}

	if page.backButtonW.Clicked(common.gtx) {
		page.clear(common)
		*common.page = PageWallet
	}
}

func (page *walletPassphrasePage) validateInputs(common pageCommon) bool {
	page.errorLabel.Text = ""
	page.oldPass.ErrorLabel.Text = ""
	page.newPass.ErrorLabel.Text = ""
	page.confPass.ErrorLabel.Text = ""
	page.savePassword.Background = common.theme.Color.Hint

	if page.oldPassW.Text() == "" {
		page.oldPass.ErrorLabel.Text = "Please wallet old password"
		return false
	}
	if page.newPassW.Text() == "" {
		page.newPass.ErrorLabel.Text = "Please wallet new password"
		return false
	}
	if page.confPassW.Text() == "" {
		page.confPass.ErrorLabel.Text = "Please wallet new password again"
		return false
	}

	if page.confPassW.Text() != page.newPassW.Text() {
		page.errorLabel.Text = "New wallet passwords does no match. Try again."
		return false
	}

	page.savePassword.Background = common.theme.Color.Primary
	return true
}

func (page *walletPassphrasePage) handlerEditorEvents(common pageCommon, w *widget.Editor) {
	for _, evt := range w.Events(common.gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			page.validateInputs(common)
			return
		}
	}
}

func (page *walletPassphrasePage) resetFields() {
	page.oldPassW.SetText("")
	page.newPassW.SetText("")
	page.confPassW.SetText("")
	page.passwordBar.Progress = 0
}

func (page *walletPassphrasePage) clear(common pageCommon) {
	page.savePassword.Background = common.theme.Color.Hint
	page.oldPassW.SetText("")
	page.newPassW.SetText("")
	page.confPassW.SetText("")
	page.errorLabel.Text = ""
	page.savePassword.Text = "Change"
	page.passwordBar.Progress = 0
}

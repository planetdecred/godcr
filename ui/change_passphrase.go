package ui

import (
	"image/color"
	"strings"

	"github.com/raedahgroup/godcr/ui/values"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

const PageWalletPassphrase = "walletPassphrase"

type walletPassphrasePage struct {
	container                  layout.List
	newPass, confPass, oldPass decredmaterial.Editor
	passwordBar                decredmaterial.ProgressBarStyle
	backButton                 decredmaterial.IconButton
	savePassword               decredmaterial.Button
	errorLabel                 decredmaterial.Label
}

func (win *Window) WalletPassphrasePage(common pageCommon) layout.Widget {
	page := &walletPassphrasePage{
		container: layout.List{
			Axis: layout.Vertical,
		},
		oldPass:      common.theme.Editor(new(widget.Editor), "Enter old password"),
		newPass:      common.theme.Editor(new(widget.Editor), "Enter new password"),
		confPass:     common.theme.Editor(new(widget.Editor), "Confirm new password"),
		savePassword: common.theme.Button(new(widget.Clickable), "Change"),
		errorLabel:   common.theme.Caption(""),
		backButton:   common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack),
	}
	page.oldPass.IsRequired = true
	page.oldPass.Editor.SingleLine = true
	page.newPass.IsRequired = true
	page.newPass.Editor.SingleLine = true
	page.confPass.IsRequired = true
	page.confPass.Editor.SingleLine = true
	page.savePassword.TextSize = values.TextSize12
	page.savePassword.Background = common.theme.Color.Hint
	page.passwordBar = common.theme.ProgressBar(0)
	page.errorLabel.Color = common.theme.Color.Danger

	page.backButton.Color = common.theme.Color.Hint
	page.backButton.Size = values.MarginPadding30

	return func(gtx C) D {
		page.handle(common)
		return page.Layout(common)
	}
}

// Layout lays out the widgets for the change wallet passphrase page.
func (page *walletPassphrasePage) Layout(common pageCommon) layout.Dimensions {
	gtx := common.gtx
	wdgs := []func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.W.Layout(common.gtx, func(gtx C) D {
						return page.backButton.Layout(common.gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding45}.Layout(common.gtx, func(gtx C) D {
						return common.theme.H5("Change Wallet Password").Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return common.theme.Body1("Are about changing the passphrase for").Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						txt := common.theme.H5(common.info.Wallets[*common.selectedWallet].Name)
						txt.Color = common.theme.Color.Danger
						return txt.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return page.errorLabel.Layout(gtx)
		},
		func(gtx C) D {
			return page.oldPass.Layout(gtx)
		},
		func(gtx C) D {
			// page.passwordStrength(common)
			return layout.Dimensions{}
		},
		func(gtx C) D {
			return page.newPass.Layout(gtx)
		},
		func(gtx C) D {
			return page.confPass.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
				return page.savePassword.Layout(gtx)
			})
		},
	}

	return common.Layout(gtx, func(gtx C) D {
		return page.container.Layout(gtx, len(wdgs), func(gtx C, i int) D {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, wdgs[i])
		})
	})
}

func (page *walletPassphrasePage) passwordStrength(common pageCommon) {
	layout.Inset{Top: values.MarginPadding10}.Layout(common.gtx, func(gtx C) D {
		common.gtx.Constraints.Max.Y = 20
		return page.passwordBar.Layout(common.gtx)
	})
}

// Handle handles all widget inputs on the main wallets page.
func (page *walletPassphrasePage) handle(common pageCommon) {
	page.handlerEditorEvents(common, page.confPass.Editor)
	page.handlerEditorEvents(common, page.newPass.Editor)
	page.handlerEditorEvents(common, page.oldPass.Editor)

	if page.savePassword.Button.Clicked() && page.validateInputs(common) {
		op := page.oldPass.Editor.Text()
		np := page.newPass.Editor.Text()

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

	if strings.Trim(page.newPass.Editor.Text(), " ") != "" {
		strength := dcrlibwallet.ShannonEntropy(page.newPass.Editor.Text()) / 4.0
		switch {
		case strength > 0.6:
			page.passwordBar.Color = common.theme.Color.Success
		case strength > 0.3 && strength <= 0.6:
			page.passwordBar.Color = color.RGBA{248, 231, 27, 255}
		default:
			page.passwordBar.Color = common.theme.Color.Danger
		}

		page.passwordBar.Progress = int(strength * 100)
	}

	if page.backButton.Button.Clicked() {
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

	if page.oldPass.Editor.Text() == "" {
		page.oldPass.ErrorLabel.Text = "Please wallet old password"
		return false
	}
	if page.newPass.Editor.Text() == "" {
		page.newPass.ErrorLabel.Text = "Please wallet new password"
		return false
	}
	if page.confPass.Editor.Text() == "" {
		page.confPass.ErrorLabel.Text = "Please wallet new password again"
		return false
	}

	if page.confPass.Editor.Text() != page.newPass.Editor.Text() {
		page.errorLabel.Text = "New wallet passwords does no match. Try again."
		return false
	}

	page.savePassword.Background = common.theme.Color.Primary
	return true
}

func (page *walletPassphrasePage) handlerEditorEvents(common pageCommon, w *widget.Editor) {
	for _, evt := range w.Events() {
		switch evt.(type) {
		case widget.ChangeEvent:
			page.validateInputs(common)
			return
		}
	}
}

func (page *walletPassphrasePage) resetFields() {
	page.oldPass.Editor.SetText("")
	page.newPass.Editor.SetText("")
	page.confPass.Editor.SetText("")
	page.passwordBar.Progress = 0
}

func (page *walletPassphrasePage) clear(common pageCommon) {
	page.savePassword.Background = common.theme.Color.Hint
	page.oldPass.Editor.SetText("")
	page.newPass.Editor.SetText("")
	page.confPass.Editor.SetText("")
	page.errorLabel.Text = ""
	page.savePassword.Text = "Change"
	page.passwordBar.Progress = 0
}

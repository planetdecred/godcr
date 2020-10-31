package ui

import (
	"image/color"
	"strings"

	"github.com/planetdecred/godcr/ui/values"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const PageWalletPassphrase = "WalletPassphrase"

type walletPassphrasePage struct {
	container                  layout.List
	newPass, confPass, oldPass decredmaterial.Editor
	passwordBar                decredmaterial.ProgressBarStyle
	backButton                 decredmaterial.IconButton
	savePassword               decredmaterial.Button
	errorLabel                 decredmaterial.Label
}

func (win *Window) WalletPassphrasePage(common pageCommon) layout.Widget {
	pg := &walletPassphrasePage{
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
	pg.oldPass.IsRequired = true
	pg.oldPass.Editor.SingleLine = true
	pg.newPass.IsRequired = true
	pg.newPass.Editor.SingleLine = true
	pg.confPass.IsRequired = true
	pg.confPass.Editor.SingleLine = true
	pg.savePassword.TextSize = values.TextSize12
	pg.savePassword.Background = common.theme.Color.Hint
	pg.passwordBar = common.theme.ProgressBar(0)
	pg.errorLabel.Color = common.theme.Color.Danger

	pg.backButton.Color = common.theme.Color.Hint
	pg.backButton.Size = values.MarginPadding30
	pg.backButton.Inset = layout.UniformInset(values.MarginPadding0)

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

// Layout lays out the widgets for the change wallet passphrase pg.
func (pg *walletPassphrasePage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	wdgs := []func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.W.Layout(gtx, func(gtx C) D {
						return pg.backButton.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
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
			return pg.errorLabel.Layout(gtx)
		},
		func(gtx C) D {
			return pg.oldPass.Layout(gtx)
		},
		func(gtx C) D {
			// pg.passwordStrength(common)
			return layout.Dimensions{}
		},
		func(gtx C) D {
			return pg.newPass.Layout(gtx)
		},
		func(gtx C) D {
			return pg.confPass.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
				return pg.savePassword.Layout(gtx)
			})
		},
	}

	return common.Layout(gtx, func(gtx C) D {
		return pg.container.Layout(gtx, len(wdgs), func(gtx C, i int) D {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, wdgs[i])
		})
	})
}

func (pg *walletPassphrasePage) passwordStrength(gtx layout.Context, common pageCommon) {
	layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Max.Y = 20
		return pg.passwordBar.Layout(gtx)
	})
}

// Handle handles all widget inputs on the main wallets pg.
func (pg *walletPassphrasePage) handle(common pageCommon) {
	pg.handlerEditorEvents(common, pg.confPass.Editor)
	pg.handlerEditorEvents(common, pg.newPass.Editor)
	pg.handlerEditorEvents(common, pg.oldPass.Editor)

	if pg.savePassword.Button.Clicked() && pg.validateInputs(common) {
		op := pg.oldPass.Editor.Text()
		np := pg.newPass.Editor.Text()

		err := common.wallet.ChangeWalletPassphrase(common.info.Wallets[*common.selectedWallet].ID, op, np)
		if err != nil {
			log.Debug("Error changing wallet password " + err.Error())
			if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
				pg.errorLabel.Text = "Passphrase is incorrect"
			} else {
				pg.errorLabel.Text = err.Error()
			}
			return
		}

		pg.errorLabel.Text = "Password changed successfully"
		pg.errorLabel.Color = common.theme.Color.Success
		pg.savePassword.Text = "Changed"
		pg.savePassword.Background = common.theme.Color.Success
		pg.resetFields()
	}

	if strings.Trim(pg.newPass.Editor.Text(), " ") != "" {
		strength := dcrlibwallet.ShannonEntropy(pg.newPass.Editor.Text()) / 4.0
		switch {
		case strength > 0.6:
			pg.passwordBar.Color = common.theme.Color.Success
		case strength > 0.3 && strength <= 0.6:
			pg.passwordBar.Color = color.RGBA{248, 231, 27, 255}
		default:
			pg.passwordBar.Color = common.theme.Color.Danger
		}

		pg.passwordBar.Progress = int(strength * 100)
	}

	if pg.backButton.Button.Clicked() {
		pg.clear(common)
		*common.page = PageWallet
	}
}

func (pg *walletPassphrasePage) validateInputs(common pageCommon) bool {
	pg.errorLabel.Text = ""
	pg.oldPass.ErrorLabel.Text = ""
	pg.newPass.ErrorLabel.Text = ""
	pg.confPass.ErrorLabel.Text = ""
	pg.savePassword.Background = common.theme.Color.Hint

	if pg.oldPass.Editor.Text() == "" {
		pg.oldPass.ErrorLabel.Text = "Please wallet old password"
		return false
	}
	if pg.newPass.Editor.Text() == "" {
		pg.newPass.ErrorLabel.Text = "Please wallet new password"
		return false
	}
	if pg.confPass.Editor.Text() == "" {
		pg.confPass.ErrorLabel.Text = "Please wallet new password again"
		return false
	}

	if pg.confPass.Editor.Text() != pg.newPass.Editor.Text() {
		pg.errorLabel.Text = "New wallet passwords does no match. Try again."
		return false
	}

	pg.savePassword.Background = common.theme.Color.Primary
	return true
}

func (pg *walletPassphrasePage) handlerEditorEvents(common pageCommon, w *widget.Editor) {
	for _, evt := range w.Events() {
		switch evt.(type) {
		case widget.ChangeEvent:
			pg.validateInputs(common)
			return
		}
	}
}

func (pg *walletPassphrasePage) resetFields() {
	pg.oldPass.Editor.SetText("")
	pg.newPass.Editor.SetText("")
	pg.confPass.Editor.SetText("")
	pg.passwordBar.Progress = 0
}

func (pg *walletPassphrasePage) clear(common pageCommon) {
	pg.savePassword.Background = common.theme.Color.Hint
	pg.oldPass.Editor.SetText("")
	pg.newPass.Editor.SetText("")
	pg.confPass.Editor.SetText("")
	pg.errorLabel.Text = ""
	pg.savePassword.Text = "Change"
	pg.passwordBar.Progress = 0
}

package ui

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr/ui/values"

	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
)

const PageWalletAccounts = "walletAccounts"

type walletAccountPage struct {
	walletID  int
	wallet    *wallet.Wallet
	container layout.List
	// accountNameW         widget.Editor
	accountName decredmaterial.Editor
	// backButtonW, createW widget.Clickable
	backButton    decredmaterial.IconButton
	create        decredmaterial.Button
	errorLabel    decredmaterial.Label
	passwordModal *decredmaterial.Password
	isPassword    bool
	state         bool
	errChan       chan error
}

func (win *Window) WalletAccountPage(common pageCommon) layout.Widget {
	page := &walletAccountPage{
		container: layout.List{
			Axis: layout.Vertical,
		},
		wallet:        common.wallet,
		passwordModal: common.theme.Password(),
		accountName:   common.theme.Editor(new(widget.Editor), "Enter Account Name"),
		create:        common.theme.Button(new(widget.Clickable), "Create"),
		errorLabel:    common.theme.Caption(""),
		backButton:    common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack),
		state:         common.states.creating,
		errChan:       common.errorChannels[PageWalletAccounts],
	}
	page.accountName.IsRequired = true
	page.accountName.Editor.SingleLine = true
	page.create.TextSize = values.TextSize12
	page.errorLabel.Color = common.theme.Color.Danger

	page.create.Background = common.theme.Color.Hint
	page.backButton.Color = common.theme.Color.Hint
	page.backButton.Size = values.MarginPadding30

	return func(gtx C) D {
		page.handle(common)
		return page.Layout(common)
	}
}

// Layout lays out the widgets for the change wallet passphrase page.
func (page *walletAccountPage) Layout(common pageCommon) layout.Dimensions {
	select {
	case err := <-page.errChan:
		page.state = false
		if err.Error() == "invalid_passphrase" {
			page.errorLabel.Text = "Wallet passphrase is incorrect"
		} else {
			page.errorLabel.Text = err.Error()
		}
	default:
	}

	return layout.Flex{}.Layout(common.gtx,
		layout.Flexed(1, func(gtx C) D {
			if page.state {
				return page.processing(common)(gtx)
			} else {
				if page.isPassword {
					page.walletID = common.info.Wallets[*common.selectedWallet].ID
					page.passwordModal.Layout(gtx, page.confirm, page.cancel)
				}
				return page.createAccount(common)
			}
		}),
	)
}

func (page *walletAccountPage) createAccount(common pageCommon) layout.Dimensions {
	gtx := common.gtx
	wdgs := []func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.W.Layout(gtx, func(gtx C) D {
						return page.backButton.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding45}.Layout(gtx, func(gtx C) D {
						return common.theme.H5("Create Wallet Acount").Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			dims := layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return common.theme.Body1("Are about changing an Account in").Layout(gtx)
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
			return dims
		},
		func(gtx C) D {
			return page.errorLabel.Layout(gtx)
		},
		func(gtx C) D {
			return page.accountName.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
				return page.create.Layout(gtx)
			})
		},
	}

	return common.Layout(gtx, func(gtx C) D {
		return page.container.Layout(gtx, len(wdgs), func(gtx C, i int) D {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, wdgs[i])
		})
	})
}

// Handle handles all widget inputs on the main wallets page.
func (page *walletAccountPage) handle(common pageCommon) {
	page.handleEditorEvents(common, page.accountName.Editor)

	if page.create.Button.Clicked() && page.validateInputs(common) {
		page.isPassword = true
	}

	if page.backButton.Button.Clicked() {
		page.clear(common)
		*common.page = PageWallet
	}
}

func (page *walletAccountPage) validateInputs(common pageCommon) bool {
	page.errorLabel.Text = ""
	page.accountName.ErrorLabel.Text = ""
	page.create.Background = common.theme.Color.Hint

	if page.accountName.Editor.Text() == "" {
		page.accountName.ErrorLabel.Text = "Please wallet old password"
		return false
	}

	page.create.Background = common.theme.Color.Primary
	return true
}

func (page *walletAccountPage) handleEditorEvents(common pageCommon, w *widget.Editor) {
	for _, evt := range w.Events() {
		switch evt.(type) {
		case widget.ChangeEvent:
			page.validateInputs(common)
			return
		}
	}
}

func (page *walletAccountPage) confirm(passphrase []byte) {
	page.isPassword = false

	page.wallet.AddAccount(page.walletID, page.accountName.Editor.Text(), passphrase, page.errChan)
	page.state = true
}

func (page *walletAccountPage) cancel() {
	page.isPassword = false
}

func (page *walletAccountPage) processing(common pageCommon) layout.Widget {
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					message := common.theme.H3("Creating Account...")
					message.Alignment = text.Middle
					return message.Layout(gtx)
				})
			}),
		)
	}
}

func (page *walletAccountPage) clear(common pageCommon) {
	page.create.Background = common.theme.Color.Hint
	page.accountName.Editor.SetText("")
	page.errorLabel.Text = ""
}

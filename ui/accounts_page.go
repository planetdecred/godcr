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
	walletID             int
	wallet               *wallet.Wallet
	container            layout.List
	accountNameW         widget.Editor
	accountName          decredmaterial.Editor
	backButtonW, createW widget.Button
	backButton           decredmaterial.IconButton
	create               decredmaterial.Button
	errorLabel           decredmaterial.Label
	passwordModal        *decredmaterial.Password
	isPassword           bool
	state                bool
	errChan              chan error
}

func (win *Window) WalletAccountPage(common pageCommon) layout.Widget {
	page := &walletAccountPage{
		container: layout.List{
			Axis: layout.Vertical,
		},
		wallet:        common.wallet,
		passwordModal: common.theme.Password(),
		accountName:   common.theme.Editor("Enter Account Name"),
		create:        common.theme.Button("Create"),
		errorLabel:    common.theme.Caption(""),
		backButton:    common.theme.PlainIconButton(common.icons.navigationArrowBack),
		state:         common.states.creating,
		errChan:       common.errorChannels[PageWalletAccounts],
	}
	page.accountName.IsRequired = true
	page.accountNameW.SingleLine = true
	page.create.TextSize = values.TextSize12
	page.errorLabel.Color = common.theme.Color.Danger

	page.create.Background = common.theme.Color.Hint
	page.backButton.Color = common.theme.Color.Hint
	page.backButton.Size = values.MarginPadding30

	return func() {
		page.Layout(common)
		page.handle(common)
	}
}

// Layout lays out the widgets for the change wallet passphrase page.
func (page *walletAccountPage) Layout(common pageCommon) {
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

	layout.Flex{}.Layout(common.gtx,
		layout.Flexed(1, func() {
			if page.state {
				page.processing(common)()
			} else {
				page.createAccount(common)
			}
		}),
	)
}

func (page *walletAccountPage) createAccount(common pageCommon) {
	gtx := common.gtx
	wdgs := []func(){
		func() {
			layout.W.Layout(common.gtx, func() {
				page.backButton.Layout(common.gtx, &page.backButtonW)
			})
			layout.Inset{Left: values.MarginPadding45}.Layout(common.gtx, func() {
				common.theme.H5("Create Wallet Acount").Layout(gtx)
			})
		},
		func() {
			layout.Flex{}.Layout(gtx,
				layout.Rigid(func() {
					layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func() {
						common.theme.Body1("Are about changing an Account in").Layout(gtx)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func() {
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
			page.accountName.Layout(gtx, &page.accountNameW)
		},
		func() {
			layout.Inset{Top: values.MarginPadding20}.Layout(gtx, func() {
				page.create.Layout(gtx, &page.createW)
			})
		},
	}

	common.Layout(gtx, func() {
		page.container.Layout(gtx, len(wdgs), func(i int) {
			layout.UniformInset(values.MarginPadding5).Layout(gtx, wdgs[i])
		})
	})

	if page.isPassword {
		page.walletID = common.info.Wallets[*common.selectedWallet].ID
		page.passwordModal.Layout(gtx, page.confirm, page.cancel)
	}
}

// Handle handles all widget inputs on the main wallets page.
func (page *walletAccountPage) handle(common pageCommon) {
	gtx := common.gtx

	page.handleEditorEvents(common, &page.accountNameW)

	if page.createW.Clicked(gtx) && page.validateInputs(common) {
		page.isPassword = true
	}

	if page.backButtonW.Clicked(common.gtx) {
		page.clear(common)
		*common.page = PageWallet
	}
}

func (page *walletAccountPage) validateInputs(common pageCommon) bool {
	page.errorLabel.Text = ""
	page.accountName.ErrorLabel.Text = ""
	page.create.Background = common.theme.Color.Hint

	if page.accountNameW.Text() == "" {
		page.accountName.ErrorLabel.Text = "Please wallet old password"
		return false
	}

	page.create.Background = common.theme.Color.Primary
	return true
}

func (page *walletAccountPage) handleEditorEvents(common pageCommon, w *widget.Editor) {
	for _, evt := range w.Events(common.gtx) {
		switch evt.(type) {
		case widget.ChangeEvent:
			page.validateInputs(common)
			return
		}
	}
}

func (page *walletAccountPage) confirm(passphrase []byte) {
	page.isPassword = false

	page.wallet.AddAccount(page.walletID, page.accountNameW.Text(), passphrase, page.errChan)
	page.state = true
}

func (page *walletAccountPage) cancel() {
	page.isPassword = false
}

func (page *walletAccountPage) processing(common pageCommon) layout.Widget {
	gtx := common.gtx

	return func() {
		layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(1, func() {
				layout.Center.Layout(gtx, func() {
					message := common.theme.H3("Creating Account...")
					message.Alignment = text.Middle
					message.Layout(gtx)
				})
			}),
		)
	}
}

func (page *walletAccountPage) clear(common pageCommon) {
	page.create.Background = common.theme.Color.Hint
	page.accountNameW.SetText("")
	page.errorLabel.Text = ""
}

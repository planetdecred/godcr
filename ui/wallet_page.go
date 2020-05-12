package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

const PageWallet = "wallet"

const (
	subWalletMain = iota
	subWalletRename
	subWalletDelete
	subWalletChangePass
	subWalletBackup
	subWalletAddAcct
	subWalletVerify
	subWalletSign
)

type walletPage struct {

	subPage int
	sub     struct {
		main, delete, rename    widget.Button
		mainW, deleteW, renameW decredmaterial.Button
	}
	container, accountsList                        layout.List
	rename, delete, addAcct widget.Button
	signMessageDiag, verifyMessageDiag                        *widget.Button
	deleteW, gotoSignMessagePageBtn, gotoVerifyMessagePageBtn                                        decredmaterial.Button
	renameW, beginRenameW, cancelRenameW, addAcctW decredmaterial.IconButton
	editor, password                               widget.Editor
	editorW, passwordW                             decredmaterial.Editor
	errorLabel                                     decredmaterial.Label
}

func WalletPage(common pageCommon) layout.Widget {
	page := &walletPage{
		container: layout.List{
			Axis: layout.Vertical,
		},
		accountsList: layout.List{
			Axis: layout.Vertical,
		},

		beginRenameW:  common.theme.PlainIconButton(common.icons.contentCreate),
		cancelRenameW: common.theme.PlainIconButton(common.icons.contentClear),
		renameW:       common.theme.PlainIconButton(common.icons.navigationCheck),
		editorW:       common.theme.Editor("Enter wallet name"),
		passwordW:     common.theme.Editor("Enter Wallet Password"),
		addAcctW:      common.theme.IconButton(common.icons.contentAdd),
		deleteW:       common.theme.DangerButton("Confirm Delete Wallet"),
		gotoSignMessagePageBtn:   win.outputs.signMessageDiag,
		signMessageDiag:          &win.inputs.signMessageDiag,
		gotoVerifyMessagePageBtn: win.outputs.verifyMessDiag,
		addWalletW:               common.theme.Button("Add Wallet"),
		verifyMessageDiag:        &win.inputs.verifyMessDiag,
	}

	page.sub.mainW = common.theme.Button("Back")
	page.sub.deleteW = common.theme.DangerButton("Delete Wallet")
	page.sub.renameW = common.theme.Button("Rename")

	return func() {
		page.Layout(common)
		page.Handle(common)
	}
}

// Layout lays out the widgets for the main wallets page.
func (page *walletPage) Layout(common pageCommon) {
	if common.walletsTab.Changed() { // reset everything
		page.subPage = subWalletMain
		//page.renaming = false
	}

	if *common.err != nil {
		page.errorLabel.Text = (*common.err).Error()
		*common.err = nil
	}

	switch page.subPage {
	case subWalletMain:
		page.subMain(common)
	case subWalletRename:
		page.subRename(common)
	case subWalletDelete:
		page.subDelete(common)
	}

}

func (page *walletPage) subMain(common pageCommon) {
	gtx := common.gtx
	current := common.info.Wallets[*common.selectedWallet]
	wdgs := []func(){
		func() {
common.theme.H3(current.Name).Layout(common.gtx)
		},
		func() {
			common.theme.H5("Total Balance: " + current.Balance).Layout(gtx)
		},
		func() {
			horFlex.Layout(gtx,
				rigid(func() {
					common.theme.H5("Accounts").Layout(gtx)
				}),
				rigid(func() {
					layout.S.Layout(gtx, func() {
						page.addAcctW.Layout(gtx, &page.addAcct)
					})
				}),
			)
		},
		func() {
			page.accountsList.Layout(gtx, len(current.Accounts), func(i int) {
				acct := current.Accounts[i]
				a := func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							common.theme.Body1(acct.Name).Layout(gtx)
						}),
						layout.Rigid(func() {
							common.theme.Body1(acct.TotalBalance).Layout(gtx)
						}),
						layout.Rigid(func() {
							common.theme.Body1("Keys: " + acct.Keys.External + " external, " + acct.Keys.Internal + " internal, " + acct.Keys.Imported + " imported").Layout(gtx)
						}),
						layout.Rigid(func() {
							common.theme.Body1("HD Path: " + acct.HDPath).Layout(gtx)
						}),
					)
				}
				layout.UniformInset(unit.Dp(5)).Layout(gtx, a)
			})
		},
		func() {
			layout.Flex{}.Layout(gtx,
				layout.Rigid(func() {
					page.sub.deleteW.Layout(gtx, &page.sub.delete)
				}),
				layout.Rigid(func() {
					//page.deleteW.Layout(gtx, &page.delete)
				}),
				layout.Rigid(func() {
					layout.Inset{
						Left: unit.Dp(5),
					}.Layout(gtx, func() {
						page.gotoVerifyMessagePageBtn.Layout(gtx, page.verifyMessageDiag)
					})
				}),
				layout.Rigid(func() {
					inset := layout.Inset{
						Left: unit.Dp(5),
					}
					inset.Layout(gtx, func() {
						page.addWalletW.Layout(gtx, page.addWallet)
					})
				}),
			)
		},
	}
	common.LayoutWithWallets(gtx, func() {
		page.container.Layout(common.gtx, len(wdgs), func(i int) {
			wdgs[i]()
		})
	})
}

func (page *walletPage) subRename(common pageCommon) {
	gtx := common.gtx
	current := common.info.Wallets[*common.selectedWallet]
	common.Layout(gtx, func() {
		common.theme.H3(current.Name).Layout(common.gtx)
	})
}

func (page *walletPage) subDelete(common pageCommon) {
	gtx := common.gtx
	current := common.info.Wallets[*common.selectedWallet]
	list := layout.List{Axis: layout.Vertical}
	wdgs := []func(){
		func() {
			page.sub.mainW.Layout(gtx, &page.sub.main)
		},
		func() {
			common.theme.H3(current.Name).Layout(common.gtx)
		},
		func() {
			page.passwordW.Layout(gtx, &page.password)
		},
		func() {
			page.errorLabel.Layout(gtx)
		},
		func() {
			page.deleteW.Layout(gtx, &page.delete)
		},
	}
	common.Layout(gtx, func() {
		list.Layout(gtx, len(wdgs), func(i int) {
			wdgs[i]()
		})
	})
}

// Handle handles all widget inputs on the main wallets page.
func (page *walletPage) Handle(common pageCommon) {
	gtx := common.gtx
	current := common.info.Wallets[*common.selectedWallet]

	// Subs
	if page.sub.delete.Clicked(gtx) {
		page.subPage = subWalletDelete
		return
	}

	if page.sub.main.Clicked(gtx) {
		page.subPage = subWalletMain
		return
	}

	if page.sub.rename.Clicked(gtx) {
		page.subPage = subWalletRename
		return
	}

	if page.delete.Clicked(gtx) {
		common.wallet.DeleteWallet(current.ID, page.password.Text())
	}

	if page.addWallet.Clicked(gtx) {
		*common.page = PageCreateRestore
	}
}

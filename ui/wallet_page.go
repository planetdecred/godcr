package ui

// add all the sub pages to the wallet page.
//clean up the delete wallet sub page
import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
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
	current wallet.InfoShort
	wallet  *wallet.Wallet
	sub     struct {
		main, delete, rename, sign, verify widget.Button
		renameW, signW, verifyW, deleteW   decredmaterial.Button
		mainW                              decredmaterial.IconButton
	}
	signP                                          signMessagePage
	verifyPg                                       verifyMessagePage
	deletePg                                       deleteWalletPage
	container, accountsList                        layout.List
	rename, delete, addAcct                        widget.Button
	renameW, beginRenameW, cancelRenameW, addAcctW decredmaterial.IconButton
	editor, password                               widget.Editor
	editorW, passwordW                             decredmaterial.Editor
	signatureResult                                *wallet.Signature
}

func WalletPage(common pageCommon) layout.Widget {
	page := &walletPage{
		container: layout.List{
			Axis: layout.Vertical,
		},
		accountsList: layout.List{
			Axis: layout.Vertical,
		},

		wallet:        common.wallet,
		beginRenameW:  common.theme.PlainIconButton(common.icons.contentCreate),
		cancelRenameW: common.theme.PlainIconButton(common.icons.contentClear),
		renameW:       common.theme.PlainIconButton(common.icons.navigationCheck),
		editorW:       common.theme.Editor("Enter wallet name"),
		passwordW:     common.theme.Editor("Enter Wallet Password"),
		addAcctW:      common.theme.IconButton(common.icons.contentAdd),
	}
	page.sub.mainW = common.theme.IconButton(common.icons.navigationArrowBack)
	page.sub.mainW.Background = color.RGBA{}
	page.sub.mainW.Color = common.theme.Color.Text
	page.sub.mainW.Padding = unit.Dp(0)
	page.sub.mainW.Size = unit.Dp(30)
	page.sub.deleteW = common.theme.DangerButton("Delete Wallet")
	page.sub.signW = common.theme.Button("sign Message")
	page.sub.verifyW = common.theme.Button("Verify Message")
	page.sub.renameW = common.theme.Button("Rename")

	page.SignMessagePage(common)
	page.VerifyMessagePage(common)
	page.DeleteWalletPage(common)

	return func() {
		page.Layout(common)
		page.Handle(common)
		page.handleSign(common)
		page.handleVerify(common)
		page.handleDelete(common)
	}
}

// Layout lays out the widgets for the main wallets page.
func (page *walletPage) Layout(common pageCommon) {
	if common.walletsTab.Changed() { // reset everything
		page.subPage = subWalletMain
	}

	switch page.subPage {
	case subWalletMain:
		page.subMain(common)
	case subWalletRename:
		page.subRename(common)
	case subWalletDelete:
		page.subDelete(common)
	case subWalletSign:
		page.subSign(common)
	case subWalletVerify:
		page.subVerify(common)
	}
}

func (page *walletPage) subMain(common pageCommon) {
	gtx := common.gtx
	page.current = common.info.Wallets[*common.selectedWallet]
	wdgs := []func(){
		func() {
			horFlex.Layout(gtx,
				rigid(func() {
					common.theme.H1(page.current.Name).Layout(common.gtx)
				}),
				rigid(func() {
					layout.Center.Layout(gtx, func() {
						page.rename.Layout(gtx)
					})
				}),
			)
		},
		func() {
			common.theme.H5("Total Balance: " + page.current.Balance).Layout(gtx)
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
			page.accountsList.Layout(gtx, len(page.current.Accounts), func(i int) {
				acct := page.current.Accounts[i]
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
					page.sub.signW.Layout(gtx, &page.sub.sign)
				}),
				layout.Rigid(func() {
					page.sub.verifyW.Layout(gtx, &page.sub.verify)
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
	common.Layout(gtx, func() {
		common.theme.H3(page.current.Name).Layout(common.gtx)
	})
}

// Handle handles all widget inputs on the main wallets page.
func (page *walletPage) Handle(common pageCommon) {
	gtx := common.gtx

	// Subs
	if page.sub.main.Clicked(gtx) {
		page.deletePg.errorLabel.Text = ""
		page.subPage = subWalletMain
		return
	}

	if page.sub.rename.Clicked(gtx) {
		page.subPage = subWalletRename
		return
	}

	if page.sub.delete.Clicked(gtx) {
		page.subPage = subWalletDelete
		return
	}

	if page.sub.sign.Clicked(gtx) {
		page.subPage = subWalletSign
		return
	}

	if page.sub.verify.Clicked(gtx) {
		page.subPage = subWalletVerify
		return
	}
}

func (page *walletPage) returnBtn(common pageCommon) {
	layout.NW.Layout(common.gtx, func() {
		page.sub.mainW.Layout(common.gtx, &page.sub.main)
	})
}

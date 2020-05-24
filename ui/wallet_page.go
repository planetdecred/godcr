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
		main, delete, rename, sign, verify, addWallet, changePass widget.Button
		mainW, deleteW, signW, verifyW, addWalletW, renameW,
		changePassW decredmaterial.IconButton
	}
	signP                   signMessagePage
	verifyPg                verifyMessagePage
	deletePg                deleteWalletPage
	renamePg                renameWalletPage
	container, accountsList layout.List
	rename, delete, addAcct widget.Button
	line                    *decredmaterial.Line
	addAcctW                decredmaterial.IconButton
	editor, password        widget.Editor
	editorW, passwordW      decredmaterial.Editor
	signatureResult         *wallet.Signature
}

func WalletPage(common pageCommon) layout.Widget {
	page := &walletPage{
		container: layout.List{
			Axis: layout.Vertical,
		},
		accountsList: layout.List{
			Axis: layout.Vertical,
		},
		wallet:    common.wallet,
		passwordW: common.theme.Editor("Enter Wallet Password"),
		addAcctW:  common.theme.IconButton(common.icons.contentAdd),
		line:      common.theme.Line(),
	}
	page.line.Color = common.theme.Color.Gray
	page.sub.mainW = common.theme.IconButton(common.icons.navigationArrowBack)
	page.sub.mainW.Background = color.RGBA{}
	page.sub.mainW.Color = common.theme.Color.Text
	page.sub.mainW.Padding = unit.Dp(0)
	page.sub.mainW.Size = unit.Dp(30)
	page.sub.deleteW = common.theme.IconButton(common.icons.actionDelete)
	page.sub.signW = common.theme.IconButton(common.icons.communicationComment)
	page.sub.verifyW = common.theme.IconButton(common.icons.verifyAction)
	page.sub.addWalletW = common.theme.IconButton(common.icons.contentAdd)
	page.sub.renameW = common.theme.IconButton(common.icons.editorModeEdit)
	page.sub.changePassW = common.theme.IconButton(common.icons.actionLock)
	page.sub.deleteW.Background = common.theme.Color.Danger
	page.sub.deleteW.Size, page.sub.signW.Size, page.sub.verifyW.Size = unit.Dp(30), unit.Dp(30), unit.Dp(30)
	page.sub.renameW.Size, page.sub.addWalletW.Size, page.addAcctW.Size = unit.Dp(30), unit.Dp(30), unit.Dp(30)
	page.sub.changePassW.Size = unit.Dp(30)
	page.sub.deleteW.Padding, page.sub.signW.Padding, page.sub.verifyW.Padding = unit.Dp(5), unit.Dp(5), unit.Dp(5)
	page.sub.renameW.Padding, page.sub.addWalletW.Padding, page.addAcctW.Padding = unit.Dp(5), unit.Dp(5), unit.Dp(5)
	page.sub.changePassW.Padding = unit.Dp(5)

	page.SignMessagePage(common)
	page.VerifyMessagePage(common)
	page.DeleteWalletPage(common)
	page.RenameWalletPage(common)

	return func() {
		page.Layout(common)
		page.Handle(common)
		page.handleSign(common)
		page.handleVerify(common)
		page.handleDelete(common)
		page.handleRename(common)
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

	body := func() {
		layout.Stack{}.Layout(gtx,
			layout.Expanded(func() {
				layout.Inset{Top: unit.Dp(15)}.Layout(gtx, func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Flexed(0.88, func() {
							page.topRow(common)
						}),
						layout.Flexed(0.12, func() {
							page.bottomRow(common)
						}),
					)
				})
			}),
		)
	}

	common.LayoutWithWallets(gtx, func() {
		layout.UniformInset(unit.Dp(5)).Layout(gtx, body)
	})
}

func (page *walletPage) topRow(common pageCommon) {
	gtx := common.gtx
	wdgs := []func(){
		func() {
			horFlex.Layout(gtx,
				rigid(func() {
					common.theme.H4(page.current.Name).Layout(common.gtx)
				}),
				rigid(func() {
					layout.Center.Layout(gtx, func() {
						page.rename.Layout(gtx)
					})
				}),
			)
		},
		func() {
			common.theme.H6("Total Balance: " + page.current.Balance).Layout(gtx)
		},
		func() {
			horFlex.Layout(gtx,
				rigid(func() {
					common.theme.H6("Accounts").Layout(gtx)
				}),
				rigid(func() {
					layout.Inset{Left: unit.Dp(3)}.Layout(common.gtx, func() {
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
						layout.Rigid(func() {
							gtx.Constraints.Width.Min = gtx.Px(unit.Dp(350))
							gtx.Constraints.Width.Max = gtx.Constraints.Width.Min
							page.line.Layout(gtx)
						}),
					)
				}
				layout.UniformInset(unit.Dp(5)).Layout(gtx, a)
			})
		},
	}

	page.container.Layout(gtx, len(wdgs), func(i int) {
		layout.Inset{Left: unit.Dp(3)}.Layout(gtx, wdgs[i])
	})
}

func (page *walletPage) bottomRow(common pageCommon) {
	gtx := common.gtx
	layout.UniformInset(unit.Dp(5)).Layout(gtx, func() {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(page.newItem(&common, page.sub.addWalletW, &page.sub.addWallet, "Add wallet")),
			layout.Rigid(page.newItem(&common, page.sub.renameW, &page.sub.rename, "Rename wallet")),
			layout.Rigid(page.newItem(&common, page.sub.signW, &page.sub.sign, "Sign message")),
			layout.Rigid(page.newItem(&common, page.sub.verifyW, &page.sub.verify, "Verify message")),
			layout.Rigid(page.newItem(&common, page.sub.changePassW, &page.sub.changePass, "Change passphrase")),
			layout.Rigid(page.newItem(&common, page.sub.deleteW, &page.sub.delete, "Delete wallet")),
		)
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

func (page *walletPage) newItem(common *pageCommon, out decredmaterial.IconButton, in *widget.Button, label string) layout.Widget {
	return func() {
		layout.Inset{Right: unit.Dp(15), Top: unit.Dp(5)}.Layout(common.gtx, func() {
			layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(common.gtx,
				layout.Rigid(func() {
					out.Layout(common.gtx, in)
				}),
				layout.Rigid(func() {
					layout.Inset{Top: unit.Dp(3)}.Layout(common.gtx, func() {
						common.theme.Caption(label).Layout(common.gtx)
					})
				}),
			)
		})
	}
}

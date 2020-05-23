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
		main, delete, rename, cancelDelete, sign widget.Button
		deleteW, renameW, cancelDeleteW, signW   decredmaterial.Button
		mainW                                    decredmaterial.IconButton
		pageInfo                                 decredmaterial.Label
		passwordModal                            *decredmaterial.Password
	}
	signP                                          signMessagePage
	container, accountsList                        layout.List
	rename, delete, addAcct                        widget.Button
	deleteW, gotoVerifyMessagePageBtn              decredmaterial.Button
	renameW, beginRenameW, cancelRenameW, addAcctW decredmaterial.IconButton
	editor, password                               widget.Editor
	editorW, passwordW                             decredmaterial.Editor
	errorLabel                                     decredmaterial.Label
	isPasswordModalOpen                            bool
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
		deleteW:       common.theme.DangerButton("Confirm Delete Wallet"),
	}
	page.sub.mainW = common.theme.IconButton(common.icons.navigationArrowBack)
	page.sub.mainW.Background = color.RGBA{}
	page.sub.mainW.Color = common.theme.Color.Text
	page.sub.mainW.Padding = unit.Dp(0)
	page.sub.mainW.Size = unit.Dp(30)
	page.sub.deleteW = common.theme.DangerButton("Delete Wallet")
	page.sub.signW = common.theme.Button("sign Message")
	page.sub.cancelDeleteW = common.theme.Button("Cancel Wallet Delet")
	page.sub.renameW = common.theme.Button("Rename")
	page.sub.pageInfo = common.theme.Body1("")
	page.sub.passwordModal = common.theme.Password()

	page.SignMessagePage(common)

	return func() {
		page.Layout(common)
		page.Handle(common)
		page.handleSign(common)
	}
}

// Layout lays out the widgets for the main wallets page.
func (page *walletPage) Layout(common pageCommon) {
	if common.walletsTab.Changed() { // reset everything
		page.subPage = subWalletMain
		//page.renaming = false
	}

	// if *common.err != nil {
	// 	page.errorLabel.Text = (*common.err).Error()
	// 	*common.err = nil
	// }

	switch page.subPage {
	case subWalletMain:
		page.subMain(common)
	case subWalletRename:
		page.subRename(common)
	case subWalletDelete:
		page.subDelete(common)
	case subWalletSign:
		page.subSign(common)
	}

}

func (page *walletPage) subMain(common pageCommon) {
	gtx := common.gtx
	page.current = common.info.Wallets[*common.selectedWallet]
	wdgs := []func(){
		func() {
			// <<<<<<< HEAD
			// 			if page.renaming {
			// 				horFlex.Layout(gtx,
			// 					rigid(func() {
			// 						gtx.Constraints.Width.Min = gtx.Px(unit.Dp(350))
			// 						gtx.Constraints.Width.Max = gtx.Constraints.Width.Min
			// 						page.editorW.Layout(gtx, &page.editor)
			// 					}),
			// 					rigid(func() {
			// 						page.renameW.Layout(gtx, &page.rename)
			// 					}),
			// 					rigid(func() {
			// 						layout.Center.Layout(gtx, func() {
			// 							page.cancelRenameW.Layout(gtx, &page.cancelRename)
			// 						})
			// 					}),
			// 				)
			// 			} else {
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
			// 			}
			// =======
			// common.theme.H3(current.Name).Layout(common.gtx)
			// >>>>>>> rogue/wallets
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
					//page.deleteW.Layout(gtx, &page.delete)
				}),
				layout.Rigid(func() {
					layout.Inset{
						Left: unit.Dp(5),
					}.Layout(gtx, func() {
						// page.gotoVerifyMessagePageBtn.Layout(gtx, page.signMessageDiag)
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
	// current := common.info.Wallets[*common.selectedWallet]
	common.Layout(gtx, func() {
		common.theme.H3(page.current.Name).Layout(common.gtx)
	})
}

func (page *walletPage) subDelete(common pageCommon) {
	gtx := common.gtx
	list := layout.List{Axis: layout.Vertical}
	wdgs := []func(){
		// func() {
		// 	page.sub.mainW.Layout(gtx, &page.sub.main)
		// },
		func() {
			common.theme.H3(page.current.Name).Layout(common.gtx)
		},
		func() {
			page.sub.pageInfo.Text = "Are you sure you want to delete " + page.current.Name + "?"
			page.sub.pageInfo.Layout(gtx)
		},
		func() {
			inset := layout.Inset{
				Top:    unit.Dp(20),
				Bottom: unit.Dp(5),
			}
			inset.Layout(gtx, func() {
				page.sub.cancelDeleteW.Layout(gtx, &page.sub.main)
			})
		},
		func() {
			// page.errorLabel.Layout(gtx)
		},
		func() {
			page.deleteW.Layout(gtx, &page.delete)
		},
	}
	common.Layout(gtx, func() {
		layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
			list.Layout(gtx, len(wdgs), func(i int) {
				wdgs[i]()
			})
		})
	})
	if page.isPasswordModalOpen {
		common.Layout(gtx, func() {
			page.sub.passwordModal.Layout(gtx, page.confirm, page.cancel)
		})
	}
}

// Handle handles all widget inputs on the main wallets page.
func (page *walletPage) Handle(common pageCommon) {
	gtx := common.gtx
	// current := common.info.Wallets[*common.selectedWallet]

	// Subs
	if page.sub.delete.Clicked(gtx) {
		page.subPage = subWalletDelete
		return
	}

	if page.sub.sign.Clicked(gtx) {
		page.subPage = subWalletSign
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
		page.isPasswordModalOpen = true
	}

	if page.addWallet.Clicked(gtx) {
		*common.page = PageCreateRestore
	}
}

func (pg *walletPage) confirm(password []byte) {
	pg.isPasswordModalOpen = false
	// pg.isSigningMessage = true

	// pg.signButtonMaterial.Text = "Signing..."
	pg.wallet.DeleteWallet(pg.current.ID, password)
}

func (pg *walletPage) cancel() {
	pg.isPasswordModalOpen = false
}

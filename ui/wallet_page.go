package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

const PageWallet = "wallet"

// const (
// 	subDeleteWallet = iota
// 	subAddAcct
// 	subVerifyMessage
// )

type walletPage struct {
	//subPage                                            int
	container, accountsList                            layout.List
	renaming, deleting                                 bool
	rename, delete, beginRename, cancelRename, addAcct widget.Button
	renameW, beginRenameW, cancelRenameW, addAcctW     decredmaterial.IconButton
	editor                                             widget.Editor
	editorW                                            decredmaterial.Editor
	deleteW                                            decredmaterial.Button
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
		addAcctW:      common.theme.IconButton(common.icons.contentAdd),
		deleteW:       common.theme.DangerButton("Delete Wallet"),
	}

	return func() {
		page.Layout(common)
		page.Handle(common)
	}
}

func (page *walletPage) Layout(common pageCommon) {
	gtx := common.gtx
	current := common.info.Wallets[*common.selectedWallet]
	if common.walletsTab.Changed() {
		page.renaming = false
	}
	wdgs := []func(){
		func() {
			if page.renaming {
				horFlex.Layout(gtx,
					rigid(func() {
						page.editorW.Layout(gtx, &page.editor)
					}),
					rigid(func() {
						page.renameW.Layout(gtx, &page.rename)
					}),
					rigid(func() {
						layout.Center.Layout(gtx, func() {
							page.cancelRenameW.Layout(gtx, &page.cancelRename)
						})
					}),
				)
			} else {
				horFlex.Layout(gtx,
					rigid(func() {
						common.theme.H1(current.Name).Layout(common.gtx)
					}),
					rigid(func() {
						layout.Center.Layout(gtx, func() {
							page.beginRenameW.Layout(gtx, &page.beginRename)
						})
					}),
				)
			}
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
			page.deleteW.Layout(gtx, &page.delete)
		},
	}
	common.LayoutWithWallets(gtx, func() {
		page.container.Layout(common.gtx, len(wdgs), func(i int) {
			wdgs[i]()
		})
	})
}

func (page *walletPage) Handle(common pageCommon) {
	gtx := common.gtx
	current := common.info.Wallets[*common.selectedWallet]
	if page.beginRename.Clicked(gtx) {
		page.renaming = true
		page.editor.SetText(current.Name)
	}

	if page.cancelRename.Clicked(gtx) {
		page.renaming = false
	}

	if page.rename.Clicked(gtx) {
		name := page.editor.Text()
		err := common.wallet.RenameWallet(current.ID, name)
		if err != nil {
			log.Error(err)
		} else {
			common.info.Wallets[*common.selectedWallet].Name = name
			page.renaming = false
		}
	}
}

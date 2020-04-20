package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

const PageWallet = "wallet"

type walletPage struct {
	container                        layout.List
	renaming                         bool
	rename, toggleRename, addAcct    widget.Button
	renameW, toggleRenameW, addAcctW decredmaterial.IconButton
	editor                           widget.Editor
	editorW                          decredmaterial.Editor
}

func WalletPage(common pageCommon) layout.Widget {
	page := walletPage{
		container: layout.List{
			Axis: layout.Vertical,
		},
		toggleRenameW: common.theme.IconButton(common.icons.contentCreate),
		renameW:       common.theme.IconButton(common.icons.contentClear),
		editorW:       common.theme.Editor("Enter wallet name"),
		addAcctW:      common.theme.IconButton(common.icons.contentAdd),
	}
	gtx := common.gtx

	return func() {
		current := common.info.Wallets[*common.selectedWallet]
		wdgs := []func(){
			func() {
				tRename := rigid(func() {
					layout.Center.Layout(gtx, func() {
						page.toggleRenameW.Layout(gtx, &page.toggleRename)
					})
				})
				if page.renaming {
					horFlex.Layout(gtx,
						rigid(func() {
							page.editorW.Layout(gtx, &page.editor)
						}),
						rigid(func() {
							page.renameW.Layout(gtx, &page.rename)
						}),
						tRename,
					)
				} else {
					horFlex.Layout(gtx,
						rigid(func() {
							common.theme.H1(current.Name).Layout(common.gtx)
						}),
						tRename,
					)
				}
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
		}

		page.container.Layout(common.gtx, len(wdgs), func(i int) {
			wdgs[i]()
		})

		if page.toggleRename.Clicked(gtx) {
			page.renaming = !page.renaming
			page.editor.SetText(current.Name)
		}
	}
}

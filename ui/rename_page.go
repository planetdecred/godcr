package ui

// add all the renamePg pages to the wallet page.
//clean up the delete wallet renamePg page
import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

type renameWalletPage struct {
	renameW    decredmaterial.Button
	rename     widget.Button
	errorLabel decredmaterial.Label
	editor     widget.Editor
	editorW    decredmaterial.Editor
}

func (page *walletPage) RenameWalletPage(common pageCommon) {
	page.renamePg = renameWalletPage{
		editorW:    common.theme.Editor("New wallet name"),
		renameW:    common.theme.Button("Rename Wallet"),
		errorLabel: common.theme.Body2(""),
	}
	page.renamePg.errorLabel.Color = common.theme.Color.Danger
}

func (page *walletPage) subRename(common pageCommon) {
	gtx := common.gtx
	list := layout.List{Axis: layout.Vertical}
	wdgs := []func(){
		func() {
			page.returnBtn(common)
		},
		func() {
			common.theme.H5("Rename Wallet").Layout(gtx)
		},
		func() {
			layout.Flex{}.Layout(gtx,
				layout.Rigid(func() {
					layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func() {
						common.theme.Body1("Your are about to rename").Layout(gtx)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Left: unit.Dp(5)}.Layout(gtx, func() {
						txt := common.theme.H5(page.current.Name)
						txt.Color = common.theme.Color.Danger
						txt.Layout(gtx)
					})
				}),
			)
		},
		func() {
			page.renamePg.editorW.Layout(gtx, &page.renamePg.editor)
		},
		func() {
			inset := layout.Inset{
				Top:    unit.Dp(20),
				Bottom: unit.Dp(5),
			}
			inset.Layout(gtx, func() {
				page.renamePg.renameW.Layout(gtx, &page.renamePg.rename)
			})
		},
		func() {
			layout.Center.Layout(common.gtx, func() {
				layout.Inset{Top: unit.Dp(15)}.Layout(gtx, func() {
					page.renamePg.errorLabel.Layout(gtx)
				})
			})
		},
	}
	common.Layout(gtx, func() {
		layout.UniformInset(unit.Dp(20)).Layout(gtx, func() {
			list.Layout(gtx, len(wdgs), func(i int) {
				wdgs[i]()
			})
		})
	})
}

// Handle handles all widget inputs on the main wallets page.
func (page *walletPage) handleRename(common pageCommon) {
	gtx := common.gtx
	if page.renamePg.rename.Clicked(gtx) {
		name := page.renamePg.editor.Text()
		err := common.wallet.RenameWallet(page.current.ID, name)
		if err != nil {
			log.Error(err)
			page.renamePg.errorLabel.Text = err.Error()
		} else {
			common.info.Wallets[*common.selectedWallet].Name = name
			page.subPage = subWalletMain
		}
	}
}

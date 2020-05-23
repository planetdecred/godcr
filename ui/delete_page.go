package ui

// add all the deletePg pages to the wallet page.
//clean up the delete wallet deletePg page
import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

const PageDeleteWallet = "delete_wallet"

type deleteWalletPage struct {
	deleteW, cancelDeleteW decredmaterial.Button
	errorLabel             decredmaterial.Label
	passwordModal          *decredmaterial.Password
	isPasswordModalOpen    bool
	errChannel             chan error
}

func (page *walletPage) DeleteWalletPage(common pageCommon) {
	page.deletePg = deleteWalletPage{
		deleteW:       common.theme.DangerButton("Confirm Delete Wallet"),
		cancelDeleteW: common.theme.Button("Cancel Wallet Delete"),
		errorLabel:    common.theme.Body2(""),
		passwordModal: common.theme.Password(),
		errChannel:    common.errorChannels[PageDeleteWallet],
	}
	page.deletePg.errorLabel.Color = common.theme.Color.Danger
}

func (page *walletPage) subDelete(common pageCommon) {
	gtx := common.gtx
	list := layout.List{Axis: layout.Vertical}
	wdgs := []func(){
		func() {
			common.theme.H5("Delete Wallet").Layout(gtx)
		},
		func() {
			layout.Flex{}.Layout(gtx,
				layout.Rigid(func() {
					layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func() {
						common.theme.Body1("Are you sure you want to delete ").Layout(gtx)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Left: unit.Dp(5)}.Layout(gtx, func() {
						txt := common.theme.H5(page.current.Name)
						txt.Color = common.theme.Color.Danger
						txt.Layout(gtx)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Left: unit.Dp(5)}.Layout(gtx, func() {
						common.theme.H5("?").Layout(gtx)
					})
				}),
			)
		},
		func() {
			inset := layout.Inset{
				Top:    unit.Dp(20),
				Bottom: unit.Dp(5),
			}
			inset.Layout(gtx, func() {
				page.deletePg.cancelDeleteW.Layout(gtx, &page.sub.main)
			})
		},
		func() {
			page.deletePg.deleteW.Layout(gtx, &page.delete)
		},
		func() {
			layout.Center.Layout(common.gtx, func() {
				layout.Inset{Top: unit.Dp(15)}.Layout(gtx, func() {
					page.deletePg.errorLabel.Layout(gtx)
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
	if page.deletePg.isPasswordModalOpen {
		common.Layout(gtx, func() {
			page.deletePg.passwordModal.Layout(gtx, page.confirm, page.cancel)
		})
	}
}

// Handle handles all widget inputs on the main wallets page.
func (page *walletPage) handleDelete(common pageCommon) {
	gtx := common.gtx

	if page.delete.Clicked(gtx) {
		page.deletePg.errorLabel.Text = ""
		page.deletePg.isPasswordModalOpen = true
	}
	select {
	case err := <-page.deletePg.errChannel:
		fmt.Printf("DELETEWALLET PAGE ERROR! %v", err)
		if err.Error() == "invalid_passphrase" {
			page.deletePg.errorLabel.Text = "Wallet passphrase is incorect."
		} else {
			page.deletePg.errorLabel.Text = err.Error()
		}
	default:
	}
}

func (page *walletPage) confirm(password []byte) {
	page.deletePg.isPasswordModalOpen = false
	page.deletePg.errChannel = page.wallet.DeleteWallet(page.current.ID, password)
}

func (page *walletPage) cancel() {
	page.deletePg.isPasswordModalOpen = false
}

package page

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// WalletsID is the id of the wallets page
const WalletsID = "wallets"

// Wallets contains the Wallet, Theme and various widgets
type Wallets struct {
	wal   *wallet.Wallet
	theme *materialplus.Theme

	list         layout.List
	addWalletBtn *widget.Button
	addWalletWdg material.Button

	states map[string]interface{}
}

// Init stores the theme and Wallet
func (pg *Wallets) Init(theme *materialplus.Theme, wal *wallet.Wallet, states map[string]interface{}) {
	pg.wal = wal
	pg.theme = theme

	pg.list = layout.List{Axis: layout.Vertical}

	pg.addWalletBtn = new(widget.Button)
	pg.addWalletWdg = theme.Button("Add Wallet")

	pg.states = states
}

// Draw layouts out the widgets on the given context
func (pg *Wallets) Draw(gtx *layout.Context) interface{} {
	walletInfo := pg.states[StateWalletInfo].(*wallet.MultiWalletInfo)
	if len(walletInfo.Wallets) == 0 {
		pg.theme.Label(units.Label, "Retrieving Wallet Info").Layout(gtx)
		return nil
	}

	widgets := []func(){
		func() {
			pg.theme.Label(units.Label, "Wallets").Layout(gtx)
		},
		func() {
			pg.theme.Label(unit.Dp(20), "ID\t|\tName\t|\tBalance").Layout(gtx)
		},
		func() {
			(&layout.List{
				Axis: layout.Vertical,
			}).Layout(gtx, len(walletInfo.Wallets), func(i int) {
				info := walletInfo.Wallets[i]
				pg.theme.Label(unit.Dp(18), fmt.Sprintf("%d\t|\t%s\t|\t%d atoms", i, info.Name, info.Balance)).Layout(gtx)
			})
		},
		func() {
			pg.addWalletWdg.Layout(gtx, pg.addWalletBtn)
		},
	}
	pg.list.Layout(gtx, len(widgets), layout.ListElement(func(i int) {
		layout.UniformInset(units.FlexInset).Layout(gtx, widgets[i])
	}))

	if pg.addWalletBtn.Clicked(gtx) {
		log.Debugf("{%s} AddWallet Btn Clicked", WalletsID)
		return EventNav{
			Current: WalletsID,
			Next:    LandingID,
		}
	}
	return nil
}

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

	syncBtn *widget.Button
	syncWdg material.Button
}

// Init stores the theme and Wallet
func (pg *Wallets) Init(theme *materialplus.Theme, wal *wallet.Wallet, states map[string]interface{}) {
	pg.wal = wal
	pg.theme = theme

	pg.list = layout.List{Axis: layout.Vertical}

	pg.addWalletBtn = new(widget.Button)
	pg.addWalletWdg = theme.Button("Add Wallet")

	pg.syncBtn = new(widget.Button)
	pg.syncWdg = theme.Button("Start sync")

	pg.states = states
}

// Draw layouts out the widgets on the given context
func (pg *Wallets) Draw(gtx *layout.Context) interface{} {
	walletInfo := pg.states[StateWalletInfo].(*wallet.MultiWalletInfo)
	if len(walletInfo.Wallets) == 0 {
		pg.theme.Label(units.Label, "Retrieving Wallet Info").Layout(gtx)
		return nil
	}

	if walletInfo.Synced {
		pg.syncWdg.Text = "Synced"
	} else if walletInfo.Syncing {
		pg.syncWdg.Text = "Cancel Sync"
	} else {
		pg.syncWdg.Text = "Start Sync"
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
		func() {
			pg.syncWdg.Layout(gtx, pg.syncBtn)
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
	if pg.syncBtn.Clicked(gtx) {
		if !walletInfo.Synced {
			log.Info("Starting sync")
			if err := pg.wal.StartSync(); err != nil {
				log.Error(err)
				pg.syncWdg.Text = "Error starting sync"
			} else {
				pg.syncWdg.Text = "Cancel Sync"
			}
		}
		if walletInfo.Syncing {
			log.Info("Cancelling sync")
			pg.wal.CancelSync()
		}
	}
	return nil
}

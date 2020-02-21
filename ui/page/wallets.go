package page

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr-gio/ui"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// WalletsID is the id of the wallets page
const WalletsID = "wallets"

type walletWdg struct {
	name    material.Label
	balance material.Label
	expand  material.IconButton
	expandB *widget.Button
}

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

	selected int

	expandBtns  []*widget.Button
	expandBtn   material.IconButton
	collapseBtn material.IconButton
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

	expand, err := material.NewIcon(icons.NavigationExpandMore)

	log.Debugf("Load error: %v", err) // TODO: return error

	collapse, err := material.NewIcon(icons.NavigationExpandLess)

	pg.expandBtn = pg.theme.IconButton(expand)
	//pg.expandBtn.Size = unit.Dp(70)

	pg.collapseBtn = pg.theme.IconButton(collapse)

	pg.selected = -1

	pg.states = states
	pg.expandBtns = make([]*widget.Button, 0)
}

// Draw layouts out the widgets on the given context
func (pg *Wallets) Draw(gtx *layout.Context) interface{} {
	walletInfo := pg.states[StateWalletInfo].(*wallet.MultiWalletInfo)
	numWallets := len(walletInfo.Wallets)
	if numWallets == 0 {
		pg.theme.Label(units.Label, "Retrieving Wallet Info").Layout(gtx)
		return nil
	}

	if len(pg.expandBtns) != numWallets {
		pg.expandBtns = make([]*widget.Button, numWallets)
		for i := range pg.expandBtns {
			pg.expandBtns[i] = new(widget.Button)
		}
	}

	if walletInfo.Synced {
		pg.syncWdg.Text = "Synced"
	} else if walletInfo.Syncing {
		pg.syncWdg.Text = "Cancel Sync"
	} else {
		pg.syncWdg.Text = "Start Sync"
	}

	options := []func(){
		func() {
			pg.addWalletWdg.Layout(gtx, pg.addWalletBtn)
		},
		func() {
			pg.syncWdg.Layout(gtx, pg.syncBtn)
		},
	}

	widgets := []func(){
		func() {
			pg.theme.Label(units.Label, "Wallets").Layout(gtx)
		},
		func() {
			(&layout.List{
				Axis: layout.Vertical,
			}).Layout(gtx, len(walletInfo.Wallets), func(i int) {
				layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
					pg.drawWallet(i, walletInfo.Wallets[i], gtx)
				})
			})
		},
		func() {
			layout.Align(layout.Center).Layout(gtx, func() {
				(&layout.List{}).Layout(gtx, len(options), func(i int) {
					layout.UniformInset(unit.Dp(10)).Layout(gtx, options[i])
				})
			})
		},
	}

	pg.list.Layout(gtx, len(widgets), layout.ListElement(func(i int) {
		layout.UniformInset(unit.Dp(10)).Layout(gtx, widgets[i])
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

// drawWallet draws widgets for a single wallet
func (pg *Wallets) drawWallet(id int, info wallet.InfoShort, gtx *layout.Context) {
	widgetFuncs := []func(){
		func() {
			pg.theme.Label(unit.Dp(30), info.Name).Layout(gtx)
		},
		func() {
			pg.theme.Label(unit.Dp(30), info.Balance.String()).Layout(gtx)
		},
		func() {
			if pg.selected != id {
				pg.expandBtn.Layout(gtx, pg.expandBtns[id])
				if pg.expandBtns[id].Clicked(gtx) {
					log.Debugf("Expand %d clicked", id)
					pg.selected = id
				}
			} else {
				pg.collapseBtn.Layout(gtx, pg.expandBtns[id])
				if pg.expandBtns[id].Clicked(gtx) {
					log.Debugf("Collapse %d clicked", id)
					pg.selected = -1
				}
			}

		},
	}

	draw := func() {
		layout.Align(layout.Center).Layout(gtx, func() {
			(&layout.List{}).Layout(gtx, len(widgetFuncs), func(i int) {
				layout.UniformInset(unit.Dp(5)).Layout(gtx, widgetFuncs[i])
			})
		})
	}

	if pg.selected == id {
		funcs := []func(){
			func() {
				draw()
			},
			func() {
				ui.VerticalUniformList(gtx, unit.Dp(4), []func(){
					func() {
						pg.theme.Label(unit.Dp(30), "Accounts").Layout(gtx)
					},
					func() {
						accts := make([]func(), len(info.Accounts))
						for i := range accts {
							//log.Debug("Adding acct")
							accts[i] = func() {
								pg.theme.Label(unit.Dp(20), dcrutil.Amount(info.Accounts[i].TotalBalance).String()).Layout(gtx)
							}
						}
						ui.VerticalUniformList(gtx, unit.Dp(2), accts)
					},
					func() {
						pg.theme.Button("Delete").Layout(gtx, new(widget.Button))
					},
				})
			},
		}
		(&layout.List{Axis: layout.Vertical}).Layout(gtx, len(funcs), func(i int) {
			funcs[i]()
		})
	} else {
		draw()
	}
}

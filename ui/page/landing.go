package page

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// LandingID is the id of the landing page.
const LandingID = "landing"

// Landing represents the landing page of the app.
// It should only be should shown if the app launches
// and cannot find any wallets.
type Landing struct {
	inset     layout.Inset
	container layout.List
	heading   material.Label

	restoreBtn material.Button
	restoreWdg *widget.Button
	createBtn  material.Button
	createWdg  *widget.Button
	walletsBtn material.Button
	walletsWdg *widget.Button

	states map[string]interface{}
}

// Init adds a heading and two buttons.
func (pg *Landing) Init(theme *materialplus.Theme, _ *wallet.Wallet, states map[string]interface{}) {
	pg.heading = theme.Label(units.Label, "Welcome to decred")
	pg.heading.Alignment = text.Middle

	pg.createBtn = theme.Button("Create Wallet")
	pg.createWdg = new(widget.Button)

	pg.restoreBtn = theme.Button("Restore Wallet")
	pg.restoreWdg = new(widget.Button)

	pg.walletsBtn = theme.Button("Back to Wallets")
	pg.walletsWdg = new(widget.Button)

	pg.inset = layout.UniformInset(units.FlexInset)
	pg.container = layout.List{Axis: layout.Vertical}
	pg.states = states
}

// Draw draws the page's to the given layout context.
// Does not react to any event but can return a Nav event.
func (pg *Landing) Draw(gtx *layout.Context) (res interface{}) {
	walletInfo := pg.states[StateWalletInfo].(*wallet.MultiWalletInfo)
	widgets := []func(){
		func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			pg.heading.Layout(gtx)
		},
		func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			if pg.createWdg.Clicked(gtx) {
				fmt.Println("ButtonClicked")
				// res = EventNav {
				// 	Current: LandingID,
				// 	Next: CreateID,
				// }
			}
			pg.createBtn.Layout(gtx, pg.createWdg)
		},
		func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			pg.restoreBtn.Layout(gtx, pg.restoreWdg)
		},
		func() {
			if walletInfo.LoadedWallets != 0 {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				if pg.walletsWdg.Clicked(gtx) {
					res = EventNav{
						Current: LandingID,
						Next:    WalletsID,
					}
				}
				pg.walletsBtn.Layout(gtx, pg.walletsWdg)

			}
		},
	}

	pg.container.Layout(gtx, len(widgets),
		layout.ListElement(func(i int) {
			layout.UniformInset(units.FlexInset).Layout(gtx, widgets[i])
		}),
	)
	return
}

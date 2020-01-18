package page

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/event"
	"github.com/raedahgroup/godcr-gio/ui/units"
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
	page
}

// Init adds a heading and two buttons.
func (pg *Landing) Init(theme *material.Theme, gtx *layout.Context, evt chan event.Event) {
	pg.heading = theme.Label(units.Label, "Welcome to decred")
	pg.heading.Alignment = text.Middle

	pg.createBtn = theme.Button("Create Wallet")
	pg.createWdg = new(widget.Button)

	pg.restoreBtn = theme.Button("Restore Wallet")
	pg.restoreWdg = new(widget.Button)

	pg.inset = layout.UniformInset(units.FlexInset)

	pg.gtx = gtx
	pg.event = evt
}

// Draw adds all the widgets to the stored layout context.
func (pg *Landing) Draw() {
	pg.container.Layout(pg.gtx, 3,
		layout.ListElement(func(i int) {
			switch i {
			case 1:
				pg.inset.Layout(pg.gtx, func() {
					pg.heading.Layout(pg.gtx)
				})
			case 2:
				pg.inset.Layout(pg.gtx, func() {
					if pg.createWdg.Clicked(pg.gtx) {
						fmt.Println("ButtonClicked")
						// pg.event <- event.Nav {
						// 	Current: LandingID,
						// 	Next: CreateID
						// }
					}
					pg.createBtn.Layout(pg.gtx, pg.createWdg)
				})
			default:
				pg.inset.Layout(pg.gtx, func() {
					pg.restoreBtn.Layout(pg.gtx, pg.restoreWdg)
				})
			}
		}),
	)
}

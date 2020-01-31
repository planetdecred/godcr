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
}

// Init adds a heading and two buttons.
func (pg *Landing) Init(theme *material.Theme) {
	pg.heading = theme.Label(units.Label, "Welcome to decred")
	pg.heading.Alignment = text.Middle

	pg.createBtn = theme.Button("Create Wallet")
	pg.createWdg = new(widget.Button)

	pg.restoreBtn = theme.Button("Restore Wallet")
	pg.restoreWdg = new(widget.Button)

	pg.inset = layout.UniformInset(units.FlexInset)
}

// Draw draws the page's to the given layout context.
// Does not react to any event but can return a Nav event.
func (pg *Landing) Draw(gtx *layout.Context, _ event.Event) (evt event.Event) {
	pg.container.Layout(gtx, 3,
		layout.ListElement(func(i int) {
			switch i {
			case 0:
				pg.inset.Layout(gtx, func() {
					pg.heading.Layout(gtx)
				})
			case 1:
				pg.inset.Layout(gtx, func() {
					if pg.createWdg.Clicked(gtx) {
						fmt.Println("ButtonClicked")
						//evt = event.Nav {
						// 	Current: LandingID,
						// 	Next: CreateID
						// }
					}
					pg.createBtn.Layout(gtx, pg.createWdg)
				})
			default:
				pg.inset.Layout(gtx, func() {
					pg.restoreBtn.Layout(gtx, pg.restoreWdg)
					if pg.restoreWdg.Clicked(gtx) {
						fmt.Println("ButtonClicked")
						evt = event.Nav{
							Current: LandingID,
							Next:    RestoreID,
						}
					}
				})
			}
		}),
	)
	return
}

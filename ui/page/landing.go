package page

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui/units"
)

// Landing represents the landing page of the app.
// It should only be should shown if the app launches
// and cannot find any wallets
type Landing struct {
	inset   layout.Inset
	stack   layout.Stack
	heading material.Label

	restoreBtn material.Button
	restoreWdg *widget.Button
	createBtn  material.Button
	createWdg  *widget.Button
	gtk        *layout.Context
}

// Init adds a heading and two buttons
func (page *Landing) Init(theme *material.Theme, gtk *layout.Context) {
	page.heading = theme.Label(units.Label, "Welcome to decred")
	page.heading.Alignment = text.Middle

	page.createBtn = theme.Button("Create Wallet")
	page.createWdg = new(widget.Button)

	page.restoreBtn = theme.Button("Restore Wallet")
	page.restoreWdg = new(widget.Button)

	page.inset = layout.UniformInset(units.FlexInset)

	page.stack.Alignment = layout.W

	page.gtk = gtk
}

// Draw adds all the widgets to the stored layout context
func (page *Landing) Draw() {
	page.stack.Layout(page.gtk,
		layout.Stacked(func() {
			page.inset.Layout(page.gtk, func() {
				page.heading.Layout(page.gtk)
			})
		}),
		layout.Stacked(func() {
			page.inset.Layout(page.gtk, func() {
				page.createBtn.Layout(page.gtk, page.createWdg)
			})
		}),
		layout.Stacked(func() {
			page.inset.Layout(page.gtk, func() {
				page.restoreBtn.Layout(page.gtk, page.restoreWdg)
			})
		}),
	)
}

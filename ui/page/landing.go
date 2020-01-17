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
	Heading       material.Label
	CreateWallet  material.Button
	RestoreWallet material.Button
	createBtn     *widget.Button
	gtk           *layout.Context
}

// Init adds a heading and two buttons
func (page *Landing) Init(theme *material.Theme, gtk *layout.Context) {
	heading := theme.Label(units.Label, "Welcome to decred")
	heading.Alignment = text.Middle

	create := theme.Button("Create Wallet")
	cbtn := new(widget.Button)

	page.createBtn = cbtn
	page.CreateWallet = create

	page.gtk = gtk
	page.Heading = heading
}

// Draw adds all the widgets to the stored layout context
func (page *Landing) Draw() {
	page.Heading.Layout(page.gtk)
	page.CreateWallet.Layout(page.gtk, page.createBtn)
}

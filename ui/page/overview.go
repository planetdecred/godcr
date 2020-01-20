package page

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/raedahgroup/godcr-gio/ui/units"
)

type Overview struct {
	gtk *layout.Context

	balance material.Label
}

var list = layout.List{
	Axis: layout.Vertical,
}

// Init initializes all widgets to be used on the page
func (page *Overview) Init(theme *material.Theme, gtk *layout.Context) {
	heading := theme.Label(units.Label, "Welcome to decred")
	heading.Alignment = text.Middle

	page.gtk = gtk
	page.balance = theme.H5("154.0928281")
}

// Draw adds all the widgets to the stored layout context
func (page *Overview) Draw() {
	widgets := []func(){
		func() {
			page.balance.Layout(page.gtk)
		},
	}

	list.Layout(page.gtk, len(widgets), func(i int) {
		layout.UniformInset(units.Padding).Layout(page.gtk, widgets[i])
	})
}

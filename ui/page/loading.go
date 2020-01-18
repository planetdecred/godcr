package page

import (
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui/units"
)

// LoadingID is the id of the loading page.
const LoadingID = "loading"

// Loading represents the loading page of the app.
type Loading struct {
	page
	label material.Label
}

// Init initializies the page with a label.
func (pg *Loading) Init(theme *material.Theme, gtx *layout.Context) {
	pg.gtx = gtx
	pg.label = theme.Label(units.Label, "LOADING")
}

// Draw renders the page widgets.
func (pg *Loading) Draw() {
	pg.label.Layout(pg.gtx)
}

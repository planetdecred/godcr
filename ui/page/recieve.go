package page

import (
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/event"
	"github.com/raedahgroup/godcr-gio/ui/units"
)

// ReceivingID is the id of the receiving page.
const ReceivingID = "receiving"

// Receiving represents the receiving page of the app.
type Receiving struct {
	label material.Label
}

// Init initializies the page with a label.
func (pg *Receiving) Init(theme *material.Theme) {
	pg.label = theme.Label(units.Label, "Receive DCR")
	pg.label = theme.Label(units.Label, "Receive DCR")
	pg.label = theme.Label(units.Label, "Receive DCR")
	pg.label = theme.Label(units.Label, "Receive DCR")

}

// Draw renders the page widgets.
// It does not react to nor does it generate any event.
func (pg *Receiving) Draw(gtx *layout.Context, _ event.Event) event.Event {
	pg.label.Layout(gtx)
	pg.label.Layout(gtx)
	return nil
}

// Package page provides an interface and implementations
// for creating and using pages.
package page

import (
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/event"
)

// Page represents a single page of the app.
//
// Init creates widgets with the given theme.
// Implementations can store the event Duplex for
// communication with the window
//
// Draw draws the implementation's widgets to the given
// layout context reacting to the given event if needed
// and returns any generated user event.
// Draw is only called once per frame for the active page.
type Page interface {
	Init(*material.Theme)
	Draw(*layout.Context, event.Event) event.Event
}

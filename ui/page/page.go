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
// Init creates widgets with the given theme and
// layout context.
// Implementations must store the layout context and event channel.
//
// Draw adds the implementation's widgets to the stored
// layout context
type Page interface {
	Init(*material.Theme, *layout.Context, chan event.Event)
	Draw()
}

// page encapsulates the base structure needed for
// a Page implementation.
type page struct {
	event chan event.Event
	gtx   *layout.Context
}

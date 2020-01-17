// Package page provides an interface and implementations
// for creating and using pages
package page

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

// Page represents a single page of the app
//
// Init creates widgets with the given theme and
// layout context
// Implementations must store the layout context
//
// Draw adds the implementation's widgets to the stored
// layout context
type Page interface {
	Init(*material.Theme, *layout.Context)
	Draw()
}

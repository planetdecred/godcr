package load

import "gioui.org/layout"

type Page interface {
	ID() string
	OnResume() // called when a page is starting or resuming from a paused state.
	Layout(layout.Context) layout.Dimensions
	Handle()
	OnClose()
}

type Modal interface {
	ModalID() string
	OnResume()
	Layout(gtx layout.Context) layout.Dimensions
	OnDismiss()
	Dismiss()
	Show()
	Handle()
}

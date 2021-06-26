package ui

import (
	"gioui.org/layout"
)

type Modal interface {
	ModalID() string
	OnResume()
	Layout(gtx layout.Context) layout.Dimensions
	OnDismiss()
	Dismiss()
	Show()
	Handle()
}

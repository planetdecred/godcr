package ui

import (
	"time"

	"github.com/planetdecred/godcr/ui/values"

	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

type (
	toast struct {
		text    string
		success bool
		timer   *time.Timer
	}
)

func displayToast(th *decredmaterial.Theme, gtx layout.Context, n *toast) layout.Dimensions {
	color := th.Color.Success
	if !n.success {
		color = th.Color.Danger
	}

	return decredmaterial.Card{Color: color, Rounded: true}.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Top: values.MarginPadding7, Bottom: values.MarginPadding7,
			Left: values.MarginPadding15, Right: values.MarginPadding15,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			t := th.Body1(n.text)
			t.Color = th.Color.Surface
			return t.Layout(gtx)
		})
	})
}

// Timer implements the standard time.AfterFunc. It only creates a new time.Timer when toast.timer is nil.
// Gio re-renders its UI recursively and Timer prevents multiple time.Timer from being created before the duration is
// exceeded.
func (n *toast) Timer(d time.Duration, f func()) {
	if n.timer != nil {
		return
	}
	n.timer = time.AfterFunc(d, f)
}

func (n *toast) ResetTimer() {
	n.timer = nil
}

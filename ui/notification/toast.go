package notification

import (
	"time"

	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Toast struct {
	theme   *decredmaterial.Theme
	success bool
	message string
	timer   *time.Timer

	show bool
}

func NewToast(th *decredmaterial.Theme) *Toast {
	return &Toast{
		theme: th,
	}
}

func (t *Toast) Notify(message string, success bool) {
	t.message = message
	t.success = success
	t.show = true
}

func (t *Toast) Layout(gtx layout.Context) layout.Dimensions {
	if !t.show {
		return layout.Dimensions{}
	}

	t.handleToastDisplay()
	color := t.theme.Color.Success
	if !t.success {
		color = t.theme.Color.Danger
	}

	card := t.theme.Card()
	card.Color = color
	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.Inset{Top: values.MarginPadding65}.Layout(gtx, func(gtx C) D {
			return card.Layout(gtx, func(gtx C) D {
				return layout.Inset{
					Top: values.MarginPadding7, Bottom: values.MarginPadding7,
					Left: values.MarginPadding15, Right: values.MarginPadding15,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					msg := t.theme.Body1(t.message)
					msg.Color = t.theme.Color.Surface
					return msg.Layout(gtx)
				})
			})
		})
	})
}

func (t *Toast) handleToastDisplay() {
	if t.show {
		// create a new timer if the Notify method is called by another process
		t.timer = time.NewTimer(time.Second * 3)
	}

	if t.timer == nil {
		t.timer = time.NewTimer(time.Second * 3)
	}

	select {
	case <-t.timer.C:
		t.show = false
		t.timer = nil
	default:
	}
}

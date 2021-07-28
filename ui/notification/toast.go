package notification

import (
	"sync"
	"time"

	"gioui.org/layout"
	"gioui.org/op"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Toast struct {
	sync.Mutex
	theme   *decredmaterial.Theme
	success bool
	message string
	timer   *time.Timer
}

func NewToast(th *decredmaterial.Theme) *Toast {
	return &Toast{
		theme: th,
	}
}

func (t *Toast) Notify(message string, success bool) {
	t.Lock()
	t.message = message
	t.success = success
	t.timer = time.NewTimer(time.Second * 3)
	t.Unlock()
}

func (t *Toast) Layout(gtx layout.Context) layout.Dimensions {
	if t.timer == nil {
		return layout.Dimensions{}
	}

	t.handleToastDisplay(gtx)
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

func (t *Toast) handleToastDisplay(gtx layout.Context) {
	select {
	case <-t.timer.C:
		t.timer = nil
		op.InvalidateOp{}.Add(gtx.Ops)
	default:
	}
}

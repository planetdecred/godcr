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

type duration int32

const (
	Short duration = iota
	Long
)

func NewToast(th *decredmaterial.Theme) *Toast {
	return &Toast{
		theme: th,
	}
}

func getDurationFromDelay(d duration) time.Duration {
	switch d {
	case Short:
		return 2 * time.Second
	case Long:
		return 5 * time.Second
	default:
		return 2 * time.Second
	}
}

func (t *Toast) Notify(message string, success bool, d ...duration) {
	var notificationDelay duration
	if len(d) > 0 {
		notificationDelay = d[0]
	}

	t.Lock()
	t.message = message
	t.success = success
	t.timer = time.NewTimer(getDurationFromDelay(notificationDelay))
	t.Unlock()
}

func (t *Toast) Layout(gtx layout.Context) layout.Dimensions {
	t.handleToastDisplay(gtx)
	if t.timer == nil {
		return layout.Dimensions{}
	}

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
	if t.timer == nil {
		return
	}

	select {
	case <-t.timer.C:
		t.timer = nil
		op.InvalidateOp{}.Add(gtx.Ops)
	default:
	}
}

package notification

import (
	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

func (c *Conductor) LayoutNotifications(th *decredmaterial.Theme, gtx layout.Context) layout.Dimensions {
	active := c.GetActiveNotification()

	color := th.Color.Success
	switch active.NotificationType {
	case Error:
		color = th.Color.Danger
	case Warning:
		color = th.Color.Success2 // Not sure about this
	}

	card := th.Card()
	card.Color = color
	return card.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: values.MarginPadding7, Bottom: values.MarginPadding7,
			Left: values.MarginPadding15, Right: values.MarginPadding15,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			t := th.Body1(active.text)
			t.Color = th.Color.Surface
			return t.Layout(gtx)
		})
	})
}

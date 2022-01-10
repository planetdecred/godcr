package dexclient

import (
	"decred.org/dcrdex/client/db"
	"gioui.org/layout"
	"gioui.org/text"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const notificationModalID = "dex_notifications_modal"

const dexNotificationConfigKey = "dex_notifications"

type notificationModal struct {
	*load.Load
	modal        *decredmaterial.Modal
	severityIcon *decredmaterial.Icon
}

func newNotificationModal(l *load.Load) *notificationModal {
	nmd := &notificationModal{
		Load:         l,
		modal:        l.Theme.ModalFloatTitle(),
		severityIcon: decredmaterial.NewIcon(l.Icons.ImageBrightness1),
	}
	nmd.modal.SetPadding(values.MarginPadding0)

	return nmd
}

func (nmd *notificationModal) ModalID() string {
	return notificationModalID
}

func (nmd *notificationModal) Show() {
	nmd.ShowModal(nmd)
}

func (nmd *notificationModal) Dismiss() {
	nmd.DismissModal(nmd)
}

func (nmd *notificationModal) OnDismiss() {
}

func (nmd *notificationModal) OnResume() {
}

func (nmd *notificationModal) Handle() {
	if nmd.modal.BackdropClicked(true) {
		nmd.Dismiss()
	}
}

func (nmd *notificationModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return layout.Inset{
				Top:    values.MarginPadding5,
				Bottom: values.MarginPadding5,
				Left:   values.MarginPadding10,
			}.Layout(gtx, nmd.Load.Theme.Label(values.TextSize20, values.String(values.StrNotifications)).Layout)
		},
		func(gtx C) D {
			var notifications []*db.Notification
			err := nmd.WL.MultiWallet.ReadUserConfigValue(dexNotificationConfigKey, &notifications)
			if err != nil {
				return D{}
			}

			childrens := make([]layout.FlexChild, 0, len(notifications))
			for i, ntfn := range notifications {
				n := ntfn
				index := i
				childrens = append(childrens, layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							if index == 0 {
								return D{}
							}
							return layout.Inset{
								Top:    values.MarginPadding5,
								Bottom: values.MarginPadding5,
							}.Layout(gtx, nmd.Theme.Separator().Layout)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									switch n.Severity() {
									case db.Success:
										nmd.severityIcon.Color = nmd.Theme.Color.Success
									case db.WarningLevel:
										nmd.severityIcon.Color = nmd.Theme.Color.Orange3
									case db.ErrorLevel:
										nmd.severityIcon.Color = nmd.Theme.Color.Danger
									default:
										nmd.severityIcon.Color = nmd.Theme.Color.Background
									}
									return layout.Inset{
										Top:   values.MarginPadding6,
										Right: values.MarginPadding8,
									}.Layout(gtx, func(gtx C) D {
										return nmd.severityIcon.Layout(gtx, values.MarginPadding8)
									})
								}),
								layout.Flexed(1, func(gtx C) D {
									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											textLabel := nmd.Theme.Label(values.TextSize14, n.Subject())
											textLabel.Font.Weight = text.Bold
											return textLabel.Layout(gtx)
										}),
										layout.Rigid(nmd.Theme.Label(values.TextSize14, n.Details()).Layout),
									)
								}),
								layout.Rigid(func(gtx C) D {
									gtx.Constraints.Min.X = 80
									return nmd.Theme.Label(values.TextSize14, timeSince(n.TimeStamp)).Layout(gtx)
								}),
							)
						}),
					)
				}))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childrens...)
		},
	}

	return nmd.modal.Layout(gtx, w)
}

package dexclient

import (
	"encoding/hex"
	"time"

	"decred.org/dcrdex/dex"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const notificationModalID = "dex_notifications_modal"

const dexNotificationConfigKey = "dex_notifications"

type notificationModal struct {
	*load.Load
	modal             *decredmaterial.Modal
	severityIcon      *decredmaterial.Icon
	notificationBtn   decredmaterial.Button
	recentActivityBtn decredmaterial.Button
	recentActivity    *[]*notification
	isRecentActivity  bool
}

func newNotificationModal(l *load.Load, recentActivity *[]*notification) *notificationModal {
	tabButton := func(t string, active bool) decredmaterial.Button {
		btn := l.Theme.OutlineButton(t)
		btn.CornerRadius = values.MarginPadding0
		btn.Inset = layout.Inset{
			Top:    values.MarginPadding5,
			Bottom: values.MarginPadding5,
			Left:   values.MarginPadding9,
			Right:  values.MarginPadding9,
		}
		btn.TextSize = values.TextSize14
		btn.Color = l.Theme.Color.Gray3
		if active {
			btn.Font.Weight = text.Bold
			btn.Color = l.Theme.Color.Primary
		}
		return btn
	}

	nmd := &notificationModal{
		Load:              l,
		modal:             l.Theme.ModalFloatTitle(),
		severityIcon:      decredmaterial.NewIcon(l.Icons.ImageBrightness1),
		recentActivity:    recentActivity,
		notificationBtn:   tabButton(values.String(values.StrNotifications), true),
		recentActivityBtn: tabButton(strRecentActivity, false),
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
	notifications := getNotifications(nmd.WL.MultiWallet, true)
	ids := make([]dex.Bytes, 0, len(notifications))
	for _, n := range notifications {
		b, err := hex.DecodeString(n.ID)
		if err != nil {
			continue
		}
		ids = append(ids, b)
	}

	if len(ids) > 0 {
		go func() {
			nmd.Dexc().Core().AckNotes(ids)
		}()
	}
}

func (nmd *notificationModal) Handle() {
	if nmd.notificationBtn.Clicked() {
		nmd.isRecentActivity = false
		nmd.notificationBtn.Font.Weight = text.Bold
		nmd.notificationBtn.Color = nmd.Theme.Color.Primary
		nmd.recentActivityBtn.Font.Weight = text.Normal
		nmd.recentActivityBtn.Color = nmd.Theme.Color.Gray3
	}

	if nmd.recentActivityBtn.Clicked() {
		nmd.isRecentActivity = true
		nmd.recentActivityBtn.Font.Weight = text.Bold
		nmd.recentActivityBtn.Color = nmd.Theme.Color.Primary
		nmd.notificationBtn.Font.Weight = text.Normal
		nmd.notificationBtn.Color = nmd.Theme.Color.Gray3
	}

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
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(.5, nmd.notificationBtn.Layout),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Left:  values.MarginPadding1,
							Right: values.MarginPadding1,
						}.Layout(gtx, func(gtx C) D { return D{} })
					}),
					layout.Flexed(.5, nmd.recentActivityBtn.Layout),
				)
			})
		},
		nmd.notificationLayout(),
		nmd.recentActivityLayout(),
	}

	op.InvalidateOp{At: gtx.Now.Add(1 * time.Second)}.Add(gtx.Ops)
	return nmd.modal.Layout(gtx, w)
}

func (nmd *notificationModal) notificationLayout() layout.Widget {
	return func(gtx C) D {
		gtx.Constraints.Min.Y = 300
		if nmd.isRecentActivity {
			return D{}
		}

		notifications := getNotifications(nmd.WL.MultiWallet, false)
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
								return layout.Inset{
									Top:   values.MarginPadding6,
									Right: values.MarginPadding8,
								}.Layout(gtx, func(gtx C) D {
									nmd.severityIcon.Color = severityColor(n.Severity, nmd.Theme.Color)
									return nmd.severityIcon.Layout(gtx, values.MarginPadding8)
								})
							}),
							layout.Flexed(1, func(gtx C) D {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										textLabel := nmd.Theme.Label(values.TextSize14, n.Subject)
										textLabel.Font.Weight = text.Bold
										return textLabel.Layout(gtx)
									}),
									layout.Rigid(nmd.Theme.Label(values.TextSize14, n.Details).Layout),
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
	}
}

func (nmd *notificationModal) recentActivityLayout() layout.Widget {
	return func(gtx C) D {
		gtx.Constraints.Min.Y = 300
		if !nmd.isRecentActivity {
			return D{}
		}

		notifications := *nmd.recentActivity
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
							layout.Flexed(1, func(gtx C) D {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										ts := time.Unix(int64(n.TimeStamp), 0).Format("Jan 2, 2006 15:04:05 PM")
										textLabel := nmd.Theme.Label(values.TextSize14, ts)
										return textLabel.Layout(gtx)
									}),
									layout.Rigid(nmd.Theme.Label(values.TextSize14, n.Details).Layout),
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
	}
}

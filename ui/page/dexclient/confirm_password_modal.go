package dexclient

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const confirmPasswordModalID = "confirm_password_modal"

type confirmPasswordModal struct {
	*load.Load
	modal       *decredmaterial.Modal
	submit      decredmaterial.Button
	appPassword decredmaterial.Editor
	confirmed   func([]byte)
}

func newconfirmPasswordModal(l *load.Load) *confirmPasswordModal {
	md := &confirmPasswordModal{
		Load:        l,
		modal:       l.Theme.ModalFloatTitle(),
		appPassword: l.Theme.EditorPassword(new(widget.Editor), "Password"),
		submit:      l.Theme.Button(new(widget.Clickable), "OK"),
	}

	md.submit.TextSize = values.TextSize12
	md.submit.Background = l.Theme.Color.Primary
	md.appPassword.Editor.SingleLine = true

	return md
}

func (md *confirmPasswordModal) ModalID() string {
	return confirmPasswordModalID
}

func (md *confirmPasswordModal) Show() {
	md.ShowModal(md)
}

func (md *confirmPasswordModal) Dismiss() {
	md.DismissModal(md)
}

func (md *confirmPasswordModal) OnDismiss() {
	md.appPassword.Editor.SetText("")
}

func (md *confirmPasswordModal) OnResume() {
	md.appPassword.Editor.Focus()
}

func (md *confirmPasswordModal) Handle() {
	if md.submit.Button.Clicked() {
		if strings.Trim(md.appPassword.Editor.Text(), " ") == "" {
			return
		}
		md.confirmed([]byte(md.appPassword.Editor.Text()))
		md.Dismiss()
	}
}

func (md *confirmPasswordModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return md.Load.Theme.Label(values.TextSize14, "Authorize this order with your app password.").Layout(gtx)
			})
		},
		func(gtx C) D {
			return md.appPassword.Layout(gtx)
		},
		func(gtx C) D {
			return md.submit.Layout(gtx)
		},
	}

	return md.modal.Layout(gtx, w, 900)
}

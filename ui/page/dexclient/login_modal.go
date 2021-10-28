package dexclient

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const dexLoginModalID = "dex_login_modal"

type loginModal struct {
	*load.Load
	modal    *decredmaterial.Modal
	loggedIn func()

	submit      decredmaterial.Button
	appPassword decredmaterial.Editor
}

func newloginModal(l *load.Load) *loginModal {
	md := &loginModal{
		Load:        l,
		modal:       l.Theme.ModalFloatTitle(),
		submit:      l.Theme.Button("Login"),
		appPassword: l.Theme.EditorPassword(new(widget.Editor), "App password"),
	}

	md.submit.TextSize = values.TextSize12
	md.submit.Background = l.Theme.Color.Primary
	md.appPassword.Editor.SingleLine = true

	return md
}

func (md *loginModal) ModalID() string {
	return dexLoginModalID
}

func (md *loginModal) Show() {
	md.ShowModal(md)
}

func (md *loginModal) Dismiss() {
	md.DismissModal(md)
}

func (md *loginModal) OnDismiss() {
	md.appPassword.Editor.SetText("")
}

func (md *loginModal) OnResume() {
	md.appPassword.Editor.Focus()
}

func (md *loginModal) Handle() {
	if md.submit.Button.Clicked() {
		_, err := md.Dexc.Login([]byte(md.appPassword.Editor.Text()))
		if err != nil {
			md.Toast.NotifyError(err.Error())
			return
		}

		md.loggedIn()
		md.Dismiss()
	}
}

func (md *loginModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return md.Theme.H6("Login").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return md.appPassword.Layout(gtx)
					})
				}),
				layout.Rigid(md.submit.Layout),
			)
		},
	}

	return md.modal.Layout(gtx, w, 900)
}

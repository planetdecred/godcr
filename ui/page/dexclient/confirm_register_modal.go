package dexclient

import (
	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const confirmRegisterModalID = "confirm_register_modal"

type confirmRegisterModal struct {
	*load.Load
	modal       *decredmaterial.Modal
	register    decredmaterial.Button
	appPassword decredmaterial.Editor
	isSending   bool
	cert        []byte
	ce          *core.Exchange
	confirmed   func([]byte)
}

func newConfirmRegisterModal(l *load.Load) *confirmRegisterModal {
	md := &confirmRegisterModal{
		Load:        l,
		modal:       l.Theme.ModalFloatTitle(),
		appPassword: l.Theme.EditorPassword(new(widget.Editor), "App password"),
		register:    l.Theme.Button("Register"),
	}

	md.register.TextSize = values.TextSize12
	md.register.Background = l.Theme.Color.Primary
	md.appPassword.Editor.SingleLine = true

	return md
}

func (md *confirmRegisterModal) ModalID() string {
	return confirmRegisterModalID
}

func (md *confirmRegisterModal) Show() {
	md.ShowModal(md)
}

func (md *confirmRegisterModal) Dismiss() {
	md.DismissModal(md)
}

func (md *confirmRegisterModal) OnDismiss() {
	md.appPassword.Editor.SetText("")
}

func (md *confirmRegisterModal) OnResume() {
	md.appPassword.Editor.Focus()
}

func (md *confirmRegisterModal) Handle() {
	if md.register.Button.Clicked() {
		if md.appPassword.Editor.Text() == "" || md.isSending {
			return
		}

		md.isSending = true
		go func() {
			form := &core.RegisterForm{
				AppPass: []byte(md.appPassword.Editor.Text()),
				Addr:    md.ce.Host,
				Fee:     md.ce.Fee.Amt,
				Cert:    md.cert,
			}
			_, err := md.Dexc.Register(form)

			md.isSending = false
			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
			}

			md.confirmed([]byte(md.appPassword.Editor.Text()))
			md.Dismiss()
		}()
	}

}

func (md *confirmRegisterModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return md.Load.Theme.Label(values.TextSize20, "Confirm Registration").Layout(gtx)
		},
		func(gtx C) D {
			return md.Load.Theme.Label(values.TextSize14, "Enter your app password to confirm DEX registration. When you submit this form, 1.000 DCR will be spent from your Decred wallet to pay registration fees.").Layout(gtx)
		},
		func(gtx C) D {
			return md.Load.Theme.Label(values.TextSize14, "The DCR lot size for this DEX is 1.000 DCR. All trades are in multiples of this lot size. This is the minimum possible trade amount in DCR.").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
				return md.appPassword.Layout(gtx)
			})
		},
		func(gtx C) D {
			return md.register.Layout(gtx)
		},
	}

	return md.modal.Layout(gtx, w, 900)
}

func (md *confirmRegisterModal) updateCertAndExchange(cert []byte, ce *core.Exchange) {
	md.cert = cert
	md.ce = ce
}

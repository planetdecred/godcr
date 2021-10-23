package dexclient

import (
	"fmt"
	"strconv"

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
	dex         *core.Exchange
	cert        []byte
	confirmed   func()
}

func newConfirmRegisterModal(l *load.Load, dex *core.Exchange, cert []byte) *confirmRegisterModal {
	md := &confirmRegisterModal{
		Load:        l,
		modal:       l.Theme.ModalFloatTitle(),
		appPassword: l.Theme.EditorPassword(new(widget.Editor), "App password"),
		register:    l.Theme.Button("Register"),
		dex:         dex,
		cert:        cert,
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
				Addr:    md.dex.Host,
				Fee:     md.dex.Fee.Amt,
				Cert:    md.cert,
			}
			_, err := md.Dexc.Register(form)

			md.isSending = false
			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
			}

			md.confirmed()
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
			amount := strconv.FormatFloat(float64(md.dex.Fee.Amt)/1e8, 'f', -1, 64)
			txt := fmt.Sprintf("Enter your app password to confirm DEX registration. When you submit this form, %s DCR will be spent from your Decred wallet to pay registration fees.", amount)
			return md.Load.Theme.Label(values.TextSize14, txt).Layout(gtx)
		},
		func(gtx C) D {
			// lotSize := md.dex.Markets[md.dex.Host].LotSize
			txt := fmt.Sprintf("The DCR lot size for this DEX is %d DCR. All trades are in multiples of this lot size. This is the minimum possible trade amount in DCR.", 1)
			return md.Load.Theme.Label(values.TextSize14, txt).Layout(gtx)
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

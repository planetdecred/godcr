package dexclient

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const dexPasswordModalID = "dex_password_modal"

type passwordModal struct {
	*load.Load
	modal *decredmaterial.Modal

	createPassword                decredmaterial.Button
	appPassword, appPasswordAgain decredmaterial.Editor
	marketBaseID, marketQuoteID   uint32

	isSending    bool
	appInitiated func()
}

func newPasswordModal(l *load.Load) *passwordModal {
	md := &passwordModal{
		Load:             l,
		modal:            l.Theme.ModalFloatTitle(),
		appPassword:      l.Theme.EditorPassword(new(widget.Editor), "Password"),
		appPasswordAgain: l.Theme.EditorPassword(new(widget.Editor), "Password again"),
		createPassword:   l.Theme.Button("Create password"),
	}

	md.createPassword.TextSize = values.TextSize12
	md.createPassword.Background = l.Theme.Color.Primary
	md.appPassword.Editor.SingleLine = true
	md.appPasswordAgain.Editor.SingleLine = true

	return md
}

func (md *passwordModal) ModalID() string {
	return dexPasswordModalID
}

func (md *passwordModal) Show() {
	md.ShowModal(md)
}

func (md *passwordModal) Dismiss() {
	md.DismissModal(md)
}

func (md *passwordModal) OnDismiss() {
	md.appPassword.Editor.SetText("")
}

func (md *passwordModal) OnResume() {
	md.appPassword.Editor.Focus()
}

func (md *passwordModal) Handle() {
	if md.createPassword.Button.Clicked() {
		if md.appPasswordAgain.Editor.Text() != md.appPassword.Editor.Text() || md.isSending {
			return
		}

		md.isSending = true
		go func() {
			// TODO: Generate and save a 64-byte seed and pass it to InitializeClient
			// to enable dex restores if the dex db becomes corrupted. Alternatively,
			// passing nil will cause dex to generate a random seed which can be saved
			// for later dex restoration efforts.
			err := md.DL.InitializeClient([]byte(md.appPassword.Editor.Text()), nil)
			md.isSending = false
			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
			}
			md.appInitiated()

			md.Dismiss()
		}()
	}
}

func (md *passwordModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return md.Theme.Label(values.TextSize20, "Set App Password").Layout(gtx)
		},
		func(gtx C) D {
			return md.Theme.Label(values.TextSize14, "Set your app password. This password will protect your DEX account keys and connected wallets.").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return md.appPassword.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return md.appPasswordAgain.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return md.createPassword.Layout(gtx)
		},
	}

	return md.modal.Layout(gtx, w, 900)
}

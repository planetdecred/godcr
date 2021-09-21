package dexclient

import (
	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const addDexModalID = "add_dex_modal"

type addDexModal struct {
	*load.Load
	modal            *decredmaterial.Modal
	addDexServer     decredmaterial.Button
	dexServerAddress decredmaterial.Editor
	isSending        bool
	cert             decredmaterial.Editor
	created          func([]byte, *core.Exchange)
}

func newAddDexModal(l *load.Load) *addDexModal {
	md := &addDexModal{
		Load:             l,
		modal:            l.Theme.ModalFloatTitle(),
		dexServerAddress: l.Theme.Editor(new(widget.Editor), "DEX Address"),
		addDexServer:     l.Theme.Button(new(widget.Clickable), "Submit"),
		cert:             l.Theme.Editor(new(widget.Editor), "Cert content"),
	}

	md.addDexServer.TextSize = values.TextSize12
	md.addDexServer.Background = l.Theme.Color.Primary
	md.dexServerAddress.Editor.SingleLine = true

	return md
}

func (md *addDexModal) ModalID() string {
	return dexPasswordModalID
}

func (md *addDexModal) Show() {
	md.ShowModal(md)
}

func (md *addDexModal) Dismiss() {
	md.DismissModal(md)
}

func (md *addDexModal) OnDismiss() {
	md.dexServerAddress.Editor.SetText("")
}

func (md *addDexModal) OnResume() {
	md.dexServerAddress.Editor.Focus()
}

func (md *addDexModal) Handle() {
	if md.addDexServer.Button.Clicked() {
		if md.dexServerAddress.Editor.Text() == "" || md.isSending {
			return
		}

		md.isSending = true
		go func() {
			c := []byte(md.cert.Editor.Text())
			ce, err := md.DL.GetDEXConfig(md.dexServerAddress.Editor.Text(), c)
			md.isSending = false
			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
			}

			md.created(c, ce)
			md.Dismiss()
		}()
	}

}

func (md *addDexModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return md.Load.Theme.Label(values.TextSize20, "Add a dex").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return md.dexServerAddress.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return md.cert.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return md.addDexServer.Layout(gtx)
		},
	}

	return md.modal.Layout(gtx, w, 900)
}

package dexclient

import (
	"fmt"
	"io/ioutil"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/sqweek/dialog"
)

const addDexModalID = "add_dex_modal"

const testDexHost = "127.0.0.1:7232"

type addDexModal struct {
	*load.Load
	modal                     *decredmaterial.Modal
	addDexServer, addCertFile decredmaterial.Button
	dexServerAddress          decredmaterial.Editor
	isSending                 bool
	cert                      []byte
	created                   func([]byte, *core.Exchange)
}

func newAddDexModal(l *load.Load) *addDexModal {
	md := &addDexModal{
		Load:             l,
		modal:            l.Theme.ModalFloatTitle(),
		dexServerAddress: l.Theme.Editor(new(widget.Editor), "DEX Address"),
		addDexServer:     l.Theme.Button(new(widget.Clickable), "Submit"),
		addCertFile:      l.Theme.Button(new(widget.Clickable), "Add a file"),
	}

	md.addDexServer.TextSize = values.TextSize12
	md.addDexServer.Background = l.Theme.Color.Primary
	md.dexServerAddress.Editor.SingleLine = true
	md.dexServerAddress.Editor.SetText(fmt.Sprintf("http://%s", testDexHost))

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
			ce, err := md.DL.GetDEXConfig(md.dexServerAddress.Editor.Text(), md.cert)
			md.isSending = false
			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
			}

			md.created(md.cert, ce)
			md.Dismiss()
		}()
	}

	if md.addCertFile.Button.Clicked() {
		go func() {
			filename, err := dialog.File().Filter("Select TLS Certificate", "cert").Load()

			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
			}

			content, err := ioutil.ReadFile(filename)
			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
			}
			md.cert = content
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
						return md.addCertFile.Layout(gtx)
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

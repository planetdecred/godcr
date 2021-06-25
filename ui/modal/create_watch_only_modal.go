package modal

import (
	"fmt"
	"github.com/planetdecred/godcr/ui"
	"github.com/planetdecred/godcr/ui/page"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const CreateWatchOnly = "create_watch_only_modal"

type createWatchOnlyModal struct {
	*ui.Common
	randomID string
	modal    decredmaterial.Modal

	walletName     decredmaterial.Editor
	extendedPubKey decredmaterial.Editor

	isLoading      bool
	materialLoader material.LoaderStyle

	callback   func(walletName, extPubKey string, m *createWatchOnlyModal) bool // return true to dismiss dialog
	btnPositve decredmaterial.Button

	btnNegative decredmaterial.Button
}

func NewCreateWatchOnlyModal(c *ui.Common) *createWatchOnlyModal {
	cm := &createWatchOnlyModal{
		Common:  c,
		randomID:    fmt.Sprintf("%s-%d", CreateWatchOnly, ui.GenerateRandomNumber()),
		modal:       *c.Theme.ModalFloatTitle(),
		btnPositve:  c.Theme.Button(new(widget.Clickable), values.String(values.StrImport)),
		btnNegative: c.Theme.Button(new(widget.Clickable), values.String(values.StrCancel)),
	}

	cm.btnPositve.TextSize, cm.btnNegative.TextSize = values.TextSize16, values.TextSize16
	cm.btnPositve.Font.Weight, cm.btnNegative.Font.Weight = text.Bold, text.Bold

	cm.walletName = c.Theme.Editor(new(widget.Editor), "Wallet name")
	cm.walletName.Editor.SingleLine, cm.walletName.Editor.Submit = true, true

	cm.extendedPubKey = c.Theme.EditorPassword(new(widget.Editor), "Extended public key")
	cm.extendedPubKey.Editor.Submit = true

	th := material.NewTheme(gofont.Collection())
	cm.materialLoader = material.Loader(th)

	return cm
}

func (cm *createWatchOnlyModal) ModalID() string {
	return cm.randomID
}

func (cm *createWatchOnlyModal) OnResume() {
}

func (cm *createWatchOnlyModal) OnDismiss() {

}

func (cm *createWatchOnlyModal) Show() {
	cm.ShowModal(cm)
}

func (cm *createWatchOnlyModal) Dismiss() {
	cm.DismissModal(cm)
}

func (cm *createWatchOnlyModal) setLoading(loading bool) {
	cm.isLoading = loading
}

func (cm *createWatchOnlyModal) setError(err string) {
	if err == "" {
		cm.extendedPubKey.ClearError()
	} else {
		cm.extendedPubKey.SetError(err)
	}
}

func (cm *createWatchOnlyModal) watchOnlyCreated(callback func(walletName, extPubKey string, m *createWatchOnlyModal) bool) *createWatchOnlyModal {
	cm.callback = callback
	return cm
}

func (cm *createWatchOnlyModal) Handle() {

	if ui.EditorsNotEmpty(cm.walletName.Editor, cm.extendedPubKey.Editor) ||
		ui.HandleSubmitEvent(cm.walletName.Editor, cm.extendedPubKey.Editor) {
		for cm.btnPositve.Button.Clicked() {
			cm.setLoading(true)
			if cm.callback(cm.walletName.Editor.Text(), cm.extendedPubKey.Editor.Text(), cm) {
				cm.Dismiss()
			}
		}
	}

	if cm.btnNegative.Button.Clicked() {
		if !cm.isLoading {
			cm.Dismiss()
		}
	}
}

func (cm *createWatchOnlyModal) Layout(gtx layout.Context) page.D {
	w := []layout.Widget{
		func(gtx page.C) page.D {
			t := cm.Theme.H6(values.String(values.StrImportWatchingOnlyWallet))
			t.Font.Weight = text.Bold
			return t.Layout(gtx)
		},
		func(gtx page.C) page.D {
			return cm.walletName.Layout(gtx)
		},
		func(gtx page.C) page.D {
			return cm.extendedPubKey.Layout(gtx)
		},
		func(gtx page.C) page.D {
			return layout.E.Layout(gtx, func(gtx page.C) page.D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx page.C) page.D {

						cm.btnNegative.Background = cm.Theme.Color.Surface
						cm.btnNegative.Color = cm.Theme.Color.Primary
						return cm.btnNegative.Layout(gtx)
					}),
					layout.Rigid(func(gtx page.C) page.D {
						if cm.isLoading {
							return cm.materialLoader.Layout(gtx)
						}
						cm.btnPositve.Background, cm.btnPositve.Color = cm.Theme.Color.Surface, cm.Theme.Color.Primary
						return cm.btnPositve.Layout(gtx)
					}),
				)
			})
		},
	}

	return cm.modal.Layout(gtx, w, 850)
}

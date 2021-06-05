package ui

import (
	"fmt"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

type createWatchOnlyModal struct {
	*pageCommon
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

func newCreateWatchOnlyModal(common *pageCommon) *createWatchOnlyModal {
	cm := &createWatchOnlyModal{
		pageCommon:  common,
		randomID:    fmt.Sprintf("%s-%d", ModalInfo, generateRandomNumber()),
		modal:       *common.theme.ModalFloatTitle(),
		btnPositve:  common.theme.Button(new(widget.Clickable), values.String(values.StrImport)),
		btnNegative: common.theme.Button(new(widget.Clickable), values.String(values.StrCancel)),
	}

	cm.btnPositve.TextSize, cm.btnNegative.TextSize = values.TextSize16, values.TextSize16
	cm.btnPositve.Font.Weight, cm.btnNegative.Font.Weight = text.Bold, text.Bold

	cm.walletName = common.theme.Editor(new(widget.Editor), "Wallet name")
	cm.walletName.Editor.SingleLine, cm.walletName.Editor.Submit = true, true

	cm.extendedPubKey = common.theme.EditorPassword(new(widget.Editor), "Extended public key")
	cm.extendedPubKey.Editor.Submit = true

	th := material.NewTheme(gofont.Collection())
	cm.materialLoader = material.Loader(th)

	return cm
}

func (cm *createWatchOnlyModal) modalID() string {
	return cm.randomID
}

func (cm *createWatchOnlyModal) OnResume() {
}

func (cm *createWatchOnlyModal) OnDismiss() {

}

func (cm *createWatchOnlyModal) show() {
	cm.showModal(cm)
}

func (cm *createWatchOnlyModal) dismiss() {
	cm.dismissModal(cm)
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

func (cm *createWatchOnlyModal) callbackFunc(callback func(walletName, extPubKey string, m *createWatchOnlyModal) bool) *createWatchOnlyModal {
	cm.callback = callback
	return cm
}

func (cm *createWatchOnlyModal) handle() {

	if editorsNotEmpty(cm.walletName.Editor, cm.extendedPubKey.Editor) ||
		handleSubmitEvent(cm.walletName.Editor, cm.extendedPubKey.Editor) {
		for cm.btnPositve.Button.Clicked() {
			cm.setLoading(true)
			if cm.callback(cm.walletName.Editor.Text(), cm.extendedPubKey.Editor.Text(), cm) {
				cm.dismiss()
			}
		}
	}

	if cm.btnNegative.Button.Clicked() {
		if !cm.isLoading {
			cm.dismiss()
		}
	}
}

func (cm *createWatchOnlyModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			t := cm.theme.H6(values.String(values.StrImportWatchingOnlyWallet))
			t.Font.Weight = text.Bold
			return t.Layout(gtx)
		},
		func(gtx C) D {
			return cm.walletName.Layout(gtx)
		},
		func(gtx C) D {
			return cm.extendedPubKey.Layout(gtx)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {

						cm.btnNegative.Background = cm.theme.Color.Surface
						cm.btnNegative.Color = cm.theme.Color.Primary
						return cm.btnNegative.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						if cm.isLoading {
							return cm.materialLoader.Layout(gtx)
						}
						cm.btnPositve.Background, cm.btnPositve.Color = cm.theme.Color.Surface, cm.theme.Color.Primary
						return cm.btnPositve.Layout(gtx)
					}),
				)
			})
		},
	}

	return cm.modal.Layout(gtx, w, 850)
}

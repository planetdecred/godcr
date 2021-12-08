package modal

import (
	"fmt"

	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const CreateWatchOnly = "create_watch_only_modal"

type CreateWatchOnlyModal struct {
	*load.Load

	modal          decredmaterial.Modal
	materialLoader material.LoaderStyle

	walletName     decredmaterial.Editor
	extendedPubKey decredmaterial.Editor

	btnPositve  decredmaterial.Button
	btnNegative decredmaterial.Button
	keyEvent    chan *key.Event

	randomID string

	isLoading    bool
	isCancelable bool
	isEnabled    bool

	callback func(walletName, extPubKey string, m *CreateWatchOnlyModal) bool // return true to dismiss dialog
}

func NewCreateWatchOnlyModal(l *load.Load) *CreateWatchOnlyModal {
	cm := &CreateWatchOnlyModal{
		Load:         l,
		randomID:     fmt.Sprintf("%s-%d", CreateWatchOnly, decredmaterial.GenerateRandomNumber()),
		modal:        *l.Theme.ModalFloatTitle(),
		btnPositve:   l.Theme.Button(values.String(values.StrImport)),
		btnNegative:  l.Theme.OutlineButton(values.String(values.StrCancel)),
		isCancelable: true,
		keyEvent:     l.Receiver.KeyEvents,
	}

	cm.btnPositve.Font.Weight = text.Medium

	cm.btnNegative.Font.Weight = text.Medium
	cm.btnNegative.Margin = layout.Inset{Right: values.MarginPadding8}

	cm.walletName = l.Theme.Editor(new(widget.Editor), "Wallet name")
	cm.walletName.Editor.SingleLine, cm.walletName.Editor.Submit = true, true

	cm.extendedPubKey = l.Theme.EditorPassword(new(widget.Editor), "Extended public key")
	cm.extendedPubKey.Editor.Submit = true

	th := material.NewTheme(gofont.Collection())
	cm.materialLoader = material.Loader(th)

	return cm
}

func (cm *CreateWatchOnlyModal) ModalID() string {
	return cm.randomID
}

func (cm *CreateWatchOnlyModal) OnResume() {
	cm.walletName.Editor.Focus()
	cm.Load.EnableKeyEvent = true
}

func (cm *CreateWatchOnlyModal) OnDismiss() {
	cm.Load.EnableKeyEvent = false
}

func (cm *CreateWatchOnlyModal) Show() {
	cm.ShowModal(cm)
}

func (cm *CreateWatchOnlyModal) Dismiss() {
	cm.DismissModal(cm)
}

func (cm *CreateWatchOnlyModal) SetLoading(loading bool) {
	cm.isLoading = loading
}

func (cm *CreateWatchOnlyModal) SetCancelable(min bool) *CreateWatchOnlyModal {
	cm.isCancelable = min
	return cm
}

func (cm *CreateWatchOnlyModal) SetError(err string) {
	cm.extendedPubKey.SetError(err)
}

func (cm *CreateWatchOnlyModal) WatchOnlyCreated(callback func(walletName, extPubKey string, m *CreateWatchOnlyModal) bool) *CreateWatchOnlyModal {
	cm.callback = callback
	return cm
}

func (cm *CreateWatchOnlyModal) Handle() {
	if editorsNotEmpty(cm.walletName.Editor) ||
		editorsNotEmpty(cm.extendedPubKey.Editor) {
		cm.btnPositve.Background = cm.Theme.Color.Primary
		cm.isEnabled = true
	} else {
		cm.btnPositve.Background = cm.Theme.Color.Gray3
		cm.isEnabled = false
	}

	isSubmit, isChanged := decredmaterial.HandleEditorEvents(cm.walletName.Editor, cm.extendedPubKey.Editor)
	if isChanged {
		// reset editor errors
		cm.walletName.SetError("")
		cm.extendedPubKey.SetError("")
	}

	for (cm.btnPositve.Clicked() || isSubmit) && cm.isEnabled {
		if !editorsNotEmpty(cm.walletName.Editor) {
			cm.walletName.SetError("enter wallet name")
			return
		}

		if !editorsNotEmpty(cm.extendedPubKey.Editor) {
			cm.extendedPubKey.SetError("enter a valid extendedPubKey")
			return
		}

		cm.SetLoading(true)
		if cm.callback(cm.walletName.Editor.Text(), cm.extendedPubKey.Editor.Text(), cm) {
			cm.Dismiss()
		}
	}

	cm.btnNegative.SetEnabled(!cm.isLoading)
	if cm.btnNegative.Clicked() {
		if !cm.isLoading {
			cm.Dismiss()
		}
	}

	if cm.modal.BackdropClicked(cm.isCancelable) {
		if !cm.isLoading {
			cm.Dismiss()
		}
	}
	decredmaterial.SwitchEditors(cm.keyEvent, cm.walletName.Editor, cm.extendedPubKey.Editor)
}

func (cm *CreateWatchOnlyModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			t := cm.Theme.H6(values.String(values.StrImportWatchingOnlyWallet))
			t.Font.Weight = text.SemiBold
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
						if cm.isLoading {
							return D{}
						}

						return cm.btnNegative.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						if cm.isLoading {
							return cm.materialLoader.Layout(gtx)
						}
						return cm.btnPositve.Layout(gtx)
					}),
				)
			})
		},
	}

	return cm.modal.Layout(gtx, w)
}

package modal

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const HexModal = "Hex_modal"

type HexRestoreModal struct {
	*load.Load
	randomID  string
	modal     decredmaterial.Modal
	hexEditor decredmaterial.Editor

	dialogTitle string
	description string

	isLoading    bool
	isCancelable bool

	customWidget layout.Widget

	materialLoader material.LoaderStyle

	positiveButtonText    string
	positiveButtonClicked func(password string, m *HexRestoreModal) bool
	btnPositve            decredmaterial.Button

	negativeButtonText    string
	negativeButtonClicked func()
	btnNegative           decredmaterial.Button
}

func NewHexRestoreModal(l *load.Load) *HexRestoreModal {
	hm := &HexRestoreModal{
		Load:         l,
		randomID:     fmt.Sprintf("%s-%d", HexModal, decredmaterial.GenerateRandomNumber()),
		modal:        *l.Theme.ModalFloatTitle(),
		btnPositve:   l.Theme.Button("Confirm"),
		btnNegative:  l.Theme.OutlineButton("Cancel"),
		isCancelable: true,
	}

	hm.btnPositve.Font.Weight = text.Medium

	hm.btnNegative.Font.Weight = text.Medium
	hm.btnNegative.Margin.Right = values.MarginPadding8

	hm.hexEditor = l.Theme.Editor(new(widget.Editor), "Hint")
	hm.hexEditor.Editor.SingleLine, hm.hexEditor.Editor.Submit = true, true

	hm.materialLoader = material.Loader(l.Theme.Base)

	return hm
}

func (hm *HexRestoreModal) ModalID() string {
	return hm.randomID
}

func (hm *HexRestoreModal) OnResume() {
	hm.hexEditor.Editor.Focus()
}

func (hm *HexRestoreModal) OnDismiss() {}

func (hm *HexRestoreModal) Show() {
	hm.ShowModal(hm)
	hm.btnPositve.SetEnabled(false)
}

func (hm *HexRestoreModal) Dismiss() {
	hm.DismissModal(hm)
}

func (hm *HexRestoreModal) Title(title string) *HexRestoreModal {
	hm.dialogTitle = title
	return hm
}

func (hm *HexRestoreModal) Description(description string) *HexRestoreModal {
	hm.description = description
	return hm
}

func (hm *HexRestoreModal) UseCustomWidget(layout layout.Widget) *HexRestoreModal {
	hm.customWidget = layout
	return hm
}

func (hm *HexRestoreModal) Hint(hint string) *HexRestoreModal {
	hm.hexEditor.Hint = hint
	return hm
}

func (hm *HexRestoreModal) PositiveButton(text string, clicked func(hex string, m *HexRestoreModal) bool) *HexRestoreModal {
	hm.positiveButtonText = text
	hm.positiveButtonClicked = clicked
	return hm
}

func (hm *HexRestoreModal) NegativeButton(text string, clicked func()) *HexRestoreModal {
	hm.negativeButtonText = text
	hm.negativeButtonClicked = clicked
	return hm
}

func (hm *HexRestoreModal) SetLoading(loading bool) {
	hm.isLoading = loading
	hm.modal.SetDisabled(loading)
}

func (hm *HexRestoreModal) SetCancelable(min bool) *HexRestoreModal {
	hm.isCancelable = min
	return hm
}

func (hm *HexRestoreModal) SetError(err string) {
	if err == "" {
		hm.hexEditor.ClearError()
	} else {
		hm.hexEditor.SetError(err)
	}
}

func (hm *HexRestoreModal) Handle() {
	isSubmit, isChanged := decredmaterial.HandleEditorEvents(hm.hexEditor.Editor)
	if isChanged {
		hm.hexEditor.SetError("")
		hex := hm.hexEditor.Editor.Text()
		if len(hex) >= 16 && len(hex) <= 64 {
			hm.btnPositve.SetEnabled(true)
		} else {
			hm.btnPositve.SetEnabled(false)
		}
	}

	if hm.btnPositve.Button.Clicked() || isSubmit {
		if !editorsNotEmpty(hm.hexEditor.Editor) {
			hm.hexEditor.SetError("Enter Hex")
			return
		}

		if hm.isLoading {
			return
		}

		hm.SetLoading(true)
		hm.SetError("")
		if hm.positiveButtonClicked(hm.hexEditor.Editor.Text(), hm) {
			hm.DismissModal(hm)
		}
	}

	hm.btnNegative.SetEnabled(!hm.isLoading)
	for hm.btnNegative.Clicked() {
		if !hm.isLoading {
			hm.DismissModal(hm)
			hm.negativeButtonClicked()
		}
	}

	if hm.modal.BackdropClicked(hm.isCancelable) {
		if !hm.isLoading {
			hm.Dismiss()
			hm.negativeButtonClicked()
		}
	}
}

func (hm *HexRestoreModal) Layout(gtx layout.Context) D {
	title := func(gtx C) D {
		t := hm.Theme.H6(hm.dialogTitle)
		t.Font.Weight = text.SemiBold
		return t.Layout(gtx)
	}

	description := func(gtx C) D {
		t := hm.Theme.Body2(hm.description)
		return t.Layout(gtx)
	}

	editor := func(gtx C) D {
		return hm.hexEditor.Layout(gtx)
	}

	actionButtons := func(gtx C) D {
		return layout.E.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if hm.negativeButtonText == "" || hm.isLoading {
						return layout.Dimensions{}
					}

					hm.btnNegative.Text = hm.negativeButtonText
					return hm.btnNegative.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if hm.isLoading {
						return hm.materialLoader.Layout(gtx)
					}

					if hm.positiveButtonText == "" {
						return layout.Dimensions{}
					}

					hm.btnPositve.Text = hm.positiveButtonText
					return hm.btnPositve.Layout(gtx)
				}),
			)
		})
	}
	var w []layout.Widget

	w = append(w, title)

	if hm.description != "" {
		w = append(w, description)
	}

	if hm.customWidget != nil {
		w = append(w, hm.customWidget)
	}

	w = append(w, editor)
	w = append(w, actionButtons)

	return hm.modal.Layout(gtx, w)
}

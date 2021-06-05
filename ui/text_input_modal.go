package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

type textInputModal struct {
	*infoModal

	isLoading bool

	textInput decredmaterial.Editor
	callback  func(string, *textInputModal) bool
}

func newTextInputModal(common *pageCommon) *textInputModal {
	tm := &textInputModal{
		infoModal: newInfoModal(common),
	}

	tm.textInput = common.theme.Editor(new(widget.Editor), "Hint")
	tm.textInput.Editor.SingleLine, tm.textInput.Editor.Submit = true, true

	return tm
}

func (tm *textInputModal) show() {
	tm.pageCommon.showModal(tm)
}

func (tm *textInputModal) hint(hint string) *textInputModal {
	tm.textInput.Hint = hint
	return tm
}

func (tm *textInputModal) positiveButton(text string, callback func(string, *textInputModal) bool) *textInputModal {
	tm.positiveButtonText = text
	tm.callback = callback
	return tm
}

func (tm *textInputModal) setError(err string) {
	if err == "" {
		tm.textInput.ClearError()
	} else {
		tm.textInput.SetError(err)
	}
}

func (tm *textInputModal) handle() {

	for tm.btnPositve.Button.Clicked() {
		if tm.isLoading {
			continue
		}

		tm.isLoading = true
		tm.setError("")
		if tm.callback(tm.textInput.Editor.Text(), tm) {
			tm.dismiss()
		}
	}

	for tm.btnNegative.Button.Clicked() {
		if !tm.isLoading {
			tm.dismiss()
			tm.negativeButtonClicked()
		}
	}
}

func (tm *textInputModal) Layout(gtx layout.Context) D {

	var w []layout.Widget

	if tm.dialogTitle != "" {
		w = append(w, tm.titleLayout(gtx))
	}

	w = append(w, tm.textInput.Layout)
	if tm.negativeButtonText != "" || tm.positiveButtonText != "" {
		w = append(w, tm.actionButtonsLayout(gtx))
	}

	return tm.modal.Layout(gtx, w, 850)
}

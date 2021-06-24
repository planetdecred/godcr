package modal

import (
	"fmt"
	"github.com/planetdecred/godcr/ui"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const TextInput = "text_input_modal"

type textInputModal struct {
	*infoModal

	isLoading bool

	textInput decredmaterial.Editor
	callback  func(string, *textInputModal) bool
}

func NewTextInputModal(common *ui.Common) *textInputModal {
	tm := &textInputModal{
		infoModal: NewInfoModal(common),
	}

	tm.randomID = fmt.Sprintf("%s-%d", TextInput, ui.GenerateRandomNumber())

	tm.textInput = common.Theme.Editor(new(widget.Editor), "Hint")
	tm.textInput.Editor.SingleLine, tm.textInput.Editor.Submit = true, true

	return tm
}

func (tm *textInputModal) Show() {
	tm.ShowModal(tm)
}

func (tm *textInputModal) Dismiss() {
	tm.DismissModal(tm)
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
			tm.Dismiss()
		}
	}

	for tm.btnNegative.Button.Clicked() {
		if !tm.isLoading {
			tm.Dismiss()
			tm.negativeButtonClicked()
		}
	}
}

func (tm *textInputModal) Layout(gtx layout.Context) ui.D {

	var w []layout.Widget

	if tm.dialogTitle != "" {
		w = append(w, tm.titleLayout())
	}

	w = append(w, tm.textInput.Layout)
	if tm.negativeButtonText != "" || tm.positiveButtonText != "" {
		w = append(w, tm.actionButtonsLayout())
	}

	return tm.modal.Layout(gtx, w, 850)
}

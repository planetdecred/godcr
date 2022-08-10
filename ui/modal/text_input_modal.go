package modal

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/renderers"
	"github.com/planetdecred/godcr/ui/values"
)

type TextInputModal struct {
	*InfoModal

	isLoading           bool
	showAccountWarnInfo bool
	isCancelable        bool
	isEnabled           bool

	textInput decredmaterial.Editor
	callback  func(string, *TextInputModal) bool

	positiveButtonColor color.NRGBA
	textCustomTemplate  []layout.Widget
}

func NewTextInputModal(l *load.Load) *TextInputModal {
	tm := &TextInputModal{
		InfoModal:    NewInfoModalWithKey(l, "text_input_modal", Outline),
		isCancelable: true,
	}

	tm.textInput = l.Theme.Editor(new(widget.Editor), values.String(values.StrHint))
	tm.textInput.Editor.SingleLine, tm.textInput.Editor.Submit = true, true

	return tm
}

func (tm *TextInputModal) OnResume() {
	tm.textInput.Editor.Focus()
}

func (tm *TextInputModal) Hint(hint string) *TextInputModal {
	tm.textInput.Hint = hint
	return tm
}

func (tm *TextInputModal) SetLoading(loading bool) {
	tm.isLoading = loading
	tm.Modal.SetDisabled(loading)
}

func (tm *TextInputModal) ShowAccountInfoTip(show bool) *TextInputModal {
	tm.showAccountWarnInfo = show
	return tm
}

func (tm *TextInputModal) PositiveButton(text string, callback func(string, *TextInputModal) bool) *TextInputModal {
	tm.positiveButtonText = text
	tm.callback = callback
	return tm
}

func (tm *TextInputModal) PositiveButtonStyle(background, text color.NRGBA) *TextInputModal {
	tm.positiveButtonColor, tm.btnPositive.Color = background, text
	return tm
}

func (tm *TextInputModal) SetError(err string) {
	if err == "" {
		tm.textInput.ClearError()
	} else {
		tm.textInput.SetError(err)
	}
}

func (tm *TextInputModal) SetCancelable(min bool) *TextInputModal {
	tm.isCancelable = min
	return tm
}

func (tm *TextInputModal) SetTextWithTemplate(template string) *TextInputModal {
	switch template {
	case AllowUnmixedSpendingTemplate:
		tm.textCustomTemplate = allowUnspendUnmixedAcct(tm.Load)
	}
	return tm
}

func (tm *TextInputModal) Handle() {

	if editorsNotEmpty(tm.textInput.Editor) {
		tm.btnPositive.Background = tm.positiveButtonColor
		tm.isEnabled = true
	} else {
		tm.btnPositive.Background = tm.Theme.Color.Gray3
		tm.isEnabled = false
	}

	isSubmit, isChanged := decredmaterial.HandleEditorEvents(tm.textInput.Editor)
	if isChanged {
		tm.textInput.SetError("")
	}

	if (tm.btnPositive.Clicked() || isSubmit) && tm.isEnabled {
		if tm.isLoading {
			return
		}

		tm.SetLoading(true)
		tm.SetError("")
		if tm.callback(tm.textInput.Editor.Text(), tm) {
			tm.Dismiss()
		}
	}

	for tm.btnNegative.Clicked() {
		if !tm.isLoading {
			tm.Dismiss()
			tm.negativeButtonClicked()
		}
	}

	if tm.Modal.BackdropClicked(tm.isCancelable) {
		if !tm.isLoading {
			tm.Dismiss()
			tm.negativeButtonClicked()
		}
	}
}

func (tm *TextInputModal) Layout(gtx layout.Context) D {

	var w []layout.Widget

	if tm.dialogTitle != "" {
		w = append(w, tm.titleLayout())
	}

	if tm.showAccountWarnInfo {
		l := func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					img := decredmaterial.NewIcon(tm.Theme.Icons.ActionInfo)
					img.Color = tm.Theme.Color.Gray1
					inset := layout.Inset{Right: values.MarginPadding4}
					return inset.Layout(gtx, func(gtx C) D {
						return img.Layout(gtx, values.MarginPadding20)
					})
				}),
				layout.Rigid(func(gtx C) D {
					text := values.StringF(values.StrAddAcctWarn, `<span style="text-color: grayText1">`, `<span style="font-weight: bold">`, `</span>`, `</span>`)
					return renderers.RenderHTML(text, tm.Theme).Layout(gtx)
				}),
			)
		}
		w = append(w, l)
	}

	if tm.textCustomTemplate != nil {
		w = append(w, tm.textCustomTemplate...)
	}

	w = append(w, tm.textInput.Layout)

	if tm.negativeButtonText != "" || tm.positiveButtonText != "" {
		w = append(w, tm.actionButtonsLayout())
	}

	return tm.Modal.Layout(gtx, w)
}

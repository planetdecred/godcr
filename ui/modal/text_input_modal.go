package modal

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const TextInput = "text_input_modal"

type TextInputModal struct {
	*InfoModal

	IsLoading           bool
	showAccountWarnInfo bool
	isCancelable        bool
	isEnabled           bool

	textInput decredmaterial.Editor
	callback  func(string, *TextInputModal) bool
}

func NewTextInputModal(l *load.Load) *TextInputModal {
	tm := &TextInputModal{
		InfoModal:    NewInfoModal(l),
		isCancelable: true,
	}

	tm.randomID = fmt.Sprintf("%s-%d", TextInput, generateRandomNumber())

	tm.textInput = l.Theme.Editor(new(widget.Editor), "Hint")
	tm.textInput.Editor.SingleLine, tm.textInput.Editor.Submit = true, true

	return tm
}

func (tm *TextInputModal) Show() {
	tm.ShowModal(tm)
}

func (tm *TextInputModal) OnResume() {
	tm.textInput.Editor.Focus()
}

func (tm *TextInputModal) Dismiss() {
	tm.DismissModal(tm)
}

func (tm *TextInputModal) Hint(hint string) *TextInputModal {
	tm.textInput.Hint = hint
	return tm
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
	tm.btnPositve.Background, tm.btnPositve.Color = background, text
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

func (tm *TextInputModal) Handle() {
	if editorsNotEmpty(tm.textInput.Editor) {
		tm.btnPositve.Background = tm.Theme.Color.Primary
		tm.isEnabled = true
	} else {
		tm.btnPositve.Background = tm.Theme.Color.InactiveGray
		tm.isEnabled = false
	}

	isSubmit, isChanged := decredmaterial.HandleEditorEvents(tm.textInput.Editor)
	if isChanged {
		tm.textInput.SetError("")
	}

	if (tm.btnPositve.Button.Clicked() || isSubmit) && tm.isEnabled {
		if tm.IsLoading {
			return
		}

		tm.IsLoading = true
		tm.SetError("")
		if tm.callback(tm.textInput.Editor.Text(), tm) {
			tm.Dismiss()
		}
	}

	for tm.btnNegative.Button.Clicked() {
		if !tm.IsLoading {
			tm.Dismiss()
			tm.negativeButtonClicked()
		}
	}

	if tm.modal.BackdropClicked(tm.isCancelable) {
		if !tm.IsLoading {
			tm.Dismiss()
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
					img := tm.Icons.ActionInfo
					img.Color = tm.Theme.Color.Gray3
					inset := layout.Inset{Right: values.MarginPadding4}
					return inset.Layout(gtx, func(gtx C) D {
						return img.Layout(gtx, values.MarginPadding20)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := tm.Theme.Label(values.MarginPadding16, "Accounts")
							txt.Color = tm.Theme.Color.Gray4
							return txt.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							txt := tm.Theme.Label(values.MarginPadding16, "cannot")
							txt.Font.Weight = text.Bold
							txt.Color = tm.Theme.Color.Gray4
							inset := layout.Inset{Right: values.MarginPadding2, Left: values.MarginPadding2}
							return inset.Layout(gtx, txt.Layout)
						}),
						layout.Rigid(func(gtx C) D {
							txt := tm.Theme.Label(values.MarginPadding16, "be deleted once created")
							txt.Color = tm.Theme.Color.Gray4
							return txt.Layout(gtx)
						}),
					)
				}),
			)
		}
		w = append(w, l)
	}

	w = append(w, tm.textInput.Layout)

	if tm.negativeButtonText != "" || tm.positiveButtonText != "" {
		w = append(w, tm.actionButtonsLayout())
	}

	return tm.modal.Layout(gtx, w, 850)
}

package modal

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const TextInput = "text_input_modal"

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

func NewTextInputModal(app *app.App) *TextInputModal {
	tm := &TextInputModal{
		InfoModal:    NewInfoModal(app),
		isCancelable: true,
	}

	tm.randomID = fmt.Sprintf("%s-%d", TextInput, decredmaterial.GenerateRandomNumber())

	tm.textInput = app.Theme.Editor(new(widget.Editor), "Hint")
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

func (tm *TextInputModal) SetLoading(loading bool) {
	tm.isLoading = loading
	tm.modal.SetDisabled(loading)
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
	tm.positiveButtonColor, tm.btnPositve.Color = background, text
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
		tm.textCustomTemplate = allowUnspendUnmixedAcct(tm.Theme)
	}
	return tm
}

func (tm *TextInputModal) Handle() {

	if editorsNotEmpty(tm.textInput.Editor) {
		tm.btnPositve.Background = tm.positiveButtonColor
		tm.isEnabled = true
	} else {
		tm.btnPositve.Background = tm.Theme.Color.Gray3
		tm.isEnabled = false
	}

	isSubmit, isChanged := decredmaterial.HandleEditorEvents(tm.textInput.Editor)
	if isChanged {
		tm.textInput.SetError("")
	}

	if (tm.btnPositve.Clicked() || isSubmit) && tm.isEnabled {
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

	if tm.modal.BackdropClicked(tm.isCancelable) {
		if !tm.isLoading {
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
					img := decredmaterial.NewIcon(tm.Theme.Icons.ActionInfo)
					img.Color = tm.Theme.Color.Gray1
					inset := layout.Inset{Right: values.MarginPadding4}
					return inset.Layout(gtx, func(gtx C) D {
						return img.Layout(gtx, values.MarginPadding20)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := tm.Theme.Label(values.MarginPadding16, "Accounts")
							txt.Color = tm.Theme.Color.GrayText1
							return txt.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							txt := tm.Theme.Label(values.MarginPadding16, "cannot")
							txt.Font.Weight = text.SemiBold
							txt.Color = tm.Theme.Color.GrayText1
							inset := layout.Inset{Right: values.MarginPadding2, Left: values.MarginPadding2}
							return inset.Layout(gtx, txt.Layout)
						}),
						layout.Rigid(func(gtx C) D {
							txt := tm.Theme.Label(values.MarginPadding16, "be deleted once created")
							txt.Color = tm.Theme.Color.GrayText1
							return txt.Layout(gtx)
						}),
					)
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

	return tm.modal.Layout(gtx, w)
}

package modal

import (
	"fmt"
	"image/color"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const Info = "info_modal"

type InfoModal struct {
	*app.App

	randomID        string
	modal           decredmaterial.Modal
	enterKeyPressed bool

	dialogIcon *decredmaterial.Icon

	dialogTitle    string
	subtitle       string
	customTemplate []layout.Widget
	customWidget   layout.Widget

	positiveButtonText    string
	positiveButtonClicked func(isChecked bool)
	btnPositve            decredmaterial.Button

	negativeButtonText    string
	negativeButtonClicked func()
	btnNegative           decredmaterial.Button

	checkbox      decredmaterial.CheckBoxStyle
	mustBeChecked bool

	titleAlignment, btnAlignment layout.Direction

	isCancelable bool
}

// func NewInfoModal(app *app.App) *InfoModal {
func NewInfoModal(dt interface{}) *InfoModal {
	app, ok := dt.(*app.App)
	if !ok {
		panic("want app")
	}
	in := &InfoModal{
		App:          app,
		randomID:     fmt.Sprintf("%s-%d", Info, decredmaterial.GenerateRandomNumber()),
		modal:        *app.Theme.ModalFloatTitle(),
		btnPositve:   app.Theme.OutlineButton("Yes"),
		btnNegative:  app.Theme.OutlineButton("No"),
		isCancelable: true,
		btnAlignment: layout.E,
	}

	in.btnPositve.Font.Weight = text.Medium
	in.btnNegative.Font.Weight = text.Medium

	return in
}

func (in *InfoModal) ModalID() string {
	return in.randomID
}

func (in *InfoModal) Show() {
	in.ShowModal(in)
}

func (in *InfoModal) Dismiss() {
	in.DismissModal(in)
}

func (in *InfoModal) OnResume() {}

func (in *InfoModal) OnDismiss() {}

func (in *InfoModal) SetCancelable(min bool) *InfoModal {
	in.isCancelable = min
	return in
}

func (in *InfoModal) SetContentAlignment(title, btn layout.Direction) *InfoModal {
	in.titleAlignment = title
	in.btnAlignment = btn
	return in
}

func (in *InfoModal) Icon(icon *decredmaterial.Icon) *InfoModal {
	in.dialogIcon = icon
	return in
}

func (in *InfoModal) CheckBox(checkbox decredmaterial.CheckBoxStyle, mustBeChecked bool) *InfoModal {
	in.checkbox = checkbox
	in.mustBeChecked = mustBeChecked // determine if the checkbox must be selected to proceed
	return in
}

func (in *InfoModal) Title(title string) *InfoModal {
	in.dialogTitle = title
	return in
}

func (in *InfoModal) Body(subtitle string) *InfoModal {
	in.subtitle = subtitle
	return in
}

func (in *InfoModal) PositiveButton(text string, clicked func(isChecked bool)) *InfoModal {
	in.positiveButtonText = text
	in.positiveButtonClicked = clicked
	return in
}

func (in *InfoModal) PositiveButtonStyle(background, text color.NRGBA) *InfoModal {
	in.btnPositve.Background, in.btnPositve.Color = background, text
	return in
}

func (in *InfoModal) NegativeButton(text string, clicked func()) *InfoModal {
	in.negativeButtonText = text
	in.negativeButtonClicked = clicked
	return in
}

// for backwards compatibilty
func (in *InfoModal) SetupWithTemplate(template string) *InfoModal {
	title := in.dialogTitle
	subtitle := in.subtitle
	var customTemplate []layout.Widget
	switch template {
	case TransactionDetailsInfoTemplate:
		title = "How to copy"
		customTemplate = transactionDetailsInfo(in.Theme)
	case SignMessageInfoTemplate:
		customTemplate = signMessageInfo(in.Theme)
	case VerifyMessageInfoTemplate:
		customTemplate = verifyMessageInfo(in.Theme)
	case PrivacyInfoTemplate:
		title = "How to use the mixer?"
		customTemplate = privacyInfo(in.Theme)
	case SetupMixerInfoTemplate:
		customTemplate = setupMixerInfo(in.Theme)
	case WalletBackupInfoTemplate:
		customTemplate = backupInfo(in.Theme)
	case TicketPriceErrorTemplate:
		customTemplate = ticketPriceErrorInfo(in.MultiWallet(), in.Theme)
	}

	in.dialogTitle = title
	in.subtitle = subtitle
	in.customTemplate = customTemplate
	return in
}

func (in *InfoModal) UseCustomWidget(layout layout.Widget) *InfoModal {
	in.customWidget = layout
	return in
}

// HandleKeyEvent is called when a key is pressed on the current window.
// Satisfies the load.KeyEventHandler interface for receiving key events.
func (in *InfoModal) HandleKeyEvent(evt *key.Event) {
	if (evt.Name == key.NameReturn || evt.Name == key.NameEnter) && evt.State == key.Press {
		in.btnPositve.Click()
		in.RefreshWindow()
	}
}

func (in *InfoModal) Handle() {
	for in.btnPositve.Clicked() {
		in.DismissModal(in)
		isChecked := false
		if in.checkbox.CheckBox != nil {
			isChecked = in.checkbox.CheckBox.Value
		}

		in.positiveButtonClicked(isChecked)
	}

	for in.btnNegative.Clicked() {
		in.DismissModal(in)
		in.negativeButtonClicked()
	}

	if in.modal.BackdropClicked(in.isCancelable) {
		in.Dismiss()
	}

	if in.checkbox.CheckBox != nil {
		if in.mustBeChecked {
			in.btnNegative.SetEnabled(in.checkbox.CheckBox.Value)
		}
	}
}

func (in *InfoModal) Layout(gtx layout.Context) D {
	icon := func(gtx C) D {
		if in.dialogIcon == nil {
			return layout.Dimensions{}
		}

		return layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return in.dialogIcon.Layout(gtx, values.MarginPadding50)
			})
		})
	}

	checkbox := func(gtx C) D {
		if in.checkbox.CheckBox == nil {
			return layout.Dimensions{}
		}

		return layout.Inset{Top: values.MarginPaddingMinus5, Left: values.MarginPaddingMinus5}.Layout(gtx, func(gtx C) D {
			in.checkbox.TextSize = values.TextSize14
			in.checkbox.Color = in.Theme.Color.GrayText1
			in.checkbox.IconColor = in.Theme.Color.Gray2
			if in.checkbox.CheckBox.Value {
				in.checkbox.IconColor = in.Theme.Color.Primary
			}
			return in.checkbox.Layout(gtx)
		})
	}

	subtitle := func(gtx C) D {
		text := in.Theme.Body1(in.subtitle)
		text.Color = in.Theme.Color.GrayText2
		return text.Layout(gtx)
	}

	var w []layout.Widget

	// Every section of the dialog is optional
	if in.dialogIcon != nil {
		w = append(w, icon)
	}

	if in.dialogTitle != "" {
		w = append(w, in.titleLayout())
	}

	if in.subtitle != "" {
		w = append(w, subtitle)
	}

	if in.customTemplate != nil {
		w = append(w, in.customTemplate...)
	}

	if in.checkbox.CheckBox != nil {
		w = append(w, checkbox)
	}

	if in.customWidget != nil {
		w = append(w, in.customWidget)
	}

	if in.negativeButtonText != "" || in.positiveButtonText != "" {
		w = append(w, in.actionButtonsLayout())
	}

	return in.modal.Layout(gtx, w)
}

func (in *InfoModal) titleLayout() layout.Widget {
	return func(gtx C) D {
		t := in.Theme.H6(in.dialogTitle)
		t.Font.Weight = text.SemiBold
		return in.titleAlignment.Layout(gtx, t.Layout)
	}
}

func (in *InfoModal) actionButtonsLayout() layout.Widget {
	return func(gtx C) D {
		return in.btnAlignment.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if in.negativeButtonText == "" {
						return layout.Dimensions{}
					}

					in.btnNegative.Text = in.negativeButtonText
					return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, in.btnNegative.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					if in.positiveButtonText == "" {
						return layout.Dimensions{}
					}

					in.btnPositve.Text = in.positiveButtonText
					return in.btnPositve.Layout(gtx)
				}),
			)
		})
	}
}

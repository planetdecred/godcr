package modal

import (
	"fmt"
	"image/color"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const Info = "info_modal"

type InfoModal struct {
	*load.Load
	randomID        string
	modal           decredmaterial.Modal
	keyEvent        chan *key.Event
	enterKeyPressed bool

	dialogIcon *widget.Icon

	dialogTitle    string
	subtitle       string
	customTemplate []layout.Widget

	positiveButtonText    string
	positiveButtonClicked func()
	btnPositve            decredmaterial.Button

	negativeButtonText    string
	negativeButtonClicked func()
	btnNegative           decredmaterial.Button

	isCancelable bool

	//TODO: neutral button
}

func NewInfoModal(l *load.Load) *InfoModal {
	in := &InfoModal{
		Load:         l,
		randomID:     fmt.Sprintf("%s-%d", Info, generateRandomNumber()),
		modal:        *l.Theme.ModalFloatTitle(),
		btnPositve:   l.Theme.OutlineButton(new(widget.Clickable), "Yes"),
		btnNegative:  l.Theme.OutlineButton(new(widget.Clickable), "No"),
		keyEvent:     l.Receiver.KeyEvents,
		isCancelable: true,
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

func (in *InfoModal) OnResume() {
}

func (in *InfoModal) OnDismiss() {

}

func (in *InfoModal) SetCancelable(min bool) *InfoModal {
	in.isCancelable = min
	return in
}

func (in *InfoModal) Icon(icon *widget.Icon) *InfoModal {
	in.dialogIcon = icon
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

func (in *InfoModal) PositiveButton(text string, clicked func()) *InfoModal {
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
	case SecurityToolsInfoTemplate:
		subtitle = "Various tools that help in different aspects of crypto currency security will be located here."
	case SetupMixerInfoTemplate:
		customTemplate = setupMixerInfo(in.Theme)
	}

	in.dialogTitle = title
	in.subtitle = subtitle
	in.customTemplate = customTemplate
	return in
}

func (in *InfoModal) handleEnterKeypress() {
	// Todo enter button for info modals.
	select {
	case event := <-in.keyEvent:
		if (event.Name == key.NameReturn || event.Name == key.NameEnter) && event.State == key.Press && in.customTemplate != nil {
			in.enterKeyPressed = true
		}
	default:
	}
}

func (in *InfoModal) Handle() {
	for in.btnPositve.Clicked() {
		in.DismissModal(in)
		in.positiveButtonClicked()
	}

	for in.btnNegative.Clicked() {
		in.DismissModal(in)
		in.negativeButtonClicked()
	}

	if in.modal.BackdropClicked(in.isCancelable) {
		in.Dismiss()
	}
}

func (in *InfoModal) Layout(gtx layout.Context) D {
	icon := func(gtx C) D {
		if in.dialogIcon == nil {
			return layout.Dimensions{}
		}

		return layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				in.dialogIcon.Color = in.Theme.Color.DeepBlue
				return in.dialogIcon.Layout(gtx, values.MarginPadding50)
			})
		})
	}

	subtitle := func(gtx C) D {
		text := in.Theme.Body1(in.subtitle)
		text.Color = in.Theme.Color.Gray
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

	if in.negativeButtonText != "" || in.positiveButtonText != "" {
		w = append(w, in.actionButtonsLayout())
	}

	return in.modal.Layout(gtx, w, 850)
}

func (in *InfoModal) titleLayout() layout.Widget {
	return func(gtx C) D {
		t := in.Theme.H6(in.dialogTitle)
		t.Font.Weight = text.Bold
		return t.Layout(gtx)
	}
}

func (in *InfoModal) actionButtonsLayout() layout.Widget {
	return func(gtx C) D {
		return layout.E.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if in.negativeButtonText == "" {
						return layout.Dimensions{}
					}

					in.btnNegative.Text = in.negativeButtonText
					return in.btnNegative.Layout(gtx)
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

package modal

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page"
	"github.com/planetdecred/godcr/ui/values"
)

const Info = "info_modal"

type Common interface {

}

type infoModal struct {
	*load.Load
	randomID string
	modal    decredmaterial.Modal

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

	//TODO: neutral button
}

func NewInfoModal(l *load.Load) *infoModal {
	in := &infoModal{
		Load:  l,
		randomID:    fmt.Sprintf("%s-%d", Info, ui.GenerateRandomNumber()),
		modal:       *l.Theme.ModalFloatTitle(),
		btnPositve:  l.Theme.Button(new(widget.Clickable), "Yes"),
		btnNegative: l.Theme.Button(new(widget.Clickable), "No"),
	}

	in.btnPositve.TextSize, in.btnNegative.TextSize = values.TextSize16, values.TextSize16
	in.btnPositve.Font.Weight, in.btnNegative.Font.Weight = text.Bold, text.Bold

	return in
}

func (in *infoModal) ModalID() string {
	return in.randomID
}

func (in *infoModal) Show() {
	in.ShowModal(in)
}

func (in *infoModal) Dismiss() {
	in.DismissModal(in)
}

func (in *infoModal) OnResume() {
}

func (in *infoModal) OnDismiss() {

}

func (in *infoModal) icon(icon *widget.Icon) *infoModal {
	in.dialogIcon = icon
	return in
}

func (in *infoModal) Title(title string) *infoModal {
	in.dialogTitle = title
	return in
}

func (in *infoModal) Body(subtitle string) *infoModal {
	in.subtitle = subtitle
	return in
}

func (in *infoModal) PositiveButton(text string, clicked func()) *infoModal {
	in.positiveButtonText = text
	in.positiveButtonClicked = clicked
	return in
}

func (in *infoModal) NegativeButton(text string, clicked func()) *infoModal {
	in.negativeButtonText = text
	in.negativeButtonClicked = clicked
	return in
}

// for backwards compatibilty
func (in *infoModal) SetupWithTemplate(template string) *infoModal {
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

func (in *infoModal) Handle() {

	for in.btnPositve.Button.Clicked() {
		in.DismissModal(in)
		in.positiveButtonClicked()
	}

	for in.btnNegative.Button.Clicked() {
		in.DismissModal(in)
		in.negativeButtonClicked()
	}
}

func (in *infoModal) Layout(gtx layout.Context) page.D {
	icon := func(gtx page.C) page.D {
		if in.dialogIcon == nil {
			return layout.Dimensions{}
		}

		return layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding20}.Layout(gtx, func(gtx page.C) page.D {
			return layout.Center.Layout(gtx, func(gtx page.C) page.D {
				in.dialogIcon.Color = in.Theme.Color.DeepBlue
				return in.dialogIcon.Layout(gtx, values.MarginPadding50)
			})
		})
	}

	subtitle := func(gtx page.C) page.D {
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

func (in *infoModal) titleLayout() layout.Widget {
	return func(gtx page.C) page.D {
		t := in.Theme.H6(in.dialogTitle)
		t.Font.Weight = text.Bold
		return t.Layout(gtx)
	}
}

func (in *infoModal) actionButtonsLayout() layout.Widget {
	return func(gtx page.C) page.D {
		return layout.E.Layout(gtx, func(gtx page.C) page.D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx page.C) page.D {
					if in.negativeButtonText == "" {
						return layout.Dimensions{}
					}

					in.btnNegative.Text = in.negativeButtonText
					in.btnNegative.Background = in.Theme.Color.Surface
					in.btnNegative.Color = in.Theme.Color.Primary
					return in.btnNegative.Layout(gtx)
				}),
				layout.Rigid(func(gtx page.C) page.D {
					if in.positiveButtonText == "" {
						return layout.Dimensions{}
					}

					in.btnPositve.Text = in.positiveButtonText
					in.btnPositve.Background, in.btnPositve.Color = in.Theme.Color.Surface, in.Theme.Color.Primary
					return in.btnPositve.Layout(gtx)
				}),
			)
		})
	}
}

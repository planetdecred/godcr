package ui

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const MadalInfo = "info_modal"

type infoModal struct {
	*pageCommon

	modal decredmaterial.Modal

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

func newInfoModal(common *pageCommon) *infoModal {
	in := &infoModal{
		pageCommon:  common,
		modal:       *common.theme.ModalFloatTitle(),
		btnPositve:  common.theme.Button(new(widget.Clickable), "Yes"),
		btnNegative: common.theme.Button(new(widget.Clickable), "No"),
	}

	in.btnPositve.TextSize, in.btnNegative.TextSize = values.TextSize16, values.TextSize16
	in.btnPositve.Font.Weight, in.btnNegative.Font.Weight = text.Bold, text.Bold

	return in
}

func (in *infoModal) modalID() string {
	return MadalInfo + in.dialogTitle // TODO
}

func (in *infoModal) show() {
	in.pageCommon.showModal(in)
}

func (in *infoModal) OnResume() {
}

func (in *infoModal) OnDismiss() {

}

func (in *infoModal) icon(icon *widget.Icon) *infoModal {
	in.dialogIcon = icon
	return in
}

func (in *infoModal) title(title string) *infoModal {
	in.dialogTitle = title
	return in
}

func (in *infoModal) body(subtitle string) *infoModal {
	in.subtitle = subtitle
	return in
}

func (in *infoModal) positiveButton(text string, clicked func()) *infoModal {
	in.positiveButtonText = text
	in.positiveButtonClicked = clicked
	return in
}

func (in *infoModal) negativeButton(text string, clicked func()) *infoModal {
	in.negativeButtonText = text
	in.negativeButtonClicked = clicked
	return in
}

// for backwards compatibilty
func (in *infoModal) setupWithTemplate(template string) *infoModal {
	title := in.dialogTitle
	subtitle := in.subtitle
	var customTemplate []layout.Widget
	switch template {
	case TransactionDetailsInfoTemplate:
		title = "How to copy"
		customTemplate = transactionDetailsInfo(in.theme)
	case SignMessageInfoTemplate:
		customTemplate = signMessageInfo(in.theme)
	case VerifyMessageInfoTemplate:
		customTemplate = verifyMessageInfo(in.theme)
	case PrivacyInfoTemplate:
		title = "How to use the mixer?"
		customTemplate = privacyInfo(in.theme)
	case SecurityToolsInfoTemplate:
		subtitle = "Various tools that help in different aspects of crypto currency security will be located here."
	case SetupMixerInfoTemplate:
		customTemplate = setupMixerInfo(in.theme)
	}

	in.dialogTitle = title
	in.subtitle = subtitle
	in.customTemplate = customTemplate
	return in
}

func (in *infoModal) handle() {

	for in.btnPositve.Button.Clicked() {
		in.dismissModal(in)
		in.positiveButtonClicked()
	}

	for in.btnNegative.Button.Clicked() {
		in.dismissModal(in)
		in.negativeButtonClicked()
	}
}

func (in *infoModal) Layout(gtx layout.Context) D {

	icon := func(gtx C) D {
		if in.dialogIcon == nil {
			return layout.Dimensions{}
		}

		return layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				in.dialogIcon.Color = in.theme.Color.DeepBlue
				return in.dialogIcon.Layout(gtx, values.MarginPadding50)
			})
		})
	}

	title := func(gtx C) D {
		t := in.theme.H6(in.dialogTitle)
		t.Font.Weight = text.Bold
		return t.Layout(gtx)
	}

	subtitle := func(gtx C) D {
		text := in.theme.Body1(in.subtitle)
		text.Color = in.theme.Color.Gray
		return text.Layout(gtx)
	}

	actionButtons := func(gtx C) D {
		return layout.E.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if in.negativeButtonText == "" {
						return layout.Dimensions{}
					}

					in.btnNegative.Text = in.negativeButtonText
					in.btnNegative.Background = in.theme.Color.Surface
					in.btnNegative.Color = in.theme.Color.Primary
					return in.btnNegative.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if in.positiveButtonText == "" {
						return layout.Dimensions{}
					}

					in.btnPositve.Text = in.positiveButtonText
					//TODO
					// if load.template == ConfirmRemoveTemplate {
					// 	m.confirm.Background, m.confirm.Color = th.Color.Surface, th.Color.Danger
					// }
					// if load.template == RescanWalletTemplate {
					in.btnPositve.Background, in.btnPositve.Color = in.theme.Color.Surface, in.theme.Color.Primary
					// }
					// if load.loading {
					// 	th := material.NewTheme(gofont.Collection())
					// 	return layout.Inset{Top: unit.Dp(7)}.Layout(gtx, func(gtx C) D {
					// 		return material.Loader(th).Layout(gtx)
					// 	})
					// } //TODO
					return in.btnPositve.Layout(gtx)
				}),
			)
		})
	}

	var w []layout.Widget

	// Every section of the dialog is optional
	if in.dialogIcon != nil {
		w = append(w, icon)
	}

	if in.dialogTitle != "" {
		w = append(w, title)
	}

	if in.subtitle != "" {
		w = append(w, subtitle)
	}

	if in.customTemplate != nil {
		w = append(w, in.customTemplate...)
	}

	if in.negativeButtonText != "" || in.positiveButtonText != "" {
		w = append(w, actionButtons)
	}

	return in.modal.Layout(gtx, w, 850)

	// w := m.handle(th, load)
	// w = append(title, w...)
	// w = append(w, m.actions(th, load)...)
	// return layout.Dimensions{}
}

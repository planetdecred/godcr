package modal

import (
	"image/color"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

type InfoModal struct {
	*load.Load
	*decredmaterial.Modal

	enterKeyPressed bool

	dialogIcon *decredmaterial.Icon

	dialogTitle    string
	subtitle       string
	customTemplate []layout.Widget
	customWidget   layout.Widget

	positiveButtonText    string
	positiveButtonClicked func(isChecked bool) bool
	btnPositve            decredmaterial.Button

	negativeButtonText    string
	negativeButtonClicked func()
	btnNegative           decredmaterial.Button

	checkbox      decredmaterial.CheckBoxStyle
	mustBeChecked bool

	titleAlignment, btnAlignment layout.Direction
	materialLoader               material.LoaderStyle

	isCancelable bool
	isLoading    bool
}

func NewInfoModal(l *load.Load) *InfoModal {
	return NewInfoModalWithKey(l, "info_modal")
}

func NewInfoModalWithKey(l *load.Load, key string) *InfoModal {

	in := &InfoModal{
		Load:         l,
		Modal:        l.Theme.ModalFloatTitle(key),
		btnPositve:   l.Theme.OutlineButton(values.String(values.StrYes)),
		btnNegative:  l.Theme.OutlineButton(values.String(values.StrNo)),
		isCancelable: true,
		isLoading:    false,
		btnAlignment: layout.E,
	}

	in.btnPositve.Font.Weight = text.Medium
	in.btnNegative.Font.Weight = text.Medium

	in.materialLoader = material.Loader(l.Theme.Base)

	return in
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

func (in *InfoModal) SetLoading(loading bool) {
	in.isLoading = loading
	in.Modal.SetDisabled(loading)
}

func (in *InfoModal) Title(title string) *InfoModal {
	in.dialogTitle = title
	return in
}

func (in *InfoModal) Body(subtitle string) *InfoModal {
	in.subtitle = subtitle
	return in
}

func (in *InfoModal) PositiveButton(text string, clicked func(isChecked bool) bool) *InfoModal {
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
		title = values.String(values.StrHowToCopy)
		customTemplate = transactionDetailsInfo(in.Theme)
	case SignMessageInfoTemplate:
		customTemplate = signMessageInfo(in.Theme)
	case VerifyMessageInfoTemplate:
		customTemplate = verifyMessageInfo(in.Theme)
	case PrivacyInfoTemplate:
		title = values.String(values.StrUseMixer)
		customTemplate = privacyInfo(in.Load)
	case SetupMixerInfoTemplate:
		customTemplate = setupMixerInfo(in.Theme)
	case WalletBackupInfoTemplate:
		customTemplate = backupInfo(in.Theme)
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

// KeysToHandle returns an expression that describes a set of key combinations
// that this modal wishes to capture. The HandleKeyPress() method will only be
// called when any of these key combinations is pressed.
// Satisfies the load.KeyEventHandler interface for receiving key events.
func (in *InfoModal) KeysToHandle() key.Set {
	return decredmaterial.AnyKey(key.NameReturn, key.NameEnter, key.NameEscape)
}

// HandleKeyPress is called when one or more keys are pressed on the current
// window that match any of the key combinations returned by KeysToHandle().
// Satisfies the load.KeyEventHandler interface for receiving key events.
func (in *InfoModal) HandleKeyPress(evt *key.Event) {
	in.btnPositve.Click()
	in.ParentWindow().Reload()
}

func (in *InfoModal) Handle() {
	for in.btnPositve.Clicked() {
		if in.isLoading {
			return
		}
		isChecked := false
		if in.checkbox.CheckBox != nil {
			isChecked = in.checkbox.CheckBox.Value
		}

		if in.positiveButtonClicked(isChecked) {
			in.Dismiss()
		}
	}

	for in.btnNegative.Clicked() {
		if !in.isLoading {
			in.Dismiss()
			in.negativeButtonClicked()
		}
	}

	if in.Modal.BackdropClicked(in.isCancelable) {
		if !in.isLoading {
			in.Dismiss()
		}
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

	return in.Modal.Layout(gtx, w)
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
					if in.negativeButtonText == "" || in.isLoading {
						return layout.Dimensions{}
					}

					in.btnNegative.Text = in.negativeButtonText
					gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding250)
					return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, in.btnNegative.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					if in.isLoading {
						return in.materialLoader.Layout(gtx)
					}

					if in.positiveButtonText == "" {
						return layout.Dimensions{}
					}

					in.btnPositve.Text = in.positiveButtonText
					gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding250)
					return in.btnPositve.Layout(gtx)
				}),
			)
		})
	}
}

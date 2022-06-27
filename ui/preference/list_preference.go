package preference

import (
	"image/color"
	"sort"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

type ListPreferenceModal struct {
	*load.Load
	*decredmaterial.Modal

	optionsRadioGroup *widget.Enum

	items         map[string]string //[key]str-key
	itemKeys      []string
	title         string
	preferenceKey string
	defaultValue  string // str-key
	initialValue  string
	currentValue  string

	positiveButtonText    string
	positiveButtonClicked func()
	btnPositive           decredmaterial.Button

	negativeButtonText    string
	negativeButtonClicked func()
	btnNegative           decredmaterial.Button
}

func NewListPreference(l *load.Load, preferenceKey, defaultValue string, items map[string]string) *ListPreferenceModal {

	// sort keys to keep order when refreshed
	sortedKeys := make([]string, 0)
	for k := range items {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Slice(sortedKeys, func(i int, j int) bool { return sortedKeys[i] < sortedKeys[j] })

	lp := ListPreferenceModal{
		Load:          l,
		preferenceKey: preferenceKey,
		defaultValue:  defaultValue,

		items:    items,
		itemKeys: sortedKeys,

		optionsRadioGroup: new(widget.Enum),
		Modal:             l.Theme.ModalFloatTitle("list_preference"),

		btnPositive: l.Theme.OutlineButton(values.String(values.StrSave)),
		btnNegative: l.Theme.OutlineButton(values.String(values.StrCancel)),
	}

	return &lp
}

func (lp *ListPreferenceModal) OnResume() {
	initialValue := lp.WL.MultiWallet.ReadStringConfigValueForKey(lp.preferenceKey)
	if initialValue == "" {
		initialValue = lp.defaultValue
	}

	lp.initialValue = initialValue
	lp.currentValue = initialValue

	lp.optionsRadioGroup.Value = lp.currentValue
}

func (lp *ListPreferenceModal) OnDismiss() {}

func (lp *ListPreferenceModal) Title(title string) *ListPreferenceModal {
	lp.title = title
	return lp
}

func (lp *ListPreferenceModal) PositiveButton(text string, clicked func()) *ListPreferenceModal {
	lp.positiveButtonText = text
	lp.positiveButtonClicked = clicked
	return lp
}

func (lp *ListPreferenceModal) PositiveButtonStyle(background, text color.NRGBA) *ListPreferenceModal {
	lp.btnPositive.Background, lp.btnPositive.Color = background, text
	return lp
}

func (lp *ListPreferenceModal) NegativeButton(text string, clicked func()) *ListPreferenceModal {
	lp.negativeButtonText = text
	lp.negativeButtonClicked = clicked
	return lp
}

func (lp *ListPreferenceModal) Handle() {

	for lp.optionsRadioGroup.Changed() {
		lp.currentValue = lp.optionsRadioGroup.Value
	}

	for lp.btnNegative.Clicked() {
		lp.Modal.Dismiss()
	}

	for lp.btnPositive.Clicked() {
		lp.currentValue = lp.optionsRadioGroup.Value
		lp.WL.MultiWallet.SaveUserConfigValue(lp.preferenceKey, lp.optionsRadioGroup.Value)
		lp.RefreshTheme(lp.ParentWindow())
		lp.Modal.Dismiss()
	}

	if lp.Modal.BackdropClicked(true) {
		lp.Modal.Dismiss()
	}
}

func (lp *ListPreferenceModal) Layout(gtx layout.Context) layout.Dimensions {
	w := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			txt := lp.Theme.H6(values.String(lp.title))
			txt.Color = lp.Theme.Color.Text
			return txt.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, lp.layoutItems()...)
		},
		func(gtx layout.Context) layout.Dimensions {
			return lp.actionButtonsLayout(gtx)
		},
	}

	return lp.Modal.Layout(gtx, w)
}

func (lp *ListPreferenceModal) layoutItems() []layout.FlexChild {

	items := make([]layout.FlexChild, 0)
	for _, k := range lp.itemKeys {
		radioItem := layout.Rigid(lp.Theme.RadioButton(lp.optionsRadioGroup, k, values.String(lp.items[k]), lp.Theme.Color.DeepBlue, lp.Theme.Color.Primary).Layout)

		items = append(items, radioItem)
	}

	return items
}

func (in *ListPreferenceModal) actionButtonsLayout(gtx layout.Context) layout.Dimensions {
	return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if in.negativeButtonText == "" {
					return layout.Dimensions{}
				}

				in.btnNegative.Text = in.negativeButtonText
				gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding250)
				return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, in.btnNegative.Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if in.positiveButtonText == "" {
					return layout.Dimensions{}
				}

				in.btnPositive.Text = in.positiveButtonText
				gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding250)
				return in.btnPositive.Layout(gtx)
			}),
		)
	})
}

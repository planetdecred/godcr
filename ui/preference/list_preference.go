package preference

import (
	"image/color"
	"sort"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

type ListPreference struct {
	wallet        *wallet.Wallet
	preferenceKey string
	defaultValue  string // str-key
	initialValue  string
	currentValue  string

	theme *decredmaterial.Theme

	IsShowing   bool
	titleStrKey string
	items       map[string]string //[key]str-key

	clickable         *widget.Clickable
	optionsRadioGroup *widget.Enum

	positiveButtonClicked func()
	positiveButtonStrKey  string
	positiveButton        decredmaterial.Button

	negativeButtonClicked func()
	negativeButtonStrKey  string
	negativeButton        decredmaterial.Button
}

func NewListPreference(wallet *wallet.Wallet, theme *decredmaterial.Theme, preferenceKey, defaultValue string, items map[string]string) *ListPreference {
	return &ListPreference{
		wallet:        wallet,
		preferenceKey: preferenceKey,
		defaultValue:  defaultValue,
		theme:         theme,

		items: items,

		IsShowing: false,

		clickable:         new(widget.Clickable),
		optionsRadioGroup: new(widget.Enum),

		positiveButton: theme.Button(new(widget.Clickable), ""),
		negativeButton: theme.Button(new(widget.Clickable), ""),
	}
}

func (lp *ListPreference) Title(titleStrKey string) *ListPreference {
	lp.titleStrKey = titleStrKey
	return lp
}

func (lp *ListPreference) PostiveButton(strkey string, clicked func()) *ListPreference {
	lp.positiveButtonStrKey = strkey
	lp.positiveButtonClicked = clicked

	return lp
}

func (lp *ListPreference) NegativeButton(strkey string, clicked func()) *ListPreference {
	lp.negativeButtonStrKey = strkey
	lp.negativeButtonClicked = clicked

	return lp
}

func (lp *ListPreference) Clickable() *widget.Clickable {
	return lp.clickable
}

func (lp *ListPreference) Handle() {

	for lp.clickable.Clicked() {
		initialValue := lp.wallet.ReadStringConfigValueForKey(lp.preferenceKey)
		if initialValue == "" {
			initialValue = lp.defaultValue
		}

		lp.initialValue = initialValue
		lp.currentValue = initialValue
		lp.IsShowing = true
	}

	for lp.negativeButton.Button.Clicked() {
		lp.setValue(lp.initialValue) // reset value
		lp.IsShowing = false

		lp.negativeButtonClicked()
	}

	for lp.positiveButton.Button.Clicked() {
		lp.setValue(lp.optionsRadioGroup.Value) // set value
		lp.IsShowing = false

		lp.positiveButtonClicked()
	}

	for lp.optionsRadioGroup.Changed() {
		lp.currentValue = lp.optionsRadioGroup.Value
	}
}

func (lp *ListPreference) setValue(value string) {
	lp.wallet.SaveConfigValueForKey(lp.preferenceKey, value)
	values.SetUserLanguage(value)
}

func (lp *ListPreference) Layout(gtx layout.Context, body layout.Dimensions) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return body
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return lp.modal(gtx)
		}),
	)
}

func (lp *ListPreference) modal(gtx layout.Context) layout.Dimensions {
	w := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			txt := lp.theme.H6(values.String(lp.titleStrKey))
			txt.Color = lp.theme.Color.Text
			return txt.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, lp.layoutItems()...)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, lp.layoutButtons()...)
			})
		},
	}

	lp.optionsRadioGroup.Value = lp.currentValue

	return lp.theme.Modal().Layout(gtx, w, 1050)
}

func (lp *ListPreference) layoutItems() []layout.FlexChild {

	// sort keys to keep order when refreshed
	keys := make([]string, 0)
	for k := range lp.items {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i int, j int) bool { return keys[i] < keys[j] })

	items := make([]layout.FlexChild, 0)
	for _, k := range keys {
		radioItem := layout.Rigid(lp.theme.RadioButton(lp.optionsRadioGroup, k, values.String(lp.items[k])).Layout)

		items = append(items, radioItem)
	}

	return items
}

func (lp *ListPreference) layoutButtons() []layout.FlexChild {

	buttonLayout := func(button decredmaterial.Button) layout.FlexChild {

		l := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				button.Background, button.Color = color.NRGBA{}, lp.theme.Color.Primary
				return button.Layout(gtx)
			})
		})

		return l
	}

	buttons := make([]layout.FlexChild, 0)

	if lp.positiveButtonStrKey != "" {
		positiveButtonLayout := buttonLayout(lp.positiveButton)
		lp.positiveButton.Text = values.String(lp.positiveButtonStrKey)

		buttons = append(buttons, positiveButtonLayout)
	}

	if lp.negativeButtonStrKey != "" {
		negativeButtonLayout := buttonLayout(lp.negativeButton)
		lp.negativeButton.Text = values.String(lp.negativeButtonStrKey)

		buttons = append(buttons, negativeButtonLayout)
	}

	return buttons
}

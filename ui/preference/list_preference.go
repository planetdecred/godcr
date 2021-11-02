package preference

import (
	"sort"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
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
	modal *decredmaterial.Modal

	IsShowing   bool
	titleStrKey string
	items       map[string]string //[key]str-key
	itemKeys    []string

	clickable         *decredmaterial.Clickable
	optionsRadioGroup *widget.Enum

	cancelButton decredmaterial.IconButton

	updateButtonClicked func()
}

func NewListPreference(wallet *wallet.Wallet, l *load.Load, preferenceKey, defaultValue string, items map[string]string) *ListPreference {

	// sort keys to keep order when refreshed
	sortedKeys := make([]string, 0)
	for k := range items {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Slice(sortedKeys, func(i int, j int) bool { return sortedKeys[i] < sortedKeys[j] })

	lp := ListPreference{
		wallet:        wallet,
		preferenceKey: preferenceKey,
		defaultValue:  defaultValue,
		theme:         l.Theme,

		items:    items,
		itemKeys: sortedKeys,

		IsShowing: false,

		clickable:         l.Theme.NewClickable(false),
		optionsRadioGroup: new(widget.Enum),
		modal:             l.Theme.Modal(),
	}

	lp.cancelButton = l.Theme.PlainIconButton(l.Icons.ContentClear)
	lp.cancelButton.Color = l.Theme.Color.Gray3
	lp.cancelButton.Size = values.MarginPadding24
	lp.cancelButton.Inset = layout.UniformInset(values.MarginPadding4)

	return &lp
}

func (lp *ListPreference) Title(titleStrKey string) *ListPreference {
	lp.titleStrKey = titleStrKey
	return lp
}

func (lp *ListPreference) UpdateValues(clicked func()) *ListPreference {
	lp.updateButtonClicked = clicked
	return lp
}

func (lp *ListPreference) Clickable() *decredmaterial.Clickable {
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

	for lp.optionsRadioGroup.Changed() {
		lp.currentValue = lp.optionsRadioGroup.Value
		lp.setValue(lp.optionsRadioGroup.Value)
		lp.IsShowing = false
		lp.updateButtonClicked()
	}

	for lp.cancelButton.Button.Clicked() {
		lp.IsShowing = false
	}

	if lp.modal.BackdropClicked(true) {
		lp.IsShowing = false
	}
}

func (lp *ListPreference) setValue(value string) {
	lp.wallet.SaveConfigValueForKey(lp.preferenceKey, value)
}

func (lp *ListPreference) Layout(gtx layout.Context, body layout.Dimensions) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return body
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return lp.printModal(gtx)
		}),
	)
}

func (lp *ListPreference) printModal(gtx layout.Context) layout.Dimensions {
	w := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			txt := lp.theme.H6(values.String(lp.titleStrKey))
			txt.Color = lp.theme.Color.Text
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.
				Layout(gtx, layout.Rigid(txt.Layout), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top: values.MarginPaddingMinus2,
					}.Layout(gtx, lp.cancelButton.Layout)
				}))
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, lp.layoutItems()...)
		},
	}

	lp.optionsRadioGroup.Value = lp.currentValue

	return lp.modal.Layout(gtx, w)
}

func (lp *ListPreference) layoutItems() []layout.FlexChild {

	items := make([]layout.FlexChild, 0)
	for _, k := range lp.itemKeys {
		radioItem := layout.Rigid(lp.theme.RadioButton(lp.optionsRadioGroup, k, values.String(lp.items[k]), lp.theme.Color.DeepBlue).Layout)

		items = append(items, radioItem)
	}

	return items
}

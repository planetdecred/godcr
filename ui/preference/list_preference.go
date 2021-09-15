package preference

import (
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
	itemKeys    []string

	clickable         *decredmaterial.Clickable
	optionsRadioGroup *widget.Enum

	updateButtonClicked func()
}

func NewListPreference(wallet *wallet.Wallet, theme *decredmaterial.Theme, preferenceKey, defaultValue string, items map[string]string) *ListPreference {

	// sort keys to keep order when refreshed
	sortedKeys := make([]string, 0)
	for k := range items {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Slice(sortedKeys, func(i int, j int) bool { return sortedKeys[i] < sortedKeys[j] })

	return &ListPreference{
		wallet:        wallet,
		preferenceKey: preferenceKey,
		defaultValue:  defaultValue,
		theme:         theme,

		items:    items,
		itemKeys: sortedKeys,

		IsShowing: false,

		clickable:         theme.NewClickable(false),
		optionsRadioGroup: new(widget.Enum),
	}
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
	}

	lp.optionsRadioGroup.Value = lp.currentValue

	return lp.theme.Modal().Layout(gtx, w, 1050)
}

func (lp *ListPreference) layoutItems() []layout.FlexChild {

	items := make([]layout.FlexChild, 0)
	for _, k := range lp.itemKeys {
		radioItem := layout.Rigid(lp.theme.RadioButton(lp.optionsRadioGroup, k, values.String(lp.items[k]), lp.theme.Color.DeepBlue).Layout)

		items = append(items, radioItem)
	}

	return items
}

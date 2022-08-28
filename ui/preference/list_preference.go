package preference

import (
	"sort"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/renderers"
	"github.com/planetdecred/godcr/ui/values"
)

type ListPreferenceModal struct {
	*load.Load
	*decredmaterial.Modal

	optionsRadioGroup *widget.Enum

	btnSave   decredmaterial.Button
	btnCancel decredmaterial.Button

	items         map[string]string //[key]str-key
	itemKeys      []string
	title         string
	subtitle      string
	preferenceKey string
	defaultValue  string // str-key
	initialValue  string
	currentValue  string

	updateButtonClicked func()
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

		btnSave:   l.Theme.Button(values.String(values.StrSave)),
		btnCancel: l.Theme.OutlineButton(values.String(values.StrCancel)),

		items:    items,
		itemKeys: sortedKeys,

		optionsRadioGroup: new(widget.Enum),
		Modal:             l.Theme.ModalFloatTitle("list_preference"),
	}

	lp.btnSave.Font.Weight = text.Medium
	lp.btnCancel.Font.Weight = text.Medium

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

func (lp *ListPreferenceModal) Subtitle(subtitle string) *ListPreferenceModal {
	lp.subtitle = subtitle
	return lp
}

func (lp *ListPreferenceModal) UpdateValues(clicked func()) *ListPreferenceModal {
	lp.updateButtonClicked = clicked
	return lp
}

func (lp *ListPreferenceModal) Handle() {
	for lp.btnSave.Button.Clicked() {
		lp.currentValue = lp.optionsRadioGroup.Value
		lp.WL.MultiWallet.SaveUserConfigValue(lp.preferenceKey, lp.optionsRadioGroup.Value)
		lp.updateButtonClicked()
		lp.RefreshTheme(lp.ParentWindow())
		lp.Dismiss()
	}

	for lp.btnCancel.Button.Clicked() {
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
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					text := values.StringF(lp.subtitle, `<span style="text-color: text">`, `<span style="font-weight: bold">`, `</span><span style="text-color: danger">`, `</span></span>`)
					return renderers.RenderHTML(text, lp.Load.Theme).Layout(gtx)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, lp.layoutItems()...)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(lp.btnCancel.Layout),
					layout.Rigid(lp.btnSave.Layout),
				)
			})
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

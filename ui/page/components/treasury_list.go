package components

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

type TreasuryItem struct {
	Policy            dcrlibwallet.TreasuryKeyPolicy
	OptionsRadioGroup *widget.Enum
	VoteChoices       [3]string
	SetChoiceButton   decredmaterial.Button
}

func (t *TreasuryItem) SetVoteChoices(voteChoices [3]string) {
	t.VoteChoices = voteChoices
}

func TreasuryItemWidget(gtx C, l *load.Load, treasuryItem *TreasuryItem) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layoutPiKey(gtx, l, treasuryItem.Policy)
		}),
		layout.Rigid(layoutVoteChoice(l, treasuryItem)),
		layout.Rigid(func(gtx C) D {
			return layoutPolicyVoteAction(gtx, l, treasuryItem)
		}),
	)
}

func layoutPiKey(gtx C, l *load.Load, treasuryKeyPolicy dcrlibwallet.TreasuryKeyPolicy) D {

	statusLabel := l.Theme.Label(values.TextSize14, treasuryKeyPolicy.PiKey)
	backgroundColor := l.Theme.Color.LightBlue

	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			lbl := l.Theme.Label(values.TextSize20, values.String(values.StrPiKey))
			lbl.Font.Weight = text.SemiBold
			return lbl.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return decredmaterial.LinearLayout{
				Background: backgroundColor,
				Width:      decredmaterial.WrapContent,
				Height:     decredmaterial.WrapContent,
				Direction:  layout.Center,
				Alignment:  layout.Middle,
				Border: decredmaterial.Border{
					Color:  backgroundColor,
					Width:  values.MarginPadding1,
					Radius: decredmaterial.Radius(4),
				},
				Padding: layout.Inset{
					Top:    values.MarginPadding3,
					Bottom: values.MarginPadding3,
					Left:   values.MarginPadding8,
					Right:  values.MarginPadding8},
				Margin: layout.Inset{Left: values.MarginPadding10},
			}.Layout2(gtx, statusLabel.Layout)
		}),
	)
}

func layoutVoteChoice(l *load.Load, treasuryItem *TreasuryItem) layout.Widget {
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				lbl := l.Theme.Label(values.TextSize16, values.String(values.StrSetTreasuryPolicy))
				lbl.Font.Weight = text.SemiBold
				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, lbl.Layout)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding10, Left: values.MarginPadding0}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, layoutItems(l, treasuryItem)...)
				})
			}),
		)
	}
}

func layoutItems(l *load.Load, treasuryItem *TreasuryItem) []layout.FlexChild {
	voteChoices := [...]string{
		strings.ToLower(values.String(values.StrYes)),
		strings.ToLower(values.String(values.StrNo)),
		strings.ToLower(values.String(values.StrAbstain)),
	}
	items := make([]layout.FlexChild, 0)
	for _, voteChoice := range voteChoices {
		radioBtn := l.Theme.RadioButton(treasuryItem.OptionsRadioGroup, voteChoice, voteChoice, l.Theme.Color.DeepBlue, l.Theme.Color.Primary)
		radioItem := layout.Rigid(radioBtn.Layout)
		items = append(items, radioItem)
	}

	return items
}

func layoutPolicyVoteAction(gtx C, l *load.Load, treasuryItem *TreasuryItem) D {
	gtx.Constraints.Min.X, gtx.Constraints.Max.X = gtx.Dp(unit.Dp(150)), gtx.Dp(unit.Dp(200))
	treasuryItem.SetChoiceButton.Background = l.Theme.Color.Gray3
	treasuryItem.SetChoiceButton.SetEnabled(false)

	if treasuryItem.OptionsRadioGroup.Value != "" && treasuryItem.OptionsRadioGroup.Value != treasuryItem.Policy.Policy {
		treasuryItem.SetChoiceButton.Background = l.Theme.Color.Primary
		treasuryItem.SetChoiceButton.SetEnabled(true)
	}
	return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, treasuryItem.SetChoiceButton.Layout)
}

func LayoutNoPoliciesFound(gtx C, l *load.Load, syncing bool) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	text := values.String(values.StrNoPoliciesYet)
	if syncing {
		text = values.String(values.StrFetchingPolicies)
	}
	return layout.Center.Layout(gtx, func(gtx C) D {
		lbl := l.Theme.Body1(text)
		lbl.Color = l.Theme.Color.GrayText3
		return layout.Inset{
			Top:    values.MarginPadding10,
			Bottom: values.MarginPadding10,
		}.Layout(gtx, lbl.Layout)
	})
}

func LoadPolicies(l *load.Load, selectedWallet *dcrlibwallet.Wallet, pikey string) []*TreasuryItem {
	policies, err := selectedWallet.TreasuryPolicies(pikey, "")
	if err != nil {
		return nil
	}

	treasuryItems := make([]*TreasuryItem, len(policies))
	for i := 0; i < len(policies); i++ {
		treasuryItems[i] = &TreasuryItem{
			Policy:            *policies[i],
			OptionsRadioGroup: new(widget.Enum),
			SetChoiceButton:   l.Theme.Button(values.String(values.StrUpdatePreference)),
		}

		treasuryItems[i].OptionsRadioGroup.Value = treasuryItems[i].Policy.Policy
	}
	return treasuryItems
}

package components

import (
	"image/color"
	// "strings"

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
	Policy          dcrlibwallet.TreasuryKeyPolicy
	VoteChoices     [3]string
	SetChoiceButton decredmaterial.Button
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

	// var statusLabel decredmaterial.Label
	var backgroundColor color.NRGBA

	statusLabel := l.Theme.Label(values.MarginPadding14, treasuryKeyPolicy.Key)
	backgroundColor = l.Theme.Color.LightBlue

	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			lbl := l.Theme.Label(values.MarginPadding20, "Pi key")
			lbl.Font.Weight = text.SemiBold
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(lbl.Layout),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return decredmaterial.LinearLayout{
				Background: backgroundColor,
				Width:      decredmaterial.WrapContent,
				Height:     decredmaterial.WrapContent,
				Direction:  layout.Center,
				Alignment:  layout.Middle,
				Border:     decredmaterial.Border{Color: backgroundColor, Width: values.MarginPadding1, Radius: decredmaterial.Radius(10)},
				Padding:    layout.Inset{Top: values.MarginPadding3, Bottom: values.MarginPadding3, Left: values.MarginPadding8, Right: values.MarginPadding8},
				Margin:     layout.Inset{Left: values.MarginPadding10},
			}.Layout(gtx,
				layout.Rigid(statusLabel.Layout))
		}),
	)
}

func layoutVoteChoice(l *load.Load, treasuryItem *TreasuryItem) layout.Widget {
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				lbl := l.Theme.Label(values.MarginPadding16, "Set vote choice")
				lbl.Font.Weight = text.SemiBold
				return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, lbl.Layout)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, layoutItems(l, treasuryItem)...)
			}),
		)
	}
}

func layoutItems(l *load.Load, treasuryItem *TreasuryItem) []layout.FlexChild {
	voteChoices := [...]string{"yes", "no", "abstain"}
	initialValue := treasuryItem.Policy.Policy
	println("[][][] initial", initialValue)
	treasuryItem.VoteChoices = voteChoices

	optionsRadioGroup := new(widget.Enum)
	optionsRadioGroup.Value = initialValue
	items := make([]layout.FlexChild, 0)
	for _, voteChoice := range voteChoices {
		println("[][][] inivoteChoicetial", voteChoice)

		radioBtn := l.Theme.RadioButton(optionsRadioGroup, voteChoice, voteChoice, l.Theme.Color.DeepBlue, l.Theme.Color.Primary)
		radioItem := layout.Rigid(radioBtn.Layout)
		items = append(items, radioItem)
	}

	return items
}

func layoutPolicyVoteAction(gtx C, l *load.Load, item *TreasuryItem) D {
	gtx.Constraints.Min.X, gtx.Constraints.Max.X = gtx.Px(unit.Dp(150)), gtx.Px(unit.Dp(200))
	// item.SetChoiceButton.Background = l.Theme.Color.Gray3
	// item.SetChoiceButton.SetEnabled(false)
	// if item.Agenda.Status == dcrlibwallet.AgendaStatusUpcoming.String() || item.Agenda.Status == dcrlibwallet.AgendaStatusInProgress.String() {
	item.SetChoiceButton.Background = l.Theme.Color.Primary
	item.SetChoiceButton.SetEnabled(true)
	// }
	return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, item.SetChoiceButton.Layout)
}

func LayoutNoPoliciesFound(gtx C, l *load.Load, syncing bool) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	text := l.Theme.Body1("No Policies yet")
	text.Color = l.Theme.Color.GrayText3
	if syncing {
		text = l.Theme.Body1("Fetching Policies")
	}
	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Top:    values.MarginPadding10,
			Bottom: values.MarginPadding10,
		}.Layout(gtx, text.Layout)
	})
}

func LoadPolicies(l *load.Load, selectedWallet *dcrlibwallet.Wallet, pikey string) []*TreasuryItem {
	// println("[][][][] wallet", selectedWallet.Name)
	policies, err := selectedWallet.AllTreasuryPolicies(pikey, "")
	println("[][][][] length", len(policies))
	if err != nil {
		return nil
	}
	treasuryItems := make([]*TreasuryItem, len(policies))
	for i := 0; i < len(policies); i++ {
		treasuryItems[i] = &TreasuryItem{
			Policy:          *policies[i],
			SetChoiceButton: l.Theme.Button("Set Choice"),
		}
	}
	return treasuryItems
}

package governance

import (
	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/renderers"
	"github.com/planetdecred/godcr/ui/values"
)

func (pg *GovernancePage) splashScreenLayout(gtx layout.Context) layout.Dimensions {
	return decredmaterial.LinearLayout{
		Orientation: layout.Vertical,
		Width:       decredmaterial.MatchParent,
		Height:      decredmaterial.WrapContent,
		Background:  pg.Theme.Color.Surface,
		Direction:   layout.Center,
		Alignment:   layout.Middle,
		Border:      decredmaterial.Border{Radius: decredmaterial.Radius(14)},
		Padding:     layout.UniformInset(values.MarginPadding24)}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.Icons.GovernanceActiveIcon.LayoutSize(gtx, values.MarginPadding150)
					}),
					layout.Rigid(func(gtx C) D {
						txt := pg.Theme.Label(values.MarginPadding24, "How does Governance Work?")
						txt.Font.Weight = text.SemiBold

						return layout.Inset{
							Top:    values.MarginPadding30,
							Bottom: values.MarginPadding16,
						}.Layout(gtx, txt.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						text := `<span style="text-color: gray">
						The Decred community can participate in proposal discussions 
						for new initiatives and request funding for these initiatives. 
						Decred stakeholders can vote if these proposals should be approved 
						and paid for by the Decred Treasury.

						Would you like to fetch and view the proposals?
					</span>`

						return renderers.RenderHTML(text, pg.Theme).Layout(gtx)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding350)
			return layout.Inset{
				Top:   values.MarginPadding24,
				Right: values.MarginPadding16,
			}.Layout(gtx, pg.fetchProposals.Layout)
		}),
	)
}

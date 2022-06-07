package governance

import (
	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/renderers"
	"github.com/planetdecred/godcr/ui/values"
)

func (pg *Page) initSplashScreenWidgets() {
	_, pg.splashScreenInfoButton = components.SubpageHeaderButtons(pg.Load)
	pg.enableGovernanceBtn = pg.Theme.Button(values.String(values.StrFetchProposals))
}

func (pg *Page) splashScreenLayout(gtx layout.Context) layout.Dimensions {
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
			return layout.Stack{Alignment: layout.NE}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return pg.Theme.Icons.GovernanceActiveIcon.LayoutSize(gtx, values.MarginPadding150)
						}),
						layout.Rigid(func(gtx C) D {
							txt := pg.Theme.Label(values.TextSize24, values.String(values.StrHowGovernanceWork))
							txt.Font.Weight = text.SemiBold

							return layout.Inset{
								Top:    values.MarginPadding30,
								Bottom: values.MarginPadding16,
							}.Layout(gtx, txt.Layout)
						}),
						layout.Rigid(func(gtx C) D {
							text := values.StringF(values.StrGovernanceInfo, `<span style="text-color: gray">`, `<br>`, `</span>`)
							return renderers.RenderHTML(text, pg.Theme).Layout(gtx)
						}),
					)
				}),
				layout.Stacked(pg.splashScreenInfoButton.Layout),
			)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Dp(values.MarginPadding350)
			return layout.Inset{
				Top:   values.MarginPadding24,
				Right: values.MarginPadding16,
			}.Layout(gtx, pg.enableGovernanceBtn.Layout)
		}),
	)
}

func (pg *Page) showInfoModal() {
	info := modal.NewInfoModal(pg.Load).
		Title(values.String(values.StrGovernance)).
		Body(values.String(values.StrProposalInfo)).
		PositiveButton(values.String(values.StrGotIt), func(isChecked bool) bool {
			return true
		})
	pg.ParentWindow().ShowModal(info)
}

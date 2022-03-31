package staking

import (
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

func (pg *Page) walletBalanceLayout(gtx C) D {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						txt := pg.Theme.Label(values.TextSize14, "Balance:")
						txt.Color = pg.Theme.Color.GrayText2
						return txt.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						txt := pg.Theme.Label(values.TextSize14, "")
						txt.Color = pg.Theme.Color.GrayText2

						totalBalance, _, err := components.CalculateTotalWalletsBalance(pg.Load)
						if err == nil {
							txt.Text = totalBalance.String()
						} else {
							txt.Text = err.Error()
						}
						return layout.Inset{
							Left:  values.MarginPadding5,
							Right: values.MarginPadding16,
						}.Layout(gtx, txt.Layout)
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				totalBalance, spendable, _ := components.CalculateTotalWalletsBalance(pg.Load)
				locked := totalBalance - spendable
				items := []decredmaterial.ProgressBarItem{
					{
						Value:   float32(spendable.ToCoin()),
						Color:   pg.Theme.Color.Primary,
						SubText: "Spendable",
					},
					{
						Value:   float32(locked.ToCoin()),
						Color:   pg.Theme.Color.Danger,
						SubText: "Locked",
					},
				}
				return pg.Theme.MultiLayerProgressBar(float32(totalBalance.ToCoin()), items).Layout(gtx)
			}),
		)
	})
}

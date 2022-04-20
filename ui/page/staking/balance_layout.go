package staking

import (
	"image/color"

	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

func (pg *Page) walletBalanceLayout(gtx C) D {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{
			Spacing:   layout.SpaceBetween,
			Alignment: layout.End,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
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
				})
			}),
			layout.Rigid(func(gtx C) D {
				totalBalance, spendable, _ := components.CalculateTotalWalletsBalance(pg.Load)
				locked := totalBalance - spendable
				items := []decredmaterial.ProgressBarItem{
					{
						Value:   float32(spendable.ToCoin()),
						Color:   pg.Theme.Color.Turquoise300,
						SubText: "Spendable",
					},
					{
						Value:   float32(locked.ToCoin()),
						Color:   pg.Theme.Color.Primary,
						SubText: "Locked",
					},
				}

				labelWdg := func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return pg.layoutIconAndText(gtx, "Spendable: ", spendable.String(), items[0].Color)
						}),
						layout.Rigid(func(gtx C) D {
							return pg.layoutIconAndText(gtx, "Locked: ", locked.String(), items[1].Color)
						}),
					)
				}
				return pg.Theme.MultiLayerProgressBar(float32(totalBalance.ToCoin()), items).Layout(gtx, labelWdg)
			}),
		)
	})
}

func (pg *Page) layoutIconAndText(gtx C, title string, val string, col color.NRGBA) D {
	return layout.Inset{Right: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: values.MarginPadding5, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					ic := decredmaterial.NewIcon(pg.Theme.Icons.ImageBrightness1)
					ic.Color = col
					return ic.Layout(gtx, values.MarginPadding8)
				})
			}),
			layout.Rigid(func(gtx C) D {
				txt := pg.Theme.Label(values.TextSize14, title)
				txt.Color = pg.Theme.Color.GrayText2
				return txt.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				txt := pg.Theme.Label(values.TextSize14, val)
				txt.Color = pg.Theme.Color.GrayText2
				return txt.Layout(gtx)
			}),
		)
	})
}

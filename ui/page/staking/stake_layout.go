package staking

import (
	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

func (pg *Page) initStakePriceWidget() *Page {
	pg.stakeBtn = pg.Theme.Button(values.String(values.StrStake))
	pg.autoPurchaseSettings = pg.Theme.NewClickable(false)
	pg.autoPurchase = pg.Theme.Switch()
	return pg
}

func (pg *Page) stakePriceSection(gtx C) D {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding11,
				}.Layout(gtx, func(gtx C) D {
					leftWg := func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								title := pg.Theme.Label(values.TextSize14, values.String(values.StrTicketPrice))
								title.Color = pg.Theme.Color.GrayText2
								return title.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Left:  values.MarginPadding8,
									Right: values.MarginPadding4,
								}.Layout(gtx, pg.Theme.Icons.TimerIcon.Layout12dp)
							}),
							layout.Rigid(func(gtx C) D {
								secs, _ := pg.WL.MultiWallet.NextTicketPriceRemaining()
								txt := pg.Theme.Label(values.TextSize14, nextTicketRemaining(int(secs)))
								txt.Color = pg.Theme.Color.GrayText2

								if pg.WL.MultiWallet.IsSyncing() {
									txt.Text = values.String(values.StrSyncingState)
								}
								return txt.Layout(gtx)
							}),
						)
					}

					rightWg := func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								icon := pg.Theme.Icons.SettingsActiveIcon
								if pg.ticketBuyerWallet.IsAutoTicketsPurchaseActive() {
									icon = pg.Theme.Icons.SettingsInactiveIcon
								}
								return pg.autoPurchaseSettings.Layout(gtx, icon.Layout24dp)
							}),
							layout.Rigid(func(gtx C) D {
								title := pg.Theme.Label(values.TextSize14, values.String(values.StrAutoTicketPurchase))
								title.Color = pg.Theme.Color.GrayText2
								return layout.Inset{
									Left:  values.MarginPadding4,
									Right: values.MarginPadding4,
								}.Layout(gtx, title.Layout)
							}),
							layout.Rigid(pg.autoPurchase.Layout),
						)
					}
					return pg.titleRow(gtx, leftWg, rightWg)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding8,
				}.Layout(gtx, func(gtx C) D {
					ic := pg.Theme.Icons.NewStakeIcon
					return layout.Center.Layout(gtx, ic.Layout48dp)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding16,
				}.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return components.LayoutBalanceSize(gtx, pg.Load, pg.ticketPrice, values.TextSize28)
					})
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Center.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Px(values.MarginPadding150)
					return pg.stakeBtn.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				if pg.WL.MultiWallet.IsSynced() {
					return D{}
				}

				notSynced := pg.Theme.Label(values.TextSize10, values.String(values.StrWalletNotSynced))
				notSynced.Color = pg.Theme.Color.Danger
				return layout.Center.Layout(gtx, notSynced.Layout)
			}),
		)
	})
}

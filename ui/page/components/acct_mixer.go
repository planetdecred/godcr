package components

import (
	"gioui.org/layout"

	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

func mixerInfoStatusTextLayout(gtx C, l *load.Load, mixerActive bool) D {
	txt := l.Theme.H6("Mixer")
	subtxt := l.Theme.Body2("Ready to mix")
	subtxt.Color = l.Theme.Color.Gray
	iconVisibility := false

	if mixerActive {
		txt.Text = "Mixer is running..."
		subtxt.Text = "Keep this app opened"
		iconVisibility = true
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(txt.Layout),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if !iconVisibility {
						return layout.Dimensions{}
					}

					return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, l.Icons.AlertGray.Layout16dp)
				}),
				layout.Rigid(func(gtx C) D {
					return subtxt.Layout(gtx)
				}),
			)
		}),
	)
}

func MixerInfoLayout(gtx C, l *load.Load, wallet *dcrlibwallet.Wallet, overview bool, button layout.Widget) D {
	return l.Theme.Card().Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								ic := l.Icons.Mixer
								return ic.Layout24dp(gtx)
							}),
							layout.Flexed(1, func(gtx C) D {
								return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
									return mixerInfoStatusTextLayout(gtx, l, wallet.IsAccountMixerActive())
								})
							}),
							layout.Rigid(button),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					content := l.Theme.Card()
					content.Color = l.Theme.Color.LightGray
					return content.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
							mixedBalance := "0.00"
							unmixedBalance := "0.00"
							accounts, _ := wallet.GetAccountsRaw()
							for _, acct := range accounts.Acc {
								if acct.Number == wallet.MixedAccountNumber() {
									mixedBalance = dcrutil.Amount(acct.TotalBalance).String()
								} else if acct.Number == wallet.UnmixedAccountNumber() {
									unmixedBalance = dcrutil.Amount(acct.TotalBalance).String()
								}
							}

							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											txt := l.Theme.Label(values.TextSize14, "Unmixed balance")
											txt.Color = l.Theme.Color.Gray
											return txt.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return LayoutBalance(gtx, l, unmixedBalance)
										}),
									)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Center.Layout(gtx, l.Icons.ArrowDownIcon.Layout24dp)
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											t := l.Theme.Label(values.TextSize14, "Mixed balance")
											t.Color = l.Theme.Color.Gray
											return t.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return LayoutBalance(gtx, l, mixedBalance)
										}),
									)
								}),
							)
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					if wallet.IsAccountMixerActive() {
						txt := l.Theme.Body2("The mixer will automatically stop when unmixed balance are fully mixed.")
						txt.Color = l.Theme.Color.Gray
						return txt.Layout(gtx)
					}
					return D{}
				}),
			)
		})
	})
}

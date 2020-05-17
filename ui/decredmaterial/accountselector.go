package decredmaterial

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/wallet"
)

type AccountSelector struct {
	theme        *Theme
	wallets      []walletInfo
	collapsibles []*Collapsible
}

type walletInfo struct {
	info    wallet.InfoShort
	buttons []*widget.Button
}

func (t *Theme) AccountSelector(wallets []wallet.InfoShort) *AccountSelector {
	a := &AccountSelector{
		theme:        t,
		collapsibles: make([]*Collapsible, len(wallets)),
	}

	for i := range wallets {
		a.collapsibles[i] = t.Collapsible()

		walletInfo := walletInfo{
			info:    wallets[i],
			buttons: make([]*widget.Button, len(wallets[i].Accounts)),
		}

		for k := range wallets[i].Accounts {
			walletInfo.buttons[k] = new(widget.Button)
		}

		a.wallets = append(a.wallets, walletInfo)
	}

	return a
}

func (a *AccountSelector) handleSelection(gtx *layout.Context, walletInfo walletInfo, onSelect func(wallet.InfoShort, wallet.Account)) {
	for i := range walletInfo.buttons {
		for walletInfo.buttons[i].Clicked(gtx) {
			onSelect(walletInfo.info, walletInfo.info.Accounts[i])
		}
	}
}

func (a *AccountSelector) Layout(gtx *layout.Context, onSelect func(wallet.InfoShort, wallet.Account)) {
	w := []func(){}
	for index := range a.wallets {
		i := index

		a.handleSelection(gtx, a.wallets[i], onSelect)

		header := func(gtx *layout.Context) {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					inset := layout.Inset{
						Left: unit.Dp(30),
					}
					inset.Layout(gtx, func() {
						a.theme.walletIcon.Layout(gtx, unit.Dp(20))
					})
				}),
				layout.Rigid(func() {
					inset := layout.Inset{
						Left: unit.Dp(30),
					}
					inset.Layout(gtx, func() {
						a.theme.Body1(a.wallets[i].info.Name).Layout(gtx)
					})
				}),
				layout.Rigid(func() {
					gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
					layout.E.Layout(gtx, func() {
						a.theme.H6(dcrutil.Amount(a.wallets[i].info.SpendableBalance).String()).Layout(gtx)
					})
				}),
			)
		}

		bd := func(gtx *layout.Context) {
			list := layout.List{Axis: layout.Vertical}
			list.Layout(gtx, len(a.wallets[i].info.Accounts), func(k int) {
				inset := layout.Inset{
					Top:  unit.Dp(5),
					Left: unit.Dp(60),
				}
				inset.Layout(gtx, func() {
					layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func() {
							a.theme.walletIcon.Layout(gtx, unit.Dp(15))
						}),
						layout.Rigid(func() {
							inset := layout.Inset{
								Left: unit.Dp(30),
							}
							inset.Layout(gtx, func() {
								a.theme.Body1(a.wallets[i].info.Accounts[k].Name).Layout(gtx)
							})
						}),
						layout.Rigid(func() {
							gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
							layout.E.Layout(gtx, func() {
								a.theme.Body1(dcrutil.Amount(a.wallets[i].info.Accounts[k].SpendableBalance).String()).Layout(gtx)
							})
						}),
					)
				})
				pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
				a.wallets[i].buttons[k].Layout(gtx)
			})
		}

		w = append(w, func() {
			a.collapsibles[i].Layout(gtx, header, bd)
		})
	}
	a.theme.Modal(gtx, "Select Account", w)
}

// NOTE: This file along with its content will be deleted when all pages have been updated to the new UI designs

package ui

import (
	"gioui.org/layout"
	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

func (page pageCommon) LayoutWithWallets(gtx layout.Context, body layout.Widget) layout.Dimensions {
	bd := func(gtx C) D {
		if page.walletTabs.ChangeEvent() {
			*page.selectedWallet = page.walletTabs.Selected
			*page.selectedAccount = 0
			page.accountTabs.Selected = 0

			accounts := make([]decredmaterial.TabItem, len(page.info.Wallets[*page.selectedWallet].Accounts))
			for i, account := range page.info.Wallets[*page.selectedWallet].Accounts {
				if account.Name == "imported" {
					continue
				}
				accounts[i] = decredmaterial.TabItem{
					Title: page.info.Wallets[*page.selectedWallet].Accounts[i].Name,
				}
			}
			page.accountTabs.SetTabs(accounts)
		}
		return page.walletTabs.Layout(gtx, body)
	}
	return page.Layout(gtx, bd)
}

func (page pageCommon) LayoutWithAccounts(gtx layout.Context, body layout.Widget) layout.Dimensions {
	if page.accountTabs.ChangeEvent() {
		*page.selectedAccount = page.accountTabs.Selected
	}

	if page.selectedUTXO[page.info.Wallets[*page.selectedWallet].ID] == nil {
		current := page.info.Wallets[*page.selectedWallet]
		account := page.info.Wallets[*page.selectedWallet].Accounts[*page.selectedAccount]
		page.selectedUTXO[current.ID] = make(map[int32]map[string]*wallet.UnspentOutput)
		page.selectedUTXO[current.ID][account.Number] = make(map[string]*wallet.UnspentOutput)
	}

	return page.LayoutWithWallets(gtx, func(gtx C) D {
		return page.accountTabs.Layout(gtx, body)
	})
}

func (page pageCommon) SelectedAccountLayout(gtx layout.Context) layout.Dimensions {
	current := page.info.Wallets[*page.selectedWallet]
	account := page.info.Wallets[*page.selectedWallet].Accounts[*page.selectedAccount]

	selectedDetails := func(gtx C) D {
		return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
								return page.theme.H6(account.Name).Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
								return page.theme.H6(dcrutil.Amount(account.SpendableBalance).String()).Layout(gtx)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
								return page.theme.Body2(current.Name).Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
								return page.theme.Body2(current.Balance).Layout(gtx)
							})
						}),
					)
				}),
			)
		})
	}

	card := page.theme.Card()
	card.Radius = decredmaterial.CornerRadius{
		NE: 0,
		NW: 0,
		SE: 0,
		SW: 0,
	}
	return card.Layout(gtx, selectedDetails)
}

func (page pageCommon) Modal(gtx layout.Context, body layout.Dimensions, modal layout.Dimensions) layout.Dimensions {
	dims := layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return body
		}),
		layout.Stacked(func(gtx C) D {
			return modal
		}),
	)
	return dims
}

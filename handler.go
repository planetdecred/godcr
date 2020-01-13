package gio

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/helper"
	"github.com/raedahgroup/godcr-gio/pages"
	"github.com/raedahgroup/godcr-gio/pages/common"
	"github.com/raedahgroup/godcr-gio/pages/wallet"
	"github.com/raedahgroup/godcr-gio/widgets"
)

type standalonePageHandler interface {
	Render(ctx *layout.Context, refreshWindowFunc func(), changePageFunc func(string))
}

type navPageHandler interface {
	BeforeRender(*common.Syncer, *helper.MultiWallet)
	Render(*layout.Context, func(string))
}

type navPage struct {
	name    string
	label   string
	icon    material.Image
	button  *widgets.ClickableLabel
	handler navPageHandler
}

func getStandalonePages(multiWallet *helper.MultiWallet) map[string]standalonePageHandler {
	return map[string]standalonePageHandler{
		"welcome":       wallet.NewWelcomePage(multiWallet),
		"createwallet":  wallet.NewCreateWalletPage(multiWallet),
		"restorewallet": wallet.NewRestoreWalletPage(multiWallet),
	}
}

func getNavPages() []navPage {
	return []navPage{
		{
			name:    "overview",
			label:   "Overview",
			icon:    helper.OverviewImage,
			button:  widgets.NewClickableLabel("Overview"),
			handler: pages.NewOverviewPage(),
		},
		{
			name:    "transactions",
			label:   "Transactions",
			icon:    helper.TransactionsImage,
			button:  widgets.NewClickableLabel("Transactions"),
			handler: &notImplementedNavPageHandler{"History"},
		},
		{
			name:    "wallets",
			label:   "Wallets",
			icon:    helper.WalletsImage,
			button:  widgets.NewClickableLabel("Wallets"),
			handler: &notImplementedNavPageHandler{"Wallets"},
		},
		{
			name:    "more",
			label:   "More",
			icon:    helper.MoreImage,
			button:  widgets.NewClickableLabel("More"),
			handler: &notImplementedNavPageHandler{"More"},
		},
	}
}

type notImplementedNavPageHandler struct {
	pageTitle string
}

func (_ *notImplementedNavPageHandler) BeforeRender(_ *common.Syncer, _ *helper.MultiWallet) {

}

func (p *notImplementedNavPageHandler) Render(ctx *layout.Context, _ func(string)) {
	widgets.NewLabel("Page not yet implemented").SetSize(7).SetColor(helper.GrayColor).SetWeight(text.Bold).SetStyle(text.Italic).Draw(ctx)
}

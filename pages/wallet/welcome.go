package wallet

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/raedahgroup/godcr-gio/helper"
	"github.com/raedahgroup/godcr-gio/widgets"
)

type (
	WelcomePage struct {
		multiWallet         *helper.MultiWallet
		createWalletButton  *widgets.Button
		restoreWalletButton *widgets.Button
	}
)

func NewWelcomePage(multiWallet *helper.MultiWallet) *WelcomePage {
	return &WelcomePage{
		multiWallet:         multiWallet,
		createWalletButton:  widgets.NewButton("Create a new wallet", widgets.AddIcon).SetBackgroundColor(helper.DecredLightBlueColor),
		restoreWalletButton: widgets.NewButton("Restore an existing wallet", widgets.ReturnIcon).SetBackgroundColor(helper.DecredGreenColor),
	}
}

func (w *WelcomePage) Render(ctx *layout.Context, refreshWindowFunc func(), changePageFunc func(page string)) {
	helper.DrawLogo(ctx)

	inset := layout.Inset{
		Left:  unit.Dp(helper.StandaloneScreenPadding),
		Right: unit.Dp(helper.StandaloneScreenPadding),
	}
	inset.Layout(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(55),
		}
		inset.Layout(ctx, func() {
			widgets.NewLabel("Welcome to", 6).Draw(ctx)
		})

		inset = layout.Inset{
			Top: unit.Dp(85),
		}
		inset.Layout(ctx, func() {
			widgets.NewLabel("Decred Desktop Wallet", 6).Draw(ctx)
		})

		// create button section
		inset = layout.Inset{
			Top: unit.Dp(335),
		}
		inset.Layout(ctx, func() {
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
			ctx.Constraints.Height.Min = 50

			w.createWalletButton.Draw(ctx, func() {
				changePageFunc("createwallet")
			})
		})

		// restore button section
		inset = layout.Inset{
			Top: unit.Dp(395),
		}
		inset.Layout(ctx, func() {
			ctx.Constraints.Width.Min = ctx.Constraints.Width.Max
			ctx.Constraints.Height.Min = 50

			w.restoreWalletButton.Draw(ctx, func() {
				changePageFunc("restorewallet")
			})
		})
	})
}

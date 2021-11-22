package dexclient

import (
	"fmt"
	"strconv"

	"decred.org/dcrdex/client/asset/btc"
	"decred.org/dcrdex/client/asset/dcr"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const DexCreateWalletPageID = "DexCreateWalletPage"

type DexCreateWallet struct {
	*load.Load
	sourceAccountSelector *components.AccountSelector
	backButton            decredmaterial.IconButton
	createNewWallet       decredmaterial.Button
	walletPassword        decredmaterial.Editor
	appPassword           decredmaterial.Editor
	walletInfoWidget      *walletInfoWidget
	isSending             bool
	walletCreated         func()
}

func NewDexCreateWallet(l *load.Load, wallInfo *walletInfoWidget, walletCreated func()) *DexCreateWallet {
	pg := &DexCreateWallet{
		Load:             l,
		walletPassword:   l.Theme.EditorPassword(new(widget.Editor), "Wallet Password"),
		appPassword:      l.Theme.EditorPassword(new(widget.Editor), "App Password"),
		createNewWallet:  l.Theme.Button("Add"),
		walletInfoWidget: wallInfo,
		walletCreated:    walletCreated,
	}

	pg.createNewWallet.TextSize = values.TextSize12
	pg.createNewWallet.Background = l.Theme.Color.Primary
	pg.appPassword.Editor.SingleLine = true
	pg.appPassword.Editor.SetText("")

	pg.sourceAccountSelector = components.NewAccountSelector(pg.Load).
		Title("Select DCR account to use with DEX").
		AccountSelected(func(selectedAccount *dcrlibwallet.Account) {}).
		AccountValidator(func(account *dcrlibwallet.Account) bool {
			// Filter out imported account and mixed.
			wal := pg.WL.MultiWallet.WalletWithID(account.WalletID)
			if account.Number == load.MaxInt32 ||
				account.Number == wal.MixedAccountNumber() {
				return false
			}
			return true
		})
	err := pg.sourceAccountSelector.SelectFirstWalletValidAccount()
	if err != nil {
		pg.Toast.NotifyError(err.Error())
	}

	pg.backButton, _ = components.SubpageHeaderButtons(l)

	return pg
}

func (pg *DexCreateWallet) ID() string {
	return DexCreateWalletPageID
}

func (pg *DexCreateWallet) OnClose() {}

func (pg *DexCreateWallet) OnResume() {}

func (pg *DexCreateWallet) Handle() {
	if pg.createNewWallet.Button.Clicked() {
		if pg.appPassword.Editor.Text() == "" || pg.isSending {
			return
		}

		pg.isSending = true
		go func() {
			defer func() {
				pg.isSending = false
			}()

			coinID := pg.walletInfoWidget.coinID
			coinName := pg.walletInfoWidget.coinName
			if pg.Dexc().HasWallet(int32(coinID)) {
				pg.Toast.NotifyError(fmt.Sprintf("already connected a %s wallet", coinName))
				return
			}

			settings := make(map[string]string)
			var walletType string
			appPass := []byte(pg.appPassword.Editor.Text())
			walletPass := []byte(pg.walletPassword.Editor.Text())

			switch coinID {
			case dcr.BipID:
				selectedAccount := pg.sourceAccountSelector.SelectedAccount()
				settings[dcrlibwallet.DexDcrWalletIDConfigKey] = strconv.Itoa(selectedAccount.WalletID)
				settings["account"] = selectedAccount.Name
				settings["password"] = pg.walletPassword.Editor.Text()
				walletType = dcrlibwallet.CustomDexDcrWalletType
			case btc.BipID:
				walletType = "SPV" // decred.org/dcrdex/client/asset/btc.walletTypeSPV
				walletPass = nil   // Core doesn't accept wallet passwords for dex-managed spv wallets.
			}

			err := pg.Dexc().AddWallet(coinID, walletType, settings, appPass, walletPass)
			if err != nil {
				pg.Toast.NotifyError(err.Error())
				return
			}

			pg.walletCreated()
		}()
	}
}

func (pg *DexCreateWallet) Layout(gtx layout.Context) D {
	body := func(gtx C) D {
		page := components.SubPage{
			Load:       pg.Load,
			Title:      "Add Wallet",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return pg.Theme.Card().Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
									return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return pg.Load.Theme.Label(values.TextSize20, "Add a").Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Left: values.MarginPadding8, Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
												ic := pg.walletInfoWidget.image
												ic.Scale = 0.2
												return pg.walletInfoWidget.image.Layout(gtx)
											})
										}),
										layout.Rigid(func(gtx C) D {
											return pg.Load.Theme.Label(values.TextSize20, fmt.Sprintf("%s Wallet", pg.walletInfoWidget.coinName)).Layout(gtx)
										}),
									)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return pg.Load.Theme.Label(values.TextSize14, "Your wallet is required to pay registration fees.").Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								if pg.walletInfoWidget.coinID == dcr.BipID {
									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
												return pg.sourceAccountSelector.Layout(gtx)
											})
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
												return pg.walletPassword.Layout(gtx)
											})
										}),
									)
								}
								return D{}
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
									return pg.appPassword.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
									return pg.createNewWallet.Layout(gtx)
								})
							}),
						)
					})
				})
			},
		}

		return page.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

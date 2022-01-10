package dexclient

import (
	"context"
	"fmt"
	"image/color"
	"sort"
	"strings"

	"decred.org/dcrdex/client/asset/btc"
	"decred.org/dcrdex/client/asset/dcr"
	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const DexWalletsPageID = "DexWallets"

type DexWalletsPage struct {
	*load.Load
	ctx          context.Context
	ctxCancel    context.CancelFunc
	list         *widget.List
	backButton   decredmaterial.IconButton
	assetWidgets []*assetWidget
}

type assetWidget struct {
	withdrawBtn     *decredmaterial.Clickable
	depositBtn      *decredmaterial.Clickable
	lockBtn         *decredmaterial.Clickable
	unLockBtn       *decredmaterial.Clickable
	createWalletBtn *decredmaterial.Clickable
	asset           *core.SupportedAsset
}

func NewDexWalletsPage(l *load.Load) *DexWalletsPage {
	pg := &DexWalletsPage{
		Load: l,
		list: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}

	pg.backButton, _ = components.SubpageHeaderButtons(pg.Load)

	return pg
}

func (pg *DexWalletsPage) ID() string {
	return DexWalletsPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *DexWalletsPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	go pg.readNotifications()
	pg.assetWidgets = pg.initAssetWidgets()
}

func (pg *DexWalletsPage) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrWallets),
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx layout.Context) layout.Dimensions {
				return pg.Theme.List(pg.list).Layout(gtx, len(pg.assetWidgets), func(gtx C, i int) D {
					return pg.assetRowLayout(gtx, pg.assetWidgets[i])
				})
			},
		}
		return sp.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

func (pg *DexWalletsPage) assetRowLayout(gtx C, assetW *assetWidget) D {
	asset := assetW.asset
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Inset{
				Top:    values.MarginPadding8,
				Bottom: values.MarginPadding8,
				Left:   values.MarginPadding16,
				Right:  values.MarginPadding16,
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										ic := components.CoinImageBySymbol(&pg.Icons, asset.Symbol)
										ic.Scale = 0.2
										return ic.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{
											Left:  values.MarginPadding8,
											Right: values.MarginPadding8,
										}.Layout(gtx, pg.Theme.Body2(asset.Info.Name).Layout)
									}),
									layout.Rigid(pg.Theme.Body2(fmt.Sprintf("(%s)", strings.ToUpper(asset.Symbol))).Layout),
								)
							}),
							layout.Rigid(func(gtx C) D {
								var t string
								if asset.Wallet != nil {
									t = formatAmount(asset.Wallet.Balance.Available, &asset.Info.UnitInfo)
								} else {
									t = "0.00000000"
								}
								return pg.Theme.Label(values.MarginPadding14, t).Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								var t string
								var c color.NRGBA
								switch {
								case asset.Wallet == nil:
									t = strNoWallet
									c = pg.Theme.Color.Danger
								case asset.Wallet.Open:
									t = strReady
									c = pg.Theme.Color.Success
								case asset.Wallet.Running:
									t = strLocked
									c = pg.Theme.Color.Text
								default:
									t = strOff
									c = pg.Theme.Color.Text
								}
								label := pg.Theme.Label(values.TextSize16, t)
								label.Color = c
								return label.Layout(gtx)
							}),
						)
					}),

					layout.Rigid(func(gtx C) D {
						if asset.Wallet != nil {
							if asset.Wallet.Running && !asset.Wallet.Synced {
								return pg.Theme.Body2(walletSyncPercentage(asset.Wallet)).Layout(gtx)
							}
						}
						return D{}
					}),
					layout.Rigid(func(gtx C) D {
						btn := func(b *decredmaterial.Clickable, label string) D {
							return layout.Inset{Left: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
								return widget.Border{
									Color:        pg.Theme.Color.Gray2,
									CornerRadius: values.MarginPadding4,
									Width:        values.MarginPadding1,
								}.Layout(gtx, func(gtx C) D {
									return b.Layout(gtx, func(gtx C) D {
										return layout.Inset{
											Top:    values.MarginPadding4,
											Bottom: values.MarginPadding4,
											Left:   values.MarginPadding8,
											Right:  values.MarginPadding8,
										}.Layout(gtx, pg.Theme.Label(values.MarginPadding12, label).Layout)
									})
								})
							})
						}

						if asset.Wallet != nil {
							if !asset.Wallet.Open {
								return btn(assetW.unLockBtn, values.String(values.StrUnlock))
							}
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return btn(assetW.withdrawBtn, strWithdraw)
								}),
								layout.Rigid(func(gtx C) D {
									return btn(assetW.depositBtn, strDeposit)
								}),
								layout.Rigid(func(gtx C) D {
									return btn(assetW.lockBtn, strLock)
								}),
							)
						}
						return btn(assetW.createWalletBtn, fmt.Sprintf(nStrCreateAWallet, asset.Info.Name))
					}),
				)
			})
		})
	})
}

func (pg *DexWalletsPage) initAssetWidgets() []*assetWidget {
	assetMap := pg.Dexc().Core().SupportedAssets()
	assets := make([]*assetWidget, 0, len(assetMap))
	nowallets := make([]*assetWidget, 0, len(assetMap))
	for _, asset := range assetMap {
		// TODO: only support dcr and btc assets for now, remove this when support more
		if asset.ID != btc.BipID && asset.ID != dcr.BipID {
			continue
		}

		clickable := func() *decredmaterial.Clickable {
			cl := pg.Theme.NewClickable(true)
			cl.Radius = decredmaterial.Radius(0)
			return cl
		}

		aw := &assetWidget{
			withdrawBtn:     clickable(),
			depositBtn:      clickable(),
			createWalletBtn: clickable(),
			lockBtn:         clickable(),
			unLockBtn:       clickable(),
			asset:           asset,
		}

		if asset.Wallet == nil {
			nowallets = append(nowallets, aw)
		} else {
			assets = append(assets, aw)
		}
	}
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].asset.Info.Name < assets[j].asset.Info.Name
	})
	sort.Slice(nowallets, func(i, j int) bool {
		return nowallets[i].asset.Info.Name < nowallets[j].asset.Info.Name
	})
	return append(assets, nowallets...)
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *DexWalletsPage) HandleUserInteractions() {
	for _, assetW := range pg.assetWidgets {
		if assetW.createWalletBtn.Clicked() {
			newCreateWalletModal(pg.Load,
				&walletInfoWidget{
					image:    components.CoinImageBySymbol(&pg.Icons, assetW.asset.Symbol),
					coinName: assetW.asset.Symbol,
					coinID:   assetW.asset.ID,
				},
				"",
				func(_ *createWalletModal) {
					pg.assetWidgets = pg.initAssetWidgets()
				}).Show()
		}

		if assetW.depositBtn.Clicked() {
			newDepositModal(pg.Load, &walletInfoWidget{
				image:    components.CoinImageBySymbol(&pg.Icons, assetW.asset.Symbol),
				coinName: assetW.asset.Symbol,
				coinID:   assetW.asset.ID,
			}, assetW.asset.Wallet.Address).Show()
		}

		if assetW.withdrawBtn.Clicked() {
			newWithdrawModal(pg.Load, &walletInfoWidget{
				image:    components.CoinImageBySymbol(&pg.Icons, assetW.asset.Symbol),
				coinName: assetW.asset.Symbol,
				coinID:   assetW.asset.ID,
			}, assetW.asset).Show()
		}

		if assetW.unLockBtn.Clicked() {
			a := assetW.asset
			modal.NewPasswordModal(pg.Load).
				Title(fmt.Sprintf(nStrUnlockWall, a.Info.Name)).
				Hint(strAppPassword).
				NegativeButton(values.String(values.StrCancel), func() {
				}).
				PositiveButton(strUnLock, func(password string, m *modal.PasswordModal) bool {
					go func() {
						err := pg.Dexc().Core().OpenWallet(a.ID, []byte(password))
						if err != nil {
							pg.Toast.NotifyError(err.Error())
							m.SetLoading(false)
							return
						}

						m.Dismiss()
					}()
					return false
				}).Show()
		}

		if assetW.lockBtn.Clicked() {
			id := assetW.asset.ID
			err := pg.Dexc().Core().CloseWallet(id)
			if err != nil {
				pg.Toast.NotifyError(err.Error())
			}
		}
	}
}

func (pg *DexWalletsPage) readNotifications() {
	ch := pg.Dexc().Core().NotificationFeed()
	for {
		select {
		case n := <-ch:
			switch n.Type() {
			case core.NoteTypeBalance:
				wallB := n.(*core.BalanceNote)
				for i, aw := range pg.assetWidgets {
					if aw.asset.ID == wallB.AssetID {
						pg.assetWidgets[i].asset.Wallet.Balance = wallB.Balance
						pg.RefreshWindow()
						break
					}
				}
			case core.NoteTypeWalletState, core.NoteTypeWalletConfig:
				wallS := n.(*core.WalletStateNote)
				for i, aw := range pg.assetWidgets {
					if aw.asset.ID == wallS.Wallet.AssetID {
						pg.assetWidgets[i].asset.Wallet = wallS.Wallet
						pg.RefreshWindow()
						break
					}
				}
			default:
			}
		case <-pg.ctx.Done():
			return
		}
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *DexWalletsPage) OnNavigatedFrom() {
	pg.ctxCancel()
}

package dexclient

import (
	"fmt"
	"strconv"
	"strings"

	"decred.org/dcrdex/client/asset/btc"
	"decred.org/dcrdex/client/asset/dcr"
	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const dexCreateWalletModalID = "dex_create_wallet_modal"

type createWalletModal struct {
	*load.Load
	sourceAccountSelector *components.AccountSelector
	modal                 *decredmaterial.Modal
	submitBtn             decredmaterial.Button
	cancelBtn             decredmaterial.Button
	walletPassword        decredmaterial.Editor
	appPassword           decredmaterial.Editor
	walletInfoWidget      *walletInfoWidget
	materialLoader        material.LoaderStyle
	isSending             bool
	appPass               string
	isRegisterAction      bool
	walletCreated         func(md *createWalletModal)
}

type walletInfoWidget struct {
	image    *decredmaterial.Image
	coinName string
	coinID   uint32
}

func newCreateWalletModal(l *load.Load, wallInfo *walletInfoWidget, appPass string, walletCreated func(md *createWalletModal)) *createWalletModal {
	md := &createWalletModal{
		Load:             l,
		modal:            l.Theme.ModalFloatTitle(),
		walletPassword:   l.Theme.EditorPassword(new(widget.Editor), strWalletPassword),
		appPassword:      l.Theme.EditorPassword(new(widget.Editor), strAppPassword),
		submitBtn:        l.Theme.Button(strSubmit),
		cancelBtn:        l.Theme.OutlineButton(values.String(values.StrCancel)),
		materialLoader:   material.Loader(material.NewTheme(gofont.Collection())),
		walletInfoWidget: wallInfo,
		walletCreated:    walletCreated,
		appPass:          appPass,
	}

	md.appPassword.Editor.SingleLine = true
	md.appPassword.Editor.SetText("")

	md.sourceAccountSelector = components.NewAccountSelector(md.Load).
		Title(strSellectWallet).
		AccountSelected(func(selectedAccount *dcrlibwallet.Account) {}).
		AccountValidator(func(account *dcrlibwallet.Account) bool {
			// Filter out imported account and mixed.
			wal := md.WL.MultiWallet.WalletWithID(account.WalletID)
			if account.Number == load.MaxInt32 ||
				account.Number == wal.MixedAccountNumber() {
				return false
			}
			return true
		})
	err := md.sourceAccountSelector.SelectFirstWalletValidAccount()
	if err != nil {
		md.Toast.NotifyError(err.Error())
	}

	return md
}

func (md *createWalletModal) ModalID() string {
	return dexCreateWalletModalID
}

func (md *createWalletModal) Show() {
	md.ShowModal(md)
}

func (md *createWalletModal) Dismiss() {
	md.DismissModal(md)
}

func (md *createWalletModal) OnDismiss() {
}

func (md *createWalletModal) OnResume() {
}

func (md *createWalletModal) SetRegisterAction(registerAction bool) *createWalletModal {
	md.isRegisterAction = registerAction
	return md
}

func (md *createWalletModal) Handle() {
	if md.cancelBtn.Button.Clicked() && !md.isSending {
		md.Dismiss()
	}

	if md.submitBtn.Button.Clicked() {
		if md.isSending {
			return
		}

		md.isSending = true
		go func() {
			defer func() {
				md.isSending = false
			}()

			coinID := md.walletInfoWidget.coinID
			coinName := md.walletInfoWidget.coinName
			if md.Dexc().HasWallet(int32(coinID)) {
				md.Toast.NotifyError(fmt.Sprintf(nStrAlreadyConnectWallet, coinName))
				return
			}

			settings := make(map[string]string)
			var walletType string
			appPass := md.appPass
			if appPass == "" {
				appPass = md.appPassword.Editor.Text()
			}
			walletPass := []byte(md.walletPassword.Editor.Text())

			switch coinID {
			case dcr.BipID:
				selectedAccount := md.sourceAccountSelector.SelectedAccount()
				settings[dcrlibwallet.DexDcrWalletIDConfigKey] = strconv.Itoa(selectedAccount.WalletID)
				settings["account"] = selectedAccount.Name
				settings["password"] = md.walletPassword.Editor.Text()
				walletType = dcrlibwallet.CustomDexDcrWalletType
			case btc.BipID:
				walletType = "SPV" // decred.org/dcrdex/client/asset/btc.walletTypeSPV
				walletPass = nil   // Core doesn't accept wallet passwords for dex-managed spv wallets.
			}

			err := md.Dexc().AddWallet(coinID, walletType, settings, []byte(appPass), walletPass)
			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
			}

			md.Dismiss()
			md.walletCreated(md)
		}()
	}
}

func (md *createWalletModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return md.Load.Theme.Label(values.TextSize20, strAddA).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding8, Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						ic := md.walletInfoWidget.image
						ic.Scale = 0.2
						return md.walletInfoWidget.image.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return md.Load.Theme.Label(values.TextSize20, fmt.Sprintf(nStrNameWallet, md.walletInfoWidget.coinName)).Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if !md.isRegisterAction {
						return D{}
					}
					return md.Load.Theme.Label(values.TextSize14, strRequireWalletPayFee).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if md.walletInfoWidget.coinID == dcr.BipID {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
									return md.sourceAccountSelector.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
									return md.walletPassword.Layout(gtx)
								})
							}),
						)
					}
					return D{}
				}),
				layout.Rigid(func(gtx C) D {
					if md.appPass != "" {
						return D{}
					}
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return md.appPassword.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if md.isSending {
							return D{}
						}
						return layout.Inset{
							Right:  values.MarginPadding4,
							Bottom: values.MarginPadding15,
						}.Layout(gtx, md.cancelBtn.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						if md.isSending {
							return layout.Inset{
								Top:    values.MarginPadding10,
								Bottom: values.MarginPadding15,
							}.Layout(gtx, md.materialLoader.Layout)
						}
						return md.submitBtn.Layout(gtx)
					}),
				)
			})
		},
	}

	return md.modal.Layout(gtx, w)
}

const DexassetSelectorModalID = "dex_asset_selector_modal"

type assetSelectorModal struct {
	*load.Load
	*core.Exchange
	listWallet    *widget.List
	listWalletWdg []*assetSelectWidget
	modal         *decredmaterial.Modal
	cancelBtn     decredmaterial.Button
	assetSelected func(*core.SupportedAsset)
}

type assetSelectWidget struct {
	selectBtn *decredmaterial.Clickable
	feeAsset  *core.FeeAsset
	asset     *core.SupportedAsset
}

func newAssetSelectorModal(l *load.Load, d *core.Exchange) *assetSelectorModal {
	amd := &assetSelectorModal{
		Load:      l,
		Exchange:  d,
		modal:     l.Theme.ModalFloatTitle(),
		cancelBtn: l.Theme.OutlineButton(values.String(values.StrCancel)),
		listWallet: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}

	return amd
}

func (amd *assetSelectorModal) ModalID() string {
	return DexassetSelectorModalID
}

func (amd *assetSelectorModal) Show() {
	amd.ShowModal(amd)
}

func (amd *assetSelectorModal) Dismiss() {
	amd.DismissModal(amd)
}

func (amd *assetSelectorModal) OnDismiss() {
}

func (amd *assetSelectorModal) AssetSelected(callback func(*core.SupportedAsset)) *assetSelectorModal {
	amd.assetSelected = callback
	return amd
}

func (amd *assetSelectorModal) OnResume() {
	amd.listWalletWdg = make([]*assetSelectWidget, 0)
	assetMap := amd.Dexc().Core().SupportedAssets()
	for _, feeAsset := range amd.RegFees {
		asset, found := assetMap[feeAsset.ID]
		if !found {
			continue
		}
		cl := amd.Theme.NewClickable(true)
		cl.Radius = decredmaterial.Radius(0)
		amd.listWalletWdg = append(amd.listWalletWdg, &assetSelectWidget{
			selectBtn: cl,
			feeAsset:  feeAsset,
			asset:     asset,
		})
	}
}

func (amd *assetSelectorModal) Handle() {
	if amd.cancelBtn.Button.Clicked() {
		amd.Dismiss()
	}

	for _, walWdg := range amd.listWalletWdg {
		if walWdg.selectBtn.Clicked() {
			amd.assetSelected(walWdg.asset)
			amd.Dismiss()
			return
		}
	}
}

func (amd *assetSelectorModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		amd.Load.Theme.Label(values.TextSize20, strConfirmSelectAssetPayFee).Layout,
		amd.assetsInfoLayout(),
		amd.Theme.Separator().Layout,
		amd.marketSummaryLayout(),
		func(gtx C) D {
			return layout.E.Layout(gtx, amd.cancelBtn.Layout)
		},
	}

	return amd.modal.Layout(gtx, w)
}

func (amd *assetSelectorModal) assetsInfoLayout() layout.Widget {
	return func(gtx C) D {
		return amd.Theme.List(amd.listWallet).Layout(gtx, len(amd.listWalletWdg), func(gtx C, i int) D {
			wallWdg := amd.listWalletWdg[i]
			return wallWdg.selectBtn.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Inset{
					Top:    values.MarginPadding5,
					Bottom: values.MarginPadding5,
					Left:   values.MarginPadding8,
					Right:  values.MarginPadding8,
				}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							ic := components.CoinImageBySymbol(&amd.Icons, wallWdg.asset.Symbol)
							ic.Scale = .35
							return ic.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									convertedAmount := formatAmount(wallWdg.feeAsset.Amt, &wallWdg.asset.Info.UnitInfo)
									convertedAmountSymbol := fmt.Sprintf("%s %s", convertedAmount, wallWdg.asset.Info.UnitInfo.Conventional.Unit)
									return amd.Theme.Label(values.TextSize16, strings.ToUpper(convertedAmountSymbol)).Layout(gtx)
								}),
								layout.Rigid(amd.Theme.Label(values.TextSize12, fmt.Sprintf(nStrNumberConfirmations, wallWdg.feeAsset.Confs)).Layout),
								layout.Rigid(func(gtx C) D {
									walletReady := amd.Theme.Label(values.TextSize12, strSetupNeeded)
									walletReady.Color = amd.Theme.Color.Yellow
									wallet := wallWdg.asset.Wallet
									if wallet != nil {
										walletReady.Text = strWalletReady
										walletReady.Color = amd.Theme.Color.Success
									}
									return walletReady.Layout(gtx)
								}),
							)
						}),
						layout.Rigid(amd.marketInfoLayout(wallWdg.feeAsset)),
						layout.Rigid(amd.lotSizeLayout()),
					)
				})
			})
		})
	}
}

func (amd *assetSelectorModal) marketInfoLayout(feeAsset *core.FeeAsset) layout.Widget {
	return func(gtx C) D {
		childrens := []layout.FlexChild{
			layout.Rigid(func(gtx C) D {
				return amd.Theme.Label(values.TextSize12, strMarket).Layout(gtx)
			}),
		}

		for _, mkt := range amd.Markets {
			if !supportedMarket(mkt) {
				continue
			}
			if mkt.BaseID != feeAsset.ID && mkt.QuoteID != feeAsset.ID {
				continue
			}

			mkt := mkt
			var ic *decredmaterial.Image
			if excludeBase := feeAsset.ID == mkt.BaseID; excludeBase {
				ic = components.CoinImageBySymbol(&amd.Icons, dex.BipIDSymbol(mkt.QuoteID))
			} else {
				ic = components.CoinImageBySymbol(&amd.Icons, dex.BipIDSymbol(mkt.BaseID))
			}
			ic.Scale = .11
			childrens = append(childrens, layout.Rigid(func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, ic.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						txt := fmt.Sprintf("%s-%s", mkt.BaseSymbol, mkt.QuoteSymbol)
						return amd.Theme.Label(values.TextSize10, strings.ToUpper(txt)).Layout(gtx)
					}),
				)
			}))
		}

		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childrens...)
	}
}

func (amd *assetSelectorModal) marketSummaryLayout() layout.Widget {
	return func(gtx C) D {
		return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(amd.Theme.Label(values.TextSize14, strAllMarketAt).Layout),
					layout.Rigid(amd.Theme.Label(values.TextSize12, amd.Host).Layout),
				)
			}),
			layout.Rigid(func(gtx C) D {
				childrens := []layout.FlexChild{
					layout.Rigid(amd.Theme.Label(values.TextSize12, strMarket).Layout),
				}
				for _, mkt := range amd.Markets {
					if !supportedMarket(mkt) {
						continue
					}
					mkt := mkt
					baseIc := components.CoinImageBySymbol(&amd.Icons, dex.BipIDSymbol(mkt.BaseID))
					quoteIc := components.CoinImageBySymbol(&amd.Icons, dex.BipIDSymbol(mkt.QuoteID))
					baseIc.Scale, quoteIc.Scale = .11, .11
					childrens = append(childrens, layout.Rigid(func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, baseIc.Layout)
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, quoteIc.Layout)
									}),
								)
							}),
							layout.Rigid(func(gtx C) D {
								txt := fmt.Sprintf("%s-%s", mkt.BaseSymbol, mkt.QuoteSymbol)
								return amd.Theme.Label(values.TextSize10, strings.ToUpper(txt)).Layout(gtx)
							}),
						)
					}))
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childrens...)
			}),
			layout.Rigid(amd.lotSizeLayout()),
		)
	}
}

func (amd *assetSelectorModal) lotSizeLayout() layout.Widget {
	return func(gtx C) D {
		childrens := []layout.FlexChild{
			layout.Rigid(amd.Theme.Label(values.TextSize12, strLotSize).Layout),
		}

		for _, mkt := range amd.Markets {
			if !supportedMarket(mkt) {
				continue
			}
			asset := amd.Dexc().Core().SupportedAssets()[mkt.BaseID]
			if asset == nil {
				continue
			}
			baseUnitInfo := asset.Info.UnitInfo
			txt := fmt.Sprintf("%s %s", formatAmount(mkt.LotSize, &baseUnitInfo), baseUnitInfo.Conventional.Unit)
			childrens = append(childrens, layout.Rigid(amd.Theme.Label(values.TextSize10, strings.ToUpper(txt)).Layout))
		}

		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childrens...)
	}
}

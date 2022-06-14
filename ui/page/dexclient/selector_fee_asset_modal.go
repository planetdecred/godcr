package dexclient

import (
	"fmt"
	"strings"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex"
	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const dexAssetSelectorModalID = "dex_asset_selector_modal"

type assetSelectorModal struct {
	*load.Load
	dexServer             *core.Exchange
	listFeeAssetClickable map[uint32]*decredmaterial.Clickable
	modal                 *decredmaterial.Modal
	cancelBtn             decredmaterial.Button
	onAssetSelected       func(*core.SupportedAsset)
}

func newFeeAssetSelectorModal(l *load.Load, d *core.Exchange) *assetSelectorModal {
	amd := &assetSelectorModal{
		Load:      l,
		dexServer: d,
		modal:     l.Theme.ModalFloatTitle(),
		cancelBtn: l.Theme.OutlineButton(values.String(values.StrCancel)),
	}

	return amd
}

func (amd *assetSelectorModal) ModalID() string {
	return dexAssetSelectorModalID
}

func (amd *assetSelectorModal) Show() {
	amd.ShowModal(amd)
}

func (amd *assetSelectorModal) Dismiss() {
	amd.DismissModal(amd)
}

func (amd *assetSelectorModal) OnDismiss() {}

func (amd *assetSelectorModal) OnAssetSelected(callback func(*core.SupportedAsset)) *assetSelectorModal {
	amd.onAssetSelected = callback
	return amd
}

func (amd *assetSelectorModal) OnResume() {
	listFeeAsset := sortFeeAsset(amd.dexServer.RegFees)
	amd.listFeeAssetClickable = make(map[uint32]*decredmaterial.Clickable, len(listFeeAsset))
	assetMap := amd.Dexc().Core().SupportedAssets()

	for _, feeAsset := range listFeeAsset {
		_, found := assetMap[feeAsset.ID]
		if !found {
			continue
		}
		cl := amd.Theme.NewClickable(true)
		cl.Radius = decredmaterial.Radius(0)
		amd.listFeeAssetClickable[feeAsset.ID] = cl
	}
}

func (amd *assetSelectorModal) Handle() {
	if amd.cancelBtn.Button.Clicked() {
		amd.Dismiss()
	}

	for assetID, cl := range amd.listFeeAssetClickable {
		if cl.Clicked() {
			amd.onAssetSelected(amd.Dexc().Core().SupportedAssets()[assetID])
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
		listFeeAsset := sortFeeAsset(amd.dexServer.RegFees)
		var childrens = make([]layout.FlexChild, 0, len(listFeeAsset))

		for i := 0; i < len(listFeeAsset); i++ {
			feeAsset := listFeeAsset[i]
			asset := amd.Dexc().Core().SupportedAssets()[feeAsset.ID]

			childrens = append(childrens, layout.Rigid(func(gtx C) D {
				return amd.listFeeAssetClickable[feeAsset.ID].Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Inset{
						Top:    values.MarginPadding5,
						Bottom: values.MarginPadding5,
						Left:   values.MarginPadding8,
						Right:  values.MarginPadding8,
					}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								ic := components.CoinImageBySymbol(&amd.Icons, asset.Symbol)
								ic.Scale = .35
								return ic.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										convertedAmount := formatAmount(feeAsset.Amt, &asset.Info.UnitInfo)
										convertedAmountSymbol := fmt.Sprintf("%s %s", convertedAmount, asset.Info.UnitInfo.Conventional.Unit)
										return amd.Theme.Label(values.TextSize16, strings.ToUpper(convertedAmountSymbol)).Layout(gtx)
									}),
									layout.Rigid(amd.Theme.Label(values.TextSize12, fmt.Sprintf(nStrNumberConfirmations, feeAsset.Confs)).Layout),
									layout.Rigid(func(gtx C) D {
										walletReady := amd.Theme.Label(values.TextSize12, strSetupNeeded)
										walletReady.Color = amd.Theme.Color.Yellow
										if asset.Wallet != nil {
											walletReady.Text = strWalletReady
											walletReady.Color = amd.Theme.Color.Success
										}
										return walletReady.Layout(gtx)
									}),
								)
							}),
							layout.Rigid(amd.marketInfoLayout(feeAsset)),
							layout.Rigid(amd.lotSizeLayout()),
						)
					})
				})
			}))
		}

		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childrens...)
	}
}

func (amd *assetSelectorModal) marketInfoLayout(feeAsset *core.FeeAsset) layout.Widget {
	return func(gtx C) D {
		childrens := []layout.FlexChild{
			layout.Rigid(amd.Theme.Label(values.TextSize12, strMarket).Layout),
		}

		for _, mkt := range amd.dexServer.Markets {
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
					layout.Rigid(amd.Theme.Label(values.TextSize12, amd.dexServer.Host).Layout),
				)
			}),
			layout.Rigid(func(gtx C) D {
				childrens := []layout.FlexChild{
					layout.Rigid(amd.Theme.Label(values.TextSize12, strMarket).Layout),
				}
				for _, mkt := range amd.dexServer.Markets {
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

		for _, mkt := range amd.dexServer.Markets {
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

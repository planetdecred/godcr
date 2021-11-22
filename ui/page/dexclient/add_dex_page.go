package dexclient

import (
	"fmt"
	"strconv"
	"strings"

	"decred.org/dcrdex/client/asset"
	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const AddDexPageID = "AddDex"

const testDexHost = "dex-test.ssgen.io:7232"

type AddDexPage struct {
	*load.Load
	addDexServer     decredmaterial.Button
	dexServerAddress decredmaterial.Editor
	isSending        bool
	cert             decredmaterial.Editor
}

type walletInfoWidget struct {
	image    *decredmaterial.Image
	coinName string
	coinID   uint32
}

func NewAddDexPage(l *load.Load) *AddDexPage {
	pg := &AddDexPage{
		Load:             l,
		dexServerAddress: l.Theme.Editor(new(widget.Editor), "DEX Address"),
		addDexServer:     l.Theme.Button("Submit"),
		cert:             l.Theme.Editor(new(widget.Editor), "Cert content"),
	}

	pg.addDexServer.TextSize = values.TextSize12
	pg.addDexServer.Background = l.Theme.Color.Primary

	pg.dexServerAddress.Editor.SetText(testDexHost)
	pg.dexServerAddress.Editor.SingleLine = true

	return pg
}

func (pg *AddDexPage) ID() string {
	return AddDexPageID
}

func (pg *AddDexPage) OnClose() {}

func (pg *AddDexPage) OnResume() {
}

func (pg *AddDexPage) Handle() {
	if pg.addDexServer.Button.Clicked() {
		if pg.dexServerAddress.Editor.Text() == "" || pg.isSending {
			return
		}

		pg.isSending = true
		go func() {
			cert := []byte(pg.cert.Editor.Text())
			dex, err := pg.Dexc().DEXServerInfo(pg.dexServerAddress.Editor.Text(), cert)
			pg.isSending = false
			if err != nil {
				pg.Toast.NotifyError(err.Error())
				return
			}

			// Ensure a wallet is connected that can be used to pay the fees.
			// TODO: This automatically selects the dcr wallet if the DEX
			// supports it for fee payment, otherwise picks a random wallet
			// to use for fee payment. Should instead update the modal UI
			// to show the options and let the user choose which wallet to
			// set up and use for fee payment.
			feeAssetName := "dcr"
			feeAsset := dex.RegFees[feeAssetName]
			if feeAsset == nil {
				for feeAssetName, feeAsset = range dex.RegFees {
					break
				}
			}

			completeRegistration := func() {
				modal.NewPasswordModal(pg.Load).
					Title("Confirm Registration").
					Hint("App password").
					Description(confirmRegisterModalDesc(dex, feeAssetName)).
					NegativeButton(values.String(values.StrCancel), func() {}).
					PositiveButton("Register", func(password string, pm *modal.PasswordModal) bool {
						go func() {
							_, err := pg.Dexc().RegisterWithDEXServer(dex.Host,
								cert,
								int64(dex.Fee.Amt),
								int32(dex.Fee.ID),
								[]byte(password))
							if err != nil {
								pm.SetError(err.Error())
								pm.SetLoading(false)
								return
							}
							pg.ChangeFragment(NewMarketPage(pg.Load))
							pm.Dismiss()
						}()

						return false
					}).Show()
			}

			if pg.Dexc().HasWallet(int32(feeAsset.ID)) {
				completeRegistration()
			} else {
				wallInfoWdg := &walletInfoWidget{
					image:    components.CoinImageBySymbol(&pg.Load.Icons, feeAssetName),
					coinName: feeAssetName,
					coinID:   feeAsset.ID,
				}
				pg.ChangeFragment(NewDexCreateWallet(pg.Load, wallInfoWdg, completeRegistration))
			}
		}()
	}
}

func (pg *AddDexPage) Layout(gtx layout.Context) D {
	body := func(gtx C) D {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
							return pg.Load.Theme.Label(values.TextSize20, "Add a dex").Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
							return pg.dexServerAddress.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
							gtx.Constraints.Max.Y = 400
							return pg.cert.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
							return pg.addDexServer.Layout(gtx)
						})
					}),
				)
			})
		})
	}

	return components.UniformPadding(gtx, body)
}

func confirmRegisterModalDesc(dex *core.Exchange, selectedFeeAsset string) string {
	feeAsset := dex.RegFees[selectedFeeAsset]
	feeAmt := formatAmount(feeAsset.ID, selectedFeeAsset, feeAsset.Amt)
	txt := fmt.Sprintf("Enter your app password to confirm DEX registration. When you submit this form, %s will be spent from your wallet to pay registration fees.", feeAmt)
	markets := make([]string, 0, len(dex.Markets))
	for _, mkt := range dex.Markets {
		lotSize := formatAmount(mkt.BaseID, mkt.BaseSymbol, mkt.LotSize)
		markets = append(markets, fmt.Sprintf("Base: %s\tQuote: %s\tLot Size: %s", strings.ToUpper(mkt.BaseSymbol), strings.ToUpper(mkt.QuoteSymbol), lotSize))
	}
	return fmt.Sprintf("%s\n\nThis DEX supports the following markets. All trades are in multiples of each market's lot size.\n\n%s", txt, strings.Join(markets, "\n"))
}

func formatAmount(assetID uint32, assetName string, amount uint64) string {
	assetInfo, err := asset.Info(assetID)
	if err != nil {
		return fmt.Sprintf("%d [%s units]", amount, assetName)
	}
	unitInfo := assetInfo.UnitInfo
	convertedLotSize := float64(amount) / float64(unitInfo.Conventional.ConversionFactor)
	return fmt.Sprintf("%s %s", strconv.FormatFloat(convertedLotSize, 'f', -1, 64), unitInfo.Conventional.Unit)
}

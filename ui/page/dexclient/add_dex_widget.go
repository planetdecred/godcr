package dexclient

import (
	"fmt"
	"strconv"
	"strings"

	"decred.org/dcrdex/client/asset"
	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const testDexHost = "dex-test.ssgen.io:7232"

type addDexWidget struct {
	*load.Load
	addDexServer     decredmaterial.Button
	dexServerAddress decredmaterial.Editor
	isSending        bool
	cert             decredmaterial.Editor
	createWalletWdg  *dexCreateWalletWidget
}

func newAddDexWidget(l *load.Load) *addDexWidget {
	dwg := &addDexWidget{
		Load:             l,
		dexServerAddress: l.Theme.Editor(new(widget.Editor), "DEX Address"),
		addDexServer:     l.Theme.Button("Submit"),
		cert:             l.Theme.Editor(new(widget.Editor), "Cert content"),
	}

	dwg.addDexServer.TextSize = values.TextSize12
	dwg.addDexServer.Background = l.Theme.Color.Primary

	dwg.dexServerAddress.Editor.SingleLine = true
	if l.WL.MultiWallet.NetType() == dcrlibwallet.Testnet3 {
		dwg.dexServerAddress.Editor.SetText(testDexHost)
	}

	return dwg
}

func (dwg *addDexWidget) handle() {
	if dwg.addDexServer.Button.Clicked() {
		if dwg.dexServerAddress.Editor.Text() == "" || dwg.isSending {
			return
		}

		dwg.isSending = true
		go func() {
			cert := []byte(dwg.cert.Editor.Text())
			dex, err := dwg.Dexc().DEXServerInfo(dwg.dexServerAddress.Editor.Text(), cert)
			dwg.isSending = false
			if err != nil {
				dwg.Toast.NotifyError(err.Error())
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

			if dwg.Dexc().HasWallet(int32(feeAsset.ID)) {
				dwg.completeRegistration(dex, feeAssetName, cert)
				return
			}

			dwg.createWalletWdg = newDexCreateWalletWidget(dwg.Load,
				&walletInfoWidget{
					image:    components.CoinImageBySymbol(&dwg.Load.Icons, feeAssetName),
					coinName: feeAssetName,
					coinID:   feeAsset.ID,
				},
				func() { dwg.completeRegistration(dex, feeAssetName, cert) })
		}()
	}

	if dwg.createWalletWdg != nil {
		dwg.createWalletWdg.handle()
	}
}

func (dwg *addDexWidget) layout(gtx layout.Context) D {
	if dwg.createWalletWdg != nil {
		return dwg.createWalletWdg.layout(gtx)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
				return dwg.Load.Theme.Label(values.TextSize20, "Add a dex").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
				return dwg.dexServerAddress.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Max.Y = 400
				return dwg.cert.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
				if dwg.isSending {
					dwg.addDexServer.Background = dwg.Theme.Color.Hint
				} else {
					dwg.addDexServer.Background = dwg.Theme.Color.Primary
				}
				return dwg.addDexServer.Layout(gtx)
			})
		}),
	)
}

func (dwg *addDexWidget) completeRegistration(dex *core.Exchange, feeAssetName string, cert []byte) {
	modal.NewPasswordModal(dwg.Load).
		Title("Confirm Registration").
		Hint("App password").
		Description(confirmRegisterModalDesc(dex, feeAssetName)).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton("Register", func(password string, pm *modal.PasswordModal) bool {
			go func() {
				_, err := dwg.Dexc().RegisterWithDEXServer(dex.Host,
					cert,
					int64(dex.Fee.Amt),
					int32(dex.Fee.ID),
					[]byte(password))
				if err != nil {
					pm.SetError(err.Error())
					pm.SetLoading(false)
					return
				}
				pm.Dismiss()
				pm.ChangeFragment(NewMarketPage(pm.Load))
			}()

			return false
		}).Show()
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

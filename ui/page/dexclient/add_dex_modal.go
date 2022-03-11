package dexclient

import (
	"fmt"
	"strconv"
	"strings"

	"decred.org/dcrdex/client/asset"
	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const addDexModalID = "add_dex_modal"

const testDexHost = "dex-test.ssgen.io:7232"

type addDexModal struct {
	*load.Load
	modal            *decredmaterial.Modal
	addDexServer     decredmaterial.Button
	dexServerAddress decredmaterial.Editor
	isSending        bool
	cert             decredmaterial.Editor
	cancel           decredmaterial.Button
	materialLoader   material.LoaderStyle
}

func newAddDexModal(l *load.Load) *addDexModal {
	md := &addDexModal{
		Load:             l,
		modal:            l.Theme.ModalFloatTitle(),
		dexServerAddress: l.Theme.Editor(new(widget.Editor), "DEX Address"),
		cert:             l.Theme.Editor(new(widget.Editor), "Cert content"),
		addDexServer:     l.Theme.Button("Submit"),
		cancel:           l.Theme.OutlineButton("Cancel"),
		materialLoader:   material.Loader(l.Theme.Base),
	}

	md.dexServerAddress.Editor.SingleLine = true
	if l.WL.MultiWallet.NetType() == dcrlibwallet.Testnet3 {
		md.dexServerAddress.Editor.SetText(testDexHost)
	}

	return md
}

func (md *addDexModal) ModalID() string {
	return addDexModalID
}

func (md *addDexModal) Show() {
	md.ShowModal(md)
}

func (md *addDexModal) Dismiss() {
	md.DismissModal(md)
}

func (md *addDexModal) OnDismiss() {
	md.dexServerAddress.Editor.SetText("")
}

func (md *addDexModal) OnResume() {
	md.dexServerAddress.Editor.Focus()
}

func (md *addDexModal) Handle() {
	if md.cancel.Button.Clicked() && !md.isSending {
		md.Dismiss()
	}

	if md.addDexServer.Button.Clicked() {
		if md.dexServerAddress.Editor.Text() == "" || md.isSending {
			return
		}

		md.isSending = true
		md.modal.SetDisabled(true)
		go func() {
			cert := []byte(md.cert.Editor.Text())
			dex, err := md.Dexc().DEXServerInfo(md.dexServerAddress.Editor.Text(), cert)
			md.isSending = false
			md.modal.SetDisabled(false)

			if err != nil {
				md.Toast.NotifyError(err.Error())
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

			// Dismiss this modal before displaying a new one for adding a wallet
			// or completing the registration.
			md.Dismiss()
			if md.Dexc().HasWallet(int32(feeAsset.ID)) {
				md.completeRegistration(dex, feeAssetName, cert)
				return
			}

			newCreateWalletModal(md.Load,
				&walletInfoWidget{
					image:    components.CoinImageBySymbol(&md.Load.Icons, feeAssetName),
					coinName: feeAssetName,
					coinID:   feeAsset.ID,
				},
				func() {
					md.completeRegistration(dex, feeAssetName, cert)
				}).Show()
		}()
	}
}

func (md *addDexModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return md.Load.Theme.Label(values.TextSize20, "Add a dex").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return md.dexServerAddress.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Max.Y = 350
						return md.cert.Layout(gtx)
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
						}.Layout(gtx, md.cancel.Layout)
					}),
					layout.Rigid(func(gtx C) D {
						if md.isSending {
							return layout.Inset{
								Top:    values.MarginPadding10,
								Bottom: values.MarginPadding15,
							}.Layout(gtx, md.materialLoader.Layout)
						}
						return md.addDexServer.Layout(gtx)
					}),
				)
			})
		},
	}

	return md.modal.Layout(gtx, w)
}

func (md *addDexModal) completeRegistration(dex *core.Exchange, feeAssetName string, cert []byte) {
	modal.NewPasswordModal(md.Load).
		Title("Confirm Registration").
		Hint("App password").
		Description(confirmRegisterModalDesc(dex, feeAssetName)).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton("Register", func(password string, pm *modal.PasswordModal) bool {
			go func() {
				_, err := md.Dexc().RegisterWithDEXServer(dex.Host,
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

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
	"github.com/planetdecred/godcr/ui/values"
)

const confirmRegisterModalID = "confirm_register_modal"

type confirmRegisterModal struct {
	*load.Load
	modal            *decredmaterial.Modal
	register         decredmaterial.Button
	appPassword      decredmaterial.Editor
	isSending        bool
	dex              *core.Exchange
	cert             []byte
	selectedFeeAsset string
	confirmed        func()
}

func newConfirmRegisterModal(l *load.Load, dex *core.Exchange, cert []byte, selectedFeeAsset string) *confirmRegisterModal {
	md := &confirmRegisterModal{
		Load:             l,
		modal:            l.Theme.ModalFloatTitle(),
		appPassword:      l.Theme.EditorPassword(new(widget.Editor), "App password"),
		register:         l.Theme.Button("Register"),
		dex:              dex,
		cert:             cert,
		selectedFeeAsset: selectedFeeAsset,
	}

	md.register.TextSize = values.TextSize12
	md.register.Background = l.Theme.Color.Primary
	md.appPassword.Editor.SingleLine = true

	return md
}

func (md *confirmRegisterModal) ModalID() string {
	return confirmRegisterModalID
}

func (md *confirmRegisterModal) Show() {
	md.ShowModal(md)
}

func (md *confirmRegisterModal) Dismiss() {
	md.DismissModal(md)
}

func (md *confirmRegisterModal) OnDismiss() {
	md.appPassword.Editor.SetText("")
}

func (md *confirmRegisterModal) OnResume() {
	md.appPassword.Editor.Focus()
}

func (md *confirmRegisterModal) Handle() {
	if md.register.Button.Clicked() {
		if md.appPassword.Editor.Text() == "" || md.isSending {
			return
		}

		md.isSending = true
		go func() {
			form := &core.RegisterForm{
				AppPass: []byte(md.appPassword.Editor.Text()),
				Addr:    md.dex.Host,
				Fee:     md.dex.Fee.Amt,
				Cert:    md.cert,
			}
			_, err := md.Dexc.Register(form)

			md.isSending = false
			if err != nil {
				md.Toast.NotifyError(err.Error())
				return
			}

			md.confirmed()
			md.Dismiss()
		}()
	}

}

func (md *confirmRegisterModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return md.Load.Theme.Label(values.TextSize20, "Confirm Registration").Layout(gtx)
		},
		func(gtx C) D {
			feeAsset := md.dex.RegFees[md.selectedFeeAsset]
			feeAmt := formatAmount(feeAsset.ID, md.selectedFeeAsset, feeAsset.Amt)
			txt := fmt.Sprintf("Enter your app password to confirm DEX registration. When you submit this form, %s will be spent from your wallet to pay registration fees.", feeAmt)
			return md.Load.Theme.Label(values.TextSize14, txt).Layout(gtx)
		},
		func(gtx C) D {
			markets := make([]string, 0, len(md.dex.Markets))
			for _, mkt := range md.dex.Markets {
				lotSize := formatAmount(mkt.BaseID, mkt.BaseSymbol, mkt.LotSize)
				markets = append(markets, fmt.Sprintf("Base: %s\tQuote: %s\tLot Size: %s", strings.ToUpper(mkt.BaseSymbol), strings.ToUpper(mkt.QuoteSymbol), lotSize))
			}
			txt := fmt.Sprintf("This DEX supports the following markets. All trades are in multiples of each market's lot size.\n\n%s", strings.Join(markets, "\n"))
			return md.Load.Theme.Label(values.TextSize14, txt).Layout(gtx)
		},
		func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
				return md.appPassword.Layout(gtx)
			})
		},
		func(gtx C) D {
			return md.register.Layout(gtx)
		},
	}

	return md.modal.Layout(gtx, w, 900)
}

func formatAmount(assetID uint32, assetName string, amount uint64) string {
	assetInfo, err := asset.Info(assetID)
	if err != nil {
		return fmt.Sprintf("%d [%s units]", amount, assetName)
	} else {
		unitInfo := assetInfo.UnitInfo
		convertedLotSize := float64(amount) / float64(unitInfo.Conventional.ConversionFactor)
		return fmt.Sprintf("%s %s", strconv.FormatFloat(convertedLotSize, 'f', -1, 64), unitInfo.Conventional.Unit)
	}
}

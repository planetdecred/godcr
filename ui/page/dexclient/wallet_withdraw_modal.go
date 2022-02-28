package dexclient

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex"
	"decred.org/dcrdex/dex/encode"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const withdraweModalID = "dex_withdraw_modal"

type withdrawModal struct {
	*load.Load
	modal                        *decredmaterial.Modal
	walletInfoWidget             *walletInfoWidget
	as                           *core.SupportedAsset
	address, amount, appPassword decredmaterial.Editor
	qrImage                      *image.Image
	cancelBtn                    decredmaterial.Button
	submitBtn                    decredmaterial.Button
	isSending                    bool
	materialLoader               material.LoaderStyle
}

// withdrawForm is sent to initiate a withdraw.
type withdrawForm struct {
	AssetID uint32           `json:"assetID"`
	Value   uint64           `json:"value"`
	Address string           `json:"address"`
	Pass    encode.PassBytes `json:"pw"`
}

func newWithdrawModal(l *load.Load, wallInfo *walletInfoWidget, as *core.SupportedAsset) *withdrawModal {
	md := &withdrawModal{
		Load:             l,
		walletInfoWidget: wallInfo,
		modal:            l.Theme.ModalFloatTitle(),
		cancelBtn:        l.Theme.OutlineButton(values.String(values.StrCancel)),
		submitBtn:        l.Theme.Button(strWithdraw),
		address:          l.Theme.Editor(&widget.Editor{SingleLine: true}, strAddress),
		amount:           l.Theme.Editor(&widget.Editor{SingleLine: true}, strAmount),
		appPassword:      l.Theme.EditorPassword(new(widget.Editor), strAppPassword),
		materialLoader:   material.Loader(material.NewTheme(gofont.Collection())),
		as:               as,
	}

	return md
}

func (md *withdrawModal) ModalID() string {
	return withdraweModalID
}

func (md *withdrawModal) Show() {
	md.ShowModal(md)
}

func (md *withdrawModal) Dismiss() {
	md.DismissModal(md)
}

func (md *withdrawModal) OnDismiss() {
}

func (md *withdrawModal) OnResume() {
}

func (md *withdrawModal) Handle() {
	if md.cancelBtn.Button.Clicked() {
		md.Dismiss()
	}

	if md.submitBtn.Button.Clicked() {
		if md.isSending {
			return
		}

		md.isSending = true
		md.modal.SetDisabled(true)
		if ok := md.doWithdraw(); !ok {
			md.isSending = false
			md.modal.SetDisabled(false)
			return
		}

		md.Dismiss()
	}
}

func (md *withdrawModal) doWithdraw() bool {
	amount, err := strconv.ParseFloat(md.amount.Editor.Text(), 64)
	if err != nil {
		md.Toast.NotifyError(err.Error())
		return false
	}
	v := uint64(amount * float64(md.as.Info.UnitInfo.Conventional.ConversionFactor))

	form := &withdrawForm{
		AssetID: md.walletInfoWidget.coinID,
		Value:   v,
		Address: md.address.Editor.Text(),
		Pass:    []byte(md.appPassword.Editor.Text()),
	}

	ok := md.Dexc().HasWallet(int32(form.AssetID))
	if !ok {
		md.Toast.NotifyError(fmt.Sprintf(nStrNoWalletFound, dex.BipIDSymbol(form.AssetID)))
		return false
	}

	_, err = md.Dexc().Core().Withdraw(form.Pass, form.AssetID, form.Value, form.Address)
	if err != nil {
		md.Toast.NotifyError(fmt.Sprintf(nStrWithdrawErr, err.Error()))
		return false
	}

	return true
}

func (md *withdrawModal) Layout(gtx layout.Context) D {
	w := []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(md.Load.Theme.Label(values.TextSize20, strWithdraw).Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding8, Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						ic := md.walletInfoWidget.image
						ic.Scale = 0.2
						return md.walletInfoWidget.image.Layout(gtx)
					})
				}),
				layout.Rigid(md.Load.Theme.Label(values.TextSize20, strings.ToUpper(md.walletInfoWidget.coinName)).Layout),
			)
		},
		func(gtx C) D {
			amt := formatAmount(md.as.Wallet.Balance.Available, &md.as.Info.UnitInfo)
			return md.Load.Theme.Label(values.TextSize14, fmt.Sprintf(nStrAmountAvailable, amt)).Layout(gtx)
		},
		md.address.Layout,
		md.amount.Layout,
		md.appPassword.Layout,
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

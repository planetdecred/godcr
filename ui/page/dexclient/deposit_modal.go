package dexclient

import (
	"bytes"
	"fmt"
	"image"
	"image/color"

	"decred.org/dcrdex/client/asset/btc"
	"decred.org/dcrdex/client/asset/dcr"
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/assets"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/yeqown/go-qrcode"
)

const depositeModalID = "deposit_modal"

type depositModal struct {
	*load.Load
	gtx              *layout.Context
	modal            *decredmaterial.Modal
	walletInfoWidget *walletInfoWidget
	wallAddress      string
	qrImage          *image.Image
	newAddr, copy    decredmaterial.Button
	cancel           decredmaterial.Button
}

func newDepositModal(l *load.Load, wallInfo *walletInfoWidget, wallAddress string) *depositModal {
	md := &depositModal{
		Load:             l,
		walletInfoWidget: wallInfo,
		wallAddress:      wallAddress,
		modal:            l.Theme.ModalFloatTitle(),
		cancel:           l.Theme.OutlineButton("Cancel"),
		copy:             l.Theme.Button("Copy"),
		newAddr:          l.Theme.Button("New Address"),
	}

	md.copy.Background = color.NRGBA{}
	md.copy.HighlightColor = md.Theme.Color.SurfaceHighlight
	md.copy.Color = md.Theme.Color.Primary

	md.newAddr.Background = md.Theme.Color.Surface
	md.newAddr.HighlightColor = md.Theme.Color.SurfaceHighlight
	md.newAddr.Color = md.Theme.Color.Primary

	md.generateQRForAddress(wallInfo.coinID)

	return md
}

func (md *depositModal) ModalID() string {
	return depositeModalID
}

func (md *depositModal) Show() {
	md.ShowModal(md)
}

func (md *depositModal) Dismiss() {
	md.DismissModal(md)
}

func (md *depositModal) OnDismiss() {
}

func (md *depositModal) OnResume() {
}

func (md *depositModal) Handle() {
	gtx := md.gtx

	if md.cancel.Button.Clicked() {
		md.Dismiss()
	}

	if md.copy.Clicked() {
		clipboard.WriteOp{Text: md.wallAddress}.Add(gtx.Ops)
		md.Toast.Notify(fmt.Sprintf("Copied %s address to clipboard", md.walletInfoWidget.coinName))
		return
	}

	if md.newAddr.Clicked() {
		newAddr, err := md.generateNewAddress(md.walletInfoWidget.coinID)
		if err != nil {
			fmt.Println("Error generating new address" + err.Error())
			return
		}

		md.wallAddress = newAddr
		md.generateQRForAddress(md.walletInfoWidget.coinID)
	}
}

func (md *depositModal) Layout(gtx layout.Context) D {
	if md.gtx == nil {
		md.gtx = &gtx
	}
	w := []layout.Widget{
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return md.Load.Theme.Label(values.TextSize20, "Deposit").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding8, Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						ic := md.walletInfoWidget.image
						ic.Scale = 0.2
						return md.walletInfoWidget.image.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return md.Load.Theme.Label(values.TextSize20, fmt.Sprintf("%s", md.walletInfoWidget.coinName)).Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Flex{
					Axis:      layout.Vertical,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if md.qrImage == nil {
							return layout.Dimensions{}
						}
						return md.Theme.ImageIcon(gtx, *md.qrImage, 450)
					}),
					layout.Rigid(func(gtx C) D {
						if md.wallAddress != "" {
							return md.addressLayout(gtx)
						}
						return layout.Dimensions{}
					}),
				)
			})
		},

		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return md.cancel.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return md.copy.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return md.newAddr.Layout(gtx)
					}),
				)
			})
		},
	}

	return md.modal.Layout(gtx, w)
}

func (md *depositModal) addressLayout(gtx layout.Context) layout.Dimensions {
	card := decredmaterial.Card{
		Color: md.Theme.Color.Gray4,
	}
	return layout.Inset{
		Top:    values.MarginPadding14,
		Bottom: values.MarginPadding16,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				card.Radius = decredmaterial.CornerRadius{TopRight: 8, TopLeft: 8, BottomRight: 8, BottomLeft: 8}
				return card.Layout(gtx, func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding16,
						Bottom: values.MarginPadding16,
						Left:   values.MarginPadding16,
						Right:  values.MarginPadding16,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return md.Theme.Label(values.TextSize14, md.wallAddress).Layout(gtx)
					})
				})
			}),
		)
	})
}

func (md *depositModal) generateQRForAddress(coinID uint32) {
	imgName := ""

	switch coinID {
	case dcr.BipID:
		imgName = "qrcodeSymbol"
	case btc.BipID:
		imgName = "dex_btc"
	}

	opt := qrcode.WithLogoImage(assets.DecredIcons[imgName])
	qrCode, err := qrcode.New(md.wallAddress, opt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var buff bytes.Buffer
	err = qrCode.SaveTo(&buff)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	imgdec, _, err := image.Decode(bytes.NewReader(buff.Bytes()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	md.qrImage = &imgdec
}

func (md *depositModal) generateNewAddress(assetID uint32) (string, error) {
	addr, err := md.Dexc().Core().NewDepositAddress(assetID)
	if err != nil {
		return "", err
	}

	return addr, nil
}

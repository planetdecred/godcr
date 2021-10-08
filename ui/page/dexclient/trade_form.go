package dexclient

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	StrPrice          = "Price"
	StrQuantity       = "Quantity"
	StrLots           = "Lots"
	StrDCR            = "DCR"
	StrPlaceBuyOrder  = "Place orde to buy DCR"
	StrPlaceSellOrder = "Place orde to sell DCR"
	StrTifNow         = "Immediate or cancel"
	StrSell           = "Sell"
	StrBuy            = "Buy"
)

type TradeFormWidget struct {
	isLimit                     bool
	isSell                      bool
	tifNow                      decredmaterial.CheckBoxStyle
	buy, sell                   decredmaterial.Button
	reviewOrder                 decredmaterial.Button
	host                        string
	price                       decredmaterial.Label
	quantity                    decredmaterial.Label
	rateField                   decredmaterial.Editor
	lotField                    decredmaterial.Editor
	qtyField                    decredmaterial.Editor
	marketBaseID, marketQuoteID uint32
	marketQuote, marketBase     string
	buyColor                    color.NRGBA
	sellColor                   color.NRGBA
	inactiveColor               color.NRGBA
}

func NewTradeFormWidget(l *load.Load) *TradeFormWidget {
	form := &TradeFormWidget{
		price:       l.Theme.Label(values.TextSize14, StrPrice),
		quantity:    l.Theme.Label(values.TextSize14, StrQuantity),
		rateField:   l.Theme.Editor(new(widget.Editor), StrPrice),
		lotField:    l.Theme.Editor(new(widget.Editor), StrLots),
		qtyField:    l.Theme.Editor(new(widget.Editor), StrDCR),
		reviewOrder: l.Theme.Button(StrPlaceBuyOrder),
		tifNow:      l.Theme.CheckBox(new(widget.Bool), StrTifNow),
		sell:        l.Theme.Button(StrSell),
		buy:         l.Theme.Button(StrBuy),
		isSell:      false,
	}
	form.rateField.Editor.SingleLine = true
	form.rateField.Editor.SetText("0")
	form.lotField.Editor.SingleLine = true
	form.qtyField.Editor.SingleLine = true
	form.reviewOrder.Background = l.Theme.Color.Buy

	inset := layout.Inset{
		Left:   values.MarginPadding8,
		Right:  values.MarginPadding8,
		Top:    values.MarginPadding5,
		Bottom: values.MarginPadding5,
	}
	form.reviewOrder.Inset = inset
	form.reviewOrder.TextSize = values.TextSize14

	inset.Left = values.MarginPadding40
	inset.Right = values.MarginPadding40
	form.sell.Background = l.Theme.Color.InactiveGray
	form.sell.Inset = inset
	form.sell.TextSize = values.TextSize14

	form.buy.Background = l.Theme.Color.Buy
	form.buy.Inset = inset
	form.buy.TextSize = values.TextSize14

	form.buyColor = l.Theme.Color.Buy
	form.sellColor = l.Theme.Color.Sell
	form.inactiveColor = l.Theme.Color.InactiveGray

	return form
}

func (mktForm *TradeFormWidget) Layout(gtx C) D {
	mktForm.handler()

	inset := layout.Inset{Bottom: values.MarginPadding10}
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	return inset.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(mktForm.sellOrBuyButtonsLayout),
			layout.Rigid(func(gtx C) D {
				return inset.Layout(gtx, func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(.2, func(gtx C) D {
							return mktForm.price.Layout(gtx)
						}),
						layout.Flexed(.8, func(gtx C) D {
							txt := fmt.Sprintf("%s/%s", strings.ToUpper(mktForm.marketQuote), strings.ToUpper(mktForm.marketBase))
							mktForm.rateField.Hint = txt
							return mktForm.rateField.Layout(gtx)
						}),
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return inset.Layout(gtx, func(gtx C) D {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(.2, mktForm.quantity.Layout),
						layout.Flexed(.8, func(gtx C) D {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Flexed(0.4, func(gtx C) D {
									return layout.Inset{
										Right: values.MarginPadding20,
									}.Layout(gtx, mktForm.lotField.Layout)
								}),
								layout.Flexed(0.6, mktForm.qtyField.Layout),
							)
						}),
					)
				})
			}),
			layout.Rigid(mktForm.tifNow.Layout),
			layout.Rigid(func(gtx C) D {
				return layout.E.Layout(gtx, mktForm.reviewOrder.Layout)
			}),
		)
	})
}

func (mktForm *TradeFormWidget) sellOrBuyButtonsLayout(gtx C) D {
	return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
					return mktForm.buy.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return mktForm.sell.Layout(gtx)
			}),
		)
	})
}

func (mktForm *TradeFormWidget) SetMarketBaseID(mktBaseID uint32) *TradeFormWidget {
	mktForm.marketBaseID = mktBaseID
	return mktForm
}

func (mktForm *TradeFormWidget) SetMarketQuoteID(mktQuoteID uint32) *TradeFormWidget {
	mktForm.marketQuoteID = mktQuoteID
	return mktForm
}

func (mktForm *TradeFormWidget) SetMarketBase(mktBase string) *TradeFormWidget {
	mktForm.marketBase = mktBase
	return mktForm
}

func (mktForm *TradeFormWidget) SetMarketQuote(mktQuote string) *TradeFormWidget {
	mktForm.marketQuote = mktQuote
	return mktForm
}

func (mktForm *TradeFormWidget) SetHost(host string) *TradeFormWidget {
	mktForm.host = host
	return mktForm
}

func (mktForm *TradeFormWidget) ShowVerifyOrder() bool {
	return mktForm.reviewOrder.Button.Clicked()
}

func (mktForm *TradeFormWidget) handler() {
	if mktForm.buy.Button.Clicked() {
		mktForm.isSell = false
		mktForm.sell.Background = mktForm.inactiveColor
		mktForm.buy.Background = mktForm.buyColor
		mktForm.reviewOrder.Text = StrPlaceBuyOrder
		mktForm.reviewOrder.Background = mktForm.buyColor
	}

	if mktForm.sell.Button.Clicked() {
		mktForm.isSell = true
		mktForm.sell.Background = mktForm.sellColor
		mktForm.buy.Background = mktForm.inactiveColor
		mktForm.reviewOrder.Text = StrPlaceSellOrder
		mktForm.reviewOrder.Background = mktForm.sellColor
	}
}

// TradeForm get trade form values to buy or sell an order
func (mktForm *TradeFormWidget) TradeForm() (error, *core.TradeForm) {
	qtyField, err := strconv.ParseUint(mktForm.qtyField.Editor.Text(), 10, 64)
	if err != nil {
		return err, nil
	}

	rateField, err := strconv.ParseFloat(mktForm.rateField.Editor.Text(), 64)
	if err != nil {
		return err, nil
	}

	return nil, &core.TradeForm{
		Host:    mktForm.host,
		Base:    mktForm.marketBaseID,
		Quote:   mktForm.marketQuoteID,
		Qty:     qtyField * 1e8,
		Rate:    uint64(rateField * 1e8),
		IsLimit: mktForm.isLimit,
		Sell:    mktForm.isSell,
		TifNow:  mktForm.tifNow.CheckBox.Value,
	}
}

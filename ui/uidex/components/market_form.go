package components

import (
	"fmt"
	"strconv"
	"strings"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

type TradeFormWidget struct {
	isLimit                     bool
	sell                        bool
	tifNow                      decredmaterial.CheckBoxStyle
	reviewOrder                 decredmaterial.Button
	host                        string
	price                       decredmaterial.Label
	quantity                    decredmaterial.Label
	rateField                   decredmaterial.Editor
	lotField                    decredmaterial.Editor
	qtyField                    decredmaterial.Editor
	marketBaseID, marketQuoteID uint32
	marketQuote, marketBase     string
}

func NewTradeFormWidget(theme *decredmaterial.Theme) *TradeFormWidget {
	form := &TradeFormWidget{
		price:       theme.Label(values.TextSize14, "Price"),
		quantity:    theme.Label(values.TextSize14, "Quantity"),
		rateField:   theme.Editor(new(widget.Editor), "Price"),
		lotField:    theme.Editor(new(widget.Editor), "Lots"),
		qtyField:    theme.Editor(new(widget.Editor), "DCR"),
		reviewOrder: theme.Button(new(widget.Clickable), "Review order"),
		tifNow:      theme.CheckBox(new(widget.Bool), "Immediate or cancel"),
	}
	form.rateField.Editor.SingleLine = true
	form.rateField.Editor.SetText("0")
	form.lotField.Editor.SingleLine = true
	form.qtyField.Editor.SingleLine = true
	form.reviewOrder.Background = theme.Color.Success

	return form
}

func (mktForm *TradeFormWidget) Layout(gtx layout.Context) layout.Dimensions {
	inset := layout.Inset{Bottom: values.MarginPadding10}
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(.2, func(gtx layout.Context) layout.Dimensions {
							return mktForm.price.Layout(gtx)
						}),
						layout.Flexed(.8, func(gtx layout.Context) layout.Dimensions {
							txt := fmt.Sprintf("%s/%s", strings.ToUpper(mktForm.marketQuote), strings.ToUpper(mktForm.marketBase))
							mktForm.rateField.Hint = txt
							return mktForm.rateField.Layout(gtx)
						}),
					)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(.2, func(gtx layout.Context) layout.Dimensions {
							return mktForm.quantity.Layout(gtx)
						}),
						layout.Flexed(.8, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Flexed(0.4, func(gtx layout.Context) layout.Dimensions {
									return layout.Inset{Right: values.MarginPadding20}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return mktForm.lotField.Layout(gtx)
									})
								}),
								layout.Flexed(0.6, func(gtx layout.Context) layout.Dimensions {
									return mktForm.qtyField.Layout(gtx)
								}),
							)
						}),
					)
				})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return mktForm.tifNow.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				mktForm.reviewOrder.Text = "Place order to buy  DCR"
				return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return mktForm.reviewOrder.Layout(gtx)
				})
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
		Sell:    mktForm.sell,
		TifNow:  mktForm.tifNow.CheckBox.Value,
	}
}

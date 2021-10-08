package dexclient

import (
	"image/color"
	"strconv"
	"strings"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/dexc"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

type miniTradeFormWidget struct {
	*load.Load
	isSell                        bool
	submit                        decredmaterial.Button
	direction                     decredmaterial.IconButton
	invoicedAmount, orderedAmount decredmaterial.Editor
	orderBook                     *core.OrderBook
	mkt                           *selectedMaket
}

func newMiniTradeFormWidget(l *load.Load) *miniTradeFormWidget {
	miniTradeFormWdg := &miniTradeFormWidget{
		Load:           l,
		submit:         l.Theme.Button("OK"),
		invoicedAmount: l.Theme.Editor(new(widget.Editor), "I have"),
		orderedAmount:  l.Theme.Editor(new(widget.Editor), "I get"),
		direction:      l.Theme.IconButton(l.Icons.ActionSwapHoriz),
		orderBook:      new(core.OrderBook),
		isSell:         true,

		// TODO: will add colapsible to select trading pairs
		mkt: mkt,
	}

	miniTradeFormWdg.direction.Background = color.NRGBA{}
	miniTradeFormWdg.direction.Size = values.MarginPadding20
	miniTradeFormWdg.direction.Color = l.Theme.Color.Primary

	miniTradeFormWdg.invoicedAmount.Editor.SingleLine = true
	miniTradeFormWdg.invoicedAmount.HasCustomButton = true
	miniTradeFormWdg.invoicedAmount.CustomButton.Inset = layout.UniformInset(values.MarginPadding6)

	miniTradeFormWdg.orderedAmount.Editor.SingleLine = true
	miniTradeFormWdg.orderedAmount.HasCustomButton = true
	miniTradeFormWdg.orderedAmount.CustomButton.Inset = layout.UniformInset(values.MarginPadding6)
	miniTradeFormWdg.updateInputStatus()

	miniTradeFormWdg.submit.TextSize = values.TextSize12
	miniTradeFormWdg.submit.Background = l.Theme.Color.Primary

	return miniTradeFormWdg
}

func (miniTradeFormWdg *miniTradeFormWidget) layout(gtx C) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Flexed(.5, miniTradeFormWdg.invoicedAmount.Layout),
				layout.Rigid(miniTradeFormWdg.direction.Layout),
				layout.Flexed(.5, miniTradeFormWdg.orderedAmount.Layout),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
				return miniTradeFormWdg.submit.Layout(gtx)
			})
		}),
	)
}

func (miniTradeFormWdg *miniTradeFormWidget) updateInputStatus() {
	if miniTradeFormWdg.isSell {
		miniTradeFormWdg.invoicedAmount.Editor.SetText("0")
		miniTradeFormWdg.invoicedAmount.CustomButton.Text = strings.ToUpper(miniTradeFormWdg.mkt.marketBase)
		miniTradeFormWdg.invoicedAmount.CustomButton.Background = miniTradeFormWdg.Theme.Color.Primary

		miniTradeFormWdg.orderedAmount.Editor.SetText("0")
		miniTradeFormWdg.orderedAmount.CustomButton.Text = strings.ToUpper(miniTradeFormWdg.mkt.marketQuote)
		miniTradeFormWdg.orderedAmount.CustomButton.Background = miniTradeFormWdg.Theme.Color.Success
	} else {
		miniTradeFormWdg.invoicedAmount.Editor.SetText("0")
		miniTradeFormWdg.invoicedAmount.CustomButton.Text = strings.ToUpper(miniTradeFormWdg.mkt.marketQuote)
		miniTradeFormWdg.invoicedAmount.CustomButton.Background = miniTradeFormWdg.Theme.Color.Success

		miniTradeFormWdg.orderedAmount.Editor.SetText("0")
		miniTradeFormWdg.orderedAmount.CustomButton.Text = strings.ToUpper(miniTradeFormWdg.mkt.marketBase)
		miniTradeFormWdg.orderedAmount.CustomButton.Background = miniTradeFormWdg.Theme.Color.Primary
	}
}

func (miniTradeFormWdg *miniTradeFormWidget) handle() {
	if miniTradeFormWdg.direction.Button.Clicked() {
		miniTradeFormWdg.isSell = !miniTradeFormWdg.isSell
		miniTradeFormWdg.updateInputStatus()
	}

	ord := miniTradeFormWdg.orderBook
	if ord != nil {
		var rate float64

		bitSize := 64

		if miniTradeFormWdg.isSell {
			_, rate = minMaxRateOrderBook(ord.Buys)
		} else {
			rate, _ = minMaxRateOrderBook(ord.Sells)
		}

		for _, evt := range miniTradeFormWdg.invoicedAmount.Editor.Events() {
			if miniTradeFormWdg.invoicedAmount.Editor.Focused() {
				switch evt.(type) {
				case widget.ChangeEvent:
					qty, err := strconv.ParseFloat(miniTradeFormWdg.invoicedAmount.Editor.Text(), bitSize)
					if err != nil || rate <= 0 {
						return
					}
					var value string
					if miniTradeFormWdg.isSell {
						value = strconv.FormatFloat(rate*qty, 'f', -1, bitSize)
					} else {
						value = strconv.FormatFloat(qty/rate, 'f', -1, bitSize)
					}
					miniTradeFormWdg.orderedAmount.Editor.SetText(value)
				}
			}
		}

		for _, evt := range miniTradeFormWdg.orderedAmount.Editor.Events() {
			if miniTradeFormWdg.orderedAmount.Editor.Focused() {
				switch evt.(type) {
				case widget.ChangeEvent:
					qty, err := strconv.ParseFloat(miniTradeFormWdg.orderedAmount.Editor.Text(), bitSize)
					if err != nil || rate <= 0 {
						return
					}
					var value string
					if miniTradeFormWdg.isSell {
						value = strconv.FormatFloat(qty/rate, 'f', -1, bitSize)
					} else {
						value = strconv.FormatFloat(rate*qty, 'f', -1, bitSize)
					}
					miniTradeFormWdg.invoicedAmount.Editor.SetText(value)
				}
			}
		}
	}

	if miniTradeFormWdg.submit.Button.Clicked() {
		var qty uint64
		if miniTradeFormWdg.isSell {
			v, err := strconv.ParseUint(miniTradeFormWdg.invoicedAmount.Editor.Text(), 10, 64)
			if err != nil {
				miniTradeFormWdg.Toast.NotifyError(err.Error())
				return
			}
			qty = v * dexc.ConversionFactor
		} else {
			v, err := strconv.ParseFloat(miniTradeFormWdg.invoicedAmount.Editor.Text(), 64)
			if err != nil {
				miniTradeFormWdg.Toast.NotifyError(err.Error())
				return
			}
			qty = uint64(v * dexc.ConversionFactor)
		}

		md := newconfirmPasswordModal(miniTradeFormWdg.Load)
		md.confirmed = func(password []byte) {
			form := core.TradeForm{
				Host:    miniTradeFormWdg.mkt.host,
				Base:    miniTradeFormWdg.mkt.marketBaseID,
				Quote:   miniTradeFormWdg.mkt.marketQuoteID,
				Qty:     qty,
				IsLimit: false,
				Sell:    miniTradeFormWdg.isSell,
				TifNow:  false,
			}
			_, err := miniTradeFormWdg.DL.Trade(password, &form)
			if err != nil {
				miniTradeFormWdg.Toast.NotifyError(err.Error())
				return
			}
		}

		md.Show()
	}
}

func minMaxRateOrderBook(orders []*core.MiniOrder) (float64, float64) {
	if len(orders) == 0 {
		return 0, 0
	}
	var max = orders[0].Rate
	var min = orders[0].Rate
	for _, value := range orders {
		if max < value.Rate {
			max = value.Rate
		}
		if min > value.Rate {
			min = value.Rate
		}
	}
	return min, max
}

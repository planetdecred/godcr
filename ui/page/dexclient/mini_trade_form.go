package dexclient

import (
	"fmt"
	"image/color"
	"strconv"

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
	host                          string
	invoicedAmount, orderedAmount decredmaterial.Editor
	marketBaseID, marketQuoteID   uint32
	marketQuote, marketBase       string
	orderBook                     *core.OrderBook
}

func newMiniTradeFormWidget(l *load.Load) *miniTradeFormWidget {
	miniTradeFormWdg := &miniTradeFormWidget{
		Load:           l,
		submit:         l.Theme.Button(new(widget.Clickable), "SWAP"),
		invoicedAmount: l.Theme.Editor(new(widget.Editor), "I have"),
		orderedAmount:  l.Theme.Editor(new(widget.Editor), "I get"),
		direction:      l.Theme.IconButton(new(widget.Clickable), l.Icons.ActionSwapHoriz),
		orderBook:      new(core.OrderBook),
		isSell:         true,

		// TODO: will add colapsible to select trading pairs
		host:          testDexHost,
		marketBaseID:  dexc.DefaultAssetID,
		marketQuoteID: 0,
	}

	miniTradeFormWdg.direction.Background = color.NRGBA{}
	miniTradeFormWdg.direction.Size = values.MarginPadding20
	miniTradeFormWdg.direction.Color = l.Theme.Color.DeepBlue

	miniTradeFormWdg.invoicedAmount.Editor.SingleLine = true
	miniTradeFormWdg.invoicedAmount.Editor.SetText("0")
	miniTradeFormWdg.invoicedAmount.HasCustomButton = true
	miniTradeFormWdg.invoicedAmount.CustomButton.Inset = layout.UniformInset(values.MarginPadding6)
	miniTradeFormWdg.invoicedAmount.CustomButton.Text = "DCR"

	miniTradeFormWdg.orderedAmount.Editor.SingleLine = true
	miniTradeFormWdg.orderedAmount.Editor.SetText("0")
	miniTradeFormWdg.orderedAmount.HasCustomButton = true
	miniTradeFormWdg.orderedAmount.CustomButton.Background = l.Theme.Color.Success
	miniTradeFormWdg.orderedAmount.CustomButton.Inset = layout.UniformInset(values.MarginPadding6)
	miniTradeFormWdg.orderedAmount.CustomButton.Text = "BTC"

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

func (miniTradeFormWdg *miniTradeFormWidget) handle() {
	// if miniTradeFormWdg.direction.Button.Clicked() {
	// 	miniTradeFormWdg.isSell = !miniTradeFormWdg.isSell
	// }

	ord := miniTradeFormWdg.orderBook
	var rate float64 = 0

	if miniTradeFormWdg.isSell {
		_, rate = minMaxRateOrderBook(ord.Buys)
	} else {
		rate, _ = minMaxRateOrderBook(ord.Sells)
	}

	for _, evt := range miniTradeFormWdg.invoicedAmount.Editor.Events() {
		if miniTradeFormWdg.invoicedAmount.Editor.Focused() {
			switch evt.(type) {
			case widget.ChangeEvent:
				qty, err := strconv.ParseFloat(miniTradeFormWdg.invoicedAmount.Editor.Text(), 64)
				if err != nil {
					return
				}
				miniTradeFormWdg.orderedAmount.Editor.SetText(fmt.Sprintf("%f", rate*qty))
			}
		}
	}

	for _, evt := range miniTradeFormWdg.orderedAmount.Editor.Events() {
		if miniTradeFormWdg.orderedAmount.Editor.Focused() {
			switch evt.(type) {
			case widget.ChangeEvent:
				qty, err := strconv.ParseFloat(miniTradeFormWdg.orderedAmount.Editor.Text(), 64)
				if err != nil {
					return
				}
				miniTradeFormWdg.invoicedAmount.Editor.SetText(fmt.Sprintf("%f", qty/rate))
			}
		}
	}

	if miniTradeFormWdg.submit.Button.Clicked() {
		inAmount, err := strconv.ParseUint(miniTradeFormWdg.invoicedAmount.Editor.Text(), 10, 64)
		if err != nil {
			miniTradeFormWdg.Toast.NotifyError(err.Error())
			return
		}

		md := newconfirmPasswordModal(miniTradeFormWdg.Load)

		md.confirmed = func(password []byte) {
			form := core.TradeForm{
				Host:    miniTradeFormWdg.host,
				Base:    miniTradeFormWdg.marketBaseID,
				Quote:   miniTradeFormWdg.marketQuoteID,
				Qty:     inAmount * dexc.ConversionFactor,
				IsLimit: false,
				Sell:    miniTradeFormWdg.isSell,
				TifNow:  false,
			}

			_, err = miniTradeFormWdg.DL.Trade(password, &form)

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

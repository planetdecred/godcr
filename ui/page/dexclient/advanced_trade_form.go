package dexclient

import (
	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/load"
)

type advancedTradeFormWidget struct {
	*load.Load
	depthChart *depthChart
}

func newAdvancedTradeFormWidget(l *load.Load) *advancedTradeFormWidget {
	a := &advancedTradeFormWidget{
		Load:       l,
		depthChart: newDepthChart(l),
	}

	return a
}

func (a *advancedTradeFormWidget) layout(orderBook *core.OrderBook) layout.Widget {
	return func(gtx C) D {
		if orderBook == nil {
			return D{}
		}

		return a.depthChart.layout(gtx, orderBook.Buys, orderBook.Sells)
	}
}

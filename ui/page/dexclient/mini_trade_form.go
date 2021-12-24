package dexclient

import (
	"image/color"
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
	"github.com/planetdecred/godcr/ui/values"
)

type miniTradeFormWidget struct {
	*load.Load
	isSell                        bool
	submitBtn                     decredmaterial.Button
	directionBtn                  decredmaterial.IconButton
	invoicedAmount, orderedAmount decredmaterial.Editor
	host                          string
	mkt                           *core.Market
}

func newMiniTradeFormWidget(l *load.Load, host string, mkt *core.Market) *miniTradeFormWidget {
	m := &miniTradeFormWidget{
		Load:           l,
		submitBtn:      l.Theme.Button("OK"),
		invoicedAmount: l.Theme.Editor(new(widget.Editor), "I have"),
		orderedAmount:  l.Theme.Editor(new(widget.Editor), "I get"),
		directionBtn:   l.Theme.IconButton(l.Icons.ExchangeIcon),
		isSell:         true,
		host:           host,
		mkt:            mkt,
	}

	m.directionBtn.Size = values.MarginPadding20
	m.directionBtn.ChangeColorStyle(&values.ColorStyle{Background: color.NRGBA{}})

	m.invoicedAmount.Editor.SingleLine = true
	m.invoicedAmount.HasCustomButton = true
	m.invoicedAmount.CustomButton.Inset = layout.UniformInset(values.MarginPadding6)

	m.orderedAmount.Editor.SingleLine = true
	m.orderedAmount.HasCustomButton = true
	m.orderedAmount.CustomButton.Inset = layout.UniformInset(values.MarginPadding6)
	m.changeDirection()

	m.submitBtn.TextSize = values.TextSize12
	m.submitBtn.Background = l.Theme.Color.Primary

	return m
}

func (m *miniTradeFormWidget) layout(gtx C) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Flexed(.5, m.invoicedAmount.Layout),
				layout.Rigid(m.directionBtn.Layout),
				layout.Flexed(.5, m.orderedAmount.Layout),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
				return m.submitBtn.Layout(gtx)
			})
		}),
	)
}

func (m *miniTradeFormWidget) changeDirection() {
	m.orderedAmount.Editor.SetText("0")
	m.invoicedAmount.Editor.SetText("0")
	if m.mkt == nil {
		return
	}
	if m.isSell {
		m.invoicedAmount.CustomButton.Text = strings.ToUpper(m.mkt.BaseSymbol)
		m.invoicedAmount.CustomButton.Background = m.Theme.Color.Primary
		m.orderedAmount.CustomButton.Text = strings.ToUpper(m.mkt.QuoteSymbol)
		m.orderedAmount.CustomButton.Background = m.Theme.Color.Success
	} else {
		m.invoicedAmount.CustomButton.Text = strings.ToUpper(m.mkt.QuoteSymbol)
		m.invoicedAmount.CustomButton.Background = m.Theme.Color.Success
		m.orderedAmount.CustomButton.Text = strings.ToUpper(m.mkt.BaseSymbol)
		m.orderedAmount.CustomButton.Background = m.Theme.Color.Primary
	}
}

func (m miniTradeFormWidget) suggestValue(ord *core.OrderBook, amount string, isSell, invoicedChange bool) string {
	var rate float64
	bitSize := 64
	if m.isSell {
		_, rate = minMaxRateOrderBook(ord.Buys)
	} else {
		rate, _ = minMaxRateOrderBook(ord.Sells)
	}

	qty, err := strconv.ParseFloat(amount, bitSize)
	if err != nil || rate <= 0 {
		return ""
	}

	if invoicedChange {
		if isSell {
			return strconv.FormatFloat(rate*qty, 'f', -1, bitSize)
		}
		return strconv.FormatFloat(qty/rate, 'f', -1, bitSize)
	}

	if isSell {
		return strconv.FormatFloat(qty/rate, 'f', -1, bitSize)
	}
	return strconv.FormatFloat(rate*qty, 'f', -1, bitSize)
}

func (m *miniTradeFormWidget) handle(ord *core.OrderBook) {
	if m.directionBtn.Button.Clicked() {
		m.isSell = !m.isSell
		m.changeDirection()
	}

	if ord != nil {
		if m.invoicedAmount.Editor.Focused() {
			_, change := decredmaterial.HandleEditorEvents(m.invoicedAmount.Editor)
			if change {
				value := m.suggestValue(ord, m.invoicedAmount.Editor.Text(), m.isSell, true)
				m.orderedAmount.Editor.SetText(value)
			}
		}

		if m.orderedAmount.Editor.Focused() {
			_, change := decredmaterial.HandleEditorEvents(m.orderedAmount.Editor)
			if change {
				value := m.suggestValue(ord, m.orderedAmount.Editor.Text(), m.isSell, false)
				m.invoicedAmount.Editor.SetText(value)
			}
		}
	}

	if m.submitBtn.Button.Clicked() {
		if m.mkt == nil {
			return
		}
		var qty uint64
		if m.isSell {
			assetInfo, err := asset.Info(m.mkt.BaseID)
			if err != nil {
				m.Toast.NotifyError(err.Error())
				return
			}
			v, err := strconv.ParseUint(m.invoicedAmount.Editor.Text(), 10, 64)
			if err != nil {
				m.Toast.NotifyError(err.Error())
				return
			}
			qty = v * assetInfo.UnitInfo.Conventional.ConversionFactor
		} else {
			assetInfo, err := asset.Info(m.mkt.QuoteID)
			if err != nil {
				m.Toast.NotifyError(err.Error())
				return
			}
			v, err := strconv.ParseFloat(m.invoicedAmount.Editor.Text(), 64)
			if err != nil {
				m.Toast.NotifyError(err.Error())
				return
			}
			qty = uint64(v * float64(assetInfo.UnitInfo.Conventional.ConversionFactor))
		}

		modal.NewPasswordModal(m.Load).
			Title("App password").
			Hint("Authorize this order with your app password.").
			Description("IMPORTANT: Trades take time to settle, and you cannot turn off the DEX client software, or the BTC or DCR blockchain and/or wallet software, until settlement is complete. Settlement can complete as quickly as a few minutes or take as long as a few hours.").
			NegativeButton(values.String(values.StrCancel), func() {}).
			PositiveButton("Ok", func(password string, pm *modal.PasswordModal) bool {
				go func() {
					form := dcrlibwallet.FreshOrder{
						BaseAssetID:  m.mkt.BaseID,
						QuoteAssetID: m.mkt.QuoteID,
						Qty:          qty,
						IsLimit:      false,
						Sell:         m.isSell,
						TifNow:       false,
					}
					_, err := m.Dexc().PlaceOrderWithServer(m.host, &form, []byte(password))
					if err != nil {
						pm.SetError(err.Error())
						pm.SetLoading(false)
						return
					}
					m.orderedAmount.Editor.SetText("0")
					m.invoicedAmount.Editor.SetText("0")
					m.Toast.Notify("Successfully!")
					pm.Dismiss()
				}()
				return false
			}).Show()
	}
}

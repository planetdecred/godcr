package dexclient

import (
	"fmt"
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
	"github.com/planetdecred/godcr/ui/page/components"
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

func newMiniTradeFormWidget(l *load.Load) *miniTradeFormWidget {
	m := &miniTradeFormWidget{
		Load:           l,
		submitBtn:      l.Theme.Button(strSubmit),
		invoicedAmount: l.Theme.Editor(new(widget.Editor), strIHave),
		orderedAmount:  l.Theme.Editor(new(widget.Editor), strIGet),
		directionBtn:   l.Theme.IconButton(l.Icons.ExchangeIcon),
		isSell:         true,
	}

	m.directionBtn.Size = values.MarginPadding20
	m.directionBtn.ChangeColorStyle(&values.ColorStyle{Background: color.NRGBA{}})

	m.invoicedAmount.Editor.SingleLine = true
	m.invoicedAmount.HasCustomButton = true
	m.invoicedAmount.CustomButton.Margin.Top = values.MarginPaddingMinus2
	m.invoicedAmount.CustomButton.Inset = layout.UniformInset(values.MarginPadding6)

	m.orderedAmount.Editor.SingleLine = true
	m.orderedAmount.HasCustomButton = true
	m.orderedAmount.CustomButton.Margin.Top = values.MarginPaddingMinus2
	m.orderedAmount.CustomButton.Inset = layout.UniformInset(values.MarginPadding6)

	m.submitBtn.TextSize = values.TextSize12
	m.submitBtn.SetEnabled(false)

	return m
}

func (m *miniTradeFormWidget) setHostAndMarket(host string, market *core.Market) *miniTradeFormWidget {
	if m.host != host || m.mkt == nil || (m.mkt != nil && m.mkt.Name != market.Name) {
		m.host = host
		m.mkt = market
		m.isSell = true
		m.changeDirection()
	}

	return m
}

func (m *miniTradeFormWidget) layout(gtx C) D {
	if m.mkt == nil {
		return D{}
	}

	mapAssets := m.Dexc().Core().SupportedAssets()
	baseAsset := mapAssets[m.mkt.BaseID]
	quoteAsset := mapAssets[m.mkt.QuoteID]
	baseUnitInfo := baseAsset.Info.UnitInfo
	quoteUnitInfo := quoteAsset.Info.UnitInfo
	lotSize := fmt.Sprintf("Lot Size: %s %s", formatAmount(m.mkt.LotSize, &baseUnitInfo), baseUnitInfo.Conventional.Unit)
	rateStep := fmt.Sprintf("Rate Step: %s %s", formatAmount(m.mkt.RateStep, &quoteUnitInfo), quoteUnitInfo.Conventional.Unit)
	baseSyncPercentage := ""
	quoteSyncPercentage := ""
	if baseAsset.Wallet != nil && !baseAsset.Wallet.Synced {
		baseSyncPercentage = walletSyncPercentage(baseAsset.Wallet)
	}
	if quoteAsset.Wallet != nil && !quoteAsset.Wallet.Synced {
		quoteSyncPercentage = walletSyncPercentage(quoteAsset.Wallet)
	}

	row := func(e *decredmaterial.Editor, isInvoiced, isSell bool, baseSyncPercentage, quoteSyncPercentage string) layout.Widget {
		syncPercentage := ""
		if isInvoiced {
			if isSell {
				syncPercentage = baseSyncPercentage
			} else {
				syncPercentage = quoteSyncPercentage
			}
		} else {
			if isSell {
				syncPercentage = quoteSyncPercentage
			} else {
				syncPercentage = baseSyncPercentage
			}
		}

		return func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(e.Layout),
				layout.Rigid(func(gtx C) D {
					t := m.Theme.Label(values.TextSize12, syncPercentage)
					t.Color = m.Theme.Color.Orange
					return layout.E.Layout(gtx, t.Layout)
				}),
			)
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Flexed(.5, row(&m.invoicedAmount, true, m.isSell, baseSyncPercentage, quoteSyncPercentage)),
				layout.Rigid(m.directionBtn.Layout),
				layout.Flexed(.5, row(&m.orderedAmount, false, m.isSell, baseSyncPercentage, quoteSyncPercentage)),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding12, Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(m.Theme.Label(values.TextSize14, lotSize).Layout),
					layout.Rigid(m.Theme.Label(values.TextSize14, rateStep).Layout),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, m.submitBtn.Layout)
		}),
	)
}

func (m *miniTradeFormWidget) changeDirection() {
	t1 := m.orderedAmount.Editor.Text()
	t2 := m.invoicedAmount.Editor.Text()
	m.orderedAmount.Editor.SetText(t2)
	m.invoicedAmount.Editor.SetText(t1)

	if m.mkt == nil {
		return
	}
	m.validateForm()

	mapAssets := m.Dexc().Core().SupportedAssets()
	baseAsset := mapAssets[m.mkt.BaseID]
	quoteAsset := mapAssets[m.mkt.QuoteID]
	baseSymbol := m.mkt.BaseSymbol
	quoteSymbol := m.mkt.QuoteSymbol
	if baseAsset.Wallet == nil {
		baseSymbol = fmt.Sprintf("Add %s wallet", baseAsset.Symbol)
	}
	if quoteAsset.Wallet == nil {
		quoteSymbol = fmt.Sprintf("Add %s wallet", quoteAsset.Symbol)
	}

	if m.isSell {
		m.invoicedAmount.CustomButton.Text = strings.ToUpper(baseSymbol)
		m.invoicedAmount.CustomButton.Background = m.Theme.Color.Primary
		m.orderedAmount.CustomButton.Text = strings.ToUpper(quoteSymbol)
		m.orderedAmount.CustomButton.Background = m.Theme.Color.Success
	} else {
		m.invoicedAmount.CustomButton.Text = strings.ToUpper(quoteSymbol)
		m.invoicedAmount.CustomButton.Background = m.Theme.Color.Success
		m.orderedAmount.CustomButton.Text = strings.ToUpper(baseSymbol)
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

func (m *miniTradeFormWidget) validateForm() {
	m.submitBtn.SetEnabled(false)
	m.invoicedAmount.SetError("")
	if m.invoicedAmount.Editor.Text() == "" ||
		m.orderedAmount.Editor.Text() == "" ||
		!m.Dexc().HasWallet(int32(m.mkt.QuoteID)) ||
		!m.Dexc().HasWallet(int32(m.mkt.BaseID)) {
		return
	}

	baseWallState := m.Dexc().Core().WalletState(m.mkt.BaseID)
	quoteWallState := m.Dexc().Core().WalletState(m.mkt.QuoteID)

	if !baseWallState.Synced || !quoteWallState.Synced {
		return
	}

	if m.isSell {
		amount, err := strconv.ParseUint(m.invoicedAmount.Editor.Text(), 10, 64)
		if err != nil {
			m.invoicedAmount.SetError("Invalid amount")
			return
		}

		assetInfo, err := asset.Info(m.mkt.BaseID)
		if err != nil {
			m.invoicedAmount.SetError(err.Error())
			return
		}

		if amount%(m.mkt.LotSize/assetInfo.UnitInfo.Conventional.ConversionFactor) != 0 {
			m.invoicedAmount.SetError("Invalid amount")
			return
		}
	} else {
		_, err := strconv.ParseFloat(m.invoicedAmount.Editor.Text(), 64)
		if err != nil {
			m.invoicedAmount.SetError("Invalid amount")
			return
		}
	}

	m.submitBtn.SetEnabled(true)
}

func (m *miniTradeFormWidget) handle(ord *core.OrderBook) {
	if m.mkt == nil {
		return
	}

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
				m.validateForm()
			}
		}

		if m.orderedAmount.Editor.Focused() {
			_, change := decredmaterial.HandleEditorEvents(m.orderedAmount.Editor)
			if change {
				value := m.suggestValue(ord, m.orderedAmount.Editor.Text(), m.isSell, false)
				m.invoicedAmount.Editor.SetText(value)
				m.validateForm()
			}
		}
	}

	if m.submitBtn.Button.Clicked() {
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
			Title(strAppPassword).
			Hint(strAuthOrderAppPassword).
			Description(strNoteConfirmTradeMessage).
			NegativeButton(values.String(values.StrCancel), func() {}).
			PositiveButton(strOk, func(password string, pm *modal.PasswordModal) bool {
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
					m.clearInputs()
					m.Toast.Notify(strSuccessful)
					pm.Dismiss()
				}()
				return false
			}).Show()
	}

	if m.invoicedAmount.CustomButton.Button.Clicked() {
		// If not have wallet go to create wallet
		// TODO: If have wallet show modal to pick other coin to trade
		if m.isSell {
			if !m.Dexc().HasWallet(int32(m.mkt.BaseID)) {
				m.doCreateWallet(m.mkt.BaseSymbol, m.mkt.BaseID)
			}
		} else {
			if !m.Dexc().HasWallet(int32(m.mkt.QuoteID)) {
				m.doCreateWallet(m.mkt.QuoteSymbol, m.mkt.QuoteID)
			}
		}
	}

	if m.orderedAmount.CustomButton.Button.Clicked() {
		if m.isSell {
			if !m.Dexc().HasWallet(int32(m.mkt.QuoteID)) {
				m.doCreateWallet(m.mkt.QuoteSymbol, m.mkt.QuoteID)
			}
		} else {
			if !m.Dexc().HasWallet(int32(m.mkt.BaseID)) {
				m.doCreateWallet(m.mkt.BaseSymbol, m.mkt.BaseID)
			}
		}
	}
}

func (m *miniTradeFormWidget) clearInputs() {
	m.orderedAmount.Editor.SetText("")
	m.invoicedAmount.Editor.SetText("")
	m.orderedAmount.SetError("")
	m.invoicedAmount.SetError("")
	m.submitBtn.SetEnabled(false)
}

func (m *miniTradeFormWidget) doCreateWallet(symbol string, coinID uint32) {
	newCreateWalletModal(m.Load,
		&walletInfoWidget{
			image:    components.CoinImageBySymbol(&m.Icons, symbol),
			coinName: symbol,
			coinID:   coinID,
		},
		"",
		func(_ *createWalletModal) {
			m.RefreshWindow()
		}).Show()
}

package ui

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"time"

	"github.com/raedahgroup/godcr/ui/values"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/atotto/clipboard"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type amountValue struct {
	activeTotalAmount           string
	inactiveTotalAmount         string
	activeTransactionFeeValue   string
	inactiveTransactionFeeValue string
	activeTotalCostValue        string
	inactiveTotalCostValue      string
}

type SendPage struct {
	pageContainer layout.List
	theme         *decredmaterial.Theme

	txAuthor        *dcrlibwallet.TxAuthor
	broadcastResult *wallet.Broadcast

	wallet          *wallet.Wallet
	selectedWallet  wallet.InfoShort
	selectedAccount wallet.Account

	destinationAddressEditor     decredmaterial.Editor
	sendAmountEditor             decredmaterial.Editor
	nextButton                   decredmaterial.Button
	closeConfirmationModalButton decredmaterial.Button
	confirmButton                decredmaterial.Button

	confirmModal *decredmaterial.Modal

	copyIcon     decredmaterial.IconButton
	currencySwap decredmaterial.IconButton

	remainingBalance int64
	amountAtoms      int64
	totalCostDCR     int64
	txFee            int64

	usdExchangeRate float64
	inputAmount     float64
	amountUSDtoDCR  float64
	amountDCRtoUSD  float64

	count int

	sendErrorText      string
	txHashText         string
	txHash             string
	calculateErrorText string

	activeTotalAmount   string
	inactiveTotalAmount string

	activeExchange   string
	inactiveExchange string

	activeTransactionFeeValue   string
	inactiveTransactionFeeValue string

	activeTotalCostValue   string
	inactiveTotalCostValue string

	balanceAfterSendValue string

	LastTradeRate string

	passwordModal *decredmaterial.Password

	isConfirmationModalOpen   bool
	isPasswordModalOpen       bool
	isBroadcastingTransaction bool
	shouldInitializeTxAuthor  bool
	hasCopiedTxHash           bool

	txAuthorErrChan  chan error
	broadcastErrChan chan error
}

const (
	PageSend = "send"
)

func (win *Window) SendPage(common pageCommon) layout.Widget {
	page := &SendPage{
		pageContainer: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},

		theme:           common.theme,
		wallet:          common.wallet,
		txAuthor:        &win.txAuthor,
		broadcastResult: &win.broadcastResult,

		sendErrorText: "",
		txHashText:    "",

		activeExchange:   "DCR",
		inactiveExchange: "USD",

		closeConfirmationModalButton: common.theme.Button(new(widget.Clickable), "Close"),
		nextButton:                   common.theme.Button(new(widget.Clickable), "Next"),
		confirmButton:                common.theme.Button(new(widget.Clickable), "Confirm"),

		confirmModal: common.theme.Modal("Confirm Send Transaction"),

		copyIcon: common.theme.IconButton(new(widget.Clickable), mustIcon(widget.NewIcon(icons.ContentContentCopy))),

		isConfirmationModalOpen:   false,
		isPasswordModalOpen:       false,
		hasCopiedTxHash:           false,
		isBroadcastingTransaction: false,

		passwordModal:    common.theme.Password(),
		broadcastErrChan: make(chan error),
		txAuthorErrChan:  make(chan error),
	}

	page.balanceAfterSendValue = "- DCR"

	page.destinationAddressEditor = common.theme.Editor(new(widget.Editor), "Destination Address")
	page.destinationAddressEditor.SetRequiredErrorText("")
	page.destinationAddressEditor.IsRequired = true
	page.destinationAddressEditor.IsVisible = true

	page.sendAmountEditor = common.theme.Editor(new(widget.Editor), "Amount to be sent")
	page.sendAmountEditor.SetRequiredErrorText("")
	page.sendAmountEditor.IsRequired = true
	page.sendAmountEditor.IsTitleLabel = false

	page.closeConfirmationModalButton.Background = common.theme.Color.Gray
	page.destinationAddressEditor.Editor.SetText("")

	page.copyIcon.Background = common.theme.Color.Background
	page.copyIcon.Color = common.theme.Color.Text
	page.copyIcon.Size = values.MarginPadding35
	page.copyIcon.Inset = layout.UniformInset(values.MarginPadding5)

	page.currencySwap = common.theme.IconButton(new(widget.Clickable), common.icons.actionSwapVert)
	page.currencySwap.Background = color.RGBA{}
	page.currencySwap.Color = common.theme.Color.Text
	page.currencySwap.Inset = layout.UniformInset(values.MarginPadding0)
	page.currencySwap.Size = values.MarginPadding30
	go common.wallet.GetUSDExchangeValues(&page)

	return func(gtx C) D {
		page.Handle(common)
		return page.Layout(gtx, common)
	}
}

func (pg *SendPage) Handle(c pageCommon) {
	if len(c.info.Wallets) == 0 {
		return
	}

	if pg.LastTradeRate == "" && pg.count == 0 {
		pg.count = 1
		pg.calculateValues()
	}

	if (pg.LastTradeRate != "" && pg.count == 0) || (pg.LastTradeRate != "" && pg.count == 1) {
		pg.count = 2
		pg.calculateValues()
	}

	if pg.selectedAccount.CurrentAddress != c.info.Wallets[*c.selectedWallet].Accounts[*c.selectedAccount].CurrentAddress {
		pg.shouldInitializeTxAuthor = true
		pg.selectedAccount = c.info.Wallets[*c.selectedWallet].Accounts[*c.selectedAccount]
	}

	if pg.selectedWallet.ID != c.info.Wallets[*c.selectedWallet].ID {
		pg.shouldInitializeTxAuthor = true
		pg.selectedWallet = c.info.Wallets[*c.selectedWallet]
	}

	if pg.shouldInitializeTxAuthor {
		pg.shouldInitializeTxAuthor = false
		pg.sendAmountEditor.Editor.SetText("")
		pg.calculateErrorText = ""
		pg.sendErrorText = ""
		c.wallet.CreateTransaction(pg.selectedWallet.ID, pg.selectedAccount.Number, pg.txAuthorErrChan)
	}

	pg.validate(true)
	pg.watchForBroadcastResult()

	if pg.hasCopiedTxHash {
		time.AfterFunc(3*time.Second, func() {
			pg.hasCopiedTxHash = false
		})
	}

	if pg.isBroadcastingTransaction {
		col := pg.theme.Color.Gray
		col.A = 150
		pg.confirmButton.Text = "Sending..."
		pg.confirmButton.Background = col
	} else {
		pg.confirmButton.Text = "Send"
		pg.confirmButton.Background = pg.theme.Color.Primary
	}

	for pg.nextButton.Button.Clicked() {
		if pg.validate(false) && pg.calculateErrorText == "" {
			pg.isConfirmationModalOpen = true
		}
	}

	for pg.confirmButton.Button.Clicked() {
		pg.sendErrorText = ""
		pg.isPasswordModalOpen = true
	}

	for pg.closeConfirmationModalButton.Button.Clicked() {
		pg.sendErrorText = ""
		pg.isConfirmationModalOpen = false
	}

	for pg.currencySwap.Button.Clicked() {
		if pg.LastTradeRate != "" {
			if pg.activeExchange == "DCR" {
				pg.activeExchange = "USD"
				pg.inactiveExchange = "DCR"
			} else {
				pg.activeExchange = "DCR"
				pg.inactiveExchange = "USD"
			}
		}

		pg.calculateValues()
	}

	for _, evt := range pg.destinationAddressEditor.Editor.Events() {
		go pg.calculateValues()
		pg.changeEvt(evt)
	}

	if pg.destinationAddressEditor.Editor.Len() == 0 || pg.sendAmountEditor.Editor.Len() == 0 {
		pg.balanceAfterSend(pg.selectedAccount.SpendableBalance)
	}

	for _, evt := range pg.sendAmountEditor.Editor.Events() {
		go pg.calculateValues()
		pg.changeEvt(evt)
	}

	for pg.copyIcon.Button.Clicked() {
		clipboard.WriteAll(pg.txHash)
		pg.hasCopiedTxHash = true
	}

	select {
	case err := <-pg.txAuthorErrChan:
		pg.calculateErrorText = err.Error()
	case err := <-pg.broadcastErrChan:
		pg.sendErrorText = err.Error()
		pg.isBroadcastingTransaction = false
	default:
	}
}

func (pg *SendPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	if len(common.info.Wallets) == 0 {
		return layout.Dimensions{}
	}

	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.drawSuccessSection(gtx)
		},
		func(gtx C) D {
			return pg.drawCopiedLabelSection(gtx)
		},
		func(gtx C) D {
			return pg.drawSelectedAccountSection(gtx)
		},
		func(gtx C) D {
			return pg.destinationAddressEditor.Layout(gtx)
		},
		func(gtx C) D {
			return pg.sendAmountLayout(gtx)
		},
		func(gtx C) D {
			return pg.drawTransactionDetailWidgets(gtx)
		},
		func(gtx C) D {
			if pg.calculateErrorText != "" {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return pg.theme.ErrorAlert(gtx, pg.calculateErrorText)
			}
			return layout.Dimensions{}
		},
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return pg.nextButton.Layout(gtx)
		},
	}

	dims := common.LayoutWithAccounts(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
			return pg.pageContainer.Layout(gtx, len(pageContent), func(gtx C, i int) D {
				return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, pageContent[i])
			})
		})
	})

	if pg.isConfirmationModalOpen && pg.isPasswordModalOpen {
		return common.Modal(gtx, dims, pg.drawPasswordModal(gtx))
	} else if pg.isConfirmationModalOpen {
		return common.Modal(gtx, dims, pg.drawConfirmationModal(gtx))
	}

	return dims
}

func (pg *SendPage) drawSuccessSection(gtx layout.Context) layout.Dimensions {
	if pg.txHashText != "" {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(0.99, func(gtx C) D {
				return pg.theme.SuccessAlert(gtx, pg.txHashText)
			}),
			layout.Rigid(func(gtx C) D {
				inset := layout.Inset{
					Left: values.MarginPadding5,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return pg.copyIcon.Layout(gtx)
				})
			}),
		)
	}
	return layout.Dimensions{}
}

func (pg *SendPage) drawCopiedLabelSection(gtx layout.Context) layout.Dimensions {
	if pg.hasCopiedTxHash {
		return pg.theme.Caption("copied").Layout(gtx)
	}
	return layout.Dimensions{}
}

func (pg *SendPage) drawSelectedAccountSection(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				selectedDetails := func(gtx C) D {
					return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
											return pg.theme.Body2(pg.selectedAccount.Name).Layout(gtx)
										})
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
											return pg.theme.Body2(dcrutil.Amount(pg.selectedAccount.SpendableBalance).String()).Layout(gtx)
										})
									}),
								)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									return pg.theme.Body2(pg.selectedWallet.Name).Layout(gtx)
								})
							}),
						)
					})
				}
				return decredmaterial.Card{Color: pg.theme.Color.Surface}.Layout(gtx, selectedDetails)
			}),
		)
	})
}

func (pg *SendPage) drawTransactionDetailWidgets(gtx layout.Context) layout.Dimensions {
	w := []func(gtx C) D{
		func(gtx C) D {
			return pg.tableLayout(gtx, pg.theme.Body2("Transaction Fee"), pg.activeTransactionFeeValue, pg.inactiveTransactionFeeValue)
		},
		func(gtx C) D {
			return pg.tableLayout(gtx, pg.theme.Body2("Total Cost"), pg.activeTotalCostValue, pg.inactiveTotalCostValue)
		},
		func(gtx C) D {
			return pg.tableLayout(gtx, pg.theme.Body2("Balance after send"), pg.balanceAfterSendValue, "")
		},
	}

	list := layout.List{Axis: layout.Vertical}
	return list.Layout(gtx, len(w), func(gtx C, i int) D {
		inset := layout.Inset{
			Top: values.MarginPadding10,
		}
		return inset.Layout(gtx, w[i])
	})
}

func (pg *SendPage) tableLayout(gtx layout.Context, leftLabel decredmaterial.Label, active, inactive string) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return leftLabel.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.theme.Body1(active).Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						b := pg.theme.Body1(inactive)
						b.Color = pg.theme.Color.Hint
						inset := layout.Inset{
							Left: values.MarginPadding5,
						}
						return inset.Layout(gtx, func(gtx C) D {
							return b.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

func (pg *SendPage) sendAmountLayout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.W.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.theme.H6(pg.activeTotalAmount).Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								m := values.MarginPadding10
								return layout.Inset{Left: m, Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
									return pg.currencySwap.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									return pg.sendAmountEditor.Layout(gtx)
								})
							}),
						)
					}),
					layout.Rigid(func(gtx C) D {
						txt := pg.theme.Body2(pg.inactiveTotalAmount)
						if pg.LastTradeRate == "" {
							txt.Color = pg.theme.Color.Danger
						}
						return txt.Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (pg *SendPage) drawConfirmationModal(gtx layout.Context) layout.Dimensions {
	if !pg.isConfirmationModalOpen {
		return layout.Dimensions{}
	}

	w := []func(gtx C) D{
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			if pg.sendErrorText != "" {
				return pg.theme.ErrorAlert(gtx, pg.sendErrorText)
			}
			return layout.Dimensions{
				Size: image.Point{
					X: gtx.Constraints.Max.X,
				},
			}
		},
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return pg.theme.Body1(fmt.Sprintf("Sending from %s (%s)", pg.selectedAccount.Name, pg.selectedWallet.Name)).Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return pg.theme.Body2("To destination address").Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return pg.theme.Body1(pg.destinationAddressEditor.Editor.Text()).Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return pg.drawTransactionDetailWidgets(gtx)
		},
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.Y
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return pg.theme.Caption("Your DCR will be sent and CANNOT be undone").Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.confirmButton.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					inset := layout.Inset{
						Left: values.MarginPadding5,
					}
					return inset.Layout(gtx, func(gtx C) D {
						return pg.closeConfirmationModalButton.Layout(gtx)
					})
				}),
			)
		},
	}
	return pg.confirmModal.Layout(gtx, w, 850)
}

func (pg *SendPage) drawPasswordModal(gtx layout.Context) layout.Dimensions {
	return pg.passwordModal.Layout(gtx, func(password []byte) {
		pg.isBroadcastingTransaction = true
		pg.isPasswordModalOpen = false

		pg.wallet.BroadcastTransaction(pg.txAuthor, password, pg.broadcastErrChan)
	}, func() {
		pg.isPasswordModalOpen = false
	})
}

func (pg *SendPage) validate(ignoreEmpty bool) bool {
	isAddressValid := pg.validateDestinationAddress(ignoreEmpty)
	isAmountValid := pg.validateAmount(ignoreEmpty)

	if !isAddressValid || !isAmountValid || pg.calculateErrorText != "" {
		pg.nextButton.Background = pg.theme.Color.Hint
		return false
	}

	pg.nextButton.Background = pg.theme.Color.Primary
	return true
}

func (pg *SendPage) validateDestinationAddress(ignoreEmpty bool) bool {
	pg.destinationAddressEditor.ClearError()
	destinationAddress := pg.destinationAddressEditor.Editor.Text()
	if destinationAddress == "" && !ignoreEmpty {
		pg.destinationAddressEditor.SetError("please enter a destination address")
		return false
	}

	if destinationAddress != "" {
		isValid, _ := pg.wallet.IsAddressValid(destinationAddress)
		if !isValid {
			pg.destinationAddressEditor.SetError("invalid address")
			return false
		}
	}
	return true
}

func (pg *SendPage) validateAmount(ignoreEmpty bool) bool {
	pg.sendAmountEditor.ClearError()
	amount := pg.sendAmountEditor.Editor.Text()
	if amount == "" {
		if !ignoreEmpty {
			pg.sendAmountEditor.SetError("please enter a send amount")
		}
		return false
	}

	if amount != "" {
		_, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			pg.sendAmountEditor.SetError("please enter a valid amount")
			return false
		}
	}

	return true
}

func (pg *SendPage) calculateValues() {
	defaultActiveValues := fmt.Sprintf("- %s", pg.activeExchange)
	defaultInactiveValues := fmt.Sprintf("(- %s)", pg.inactiveExchange)
	noExchangeText := "Exchange rate not fetched"

	pg.activeTransactionFeeValue = defaultActiveValues
	pg.activeTotalCostValue = defaultActiveValues
	pg.inactiveTransactionFeeValue = defaultInactiveValues
	pg.inactiveTotalCostValue = defaultInactiveValues

	pg.calculateErrorText = ""
	pg.activeTotalAmount = defaultActiveValues
	pg.inactiveTotalAmount = fmt.Sprintf("- %s", pg.inactiveExchange)

	// default values when exchange is not available
	if pg.LastTradeRate == "" {
		pg.activeTransactionFeeValue = defaultActiveValues
		pg.activeTotalCostValue = defaultActiveValues
		pg.inactiveTransactionFeeValue = ""
		pg.inactiveTotalCostValue = ""
		pg.activeTotalAmount = defaultActiveValues
		pg.inactiveTotalAmount = noExchangeText
	}

	if pg.txAuthor == nil || !pg.validate(true) {
		return
	}

	pg.inputAmount, _ = strconv.ParseFloat(pg.sendAmountEditor.Editor.Text(), 64)

	if pg.LastTradeRate != "" {
		pg.usdExchangeRate, _ = strconv.ParseFloat(pg.LastTradeRate, 64)
		pg.amountUSDtoDCR = pg.inputAmount / pg.usdExchangeRate
		pg.amountDCRtoUSD = pg.inputAmount * pg.usdExchangeRate
	}

	if pg.activeExchange == "USD" && pg.LastTradeRate != "" {
		pg.amountAtoms = pg.setDestinationAddr(pg.amountUSDtoDCR)
		if pg.amountAtoms == 0 {
			return
		}
	} else {
		pg.amountAtoms = pg.setDestinationAddr(pg.inputAmount)
		if pg.amountAtoms == 0 {
			return
		}
	}

	pg.txFee = pg.getTxFee()
	if pg.txFee == 0 {
		return
	}

	pg.totalCostDCR = pg.txFee + pg.amountAtoms

	pg.updateDefaultValues()
	pg.balanceAfterSend(pg.totalCostDCR)
}

func (pg *SendPage) setDestinationAddr(sendAmount float64) int64 {
	amount, err := dcrutil.NewAmount(sendAmount)
	if err != nil {
		pg.calculateErrorText = fmt.Sprintf("error estimating transaction fee: %s", err)
		return 0
	}

	pg.amountAtoms = int64(amount)
	pg.txAuthor.RemoveSendDestination(0)
	pg.txAuthor.AddSendDestination(pg.destinationAddressEditor.Editor.Text(), pg.amountAtoms, false)
	return pg.amountAtoms
}

func (pg *SendPage) amountValues() amountValue {
	txFeeValueUSD := dcrutil.Amount(pg.txFee).ToCoin() * pg.usdExchangeRate
	switch {
	case pg.activeExchange == "USD" && pg.LastTradeRate != "":
		return amountValue{
			activeTotalAmount:           fmt.Sprintf("%s USD", pg.sendAmountEditor.Editor.Text()),
			inactiveTotalAmount:         dcrutil.Amount(pg.amountAtoms).String(),
			activeTransactionFeeValue:   fmt.Sprintf("%f USD", txFeeValueUSD),
			inactiveTransactionFeeValue: fmt.Sprintf("(%s)", dcrutil.Amount(pg.txFee).String()),
			activeTotalCostValue:        fmt.Sprintf("%s USD", strconv.FormatFloat(pg.inputAmount+txFeeValueUSD, 'f', 7, 64)),
			inactiveTotalCostValue:      fmt.Sprintf("(%s )", dcrutil.Amount(pg.totalCostDCR).String()),
		}
	case pg.activeExchange == "DCR" && pg.LastTradeRate != "":
		return amountValue{
			activeTotalAmount:           dcrutil.Amount(pg.amountAtoms).String(),
			inactiveTotalAmount:         fmt.Sprintf("%s USD", strconv.FormatFloat(pg.amountDCRtoUSD, 'f', 7, 64)),
			activeTransactionFeeValue:   dcrutil.Amount(pg.txFee).String(),
			inactiveTransactionFeeValue: fmt.Sprintf("(%f USD)", txFeeValueUSD),
			activeTotalCostValue:        dcrutil.Amount(pg.totalCostDCR).String(),
			inactiveTotalCostValue:      fmt.Sprintf("(%s USD)", strconv.FormatFloat(pg.amountDCRtoUSD+txFeeValueUSD, 'f', 7, 64)),
		}
	default:
		return amountValue{
			activeTotalAmount:         dcrutil.Amount(pg.amountAtoms).String(),
			inactiveTotalAmount:       "Exchange rate not fetched",
			activeTransactionFeeValue: dcrutil.Amount(pg.txFee).String(),
			activeTotalCostValue:      dcrutil.Amount(pg.totalCostDCR).String(),
		}
	}
}

func (pg *SendPage) updateDefaultValues() {
	v := pg.amountValues()
	pg.activeTotalAmount = v.activeTotalAmount
	pg.inactiveTotalAmount = v.inactiveTotalAmount
	pg.activeTransactionFeeValue = v.activeTransactionFeeValue
	pg.inactiveTransactionFeeValue = v.inactiveTransactionFeeValue
	pg.activeTotalCostValue = v.activeTotalCostValue
	pg.inactiveTotalCostValue = v.inactiveTotalCostValue
}

func (pg *SendPage) getTxFee() int64 {
	// calculate transaction fee
	feeAndSize, err := pg.txAuthor.EstimateFeeAndSize()
	if err != nil {
		pg.calculateErrorText = fmt.Sprintf("error estimating transaction fee: %s", err)
		return 0
	}

	return feeAndSize.Fee.AtomValue
}

func (pg *SendPage) balanceAfterSend(totalCost int64) {
	pg.remainingBalance = pg.selectedWallet.SpendableBalance - totalCost
	pg.balanceAfterSendValue = dcrutil.Amount(pg.remainingBalance).String()
}

func (pg *SendPage) watchForBroadcastResult() {
	if pg.broadcastResult == nil {
		return
	}

	if pg.broadcastResult.TxHash != "" {
		if pg.remainingBalance != -1 {
			pg.selectedAccount.SpendableBalance = pg.remainingBalance
		}
		pg.remainingBalance = -1

		pg.txHash = pg.broadcastResult.TxHash
		pg.txHashText = fmt.Sprintf("Successful. Hash: %s", pg.broadcastResult.TxHash)
		pg.destinationAddressEditor.Editor.SetText("")
		pg.sendAmountEditor.Editor.SetText("")
		pg.isConfirmationModalOpen = false
		pg.isBroadcastingTransaction = false
		pg.broadcastResult.TxHash = ""
	}
}

func (pg *SendPage) changeEvt(evt widget.EditorEvent) {
	switch evt.(type) {
	case widget.ChangeEvent:
		go pg.wallet.GetUSDExchangeValues(&pg)
	}
}

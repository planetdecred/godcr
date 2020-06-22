package ui

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
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

	destinationAddressEditor           *widget.Editor
	sendAmountEditor                   *widget.Editor
	nextButtonWidget                   *widget.Button
	closeConfirmationModalButtonWidget *widget.Button
	confirmButtonWidget                *widget.Button
	copyIconWidget                     *widget.Button
	currencySwapWidget                 widget.Button

	destinationAddressEditorMaterial     decredmaterial.Editor
	sendAmountEditorMaterial             decredmaterial.Editor
	nextButtonMaterial                   decredmaterial.Button
	closeConfirmationModalButtonMaterial decredmaterial.Button
	confirmButtonMaterial                decredmaterial.Button

	copyIconMaterial decredmaterial.IconButton
	currencySwap     decredmaterial.IconButton

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

		destinationAddressEditor:           new(widget.Editor),
		sendAmountEditor:                   new(widget.Editor),
		nextButtonWidget:                   new(widget.Button),
		closeConfirmationModalButtonWidget: new(widget.Button),
		confirmButtonWidget:                new(widget.Button),
		copyIconWidget:                     new(widget.Button),

		sendErrorText: "",
		txHashText:    "",

		activeExchange:   "DCR",
		inactiveExchange: "USD",

		closeConfirmationModalButtonMaterial: common.theme.Button("Close"),
		nextButtonMaterial:                   common.theme.Button("Next"),
		confirmButtonMaterial:                common.theme.Button("Confirm"),

		copyIconMaterial: common.theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentContentCopy))),

		isConfirmationModalOpen:   false,
		isPasswordModalOpen:       false,
		hasCopiedTxHash:           false,
		isBroadcastingTransaction: false,

		passwordModal:    common.theme.Password(),
		broadcastErrChan: make(chan error),
		txAuthorErrChan:  make(chan error),
	}

	page.balanceAfterSendValue = "- DCR"

	page.destinationAddressEditorMaterial = common.theme.Editor("Destination Address")
	page.destinationAddressEditorMaterial.SetRequiredErrorText("")
	page.destinationAddressEditorMaterial.IsRequired = true
	page.destinationAddressEditorMaterial.IsVisible = true

	page.sendAmountEditorMaterial = common.theme.Editor("Amount to be sent")
	page.sendAmountEditorMaterial.SetRequiredErrorText("")
	page.sendAmountEditorMaterial.IsRequired = true
	page.sendAmountEditorMaterial.IsTitleLabel = false

	page.closeConfirmationModalButtonMaterial.Background = common.theme.Color.Gray
	page.destinationAddressEditor.SetText("")

	page.copyIconMaterial.Background = common.theme.Color.Background
	page.copyIconMaterial.Color = common.theme.Color.Text
	page.copyIconMaterial.Size = unit.Dp(35)
	page.copyIconMaterial.Padding = unit.Dp(5)

	page.currencySwap = common.theme.IconButton(common.icons.actionSwapVert)
	page.currencySwap.Background = color.RGBA{}
	page.currencySwap.Color = common.theme.Color.Text
	page.currencySwap.Padding = unit.Dp(0)
	page.currencySwap.Size = unit.Dp(30)
	go common.wallet.GetUSDExchangeValues(&page)

	return func() {
		page.Layout(common)
		page.drawConfirmationModal(common)
		page.drawPasswordModal(common)
		page.Handle(common)
	}
}

func (pg *SendPage) Handle(c pageCommon) {
	if len(c.info.Wallets) == 0 {
		return
	}

	gtx := c.gtx

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
		pg.sendAmountEditor.SetText("")
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
		pg.confirmButtonMaterial.Text = "Sending..."
		pg.confirmButtonMaterial.Background = col
	} else {
		pg.confirmButtonMaterial.Text = "Send"
		pg.confirmButtonMaterial.Background = pg.theme.Color.Primary
	}

	for pg.nextButtonWidget.Clicked(gtx) {
		if pg.validate(false) && pg.calculateErrorText == "" {
			pg.isConfirmationModalOpen = true
		}
	}

	for pg.confirmButtonWidget.Clicked(gtx) {
		pg.sendErrorText = ""
		pg.isPasswordModalOpen = true
	}

	for pg.closeConfirmationModalButtonWidget.Clicked(gtx) {
		pg.sendErrorText = ""
		pg.isConfirmationModalOpen = false
	}

	for pg.currencySwapWidget.Clicked(gtx) {
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

	for _, evt := range pg.destinationAddressEditor.Events(gtx) {
		go pg.calculateValues()
		pg.changeEvt(evt)
	}

	if pg.destinationAddressEditor.Len() == 0 || pg.sendAmountEditor.Len() == 0 {
		pg.balanceAfterSend(pg.selectedAccount.SpendableBalance)
	}

	for _, evt := range pg.sendAmountEditor.Events(gtx) {
		go pg.calculateValues()
		pg.changeEvt(evt)
	}

	for pg.copyIconWidget.Clicked(gtx) {
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

func (pg *SendPage) Layout(common pageCommon) {
	if len(common.info.Wallets) == 0 {
		return
	}

	pageContent := []func(){
		func() {
			pg.drawSuccessSection(common.gtx)
		},
		func() {
			pg.drawCopiedLabelSection(common.gtx)
		},
		func() {
			pg.drawSelectedAccountSection(common.gtx)
		},
		func() {
			pg.destinationAddressEditorMaterial.Layout(common.gtx, pg.destinationAddressEditor)
		},
		func() {
			pg.sendAmountLayout(common.gtx)
		},
		func() {
			pg.drawTransactionDetailWidgets(common.gtx)
		},
		func() {
			if pg.calculateErrorText != "" {
				common.gtx.Constraints.Width.Min = common.gtx.Constraints.Width.Max
				pg.theme.ErrorAlert(common.gtx, pg.calculateErrorText)
			}
		},
		func() {
			common.gtx.Constraints.Width.Min = common.gtx.Constraints.Width.Max
			pg.nextButtonMaterial.Layout(common.gtx, pg.nextButtonWidget)
		},
	}

	common.LayoutWithAccounts(common.gtx, func() {
		layout.Inset{Right: unit.Dp(110)}.Layout(common.gtx, func() {
			pg.pageContainer.Layout(common.gtx, len(pageContent), func(i int) {
				layout.Inset{Top: unit.Dp(5)}.Layout(common.gtx, pageContent[i])
			})
		})
	})
}

func (pg *SendPage) drawSuccessSection(gtx *layout.Context) {
	if pg.txHashText != "" {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(0.99, func() {
				pg.theme.SuccessAlert(gtx, pg.txHashText)
			}),
			layout.Rigid(func() {
				inset := layout.Inset{
					Left: unit.Dp(3),
				}
				inset.Layout(gtx, func() {
					pg.copyIconMaterial.Layout(gtx, pg.copyIconWidget)
				})
			}),
		)
	}
}

func (pg *SendPage) drawCopiedLabelSection(gtx *layout.Context) {
	if pg.hasCopiedTxHash {
		pg.theme.Caption("copied").Layout(gtx)
	}
}

func (pg *SendPage) drawSelectedAccountSection(gtx *layout.Context) {
	layout.Center.Layout(gtx, func() {
		layout.Stack{}.Layout(gtx,
			layout.Stacked(func() {
				selectedDetails := func() {
					layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
						layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func() {
								layout.Flex{}.Layout(gtx,
									layout.Rigid(func() {
										layout.Inset{Bottom: unit.Dp(5)}.Layout(gtx, func() {
											pg.theme.Body2(pg.selectedAccount.Name).Layout(gtx)
										})
									}),
									layout.Rigid(func() {
										layout.Inset{Left: unit.Dp(20)}.Layout(gtx, func() {
											pg.theme.Body2(dcrutil.Amount(pg.selectedAccount.SpendableBalance).String()).Layout(gtx)
										})
									}),
								)
							}),
							layout.Rigid(func() {
								layout.Inset{Bottom: unit.Dp(5)}.Layout(gtx, func() {
									pg.theme.Body2(pg.selectedWallet.Name).Layout(gtx)
								})
							}),
						)
					})
				}
				decredmaterial.Card{}.Layout(gtx, selectedDetails)
			}),
		)
	})
}

func (pg *SendPage) drawTransactionDetailWidgets(gtx *layout.Context) {
	w := []func(){
		func() {
			pg.tableLayout(gtx, pg.theme.Body2("Transaction Fee"), pg.activeTransactionFeeValue, pg.inactiveTransactionFeeValue)
		},
		func() {
			pg.tableLayout(gtx, pg.theme.Body2("Total Cost"), pg.activeTotalCostValue, pg.inactiveTotalCostValue)
		},
		func() {
			pg.tableLayout(gtx, pg.theme.Body2("Balance after send"), pg.balanceAfterSendValue, "")
		},
	}

	list := layout.List{Axis: layout.Vertical}
	list.Layout(gtx, len(w), func(i int) {
		inset := layout.Inset{
			Top: unit.Dp(10),
		}
		inset.Layout(gtx, w[i])
	})
}

func (pg *SendPage) tableLayout(gtx *layout.Context, leftLabel decredmaterial.Label, active, inactive string) {
	layout.Flex{}.Layout(gtx,
		layout.Rigid(func() {
			leftLabel.Layout(gtx)
		}),
		layout.Flexed(1, func() {
			layout.E.Layout(gtx, func() {
				layout.Flex{}.Layout(gtx,
					layout.Rigid(func() {
						pg.theme.Body1(active).Layout(gtx)
					}),
					layout.Rigid(func() {
						b := pg.theme.Body1(inactive)
						b.Color = pg.theme.Color.Hint
						inset := layout.Inset{
							Left: unit.Dp(5),
						}
						inset.Layout(gtx, func() {
							b.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

func (pg *SendPage) sendAmountLayout(gtx *layout.Context) {
	layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func() {
			layout.W.Layout(gtx, func() {
				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func() {
						pg.theme.H6(pg.activeTotalAmount).Layout(gtx)
					}),
					layout.Rigid(func() {
						layout.Flex{}.Layout(gtx,
							layout.Rigid(func() {
								layout.Inset{Left: unit.Dp(10), Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(gtx, func() {
									pg.currencySwap.Layout(gtx, &pg.currencySwapWidget)
								})
							}),
							layout.Rigid(func() {
								layout.Inset{Left: unit.Dp(7)}.Layout(gtx, func() {
									pg.sendAmountEditorMaterial.Layout(gtx, pg.sendAmountEditor)
								})
							}),
						)
					}),
					layout.Rigid(func() {
						txt := pg.theme.Body2(pg.inactiveTotalAmount)
						if pg.LastTradeRate == "" {
							txt.Color = pg.theme.Color.Danger
						}
						txt.Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (pg *SendPage) drawConfirmationModal(c pageCommon) {
	if !pg.isConfirmationModalOpen {
		return
	}
	gtx := c.gtx
	w := []func(){
		func() {
			if pg.sendErrorText != "" {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				pg.theme.ErrorAlert(gtx, pg.sendErrorText)
			}
		},
		func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func() {
					pg.theme.Body1(fmt.Sprintf("Sending from %s (%s)", pg.selectedAccount.Name, pg.selectedWallet.Name)).Layout(gtx)
				}),
			)
		},
		func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func() {
					pg.theme.Body2("To destination address").Layout(gtx)
				}),
			)
		},
		func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func() {
					pg.theme.Body1(pg.destinationAddressEditor.Text()).Layout(gtx)
				}),
			)
		},
		func() {
			pg.drawTransactionDetailWidgets(gtx)
		},
		func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func() {
					pg.theme.Caption("Your DCR will be sent and CANNOT be undone").Layout(gtx)
				}),
			)
		},
		func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					pg.confirmButtonMaterial.Layout(gtx, pg.confirmButtonWidget)
				}),
				layout.Rigid(func() {
					inset := layout.Inset{
						Left: unit.Dp(5),
					}
					inset.Layout(gtx, func() {
						pg.closeConfirmationModalButtonMaterial.Layout(gtx, pg.closeConfirmationModalButtonWidget)
					})
				}),
			)
		},
	}
	pg.theme.Modal(gtx, "Confirm Send Transaction", w)
}

func (pg *SendPage) drawPasswordModal(c pageCommon) {
	if !(pg.isConfirmationModalOpen && pg.isPasswordModalOpen) {
		return
	}
	gtx := c.gtx
	pg.passwordModal.Layout(gtx, func(password []byte) {
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
		pg.nextButtonMaterial.Background = pg.theme.Color.Hint
		return false
	}

	pg.nextButtonMaterial.Background = pg.theme.Color.Primary
	return true
}

func (pg *SendPage) validateDestinationAddress(ignoreEmpty bool) bool {
	pg.destinationAddressEditorMaterial.ClearError()
	destinationAddress := pg.destinationAddressEditor.Text()
	if destinationAddress == "" && !ignoreEmpty {
		pg.destinationAddressEditorMaterial.SetError("please enter a destination address")
		return false
	}

	if destinationAddress != "" {
		isValid, _ := pg.wallet.IsAddressValid(destinationAddress)
		if !isValid {
			pg.destinationAddressEditorMaterial.SetError("invalid address")
			return false
		}
	}
	return true
}

func (pg *SendPage) validateAmount(ignoreEmpty bool) bool {
	pg.sendAmountEditorMaterial.ClearError()
	amount := pg.sendAmountEditor.Text()
	if amount == "" {
		if !ignoreEmpty {
			pg.sendAmountEditorMaterial.SetError("please enter a send amount")
		}
		return false
	}

	if amount != "" {
		_, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			pg.sendAmountEditorMaterial.SetError("please enter a valid amount")
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

	pg.inputAmount, _ = strconv.ParseFloat(pg.sendAmountEditor.Text(), 64)

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
	pg.txAuthor.AddSendDestination(pg.destinationAddressEditor.Text(), pg.amountAtoms, false)
	return pg.amountAtoms
}

func (pg *SendPage) amountValues() amountValue {
	txFeeValueUSD := dcrutil.Amount(pg.txFee).ToCoin() * pg.usdExchangeRate
	switch {
	case pg.activeExchange == "USD" && pg.LastTradeRate != "":
		return amountValue{
			activeTotalAmount:           fmt.Sprintf("%s USD", pg.sendAmountEditor.Text()),
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
		pg.destinationAddressEditor.SetText("")
		pg.sendAmountEditor.SetText("")
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

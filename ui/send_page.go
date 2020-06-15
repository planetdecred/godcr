package ui

import (
	"fmt"
	"image/color"
	"net/http"
	"strconv"
	"time"
	// "io/ioutil"
	"encoding/json"

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

type data struct {
	LastTradeRate string
}

type SendPage struct {
	pageContainer   layout.List
	theme           *decredmaterial.Theme
	wallet          *wallet.Wallet
	wallets         []wallet.InfoShort
	txAuthor        *dcrlibwallet.TxAuthor
	broadcastResult *wallet.Broadcast

	selectedWallet  wallet.InfoShort
	selectedAccount wallet.Account

	destinationAddressEditor           *widget.Editor
	sendAmountEditor                   *widget.Editor
	nextButtonWidget                   *widget.Button
	closeConfirmationModalButtonWidget *widget.Button
	confirmButtonWidget                *widget.Button
	copyIconWidget                     *widget.Button
	currencySwapWidget                 widget.Button

	activeTransactionFeeValueLabel   decredmaterial.Label
	inactiveTransactionFeeValueLabel decredmaterial.Label

	activeTransactionFeeValue   string
	inactiveTransactionFeeValue string

	activeTotalCostValue   string
	inactiveTotalCostValue string

	balanceAfterSendValue string

	activeTotalCostValueLabel   decredmaterial.Label
	inactiveTotalCostValueLabel decredmaterial.Label

	balanceAfterSendValueLabel decredmaterial.Label

	// activeExchLabel decredmaterial.Label
	// inactiveExchLabel decredmaterial.Label

	// transactionFeeValueLabel   decredmaterial.Label
	// totalCostValueLabel        decredmaterial.Label
	// balanceAfterSendValueLabel decredmaterial.Label
	// txFeeLabelUSD decredmaterial.Label
	// totalCostLabelUSD decredmaterial.Label

	destinationAddressEditorMaterial     decredmaterial.Editor
	sendAmountEditorMaterial             decredmaterial.Editor
	nextButtonMaterial                   decredmaterial.Button
	closeConfirmationModalButtonMaterial decredmaterial.Button
	confirmButtonMaterial                decredmaterial.Button
	accountsTab                          *decredmaterial.Tabs
	walletsTab                           *decredmaterial.Tabs

	copyIconMaterial decredmaterial.IconButton
	currencySwap     decredmaterial.IconButton

	remainingBalance   int64
	sendErrorText      string
	txHashText         string
	txHash             string
	calculateErrorText string
	activeAmount       string
	inactiveAmount     string
	selectedExchange   string
	unselectedExchange string

	passwordModal *decredmaterial.Password

	isConfirmationModalOpen   bool
	isPasswordModalOpen       bool
	isBroadcastingTransaction bool
	hasInitializedTxAuthor    bool
	hasCopiedTxHash           bool

	txAuthorErrChan  chan error
	broadcastErrChan chan error

	data data
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
		wallets:         common.info.Wallets,
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

		selectedExchange:   "DCR",
		unselectedExchange: "USD",

		closeConfirmationModalButtonMaterial: common.theme.Button("Close"),
		nextButtonMaterial:                   common.theme.Button("Next"),
		confirmButtonMaterial:                common.theme.Button("Confirm"),
		activeTransactionFeeValueLabel:       common.theme.Body2("- DCR"),
		activeTotalCostValueLabel:            common.theme.Body2("- DCR"),
		balanceAfterSendValueLabel:           common.theme.Body2("- DCR"),
		inactiveTransactionFeeValueLabel:     common.theme.Body2("(- USD)"),
		inactiveTotalCostValueLabel:          common.theme.Body2("(- USD)"),

		copyIconMaterial: common.theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentContentCopy))),

		isConfirmationModalOpen:   false,
		isPasswordModalOpen:       false,
		hasInitializedTxAuthor:    false,
		hasCopiedTxHash:           false,
		isBroadcastingTransaction: false,

		passwordModal: common.theme.Password(),
		accountsTab:   decredmaterial.NewTabs(),
		walletsTab:    decredmaterial.NewTabs(),

		broadcastErrChan: make(chan error),
		txAuthorErrChan:  make(chan error),
	}

	page.walletsTab.Position = decredmaterial.Top
	page.accountsTab.Position = decredmaterial.Top

	page.activeAmount = "- " + page.selectedExchange
	page.inactiveAmount = "- " + page.unselectedExchange

	page.activeTransactionFeeValue = "- " + page.selectedExchange
	page.inactiveTransactionFeeValue = fmt.Sprintf("(- %s)", page.unselectedExchange)

	page.activeTotalCostValue = "- " + page.selectedExchange
	page.inactiveTotalCostValue = fmt.Sprintf("(- %s)", page.unselectedExchange)

	page.balanceAfterSendValue = "- " + page.selectedExchange

	page.inactiveTransactionFeeValueLabel.Color = common.theme.Color.Hint
	page.inactiveTotalCostValueLabel.Color = common.theme.Color.Hint

	page.destinationAddressEditorMaterial = common.theme.Editor("Destination Address")
	page.destinationAddressEditorMaterial.SetRequiredErrorText("")
	page.destinationAddressEditorMaterial.IsRequired = true
	page.destinationAddressEditorMaterial.IsVisible = true

	page.sendAmountEditorMaterial = common.theme.Editor("Amount to be sent")
	page.sendAmountEditorMaterial.SetRequiredErrorText("")
	page.sendAmountEditorMaterial.IsRequired = true
	// page.sendAmountEditorMaterial.IsTitleLabel = false

	page.closeConfirmationModalButtonMaterial.Background = common.theme.Color.Gray
	page.destinationAddressEditor.SetText("")

	page.copyIconMaterial.Background = common.theme.Color.Background
	page.copyIconMaterial.Color = common.theme.Color.Text
	page.copyIconMaterial.Size = unit.Dp(35)
	page.copyIconMaterial.Padding = unit.Dp(5)

	page.currencySwap = common.theme.IconButton(common.icons.actionSwapHoriz)
	page.currencySwap.Background = color.RGBA{}
	page.currencySwap.Color = common.theme.Color.Text
	page.currencySwap.Padding = unit.Dp(0)
	page.currencySwap.Size = unit.Dp(30)

	return func() {
		page.Layout(common)
		page.Handle(common)
	}
}

func (pg *SendPage) Handle(common pageCommon) {
	pg.validate(true)
	pg.watchForBroadcastResult()

	if pg.walletsTab.Changed() {
		pg.selectedWallet = pg.wallets[pg.walletsTab.Selected]
		pg.selectedAccount = pg.selectedWallet.Accounts[0]
		pg.accountsTab.Selected = 0

		pg.setAccountTabs()
		pg.wallet.CreateTransaction(pg.selectedWallet.ID, pg.selectedAccount.Number, pg.txAuthorErrChan)

		pg.sendAmountEditor.SetText("")
		pg.calculateErrorText = ""
		pg.sendErrorText = ""
	}

	if pg.accountsTab.Changed() {
		pg.selectedAccount = pg.selectedWallet.Accounts[pg.accountsTab.Selected]
		pg.wallet.CreateTransaction(pg.selectedWallet.ID, pg.selectedAccount.Number, pg.txAuthorErrChan)
	}

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

	for pg.nextButtonWidget.Clicked(common.gtx) {
		if pg.validate(false) && pg.calculateErrorText == "" {
			pg.isConfirmationModalOpen = true
		}
	}

	for pg.confirmButtonWidget.Clicked(common.gtx) {
		pg.sendErrorText = ""
		pg.isPasswordModalOpen = true
	}

	for pg.closeConfirmationModalButtonWidget.Clicked(common.gtx) {
		pg.sendErrorText = ""
		pg.isConfirmationModalOpen = false
	}

	for pg.currencySwapWidget.Clicked(common.gtx) {
		if pg.selectedExchange == "DCR" {
			pg.selectedExchange = "USD"
			pg.unselectedExchange = "DCR"
		} else {
			pg.selectedExchange = "DCR"
			pg.unselectedExchange = "USD"
		}
		pg.calculateValues()
	}

	for range pg.destinationAddressEditor.Events(common.gtx) {
		go pg.calculateValues()
	}

	if pg.destinationAddressEditor.Len() == 0 || pg.sendAmountEditor.Len() == 0 {
		pg.balanceAfterSend(pg.selectedAccount.SpendableBalance)
	}

	for range pg.sendAmountEditor.Events(common.gtx) {
		go pg.calculateValues()
	}

	for pg.copyIconWidget.Clicked(common.gtx) {
		clipboard.WriteAll(pg.txHash)
		pg.hasCopiedTxHash = true
	}
	// switch {
	// case pg.selectedExchange == "DCR":
	// 	pg.activeAmount = pg.inactiveAmount
	// 	pg.inactiveAmount = pg.inactiveAmount
	// case pg.selectedExchange == "USD":
	// 	pg.activeAmount = "0"
	// 	pg.inactiveAmount = "0"
	// }

	select {
	case err := <-pg.txAuthorErrChan:
		pg.calculateErrorText = err.Error()
	case err := <-pg.broadcastErrChan:
		pg.sendErrorText = err.Error()
		pg.isBroadcastingTransaction = false
	default:
	}
}

func (pg *SendPage) drawWalletsTab(common pageCommon, body func()) {
	wallets := make([]decredmaterial.TabItem, len(pg.wallets))
	for i := range pg.wallets {
		wallets[i] = decredmaterial.TabItem{
			Label: pg.theme.Body1(pg.wallets[i].Name),
		}
	}
	pg.walletsTab.SetTabs(wallets)

	pg.setAccountTabs()
	pg.walletsTab.Layout(common.gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(common.gtx,
			layout.Rigid(func() {
				layout.Inset{Top: unit.Dp(10), Right: unit.Dp(10)}.Layout(common.gtx, func() {
					pg.theme.H6("Accounts: ").Layout(common.gtx)
				})
			}),
			layout.Rigid(func() {
				pg.accountsTab.Layout(common.gtx, body)
			}),
		)
	})
}

func (pg *SendPage) setAccountTabs() {
	accounts := make([]decredmaterial.TabItem, len(pg.selectedWallet.Accounts))
	for i := range pg.selectedWallet.Accounts {
		if pg.selectedWallet.Accounts[i].Name == "imported" {
			continue
		}
		accounts[i] = decredmaterial.TabItem{
			Label: pg.theme.Body1(pg.selectedWallet.Accounts[i].Name),
		}
	}
	pg.accountsTab.SetTabs(accounts)
}

func (pg *SendPage) Layout(common pageCommon) {
	go pg.getUSDValues(&pg.data)

	if len(common.info.Wallets) == 0 {
		// show no wallets message
		return
	}

	pg.wallets = common.info.Wallets
	if !pg.hasInitializedTxAuthor {
		pg.selectedWallet = pg.wallets[*common.selectedWallet]
		pg.selectedAccount = pg.selectedWallet.Accounts[0]

		pg.wallet.CreateTransaction(pg.selectedWallet.ID, pg.selectedAccount.Number, pg.txAuthorErrChan)
		pg.hasInitializedTxAuthor = true
	}

	common.Layout(common.gtx, func() {
		pg.drawPageContents(common)
	})

	if pg.isConfirmationModalOpen {
		pg.drawConfirmationModal(common.gtx)

		if pg.isPasswordModalOpen {
			pg.drawPasswordModal(common.gtx)
		}
	}
}

func (pg *SendPage) drawPageContents(common pageCommon) {
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
			pg.DCRToUSDLayout(common.gtx)
		},
		func() {
			pg.sendAmountEditorMaterial.Layout(common.gtx, pg.sendAmountEditor)
		},
		// func() {
		// 	layout.Flex{}.Layout(common.gtx,
		// 		layout.Rigid(func() {
		// 			pg.DCRToUSDLayout(common.gtx)
		// 		}),
		// 		layout.Rigid(func() {
		// 			if pg.sendAmountEditor.Focused() {
		// 				pg.sendAmountEditorMaterial.Hint = ""
		// 			} else {
		// 				pg.sendAmountEditorMaterial.Hint = "Amount to be sent"
		// 			}
		// 			pg.sendAmountEditorMaterial.Layout(common.gtx, pg.sendAmountEditor)
		// 		}),
		// 	)
		// },
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

	w := func() {
		inset := layout.Inset{
			Left: unit.Dp(-110),
		}
		inset.Layout(common.gtx, func() {
			layout.Flex{Axis: layout.Vertical}.Layout(common.gtx,
				layout.Rigid(func() {
					layout.UniformInset(unit.Dp(7)).Layout(common.gtx, func() {
						pg.pageContainer.Layout(common.gtx, len(pageContent), func(i int) {
							layout.Inset{Top: unit.Dp(5)}.Layout(common.gtx, pageContent[i])
						})
					})
				}),
			)
		})
	}
	pg.drawWalletsTab(common, w)
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
			pg.tableLayoutFunc(gtx, pg.theme.Body2("Transaction Fee"), pg.activeTransactionFeeValue, pg.inactiveTransactionFeeValue)
		},
		func() {
			pg.tableLayoutFunc(gtx, pg.theme.Body2("Total Cost"), pg.activeTotalCostValue, pg.inactiveTotalCostValue)
		},
		func() {
			pg.tableLayoutFunc(gtx, pg.theme.Body2("Balance after send"), "", "")
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

func (pg *SendPage) tableLayoutFunc(gtx *layout.Context, leftLabel decredmaterial.Label, active, inactive string) {
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
						b.Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (pg *SendPage) DCRToUSDLayout(gtx *layout.Context) {
	layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func() {
			layout.W.Layout(gtx, func() {
				layout.Flex{}.Layout(gtx,
					layout.Rigid(func() {
						pg.theme.Body1(pg.activeAmount).Layout(gtx)

					}),
					layout.Rigid(func() {
						layout.Inset{Left: unit.Dp(3)}.Layout(gtx, func() {
							pg.currencySwap.Layout(gtx, &pg.currencySwapWidget)
						})
					}),
					layout.Rigid(func() {
						pg.theme.Body2(pg.inactiveAmount).Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (pg *SendPage) drawConfirmationModal(gtx *layout.Context) {
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

func (pg *SendPage) drawPasswordModal(gtx *layout.Context) {
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
	// pg.activeTransactionFeeValueLabel.Text = "- DCR"
	// pg.activeTotalCostValueLabel.Text = "- DCR"
	// pg.inactiveTransactionFeeValueLabel.Text = "(- USD)"
	// pg.inactiveTotalCostValueLabel.Text = "(- USD)"

	pg.activeTransactionFeeValue = fmt.Sprintf("- %s", pg.selectedExchange)
	pg.activeTotalCostValue = fmt.Sprintf("- %s", pg.selectedExchange)
	pg.inactiveTransactionFeeValue = fmt.Sprintf("(- %s)", pg.unselectedExchange)
	pg.inactiveTotalCostValue = fmt.Sprintf("(- %s)", pg.unselectedExchange)

	pg.calculateErrorText = ""
	pg.activeAmount = fmt.Sprintf("- %s", pg.selectedExchange)
	pg.inactiveAmount = fmt.Sprintf("- %s", pg.unselectedExchange)

	if pg.txAuthor == nil || !pg.validate(true) {
		return
	}

	amountDCR, _ := strconv.ParseFloat(pg.sendAmountEditor.Text(), 64)
	amount, err := dcrutil.NewAmount(amountDCR)
	if err != nil {
		pg.calculateErrorText = fmt.Sprintf("error estimating transaction fee: %s", err)
		return
	}
	amountAtoms := int64(amount)

	// set destination address
	pg.txAuthor.RemoveSendDestination(0)
	pg.txAuthor.AddSendDestination(pg.destinationAddressEditor.Text(), amountAtoms, false)

	// calculate transaction fee
	feeAndSize, err := pg.txAuthor.EstimateFeeAndSize()
	if err != nil {
		pg.calculateErrorText = fmt.Sprintf("error estimating transaction fee: %s", err)
		return
	}

	txFee := feeAndSize.Fee.AtomValue
	totalCost := txFee + amountAtoms

	// pg.calculatedUSDValues(amountAtoms, txFee)

	pg.remainingBalance = pg.selectedWallet.SpendableBalance - totalCost
	pg.transactionFeeValueLabel.Text = dcrutil.Amount(txFee).String()
	pg.totalCostValueLabel.Text = dcrutil.Amount(totalCost).String()
	pg.balanceAfterSend(pg.remainingBalance)
}

func (pg *SendPage) balanceAfterSend(balance int64) {
	pg.balanceAfterSendValueLabel.Text = dcrutil.Amount(balance).String()
}

func (pg *SendPage) calculatedUSDValues(amountAtoms, txFee int64) {
 pg.getUSDValues(&pg.data)
	amount, err := strconv.ParseFloat(pg.data.LastTradeRate, 64)
	if err != nil {
		panic(err)
	}

	totalAmountUSD := amountDCR * amountUSD
	txFeeValueUSD := dcrutil.Amount(txFee).ToCoin() * amountDCR

	switch {
	case pg.selectedExchange == "DCR":
		pg.activeAmount = dcrutil.Amount(amountAtoms).String()
		pg.inactiveAmount = fmt.Sprintf("%s USD", strconv.FormatFloat(totalAmountUSD, 'f', 5, 64))
		pg.activeTransactionFeeValue = dcrutil.Amount(txFee).String()
		pg.inactiveTransactionFeeValue = fmt.Sprintf("(%f USD)", txFeeValueUSD)
		pg.activeTotalCostValue = dcrutil.Amount(totalCost).String()
		pg.inactiveTotalCostValue = fmt.Sprintf("(%s USD)", strconv.FormatFloat(totalAmountUSD+txFeeValueUSD, 'f', 5, 64))

	case pg.selectedExchange == "USD":
		pg.activeAmount = fmt.Sprintf("%s USD", strconv.FormatFloat(totalAmountUSD, 'f', 5, 64))
		pg.inactiveAmount = dcrutil.Amount(amountAtoms).String()
		pg.activeTransactionFeeValue = fmt.Sprintf("%f USD", txFeeValueUSD)
		pg.inactiveTransactionFeeValue = fmt.Sprintf("(%s)", dcrutil.Amount(txFee).String())
		pg.activeTotalCostValue = fmt.Sprintf("%s USD", strconv.FormatFloat(totalAmountUSD+txFeeValueUSD, 'f', 5, 64))
		pg.inactiveTotalCostValue = fmt.Sprintf("(%s )", dcrutil.Amount(totalCost).String())
	}

	pg.remainingBalance = pg.selectedWallet.SpendableBalance - totalCost
	pg.balanceAfterSendValueLabel.Text = dcrutil.Amount(pg.remainingBalance).String()
}

// func (pg *SendPage) calculatedUSDValues(amountAtoms, txFee int64) {
// 	amount, err := strconv.ParseFloat("16.0000", 64)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// usd, err := dcrutil.NewAmount(amount)
// 	// if err != nil {
// 	// 	pg.calculateErrorText = fmt.Sprintf("error estimating transaction fee: %s", err)
// 	// 	return
// 	// }
// 	// usdex := int64(usd)

// 	// if pg.selectedExchange == "DCR" {
// 	// 	pg.activeAmount = dcrutil.Amount(amountAtoms).String()
// 	// 	pg.inactiveAmount = fmt.Sprintf("%s USD", strconv.Itoa(int(usdex*amountAtoms)))
// 	// } else {
// 	// 	pg.activeAmount = fmt.Sprintf("%s USD", strconv.Itoa(int(usdex*amountAtoms)))
// 	// 	pg.inactiveAmount = strconv.Itoa(int(amountAtoms))
// 	// }

// 	// usdRate := usdex * txFee
// 	// pg.inactiveTransactionFeeValueLabel.Text = fmt.Sprintf("(%s USD)", strconv.Itoa(int(usdRate)))
// 	// pg.inactiveTotalCostValueLabel.Text = fmt.Sprintf("(%s USD)", strconv.Itoa(int(usdRate+usdex)))
// }

func (pg *SendPage) getUSDValues(target interface{}) {
	url := "https://api.bittrex.com/v3/markets/DCR-USDT/ticker"
	resp, err := http.Get(url)
	if err != nil {
		// fmt.Println(err.Error())
		return
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(target)
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

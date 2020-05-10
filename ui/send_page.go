package ui

import (
	"fmt"
	"strconv"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
)

const PageSend = "send"

var sendPage *Send

// Send represents the send page of the app.
// It should only be accessible if the app finds
// at least one wallet.
type Send struct {
	theme    *decredmaterial.Theme
	txAuthor *dcrlibwallet.TxAuthor
	wallet   *wallet.Wallet

	selectedWallet  wallet.InfoShort
	selectedAccount wallet.Account

	isConfirmationModalOpen   bool
	isPasswordModalOpen       bool
	isBroadcastingTransaction bool

	// editors
	destinationAddressEditor *widget.Editor
	sendAmountEditor         *widget.Editor

	destinationAddressEditorMaterial decredmaterial.Editor
	sendAmountEditorMaterial         decredmaterial.Editor

	// lines
	editorLine *decredmaterial.Line

	// labels
	titleLabel                     decredmaterial.Label
	selectedAccountLabel           decredmaterial.Label
	selectedAccountBalanceLabel    decredmaterial.Label
	selectedWalletLabel            decredmaterial.Label
	sendAmountLabel                decredmaterial.Label
	sourceAccountLabel             decredmaterial.Label
	destinationAddressLabel        decredmaterial.Label
	transactionFeeLabel            decredmaterial.Label
	transactionFeeValueLabel       decredmaterial.Label
	totalCostLabel                 decredmaterial.Label
	totalCostValueLabel            decredmaterial.Label
	balanceAfterSendLabel          decredmaterial.Label
	balanceAfterSendValueLabel     decredmaterial.Label
	confirmModalTitleLabel         decredmaterial.Label
	confirmSourceAccountLabel      decredmaterial.Label
	confirmDestinationAddressLabel decredmaterial.Label
	confirmWarningLabel            decredmaterial.Label
	accountsModalTitleLabel        decredmaterial.Label
	txHashLabel                    decredmaterial.Label

	// error labels
	destinationAddressErrorLabel decredmaterial.Label
	sendAmountErrorLabel         decredmaterial.Label
	calculateErrorLabel          decredmaterial.Label
	sendErrorLabel               decredmaterial.Label

	// buttons
	nextButtonWidget                   *widget.Button
	confirmButtonWidget                *widget.Button
	closeConfirmationModalButtonWidget *widget.Button
	pasteAddressButtonWidget           *widget.Button

	nextButtonMaterial                   decredmaterial.Button
	confirmButtonMaterial                decredmaterial.Button
	closeConfirmationModalButtonMaterial decredmaterial.Button
	pasteAddressButtonMaterial           decredmaterial.Button

	passwordModal *decredmaterial.Password
}

// NewSendPage initializes this page's widgets
func (win *Window) NewSendPage() *Send {
	pg := &Send{}
	pg.theme = win.theme
	pg.wallet = win.wallet
	pg.isConfirmationModalOpen = false
	pg.isBroadcastingTransaction = false
	pg.isPasswordModalOpen = false
	pg.txAuthor = win.txAuthor

	// labels
	pg.titleLabel = pg.theme.H5("Send DCR")
	pg.selectedAccountLabel = pg.theme.Body2("")
	pg.selectedAccountBalanceLabel = pg.theme.Body2("")
	pg.selectedWalletLabel = pg.theme.Body2("")
	pg.sourceAccountLabel = pg.theme.Body2("Source Account:")
	pg.destinationAddressLabel = pg.theme.Body2("Destination Address")
	pg.transactionFeeLabel = pg.theme.Body2("Transaction Fee:")
	pg.transactionFeeValueLabel = pg.theme.Body2("0 DCR")
	pg.sendAmountLabel = pg.theme.Body2("Amount")
	pg.totalCostLabel = pg.theme.Body2("Total Cost")
	pg.totalCostValueLabel = pg.theme.Body2("0 DCR")
	pg.balanceAfterSendLabel = pg.theme.Body2("Balance after send")
	pg.balanceAfterSendValueLabel = pg.theme.Body2("0 DCR")
	pg.confirmModalTitleLabel = pg.theme.H5("Confirm to send")
	pg.confirmSourceAccountLabel = pg.theme.Body2("")
	pg.confirmDestinationAddressLabel = pg.theme.Body2("To destination address")
	pg.confirmWarningLabel = pg.theme.Caption("Your DCR will be sent and CANNOT be undone")
	pg.accountsModalTitleLabel = pg.theme.H5("Choose sending account")
	pg.txHashLabel = pg.theme.Body2("")
	pg.txHashLabel.Color = pg.theme.Color.Success

	pg.destinationAddressErrorLabel = pg.theme.Caption("")
	pg.sendErrorLabel = pg.theme.Caption("")
	pg.sendAmountErrorLabel = pg.theme.Caption("")
	pg.calculateErrorLabel = pg.theme.Caption("")
	pg.destinationAddressErrorLabel.Color = pg.theme.Color.Danger
	pg.sendErrorLabel.Color = pg.theme.Color.Danger
	pg.sendAmountErrorLabel.Color = pg.theme.Color.Danger
	pg.calculateErrorLabel.Color = pg.theme.Color.Danger

	pg.passwordModal = pg.theme.Password()

	pg.destinationAddressEditor = new(widget.Editor)
	pg.sendAmountEditor = new(widget.Editor)

	pg.destinationAddressEditorMaterial = pg.theme.Editor("Destination Address")
	pg.sendAmountEditorMaterial = pg.theme.Editor("Amount to be sent")

	pg.editorLine = pg.theme.Line()

	pg.nextButtonWidget = new(widget.Button)
	pg.confirmButtonWidget = new(widget.Button)
	pg.closeConfirmationModalButtonWidget = new(widget.Button)
	pg.pasteAddressButtonWidget = new(widget.Button)

	pg.nextButtonMaterial = pg.theme.Button("Next")
	pg.confirmButtonMaterial = pg.theme.Button("Confirm")
	pg.closeConfirmationModalButtonMaterial = pg.theme.Button("Close")
	pg.pasteAddressButtonMaterial = pg.theme.Button("Paste")

	pg.closeConfirmationModalButtonMaterial.Background = pg.theme.Color.Surface
	pg.closeConfirmationModalButtonMaterial.Color = pg.theme.Color.Primary

	return pg
}

// Draw renders all of this page's widgets
func (pg *Send) Draw(gtx *layout.Context) interface{} {
	pg.handleEvents(gtx)

	widgetsWrapper := [][]func(){
		{
			func() {
				pg.titleLabel.Layout(gtx)

				if pg.txHashLabel.Text != "" {
					inset := layout.Inset{
						Top: unit.Dp(40),
					}
					inset.Layout(gtx, func() {
						pg.txHashLabel.Layout(gtx)
					})
				}
			},
			func() {
				pg.drawSelectedAccountSection(gtx)
			},
			func() {
				pg.drawDestinationAddressSection(gtx)
			},
			func() {
				pg.drawSendAmountSection(gtx)
			},
		},
		pg.getTransactionDetailWidgets(gtx),
		{
			func() {
				pg.calculateErrorLabel.Layout(gtx)
			},
			func() {
				pg.nextButtonMaterial.Layout(gtx, pg.nextButtonWidget)
			},
		},
	}

	widgets := collapseWidgetsWrapper(widgetsWrapper)

	list := layout.List{Axis: layout.Vertical}
	list.Layout(gtx, len(widgets), func(i int) {
		layout.UniformInset(unit.Dp(7)).Layout(gtx, widgets[i])
	})

	if pg.isConfirmationModalOpen {
		pg.drawConfirmationModal(gtx)

		if pg.isPasswordModalOpen {
			pg.drawPasswordModal(gtx)
		}
	}
	return nil
}

func (pg *Send) handleEvents(gtx *layout.Context) {
	for pg.pasteAddressButtonWidget.Clicked(gtx) {
		pg.destinationAddressEditor.Insert(GetClipboardContent())
	}

	for pg.nextButtonWidget.Clicked(gtx) {
		if pg.validate(false) && pg.calculateErrorLabel.Text == "" {
			pg.isConfirmationModalOpen = true
		}
	}

	for pg.confirmButtonWidget.Clicked(gtx) {
		pg.isPasswordModalOpen = true
	}

	for pg.closeConfirmationModalButtonWidget.Clicked(gtx) {
		pg.isConfirmationModalOpen = false
	}
}

func (pg *Send) drawDestinationAddressSection(gtx *layout.Context) {
	pg.destinationAddressLabel.Layout(gtx)
	inset := layout.Inset{
		Top: unit.Dp(20),
	}
	inset.Layout(gtx, func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(0.95, func() {
				pg.destinationAddressEditorMaterial.Layout(gtx, pg.destinationAddressEditor)
				inset := layout.Inset{
					Top: unit.Dp(25),
				}
				inset.Layout(gtx, func() {
					pg.editorLine.Draw(gtx)
				})
			}),
			layout.Rigid(func() {
				pg.pasteAddressButtonMaterial.Layout(gtx, pg.pasteAddressButtonWidget)
			}),
		)
	})

	inset = layout.Inset{
		Top: unit.Dp(40),
	}
	inset.Layout(gtx, func() {
		pg.destinationAddressErrorLabel.Layout(gtx)
	})
}

func (pg *Send) drawSendAmountSection(gtx *layout.Context) {
	pg.sendAmountLabel.Layout(gtx)
	inset := layout.Inset{
		Top: unit.Dp(20),
	}
	inset.Layout(gtx, func() {
		pg.sendAmountEditorMaterial.Layout(gtx, pg.sendAmountEditor)

		inset := layout.Inset{
			Top: unit.Dp(25),
		}
		inset.Layout(gtx, func() {
			pg.editorLine.Draw(gtx)
		})
	})

	if pg.sendAmountErrorLabel.Text != "" {
		inset := layout.Inset{
			Top: unit.Dp(40),
		}
		inset.Layout(gtx, func() {
			pg.sendAmountErrorLabel.Layout(gtx)
		})
	}
}

func (pg *Send) drawSelectedAccountSection(gtx *layout.Context) {
	layout.Flex{}.Layout(gtx,
		layout.Flexed(0.22, func() {
		}),
		layout.Flexed(1, func() {
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
		}),
	)
}

func (pg *Send) drawConfirmationModal(gtx *layout.Context) {
	widgetsWrapper := [][]func(){
		{
			func() {
				pg.theme.H4("Confirm to send").Layout(gtx)
			},
			func() {
				if pg.sendErrorLabel.Text != "" {
					pg.sendErrorLabel.Layout(gtx)
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
						pg.confirmDestinationAddressLabel.Layout(gtx)
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
		},
		pg.getTransactionDetailWidgets(gtx),
		{
			func() {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				layout.Stack{Alignment: layout.Center}.Layout(gtx,
					layout.Expanded(func() {
						pg.confirmWarningLabel.Layout(gtx)
					}),
				)
			},
			func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func() {
						pg.confirmButtonMaterial.Layout(gtx, pg.confirmButtonWidget)
					}),
					layout.Rigid(func() {
						pg.closeConfirmationModalButtonMaterial.Layout(gtx, pg.closeConfirmationModalButtonWidget)
					}),
				)
			},
		},
	}

	widgets := collapseWidgetsWrapper(widgetsWrapper)
	pg.theme.Modal(gtx, func() {
		(&layout.List{Axis: layout.Vertical}).Layout(gtx, len(widgets), func(i int) {
			layout.UniformInset(unit.Dp(10)).Layout(gtx, widgets[i])
		})
	})
}

func (pg *Send) drawPasswordModal(gtx *layout.Context) {
	pg.theme.Modal(gtx, func() {
		pg.passwordModal.Layout(gtx, func(password []byte) {
			pg.isBroadcastingTransaction = true
			pg.isPasswordModalOpen = false

			pg.wallet.BroadcastTransaction(pg.txAuthor, password)
		}, func() {
			pg.isPasswordModalOpen = false
		})
	})
}

func (pg *Send) getTransactionDetailWidgets(gtx *layout.Context) []func() {
	return []func(){
		pg.tableLayoutFunc(gtx, pg.transactionFeeLabel, pg.transactionFeeValueLabel),
		pg.tableLayoutFunc(gtx, pg.totalCostLabel, pg.totalCostValueLabel),
		pg.tableLayoutFunc(gtx, pg.balanceAfterSendLabel, pg.balanceAfterSendValueLabel),
	}
}

func (pg *Send) tableLayoutFunc(gtx *layout.Context, leftLabel, rightLabel decredmaterial.Label) func() {
	return func() {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func() {
				leftLabel.Layout(gtx)
			}),
			layout.Flexed(1, func() {
				layout.Stack{Alignment: layout.NE}.Layout(gtx,
					layout.Stacked(func() {
						rightLabel.Layout(gtx)
					}),
				)
			}),
		)
	}
}

// collapseWidggetsWrapper receives a two dimensional slice and collapses it into a single slice
func collapseWidgetsWrapper(widgetsWrapper [][]func()) []func() {
	widgets := []func(){}
	for i := range widgetsWrapper {
		widgets = append(widgets, widgetsWrapper[i]...)
	}
	return widgets
}

func (pg *Send) validateDestinationAddress(ignoreEmpty bool) bool {
	pg.destinationAddressErrorLabel.Text = ""

	destinationAddress := pg.destinationAddressEditor.Text()
	if destinationAddress == "" && !ignoreEmpty {
		pg.destinationAddressErrorLabel.Text = "please enter a destination address"
		return false
	}

	if destinationAddress != "" {
		isValid, _ := pg.wallet.IsAddressValid(destinationAddress)
		if !isValid {
			pg.destinationAddressErrorLabel.Text = "invalid address"
			return false
		}
	}
	return true
}

func (pg *Send) validateAmount(ignoreEmpty bool) bool {
	pg.sendAmountErrorLabel.Text = ""

	amount := pg.sendAmountEditor.Text()
	if amount == "" {
		if !ignoreEmpty {
			pg.sendAmountErrorLabel.Text = "please enter a send amount"
		}
		return false
	}

	if amount != "" {
		_, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			pg.sendAmountErrorLabel.Text = "please enter a valid amount"
			return false
		}
	}

	return true
}

func (pg *Send) validate(ignoreEmpty bool) bool {
	isAddressValid := pg.validateDestinationAddress(ignoreEmpty)
	isAmountValid := pg.validateAmount(ignoreEmpty)

	if !isAddressValid || !isAmountValid || pg.calculateErrorLabel.Text != "" {
		pg.nextButtonMaterial.Background = pg.theme.Color.Hint
		return false
	}

	pg.nextButtonMaterial.Background = pg.theme.Color.Primary
	return true
}

func (pg *Send) watchForBroadcastResult(broadcastResult *wallet.Broadcast) {
	if broadcastResult == nil {
		return
	}

	txHash := broadcastResult.TxHash
	err := broadcastResult.Err

	if err != nil {
		pg.sendErrorLabel.Text = fmt.Sprintf("error broadcasting transaction: %s", err.Error())
	} else if txHash != "" {
		pg.txHashLabel.Text = fmt.Sprintf("The transaction was published successfully. Hash: %s", txHash)
		pg.destinationAddressEditor.SetText("")
		pg.sendAmountEditor.SetText("")
		pg.isConfirmationModalOpen = false
	}
	pg.isBroadcastingTransaction = false
}

func (pg *Send) watchAndUpdateValues(gtx *layout.Context, broadcastResult *wallet.Broadcast) {
	pg.watchForBroadcastResult(broadcastResult)

	pg.confirmButtonMaterial.Text = "Send"
	pg.confirmButtonMaterial.Background = pg.theme.Color.Primary

	if pg.isBroadcastingTransaction {
		pg.confirmButtonMaterial.Text = "Sending..."
		pg.confirmButtonMaterial.Background = pg.theme.Color.Background
	}

	for range pg.destinationAddressEditor.Events(gtx) {
		go pg.calculateValues()
	}

	for range pg.sendAmountEditor.Events(gtx) {
		go pg.calculateValues()
	}
}

func (pg *Send) calculateValues() {
	pg.transactionFeeValueLabel.Text = "0 DCR"
	pg.totalCostValueLabel.Text = "0 DCR "
	pg.calculateErrorLabel.Text = ""

	pg.balanceAfterSendValueLabel.Text = dcrutil.Amount(pg.selectedWallet.SpendableBalance).String()

	if pg.txAuthor == nil || !pg.validate(true) {
		return
	}

	amountDCR, _ := strconv.ParseFloat(pg.sendAmountEditor.Text(), 64)
	if amountDCR <= 0 {
		return
	}

	amount, err := dcrutil.NewAmount(amountDCR)
	if err != nil {
		pg.calculateErrorLabel.Text = fmt.Sprintf("error estimating transaction fee: %s", err)
		return
	}

	amountAtoms := int64(amount)

	// set destination address
	pg.txAuthor.RemoveSendDestination(0)
	pg.txAuthor.AddSendDestination(pg.destinationAddressEditor.Text(), amountAtoms, false)

	// calculate transaction fee
	feeAndSize, err := pg.txAuthor.EstimateFeeAndSize()
	if err != nil {
		pg.calculateErrorLabel.Text = fmt.Sprintf("error estimating transaction fee: %s", err)
		return
	}

	txFee := feeAndSize.Fee.AtomValue
	totalCost := txFee + amountAtoms
	remainingBalance := pg.selectedWallet.SpendableBalance - totalCost

	pg.transactionFeeValueLabel.Text = dcrutil.Amount(txFee).String()
	pg.totalCostValueLabel.Text = dcrutil.Amount(totalCost).String()
	pg.balanceAfterSendValueLabel.Text = dcrutil.Amount(remainingBalance).String()
}

func (win *Window) SendPage() {
	if win.walletInfo.LoadedWallets == 0 {
		win.Page(func() {
			win.outputs.noWallet.Layout(win.gtx)
		})
		return
	}

	selectedWallet := win.walletInfo.Wallets[win.selected]
	selectedAccount := win.walletInfo.Wallets[win.selected].Accounts[win.selectedAccount]

	if sendPage == nil {
		sendPage = win.NewSendPage()
		sendPage.wallet = win.wallet

		sendPage.selectedWallet = selectedWallet
		sendPage.selectedAccount = selectedAccount

		sendPage.wallet.CreateTransaction(selectedWallet.ID, selectedAccount.Number)
	}

	if sendPage.selectedWallet.ID != selectedWallet.ID || sendPage.selectedAccount.Number != selectedAccount.Number {
		sendPage.selectedWallet = selectedWallet
		sendPage.selectedAccount = selectedAccount

		sendPage.wallet.CreateTransaction(selectedWallet.ID, selectedAccount.Number)
	}

	if win.txAuthor != nil {
		sendPage.txAuthor = win.txAuthor
		win.txAuthor = nil
	}

	sendPage.validate(true)
	sendPage.watchAndUpdateValues(win.gtx, win.broadcastResult)
	if win.broadcastResult != nil {
		win.broadcastResult = nil
	}

	body := func() {
		layout.Stack{}.Layout(win.gtx,
			layout.Expanded(func() {
				layout.Flex{}.Layout(win.gtx,
					layout.Rigid(func() {
						win.combined.sel.Layout(win.gtx, func() {

						})
					}),
					layout.Rigid(func() {
						sendPage.Draw(win.gtx)
					}),
				)
			}),
		)
	}
	win.TabbedPage(body)
}

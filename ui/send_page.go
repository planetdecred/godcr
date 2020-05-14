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

type SendPage struct {
	theme           *decredmaterial.Theme
	wallet          *wallet.Wallet
	txAuthor        *dcrlibwallet.TxAuthor
	broadcastResult *wallet.Broadcast

	selectedWallet  wallet.InfoShort
	selectedAccount wallet.Account

	destinationAddressEditor           *widget.Editor
	sendAmountEditor                   *widget.Editor
	pasteAddressButtonWidget           *widget.Button
	nextButtonWidget                   *widget.Button
	closeConfirmationModalButtonWidget *widget.Button
	confirmButtonWidget                *widget.Button

	copyIconMaterial decredmaterial.IconButton
	copyIconWidget   *widget.Button

	destinationAddressErrorLabel decredmaterial.Label
	sendAmountErrorLabel         decredmaterial.Label
	transactionFeeValueLabel     decredmaterial.Label
	totalCostValueLabel          decredmaterial.Label
	balanceAfterSendValueLabel   decredmaterial.Label
	calculateErrorLabel          decredmaterial.Label

	destinationAddressEditorMaterial     decredmaterial.Editor
	sendAmountEditorMaterial             decredmaterial.Editor
	pasteAddressButtonMaterial           decredmaterial.IconButton
	nextButtonMaterial                   decredmaterial.Button
	closeConfirmationModalButtonMaterial decredmaterial.Button
	confirmButtonMaterial                decredmaterial.Button

	sendErrorText string
	txHashText    string
	txHash        string

	editorLine    *decredmaterial.Line
	passwordModal *decredmaterial.Password

	isConfirmationModalOpen   bool
	isPasswordModalOpen       bool
	isBroadcastingTransaction bool
	hasSetWallet              bool
	hasCopiedTxHash           bool
}

const PageSend = "send"

func (win *Window) SendPage(common pageCommon) layout.Widget {
	page := &SendPage{
		theme:           common.theme,
		wallet:          common.wallet,
		txAuthor:        &win.txAuthor,
		broadcastResult: &win.broadcastResult,

		destinationAddressEditor:           new(widget.Editor),
		sendAmountEditor:                   new(widget.Editor),
		pasteAddressButtonWidget:           new(widget.Button),
		nextButtonWidget:                   new(widget.Button),
		closeConfirmationModalButtonWidget: new(widget.Button),
		confirmButtonWidget:                new(widget.Button),
		copyIconWidget:                     new(widget.Button),

		destinationAddressErrorLabel: common.theme.ErrorLabel(""),
		sendAmountErrorLabel:         common.theme.ErrorLabel(""),
		sendErrorText:                "",
		txHashText:                   "",

		destinationAddressEditorMaterial:     common.theme.Editor("Destination Address"),
		closeConfirmationModalButtonMaterial: common.theme.Button("Close"),
		sendAmountEditorMaterial:             common.theme.Editor("Amount to be sent"),
		editorLine:                           common.theme.Line(),
		nextButtonMaterial:                   common.theme.Button("Next"),
		confirmButtonMaterial:                common.theme.Button("Confirm"),
		transactionFeeValueLabel:             common.theme.Body2("0 DCR"),
		totalCostValueLabel:                  common.theme.Body2("0 DCR"),
		balanceAfterSendValueLabel:           common.theme.Body2("0 DCR"),
		copyIconMaterial:                     common.theme.IconButton(mustIcon(decredmaterial.NewIcon(icons.ContentContentCopy))),

		isConfirmationModalOpen: false,
		isPasswordModalOpen:     false,
		hasSetWallet:            false,
		hasCopiedTxHash:         false,

		passwordModal: common.theme.Password(),

		pasteAddressButtonMaterial: decredmaterial.IconButton{
			Icon:       mustIcon(decredmaterial.NewIcon(icons.ContentContentPaste)),
			Size:       unit.Dp(30),
			Background: color.RGBA{},
			Color:      common.theme.Color.Text,
			Padding:    unit.Dp(5),
		},
	}
	page.closeConfirmationModalButtonMaterial.Background = common.theme.Color.Gray
	page.copyIconMaterial.Background = common.theme.Color.Background
	page.copyIconMaterial.Size = unit.Dp(30)
	page.copyIconMaterial.Color = common.theme.Color.Text
	page.copyIconMaterial.Padding = unit.Dp(5)

	return func() {
		page.Layout(common)
		page.Handle(common)
	}
}

func (pg *SendPage) Handle(common pageCommon) {
	pg.validate(true)

	pg.watchForBroadcastResult()
	pg.confirmButtonMaterial.Text = "Send"
	pg.confirmButtonMaterial.Background = pg.theme.Color.Primary

	if pg.hasCopiedTxHash {
		time.AfterFunc(3*time.Second, func() {
			pg.hasCopiedTxHash = false
		})
	}

	if pg.isBroadcastingTransaction {
		pg.confirmButtonMaterial.Text = "Sending..."
		pg.confirmButtonMaterial.Background = pg.theme.Color.Background
	}

	for pg.copyIconWidget.Clicked(common.gtx) {
		clipboard.WriteAll(pg.txHash)
	}

	for pg.pasteAddressButtonWidget.Clicked(common.gtx) {
		pg.destinationAddressEditor.Insert(GetClipboardContent())
	}

	for pg.nextButtonWidget.Clicked(common.gtx) {
		if pg.validate(false) && pg.calculateErrorLabel.Text == "" {
			pg.isConfirmationModalOpen = true
		}
	}

	for pg.confirmButtonWidget.Clicked(common.gtx) {
		pg.isPasswordModalOpen = true
	}

	for pg.closeConfirmationModalButtonWidget.Clicked(common.gtx) {
		pg.sendErrorText = ""
		pg.isConfirmationModalOpen = false
	}

	for range pg.destinationAddressEditor.Events(common.gtx) {
		go pg.calculateValues()
	}

	for range pg.sendAmountEditor.Events(common.gtx) {
		go pg.calculateValues()
	}
}

func (pg *SendPage) Layout(common pageCommon) {
	if len(common.info.Wallets) == 0 {
		// show no wallets message
		return
	}

	selectedWallet := common.info.Wallets[*common.selectedWallet]
	selectedAccount := selectedWallet.Accounts[0]

	if !pg.hasSetWallet || (pg.selectedWallet.ID != selectedWallet.ID || pg.selectedAccount.Number != selectedAccount.Number) {
		pg.selectedWallet = selectedWallet
		pg.selectedAccount = selectedAccount

		pg.wallet.CreateTransaction(pg.selectedWallet.ID, pg.selectedAccount.Number)
		pg.hasSetWallet = true
	}

	w := []func(){
		func() {
			pg.theme.H5("Send DCR").Layout(common.gtx)
		},
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
			pg.drawDestinationAddressSection(common.gtx)
		},
		func() {
			pg.drawSendAmountSection(common.gtx)
		},
		func() {
			pg.drawTransactionDetailWidgets(common.gtx)
		},
		func() {

		},
		func() {
			pg.nextButtonMaterial.Layout(common.gtx, pg.nextButtonWidget)
		},
	}

	common.LayoutWithWallets(common.gtx, func() {
		list := layout.List{Axis: layout.Vertical}
		list.Layout(common.gtx, len(w), func(i int) {
			layout.UniformInset(unit.Dp(7)).Layout(common.gtx, w[i])
		})
	})

	if pg.isConfirmationModalOpen {
		pg.drawConfirmationModal(common.gtx)

		if pg.isPasswordModalOpen {
			pg.drawPasswordModal(common.gtx)
		}
	}
}

func (pg *SendPage) drawSuccessSection(gtx *layout.Context) {
	if pg.txHashText != "" {
		layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func() {
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
		layout.Flex{}.Layout(gtx,
			layout.Flexed(0.35, func() {
			}),
			layout.Flexed(1, func() {
				pg.theme.Caption("copied").Layout(gtx)
			}),
		)
	}
}

func (pg *SendPage) drawSelectedAccountSection(gtx *layout.Context) {
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

func (pg *SendPage) drawDestinationAddressSection(gtx *layout.Context) {
	pg.theme.Body2("Destination Address").Layout(gtx)

	inset := layout.Inset{
		Top: unit.Dp(20),
	}

	inset.Layout(gtx, func() {
		pg.destinationAddressEditorMaterial.Layout(gtx, pg.destinationAddressEditor)
		inset := layout.Inset{
			Left: unit.Dp(float32(gtx.Constraints.Width.Max - 30)),
		}
		inset.Layout(gtx, func() {
			pg.pasteAddressButtonMaterial.Layout(gtx, pg.pasteAddressButtonWidget)
		})

		inset = layout.Inset{
			Top: unit.Dp(25),
		}
		inset.Layout(gtx, func() {
			pg.editorLine.Draw(gtx)
		})
	})

	inset = layout.Inset{
		Top: unit.Dp(48),
	}
	inset.Layout(gtx, func() {
		pg.destinationAddressErrorLabel.Layout(gtx)
	})
}

func (pg *SendPage) drawSendAmountSection(gtx *layout.Context) {
	pg.theme.Body2("Amount").Layout(gtx)
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
			Top: unit.Dp(48),
		}
		inset.Layout(gtx, func() {
			pg.sendAmountErrorLabel.Layout(gtx)
		})
	}
}

func (pg *SendPage) drawTransactionDetailWidgets(gtx *layout.Context) {
	w := []func(){
		func() {
			pg.tableLayoutFunc(gtx, pg.theme.Body2("Transaction Fee"), pg.transactionFeeValueLabel)
		},
		func() {
			pg.tableLayoutFunc(gtx, pg.theme.Body2("Total Cost"), pg.totalCostValueLabel)
		},
		func() {
			pg.tableLayoutFunc(gtx, pg.theme.Body2("Balance after send"), pg.balanceAfterSendValueLabel)
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

func (pg *SendPage) tableLayoutFunc(gtx *layout.Context, leftLabel, rightLabel decredmaterial.Label) {
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

func (pg *SendPage) drawConfirmationModal(gtx *layout.Context) {
	w := []func(){
		func() {
			pg.theme.H4("Confirm to send").Layout(gtx)
		},
		func() {
			if pg.sendErrorText != "" {
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
						Left: unit.Dp(15),
					}
					inset.Layout(gtx, func() {
						pg.closeConfirmationModalButtonMaterial.Layout(gtx, pg.closeConfirmationModalButtonWidget)
					})
				}),
			)
		},
	}

	pg.theme.Modal(gtx, func() {
		(&layout.List{Axis: layout.Vertical}).Layout(gtx, len(w), func(i int) {
			layout.UniformInset(unit.Dp(10)).Layout(gtx, w[i])
		})
	})
}

func (pg *SendPage) drawPasswordModal(gtx *layout.Context) {
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

func (pg *SendPage) validate(ignoreEmpty bool) bool {
	isAddressValid := pg.validateDestinationAddress(ignoreEmpty)
	isAmountValid := pg.validateAmount(ignoreEmpty)

	if !isAddressValid || !isAmountValid || pg.calculateErrorLabel.Text != "" {
		pg.nextButtonMaterial.Background = pg.theme.Color.Hint
		return false
	}

	pg.nextButtonMaterial.Background = pg.theme.Color.Primary
	return true
}

func (pg *SendPage) validateDestinationAddress(ignoreEmpty bool) bool {
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

func (pg *SendPage) validateAmount(ignoreEmpty bool) bool {
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

func (pg *SendPage) calculateValues() {
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

func (pg *SendPage) watchForBroadcastResult() {
	if pg.broadcastResult == nil {
		return
	}

	txHash := pg.broadcastResult.TxHash
	err := pg.broadcastResult.Err

	if err != nil {
		pg.sendErrorText = fmt.Sprintf("error broadcasting transaction: %s", err.Error())
	} else if txHash != "" {
		pg.txHash = txHash
		pg.txHashText = fmt.Sprintf("Successfull. Hash: %s", txHash)
		pg.destinationAddressEditor.SetText("")
		pg.sendAmountEditor.SetText("")
		pg.isConfirmationModalOpen = false
	}
	pg.isBroadcastingTransaction = false
}

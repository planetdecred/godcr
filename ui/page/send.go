package page

import (
	"fmt"
	"image"
	"strconv"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"

	"github.com/raedahgroup/godcr-gio/ui"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// SendID is the id of the send page
const SendID = "send"

type accountModalWidgets struct {
	titleLabel material.Label
	titleLine  *materialplus.Line
}

type editor struct {
	editor   *widget.Editor
	material material.Editor
	line     *materialplus.Line
}

type button struct {
	button   *widget.Button
	material material.Button
}

// Send represents the send page of the app.
// It should only be accessible if the app finds
// at least one wallet.
type Send struct {
	theme    *materialplus.Theme
	states   map[string]interface{}
	txAuthor *dcrlibwallet.TxAuthor
	wallet   *wallet.Wallet
	wallets  []wallet.InfoShort

	selectedWallet  *wallet.InfoShort
	selectedAccount *wallet.Account

	isAccountModalOpen        bool
	isConfirmationModalOpen   bool
	isPasswordModalOpen       bool
	isBroadcastingTransaction bool

	// editors
	destinationAddressEditor *editor
	sendAmountEditor         *editor

	// labels
	titleLabel                     material.Label
	selectedAccountLabel           material.Label
	selectedWalletLabel            material.Label
	sendAmountLabel                material.Label
	sourceAccountLabel             material.Label
	destinationAddressLabel        material.Label
	transactionFeeLabel            material.Label
	transactionFeeValueLabel       material.Label
	totalCostLabel                 material.Label
	totalCostValueLabel            material.Label
	balanceAfterSendLabel          material.Label
	balanceAfterSendValueLabel     material.Label
	confirmModalTitleLabel         material.Label
	confirmSourceAccountLabel      material.Label
	confirmDestinationAddressLabel material.Label
	confirmWarningLabel            material.Label
	accountsModalTitleLabel        material.Label
	txHashLabel                    material.Label

	// error labels
	destinationAddressErrorLabel material.Label
	sendAmountErrorLabel         material.Label
	calculateErrorLabel          material.Label
	sendErrorLabel               material.Label

	// buttons
	selectAccountButton *button
	nextButton          *button
	confirmButton       *button

	passwordModal *materialplus.Password

	accountSelectorButtons map[string]*widget.Button
}

// Init initializes this page's widgets
func (pg *Send) Init(theme *materialplus.Theme, wal *wallet.Wallet, states map[string]interface{}) {
	pg.theme = theme
	pg.wallet = wal
	pg.states = states
	pg.isAccountModalOpen = false
	pg.isConfirmationModalOpen = false
	pg.isBroadcastingTransaction = false
	pg.isPasswordModalOpen = false

	// labels
	pg.titleLabel = theme.H5("Send DCR")
	pg.selectedAccountLabel = pg.theme.Body2("")
	pg.selectedWalletLabel = pg.theme.Body2("")
	pg.sourceAccountLabel = pg.theme.Body2("Source Account:")
	pg.destinationAddressLabel = pg.theme.Body2("Destination Address")
	pg.transactionFeeLabel = pg.theme.Body2("Transaction Fee:")
	pg.transactionFeeValueLabel = pg.theme.Body2("0 DCR")
	pg.sendAmountLabel = pg.theme.Body2("Amount")
	pg.totalCostLabel = pg.theme.Body2("Total Cost")
	pg.totalCostValueLabel = pg.theme.Body2("O DCR")
	pg.balanceAfterSendLabel = pg.theme.Body2("Balance after send")
	pg.balanceAfterSendValueLabel = pg.theme.Body2("0 DCR")
	pg.confirmModalTitleLabel = pg.theme.H5("Confirm to send")
	pg.confirmSourceAccountLabel = pg.theme.Body2("")
	pg.confirmDestinationAddressLabel = pg.theme.Body2("To destination address")
	pg.confirmWarningLabel = pg.theme.Caption("Your DCR will be sent and CANNOT be undone")
	pg.accountsModalTitleLabel = pg.theme.H5("Choose sending account")
	pg.txHashLabel = pg.theme.H5("dddd")
	pg.txHashLabel.Color = pg.theme.Success

	pg.destinationAddressErrorLabel = pg.theme.Caption("")
	pg.sendErrorLabel = pg.theme.Caption("")
	pg.sendAmountErrorLabel = pg.theme.Caption("")
	pg.calculateErrorLabel = pg.theme.Caption("")
	pg.destinationAddressErrorLabel.Color = ui.DangerColor
	pg.sendErrorLabel.Color = ui.DangerColor
	pg.sendAmountErrorLabel.Color = ui.DangerColor
	pg.calculateErrorLabel.Color = ui.DangerColor

	pg.accountSelectorButtons = map[string]*widget.Button{}

	pg.passwordModal = theme.Password()
	pg.destinationAddressEditor = &editor{
		editor:   new(widget.Editor),
		material: pg.theme.Editor("Destination Address"),
		line:     pg.theme.Line(),
	}

	// TODO remvoe following line
	pg.destinationAddressEditor.editor.SetText("TseCXEcbPdSbDY2ZU97XuUuQCNLHcpyq3iH")

	pg.sendAmountEditor = &editor{
		editor:   new(widget.Editor),
		material: pg.theme.Editor("Amount to be sent"),
		line:     pg.theme.Line(),
	}

	pg.selectAccountButton = &button{
		button:   new(widget.Button),
		material: pg.theme.Button("Change account"),
	}
	pg.selectAccountButton.material.Background = ui.WhiteColor
	pg.selectAccountButton.material.Color = ui.LightBlueColor

	pg.nextButton = &button{
		button:   new(widget.Button),
		material: pg.theme.Button("Next"),
	}

	pg.confirmButton = &button{
		button:   new(widget.Button),
		material: pg.theme.Button("Confirm"),
	}
}

// Draw renders all of this page's widgets
func (pg *Send) Draw(gtx *layout.Context) interface{} {
	if pg.wallets == nil {
		pg.waitAndSetWalletInfo()
	}

	pg.validate(true)
	pg.watchAndUpdateValues(gtx)

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
				pg.destinationAddressLabel.Layout(gtx)
				inset := layout.Inset{
					Top: unit.Dp(20),
				}
				inset.Layout(gtx, func() {
					pg.destinationAddressEditor.layout(gtx)
				})

				if pg.destinationAddressErrorLabel.Text != "" {
					inset := layout.Inset{
						Top: unit.Dp(40),
					}
					inset.Layout(gtx, func() {
						pg.destinationAddressErrorLabel.Layout(gtx)
					})
				}
			},
			func() {
				pg.sendAmountLabel.Layout(gtx)
				inset := layout.Inset{
					Top: unit.Dp(20),
				}
				inset.Layout(gtx, func() {
					pg.sendAmountEditor.layout(gtx)
				})

				if pg.sendAmountErrorLabel.Text != "" {
					inset := layout.Inset{
						Top: unit.Dp(40),
					}
					inset.Layout(gtx, func() {
						pg.sendAmountErrorLabel.Layout(gtx)
					})
				}
			},
		},
		pg.getTransactionDetailWidgets(gtx),
		{
			func() {
				pg.calculateErrorLabel.Layout(gtx)
			},
			func() {
				for pg.nextButton.button.Clicked(gtx) {
					if pg.validate(false) && pg.calculateErrorLabel.Text == "" {
						pg.isConfirmationModalOpen = true
					}
				}

				pg.nextButton.layout(gtx)
			},
		},
	}

	widgets := reduceWidgetsWrapper(widgetsWrapper)

	list := layout.List{Axis: layout.Vertical}
	list.Layout(gtx, len(widgets), func(i int) {
		layout.UniformInset(unit.Dp(10)).Layout(gtx, widgets[i])
	})

	if pg.isAccountModalOpen {
		pg.drawAccountsModal(gtx)
	} else if pg.isConfirmationModalOpen {
		pg.drawConfirmationModal(gtx)

		if pg.isPasswordModalOpen {
			pg.drawPasswordModal(gtx)
		}
	}

	return nil
}

func (pg *Send) drawPasswordModal(gtx *layout.Context) {
	pg.theme.Modal(gtx, func() {
		pg.passwordModal.Draw(gtx, pg.confirm, pg.cancel)
	})
}

func (pg *Send) confirm(password string) {
	pg.isBroadcastingTransaction = true
	pg.isPasswordModalOpen = false

	pg.wallet.BroadcastTransaction(pg.txAuthor, password)
}

func (pg *Send) cancel() {
	pg.isPasswordModalOpen = false
}

func reduceWidgetsWrapper(widgetsWrapper [][]func()) []func() {
	widgets := []func(){}
	for i := range widgetsWrapper {
		widgets = append(widgets, widgetsWrapper[i]...)
	}

	return widgets
}

func (pg *Send) getTransactionDetailWidgets(gtx *layout.Context) []func() {
	widgets := []func(){
		func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					pg.transactionFeeLabel.Layout(gtx)
				}),
				layout.Flexed(1, func() {
					layout.Stack{Alignment: layout.NE}.Layout(gtx,
						layout.Stacked(func() {
							layout.Align(layout.Center).Layout(gtx, func() {
								pg.transactionFeeValueLabel.Layout(gtx)
							})
						}),
					)
				}),
			)
		},
		func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					pg.totalCostLabel.Layout(gtx)
				}),
				layout.Flexed(1, func() {
					layout.Stack{Alignment: layout.NE}.Layout(gtx,
						layout.Stacked(func() {
							layout.Align(layout.Center).Layout(gtx, func() {
								pg.totalCostValueLabel.Layout(gtx)
							})
						}),
					)
				}),
			)
		},
		func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					pg.balanceAfterSendLabel.Layout(gtx)
				}),
				layout.Flexed(1, func() {
					layout.Stack{Alignment: layout.NE}.Layout(gtx,
						layout.Stacked(func() {
							layout.Align(layout.Center).Layout(gtx, func() {
								pg.balanceAfterSendValueLabel.Layout(gtx)
							})
						}),
					)
				}),
			)
		},
	}

	return widgets
}

func (pg *Send) drawSelectedAccountSection(gtx *layout.Context) {
	widgets := []func(){
		func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					inset := layout.Inset{
						Top: unit.Dp(8.5),
					}
					inset.Layout(gtx, func() {
						pg.sourceAccountLabel.Layout(gtx)
					})
				}),
				layout.Rigid(func() {
					for pg.selectAccountButton.button.Clicked(gtx) {
						pg.isAccountModalOpen = true
					}
					pg.selectAccountButton.layout(gtx)
				}),
			)
		},
		func() {
			pg.selectedAccountLabel.Layout(gtx)
		},
		func() {
			pg.selectedWalletLabel.Layout(gtx)
		},
	}

	(&layout.List{Axis: layout.Vertical}).Layout(gtx, len(widgets), func(i int) {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, widgets[i])
	})
}

func (pg *Send) drawAccountsModal(gtx *layout.Context) {
	pg.theme.Modal(gtx, func() {
		pg.theme.H4("Choose sending account:").Layout(gtx)

		inset := layout.Inset{
			Top: unit.Dp(40),
		}
		inset.Layout(gtx, func() {
			(&layout.List{Axis: layout.Vertical}).Layout(gtx, len(pg.wallets), func(i int) {
				wallet := pg.wallets[i]
				pg.theme.H5(fmt.Sprintf("%s - %s", wallet.Name, dcrutil.Amount(wallet.TotalBalance).String())).Layout(gtx)

				(&layout.List{Axis: layout.Vertical}).Layout(gtx, len(wallet.Accounts), func(k int) {
					account := pg.wallets[i].Accounts[k]

					buttonKey := wallet.Name + account.Name
					pg.registerAccountSelectorButton(buttonKey)

					for pg.accountSelectorButtons[buttonKey].Clicked(gtx) {
						pg.setSelectedAccount(wallet, account)
						pg.isAccountModalOpen = false
					}

					layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func() {
							inset := layout.Inset{
								Left: unit.Dp(10),
							}
							inset.Layout(gtx, func() {
								inset := layout.Inset{
									Top: unit.Dp(25),
								}
								inset.Layout(gtx, func() {
									pg.theme.H6(fmt.Sprintf("%s %s", account.Name, dcrutil.Amount(account.TotalBalance).String())).Layout(gtx)
								})

								inset = layout.Inset{
									Top: unit.Dp(50),
								}
								inset.Layout(gtx, func() {
									pg.theme.Body1(fmt.Sprintf("Spendable: %s", dcrutil.Amount(account.SpendableBalance).String())).Layout(gtx)
								})
							})
						}),
					)
					pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
					pg.accountSelectorButtons[buttonKey].Layout(gtx)
				})
			})
		})
	})
}

func (pg *Send) drawConfirmationModal(gtx *layout.Context) {
	widgetsWrapper := [][]func(){
		{
			func() {
				pg.confirmModalTitleLabel.Layout(gtx)
			},
			func() {
				if pg.sendErrorLabel.Text != "" {
					pg.sendErrorLabel.Layout(gtx)
				} else {

				}
			},
			func() {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				layout.Align(layout.Center).Layout(gtx, func() {
					pg.theme.Body1(fmt.Sprintf("Sending from %s (%s)", pg.selectedAccountLabel.Text, pg.selectedWalletLabel.Text)).Layout(gtx)
				})
			},
			func() {},
			func() {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				layout.Align(layout.Center).Layout(gtx, func() {
					pg.confirmDestinationAddressLabel.Layout(gtx)
				})
			},
			func() {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				layout.Align(layout.Center).Layout(gtx, func() {
					pg.theme.Body1(pg.destinationAddressEditor.editor.Text()).Layout(gtx)
				})
			},
		},
		pg.getTransactionDetailWidgets(gtx),
		{
			func() {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				layout.Align(layout.Center).Layout(gtx, func() {
					pg.confirmWarningLabel.Layout(gtx)
				})
			},
			func() {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max

				for pg.confirmButton.button.Clicked(gtx) {
					pg.isPasswordModalOpen = true
				}
				pg.confirmButton.layout(gtx)
			},
		},
	}

	widgets := reduceWidgetsWrapper(widgetsWrapper)

	pg.theme.Modal(gtx, func() {
		(&layout.List{Axis: layout.Vertical}).Layout(gtx, len(widgets), func(i int) {
			layout.UniformInset(unit.Dp(10)).Layout(gtx, widgets[i])
		})
	})
}

func (pg *Send) waitAndSetWalletInfo() {
	walletInfoState := pg.states[StateWalletInfo]
	if walletInfoState == nil {
		return
	}

	walletInfo := walletInfoState.(*wallet.MultiWalletInfo)
	pg.setDefaultSendAccount(walletInfo.Wallets)
}

func (pg *Send) registerAccountSelectorButton(accountName string) {
	if _, ok := pg.accountSelectorButtons[accountName]; !ok {
		pg.accountSelectorButtons[accountName] = new(widget.Button)
	}
}

func (pg *Send) setDefaultSendAccount(wallets []wallet.InfoShort) {
	pg.wallets = wallets

	for i := range wallets {
		if len(wallets[i].Accounts) == 0 {
			continue
		}

		pg.setSelectedAccount(wallets[i], wallets[i].Accounts[0])
		pg.balanceAfterSendValueLabel.Text = dcrutil.Amount(pg.selectedWallet.SpendableBalance).String()
		break
	}
}

func (pg *Send) setSelectedAccount(wallet wallet.InfoShort, account wallet.Account) {
	pg.selectedWallet = &wallet
	pg.selectedAccount = &account

	pg.selectedWalletLabel.Text = wallet.Name
	pg.selectedAccountLabel.Text = fmt.Sprintf("%s - %s", account.Name, dcrutil.Amount(account.SpendableBalance).String())

	// create a mew transaction everytime a new account is chosen
	pg.wallet.CreateTransaction(wallet.ID, account.Number)
	pg.calculateValues()
}

func (pg *Send) validateDestinationAddress(ignoreEmpty bool) bool {
	pg.destinationAddressErrorLabel.Text = ""

	destinationAddress := pg.destinationAddressEditor.editor.Text()
	if destinationAddress == "" {
		if !ignoreEmpty {
			pg.destinationAddressErrorLabel.Text = "please enter a destination address"
		}
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

	amount := pg.sendAmountEditor.editor.Text()
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
		pg.nextButton.material.Background = ui.GrayColor
		return false
	}

	pg.nextButton.material.Background = ui.LightBlueColor
	return true
}

func (pg *Send) watchForBroadcastResult(gtx *layout.Context) {
	err := pg.states[StateError]
	hash := pg.states[StateTxHash]

	if err == nil && hash == nil {
		return
	}

	if err != nil {
		pg.sendErrorLabel.Text = fmt.Sprintf("error broadcasting transaction: %s", err.(error).Error())
		delete(pg.states, StateError)
	} else if hash != nil {
		pg.txHashLabel.Text = fmt.Sprintf("The transaction was published successfully. Hash: %s", hash.(*wallet.TxHash).Hash)
		pg.destinationAddressEditor.editor.SetText("")
		pg.sendAmountEditor.editor.SetText("")
		pg.isConfirmationModalOpen = false

		delete(pg.states, StateTxHash)
	}

	pg.isBroadcastingTransaction = false
}

func (pg *Send) watchAndUpdateValues(gtx *layout.Context) {
	pg.watchForBroadcastResult(gtx)

	pg.confirmButton.material.Text = "Send"
	pg.confirmButton.material.Background = ui.LightBlueColor

	txAuthor := pg.states[StateTxAuthor]
	if txAuthor != nil {
		pg.txAuthor = txAuthor.(*dcrlibwallet.TxAuthor)
		delete(pg.states, StateTxAuthor)
	}

	if pg.isBroadcastingTransaction {
		pg.confirmButton.material.Text = "Sending..."
		pg.confirmButton.material.Background = ui.GrayColor
	}

	for range pg.destinationAddressEditor.editor.Events(gtx) {
		pg.calculateValues()
	}

	for range pg.sendAmountEditor.editor.Events(gtx) {
		pg.calculateValues()
	}
}

func (pg *Send) calculateValues() {
	pg.transactionFeeValueLabel.Text = "0 DCR"
	pg.totalCostValueLabel.Text = "0 DCR "
	pg.calculateErrorLabel.Text = ""

	if pg.selectedWallet != nil {
		pg.balanceAfterSendValueLabel.Text = dcrutil.Amount(pg.selectedWallet.SpendableBalance).String()
	}

	if pg.txAuthor == nil || !pg.validate(true) {
		return
	}

	amountDCR, _ := strconv.ParseFloat(pg.sendAmountEditor.editor.Text(), 64)
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
	pg.txAuthor.AddSendDestination(pg.destinationAddressEditor.editor.Text(), amountAtoms, false)

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

func (e *editor) layout(gtx *layout.Context) {
	e.material.Layout(gtx, e.editor)
	inset := layout.Inset{
		Top: unit.Dp(float32(gtx.Constraints.Height.Min)),
	}
	inset.Layout(gtx, func() {
		e.line.Draw(gtx)
	})
}

func (e *button) layout(gtx *layout.Context) {
	e.material.Layout(gtx, e.button)
}

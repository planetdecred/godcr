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
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// SendID is the id of the send page
const SendID = "send"

type confirmModalWidgets struct {
	line       *materialplus.Line
	sendButton *widget.Button

	confirmLabel       material.Label
	sendingFromLabel   material.Label
	toDestinationLabel material.Label
	sendWarningLabel   material.Label
}

type accountModalWidgets struct {
	titleLabel material.Label
	titleLine  *materialplus.Line
}

// Send represents the send page of the app.
// It should only be accessible if the app finds
// at least one wallet.
type Send struct {
	theme                      *materialplus.Theme
	container                  layout.List
	wallets                    []wallet.InfoShort
	isShowingConfirmationModal bool
	isAccountModalOpen         bool
	wallet                     *wallet.Wallet
	transaction                *dcrlibwallet.TxAuthor

	// selected account values
	selectedWalletID  int
	selectedAccountID int32

	selectedWallet  *wallet.InfoShort
	selectedAccount *wallet.Account

	// labels
	loadingLabel               material.Label
	titleLabel                 material.Label
	fromLabel                  material.Label
	destinationAddressLabel    material.Label
	sendAmountLabel            material.Label
	txFeeLabel                 material.Label
	txFeeValueLabel            material.Label
	totalCostLabel             material.Label
	totalCostValueLabel        material.Label
	remainingBalanceLabel      material.Label
	remainingBalanceValueLabel material.Label
	amountLabel                material.Label
	destinationLabel           material.Label

	// confirm modal labels
	confirmLabel                   material.Label
	confirmSendingFromLabel        material.Label
	confirmDestinationAddressLabel material.Label
	confirmWarningLabel            material.Label

	// error labels
	destinationAddressErrorLabel material.Label
	amountErrorLabel             material.Label

	// selected account labels
	sendWalletNameLabel              material.Label
	sendAccountNameLabel             material.Label
	sendAccountSpendableBalanceLabel material.Label

	// editors
	destinationAddressEditor *widget.Editor
	amountEditor             *widget.Editor

	// buttons
	openAccountModalButton *widget.Button
	nextButton             *widget.Button
	confirmButton          *widget.Button

	accountSelectorButtons map[string]*widget.Button

	// modals
	accountModalWidgets *accountModalWidgets
	confirmModalWidgets *confirmModalWidgets

	txAuthor *dcrlibwallet.TxAuthor

	// state
	states map[string]interface{}
}

// Init initializes this page's widgets
func (pg *Send) Init(theme *materialplus.Theme, wal *wallet.Wallet, states map[string]interface{}) {
	pg.theme = theme
	pg.states = states
	pg.container.Axis = layout.Vertical
	pg.isShowingConfirmationModal = false
	pg.isAccountModalOpen = false
	pg.wallet = wal

	// main labels
	pg.titleLabel = theme.Label(units.Label, "Send DCR")
	pg.loadingLabel = theme.Caption("loading...")
	pg.fromLabel = theme.Body1("From:")
	pg.destinationAddressLabel = theme.Body1("Destination Address:")
	pg.sendAmountLabel = theme.Body1("Amount")
	pg.txFeeLabel = theme.Body1("Transaction Fee:")
	pg.txFeeValueLabel = theme.Body1("0 DCR")
	pg.totalCostLabel = theme.Body1("Total Cost:")
	pg.totalCostValueLabel = theme.Body1("0 DCR")
	pg.txFeeValueLabel = theme.Body1("0 DCR")
	pg.remainingBalanceLabel = theme.Body1("Balance after send")
	pg.remainingBalanceValueLabel = theme.Body1("0 DCR")
	pg.amountLabel = theme.Body1("0 DCR")
	pg.destinationLabel = theme.Body1("")

	// confirm modal labels
	pg.confirmLabel = theme.Body1("Confirm to send")
	pg.confirmSendingFromLabel = theme.Body1("")
	pg.confirmDestinationAddressLabel = theme.Body1("To destination address")
	pg.confirmWarningLabel = theme.Caption("Your DCR will be sent and CANNOT be undone.")

	// error labels
	pg.destinationAddressErrorLabel = theme.Caption("")
	pg.destinationAddressErrorLabel.Color = ui.DangerColor
	pg.amountErrorLabel = theme.Caption("")
	pg.amountErrorLabel.Color = ui.DangerColor

	// selected account labels
	pg.sendWalletNameLabel = theme.Body1("")
	pg.sendAccountNameLabel = theme.Body1("")
	pg.sendAccountSpendableBalanceLabel = theme.Body1("")

	// editors
	pg.destinationAddressEditor = new(widget.Editor)
	pg.destinationAddressEditor.SetText("TseCXEcbPdSbDY2ZU97XuUuQCNLHcpyq3iH")
	pg.amountEditor = new(widget.Editor)

	// buttons
	pg.openAccountModalButton = new(widget.Button)
	pg.nextButton = new(widget.Button)
	pg.confirmButton = new(widget.Button)

	pg.accountSelectorButtons = map[string]*widget.Button{}

	// accountModalWidgets
	pg.accountModalWidgets = &accountModalWidgets{
		titleLabel: theme.H3("Choose a sending account"),
		titleLine:  theme.Line(),
	}
}

func (pg *Send) initModalWidgets(theme *materialplus.Theme) {
	pg.confirmModalWidgets = &confirmModalWidgets{
		line:               theme.Line(),
		confirmLabel:       theme.Body1("Confirm to send"),
		sendingFromLabel:   theme.Body1("Sending from Default (wallet-2)"),
		toDestinationLabel: theme.Body2("To destination address"),
		sendWarningLabel:   theme.Caption("Your DCR will be sent and CANNOT be undone."),
		sendButton:         new(widget.Button),
	}
}

// Draw renders all of this page's widgets
func (pg *Send) Draw(gtx *layout.Context) interface{} {
	pg.validate(true)
	pg.watchAndUpdateValues(gtx)

	// set wallet options
	if pg.wallets == nil {
		walletInfo := pg.states[StateWalletInfo].(*wallet.MultiWalletInfo)
		if len(walletInfo.Wallets) > 0 {
			pg.setDefaultSendAccount(walletInfo.Wallets)
		}
	}

	widgetFuncs := []func(){
		func() {
			pg.titleLabel.Layout(gtx)
		},
		func() {
			layout.UniformInset(unit.Dp(5)).Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func() {
						inset := layout.Inset{
							Top: unit.Dp(8.5),
						}
						inset.Layout(gtx, func() {
							pg.fromLabel.Layout(gtx)
						})
					}),
					layout.Rigid(func() {
						btn := pg.theme.Button("(select)")
						btn.Background = ui.WhiteColor
						btn.Color = ui.LightBlueColor

						for pg.openAccountModalButton.Clicked(gtx) {
							pg.isAccountModalOpen = true
						}

						gtx.Constraints.Height.Max = 35
						btn.Layout(gtx, pg.openAccountModalButton)
					}),
				)

				inset := layout.Inset{
					Top: unit.Dp(35),
				}
				inset.Layout(gtx, func() {
					layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func() {
							pg.sendAccountNameLabel.Layout(gtx)
						}),
						layout.Rigid(func() {
							inset := layout.Inset{
								Left: unit.Dp(30),
							}
							inset.Layout(gtx, func() {
								pg.sendAccountSpendableBalanceLabel.Layout(gtx)
							})
						}),
					)
				})

				inset = layout.Inset{
					Top: unit.Dp(55),
				}
				inset.Layout(gtx, func() {
					pg.sendWalletNameLabel.Layout(gtx)
				})
			})
		},
		func() {
			layout.UniformInset(unit.Dp(5)).Layout(gtx, func() {
				pg.destinationAddressLabel.Layout(gtx)
				inset := layout.Inset{
					Top: unit.Dp(20),
				}
				inset.Layout(gtx, func() {
					pg.theme.Editor("Destination Address").Layout(gtx, pg.destinationAddressEditor)

					inset := layout.Inset{
						Top: unit.Dp(15),
					}
					inset.Layout(gtx, func() {
						pg.destinationAddressErrorLabel.Layout(gtx)
					})
				})
			})
		},
		func() {
			layout.UniformInset(unit.Dp(5)).Layout(gtx, func() {
				pg.sendAmountLabel.Layout(gtx)
				inset := layout.Inset{
					Top: unit.Dp(20),
				}
				inset.Layout(gtx, func() {
					pg.theme.Editor("Amount").Layout(gtx, pg.amountEditor)

					inset := layout.Inset{
						Top: unit.Dp(15),
					}
					inset.Layout(gtx, func() {
						pg.amountErrorLabel.Layout(gtx)
					})
				})
			})
		},
		func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					pg.txFeeLabel.Layout(gtx)
				}),
				layout.Flexed(1, func() {
					layout.Stack{Alignment: layout.NE}.Layout(gtx,
						layout.Stacked(func() {
							layout.Align(layout.Center).Layout(gtx, func() {
								pg.txFeeValueLabel.Layout(gtx)
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

			inset := layout.Inset{
				Top: unit.Dp(30),
			}
			inset.Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func() {
						pg.remainingBalanceLabel.Layout(gtx)
					}),
					layout.Flexed(1, func() {
						layout.Stack{Alignment: layout.NE}.Layout(gtx,
							layout.Stacked(func() {
								pg.remainingBalanceValueLabel.Layout(gtx)
							}),
						)
					}),
				)
			})
		},
		func() {
			for pg.nextButton.Clicked(gtx) {
				if pg.validate(false) {
					pg.isShowingConfirmationModal = true
				}
			}

			btn := pg.theme.Button("Next")
			if pg.validate(true) {
				btn.Background = ui.LightBlueColor
			} else {
				btn.Background = ui.GrayColor
			}
			btn.Layout(gtx, pg.nextButton)
		},
	}

	pg.container.Layout(gtx, len(widgetFuncs), func(i int) {
		layout.UniformInset(unit.Dp(10)).Layout(gtx, widgetFuncs[i])
	})

	if pg.isAccountModalOpen {
		pg.drawAccountsModal(gtx)
	} else if pg.isShowingConfirmationModal {
		pg.drawConfirmationModal(gtx)
	}

	return nil
}

func (pg *Send) addAccountSelectorButton(accountName string) {
	if _, ok := pg.accountSelectorButtons[accountName]; !ok {
		pg.accountSelectorButtons[accountName] = new(widget.Button)
	}
}

func (pg *Send) drawAccountsModal(gtx *layout.Context) {
	pg.theme.Modal(gtx, func() {
		list := layout.List{Axis: layout.Vertical}
		list.Layout(gtx, len(pg.wallets), func(i int) {
			wallet := pg.wallets[i]
			pg.theme.H5(fmt.Sprintf("%s - %s", wallet.Name, dcrutil.Amount(wallet.TotalBalance).String())).Layout(gtx)

			list := layout.List{Axis: layout.Vertical}
			list.Layout(gtx, len(wallet.Accounts), func(k int) {
				account := pg.wallets[i].Accounts[k]

				pg.addAccountSelectorButton(wallet.Name + account.Name)

				for pg.accountSelectorButtons[wallet.Name+account.Name].Clicked(gtx) {
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
				pg.accountSelectorButtons[wallet.Name+account.Name].Layout(gtx)
			})
		})
	})
}

func (pg *Send) drawConfirmationModal(gtx *layout.Context) {
	modalWidgetFuncs := []func(){
		func() {
			layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
				pg.confirmLabel.Layout(gtx)
			})
		},
		func() {
			inset := layout.Inset{
				Top:    unit.Dp(1),
				Bottom: unit.Dp(1),
			}
			inset.Layout(gtx, func() {
				//pg.confirmModalWidgets.line.Draw(gtx)
			})
		},
		func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			layout.Align(layout.Center).Layout(gtx, func() {
				inset := layout.Inset{
					Top: unit.Dp(5),
				}
				inset.Layout(gtx, func() {
					pg.confirmSendingFromLabel.Layout(gtx)
				})
			})
		},
		func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			layout.Align(layout.Center).Layout(gtx, func() {
				inset := layout.Inset{
					Top: unit.Dp(5),
				}
				inset.Layout(gtx, func() {
					pg.theme.H5(pg.amountEditor.Text() + " DCR").Layout(gtx)
				})
			})
		},
		func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			layout.Align(layout.Center).Layout(gtx, func() {
				inset := layout.Inset{
					Top: unit.Dp(40),
				}
				inset.Layout(gtx, func() {
					pg.confirmDestinationAddressLabel.Layout(gtx)
				})
			})
		},
		func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			layout.Align(layout.Center).Layout(gtx, func() {
				inset := layout.Inset{
					Top: unit.Dp(6),
				}
				inset.Layout(gtx, func() {
					pg.theme.Body1(pg.destinationAddressEditor.Text()).Layout(gtx)
				})
			})
		},
		func() {
			inset := layout.Inset{
				Top: unit.Dp(25),
			}
			inset.Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func() {
						pg.txFeeLabel.Layout(gtx)
					}),
					layout.Flexed(1, func() {
						layout.Stack{Alignment: layout.NE}.Layout(gtx,
							layout.Stacked(func() {
								pg.txFeeValueLabel.Layout(gtx)
							}),
						)
					}),
				)
			})
		},
		func() {
			inset := layout.Inset{
				Top: unit.Dp(15),
			}
			inset.Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func() {
						pg.totalCostLabel.Layout(gtx)
					}),
					layout.Flexed(1, func() {
						layout.Stack{Alignment: layout.NE}.Layout(gtx,
							layout.Stacked(func() {
								pg.totalCostValueLabel.Layout(gtx)
							}),
						)
					}),
				)
			})
		},
		func() {
			inset := layout.Inset{
				Top: unit.Dp(15),
			}
			inset.Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func() {
						pg.remainingBalanceLabel.Layout(gtx)
					}),
					layout.Flexed(1, func() {
						layout.Stack{Alignment: layout.NE}.Layout(gtx,
							layout.Stacked(func() {
								pg.totalCostValueLabel.Layout(gtx)
							}),
						)
					}),
				)
			})
		},
		func() {
			inset := layout.Inset{
				Top: unit.Dp(35),
			}
			inset.Layout(gtx, func() {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				layout.Align(layout.Center).Layout(gtx, func() {
					pg.confirmWarningLabel.Layout(gtx)
				})
			})
		},
		func() {
			inset := layout.Inset{
				Top: unit.Dp(7),
			}
			inset.Layout(gtx, func() {
				btn := pg.theme.Button("Send " + pg.amountEditor.Text() + " DCR")

				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				gtx.Constraints.Height.Min = 50
				btn.Layout(gtx, pg.confirmButton)
			})
		},
	}

	inset := layout.Inset{
		Top:  unit.Dp(0),
		Left: unit.Dp(0),
	}
	inset.Layout(gtx, func() {
		pg.theme.Modal(gtx, func() {
			list := layout.List{Axis: layout.Vertical}
			list.Layout(gtx, len(modalWidgetFuncs), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, modalWidgetFuncs[i])
			})
		})
	})
}

func (pg *Send) setDefaultSendAccount(wallets []wallet.InfoShort) {
	pg.wallets = wallets

	for i := range wallets {
		if len(wallets[i].Accounts) == 0 {
			continue
		}

		pg.setSelectedAccount(wallets[i], wallets[i].Accounts[0])
		break
	}

	pg.remainingBalanceValueLabel.Text = dcrutil.Amount(pg.selectedWallet.SpendableBalance).String()
}

func (pg *Send) setSelectedAccount(wallet wallet.InfoShort, account wallet.Account) {
	pg.selectedWallet = &wallet
	pg.selectedAccount = &account

	pg.sendWalletNameLabel.Text = wallet.Name
	pg.sendAccountNameLabel.Text = account.Name
	pg.sendAccountSpendableBalanceLabel.Text = dcrutil.Amount(account.SpendableBalance).String()

	// create a mew transaction everytime a new account is chosen
	pg.wallet.CreateTransaction(wallet.ID, account.Number, dcrlibwallet.DefaultRequiredConfirmations)
	pg.calculateValues()
}

func (pg *Send) validateDestinationAddress(ignoreEmpty bool) bool {
	destinationAddress := pg.destinationAddressEditor.Text()
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

	pg.destinationAddressErrorLabel.Text = ""
	return true
}

func (pg *Send) validateAmount(ignoreEmpty bool) bool {
	amount := pg.amountEditor.Text()
	if amount == "" {
		if !ignoreEmpty {
			pg.amountErrorLabel.Text = "please enter a send amount"
		}
		return false
	}

	if amount != "" {
		_, err := strconv.ParseFloat(pg.amountEditor.Text(), 64)
		if err != nil {
			pg.amountErrorLabel.Text = "please enter a valid amount"
			return false
		}
	}

	pg.amountErrorLabel.Text = ""
	return true
}

func (pg *Send) validate(ignoreEmpty bool) bool {
	isAddressValid := pg.validateDestinationAddress(ignoreEmpty)
	isAmountValid := pg.validateAmount(ignoreEmpty)

	if !isAddressValid || !isAmountValid {
		return false
	}

	return true
}

func (pg *Send) watchAndUpdateValues(gtx *layout.Context) {
	txAuthor := pg.states[StateTxAuthor]
	if txAuthor != nil {
		pg.txAuthor = txAuthor.(*dcrlibwallet.TxAuthor)
		delete(pg.states, StateTxAuthor)
	}

	for range pg.destinationAddressEditor.Events(gtx) {
		pg.calculateValues()
	}

	for range pg.amountEditor.Events(gtx) {
		pg.calculateValues()
	}
}

func (pg *Send) calculateValues() {
	pg.txFeeValueLabel.Text = "0 DCR"
	pg.totalCostValueLabel.Text = "0 DCR"

	if pg.selectedWallet != nil {
		pg.remainingBalanceValueLabel.Text = dcrutil.Amount(pg.selectedWallet.SpendableBalance).String()
	}

	if pg.txAuthor == nil || !pg.validate(true) {
		return
	}

	amountDCR, _ := strconv.ParseFloat(pg.amountEditor.Text(), 64)
	if amountDCR <= 0 {
		return
	}

	amount, err := dcrutil.NewAmount(amountDCR)
	if err != nil {
		fmt.Println(err)
		return
	}

	amountAtoms := int64(amount)

	// set destination address
	pg.txAuthor.RemoveSendDestination(0)
	pg.txAuthor.AddSendDestination(pg.destinationAddressEditor.Text(), amountAtoms, false)

	// calculate transaction fee
	feeAndSize, err := pg.txAuthor.EstimateFeeAndSize()
	if err != nil {
		fmt.Println(err)
		return
	}

	txFee := feeAndSize.Fee.AtomValue
	totalCost := txFee + amountAtoms
	remainingBalance := pg.selectedWallet.SpendableBalance - totalCost

	pg.txFeeValueLabel.Text = dcrutil.Amount(txFee).String()
	pg.totalCostValueLabel.Text = dcrutil.Amount(totalCost).String()
	pg.remainingBalanceValueLabel.Text = dcrutil.Amount(remainingBalance).String()
}

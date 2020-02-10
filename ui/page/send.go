package page

import (
	//"fmt"
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

type modalWidgets struct {
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

	// calculated values
	transactionFee             int64
	totalCost                  int64
	remainingBalance           int64
	previousAmount             int64
	previousDestinationAddress string

	// labels
	loadingLabel            material.Label
	titleLabel              material.Label
	fromLabel               material.Label
	destinationAddressLabel material.Label
	sendAmountLabel         material.Label
	txFeeLabel              material.Label
	totalCostLabel          material.Label
	remainingBalanceLabel   material.Label

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

	accountSelectorButtons map[string]*widget.Button
	// modals
	accountModalWidgets *accountModalWidgets

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

	// calculated values
	pg.transactionFee = 0
	pg.totalCost = 0
	pg.previousAmount = -1
	pg.previousDestinationAddress = ""

	// labels
	pg.titleLabel = theme.Label(units.Label, "Send DCR")
	pg.loadingLabel = theme.Caption("loading...")
	pg.fromLabel = theme.Body1("From:")
	pg.destinationAddressLabel = theme.Body1("Destination Address:")
	pg.sendAmountLabel = theme.Body1("Amount")
	pg.txFeeLabel = theme.Body1("Transaction Fee:")
	pg.totalCostLabel = theme.Body1("Total Cost:")
	pg.remainingBalanceLabel = theme.Body1("Balance after send")

	// error labels
	pg.destinationAddressErrorLabel = theme.Caption("error")
	pg.destinationAddressErrorLabel.Color = ui.DangerColor
	pg.amountErrorLabel = theme.Caption("error")
	pg.amountErrorLabel.Color = ui.DangerColor

	// selected account labels
	pg.sendWalletNameLabel = theme.Body1("")
	pg.sendAccountNameLabel = theme.Body1("")
	pg.sendAccountSpendableBalanceLabel = theme.Body1("")

	// editors
	pg.destinationAddressEditor = new(widget.Editor)
	pg.amountEditor = new(widget.Editor)

	// buttons
	pg.openAccountModalButton = new(widget.Button)
	pg.nextButton = new(widget.Button)

	pg.accountSelectorButtons = map[string]*widget.Button{}

	// accountModalWidgets
	pg.accountModalWidgets = &accountModalWidgets{
		titleLabel: theme.H3("Choose a sending account"),
		titleLine:  theme.Line(),
	}
}

func (pg *Send) initModalWidgets(theme *materialplus.Theme) {
	/**pg.modalWidgets = &modalWidgets{
		line:               theme.Line(),
		confirmLabel:       theme.Body1("Confirm to send"),
		sendingFromLabel:   theme.Body1("Sending from Default (wallet-2)"),
		toDestinationLabel: theme.Body2("To destination address"),
		sendWarningLabel:   theme.Caption("Your DCR will be sent and CANNOT be undone."),
		sendButton:         new(widget.Button),
	}**/
}

func (pg *Send) initValidationWidgets(theme *materialplus.Theme) {
	/**destinationAddressLabel := theme.Body2("")
	destinationAddressLabel.Color = ui.DangerColor

	amountLabel := theme.Body2("")
	amountLabel.Color = ui.DangerColor

	pg.validationWidgets = &validationWidgets{
		destinationAddressLabel: destinationAddressLabel,
		amountLabel:             amountLabel,
	}**/
}

// Draw renders all of this page's widgets
func (pg *Send) Draw(gtx *layout.Context) interface{} {
	go pg.validate(false)
	pg.updateValuesAndTx(gtx)

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
								pg.theme.Body1(dcrutil.Amount(pg.transactionFee).String()).Layout(gtx)
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
								pg.theme.Body1(dcrutil.Amount(pg.totalCost).String()).Layout(gtx)
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
								pg.theme.Body1(dcrutil.Amount(pg.remainingBalance).String()).Layout(gtx)
							}),
						)
					}),
				)
			})
		},
		func() {
			hasPassedValidation := false

			if pg.destinationAddressEditor.Text() != "" && pg.amountEditor.Text() != "" {
				if pg.destinationAddressErrorLabel.Text == "" && pg.amountErrorLabel.Text == "" {
					hasPassedValidation = true
				}
			}

			for pg.nextButton.Clicked(gtx) {

			}

			btn := pg.theme.Button("Next")
			if hasPassedValidation {
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

func (pg *Send) drawAccountsModal(gtx *layout.Context) {
	pg.theme.Modal(gtx, func() {
		list := layout.List{Axis: layout.Vertical}
		list.Layout(gtx, len(pg.wallets)+1, func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
				if i == 0 {
					layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
						pg.accountModalWidgets.titleLabel.Layout(gtx)
					})
					return
				}

				wallet := pg.wallets[i-1]

				walletNameLabel := pg.theme.H5(wallet.Name + dcrutil.Amount(wallet.TotalBalance).String())
				walletNameLabel.Layout(gtx)

				list := layout.List{Axis: layout.Vertical}
				list.Layout(gtx, len(wallet.Accounts), func(i int) {
					account := wallet.Accounts[i]

					if _, ok := pg.accountSelectorButtons[account.Name]; !ok {
						pg.accountSelectorButtons[account.Name] = new(widget.Button)
					}

					for pg.accountSelectorButtons[account.Name].Clicked(gtx) {
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
									sendAccountNameLabel := pg.theme.H6(account.Name + "  " + dcrutil.Amount(account.TotalBalance).String())
									sendAccountNameLabel.Layout(gtx)
								})

								inset = layout.Inset{
									Top: unit.Dp(50),
								}
								inset.Layout(gtx, func() {
									spendableBalanceLabel := pg.theme.Body1("Spendable: " + dcrutil.Amount(account.SpendableBalance).String())
									spendableBalanceLabel.Layout(gtx)
								})
							})
						}),
						layout.Rigid(func() {

						}),
					)
					pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
					pg.accountSelectorButtons[account.Name].Layout(gtx)
				})
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
}

func (pg *Send) setSelectedAccount(wallet wallet.InfoShort, account wallet.Account) {
	pg.sendWalletNameLabel.Text = wallet.Name
	pg.sendAccountNameLabel.Text = account.Name
	pg.sendAccountSpendableBalanceLabel.Text = dcrutil.Amount(account.SpendableBalance).String()
}

func (pg *Send) updateValuesAndTx(gtx *layout.Context) {
	hasError := false
	if pg.destinationAddressErrorLabel.Text != "" && pg.amountErrorLabel.Text != "" {
		hasError = true
	}

	for range pg.destinationAddressEditor.Events(gtx) {
		if !hasError {
			pg.calculateTxFee(false)
		}
	}

	for range pg.amountEditor.Events(gtx) {
		if !hasError {
			pg.calculateTxFee(false)
		}
	}
}

func (pg *Send) calculateTxFee(createNewTx bool) {
	if createNewTx {
		pg.createTransaction()
	}
	/**pg.transaction = pg.wallet.CreateTransaction(pg.selectedWalletID, pg.selectedAccountID, dcrlibwallet.DefaultRequiredConfirmations)
	pg.transaction.AddSendDestination(pg.destinationAddressEditor.Text(), 0, false)

	txFee, err := pg.transaction.EstimateFeeAndSize()
	if err != nil {
		//handle error
		return
	}**/
}

func (pg *Send) createTransaction() {
	pg.wallet.CreateTransaction(pg.selectedWalletID, pg.selectedAccountID, dcrlibwallet.DefaultRequiredConfirmations)
}

func (pg *Send) validate(hasSubmitted bool) {
	destinationAddress := pg.destinationAddressEditor.Text()
	amount := pg.amountEditor.Text()

	pg.destinationAddressErrorLabel.Text = ""
	pg.amountErrorLabel.Text = ""

	if hasSubmitted {
		if destinationAddress == "" {
			pg.destinationAddressErrorLabel.Text = "please enter a destination address"
		}

		if amount == "" {
			pg.amountErrorLabel.Text = "please enter a send amount"
		}
	}

	if destinationAddress != "" {
		if isValid, _ := pg.wallet.IsAddressValid(destinationAddress); !isValid {
			pg.destinationAddressErrorLabel.Text = "invalid address"
		}
	}

	if amount != "" {
		_, err := strconv.Atoi(amount)
		if err != nil {
			pg.amountErrorLabel.Text = "please enter a valid amount"
		}
	}
}

func (pg *Send) drawConfirmationModal(gtx *layout.Context) {
	/**modalWidgetFuncs := []func(){
		func() {
			layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
				pg.modalWidgets.confirmLabel.Layout(gtx)
			})
		},
		func() {
			inset := layout.Inset{
				Top:    unit.Dp(1),
				Bottom: unit.Dp(1),
			}
			inset.Layout(gtx, func() {
				pg.modalWidgets.line.Draw(gtx)
			})
		},
		func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			layout.Align(layout.Center).Layout(gtx, func() {
				inset := layout.Inset{
					Top: unit.Dp(5),
				}
				inset.Layout(gtx, func() {
					pg.modalWidgets.sendingFromLabel.Layout(gtx)
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
					pg.modalWidgets.toDestinationLabel.Layout(gtx)
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
								pg.theme.Body1(dcrutil.Amount(pg.transactionFee).String()).Layout(gtx)
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
								// TODO get this value from wallet balance
								// i.e wallet balance - total cost
								pg.theme.Body2("4.37280441 DCR").Layout(gtx)
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
					pg.modalWidgets.sendWarningLabel.Layout(gtx)
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
				btn.Layout(gtx, pg.modalWidgets.sendButton)
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
	})**/
}

package page

import (
	"strconv"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/decred/dcrd/dcrutil"

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

type validationWidgets struct {
	destinationAddressLabel material.Label
	amountLabel             material.Label
}

// Send represents the send page of the app.
// It should only be accessible if the app finds
// at least one wallet.
type Send struct {
	theme     *materialplus.Theme
	container layout.List

	isShowingConfirmationModal bool

	headerLabel              material.Label
	fromLabel                material.Label
	toLabel                  material.Label
	amountLabel              material.Label
	txFeeLabel               material.Label
	txFeeValueLabel          material.Label
	totalCostLabel           material.Label
	totalCostValueLabel      material.Label
	remainingBalanceLabel    material.Label
	destinationAddressEditor *widget.Editor
	destinationAddressInput  material.Editor
	amountEditor             *widget.Editor
	amountInput              material.Editor
	nextButton               *widget.Button
	accountSelector          *materialplus.Select

	modalWidgets      *modalWidgets
	validationWidgets *validationWidgets

	transactionFee int64 // should be calculated by wallet backend
}

// Init initializes this page's widgets
func (pg *Send) Init(theme *materialplus.Theme, wal *wallet.Wallet) {
	pg.theme = theme

	pg.container.Axis = layout.Vertical
	pg.isShowingConfirmationModal = false

	pg.headerLabel = theme.Label(units.Label, "Send DCR")
	pg.fromLabel = theme.Body1("From")
	pg.toLabel = theme.Body1("To")
	pg.amountLabel = theme.Body1("Amount")
	pg.txFeeLabel = theme.Body1("Transaction fee")
	pg.txFeeValueLabel = theme.Body1("0.75 DCR")
	pg.totalCostLabel = theme.Body1("Total Cost")
	pg.totalCostValueLabel = theme.Body1("3.5442 DCR")
	pg.remainingBalanceLabel = theme.Body1("Balance after send")
	pg.destinationAddressEditor = new(widget.Editor)
	pg.destinationAddressInput = theme.Editor("Destination address")
	pg.amountEditor = new(widget.Editor)
	pg.amountInput = theme.Editor("0 DCR")
	pg.nextButton = new(widget.Button)

	dummyAccountsMap := map[string]string{
		"wallet-1": "100 DCR",
		"wallet-2": "7.645664DCR",
	}
	pg.accountSelector = theme.Select(dummyAccountsMap)

	// init modal widgets
	pg.initModalWidgets(theme)

	// init validation widgets
	pg.initValidationWidgets(theme)

	// set dummy data
	pg.amountEditor.SetText("3.1459265")
	pg.destinationAddressEditor.SetText("TsfDLrRkk9ciUuwfp2b8PawwnukYD7yAjGd")

	pg.transactionFee = 2510

}

func (pg *Send) initModalWidgets(theme *materialplus.Theme) {
	pg.modalWidgets = &modalWidgets{
		line:               theme.Line(),
		confirmLabel:       theme.Body1("Confirm to send"),
		sendingFromLabel:   theme.Body1("Sending from Default (wallet-2)"),
		toDestinationLabel: theme.Body2("To destination address"),
		sendWarningLabel:   theme.Caption("Your DCR will be sent and CANNOT be undone."),
		sendButton:         new(widget.Button),
	}
}

func (pg *Send) initValidationWidgets(theme *materialplus.Theme) {
	destinationAddressLabel := theme.Body2("")
	destinationAddressLabel.Color = ui.DangerColor

	amountLabel := theme.Body2("")
	amountLabel.Color = ui.DangerColor

	pg.validationWidgets = &validationWidgets{
		destinationAddressLabel: destinationAddressLabel,
		amountLabel:             amountLabel,
	}
}

// Draw renders all of this page's widgets
func (pg *Send) Draw(gtx *layout.Context, _ ...interface{}) interface{} {
	widgetFuncs := []func(){
		func() {
			pg.headerLabel.Layout(gtx)
		},
		func() {
			layout.UniformInset(unit.Dp(5)).Layout(gtx, func() {
				pg.fromLabel.Layout(gtx)

				inset := layout.Inset{
					Top: unit.Dp(25),
				}
				inset.Layout(gtx, func() {
					layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func() {
							pg.accountSelector.Draw(gtx)
						}),
						layout.Flexed(1, func() {
							layout.Stack{Alignment: layout.NE}.Layout(gtx,
								layout.Stacked(func() {
									layout.Align(layout.Center).Layout(gtx, func() {

									})
								}),
							)
						}),
					)
				})
			})
		},
		func() {
			layout.UniformInset(unit.Dp(5)).Layout(gtx, func() {
				pg.toLabel.Layout(gtx)

				inset := layout.Inset{
					Top: unit.Dp(25),
				}
				inset.Layout(gtx, func() {
					pg.destinationAddressInput.Layout(gtx, pg.destinationAddressEditor)
				})

				if pg.validationWidgets.destinationAddressLabel.Text != "" {
					inset = layout.Inset{
						Top: unit.Dp(43),
					}
					inset.Layout(gtx, func() {
						pg.validationWidgets.destinationAddressLabel.Layout(gtx)
					})
				}
			})
		},
		func() {
			layout.UniformInset(unit.Dp(5)).Layout(gtx, func() {
				layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func() {
						pg.amountLabel.Layout(gtx)

						if pg.validationWidgets.amountLabel.Text != "" {
							inset := layout.Inset{
								Top: unit.Dp(43),
							}
							inset.Layout(gtx, func() {
								pg.validationWidgets.amountLabel.Layout(gtx)
							})
						}
					}),
				)

				inset := layout.Inset{
					Top: unit.Dp(25),
				}
				inset.Layout(gtx, func() {
					pg.amountInput.Layout(gtx, pg.amountEditor)
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
			for pg.nextButton.Clicked(gtx) {
				if pg.validate(true) {
					pg.isShowingConfirmationModal = true
				}
			}
			nextBtn := pg.theme.Button("Next")
			if pg.validate(false) {
				nextBtn.Background = ui.LightBlueColor
			} else {
				nextBtn.Background = ui.GrayColor
			}
			nextBtn.Layout(gtx, pg.nextButton)
		},
	}

	pg.container.Layout(gtx, len(widgetFuncs), func(i int) {
		layout.UniformInset(unit.Dp(10)).Layout(gtx, widgetFuncs[i])
	})

	if pg.isShowingConfirmationModal {
		pg.drawConfirmationModal(gtx)
	}

	return nil
}

func (pg *Send) validate(setMessages bool) bool {
	isValid := true

	if pg.destinationAddressEditor.Text() == "" {
		if setMessages {
			pg.validationWidgets.destinationAddressLabel.Text = "Please enter a destination address"
		}
		isValid = false
	} else {
		// TODO check if destination address is a correct address for the current network
		pg.validationWidgets.destinationAddressLabel.Text = ""
	}

	amount := pg.amountEditor.Text()
	if amount == "" {
		if setMessages {
			pg.validationWidgets.amountLabel.Text = "Please enter an amount"
		}
		isValid = false
	} else {
		// check if amount is a valid number
		_, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			if setMessages {
				pg.validationWidgets.amountLabel.Text = "Please enter a valid amount"
			}
			isValid = false
		} else {
			pg.validationWidgets.amountLabel.Text = ""
		}
	}

	return isValid
}

func (pg *Send) drawConfirmationModal(gtx *layout.Context) {
	modalWidgetFuncs := []func(){
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
	})
}

package ui

import (
	"fmt"
	"image/color"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

type amountValue struct {
	inactiveTotalAmount         string
	activeTransactionFeeValue   string
	inactiveTransactionFeeValue string
	activeTotalCostValue        string
	inactiveTotalCostValue      string
}

type SendPage struct {
	pageContainer   layout.List
	theme           *decredmaterial.Theme
	txAuthor        *dcrlibwallet.TxAuthor
	broadcastResult *wallet.Broadcast

	wallet                 *wallet.Wallet
	selectedWallet         wallet.InfoShort
	selectedAccount        wallet.Account
	unspentOutputsSelected *map[int]map[int32]map[string]*wallet.UnspentOutput

	destinationAddressEditor     decredmaterial.Editor
	customChangeAddressEditor    decredmaterial.Editor
	sendAmountEditor             decredmaterial.Editor
	nextButton                   decredmaterial.Button
	closeConfirmationModalButton decredmaterial.Button
	confirmButton                decredmaterial.Button
	maxButton                    decredmaterial.Button
	sendToButton                 decredmaterial.Button

	confirmModal *decredmaterial.Modal

	currencySwap decredmaterial.IconButton

	txFeeCollapsible *decredmaterial.Collapsible
	txLine           *decredmaterial.Line

	remainingBalance int64
	amountAtoms      int64
	totalCostDCR     int64
	txFee            int64
	spendableBalance int64

	usdExchangeRate float64
	inputAmount     float64
	amountUSDtoDCR  float64
	amountDCRtoUSD  float64

	count              int
	defualtEditorWidth int
	nextEditorWidth    int

	amountErrorText    string
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
	line          *decredmaterial.Line

	isConfirmationModalOpen   bool
	isPasswordModalOpen       bool
	isBroadcastingTransaction bool
	shouldInitializeTxAuthor  bool

	txAuthorErrChan  chan error
	broadcastErrChan chan error

	borderColor color.NRGBA

	toggleCoinCtrl      *widget.Bool
	inputButtonCoinCtrl decredmaterial.Button
}

const (
	PageSend               = "Send"
	invalidPassphraseError = "error broadcasting transaction: " + dcrlibwallet.ErrInvalidPassphrase
)

func (win *Window) SendPage(common pageCommon) layout.Widget {
	pg := &SendPage{
		pageContainer: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},

		theme:                  common.theme,
		wallet:                 common.wallet,
		txAuthor:               &win.txAuthor,
		broadcastResult:        &win.broadcastResult,
		unspentOutputsSelected: &common.selectedUTXO,

		activeExchange:   "DCR",
		inactiveExchange: "USD",

		closeConfirmationModalButton: common.theme.Button(new(widget.Clickable), "Close"),
		nextButton:                   common.theme.Button(new(widget.Clickable), "Next"),
		confirmButton:                common.theme.Button(new(widget.Clickable), "Confirm"),
		maxButton:                    common.theme.Button(new(widget.Clickable), "MAX"),
		txFeeCollapsible:             common.theme.Collapsible(),
		txLine:                       common.theme.Line(),

		confirmModal:              common.theme.Modal(),
		isConfirmationModalOpen:   false,
		isPasswordModalOpen:       false,
		isBroadcastingTransaction: false,

		passwordModal:    common.theme.Password(),
		broadcastErrChan: make(chan error),
		txAuthorErrChan:  make(chan error),
		line:             common.theme.Line(),
	}
	pg.line.Color = common.theme.Color.Gray
	pg.line.Height = 2

	pg.borderColor = common.theme.Color.Hint

	pg.balanceAfterSendValue = "- DCR"

	pg.destinationAddressEditor = common.theme.Editor(new(widget.Editor), "Destination Address")
	pg.destinationAddressEditor.IsRequired = true
	pg.destinationAddressEditor.IsVisible = true
	pg.destinationAddressEditor.IsTitleLabel = false
	pg.destinationAddressEditor.Editor.SetText("")
	pg.destinationAddressEditor.Editor.SingleLine = true

	pg.customChangeAddressEditor = common.theme.Editor(new(widget.Editor), "Custom Change Address")
	pg.customChangeAddressEditor.IsVisible = true
	pg.customChangeAddressEditor.IsTitleLabel = false
	pg.customChangeAddressEditor.Editor.SetText("")
	pg.customChangeAddressEditor.Editor.SingleLine = true

	pg.sendAmountEditor = common.theme.Editor(new(widget.Editor), "Amount to be sent")
	pg.sendAmountEditor.SetRequiredErrorText("")
	pg.sendAmountEditor.IsRequired = true
	pg.sendAmountEditor.IsTitleLabel = false
	pg.sendAmountEditor.Bordered = false
	pg.sendAmountEditor.Editor.SingleLine = true
	pg.sendAmountEditor.Editor.SetText("0")
	pg.sendAmountEditor.TextSize = values.TextSize24

	pg.closeConfirmationModalButton.Background = common.theme.Color.Gray

	pg.currencySwap = common.theme.IconButton(new(widget.Clickable), common.icons.actionSwapVert)
	pg.currencySwap.Background = color.NRGBA{}
	pg.currencySwap.Color = common.theme.Color.Text
	pg.currencySwap.Inset = layout.UniformInset(values.MarginPadding0)
	pg.currencySwap.Size = values.MarginPadding30

	pg.maxButton.Background = common.theme.Color.Black
	pg.maxButton.Inset = layout.UniformInset(values.MarginPadding5)

	pg.sendToButton = common.theme.Button(new(widget.Clickable), "Send to account")
	pg.sendToButton.TextSize = values.TextSize14
	pg.sendToButton.Background = color.NRGBA{}
	pg.sendToButton.Color = common.theme.Color.Primary
	pg.sendToButton.Inset = layout.UniformInset(values.MarginPadding0)

	pg.toggleCoinCtrl = new(widget.Bool)
	pg.inputButtonCoinCtrl = common.theme.Button(new(widget.Clickable), "Inputs")
	pg.inputButtonCoinCtrl.Inset = layout.UniformInset(values.MarginPadding5)
	pg.inputButtonCoinCtrl.TextSize = values.MarginPadding10

	// defualtEditorWidth is the editor text size values.TextSize24
	pg.defualtEditorWidth = 24

	pg.txLine.Color = common.theme.Color.Gray

	go common.wallet.GetUSDExchangeValues(&pg)

	return func(gtx C) D {
		pg.Handle(common)
		return pg.Layout(gtx, common)
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

	if pg.toggleCoinCtrl.Value {
		_, spendableBalance := pg.calculateBalanceUTXO()
		pg.spendableBalance = spendableBalance
	} else {
		pg.spendableBalance = pg.selectedAccount.SpendableBalance
	}

	if pg.shouldInitializeTxAuthor {
		pg.shouldInitializeTxAuthor = false
		pg.sendAmountEditor.Editor.SetText("")
		pg.calculateErrorText = ""
		c.wallet.CreateTransaction(pg.selectedWallet.ID, pg.selectedAccount.Number, pg.txAuthorErrChan)
	}

	pg.validate(true)
	pg.watchForBroadcastResult(c)

	if pg.isBroadcastingTransaction {
		col := pg.theme.Color.Gray
		col.A = 150

		pg.nextButton.Text = "Sending..."
		pg.nextButton.Background = col
	} else {
		pg.nextButton.Text = "Next"
		pg.nextButton.Background = pg.theme.Color.Primary
	}

	for pg.nextButton.Button.Clicked() {
		if pg.validate(false) && pg.calculateErrorText == "" {
			pg.isConfirmationModalOpen = true
		}
	}

	for pg.confirmButton.Button.Clicked() {
		pg.isConfirmationModalOpen = false
		pg.isPasswordModalOpen = true
	}

	for pg.closeConfirmationModalButton.Button.Clicked() {
		pg.isConfirmationModalOpen = false
	}

	for pg.maxButton.Button.Clicked() {
		pg.activeExchange = "DCR"
		amountMax, err := pg.txAuthor.EstimateMaxSendAmount()
		if err != nil {
			return
		}
		pg.sendAmountEditor.Editor.SetText(fmt.Sprintf("%.10f", amountMax.DcrValue))
		pg.calculateValues()
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

	for range pg.destinationAddressEditor.Editor.Events() {
		go pg.calculateValues()
	}

	if pg.destinationAddressEditor.Editor.Len() == 0 || pg.sendAmountEditor.Editor.Len() == 0 {
		pg.balanceAfterSend(true)
	}

	for _, evt := range pg.sendAmountEditor.Editor.Events() {
		go pg.calculateValues()
		pg.handleEditorChange(evt)
	}

	if pg.sendAmountEditor.Editor.Focused() || pg.calculateErrorText != "" {
		if pg.calculateErrorText != "" {
			pg.borderColor = pg.theme.Color.Danger
		} else {
			pg.borderColor = pg.theme.Color.Primary
		}
	} else {
		pg.borderColor = pg.theme.Color.Hint
	}

	if pg.toggleCoinCtrl.Changed() && !pg.toggleCoinCtrl.Value {
		pg.txAuthor.UseInputs(nil)
	}

	if pg.inputButtonCoinCtrl.Button.Clicked() {
		c.wallet.AllUnspentOutputs(pg.selectedWallet.ID, pg.selectedAccount.Number)
		*c.page = PageUTXO
	}

	select {
	case err := <-pg.txAuthorErrChan:
		pg.calculateErrorText = err.Error()
		c.Notify(pg.calculateErrorText, false)
	case err := <-pg.broadcastErrChan:
		c.Notify(err.Error(), false)

		if err.Error() == invalidPassphraseError {
			time.AfterFunc(time.Second*3, func() {
				pg.isConfirmationModalOpen = true
			})
		}
		pg.isBroadcastingTransaction = false
	default:
	}
}

func (pg *SendPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return common.SelectedAccountLayout(gtx)
			})
		},
		func(gtx C) D {
			return pg.coinControlLayout(gtx, &common)
		},
		func(gtx C) D {
			return pg.destinationAddrSection(gtx)
		},
		func(gtx C) D {
			return pg.sendAmountSection(gtx)
		},
		func(gtx C) D {
			gtx.Constraints.Max.X = gtx.Px(values.MarginPadding450)
			return pg.drawTransactionDetailWidgets(gtx)
		},
		func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding450)
			return pg.nextButton.Layout(gtx)
		},
	}

	dims := layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return common.LayoutWithAccounts(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						if pg.pageContainer.Position.First > 0 {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							l := pg.theme.Line()
							l.Color = pg.theme.Color.Hint
							l.Width = gtx.Constraints.Min.X
							l.Height = 2
							return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
								return l.Layout(gtx)
							})
						}
						return layout.Dimensions{}
					}),
					layout.Rigid(func(gtx C) D {
						return pg.pageContainer.Layout(gtx, len(pageContent), func(gtx C, i int) D {
							p := values.MarginPadding10
							return layout.Inset{Left: p, Bottom: p, Right: p}.Layout(gtx, pageContent[i])
						})
					}),
				)
			})
		}),
	)

	if pg.isConfirmationModalOpen {
		return common.Modal(gtx, dims, pg.drawConfirmationModal(gtx))
	}

	if pg.isPasswordModalOpen {
		return common.Modal(gtx, dims, pg.drawPasswordModal(gtx))
	}

	return dims
}

func (pg *SendPage) drawTransactionDetailWidgets(gtx layout.Context) layout.Dimensions {
	w := []func(gtx C) D{
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

func (pg *SendPage) destinationAddrSection(gtx layout.Context) layout.Dimensions {
	return pg.centralize(gtx, func(gtx C) D {
		main := layout.UniformInset(values.MarginPadding20)
		return pg.sectionLayout(gtx, main, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.sendToAddressLayout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.destinationAddressEditor.Layout(gtx)
				}),
			)
		})
	})
}

func (pg *SendPage) sendAmountSection(gtx layout.Context) layout.Dimensions {
	return pg.centralize(gtx, func(gtx C) D {
		main := layout.UniformInset(values.MarginPadding20)
		return pg.sectionLayout(gtx, main, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.spendableBalanceLayout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.sectionBorder(gtx, values.MarginPadding10, func(gtx C) D {
						return pg.amountInputLayout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					txt := pg.theme.Body2(pg.amountErrorText)
					txt.Color = pg.theme.Color.Danger
					if pg.amountErrorText != "" {
						return txt.Layout(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return pg.txLine.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return pg.txFeeLayout(gtx)
				}),
			)
		})
	})
}

func (pg *SendPage) sendToAddressLayout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			amt := pg.theme.Body2("To")
			amt.Color = pg.theme.Color.Gray
			return amt.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
					return pg.sendToButton.Layout(gtx)
				})
			})
		}),
	)
}

func (pg *SendPage) spendableBalanceLayout(gtx layout.Context) layout.Dimensions {
	inset := layout.Inset{
		Bottom: values.MarginPadding10,
	}
	return inset.Layout(gtx, func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				amt := pg.theme.Body2("Amount")
				amt.Color = pg.theme.Color.Gray
				return amt.Layout(gtx)
			}),
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							title := pg.theme.Body2("Spendable Balance: ")
							title.Color = pg.theme.Color.Gray
							return title.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							sb := dcrutil.Amount(pg.spendableBalance).String()
							b := pg.theme.Body2(sb)
							b.Color = pg.theme.Color.Gray
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
	})
}

func (pg *SendPage) amountInputLayout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.W.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										w := pg.defualtEditorWidth
										if pg.nextEditorWidth != 0 {
											w = pg.nextEditorWidth
										}
										gtx.Constraints.Max.X = w
										return pg.sendAmountEditor.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										// this adjusts space between input values and currency symbol.
										m := values.MarginPadding5
										e := pg.sendAmountEditor.Editor.Len()
										if e > 0 {
											m = values.MarginPaddingMinus5
										}
										return layout.Inset{Left: m, Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
											return pg.theme.H6(pg.activeTotalAmount).Layout(gtx)
										})
									}),
									layout.Flexed(1, func(gtx C) D {
										return layout.E.Layout(gtx, func(gtx C) D {
											return pg.maxButton.Layout(gtx)
										})
									}),
								)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										m := values.MarginPadding10
										return layout.Inset{Left: m, Bottom: m}.Layout(gtx, func(gtx C) D {
											return pg.currencySwap.Layout(gtx)
										})
									}),
									layout.Rigid(func(gtx C) D {
										pg.line.Width = gtx.Constraints.Max.X
										return layout.Inset{Left: values.MarginPadding5, Top: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
											return pg.line.Layout(gtx)
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
		}),
	)
}

func (pg *SendPage) txFeeLayout(gtx layout.Context) layout.Dimensions {
	collapsibleHeader := func(gtx C) D {
		gtx.Constraints.Max.X = gtx.Px(values.MarginPadding390)
		return pg.tableLayout(gtx, pg.theme.Body2("Transaction Fee"), pg.activeTransactionFeeValue, pg.inactiveTransactionFeeValue)
	}

	collapsibleBody := func(gtx C) D {
		card := pg.theme.Card()
		card.Radius = decredmaterial.CornerRadius{
			NE: 0,
			NW: 0,
			SE: 0,
			SW: 0,
		}
		card.Color = pg.theme.Color.Background

		return card.Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				gtx.Constraints.Min.Y = 100
				return pg.theme.Body2("not implemented yet").Layout(gtx)
			})
		})
	}

	return pg.txFeeCollapsible.Layout(gtx, collapsibleHeader, collapsibleBody, nil)
}

func (pg *SendPage) drawConfirmationModal(gtx layout.Context) layout.Dimensions {
	w := []func(gtx C) D{
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
			return layout.Inset{Right: values.MarginPadding15, Bottom: values.MarginPaddingMinus10}.Layout(gtx, func(gtx C) D {
				return pg.txFeeLayout(gtx)
			})
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
			return layout.Center.Layout(gtx, func(gtx C) D {
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
			})
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

func (pg *SendPage) sectionBorder(gtx layout.Context, padding unit.Value, body layout.Widget) layout.Dimensions {
	border := widget.Border{Color: pg.borderColor, CornerRadius: values.MarginPadding5, Width: values.MarginPadding1}
	return border.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(padding).Layout(gtx, body)
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
			pg.sendAmountEditor.SetError("")
		}
		return false
	}

	if amount != "" {
		_, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			pg.sendAmountEditor.SetError("")
			return false
		}
	}

	return true
}

func (pg *SendPage) calculateValues() {
	defaultActiveValues := fmt.Sprintf("- %s", pg.activeExchange)
	defaultInactiveValues := fmt.Sprintf("(- %s)", pg.inactiveExchange)
	noExchangeText := "Exchange rate not fetched"
	pg.sendAmountEditor.Hint = "0"

	pg.activeTransactionFeeValue = defaultActiveValues
	pg.activeTotalCostValue = defaultActiveValues
	pg.inactiveTransactionFeeValue = defaultInactiveValues
	pg.inactiveTotalCostValue = defaultInactiveValues

	pg.calculateErrorText = ""
	pg.activeTotalAmount = pg.activeExchange
	pg.inactiveTotalAmount = fmt.Sprintf("0 %s", pg.inactiveExchange)

	// default values when exchange is not available
	if pg.LastTradeRate == "" {
		pg.activeTransactionFeeValue = defaultActiveValues
		pg.activeTotalCostValue = defaultActiveValues
		pg.inactiveTransactionFeeValue = ""
		pg.inactiveTotalCostValue = ""
		pg.activeTotalAmount = pg.activeExchange
		pg.inactiveTotalAmount = noExchangeText
	}

	if reflect.DeepEqual(pg.txAuthor, &dcrlibwallet.TxAuthor{}) || !pg.validate(true) {
		return
	}

	pg.inputAmount, _ = strconv.ParseFloat(pg.sendAmountEditor.Editor.Text(), 64)

	if pg.LastTradeRate != "" {
		pg.usdExchangeRate, _ = strconv.ParseFloat(pg.LastTradeRate, 64)
		pg.amountUSDtoDCR = pg.inputAmount / pg.usdExchangeRate
		pg.amountDCRtoUSD = pg.inputAmount * pg.usdExchangeRate
	}

	pg.setChangeDestinationAddr()
	if pg.activeExchange == "USD" && pg.LastTradeRate != "" {
		pg.amountAtoms = pg.setDestinationAddr(pg.amountUSDtoDCR)
	} else {
		pg.amountAtoms = pg.setDestinationAddr(pg.inputAmount)
	}

	if pg.amountAtoms == 0 {
		return
	}

	pg.txFee = pg.getTxFee(pg.toggleCoinCtrl.Value)
	if pg.txFee == 0 {
		return
	}

	pg.totalCostDCR = pg.txFee + pg.amountAtoms

	pg.updateDefaultValues()
	pg.balanceAfterSend(false)
}

func (pg *SendPage) setDestinationAddr(sendAmount float64) int64 {
	pg.amountErrorText = ""
	amount, err := dcrutil.NewAmount(sendAmount)
	if err != nil {
		pg.feeEstimationError(err.Error(), "amount")
		return 0
	}

	pg.amountAtoms = int64(amount)
	pg.txAuthor.RemoveSendDestination(0)
	pg.txAuthor.AddSendDestination(pg.destinationAddressEditor.Editor.Text(), pg.amountAtoms, false)
	return pg.amountAtoms
}

func (pg *SendPage) setChangeDestinationAddr() {
	if pg.customChangeAddressEditor.Editor.Len() > 0 {
		pg.txAuthor.RemoveChangeDestination()
		pg.txAuthor.SetChangeDestination(pg.customChangeAddressEditor.Editor.Text())
	}
}

func (pg *SendPage) amountValues() amountValue {
	txFeeValueUSD := dcrutil.Amount(pg.txFee).ToCoin() * pg.usdExchangeRate
	switch {
	case pg.activeExchange == "USD" && pg.LastTradeRate != "":
		return amountValue{
			inactiveTotalAmount:         dcrutil.Amount(pg.amountAtoms).String(),
			activeTransactionFeeValue:   fmt.Sprintf("%f USD", txFeeValueUSD),
			inactiveTransactionFeeValue: fmt.Sprintf("(%s)", dcrutil.Amount(pg.txFee).String()),
			activeTotalCostValue:        fmt.Sprintf("%s USD", strconv.FormatFloat(pg.inputAmount+txFeeValueUSD, 'f', 7, 64)),
			inactiveTotalCostValue:      fmt.Sprintf("(%s )", dcrutil.Amount(pg.totalCostDCR).String()),
		}
	case pg.activeExchange == "DCR" && pg.LastTradeRate != "":
		return amountValue{
			inactiveTotalAmount:         fmt.Sprintf("%s USD", strconv.FormatFloat(pg.amountDCRtoUSD, 'f', 2, 64)),
			activeTransactionFeeValue:   dcrutil.Amount(pg.txFee).String(),
			inactiveTransactionFeeValue: fmt.Sprintf("(%f USD)", txFeeValueUSD),
			activeTotalCostValue:        dcrutil.Amount(pg.totalCostDCR).String(),
			inactiveTotalCostValue:      fmt.Sprintf("(%s USD)", strconv.FormatFloat(pg.amountDCRtoUSD+txFeeValueUSD, 'f', 7, 64)),
		}
	default:
		return amountValue{
			inactiveTotalAmount:       "Exchange rate not fetched",
			activeTransactionFeeValue: dcrutil.Amount(pg.txFee).String(),
			activeTotalCostValue:      dcrutil.Amount(pg.totalCostDCR).String(),
		}
	}
}

func (pg *SendPage) updateDefaultValues() {
	v := pg.amountValues()
	pg.activeTotalAmount = pg.activeExchange
	pg.inactiveTotalAmount = v.inactiveTotalAmount
	pg.activeTransactionFeeValue = v.activeTransactionFeeValue
	pg.inactiveTransactionFeeValue = v.inactiveTransactionFeeValue
	pg.activeTotalCostValue = v.activeTotalCostValue
	pg.inactiveTotalCostValue = v.inactiveTotalCostValue
}

func (pg *SendPage) getTxFee(isCustomInputs bool) int64 {
	// calculate transaction fee
	pg.amountErrorText = ""
	if isCustomInputs {
		utxoKeys, _ := pg.calculateBalanceUTXO()
		err := pg.txAuthor.UseInputs(utxoKeys)
		if err != nil {
			log.Error(err)
			return 0
		}
	}

	feeAndSize, err := pg.txAuthor.EstimateFeeAndSize()
	if err != nil {
		pg.feeEstimationError(err.Error(), "fee")
		return 0
	}

	return feeAndSize.Fee.AtomValue
}

func (pg *SendPage) calculateBalanceUTXO() ([]string, int64) {
	utxos := (*pg.unspentOutputsSelected)[pg.selectedWallet.ID][pg.selectedAccount.Number]
	var utxoKeys []string
	var totalAmount int64
	for utxoKey, utxo := range utxos {
		utxoKeys = append(utxoKeys, utxoKey)
		totalAmount += utxo.UTXO.Amount
	}
	return utxoKeys, totalAmount
}

func (pg *SendPage) balanceAfterSend(isInputAmountEmpty bool) {
	pg.remainingBalance = 0
	if isInputAmountEmpty {
		pg.remainingBalance = pg.selectedAccount.SpendableBalance
	} else {
		pg.remainingBalance = pg.selectedAccount.SpendableBalance - pg.totalCostDCR
	}
	pg.balanceAfterSendValue = dcrutil.Amount(pg.remainingBalance).String()
}

func (pg *SendPage) feeEstimationError(err, errorPath string) {
	if err == "insufficient_balance" {
		pg.amountErrorText = "Not enough funds"
	}
	if strings.Contains(err, "invalid amount") {
		pg.amountErrorText = "Invalid amount"
	}
	pg.calculateErrorText = fmt.Sprintf("error estimating transaction %s: %s", errorPath, err)
}

func (pg *SendPage) watchForBroadcastResult(c pageCommon) {
	if pg.broadcastResult == nil {
		return
	}

	if pg.broadcastResult.TxHash != "" {
		if pg.remainingBalance != -1 {
			pg.spendableBalance = pg.remainingBalance
		}
		pg.remainingBalance = -1
		c.Notify("Transaction Sent", true)

		pg.destinationAddressEditor.Editor.SetText("")
		pg.sendAmountEditor.Editor.SetText("")
		pg.isConfirmationModalOpen = false
		pg.isBroadcastingTransaction = false
		pg.broadcastResult.TxHash = ""
		pg.calculateValues()
		(*pg.unspentOutputsSelected)[pg.selectedWallet.ID][pg.selectedAccount.Number] = make(map[string]*wallet.UnspentOutput)
	}
}

// handleEditorChange handles changes on the editor and adjust its width of the send amount input field
// it also updates the DCR - USD exchange rate value
func (pg *SendPage) handleEditorChange(evt widget.EditorEvent) {
	editorTextLength := pg.sendAmountEditor.Editor.Len()

	// calculateNextWidth use the values of the defualtEditorWidth(the editor text size) and
	// total number of text in the editor to determine the width of the amount field
	calculateNextWidth := func() {
		editorTextLength = editorTextLength + 1
		pg.nextEditorWidth = pg.defualtEditorWidth * editorTextLength
	}

	switch evt.(type) {
	case widget.ChangeEvent:
		calculateNextWidth()
		go pg.wallet.GetUSDExchangeValues(&pg)
	}
}

// drawlayout wraps the pg tx and sync section in a card layout
func (pg *SendPage) sectionLayout(gtx layout.Context, inset layout.Inset, body layout.Widget) layout.Dimensions {
	gtx.Constraints.Max.X = gtx.Px(values.MarginPadding450)
	gtx.Constraints.Min.X = gtx.Px(values.MarginPadding450)
	return pg.theme.Card().Layout(gtx, func(gtx C) D {
		return inset.Layout(gtx, body)
	})
}

func (pg *SendPage) centralize(gtx layout.Context, content layout.Widget) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx C) D {
			return layout.Center.Layout(gtx, content)
		}),
	)
}

func (pg *SendPage) coinControlLayout(gtx layout.Context, c *pageCommon) layout.Dimensions {
	main := layout.UniformInset(values.MarginPadding20)
	return pg.sectionLayout(gtx, main, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.theme.Switch(pg.toggleCoinCtrl).Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Left: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
							return pg.theme.Body1("Coin control features").Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				if !pg.toggleCoinCtrl.Value {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Top: values.MarginPadding10,
						}.Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx C) D { return pg.inputButtonCoinCtrl.Layout(gtx) }),
								layout.Rigid(func(gtx C) D {
									utxos := c.selectedUTXO[pg.selectedWallet.ID][pg.selectedAccount.Number]
									var totalAmount int64
									for _, utxo := range utxos {
										totalAmount += utxo.UTXO.Amount
									}
									txt := "Automatically selected"
									if len(utxos) > 0 {
										txt = fmt.Sprintf("Selected: %d | Amount: %s", len(utxos), dcrutil.Amount(totalAmount).String())
									}
									return layout.Inset{Left: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
										return pg.theme.Body1(txt).Layout(gtx)
									})
								}),
							)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Top: values.MarginPadding10,
						}.Layout(gtx, func(gtx C) D {
							return pg.customChangeAddressEditor.Layout(gtx)
						})
					}),
				)
			}),
		)
	})
}

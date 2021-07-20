package send

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

const (
	SendPageID = "Send"
)

type SendPage struct {
	*load.Load
	pageContainer layout.List

	sourceAccountSelector *components.AccountSelector
	sendDestination       *destination
	amount                *sendAmount

	backButton   decredmaterial.IconButton
	infoButton   decredmaterial.IconButton
	moreOption   decredmaterial.IconButton
	nextButton   decredmaterial.Button
	sendToButton decredmaterial.Button
	clearAllBtn  decredmaterial.Button

	txFeeCollapsible *decredmaterial.Collapsible

	moreOptionIsOpen bool

	exchangeRate float64

	*authoredTxData
}

type authoredTxData struct {
	txAuthor            *dcrlibwallet.TxAuthor
	destinationAddress  string
	destinationAccount  *dcrlibwallet.Account
	sourceAccount       *dcrlibwallet.Account
	txFee               string
	txFeeUSD            string
	estSignedSize       string
	totalCost           string
	totalCostUSD        string
	balanceAfterSend    string
	balanceAfterSendUSD string
	sendAmount          string
	sendAmountUSD       string
}

func NewSendPage(l *load.Load) *SendPage {
	pg := &SendPage{
		Load: l,
		pageContainer: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},
		sendDestination: newSendDestination(l),
		amount:          newSendAmount(l),

		clearAllBtn:      l.Theme.Button(new(widget.Clickable), "Clear all fields"),
		txFeeCollapsible: l.Theme.Collapsible(),

		exchangeRate: -1,

		authoredTxData: &authoredTxData{},
	}

	pg.nextButton = l.Theme.Button(new(widget.Clickable), "Next")
	pg.nextButton.Background = l.Theme.Color.InactiveGray

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(pg.Load)
	pg.backButton.Icon = pg.Icons.ContentClear

	pg.moreOption = l.Theme.PlainIconButton(new(widget.Clickable), pg.Icons.NavMoreIcon)
	pg.moreOption.Color = l.Theme.Color.Gray3
	pg.moreOption.Inset = layout.UniformInset(values.MarginPadding0)

	pg.sendToButton = l.Theme.Button(new(widget.Clickable), "Send to account")
	pg.sendToButton.TextSize = values.TextSize14
	pg.sendToButton.Background = color.NRGBA{}
	pg.sendToButton.Color = l.Theme.Color.Primary
	pg.sendToButton.Inset = layout.UniformInset(values.MarginPadding0)

	pg.clearAllBtn.Background = l.Theme.Color.Surface
	pg.clearAllBtn.Color = l.Theme.Color.Text
	pg.clearAllBtn.Inset = layout.UniformInset(values.MarginPadding15)

	// Source account picker
	pg.sourceAccountSelector = components.NewAccountSelector(l).
		Title("Sending account").
		AccountSelected(func(selectedAccount *dcrlibwallet.Account) {
			pg.validateAndConstructTx()
		}).
		AccountValidator(func(account *dcrlibwallet.Account) bool {
			wal := pg.Load.WL.MultiWallet.WalletWithID(account.WalletID)

			// Imported and watch only wallet accounts are invalid for sending
			accountIsValid := account.Number != load.MaxInt32 && !wal.IsWatchingOnlyWallet()

			if wal.ReadBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, false) {
				// privacy is enabled for selected wallet

				if pg.sendDestination.sendToAddress {
					// only mixed can send to address
					accountIsValid = account.Number == wal.MixedAccountNumber()
				} else {
					// send to account, check if selected destination account belongs to wallet
					destinationAccount := pg.sendDestination.destinationAccountSelector.SelectedAccount()
					if destinationAccount.WalletID != account.WalletID {
						accountIsValid = account.Number == wal.MixedAccountNumber()
					}
				}
			}
			return accountIsValid
		})

	pg.sendDestination.destinationAccountSelector.AccountSelected(func(selectedAccount *dcrlibwallet.Account) {
		pg.validateAndConstructTx()
		pg.sourceAccountSelector.SelectFirstWalletValidAccount() // refresh source account
	})

	pg.sendDestination.addressChanged = func() {
		pg.validateAndConstructTx()
	}

	pg.amount.amountChanged = func() {
		pg.validateAndConstructTx()
	}

	return pg
}

func (pg *SendPage) OnResume() {
	pg.sendDestination.destinationAccountSelector.SelectFirstWalletValidAccount()
	pg.sourceAccountSelector.SelectFirstWalletValidAccount()

	currencyExchangeValue := pg.WL.MultiWallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	if currencyExchangeValue == components.USDExchangeValue {
		pg.fetchExchangeValue()
	}
}

func (pg *SendPage) fetchExchangeValue() {
	go func() {
		var dcrUsdtBittrex load.DCRUSDTBittrex
		err := load.GetUSDExchangeValue(&dcrUsdtBittrex)
		if err != nil {
			// TODO: handle exchange error
			return
		}

		exchangeRate, err := strconv.ParseFloat(dcrUsdtBittrex.LastTradeRate, 64)
		if err != nil {
			// TODO: handle exchange error
			return
		}

		pg.exchangeRate = exchangeRate
		pg.amount.setExchangeRate(exchangeRate)
		pg.validateAndConstructTx() // convert estimates to usd
	}()
}

func (pg *SendPage) Layout(gtx layout.Context) layout.Dimensions {
	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.pageSections(gtx, "From", false, func(gtx C) D {
				return pg.sourceAccountSelector.Layout(gtx)
			})
		},
		func(gtx C) D {
			return pg.toSection(gtx)
		},
		func(gtx C) D {
			return pg.feeSection(gtx)
		},
	}

	dims := layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return layout.Stack{Alignment: layout.NE}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return components.UniformPadding(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
									return pg.topNav(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
									return pg.pageContainer.Layout(gtx, len(pageContent), func(gtx C, i int) D {
										return layout.Inset{Bottom: values.MarginPadding4, Top: values.MarginPadding4}.Layout(gtx, pageContent[i])
									})
								})
							}),
						)
					})
				}),
				layout.Stacked(func(gtx C) D {
					if pg.moreOptionIsOpen {
						inset := layout.Inset{
							Top:   values.MarginPadding40,
							Right: values.MarginPadding20,
						}
						return inset.Layout(gtx, func(gtx C) D {
							border := widget.Border{Color: pg.Theme.Color.Background, CornerRadius: values.MarginPadding5, Width: values.MarginPadding1}
							return border.Layout(gtx, pg.clearAllBtn.Layout)
						})
					}
					return layout.Dimensions{}
				}),
			)
		}),
		layout.Stacked(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.S.Layout(gtx, func(gtx C) D {
				return layout.Inset{Left: values.MarginPadding1}.Layout(gtx, func(gtx C) D {
					return pg.balanceSection(gtx)
				})
			})
		}),
	)

	return dims
}

func (pg *SendPage) topNav(gtx layout.Context) layout.Dimensions {
	m := values.MarginPadding20
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.backButton.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: m}.Layout(gtx, pg.Theme.H6("Send DCR").Layout)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(pg.infoButton.Layout),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Left: m}.Layout(gtx, pg.moreOption.Layout)
					}),
				)
			})
		}),
	)
}

func (pg *SendPage) toSection(gtx layout.Context) layout.Dimensions {

	return pg.pageSections(gtx, "To", true, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding16,
				}.Layout(gtx, func(gtx C) D {
					if !pg.sendDestination.sendToAddress {
						return pg.sendDestination.destinationAccountSelector.Layout(gtx)
					}
					return pg.sendDestination.destinationAddressEditor.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				if pg.exchangeRate != -1 {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(0.45, func(gtx C) D {
							return pg.amount.dcrAmountEditor.Layout(gtx)
						}),
						layout.Flexed(0.1, func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									icon := pg.Icons.CurrencySwapIcon
									icon.Scale = 0.45
									return icon.Layout(gtx)
								})
							})
						}),
						layout.Flexed(0.45, func(gtx C) D {
							return pg.amount.usdAmountEditor.Layout(gtx)
						}),
					)
				}
				return pg.amount.dcrAmountEditor.Layout(gtx)
			}),
		)
	})
}

func (pg *SendPage) feeSection(gtx layout.Context) layout.Dimensions {
	collapsibleHeader := func(gtx C) D {
		feeText := pg.txFee
		if pg.exchangeRate != -1 {
			feeText = fmt.Sprintf("%s (%s)", pg.txFee, pg.txFeeUSD)
		}
		return pg.Theme.Body1(feeText).Layout(gtx)
	}

	collapsibleBody := func(gtx C) D {
		card := pg.Theme.Card()
		card.Color = pg.Theme.Color.LightGray
		inset := layout.Inset{
			Top: values.MarginPadding10,
		}
		return inset.Layout(gtx, func(gtx C) D {
			return card.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							//TODO
							return pg.contentRow(gtx, "Estimated time", "10 minutes (2 blocks)")
						}),
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Top:    values.MarginPadding5,
								Bottom: values.MarginPadding5,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return pg.contentRow(gtx, "Estimated size", pg.estSignedSize)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return pg.contentRow(gtx, "Fee rate", "10 atoms/Byte")
						}),
					)
				})
			})
		})
	}
	inset := layout.Inset{
		Bottom: values.MarginPadding75,
	}
	return inset.Layout(gtx, func(gtx C) D {
		return pg.pageSections(gtx, "Fee", false, func(gtx C) D {
			return pg.txFeeCollapsible.Layout(gtx, collapsibleHeader, collapsibleBody)
		})
	})
}

func (pg *SendPage) balanceSection(gtx layout.Context) layout.Dimensions {
	c := pg.Theme.Card()
	c.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
	return c.Layout(gtx, func(gtx C) D {
		return components.UniformPadding(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(0.6, func(gtx C) D {
					inset := layout.Inset{
						Right: values.MarginPadding15,
					}
					return inset.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								inset := layout.Inset{
									Bottom: values.MarginPadding10,
								}
								return inset.Layout(gtx, func(gtx C) D {
									totalCostText := pg.totalCost
									if pg.exchangeRate != -1 {
										totalCostText = fmt.Sprintf("%s (%s)", pg.totalCost, pg.totalCostUSD)
									}
									return pg.contentRow(gtx, "Total cost", totalCostText)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return pg.contentRow(gtx, "Balance after send", pg.balanceAfterSend)
							}),
						)
					})
				}),
				layout.Flexed(0.3, func(gtx C) D {
					pg.nextButton.Inset = layout.Inset{Top: values.MarginPadding15, Bottom: values.MarginPadding15}
					return pg.nextButton.Layout(gtx)
				}),
			)
		})
	})
}

func (pg *SendPage) pageSections(gtx layout.Context, title string, showAccountSwitch bool, body layout.Widget) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Bottom: values.MarginPadding16,
							}
							return inset.Layout(gtx, pg.Theme.Body1(title).Layout)
						}),
						layout.Flexed(1, func(gtx C) D {
							if showAccountSwitch {
								return layout.E.Layout(gtx, func(gtx C) D {
									inset := layout.Inset{
										Top: values.MarginPaddingMinus5,
									}
									return inset.Layout(gtx, pg.sendDestination.accountSwitch.Layout)
								})
							}
							return layout.Dimensions{}
						}),
					)
				}),
				layout.Rigid(body),
			)
		})
	})
}

func (pg *SendPage) contentRow(gtx layout.Context, leftValue, rightValue string) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			txt := pg.Theme.Body2(leftValue)
			txt.Color = pg.Theme.Color.Gray
			return txt.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(pg.Theme.Body1(rightValue).Layout),
					layout.Rigid(func(gtx C) D {
						return layout.Dimensions{}
					}),
				)
			})
		}),
	)
}

func (pg *SendPage) validateAndConstructTx() {
	if pg.validate() {
		pg.constructTx()
	} else {
		pg.clearEstimates()
	}
}

func (pg *SendPage) validate() bool {

	amountIsValid := pg.amount.amountIsValid()
	addressIsValid := pg.sendDestination.validate()

	validForSending := amountIsValid && addressIsValid
	if validForSending {
		pg.nextButton.Background = pg.Theme.Color.Primary
	} else {
		pg.nextButton.Background = pg.Theme.Color.Hint
	}

	return validForSending
}

func (pg *SendPage) constructTx() {
	destinationAddress, err := pg.sendDestination.destinationAddress()
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}
	destinationAccount := pg.sendDestination.destinationAccount()

	amountAtom, sendMax, err := pg.amount.validAmount()
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}

	sourceAccount := pg.sourceAccountSelector.SelectedAccount()
	unsignedTx, err := pg.WL.MultiWallet.NewUnsignedTx(sourceAccount.WalletID, sourceAccount.Number)
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}

	err = unsignedTx.AddSendDestination(destinationAddress, amountAtom, sendMax)
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}

	feeAndSize, err := unsignedTx.EstimateFeeAndSize()
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}

	feeAtom := feeAndSize.Fee.AtomValue
	if sendMax {
		amountAtom = sourceAccount.Balance.Spendable - feeAtom
	}

	totalSendingAmount := dcrutil.Amount(amountAtom + feeAtom)
	balanceAfterSend := dcrutil.Amount(sourceAccount.Balance.Spendable - int64(totalSendingAmount))

	// populate display data
	pg.txFee = dcrutil.Amount(feeAtom).String()
	pg.estSignedSize = fmt.Sprintf("%d bytes", feeAndSize.EstimatedSignedSize)
	pg.totalCost = totalSendingAmount.String()
	pg.balanceAfterSend = balanceAfterSend.String()
	pg.sendAmount = dcrutil.Amount(amountAtom).String()
	pg.destinationAddress = destinationAddress
	pg.destinationAccount = destinationAccount
	pg.sourceAccount = sourceAccount

	if sendMax {
		// TODO: this workaround ignores the change events from the
		// amount input to avoid construct tx cycle.
		pg.amount.setAmount(amountAtom)
	}

	if pg.exchangeRate != -1 {
		pg.txFeeUSD = fmt.Sprintf("$%.4f", load.DCRToUSD(pg.exchangeRate, feeAndSize.Fee.DcrValue))
		pg.totalCostUSD = load.FormatUSDBalance(pg.Printer, load.DCRToUSD(pg.exchangeRate, totalSendingAmount.ToCoin()))
		pg.balanceAfterSendUSD = load.FormatUSDBalance(pg.Printer, load.DCRToUSD(pg.exchangeRate, balanceAfterSend.ToCoin()))

		usdAmount := load.DCRToUSD(pg.exchangeRate, dcrutil.Amount(amountAtom).ToCoin())
		pg.sendAmountUSD = load.FormatUSDBalance(pg.Printer, usdAmount)
	}

	pg.txAuthor = unsignedTx
}

func (pg *SendPage) feeEstimationError(err string) {
	if err == dcrlibwallet.ErrInsufficientBalance {
		pg.amount.setError("Not enough funds")
	} else if strings.Contains(err, invalidAmountErr) {
		pg.amount.setError(invalidAmountErr)
	} else {
		pg.amount.setError(err)
		pg.CreateToast("Error estimating transaction: "+err, false)
	}

	pg.clearEstimates()
}

func (pg *SendPage) clearEstimates() {
	pg.txAuthor = nil
	pg.txFee = " - "
	pg.txFeeUSD = " - "
	pg.estSignedSize = " - "
	pg.totalCost = " - "
	pg.totalCostUSD = " - "
	pg.balanceAfterSend = " - "
	pg.balanceAfterSendUSD = " - "
	pg.sendAmount = " - "
	pg.sendAmountUSD = " - "
}

func (pg *SendPage) resetFields() {
	pg.sendDestination.clearAddressInput()

	pg.amount.resetFields()
}

func (pg *SendPage) Handle() {

	pg.sendDestination.handle()
	pg.amount.handle()

	if pg.backButton.Button.Clicked() {
		pg.ChangePage(*pg.ReturnPage)
	}

	if pg.infoButton.Button.Clicked() {
		info := modal.NewInfoModal(pg.Load).
			Title("Send DCR").
			Body("Input or scan the destination wallet address and input the amount to send funds.").
			PositiveButton("Got it", func() {})
		pg.ShowModal(info)
	}

	for pg.moreOption.Button.Clicked() {
		pg.moreOptionIsOpen = !pg.moreOptionIsOpen
	}

	for pg.nextButton.Button.Clicked() {
		if pg.txAuthor != nil {
			confirmTxModal := newSendConfirmModal(pg.Load, pg.authoredTxData)
			confirmTxModal.exchangeRateSet = pg.exchangeRate != -1

			confirmTxModal.txSent = func() {
				pg.resetFields()
				pg.clearEstimates()
			}

			confirmTxModal.Show()
		}
	}

	for pg.clearAllBtn.Clicked() {
		pg.moreOptionIsOpen = true

		pg.sendDestination.clearAddressInput()

		pg.amount.clearAmount()
	}

}

func (pg *SendPage) OnClose() {

}

package page

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
	"github.com/planetdecred/godcr/ui/values"
)

const (
	SendPageID             = "Send"
	invalidPassphraseError = "error broadcasting transaction: " + dcrlibwallet.ErrInvalidPassphrase
)

type sendPage struct {
	*load.Load
	pageContainer layout.List

	sourceAccountSelector      *AccountSelector
	destinationAccountSelector *AccountSelector

	destinationAddressEditor decredmaterial.Editor
	dcrAmountEditor          decredmaterial.Editor
	usdAmountEditor          decredmaterial.Editor

	backButton   decredmaterial.IconButton
	infoButton   decredmaterial.IconButton
	moreOption   decredmaterial.IconButton
	nextButton   decredmaterial.Button
	maxButton    decredmaterial.Button
	sendToButton decredmaterial.Button
	clearAllBtn  decredmaterial.Button

	accountSwitch    *decredmaterial.SwitchButtonText
	txFeeCollapsible *decredmaterial.Collapsible

	moreOptionIsOpen bool
	sendToAddress    bool
	sendMax          bool

	amountErrorText string
	exchangeRate    float64

	*authoredTxData
}

type authoredTxData struct {
	txAuthor            *dcrlibwallet.TxAuthor
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

func NewSendPage(l *load.Load) *sendPage {
	pg := &sendPage{
		Load: l,
		pageContainer: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},

		maxButton:        l.Theme.Button(new(widget.Clickable), "MAX"),
		clearAllBtn:      l.Theme.Button(new(widget.Clickable), "Clear all fields"),
		txFeeCollapsible: l.Theme.Collapsible(),

		exchangeRate: -1,

		authoredTxData: &authoredTxData{},
	}

	pg.accountSwitch = l.Theme.SwitchButtonText([]decredmaterial.SwitchItem{{Text: "Address"}, {Text: "My account"}})

	pg.nextButton = l.Theme.Button(new(widget.Clickable), "Next")
	pg.nextButton.Background = l.Theme.Color.InactiveGray

	pg.dcrAmountEditor = l.Theme.Editor(new(widget.Editor), "Amount (DCR)")
	pg.dcrAmountEditor.Editor.SetText("")
	pg.dcrAmountEditor.IsCustomButton = true
	pg.dcrAmountEditor.Editor.SingleLine = true
	pg.dcrAmountEditor.CustomButton.Background = l.Theme.Color.Gray
	pg.dcrAmountEditor.CustomButton.Inset = layout.UniformInset(values.MarginPadding2)
	pg.dcrAmountEditor.CustomButton.Text = "Max"
	pg.dcrAmountEditor.CustomButton.CornerRadius = values.MarginPadding0

	pg.usdAmountEditor = l.Theme.Editor(new(widget.Editor), "Amount (USD)")
	pg.usdAmountEditor.Editor.SetText("")
	pg.usdAmountEditor.IsCustomButton = true
	pg.usdAmountEditor.Editor.SingleLine = true
	pg.usdAmountEditor.CustomButton.Background = l.Theme.Color.Gray
	pg.usdAmountEditor.CustomButton.Inset = layout.UniformInset(values.MarginPadding2)
	pg.usdAmountEditor.CustomButton.Text = "Max"
	pg.usdAmountEditor.CustomButton.CornerRadius = values.MarginPadding0

	pg.destinationAddressEditor = l.Theme.Editor(new(widget.Editor), "Address")
	pg.destinationAddressEditor.Editor.SingleLine = true
	pg.destinationAddressEditor.Editor.SetText("")

	pg.backButton, pg.infoButton = subpageHeaderButtons(pg.Load)
	pg.backButton.Icon = pg.Icons.ContentClear

	pg.moreOption = l.Theme.PlainIconButton(new(widget.Clickable), pg.Icons.NavMoreIcon)
	pg.moreOption.Color = l.Theme.Color.Gray3
	pg.moreOption.Inset = layout.UniformInset(values.MarginPadding0)

	pg.maxButton.Background = l.Theme.Color.Gray3
	pg.maxButton.Inset = layout.UniformInset(values.MarginPadding5)

	pg.sendToButton = l.Theme.Button(new(widget.Clickable), "Send to account")
	pg.sendToButton.TextSize = values.TextSize14
	pg.sendToButton.Background = color.NRGBA{}
	pg.sendToButton.Color = l.Theme.Color.Primary
	pg.sendToButton.Inset = layout.UniformInset(values.MarginPadding0)

	pg.clearAllBtn.Background = l.Theme.Color.Surface
	pg.clearAllBtn.Color = l.Theme.Color.Text
	pg.clearAllBtn.Inset = layout.UniformInset(values.MarginPadding15)

	// Source account picker
	pg.sourceAccountSelector = NewAccountSelector(l).
		title("Sending account").
		accountSelected(func(selectedAccount *dcrlibwallet.Account) {
			pg.validateAndConstructTx()
		}).
		accountValidator(func(account *dcrlibwallet.Account) bool {
			wal := pg.Load.WL.MultiWallet.WalletWithID(account.WalletID)

			// Imported and watch only wallet accounts are invalid for sending
			accountIsValid := account.Number != MaxInt32 && !wal.IsWatchingOnlyWallet()

			if wal.ReadBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, false) {
				// privacy is enabled for selected wallet

				if pg.sendToAddress {
					// only mixed can send to address
					accountIsValid = account.Number == wal.MixedAccountNumber()
				} else {
					// send to account, check if selected destination account belongs to wallet
					destinationAccount := pg.destinationAccountSelector.selectedAccount
					if destinationAccount.WalletID != account.WalletID {
						accountIsValid = account.Number == wal.MixedAccountNumber()
					}
				}
			}
			return accountIsValid
		})

	// Destination account picker
	pg.destinationAccountSelector = NewAccountSelector(pg.Load).
		title("Receiving account").
		accountSelected(func(selectedAccount *dcrlibwallet.Account) {
			pg.validateAndConstructTx()

			pg.sourceAccountSelector.selectFirstWalletValidAccount() // refresh source account
		}).
		accountValidator(func(account *dcrlibwallet.Account) bool {

			// Filter out imported account and mixed.
			wal := pg.Load.WL.MultiWallet.WalletWithID(account.WalletID)
			if account.Number == MaxInt32 ||
				account.Number == wal.MixedAccountNumber() {
				return false
			}

			return true
		})

	return pg
}

func (pg *sendPage) OnResume() {
	pg.destinationAccountSelector.selectFirstWalletValidAccount()
	pg.sourceAccountSelector.selectFirstWalletValidAccount()

	currencyExchangeValue := pg.WL.MultiWallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	if currencyExchangeValue == USDExchangeValue {
		pg.fetchExchangeValue()
	}
}

func (pg *sendPage) fetchExchangeValue() {
	go func() {
		var dcrUsdtBittrex load.DCRUSDTBittrex
		err := load.GetUSDExchangeValue(&dcrUsdtBittrex)
		if err != nil {
			pg.usdAmountEditor.SetError(err.Error())
			return
		}

		exchangeRate, err := strconv.ParseFloat(dcrUsdtBittrex.LastTradeRate, 64)
		if err != nil {
			pg.usdAmountEditor.SetError(err.Error())
			return
		}

		pg.exchangeRate = exchangeRate
		pg.usdAmountEditor.SetError("")
	}()
}

func (pg *sendPage) Layout(gtx layout.Context) layout.Dimensions {
	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.pageSections(gtx, "From", func(gtx C) D {
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
					return uniformPadding(gtx, func(gtx C) D {
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

func (pg *sendPage) topNav(gtx layout.Context) layout.Dimensions {
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

func (pg *sendPage) toSection(gtx layout.Context) layout.Dimensions {
	pg.dcrAmountEditor.SetError(pg.amountErrorText)

	if pg.amountErrorText != "" {
		pg.dcrAmountEditor.LineColor, pg.dcrAmountEditor.TitleLabelColor = pg.Theme.Color.Danger, pg.Theme.Color.Danger
		pg.usdAmountEditor.LineColor, pg.usdAmountEditor.TitleLabelColor = pg.Theme.Color.Danger, pg.Theme.Color.Danger
	} else {
		pg.dcrAmountEditor.LineColor, pg.dcrAmountEditor.TitleLabelColor = pg.Theme.Color.Gray1, pg.Theme.Color.Gray3
		pg.usdAmountEditor.LineColor, pg.usdAmountEditor.TitleLabelColor = pg.Theme.Color.Gray1, pg.Theme.Color.Gray3
	}

	return pg.pageSections(gtx, "To", func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding16,
				}.Layout(gtx, func(gtx C) D {
					if !pg.sendToAddress {
						return pg.destinationAccountSelector.Layout(gtx)
					}
					return pg.destinationAddressEditor.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				if pg.exchangeRate != -1 {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(0.45, func(gtx C) D {
							return pg.dcrAmountEditor.Layout(gtx)
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
							return pg.usdAmountEditor.Layout(gtx)
						}),
					)
				}
				return pg.dcrAmountEditor.Layout(gtx)
			}),
		)
	})
}

func (pg *sendPage) feeSection(gtx layout.Context) layout.Dimensions {
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
		return pg.pageSections(gtx, "Fee", func(gtx C) D {
			return pg.txFeeCollapsible.Layout(gtx, collapsibleHeader, collapsibleBody)
		})
	})
}

func (pg *sendPage) balanceSection(gtx layout.Context) layout.Dimensions {
	c := pg.Theme.Card()
	c.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
	return c.Layout(gtx, func(gtx C) D {
		return uniformPadding(gtx, func(gtx C) D {
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

func (pg *sendPage) pageSections(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
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
							if title == "To" { //TODO
								return layout.E.Layout(gtx, func(gtx C) D {
									inset := layout.Inset{
										Top: values.MarginPaddingMinus5,
									}
									return inset.Layout(gtx, pg.accountSwitch.Layout)
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

func (pg *sendPage) contentRow(gtx layout.Context, leftValue, rightValue string) layout.Dimensions {
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

func (pg *sendPage) inputsNotEmpty(editors ...*widget.Editor) bool {
	for _, e := range editors {
		if e.Text() == "" {
			return false
		}
	}
	return true
}

func (pg *sendPage) resetErrorText() {
	pg.amountErrorText = ""
	pg.destinationAddressEditor.SetError("")
	pg.dcrAmountEditor.SetError("")
	pg.usdAmountEditor.SetError("")
}

func (pg *sendPage) validateAndConstructTx() {
	if pg.validate() {
		pg.constructTx()
	} else {
		pg.clearEstimates()
	}
}

func (pg *sendPage) validate() bool {

	_, err := strconv.ParseFloat(pg.dcrAmountEditor.Editor.Text(), 64)
	amountIsValid := err == nil
	addressIsValid := true // default send to account
	if pg.sendToAddress {
		addressIsValid, _ = pg.validateDestinationAddress()
	}

	validForSending := (amountIsValid || pg.sendMax) && addressIsValid
	if validForSending {
		pg.amountErrorText = ""
		pg.destinationAddressEditor.SetError("")
		pg.nextButton.Background = pg.Theme.Color.Primary
	} else {
		pg.nextButton.Background = pg.Theme.Color.Hint
	}

	return validForSending
}

func (pg *sendPage) destinationAddress() (string, error) {
	if pg.sendToAddress {
		valid, address := pg.validateDestinationAddress()
		if valid {
			return address, nil
		}

		return "", fmt.Errorf("invalid address")
	}

	destinationAccount := pg.destinationAccountSelector.selectedAccount
	wal := pg.WL.MultiWallet.WalletWithID(destinationAccount.WalletID)

	return wal.CurrentAddress(destinationAccount.Number)
}

func (pg *sendPage) validateDestinationAddress() (bool, string) {

	address := pg.destinationAddressEditor.Editor.Text()
	address = strings.TrimSpace(address)

	if len(address) == 0 {
		pg.destinationAddressEditor.SetError("")
		return false, address
	}

	if pg.WL.MultiWallet.IsAddressValid(address) {
		pg.destinationAddressEditor.SetError("")
		return true, address
	}

	pg.destinationAddressEditor.SetError("Invalid address")
	return false, address
}

func (pg *sendPage) validateDCRAmount() bool {
	pg.amountErrorText = ""
	if pg.inputsNotEmpty(pg.dcrAmountEditor.Editor) {
		dcrAmount, err := strconv.ParseFloat(pg.dcrAmountEditor.Editor.Text(), 64)
		if err != nil {
			// empty usd input
			pg.usdAmountEditor.Editor.SetText("")
			pg.amountErrorText = "Invalid amount"
			// todo: invalid decimal places error
			return false
		}

		if pg.exchangeRate != -1 {
			usdAmount := dcrTOUSD(pg.exchangeRate, dcrAmount)
			pg.usdAmountEditor.Editor.SetText(fmt.Sprintf("%.2f", usdAmount)) // 2 decimal places
		}

		return true
	}

	// empty usd input since this is empty
	pg.usdAmountEditor.Editor.SetText("")
	return false
}

// validateUSDAmount is called when usd text changes
func (pg *sendPage) validateUSDAmount() bool {
	pg.amountErrorText = ""
	if pg.inputsNotEmpty(pg.usdAmountEditor.Editor) {
		usdAmount, err := strconv.ParseFloat(pg.usdAmountEditor.Editor.Text(), 64)
		if err != nil {
			// empty dcr input
			pg.dcrAmountEditor.Editor.SetText("")
			pg.amountErrorText = "Invalid amount"
			return false
		}

		if pg.exchangeRate != -1 { //TODO usd amount should not be visible.
			dcrAmount := usdToDCR(pg.exchangeRate, usdAmount)
			pg.dcrAmountEditor.Editor.SetText(fmt.Sprintf("%.8f", dcrAmount)) // 8 decimal places
		}

		return true
	}

	// empty dcr input since this is empty
	pg.dcrAmountEditor.Editor.SetText("")
	return false
}

func (pg *sendPage) constructTx() {
	address, err := pg.destinationAddress()
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}

	amountAtom := int64(0)

	if !pg.sendMax {
		amount, err := strconv.ParseFloat(pg.dcrAmountEditor.Editor.Text(), 64)
		if err != nil {
			pg.feeEstimationError(err.Error())
			return
		}
		amountAtom = dcrlibwallet.AmountAtom(amount)
	}

	sourceAccount := pg.sourceAccountSelector.selectedAccount
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}

	unsignedTx, err := pg.WL.MultiWallet.NewUnsignedTx(sourceAccount.WalletID, sourceAccount.Number)
	if err != nil {
		pg.feeEstimationError(err.Error())
		return
	}

	err = unsignedTx.AddSendDestination(address, amountAtom, pg.sendMax)
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
	if pg.sendMax {
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

	if pg.sendMax {
		pg.dcrAmountEditor.Editor.SetText(fmt.Sprintf("%.8f", dcrutil.Amount(amountAtom).ToCoin()))
	}

	if pg.exchangeRate != -1 {
		pg.txFeeUSD = fmt.Sprintf("$%.4f", dcrTOUSD(pg.exchangeRate, feeAndSize.Fee.DcrValue))
		pg.totalCostUSD = formatUSDBalance(pg.Printer, dcrTOUSD(pg.exchangeRate, totalSendingAmount.ToCoin()))
		pg.balanceAfterSendUSD = formatUSDBalance(pg.Printer, dcrTOUSD(pg.exchangeRate, balanceAfterSend.ToCoin()))
		pg.sendAmountUSD = formatUSDBalance(pg.Printer, dcrTOUSD(pg.exchangeRate, dcrutil.Amount(amountAtom).ToCoin()))

		if pg.sendMax {
			pg.usdAmountEditor.Editor.SetText(formatUSDBalance(pg.Printer, dcrutil.Amount(amountAtom).ToCoin()))
		}
	}

	pg.txAuthor = unsignedTx
}

func (pg *sendPage) feeEstimationError(err string) {
	if err == dcrlibwallet.ErrInsufficientBalance {
		pg.amountErrorText = "Not enough funds"
	} else if strings.Contains(err, "invalid amount") {
		pg.amountErrorText = "Invalid amount"
	} else {
		pg.amountErrorText = err
		pg.CreateToast("Error estimating transaction: "+err, false)
	}

	pg.clearEstimates()
}

func (pg *sendPage) clearEstimates() {
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

func (pg *sendPage) resetFields() {
	pg.destinationAddressEditor.SetError("")
	pg.destinationAddressEditor.Editor.SetText("")

	pg.amountErrorText = ""
	pg.dcrAmountEditor.Editor.SetText("")
	pg.usdAmountEditor.Editor.SetText("")
}

func (pg *sendPage) Handle() {
	sendToAddress := pg.accountSwitch.SelectedIndex() == 1
	if sendToAddress != pg.sendToAddress { // switch changed
		pg.sendToAddress = sendToAddress
		pg.validateAndConstructTx()
	}

	if pg.backButton.Button.Clicked() {
		pg.resetErrorText()
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

	for _, evt := range pg.destinationAddressEditor.Editor.Events() {
		if pg.destinationAddressEditor.Editor.Focused() {
			switch evt.(type) {
			case widget.ChangeEvent:
				pg.sendMax = false
				pg.validateAndConstructTx()
			}
		}
	}

	for _, evt := range pg.dcrAmountEditor.Editor.Events() {
		if pg.dcrAmountEditor.Editor.Focused() {
			switch evt.(type) {
			case widget.ChangeEvent:
				pg.sendMax = false
				pg.validateDCRAmount()
				pg.validateAndConstructTx()

			}
		}
	}

	for _, evt := range pg.usdAmountEditor.Editor.Events() {
		if pg.usdAmountEditor.Editor.Focused() {
			switch evt.(type) {
			case widget.ChangeEvent:
				pg.sendMax = false
				pg.validateUSDAmount()
				pg.validateAndConstructTx()
			}
		}
	}

	for pg.nextButton.Button.Clicked() {
		if pg.validate() {
			confirmTxModal := newSendConfirmModal(pg.Load, pg.authoredTxData)
			confirmTxModal.exchangeRateSet = pg.exchangeRate != -1
			confirmTxModal.sourceAccount = pg.sourceAccountSelector.selectedAccount
			if sendToAddress {
				confirmTxModal.destinationAddress = pg.destinationAddressEditor.Editor.Text()
			} else {
				confirmTxModal.destinationAccount = pg.destinationAccountSelector.selectedAccount
			}

			confirmTxModal.txSent = func() {
				pg.resetFields()
				pg.clearEstimates()
			}

			confirmTxModal.Show()
		}
	}

	for pg.clearAllBtn.Button.Clicked() {
		pg.moreOptionIsOpen = true

		pg.destinationAddressEditor.SetError("")
		pg.destinationAddressEditor.Editor.SetText("")

		pg.amountErrorText = ""
		pg.dcrAmountEditor.Editor.SetText("")
		pg.usdAmountEditor.Editor.SetText("")
	}
}

func (pg *sendPage) OnClose() {

}

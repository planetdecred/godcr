package ui

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const (
	PageSend               = "Send"
	invalidPassphraseError = "error broadcasting transaction: " + dcrlibwallet.ErrInvalidPassphrase
)

type sendPage struct {
	pageContainer layout.List
	common        pageCommon
	theme         *decredmaterial.Theme

	// txAuthor        *dcrlibwallet.TxAuthor
	broadcastResult *wallet.Broadcast
	wallet          *wallet.Wallet

	destinationAddressEditor decredmaterial.Editor
	dcrAmountEditor          decredmaterial.Editor
	usdAmountEditor          decredmaterial.Editor
	passwordEditor           decredmaterial.Editor

	moreOption decredmaterial.IconButton

	nextButton                   decredmaterial.Button
	closeConfirmationModalButton decredmaterial.Button
	confirmButton                decredmaterial.Button
	maxButton                    decredmaterial.Button
	sendToButton                 decredmaterial.Button
	clearAllBtn                  decredmaterial.Button

	accountSwitch    *decredmaterial.SwitchButtonText
	confirmModal     *decredmaterial.Modal
	txFeeCollapsible *decredmaterial.Collapsible

	remainingBalance int64
	totalCostDCR     int64
	txFee            int64
	spendableBalance int64

	usdExchangeRate float64
	inputAmount     float64
	amountUSDtoDCR  float64
	amountDCRtoUSD  float64

	txFeeSize          string
	txFeeEstimatedTime string

	leftTransactionFeeValue  string
	rightTransactionFeeValue string

	leftTotalCostValue  string
	rightTotalCostValue string

	sendAmountDCR string
	sendAmountUSD string

	balanceAfterSendValue string
	activeTotalAmount     string

	exchangeRate       string
	exchangeErr        string
	sendToAddress      bool
	noExchangeErrMsg   string
	amountErrorText    string
	calculateErrorText string

	isConfirmationModalOpen   bool
	isBroadcastingTransaction bool
	isMoreOption              bool

	shouldInitializeTxAuthor bool
	usdExchangeSet           bool

	txAuthorErrChan  chan error
	broadcastErrChan chan error

	walletSelected int
}

func (win *Window) SendPage(common pageCommon) Page {
	pg := &sendPage{
		pageContainer: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},

		common: common,
		theme:  common.theme,
		wallet: common.wallet,
		// txAuthor:        &win.txAuthor,
		broadcastResult: &win.broadcastResult,

		noExchangeErrMsg:             "Exchange rate not fetched",
		closeConfirmationModalButton: common.theme.Button(new(widget.Clickable), "Cancel"),
		confirmButton:                common.theme.Button(new(widget.Clickable), ""),
		maxButton:                    common.theme.Button(new(widget.Clickable), "MAX"),
		clearAllBtn:                  common.theme.Button(new(widget.Clickable), "Clear all fields"),
		txFeeCollapsible:             common.theme.Collapsible(),

		confirmModal:              common.theme.Modal(),
		isConfirmationModalOpen:   false,
		isBroadcastingTransaction: false,

		broadcastErrChan: make(chan error),
		txAuthorErrChan:  make(chan error),
	}

	pg.walletSelected = common.wallAcctSelector.selectedSendWallet

	pg.accountSwitch = common.theme.SwitchButtonText([]decredmaterial.SwitchItem{{Text: "Address"}, {Text: "My account"}})

	pg.balanceAfterSendValue = "- DCR"

	pg.nextButton = common.theme.Button(new(widget.Clickable), "Next")
	pg.nextButton.Background = pg.theme.Color.InactiveGray

	pg.dcrAmountEditor = common.theme.Editor(new(widget.Editor), "Amount (DCR)")
	pg.dcrAmountEditor.Editor.SetText("")
	pg.dcrAmountEditor.IsCustomButton = true
	pg.dcrAmountEditor.Editor.SingleLine = true
	pg.dcrAmountEditor.CustomButton.Background = common.theme.Color.Gray
	pg.dcrAmountEditor.CustomButton.Inset = layout.UniformInset(values.MarginPadding2)
	pg.dcrAmountEditor.CustomButton.Text = "Max"
	pg.dcrAmountEditor.CustomButton.CornerRadius = values.MarginPadding0

	pg.usdAmountEditor = common.theme.Editor(new(widget.Editor), "Amount (USD)")
	pg.usdAmountEditor.Editor.SetText("")
	pg.usdAmountEditor.IsCustomButton = true
	pg.usdAmountEditor.Editor.SingleLine = true
	pg.usdAmountEditor.CustomButton.Background = common.theme.Color.Gray
	pg.usdAmountEditor.CustomButton.Inset = layout.UniformInset(values.MarginPadding2)
	pg.usdAmountEditor.CustomButton.Text = "Max"
	pg.usdAmountEditor.CustomButton.CornerRadius = values.MarginPadding0

	pg.passwordEditor = win.theme.EditorPassword(new(widget.Editor), "Spending password")
	pg.passwordEditor.Editor.SetText("")
	pg.passwordEditor.Editor.SingleLine = true
	pg.passwordEditor.Editor.Submit = true

	pg.destinationAddressEditor = common.theme.Editor(new(widget.Editor), "Address")
	pg.destinationAddressEditor.Editor.SingleLine = true
	pg.destinationAddressEditor.Editor.SetText("")

	pg.closeConfirmationModalButton.Background = color.NRGBA{}
	pg.closeConfirmationModalButton.Color = common.theme.Color.Primary

	pg.moreOption = common.theme.PlainIconButton(new(widget.Clickable), common.icons.navMoreIcon)
	pg.moreOption.Color = common.theme.Color.Gray3
	pg.moreOption.Inset = layout.UniformInset(values.MarginPadding0)

	pg.maxButton.Background = common.theme.Color.Gray3
	pg.maxButton.Inset = layout.UniformInset(values.MarginPadding5)

	pg.sendToButton = common.theme.Button(new(widget.Clickable), "Send to account")
	pg.sendToButton.TextSize = values.TextSize14
	pg.sendToButton.Background = color.NRGBA{}
	pg.sendToButton.Color = common.theme.Color.Primary
	pg.sendToButton.Inset = layout.UniformInset(values.MarginPadding0)

	pg.clearAllBtn.Background = common.theme.Color.Surface
	pg.clearAllBtn.Color = common.theme.Color.Text
	pg.clearAllBtn.Inset = layout.UniformInset(values.MarginPadding15)

	pg.fetchExchangeValue()

	return pg
}

func (pg *sendPage) pageID() string {
	return PageSend
}

func (pg *sendPage) Layout(gtx layout.Context) layout.Dimensions {
	common := pg.common
	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.pageSections(gtx, "From", func(gtx C) D {
				return common.accountSelectorLayout(gtx, "send", pg.sendToAddress)
			})
		},
		func(gtx C) D {
			return pg.toSection(gtx, common)
		},
		func(gtx C) D {
			return pg.feeSection(gtx)
		},
	}

	dims := common.Layout(gtx, func(gtx C) D {
		return layout.Stack{Alignment: layout.S}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return layout.Stack{Alignment: layout.NE}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return common.UniformPadding(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
										return pg.topNav(gtx, common)
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
						if pg.isMoreOption {
							inset := layout.Inset{
								Top:   values.MarginPadding40,
								Right: values.MarginPadding20,
							}
							return inset.Layout(gtx, func(gtx C) D {
								border := widget.Border{Color: pg.theme.Color.Background, CornerRadius: values.MarginPadding5, Width: values.MarginPadding1}
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
						return pg.balanceSection(gtx, common)
					})
				})
			}),
		)
	})

	if pg.isConfirmationModalOpen {
		return common.Modal(gtx, dims, pg.confirmationModal(gtx, common))
	}

	return dims
}

func (pg *sendPage) topNav(gtx layout.Context, common pageCommon) layout.Dimensions {
	m := values.MarginPadding20
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					common.subPageBackButton.Icon = common.icons.contentClear
					return common.subPageBackButton.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: m}.Layout(gtx, pg.theme.H6("Send DCR").Layout)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(common.subPageInfoButton.Layout),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Left: m}.Layout(gtx, pg.moreOption.Layout)
					}),
				)
			})
		}),
	)
}

func (pg *sendPage) toSection(gtx layout.Context, common pageCommon) layout.Dimensions {
	return pg.pageSections(gtx, "To", func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding16,
				}.Layout(gtx, func(gtx C) D {
					if !pg.sendToAddress {
						return common.accountSelectorLayout(gtx, "receive", pg.sendToAddress)
					}
					return pg.destinationAddressEditor.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				if pg.usdExchangeSet {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(0.45, func(gtx C) D {
							return pg.dcrAmountEditor.Layout(gtx)
						}),
						layout.Flexed(0.1, func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									return decredmaterial.Clickable(gtx, new(widget.Clickable), func(gtx C) D {
										//TODO nil clickable
										icon := common.icons.currencySwapIcon
										icon.Scale = 0.45
										return icon.Layout(gtx)
									})
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
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.theme.Body1(pg.leftTransactionFeeValue).Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				b := pg.theme.Body1(pg.rightTransactionFeeValue)
				b.Color = pg.theme.Color.Gray
				inset := layout.Inset{
					Left: values.MarginPadding5,
				}
				if pg.usdExchangeSet {
					return inset.Layout(gtx, b.Layout)
				}
				return layout.Dimensions{}
			}),
		)
	}

	collapsibleBody := func(gtx C) D {
		card := pg.theme.Card()
		card.Color = pg.theme.Color.LightGray
		inset := layout.Inset{
			Top: values.MarginPadding10,
		}
		return inset.Layout(gtx, func(gtx C) D {
			return card.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return pg.contentRow(gtx, "Estimated time", "-", "")
						}),
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Top:    values.MarginPadding5,
								Bottom: values.MarginPadding5,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return pg.contentRow(gtx, "Estimated size", pg.txFeeSize, "")
							})
						}),
						layout.Rigid(func(gtx C) D {
							return pg.contentRow(gtx, "Fee rate", "10 atoms/Byte", "")
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

func (pg *sendPage) balanceSection(gtx layout.Context, common pageCommon) layout.Dimensions {
	c := pg.theme.Card()
	c.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
	return c.Layout(gtx, func(gtx C) D {
		return common.UniformPadding(gtx, func(gtx C) D {
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
									if pg.usdExchangeSet {
										return pg.contentRow(gtx, "Total cost", pg.leftTotalCostValue+" "+pg.rightTotalCostValue, "")
									}
									return pg.contentRow(gtx, "Total cost", pg.leftTotalCostValue, "")
								})
							}),
							layout.Rigid(func(gtx C) D {
								return pg.contentRow(gtx, "Balance after send", pg.balanceAfterSendValue, "")
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

func (pg *sendPage) confirmationModal(gtx layout.Context, common pageCommon) layout.Dimensions {
	receiveWallet := common.wallet.AllWallets()[common.wallAcctSelector.selectedReceiveWallet] // TODO
	receiveAcct, _ := receiveWallet.GetAccount(int32(common.wallAcctSelector.selectedReceiveAccount))
	sendWallet := common.wallet.AllWallets()[common.wallAcctSelector.selectedSendWallet] // TODO
	sendAcct, _ := sendWallet.GetAccount(int32(common.wallAcctSelector.selectedSendAccount))
	w := []layout.Widget{
		func(gtx C) D {
			return pg.theme.H6("Confim to send").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					icon := common.icons.sendIcon
					icon.Scale = 0.7
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Top: values.MarginPadding2, Right: values.MarginPadding16}.Layout(gtx, icon.Layout)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return common.layoutBalance(gtx, pg.sendAmountDCR, true)
								}),
								layout.Flexed(1, func(gtx C) D {
									if pg.usdExchangeSet {
										return layout.E.Layout(gtx, func(gtx C) D {
											txt := pg.theme.Body1(pg.sendAmountUSD)
											txt.Color = pg.theme.Color.Gray
											return txt.Layout(gtx)
										})
									}
									return layout.Dimensions{}
								}),
							)
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							icon := common.icons.navigationArrowForward
							icon.Color = pg.theme.Color.Gray3
							return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
								return icon.Layout(gtx, values.MarginPadding15)
							})
						}),
						layout.Rigid(func(gtx C) D {
							if !pg.sendToAddress {
								return layout.E.Layout(gtx, func(gtx C) D {
									return layout.Flex{}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return pg.theme.Body2(receiveAcct.Name).Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											card := pg.theme.Card()
											card.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
											card.Color = pg.theme.Color.LightGray
											inset := layout.Inset{
												Left: values.MarginPadding5,
											}
											return inset.Layout(gtx, func(gtx C) D {
												return card.Layout(gtx, func(gtx C) D {
													return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
														txt := pg.theme.Caption(receiveWallet.Name)
														txt.Color = pg.theme.Color.Gray
														return txt.Layout(gtx)
													})
												})
											})
										}),
									)
								})
							}
							return pg.theme.Body2(pg.destinationAddressEditor.Editor.Text()).Layout(gtx)
						}),
					)
				}),
			)
		},
		func(gtx C) D {
			return pg.theme.Separator().Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.contentRow(gtx, "Sending from", sendAcct.Name, sendWallet.Name)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding8, Bottom: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						if pg.usdExchangeSet {
							return pg.contentRow(gtx, "Fee", pg.leftTransactionFeeValue+" "+pg.rightTransactionFeeValue, "")
						}
						return pg.contentRow(gtx, "Fee", pg.leftTransactionFeeValue, "")
					})
				}),
				layout.Rigid(func(gtx C) D {
					if pg.usdExchangeSet {
						return pg.contentRow(gtx, "Total cost", pg.leftTotalCostValue+" "+pg.rightTotalCostValue, "")
					}
					return pg.contentRow(gtx, "Total cost", pg.leftTotalCostValue, "")
				}),
			)
		},
		func(gtx C) D {
			return pg.passwordEditor.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					icon := common.icons.actionInfo
					icon.Color = pg.theme.Color.Gray
					return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return icon.Layout(gtx, values.MarginPadding20)
					})
				}),
				layout.Rigid(func(gtx C) D {
					txt := pg.theme.Body2("Your DCR will be sent after this step.")
					txt.Color = pg.theme.Color.Gray3
					return txt.Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						inset := layout.Inset{
							Left: values.MarginPadding5,
						}
						return inset.Layout(gtx, func(gtx C) D {
							return pg.closeConfirmationModalButton.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						if common.modalLoad.loading {
							th := material.NewTheme(gofont.Collection())
							return layout.Inset{Top: unit.Dp(7)}.Layout(gtx, func(gtx C) D {
								return material.Loader(th).Layout(gtx)
							})
						}
						pg.confirmButton.Text = fmt.Sprintf("Send %s", dcrutil.Amount(pg.totalCostDCR).String())
						return pg.confirmButton.Layout(gtx)
					}),
				)
			})
		},
	}

	return pg.confirmModal.Layout(gtx, w, 900)
}

func (pg *sendPage) pageSections(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
	return pg.theme.Card().Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Bottom: values.MarginPadding16,
							}
							return inset.Layout(gtx, pg.theme.Body1(title).Layout)
						}),
						layout.Flexed(1, func(gtx C) D {
							if title == "To" {
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

func (pg *sendPage) contentRow(gtx layout.Context, leftValue, rightValue, walletName string) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			txt := pg.theme.Body2(leftValue)
			txt.Color = pg.theme.Color.Gray
			return txt.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(pg.theme.Body1(rightValue).Layout),
					layout.Rigid(func(gtx C) D {
						if walletName != "" {
							card := pg.theme.Card()
							card.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
							card.Color = pg.theme.Color.LightGray
							inset := layout.Inset{
								Left: values.MarginPadding5,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return card.Layout(gtx, func(gtx C) D {
									return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
										txt := pg.theme.Caption(walletName)
										txt.Color = pg.theme.Color.Gray
										return txt.Layout(gtx)
									})
								})
							})
						}
						return layout.Dimensions{}
					}),
				)
			})
		}),
	)
}

func (pg *sendPage) validate(c pageCommon) bool {
	if pg.sendToAddress {
		isAmountValid := pg.validateLeftAmount()
		if pg.usdAmountEditor.Editor.Focused() {
			isAmountValid = pg.validateRightAmount()
		}

		if pg.usdExchangeSet && !isAmountValid {
			if pg.usdAmountEditor.Editor.Focused() {
				pg.dcrAmountEditor.Editor.SetText("")
			} else {
				pg.usdAmountEditor.Editor.SetText("")
			}
		}

		if !pg.validateDestinationAddress(c) || !isAmountValid || pg.calculateErrorText != "" {
			pg.nextButton.Background = pg.theme.Color.Hint
			return false
		}
	}

	pg.nextButton.Background = pg.theme.Color.Primary
	return true
}

func (pg *sendPage) validateDestinationAddress(c pageCommon) bool {
	if pg.inputsNotEmpty(pg.destinationAddressEditor.Editor) {
		isValid, _ := pg.wallet.IsAddressValid(pg.destinationAddressEditor.Editor.Text())
		if !isValid {
			pg.destinationAddressEditor.SetError("Invalid address")
			return false
		}

		pg.destinationAddressEditor.SetError("")
		return true
	}

	pg.balanceAfterSend(true, c)
	pg.destinationAddressEditor.SetError("Input address")
	return false
}

func (pg *sendPage) validateLeftAmount() bool {
	if pg.inputsNotEmpty(pg.dcrAmountEditor.Editor) {
		_, err := strconv.ParseFloat(pg.dcrAmountEditor.Editor.Text(), 64)
		if err != nil {
			pg.dcrAmountEditor.SetError("Invalid amount")
			return false
		}
		pg.dcrAmountEditor.SetError("")
		return true
	}
	pg.dcrAmountEditor.SetError("")
	return false
}

func (pg *sendPage) validateAmount() bool {
	if pg.inputsNotEmpty(pg.dcrAmountEditor.Editor) {
		_, err := strconv.ParseFloat(pg.dcrAmountEditor.Editor.Text(), 64)
		if err != nil {
			pg.dcrAmountEditor.SetError("Invalid amount")
			return false
		}
		pg.dcrAmountEditor.SetError("")
		return true
	}

	pg.dcrAmountEditor.SetError("")
	return false
}

func (pg *sendPage) validateRightAmount() bool {
	if pg.inputsNotEmpty(pg.usdAmountEditor.Editor) {
		_, err := strconv.ParseFloat(pg.usdAmountEditor.Editor.Text(), 64)
		if err != nil {
			pg.usdAmountEditor.SetError("Invalid amount")
			return false
		}
		pg.usdAmountEditor.SetError("")
		return true
	}
	pg.usdAmountEditor.SetError("")
	return false
}

func (pg *sendPage) inputsNotEmpty(editors ...*widget.Editor) bool {
	for _, e := range editors {
		if e.Text() == "" {
			return false
		}
	}
	return true
}

func (pg *sendPage) updateExchangeError() {
	pg.usdAmountEditor.SetError("")
	if pg.exchangeRate == "" && pg.usdExchangeSet {
		println("Updating exhange rate eror")
		pg.usdAmountEditor.SetError(pg.noExchangeErrMsg)
	}
}

func (pg *sendPage) setDestinationAddr(sendAmount float64, common pageCommon) {
	// receiveWallet := common.info.Wallets[common.wallAcctSelector.selectedReceiveWallet]
	// receiveAcct := receiveWallet.Accounts[common.wallAcctSelector.selectedReceiveAccount]

	pg.amountErrorText = ""
	// amount, err := dcrutil.NewAmount(sendAmount)
	// if err != nil {
	// 	pg.feeEstimationError(err.Error(), "amount")
	// 	return
	// }

	// pg.txAuthor.RemoveSendDestination(0)
	// addr := pg.destinationAddressEditor.Editor.Text()
	// if pg.sendToOption == "My account" {
	// 	addr = receiveAcct.CurrentAddress
	// }
	// pg.txAuthor.AddSendDestination(addr, pg.amountAtoms, false)
	//TODO
}

func (pg *sendPage) balanceAfterSend(isInputAmountEmpty bool, c pageCommon) {
	sendWallet := c.wallet.AllWallets()[c.wallAcctSelector.selectedSendWallet] // TODO
	sendAcct, _ := sendWallet.GetAccount(int32(c.wallAcctSelector.selectedSendAccount))

	pg.remainingBalance = 0
	if isInputAmountEmpty {
		pg.remainingBalance = sendAcct.Balance.Spendable
	} else {
		pg.remainingBalance = sendAcct.Balance.Spendable - pg.totalCostDCR
	}
	pg.balanceAfterSendValue = dcrutil.Amount(pg.remainingBalance).String()
}

func (pg *sendPage) feeEstimationError(err, errorPath string) {
	if err == dcrlibwallet.ErrInsufficientBalance {
		pg.amountErrorText = "Not enough funds"
		return
	}
	if strings.Contains(err, "invalid amount") {
		pg.amountErrorText = "Invalid amount"
		return
	}
	pg.calculateErrorText = fmt.Sprintf("error estimating transaction %s: %s", errorPath, err)
}

func (pg *sendPage) watchForBroadcastResult(c pageCommon) {
	if pg.broadcastResult == nil {
		return
	}

	if pg.broadcastResult.TxHash != "" {
		c.popPage() // confirm TODO
		c.notify("1 Transaction Sent", true)

		if pg.remainingBalance != -1 {
			pg.spendableBalance = pg.remainingBalance
		}
		pg.remainingBalance = -1

		pg.isConfirmationModalOpen = false
		pg.isBroadcastingTransaction = false
		pg.resetFields()
		c.modalLoad.setLoading(false)
		pg.broadcastResult.TxHash = ""
		// pg.calculateValues(c, true)
		pg.destinationAddressEditor.Editor.SetText("")
	}
}

func (pg *sendPage) handleEditorChange(evt widget.EditorEvent, c pageCommon) {
	switch evt.(type) {
	case widget.ChangeEvent:
		pg.fetchExchangeValue()
	case widget.SubmitEvent:
		pg.sendFund(c)
	}
}

func (pg *sendPage) resetFields() {
	pg.destinationAddressEditor.SetError("")
	pg.dcrAmountEditor.Editor.SetText("")
	pg.usdAmountEditor.Editor.SetText("")
	pg.passwordEditor.Editor.SetText("")
	pg.leftTotalCostValue = ""
	pg.rightTotalCostValue = ""
}

func (pg *sendPage) resetErrorText() {
	pg.amountErrorText = ""
	pg.calculateErrorText = ""
	pg.destinationAddressEditor.SetError("")
	pg.dcrAmountEditor.SetError("")
	pg.usdAmountEditor.SetError("")
	pg.passwordEditor.SetError("")
}

func (pg *sendPage) fetchExchangeValue() {
	go func() {
		err := pg.wallet.GetUSDExchangeValues(&pg)
		if err != nil {
			pg.updateExchangeError()
		}
	}()
}

func (pg *sendPage) sendFund(c pageCommon) {
	if !pg.inputsNotEmpty(pg.passwordEditor.Editor) {
		return
	}
	c.modalLoad.setLoading(true)
	pg.isBroadcastingTransaction = true
	// TODO
	// pg.wallet.BroadcastTransaction(pg.txAuthor, []byte(pg.passwordEditor.Editor.Text()), pg.broadcastErrChan)
}

func (pg *sendPage) handle() {
	c := pg.common
	sendWallet := c.wallet.AllWallets()[c.wallAcctSelector.selectedSendWallet] // TODO
	sendAcct, _ := sendWallet.GetAccount(int32(c.wallAcctSelector.selectedSendAccount))

	if pg.exchangeErr != "" {
		c.notify(pg.exchangeErr, false)
	}

	pg.sendToAddress = pg.accountSwitch.SelectedIndex() == 1

	if c.subPageBackButton.Button.Clicked() {
		pg.resetErrorText()
		pg.resetFields()
		c.popPage()
	}

	if c.subPageInfoButton.Button.Clicked() {
		go func() {
			c.modalReceiver <- &modalLoad{
				template:   SendInfoTemplate,
				title:      "Send DCR",
				cancel:     c.closeModal,
				cancelText: "Got it",
			}
		}()
	}

	for pg.moreOption.Button.Clicked() {
		pg.isMoreOption = !pg.isMoreOption
	}

	currencyExchangeValue := pg.wallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	pg.usdExchangeSet = false
	if currencyExchangeValue == USDExchangeValue {
		pg.usdExchangeSet = true
	}

	for range pg.destinationAddressEditor.Editor.Events() {
		// pg.calculateValues(c, true)
		// construct tx
	}

	for _, evt := range pg.dcrAmountEditor.Editor.Events() {
		if pg.dcrAmountEditor.Editor.Focused() {
			pg.handleEditorChange(evt, c)
		}
	}

	for _, evt := range pg.usdAmountEditor.Editor.Events() {
		if pg.usdAmountEditor.Editor.Focused() {
			pg.handleEditorChange(evt, c)
		}
	}

	for _, evt := range pg.passwordEditor.Editor.Events() {
		if pg.passwordEditor.Editor.Focused() {
			pg.handleEditorChange(evt, c)
		}
	}

	if pg.calculateErrorText != "" {
		pg.dcrAmountEditor.LineColor, pg.dcrAmountEditor.TitleLabelColor = pg.theme.Color.Danger, pg.theme.Color.Danger
		pg.usdAmountEditor.LineColor, pg.usdAmountEditor.TitleLabelColor = pg.theme.Color.Danger, pg.theme.Color.Danger
		c.notify(pg.calculateErrorText, false)
	} else {
		pg.dcrAmountEditor.LineColor, pg.dcrAmountEditor.TitleLabelColor = pg.theme.Color.Gray1, pg.theme.Color.Gray3
		pg.usdAmountEditor.LineColor, pg.usdAmountEditor.TitleLabelColor = pg.theme.Color.Gray1, pg.theme.Color.Gray3
	}

	if pg.amountErrorText != "" {
		pg.dcrAmountEditor.SetError(pg.amountErrorText)
	}

	if pg.walletSelected != c.wallAcctSelector.selectedSendWallet {
		pg.shouldInitializeTxAuthor = true
		pg.walletSelected = c.wallAcctSelector.selectedSendWallet
	}

	if pg.shouldInitializeTxAuthor {
		pg.shouldInitializeTxAuthor = false
		pg.dcrAmountEditor.Editor.SetText("")
		pg.usdAmountEditor.Editor.SetText("")
		pg.calculateErrorText = ""
		c.wallet.CreateTransaction(sendWallet.ID, sendAcct.Number, pg.txAuthorErrChan)
	}

	activeAmountEditor := pg.dcrAmountEditor.Editor
	if pg.usdAmountEditor.Editor.Focused() {
		activeAmountEditor = pg.usdAmountEditor.Editor
	}
	if !pg.inputsNotEmpty(pg.destinationAddressEditor.Editor, activeAmountEditor) {
		pg.balanceAfterSend(true, c)
	}

	pg.watchForBroadcastResult(c)

	for pg.confirmButton.Button.Clicked() {
		pg.sendFund(c)
	}

	for pg.nextButton.Button.Clicked() {
		if pg.validate(c) && pg.calculateErrorText == "" {
			pg.isConfirmationModalOpen = true
			pg.passwordEditor.Editor.Focus()
		}
	}

	if pg.isConfirmationModalOpen {
		pg.confirmButton.Background = pg.theme.Color.Primary
		if !pg.inputsNotEmpty(pg.passwordEditor.Editor) {
			pg.confirmButton.Background = pg.theme.Color.InactiveGray
		}
	}

	for pg.closeConfirmationModalButton.Button.Clicked() {
		c.modalLoad.setLoading(false)
		pg.isConfirmationModalOpen = false
	}

	for pg.clearAllBtn.Button.Clicked() {
		pg.resetFields()
	}

	select {
	case err := <-pg.txAuthorErrChan:
		pg.calculateErrorText = err.Error()
		c.notify(pg.calculateErrorText, false)
	case err := <-pg.broadcastErrChan:
		if err.Error() == invalidPassphraseError {
			pg.passwordEditor.SetError("Wrong password")
		} else {
			c.notify(err.Error(), false)
			pg.isConfirmationModalOpen = false
		}
		c.modalLoad.setLoading(false)
		pg.isBroadcastingTransaction = false
	default:
	}

	// TODO
	if pg.dcrAmountEditor.CustomButton.Button.Clicked() {
		pg.dcrAmountEditor.Editor.Focus()
		// send max
	}
	if pg.usdAmountEditor.CustomButton.Button.Clicked() {
		pg.usdAmountEditor.Editor.Focus()
		// send max
	}
}

func (pg *sendPage) onClose() {}

package ui

import (
	"fmt"
	"image/color"
	// "reflect"
	// "strconv"
	"strings"
	// "time"
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	// "gioui.org/unit"
	"gioui.org/gesture"
	"gioui.org/widget"
	// "gioui.org/op/paint"

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

type sendPage struct {
	pageContainer, walletsList,
	accountsList layout.List
	theme *decredmaterial.Theme

	fromAccountBtn *widget.Clickable

	txAuthor        *dcrlibwallet.TxAuthor
	broadcastResult *wallet.Broadcast

	wallet          *wallet.Wallet
	selectedWallet  wallet.InfoShort
	selectedAccount **wallet.Account

	toAddress *widget.Bool

	unspentOutputsSelected *map[int]map[int32]map[string]*wallet.UnspentOutput

	destinationAddressEditor decredmaterial.Editor
	leftAmountEditor         decredmaterial.Editor
	rightAmountEditor        decredmaterial.Editor

	currencySwap, moreOption decredmaterial.IconButton

	customChangeAddressEditor    decredmaterial.Editor
	sendAmountEditor             decredmaterial.Editor
	nextButton                   decredmaterial.Button
	closeConfirmationModalButton decredmaterial.Button
	confirmButton                decredmaterial.Button
	maxButton                    decredmaterial.Button
	sendToButton                 decredmaterial.Button

	accountSwitch *decredmaterial.SwitchButtonText

	confirmModal       *decredmaterial.Modal
	walletAccountModal *decredmaterial.Modal

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
	currencyValue       string

	activeExchangeValue   string
	inactiveExchangeValue string

	activeTransactionFeeValue   string
	inactiveTransactionFeeValue string

	activeTotalCostValue   string
	inactiveTotalCostValue string

	// walletName     string
	// accountName    string
	// accountBalance string

	balanceAfterSendValue string

	LastTradeRate string

	passwordModal *decredmaterial.Password
	line          *decredmaterial.Line

	isConfirmationModalOpen   bool
	isPasswordModalOpen       bool
	isBroadcastingTransaction bool

	isWalletAccountModalOpen bool

	shouldInitializeTxAuthor bool

	txAuthorErrChan  chan error
	broadcastErrChan chan error

	borderColor color.NRGBA

	toggleCoinCtrl      *widget.Bool
	inputButtonCoinCtrl decredmaterial.Button

	toAcctDetails []*gesture.Click
}

const (
	PageSend               = "Send"
	invalidPassphraseError = "error broadcasting transaction: " + dcrlibwallet.ErrInvalidPassphrase
)

func (win *Window) SendPage(common pageCommon) layout.Widget {
	pg := &sendPage{
		pageContainer: layout.List{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		},
		accountsList: layout.List{
			Axis: layout.Vertical,
		},

		walletsList: layout.List{
			Axis: layout.Vertical,
		},

		theme:                  common.theme,
		wallet:                 common.wallet,
		txAuthor:               &win.txAuthor,
		broadcastResult:        &win.broadcastResult,
		unspentOutputsSelected: &common.selectedUTXO,
		selectedAccount:        &win.walletAccount,

		fromAccountBtn: new(widget.Clickable),
		toAddress:      new(widget.Bool),

		accountSwitch:         common.theme.SwitchButtonText("Address", "My Account", new(widget.Clickable), new(widget.Clickable)),
		activeExchangeValue:   "DCR",
		inactiveExchangeValue: "USD",

		closeConfirmationModalButton: common.theme.Button(new(widget.Clickable), "Close"),
		nextButton:                   common.theme.Button(new(widget.Clickable), "Next"),
		confirmButton:                common.theme.Button(new(widget.Clickable), "Confirm"),
		maxButton:                    common.theme.Button(new(widget.Clickable), "MAX"),
		txFeeCollapsible:             common.theme.Collapsible(),
		txLine:                       common.theme.Line(),

		confirmModal:              common.theme.Modal(),
		walletAccountModal:        common.theme.Modal(),
		isConfirmationModalOpen:   false,
		isPasswordModalOpen:       false,
		isBroadcastingTransaction: false,
		isWalletAccountModalOpen:  false,

		passwordModal:    common.theme.Password(),
		broadcastErrChan: make(chan error),
		txAuthorErrChan:  make(chan error),
		line:             common.theme.Line(),
	}

	pg.toAcctDetails = make([]*gesture.Click, 0)

	pg.line.Color = common.theme.Color.Gray
	pg.line.Height = 2

	pg.borderColor = common.theme.Color.Hint

	pg.balanceAfterSendValue = "- DCR"

	activeEditorHint := fmt.Sprintf("Amount (%s)", pg.activeExchangeValue)
	pg.leftAmountEditor = common.theme.Editor(new(widget.Editor), activeEditorHint)
	pg.leftAmountEditor.Editor.SetText("")
	pg.leftAmountEditor.IsCustomButton = true
	pg.leftAmountEditor.Editor.SingleLine = true
	pg.leftAmountEditor.CustomButton.Background = common.theme.Color.Gray
	pg.leftAmountEditor.CustomButton.Inset = layout.UniformInset(values.MarginPadding2)
	pg.leftAmountEditor.CustomButton.Text = "Max"
	pg.leftAmountEditor.CustomButton.CornerRadius = values.MarginPadding0

	inactiveEditorHint := fmt.Sprintf("Amount (%s)", pg.inactiveExchangeValue)
	pg.rightAmountEditor = common.theme.Editor(new(widget.Editor), inactiveEditorHint)
	pg.rightAmountEditor.Editor.SetText("")
	pg.rightAmountEditor.IsCustomButton = true
	pg.rightAmountEditor.Editor.SingleLine = true
	pg.rightAmountEditor.CustomButton.Background = common.theme.Color.Gray
	pg.rightAmountEditor.CustomButton.Inset = layout.UniformInset(values.MarginPadding2)
	pg.rightAmountEditor.CustomButton.Text = "Max"
	pg.rightAmountEditor.CustomButton.CornerRadius = values.MarginPadding0

	pg.destinationAddressEditor = common.theme.Editor(new(widget.Editor), "Address")
	pg.destinationAddressEditor.Editor.SingleLine, pg.destinationAddressEditor.IsVisible = true, true
	pg.destinationAddressEditor.Editor.SetText("")

	pg.customChangeAddressEditor = common.theme.Editor(new(widget.Editor), "Custom Change Address")
	pg.customChangeAddressEditor.IsVisible, pg.customChangeAddressEditor.IsTitleLabel = true, true
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

	pg.currencySwap = common.theme.IconButton(new(widget.Clickable), common.icons.actionSwapHoriz)
	pg.currencySwap.Background = color.NRGBA{}
	pg.currencySwap.Color = common.theme.Color.Text
	pg.currencySwap.Inset = layout.UniformInset(values.MarginPadding0)
	pg.currencySwap.Size = values.MarginPadding25

	pg.moreOption = common.theme.IconButton(new(widget.Clickable), common.icons.navMoreIcon)
	pg.moreOption.Background = color.NRGBA{}
	pg.moreOption.Color = common.theme.Color.Text
	pg.moreOption.Inset = layout.UniformInset(values.MarginPadding0)

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

func (pg *sendPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	if *pg.selectedAccount == nil {
		pg.selectedWallet = common.info.Wallets[0]
		*pg.selectedAccount = &common.info.Wallets[0].Accounts[0]
	}

	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.topNav(gtx, common)
		},
		func(gtx C) D {
			return pg.fromSection(gtx, common)
		},
		func(gtx C) D {
			return pg.toSection(gtx, common)
		},
		func(gtx C) D {
			return pg.feeSection(gtx)
		},
		// func(gtx C) D {
		// 	return pg.sendAmountSection(gtx)
		// },
		// func(gtx C) D {
		// 	gtx.Constraints.Max.X = gtx.Px(values.MarginPadding450)
		// 	return pg.drawTransactionDetailWidgets(gtx)
		// },
		// func(gtx C) D {
		// 	gtx.Constraints.Min.X = gtx.Px(values.MarginPadding450)
		// 	return pg.nextButton.Layout(gtx)
		// },
	}

	dims := layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return common.Layout(gtx, func(gtx C) D {
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

	// if pg.isConfirmationModalOpen {
	// 	return common.Modal(gtx, dims, pg.drawConfirmationModal(gtx))
	// }

	if pg.isWalletAccountModalOpen {
		return common.Modal(gtx, dims, pg.walletAccountSection(gtx, common))
	}

	// if pg.isPasswordModalOpen {
	// 	return common.Modal(gtx, dims, pg.drawPasswordModal(gtx))
	// }

	return dims
}

func (pg *sendPage) topNav(gtx layout.Context, common pageCommon) layout.Dimensions {
	m := values.MarginPadding15
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					common.subPageBackButton.Icon = common.icons.contentClear
					return common.subPageBackButton.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					txt := pg.theme.H6("Send DCR")
					txt.Color = pg.theme.Color.Text
					return layout.Inset{Left: m}.Layout(gtx, func(gtx C) D {
						return txt.Layout(gtx)
					})
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return common.subPageInfoButton.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Left: m}.Layout(gtx, func(gtx C) D {
							return pg.moreOption.Layout(gtx)
						})
					}),
				)
			})
		}),
	)
}

func (pg *sendPage) fromSection(gtx layout.Context, common pageCommon) layout.Dimensions {
	return pg.pageSections(gtx, "From", func(gtx C) D {
		border := widget.Border{Color: pg.theme.Color.Background, CornerRadius: values.MarginPadding5, Width: values.MarginPadding2}
		return border.Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
				return decredmaterial.Clickable(gtx, pg.fromAccountBtn, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							accountIcon := common.icons.accountIcon
							accountIcon.Scale = 0.9

							inset := layout.Inset{
								Right: values.MarginPadding10,
								Top:   values.MarginPadding2,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return accountIcon.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return pg.theme.Body1((*pg.selectedAccount).Name).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Left: values.MarginPadding5,
								Top:  values.MarginPadding2,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return decredmaterial.Card{
									Color: pg.theme.Color.Background,
								}.Layout(gtx, func(gtx C) D {
									return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
										text := pg.theme.Caption(pg.selectedWallet.Name)
										text.Color = pg.theme.Color.Gray
										return text.Layout(gtx)
									})
								})
							})
						}),
						layout.Flexed(1, func(gtx C) D {
							return layout.E.Layout(gtx, func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return pg.theme.Body1((*pg.selectedAccount).TotalBalance).Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {

										inset := layout.Inset{
											Left: values.MarginPadding5,
										}
										return inset.Layout(gtx, func(gtx C) D {
											icon := common.icons.collapseIcon
											icon.Scale = 0.25
											return icon.Layout(gtx)
										})
									}),
								)
							})
						}),
					)
				})
			})
		})
	})
}

func (pg *sendPage) toSection(gtx layout.Context, common pageCommon) layout.Dimensions {
	return pg.pageSections(gtx, "To", func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Top:    values.MarginPadding10,
					Bottom: values.MarginPadding10,
				}.Layout(gtx, func(gtx C) D {
					return pg.destinationAddressEditor.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				if strings.Contains(pg.currencyValue, "USD") {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(0.45, func(gtx C) D {
							pg.leftAmountEditor.Hint = fmt.Sprintf("Amount (%s)", pg.activeExchangeValue)
							return pg.leftAmountEditor.Layout(gtx)
						}),
						layout.Flexed(0.1, func(gtx C) D {
							return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
								return pg.currencySwap.Layout(gtx)
							})
						}),
						layout.Flexed(0.45, func(gtx C) D {
							pg.rightAmountEditor.Hint = fmt.Sprintf("Amount (%s)", pg.inactiveExchangeValue)
							return pg.rightAmountEditor.Layout(gtx)
						}),
					)
				}
				return pg.leftAmountEditor.Layout(gtx)
			}),
		)
	})
}

func (pg *sendPage) feeSection(gtx layout.Context) layout.Dimensions {
	collapsibleHeader := func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				b := pg.theme.Body1(pg.activeTransactionFeeValue)
				b.Color = pg.theme.Color.Text
				return b.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				b := pg.theme.Body1(pg.inactiveTransactionFeeValue)
				b.Color = pg.theme.Color.Hint
				inset := layout.Inset{
					Left: values.MarginPadding5,
				}
				if strings.Contains(pg.currencyValue, "USD") {
					return inset.Layout(gtx, func(gtx C) D {
						return b.Layout(gtx)
					})
				}
				return layout.Dimensions{}
			}),
		)
	}

	collapsibleBody := func(gtx C) D {
		card := pg.theme.Card()
		card.Color = pg.theme.Color.Background
		inset := layout.Inset{
			Top: values.MarginPadding10,
		}
		return inset.Layout(gtx, func(gtx C) D {
			return card.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return pg.contentRow(gtx, "Estimated time", "-")
						}),
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Top:    values.MarginPadding5,
								Bottom: values.MarginPadding5,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return pg.contentRow(gtx, "Estimated size", "-")
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
	return pg.pageSections(gtx, "Fee", func(gtx C) D {
		return pg.txFeeCollapsible.Layout(gtx, collapsibleHeader, collapsibleBody)
	})
}

func (pg *sendPage) walletAccountSection(gtx layout.Context, common pageCommon) layout.Dimensions {
	sections := func(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				txt := pg.theme.Body2(title)
				txt.Color = pg.theme.Color.Text
				inset := layout.Inset{
					Bottom: values.MarginPadding15,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return txt.Layout(gtx)
				})
			}),
			layout.Rigid(body),
		)
	}

	w := []func(gtx C) D{
		func(gtx C) D {
			txt := pg.theme.H6("Sending account")
			txt.Color = pg.theme.Color.Text
			return txt.Layout(gtx)
		},
		func(gtx C) D {
			return pg.walletsList.Layout(gtx, len(common.info.Wallets), func(gtx C, i int) D {
				wn := common.info.Wallets[i].Name
				accounts := common.info.Wallets[i].Accounts
				wIndex := i

				pg.updateAcctDetailsButtons(&accounts)
				inset := layout.Inset{
					Bottom: values.MarginPadding10,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return sections(gtx, wn, func(gtx C) D {
						return pg.accountsList.Layout(gtx, len(accounts), func(gtx C, x int) D {
							accountsName := accounts[x].Name
							totalBalance := accounts[x].TotalBalance
							spendable := dcrutil.Amount(accounts[x].SpendableBalance).String()
							aIndex := x

							click := pg.toAcctDetails[x]
							pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
							click.Add(gtx.Ops)
							pg.goToAcctDetails(gtx, common, &accounts[x], wIndex, aIndex, click)

							return pg.walletAccountsLayout(gtx, accountsName, totalBalance, spendable, common, wIndex, aIndex)
						})
					})
				})
			})
		},
	}

	return pg.walletAccountModal.Layout(gtx, w, 850)
}

func (pg *sendPage) walletAccountsLayout(gtx layout.Context, name, totalBal, spendableBal string, common pageCommon, wIndex, aIndex int) layout.Dimensions {
	accountIcon := common.icons.accountIcon
	accountIcon.Scale = 0.8

	inset := layout.Inset{
		Bottom: values.MarginPadding10,
	}
	return inset.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(0.1, func(gtx C) D {
						inset := layout.Inset{
							Right: values.MarginPadding10,
							Top:   values.MarginPadding15,
						}
						return inset.Layout(gtx, func(gtx C) D {
							return accountIcon.Layout(gtx)
						})
					}),
					layout.Flexed(0.8, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								accountLabel := pg.theme.Body2(name)
								accountLabel.Color = pg.theme.Color.Text
								accountBalLabel := pg.theme.Body2(totalBal)
								accountBalLabel.Color = pg.theme.Color.Text
								return pg.accountTableLayout(gtx, accountLabel, accountBalLabel)
							}),
							layout.Rigid(func(gtx C) D {
								spendibleLabel := pg.theme.Body2("Spendable")
								spendibleLabel.Color = pg.theme.Color.Gray
								spendibleBalLabel := pg.theme.Body2(spendableBal)
								spendibleBalLabel.Color = pg.theme.Color.Gray
								return pg.accountTableLayout(gtx, spendibleLabel, spendibleBalLabel)
							}),
						)
					}),
					layout.Flexed(0.1, func(gtx C) D {
						// 			inset := layout.Inset{
						// 				Right: values.MarginPadding10,
						// 				Top:   values.MarginPadding15,
						// 			}

						// 			// fmt.Println(*common.selectedWallet ,"==", wIndex)
						// 				// if *common.selectedWallet == wIndex {
						// 			if *common.selectedAccount == aIndex {
						// return layout.E.Layout(gtx, func(gtx C) D {
						// 				return inset.Layout(gtx, func(gtx C) D {
						// 					return common.icons.navigationCheck.Layout(gtx, values.MarginPadding20)
						// 				})
						// 			})
						// 			// }
						// 		}
						return layout.Dimensions{}

					}),
				)
			}),
		)
	})
}

func (pg *sendPage) accountTableLayout(gtx layout.Context, leftLabel, rightLabel decredmaterial.Label) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			inset := layout.Inset{
				Top: values.MarginPadding2,
			}
			return inset.Layout(gtx, func(gtx C) D {
				return leftLabel.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return rightLabel.Layout(gtx)
			})
		}),
	)
}

func (pg *sendPage) updateAcctDetailsButtons(walAcct *[]wallet.Account) {
	if len(*walAcct) != len(pg.toAcctDetails) {
		for i := 0; i < len(*walAcct); i++ {
			pg.toAcctDetails = append(pg.toAcctDetails, &gesture.Click{})
		}
	}
}

func (pg *sendPage) goToAcctDetails(gtx layout.Context, common pageCommon, acct *wallet.Account, wIndex, aIndex int, click *gesture.Click) {
	for _, e := range click.Events(gtx) {
		if e.Type == gesture.TypeClick {
			fmt.Println(wIndex, "  ", aIndex)
			*pg.selectedAccount = acct
			// *common.selectedAccount = aIndex
			*common.selectedWallet = wIndex
			pg.selectedWallet = common.info.Wallets[wIndex]
			pg.isWalletAccountModalOpen = false
		}
	}
}

func (pg *sendPage) pageSections(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
	return pg.theme.Card().Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := pg.theme.Body1(title)
							txt.Color = pg.theme.Color.Text
							inset := layout.Inset{
								Bottom: values.MarginPadding10,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return txt.Layout(gtx)
							})
						}),
						layout.Flexed(1, func(gtx C) D {
							if title == "To" {
								return layout.E.Layout(gtx, func(gtx C) D {
									return pg.accountSwitch.Layout(gtx)
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
			return pg.theme.Body2(leftValue).Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.theme.Body2(rightValue).Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (pg *sendPage) calculateValues() {
	defaultActiveValues := fmt.Sprintf("- %s", "DCR" /*pg.activeExchangeValue*/)
	defaultInactiveValues := "($ -)"
	// noExchangeText := "Exchange rate not fetched"
	// pg.sendAmountEditor.Hint = "0"

	pg.activeTransactionFeeValue = defaultActiveValues
	// pg.activeTotalCostValue = defaultActiveValues
	pg.inactiveTransactionFeeValue = defaultInactiveValues
	// pg.inactiveTotalCostValue = defaultInactiveValues

	// pg.calculateErrorText = ""
	// pg.activeTotalAmount = pg.activeExchangeValue
	// pg.inactiveTotalAmount = fmt.Sprintf("0 %s", pg.inactiveExchangeValue)

	// // default values when exchange is not available
	// if pg.LastTradeRate == "" {
	// 	pg.activeTransactionFeeValue = defaultActiveValues
	// 	pg.activeTotalCostValue = defaultActiveValues
	// 	pg.inactiveTransactionFeeValue = ""
	// 	pg.inactiveTotalCostValue = ""
	// 	pg.activeTotalAmount = pg.activeExchangeValue
	// 	pg.inactiveTotalAmount = noExchangeText
	// }

	// if reflect.DeepEqual(pg.txAuthor, &dcrlibwallet.TxAuthor{}) || !pg.validate(true) {
	// 	return
	// }

	// pg.inputAmount, _ = strconv.ParseFloat(pg.sendAmountEditor.Editor.Text(), 64)

	// if pg.LastTradeRate != "" {
	// 	pg.usdExchangeRate, _ = strconv.ParseFloat(pg.LastTradeRate, 64)
	// 	pg.amountUSDtoDCR = pg.inputAmount / pg.usdExchangeRate
	// 	pg.amountDCRtoUSD = pg.inputAmount * pg.usdExchangeRate
	// }

	// pg.setChangeDestinationAddr()
	// if pg.activeExchangeValue == "USD" && pg.LastTradeRate != "" {
	// 	pg.amountAtoms = pg.setDestinationAddr(pg.amountUSDtoDCR)
	// } else {
	// 	pg.amountAtoms = pg.setDestinationAddr(pg.inputAmount)
	// }

	// if pg.amountAtoms == 0 {
	// 	return
	// }

	// pg.txFee = pg.getTxFee(pg.toggleCoinCtrl.Value)
	// if pg.txFee == 0 {
	// 	return
	// }

	// pg.totalCostDCR = pg.txFee + pg.amountAtoms

	// pg.updateDefaultValues()
	// pg.balanceAfterSend(false)
}

// func (pg *SendPage) updateDefaultValues() {
// 	v := pg.amountValues()
// 	pg.activeTotalAmount = pg.activeExchangeValue
// 	pg.inactiveTotalAmount = v.inactiveTotalAmount
// 	pg.activeTransactionFeeValue = v.activeTransactionFeeValue
// 	pg.inactiveTransactionFeeValue = v.inactiveTransactionFeeValue
// 	pg.activeTotalCostValue = v.activeTotalCostValue
// 	pg.inactiveTotalCostValue = v.inactiveTotalCostValue
// }

func (pg *sendPage) Handle(c pageCommon) {
	pg.calculateValues()

	if c.subPageBackButton.Button.Clicked() {
		*c.page = PageOverview
	}

	if pg.fromAccountBtn.Clicked() {
		pg.isWalletAccountModalOpen = true
	}

	if len(c.info.Wallets) == 0 {
		return
	}

	pg.currencyValue = pg.wallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	if pg.currencyValue == "" {
		pg.currencyValue = "None"
	}

	for pg.currencySwap.Button.Clicked() {
		// if pg.LastTradeRate != "" {
		if pg.activeExchangeValue == "DCR" {
			pg.activeExchangeValue = "USD"
			pg.inactiveExchangeValue = "DCR"
		} else {
			pg.activeExchangeValue = "DCR"
			pg.inactiveExchangeValue = "USD"
		}
		// }

		pg.calculateValues()
	}
}

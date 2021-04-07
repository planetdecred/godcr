package ui

import (
	"fmt"
	"image"
	"image/color"
	"reflect"
	"strconv"
	"strings"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"

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

type walletAccount struct {
	selectAccount []*gesture.Click
}

type amountValue struct {
	sendAmountDCR            string
	sendAmountUSD            string
	leftTransactionFeeValue  string
	rightTransactionFeeValue string
	leftTotalCostValue       string
	rightTotalCostValue      string
}

type sendPage struct {
	pageContainer, walletsList,
	accountsList layout.List
	theme *decredmaterial.Theme

	fromAccountBtn, currencySwap *widget.Clickable

	txAuthor        *dcrlibwallet.TxAuthor
	broadcastResult *wallet.Broadcast

	wallet          *wallet.Wallet
	selectedWallet  wallet.InfoShort
	selectedAccount **wallet.Account

	toAddress    *widget.Bool
	isMoreOption bool

	walletAccount  walletAccount
	walletAccounts map[int]walletAccount

	destinationAddressEditor decredmaterial.Editor
	leftAmountEditor         decredmaterial.Editor
	rightAmountEditor        decredmaterial.Editor
	passwordEditor           decredmaterial.Editor

	moreOption decredmaterial.IconButton

	nextButton                   decredmaterial.Button
	closeConfirmationModalButton decredmaterial.Button
	confirmButton                decredmaterial.Button
	maxButton                    decredmaterial.Button
	clearAllBtn                  decredmaterial.Button

	accountSwitch *decredmaterial.SwitchButtonText

	confirmModal       *decredmaterial.Modal
	walletAccountModal *decredmaterial.Modal

	txFeeCollapsible *decredmaterial.Collapsible

	count              int
	remainingBalance   int64
	amountAtoms        int64
	totalCostDCR       int64
	txFee              int64
	txFeeSize          string
	txFeeEstimatedTime string
	spendableBalance   int64

	usdExchangeRate float64
	inputAmount     float64
	amountUSDtoDCR  float64
	amountDCRtoUSD  float64

	amountErrorText    string
	calculateErrorText string

	activeTotalAmount string

	leftExchangeValue  string
	rightExchangeValue string

	leftTransactionFeeValue  string
	rightTransactionFeeValue string

	leftTotalCostValue  string
	rightTotalCostValue string

	sendAmountDCR string
	sendAmountUSD string

	balanceAfterSendValue string
	LastTradeRate         string

	passwordModal *decredmaterial.Password

	isConfirmationModalOpen   bool
	isBroadcastingTransaction bool
	usdExchangeSet            bool
	isWalletAccountModalOpen  bool

	shouldInitializeTxAuthor bool

	txAuthorErrChan chan error

	broadcastErrChan chan error
}

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

		theme:           common.theme,
		wallet:          common.wallet,
		txAuthor:        &win.txAuthor,
		broadcastResult: &win.broadcastResult,
		selectedAccount: &win.walletAccount,

		fromAccountBtn: new(widget.Clickable),
		currencySwap:   new(widget.Clickable),
		toAddress:      new(widget.Bool),

		accountSwitch:                common.theme.SwitchButtonText("Address", "My Account", new(widget.Clickable), new(widget.Clickable)),
		leftExchangeValue:            "DCR",
		rightExchangeValue:           "USD",
		closeConfirmationModalButton: common.theme.Button(new(widget.Clickable), "Cancel"),
		confirmButton:                common.theme.Button(new(widget.Clickable), ""),
		maxButton:                    common.theme.Button(new(widget.Clickable), "MAX"),
		clearAllBtn:                  common.theme.Button(new(widget.Clickable), "Clear all fields"),
		txFeeCollapsible:             common.theme.Collapsible(),

		confirmModal:              common.theme.Modal(),
		walletAccountModal:        common.theme.Modal(),
		isConfirmationModalOpen:   false,
		isBroadcastingTransaction: false,
		isWalletAccountModalOpen:  false,

		passwordModal:    common.theme.Password(),
		broadcastErrChan: make(chan error),
		txAuthorErrChan:  make(chan error),
	}

	pg.walletAccounts = make(map[int]walletAccount)

	pg.balanceAfterSendValue = "- DCR"

	pg.nextButton = common.theme.Button(new(widget.Clickable), "Next")
	pg.nextButton.Background = pg.theme.Color.InactiveGray

	activeEditorHint := fmt.Sprintf("Amount (%s)", pg.leftExchangeValue)
	pg.leftAmountEditor = common.theme.Editor(new(widget.Editor), activeEditorHint)
	pg.leftAmountEditor.Editor.SetText("")
	pg.leftAmountEditor.IsCustomButton = true
	pg.leftAmountEditor.Editor.SingleLine = true
	pg.leftAmountEditor.CustomButton.Background = common.theme.Color.Gray
	pg.leftAmountEditor.CustomButton.Inset = layout.UniformInset(values.MarginPadding2)
	pg.leftAmountEditor.CustomButton.Text = "Max"
	pg.leftAmountEditor.CustomButton.CornerRadius = values.MarginPadding0

	inactiveEditorHint := fmt.Sprintf("Amount (%s)", pg.rightExchangeValue)
	pg.rightAmountEditor = common.theme.Editor(new(widget.Editor), inactiveEditorHint)
	pg.rightAmountEditor.Editor.SetText("")
	pg.rightAmountEditor.IsCustomButton = true
	pg.rightAmountEditor.Editor.SingleLine = true
	pg.rightAmountEditor.CustomButton.Background = common.theme.Color.Gray
	pg.rightAmountEditor.CustomButton.Inset = layout.UniformInset(values.MarginPadding2)
	pg.rightAmountEditor.CustomButton.Text = "Max"
	pg.rightAmountEditor.CustomButton.CornerRadius = values.MarginPadding0

	editorColor := common.theme.Color.Primary
	pg.passwordEditor = common.theme.Editor(new(widget.Editor), "Spending password")
	pg.passwordEditor.IsTitleLabel = true
	pg.passwordEditor.Editor.SetText("")
	pg.passwordEditor.Editor.SingleLine = true
	pg.passwordEditor.Editor.Mask = '*'
	pg.passwordEditor.LineColor, pg.passwordEditor.TitleLabelColor = editorColor, editorColor

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

	pg.clearAllBtn.Background = common.theme.Color.Surface
	pg.clearAllBtn.Color = common.theme.Color.Text
	pg.clearAllBtn.Inset = layout.UniformInset(values.MarginPadding15)

	pg.fetchExchangeValue()

	return func(gtx C) D {
		pg.Handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *sendPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	for index := 0; index < common.info.LoadedWallets; index++ {
		if _, ok := pg.walletAccounts[index]; !ok {
			accounts := common.info.Wallets[index].Accounts
			if len(accounts) != len(pg.walletAccount.selectAccount) {
				pg.walletAccount.selectAccount = make([]*gesture.Click, len(accounts))
				for index := range accounts {
					pg.walletAccount.selectAccount[index] = &gesture.Click{}
				}
			}

			pg.walletAccounts[index] = pg.walletAccount
		}
	}

	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.fromSection(gtx, common)
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
								return border.Layout(gtx, func(gtx C) D {
									return pg.clearAllBtn.Layout(gtx)
								})
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

	if pg.isWalletAccountModalOpen {
		return common.Modal(gtx, dims, pg.walletAccountSection(gtx, common))
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
					return layout.Inset{Left: m}.Layout(gtx, func(gtx C) D {
						return pg.theme.H6("Send DCR").Layout(gtx)
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
		border := widget.Border{Color: pg.theme.Color.Gray1, CornerRadius: values.MarginPadding8, Width: values.MarginPadding2}
		return border.Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
				return decredmaterial.Clickable(gtx, pg.fromAccountBtn, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							accountIcon := common.icons.accountIcon
							accountIcon.Scale = 0.9

							inset := layout.Inset{
								Right: values.MarginPadding8,
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
								Left: values.MarginPadding4,
								Top:  values.MarginPadding2,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return decredmaterial.Card{
									Color: pg.theme.Color.LightGray,
								}.Layout(gtx, func(gtx C) D {
									m2 := values.MarginPadding2
									m4 := values.MarginPadding4
									inset := layout.Inset{
										Left:   m4,
										Top:    m2,
										Bottom: m2,
										Right:  m4,
									}
									return inset.Layout(gtx, func(gtx C) D {
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
											Left: values.MarginPadding15,
										}
										return inset.Layout(gtx, func(gtx C) D {
											return common.icons.dropDownIcon.Layout(gtx, values.MarginPadding20)
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
					Bottom: values.MarginPadding16,
				}.Layout(gtx, func(gtx C) D {
					return pg.destinationAddressEditor.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				if pg.usdExchangeSet {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(0.45, func(gtx C) D {
							pg.leftAmountEditor.Hint = fmt.Sprintf("Amount (%s)", pg.leftExchangeValue)
							return pg.leftAmountEditor.Layout(gtx)
						}),
						layout.Flexed(0.1, func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									return decredmaterial.Clickable(gtx, pg.currencySwap, func(gtx C) D {
										icon := common.icons.currencySwapIcon
										icon.Scale = 0.45
										return icon.Layout(gtx)
									})
								})
							})
						}),
						layout.Flexed(0.45, func(gtx C) D {
							pg.rightAmountEditor.Hint = fmt.Sprintf("Amount (%s)", pg.rightExchangeValue)
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
				return pg.theme.Body1(pg.leftTransactionFeeValue).Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				b := pg.theme.Body1(pg.rightTransactionFeeValue)
				b.Color = pg.theme.Color.Gray
				inset := layout.Inset{
					Left: values.MarginPadding5,
				}
				if pg.usdExchangeSet {
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

func (pg *sendPage) walletAccountSection(gtx layout.Context, common pageCommon) layout.Dimensions {
	sections := func(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
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

				inset := layout.Inset{
					Bottom: values.MarginPadding10,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return sections(gtx, wn, func(gtx C) D {
						return pg.accountsList.Layout(gtx, len(accounts), func(gtx C, x int) D {
							accountsName := accounts[x].Name
							totalBalance := accounts[x].TotalBalance
							spendable := dcrutil.Amount(accounts[x].SpendableBalance).String()

							click := pg.walletAccounts[wIndex].selectAccount[x]
							pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
							click.Add(gtx.Ops)
							pg.walletAccountsHandler(gtx, common, &accounts[x], wIndex, x, click)

							return pg.walletAccountsPopupLayout(gtx, accountsName, totalBalance, spendable, common, wIndex, x)
						})
					})
				})
			})
		},
	}

	return pg.walletAccountModal.Layout(gtx, w, 850)
}

func (pg *sendPage) walletAccountsPopupLayout(gtx layout.Context, name, totalBal, spendableBal string, common pageCommon, wIndex, aIndex int) layout.Dimensions {
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
						inset := layout.Inset{
							Right: values.MarginPadding10,
							Top:   values.MarginPadding10,
						}

						if *common.selectedWallet == wIndex {
							if *common.selectedAccount == aIndex {
								return layout.E.Layout(gtx, func(gtx C) D {
									return inset.Layout(gtx, func(gtx C) D {
										return common.icons.navigationCheck.Layout(gtx, values.MarginPadding20)
									})
								})
							}
						}
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

func (pg *sendPage) walletAccountsHandler(gtx layout.Context, common pageCommon, acct *wallet.Account, wIndex, aIndex int, click *gesture.Click) {
	for _, e := range click.Events(gtx) {
		if e.Type == gesture.TypeClick {
			*pg.selectedAccount = acct
			*common.selectedAccount = aIndex
			*common.selectedWallet = wIndex
			pg.selectedWallet = common.info.Wallets[wIndex]
			pg.isWalletAccountModalOpen = false
		}
	}
}

func (pg *sendPage) confirmationModal(gtx layout.Context, common pageCommon) layout.Dimensions {
	w := []func(gtx C) D{
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
									return common.layoutBalance(gtx, pg.sendAmountDCR)
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
					return pg.contentRow(gtx, "Sending from", (*pg.selectedAccount).Name, pg.selectedWallet.Name)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding8, Bottom: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						return pg.contentRow(gtx, "Fee", pg.leftTransactionFeeValue, "")
					})
				}),
				layout.Rigid(func(gtx C) D {
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
							return inset.Layout(gtx, func(gtx C) D {
								return pg.theme.Body1(title).Layout(gtx)
							})
						}),
						layout.Flexed(1, func(gtx C) D {
							if title == "To" {
								return layout.E.Layout(gtx, func(gtx C) D {
									inset := layout.Inset{
										Top: values.MarginPaddingMinus5,
									}
									return inset.Layout(gtx, func(gtx C) D {
										return pg.accountSwitch.Layout(gtx)
									})
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
					layout.Rigid(func(gtx C) D {
						return pg.theme.Body1(rightValue).Layout(gtx)
					}),
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

func (pg *sendPage) validate() bool {
	isAmountValid := pg.validateLeftAmount()
	if pg.rightAmountEditor.Editor.Focused() {
		isAmountValid = pg.validateRightAmount()
	}

	if !pg.validateDestinationAddress() || !isAmountValid || pg.calculateErrorText != "" {
		pg.nextButton.Background = pg.theme.Color.Hint
		return false
	}

	pg.nextButton.Background = pg.theme.Color.Primary
	return true
}

func (pg *sendPage) validateDestinationAddress() bool {
	if pg.inputsNotEmpty(pg.destinationAddressEditor.Editor) {
		isValid, _ := pg.wallet.IsAddressValid(pg.destinationAddressEditor.Editor.Text())
		if !isValid {
			pg.destinationAddressEditor.SetError("Invalid address")
			return false
		}

		pg.destinationAddressEditor.SetError("")
		return true
	}

	pg.balanceAfterSend(true)
	pg.destinationAddressEditor.SetError("Input address")
	return false
}

func (pg *sendPage) validateLeftAmount() bool {
	if pg.inputsNotEmpty(pg.leftAmountEditor.Editor) {
		_, err := strconv.ParseFloat(pg.leftAmountEditor.Editor.Text(), 64)
		if err != nil {
			if strings.Contains(err.Error(), "invalid") {
				pg.leftAmountEditor.SetError("Invalid amount")
			}
			return false
		}
		pg.leftAmountEditor.SetError("")
		return true
	}
	pg.leftAmountEditor.SetError("")
	return false
}

func (pg *sendPage) validateRightAmount() bool {
	if pg.inputsNotEmpty(pg.rightAmountEditor.Editor) {
		_, err := strconv.ParseFloat(pg.rightAmountEditor.Editor.Text(), 64)
		if err != nil {
			if strings.Contains(err.Error(), "invalid") {
				pg.rightAmountEditor.SetError("Invalid amount")
			}
			return false
		}
		pg.rightAmountEditor.SetError("")
		return true
	}
	pg.rightAmountEditor.SetError("")
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

func (pg *sendPage) calculateValues() {
	defaultLeftValues := fmt.Sprintf("- %s", "DCR")
	defaultRightValues := "($ -)"

	pg.leftTransactionFeeValue = defaultLeftValues
	pg.rightTransactionFeeValue = defaultRightValues

	pg.leftTotalCostValue = defaultLeftValues
	pg.rightTotalCostValue = defaultRightValues
	pg.calculateErrorText = ""
	pg.txFeeSize = "-"
	pg.txFeeEstimatedTime = "-"
	pg.sendAmountDCR = defaultLeftValues
	pg.sendAmountUSD = defaultRightValues

	if reflect.DeepEqual(pg.txAuthor, &dcrlibwallet.TxAuthor{}) || !pg.validate() {
		return
	}

	pg.inputAmount, _ = strconv.ParseFloat(pg.leftAmountEditor.Editor.Text(), 64)
	if pg.usdExchangeSet && pg.rightAmountEditor.Editor.Focused() {
		pg.inputAmount, _ = strconv.ParseFloat(pg.rightAmountEditor.Editor.Text(), 64)
	}

	if pg.usdExchangeSet && pg.LastTradeRate != "" {
		pg.usdExchangeRate, _ = strconv.ParseFloat(pg.LastTradeRate, 64)
		pg.amountUSDtoDCR = pg.inputAmount / pg.usdExchangeRate
		pg.amountDCRtoUSD = pg.inputAmount * pg.usdExchangeRate
	}

	pg.updateAmountInputsValues()
	pg.getTxFee()
	pg.updateDefaultValues()
	pg.balanceAfterSend(false)
}

func (pg *sendPage) updateAmountInputsValues() {
	switch {
	case pg.leftExchangeValue == "USD" && pg.LastTradeRate != "" && pg.leftAmountEditor.Editor.Focused():
		pg.rightAmountEditor.Editor.SetText(fmt.Sprintf("%f", pg.amountUSDtoDCR))
		pg.setDestinationAddr(pg.amountUSDtoDCR)
	case pg.leftExchangeValue == "USD" && pg.LastTradeRate != "" && pg.rightAmountEditor.Editor.Focused():
		pg.leftAmountEditor.Editor.SetText(fmt.Sprintf("%f", pg.amountDCRtoUSD))
		pg.setDestinationAddr(pg.inputAmount)
	case pg.leftExchangeValue == "DCR" && pg.LastTradeRate != "" && pg.rightAmountEditor.Editor.Focused():
		pg.leftAmountEditor.Editor.SetText(fmt.Sprintf("%f", pg.amountUSDtoDCR))
		pg.setDestinationAddr(pg.amountUSDtoDCR)
	case pg.leftExchangeValue == "DCR" && pg.LastTradeRate != "" && pg.leftAmountEditor.Editor.Focused():
		pg.rightAmountEditor.Editor.SetText(fmt.Sprintf("%f", pg.amountDCRtoUSD))
		pg.setDestinationAddr(pg.inputAmount)
	default:
		if pg.rightAmountEditor.Editor.Focused() {
			pg.leftAmountEditor.Editor.SetText(pg.rightAmountEditor.Editor.Text())
		} else {
			pg.rightAmountEditor.Editor.SetText(pg.leftAmountEditor.Editor.Text())
		}
		pg.setDestinationAddr(pg.inputAmount)
	}
}

func (pg *sendPage) updateExchangeError() {
	pg.rightAmountEditor.SetError("")
	if pg.LastTradeRate == "" && pg.usdExchangeSet {
		pg.rightAmountEditor.SetError("Exchange rate not fetched")
	}
}

func (pg *sendPage) setDestinationAddr(sendAmount float64) {
	pg.amountErrorText = ""
	amount, err := dcrutil.NewAmount(sendAmount)
	if err != nil {
		pg.feeEstimationError(err.Error(), "amount")
		return
	}

	pg.amountAtoms = int64(amount)
	if pg.amountAtoms == 0 {
		return
	}

	pg.txAuthor.RemoveSendDestination(0)
	pg.txAuthor.AddSendDestination(pg.destinationAddressEditor.Editor.Text(), pg.amountAtoms, false)
}

func (pg *sendPage) amountValues() amountValue {
	pg.totalCostDCR = pg.txFee + pg.amountAtoms
	txFeeValueUSD := dcrutil.Amount(pg.txFee).ToCoin() * pg.usdExchangeRate
	switch {
	case pg.leftExchangeValue == "USD" && pg.LastTradeRate != "":
		return amountValue{
			sendAmountDCR:            dcrutil.Amount(pg.amountAtoms).String(),
			sendAmountUSD:            fmt.Sprintf("$ %f", dcrutil.Amount(pg.amountAtoms).ToCoin()*pg.usdExchangeRate),
			leftTransactionFeeValue:  fmt.Sprintf("%f USD", txFeeValueUSD),
			rightTransactionFeeValue: fmt.Sprintf("(%s)", dcrutil.Amount(pg.txFee).String()),
			leftTotalCostValue:       fmt.Sprintf("%s USD", strconv.FormatFloat(pg.inputAmount+txFeeValueUSD, 'f', 7, 64)),
			rightTotalCostValue:      fmt.Sprintf("(%s )", dcrutil.Amount(pg.totalCostDCR).String()),
		}
	case pg.leftExchangeValue == "DCR" && pg.LastTradeRate != "":
		return amountValue{
			sendAmountDCR:            dcrutil.Amount(pg.amountAtoms).String(),
			sendAmountUSD:            fmt.Sprintf("$ %s", strconv.FormatFloat(pg.amountDCRtoUSD, 'f', 2, 64)),
			leftTransactionFeeValue:  dcrutil.Amount(pg.txFee).String(),
			rightTransactionFeeValue: fmt.Sprintf("(%f USD)", txFeeValueUSD),
			leftTotalCostValue:       dcrutil.Amount(pg.totalCostDCR).String(),
			rightTotalCostValue:      fmt.Sprintf("(%s USD)", strconv.FormatFloat(pg.amountDCRtoUSD+txFeeValueUSD, 'f', 7, 64)),
		}
	default:
		return amountValue{
			sendAmountDCR:           dcrutil.Amount(pg.amountAtoms).String(),
			sendAmountUSD:           "",
			leftTransactionFeeValue: dcrutil.Amount(pg.txFee).String(),
			leftTotalCostValue:      dcrutil.Amount(pg.totalCostDCR).String(),
		}
	}
}

func (pg *sendPage) getTxFee() {
	// calculate transaction fee
	feeAndSize, err := pg.txAuthor.EstimateFeeAndSize()
	if err != nil {
		pg.feeEstimationError(err.Error(), "fee")
		return
	}

	pg.txFee = feeAndSize.Fee.AtomValue
	pg.txFeeSize = fmt.Sprintf("%v Bytes", feeAndSize.EstimatedSignedSize)
}

func (pg *sendPage) balanceAfterSend(isInputAmountEmpty bool) {
	pg.remainingBalance = 0
	if isInputAmountEmpty {
		pg.remainingBalance = (*pg.selectedAccount).SpendableBalance
	} else {
		pg.remainingBalance = (*pg.selectedAccount).SpendableBalance - pg.totalCostDCR
	}
	pg.balanceAfterSendValue = dcrutil.Amount(pg.remainingBalance).String()
}

func (pg *sendPage) feeEstimationError(err, errorPath string) {
	if err == "insufficient_balance" {
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
		*c.page = PageOverview
		c.notify("1 Transaction Sent", true)

		if pg.remainingBalance != -1 {
			pg.spendableBalance = pg.remainingBalance
		}
		pg.remainingBalance = -1

		pg.isConfirmationModalOpen = false
		pg.isBroadcastingTransaction = false
		pg.resetFields()
		pg.broadcastResult.TxHash = ""
		pg.calculateValues()
		pg.destinationAddressEditor.Editor.SetText("")
	}
}

func (pg *sendPage) handleEditorChange(evt widget.EditorEvent) {
	switch evt.(type) {
	case widget.ChangeEvent:
		pg.fetchExchangeValue()
		pg.calculateValues()
	}
}

func (pg *sendPage) updateDefaultValues() {
	v := pg.amountValues()
	pg.sendAmountDCR = v.sendAmountDCR
	pg.sendAmountUSD = v.sendAmountUSD
	pg.activeTotalAmount = pg.leftExchangeValue
	pg.leftTransactionFeeValue = v.leftTransactionFeeValue
	pg.rightTransactionFeeValue = v.rightTransactionFeeValue
	pg.leftTotalCostValue = v.leftTotalCostValue
	pg.rightTotalCostValue = v.rightTotalCostValue
}

func (pg *sendPage) resetFields() {
	pg.destinationAddressEditor.SetError("")
	pg.leftAmountEditor.Editor.SetText("")
	pg.rightAmountEditor.Editor.SetText("")
}

func (pg *sendPage) fetchExchangeValue() {
	go func() {
		err := pg.wallet.GetUSDExchangeValues(&pg)
		if err != nil {
			pg.updateExchangeError()
		}
	}()
}

func (pg *sendPage) Handle(c pageCommon) {
	if len(c.info.Wallets) == 0 {
		return
	}

	if *pg.selectedAccount == nil {
		pg.selectedWallet = c.info.Wallets[0]
		*pg.selectedAccount = &c.info.Wallets[0].Accounts[0]
		pg.shouldInitializeTxAuthor = true
	}

	if pg.LastTradeRate == "" && pg.count == 0 {
		pg.count = 1
		pg.calculateValues()
	}

	if (pg.LastTradeRate != "" && pg.count == 0) || (pg.LastTradeRate != "" && pg.count == 1) {
		pg.count = 2
		pg.calculateValues()
	}

	pg.updateExchangeError()

	if c.subPageBackButton.Button.Clicked() {
		*c.page = PageOverview
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

	if pg.fromAccountBtn.Clicked() {
		pg.isWalletAccountModalOpen = true
	}

	currencyExchangeValue := pg.wallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
	pg.usdExchangeSet = false
	if strings.Contains(currencyExchangeValue, "USD") {
		pg.usdExchangeSet = true
	}

	for range pg.destinationAddressEditor.Editor.Events() {
		pg.calculateValues()
	}

	for pg.currencySwap.Clicked() {
		if pg.LastTradeRate != "" {
			if pg.leftExchangeValue == "DCR" {
				pg.leftExchangeValue = "USD"
				pg.rightExchangeValue = "DCR"
			} else {
				pg.leftExchangeValue = "DCR"
				pg.rightExchangeValue = "USD"
			}
		}

		pg.calculateValues()
	}

	for _, evt := range pg.leftAmountEditor.Editor.Events() {
		if pg.leftAmountEditor.Editor.Focused() {
			pg.handleEditorChange(evt)
		}
	}

	for _, evt := range pg.rightAmountEditor.Editor.Events() {
		if pg.rightAmountEditor.Editor.Focused() {
			pg.handleEditorChange(evt)
		}
	}

	if pg.calculateErrorText != "" {
		pg.leftAmountEditor.LineColor, pg.leftAmountEditor.TitleLabelColor = pg.theme.Color.Danger, pg.theme.Color.Danger
		pg.rightAmountEditor.LineColor, pg.rightAmountEditor.TitleLabelColor = pg.theme.Color.Danger, pg.theme.Color.Danger
		c.notify(pg.calculateErrorText, false)
	} else {
		pg.leftAmountEditor.LineColor, pg.leftAmountEditor.TitleLabelColor = pg.theme.Color.Gray1, pg.theme.Color.Gray3
		pg.rightAmountEditor.LineColor, pg.rightAmountEditor.TitleLabelColor = pg.theme.Color.Gray1, pg.theme.Color.Gray3
	}

	if pg.amountErrorText != "" {
		pg.leftAmountEditor.SetError(pg.amountErrorText)
	}

	if pg.shouldInitializeTxAuthor {
		pg.shouldInitializeTxAuthor = false
		pg.leftAmountEditor.Editor.SetText("")
		pg.rightAmountEditor.Editor.SetText("")
		pg.calculateErrorText = ""
		c.wallet.CreateTransaction(pg.selectedWallet.ID, (*pg.selectedAccount).Number, pg.txAuthorErrChan)
	}

	activeAmountEditor := pg.leftAmountEditor.Editor
	if pg.rightAmountEditor.Editor.Focused() {
		activeAmountEditor = pg.rightAmountEditor.Editor
	}
	if !pg.inputsNotEmpty(pg.destinationAddressEditor.Editor, activeAmountEditor) {
		pg.balanceAfterSend(true)
	}

	pg.watchForBroadcastResult(c)

	for pg.confirmButton.Button.Clicked() {
		if !pg.inputsNotEmpty(pg.passwordEditor.Editor) {
			return
		}
		pg.isBroadcastingTransaction = true
		pg.wallet.BroadcastTransaction(pg.txAuthor, []byte(pg.passwordEditor.Editor.Text()), pg.broadcastErrChan)
	}

	for pg.nextButton.Button.Clicked() {
		if pg.validate() && pg.calculateErrorText == "" {
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
		pg.isBroadcastingTransaction = false
	default:
	}
}

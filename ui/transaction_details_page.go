package ui

import (
	"fmt"
	"image/color"
	"strings"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PageTransactionDetails = "TransactionDetails"

type transactionDetailsPage struct {
	theme                           *decredmaterial.Theme
	common                          *pageCommon
	transactionDetailsPageContainer layout.List
	transactionInputsContainer      layout.List
	transactionOutputsContainer     layout.List
	hashClickable                   *widget.Clickable
	destAddressClickable            *widget.Clickable
	copyTextBtn                     []decredmaterial.Button
	dot                             *widget.Icon
	toDcrdata                       *widget.Clickable
	outputsCollapsible              *decredmaterial.Collapsible
	inputsCollapsible               *decredmaterial.Collapsible
	gtx                             *layout.Context

	transaction *dcrlibwallet.Transaction
	wallet      *dcrlibwallet.Wallet

	txSourceAccount      string
	txDestinationAddress string
}

func TransactionDetailsPage(common *pageCommon, transaction *dcrlibwallet.Transaction) Page {
	pg := &transactionDetailsPage{
		transactionDetailsPageContainer: layout.List{
			Axis: layout.Vertical,
		},
		transactionInputsContainer: layout.List{
			Axis: layout.Vertical,
		},
		transactionOutputsContainer: layout.List{
			Axis: layout.Vertical,
		},

		theme:  common.theme,
		common: common,

		outputsCollapsible: common.theme.Collapsible(),
		inputsCollapsible:  common.theme.Collapsible(),

		hashClickable:        new(widget.Clickable),
		destAddressClickable: new(widget.Clickable),
		toDcrdata:            new(widget.Clickable),

		transaction: transaction,
		wallet:      common.multiWallet.WalletWithID(transaction.WalletID),
	}

	pg.copyTextBtn = make([]decredmaterial.Button, 0)

	pg.dot = common.icons.imageBrightness1
	pg.dot.Color = common.theme.Color.Gray

	// find source account
	if transaction.Direction == dcrlibwallet.TxDirectionSent ||
		transaction.Direction == dcrlibwallet.TxDirectionTransferred {
		for _, input := range transaction.Inputs {
			if input.AccountNumber != -1 {
				accountName, err := pg.wallet.AccountName(input.AccountNumber)
				if err != nil {
					log.Error(err)
				} else {
					pg.txSourceAccount = accountName
				}
			}
		}
	}

	//	find destination address
	if transaction.Direction == dcrlibwallet.TxDirectionSent {
		for _, output := range transaction.Outputs {
			if output.AccountNumber == -1 {
				pg.txDestinationAddress = output.Address
			}
		}
	}

	return pg
}

func (pg *transactionDetailsPage) OnResume() {

}

func (pg *transactionDetailsPage) Layout(gtx layout.Context) layout.Dimensions {
	common := pg.common
	if pg.gtx == nil {
		pg.gtx = &gtx
	}
	body := func(gtx C) D {
		page := SubPage{
			title: dcrlibwallet.TransactionDirectionName(pg.transaction.Direction),
			back: func() {
				common.changePage(*common.returnPage)
			},
			body: func(gtx layout.Context) layout.Dimensions {
				widgets := []func(gtx C) D{
					func(gtx C) D {
						return pg.txnBalanceAndStatus(gtx, common)
					},
					func(gtx C) D {
						return pg.separator(gtx)
					},
					func(gtx C) D {
						return pg.txnTypeAndID(gtx)
					},
					func(gtx C) D {
						return pg.separator(gtx)
					},
					func(gtx C) D {
						return pg.txnInputs(gtx)
					},
					func(gtx C) D {
						return pg.separator(gtx)
					},
					func(gtx C) D {
						return pg.txnOutputs(gtx, common)
					},
					func(gtx C) D {
						return pg.separator(gtx)
					},
					func(gtx C) D {
						return pg.viewTxn(gtx, common)
					},
				}
				return common.theme.Card().Layout(gtx, func(gtx C) D {
					return pg.transactionDetailsPageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
						return layout.Inset{}.Layout(gtx, widgets[i])
					})
				})
			},
			infoTemplate: TransactionDetailsInfoTemplate,
		}
		return common.SubPageLayout(gtx, page)
	}

	return common.UniformPadding(gtx, body)
}

func (pg *transactionDetailsPage) txnBalanceAndStatus(gtx layout.Context, common *pageCommon) layout.Dimensions {
	txnWidgets := initTxnWidgets(common, pg.transaction)
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Right: values.MarginPadding15,
					Top:   values.MarginPadding10,
				}.Layout(gtx, txnWidgets.direction.Layout)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						amount := strings.Split(dcrutil.Amount(pg.transaction.Amount).String(), " ")
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, common.theme.H4(amount[0]).Layout)
							}),
							layout.Rigid(common.theme.H6(amount[1]).Layout),
						)
					}),
					layout.Rigid(func(gtx C) D {
						m := values.MarginPadding10
						return layout.Inset{
							Top:    m,
							Bottom: m,
						}.Layout(gtx, func(gtx C) D {
							txnWidgets.time.Color = common.theme.Color.Gray
							return txnWidgets.time.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Right: values.MarginPadding5,
									Top:   values.MarginPadding2,
								}.Layout(gtx, txnWidgets.statusIcon.Layout)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								txt := common.theme.Body1("")
								if pg.txConfirmations() > 1 {
									txt.Text = strings.Title("confirmed")
									txt.Color = common.theme.Color.Success
								} else {
									txt.Color = common.theme.Color.Gray
								}
								return txt.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								m := values.MarginPadding10
								return layout.Inset{
									Left:  m,
									Right: m,
									Top:   m,
								}.Layout(gtx, func(gtx C) D {
									return pg.dot.Layout(gtx, values.MarginPadding2)
								})
							}),
							layout.Rigid(func(gtx C) D {
								txt := common.theme.Body1(values.StringF(values.StrNConfirmations, pg.txConfirmations()))
								txt.Color = common.theme.Color.Gray
								return txt.Layout(gtx)
							}),
						)
					}),
				)
			}),
		)
	})
}

//TODO: do this at startup
func (pg *transactionDetailsPage) txConfirmations() int32 {
	transaction := pg.transaction
	if transaction.BlockHeight != -1 {
		return (pg.common.multiWallet.WalletWithID(transaction.WalletID).GetBestBlock() - transaction.BlockHeight) + 1
	}

	return 0
}

func (pg *transactionDetailsPage) txnTypeAndID(gtx layout.Context) layout.Dimensions {
	transaction := pg.transaction
	return pg.pageSections(gtx, func(gtx C) D {
		m := values.MarginPadding10
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.txnInfoSection(gtx, values.String(values.StrFrom), pg.txSourceAccount, true, nil)
			}),
			layout.Rigid(func(gtx C) D {
				if transaction.Direction == dcrlibwallet.TxDirectionSent {
					return layout.Inset{Top: m}.Layout(gtx, func(gtx C) D {
						return pg.txnInfoSection(gtx, values.String(values.StrTo), pg.txDestinationAddress, false, pg.destAddressClickable)
					})
				}
				return layout.Dimensions{}
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: m, Top: m}.Layout(gtx, func(gtx C) D {
					return pg.txnInfoSection(gtx, values.String(values.StrFee), dcrutil.Amount(transaction.Fee).String(), false, nil)
				})
			}),
			layout.Rigid(func(gtx C) D {
				if transaction.BlockHeight != -1 {
					return pg.txnInfoSection(gtx, values.String(values.StrIncludedInBlock), fmt.Sprintf("%d", transaction.BlockHeight), false, nil)
				}
				return layout.Dimensions{}
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: m, Top: m}.Layout(gtx, func(gtx C) D {
					return pg.txnInfoSection(gtx, values.String(values.StrType), transaction.Type, false, nil)
				})
			}),
			layout.Rigid(func(gtx C) D {
				trimmedHash := transaction.Hash[:24] + "..." + transaction.Hash[len(transaction.Hash)-24:]
				return layout.Inset{Bottom: m}.Layout(gtx, func(gtx C) D {
					return pg.txnInfoSection(gtx, values.String(values.StrTransactionID), trimmedHash, false, pg.hashClickable)
				})
			}),
		)
	})
}

func (pg *transactionDetailsPage) txnInfoSection(gtx layout.Context, label, value string, showWalletBadge bool, clickable *widget.Clickable) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			t := pg.theme.Body1(label)
			t.Color = pg.theme.Color.Gray
			return t.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if showWalletBadge {
						card := pg.theme.Card()
						card.Radius = decredmaterial.CornerRadius{
							NE: 0,
							NW: 0,
							SE: 0,
							SW: 0,
						}
						card.Color = pg.theme.Color.LightGray
						return card.Layout(gtx, func(gtx C) D {
							return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
								txt := pg.theme.Body2(pg.wallet.Name)
								txt.Color = pg.theme.Color.Gray
								return txt.Layout(gtx)
							})
						})
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						if clickable == nil {
							txt := pg.theme.Body1(value)
							return txt.Layout(gtx)
						}

						btn := pg.theme.Button(clickable, value)
						btn.Color = pg.theme.Color.Primary
						btn.Background = color.NRGBA{}
						btn.Inset = layout.UniformInset(values.MarginPadding0)
						return btn.Layout(gtx)
					})
				}),
			)
		}),
	)
}

func (pg *transactionDetailsPage) txnInputs(gtx layout.Context) layout.Dimensions {
	transaction := pg.transaction
	x := len(transaction.Inputs) + len(transaction.Outputs)
	for i := 0; i < x; i++ {
		pg.copyTextBtn = append(pg.copyTextBtn, pg.theme.Button(new(widget.Clickable), ""))
	}

	collapsibleHeader := func(gtx C) D {
		t := pg.theme.Body1(values.StringF(values.StrXInputsConsumed, len(transaction.Inputs)))
		t.Color = pg.theme.Color.Gray
		return t.Layout(gtx)
	}

	collapsibleBody := func(gtx C) D {
		return pg.transactionInputsContainer.Layout(gtx, len(transaction.Inputs), func(gtx C, i int) D {
			input := transaction.Inputs[i]
			return pg.txnIORow(gtx, input.Amount, input.AccountNumber, input.PreviousOutpoint, i)
		})
	}
	return pg.pageSections(gtx, func(gtx C) D {
		return pg.inputsCollapsible.Layout(gtx, collapsibleHeader, collapsibleBody)
	})
}

func (pg *transactionDetailsPage) txnOutputs(gtx layout.Context, common *pageCommon) layout.Dimensions {
	transaction := pg.transaction

	collapsibleHeader := func(gtx C) D {
		t := common.theme.Body1(values.StringF(values.StrXOutputCreated, len(transaction.Outputs)))
		t.Color = common.theme.Color.Gray
		return t.Layout(gtx)
	}

	collapsibleBody := func(gtx C) D {
		return pg.transactionOutputsContainer.Layout(gtx, len(transaction.Outputs), func(gtx C, i int) D {
			output := transaction.Outputs[i]
			x := len(transaction.Inputs)
			return pg.txnIORow(gtx, output.Amount, output.AccountNumber, output.Address, i+x)
		})
	}
	return pg.pageSections(gtx, func(gtx C) D {
		return pg.outputsCollapsible.Layout(gtx, collapsibleHeader, collapsibleBody)
	})
}

func (pg *transactionDetailsPage) txnIORow(gtx layout.Context, amount int64, acctNum int32, address string, i int) layout.Dimensions {

	accountName := "external"
	walletName := ""
	if acctNum != -1 {
		name, err := pg.wallet.AccountName(acctNum)
		if err == nil {
			accountName = name
			walletName = pg.wallet.Name
		}
	}

	accountName = fmt.Sprintf("(%s)", accountName)
	amt := dcrutil.Amount(amount).String()

	return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
		card := pg.theme.Card()
		card.Color = pg.theme.Color.LightGray
		return card.Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(pg.theme.Body1(amt).Layout),
							layout.Rigid(func(gtx C) D {
								m := values.MarginPadding5
								return layout.Inset{
									Left:  m,
									Right: m,
								}.Layout(gtx, pg.theme.Body1(accountName).Layout)
							}),
							layout.Rigid(func(gtx C) D {
								card := pg.theme.Card()
								card.Radius = decredmaterial.CornerRadius{
									NE: 0,
									NW: 0,
									SE: 0,
									SW: 0,
								}
								card.Color = pg.theme.Color.LightGray
								return card.Layout(gtx, func(gtx C) D {
									return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
										txt := pg.theme.Body2(walletName)
										txt.Color = pg.theme.Color.Gray
										return txt.Layout(gtx)
									})
								})
							}),
						)
					}),
					layout.Rigid(func(gtx C) D {
						pg.copyTextBtn[i].Color = pg.theme.Color.Primary
						pg.copyTextBtn[i].Background = color.NRGBA{}
						pg.copyTextBtn[i].Text = address
						pg.copyTextBtn[i].Inset = layout.UniformInset(values.MarginPadding0)

						return layout.W.Layout(gtx, pg.copyTextBtn[i].Layout)
					}),
				)
			})
		})
	})
}

func (pg *transactionDetailsPage) viewTxn(gtx layout.Context, common *pageCommon) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(pg.theme.Body1(values.String(values.StrViewOnDcrdata)).Layout),
			layout.Rigid(func(gtx C) D {
				redirect := common.icons.redirectIcon
				redirect.Scale = 1.0
				return decredmaterial.Clickable(gtx, pg.toDcrdata, redirect.Layout)
			}),
		)
	})
}

func (pg *transactionDetailsPage) pageSections(gtx layout.Context, body layout.Widget) layout.Dimensions {
	m := values.MarginPadding20
	mtb := values.MarginPadding5
	return layout.Inset{Left: m, Right: m, Top: mtb, Bottom: mtb}.Layout(gtx, body)
}

func (pg *transactionDetailsPage) separator(gtx layout.Context) layout.Dimensions {
	m := values.MarginPadding5
	return layout.Inset{Top: m, Bottom: m}.Layout(gtx, pg.theme.Separator().Layout)
}

func (pg *transactionDetailsPage) handle() {
	common := pg.common
	gtx := pg.gtx
	if pg.toDcrdata.Clicked() {
		goToURL(common.wallet.GetBlockExplorerURL(pg.transaction.Hash))
	}

	for _, b := range pg.copyTextBtn {
		for b.Button.Clicked() {
			clipboard.WriteOp{Text: b.Text}.Add(gtx.Ops)
		}
	}

	for pg.hashClickable.Clicked() {
		clipboard.WriteOp{Text: pg.transaction.Hash}.Add(gtx.Ops)
	}

	for pg.destAddressClickable.Clicked() {
		clipboard.WriteOp{Text: pg.txDestinationAddress}.Add(gtx.Ops)
	}
}

func (pg *transactionDetailsPage) onClose() {}

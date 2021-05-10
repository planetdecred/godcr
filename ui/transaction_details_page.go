package ui

import (
	"fmt"
	"image/color"
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/utils"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageTransactionDetails = "TransactionDetails"

type transactionDetailsPage struct {
	theme                           *decredmaterial.Theme
	transactionDetailsPageContainer layout.List
	transactionInputsContainer      layout.List
	transactionOutputsContainer     layout.List
	txnInfo                         **wallet.Transaction
	hashBtn                         decredmaterial.Button
	copyTextBtn                     []decredmaterial.Button
	dot                             *widget.Icon
	toDcrdata                       *widget.Clickable
	outputsCollapsible              *decredmaterial.Collapsible
	inputsCollapsible               *decredmaterial.Collapsible
}

func (win *Window) TransactionDetailsPage(common pageCommon) layout.Widget {
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

		txnInfo: &win.walletTransaction,
		theme:   common.theme,

		outputsCollapsible: common.theme.Collapsible(),
		inputsCollapsible:  common.theme.Collapsible(),

		hashBtn:   common.theme.Button(new(widget.Clickable), ""),
		toDcrdata: new(widget.Clickable),
	}

	pg.copyTextBtn = make([]decredmaterial.Button, 0)

	pg.dot = common.icons.imageBrightness1
	pg.dot.Color = common.theme.Color.Gray

	return func(gtx C) D {
		pg.Handler(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *transactionDetailsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	body := func(gtx C) D {
		page := SubPage{
			title: dcrlibwallet.TransactionDirectionName((*pg.txnInfo).Txn.Direction),
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
						return pg.txnOutputs(gtx, &common)
					},
					func(gtx C) D {
						return pg.separator(gtx)
					},
					func(gtx C) D {
						if *pg.txnInfo == nil {
							return layout.Dimensions{}
						}
						return pg.viewTxn(gtx, &common)
					},
				}
				return common.theme.Card().Layout(gtx, func(gtx C) D {
					if *pg.txnInfo == nil {
						return layout.Dimensions{}
					}
					return pg.transactionDetailsPageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
						return layout.Inset{}.Layout(gtx, widgets[i])
					})
				})
			},
			infoTemplate: TransactionDetailsInfoTemplate,
		}
		return common.SubPageLayout(gtx, page)
	}

	return common.Layout(gtx, func(gtx C) D {
		return common.UniformPadding(gtx, body)
	})
}

func (pg *transactionDetailsPage) txnBalanceAndStatus(gtx layout.Context, common pageCommon) layout.Dimensions {
	txnWidgets := initTxnWidgets(common, **pg.txnInfo)
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Right: values.MarginPadding15,
					Top:   values.MarginPadding10,
				}.Layout(gtx, func(gtx C) D {
					return txnWidgets.direction.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						amount := strings.Split((*pg.txnInfo).Balance, " ")
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
									return common.theme.H4(amount[0]).Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return common.theme.H6(amount[1]).Layout(gtx)
							}),
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
								}.Layout(gtx, func(gtx C) D {
									return txnWidgets.statusIcon.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								txt := common.theme.Body1("")
								if (*pg.txnInfo).Status == "confirmed" {
									txt.Text = strings.Title(strings.ToLower((*pg.txnInfo).Status))
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
								txt := common.theme.Body1(fmt.Sprintf("%d Confirmations", (*pg.txnInfo).Confirmations))
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

func (pg *transactionDetailsPage) txnTypeAndID(gtx layout.Context) layout.Dimensions {
	transaction := *pg.txnInfo
	return pg.pageSections(gtx, func(gtx C) D {
		m := values.MarginPadding10
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.txnInfoSection(gtx, "From", transaction.WalletName, transaction.AccountName, true, false)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: m, Top: m}.Layout(gtx, func(gtx C) D {
					return pg.txnInfoSection(gtx, "Fee", "", dcrutil.Amount(transaction.Txn.Fee).String(), false, false)
				})
			}),
			layout.Rigid(func(gtx C) D {
				if transaction.Txn.BlockHeight != -1 {
					return pg.txnInfoSection(gtx, "Included in block", "", fmt.Sprintf("%d", transaction.Txn.BlockHeight), false, false)
				}
				return layout.Dimensions{}
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: m, Top: m}.Layout(gtx, func(gtx C) D {
					return pg.txnInfoSection(gtx, "Type", "", transaction.Txn.Type, false, false)
				})
			}),
			layout.Rigid(func(gtx C) D {
				trimmedHash := transaction.Txn.Hash[:24] + "..." + transaction.Txn.Hash[len(transaction.Txn.Hash)-24:]
				return layout.Inset{Bottom: m}.Layout(gtx, func(gtx C) D {
					return pg.txnInfoSection(gtx, "Transaction ID", "", trimmedHash, false, true)
				})
			}),
		)
	})
}

func (pg *transactionDetailsPage) txnInfoSection(gtx layout.Context, t1, t2, t3 string, first, copy bool) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			t := pg.theme.Body1(t1)
			t.Color = pg.theme.Color.Gray
			return t.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if t2 != "" {
						if first {
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
									txt := pg.theme.Body2(strings.Title(strings.ToLower(t2)))
									txt.Color = pg.theme.Color.Gray
									return txt.Layout(gtx)
								})
							})
						}
						return pg.theme.Body1(t2).Layout(gtx)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						if first || !copy {
							txt := pg.theme.Body1(strings.Title(strings.ToLower(t3)))
							return txt.Layout(gtx)
						}

						pg.hashBtn.Color = pg.theme.Color.Primary
						pg.hashBtn.Background = color.NRGBA{}
						pg.hashBtn.Text = t3
						pg.hashBtn.Inset = layout.UniformInset(values.MarginPadding0)
						return pg.hashBtn.Layout(gtx)
					})
				}),
			)
		}),
	)
}

func (pg *transactionDetailsPage) txnInputs(gtx layout.Context) layout.Dimensions {
	transaction := *pg.txnInfo
	x := len(transaction.Txn.Inputs) + len(transaction.Txn.Outputs)
	for i := 0; i < x; i++ {
		pg.copyTextBtn = append(pg.copyTextBtn, pg.theme.Button(new(widget.Clickable), ""))
	}

	collapsibleHeader := func(gtx C) D {
		t := pg.theme.Body1(fmt.Sprintf("%d Inputs consumed", len(transaction.Txn.Inputs)))
		t.Color = pg.theme.Color.Gray
		return t.Layout(gtx)
	}

	collapsibleBody := func(gtx C) D {
		return pg.transactionInputsContainer.Layout(gtx, len(transaction.Txn.Inputs), func(gtx C, i int) D {
			amount := dcrutil.Amount(transaction.Txn.Inputs[i].Amount).String()
			acctName := fmt.Sprintf("(%s)", transaction.AccountName)
			walName := transaction.WalletName
			hashAcct := transaction.Txn.Inputs[i].PreviousOutpoint
			return pg.txnIORow(gtx, amount, acctName, walName, hashAcct, i)
		})
	}
	return pg.pageSections(gtx, func(gtx C) D {
		return pg.inputsCollapsible.Layout(gtx, collapsibleHeader, collapsibleBody)
	})
}

func (pg *transactionDetailsPage) txnOutputs(gtx layout.Context, common *pageCommon) layout.Dimensions {
	transaction := *pg.txnInfo

	collapsibleHeader := func(gtx C) D {
		t := common.theme.Body1(fmt.Sprintf("%d Outputs created", len(transaction.Txn.Outputs)))
		t.Color = common.theme.Color.Gray
		return t.Layout(gtx)
	}

	collapsibleBody := func(gtx C) D {
		return pg.transactionOutputsContainer.Layout(gtx, len(transaction.Txn.Outputs), func(gtx C, i int) D {
			amount := dcrutil.Amount(transaction.Txn.Outputs[i].Amount).String()
			acctName := fmt.Sprintf("(%s)", transaction.AccountName)
			walName := transaction.WalletName
			hashAcct := transaction.Txn.Outputs[i].Address
			x := len(transaction.Txn.Inputs)
			return pg.txnIORow(gtx, amount, acctName, walName, hashAcct, i+x)
		})
	}
	return pg.pageSections(gtx, func(gtx C) D {
		return pg.outputsCollapsible.Layout(gtx, collapsibleHeader, collapsibleBody)
	})
}

func (pg *transactionDetailsPage) txnIORow(gtx layout.Context, amount, acctName, walName, hashAcct string, i int) layout.Dimensions {
	return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
		card := pg.theme.Card()
		card.Color = pg.theme.Color.LightGray
		return card.Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return pg.theme.Body1(amount).Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								m := values.MarginPadding5
								return layout.Inset{
									Left:  m,
									Right: m,
								}.Layout(gtx, func(gtx C) D {
									return pg.theme.Body1(acctName).Layout(gtx)
								})
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
										txt := pg.theme.Body2(walName)
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
						pg.copyTextBtn[i].Text = hashAcct
						pg.copyTextBtn[i].Inset = layout.UniformInset(values.MarginPadding0)

						return layout.W.Layout(gtx, func(gtx C) D {
							return pg.copyTextBtn[i].Layout(gtx)
						})
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
			layout.Rigid(func(gtx C) D {
				return pg.theme.Body1("View on dcrdata").Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				redirect := common.icons.redirectIcon
				redirect.Scale = 1.0
				return decredmaterial.Clickable(gtx, pg.toDcrdata, func(gtx C) D {
					return redirect.Layout(gtx)
				})
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
	return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
		return pg.theme.Separator().Layout(gtx)
	})
}

func (pg *transactionDetailsPage) Handler(common pageCommon) {
	if pg.toDcrdata.Clicked() {
		utils.GoToURL(common.wallet.GetBlockExplorerURL((*pg.txnInfo).Txn.Hash))
	}

	for _, b := range pg.copyTextBtn {
		for b.Button.Clicked() {
			t := b.Text
			go func() {
				common.clipboard <- WriteClipboard{Text: t}
			}()
		}
	}

	for pg.hashBtn.Button.Clicked() {
		go func() {
			common.clipboard <- WriteClipboard{Text: (*pg.txnInfo).Txn.Hash}
		}()
	}
}

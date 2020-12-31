package ui

import (
	"fmt"
	"image/color"
	"os/exec"
	"runtime"
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageTransactionDetails = "TransactionDetails"

type transactionDetailsPage struct {
	theme                           *decredmaterial.Theme
	transactionDetailsPageContainer layout.List
	transactionInputsContainer      layout.List
	transactionOutputsContainer     layout.List
	autoCompleteList                *layout.List
	backButton                      decredmaterial.IconButton
	txnInfo                         **wallet.Transaction
	minInfoBtn                      decredmaterial.Button
	hashBtn                         decredmaterial.Button
	copyTextBtn                     []decredmaterial.Button
	infoBtn                         decredmaterial.IconButton
	dot                             *widget.Icon
	toDcrdata                       *widget.Clickable
	outputsCollapsible              *decredmaterial.Collapsible
	inputsCollapsible               *decredmaterial.Collapsible
	line                            *decredmaterial.Line
	infoModal                       *decredmaterial.Modal
	showInfo                        bool
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
		autoCompleteList: &layout.List{Axis: layout.Horizontal},

		txnInfo:  &win.walletTransaction,
		theme:    common.theme,
		showInfo: false,

		outputsCollapsible: common.theme.Collapsible(),
		inputsCollapsible:  common.theme.Collapsible(),

		backButton: common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack),
		minInfoBtn: common.theme.Button(new(widget.Clickable), "Got it"),
		hashBtn:    common.theme.Button(new(widget.Clickable), ""),
		toDcrdata:  new(widget.Clickable),
		line:       common.theme.Line(),
		infoModal:  common.theme.Modal(),
	}

	pg.copyTextBtn = make([]decredmaterial.Button, 0)

	pg.line.Color = common.theme.Color.Background
	pg.backButton.Color = common.theme.Color.Text
	pg.backButton.Inset = layout.UniformInset(values.MarginPadding0)
	pg.minInfoBtn.Background = color.RGBA{}
	pg.minInfoBtn.Color = common.theme.Color.Primary
	pg.minInfoBtn.TextSize = values.MarginPadding20
	pg.infoBtn = common.theme.IconButton(new(widget.Clickable), common.icons.actionInfo)
	pg.infoBtn.Color = common.theme.Color.Gray
	pg.infoBtn.Background = common.theme.Color.Surface
	pg.infoBtn.Inset = layout.UniformInset(values.MarginPadding0)
	pg.dot = common.icons.imageBrightness1
	pg.dot.Color = common.theme.Color.Gray

	return func(gtx C) D {
		pg.Handler(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *transactionDetailsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	widgets := []func(gtx C) D{
		func(gtx C) D {
			return pg.header(gtx)
		},
		func(gtx C) D {
			return pg.txnBalanceAndStatus(gtx, &common)
		},
		func(gtx C) D {
			return pg.divide(gtx)
		},
		func(gtx C) D {
			return pg.txnTypeAndID(gtx)
		},
		func(gtx C) D {
			return pg.divide(gtx)
		},
		func(gtx C) D {
			return pg.txnInputs(gtx)
		},
		func(gtx C) D {
			return pg.divide(gtx)
		},
		func(gtx C) D {
			return pg.txnOutputs(gtx, &common)
		},
		func(gtx C) D {
			return pg.divide(gtx)
		},
		func(gtx C) D {
			if *pg.txnInfo == nil {
				return layout.Dimensions{}
			}
			return pg.viewTxn(gtx, &common)
		},
	}

	body := common.Layout(gtx, func(gtx C) D {
		return common.theme.Card().Layout(gtx, func(gtx C) D {
			if *pg.txnInfo == nil {
				return layout.Dimensions{}
			}
			return pg.transactionDetailsPageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
				return layout.Inset{}.Layout(gtx, widgets[i])
			})
		})
	})

	if pg.showInfo {
		info := []func(gtx C) D{
			func(gtx C) D {
				return pg.infoModalLayout(gtx, &common)
			},
		}

		return pg.infoModal.Layout(gtx, info, 1300)
	}
	return body
}

func (pg *transactionDetailsPage) header(gtx layout.Context) layout.Dimensions {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.W.Layout(gtx, func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
							return pg.backButton.Layout(gtx)
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					txt := pg.theme.H6("")
					if *pg.txnInfo != nil {
						txt.Text = dcrlibwallet.TransactionDirectionName((*pg.txnInfo).Txn.Direction)
					} else {
						txt.Text = "Not found"
					}

					txt.Alignment = text.Middle
					return txt.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.infoBtn.Layout(gtx)
					})
				}),
			)
		})
	})
}

func (pg *transactionDetailsPage) txnBalanceAndStatus(gtx layout.Context, common *pageCommon) layout.Dimensions {
	txnWidgets := transactionWdg{}
	initTxnWidgets(common, *pg.txnInfo, &txnWidgets)

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
				return pg.txnInfoSection(gtx, "From", transaction.WalletName, transaction.Txn.Inputs[0].AccountName, true, false)
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
							card.Color = pg.theme.Color.Background
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
						pg.hashBtn.Background = color.RGBA{}
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
			acctName := fmt.Sprintf("(%s)", transaction.Txn.Inputs[i].AccountName)
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
			acctName := fmt.Sprintf("(%s)", transaction.Txn.Outputs[i].AccountName)
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
		card.Color = pg.theme.Color.Background
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
								card.Color = pg.theme.Color.Background
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
						pg.copyTextBtn[i].Background = color.RGBA{}
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
				redirect.Scale = 0.26
				return decredmaterial.Clickable(gtx, pg.toDcrdata, func(gtx C) D {
					return redirect.Layout(gtx)
				})
			}),
		)
	})
}

func (pg *transactionDetailsPage) viewTxnOnBrowser(common *pageCommon) {
	var err error
	url := common.wallet.GetBlockExplorerURL((*pg.txnInfo).Txn.Hash)

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Error(err)
	}
}

func (pg *transactionDetailsPage) infoModalLayout(gtx layout.Context, common *pageCommon) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.NW.Layout(gtx, func(gtx C) D {
				t := pg.theme.Body1("Tap on")
				t.Color = common.theme.Color.Text
				return pg.theme.H6("How to copy").Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				inset := layout.Inset{
					Top:    values.MarginPadding20,
					Bottom: values.MarginPadding30,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							t := pg.theme.Body1("Tap on")
							t.Color = common.theme.Color.Gray
							return t.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							t := pg.theme.Body1("blue text")
							t.Color = common.theme.Color.Primary
							m := values.MarginPadding2
							return layout.Inset{
								Left:  m,
								Right: m,
							}.Layout(gtx, func(gtx C) D {
								return t.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							t := pg.theme.Body1("to copy the item.")
							t.Color = common.theme.Color.Gray
							return t.Layout(gtx)
						}),
					)
				})
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.SE.Layout(gtx, func(gtx C) D {
				return pg.minInfoBtn.Layout(gtx)
			})
		}),
	)
}

func (pg *transactionDetailsPage) pageSections(gtx layout.Context, body layout.Widget) layout.Dimensions {
	m := values.MarginPadding20
	mtb := values.MarginPadding5
	return layout.Inset{Left: m, Right: m, Top: mtb, Bottom: mtb}.Layout(gtx, body)
}

func (pg *transactionDetailsPage) divide(gtx layout.Context) layout.Dimensions {
	pg.line.Width = gtx.Constraints.Max.X
	pg.line.Height = 2
	m := values.MarginPadding5
	return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
		return pg.line.Layout(gtx)
	})
}

func (pg *transactionDetailsPage) Handler(common pageCommon) {
	if pg.toDcrdata.Clicked() {
		pg.viewTxnOnBrowser(&common)
	}

	if pg.infoBtn.Button.Clicked() {
		pg.showInfo = true
	}

	if pg.minInfoBtn.Button.Clicked() {
		pg.showInfo = false
	}

	if pg.backButton.Button.Clicked() {
		*common.page = PageTransactions
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

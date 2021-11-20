package transaction

import (
	"fmt"
	"strings"
	"time"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const TransactionDetailsPageID = "TransactionDetails"

type transactionWdg struct {
	confirmationIcons    *decredmaterial.Image
	icon                 *decredmaterial.Image
	title                string
	time, status, wallet decredmaterial.Label

	copyTextButtons []decredmaterial.Button
}

type TxDetailsPage struct {
	*load.Load

	list *widget.List

	transactionDetailsPageContainer layout.List
	transactionInputsContainer      layout.List
	transactionOutputsContainer     layout.List
	associatedTicketClickable       *decredmaterial.Clickable
	hashClickable                   *widget.Clickable
	destAddressClickable            *widget.Clickable
	dot                             *decredmaterial.Icon
	toDcrdata                       *decredmaterial.Clickable
	outputsCollapsible              *decredmaterial.Collapsible
	inputsCollapsible               *decredmaterial.Collapsible
	backButton                      decredmaterial.IconButton
	infoButton                      decredmaterial.IconButton
	gtx                             *layout.Context

	txnWidgets    transactionWdg
	transaction   *dcrlibwallet.Transaction
	ticketSpender *dcrlibwallet.Transaction // vote or revoke ticket
	ticketSpent   *dcrlibwallet.Transaction // ticket spent in a vote or revoke
	wallet        *dcrlibwallet.Wallet

	txSourceAccount      string
	txDestinationAddress string
}

func NewTransactionDetailsPage(l *load.Load, transaction *dcrlibwallet.Transaction) *TxDetailsPage {
	pg := &TxDetailsPage{
		Load: l,
		list: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		transactionDetailsPageContainer: layout.List{
			Axis: layout.Vertical,
		},
		transactionInputsContainer: layout.List{
			Axis: layout.Vertical,
		},
		transactionOutputsContainer: layout.List{
			Axis: layout.Vertical,
		},

		outputsCollapsible: l.Theme.Collapsible(),
		inputsCollapsible:  l.Theme.Collapsible(),

		associatedTicketClickable: l.Theme.NewClickable(true),
		hashClickable:             new(widget.Clickable),
		destAddressClickable:      new(widget.Clickable),
		toDcrdata:                 l.Theme.NewClickable(true),

		transaction: transaction,
		wallet:      l.WL.MultiWallet.WalletWithID(transaction.WalletID),
	}

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(pg.Load)
	pg.dot = decredmaterial.NewIcon(l.Icons.ImageBrightness1)
	pg.dot.Color = l.Theme.Color.Gray

	// find source account
	if transaction.Direction == dcrlibwallet.TxDirectionSent ||
		transaction.Direction == dcrlibwallet.TxDirectionTransferred {
		for _, input := range transaction.Inputs {
			if input.AccountNumber != -1 {
				accountName, err := pg.wallet.AccountName(input.AccountNumber)
				if err != nil {
					// log.Error(err)
				} else {
					pg.txSourceAccount = accountName
				}
			}
		}
	}
	//
	//	find destination address
	if transaction.Direction == dcrlibwallet.TxDirectionSent {
		for _, output := range transaction.Outputs {
			if output.AccountNumber == -1 {
				pg.txDestinationAddress = output.Address
			}
		}
	}

	pg.txnWidgets = initTxnWidgets(pg.Load, pg.transaction)

	return pg
}

func (pg *TxDetailsPage) ID() string {
	return TransactionDetailsPageID
}

func (pg *TxDetailsPage) OnResume() {
	if pg.transaction.TicketSpentHash != "" {
		pg.ticketSpent, _ = pg.wallet.GetTransactionRaw(pg.transaction.TicketSpentHash)
	}

	if ok, _ := pg.wallet.TicketHasVotedOrRevoked(pg.transaction.Hash); ok {
		pg.ticketSpender, _ = pg.wallet.TicketSpender(pg.transaction.Hash)
	}
}

func (pg *TxDetailsPage) Layout(gtx layout.Context) layout.Dimensions {
	if pg.gtx == nil {
		pg.gtx = &gtx
	}

	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      pg.txnWidgets.title,
			BackButton: pg.backButton,
			InfoButton: pg.infoButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx layout.Context) layout.Dimensions {
				widgets := []func(gtx C) D{
					func(gtx C) D {
						return pg.txnBalanceAndStatus(gtx)
					},
					func(gtx C) D {
						return pg.Theme.Separator().Layout(gtx)
					},
					func(gtx C) D {
						return pg.ticketDetails(gtx)
					},
					func(gtx C) D {
						return pg.associatedTicket(gtx)
					},
					func(gtx C) D {
						return pg.txnTypeAndID(gtx)
					},
					func(gtx C) D {
						return pg.Theme.Separator().Layout(gtx)
					},
					func(gtx C) D {
						return pg.txnInputs(gtx)
					},
					func(gtx C) D {
						return pg.Theme.Separator().Layout(gtx)
					},
					func(gtx C) D {
						return pg.txnOutputs(gtx)
					},
					func(gtx C) D {
						return pg.Theme.Separator().Layout(gtx)
					},
					func(gtx C) D {
						return pg.viewTxn(gtx)
					},
				}
				return pg.Theme.List(pg.list).Layout(gtx, 1, func(gtx C, i int) D {
					return pg.Theme.Card().Layout(gtx, func(gtx C) D {
						return pg.transactionDetailsPageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
							return layout.Inset{}.Layout(gtx, widgets[i])
						})
					})
				})
			},
			InfoTemplate: modal.TransactionDetailsInfoTemplate,
		}
		return sp.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

func (pg *TxDetailsPage) txnBalanceAndStatus(gtx layout.Context) layout.Dimensions {
	return decredmaterial.LinearLayout{
		Width:       decredmaterial.MatchParent,
		Height:      decredmaterial.WrapContent,
		Orientation: layout.Horizontal,
		Padding:     layout.UniformInset(values.MarginPadding16),
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Right: values.MarginPadding16,
				Top:   values.MarginPadding12,
			}.Layout(gtx, pg.txnWidgets.icon.Layout24dp)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					amount := dcrutil.Amount(pg.transaction.Amount).String()
					if pg.transaction.Type == dcrlibwallet.TxTypeMixed {
						amount = dcrutil.Amount(pg.transaction.MixDenomination).String()
					} else if pg.transaction.Type == dcrlibwallet.TxTypeRegular && pg.transaction.Direction == dcrlibwallet.TxDirectionSent {
						amount = "-" + amount
					}
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return components.LayoutBalanceSize(gtx, pg.Load, amount, values.TextSize34)
						}),
						layout.Rigid(func(gtx C) D {
							if pg.transaction.Type == dcrlibwallet.TxTypeMixed && pg.transaction.MixCount > 1 {

								label := pg.Theme.H5(fmt.Sprintf("x%d", pg.transaction.MixCount))
								label.Color = pg.Theme.Color.Gray
								return layout.Inset{
									Left: values.MarginPadding8,
								}.Layout(gtx, label.Layout)
							}
							return D{}
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					m := values.MarginPadding10
					return layout.Inset{
						Top:    m,
						Bottom: m,
					}.Layout(gtx, func(gtx C) D {
						pg.txnWidgets.time.Color = pg.Theme.Color.Gray
						return pg.txnWidgets.time.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{
								Right: values.MarginPadding4,
								Top:   values.MarginPadding4,
							}.Layout(gtx, func(gtx C) D {
								return pg.txnWidgets.confirmationIcons.Layout12dp(gtx)
							})
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							txt := pg.Theme.Body1("")
							if pg.txConfirmations() > 1 {
								txt.Text = strings.Title("confirmed")
								txt.Color = pg.Theme.Color.Success
							} else {
								txt.Text = strings.Title("pending")
								txt.Color = pg.Theme.Color.Gray
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
							txt := pg.Theme.Body1(values.StringF(values.StrNConfirmations, pg.txConfirmations()))
							txt.Color = pg.Theme.Color.Gray
							return txt.Layout(gtx)
						}),
					)
				}),
			)
		}),
	)
}

func (pg *TxDetailsPage) maturityProgressBar(gtx C) D {
	return decredmaterial.LinearLayout{
		Width:       decredmaterial.MatchParent,
		Height:      decredmaterial.WrapContent,
		Orientation: layout.Horizontal,
		Margin:      layout.Inset{Top: values.MarginPadding12},
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			t := pg.Theme.Label(values.TextSize14, "Maturity")
			t.Color = pg.Theme.Color.Gray
			return t.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {

			percentageLabel := pg.Theme.Label(values.TextSize14, "25%")
			percentageLabel.Color = pg.Theme.Color.Gray

			progress := pg.Theme.ProgressBar(40)
			progress.Color = pg.Theme.Color.LightBlue
			progress.TrackColor = pg.Theme.Color.BlueProgressTint
			progress.Height = values.MarginPadding8
			progress.Width = values.MarginPadding80
			progress.Radius = decredmaterial.Radius(values.MarginPadding8.V)

			timeLeft := pg.Theme.Label(values.TextSize16, "18 hours")
			timeLeft.Color = pg.Theme.Color.DeepBlue

			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(percentageLabel.Layout),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Left: values.MarginPadding6, Right: values.MarginPadding6}.Layout(gtx, progress.Layout)
					}),
					layout.Rigid(timeLeft.Layout),
				)
			})
		}),
	)
}

func (pg *TxDetailsPage) ticketDetails(gtx C) D {
	if !pg.wallet.TxMatchesFilter(pg.transaction, dcrlibwallet.TxFilterStaking) ||
		pg.transaction.Type == dcrlibwallet.TxTypeRevocation {
		return D{}
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return decredmaterial.LinearLayout{
				Width:       decredmaterial.MatchParent,
				Height:      decredmaterial.WrapContent,
				Orientation: layout.Vertical,
				Padding:     layout.Inset{Left: values.MarginPadding16, Right: values.MarginPadding16, Bottom: values.MarginPadding12},
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if pg.transaction.Type == dcrlibwallet.TxTypeTicketPurchase {
						var status string
						if pg.ticketSpender != nil {
							if pg.ticketSpender.Type == dcrlibwallet.TxTypeVote {
								status = "Voted"
							} else {
								status = "Revoked"
							}
						} else if pg.wallet.TxMatchesFilter(pg.transaction, dcrlibwallet.TxFilterLive) {
							status = "Live"
						} else if pg.wallet.TxMatchesFilter(pg.transaction, dcrlibwallet.TxFilterImmature) {
							status = "Immature"
						} else if pg.wallet.TxMatchesFilter(pg.transaction, dcrlibwallet.TxFilterUnmined) {
							status = "Unmined"
						} else if pg.wallet.TxMatchesFilter(pg.transaction, dcrlibwallet.TxFilterExpired) {
							status = "Expired"
						} else {
							status = "Unknown"
						}

						return layout.Inset{Top: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
							return pg.txnInfoSection(gtx, "Status", status, false, nil)
						})
					}

					return D{}
				}),
				layout.Rigid(func(gtx C) D {
					// TODO spendable progress bar

					if false {
						return pg.maturityProgressBar(gtx)
					}

					return D{}
				}),
				layout.Rigid(func(gtx C) D {
					if pg.transaction.Type == dcrlibwallet.TxTypeVote {
						return layout.Inset{Top: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
							return pg.txnInfoSection(gtx, "Days to vote", fmt.Sprintf("%d days", pg.transaction.DaysToVoteOrRevoke), false, nil)
						})
					}

					return D{}
				}),
				layout.Rigid(func(gtx C) D {
					if pg.transaction.Type == dcrlibwallet.TxTypeVote {
						return layout.Inset{Top: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
							return pg.txnInfoSection(gtx, "Reward", dcrutil.Amount(pg.transaction.VoteReward).String(), false, nil)
						})
					}
					return D{}
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			return pg.Theme.Separator().Layout(gtx)
		}),
	)
}

func (pg *TxDetailsPage) associatedTicket(gtx C) D {
	if pg.transaction.Type != dcrlibwallet.TxTypeVote && pg.transaction.Type != dcrlibwallet.TxTypeRevocation {
		return D{}
	}

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.associatedTicketClickable.Layout(gtx, func(gtx C) D {
				return decredmaterial.LinearLayout{
					Width:       decredmaterial.MatchParent,
					Height:      decredmaterial.WrapContent,
					Orientation: layout.Horizontal,
					Padding:     layout.Inset{Left: values.MarginPadding16, Top: values.MarginPadding12, Right: values.MarginPadding16, Bottom: values.MarginPadding12},
				}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						label := pg.Theme.Label(values.TextSize16, "View associated ticket")
						label.Color = pg.Theme.Color.DeepBlue
						return label.Layout(gtx)
					}),
					layout.Flexed(1, func(gtx C) D {
						icon := pg.Icons.Next
						return layout.E.Layout(gtx, icon.Layout24dp)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return pg.Theme.Separator().Layout(gtx)
		}),
	)
}

//TODO: do this at startup
func (pg *TxDetailsPage) txConfirmations() int32 {
	transaction := pg.transaction
	if transaction.BlockHeight != -1 {
		return (pg.WL.MultiWallet.WalletWithID(transaction.WalletID).GetBestBlock() - transaction.BlockHeight) + 1
	}

	return 0
}

func (pg *TxDetailsPage) txnTypeAndID(gtx layout.Context) layout.Dimensions {
	transaction := pg.transaction
	m := values.MarginPadding12
	return decredmaterial.LinearLayout{
		Width:       decredmaterial.MatchParent,
		Height:      decredmaterial.WrapContent,
		Orientation: layout.Vertical,
		Padding:     layout.UniformInset(values.MarginPadding16),
	}.Layout(gtx,
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
			return layout.Inset{Top: m}.Layout(gtx, func(gtx C) D {
				return pg.txnInfoSection(gtx, values.String(values.StrFee), dcrutil.Amount(transaction.Fee).String(), false, nil)
			})
		}),
		layout.Rigid(func(gtx C) D {
			if transaction.BlockHeight != -1 {
				return layout.Inset{Top: m}.Layout(gtx, func(gtx C) D {
					return pg.txnInfoSection(gtx, values.String(values.StrIncludedInBlock), fmt.Sprintf("%d", transaction.BlockHeight), false, nil)
				})
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: m}.Layout(gtx, func(gtx C) D {
				return pg.txnInfoSection(gtx, values.String(values.StrType), transaction.Type, false, nil)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: m}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						t := pg.Theme.Label(values.TextSize14, values.String(values.StrTransactionID))
						t.Color = pg.Theme.Color.Gray
						return t.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						btn := pg.Theme.OutlineButton(transaction.Hash)
						btn.TextSize = values.MarginPadding14
						btn.SetClickable(pg.hashClickable)
						btn.Inset = layout.UniformInset(values.MarginPadding0)
						return btn.Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (pg *TxDetailsPage) txnInfoSection(gtx layout.Context, label, value string, showWalletBadge bool, clickable *widget.Clickable) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			t := pg.Theme.Label(values.TextSize14, label)
			t.Color = pg.Theme.Color.Gray
			return t.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if showWalletBadge {
						card := pg.Theme.Card()
						card.Radius = decredmaterial.Radius(0)
						card.Color = pg.Theme.Color.LightGray
						return card.Layout(gtx, func(gtx C) D {
							return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
								txt := pg.Theme.Body2(pg.wallet.Name)
								txt.Color = pg.Theme.Color.Gray
								return txt.Layout(gtx)
							})
						})
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						if clickable == nil {
							txt := pg.Theme.Body1(value)
							return txt.Layout(gtx)
						}

						btn := pg.Theme.OutlineButton(value)
						btn.TextSize = values.MarginPadding14
						btn.SetClickable(clickable)
						btn.Inset = layout.UniformInset(values.MarginPadding0)
						return btn.Layout(gtx)
					})
				}),
			)
		}),
	)
}

func (pg *TxDetailsPage) txnInputs(gtx layout.Context) layout.Dimensions {
	transaction := pg.transaction

	collapsibleHeader := func(gtx C) D {
		t := pg.Theme.Body1(values.StringF(values.StrXInputsConsumed, len(transaction.Inputs)))
		t.Color = pg.Theme.Color.Gray
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

func (pg *TxDetailsPage) txnOutputs(gtx layout.Context) layout.Dimensions {
	transaction := pg.transaction

	collapsibleHeader := func(gtx C) D {
		t := pg.Theme.Body1(values.StringF(values.StrXOutputCreated, len(transaction.Outputs)))
		t.Color = pg.Theme.Color.Gray
		return t.Layout(gtx)
	}

	collapsibleBody := func(gtx C) D {
		x := len(transaction.Inputs)
		return pg.transactionOutputsContainer.Layout(gtx, len(transaction.Outputs), func(gtx C, i int) D {
			output := transaction.Outputs[i]
			return pg.txnIORow(gtx, output.Amount, output.AccountNumber, output.Address, i+x)
		})
	}
	return pg.pageSections(gtx, func(gtx C) D {
		return pg.outputsCollapsible.Layout(gtx, collapsibleHeader, collapsibleBody)
	})
}

func (pg *TxDetailsPage) txnIORow(gtx layout.Context, amount int64, acctNum int32, address string, i int) layout.Dimensions {

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

	return layout.Inset{Top: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
		card := pg.Theme.Card()
		card.Color = pg.Theme.Color.LightGray
		return card.Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(pg.Theme.Body1(amt).Layout),
							layout.Rigid(func(gtx C) D {
								m := values.MarginPadding5
								return layout.Inset{
									Left:  m,
									Right: m,
								}.Layout(gtx, pg.Theme.Body1(accountName).Layout)
							}),
							layout.Rigid(func(gtx C) D {
								card := pg.Theme.Card()
								card.Radius = decredmaterial.Radius(0)
								card.Color = pg.Theme.Color.LightGray
								return card.Layout(gtx, func(gtx C) D {
									return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
										txt := pg.Theme.Body2(walletName)
										txt.Color = pg.Theme.Color.Gray
										return txt.Layout(gtx)
									})
								})
							}),
						)
					}),
					layout.Rigid(func(gtx C) D {
						pg.txnWidgets.copyTextButtons[i].Text = address

						return layout.W.Layout(gtx, pg.txnWidgets.copyTextButtons[i].Layout)
					}),
				)
			})
		})
	})
}

func (pg *TxDetailsPage) viewTxn(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(pg.Theme.Body1(values.String(values.StrViewOnDcrdata)).Layout),
			layout.Rigid(func(gtx C) D {
				redirect := pg.Icons.RedirectIcon
				return pg.toDcrdata.Layout(gtx, redirect.Layout24dp)
			}),
		)
	})
}

func (pg *TxDetailsPage) pageSections(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, body)
}

func (pg *TxDetailsPage) Handle() {
	gtx := pg.gtx
	if pg.toDcrdata.Clicked() {
		components.GoToURL(pg.WL.Wallet.GetBlockExplorerURL(pg.transaction.Hash))
	}

	for _, b := range pg.txnWidgets.copyTextButtons {
		for b.Clicked() {
			clipboard.WriteOp{Text: b.Text}.Add(gtx.Ops)
			pg.Toast.Notify("Copied")
		}
	}

	for pg.hashClickable.Clicked() {
		clipboard.WriteOp{Text: pg.transaction.Hash}.Add(gtx.Ops)
		pg.Toast.Notify("Transaction Hash copied")
	}

	for pg.destAddressClickable.Clicked() {
		clipboard.WriteOp{Text: pg.txDestinationAddress}.Add(gtx.Ops)
		pg.Toast.Notify("Address copied")
	}

	for pg.associatedTicketClickable.Clicked() {
		if pg.ticketSpent != nil {
			pg.ChangeFragment(NewTransactionDetailsPage(pg.Load, pg.ticketSpent))
		}
	}
}

func (pg *TxDetailsPage) OnClose() {}

func initTxnWidgets(l *load.Load, transaction *dcrlibwallet.Transaction) transactionWdg {

	var txn transactionWdg
	wal := l.WL.MultiWallet.WalletWithID(transaction.WalletID)

	t := time.Unix(transaction.Timestamp, 0).UTC()
	txn.time = l.Theme.Body1(t.Format(time.UnixDate))
	txn.status = l.Theme.Body1("")
	txn.wallet = l.Theme.Body2(wal.Name)

	if components.TxConfirmations(l, *transaction) > 1 {
		txn.status.Text = components.FormatDateOrTime(transaction.Timestamp)
		txn.confirmationIcons = l.Icons.ConfirmIcon
	} else {
		txn.status.Text = "pending"
		txn.status.Color = l.Theme.Color.Gray
		txn.confirmationIcons = l.Icons.PendingIcon
	}

	var ticketSpender *dcrlibwallet.Transaction
	if wal.TxMatchesFilter(transaction, dcrlibwallet.TxFilterStaking) {
		ticketSpender, _ = wal.TicketSpender(transaction.Hash)
	}
	txStatus := components.TransactionTitleIcon(l, wal, transaction, ticketSpender)

	txn.title = txStatus.Title
	txn.icon = txStatus.Icon

	x := len(transaction.Inputs) + len(transaction.Outputs)
	txn.copyTextButtons = make([]decredmaterial.Button, x)
	for i := 0; i < x; i++ {
		btn := l.Theme.OutlineButton("")
		btn.TextSize = values.MarginPadding14
		btn.Inset = layout.UniformInset(values.MarginPadding0)
		txn.copyTextButtons[i] = btn
	}

	return txn
}

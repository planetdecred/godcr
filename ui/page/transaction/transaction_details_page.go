package transaction

import (
	"fmt"
	"gioui.org/op"
	"strings"
	"time"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	TransactionDetailsPageID = "TransactionDetails"
	viewBlockID              = "viewBlock"
	copyBlockID              = "copyBlock"
)

type transactionWdg struct {
	confirmationIcons    *decredmaterial.Image
	time, status, wallet decredmaterial.Label

	copyTextButtons []decredmaterial.Button
	txStatus        *components.TxStatus
}

type moreItem struct {
	text   string
	id     string
	button *decredmaterial.Clickable
}

type TxDetailsPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	list *widget.List

	transactionDetailsPageContainer layout.List
	transactionInputsContainer      layout.List
	transactionOutputsContainer     layout.List
	associatedTicketClickable       *decredmaterial.Clickable
	hashClickable                   *widget.Clickable
	destAddressClickable            *widget.Clickable
	dot                             *decredmaterial.Icon
	outputsCollapsible              *decredmaterial.Collapsible
	inputsCollapsible               *decredmaterial.Collapsible
	backButton                      decredmaterial.IconButton
	rebroadcast                     decredmaterial.Label
	rebroadcastClickable            *decredmaterial.Clickable
	rebroadcastIcon                 *decredmaterial.Image
	moreOption                      *decredmaterial.Clickable

	shadowBox *decredmaterial.Shadow

	txnWidgets    transactionWdg
	transaction   *dcrlibwallet.Transaction
	ticketSpender *dcrlibwallet.Transaction // vote or revoke ticket
	ticketSpent   *dcrlibwallet.Transaction // ticket spent in a vote or revoke
	txBackStack   *dcrlibwallet.Transaction // track original transaction
	wallet        *dcrlibwallet.Wallet

	moreItems []moreItem

	txSourceAccount      string
	txDestinationAddress string
	title                string

	moreOptionIsOpen bool
}

func NewTransactionDetailsPage(l *load.Load, transaction *dcrlibwallet.Transaction, isTicket bool) *TxDetailsPage {
	rebroadcast := l.Theme.Label(values.TextSize14, values.String(values.StrRebroadcast))
	rebroadcast.TextSize = values.TextSize14
	rebroadcast.Color = l.Theme.Color.Text
	pg := &TxDetailsPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(TransactionDetailsPageID),
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
		moreOption:                l.Theme.NewClickable(false),
		shadowBox:                 l.Theme.Shadow(),

		transaction:          transaction,
		wallet:               l.WL.MultiWallet.WalletWithID(transaction.WalletID),
		rebroadcast:          rebroadcast,
		rebroadcastClickable: l.Theme.NewClickable(true),
		rebroadcastIcon:      l.Theme.Icons.Rebroadcast,
	}

	pg.backButton, _ = components.SubpageHeaderButtons(pg.Load)

	pg.dot = decredmaterial.NewIcon(l.Theme.Icons.ImageBrightness1)
	pg.dot.Color = l.Theme.Color.Gray1

	pg.moreItems = pg.getMoreItem()

	return pg
}

func (pg *TxDetailsPage) getTXSourceAccountAndDirection() {
	// find source account
	if pg.transaction.Direction == dcrlibwallet.TxDirectionSent ||
		pg.transaction.Direction == dcrlibwallet.TxDirectionTransferred {
		for _, input := range pg.transaction.Inputs {
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
	//	find destination address
	if pg.transaction.Direction == dcrlibwallet.TxDirectionSent {
		for _, output := range pg.transaction.Outputs {
			if output.AccountNumber == -1 {
				pg.txDestinationAddress = output.Address
			}
		}
	}
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *TxDetailsPage) OnNavigatedTo() {
	if pg.transaction.TicketSpentHash != "" {
		pg.ticketSpent, _ = pg.wallet.GetTransactionRaw(pg.transaction.TicketSpentHash)
	}

	if ok, _ := pg.wallet.TicketHasVotedOrRevoked(pg.transaction.Hash); ok {
		pg.ticketSpender, _ = pg.wallet.TicketSpender(pg.transaction.Hash)
	}

	pg.title = "Transaction Details"
	if pg.transaction.Type == "Ticket" {
		pg.title = "Ticket Details"
	}

	pg.getTXSourceAccountAndDirection()
	pg.txnWidgets = initTxnWidgets(pg.Load, pg.transaction)
}

func (pg *TxDetailsPage) getMoreItem() []moreItem {
	return []moreItem{
		{
			text:   values.String(values.StrViewOnExplorer),
			button: pg.Theme.NewClickable(true),
			id:     viewBlockID,
		},
		{
			text:   values.String(values.StrCopyBlockLink),
			button: pg.Theme.NewClickable(true),
			id:     copyBlockID,
		},
	}
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *TxDetailsPage) Layout(gtx C) D {
	pg.handleTextCopyEvent(gtx)

	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      pg.title,
			BackButton: pg.backButton,
			ExtraItem:  pg.moreOption,
			Extra: func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(pg.Theme.Icons.EllipseHoriz.Layout24dp),
						layout.Rigid(func(gtx C) D {
							if pg.moreOptionIsOpen {
								pg.layoutOptionsMenu(gtx)
							}
							return D{}
						}),
					)
				})
			},
			Back: func() {
				if pg.txBackStack == nil {
					pg.ParentNavigator().CloseCurrentPage()
					return
				}
				pg.transaction = pg.txBackStack
				pg.getTXSourceAccountAndDirection()
				pg.txnWidgets = initTxnWidgets(pg.Load, pg.transaction)
				pg.txBackStack = nil
				pg.ParentWindow().Reload()
			},
			Body: func(gtx C) D {
				widgets := []func(gtx C) D{
					pg.txDetailsHeader,
					func(gtx C) D {
						return pg.Theme.Separator().Layout(gtx)
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
				}
				return pg.Theme.List(pg.list).Layout(gtx, 1, func(gtx C, i int) D {
					return pg.Theme.Card().Layout(gtx, func(gtx C) D {
						return pg.transactionDetailsPageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
							return layout.Inset{}.Layout(gtx, widgets[i])
						})
					})
				})
			},
		}
		return sp.CombinedLayout(pg.ParentWindow(), gtx)
	}

	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx, body)
	}
	return pg.layoutDesktop(gtx, body)
}

func (pg *TxDetailsPage) layoutDesktop(gtx C, body layout.Widget) D {
	return components.UniformPadding(gtx, body)
}

func (pg *TxDetailsPage) layoutMobile(gtx C, body layout.Widget) D {
	return components.UniformMobile(gtx, false, false, body)
}

func (pg *TxDetailsPage) txDetailsHeader(gtx C) D {
	return decredmaterial.LinearLayout{
		Width:       decredmaterial.MatchParent,
		Height:      decredmaterial.WrapContent,
		Orientation: layout.Horizontal,
		Padding: layout.Inset{
			Top:    values.MarginPadding24,
			Bottom: values.MarginPadding30,
		},
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Right: values.MarginPadding22,
			}.Layout(gtx, pg.txnWidgets.txStatus.Icon.Layout24dp)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(pg.Theme.Label(values.TextSize16, values.String(values.StrStatus)+": ").Layout),
								layout.Rigid(pg.Theme.Label(values.TextSize16, pg.txnWidgets.txStatus.Title).Layout),
								layout.Rigid(func(gtx C) D {
									if pg.txnWidgets.txStatus.TicketStatus == dcrlibwallet.TicketStatusImmature {
										confs := pg.transaction.Confirmations(pg.WL.SelectedWallet.Wallet.GetBestBlock())
										progress := (float32(confs) / float32(pg.WL.MultiWallet.TicketMaturity())) * 100

										p := pg.Theme.ProgressBarCirle(int(progress))
										p.Color = pg.txnWidgets.txStatus.ProgressBarColor
										return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
											gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding22)
											gtx.Constraints.Min.X = gtx.Constraints.Max.X
											gtx.Constraints.Max.Y = gtx.Dp(values.MarginPadding22)
											gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
											return p.Layout(gtx)
										})
									}
									return D{}
								}),
							)
						}),
						layout.Rigid(func(gtx C) D {
							col := pg.Theme.Color.GrayText2
							switch pg.txnWidgets.txStatus.TicketStatus {
							case dcrlibwallet.TicketStatusImmature:
								maturity := pg.WL.MultiWallet.TicketMaturity()
								blockTime := pg.WL.MultiWallet.TargetTimePerBlockMinutes()
								maturityDuration := time.Duration(maturity*int32(blockTime)) * time.Minute

								lbl := pg.Theme.Label(values.TextSize16, values.StringF(values.StrImmatureInfo, pg.transaction.BlockHeight, maturity,
									maturityDuration.String()))
								lbl.Color = col
								return lbl.Layout(gtx)

							case dcrlibwallet.TicketStatusLive:
								return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										lbl := pg.Theme.Label(values.TextSize16, "Life Span: ")
										lbl.Color = col
										return lbl.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										lbl := pg.Theme.Label(values.TextSize16, values.String(values.StrLiveInfoDisc))
										lbl.Color = col
										return lbl.Layout(gtx)
									}),
								)

							case dcrlibwallet.TicketStatusVotedOrRevoked:
								if pg.ticketSpender.Type == dcrlibwallet.TxTypeVote {
									return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											lbl := pg.Theme.Label(values.TextSize16, "Reward: ")
											lbl.Color = col
											return lbl.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											lbl := pg.Theme.Label(values.TextSize16, dcrutil.Amount(pg.transaction.VoteReward).String())
											lbl.Color = col
											return lbl.Layout(gtx)
										}),
									)
								}

								return D{}

							default:
								if pg.ticketSpender != nil { // voted or revoked
									if pg.ticketSpender.Type == dcrlibwallet.TxTypeVote {
										maturity := pg.WL.MultiWallet.TicketMaturity()
										confirmations := pg.transaction.Confirmations(pg.WL.SelectedWallet.Wallet.GetBestBlock())

										timeRemaining := time.Duration(float64(maturity-confirmations)*pg.WL.MultiWallet.TargetTimePerBlockMinutes()) * time.Minute
										maturityDuration := components.TimeFormat(int(timeRemaining.Seconds()), false)

										return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												lbl := pg.Theme.Label(values.TextSize16, values.String(values.StrDaysToVote+": "))
												lbl.Color = col
												return lbl.Layout(gtx)
											}),
											layout.Rigid(func(gtx C) D {
												lbl := pg.Theme.Label(values.TextSize16, fmt.Sprintf("%v", maturityDuration))
												lbl.Color = col
												return lbl.Layout(gtx)
											}),
										)

										// lbl := pg.Theme.Label(values.TextSize16, dcrutil.Amount(pg.transaction.VoteReward).String())
										// lbl.Color = col
										// return lbl.Layout(gtx)
									}
								}
								return D{}
							}
							// todo== ticket vote statistics
						}),
						layout.Rigid(func(gtx C) D {
							if pg.transaction.BlockHeight == -1 {
								if !pg.rebroadcastClickable.Enabled() {
									gtx = pg.rebroadcastClickable.SetEnabled(false, &gtx)
								}
								return decredmaterial.LinearLayout{
									Width:     decredmaterial.WrapContent,
									Height:    decredmaterial.WrapContent,
									Clickable: pg.rebroadcastClickable,
									Direction: layout.Center,
									Alignment: layout.Middle,
									Border:    decredmaterial.Border{Color: pg.Theme.Color.Gray2, Width: values.MarginPadding1, Radius: decredmaterial.Radius(10)},
									Padding:   layout.Inset{Top: values.MarginPadding3, Bottom: values.MarginPadding3, Left: values.MarginPadding8, Right: values.MarginPadding8},
									Margin:    layout.Inset{Left: values.MarginPadding10},
								}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
											return pg.rebroadcastIcon.Layout16dp(gtx)
										})
									}),
									layout.Rigid(func(gtx C) D {
										return pg.rebroadcast.Layout(gtx)
									}))
							}
							return D{}
						}),
					)
				}),
			)
		}),
	)
}

func (pg *TxDetailsPage) txnBalanceAndStatus(gtx C) D {
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
			}.Layout(gtx, pg.txnWidgets.txStatus.Icon.Layout24dp)
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
								label.Color = pg.Theme.Color.GrayText2
								return layout.Inset{
									Left: values.MarginPadding8,
								}.Layout(gtx, label.Layout)
							}
							return D{}
						}),
						layout.Rigid(func(gtx C) D {
							if pg.transaction.BlockHeight == -1 {
								if !pg.rebroadcastClickable.Enabled() {
									gtx = pg.rebroadcastClickable.SetEnabled(false, &gtx)
								}
								return decredmaterial.LinearLayout{
									Width:     decredmaterial.WrapContent,
									Height:    decredmaterial.WrapContent,
									Clickable: pg.rebroadcastClickable,
									Direction: layout.Center,
									Alignment: layout.Middle,
									Border:    decredmaterial.Border{Color: pg.Theme.Color.Gray2, Width: values.MarginPadding1, Radius: decredmaterial.Radius(10)},
									Padding:   layout.Inset{Top: values.MarginPadding3, Bottom: values.MarginPadding3, Left: values.MarginPadding8, Right: values.MarginPadding8},
									Margin:    layout.Inset{Left: values.MarginPadding10},
								}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
											return pg.rebroadcastIcon.Layout16dp(gtx)
										})
									}),
									layout.Rigid(func(gtx C) D {
										return pg.rebroadcast.Layout(gtx)
									}))
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
						pg.txnWidgets.time.Color = pg.Theme.Color.Gray1
						return pg.txnWidgets.time.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{
								Right: values.MarginPadding4,
								Top:   values.MarginPadding4,
							}.Layout(gtx, func(gtx C) D {
								return pg.txnWidgets.confirmationIcons.Layout12dp(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							txt := pg.Theme.Body1("")
							if pg.txConfirmations() > 1 {
								txt.Text = strings.Title(values.String(values.StrConfirmed))
								txt.Color = pg.Theme.Color.Success
							} else {
								txt.Text = strings.Title(values.String(values.StrPending))
								txt.Color = pg.Theme.Color.GrayText2
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
							txt.Color = pg.Theme.Color.GrayText2
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
			t := pg.Theme.Label(values.TextSize14, values.String(values.StrMaturity))
			t.Color = pg.Theme.Color.GrayText2
			return t.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {

			percentageLabel := pg.Theme.Label(values.TextSize14, "25%")
			percentageLabel.Color = pg.Theme.Color.GrayText2

			progress := pg.Theme.ProgressBar(40)
			progress.Color = pg.Theme.Color.LightBlue
			progress.TrackColor = pg.Theme.Color.BlueProgressTint
			progress.Height = values.MarginPadding8
			progress.Width = values.MarginPadding80
			progress.Radius = decredmaterial.Radius(8)

			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(percentageLabel.Layout),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Left: values.MarginPadding6, Right: values.MarginPadding6}.Layout(gtx, progress.Layout)
					}),
					layout.Rigid(pg.Theme.Label(values.TextSize16, fmt.Sprintf("%d %s", 18, values.String(values.StrHours))).Layout),
				)
			})
		}),
	)
}

func (pg *TxDetailsPage) keyValue(gtx C, key string, value layout.Widget) D {
	return layout.Inset{Bottom: values.MarginPadding18}.Layout(gtx, func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Flexed(.4, func(gtx C) D {
				return layout.Inset{Right: values.MarginPadding35}.Layout(gtx, func(gtx C) D {
					lbl := pg.Theme.Label(values.TextSize14, key)
					lbl.Color = pg.Theme.Color.GrayText2
					return lbl.Layout(gtx)
				})
			}),
			layout.Flexed(.6, value),
		)
	})
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
								status = values.String(values.StrVoted)
							} else {
								status = values.String(values.StrRevoked)
							}
						} else if pg.wallet.TxMatchesFilter(pg.transaction, dcrlibwallet.TxFilterLive) {
							status = values.String(values.StrLive)
						} else if pg.wallet.TxMatchesFilter(pg.transaction, dcrlibwallet.TxFilterImmature) {
							status = values.String(values.StrImmature)
						} else if pg.wallet.TxMatchesFilter(pg.transaction, dcrlibwallet.TxFilterUnmined) {
							status = values.String(values.StrUmined)
						} else if pg.wallet.TxMatchesFilter(pg.transaction, dcrlibwallet.TxFilterExpired) {
							status = values.String(values.StrExpired)
						} else {
							status = values.String(values.StrUnknown)
						}

						return layout.Inset{Top: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
							return pg.txnInfoSection(gtx, values.String(values.StrStatus), status, false, nil)
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
							txt := values.String(values.StrDaysToVote)
							return pg.txnInfoSection(gtx, txt, fmt.Sprintf("%d %s", pg.transaction.DaysToVoteOrRevoke, values.String(values.StrDays)), false, nil)
						})
					}

					return D{}
				}),
				layout.Rigid(func(gtx C) D {
					if pg.transaction.Type == dcrlibwallet.TxTypeVote {
						return layout.Inset{Top: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
							txt := values.String(values.StrReward)
							return pg.txnInfoSection(gtx, txt, dcrutil.Amount(pg.transaction.VoteReward).String(), false, nil)
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
					layout.Rigid(pg.Theme.Label(values.TextSize16, values.String(values.StrViewTicket)).Layout),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, pg.Theme.Icons.Next.Layout24dp)
					}),
				)
			})
		}),
		layout.Rigid(pg.Theme.Separator().Layout),
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

func (pg *TxDetailsPage) txnTypeAndID(gtx C) D {
	transaction := pg.transaction
	return decredmaterial.LinearLayout{
		Width:       decredmaterial.MatchParent,
		Height:      decredmaterial.WrapContent,
		Orientation: layout.Vertical,
		Padding:     layout.UniformInset(values.MarginPadding16),
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.keyValue(gtx, values.String(values.StrAccount), pg.Theme.Label(values.TextSize14, pg.txSourceAccount).Layout)
		}),
		layout.Rigid(func(gtx C) D {
			key := values.String(values.StrTicketPrice)
			if pg.transaction.Type == "Ticket" {
				key = values.String(values.StrAmount)
			}

			amount := dcrutil.Amount(pg.transaction.Amount).String()
			if pg.transaction.Type == dcrlibwallet.TxTypeMixed {
				amount = dcrutil.Amount(pg.transaction.MixDenomination).String()
			} else if pg.transaction.Type == dcrlibwallet.TxTypeRegular && pg.transaction.Direction == dcrlibwallet.TxDirectionSent {
				amount = "-" + amount
			}
			return pg.keyValue(gtx, key, pg.Theme.Label(values.TextSize14, amount).Layout)
		}),
		layout.Rigid(func(gtx C) D {
			if transaction.BlockHeight != -1 {
				return pg.keyValue(gtx, values.String(values.StrIncludedInBlock), pg.Theme.Label(values.TextSize14, fmt.Sprintf("%d", transaction.BlockHeight)).Layout)
			}
			return D{}
		}),
		layout.Rigid(func(gtx C) D {
			if pg.ticketSpender != nil { // voted or revoked
				if pg.ticketSpender.Type == dcrlibwallet.TxTypeVote {
					maturity := pg.WL.MultiWallet.TicketMaturity()
					confirmations := pg.transaction.Confirmations(pg.WL.SelectedWallet.Wallet.GetBestBlock())

					timeRemaining := time.Duration(float64(maturity-confirmations)*pg.WL.MultiWallet.TargetTimePerBlockMinutes()) * time.Minute
					// maturityDuration := components.TimeFormat(int(timeRemaining.Seconds()), false)

					progressBar := func(gtx C) D {
						percentageLabel := pg.Theme.Label(values.TextSize14, "25%")
						percentageLabel.Color = pg.Theme.Color.GrayText2

						progress := pg.Theme.ProgressBar(int(timeRemaining.Seconds()))
						progress.Color = pg.Theme.Color.Success2
						progress.TrackColor = pg.Theme.Color.Success2
						progress.Height = values.MarginPadding8
						progress.Width = values.MarginPadding80
						progress.Radius = decredmaterial.Radius(8)

						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Flex{
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(percentageLabel.Layout),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Left: values.MarginPadding6, Right: values.MarginPadding6}.Layout(gtx, progress.Layout)
								}),
								layout.Rigid(pg.Theme.Label(values.TextSize16, fmt.Sprintf("%d %s", 18, values.String(values.StrHours))).Layout),
							)
						})
					}

					return pg.keyValue(gtx, "Spendable In:", progressBar)
				}
			}
			return pg.keyValue(gtx, "Purchased On", pg.txnWidgets.time.Layout)
		}),
		layout.Rigid(func(gtx C) D {
			stat := func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Right: values.MarginPadding4,
							Top:   values.MarginPadding4,
						}.Layout(gtx, func(gtx C) D {
							return pg.txnWidgets.confirmationIcons.Layout12dp(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						txt := pg.Theme.Body2("")
						if pg.txConfirmations() > 1 {
							txt.Text = strings.Title(values.String(values.StrConfirmed))
							txt.Color = pg.Theme.Color.Success
						} else {
							txt.Text = strings.Title(values.String(values.StrPending))
							txt.Color = pg.Theme.Color.GrayText2
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
						txt := pg.Theme.Body2(values.StringF(values.StrNConfirmations, pg.txConfirmations()))
						txt.Color = pg.Theme.Color.GrayText2
						return txt.Layout(gtx)
					}),
				)
			}
			return pg.keyValue(gtx, "Confirmation Status", stat)
		}),
		layout.Rigid(func(gtx C) D {
			return pg.keyValue(gtx, "Transaction Fee", pg.Theme.Label(values.TextSize14, dcrutil.Amount(transaction.Fee).String()).Layout)
		}),
		layout.Rigid(func(gtx C) D {
			// todo -- from dcrlibwallet
			return pg.keyValue(gtx, "VSP", pg.Theme.Label(values.TextSize14, "vsp.stakeminer.com").Layout)
		}),
		layout.Rigid(func(gtx C) D {
			// todo -- from dcrlibwallet
			return pg.keyValue(gtx, "VSP Fee", pg.Theme.Label(values.TextSize14, "0.0000121 DCR").Layout)
		}),
		layout.Rigid(func(gtx C) D {
			first := transaction.Hash[0 : len(transaction.Hash)-20]
			second := transaction.Hash[len(transaction.Hash)-20:]

			btn := pg.Theme.OutlineButton(first + " " + second)
			btn.TextSize = values.TextSize14
			btn.SetClickable(pg.hashClickable)
			btn.Inset = layout.UniformInset(values.MarginPadding0)
			return pg.keyValue(gtx, values.String(values.StrTransactionID), btn.Layout)
		}),
		// layout.Rigid(func(gtx C) D {
		// 	return layout.Inset{Top: m}.Layout(gtx, func(gtx C) D {
		// 		return pg.txnInfoSection(gtx, values.String(values.StrType), transaction.Type, false, nil)
		// 	})
		// }),
		// layout.Rigid(func(gtx C) D {
		// 	return layout.Inset{Top: m}.Layout(gtx, func(gtx C) D {
		// 		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// 			layout.Rigid(func(gtx C) D {
		// 				t := pg.Theme.Label(values.TextSize14, values.String(values.StrTransactionID))
		// 				t.Color = pg.Theme.Color.GrayText2
		// 				return t.Layout(gtx)
		// 			}),
		// 			layout.Rigid(func(gtx C) D {
		// 				btn := pg.Theme.OutlineButton(transaction.Hash)
		// 				btn.TextSize = values.TextSize14
		// 				btn.SetClickable(pg.hashClickable)
		// 				btn.Inset = layout.UniformInset(values.MarginPadding0)
		// 				return btn.Layout(gtx)
		// 			}),
		// 		)
		// 	})
		// }),
	)
}

func (pg *TxDetailsPage) txnInfoSection(gtx C, label, value string, showWalletBadge bool, clickable *widget.Clickable) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			t := pg.Theme.Label(values.TextSize14, label)
			t.Color = pg.Theme.Color.GrayText2
			return t.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if showWalletBadge {
						card := pg.Theme.Card()
						card.Radius = decredmaterial.Radius(0)
						card.Color = pg.Theme.Color.Gray4
						return card.Layout(gtx, func(gtx C) D {
							return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
								txt := pg.Theme.Body2(pg.wallet.Name)
								txt.Color = pg.Theme.Color.GrayText2
								return txt.Layout(gtx)
							})
						})
					}
					return D{}
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						if clickable == nil {
							txt := pg.Theme.Body1(value)
							return txt.Layout(gtx)
						}

						btn := pg.Theme.OutlineButton(value)
						btn.TextSize = values.TextSize14
						btn.SetClickable(clickable)
						btn.Inset = layout.UniformInset(values.MarginPadding0)
						return btn.Layout(gtx)
					})
				}),
			)
		}),
	)
}

func (pg *TxDetailsPage) txnInputs(gtx C) D {
	transaction := pg.transaction

	collapsibleHeader := func(gtx C) D {
		t := pg.Theme.Body1(values.StringF(values.StrXInputsConsumed, len(transaction.Inputs)))
		t.Color = pg.Theme.Color.GrayText2
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

func (pg *TxDetailsPage) txnOutputs(gtx C) D {
	transaction := pg.transaction

	collapsibleHeader := func(gtx C) D {
		t := pg.Theme.Body1(values.StringF(values.StrXOutputCreated, len(transaction.Outputs)))
		t.Color = pg.Theme.Color.GrayText2
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

func (pg *TxDetailsPage) txnIORow(gtx C, amount int64, acctNum int32, address string, i int) D {

	accountName := values.String(values.StrExternal)
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
		card.Color = pg.Theme.Color.Gray4
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
								card.Color = pg.Theme.Color.Gray4
								return card.Layout(gtx, func(gtx C) D {
									return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
										txt := pg.Theme.Body2(walletName)
										txt.Color = pg.Theme.Color.GrayText2
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

func (pg *TxDetailsPage) layoutOptionsMenu(gtx C) {
	inset := layout.Inset{
		Left: values.MarginPaddingMinus150,
	}

	m := op.Record(gtx.Ops)
	inset.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding168)
		return pg.shadowBox.Layout(gtx, func(gtx C) D {
			optionsMenuCard := decredmaterial.Card{Color: pg.Theme.Color.Surface}
			optionsMenuCard.Radius = decredmaterial.Radius(5)
			return optionsMenuCard.Layout(gtx, func(gtx C) D {
				return (&layout.List{Axis: layout.Vertical}).Layout(gtx, len(pg.moreItems), func(gtx C, i int) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return pg.moreItems[i].button.Layout(gtx, func(gtx C) D {
								return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
									gtx.Constraints.Min.X = gtx.Constraints.Max.X
									return pg.Theme.Label(values.TextSize14, pg.moreItems[i].text).Layout(gtx)
								})
							})
						}),
					)
				})
			})
		})
	})
	op.Defer(gtx.Ops, m.Stop())
}

func (pg *TxDetailsPage) pageSections(gtx C, body layout.Widget) D {
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, body)
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *TxDetailsPage) HandleUserInteractions() {
	for pg.moreOption.Clicked() {
		pg.moreOptionIsOpen = !pg.moreOptionIsOpen
	}

	for pg.associatedTicketClickable.Clicked() {
		if pg.ticketSpent != nil {
			pg.txBackStack = pg.transaction
			pg.transaction = pg.ticketSpent
			pg.getTXSourceAccountAndDirection()
			pg.txnWidgets = initTxnWidgets(pg.Load, pg.transaction)
			pg.ParentWindow().Reload()
		}
	}

	if pg.rebroadcastClickable.Clicked() {
		go func() {
			pg.rebroadcastClickable.SetEnabled(false, nil)
			if !pg.Load.WL.MultiWallet.IsConnectedToDecredNetwork() {
				// if user is not conected to the network, notify the user
				pg.Toast.NotifyError(values.String(values.StrNotConnected))
				if !pg.rebroadcastClickable.Enabled() {
					pg.rebroadcastClickable.SetEnabled(true, nil)
				}
				return
			}

			err := pg.wallet.PublishUnminedTransactions()
			if err != nil {
				// If transactions are not published, notify the user
				pg.Toast.NotifyError(err.Error())
			} else {
				pg.Toast.Notify(values.String(values.StrRepublished))
			}
			if !pg.rebroadcastClickable.Enabled() {
				pg.rebroadcastClickable.SetEnabled(true, nil)
			}
		}()
	}

	redirectURL := pg.WL.Wallet.GetBlockExplorerURL(pg.transaction.Hash)
	for _, menu := range pg.moreItems {
		if menu.button.Clicked() && menu.id == viewBlockID {
			components.GoToURL(redirectURL)
			pg.moreOptionIsOpen = false
			break
		}

	}
}

func (pg *TxDetailsPage) handleTextCopyEvent(gtx C) {
	for _, b := range pg.txnWidgets.copyTextButtons {
		for b.Clicked() {
			clipboard.WriteOp{Text: b.Text}.Add(gtx.Ops)
			pg.Toast.Notify(values.String(values.StrCopied))
			break
		}
		break
	}

	for pg.hashClickable.Clicked() {
		clipboard.WriteOp{Text: pg.transaction.Hash}.Add(gtx.Ops)
		pg.Toast.Notify(values.String(values.StrTxHashCopied))
		break
	}

	for pg.destAddressClickable.Clicked() {
		clipboard.WriteOp{Text: pg.txDestinationAddress}.Add(gtx.Ops)
		pg.Toast.Notify(values.String(values.StrAddressCopied))
		break
	}

	redirectURL := pg.WL.Wallet.GetBlockExplorerURL(pg.transaction.Hash)
	for _, menu := range pg.moreItems {
		if menu.button.Clicked() && menu.id == copyBlockID {
			clipboard.WriteOp{Text: redirectURL}.Add(gtx.Ops)
			pg.Toast.Notify("URL copied")
			pg.moreOptionIsOpen = false
			break
		}
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *TxDetailsPage) OnNavigatedFrom() {}

func initTxnWidgets(l *load.Load, transaction *dcrlibwallet.Transaction) transactionWdg {

	var txn transactionWdg
	wal := l.WL.SelectedWallet.Wallet

	t := time.Unix(transaction.Timestamp, 0).UTC()
	txn.time = l.Theme.Body2(t.Format(time.UnixDate))
	txn.status = l.Theme.Body1("")
	txn.wallet = l.Theme.Body2(wal.Name)

	if components.TxConfirmations(l, *transaction) > 1 {
		txn.status.Text = components.FormatDateOrTime(transaction.Timestamp)
		txn.confirmationIcons = l.Theme.Icons.ConfirmIcon
	} else {
		txn.status.Text = values.String(values.StrPending)
		txn.status.Color = l.Theme.Color.GrayText2
		txn.confirmationIcons = l.Theme.Icons.PendingIcon
	}

	var ticketSpender *dcrlibwallet.Transaction
	if wal.TxMatchesFilter(transaction, dcrlibwallet.TxFilterStaking) {
		ticketSpender, _ = wal.TicketSpender(transaction.Hash)
	}
	txStatus := components.TransactionTitleIcon(l, wal, transaction, ticketSpender)
	txn.txStatus = txStatus

	x := len(transaction.Inputs) + len(transaction.Outputs)
	txn.copyTextButtons = make([]decredmaterial.Button, x)
	for i := 0; i < x; i++ {
		btn := l.Theme.OutlineButton("")
		btn.TextSize = values.TextSize14
		btn.Inset = layout.UniformInset(values.MarginPadding0)
		txn.copyTextButtons[i] = btn
	}

	return txn
}

package transaction

import (
	"context"
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const TransactionsPageID = "Transactions"

type (
	C = layout.Context
	D = layout.Dimensions
)

type TransactionsPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	*listeners.TxAndBlockNotificationListener
	ctx       context.Context // page context
	ctxCancel context.CancelFunc
	separator decredmaterial.Line

	walletTabList         *decredmaterial.ClickableList // Tab list of all loaded wallets.
	selectedCategoryIndex int
	walletTabTitles       []string
	changed               bool

	orderDropDown   *decredmaterial.DropDown
	txTypeDropDown  *decredmaterial.DropDown
	walletDropDown  *decredmaterial.DropDown
	transactionList *decredmaterial.ClickableList
	container       *widget.List
	transactions    []dcrlibwallet.Transaction
	wallets         []*dcrlibwallet.Wallet
}

func NewTransactionsPage(l *load.Load) *TransactionsPage {
	pg := &TransactionsPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(TransactionsPageID),
		container: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		separator:       l.Theme.Separator(),
		transactionList: l.Theme.NewClickableList(layout.Vertical),
		walletTabList:   l.Theme.NewClickableList(layout.Horizontal),
	}

	pg.walletTabList.IsHoverable = false
	pg.transactionList.Radius = decredmaterial.Radius(14)
	pg.transactionList.IsShadowEnabled = true

	pg.orderDropDown = components.CreateOrderDropDown(l, values.TxDropdownGroup, 1)
	pg.wallets = pg.WL.SortedWalletList()
	pg.walletDropDown = components.CreateOrUpdateWalletDropDown(pg.Load, &pg.walletDropDown, pg.wallets, values.TxDropdownGroup, 0)
	pg.txTypeDropDown = l.Theme.DropDown([]decredmaterial.DropDownItem{
		{
			Text: values.String(values.StrAll),
		},
		{
			Text: values.String(values.StrSent),
		},
		{
			Text: values.String(values.StrReceived),
		},
		{
			Text: values.String(values.StrYourself),
		},
		{
			Text: values.String(values.StrMixed),
		},
		{
			Text: values.String(values.StrStaking),
		},
	}, values.TxDropdownGroup, 2)

	walletTitles := make([]string, 0)
	for _, wallet := range pg.wallets {
		walletTitles = append(walletTitles, wallet.Name)
	}

	pg.walletTabTitles = walletTitles

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *TransactionsPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())

	pg.listenForTxNotifications()
	pg.loadTransactions(pg.walletDropDown.SelectedIndex())
}

func (pg *TransactionsPage) loadTransactions(selectedWalletIndex int) {
	selectedWallet := pg.wallets[selectedWalletIndex]
	newestFirst := pg.orderDropDown.SelectedIndex() == 0

	txFilter := dcrlibwallet.TxFilterAll
	switch pg.txTypeDropDown.SelectedIndex() {
	case 1:
		txFilter = dcrlibwallet.TxFilterSent
	case 2:
		txFilter = dcrlibwallet.TxFilterReceived
	case 3:
		txFilter = dcrlibwallet.TxFilterTransferred
	case 4:
		txFilter = dcrlibwallet.TxFilterMixed
	case 5:
		txFilter = dcrlibwallet.TxFilterStaking
	}

	wallTxs, err := selectedWallet.GetTransactionsRaw(0, 0, txFilter, newestFirst) //TODO
	if err != nil {
		// log.Error("Error loading transactions:", err)
	} else {
		pg.transactions = wallTxs
	}
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *TransactionsPage) Layout(gtx layout.Context) layout.Dimensions {
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx)
	}
	return pg.layoutDesktop(gtx)
}

func (pg *TransactionsPage) layoutDesktop(gtx layout.Context) layout.Dimensions {
	container := func(gtx C) D {
		wallTxs := pg.transactions
		return layout.Stack{Alignment: layout.N}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return layout.Inset{
					Top: values.MarginPadding60,
				}.Layout(gtx, func(gtx C) D {
					return pg.Theme.List(pg.container).Layout(gtx, 1, func(gtx C, i int) D {
						return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
							return pg.Theme.Card().Layout(gtx, func(gtx C) D {

								// return "No transactions yet" text if there are no transactions
								if len(wallTxs) == 0 {
									padding := values.MarginPadding16
									txt := pg.Theme.Body1(values.String(values.StrNoTransactions))
									txt.Color = pg.Theme.Color.GrayText3
									gtx.Constraints.Min.X = gtx.Constraints.Max.X
									return layout.Center.Layout(gtx, func(gtx C) D {
										return layout.Inset{Top: padding, Bottom: padding}.Layout(gtx, txt.Layout)
									})
								}

								return pg.transactionList.Layout(gtx, len(wallTxs), func(gtx C, index int) D {
									var row = components.TransactionRow{
										Transaction: wallTxs[index],
										Index:       index,
										ShowBadge:   false,
									}

									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return components.LayoutTransactionRow(gtx, pg.Load, row)
										}),
										layout.Rigid(func(gtx C) D {
											// No divider for last row
											if row.Index == len(wallTxs)-1 {
												return layout.Dimensions{}
											}

											gtx.Constraints.Min.X = gtx.Constraints.Max.X
											separator := pg.Theme.Separator()
											return layout.E.Layout(gtx, func(gtx C) D {
												// Show bottom divider for all rows except last
												return layout.Inset{Left: values.MarginPadding56}.Layout(gtx, separator.Layout)
											})
										}),
									)
								})
							})
						})
					})
				})
			}),
			layout.Expanded(func(gtx C) D {
				return pg.walletDropDown.Layout(gtx, 0, false)
			}),
			layout.Expanded(func(gtx C) D {
				return pg.orderDropDown.Layout(gtx, 0, true)
			}),
			layout.Expanded(func(gtx C) D {
				return pg.txTypeDropDown.Layout(gtx, pg.orderDropDown.Width-4, true)
			}),
		)
	}
	return components.UniformPadding(gtx, container)
}

func (pg *TransactionsPage) layoutMobile(gtx layout.Context) layout.Dimensions {
	container := func(gtx C) D {
		wallTxs := pg.transactions
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				if len(pg.wallets) > 1 {
					return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
						return pg.layoutTabs(gtx)
					})
				}
				return D{}
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Stack{Alignment: layout.N}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return layout.Inset{
							Top: values.MarginPadding60,
						}.Layout(gtx, func(gtx C) D {
							return pg.Theme.List(pg.container).Layout(gtx, 1, func(gtx C, i int) D {
								return pg.Theme.Card().Layout(gtx, func(gtx C) D {

									// return "No transactions yet" text if there are no transactions
									if len(wallTxs) == 0 {
										padding := values.MarginPadding16
										txt := pg.Theme.Body1(values.String(values.StrNoTransactions))
										txt.Color = pg.Theme.Color.GrayText3
										gtx.Constraints.Min.X = gtx.Constraints.Max.X
										return layout.Center.Layout(gtx, func(gtx C) D {
											return layout.Inset{Top: padding, Bottom: padding}.Layout(gtx, txt.Layout)
										})
									}

									return pg.transactionList.Layout(gtx, len(wallTxs), func(gtx C, index int) D {
										var row = components.TransactionRow{
											Transaction: wallTxs[index],
											Index:       index,
											ShowBadge:   false,
										}

										return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												return components.LayoutTransactionRow(gtx, pg.Load, row)
											}),
											layout.Rigid(func(gtx C) D {
												// No divider for last row
												if row.Index == len(wallTxs)-1 {
													return layout.Dimensions{}
												}

												gtx.Constraints.Min.X = gtx.Constraints.Max.X
												separator := pg.Theme.Separator()
												return layout.E.Layout(gtx, func(gtx C) D {
													// Show bottom divider for all rows except last
													return layout.Inset{Left: values.MarginPadding56}.Layout(gtx, separator.Layout)
												})
											}),
										)
									})
								})
							})
						})
					}),
					layout.Expanded(func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
							return pg.orderDropDown.Layout(gtx, 0, true)
						})
					}),
					layout.Expanded(func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
							return pg.txTypeDropDown.Layout(gtx, pg.orderDropDown.Width-4, true)
						})
					}),
				)
			}),
		)

	}
	return components.UniformMobile(gtx, true, container)
}

func (pg *TransactionsPage) layoutTabs(gtx C) D {
	var dims layout.Dimensions

	return pg.walletTabList.Layout(gtx, len(pg.walletTabTitles), func(gtx C, i int) D {
		return layout.Stack{Alignment: layout.S}.Layout(gtx,
			layout.Stacked(func(gtx C) D {
				return layout.Inset{
					Right:  values.MarginPadding24,
					Bottom: values.MarginPadding8,
				}.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						lbl := pg.Theme.Label(values.TextSize16, pg.walletTabTitles[i])
						lbl.Color = pg.Theme.Color.GrayText1
						if pg.selectedCategoryIndex == i {
							lbl.Color = pg.Theme.Color.Primary
							dims = lbl.Layout(gtx)
						}

						return lbl.Layout(gtx)
					})
				})
			}),
			layout.Stacked(func(gtx C) D {
				if pg.selectedCategoryIndex != i {
					return D{}
				}

				tabHeight := gtx.Dp(values.MarginPadding2)
				tabRect := image.Rect(0, 0, dims.Size.X, tabHeight)

				return layout.Inset{
					Left: values.MarginPaddingMinus22,
				}.Layout(gtx, func(gtx C) D {
					paint.FillShape(gtx.Ops, pg.Theme.Color.Primary, clip.Rect(tabRect).Op())
					return layout.Dimensions{
						Size: image.Point{X: dims.Size.X, Y: tabHeight},
					}
				})
			}),
		)
	})
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *TransactionsPage) HandleUserInteractions() {

	for pg.txTypeDropDown.Changed() {
		pg.loadTransactions(pg.walletDropDown.SelectedIndex())
	}

	for pg.orderDropDown.Changed() {
		pg.loadTransactions(pg.walletDropDown.SelectedIndex())
	}

	for pg.walletDropDown.Changed() {
		pg.loadTransactions(pg.walletDropDown.SelectedIndex())
	}

	if clicked, selectedItem := pg.transactionList.ItemClicked(); clicked {
		pg.ParentNavigator().Display(NewTransactionDetailsPage(pg.Load, &pg.transactions[selectedItem]))
	}
	decredmaterial.DisplayOneDropdown(pg.walletDropDown, pg.txTypeDropDown, pg.orderDropDown)

	if clicked, selectedItem := pg.walletTabList.ItemClicked(); clicked {
		if pg.selectedCategoryIndex != selectedItem {
			pg.selectedCategoryIndex = selectedItem
			pg.changed = true
		}

		pg.loadTransactions(pg.selectedCategoryIndex)
		pg.changed = false
	}
}

func (pg *TransactionsPage) listenForTxNotifications() {
	if pg.TxAndBlockNotificationListener != nil {
		return
	}
	pg.TxAndBlockNotificationListener = listeners.NewTxAndBlockNotificationListener()
	err := pg.WL.MultiWallet.AddTxAndBlockNotificationListener(pg.TxAndBlockNotificationListener, true, TransactionsPageID)
	if err != nil {
		log.Errorf("Error adding tx and block notification listener: %v", err)
		return
	}

	go func() {
		for {
			select {
			case n := <-pg.TxAndBlockNotifChan:
				if n.Type == listeners.NewTransaction {

					selectedWallet := pg.wallets[pg.walletDropDown.SelectedIndex()]
					if selectedWallet.ID == n.Transaction.WalletID {
						pg.loadTransactions(pg.walletDropDown.SelectedIndex())
						pg.ParentWindow().Reload()
					}
				}
			case <-pg.ctx.Done():
				pg.WL.MultiWallet.RemoveTxAndBlockNotificationListener(TransactionsPageID)
				close(pg.TxAndBlockNotifChan)
				pg.TxAndBlockNotificationListener = nil

				return
			}
		}
	}()
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *TransactionsPage) OnNavigatedFrom() {
	pg.ctxCancel()
}

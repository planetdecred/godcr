package staking

import (
	"context"
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	tpage "github.com/planetdecred/godcr/ui/page/transaction"
	"github.com/planetdecred/godcr/ui/values"
)

const listPageID = "StakingList"

type txType int

const (
	All txType = iota
	Unmined
	Immature
	Live
	Voted
	Expired
	Revoked
)

type ListPage struct {
	*load.Load
	*listeners.TxAndBlockNotificationListener
	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	tickets     []*transactionItem
	ticketsList *decredmaterial.ClickableList
	scrollBar   *widget.List

	orderDropDown      *decredmaterial.DropDown
	ticketTypeDropDown *decredmaterial.DropDown
	walletDropDown     *decredmaterial.DropDown
	backButton         decredmaterial.IconButton

	wallets []*dcrlibwallet.Wallet
}

func newListPage(l *load.Load) *ListPage {
	pg := &ListPage{
		Load:        l,
		ticketsList: l.Theme.NewClickableList(layout.Vertical),
		scrollBar: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
	}
	pg.backButton, _ = components.SubpageHeaderButtons(pg.Load)

	pg.orderDropDown = createOrderDropDown(l.Theme)
	pg.wallets = pg.WL.SortedWalletList()
	components.CreateOrUpdateWalletDropDown(pg.Load, &pg.walletDropDown, pg.wallets, values.StakingDropdownGroup, 0) // first in the set during layout pos 0
	pg.ticketTypeDropDown = l.Theme.DropDown([]decredmaterial.DropDownItem{
		{Text: "All"},
		{Text: "Unmined"},
		{Text: "Immature"},
		{Text: "Live"},
		{Text: "Voted"},
		{Text: "Expired"},
		{Text: "Revoked"},
	}, values.StakingDropdownGroup, 2)

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *ListPage) ID() string {
	return listPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *ListPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
	pg.listenForTxNotifications()
	pg.fetchTickets()
}

func (pg *ListPage) listenForTxNotifications() {
	if pg.TxAndBlockNotificationListener == nil {
		pg.TxAndBlockNotificationListener = listeners.NewTxAndBlockNotificationListener()
	}
	err := pg.WL.MultiWallet.AddTxAndBlockNotificationListener(pg.TxAndBlockNotificationListener, true, listPageID)
	if err != nil {
		log.Errorf("Error adding tx and block notification listener: %v", err)
		return
	}

	go func() {
		for {

			select {
			case n := <-pg.TxAndBlockNotifChan:
				if n.Type == listeners.BlockAttached {
					selectedWallet := pg.wallets[pg.walletDropDown.SelectedIndex()]
					if selectedWallet.ID == n.WalletID {
						pg.fetchTickets()
						pg.RefreshWindow()
					}
				} else if n.Type == listeners.NewTransaction {
					selectedWallet := pg.wallets[pg.walletDropDown.SelectedIndex()]
					if selectedWallet.ID == n.Transaction.WalletID {
						pg.fetchTickets()
						pg.RefreshWindow()
					}
				}
			case <-pg.ctx.Done():
				pg.WL.MultiWallet.RemoveTxAndBlockNotificationListener(listPageID) // Remove listener
				close(pg.TxAndBlockNotifChan)                                      // Close channel

				return
			}
		}
	}()
}

func (pg *ListPage) fetchTickets() {
	var txFilter int32
	var ticketTypeDropdown = txType(pg.ticketTypeDropDown.SelectedIndex())
	switch ticketTypeDropdown {
	case Unmined:
		txFilter = dcrlibwallet.TxFilterUnmined
	case Immature:
		txFilter = dcrlibwallet.TxFilterImmature
	case Live:
		txFilter = dcrlibwallet.TxFilterLive
	case Voted:
		txFilter = dcrlibwallet.TxFilterTickets
	case Expired:
		txFilter = dcrlibwallet.TxFilterExpired
	case Revoked:
		txFilter = dcrlibwallet.TxFilterTickets
	default:
		txFilter = dcrlibwallet.TxFilterTickets
	}

	newestFirst := pg.orderDropDown.SelectedIndex() == 0
	selectedWalletID := pg.wallets[pg.walletDropDown.SelectedIndex()].ID
	multiWallet := pg.WL.MultiWallet
	w := multiWallet.WalletWithID(selectedWalletID)
	txs, err := w.GetTransactionsRaw(0, 0, txFilter, newestFirst)
	if err != nil {
		pg.Toast.NotifyError(err.Error())
		return
	}

	tickets, err := stakeToTransactionItems(pg.Load, txs, newestFirst, func(filter int32) bool {
		switch filter {
		case dcrlibwallet.TxFilterVoted:
			return ticketTypeDropdown == Voted
		case dcrlibwallet.TxFilterRevoked:
			return ticketTypeDropdown == Revoked
		}

		return filter == txFilter
	})
	if err != nil {
		pg.Toast.NotifyError(err.Error())
		return
	}

	pg.tickets = tickets
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *ListPage) Layout(gtx C) D {
	body := func(gtx C) D {
		page := components.SubPage{
			Load:       pg.Load,
			Title:      "All tickets",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return layout.Stack{Alignment: layout.N}.Layout(gtx,
					layout.Expanded(func(gtx C) D {
						return layout.Inset{Top: values.MarginPadding60}.Layout(gtx, func(gtx C) D {
							return pg.Theme.List(pg.scrollBar).Layout(gtx, 1, func(gtx C, index int) D {
								return pg.Theme.Card().Layout(gtx, func(gtx C) D {
									tickets := pg.tickets

									if len(tickets) == 0 {
										gtx.Constraints.Min.X = gtx.Constraints.Max.X

										txt := pg.Theme.Body1("No tickets yet")
										txt.Color = pg.Theme.Color.GrayText3
										txt.Alignment = text.Middle
										return layout.Inset{Top: values.MarginPadding15, Bottom: values.MarginPadding16}.Layout(gtx, txt.Layout)
									}

									return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
										return pg.ticketListLayout(gtx, tickets)
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
						return pg.ticketTypeDropDown.Layout(gtx, pg.orderDropDown.Width-4, true)
					}),
				)
			},
		}
		return page.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

func (pg *ListPage) ticketListLayout(gtx layout.Context, tickets []*transactionItem) layout.Dimensions {
	return pg.ticketsList.Layout(gtx, len(tickets), func(gtx C, index int) D {
		var ticket = tickets[index]

		return ticketListLayout(gtx, pg.Load, ticket, index, false)
	})
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *ListPage) HandleUserInteractions() {
	for pg.orderDropDown.Changed() {
		pg.fetchTickets()
	}

	for pg.walletDropDown.Changed() {
		pg.fetchTickets()
	}

	for pg.ticketTypeDropDown.Changed() {
		pg.fetchTickets()
	}

	if clicked, selectedItem := pg.ticketsList.ItemClicked(); clicked {
		ticketTx := pg.tickets[selectedItem].transaction
		pg.ChangeFragment(tpage.NewTransactionDetailsPage(pg.Load, ticketTx))

		// Check if this ticket is fully registered with a VSP
		// and log any discrepancies.
		// NOTE: Wallet needs to be unlocked to get the ticket status
		// from the vsp. Otherwise, only the wallet-stored info will
		// be retrieved. This is fine because we're only just logging
		// but where it is necessary to display vsp-stored info, the
		// wallet passphrase should be requested and used to unlock
		// the wallet before calling this method.
		// TODO: Use log.Errorf and log.Warnf instead of fmt.Printf.
		ticketInfo, err := pg.WL.MultiWallet.VSPTicketInfo(ticketTx.WalletID, ticketTx.Hash)
		if err != nil {
			fmt.Printf("VSPTicketInfo error: %v\n", err)
		} else {
			if ticketInfo.FeeTxStatus != dcrlibwallet.VSPFeeProcessConfirmed {
				fmt.Printf("[WARN] Ticket %s has unconfirmed fee tx %s with status %q, vsp %s \n",
					ticketTx.Hash, ticketInfo.FeeTxHash, ticketInfo.FeeTxStatus.String(), ticketInfo.VSP)
			}
			if ticketInfo.ConfirmedByVSP == nil || !*ticketInfo.ConfirmedByVSP {
				fmt.Printf("[WARN] Ticket %s is not confirmed by VSP %s. Fee tx %s, status %q \n",
					ticketTx.Hash, ticketInfo.VSP, ticketInfo.FeeTxHash, ticketInfo.FeeTxStatus.String())
			}
		}
	}

	decredmaterial.DisplayOneDropdown(pg.ticketTypeDropDown, pg.orderDropDown, pg.walletDropDown)
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *ListPage) OnNavigatedFrom() {
	pg.ctxCancel()
}

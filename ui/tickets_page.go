package ui

import (
	"fmt"
	"strconv"

	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"

	"gioui.org/layout"
	"gioui.org/text"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const PageTickets = "Tickets"

type ticketPage struct {
	th   *decredmaterial.Theme
	wal  *wallet.Wallet
	vspd *dcrlibwallet.VSPD

	ticketPageContainer layout.List
	ticketList          layout.List
	ticketPurchaseList  layout.List
	unconfirmedList     layout.List

	purchaseTicketButton decredmaterial.Button
	inputNumbTickets     decredmaterial.Editor
	inputExpiryBlocks    decredmaterial.Editor
	tickets              **wallet.Tickets
	transactions         **wallet.Transactions
	ticketPrice          string

	walletsDropdown  *decredmaterial.DropDown
	accountsDropdown *decredmaterial.DropDown
}

func (win *Window) TicketPage(common pageCommon) layout.Widget {
	pg := &ticketPage{
		th:           common.theme,
		wal:          win.wallet,
		tickets:      &win.walletTickets,
		transactions: &win.walletTransactions,

		ticketList:           layout.List{Axis: layout.Vertical},
		unconfirmedList:      layout.List{Axis: layout.Vertical},
		ticketPageContainer:  layout.List{Axis: layout.Vertical},
		ticketPurchaseList:   layout.List{Axis: layout.Vertical},
		inputNumbTickets:     common.theme.Editor(new(widget.Editor), "Enter tickets amount"),
		inputExpiryBlocks:    common.theme.Editor(new(widget.Editor), "Expiry in blocks"),
		purchaseTicketButton: common.theme.Button(new(widget.Clickable), "Purchase"),
	}
	pg.purchaseTicketButton.TextSize = values.TextSize12
	pg.inputNumbTickets.Editor.SingleLine, pg.inputExpiryBlocks.Editor.SingleLine = true, true

	return func(gtx C) D {
		pg.Handler(common)
		return pg.layout(gtx, common)
	}
}

func (pg *ticketPage) layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	return c.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		pg.setWallets(c)
		pg.setAccounts(c)

		sections := []func(gtx C) D{
			func(ctx layout.Context) layout.Dimensions {
				return pg.LayoutTicketPurchase(gtx, c)
			},
			func(ctx layout.Context) layout.Dimensions {
				return pg.ticketRowHeader(gtx, c)
			},
			func(ctx layout.Context) layout.Dimensions {
				return pg.LayoutTicketList(gtx, c)
			},
			func(ctx layout.Context) layout.Dimensions {
				walletID := c.info.Wallets[pg.walletsDropdown.SelectedIndex()].ID
				if pg.tickets != nil {
					if len((*pg.tickets).Unconfirmed[walletID]) > 0 {
						return pg.LayoutUnconfirmedPurchased(gtx, c)
					}
				}
				return layout.Dimensions{}
			},
		}

		return layout.Stack{Alignment: layout.N}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: values.MarginPadding60}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return pg.th.Card().Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return pg.ticketPageContainer.Layout(gtx, len(sections), func(gtx C, i int) D {
							return layout.Inset{}.Layout(gtx, sections[i])
						})
					})
				})
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
									return pg.walletsDropdown.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return pg.accountsDropdown.Layout(gtx)
							}),
						)
					}),
				)
			}),
		)
	})
}

func (pg *ticketPage) LayoutTicketPurchase(gtx layout.Context, c pageCommon) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding20).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return pg.walletBalance(gtx, c)
							})
						})
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx C) D {
								return c.theme.H6(pg.ticketPrice).Layout(gtx)
							})
						})
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(.5, func(gtx C) D {
							return pg.inputNumbTickets.Layout(gtx)
						}),
						layout.Flexed(.5, func(gtx C) D {
							return pg.inputExpiryBlocks.Layout(gtx)
						}),
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx C) D {
						return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
							return pg.purchaseTicketButton.Layout(gtx)
						})
					}),
				)
			}),
		)
	})
}

func (pg *ticketPage) walletBalance(gtx layout.Context, common pageCommon) layout.Dimensions {
	selectedWallet := common.info.Wallets[pg.walletsDropdown.SelectedIndex()]
	selectedAccount := selectedWallet.Accounts[pg.accountsDropdown.SelectedIndex()]

	selectedDetails := func(gtx C) D {
		return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
								return pg.th.H6(selectedAccount.Name).Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
								return pg.th.H6(dcrutil.Amount(selectedAccount.SpendableBalance).String()).Layout(gtx)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
								return pg.th.Body2(selectedAccount.Name).Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
								return pg.th.Body2(selectedWallet.Balance).Layout(gtx)
							})
						}),
					)
				}),
			)
		})
	}

	card := pg.th.Card()
	card.Radius = decredmaterial.CornerRadius{
		NE: 0,
		NW: 0,
		SE: 0,
		SW: 0,
	}
	return card.Layout(gtx, selectedDetails)
}

func (pg *ticketPage) LayoutTicketList(gtx layout.Context, common pageCommon) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding0).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				margin := values.MarginPadding20
				return layout.Inset{Left: margin, Right: margin, Bottom: margin}.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if *pg.tickets == nil {
							return layout.Dimensions{}
						}
						walletID := common.info.Wallets[pg.walletsDropdown.SelectedIndex()].ID
						tickets := (*pg.tickets).Confirmed[walletID]
						if len(tickets) == 0 {
							return common.theme.Body2("No ticket").Layout(gtx)
						}

						return pg.ticketList.Layout(gtx, len(tickets), func(gtx C, index int) D {
							return pg.ticketRowInfo(gtx, common, tickets[index])
						})
					})
				})
			}),
		)
	})
}

func (pg *ticketPage) LayoutUnconfirmedPurchased(gtx layout.Context, common pageCommon) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			tickets := (*pg.tickets).Unconfirmed[common.info.Wallets[pg.walletsDropdown.SelectedIndex()].ID]
			margin := values.MarginPadding20
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Left: margin, Right: margin, Bottom: margin}.Layout(gtx, func(gtx C) D {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return pg.unconfirmedList.Layout(gtx, len(tickets), func(gtx C, index int) D {
								return pg.ticketRowInfo(gtx, common, tickets[index])
							})
						})
					})
				}),
			)
		}),
	)
}

func (pg *ticketPage) ticketRowHeader(gtx layout.Context, common pageCommon) layout.Dimensions {
	txt := common.theme.Label(values.MarginPadding15, "#")
	txt.MaxLines = 1
	txt.Color = common.theme.Color.Hint
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			margin := values.MarginPadding20
			return layout.Inset{Right: margin, Left: margin}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Px(values.MarginPadding60)
							return txt.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Px(values.MarginPadding120)
							txt.Text = "Date (UTC)"
							return txt.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Px(values.MarginPadding120)
							txt.Text = "Status"
							return txt.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Px(values.MarginPadding150)
							txt.Text = "Hash"
							return txt.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Px(values.MarginPadding150)
							txt.Text = "Amount"
							txt.Alignment = text.End
							return txt.Layout(gtx)
						}),
					)
				})
			})
		}),
	)
}

// func (pg *ticketPage) ticketRowInfo(gtx layout.Context, c *pageCommon, ticket wallet.Ticket) layout.Dimensions {
func (pg *ticketPage) ticketRowInfo(gtx layout.Context, c pageCommon, ticket interface{}) layout.Dimensions {
	txt := c.theme.Label(values.MarginPadding15, "")
	txt.MaxLines = 1

	var dateTime, status, hash, amount string
	var blockHeight int32
	switch t := ticket.(type) {
	case wallet.Ticket:
		info := t.Info
		blockHeight = info.BlockHeight
		dateTime, status, hash, amount = t.DateTime, info.Status, info.Ticket.Hash.String(), t.Amount
	case wallet.UnconfirmedPurchase:
		blockHeight, dateTime, status, hash, amount = t.BlockHeight, t.DateTime, t.Status, t.Hash, t.Amount
	}

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding60)
			return c.theme.Label(values.MarginPadding15, fmt.Sprintf("%d", blockHeight)).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding120)
			txt.Text = dateTime
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding120)
			txt.Text = status
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding150)
			txt.Text = fmt.Sprintf("%s...%s", hash[:8], hash[56:])
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding150)
			txt.Alignment = text.End
			txt.Text = amount
			return txt.Layout(gtx)
		}),
	)
}

func (pg *ticketPage) purchaseTicket(c pageCommon, password []byte) {
	selectedWallet := c.info.Wallets[pg.walletsDropdown.SelectedIndex()]
	selectedAccount := selectedWallet.Accounts[pg.accountsDropdown.SelectedIndex()]

	pg.vspd = c.wallet.NewVSPD(selectedWallet.ID, selectedAccount.Number)
	_, err := pg.vspd.GetInfo()
	if err != nil {
		log.Error("[GetInfo] err:", err)
		return
	}

	numbTicketsStr := pg.inputNumbTickets.Editor.Text()
	numbTickets, err := strconv.Atoi(numbTicketsStr)
	if err != nil {
		c.Notify(err.Error(), false)
		return
	}

	expiryBlocksStr := pg.inputExpiryBlocks.Editor.Text()
	expiryBlocks, err := strconv.Atoi(expiryBlocksStr)
	if err != nil {
		c.Notify(err.Error(), false)
		return
	}

	hash, err := c.wallet.PurchaseTicket(selectedWallet.ID, selectedAccount.Number, uint32(numbTickets), password, uint32(expiryBlocks))
	if err != nil {
		log.Error("[PurchaseTicket] err:", err)
		c.Notify(err.Error(), false)
		return
	}
	go c.wallet.GetAllTransactions(0, 0, 0)

	transactionResponse, err := pg.getTicketFee(hash, password)
	if err != nil {
		log.Error("[GetTicketFee] err:", err)
		c.Notify(err.Error(), false)
		return
	}

	_, err = pg.vspd.PayVSPFee(transactionResponse, hash, "", password)
	if err != nil {
		log.Error("[PayVSPFee] err:", err)
		c.Notify(err.Error(), false)
		return
	}

	c.Notify("success", true)
	c.closeModal()
}

func (pg *ticketPage) getTicketFee(ticketHash string, password []byte) (feeTransaction string, err error) {
	resp, err := pg.vspd.GetVSPFeeAddress(ticketHash, password)
	if err != nil {
		return
	}

	feeTransaction, err = pg.vspd.CreateTicketFeeTx(resp.FeeAmount, ticketHash, resp.FeeAddress, password)
	if err != nil {
		return
	}

	return
}

func (pg *ticketPage) setWallets(common pageCommon) {
	if len(common.info.Wallets) == 0 || pg.walletsDropdown != nil {
		return
	}

	var walletDropdownItems []decredmaterial.DropDownItem
	for i := range common.info.Wallets {
		item := decredmaterial.DropDownItem{
			Text: common.info.Wallets[i].Name,
			Icon: common.icons.walletIcon,
		}
		walletDropdownItems = append(walletDropdownItems, item)
	}
	pg.walletsDropdown = common.theme.DropDown(walletDropdownItems, 0)
	tkPrice := common.wallet.TicketPrice(common.info.Wallets[pg.walletsDropdown.SelectedIndex()].ID)
	pg.ticketPrice = fmt.Sprintf("Current Ticket Price: %s", tkPrice)
}

func (pg *ticketPage) setAccounts(common pageCommon) {
	if pg.accountsDropdown != nil {
		return
	}

	var accountsDropdownItems []decredmaterial.DropDownItem
	selectedWallet := pg.walletsDropdown.SelectedIndex()
	for i := range common.info.Wallets[selectedWallet].Accounts {
		item := decredmaterial.DropDownItem{
			Text: common.info.Wallets[selectedWallet].Accounts[i].Name,
		}
		accountsDropdownItems = append(accountsDropdownItems, item)
	}
	pg.accountsDropdown = common.theme.DropDown(accountsDropdownItems, 0)
}

func (pg *ticketPage) Handler(c pageCommon) {
	for _, evt := range pg.inputNumbTickets.Editor.Events() {
		switch evt.(type) {
		case widget.ChangeEvent:
			tkPrice := c.wallet.TicketPrice(c.info.Wallets[pg.walletsDropdown.SelectedIndex()].ID)
			pg.ticketPrice = fmt.Sprintf("Current Ticket Price: %s", tkPrice)
		}
	}

	if pg.purchaseTicketButton.Button.Clicked() {
		go func() {
			c.modalReceiver <- &modalLoad{
				template: PasswordTemplate,
				title:    "Confirm to purchase",
				confirm: func(pass string) {
					go pg.purchaseTicket(c, []byte(pass))
				},
				confirmText: "Confirm",
				cancel:      c.closeModal,
				cancelText:  "Cancel",
			}
		}()
	}
}

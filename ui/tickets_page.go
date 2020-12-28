package ui

import (
	"fmt"
	"strconv"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const (
	PageTickets = "Tickets"
	ticketsView = iota
	payTicketView
)

type ticketPage struct {
	vspd                 *dcrlibwallet.VSPD
	activeView           int
	ticketPageContainer  layout.List
	ticketList           layout.List
	purchaseTicketButton decredmaterial.Button
	getFeeButton         decredmaterial.Button
	payFeeButton         decredmaterial.Button
	backToTicketsButton  decredmaterial.Button
	selectedWallet       wallet.InfoShort
	selectedAccount      wallet.Account
	shouldInitializeVSPD bool
	inputNumbTickets     decredmaterial.Editor
	inputExpiryBlocks    decredmaterial.Editor
	passwordModal        *decredmaterial.Password
	feeTx                string
	isModalPurchase      bool
	isModalGetTicketFee  bool
	isModalPayFee        bool
	tickets              **wallet.Tickets
	tiketHash            string
	ticketPrice          string
}

func (win *Window) TicketPage(common pageCommon) layout.Widget {
	pg := &ticketPage{
		tickets:              &win.walletTickets,
		activeView:           ticketsView,
		ticketList:           layout.List{Axis: layout.Vertical},
		ticketPageContainer:  layout.List{Axis: layout.Vertical},
		getFeeButton:         common.theme.Button(new(widget.Clickable), "Get Fee"),
		payFeeButton:         common.theme.Button(new(widget.Clickable), "Pay Fee"),
		inputNumbTickets:     common.theme.Editor(new(widget.Editor), "Enter tickets amount"),
		inputExpiryBlocks:    common.theme.Editor(new(widget.Editor), "Expiry in blocks"),
		purchaseTicketButton: common.theme.Button(new(widget.Clickable), "Purchase"),
		backToTicketsButton:  common.theme.Button(new(widget.Clickable), "Tickets"),
		passwordModal:        common.theme.Password(),
		tiketHash:            "",
	}
	pg.purchaseTicketButton.TextSize = values.TextSize12
	pg.backToTicketsButton.TextSize = values.TextSize12
	pg.getFeeButton.TextSize = values.TextSize12
	pg.payFeeButton.TextSize = values.TextSize12

	pg.inputNumbTickets.IsRequired = true
	pg.inputNumbTickets.Editor.SingleLine = true
	pg.inputExpiryBlocks.IsRequired = true
	pg.inputExpiryBlocks.Editor.SingleLine = true

	return func(gtx C) D {
		pg.Handler(common)
		return pg.layout(gtx, common)
	}
}

func (pg *ticketPage) Handler(c pageCommon) {
	if pg.selectedAccount.CurrentAddress != c.info.Wallets[*c.selectedWallet].Accounts[*c.selectedAccount].CurrentAddress {
		pg.shouldInitializeVSPD = true
		pg.selectedAccount = c.info.Wallets[*c.selectedWallet].Accounts[*c.selectedAccount]
	}
	if pg.selectedWallet.ID != c.info.Wallets[*c.selectedWallet].ID {
		pg.shouldInitializeVSPD = true
		pg.selectedWallet = c.info.Wallets[*c.selectedWallet]
	}

	if pg.shouldInitializeVSPD {
		pg.shouldInitializeVSPD = false
		tkPrice := c.wallet.TicketPrice(pg.selectedWallet.ID)
		pg.ticketPrice = fmt.Sprintf("Current Ticket Price: %s", tkPrice)
		pg.vspd = c.wallet.NewVSPD(pg.selectedWallet.ID, pg.selectedAccount.Number)
		// pg.password = nil
		_, err := pg.vspd.GetInfo()
		if err != nil {
			log.Error("[GetInfo] err:", err)
			return
		}
	}

	if pg.purchaseTicketButton.Button.Clicked() {
		if pg.activeView == ticketsView {
			pg.activeView = payTicketView
			return
		}
		pg.isModalPurchase = true
	}

	if pg.backToTicketsButton.Button.Clicked() {
		pg.activeView = ticketsView
	}

	if pg.getFeeButton.Button.Clicked() {
		pg.isModalGetTicketFee = true
	}

	if pg.payFeeButton.Button.Clicked() {
		pg.isModalPayFee = true
	}
}

func (pg *ticketPage) layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	if pg.activeView == ticketsView {
		return pg.LayoutTicketList(gtx, &c)
	}
	return pg.LayoutTicketPurchase(gtx, &c)
}

func (pg *ticketPage) LayoutTicketPurchase(gtx layout.Context, c *pageCommon) layout.Dimensions {
	marginTop := layout.Inset{Top: values.MarginPadding20}.Layout(gtx, func(gtx C) D { return layout.Dimensions{} })
	widgets := []func(gtx C) D{
		func(gtx C) D { return marginTop },
		func(gtx C) D {
			return c.theme.H6(pg.ticketPrice).Layout(gtx)
		},
		func(gtx C) D { return marginTop },
		func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(.5, func(gtx C) D {
						return pg.inputNumbTickets.Layout(gtx)
					}),
					layout.Flexed(.5, func(gtx C) D {
						return pg.inputExpiryBlocks.Layout(gtx)
					}),
				)
			})
		},
		func(gtx C) D { return marginTop },
		func(gtx C) D {
			if pg.tiketHash == "" {
				return layout.Dimensions{}
			}
			return c.theme.Label(values.MarginPadding15, fmt.Sprintf("Ticket hash: %s", pg.tiketHash)).Layout(gtx)
		},
		func(gtx C) D { return marginTop },
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return pg.purchaseTicketButton.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return pg.getFeeButton.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return pg.payFeeButton.Layout(gtx)
				}),
			)
		},
		func(gtx C) D { return marginTop },
		func(gtx C) D { return pg.backToTicketsButton.Layout(gtx) },
	}

	dims := c.LayoutWithAccounts(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding20).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, c.SelectedAccountLayout)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.ticketPageContainer.Layout(gtx, len(widgets), func(gtx C, i int) D {
						return layout.Inset{}.Layout(gtx, widgets[i])
					})
				}),
			)
		})
	})

	if pg.isModalPurchase {
		return c.Modal(gtx, dims, pg.drawPasswordModalPurchase(gtx, c))
	}

	if pg.isModalGetTicketFee {
		return c.Modal(gtx, dims, pg.drawPasswordModalGetTicketFee(gtx))
	}

	if pg.isModalPayFee {
		return c.Modal(gtx, dims, pg.drawPasswordModalPayFee(gtx))
	}

	return dims
}

func (pg *ticketPage) LayoutTicketList(gtx layout.Context, common *pageCommon) layout.Dimensions {
	dims := common.LayoutWithAccounts(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding20).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, common.SelectedAccountLayout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
						return pg.ticketRowHeader(gtx, common)
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						if (*pg.tickets).List == nil {
							return layout.Dimensions{}
						}
						walletID := common.info.Wallets[*common.selectedWallet].ID
						tickets := (*pg.tickets).List[walletID]
						if len(tickets) == 0 {
							return common.theme.Body2("No ticket").Layout(gtx)
						}
						return pg.ticketList.Layout(gtx, len(tickets), func(gtx C, index int) D {
							return pg.ticketRowInfo(gtx, common, tickets[index])
						})
					})
				}),
				layout.Rigid(func(gtx C) D {
					return pg.purchaseTicketButton.Layout(gtx)
				}),
			)
		})
	})
	return dims
}

func (pg *ticketPage) ticketRowHeader(gtx layout.Context, c *pageCommon) layout.Dimensions {
	txt := c.theme.Label(values.MarginPadding15, "#")
	txt.MaxLines = 1
	txt.Color = c.theme.Color.Hint
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
}

func (pg *ticketPage) ticketRowInfo(gtx layout.Context, c *pageCommon, ticket wallet.Ticket) layout.Dimensions {
	txt := c.theme.Label(values.MarginPadding15, "")
	txt.MaxLines = 1
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding60)
			return c.theme.Label(values.MarginPadding15, fmt.Sprintf("%d", ticket.Info.BlockHeight)).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding120)
			txt.Text = ticket.DateTime
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding120)
			txt.Text = ticket.Info.Status
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding150)
			hash := ticket.Info.Ticket.Hash.String()
			txt.Text = fmt.Sprintf("%s...%s", hash[:8], hash[56:])
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Px(values.MarginPadding150)
			txt.Alignment = text.End
			txt.Text = ticket.Amount
			return txt.Layout(gtx)
		}),
	)
}

func (pg *ticketPage) drawPasswordModalPurchase(gtx layout.Context, c *pageCommon) layout.Dimensions {
	return pg.passwordModal.Layout(gtx, func(password []byte) {
		numbTicketsStr := pg.inputNumbTickets.Editor.Text()
		numbTickets, err := strconv.Atoi(numbTicketsStr)
		if err != nil {
			pg.passwordModal.WithError(err.Error())
			return
		}

		expiryBlocksStr := pg.inputExpiryBlocks.Editor.Text()
		expiryBlocks, err := strconv.Atoi(expiryBlocksStr)
		if err != nil {
			pg.passwordModal.WithError(err.Error())
			return
		}

		h, err := c.wallet.PurchaseTicket(pg.selectedWallet.ID, pg.selectedAccount.Number, uint32(numbTickets), password, uint32(expiryBlocks))
		if err != nil {
			pg.passwordModal.WithError(err.Error())
			return
		}
		pg.tiketHash = h
		pg.isModalPurchase = false
	}, func() {
		pg.isModalPurchase = false
	})
}

func (pg *ticketPage) drawPasswordModalGetTicketFee(gtx layout.Context) layout.Dimensions {
	return pg.passwordModal.Layout(gtx, func(password []byte) {
		resp, err := pg.vspd.GetVSPFeeAddress(pg.tiketHash, password)
		if err != nil {
			pg.passwordModal.WithError(err.Error())
			return
		}
		pg.feeTx, err = pg.vspd.CreateTicketFeeTx(resp.FeeAmount, pg.tiketHash, resp.FeeAddress, password)
		if err != nil {
			pg.passwordModal.WithError(err.Error())
			return
		}
		pg.isModalGetTicketFee = false
	}, func() {
		pg.isModalGetTicketFee = false
	})
}

func (pg *ticketPage) drawPasswordModalPayFee(gtx layout.Context) layout.Dimensions {
	return pg.passwordModal.Layout(gtx, func(password []byte) {
		msg, err := pg.vspd.PayVSPFee(pg.feeTx, pg.tiketHash, "", password)
		if err != nil {
			pg.passwordModal.WithError(err.Error())
			return
		}
		log.Info("5: [Done]", msg.Request.VoteChoices)
		pg.isModalPayFee = false
	}, func() {
		pg.isModalPayFee = false
	})
}

package ui

import (
	"fmt"
	"strings"

	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"

	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const PageTickets = "Tickets"

type ticketPage struct {
	th   *decredmaterial.Theme
	wal  *wallet.Wallet
	vspd *dcrlibwallet.VSPD

	ticketPageContainer      layout.List
	ticketsLive              layout.List
	ticketsLiveCounterStatus layout.List
	ticketsActivity          layout.List
	ticketPurchaseList       layout.List
	unconfirmedList          layout.List

	purchaseTicketButton decredmaterial.Button
	tickets              **wallet.Tickets
	ticketPrice          string

	// walletsDropdown  *decredmaterial.DropDown
	// accountsDropdown *decredmaterial.DropDown

	autoPurchaseEnabled *widget.Bool
	toTickets           decredmaterial.TextAndIconButton
	toTicketsActivity   decredmaterial.TextAndIconButton
	ticketStatusIc      map[string]map[string]*widget.Image
}

func (win *Window) TicketPage(c pageCommon) layout.Widget {
	pg := &ticketPage{
		th:      c.theme,
		wal:     win.wallet,
		tickets: &win.walletTickets,

		ticketsLive:              layout.List{Axis: layout.Horizontal},
		ticketsLiveCounterStatus: layout.List{Axis: layout.Horizontal},
		ticketsActivity:          layout.List{Axis: layout.Vertical},
		unconfirmedList:          layout.List{Axis: layout.Vertical},
		ticketPageContainer:      layout.List{Axis: layout.Vertical},
		ticketPurchaseList:       layout.List{Axis: layout.Vertical},
		purchaseTicketButton:     c.theme.Button(new(widget.Clickable), "Purchase"),
		autoPurchaseEnabled:      new(widget.Bool),
		toTickets:                c.theme.TextAndIconButton(new(widget.Clickable), "See All", c.icons.navigationArrowForward),
		toTicketsActivity:        c.theme.TextAndIconButton(new(widget.Clickable), "See All", c.icons.navigationArrowForward),
	}
	pg.purchaseTicketButton.TextSize = values.TextSize12

	pg.toTickets.Color = c.theme.Color.Primary
	pg.toTickets.BackgroundColor = c.theme.Color.Surface

	pg.toTicketsActivity.Color = c.theme.Color.Primary
	pg.toTicketsActivity.BackgroundColor = c.theme.Color.Surface

	pg.ticketStatusIc = map[string]map[string]*widget.Image{
		"UNKNOWN": {
			"head":     nil,
			"live":     nil,
			"activity": nil,
		},
		"UNMINED": {
			"head":     c.icons.ti.ticketUnminedIcon3,
			"live":     c.icons.ti.ticketUnminedIcon1,
			"activity": c.icons.ti.ticketUnminedIcon2,
		},
		"IMMATURE": {
			"head":     c.icons.ti.ticketImmatureIcon3,
			"live":     c.icons.ti.ticketImmatureIcon1,
			"activity": c.icons.ti.ticketImmatureIcon2,
		},
		"LIVE": {
			"head":     c.icons.ti.ticketLiveIcon3,
			"live":     c.icons.ti.ticketLiveIcon1,
			"activity": c.icons.ti.ticketLiveIcon2,
		},
		"VOTED": {
			"head":     c.icons.ti.ticketVotedIcon3,
			"live":     c.icons.ti.ticketVotedIcon1,
			"activity": c.icons.ti.ticketVotedIcon2,
		},
		"MISSED": {
			"head":     c.icons.ti.ticketMissedIcon3,
			"live":     c.icons.ti.ticketMissedIcon1,
			"activity": c.icons.ti.ticketMissedIcon2,
		},
		"EXPIRED": {
			"head":     c.icons.ti.ticketExpiredIcon3,
			"live":     c.icons.ti.ticketExpiredIcon1,
			"activity": c.icons.ti.ticketExpiredIcon2,
		},
		"REVOKED": {
			"head":     c.icons.ti.ticketRevokedIcon3,
			"live":     c.icons.ti.ticketRevokedIcon1,
			"activity": c.icons.ti.ticketRevokedIcon2,
		},
	}

	return func(gtx C) D {
		pg.Handler(c)
		return pg.layout(gtx, c)
	}
}

func (pg *ticketPage) layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	return c.Layout(gtx, func(gtx C) D {
		return c.UniformPadding(gtx, func(gtx layout.Context) layout.Dimensions {
			sections := []func(gtx C) D{
				func(ctx layout.Context) layout.Dimensions {
					return pg.ticketPriceSection(gtx, c)
				},
				func(ctx layout.Context) layout.Dimensions {
					return pg.ticketsLiveSection(gtx, c)
				},
				func(ctx layout.Context) layout.Dimensions {
					return pg.ticketsActivitySection(gtx, c)
				},
				func(ctx layout.Context) layout.Dimensions {
					return pg.stackingRecordSection(gtx, c)
				},
			}

			return pg.ticketPageContainer.Layout(gtx, len(sections), func(gtx C, i int) D {
				return sections[i](gtx)
			})
		})
	})
}

func (pg *ticketPage) pageSections(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return layout.Inset{
		Bottom: values.MarginPadding8,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return pg.th.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.UniformInset(values.MarginPadding16).Layout(gtx, body)
		})
	})
}

func (pg *ticketPage) titleRow(gtx layout.Context, leftWidget, rightWidget func(C) D) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return leftWidget(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return rightWidget(gtx)
		}),
	)
}

func (pg *ticketPage) ticketPriceSection(gtx layout.Context, c pageCommon) layout.Dimensions {
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding11,
				}.Layout(gtx, func(gtx C) D {
					return pg.titleRow(gtx, c.theme.Label(values.TextSize14, "Ticket Price").Layout,
						material.Switch(pg.th.Base, pg.autoPurchaseEnabled).Layout,
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding8,
				}.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						ic := c.icons.ti.ticketPurchasedIcon
						ic.Scale = 1.0
						return ic.Layout(gtx)
					})
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding16,
				}.Layout(gtx, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						mainText, subText := breakBalance(c.printer, pg.ticketPrice)
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return pg.th.Label(values.TextSize28, mainText).Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return pg.th.Label(values.TextSize16, subText).Layout(gtx)
							}),
						)
					})
				})
			}),
			layout.Rigid(func(gtx C) D {
				return pg.purchaseTicketButton.Layout(gtx)
			}),
		)
	})
}

func (pg *ticketPage) ticketsLiveSection(gtx layout.Context, c pageCommon) layout.Dimensions {
	if *pg.tickets == nil {
		return layout.Dimensions{}
	}
	tickets := (*pg.tickets).Recent
	if len(tickets) == 0 {
		return layout.Dimensions{}
	}

	type item struct {
		status string
		count  int
	}
	var counterTicketStatus = make(map[string]item)
	for _, ticket := range tickets {
		prev, ok := counterTicketStatus[ticket.Info.Status]
		if ok {
			prev.count += 1
			counterTicketStatus[ticket.Info.Status] = prev
		} else {
			counterTicketStatus[ticket.Info.Status] = item{status: ticket.Info.Status, count: 1}
		}
	}

	var childrens []layout.FlexChild
	for _, item := range counterTicketStatus {
		childrens = append(childrens, layout.Rigid(func(gtx C) D {
			return layout.Inset{Right: values.MarginPadding14}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						ic := pg.ticketStatusIc[item.status]["head"]
						ic.Scale = 1.0
						return ic.Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
							return pg.th.Label(values.TextSize14, fmt.Sprintf("%d", item.count)).Layout(gtx)
						})
					}),
				)
			})
		}))
	}
	childrens = append(childrens, layout.Rigid(func(gtx C) D {
		return pg.toTickets.Layout(gtx)
	}))

	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding14}.Layout(gtx, func(gtx C) D {
					return pg.titleRow(gtx, c.theme.Label(values.TextSize14, "Live Tickets").Layout,
						func(gtx C) D {
							return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
								childrens...,
							)
						},
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return pg.ticketsLive.Layout(gtx, len(tickets), func(gtx C, index int) D {
					return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
						return pg.ticketLiveItemnInfo(gtx, c, tickets[index])
					})
				})
			}),
		)
	})
}

func (pg *ticketPage) ticketLiveItemnInfo(gtx layout.Context, c pageCommon, t wallet.Ticket) layout.Dimensions {
	var itemWidth int
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			wrap := pg.th.Card()
			wrap.Radius.NE = 8 // top - left
			wrap.Radius.SW = 0 // bottom - left
			wrap.Radius.NW = 8 // top - right
			wrap.Radius.SE = 0 // bottom - right
			wrap.Color = c.theme.Color.LightBlue
			return wrap.Layout(gtx, func(gtx C) D {
				return layout.Stack{Alignment: layout.S}.Layout(gtx,

					layout.Expanded(func(gtx C) D {
						return layout.NE.Layout(gtx, func(gtx C) D {
							wTimeLabel := pg.th.Card()
							wTimeLabel.Radius.NE = 0
							wTimeLabel.Radius.SW = 8
							wTimeLabel.Radius.NW = 8
							wTimeLabel.Radius.SE = 0
							return wTimeLabel.Layout(gtx, func(gtx C) D {
								return layout.Inset{
									Top:    values.MarginPadding4,
									Bottom: values.MarginPadding4,
									Right:  values.MarginPadding8,
									Left:   values.MarginPadding8,
								}.Layout(gtx, func(gtx C) D {
									return pg.th.Label(values.TextSize14, "10h 47m").Layout(gtx)
								})
							})
						})
					}),

					layout.Stacked(func(gtx C) D {
						content := layout.Inset{
							Top:    values.MarginPadding24,
							Right:  values.MarginPadding62,
							Left:   values.MarginPadding62,
							Bottom: values.MarginPadding24,
						}.Layout(gtx, func(gtx C) D {
							ic := pg.ticketStatusIc[t.Info.Status]["live"]
							ic.Scale = 1.0
							return ic.Layout(gtx)
						})
						itemWidth = content.Size.X
						return content
					}),

					layout.Stacked(func(gtx C) D {
						return layout.Center.Layout(gtx, func(gtx C) D {
							return layout.Inset{Top: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
								gtx.Constraints.Max.X = itemWidth
								p := pg.th.ProgressBar(20)
								p.Height, p.Radius = values.MarginPadding4, values.MarginPadding1
								p.Color = pg.th.Color.Success
								return p.Layout(gtx)
							})
						})
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			wrap := pg.th.Card()
			wrap.Radius.NE = 0 // top - left
			wrap.Radius.SW = 8 // bottom - left
			wrap.Radius.NW = 0 // top - right
			wrap.Radius.SE = 8 // bottom - right
			return wrap.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X, gtx.Constraints.Max.X = itemWidth, itemWidth
				return layout.Inset{
					Left:   values.MarginPadding12,
					Right:  values.MarginPadding12,
					Bottom: values.MarginPadding8,
				}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{
								Top: values.MarginPadding16,
							}.Layout(gtx, func(gtx C) D {
								return c.layoutBalance(gtx, t.Amount)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return pg.th.Label(values.MarginPadding14, t.WalletName).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{
								Top:    values.MarginPadding16,
								Bottom: values.MarginPadding16,
							}.Layout(gtx, func(gtx C) D {
								txt := pg.th.Label(values.TextSize14, t.MonthDay)
								txt.Color = pg.th.Color.Gray2
								return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return txt.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{
											Left:  values.MarginPadding4,
											Right: values.MarginPadding4,
										}.Layout(gtx, func(gtx C) D {
											ic := c.icons.imageBrightness1
											ic.Color = pg.th.Color.Gray2
											return c.icons.imageBrightness1.Layout(gtx, values.MarginPadding5)
										})
									}),
									layout.Rigid(func(gtx C) D {
										txt.Text = t.DaysBehind
										return txt.Layout(gtx)
									}),
								)
							})
						}),
					)
				})
			})
		}),
	)
}

func (pg *ticketPage) ticketsActivitySection(gtx layout.Context, c pageCommon) layout.Dimensions {
	if *pg.tickets == nil {
		return layout.Dimensions{}
	}
	tickets := (*pg.tickets).Recent
	if len(tickets) == 0 {
		return layout.Dimensions{}
	}

	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding14,
				}.Layout(gtx, func(gtx C) D {
					return pg.titleRow(gtx, c.theme.Label(values.TextSize14, "Ticket Activity").Layout,
						pg.toTicketsActivity.Layout,
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return pg.ticketsActivity.Layout(gtx, len(tickets), func(gtx C, index int) D {
					return pg.ticketActivityItemnInfo(gtx, c, tickets[index], index)
				})
			}),
		)
	})
}

func (pg *ticketPage) ticketActivityItemnInfo(gtx layout.Context, c pageCommon, t wallet.Ticket, index int) layout.Dimensions {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Right: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
				ic := pg.ticketStatusIc[t.Info.Status]["activity"]
				ic.Scale = 1.0
				return ic.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if index == 0 {
						return layout.Dimensions{}
					}
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					separator := pg.th.Separator()
					separator.Width = gtx.Constraints.Max.X
					return layout.E.Layout(gtx, func(gtx C) D {
						return separator.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding8,
						Bottom: values.MarginPadding8,
					}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Flex{
									Spacing:   layout.SpaceBetween,
									Alignment: layout.Middle,
								}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return pg.th.Label(values.TextSize18, strings.ToLower(t.Info.Status)).Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return pg.th.Label(values.TextSize14, t.DaysBehind).Layout(gtx)
									}),
								)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Flex{
									Alignment: layout.Middle,
								}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										txt := pg.th.Label(values.TextSize14, t.WalletName)
										txt.Color = pg.th.Color.Gray2
										return txt.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										ic := c.icons.ticketIconInactive
										ic.Scale = 0.5
										return ic.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										ic := c.icons.imageBrightness1
										ic.Color = pg.th.Color.Gray2
										return c.icons.imageBrightness1.Layout(gtx, values.MarginPadding5)
									}),
									layout.Rigid(func(gtx C) D {
										txt := pg.th.Label(values.TextSize14, t.Amount)
										txt.Color = pg.th.Color.Gray2
										return txt.Layout(gtx)
									}),
								)
							}),
						)
					})
				}),
			)
		}),
	)
}

func (pg *ticketPage) stackingRecordSection(gtx layout.Context, c pageCommon) layout.Dimensions {
	if *pg.tickets == nil {
		return layout.Dimensions{}
	}
	tickets := (*pg.tickets).Recent
	if len(tickets) == 0 {
		return layout.Dimensions{}
	}

	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding14,
				}.Layout(gtx, func(gtx C) D {
					return pg.titleRow(gtx, c.theme.Label(values.TextSize14, "Staking Record").Layout,
						func(gtx C) D { return layout.Dimensions{} })
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.stackingRecordItem(gtx, pg.ticketStatusIc["UNMINED"]["activity"])
					}),
					layout.Rigid(func(gtx C) D {
						return pg.stackingRecordItem(gtx, pg.ticketStatusIc["IMMATURE"]["activity"])
					}),
					layout.Rigid(func(gtx C) D {
						return pg.stackingRecordItem(gtx, pg.ticketStatusIc["LIVE"]["activity"])
					}),
					layout.Rigid(func(gtx C) D {
						return pg.stackingRecordItem(gtx, pg.ticketStatusIc["VOTED"]["activity"])
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.stackingRecordItem(gtx, pg.ticketStatusIc["MISSED"]["activity"])
					}),
					layout.Rigid(func(gtx C) D {
						return pg.stackingRecordItem(gtx, pg.ticketStatusIc["EXPIRED"]["activity"])
					}),
					layout.Rigid(func(gtx C) D {
						return pg.stackingRecordItem(gtx, pg.ticketStatusIc["REVOKED"]["activity"])
					}),
				)
			}),
		)
	})
}

func (pg *ticketPage) stackingRecordItem(gtx layout.Context, ic *widget.Image) D {
	if ic == nil {
		return layout.Dimensions{}
	}
	ic.Scale = 1.0
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return ic.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.th.Label(values.TextSize16, "2").Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						txt := pg.th.Label(values.TextSize12, "Unmined")
						txt.Color = pg.th.Color.Gray2
						return txt.Layout(gtx)
					}),
				)
			})
		}),
	)
}

func (pg *ticketPage) purchaseTicket(c pageCommon, password []byte) {
	// TODO: automatically purchase
	selectedWallet := c.info.Wallets[1]
	selectedAccount := selectedWallet.Accounts[1]

	pg.vspd = c.wallet.NewVSPD(selectedWallet.ID, selectedAccount.Number)
	_, err := pg.vspd.GetInfo()
	if err != nil {
		log.Error("[GetInfo] err:", err)
		return
	}

	hashes, err := c.wallet.PurchaseTicket(selectedWallet.ID, selectedAccount.Number, 1, password, 256)
	if err != nil {
		log.Error("[PurchaseTicket] err:", err)
		c.notify(err.Error(), false)
		return
	}

	for _, hash := range hashes {
		resp, err := pg.vspd.GetVSPFeeAddress(hash, password)
		if err != nil {
			log.Error("[CreateTicketFeeTx] err:", err)
			return
		}

		transactionResponse, err := pg.vspd.CreateTicketFeeTx(resp.FeeAmount, hash, resp.FeeAddress, password)
		if err != nil {
			log.Error("[CreateTicketFeeTx] err:", err)
			c.notify(err.Error(), false)
			return
		}

		_, err = pg.vspd.PayVSPFee(transactionResponse, hash, "", password)
		if err != nil {
			log.Error("[PayVSPFee] err:", err)
			c.notify(err.Error(), false)
			return
		}
	}

	c.notify("success", true)
	c.closeModal()
}

func (pg *ticketPage) Handler(c pageCommon) {
	if len(c.info.Wallets) > 0 && pg.ticketPrice == "" {
		pg.ticketPrice = c.wallet.TicketPrice()
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

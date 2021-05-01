package ui

import (
	"fmt"
	"image/color"
	"strings"

	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"

	"gioui.org/layout"
	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const PageTickets = "Tickets"

type ticketPage struct {
	th  *decredmaterial.Theme
	wal *wallet.Wallet

	ticketPageContainer layout.List
	ticketsLive         layout.List
	ticketsActivity     layout.List
	vspHosts            layout.List

	purchaseTicketButton  decredmaterial.Button
	cancelPurchase        decredmaterial.Button
	reviewPurchase        decredmaterial.Button
	cancelConfirmPurchase decredmaterial.Button
	purchaseTicket        decredmaterial.Button
	tickets               **wallet.Tickets
	ticketPrice           string
	showPurchaseOptions   bool
	showPurchaseConfirm   bool

	purchaseOptions     *decredmaterial.Modal
	autoPurchaseEnabled *widget.Bool
	toTickets           decredmaterial.TextAndIconButton
	toTicketsActivity   decredmaterial.TextAndIconButton
	ticketStatusIc      map[string]map[string]*widget.Image

	rememberVSP      decredmaterial.CheckBoxStyle
	showVSPHosts     bool
	selectVSP        *widget.Clickable
	inputNewVSP      decredmaterial.Editor
	spendingPassword decredmaterial.Editor
	addNewVSP        decredmaterial.Button
}

func (win *Window) TicketPage(c pageCommon) layout.Widget {
	pg := &ticketPage{
		th:      c.theme,
		wal:     win.wallet,
		tickets: &win.walletTickets,

		ticketsLive:           layout.List{Axis: layout.Horizontal},
		ticketsActivity:       layout.List{Axis: layout.Vertical},
		ticketPageContainer:   layout.List{Axis: layout.Vertical},
		vspHosts:              layout.List{Axis: layout.Vertical},
		purchaseTicketButton:  c.theme.Button(new(widget.Clickable), "Purchase"),
		cancelPurchase:        c.theme.Button(new(widget.Clickable), "Cancel"),
		cancelConfirmPurchase: c.theme.Button(new(widget.Clickable), "Cancel"),
		purchaseTicket:        c.theme.Button(new(widget.Clickable), "Purchase 1 ticket"),
		reviewPurchase:        c.theme.Button(new(widget.Clickable), "Review purchase"),
		selectVSP:             new(widget.Clickable),
		autoPurchaseEnabled:   new(widget.Bool),
		toTickets:             c.theme.TextAndIconButton(new(widget.Clickable), "See All", c.icons.navigationArrowForward),
		toTicketsActivity:     c.theme.TextAndIconButton(new(widget.Clickable), "See All", c.icons.navigationArrowForward),
		purchaseOptions:       c.theme.Modal(),
		rememberVSP:           c.theme.CheckBox(new(widget.Bool), "Remember VSP"),
		inputNewVSP:           c.theme.Editor(new(widget.Editor), "Add a new VSP..."),
		addNewVSP:             c.theme.Button(new(widget.Clickable), "Save"),
		spendingPassword:      c.theme.EditorPassword(new(widget.Editor), "Spending password"),
	}
	pg.purchaseTicketButton.TextSize = values.TextSize12

	pg.cancelPurchase.Background = color.NRGBA{}
	pg.cancelPurchase.Color = c.theme.Color.Primary
	pg.cancelConfirmPurchase.Background = color.NRGBA{}
	pg.cancelConfirmPurchase.Color = c.theme.Color.Primary

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
		pg.handler(c)
		return pg.layout(gtx, c)
	}
}

func (pg *ticketPage) layout(gtx layout.Context, c pageCommon) layout.Dimensions {
	dims := c.Layout(gtx, func(gtx C) D {
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

	if pg.showPurchaseConfirm {
		return pg.confirmPurchaseModal(gtx)
	}

	if pg.showVSPHosts {
		return pg.vspHostModalLayout(gtx)
	}

	if pg.showPurchaseOptions && !c.wallAcctSelector.isWalletAccountModalOpen {
		return pg.purchaseModal(gtx, c)
	}

	return dims
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
					tit := c.theme.Label(values.TextSize14, "Ticket Price")
					tit.Color = c.theme.Color.Gray2
					return pg.titleRow(gtx, tit.Layout, material.Switch(pg.th.Base, pg.autoPurchaseEnabled).Layout)
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
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Bottom: values.MarginPadding14}.Layout(gtx, func(gtx C) D {
					tit := c.theme.Label(values.TextSize14, "Live Tickets")
					tit.Color = c.theme.Color.Gray2
					return pg.titleRow(gtx, tit.Layout, func(gtx C) D {
						if *pg.tickets == nil {
							return layout.Dimensions{}
						}
						ticketLiveCounter := (*pg.tickets).LiveCounter
						var elements []layout.FlexChild
						for i := 0; i < len(ticketLiveCounter); i++ {
							item := ticketLiveCounter[i]
							elements = append(elements, layout.Rigid(func(gtx C) D {
								return layout.Inset{Right: values.MarginPadding14}.Layout(gtx, func(gtx C) D {
									return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											ic := pg.ticketStatusIc[item.Status]["head"]
											ic.Scale = 1.0
											return ic.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
												return pg.th.Label(values.TextSize14, fmt.Sprintf("%d", item.Count)).Layout(gtx)
											})
										}),
									)
								})
							}))
						}
						elements = append(elements, layout.Rigid(func(gtx C) D {
							return pg.toTickets.Layout(gtx)
						}))

						return layout.Flex{Alignment: layout.Middle}.Layout(gtx, elements...)
					})
				})
			}),
			layout.Rigid(func(gtx C) D {
				if *pg.tickets == nil {
					return layout.Dimensions{}
				}
				tickets := (*pg.tickets).LiveRecent
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
	tickets := (*pg.tickets).RecentActivity
	if len(tickets) == 0 {
		return layout.Dimensions{}
	}

	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding14,
				}.Layout(gtx, func(gtx C) D {
					tit := c.theme.Label(values.TextSize14, "Recent Activity")
					tit.Color = c.theme.Color.Gray2
					return pg.titleRow(gtx, tit.Layout, pg.toTicketsActivity.Layout)
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
								return endToEndRow(gtx,
									pg.th.Label(values.TextSize18, strings.Title(strings.ToLower(t.Info.Status))).Layout,
									pg.th.Label(values.TextSize14, t.DaysBehind).Layout)
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
										return layout.Inset{
											Right: values.MarginPadding4,
										}.Layout(gtx, func(gtx C) D {
											ic := c.icons.ticketIconInactive
											ic.Scale = 0.5
											return ic.Layout(gtx)
										})
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
	return pg.pageSections(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding14,
				}.Layout(gtx, func(gtx C) D {
					tit := c.theme.Label(values.TextSize14, "Staking Record")
					tit.Color = c.theme.Color.Gray2
					return pg.titleRow(gtx, tit.Layout, func(gtx C) D { return layout.Dimensions{} })
				})
			}),
			layout.Rigid(func(gtx C) D {
				if *pg.tickets == nil {
					return layout.Dimensions{}
				}
				stackingRecords := (*pg.tickets).StackingRecordCounter
				return decredmaterial.GridWrap{
					Axis:      layout.Horizontal,
					Alignment: layout.End,
				}.Layout(gtx, len(stackingRecords), func(gtx layout.Context, i int) layout.Dimensions {
					item := stackingRecords[i]
					gtx.Constraints.Min.X = int(gtx.Metric.PxPerDp) * 118

					return layout.Inset{Bottom: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								ic := pg.ticketStatusIc[item.Status]["activity"]
								if ic == nil {
									return layout.Dimensions{}
								}
								ic.Scale = 1.0
								return ic.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return pg.th.Label(values.TextSize16, fmt.Sprintf("%d", item.Count)).Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Right: values.MarginPadding40}.Layout(gtx, func(gtx C) D {
												txt := pg.th.Label(values.TextSize12, strings.Title(strings.ToLower(item.Status)))
												txt.Color = pg.th.Color.Gray2
												return txt.Layout(gtx)
											})
										}),
									)
								})
							}),
						)
					})
				})
			}),
			layout.Rigid(func(gtx C) D {
				wrapper := pg.th.Card()
				wrapper.Color = pg.th.Color.Success2
				return wrapper.Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Inset{
							Top:    values.MarginPadding16,
							Bottom: values.MarginPadding16,
						}.Layout(gtx, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Bottom: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
										txt := pg.th.Label(values.TextSize14, "Rewards Earned")
										txt.Color = pg.th.Color.Success
										return txt.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											ic := c.icons.stakeyIcon
											ic.Scale = 1.0
											return ic.Layout(gtx)
										}),
										layout.Rigid(func(gtx C) D {
											return c.layoutBalance(gtx, "16.5112316")
										}),
									)
								}),
							)
						})
					})
				})
			}),
		)
	})
}

func (pg *ticketPage) purchaseModal(gtx layout.Context, c pageCommon) layout.Dimensions {
	return pg.purchaseOptions.Layout(gtx, []func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								ic := c.icons.ti.ticketPurchasedIcon
								ic.Scale = 1.0
								return ic.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
									return c.layoutBalance(gtx, pg.ticketPrice)
								})
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							tit := pg.th.Label(values.TextSize14, "Total")
							tit.Color = pg.th.Color.Gray2
							return tit.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return pg.th.Label(values.TextSize16, pg.ticketPrice).Layout(gtx)
						}),
					)
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return c.accountSelectorLayout(gtx, "purchase", "")
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
						return pg.vspHostSelectorLayout(gtx, c)
					})
				}),
			)
		},
		func(gtx C) D {
			return pg.rememberVSP.Layout(gtx)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
							return pg.cancelPurchase.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return pg.reviewPurchase.Layout(gtx)
					}),
				)
			})
		},
	}, 900)
}

func (pg *ticketPage) confirmPurchaseModal(gtx layout.Context) layout.Dimensions {
	return pg.purchaseOptions.Layout(gtx, []func(gtx C) D{
		func(gtx C) D {
			return pg.th.Label(values.TextSize20, "Confirm to purchase tickets").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					tleft := pg.th.Label(values.TextSize14, "Amount")
					tleft.Color = pg.th.Color.Gray2
					tright := pg.th.Label(values.TextSize14, "1")
					return endToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					tleft := pg.th.Label(values.TextSize14, "Total cost")
					tleft.Color = pg.th.Color.Gray2
					tright := pg.th.Label(values.TextSize14, "122")
					return endToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding16,
						Bottom: values.MarginPadding16,
					}.Layout(gtx, pg.th.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					tleft := pg.th.Label(values.TextSize14, "Account")
					tleft.Color = pg.th.Color.Gray2
					tright := pg.th.Label(values.TextSize14, "Default")
					return endToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					tleft := pg.th.Label(values.TextSize14, "Remaining")
					tleft.Color = pg.th.Color.Gray2
					tright := pg.th.Label(values.TextSize14, "122")
					return endToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding16,
						Bottom: values.MarginPadding16,
					}.Layout(gtx, pg.th.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					tleft := pg.th.Label(values.TextSize14, "VSP")
					tleft.Color = pg.th.Color.Gray2
					tright := pg.th.Label(values.TextSize14, "stakey.net")
					return endToEndRow(gtx, tleft.Layout, tright.Layout)
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.spendingPassword.Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
							return pg.cancelConfirmPurchase.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return pg.purchaseTicket.Layout(gtx)
					}),
				)
			})
		},
	}, 900)
}

func (pg *ticketPage) vspHostSelectorLayout(gtx C, c pageCommon) layout.Dimensions {
	border := widget.Border{
		Color:        pg.th.Color.Gray1,
		CornerRadius: values.MarginPadding8,
		Width:        values.MarginPadding2,
	}
	return border.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
			return decredmaterial.Clickable(gtx, pg.selectVSP, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.th.Body1("http://dev.planetdecred.org:23125").Layout(gtx)
					}),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									txt := pg.th.Body1("1%")
									txt.Color = pg.th.Color.DeepBlue
									return txt.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									inset := layout.Inset{
										Left: values.MarginPadding15,
									}
									return inset.Layout(gtx, func(gtx C) D {
										return c.icons.dropDownIcon.Layout(gtx, values.MarginPadding20)
									})
								}),
							)
						})
					}),
				)
			})
		})
	})
}

var hostList = []string{"vsp.stakeminer.com", "dcrvsp.ubiqsmart.com", "vsp.coinmine.pl"}

func (pg *ticketPage) vspHostModalLayout(gtx C) layout.Dimensions {
	return pg.purchaseOptions.Layout(gtx, []func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.th.Label(values.TextSize20, "Voting service provider").Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					txt := pg.th.Label(values.TextSize14, "Address")
					txt.Color = pg.th.Color.Gray2
					txtFee := pg.th.Label(values.TextSize14, "Fee")
					txtFee.Color = pg.th.Color.Gray2
					return endToEndRow(gtx, txt.Layout, txtFee.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.vspHosts.Layout(gtx, len(hostList), func(gtx C, i int) D {
						return layout.Inset{Top: values.MarginPadding12, Bottom: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
							txt := pg.th.Label(values.TextSize14, "1%")
							txt.Color = pg.th.Color.Gray2
							return endToEndRow(gtx, pg.th.Label(values.TextSize16, hostList[i]).Layout, txt.Layout)
						})
					})
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return pg.inputNewVSP.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return pg.addNewVSP.Layout(gtx)
				}),
			)
		},
	}, 900)
}

func (pg *ticketPage) doPurchaseTicket(c pageCommon, password []byte) {
	selectedWallet := c.info.Wallets[c.wallAcctSelector.selectedPurchaseTicketWallet]
	selectedAccount := selectedWallet.Accounts[c.wallAcctSelector.selectedPurchaseTicketAccount]

	vspd := c.wallet.NewVSPD(selectedWallet.ID, selectedAccount.Number, "http://dev.planetdecred.org:23125")
	_, err := vspd.GetInfo()
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
		resp, err := vspd.GetVSPFeeAddress(hash, password)
		if err != nil {
			log.Error("[CreateTicketFeeTx] err:", err)
			return
		}

		transactionResponse, err := vspd.CreateTicketFeeTx(resp.FeeAmount, hash, resp.FeeAddress, password)
		if err != nil {
			log.Error("[CreateTicketFeeTx] err:", err)
			c.notify(err.Error(), false)
			return
		}

		_, err = vspd.PayVSPFee(transactionResponse, hash, "", password)
		if err != nil {
			log.Error("[PayVSPFee] err:", err)
			c.notify(err.Error(), false)
			return
		}
	}

	c.notify("success", true)
	c.closeModal()
}

func (pg *ticketPage) handler(c pageCommon) {
	if len(c.info.Wallets) > 0 && pg.ticketPrice == "" {
		pg.ticketPrice = c.wallet.TicketPrice()
	}

	if pg.purchaseTicketButton.Button.Clicked() {
		if pg.autoPurchaseEnabled.Value {
			pg.showPurchaseConfirm = true
			return
		}
		pg.showPurchaseOptions = true
	}

	if pg.cancelConfirmPurchase.Button.Clicked() {
		pg.showPurchaseConfirm = false
	}

	if pg.purchaseTicket.Button.Clicked() {
		pg.doPurchaseTicket(c, []byte(pg.spendingPassword.Editor.Text()))
	}

	if pg.cancelPurchase.Button.Clicked() {
		pg.showPurchaseOptions = false
	}

	if pg.reviewPurchase.Button.Clicked() {
		pg.showPurchaseConfirm = true
	}

	if pg.selectVSP.Clicked() {
		pg.showVSPHosts = true
	}

	if pg.addNewVSP.Button.Clicked() {
		pg.showVSPHosts = false
	}
}

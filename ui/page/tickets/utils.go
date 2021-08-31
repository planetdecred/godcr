package tickets

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const (
	uint32Size = 32 << (^uint32(0) >> 32 & 1) // 32 or 64
	maxInt32   = 1<<(uint32Size-1) - 1

	ticketAge   = "Ticket age"
	durationMsg = "10 hrs 47 mins (118/256 blocks)"
)

type tooltips struct {
	statusTooltip     *decredmaterial.Tooltip
	walletNameTooltip *decredmaterial.Tooltip
	dateTooltip       *decredmaterial.Tooltip
	daysBehindTooltip *decredmaterial.Tooltip
	durationTooltip   *decredmaterial.Tooltip
}

var (
	title         = ""
	mainMsg       = ""
	mainMsgDesc   = ""
	dayBehind     = ""
	durationTitle = ""
	durationDesc  = ""
)

func ticketStatusIcon(l *load.Load, ticketStatus string) *struct {
	icon       *widget.Image
	color      color.NRGBA
	background color.NRGBA
} {
	m := map[string]struct {
		icon       *widget.Image
		color      color.NRGBA
		background color.NRGBA
	}{
		"UNMINED": {
			l.Icons.TicketUnminedIcon,
			l.Theme.Color.DeepBlue,
			l.Theme.Color.LightBlue,
		},
		"IMMATURE": {
			l.Icons.TicketImmatureIcon,
			l.Theme.Color.DeepBlue,
			l.Theme.Color.LightBlue,
		},
		"LIVE": {
			l.Icons.TicketLiveIcon,
			l.Theme.Color.Primary,
			l.Theme.Color.LightBlue,
		},
		"VOTED": {
			l.Icons.TicketVotedIcon,
			l.Theme.Color.Success,
			l.Theme.Color.Success2,
		},
		"MISSED": {
			l.Icons.TicketMissedIcon,
			l.Theme.Color.Gray,
			l.Theme.Color.LightGray,
		},
		"EXPIRED": {
			l.Icons.TicketExpiredIcon,
			l.Theme.Color.Gray,
			l.Theme.Color.LightGray,
		},
		"REVOKED": {
			l.Icons.TicketRevokedIcon,
			l.Theme.Color.Orange,
			l.Theme.Color.Orange2,
		},
	}
	st, ok := m[ticketStatus]
	if !ok {
		return nil
	}
	return &st
}

func setText(t string) {
	switch t {
	case "UNMINED":
		title = "This ticket is waiting in mempool to be included in a block."
		mainMsg, mainMsgDesc, dayBehind, durationTitle, durationDesc = "", "", ticketAge, "Live in", durationMsg
	case "IMMATURE":
		title = "This ticket will enter the ticket pool and become a live ticket after 256 blocks (~20 hrs)."
		mainMsg, mainMsgDesc, dayBehind, durationTitle, durationDesc = "", "", ticketAge, "Live in", durationMsg
	case "LIVE":
		title = "Waiting to be chosen to vote."
		mainMsg = "The average vote time is 28 days, but can take up to 142 days."
		mainMsgDesc = "There is a 0.5% chance of expiring before being chosen to vote (this expiration returns the original ticket price without a reward)."
		dayBehind, durationTitle, durationDesc = ticketAge, "Live in", durationMsg
	case "VOTED":
		title = "Congratulations! This ticket has voted."
		mainMsg = "The ticket price + reward will become spendable after 256 blocks (~20 hrs)."
		dayBehind, durationTitle, durationDesc = "Days to vote", "Spendable in", durationMsg
	case "MISSED":
		title = "This ticket was chosen to vote, but missed the voting window."
		mainMsg = "Missed tickets will be revoked to return the original ticket price to you."
		mainMsgDesc = "If a ticket is not revoked automatically, use the revoke button."
		dayBehind, durationTitle, durationDesc = "Days to miss", "Miss in", durationMsg
	case "EXPIRED":
		title = "This ticket has not been chosen to vote within 40960 blocks, and thus expired. "
		mainMsg = "Expired tickets will be revoked to return the original ticket price to you."
		mainMsgDesc = "If a ticket is not revoked automatically, use the revoke button."
		dayBehind, durationTitle, durationDesc = "Days to expire", "Expired in", durationMsg
	case "REVOKED":
		title = "This ticket has been revoked."
		dayBehind, durationTitle, durationDesc = ticketAge, "Spendable in", durationMsg
	}
}

func ticketStatusTooltip(gtx C, l *load.Load, t *wallet.Ticket) layout.Dimensions {
	setText(t.Info.Status)
	st := ticketStatusIcon(l, t.Info.Status)
	status := l.Theme.Body2(t.Info.Status)
	status.Color = st.color
	st.icon.Scale = .5

	titleLabel, mainMsgLabel, mainMsgLabel2 := l.Theme.Body2(title), l.Theme.Body2(mainMsg), l.Theme.Body2(mainMsgDesc)
	mainMsgLabel.Color, mainMsgLabel2.Color = l.Theme.Color.Gray, l.Theme.Color.Gray
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(st.icon.Layout),
				layout.Rigid(toolTipContent(layout.Inset{Left: values.MarginPadding4}, status.Layout)),
			)
		}),
		layout.Rigid(toolTipContent(layout.Inset{Top: values.MarginPadding8}, titleLabel.Layout)),
		layout.Rigid(toolTipContent(layout.Inset{Top: values.MarginPadding8}, mainMsgLabel.Layout)),
		layout.Rigid(func(gtx C) D {
			if mainMsgDesc != "" {
				toolTipContent(layout.Inset{Top: values.MarginPadding8}, mainMsgLabel2.Layout)
			}
			return layout.Dimensions{}
		}),
	)
}

func ticketCardTooltip(gtx C, rectLayout layout.Dimensions, tooltip *decredmaterial.Tooltip, body layout.Widget) {
	inset := layout.Inset{
		Top:   values.MarginPadding15,
		Right: unit.Dp(-150),
		Left:  values.MarginPadding15,
	}

	rect := image.Rectangle{
		Max: image.Point{
			X: rectLayout.Size.X,
			Y: rectLayout.Size.Y,
		},
	}

	tooltip.Layout(gtx, rect, inset, body)
}

func walletNameDateTimeTooltip(gtx C, l *load.Load, title string, body layout.Widget) layout.Dimensions {
	walletNameLabel := l.Theme.Body2(title)
	walletNameLabel.Color = l.Theme.Color.Gray

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(walletNameLabel.Layout),
		layout.Rigid(body),
	)
}

func toolTipContent(inset layout.Inset, body layout.Widget) layout.Widget {
	return func(gtx C) D {
		return inset.Layout(gtx, body)
	}
}

// ticketCard layouts out ticket info with the shadow box, use for list horizontal or list grid
func ticketCard(gtx layout.Context, l *load.Load, t *wallet.Ticket, tooltip tooltips) layout.Dimensions {
	var itemWidth int
	st := ticketStatusIcon(l, t.Info.Status)
	if st == nil {
		return layout.Dimensions{}
	}
	st.icon.Scale = 1.0
	return l.Theme.Shadow().Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				wrap := l.Theme.Card()
				wrap.Radius = decredmaterial.CornerRadius{TopRight: 8, TopLeft: 8, BottomRight: 0, BottomLeft: 0}
				wrap.Color = st.background
				return wrap.Layout(gtx, func(gtx C) D {
					return layout.Stack{Alignment: layout.S}.Layout(gtx,
						layout.Expanded(func(gtx C) D {
							return layout.NE.Layout(gtx, func(gtx C) D {
								wTimeLabel := l.Theme.Card()
								wTimeLabel.Radius = decredmaterial.CornerRadius{TopRight: 8, TopLeft: 0, BottomRight: 0, BottomLeft: 8}
								return wTimeLabel.Layout(gtx, func(gtx C) D {
									return layout.Inset{
										Top:    values.MarginPadding4,
										Bottom: values.MarginPadding4,
										Right:  values.MarginPadding8,
										Left:   values.MarginPadding8,
									}.Layout(gtx, func(gtx C) D {
										txt := l.Theme.Label(values.TextSize14, "10h 47m")
										txtLayout := txt.Layout(gtx)
										ticketCardTooltip(gtx, txtLayout, tooltip.durationTooltip, func(gtx C) D {
											setText(t.Info.Status)
											return walletNameDateTimeTooltip(gtx, l, durationTitle,
												toolTipContent(layout.Inset{Top: values.MarginPadding8}, l.Theme.Body2(durationMsg).Layout))
										})
										return txtLayout
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
								return st.icon.Layout(gtx)
							})
							itemWidth = content.Size.X
							return content
						}),

						layout.Stacked(func(gtx C) D {
							return layout.Center.Layout(gtx, func(gtx C) D {
								return layout.Inset{Top: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
									gtx.Constraints.Max.X = itemWidth
									p := l.Theme.ProgressBar(20)
									p.Height, p.Radius = values.MarginPadding4, values.MarginPadding1
									p.Color = st.color
									return p.Layout(gtx)
								})
							})
						}),
					)
				})
			}),
			layout.Rigid(func(gtx C) D {
				wrap := l.Theme.Card()
				wrap.Radius = decredmaterial.CornerRadius{TopRight: 0, TopLeft: 0, BottomRight: 8, BottomLeft: 8}
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
									return components.LayoutBalance(gtx, l, t.Amount)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										txt := l.Theme.Label(values.MarginPadding14, t.Info.Status)
										txt.Color = st.color
										txtLayout := txt.Layout(gtx)
										ticketCardTooltip(gtx, txtLayout, tooltip.statusTooltip, func(gtx C) D {
											setText(t.Info.Status)
											return ticketStatusTooltip(gtx, l, t)
										})
										return txtLayout
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{
											Left:  values.MarginPadding4,
											Right: values.MarginPadding4,
										}.Layout(gtx, func(gtx C) D {
											ic := l.Icons.ImageBrightness1
											ic.Color = l.Theme.Color.Gray2
											return l.Icons.ImageBrightness1.Layout(gtx, values.MarginPadding5)
										})
									}),
									layout.Rigid(func(gtx C) D {
										txt := l.Theme.Label(values.MarginPadding14, t.WalletName)
										txt.Color = l.Theme.Color.Gray
										txtLayout := txt.Layout(gtx)
										ticketCardTooltip(gtx, txtLayout, tooltip.walletNameTooltip, func(gtx C) D {
											return walletNameDateTimeTooltip(gtx, l, "Wallet name",
												toolTipContent(layout.Inset{Top: values.MarginPadding8}, l.Theme.Body2(t.WalletName).Layout))
										})
										return txtLayout
									}),
								)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Top:    values.MarginPadding16,
									Bottom: values.MarginPadding16,
								}.Layout(gtx, func(gtx C) D {
									txt := l.Theme.Label(values.TextSize14, t.MonthDay)
									txt.Color = l.Theme.Color.Gray2
									return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											txtLayout := txt.Layout(gtx)
											ticketCardTooltip(gtx, txtLayout, tooltip.dateTooltip, func(gtx C) D {
												dt := strings.Split(t.DateTime, " ")
												s1 := []string{dt[0], dt[1], dt[2]}
												date := strings.Join(s1, " ")
												s2 := []string{dt[3], dt[4]}
												time := strings.Join(s2, " ")
												dateTime := fmt.Sprintf("%s at %s", date, time)
												return walletNameDateTimeTooltip(gtx, l, "Purchased",
													toolTipContent(layout.Inset{Top: values.MarginPadding8}, l.Theme.Body2(dateTime).Layout))
											})
											return txtLayout
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{
												Left:  values.MarginPadding4,
												Right: values.MarginPadding4,
											}.Layout(gtx, func(gtx C) D {
												ic := l.Icons.ImageBrightness1
												ic.Color = l.Theme.Color.Gray2
												return l.Icons.ImageBrightness1.Layout(gtx, values.MarginPadding5)
											})
										}),
										layout.Rigid(func(gtx C) D {
											txt.Text = t.DaysBehind
											txtLayout := txt.Layout(gtx)
											ticketCardTooltip(gtx, txtLayout, tooltip.daysBehindTooltip, func(gtx C) D {
												setText(t.Info.Status)
												return walletNameDateTimeTooltip(gtx, l, dayBehind,
													toolTipContent(layout.Inset{Top: values.MarginPadding8}, l.Theme.Body2(t.DaysBehind).Layout))
											})
											return txtLayout
										}),
									)
								})
							}),
						)
					})
				})
			}),
		)
	})
}

// ticketActivityRow layouts out ticket info, display ticket activities on the tickets_page and tickets_activity_page
func ticketActivityRow(gtx layout.Context, l *load.Load, t wallet.Ticket, index int) layout.Dimensions {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Right: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
				st := ticketStatusIcon(l, t.Info.Status)
				if st == nil {
					return layout.Dimensions{}
				}
				st.icon.Scale = 0.6
				return st.icon.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if index == 0 {
						return layout.Dimensions{}
					}
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					separator := l.Theme.Separator()
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
								labelStatus := l.Theme.Label(values.TextSize18, strings.Title(strings.ToLower(t.Info.Status)))
								labelStatus.Color = l.Theme.Color.DeepBlue

								labelDaysBehind := l.Theme.Label(values.TextSize14, t.DaysBehind)
								labelDaysBehind.Color = l.Theme.Color.DeepBlue

								return components.EndToEndRow(gtx,
									labelStatus.Layout,
									labelDaysBehind.Layout)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Flex{
									Alignment: layout.Middle,
								}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										txt := l.Theme.Label(values.TextSize14, t.WalletName)
										txt.Color = l.Theme.Color.Gray2
										return txt.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{
											Left:  values.MarginPadding4,
											Right: values.MarginPadding4,
										}.Layout(gtx, func(gtx C) D {
											ic := l.Icons.ImageBrightness1
											ic.Color = l.Theme.Color.Gray2
											return l.Icons.ImageBrightness1.Layout(gtx, values.MarginPadding5)
										})
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{
											Right: values.MarginPadding4,
										}.Layout(gtx, func(gtx C) D {
											ic := l.Icons.TicketIconInactive
											ic.Scale = 0.5
											return ic.Layout(gtx)
										})
									}),
									layout.Rigid(func(gtx C) D {
										txt := l.Theme.Label(values.TextSize14, t.Amount)
										txt.Color = l.Theme.Color.Gray2
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

// todo: cleanup
func createOrderDropDown(th *decredmaterial.Theme) *decredmaterial.DropDown {
	return th.DropDown([]decredmaterial.DropDownItem{{Text: values.String(values.StrNewest)},
		{Text: values.String(values.StrOldest)}}, 1)
}

// todo: cleanup
// createClickGestures returns a slice of click gestures
func createClickGestures(count int) []*gesture.Click {
	var gestures = make([]*gesture.Click, count)
	for i := 0; i < count; i++ {
		gestures[i] = &gesture.Click{}
	}
	return gestures
}

func nextTicketRemaining(allsecs int) string {
	if allsecs == 0 {
		return "imminent"
	}
	str := ""
	if allsecs > 604799 {
		weeks := allsecs / 604800
		allsecs %= 604800
		str += fmt.Sprintf("%dw ", weeks)
	}
	if allsecs > 86399 {
		days := allsecs / 86400
		allsecs %= 86400
		str += fmt.Sprintf("%dd ", days)
	}
	if allsecs > 3599 {
		hours := allsecs / 3600
		allsecs %= 3600
		str += fmt.Sprintf("%dh ", hours)
	}
	if allsecs > 59 {
		mins := allsecs / 60
		allsecs %= 60
		str += fmt.Sprintf("%dm ", mins)
	}
	if allsecs > 0 {
		str += fmt.Sprintf("%ds ", allsecs)
	}
	return str
}

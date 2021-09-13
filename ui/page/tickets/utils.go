package tickets

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"
	"time"

	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	uint32Size = 32 << (^uint32(0) >> 32 & 1) // 32 or 64
	maxInt32   = 1<<(uint32Size-1) - 1
)

type transactionItem struct {
	transaction   *dcrlibwallet.Transaction
	ticketSpender *dcrlibwallet.Transaction
	status        *components.TxStatus
	confirmations int32
	progress      float32
	showProgress  bool
	showTime      bool
	purchaseTime  string
	ticketAge     string

	statusTooltip     *decredmaterial.Tooltip
	dateTooltip       *decredmaterial.Tooltip
	daysBehindTooltip *decredmaterial.Tooltip
	durationTooltip   *decredmaterial.Tooltip
}

type Ticket struct {
	Status     string
	Fee        string
	Amount     string
	DateTime   string
	MonthDay   string
	DaysBehind string
	WalletName string
}

const (
	StakingLive     = "LIVE"
	StakingUnmined  = "UNMINED"
	StakingImmature = "IMMATURE"
	StakingRevoked  = "REVOKED"
	StakingVoted    = "VOTED"
	StakingExpired  = "EXPIRED"
)

func calculateDaysBehind(lastHeaderTime int64) string {
	diff := time.Since(time.Unix(lastHeaderTime, 0))
	daysBehind := int(math.Round(diff.Hours() / 24))
	if daysBehind < 1 {
		return "<1 day"
	} else if daysBehind == 1 {
		return "1 day"
	} else {
		return fmt.Sprintf("%d days", daysBehind)
	}
}

func transactionToTicket(tx dcrlibwallet.Transaction, w *dcrlibwallet.Wallet, maturity, expiry, bestBlock int32) Ticket {
	return Ticket{
		Status:     getTicketStatus(tx, w, maturity, expiry, bestBlock),
		Amount:     dcrutil.Amount(tx.Amount).String(),
		DateTime:   time.Unix(tx.Timestamp, 0).Format("Jan 2, 2006 03:04:05 PM"),
		MonthDay:   time.Unix(tx.Timestamp, 0).Format("Jan 2"),
		DaysBehind: calculateDaysBehind(tx.Timestamp),
		Fee:        dcrutil.Amount(tx.Fee).String(),
		WalletName: w.Name,
	}
}

func getTicketStatus(txn dcrlibwallet.Transaction, w *dcrlibwallet.Wallet, ticketMaturity, ticketExpiry, bestBlock int32) string {
	if txn.Type == dcrlibwallet.TxTypeVote {
		return StakingVoted
	}

	if txn.Type == dcrlibwallet.TxTypeRevocation {
		return StakingRevoked
	}

	s := txn.TicketStatus(ticketMaturity, ticketExpiry, bestBlock)
	switch s {
	case dcrlibwallet.TicketStatusUnmined:
		return StakingUnmined
	case dcrlibwallet.TicketStatusImmature:
		return StakingImmature
	case dcrlibwallet.TicketStatusLive:
		return StakingLive
	case dcrlibwallet.TicketStatusVotedOrRevoked:
		// handle revocation and voted tickets that have the type "TicketPurchase"
		tx, _ := w.TicketSpender(txn.Hash)
		if tx.Type == dcrlibwallet.TxTypeVote {
			return StakingVoted
		}

		if tx.Type == dcrlibwallet.TxTypeRevocation {
			return StakingRevoked
		}
	}

	return ""
}

func allLiveTickets(mw *dcrlibwallet.MultiWallet) ([]dcrlibwallet.Transaction, error) {
	var tickets []dcrlibwallet.Transaction
	liveTicketFilters := []int32{dcrlibwallet.TxFilterUnmined, dcrlibwallet.TxFilterImmature, dcrlibwallet.TxFilterLive}
	for _, filter := range liveTicketFilters {
		tx, err := mw.GetTransactionsRaw(0, 0, filter, true)
		if err != nil {
			return nil, err
		}

		tickets = append(tickets, tx...)
	}

	return tickets, nil
}

func ticketStatusProfile(l *load.Load, ticketStatus string) *struct {
	icon       *decredmaterial.Image
	color      color.NRGBA
	background color.NRGBA
} {
	m := map[string]struct {
		icon       *decredmaterial.Image
		color      color.NRGBA
		background color.NRGBA
	}{
		StakingUnmined: {
			l.Icons.TicketUnminedIcon,
			l.Theme.Color.DeepBlue,
			l.Theme.Color.LightBlue,
		},
		StakingImmature: {
			l.Icons.TicketImmatureIcon,
			l.Theme.Color.DeepBlue,
			l.Theme.Color.LightBlue,
		},
		StakingLive: {
			l.Icons.TicketLiveIcon,
			l.Theme.Color.Primary,
			l.Theme.Color.LightBlue,
		},
		StakingVoted: {
			l.Icons.TicketVotedIcon,
			l.Theme.Color.Success,
			l.Theme.Color.Success2,
		},
		"MISSED": {
			l.Icons.TicketMissedIcon,
			l.Theme.Color.Gray,
			l.Theme.Color.LightGray,
		},
		StakingExpired: {
			l.Icons.TicketExpiredIcon,
			l.Theme.Color.Gray,
			l.Theme.Color.LightGray,
		},
		StakingRevoked: {
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

func ticketStatusTooltip(gtx C, l *load.Load, tx *transactionItem) layout.Dimensions {
	status := l.Theme.Label(values.MarginPadding14, strings.ToUpper(tx.status.Title))
	status.Font.Weight = text.Medium
	status.Color = tx.status.Color

	maturity := l.WL.MultiWallet.TicketMaturity()
	blockTime := l.WL.MultiWallet.TargetTimePerBlockMinutes()
	maturityDuration := time.Duration(maturity*int32(blockTime)) * time.Minute
	maturityTime := ticketAgeTimeFormat(int(maturityDuration.Seconds()))
	var title, mainDesc, subDesc string
	switch tx.status.TicketStatus {
	case dcrlibwallet.TicketStatusUnmined:
		title = "This ticket is waiting in mempool to be included in a block."
	case dcrlibwallet.TicketStatusImmature:
		title = fmt.Sprintf("This ticket will enter the ticket pool and become a live ticket after %d blocks (~%s).", maturity, maturityTime)
	case dcrlibwallet.TicketStatusLive:
		title = "Waiting to be chosen to vote."
		mainDesc = "The average vote time is 28 days, but can take up to 142 days."
		subDesc = "There is a 0.5% chance of expiring before being chosen to vote (this expiration returns the original ticket price without a reward)."
	case dcrlibwallet.TicketStatusVotedOrRevoked:
		if tx.ticketSpender.Type == dcrlibwallet.TxTypeVote {
			title = "Congratulations! This ticket has voted."
			mainDesc = "The ticket price + reward will become spendable after %d blocks (~%s)."
		} else {
			title = "This ticket has been revoked."
			mainDesc = "The ticket price will become spendable after %d blocks (~%s)."
		}

		if tx.ticketSpender.Confirmations(l.WL.MultiWallet.GetBestBlock().Height) > maturity {
			mainDesc = ""
		} else {

			mainDesc = fmt.Sprintf(mainDesc, maturity, maturityTime)
		}
	case dcrlibwallet.TicketStatusExpired:
		title = fmt.Sprintf("This ticket has not been chosen to vote within %d blocks, and thus expired.", l.WL.MultiWallet.TicketMaturity())
		mainDesc = "Expired tickets will be revoked to return the original ticket price to you."
		subDesc = "If a ticket is not revoked automatically, use the revoke button."
	}

	titleLabel := l.Theme.Label(values.MarginPadding14, title)
	titleLabel.Color = l.Theme.Color.DeepBlue

	mainDescLabel := l.Theme.Label(values.MarginPadding14, mainDesc)
	mainDescLabel.Color = l.Theme.Color.Gray
	subDescLabel := l.Theme.Label(values.MarginPadding14, subDesc)
	subDescLabel.Color = l.Theme.Color.Gray
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(tx.status.Icon.Layout16dp),
				layout.Rigid(toolTipContent(layout.Inset{Left: values.MarginPadding4}, status.Layout)),
			)
		}),
		layout.Rigid(toolTipContent(layout.Inset{Top: values.MarginPadding8}, titleLabel.Layout)),
		layout.Rigid(func(gtx C) D {
			if mainDesc != "" {
				return toolTipContent(layout.Inset{Top: values.MarginPadding8}, mainDescLabel.Layout)(gtx)
			}

			return D{}
		}),
		layout.Rigid(func(gtx C) D {
			if subDesc != "" {
				return toolTipContent(layout.Inset{Top: values.MarginPadding8}, subDescLabel.Layout)(gtx)
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

func titleDescTooltip(gtx C, l *load.Load, title string, desc string) layout.Dimensions {
	titleLabel := l.Theme.Label(values.MarginPadding14, title)
	titleLabel.Color = l.Theme.Color.Gray

	descLabel := l.Theme.Label(values.MarginPadding14, desc)
	descLabel.Color = l.Theme.Color.DeepBlue

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(titleLabel.Layout),
		layout.Rigid(toolTipContent(layout.Inset{Top: values.MarginPadding4}, descLabel.Layout)),
	)
}

func toolTipContent(inset layout.Inset, body layout.Widget) layout.Widget {
	return func(gtx C) D {
		return inset.Layout(gtx, body)
	}
}

// ticketCard layouts out ticket info with the shadow box, use for list horizontal or list grid
func ticketCard(gtx layout.Context, l *load.Load, tx *transactionItem, showWalletName bool) layout.Dimensions {
	wal := l.WL.MultiWallet.WalletWithID(tx.transaction.WalletID)
	txStatus := tx.status

	// add this data to transactionItem so it can be shared with list
	maturity := l.WL.MultiWallet.TicketMaturity()

	return decredmaterial.LinearLayout{
		Width:       gtx.Px(values.MarginPadding168),
		Height:      decredmaterial.WrapContent,
		Orientation: layout.Vertical,
		Shadow:      l.Theme.Shadow(),
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return decredmaterial.LinearLayout{
				Width:      decredmaterial.MatchParent,
				Height:     decredmaterial.WrapContent,
				Background: txStatus.Background,
				Border:     decredmaterial.Border{Radius: decredmaterial.TopRadius(8)},
			}.Layout2(gtx, func(gtx C) D {

				return layout.Stack{}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Center.Layout(gtx, func(gtx C) D {
							return layout.Inset{
								Top:    values.MarginPadding24,
								Bottom: values.MarginPadding24,
							}.Layout(gtx, txStatus.Icon.Layout36dp)
						})
					}),
					layout.Expanded(func(gtx C) D {
						if !tx.showTime {
							return D{}
						}

						return layout.NE.Layout(gtx, func(gtx C) D {
							timeWrapper := l.Theme.Card()
							timeWrapper.Radius = decredmaterial.CornerRadius{TopRight: 8, TopLeft: 0, BottomRight: 0, BottomLeft: 8}
							return timeWrapper.Layout(gtx, func(gtx C) D {
								return layout.Inset{
									Top:    values.MarginPadding4,
									Bottom: values.MarginPadding4,
									Right:  values.MarginPadding8,
									Left:   values.MarginPadding8,
								}.Layout(gtx, func(gtx C) D {

									confirmations := tx.confirmations
									if tx.ticketSpender != nil {
										confirmations = tx.ticketSpender.Confirmations(wal.GetBestBlock())
									}

									timeRemaining := time.Duration(float64(maturity-confirmations)*l.WL.MultiWallet.TargetTimePerBlockMinutes()) * time.Minute
									maturityDuration := ticketAgeTimeFormat(int(timeRemaining.Seconds()))
									txt := l.Theme.Label(values.TextSize14, maturityTimeFormat(int(timeRemaining.Minutes())))

									durationLayout := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, l.Icons.TimerIcon.Layout)
										}),
										layout.Rigid(txt.Layout),
									)

									var durationTitle = "Live in" // immature
									if txStatus.TicketStatus != dcrlibwallet.TicketStatusImmature {
										// voted or revoked
										durationTitle = "Spendable in"
									}

									ticketCardTooltip(gtx, durationLayout, tx.durationTooltip, func(gtx C) D {
										return titleDescTooltip(gtx, l, durationTitle, fmt.Sprintf("%s (%d/%d blocks)", maturityDuration, confirmations, maturity))
									})

									return durationLayout
								})
							})
						})
					}),
					layout.Expanded(func(gtx C) D {
						return layout.S.Layout(gtx, func(gtx C) D {

							if !tx.showProgress {
								return D{}
							}

							p := l.Theme.ProgressBar(int(tx.progress))
							p.Height = values.MarginPadding4
							p.Color = txStatus.ProgressBarColor
							p.TrackColor = txStatus.ProgressTrackColor
							return p.Layout2(gtx)
						})
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return decredmaterial.LinearLayout{
				Width:       decredmaterial.MatchParent,
				Height:      decredmaterial.WrapContent,
				Orientation: layout.Vertical,
				Padding:     layout.UniformInset(values.MarginPadding16),
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return components.LayoutBalance(gtx, l, dcrutil.Amount(tx.transaction.Amount).String())
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding4,
						Bottom: values.MarginPadding8,
					}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								txt := l.Theme.Label(values.MarginPadding14, txStatus.Title)
								txt.Color = txStatus.Color
								txt.Font.Weight = text.Medium
								txtLayout := txt.Layout(gtx)
								ticketCardTooltip(gtx, txtLayout, tx.statusTooltip, func(gtx C) D {
									return ticketStatusTooltip(gtx, l, tx)
								})
								return txtLayout
							}),
							layout.Rigid(func(gtx C) D {
								if !showWalletName {
									return D{}
								}

								return layout.Inset{
									Left:  values.MarginPadding4,
									Right: values.MarginPadding4,
								}.Layout(gtx, func(gtx C) D {
									txt := l.Theme.Label(values.MarginPadding14, "•")
									txt.Color = l.Theme.Color.Gray

									return txt.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								if !showWalletName {
									return D{}
								}

								txt := l.Theme.Label(values.MarginPadding14, wal.Name)
								txt.Color = l.Theme.Color.Gray
								return txt.Layout(gtx)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					txt := l.Theme.Label(values.TextSize14, time.Unix(tx.transaction.Timestamp, 0).Format("Jan 2"))
					txt.Color = l.Theme.Color.Gray2
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txtLayout := txt.Layout(gtx)
							ticketCardTooltip(gtx, txtLayout, tx.dateTooltip, func(gtx C) D {
								dateTime := time.Unix(tx.transaction.Timestamp, 0).Format("Jan 2, 2006 at 03:04:05 PM")
								return titleDescTooltip(gtx, l, "Purchased", dateTime)
							})
							return txtLayout
						}),
						layout.Rigid(func(gtx C) D {

							var tooltipTitle string
							if tx.ticketSpender != nil { // voted or revoked
								if tx.ticketSpender.Type == dcrlibwallet.TxTypeVote {
									tooltipTitle = "Days to vote"
								} else {
									tooltipTitle = "Days to miss"
								}
							} else if txStatus.TicketStatus == dcrlibwallet.TicketStatusImmature ||
								txStatus.TicketStatus == dcrlibwallet.TicketStatusLive {
								tooltipTitle = "Ticket age"
							} else {
								return D{}
							}

							return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
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

									txt.Text = tx.ticketAge
									txtLayout := txt.Layout(gtx)
									ticketCardTooltip(gtx, txtLayout, tx.daysBehindTooltip, func(gtx C) D {
										return titleDescTooltip(gtx, l, tooltipTitle, tx.ticketAge)
									})
									return txtLayout
								}),
							)
						}),
					)
				}),
			)
		}),
	)
}

// ticketActivityRow layouts out ticket info, display ticket activities on the tickets_page and tickets_activity_page
func ticketActivityRow(gtx layout.Context, l *load.Load, t Ticket, index int) layout.Dimensions {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Right: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
				st := ticketStatusProfile(l, t.Status)
				if st == nil {
					return layout.Dimensions{}
				}
				return st.icon.Layout24dp(gtx)
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
								labelStatus := l.Theme.Label(values.TextSize18, strings.Title(strings.ToLower(t.Status)))
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
											return ic.Layout12dp(gtx)
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

func ticketAgeTimeFormat(secs int) string {
	if secs > 86399 {
		days := secs / 86400
		return fmt.Sprintf("%dd", days)
	} else if secs > 3599 {
		hours := secs / 3600
		return fmt.Sprintf("%dh", hours)
	} else if secs > 59 {
		mins := secs / 60
		return fmt.Sprintf("%dm", mins)
	}

	return fmt.Sprintf("%ds", secs)

}

func maturityTimeFormat(maturityTimeMinutes int) string {
	return fmt.Sprintf("%02d:%02d", maturityTimeMinutes/60, maturityTimeMinutes%60)
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

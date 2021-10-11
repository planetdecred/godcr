package tickets

import (
	"fmt"
	"image"
	"sort"
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
	walletNameTooltip *decredmaterial.Tooltip
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

func ticketsToTransactionItems(l *load.Load, txs []dcrlibwallet.Transaction, newestFirst bool, hasFilter func(int32) bool) ([]*transactionItem, error) {
	tickets := make([]*transactionItem, 0)
	multiWallet := l.WL.MultiWallet
	for _, tx := range txs {
		w := multiWallet.WalletWithID(tx.WalletID)

		ticketSpender, err := w.TicketSpender(tx.Hash)
		if err != nil {
			return nil, err
		}

		// Apply voted and revoked tx filter
		if (hasFilter(dcrlibwallet.TxFilterVoted) || hasFilter(dcrlibwallet.TxFilterRevoked)) && ticketSpender == nil {
			continue
		}

		if hasFilter(dcrlibwallet.TxFilterVoted) && ticketSpender.Type != dcrlibwallet.TxTypeVote {
			continue
		}

		if hasFilter(dcrlibwallet.TxFilterRevoked) && ticketSpender.Type != dcrlibwallet.TxTypeRevocation {
			continue
		}

		// This fixes a dcrlibwallet bug where live tickets transactions
		// do not have updated data of ticket spender.
		if hasFilter(dcrlibwallet.TxFilterLive) && ticketSpender != nil {
			continue
		}

		ticketCopy := tx
		txStatus := components.TransactionTitleIcon(l, w, &tx, ticketSpender)
		confirmations := tx.Confirmations(w.GetBestBlock())
		var ticketAge string

		showProgress := txStatus.TicketStatus == dcrlibwallet.TicketStatusImmature || txStatus.TicketStatus == dcrlibwallet.TicketStatusLive
		if ticketSpender != nil { /// voted or revoked
			showProgress = ticketSpender.Confirmations(w.GetBestBlock()) <= multiWallet.TicketMaturity()
			ticketAge = fmt.Sprintf("%d days", ticketSpender.DaysToVoteOrRevoke)
		} else if txStatus.TicketStatus == dcrlibwallet.TicketStatusImmature ||
			txStatus.TicketStatus == dcrlibwallet.TicketStatusLive {

			ticketAgeDuration := time.Since(time.Unix(tx.Timestamp, 0)).Seconds()
			ticketAge = components.TimeFormat(int(ticketAgeDuration), false)
		}

		showTime := showProgress && txStatus.TicketStatus != dcrlibwallet.TicketStatusLive

		var progress float32
		if showProgress {
			progressMax := multiWallet.TicketMaturity()
			if txStatus.TicketStatus == dcrlibwallet.TicketStatusLive {
				progressMax = multiWallet.TicketExpiry()
			}

			confs := confirmations
			if ticketSpender != nil {
				confs = ticketSpender.Confirmations(w.GetBestBlock())
			}

			progress = (float32(confs) / float32(progressMax)) * 100
		}

		tickets = append(tickets, &transactionItem{
			transaction:   &ticketCopy,
			ticketSpender: ticketSpender,
			status:        txStatus,
			confirmations: tx.Confirmations(w.GetBestBlock()),
			progress:      progress,
			showProgress:  showProgress,
			showTime:      showTime,
			purchaseTime:  time.Unix(tx.Timestamp, 0).Format("Jan 2"),
			ticketAge:     ticketAge,

			statusTooltip:     l.Theme.Tooltip(),
			walletNameTooltip: l.Theme.Tooltip(),
			dateTooltip:       l.Theme.Tooltip(),
			daysBehindTooltip: l.Theme.Tooltip(),
			durationTooltip:   l.Theme.Tooltip(),
		})
	}

	// bring vote and revoke tx forward
	sort.Slice(tickets[:], func(i, j int) bool {
		var timeStampI = tickets[i].transaction.Timestamp
		var timeStampJ = tickets[j].transaction.Timestamp

		if tickets[i].ticketSpender != nil {
			timeStampI = tickets[i].ticketSpender.Timestamp
		}

		if tickets[j].ticketSpender != nil {
			timeStampJ = tickets[j].ticketSpender.Timestamp
		}

		if newestFirst {
			return timeStampI > timeStampJ
		}
		return timeStampI < timeStampJ
	})

	return tickets, nil
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

func ticketStatusTooltip(gtx C, l *load.Load, tx *transactionItem) layout.Dimensions {
	status := l.Theme.Label(values.MarginPadding14, strings.ToUpper(tx.status.Title))
	status.Font.Weight = text.Medium
	status.Color = tx.status.Color

	maturity := l.WL.MultiWallet.TicketMaturity()
	blockTime := l.WL.MultiWallet.TargetTimePerBlockMinutes()
	maturityDuration := time.Duration(maturity*int32(blockTime)) * time.Minute
	maturityTime := components.TimeFormat(int(maturityDuration.Seconds()), false)
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
									maturityDuration := components.TimeFormat(int(timeRemaining.Seconds()), false)
									txt := l.Theme.Label(values.TextSize14, maturityTimeFormat(int(timeRemaining.Minutes())))

									durationLayout := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, l.Icons.TimerIcon.Layout12dp)
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

								txt := l.Theme.Label(values.TextSize14, wal.Name)
								txt.Color = l.Theme.Color.Gray
								txtLayout := txt.Layout(gtx)
								ticketCardTooltip(gtx, txtLayout, tx.walletNameTooltip, func(gtx C) D {
									return titleDescTooltip(gtx, l, "Wallet name", txt.Text)
								})
								return txtLayout
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
										ic := decredmaterial.NewIcon(l.Icons.ImageBrightness1)
										ic.Size = 5
										return ic.Layout(gtx)
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

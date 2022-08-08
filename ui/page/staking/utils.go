package staking

import (
	"fmt"
	"image"
	"sort"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/decred/dcrd/dcrutil/v4"
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

type Stake struct {
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

func stakeToTransactionItems(l *load.Load, txs []dcrlibwallet.Transaction, newestFirst bool, hasFilter func(int32) bool) ([]*transactionItem, error) {
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
		// do not have updated data of Stake spender.
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
			ticketAge = components.TimeFormat(int(ticketAgeDuration), true)
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
			purchaseTime:  time.Unix(tx.Timestamp, 0).Format("Jan 2, 2006 15:04:05 PM"),
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

func ticketStatusTooltip(gtx C, l *load.Load, tx *transactionItem) D {
	status := l.Theme.Label(values.TextSize14, strings.ToUpper(tx.status.Title))
	status.Font.Weight = text.Medium
	status.Color = tx.status.Color

	maturity := l.WL.MultiWallet.TicketMaturity()
	blockTime := l.WL.MultiWallet.TargetTimePerBlockMinutes()
	maturityDuration := time.Duration(maturity*int32(blockTime)) * time.Minute

	var title, mainDesc, subDesc string
	switch tx.status.TicketStatus {
	case dcrlibwallet.TicketStatusUnmined:
		title = values.String(values.StrUnminedInfo)
	case dcrlibwallet.TicketStatusImmature:
		title = values.StringF(values.StrImmatureInfo, maturity, maturityDuration)
	case dcrlibwallet.TicketStatusLive:
		title = values.String(values.StrLiveInfo)
		mainDesc = values.String(values.StrLiveInfoDisc)
		subDesc = values.String(values.StrLiveInfoDiscSub)
	case dcrlibwallet.TicketStatusVotedOrRevoked:
		if tx.ticketSpender.Type == dcrlibwallet.TxTypeVote {
			title = values.String(values.StrVotedInfo)
			mainDesc = values.String(values.StrVotedInfoDisc)
		} else {
			title = values.String(values.StrRevokeInfo)
			mainDesc = values.String(values.StrRevokeInfoDisc)
		}

		if tx.ticketSpender.Confirmations(l.WL.MultiWallet.GetBestBlock().Height) > maturity {
			mainDesc = ""
		} else {
			mainDesc = values.StringF(mainDesc, maturity, maturityDuration)
		}
	case dcrlibwallet.TicketStatusExpired:
		title = values.StringF(values.StrExpiredInfo, l.WL.MultiWallet.TicketMaturity())
		mainDesc = values.String(values.StrExpiredInfoDisc)
		subDesc = values.String(values.StrExpiredInfoDiscSub)
	}

	titleLabel := l.Theme.Label(values.TextSize14, title)

	mainDescLabel := l.Theme.Label(values.TextSize14, mainDesc)
	mainDescLabel.Color = l.Theme.Color.GrayText2
	subDescLabel := l.Theme.Label(values.TextSize14, subDesc)
	subDescLabel.Color = l.Theme.Color.GrayText2
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

func ticketStatusDetails(gtx C, l *load.Load, tx *transactionItem) D {
	date := time.Unix(tx.transaction.Timestamp, 0).Format("Jan 2, 2006")
	timeSplit := time.Unix(tx.transaction.Timestamp, 0).Format("03:04:05 PM")
	dateTime := fmt.Sprintf("%v at %v", date, timeSplit)
	col := l.Theme.Color.GrayText3

	switch tx.status.TicketStatus {
	case dcrlibwallet.TicketStatusUnmined:
		lbl := l.Theme.Label(values.TextSize16, values.StringF(values.StrUnminedInfo, components.TimeAgo(tx.transaction.Timestamp)))
		lbl.Color = col
		return lbl.Layout(gtx)
	case dcrlibwallet.TicketStatusImmature:
		maturity := l.WL.MultiWallet.TicketMaturity()
		blockTime := l.WL.MultiWallet.TargetTimePerBlockMinutes()
		maturityDuration := time.Duration(maturity*int32(blockTime)) * time.Minute

		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				lbl := l.Theme.Label(values.TextSize16, values.StringF(values.StrImmatureInfo, tx.transaction.BlockHeight, maturity,
					maturityDuration.String()))
				lbl.Color = col
				return lbl.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				p := l.Theme.ProgressBarCirle(int(tx.progress))
				p.Color = tx.status.Background
				return p.Layout(gtx)
			}),
		)
	case dcrlibwallet.TicketStatusLive:
		lbl := l.Theme.Label(values.TextSize16, values.String(values.StrLiveInfoDisc))
		lbl.Color = col
		return lbl.Layout(gtx)
	case dcrlibwallet.TicketStatusVotedOrRevoked:
		if tx.ticketSpender.Type == dcrlibwallet.TxTypeVote {
			return multiContent(gtx, l, dateTime, fmt.Sprintf("%s %v", values.String(values.StrVoted), components.TimeAgo(tx.transaction.Timestamp)))
		}
		lbl := l.Theme.Label(values.TextSize16, dateTime)
		lbl.Color = col
		return lbl.Layout(gtx)
	case dcrlibwallet.TicketStatusExpired:
		return multiContent(gtx, l, dateTime, fmt.Sprintf("%s %v", values.String(values.StrExpired), components.TimeAgo(tx.transaction.Timestamp)))
	default:
		return D{}
	}
}

func multiContent(gtx C, l *load.Load, leftText, rightText string) D {
	col := l.Theme.Color.GrayText3
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			lbl := l.Theme.Label(values.TextSize16, leftText)
			lbl.Color = col
			return lbl.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Right: values.MarginPadding5,
				Left:  values.MarginPadding5,
			}.Layout(gtx, func(gtx C) D {
				ic := decredmaterial.NewIcon(l.Theme.Icons.ImageBrightness1)
				ic.Color = col
				return ic.Layout(gtx, values.MarginPadding6)
			})
		}),
		layout.Rigid(func(gtx C) D {
			lbl := l.Theme.Label(values.TextSize16, rightText)
			lbl.Color = col
			return lbl.Layout(gtx)
		}),
	)
}

func ticketCardTooltip(gtx C, rectLayout layout.Dimensions, tooltip *decredmaterial.Tooltip, leftInset unit.Dp, body layout.Widget) {
	inset := layout.Inset{
		Top:  values.MarginPadding15,
		Left: leftInset,
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
	titleLabel := l.Theme.Label(values.TextSize14, title)
	titleLabel.Color = l.Theme.Color.GrayText2

	descLabel := l.Theme.Label(values.TextSize14, desc)

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

// deprecated to be removed in subsequent code clean uo
// ticketCard layouts out Stake info with the shadow box, use for list horizontal or list grid
func ticketCard(gtx layout.Context, l *load.Load, tx *transactionItem, showWalletName bool) layout.Dimensions {
	wal := l.WL.MultiWallet.WalletWithID(tx.transaction.WalletID)
	txStatus := tx.status

	// add this data to transactionItem so it can be shared with list
	maturity := l.WL.MultiWallet.TicketMaturity()

	return decredmaterial.LinearLayout{
		Width:       gtx.Dp(values.MarginPadding168),
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
											return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, l.Theme.Icons.TimerIcon.Layout12dp)
										}),
										layout.Rigid(txt.Layout),
									)

									var durationTitle = values.String(values.StrLiveIn) // immature
									if txStatus.TicketStatus != dcrlibwallet.TicketStatusImmature {
										// voted or revoked
										durationTitle = values.String(values.StrSpendableIn)
									}

									ticketCardTooltip(gtx, durationLayout, tx.durationTooltip, values.MarginPadding0, func(gtx C) D {
										return titleDescTooltip(gtx, l, durationTitle, values.StringF(values.StrDuration, maturityDuration, confirmations, maturity))
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
								txt := l.Theme.Label(values.TextSize14, txStatus.Title)
								txt.Color = txStatus.Color
								txt.Font.Weight = text.Medium
								txtLayout := txt.Layout(gtx)
								ticketCardTooltip(gtx, txtLayout, tx.statusTooltip, values.MarginPadding0, func(gtx C) D {
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
									txt := l.Theme.Label(values.TextSize14, "•")
									txt.Color = l.Theme.Color.GrayText2

									return txt.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								if !showWalletName {
									return D{}
								}

								txt := l.Theme.Label(values.TextSize14, wal.Name)
								txt.Color = l.Theme.Color.GrayText2
								txtLayout := txt.Layout(gtx)
								ticketCardTooltip(gtx, txtLayout, tx.walletNameTooltip, values.MarginPadding0, func(gtx C) D {
									return titleDescTooltip(gtx, l, values.String(values.StrWalletName), txt.Text)
								})
								return txtLayout
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					txt := l.Theme.Label(values.TextSize14, time.Unix(tx.transaction.Timestamp, 0).Format("Jan 2"))
					txt.Color = l.Theme.Color.GrayText3
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txtLayout := txt.Layout(gtx)
							ticketCardTooltip(gtx, txtLayout, tx.dateTooltip, values.MarginPaddingMinus10, func(gtx C) D {
								dateTime := time.Unix(tx.transaction.Timestamp, 0).Format("Jan 2, 2006 at 03:04:05 PM")
								return titleDescTooltip(gtx, l, values.String(values.StrPurchased), dateTime)
							})
							return txtLayout
						}),
						layout.Rigid(func(gtx C) D {

							var tooltipTitle string
							if tx.ticketSpender != nil { // voted or revoked
								if tx.ticketSpender.Type == dcrlibwallet.TxTypeVote {
									tooltipTitle = values.String(values.StrDaysToVote)
								} else {
									tooltipTitle = values.String(values.StrDaysToMiss)
								}
							} else if txStatus.TicketStatus == dcrlibwallet.TicketStatusImmature ||
								txStatus.TicketStatus == dcrlibwallet.TicketStatusLive {
								tooltipTitle = values.String(values.StrStakeAge)
							} else {
								return D{}
							}

							return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Inset{
										Left:  values.MarginPadding4,
										Right: values.MarginPadding4,
									}.Layout(gtx, func(gtx C) D {
										ic := decredmaterial.NewIcon(l.Theme.Icons.ImageBrightness1)
										return ic.Layout(gtx, values.MarginPadding5)
									})
								}),
								layout.Rigid(func(gtx C) D {
									txt.Text = tx.ticketAge
									txtLayout := txt.Layout(gtx)
									ticketCardTooltip(gtx, txtLayout, tx.daysBehindTooltip, values.MarginPaddingMinus10, func(gtx C) D {
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

func ticketListLayout(gtx C, l *load.Load, ticket *transactionItem) layout.Dimensions {
	return layout.Inset{
		Right: values.MarginPadding26,
	}.Layout(gtx, func(gtx C) D {
		return components.EndToEndRow(gtx,
			func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						wrapIcon := l.Theme.Card()
						wrapIcon.Color = ticket.status.Background
						wrapIcon.Radius = decredmaterial.Radius(8)
						dims := wrapIcon.Layout(gtx, func(gtx C) D {
							return layout.UniformInset(values.MarginPadding10).Layout(gtx, ticket.status.Icon.Layout24dp)
						})

						return layout.Inset{
							Right: values.MarginPadding16,
						}.Layout(gtx, func(gtx C) D {
							return dims
						})
					}),
					layout.Rigid(l.Theme.Label(values.TextSize18, ticket.status.Title).Layout),
				)
			},
			func(gtx C) D {
				return ticketStatusDetails(gtx, l, ticket)
			})
	})
}

// todo: cleanup
func createOrderDropDown(th *decredmaterial.Theme) *decredmaterial.DropDown {
	return th.DropDown([]decredmaterial.DropDownItem{{Text: values.String(values.StrNewest)},
		{Text: values.String(values.StrOldest)}}, values.StakingDropdownGroup, 1)
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

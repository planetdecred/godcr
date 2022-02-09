package staking

import (
	"fmt"
	"image"
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

func AllLiveTickets(mw *dcrlibwallet.MultiWallet) ([]dcrlibwallet.Transaction, error) {
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

func WalletLiveTickets(w *dcrlibwallet.Wallet) ([]dcrlibwallet.Transaction, error) {
	var tickets []dcrlibwallet.Transaction
	liveTicketFilters := []int32{dcrlibwallet.TxFilterUnmined, dcrlibwallet.TxFilterImmature, dcrlibwallet.TxFilterLive}
	for _, filter := range liveTicketFilters {
		tx, err := w.GetTransactionsRaw(0, 0, filter, true)
		if err != nil {
			return nil, err
		}

		tickets = append(tickets, tx...)
	}

	return tickets, nil
}

func ticketStatusTooltip(gtx C, l *load.Load, tx *components.TransactionItem) layout.Dimensions {
	status := l.Theme.Label(values.MarginPadding14, strings.ToUpper(tx.Status.Title))
	status.Font.Weight = text.Medium
	status.Color = tx.Status.Color

	maturity := l.WL.MultiWallet.TicketMaturity()
	blockTime := l.WL.MultiWallet.TargetTimePerBlockMinutes()
	maturityDuration := time.Duration(maturity*int32(blockTime)) * time.Minute

	// maturityTime := components.TimeFormat(int(maturityDuration.Seconds()), false)
	// fmt.Println(maturityTime)
	var title, mainDesc, subDesc string
	switch tx.Status.TicketStatus {
	case dcrlibwallet.TicketStatusUnmined:
		title = "This Stake is waiting in mempool to be included in a block."
	case dcrlibwallet.TicketStatusImmature:
		title = fmt.Sprintf("This Stake will enter the Stake pool and become a live Stake after %d blocks (~%s).", maturity, maturityDuration)
	case dcrlibwallet.TicketStatusLive:
		title = "Waiting to be chosen to vote."
		mainDesc = "The average vote time is 28 days, but can take up to 142 days."
		subDesc = "There is a 0.5% chance of expiring before being chosen to vote (this expiration returns the original Stake price without a reward)."
	case dcrlibwallet.TicketStatusVotedOrRevoked:
		if tx.TicketSpender.Type == dcrlibwallet.TxTypeVote {
			title = "Congratulations! This Stake has voted."
			mainDesc = "The Stake price + reward will become spendable after %d blocks (~%s)."
		} else {
			title = "This Stake has been revoked."
			mainDesc = "The Stake price will become spendable after %d blocks (~%s)."
		}

		if tx.TicketSpender.Confirmations(l.WL.MultiWallet.GetBestBlock().Height) > maturity {
			mainDesc = ""
		} else {

			mainDesc = fmt.Sprintf(mainDesc, maturity, maturityDuration)
		}
	case dcrlibwallet.TicketStatusExpired:
		title = fmt.Sprintf("This Stake has not been chosen to vote within %d blocks, and thus expired.", l.WL.MultiWallet.TicketMaturity())
		mainDesc = "Expired tickets will be revoked to return the original Stake price to you."
		subDesc = "If a Stake is not revoked automatically, use the revoke button."
	}

	titleLabel := l.Theme.Label(values.MarginPadding14, title)

	mainDescLabel := l.Theme.Label(values.MarginPadding14, mainDesc)
	mainDescLabel.Color = l.Theme.Color.GrayText2
	subDescLabel := l.Theme.Label(values.MarginPadding14, subDesc)
	subDescLabel.Color = l.Theme.Color.GrayText2
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(tx.Status.Icon.Layout16dp),
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

func ticketCardTooltip(gtx C, rectLayout layout.Dimensions, tooltip *decredmaterial.Tooltip, leftInset unit.Value, body layout.Widget) {
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
	titleLabel := l.Theme.Label(values.MarginPadding14, title)
	titleLabel.Color = l.Theme.Color.GrayText2

	descLabel := l.Theme.Label(values.MarginPadding14, desc)

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

// ticketCard layouts out Stake info with the shadow box, use for list horizontal or list grid
func ticketCard(gtx layout.Context, l *load.Load, tx *components.TransactionItem, showWalletName bool) layout.Dimensions {
	wal := l.WL.MultiWallet.WalletWithID(tx.Transaction.WalletID)
	txStatus := tx.Status

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
						if !tx.ShowTime {
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

									confirmations := tx.Confirmations
									if tx.TicketSpender != nil {
										confirmations = tx.TicketSpender.Confirmations(wal.GetBestBlock())
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

									ticketCardTooltip(gtx, durationLayout, tx.DurationTooltip, values.MarginPadding0, func(gtx C) D {
										return titleDescTooltip(gtx, l, durationTitle, fmt.Sprintf("%s (%d/%d blocks)", maturityDuration, confirmations, maturity))
									})

									return durationLayout
								})
							})
						})
					}),
					layout.Expanded(func(gtx C) D {
						return layout.S.Layout(gtx, func(gtx C) D {

							if !tx.ShowProgress {
								return D{}
							}

							p := l.Theme.ProgressBar(int(tx.Progress))
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
					return components.LayoutBalance(gtx, l, dcrutil.Amount(tx.Transaction.Amount).String())
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
								ticketCardTooltip(gtx, txtLayout, tx.StatusTooltip, values.MarginPadding0, func(gtx C) D {
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
								ticketCardTooltip(gtx, txtLayout, tx.WalletNameTooltip, values.MarginPadding0, func(gtx C) D {
									return titleDescTooltip(gtx, l, "Wallet name", txt.Text)
								})
								return txtLayout
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					txt := l.Theme.Label(values.TextSize14, time.Unix(tx.Transaction.Timestamp, 0).Format("Jan 2"))
					txt.Color = l.Theme.Color.GrayText3
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txtLayout := txt.Layout(gtx)
							ticketCardTooltip(gtx, txtLayout, tx.DateTooltip, values.MarginPaddingMinus10, func(gtx C) D {
								dateTime := time.Unix(tx.Transaction.Timestamp, 0).Format("Jan 2, 2006 at 03:04:05 PM")
								return titleDescTooltip(gtx, l, "Purchased", dateTime)
							})
							return txtLayout
						}),
						layout.Rigid(func(gtx C) D {

							var tooltipTitle string
							if tx.TicketSpender != nil { // voted or revoked
								if tx.TicketSpender.Type == dcrlibwallet.TxTypeVote {
									tooltipTitle = "Days to vote"
								} else {
									tooltipTitle = "Days to miss"
								}
							} else if txStatus.TicketStatus == dcrlibwallet.TicketStatusImmature ||
								txStatus.TicketStatus == dcrlibwallet.TicketStatusLive {
								tooltipTitle = "Stake age"
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
										return ic.Layout(gtx, values.MarginPadding5)
									})
								}),
								layout.Rigid(func(gtx C) D {
									txt.Text = tx.TicketAge
									txtLayout := txt.Layout(gtx)
									ticketCardTooltip(gtx, txtLayout, tx.DaysBehindTooltip, values.MarginPaddingMinus10, func(gtx C) D {
										return titleDescTooltip(gtx, l, tooltipTitle, tx.TicketAge)
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

func ticketListLayout(gtx C, l *load.Load, ticket *components.TransactionItem, i int, showWalletName bool) layout.Dimensions {
	wal := l.WL.MultiWallet.WalletWithID(ticket.Transaction.WalletID)
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Right: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
				return layout.Stack{Alignment: layout.S}.Layout(gtx,
					layout.Stacked(func(gtx C) D {
						wrapIcon := l.Theme.Card()
						wrapIcon.Color = ticket.Status.Background
						wrapIcon.Radius = decredmaterial.Radius(8)
						dims := wrapIcon.Layout(gtx, func(gtx C) D {
							return layout.UniformInset(values.MarginPadding10).Layout(gtx, ticket.Status.Icon.Layout24dp)
						})
						return dims
					}),
					layout.Expanded(func(gtx C) D {
						if !ticket.ShowProgress {
							return D{}
						}
						p := l.Theme.ProgressBar(int(ticket.Progress))
						p.Width = values.MarginPadding44
						p.Height = values.MarginPadding4
						p.Direction = layout.SW
						p.Radius = decredmaterial.BottomRadius(8)
						p.Color = ticket.Status.ProgressBarColor
						p.TrackColor = ticket.Status.ProgressTrackColor
						return p.Layout2(gtx)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if i == 0 {
						return D{}
					}
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					separator := l.Theme.Separator()
					separator.Width = gtx.Constraints.Max.X
					return layout.E.Layout(gtx, separator.Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding6,
						Bottom: values.MarginPadding10,
					}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {

								dtime := l.Theme.Label(values.TextSize14, ticket.PurchaseTime)
								dtime.Color = l.Theme.Color.GrayText3

								return components.EndToEndRow(gtx, func(gtx C) D {
									return components.LayoutBalance(gtx, l, dcrutil.Amount(ticket.Transaction.Amount).String())
								}, func(gtx C) D {
									txtLayout := dtime.Layout(gtx)
									ticketCardTooltip(gtx, txtLayout, ticket.DateTooltip, values.MarginPaddingMinus10, func(gtx C) D {
										return titleDescTooltip(gtx, l, "Purchased", ticket.PurchaseTime)
									})
									return txtLayout
								})
							}),
							layout.Rigid(func(gtx C) D {
								left := func(gtx C) D {
									return layout.Flex{}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											txt := l.Theme.Label(values.TextSize14, ticket.Status.Title)
											txt.Color = ticket.Status.Color
											txt.Font.Weight = text.Medium
											txtLayout := txt.Layout(gtx)
											ticketCardTooltip(gtx, txtLayout, ticket.StatusTooltip, values.MarginPadding0, func(gtx C) D {
												return ticketStatusTooltip(gtx, l, ticket)
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
											ticketCardTooltip(gtx, txtLayout, ticket.WalletNameTooltip, values.MarginPadding0, func(gtx C) D {
												return titleDescTooltip(gtx, l, "Wallet name", txt.Text)
											})
											return txtLayout
										}),
									)
								}

								right := func(gtx C) D {
									return layout.Flex{}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											if ticket.Status.TicketStatus == dcrlibwallet.TicketStatusImmature {
												return layout.Flex{}.Layout(gtx,
													layout.Rigid(func(gtx C) D {
														return layout.Inset{
															Top:   values.MarginPadding4,
															Left:  values.MarginPadding8,
															Right: values.MarginPadding4,
														}.Layout(gtx, func(gtx C) D {
															ic := l.Icons.TimerIcon
															return ic.Layout12dp(gtx)
														})
													}),
													layout.Rigid(func(gtx C) D {
														maturity := l.WL.MultiWallet.TicketMaturity()
														blockTime := l.WL.MultiWallet.TargetTimePerBlockMinutes()
														maturityDuration := time.Duration(maturity*int32(blockTime)) * time.Minute
														return l.Theme.Body2(maturityDuration.String()).Layout(gtx)
													}),
													layout.Rigid(func(gtx C) D {
														return layout.Inset{
															Left:  values.MarginPadding4,
															Right: values.MarginPadding4,
														}.Layout(gtx, func(gtx C) D {
															txt := l.Theme.Label(values.MarginPadding14, "•")
															txt.Color = l.Theme.Color.GrayText2
															return txt.Layout(gtx)
														})
													}),
												)
											}
											return D{}
										}),
										layout.Rigid(func(gtx C) D {
											if ticket.TicketAge == "" {
												return D{}
											}

											var tooltipTitle string
											if ticket.TicketSpender != nil { // voted or revoked
												if ticket.TicketSpender.Type == dcrlibwallet.TxTypeVote {
													tooltipTitle = "Days to vote"
												} else {
													tooltipTitle = "Days to miss"
												}
											} else if ticket.Status.TicketStatus == dcrlibwallet.TicketStatusImmature ||
												ticket.Status.TicketStatus == dcrlibwallet.TicketStatusLive {
												tooltipTitle = "Stake age"
											} else {
												return D{}
											}

											txt := l.Theme.Label(values.TextSize14, ticket.TicketAge)
											txt.Color = l.Theme.Color.GrayText3
											txtLayout := txt.Layout(gtx)
											ticketCardTooltip(gtx, txtLayout, ticket.DaysBehindTooltip, values.MarginPaddingMinus75, func(gtx C) D {
												return titleDescTooltip(gtx, l, tooltipTitle, ticket.TicketAge)
											})
											return txtLayout
										}),
									)
								}
								return components.EndToEndRow(gtx, left, right)
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

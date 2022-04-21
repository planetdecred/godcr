// components contain layout code that are shared by multiple pages but aren't widely used enough to be defined as
// widgets

package components

import (
	"fmt"
	"image/color"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/ararog/timeago"
	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	Uint32Size    = 32 // 32 or 64 ? shifting 32-bit value by 32 bits will always clear it
	MaxInt32      = 1<<(Uint32Size-1) - 1
	WalletsPageID = "Wallets"
)

var MaxWidth = unit.Dp(800)

type (
	C              = layout.Context
	D              = layout.Dimensions
	TransactionRow struct {
		Transaction dcrlibwallet.Transaction
		Index       int
		ShowBadge   bool
	}

	TxStatus struct {
		Title string
		Icon  *decredmaterial.Image

		// tx purchase only
		TicketStatus       string
		Color              color.NRGBA
		ProgressBarColor   color.NRGBA
		ProgressTrackColor color.NRGBA
		Background         color.NRGBA
	}

	// CummulativeWalletsBalance defines total balance for all available wallets.
	CummulativeWalletsBalance struct {
		Total                   dcrutil.Amount
		Spendable               dcrutil.Amount
		ImmatureReward          dcrutil.Amount
		ImmatureStakeGeneration dcrutil.Amount
		LockedByTickets         dcrutil.Amount
		VotingAuthority         dcrutil.Amount
		UnConfirmed             dcrutil.Amount
	}
)

// Container is simply a wrapper for the Inset type. Its purpose is to differentiate the use of an inset as a padding or
// margin, making it easier to visualize the structure of a layout when reading UI code.
type Container struct {
	Padding layout.Inset
}

func (c Container) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return c.Padding.Layout(gtx, w)
}

func UniformPadding(gtx layout.Context, body layout.Widget) layout.Dimensions {
	width := gtx.Constraints.Max.X

	padding := values.MarginPadding24

	if (width - 2*gtx.Dp(padding)) > gtx.Dp(MaxWidth) {
		paddingValue := float32(width-gtx.Dp(MaxWidth)) / 2
		padding = unit.Dp(paddingValue)
	}

	return layout.Inset{
		Top:    values.MarginPadding24,
		Right:  padding,
		Bottom: values.MarginPadding24,
		Left:   padding,
	}.Layout(gtx, body)
}

func UniformHorizontalPadding(gtx layout.Context, body layout.Widget) layout.Dimensions {
	width := gtx.Constraints.Max.X

	padding := values.MarginPadding24

	if (width - 2*gtx.Dp(padding)) > gtx.Dp(MaxWidth) {
		paddingValue := float32(width-gtx.Dp(MaxWidth)) / 2
		padding = unit.Dp(paddingValue)
	}

	return layout.Inset{
		Right: padding,
		Left:  padding,
	}.Layout(gtx, body)
}

func UniformMobile(gtx layout.Context, isHorizontal, withList bool, body layout.Widget) layout.Dimensions {
	insetRight := values.MarginPadding10
	if withList {
		insetRight = values.MarginPadding0
	}

	insetTop := values.MarginPadding24
	if isHorizontal {
		insetTop = values.MarginPadding0
	}

	return layout.Inset{
		Top:   insetTop,
		Right: insetRight,
		Left:  values.MarginPadding10,
	}.Layout(gtx, body)
}

func TransactionTitleIcon(l *load.Load, wal *dcrlibwallet.Wallet, tx *dcrlibwallet.Transaction, ticketSpender *dcrlibwallet.Transaction) *TxStatus {
	var txStatus TxStatus

	setTransactionIcon(l, wal, tx, &txStatus, ticketSpender)
	// setTransactionTitle(l, wal, tx, &txStatus)

	// if tx.Type == dcrlibwallet.TxTypeRegular {
	// 	if tx.Direction == dcrlibwallet.TxDirectionSent {
	// 		txStatus.Title = "Sent"
	// 		txStatus.Icon = l.Theme.Icons.SendIcon
	// 	} else if tx.Direction == dcrlibwallet.TxDirectionReceived {
	// 		txStatus.Title = "Received"
	// 		txStatus.Icon = l.Theme.Icons.ReceiveIcon
	// 	} else if tx.Direction == dcrlibwallet.TxDirectionTransferred {
	// 		txStatus.Title = "Yourself"
	// 		txStatus.Icon = l.Theme.Icons.Transferred
	// 	}
	// } else if tx.Type == dcrlibwallet.TxTypeMixed {
	// 	txStatus.Title = "Mixed"
	// 	txStatus.Icon = l.Theme.Icons.MixedTx
	// } else if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterStaking) {

	// 	if tx.Type == dcrlibwallet.TxTypeTicketPurchase {
	// 		if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterUnmined) {
	// 			txStatus.Title = "Unmined"
	// 			txStatus.Icon = l.Theme.Icons.TicketUnminedIcon
	// 			txStatus.TicketStatus = dcrlibwallet.TicketStatusUnmined
	// 			txStatus.Color = l.Theme.Color.LightBlue6
	// 			txStatus.Background = l.Theme.Color.LightBlue
	// 		} else if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterImmature) {
	// 			txStatus.Title = "Immature"
	// 			txStatus.Icon = l.Theme.Icons.TicketImmatureIcon
	// 			txStatus.Color = l.Theme.Color.LightBlue6
	// 			txStatus.TicketStatus = dcrlibwallet.TicketStatusImmature
	// 			txStatus.ProgressBarColor = l.Theme.Color.LightBlue5
	// 			txStatus.ProgressTrackColor = l.Theme.Color.LightBlue3
	// 			txStatus.Background = l.Theme.Color.LightBlue
	// 		} else if ticketSpender != nil {
	// 			if ticketSpender.Type == dcrlibwallet.TxTypeVote {
	// 				txStatus.Title = "Voted"
	// 				txStatus.Icon = l.Theme.Icons.TicketVotedIcon
	// 				txStatus.Color = l.Theme.Color.Turquoise700
	// 				txStatus.TicketStatus = dcrlibwallet.TicketStatusVotedOrRevoked
	// 				txStatus.ProgressBarColor = l.Theme.Color.Turquoise300
	// 				txStatus.ProgressTrackColor = l.Theme.Color.Turquoise100
	// 				txStatus.Background = l.Theme.Color.Success2
	// 			} else {
	// 				txStatus.Title = "Revoked"
	// 				txStatus.Icon = l.Theme.Icons.TicketRevokedIcon
	// 				txStatus.Color = l.Theme.Color.Orange
	// 				txStatus.TicketStatus = dcrlibwallet.TicketStatusVotedOrRevoked
	// 				txStatus.ProgressBarColor = l.Theme.Color.Danger
	// 				txStatus.ProgressTrackColor = l.Theme.Color.Orange3
	// 				txStatus.Background = l.Theme.Color.Orange2
	// 			}
	// 		} else if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterLive) {
	// 			txStatus.Title = "Live"
	// 			txStatus.Icon = l.Theme.Icons.TicketLiveIcon
	// 			txStatus.Color = l.Theme.Color.Primary
	// 			txStatus.TicketStatus = dcrlibwallet.TicketStatusLive
	// 			txStatus.ProgressBarColor = l.Theme.Color.Primary
	// 			txStatus.ProgressTrackColor = l.Theme.Color.LightBlue4
	// 			txStatus.Background = l.Theme.Color.Primary50
	// 		} else if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterExpired) {
	// 			txStatus.Title = "Expired"
	// 			txStatus.Icon = l.Theme.Icons.TicketExpiredIcon
	// 			txStatus.Color = l.Theme.Color.GrayText2
	// 			txStatus.TicketStatus = dcrlibwallet.TicketStatusExpired
	// 			txStatus.Background = l.Theme.Color.Gray4
	// 		} else {
	// 			txStatus.Title = "Purchased"
	// 			txStatus.Icon = l.Theme.Icons.NewStakeIcon
	// 			txStatus.Color = l.Theme.Color.Text
	// 			txStatus.Background = l.Theme.Color.LightBlue
	// 		}
	// 	} else if tx.Type == dcrlibwallet.TxTypeVote {
	// 		txStatus.Title = "Vote"
	// 		txStatus.Icon = l.Theme.Icons.TicketVotedIcon
	// 	} else if tx.Type == dcrlibwallet.TxTypeRevocation {
	// 		txStatus.Title = "Revocation"
	// 		txStatus.Icon = l.Theme.Icons.TicketRevokedIcon
	// 	}
	// }

	return &txStatus
}

func setTransactionIcon(l *load.Load, wal *dcrlibwallet.Wallet, tx *dcrlibwallet.Transaction, txStatus *TxStatus, ticketSpender *dcrlibwallet.Transaction) {
	switch tx.Direction {
	case dcrlibwallet.TxDirectionSent:
		txStatus.Title = "Sent"
		txStatus.Icon = l.Theme.Icons.SendIcon
	case dcrlibwallet.TxDirectionReceived:
		txStatus.Title = "Received"
		txStatus.Icon = l.Theme.Icons.ReceiveIcon
	default:
		txStatus.Title = "Transferred"
		txStatus.Icon = l.Theme.Icons.Transferred
	}

	// replace icon for staking tx types
	if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterStaking) {
		switch tx.Type {
		case dcrlibwallet.TxTypeTicketPurchase:
			{
				if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterUnmined) {
					txStatus.Title = "Unmined"
					txStatus.Icon = l.Theme.Icons.TicketUnminedIcon
					txStatus.TicketStatus = dcrlibwallet.TicketStatusUnmined
					txStatus.Color = l.Theme.Color.LightBlue6
					txStatus.Background = l.Theme.Color.LightBlue
				} else if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterImmature) {
					txStatus.Title = "Immature"
					txStatus.Icon = l.Theme.Icons.TicketImmatureIcon
					txStatus.Color = l.Theme.Color.LightBlue6
					txStatus.TicketStatus = dcrlibwallet.TicketStatusImmature
					txStatus.ProgressBarColor = l.Theme.Color.LightBlue5
					txStatus.ProgressTrackColor = l.Theme.Color.LightBlue3
					txStatus.Background = l.Theme.Color.LightBlue
				} else if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterLive) {
					txStatus.Title = "Live"
					txStatus.Icon = l.Theme.Icons.TicketLiveIcon
					txStatus.Color = l.Theme.Color.Primary
					txStatus.TicketStatus = dcrlibwallet.TicketStatusLive
					txStatus.ProgressBarColor = l.Theme.Color.Primary
					txStatus.ProgressTrackColor = l.Theme.Color.LightBlue4
					txStatus.Background = l.Theme.Color.Primary50
				} else if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterExpired) {
					txStatus.Title = "Expired"
					txStatus.Icon = l.Theme.Icons.TicketExpiredIcon
					txStatus.Color = l.Theme.Color.GrayText2
					txStatus.TicketStatus = dcrlibwallet.TicketStatusExpired
					txStatus.Background = l.Theme.Color.Gray4
				} else {
					if ticketSpender != nil {
						if ticketSpender.Type == dcrlibwallet.TxTypeVote {
							txStatus.Title = "Ticket, Voted"
							txStatus.Icon = l.Theme.Icons.TicketLiveIcon
							txStatus.Color = l.Theme.Color.Turquoise700
							txStatus.TicketStatus = dcrlibwallet.TicketStatusVotedOrRevoked
							txStatus.ProgressBarColor = l.Theme.Color.Turquoise300
							txStatus.ProgressTrackColor = l.Theme.Color.Turquoise100
							txStatus.Background = l.Theme.Color.Success2
						} else {
							txStatus.Title = "Ticket, Revoked"
							txStatus.Icon = l.Theme.Icons.TicketLiveIcon
							txStatus.Color = l.Theme.Color.Orange
							txStatus.TicketStatus = dcrlibwallet.TicketStatusVotedOrRevoked
							txStatus.ProgressBarColor = l.Theme.Color.Danger
							txStatus.ProgressTrackColor = l.Theme.Color.Orange3
							txStatus.Background = l.Theme.Color.Orange2
						}
					}
					// txStatus.Title = "Purchased"
					// txStatus.Icon = l.Theme.Icons.TicketLiveIcon
					// txStatus.Color = l.Theme.Color.Text
					// txStatus.Background = l.Theme.Color.LightBlue
				}
			}
		case dcrlibwallet.TxTypeVote:
			txStatus.Title = "Vote"
			txStatus.Icon = l.Theme.Icons.TicketVotedIcon
			txStatus.Color = l.Theme.Color.Turquoise700
			txStatus.TicketStatus = dcrlibwallet.TicketStatusVotedOrRevoked
			txStatus.ProgressBarColor = l.Theme.Color.Turquoise300
			txStatus.ProgressTrackColor = l.Theme.Color.Turquoise100
			txStatus.Background = l.Theme.Color.Success2
		default:
			txStatus.Title = "Revocation"
			txStatus.Icon = l.Theme.Icons.TicketRevokedIcon
			txStatus.Color = l.Theme.Color.Orange
			txStatus.TicketStatus = dcrlibwallet.TicketStatusVotedOrRevoked
			txStatus.ProgressBarColor = l.Theme.Color.Danger
			txStatus.ProgressTrackColor = l.Theme.Color.Orange3
			txStatus.Background = l.Theme.Color.Orange2
		}
	} else if tx.Type == dcrlibwallet.TxTypeMixed {
		txStatus.Title = values.String(values.StrMixed)
		txStatus.Icon = l.Theme.Icons.MixedTx
	}
}

func setTransactionTitle(l *load.Load, wal *dcrlibwallet.Wallet, tx *dcrlibwallet.Transaction, txStatus *TxStatus) {
	if tx.Type == dcrlibwallet.TxTypeMixed {
		txStatus.Title = "Mixed"
		txStatus.Icon = l.Theme.Icons.MixedTx
	} else if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterStaking) {
		switch tx.Type {
		case dcrlibwallet.TxTypeTicketPurchase:
			{
				if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterUnmined) {
					txStatus.Title = "Unmined"
				} else if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterImmature) {
					txStatus.Title = "Immature"
				} else {
					txStatus.Title = "Live"
				}
			}
		case dcrlibwallet.TxTypeVote:
			txStatus.Title = "Vote"
		case dcrlibwallet.TxTypeRevocation:
			txStatus.Title = "Revocation"
		}
	}
}

func WeekDayHourMinuteCalculator(timestamp int64) string {
	var dateTimeResult string
	timeStampNow := time.Now().Unix()
	minutesFromTxn := (timeStampNow - timestamp) / 60
	daysFromTxn := minutesFromTxn / 1440 // there are 1440 minutes in 24 hours
	weeksFromTxn := daysFromTxn / 7

	if weeksFromTxn > 0 {
		if weeksFromTxn == 1 {
			dateTimeResult = values.StringF(values.StrWeekAgo, weeksFromTxn)
			return dateTimeResult
		}

		dateTimeResult = values.StringF(values.StrWeeksAgo, weeksFromTxn)
		return dateTimeResult
	}

	if daysFromTxn > 0 {
		if daysFromTxn == 1 {
			dateTimeResult = values.StringF(values.StrDayAgo, daysFromTxn)
			return dateTimeResult
		}

		dateTimeResult = values.StringF(values.StrDaysAgo, daysFromTxn)
		return dateTimeResult
	}

	hoursFromTxn := minutesFromTxn / 60
	if hoursFromTxn > 0 {
		if hoursFromTxn == 1 {
			dateTimeResult = values.StringF(values.StrHourAgo, hoursFromTxn)
			return dateTimeResult
		}

		dateTimeResult = values.StringF(values.StrHoursAgo, hoursFromTxn)
		return dateTimeResult
	}

	if minutesFromTxn > 0 {
		if minutesFromTxn == 1 {
			dateTimeResult = values.StringF(values.StrMinuteAgo, minutesFromTxn)
			return dateTimeResult
		}

		dateTimeResult = values.StringF(values.StrMinutesAgo, minutesFromTxn)
		return dateTimeResult
	}

	dateTimeResult = values.String(values.StrJustNow)

	return dateTimeResult
}

func DurationAgo(timestamp int64) string {
	var duration string

	//Convert timestamp to date in string format (yyyy:mm:dd hr:m:s +0000 UTC)
	currentTimestamp := time.Now().UTC().String()
	txnTimestamp := time.Unix(timestamp, 0).UTC().String()

	//Split the date so we can sepparate into date and time for current time and time of txn
	currentTimeSplit := strings.Split(currentTimestamp, " ")
	txnTimeSplit := strings.Split(txnTimestamp, " ")

	//Split current date and time, and  txn date and time then store in variables
	currentDate := strings.Split(currentTimeSplit[0], "-")
	txnDate := strings.Split(txnTimeSplit[0], "-")
	yearNow, _ := strconv.Atoi(currentDate[0])
	monthNow, _ := strconv.Atoi(currentDate[1])
	txnYear, _ := strconv.Atoi(txnDate[0])
	txnMonth, _ := strconv.Atoi(txnDate[1])
	dayNow, _ := strconv.Atoi(currentDate[2])
	txnDay, _ := strconv.Atoi(txnDate[2])
	currentYearStart := 0
	txnYearEnd := 12

	if (yearNow - txnYear) > 0 {
		if (yearNow - txnYear) == 1 {
			if ((txnYearEnd - txnMonth) + (currentYearStart + monthNow)) < 12 {
				if ((txnYearEnd - txnMonth) + (currentYearStart + monthNow)) == 1 {
					if dayNow < txnDay {
						duration = WeekDayHourMinuteCalculator(timestamp)
						return duration
					}

					duration = values.StringF(values.StrMonthAgo, 1)
					return duration
				}

				monthDiff := ((txnYearEnd - txnMonth) + (currentYearStart + monthNow)) - 1
				if dayNow < txnDay {
					if monthDiff == 1 {
						duration = values.StringF(values.StrMonthAgo, monthDiff)
						return duration
					}

					duration = values.StringF(values.StrMonthsAgo, monthDiff)
					return duration
				}

				duration = values.StringF(values.StrMonthsAgo, (txnYearEnd-txnMonth)+(currentYearStart+monthNow))
				return duration
			}
			y := (yearNow - txnYear) - 1

			if dayNow < txnDay {
				if y == 0 {
					duration := values.StringF(values.StrMonthsAgo, 11)
					return duration
				}

				if y == 1 {
					duration = values.StringF(values.StrYearAgo, y)
					return duration
				}

				duration = values.StringF(values.StrYearsAgo, y)
				return duration
			}

			duration = values.StringF(values.StrYearAgo, yearNow-txnYear)
			return duration
		}

		duration = values.StringF(values.StrYearsAgo, yearNow-txnYear)
		return duration
	}

	if (monthNow - txnMonth) > 0 {
		if (monthNow - txnMonth) == 1 {
			if dayNow < txnDay {
				duration = WeekDayHourMinuteCalculator(timestamp)
				return duration
			}

			duration = values.StringF(values.StrMonthAgo, 1)
			return duration
		}

		if dayNow < txnDay {
			if (monthNow-txnMonth)-1 == 0 {
				duration := values.StringF(values.StrWeeksAgo, 3)
				return duration
			}

			if (monthNow-txnMonth)-1 == 1 {
				duration = values.StringF(values.StrMonthAgo, (monthNow-txnMonth)-1)
				return duration
			}

			duration = values.StringF(values.StrMonthsAgo, (monthNow-txnMonth)-1)
			return duration
		}

		duration = values.StringF(values.StrMonthsAgo, monthNow-txnMonth)
		return duration
	}

	duration = WeekDayHourMinuteCalculator(timestamp)

	return duration
}

// transactionRow is a single transaction row on the transactions and overview page. It lays out a transaction's
// direction, balance, status.
func LayoutTransactionRow(gtx layout.Context, l *load.Load, row TransactionRow) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	wal := l.WL.MultiWallet.WalletWithID(row.Transaction.WalletID)
	var ticketSpender *dcrlibwallet.Transaction
	if wal.TxMatchesFilter(&row.Transaction, dcrlibwallet.TxFilterStaking) {
		ticketSpender, _ = wal.TicketSpender(row.Transaction.Hash)
	}
	txStatus := TransactionTitleIcon(l, wal, &row.Transaction, ticketSpender)

	return decredmaterial.LinearLayout{
		Orientation: layout.Horizontal,
		Width:       decredmaterial.MatchParent,
		Height:      gtx.Dp(values.MarginPadding56),
		Alignment:   layout.Middle,
		Padding:     layout.Inset{Left: values.MarginPadding16, Right: values.MarginPadding16},
	}.Layout(gtx,
		layout.Rigid(txStatus.Icon.Layout24dp),
		layout.Rigid(func(gtx C) D {
			return decredmaterial.LinearLayout{
				Width:       decredmaterial.WrapContent,
				Height:      decredmaterial.MatchParent,
				Orientation: layout.Vertical,
				Padding:     layout.Inset{Left: values.MarginPadding16},
				Direction:   layout.Center,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if row.Transaction.Type == dcrlibwallet.TxTypeRegular {
						amount := dcrutil.Amount(row.Transaction.Amount).String()
						if row.Transaction.Direction == dcrlibwallet.TxDirectionSent {
							amount = "-" + amount
						}
						return LayoutBalance(gtx, l, amount)
					}

					return l.Theme.Label(values.TextSize18, txStatus.Title).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return decredmaterial.LinearLayout{
						Width:       decredmaterial.WrapContent,
						Height:      decredmaterial.WrapContent,
						Orientation: layout.Horizontal,
						Direction:   layout.W,
						Alignment:   layout.Middle,
						Margin:      layout.Inset{Top: values.MarginPadding4},
					}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							if row.ShowBadge {
								return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
									return WalletLabel(gtx, l, wal.Name)
								})
							}

							return layout.Dimensions{}
						}),
						layout.Rigid(func(gtx C) D {
							if wal.TxMatchesFilter(&row.Transaction, dcrlibwallet.TxFilterStaking) {
								ic := l.Theme.Icons.StakeIconInactive
								return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, ic.Layout12dp)
							}
							return D{}
						}),
						layout.Rigid(func(gtx C) D {
							// mix denomination or ticket price
							if row.Transaction.Type == dcrlibwallet.TxTypeMixed {
								mixedDenom := dcrutil.Amount(row.Transaction.MixDenomination).String()
								txt := l.Theme.Label(values.TextSize12, mixedDenom)
								txt.Color = l.Theme.Color.GrayText2
								return txt.Layout(gtx)
							} else if wal.TxMatchesFilter(&row.Transaction, dcrlibwallet.TxFilterStaking) {
								ticketPrice := dcrutil.Amount(row.Transaction.Amount).String()
								txt := l.Theme.Label(values.TextSize12, ticketPrice)
								txt.Color = l.Theme.Color.GrayText2
								return txt.Layout(gtx)
							}
							return layout.Dimensions{}
						}),
						layout.Rigid(func(gtx C) D {
							// Mixed outputs count
							if row.Transaction.Type == dcrlibwallet.TxTypeMixed && row.Transaction.MixCount > 1 {
								label := l.Theme.Label(values.TextSize12, fmt.Sprintf("x%d", row.Transaction.MixCount))
								label.Color = l.Theme.Color.GrayText2
								return layout.Inset{Left: values.MarginPadding4}.Layout(gtx, label.Layout)
							}
							return layout.Dimensions{}
						}),
						layout.Rigid(func(gtx C) D {
							// vote reward
							if row.Transaction.Type != dcrlibwallet.TxTypeVote {
								return D{}
							}

							return decredmaterial.LinearLayout{
								Width:       decredmaterial.WrapContent,
								Height:      decredmaterial.WrapContent,
								Orientation: layout.Horizontal,
								Margin:      layout.Inset{Left: values.MarginPadding4},
								Alignment:   layout.Middle,
							}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									label := l.Theme.Label(values.TextSize14, "+")
									label.Color = l.Theme.Color.Turquoise800
									return label.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									ic := l.Theme.Icons.DecredSymbol2

									return layout.Inset{
										Left:  values.MarginPadding4,
										Right: values.MarginPadding4,
									}.Layout(gtx, ic.Layout16dp)
								}),
								layout.Rigid(func(gtx C) D {
									label := l.Theme.Label(values.TextSize12, dcrutil.Amount(row.Transaction.VoteReward).String())
									label.Color = l.Theme.Color.Turquoise800
									return label.Layout(gtx)
								}),
							)
						}),
					)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			status := l.Theme.Body1(values.String(values.StrPending))
			if TxConfirmations(l, row.Transaction) <= 1 {
				status.Color = l.Theme.Color.GrayText1
			} else {
				status.Color = l.Theme.Color.GrayText2
				status.Text = FormatDateOrTime(row.Transaction.Timestamp)
			}
			return decredmaterial.LinearLayout{
				Width:       decredmaterial.WrapContent,
				Height:      decredmaterial.MatchParent,
				Orientation: layout.Vertical,
				Padding:     layout.Inset{Left: values.MarginPadding16},
				Alignment:   layout.End,
				Direction:   layout.E,
				Margin:      layout.Inset{Top: values.MarginPadding10},
			}.Layout(gtx,
				layout.Rigid(status.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							if row.Transaction.Type == dcrlibwallet.TxTypeVote || row.Transaction.Type == dcrlibwallet.TxTypeRevocation {
								var title string
								if row.Transaction.Type == dcrlibwallet.TxTypeVote {
									title = values.String(values.StrVote)
								} else {
									title = values.String(values.StrRevoke)
								}

								return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, func(gtx C) D {
									return WalletLabel(gtx, l, fmt.Sprintf("%dd to %s", row.Transaction.DaysToVoteOrRevoke, title))
								})
							}

							return D{}
						}),
						layout.Rigid(func(gtx C) D {
							currentTimestamp := time.Now().UTC().String()
							txnTimestamp := time.Unix(row.Transaction.Timestamp, 0).UTC().String()
							currentTimeSplit := strings.Split(currentTimestamp, " ")
							txnTimeSplit := strings.Split(txnTimestamp, " ")
							currentDate := strings.Split(currentTimeSplit[0], "-")
							txnDate := strings.Split(txnTimeSplit[0], "-")

							currentDay, _ := strconv.Atoi(currentDate[2])
							txnDay, _ := strconv.Atoi(txnDate[2])

							if currentDate[0] == txnDate[0] && currentDate[1] == txnDate[1] && currentDay-txnDay < 1 {
								return D{}
							}
							duration := l.Theme.Label(values.TextSize12, DurationAgo(row.Transaction.Timestamp))
							duration.Color = l.Theme.Color.GrayText4
							return layout.Inset{Left: values.MarginPadding2}.Layout(gtx, duration.Layout)
						}),
					)
				}),
			)
		}),
		layout.Rigid(func(gtx C) D {
			statusIcon := l.Theme.Icons.ConfirmIcon
			if TxConfirmations(l, row.Transaction) <= 1 {
				statusIcon = l.Theme.Icons.PendingIcon
			}

			return layout.Inset{Left: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
				return statusIcon.Layout12dp(gtx)
			})
		}),
	)
}

func TxConfirmations(l *load.Load, transaction dcrlibwallet.Transaction) int32 {
	if transaction.BlockHeight != -1 {
		// TODO
		return (l.WL.MultiWallet.WalletWithID(transaction.WalletID).GetBestBlock() - transaction.BlockHeight) + 1
	}

	return 0
}

func FormatDateOrTime(timestamp int64) string {
	utcTime := time.Unix(timestamp, 0).UTC()
	currentTime := time.Now().UTC()

	if strconv.Itoa(currentTime.Year()) == strconv.Itoa(utcTime.Year()) && currentTime.Month().String() == utcTime.Month().String() {
		if strconv.Itoa(currentTime.Day()) == strconv.Itoa(utcTime.Day()) {
			if strconv.Itoa(currentTime.Hour()) == strconv.Itoa(utcTime.Hour()) {
				return TimeAgo(timestamp)
			}

			return TimeAgo(timestamp)
		} else if currentTime.Day()-1 == utcTime.Day() {
			yesterday := values.String(values.StrYesterday)
			return yesterday
		}
	}

	t := strings.Split(utcTime.Format(time.UnixDate), " ")
	t2 := t[2]
	year := strconv.Itoa(utcTime.Year())
	if t[2] == "" {
		t2 = t[3]
	}
	return fmt.Sprintf("%s %s, %s", t[1], t2, year)
}

// walletLabel displays the wallet which a transaction belongs to. It is only displayed on the overview page when there
//// are transactions from multiple wallets
func WalletLabel(gtx layout.Context, l *load.Load, walletName string) D {
	return decredmaterial.Card{
		Color: l.Theme.Color.Gray4,
	}.Layout(gtx, func(gtx C) D {
		return Container{
			layout.Inset{
				Left:  values.MarginPadding4,
				Right: values.MarginPadding4,
			}}.Layout(gtx, func(gtx C) D {
			name := l.Theme.Label(values.TextSize12, walletName)
			name.Color = l.Theme.Color.GrayText2
			return name.Layout(gtx)
		})
	})
}

// EndToEndRow layouts out its content on both ends of its horizontal layout.
func EndToEndRow(gtx layout.Context, leftWidget, rightWidget func(C) D) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(leftWidget),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, rightWidget)
		}),
	)
}

func TimeAgo(timestamp int64) string {
	timeAgo, _ := timeago.TimeAgoWithTime(time.Now(), time.Unix(timestamp, 0))
	return timeAgo
}

func TruncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

func GoToURL(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}

func StringNotEmpty(texts ...string) bool {
	for _, t := range texts {
		if strings.TrimSpace(t) == "" {
			return false
		}
	}

	return true
}

func TimeFormat(secs int, long bool) string {
	var val string
	if secs > 86399 {
		val = "d"
		if long {
			val = " " + values.String(values.StrDays)
		}
		days := secs / 86400
		return fmt.Sprintf("%d%s", days, val)
	} else if secs > 3599 {
		val = "h"
		if long {
			val = " " + values.String(values.StrHours)
		}
		hours := secs / 3600
		return fmt.Sprintf("%d%s", hours, val)
	} else if secs > 59 {
		val = "s"
		if long {
			val = " " + values.String(values.StrMinutes)
		}
		mins := secs / 60
		return fmt.Sprintf("%d%s", mins, val)
	}

	val = "s"
	if long {
		val = " " + values.String(values.StrSeconds)
	}
	return fmt.Sprintf("%d %s", secs, val)
}

// createOrUpdateWalletDropDown check for len of wallets to create dropDown,
// also update the list when create, update, delete a wallet.
func CreateOrUpdateWalletDropDown(l *load.Load, dwn **decredmaterial.DropDown, wallets []*dcrlibwallet.Wallet, grp uint, pos uint) *decredmaterial.DropDown {
	var walletDropDownItems []decredmaterial.DropDownItem
	walletIcon := l.Theme.Icons.WalletIcon
	walletIcon.Scale = 1
	for _, wal := range wallets {
		item := decredmaterial.DropDownItem{
			Text: wal.Name,
			Icon: walletIcon,
		}
		walletDropDownItems = append(walletDropDownItems, item)
	}
	*dwn = l.Theme.DropDown(walletDropDownItems, grp, pos)
	return *dwn
}

func CreateOrderDropDown(l *load.Load, grp uint, pos uint) *decredmaterial.DropDown {
	return l.Theme.DropDown([]decredmaterial.DropDownItem{{Text: values.String(values.StrNewest)},
		{Text: values.String(values.StrOldest)}}, grp, pos)
}

func TranslateErr(err error) string {
	switch err.Error() {
	case dcrlibwallet.ErrInvalidPassphrase:
		return values.String(values.StrInvalidPassphrase)
	}

	return err.Error()
}

// CoinImageBySymbol returns image widget for supported asset coins.
func CoinImageBySymbol(l *load.Load, coinName string) *decredmaterial.Image {
	switch strings.ToLower(coinName) {
	case "btc":
		return l.Theme.Icons.BTC
	case "dcr":
		return l.Theme.Icons.DCR
	}
	return nil
}

func CalculateTotalWalletsBalance(l *load.Load) (*CummulativeWalletsBalance, error) {
	var totalBalance, spandableBalance, immatureReward, votingAuthority,
		immatureStakeGeneration, lockedByTickets, unConfirmed int64

	accountsResult, err := l.WL.SelectedWallet.Wallet.GetAccountsRaw()
	if err != nil {
		return nil, err
	}

	for _, account := range accountsResult.Acc {
		totalBalance += account.TotalBalance
		spandableBalance += account.Balance.Spendable
		immatureReward += account.Balance.ImmatureReward
		immatureStakeGeneration += account.Balance.ImmatureStakeGeneration
		lockedByTickets += account.Balance.LockedByTickets
		votingAuthority += account.Balance.VotingAuthority
		unConfirmed += account.Balance.UnConfirmed
	}

	cumm := &CummulativeWalletsBalance{
		Total:                   dcrutil.Amount(totalBalance),
		Spendable:               dcrutil.Amount(spandableBalance),
		ImmatureReward:          dcrutil.Amount(immatureReward),
		ImmatureStakeGeneration: dcrutil.Amount(immatureStakeGeneration),
		LockedByTickets:         dcrutil.Amount(lockedByTickets),
		VotingAuthority:         dcrutil.Amount(votingAuthority),
		UnConfirmed:             dcrutil.Amount(unConfirmed),
	}

	return cumm, nil
}

// SecondsToDays takes time in seconds and returns its string equivalent in the format ddhhmm.
func SecondsToDays(totalTimeLeft int64) string {
	q, r := divMod(totalTimeLeft, 24*60*60)
	timeLeft := time.Duration(r) * time.Second
	if q > 0 {
		return fmt.Sprintf("%dd%s", q, timeLeft.String())
	}
	return timeLeft.String()
}

// divMod divides a numerator by a denominator and returns its quotient and remainder.
func divMod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

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
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	Uint32Size       = 32 << (^uint32(0) >> 32 & 1) // 32 or 64
	MaxInt32         = 1<<(Uint32Size-1) - 1
	USDExchangeValue = "USD (Bittrex)"
	WalletsPageID    = "Wallets"
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

	if (width - 2*gtx.Px(padding)) > gtx.Px(MaxWidth) {
		paddingValue := float32(width-gtx.Px(MaxWidth)) / 2
		padding = unit.Px(paddingValue)
	}

	return layout.Inset{
		Top:    values.MarginPadding24,
		Right:  padding,
		Bottom: values.MarginPadding24,
		Left:   padding,
	}.Layout(gtx, body)
}

func TransactionTitleIcon(l *load.Load, wal *dcrlibwallet.Wallet, tx *dcrlibwallet.Transaction, ticketSpender *dcrlibwallet.Transaction) *TxStatus {
	var txStatus TxStatus

	if tx.Type == dcrlibwallet.TxTypeRegular {
		if tx.Direction == dcrlibwallet.TxDirectionSent {
			txStatus.Title = "Sent"
			txStatus.Icon = l.Icons.SendIcon
		} else if tx.Direction == dcrlibwallet.TxDirectionReceived {
			txStatus.Title = "Received"
			txStatus.Icon = l.Icons.ReceiveIcon
		} else if tx.Direction == dcrlibwallet.TxDirectionTransferred {
			txStatus.Title = "Yourself"
			txStatus.Icon = l.Icons.Transferred
		}
	} else if tx.Type == dcrlibwallet.TxTypeMixed {
		txStatus.Title = "Mixed"
		txStatus.Icon = l.Icons.MixedTx
	} else if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterStaking) {

		if tx.Type == dcrlibwallet.TxTypeTicketPurchase {
			if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterUnmined) {
				txStatus.Title = "Unmined"
				txStatus.Icon = l.Icons.TicketUnminedIcon
				txStatus.TicketStatus = dcrlibwallet.TicketStatusUnmined
				txStatus.Color = l.Theme.Color.LightBlue6
				txStatus.Background = l.Theme.Color.LightBlue
			} else if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterImmature) {
				txStatus.Title = "Immature"
				txStatus.Icon = l.Icons.TicketImmatureIcon
				txStatus.Color = l.Theme.Color.LightBlue6
				txStatus.TicketStatus = dcrlibwallet.TicketStatusImmature
				txStatus.ProgressBarColor = l.Theme.Color.LightBlue5
				txStatus.ProgressTrackColor = l.Theme.Color.LightBlue3
				txStatus.Background = l.Theme.Color.LightBlue
			} else if ticketSpender != nil {
				if ticketSpender.Type == dcrlibwallet.TxTypeVote {
					txStatus.Title = "Voted"
					txStatus.Icon = l.Icons.TicketVotedIcon
					txStatus.Color = l.Theme.Color.Turquoise700
					txStatus.TicketStatus = dcrlibwallet.TicketStatusVotedOrRevoked
					txStatus.ProgressBarColor = l.Theme.Color.Turquoise300
					txStatus.ProgressTrackColor = l.Theme.Color.Turquoise100
					txStatus.Background = l.Theme.Color.Success2
				} else {
					txStatus.Title = "Revoked"
					txStatus.Icon = l.Icons.TicketRevokedIcon
					txStatus.Color = l.Theme.Color.Orange
					txStatus.TicketStatus = dcrlibwallet.TicketStatusVotedOrRevoked
					txStatus.ProgressBarColor = l.Theme.Color.Danger
					txStatus.ProgressTrackColor = l.Theme.Color.Orange3
					txStatus.Background = l.Theme.Color.Orange2
				}
			} else if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterLive) {
				txStatus.Title = "Live"
				txStatus.Icon = l.Icons.TicketLiveIcon
				txStatus.Color = l.Theme.Color.Primary
				txStatus.TicketStatus = dcrlibwallet.TicketStatusLive
				txStatus.ProgressBarColor = l.Theme.Color.Primary
				txStatus.ProgressTrackColor = l.Theme.Color.LightBlue4
				txStatus.Background = l.Theme.Color.Primary50
			} else if wal.TxMatchesFilter(tx, dcrlibwallet.TxFilterExpired) {
				txStatus.Title = "Expired"
				txStatus.Icon = l.Icons.TicketExpiredIcon
				txStatus.Color = l.Theme.Color.Gray
				txStatus.TicketStatus = dcrlibwallet.TicketStatusExpired
				txStatus.Background = l.Theme.Color.LightGray
			} else {
				txStatus.Title = "Purchased"
				txStatus.Icon = l.Icons.TicketPurchasedIcon
				txStatus.Color = l.Theme.Color.DeepBlue
				txStatus.Background = l.Theme.Color.LightBlue
			}
		} else if tx.Type == dcrlibwallet.TxTypeVote {
			txStatus.Title = "Vote"
			txStatus.Icon = l.Icons.TicketVotedIcon
		} else if tx.Type == dcrlibwallet.TxTypeRevocation {
			txStatus.Title = "Revocation"
			txStatus.Icon = l.Icons.TicketRevokedIcon
		}
	}

	return &txStatus
}

func DurationAgo(timestamp int64) string {
	//Convert timestamp to date in string format (yyyy:mm:dd hr:m:s +0000 UTC)

	currentTimestamp := time.Now().UTC().String()
	txnTimestamp := time.Unix(timestamp, 0).UTC().String()

	//Split the date so we can sepparate into date and time for current time and time of txn

	currentTimeSplit := strings.Split(currentTimestamp, " ")
	txnTimeSplit := strings.Split(txnTimestamp, " ")

	//Split current date and time, and  txn date and time then store in variables

	currentDate := strings.Split(currentTimeSplit[0], "-")
	currentTime := strings.Split(currentTimeSplit[1], ":")
	txnDate := strings.Split(txnTimeSplit[0], "-")
	txnTime := strings.Split(txnTimeSplit[1], ":")
	var duration string

	//Group duration into thresholds

	//Yearly threshold

	if currentDate[0] != txnDate[0] {
		cy, _ := strconv.Atoi(currentDate[0])
		ty, _ := strconv.Atoi(txnDate[0])

		cm, _ := strconv.Atoi(currentDate[1])
		tm, _ := strconv.Atoi(txnDate[1])

		txnYearEnd := 13
		currentYearStart := 1

		result := cy - ty
		numMonths := (txnYearEnd - tm) + (cm - currentYearStart)

		if numMonths > 12 {
			duration = fmt.Sprintf("%s years ago", strconv.Itoa(result))
			return duration
		} else if numMonths == 12 {
			duration = "1 year ago"
			return duration
		}
		duration = fmt.Sprintf("%s months ago", strconv.Itoa(numMonths))
		return duration

	} else if currentDate[0] == txnDate[0] {
		//Monthly threshold

		if currentDate[1] != txnDate[1] {
			cm, _ := strconv.Atoi(currentDate[1])
			tm, _ := strconv.Atoi(txnDate[1])
			result := cm - tm
			if result <= 1 {
				duration = "1 month ago"
				return duration
			}
			duration = fmt.Sprintf("%s months ago", strconv.Itoa(result))
			return duration
		} else if currentDate[1] == txnDate[1] {
			//Weekly threshold

			_, currentWeek := time.Now().UTC().ISOWeek()
			_, txnWeek := time.Unix(timestamp, 0).UTC().ISOWeek()
			if currentWeek != txnWeek {
				result := currentWeek - txnWeek
				if result <= 1 {
					duration = "1 week ago"
					return duration
				}
				duration = fmt.Sprintf("%s weeks ago", strconv.Itoa(result))
				return duration
			} else if currentWeek == txnWeek {
				//Daily threshold

				if currentDate[2] != txnDate[2] {
					cd, _ := strconv.Atoi(currentDate[2])
					td, _ := strconv.Atoi(txnDate[2])
					result := cd - td
					if result <= 1 {
						duration = "1 day ago"
						return duration
					}
					duration = fmt.Sprintf("%s days ago", strconv.Itoa(result))
					return duration
				} else if currentDate[2] == txnDate[2] {
					//Hourly threshold

					if currentTime[0] != txnTime[0] {
						ch, _ := strconv.Atoi(currentTime[0])
						th, _ := strconv.Atoi(txnTime[0])
						result := ch - th
						if result <= 1 {
							duration = "Last hour"
							return duration
						}
						duration = fmt.Sprintf("%s hours ago", strconv.Itoa(result))
						return duration
					} else if currentTime[0] == txnTime[0] {
						//Minute threshold

						if currentTime[1] != txnTime[1] {
							cm, _ := strconv.Atoi(currentTime[1])
							tm, _ := strconv.Atoi(txnTime[1])
							result := cm - tm
							if result <= 1 {
								duration = "1 minute ago"
								return duration
							}
							duration = fmt.Sprintf("%s minutes ago", strconv.Itoa(result))
							return duration
						} else if currentTime[1] == txnTime[1] {
							duration = "Just now"
						}
					}
				}
			}
		}
	}

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
		Height:      gtx.Px(values.MarginPadding56),
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

					label := l.Theme.Label(values.TextSize18, txStatus.Title)
					label.Color = l.Theme.Color.DeepBlue
					return label.Layout(gtx)
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
								ic := l.Icons.TicketIconInactive
								return layout.Inset{Right: values.MarginPadding4}.Layout(gtx, ic.Layout12dp)
							}
							return D{}
						}),
						layout.Rigid(func(gtx C) D {
							// mix denomination or ticket price
							if row.Transaction.Type == dcrlibwallet.TxTypeMixed {
								mixedDenom := dcrutil.Amount(row.Transaction.MixDenomination).String()
								txt := l.Theme.Label(values.TextSize14, mixedDenom)
								txt.Color = l.Theme.Color.Gray
								return txt.Layout(gtx)
							} else if wal.TxMatchesFilter(&row.Transaction, dcrlibwallet.TxFilterStaking) {
								ticketPrice := dcrutil.Amount(row.Transaction.Amount).String()
								txt := l.Theme.Label(values.TextSize14, ticketPrice)
								txt.Color = l.Theme.Color.Gray
								return txt.Layout(gtx)
							}
							return layout.Dimensions{}
						}),
						layout.Rigid(func(gtx C) D {
							// Mixed outputs count
							if row.Transaction.Type == dcrlibwallet.TxTypeMixed && row.Transaction.MixCount > 1 {
								label := l.Theme.Label(values.TextSize14, fmt.Sprintf("x%d", row.Transaction.MixCount))
								label.Color = l.Theme.Color.Gray
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
									ic := l.Icons.DecredSymbol2

									return layout.Inset{
										Left:  values.MarginPadding4,
										Right: values.MarginPadding4,
									}.Layout(gtx, ic.Layout16dp)
								}),
								layout.Rigid(func(gtx C) D {
									label := l.Theme.Label(values.TextSize14, dcrutil.Amount(row.Transaction.VoteReward).String())
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
			status := l.Theme.Body1("pending")
			if TxConfirmations(l, row.Transaction) <= 1 {
				status.Color = l.Theme.Color.Gray5
			} else {
				status.Color = l.Theme.Color.Gray
				status.Text = FormatDateOrTime(row.Transaction.Timestamp)
			}
			return layout.E.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.End}.Layout(gtx,
					layout.Rigid(status.Layout),
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Alignment: layout.End}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								if row.Transaction.Type == dcrlibwallet.TxTypeVote || row.Transaction.Type == dcrlibwallet.TxTypeRevocation {
									var title string
									if row.Transaction.Type == dcrlibwallet.TxTypeVote {
										title = "vote"
									} else {
										title = "revoke"
									}

									return WalletLabel(gtx, l, fmt.Sprintf("%dd to %s", row.Transaction.DaysToVoteOrRevoke, title))
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
								duration := l.Theme.Label(values.TextSize14, DurationAgo(row.Transaction.Timestamp))
								duration.Color = l.Theme.Color.InactiveGray
								return layout.Inset{Left: values.MarginPadding2}.Layout(gtx, duration.Layout)
							}),
						)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			statusIcon := l.Icons.ConfirmIcon
			if TxConfirmations(l, row.Transaction) <= 1 {
				statusIcon = l.Icons.PendingIcon
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
			yesterday := "Yesterday"
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
		Color: l.Theme.Color.LightGray,
	}.Layout(gtx, func(gtx C) D {
		return Container{
			layout.Inset{
				Left:  values.MarginPadding4,
				Right: values.MarginPadding4,
			}}.Layout(gtx, func(gtx C) D {
			name := l.Theme.Label(values.TextSize12, walletName)
			name.Color = l.Theme.Color.Gray
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
			val = " day(s)"
		}
		days := secs / 86400
		return fmt.Sprintf("%d%s", days, val)
	} else if secs > 3599 {
		val = "h"
		if long {
			val = " hour(s)"
		}
		hours := secs / 3600
		return fmt.Sprintf("%d%s", hours, val)
	} else if secs > 59 {
		val = "s"
		if long {
			val = " min(s)"
		}
		mins := secs / 60
		return fmt.Sprintf("%d%s", mins, val)
	}

	val = "s"
	if long {
		val = " sec(s)"
	}
	return fmt.Sprintf("%d %s", secs, val)
}

/*
func (page *pageCommon) Modal(gtx layout.Context, body layout.Dimensions, modal layout.Dimensions) layout.Dimensions {
	dims := layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return body
		}),
		layout.Stacked(func(gtx C) D {
			return modal
		}),
	)
	return dims
}

func (page *pageCommon) initSelectAccountWidget(wallAcct map[int][]walletAccount, windex int) {
	if _, ok := wallAcct[windex]; !ok {
		accounts := page.info.Wallets[windex].Accounts
		if len(accounts) != len(wallAcct[windex]) {
			wallAcct[windex] = make([]walletAccount, len(accounts))
			for aindex := range accounts {
				if accounts[aindex].Name == "imported" {
					continue
				}

				wallAcct[windex][aindex] = walletAccount{
					walletIndex:  windex,
					accountIndex: aindex,
					evt:          &gesture.Click{},
					accountName:  accounts[aindex].Name,
					totalBalance: accounts[aindex].TotalBalance,
					spendable:    dcrutil.Amount(accounts[aindex].SpendableBalance).String(),
					number:       accounts[aindex].Number,
				}
			}
		}
	}
}*/

/*func ticketStatusTooltip(gtx C, c *pageCommon, rect image.Rectangle, t *wallet.Ticket, tooltip *decredmaterial.Tooltip) layout.Dimensions {
	inset := layout.Inset{
		Top:   values.MarginPadding15,
		Right: unit.Dp(-150),
		Left:  values.MarginPadding15,
	}
	return tooltip.Layout(gtx, rect, inset, func(gtx C) D {
		st := ticketStatusIcon(c, t.Info.Status)
		var title, message, message2 string
		switch t.Info.Status {
		case "UNMINED":
			title = "This ticket is waiting in mempool to be included in a block."
			message, message2 = "", ""
		case "IMMATURE":
			title = "This ticket will enter the ticket pool and become a live ticket after 256 blocks (~20 hrs)."
			message, message2 = "", ""
		case "LIVE":
			title = "Waiting to be chosen to vote."
			message = "The average vote time is 28 days, but can take up to 142 days."
			message2 = "There is a 0.5% chance of expiring before being chosen to vote (this expiration returns the original ticket price without a reward)."
		case "VOTED":
			title = "Congratulations! This ticket has voted."
			message = "The ticket price + reward will become spendable after 256 blocks (~20 hrs)."
			message2 = ""
		case "MISSED":
			title = "This ticket was chosen to vote, but missed the voting window."
			message = "Missed tickets will be revoked to return the original ticket price to you."
			message2 = "If a ticket is not revoked automatically, use the revoke button."
		case "EXPIRED":
			title = "This ticket has not been chosen to vote within 40960 blocks, and thus expired. "
			message = "Expired tickets will be revoked to return the original ticket price to you."
			message2 = "If a ticket is not revoked automatically, use the revoke button."
		case "REVOKED":
			title = "This ticket has been revoked."
			message = "The ticket price will become spendable after 256 blocks (~20 hrs)."
			message2 = ""
		}
		titleLabel, messageLabel, messageLabel2 := c.theme.Body2(title), c.theme.Body2(message), c.theme.Body2(message2)
		messageLabel.Color, messageLabel2.Color = c.theme.Color.Gray, c.theme.Color.Gray

		status := c.theme.Body2(t.Info.Status)
		status.Color = st.color
		st.icon.Scale = .5
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(st.icon.Layout),
					layout.Rigid(toolTipContent(layout.Inset{Left: values.MarginPadding4}, status.Layout)),
				)
			}),
			layout.Rigid(toolTipContent(layout.Inset{Top: values.MarginPadding8}, titleLabel.Layout)),
			layout.Rigid(toolTipContent(layout.Inset{Top: values.MarginPadding8}, messageLabel.Layout)),
			layout.Rigid(func(gtx C) D {
				if message2 != "" {
					toolTipContent(layout.Inset{Top: values.MarginPadding8}, messageLabel2.Layout)
				}
				return layout.Dimensions{}
			}),
		)
	})
}

func toolTipContent(inset layout.Inset, body layout.Widget) layout.Widget {
	return func(gtx C) D {
		return inset.Layout(gtx, body)
	}
}*/

// ticketCard layouts out ticket info with the shadow box, use for list horizontal or list grid
/*func ticketCard(gtx layout.Context, c *pageCommon, t *wallet.Ticket, tooltip *decredmaterial.Tooltip) layout.Dimensions {
	var itemWidth int
	st := ticketStatusIcon(c, t.Info.Status)
	if st == nil {
		return layout.Dimensions{}
	}
	st.icon.Scale = 1.0
	return c.theme.Shadow().Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				wrap := c.theme.Card()
				wrap.Radius = decredmaterial.CornerRadius{NE: 8, NW: 8, SE: 0, SW: 0}
				wrap.Color = st.background
				return wrap.Layout(gtx, func(gtx C) D {
					return layout.Stack{Alignment: layout.S}.Layout(gtx,

						layout.Expanded(func(gtx C) D {
							return layout.NE.Layout(gtx, func(gtx C) D {
								wTimeLabel := c.theme.Card()
								wTimeLabel.Radius = decredmaterial.CornerRadius{NE: 0, NW: 8, SE: 0, SW: 8}
								return wTimeLabel.Layout(gtx, func(gtx C) D {
									return layout.Inset{
										Top:    values.MarginPadding4,
										Bottom: values.MarginPadding4,
										Right:  values.MarginPadding8,
										Left:   values.MarginPadding8,
									}.Layout(gtx, func(gtx C) D {
										return c.theme.Label(values.TextSize14, "10h 47m").Layout(gtx)
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
									p := c.theme.ProgressBar(20)
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
				wrap := c.theme.Card()
				wrap.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 8, SW: 8}
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
									return c.layoutBalance(gtx, t.Amount, true)
								})
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										txt := c.theme.Label(values.MarginPadding14, t.Info.Status)
										txt.Color = st.color
										txtLayout := txt.Layout(gtx)
										rect := image.Rectangle{
											Max: image.Point{
												X: txtLayout.Size.X,
												Y: txtLayout.Size.Y,
											},
										}
										ticketStatusTooltip(gtx, c, rect, t, tooltip)
										return txtLayout
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{
											Left:  values.MarginPadding4,
											Right: values.MarginPadding4,
										}.Layout(gtx, func(gtx C) D {
											ic := c.icons.imageBrightness1
											ic.Color = c.theme.Color.Gray2
											return c.icons.imageBrightness1.Layout(gtx, values.MarginPadding5)
										})
									}),
									layout.Rigid(func(gtx C) D {
										return c.theme.Label(values.MarginPadding14, t.WalletName).Layout(gtx)
									}),
								)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Top:    values.MarginPadding16,
									Bottom: values.MarginPadding16,
								}.Layout(gtx, func(gtx C) D {
									txt := c.theme.Label(values.TextSize14, t.MonthDay)
									txt.Color = c.theme.Color.Gray2
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
												ic.Color = c.theme.Color.Gray2
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
	})
}*/

// ticketActivityRow layouts out ticket info, display ticket activities on the tickets_page and tickets_activity_page
/*func ticketActivityRow(gtx layout.Context, c *pageCommon, t wallet.Ticket, index int) layout.Dimensions {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Right: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
				st := ticketStatusIcon(c, t.Info.Status)
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
					separator := c.theme.Separator()
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
								labelStatus := c.theme.Label(values.TextSize18, strings.Title(strings.ToLower(t.Info.Status)))
								labelStatus.Color = c.theme.Color.DeepBlue

								labelDaysBehind := c.theme.Label(values.TextSize14, t.DaysBehind)
								labelDaysBehind.Color = c.theme.Color.DeepBlue

								return endToEndRow(gtx,
									labelStatus.Layout,
									labelDaysBehind.Layout)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Flex{
									Alignment: layout.Middle,
								}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										txt := c.theme.Label(values.TextSize14, t.WalletName)
										txt.Color = c.theme.Color.Gray2
										return txt.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{
											Left:  values.MarginPadding4,
											Right: values.MarginPadding4,
										}.Layout(gtx, func(gtx C) D {
											ic := c.icons.imageBrightness1
											ic.Color = c.theme.Color.Gray2
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
										txt := c.theme.Label(values.TextSize14, t.Amount)
										txt.Color = c.theme.Color.Gray2
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
}*/

/*func (page *pageCommon) handleToast() {
	if (*page.toast) == nil {
		return
	}

	if (*page.toast).timer == nil {
		(*page.toast).timer = time.NewTimer(time.Second * 3)
	}

	select {
	case <-(*page.toast).timer.C:
		*page.toast = nil
	default:
	}
}*/

// createOrUpdateWalletDropDown check for len of wallets to create dropDown,
// also update the list when create, update, delete a wallet.
func CreateOrUpdateWalletDropDown(l *load.Load, dwn **decredmaterial.DropDown, wallets []*dcrlibwallet.Wallet) {
	var walletDropDownItems []decredmaterial.DropDownItem
	walletIcon := l.Icons.WalletIcon
	walletIcon.Scale = 1
	for _, wal := range wallets {
		item := decredmaterial.DropDownItem{
			Text: wal.Name,
			Icon: walletIcon,
		}
		walletDropDownItems = append(walletDropDownItems, item)
	}
	*dwn = l.Theme.DropDown(walletDropDownItems, 1)
}

func CreateOrderDropDown(l *load.Load) *decredmaterial.DropDown {
	return l.Theme.DropDown([]decredmaterial.DropDownItem{{Text: values.String(values.StrNewest)},
		{Text: values.String(values.StrOldest)}}, 1)
}

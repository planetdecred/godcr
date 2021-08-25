// components contain layout code that are shared by multiple pages but aren't widely used enough to be defined as
// widgets

package components

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/ararog/timeago"
	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/values"
	"golang.org/x/text/message"
)

const (
	Uint32Size       = 32 << (^uint32(0) >> 32 & 1) // 32 or 64
	MaxInt32         = 1<<(Uint32Size-1) - 1
	USDExchangeValue = "usd_bittrex"
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

// breakBalance takes the balance string and returns it in two slices
func BreakBalance(p *message.Printer, balance string) (b1, b2 string) {
	var isDecimal = true
	balanceParts := strings.Split(balance, ".")
	if len(balanceParts) == 1 {
		isDecimal = false
		balanceParts = strings.Split(balance, " ")
	}

	b1 = balanceParts[0]
	if bal, err := strconv.Atoi(b1); err == nil {
		b1 = p.Sprint(bal)
	}

	b2 = balanceParts[1]
	if isDecimal {
		b1 = b1 + "." + b2[:2]
		b2 = b2[2:]
		return
	}
	b2 = " " + b2
	return
}

// LayoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func LayoutBalance(gtx layout.Context, l *load.Load, amount string) layout.Dimensions {
	// todo: make "DCR" symbols small when there are no decimals in the balance
	mainText, subText := BreakBalance(l.Printer, amount)
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			label := l.Theme.Label(values.TextSize20, mainText)
			label.Color = l.Theme.Color.DeepBlue
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			label := l.Theme.Label(values.TextSize14, subText)
			label.Color = l.Theme.Color.DeepBlue
			return label.Layout(gtx)
		}),
	)
}

// transactionRow is a single transaction row on the transactions and overview page. It lays out a transaction's
// direction, balance, status.
func LayoutTransactionRow(gtx layout.Context, l *load.Load, row TransactionRow) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X

	wal := l.WL.MultiWallet.WalletWithID(row.Transaction.WalletID)

	var title string
	var icon *widget.Image

	if row.Transaction.Type == dcrlibwallet.TxTypeRegular {
		if row.Transaction.Direction == dcrlibwallet.TxDirectionSent {
			icon = l.Icons.SendIcon
		} else if row.Transaction.Direction == dcrlibwallet.TxDirectionReceived {
			icon = l.Icons.ReceiveIcon
		} else if row.Transaction.Direction == dcrlibwallet.TxDirectionTransferred {
			icon = l.Icons.Transferred
		}
	} else if row.Transaction.Type == dcrlibwallet.TxTypeMixed {
		title = "Mixed"
		icon = l.Icons.MixedTx
	} else if wal.TxMatchesFilter(&row.Transaction, dcrlibwallet.TxFilterStaking) {

		if row.Transaction.Type == dcrlibwallet.TxTypeTicketPurchase {
			if wal.TxMatchesFilter(&row.Transaction, dcrlibwallet.TxFilterImmature) {
				title = "Immatute"
				icon = l.Icons.TicketImmatureIcon
			} else if wal.TxMatchesFilter(&row.Transaction, dcrlibwallet.TxFilterLive) {
				title = "Live"
				icon = l.Icons.TicketLiveIcon
			} else {
				title = "Purchased"
				icon = l.Icons.TicketPurchasedIcon
			}
		} else if row.Transaction.Type == dcrlibwallet.TxTypeVote {
			title = "Vote"
			icon = l.Icons.TicketVotedIcon
		} else if row.Transaction.Type == dcrlibwallet.TxTypeRevocation {
			title = "Revocation"
			icon = l.Icons.TicketRevokedIcon
		}
	}

	return decredmaterial.LinearLayout{
		Orientation: layout.Horizontal,
		Width:       decredmaterial.MatchParent,
		Height:      gtx.Px(values.MarginPadding56),
		Direction:   layout.W,
		Padding:     layout.Inset{Left: values.MarginPadding16, Right: values.MarginPadding16},
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			width := float32(icon.Src.Size().X)
			scale := 24.0 / width
			icon.Scale = scale

			return layout.W.Layout(gtx, icon.Layout)
		}),
		layout.Rigid(func(gtx C) D {
			return decredmaterial.LinearLayout{
				Width:       decredmaterial.WrapContent,
				Height:      decredmaterial.MatchParent,
				Orientation: layout.Vertical,
				Padding:     layout.Inset{Left: values.MarginPadding16},
				Direction:   layout.W,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if row.Transaction.Type == dcrlibwallet.TxTypeRegular {
						amount := dcrutil.Amount(row.Transaction.Amount).String()
						if row.Transaction.Direction == dcrlibwallet.TxDirectionSent {
							amount = "-" + amount
						}
						return LayoutBalance(gtx, l, amount)
					} else {
						label := l.Theme.Label(values.TextSize18, title)
						label.Color = l.Theme.Color.DeepBlue
						return label.Layout(gtx)
					}
				}),
				layout.Rigid(func(gtx C) D {
					if row.ShowBadge {
						return WalletLabel(gtx, l, wal.Name)
					}
					return layout.Dimensions{}
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Spacing:   layout.SpaceStart,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Right: values.MarginPadding8}.Layout(gtx,
						func(gtx C) D {
							gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
							status := l.Theme.Body1("pending")
							if TxConfirmations(l, row.Transaction) <= 1 {
								status.Color = l.Theme.Color.Gray5
							} else {
								status.Color = l.Theme.Color.Gray4
								status.Text = FormatDateOrTime(row.Transaction.Timestamp)
							}
							return layout.E.Layout(gtx, status.Layout)
						})
				}),
				layout.Rigid(func(gtx C) D {
					gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
					statusIcon := l.Icons.ConfirmIcon
					if TxConfirmations(l, row.Transaction) <= 1 {
						statusIcon = l.Icons.PendingIcon
					}
					statusIcon.Scale = 1.0
					return layout.E.Layout(gtx, statusIcon.Layout)
				}))
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
	if time.Now().UTC().Sub(utcTime).Hours() < 168 {
		return utcTime.Weekday().String()
	}

	t := strings.Split(utcTime.Format(time.UnixDate), " ")
	t2 := t[2]
	if t[2] == "" {
		t2 = t[3]
	}
	return fmt.Sprintf("%s %s", t[1], t2)
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
	for _, wal := range wallets {
		item := decredmaterial.DropDownItem{
			Text: wal.Name,
			Icon: l.Icons.WalletIcon,
		}
		walletDropDownItems = append(walletDropDownItems, item)
	}
	*dwn = l.Theme.DropDown(walletDropDownItems, 2)
}

func CreateOrderDropDown(l *load.Load) *decredmaterial.DropDown {
	return l.Theme.DropDown([]decredmaterial.DropDownItem{{Text: values.String(values.StrNewest)},
		{Text: values.String(values.StrOldest)}}, 1)
}

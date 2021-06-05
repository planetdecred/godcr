// components contain layout code that are shared by multiple pages but aren't widely used enough to be defined as
// widgets

package ui

import (
	"fmt"
	"image"
	"strings"
	"time"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const (
	purchasingAccountTitle = "Purchasing account"
	sendingAccountTitle    = "Sending account"
	receivingAccountTitle  = "Receiving account"
)

type (
	TransactionRow struct {
		transaction wallet.Transaction
		index       int
		showBadge   bool
	}
	toast struct {
		text    string
		success bool
		timer   *time.Timer
	}
)

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func (page *pageCommon) layoutBalance(gtx layout.Context, amount string, isSwitchColor bool) layout.Dimensions {
	// todo: make "DCR" symbols small when there are no decimals in the balance
	mainText, subText := breakBalance(page.printer, amount)
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			label := page.theme.Label(values.TextSize20, mainText)
			if isSwitchColor {
				label.Color = page.theme.Color.DeepBlue
			}
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			label := page.theme.Label(values.TextSize14, subText)
			if isSwitchColor {
				label.Color = page.theme.Color.DeepBlue
			}
			return label.Layout(gtx)
		}),
	)
}

var (
	navDrawerWidth          = unit.Value{U: unit.UnitDp, V: 160}
	navDrawerMinimizedWidth = unit.Value{U: unit.UnitDp, V: 72}
)

// transactionRow is a single transaction row on the transactions and overview page. It lays out a transaction's
// direction, balance, status.
func transactionRow(gtx layout.Context, common *pageCommon, row TransactionRow) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	directionIconTopMargin := values.MarginPadding16

	if row.index == 0 && row.showBadge {
		directionIconTopMargin = values.MarginPadding14
	} else if row.index == 0 {
		// todo: remove top margin from container
		directionIconTopMargin = values.MarginPadding0
	}

	return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				icon := common.icons.receiveIcon
				if row.transaction.Txn.Direction == dcrlibwallet.TxDirectionSent {
					icon = common.icons.sendIcon
				}
				icon.Scale = 1.0

				return layout.Inset{Top: directionIconTopMargin}.Layout(gtx, func(gtx C) D {
					return icon.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if row.index == 0 {
							return layout.Dimensions{}
						}
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						separator := common.theme.Separator()
						separator.Width = gtx.Constraints.Max.X - gtx.Px(unit.Dp(16))
						return layout.E.Layout(gtx, func(gtx C) D {
							// Todo: add comment
							marginBottom := values.MarginPadding16
							if row.showBadge {
								marginBottom = values.MarginPadding5
							}
							return layout.Inset{Bottom: marginBottom}.Layout(gtx,
								func(gtx C) D {
									return separator.Layout(gtx)
								})
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Inset{}.Layout(gtx, func(gtx C) D {
							return layout.Flex{
								Axis:      layout.Horizontal,
								Spacing:   layout.SpaceBetween,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Left: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
										return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												return common.layoutBalance(gtx, row.transaction.Balance, true)
											}),
											layout.Rigid(func(gtx C) D {
												if row.showBadge {
													return walletLabel(gtx, common, row.transaction.WalletName)
												}
												return layout.Dimensions{}
											}),
										)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Right: values.MarginPadding8}.Layout(gtx,
												func(gtx C) D {
													s := formatDateOrTime(row.transaction.Txn.Timestamp)
													if row.transaction.Status != "confirmed" {
														s = row.transaction.Status
													}
													status := common.theme.Body1(s)
													if row.transaction.Status != "confirmed" {
														status.Color = common.theme.Color.Gray5
													} else {
														status.Color = common.theme.Color.Gray4
													}
													status.Alignment = text.Middle
													return status.Layout(gtx)
												})
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Right: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
												statusIcon := common.icons.confirmIcon
												if row.transaction.Status != "confirmed" {
													statusIcon = common.icons.pendingIcon
												}
												statusIcon.Scale = 1.0
												return statusIcon.Layout(gtx)
											})
										}),
									)
								}),
							)
						})
					}),
				)
			}),
		)
	})
}

// walletLabel displays the wallet which a transaction belongs to. It is only displayed on the overview page when there
// are transactions from multiple wallets
func walletLabel(gtx layout.Context, c *pageCommon, walletName string) D {
	return decredmaterial.Card{
		Color: c.theme.Color.LightGray,
	}.Layout(gtx, func(gtx C) D {
		return Container{
			layout.Inset{
				Left:  values.MarginPadding4,
				Right: values.MarginPadding4,
			}}.Layout(gtx, func(gtx C) D {
			name := c.theme.Label(values.TextSize12, walletName)
			name.Color = c.theme.Color.Gray
			return name.Layout(gtx)
		})
	})
}

// endToEndRow layouts out its content on both ends of its horizontal layout.
func endToEndRow(gtx layout.Context, leftWidget, rightWidget func(C) D) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(leftWidget),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, rightWidget)
		}),
	)
}

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

func (page *pageCommon) accountSelectorLayout(gtx layout.Context, callingPage, sendOption string) layout.Dimensions {
	border := widget.Border{
		Color:        page.theme.Color.Gray1,
		CornerRadius: values.MarginPadding8,
		Width:        values.MarginPadding2,
	}
	page.wallAcctSelector.sendOption = sendOption

	d := func(gtx layout.Context, acctName, walName, bal string, btn *widget.Clickable) layout.Dimensions {
		return border.Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
				return decredmaterial.Clickable(gtx, btn, func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							accountIcon := page.icons.accountIcon
							accountIcon.Scale = 1
							inset := layout.Inset{
								Right: values.MarginPadding8,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return accountIcon.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return page.theme.Body1(acctName).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Left: values.MarginPadding4,
								Top:  values.MarginPadding2,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return decredmaterial.Card{
									Color: page.theme.Color.LightGray,
								}.Layout(gtx, func(gtx C) D {
									m2 := values.MarginPadding2
									m4 := values.MarginPadding4
									inset := layout.Inset{
										Left:   m4,
										Top:    m2,
										Bottom: m2,
										Right:  m4,
									}
									return inset.Layout(gtx, func(gtx C) D {
										text := page.theme.Caption(walName)
										text.Color = page.theme.Color.Gray
										return text.Layout(gtx)
									})
								})
							})
						}),
						layout.Flexed(1, func(gtx C) D {
							return layout.E.Layout(gtx, func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										txt := page.theme.Body1(bal)
										txt.Color = page.theme.Color.DeepBlue
										return txt.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										inset := layout.Inset{
											Left: values.MarginPadding15,
										}
										return inset.Layout(gtx, func(gtx C) D {
											return page.icons.dropDownIcon.Layout(gtx, values.MarginPadding20)
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

	switch {
	case callingPage == "send":
		wallSelect := page.info.Wallets[page.wallAcctSelector.selectedSendWallet]
		acctSelect := wallSelect.Accounts[page.wallAcctSelector.selectedSendAccount]
		return d(gtx, acctSelect.Name, wallSelect.Name, dcrutil.Amount(acctSelect.SpendableBalance).String(), page.wallAcctSelector.sendAccountBtn)
	case callingPage == "receive":
		wallSelect := page.info.Wallets[page.wallAcctSelector.selectedReceiveWallet]
		acctSelect := wallSelect.Accounts[page.wallAcctSelector.selectedReceiveAccount]
		return d(gtx, acctSelect.Name, wallSelect.Name, dcrutil.Amount(acctSelect.SpendableBalance).String(), page.wallAcctSelector.receivingAccountBtn)
	case callingPage == "purchase":
		wallSelect := page.info.Wallets[page.wallAcctSelector.selectedPurchaseTicketWallet]
		acctSelect := wallSelect.Accounts[page.wallAcctSelector.selectedPurchaseTicketAccount]
		return d(gtx, acctSelect.Name, wallSelect.Name, dcrutil.Amount(acctSelect.SpendableBalance).String(), page.wallAcctSelector.purchaseTicketAccountBtn)
	default:
		return layout.Dimensions{}
	}
}

func (page *pageCommon) walletAccountModalLayout(gtx layout.Context) layout.Dimensions {
	wallAcctGroup := func(gtx layout.Context, title string, windex int, body layout.Widget) layout.Dimensions {
		return layout.Inset{
			Bottom: values.MarginPadding10,
		}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := page.theme.Body2(title)
							txt.Color = page.theme.Color.Text
							inset := layout.Inset{
								Bottom: values.MarginPadding15,
							}
							return inset.Layout(gtx, txt.Layout)
						}),
						layout.Rigid(func(gtx C) D {
							var showInfoBtn bool
							if windex == 0 && page.wallAcctSelector.title == receivingAccountTitle {
								showInfoBtn = true
							} else if windex == 0 && page.wallAcctSelector.sendOption == "Address" {
								showInfoBtn = true
							}

							if showInfoBtn {
								inset := layout.Inset{
									Top: values.MarginPadding2,
								}
								return inset.Layout(gtx, func(gtx C) D {
									return page.wallAcctSelector.walletInfoButton.Layout(gtx)
								})
							}
							return layout.Dimensions{}
						}),
					)
				}),
				layout.Rigid(body),
			)
		})
	}
	wallAcctSelector := page.wallAcctSelector
	w := []layout.Widget{
		func(gtx C) D {
			title := page.theme.H6(wallAcctSelector.title)
			title.Color = page.theme.Color.Text
			title.Font.Weight = text.Bold
			return title.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Stack{Alignment: layout.NW}.Layout(gtx,
				layout.Expanded(func(gtx C) D {
					return wallAcctSelector.walletsList.Layout(gtx, len(page.info.Wallets), func(gtx C, windex int) D {
						if page.info.Wallets[windex].IsWatchingOnly {
							return D{}
						}
						walletID := page.info.Wallets[windex].ID
						mixedAcct := page.wallet.ReadMixerConfigValueForKey(dcrlibwallet.AccountMixerMixedAccount, walletID)
						unmixedAcct := page.wallet.ReadMixerConfigValueForKey(dcrlibwallet.AccountMixerUnmixedAccount, walletID)

						return wallAcctGroup(gtx, page.info.Wallets[windex].Name, windex, func(gtx C) D {
							return wallAcctSelector.accountsList.Layout(gtx, len(page.info.Wallets[windex].Accounts), func(gtx C, aindex int) D {
								var visibleAccount walletAccount
								fromAccount := wallAcctSelector.walletAccounts.selectSendAccount[windex][aindex]
								toAccount := wallAcctSelector.walletAccounts.selectReceiveAccount[windex][aindex]
								purchaseTicketAccount := wallAcctSelector.walletAccounts.selectPurchaseTicketAccount[windex][aindex]

								switch {
								case fromAccount.number != unmixedAcct &&
									page.wallAcctSelector.title == sendingAccountTitle && page.wallAcctSelector.sendOption == "Address":
									if fromAccount.spendable != "" || fromAccount.evt != nil {
										click := fromAccount.evt
										pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
										click.Add(gtx.Ops)
										page.walletAccountsHandler(gtx, fromAccount)
										visibleAccount = fromAccount
									}
								case toAccount.number != mixedAcct && page.wallAcctSelector.title == receivingAccountTitle:
									if toAccount.spendable != "" || toAccount.evt != nil {
										click := toAccount.evt
										pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
										click.Add(gtx.Ops)
										page.walletAccountsHandler(gtx, toAccount)
										visibleAccount = toAccount
									}
								case page.wallAcctSelector.sendOption == "My account" && page.wallAcctSelector.title == sendingAccountTitle:
									if toAccount.spendable != "" || toAccount.evt != nil {
										click := toAccount.evt
										pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
										click.Add(gtx.Ops)
										page.walletAccountsHandler(gtx, toAccount)
										visibleAccount = toAccount
									}
								case page.wallAcctSelector.title == purchasingAccountTitle:
									if purchaseTicketAccount.spendable != "" || purchaseTicketAccount.evt != nil {
										click := purchaseTicketAccount.evt
										pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
										click.Add(gtx.Ops)
										page.walletAccountsHandler(gtx, purchaseTicketAccount)
										visibleAccount = purchaseTicketAccount
									}
								default:
								}

								if visibleAccount.spendable != "" || visibleAccount.evt != nil {
									return page.walletAccountLayout(gtx, visibleAccount)
								}
								return D{}
							})
						})
					})
				}),
				layout.Stacked(func(gtx C) D {
					if page.wallAcctSelector.isWalletAccountInfo {
						inset := layout.Inset{
							Top:  values.MarginPadding20,
							Left: values.MarginPaddingMinus75,
						}
						return inset.Layout(gtx, func(gtx C) D {
							return page.walletInfoPopup(gtx)
						})
					}
					return layout.Dimensions{}
				}),
			)
		},
	}

	return wallAcctSelector.walletAccount.Layout(gtx, w, 850)
}

func (page *pageCommon) walletAccountsHandler(gtx layout.Context, wallAcct walletAccount) {
	for _, e := range wallAcct.evt.Events(gtx) {
		if e.Type == gesture.TypeClick {
			if page.wallAcctSelector.title == sendingAccountTitle {
				page.wallAcctSelector.selectedSendWallet = wallAcct.walletIndex
				page.wallAcctSelector.selectedSendAccount = wallAcct.accountIndex
			}

			if page.wallAcctSelector.title == receivingAccountTitle {
				page.wallAcctSelector.selectedReceiveWallet = wallAcct.walletIndex
				page.wallAcctSelector.selectedReceiveAccount = wallAcct.accountIndex
			}

			if page.wallAcctSelector.title == purchasingAccountTitle {
				page.wallAcctSelector.selectedPurchaseTicketWallet = wallAcct.walletIndex
				page.wallAcctSelector.selectedPurchaseTicketAccount = wallAcct.accountIndex
			}

			page.wallAcctSelector.isWalletAccountModalOpen = false
		}
	}
}

func (page *pageCommon) walletAccountLayout(gtx layout.Context, wallAcct walletAccount) layout.Dimensions {
	accountIcon := page.icons.accountIcon
	accountIcon.Scale = 1
	return layout.Inset{
		Bottom: values.MarginPadding20,
	}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(0.1, func(gtx C) D {
						return layout.Inset{
							Right: values.MarginPadding18,
						}.Layout(gtx, func(gtx C) D {
							return accountIcon.Layout(gtx)
						})
					}),
					layout.Flexed(0.8, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								acct := page.theme.Label(values.TextSize18, wallAcct.accountName)
								acct.Color = page.theme.Color.Text
								return endToEndRow(gtx, acct.Layout, func(gtx C) D {
									return page.layoutBalance(gtx, wallAcct.totalBalance, true)
								})
							}),
							layout.Rigid(func(gtx C) D {
								spendable := page.theme.Label(values.TextSize14, "Spendable")
								spendable.Color = page.theme.Color.Gray
								spendableBal := page.theme.Label(values.TextSize14, wallAcct.spendable)
								spendableBal.Color = page.theme.Color.Gray
								return endToEndRow(gtx, spendable.Layout, spendableBal.Layout)
							}),
						)
					}),
					layout.Flexed(0.1, func(gtx C) D {
						inset := layout.Inset{
							Right: values.MarginPadding10,
							Top:   values.MarginPadding10,
						}
						sections := func(gtx layout.Context) layout.Dimensions {
							return layout.E.Layout(gtx, func(gtx C) D {
								return inset.Layout(gtx, func(gtx C) D {
									return page.icons.navigationCheck.Layout(gtx, values.MarginPadding20)
								})
							})
						}

						if page.wallAcctSelector.title == sendingAccountTitle &&
							page.wallAcctSelector.selectedSendWallet == wallAcct.walletIndex &&
							page.wallAcctSelector.selectedSendAccount == wallAcct.accountIndex {
							return sections(gtx)
						}

						if page.wallAcctSelector.title == receivingAccountTitle &&
							page.wallAcctSelector.selectedReceiveWallet == wallAcct.walletIndex &&
							page.wallAcctSelector.selectedReceiveAccount == wallAcct.accountIndex {
							return sections(gtx)
						}

						if page.wallAcctSelector.title == purchasingAccountTitle &&
							page.wallAcctSelector.selectedPurchaseTicketWallet == wallAcct.walletIndex &&
							page.wallAcctSelector.selectedPurchaseTicketAccount == wallAcct.accountIndex {
							return sections(gtx)
						}

						return layout.Dimensions{}
					}),
				)
			}),
		)
	})
}

func (page *pageCommon) walletInfoPopup(gtx layout.Context) layout.Dimensions {
	acctType := "Unmixed"
	t := "from"
	txDirection := " Spending"
	if page.wallAcctSelector.title == receivingAccountTitle {
		acctType = "The mixed"
		t = "to"
		txDirection = " Receiving"
	}
	title := fmt.Sprintf("%s accounts are hidden.", acctType)
	desc := fmt.Sprintf("%s %s accounts is disabled by StakeShuffle settings to protect your privacy.", t, strings.ToLower(acctType))
	card := page.theme.Card()
	card.Radius = decredmaterial.CornerRadius{NE: 7, NW: 7, SE: 7, SW: 7}
	border := widget.Border{Color: page.theme.Color.Background, CornerRadius: values.MarginPadding7, Width: values.MarginPadding1}
	gtx.Constraints.Max.X = gtx.Px(values.MarginPadding280)
	return border.Layout(gtx, func(gtx C) D {
		return card.Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								txt := page.theme.Body2(title)
								txt.Color = page.theme.Color.DeepBlue
								txt.Font.Weight = text.Bold
								return txt.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								txt := page.theme.Body2(txDirection)
								txt.Color = page.theme.Color.Gray
								return txt.Layout(gtx)
							}),
						)
					}),
					layout.Rigid(func(gtx C) D {
						txt := page.theme.Body2(desc)
						txt.Color = page.theme.Color.Gray
						return txt.Layout(gtx)
					}),
				)
			})
		})
	})
}

func (page *pageCommon) addAccount(account walletAccount) {
	account.evt = &gesture.Click{}
	selectSendAccount := page.wallAcctSelector.walletAccounts.selectSendAccount
	selectReceiveAccount := page.wallAcctSelector.walletAccounts.selectReceiveAccount
	selectPurchaseTicketAccount := page.wallAcctSelector.walletAccounts.selectPurchaseTicketAccount

	if len(selectSendAccount) > account.walletIndex {
		account.accountIndex = len(selectSendAccount[account.walletIndex])
		page.wallAcctSelector.walletAccounts.selectSendAccount[account.walletIndex] = append(page.wallAcctSelector.walletAccounts.selectSendAccount[account.walletIndex], account)
	}

	if len(selectReceiveAccount) > account.walletIndex {
		account.accountIndex = len(selectReceiveAccount[account.walletIndex])
		page.wallAcctSelector.walletAccounts.selectReceiveAccount[account.walletIndex] = append(page.wallAcctSelector.walletAccounts.selectReceiveAccount[account.walletIndex], account)
	}

	if len(selectPurchaseTicketAccount) > account.walletIndex {
		account.accountIndex = len(selectReceiveAccount[account.walletIndex])
		page.wallAcctSelector.walletAccounts.selectPurchaseTicketAccount[account.walletIndex] = append(page.wallAcctSelector.walletAccounts.selectPurchaseTicketAccount[account.walletIndex], account)
	}
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
}

// ticketCard layouts out ticket info with the shadow box, use for list horizontal or list grid
func ticketCard(gtx layout.Context, c *pageCommon, t *wallet.Ticket) layout.Dimensions {
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
}

// ticketActivityRow layouts out ticket info, display ticket activities on the tickets_page and tickets_activity_page
func ticketActivityRow(gtx layout.Context, c *pageCommon, t wallet.Ticket, index int) layout.Dimensions {
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
}

func displayToast(th *decredmaterial.Theme, gtx layout.Context, n *toast) layout.Dimensions {
	color := th.Color.Success
	if !n.success {
		color = th.Color.Danger
	}

	card := th.Card()
	card.Color = color
	return card.Layout(gtx, func(gtx C) D {
		return layout.Inset{
			Top: values.MarginPadding7, Bottom: values.MarginPadding7,
			Left: values.MarginPadding15, Right: values.MarginPadding15,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			t := th.Body1(n.text)
			t.Color = th.Color.Surface
			return t.Layout(gtx)
		})
	})
}

func (page *pageCommon) handleToast() {
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
}

// createOrUpdateWalletDropDown check for len of wallets to create dropDown,
// also update the list when create, update, delete a wallet.
func (page *pageCommon) createOrUpdateWalletDropDown(dwn **decredmaterial.DropDown) {
	init := func() {
		var walletDropDownItems []decredmaterial.DropDownItem
		for i := range page.info.Wallets {
			item := decredmaterial.DropDownItem{
				Text: page.info.Wallets[i].Name,
				Icon: page.icons.walletIcon,
			}
			walletDropDownItems = append(walletDropDownItems, item)
		}
		*dwn = page.theme.DropDown(walletDropDownItems, 2)
	}

	if *dwn == nil && len(page.info.Wallets) > 0 {
		init()
		return
	}
	if (*dwn).Len() != len(page.info.Wallets) {
		init()
	}
}

func createOrderDropDown(c *pageCommon) *decredmaterial.DropDown {
	return c.theme.DropDown([]decredmaterial.DropDownItem{{Text: values.String(values.StrNewest)},
		{Text: values.String(values.StrOldest)}}, 1)
}

func (page *pageCommon) handler() {
	page.handleToast()

	for windex := 0; windex < page.info.LoadedWallets; windex++ {
		page.initSelectAccountWidget(page.wallAcctSelector.walletAccounts.selectSendAccount, windex)
		page.initSelectAccountWidget(page.wallAcctSelector.walletAccounts.selectReceiveAccount, windex)
		page.initSelectAccountWidget(page.wallAcctSelector.walletAccounts.selectPurchaseTicketAccount, windex)
	}

	if page.wallAcctSelector.sendAccountBtn.Clicked() {
		page.wallAcctSelector.title = sendingAccountTitle
		page.wallAcctSelector.isWalletAccountModalOpen = true
	}

	if page.wallAcctSelector.receivingAccountBtn.Clicked() {
		page.wallAcctSelector.title = receivingAccountTitle
		page.wallAcctSelector.isWalletAccountModalOpen = true
	}

	if page.wallAcctSelector.purchaseTicketAccountBtn.Clicked() {
		page.wallAcctSelector.title = purchasingAccountTitle
		page.wallAcctSelector.isWalletAccountModalOpen = true
	}

	if page.wallAcctSelector.walletInfoButton.Button.Clicked() {
		page.wallAcctSelector.isWalletAccountInfo = !page.wallAcctSelector.isWalletAccountInfo
	}
}

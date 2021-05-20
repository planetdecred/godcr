// components contain layout code that are shared by multiple pages but aren't widely used enough to be defined as
// widgets

package ui

import (
	"fmt"
	"image"
	"strconv"
	"strings"

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

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func (page pageCommon) layoutBalance(gtx layout.Context, amount string) layout.Dimensions {
	// todo: make "DCR" symbols small when there are no decimals in the balance
	mainText, subText := breakBalance(page.printer, amount)
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return page.theme.Label(values.TextSize20, mainText).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return page.theme.Label(values.TextSize14, subText).Layout(gtx)
		}),
	)
}

func (page pageCommon) layoutUSDBalance(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			currencyExchangeValue := page.wallet.ReadStringConfigValueForKey(dcrlibwallet.CurrencyConversionConfigKey)
			page.usdExchangeSet = false
			if currencyExchangeValue == USDExchangeValue {
				page.usdExchangeSet = true
			}
			if page.usdExchangeSet && page.dcrUsdtBittrex.LastTradeRate != "" {
				page.usdExchangeRate, _ = strconv.ParseFloat(page.dcrUsdtBittrex.LastTradeRate, 64)
				TotalBalanceFloat, _ := strconv.ParseFloat(page.info.TotalBalanceRaw, 64)
				page.amountDCRtoUSD = TotalBalanceFloat * page.usdExchangeRate

				inset := layout.Inset{
					Top:  values.MarginPadding3,
					Left: values.MarginPadding8,
				}
				border := widget.Border{Color: page.theme.Color.LightGray, CornerRadius: unit.Dp(8), Width: unit.Dp(0.5)}
				return inset.Layout(gtx, func(gtx C) D {
					padding := layout.Inset{
						Top:    values.MarginPadding3,
						Bottom: values.MarginPadding3,
						Left:   values.MarginPadding6,
						Right:  values.MarginPadding6,
					}
					return border.Layout(gtx, func(gtx C) D {
						return padding.Layout(gtx, func(gtx C) D {
							amountDCRtoUSDString := formatUSDBalance(page.printer, page.amountDCRtoUSD)
							return page.theme.Label(values.TextSize14, amountDCRtoUSDString).Layout(gtx)
						})
					})
				})
			}
			return D{}
		}),
	)
}

// layoutTopBar is the top horizontal bar on every page of the app. It lays out the wallet balance, receive and send
// buttons.
func (page pageCommon) layoutTopBar(gtx layout.Context) layout.Dimensions {
	card := page.theme.Card()
	card.Radius = decredmaterial.CornerRadius{}
	return card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.W.Layout(gtx, func(gtx C) D {
							h := values.MarginPadding24
							v := values.MarginPadding16
							// Balance container
							return Container{padding: layout.Inset{Right: h, Left: h, Top: v, Bottom: v}}.Layout(gtx,
								func(gtx C) D {
									return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											img := page.icons.logo
											img.Scale = 1.0
											return layout.Inset{Right: values.MarginPadding16}.Layout(gtx,
												func(gtx C) D {
													return img.Layout(gtx)
												})
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Center.Layout(gtx, func(gtx C) D {
												return page.layoutBalance(gtx, page.info.TotalBalance)
											})
										}),
										layout.Rigid(func(gtx C) D {
											return page.layoutUSDBalance(gtx)
										}),
									)
								})
						})
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
								list := layout.List{Axis: layout.Horizontal}
								return list.Layout(gtx, len(page.appBarNavItems), func(gtx C, i int) D {
									// header buttons container
									return Container{layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
										return decredmaterial.Clickable(gtx, page.appBarNavItems[i].clickable, func(gtx C) D {
											return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
												layout.Rigid(func(gtx C) D {
													return layout.Inset{Right: values.MarginPadding8}.Layout(gtx,
														func(gtx C) D {
															return layout.Center.Layout(gtx, func(gtx C) D {
																img := page.appBarNavItems[i].image
																img.Scale = 1.0
																return page.appBarNavItems[i].image.Layout(gtx)
															})
														})
												}),
												layout.Rigid(func(gtx C) D {
													return layout.Inset{
														Left: values.MarginPadding0,
													}.Layout(gtx, func(gtx C) D {
														return layout.Center.Layout(gtx, func(gtx C) D {
															return page.theme.Body1(page.appBarNavItems[i].page).Layout(gtx)
														})
													})
												}),
											)
										})
									})
								})
							})
						})
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return page.theme.Separator().Layout(gtx)
			}),
		)
	})
}

var (
	navDrawerWidth          = unit.Value{U: unit.UnitDp, V: 160}
	navDrawerMinimizedWidth = unit.Value{U: unit.UnitDp, V: 72}
)

// layoutNavDrawer is the left vertical pane on every page of the app. It vertically lays out buttons used to navigate
// to different pages.
func (page pageCommon) layoutNavDrawer(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			list := layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(page.drawerNavItems), func(gtx C, i int) D {
				background := page.theme.Color.Surface
				if page.drawerNavItems[i].page == *page.page {
					background = page.theme.Color.LightGray
				}
				txt := page.theme.Label(values.TextSize16, page.drawerNavItems[i].page)
				return decredmaterial.Clickable(gtx, page.drawerNavItems[i].clickable, func(gtx C) D {
					card := page.theme.Card()
					card.Color = background
					card.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
					return card.Layout(gtx, func(gtx C) D {
						return Container{
							layout.Inset{
								Top:    values.MarginPadding16,
								Right:  values.MarginPadding24,
								Bottom: values.MarginPadding16,
								Left:   values.MarginPadding24,
							},
						}.Layout(gtx, func(gtx C) D {
							axis := layout.Horizontal
							leftInset := values.MarginPadding15
							width := navDrawerWidth
							if *page.isNavDrawerMinimized {
								axis = layout.Vertical
								txt.TextSize = values.TextSize10
								leftInset = values.MarginPadding0
								width = navDrawerMinimizedWidth
							}

							gtx.Constraints.Min.X = gtx.Px(width)
							return layout.Flex{Axis: axis}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									img := page.drawerNavItems[i].imageInactive
									if page.drawerNavItems[i].page == *page.page {
										img = page.drawerNavItems[i].image
									}
									return layout.Center.Layout(gtx, func(gtx C) D {
										img.Scale = 1.0
										return img.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{
										Left: leftInset,
										Top:  values.MarginPadding4,
									}.Layout(gtx, func(gtx C) D {
										textColor := page.theme.Color.Gray3
										if page.drawerNavItems[i].page == *page.page {
											textColor = page.theme.Color.DeepBlue
										}
										txt.Color = textColor
										return layout.Center.Layout(gtx, txt.Layout)
									})
								}),
							)
						})
					})
				})
			})
		}),
		layout.Expanded(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.SE.Layout(gtx, func(gtx C) D {
				btn := page.minimizeNavDrawerButton
				if *page.isNavDrawerMinimized {
					btn = page.maximizeNavDrawerButton
				}
				return btn.Layout(gtx)
			})
		}),
	)
}

type TransactionRow struct {
	transaction wallet.Transaction
	index       int
	showBadge   bool
}

// transactionRow is a single transaction row on the transactions and overview page. It lays out a transaction's
// direction, balance, status.
func transactionRow(gtx layout.Context, common pageCommon, row TransactionRow) layout.Dimensions {
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
												return common.layoutBalance(gtx, row.transaction.Balance)
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
													status.Color = common.theme.Color.Gray
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
func walletLabel(gtx layout.Context, c pageCommon, walletName string) D {
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
		layout.Rigid(func(gtx C) D {
			return leftWidget(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return rightWidget(gtx)
			})
		}),
	)
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
									return page.layoutBalance(gtx, wallAcct.totalBalance)
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

func (page pageCommon) addAccount(account walletAccount) {
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

func (page pageCommon) initSelectAccountWidget(wallAcct map[int][]walletAccount, windex int) {
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

func (page pageCommon) handleNavEvents() {
	for page.minimizeNavDrawerButton.Button.Clicked() {
		*page.isNavDrawerMinimized = true
	}

	for page.maximizeNavDrawerButton.Button.Clicked() {
		*page.isNavDrawerMinimized = false
	}

	for i := range page.appBarNavItems {
		for page.appBarNavItems[i].clickable.Clicked() {
			page.setReturnPage(*page.page)
			page.changePage(page.appBarNavItems[i].page)
		}
	}

	for i := range page.drawerNavItems {
		for page.drawerNavItems[i].clickable.Clicked() {
			page.changePage(page.drawerNavItems[i].page)
		}
	}

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

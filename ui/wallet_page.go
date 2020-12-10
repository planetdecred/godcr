package ui

import (
	"strings"

	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"image/color"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageWallet = "Wallet"

type walletPage struct {
	walletInfo    *wallet.MultiWalletInfo
	subPage       int
	current       wallet.InfoShort
	wallet        *wallet.Wallet
	walletAccount **wallet.Account
	theme         *decredmaterial.Theme

	walletIcon                                 *widget.Image
	accountIcon                                *widget.Image
	addAcct, backupButton                      decredmaterial.IconButton
	container, accountsList, walletsList, list layout.List
	line                                       *decredmaterial.Line
	txFeeCollapsible                           *decredmaterial.Collapsible
	toAddWalletPage                            *widget.Clickable

	walletCollapsible []*decredmaterial.Collapsible
	toAcctDetails     []*gesture.Click
}

func (win *Window) WalletPage(common pageCommon) layout.Widget {
	pg := &walletPage{
		walletInfo: win.walletInfo,
		container: layout.List{
			Axis: layout.Vertical,
		},
		accountsList: layout.List{
			Axis: layout.Vertical,
		},

		walletsList: layout.List{
			Axis: layout.Vertical,
		},
		list: layout.List{
			Axis: layout.Vertical,
		},
		theme:            common.theme,
		wallet:           common.wallet,
		txFeeCollapsible: common.theme.Collapsible(new(widget.Clickable)),
		line:             common.theme.Line(),
		walletAccount:    &win.walletAccount,
		toAddWalletPage:  new(widget.Clickable),
	}
	pg.line.Height = 1

	// init wallet collapse
	for i := 0; i < 20; i++ {
		pg.walletCollapsible = append(pg.walletCollapsible, win.theme.Collapsible(new(widget.Clickable)))
	}

	pg.addAcct = common.theme.IconButton(new(widget.Clickable), common.icons.contentAdd)
	pg.addAcct.Inset = layout.UniformInset(values.MarginPadding0)
	pg.addAcct.Size = values.MarginPadding25
	pg.addAcct.Background = color.RGBA{}
	pg.addAcct.Color = common.theme.Color.Text

	pg.backupButton = common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowForward)
	pg.backupButton.Color = common.theme.Color.Surface
	pg.backupButton.Inset = layout.UniformInset(values.MarginPadding0)
	pg.backupButton.Size = values.MarginPadding20

	return func(gtx C) D {
		pg.Handle(common)
		return pg.Layout(gtx, common)
	}
}

// Layout lays out the widgets for the main wallets pg.
func (pg *walletPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	if common.info.LoadedWallets == 0 {
		return common.Layout(gtx, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return common.theme.H3("pg.text.noWallet").Layout(gtx)
			})
		})
	}

	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.walletSection(gtx, common)
		},
		func(gtx C) D {
			return pg.watchOnlyWalletSection(gtx, common)
		},
	}

	body := func(gtx C) D {
		return layout.Stack{Alignment: layout.SE}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return pg.container.Layout(gtx, len(pageContent), func(gtx C, i int) D {
					return layout.UniformInset(values.MarginPadding5).Layout(gtx, pageContent[i])
				})
			}),
			layout.Stacked(func(gtx C) D {
				icon := common.icons.newWalletIcon
				icon.Scale = 0.26
				gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
				return layout.SE.Layout(gtx, func(gtx C) D {
					return decredmaterial.Clickable(gtx, pg.toAddWalletPage, func(gtx C) D {
						return icon.Layout(gtx)
					})
				})
			}),
		)
	}

	return common.Layout(gtx, body)
}

func (pg *walletPage) walletSection(gtx layout.Context, common pageCommon) layout.Dimensions {
	pg.walletIcon = &widget.Image{Src: paint.NewImageOp(common.icons.walletIcon)}
	pg.walletIcon.Scale = 0.05

	return pg.walletsList.Layout(gtx, len(common.info.Wallets), func(gtx C, i int) D {
		wn := common.info.Wallets[i].Name
		wb := common.info.Wallets[i].Balance
		accounts := common.info.Wallets[i].Accounts
		seed := common.info.Wallets[i].Seed

		collapsibleHeader := func(gtx C) D {
			walName := strings.Title(strings.ToLower(wn))
			walletNameLabel := pg.theme.Body1(walName)
			walletBalLabel := pg.theme.Body1(wb)
			walletBalLabel.Color = pg.theme.Color.Gray
			return pg.tableLayout(gtx, walletNameLabel, walletBalLabel, true, len(seed))
		}

		collapsibleFooter := func(gtx C) D {
			return pg.backupSeedNotification(gtx, common)
		}

		collapsibleBody := func(gtx C) D {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				gtx.Constraints.Min.Y = 100

				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.accountsList.Layout(gtx, len(accounts), func(gtx C, i int) D {
							accountsName := accounts[i].Name
							totalBalance := accounts[i].TotalBalance
							spendable := dcrutil.Amount(accounts[i].SpendableBalance).String()

							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return pg.walletAccountsLayout(gtx, accountsName, totalBalance, spendable, common)
								}),
							)
						})
					}),
					// add to line()
					layout.Rigid(func(gtx C) D {
						pg.line.Width = gtx.Constraints.Max.X
						pg.line.Color = common.theme.Color.Background
						m := values.MarginPadding10
						return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
							return pg.line.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Bottom: values.MarginPadding5,
									Right:  values.MarginPadding10,
								}.Layout(gtx, func(gtx C) D {
									return pg.addAcct.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								txt := pg.theme.H6("Add new account")
								txt.Color = common.theme.Color.Gray
								return txt.Layout(gtx)
							}),
						)
					}),
				)
			})
		}

		return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
			pg.walletCollapsible[i].AddItems([]decredmaterial.MoreItem{
				decredmaterial.NewMoreOptionItem("Sign message"),
				decredmaterial.NewMoreOptionItem("Verify message"),
				decredmaterial.NewMoreOptionItem("View property"),
				decredmaterial.NewMoreOptionItem("Rename"),
				decredmaterial.NewMoreOptionItem("Settings"),
			})
			return pg.walletCollapsible[i].Layout(gtx, collapsibleHeader, collapsibleBody, collapsibleFooter)
		})
	})
}

func (pg *walletPage) watchOnlyWalletSection(gtx layout.Context, common pageCommon) layout.Dimensions {
	return pg.sectionLayout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				txt := pg.theme.Body1("Watch-only Wallets")
				txt.Color = pg.theme.Color.Gray
				return txt.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				pg.line.Width = gtx.Constraints.Max.X
				pg.line.Color = common.theme.Color.Hint
				m := values.MarginPadding10
				inset := layout.Inset{
					Top:    m,
					Bottom: m,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return pg.line.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return pg.theme.H6("coming soon").Layout(gtx)
			}),
		)
	})
}

func (pg *walletPage) tableLayout(gtx layout.Context, leftLabel, rightLabel decredmaterial.Label, isIcon bool, seed int) layout.Dimensions {
	m := values.MarginPadding0
	if seed > 0 {
		m = values.MarginPadding5
	}
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if isIcon {
				inset := layout.Inset{
					Right: values.MarginPadding10,
					Top:   m,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return pg.walletIcon.Layout(gtx)
				})
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return leftLabel.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if isIcon {
						if seed > 0 {
							txt := pg.theme.Caption("Not backed up")
							txt.Color = pg.theme.Color.Orange
							inset := layout.Inset{
								Bottom: values.MarginPadding5,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return txt.Layout(gtx)
							})
						}
					}
					return layout.Dimensions{}
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				inset := layout.Inset{
					Right: values.MarginPadding10,
					Top:   m,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return rightLabel.Layout(gtx)
				})
			})
		}),
	)
}

func (pg *walletPage) walletAccountsLayout(gtx layout.Context, name, totalBal, spendableBal string, common pageCommon) layout.Dimensions {
	pg.accountIcon = &widget.Image{Src: paint.NewImageOp(common.icons.accountIcon)}
	if name == "imported" {
		pg.accountIcon = &widget.Image{Src: paint.NewImageOp(common.icons.importedAccountIcon)}
	}
	pg.accountIcon.Scale = 0.8

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			pg.line.Width = gtx.Constraints.Max.X
			pg.line.Color = common.theme.Color.Background
			m := values.MarginPadding10
			return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
				return pg.line.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					inset := layout.Inset{
						Right: values.MarginPadding10,
						Top:   values.MarginPadding15,
					}
					return inset.Layout(gtx, func(gtx C) D {
						return pg.accountIcon.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Right: values.MarginPadding10,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										acctName := strings.Title(strings.ToLower(name))
										return pg.theme.H6(acctName).Layout(gtx)
									}),
									layout.Flexed(1, func(gtx C) D {
										return layout.E.Layout(gtx, func(gtx C) D {
											inset := layout.Inset{
												Right: values.MarginPadding10,
											}
											return inset.Layout(gtx, func(gtx C) D {
												return layoutBalance(gtx, totalBal, common)
											})
										})
									}),
								)
							})
						}),
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Right: values.MarginPadding10,
							}
							return inset.Layout(gtx, func(gtx C) D {
								spendibleLabel := pg.theme.Body2("Spendable")
								spendibleLabel.Color = pg.theme.Color.Gray
								spendibleBalLabel := pg.theme.Body2(spendableBal)
								spendibleBalLabel.Color = pg.theme.Color.Gray
								return pg.tableLayout(gtx, spendibleLabel, spendibleBalLabel, false, 0)
							})
						}),
					)
				}),
			)
		}),
	)
}

func (pg *walletPage) backupSeedNotification(gtx layout.Context, common pageCommon) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	color := common.theme.Color.InvText
	return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				icon := common.icons.walletAlertIcon
				icon.Scale = 0.24
				return icon.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				inset := layout.Inset{
					Left: values.MarginPadding10,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := pg.theme.Body2("Back up seed phrase")
							txt.Color = color
							return txt.Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							txt := pg.theme.Caption("Verify your seed phrase so you can recover your funds when needed.")
							txt.Color = color
							return txt.Layout(gtx)
						}),
					)
				})
			}),
			layout.Flexed(1, func(gtx C) D {
				inset := layout.Inset{
					Top: values.MarginPadding5,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.backupButton.Layout(gtx)
					})
				})
			}),
		)
	})
}

// drawlayout wraps the page tx and sync section in a card layout
func (pg *walletPage) sectionLayout(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return decredmaterial.Card{Color: pg.theme.Color.Surface, CornerStyle: decredmaterial.RoundedEdge}.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding20).Layout(gtx, body)
	})
}

// Handle handles all widget inputs on the main wallets pg.
func (pg *walletPage) Handle(common pageCommon) {
	for _, b := range pg.walletCollapsible {
		for b.Button.Clicked() {
			b.IsExpanded = !b.IsExpanded
		}
	}

	if pg.toAddWalletPage.Clicked() {
		*common.page = PageCreateRestore
	}

	if pg.backupButton.Button.Clicked() {
		*common.page = PageSeedBackup
	}
}

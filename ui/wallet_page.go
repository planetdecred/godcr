package ui

import (
	"image"
	"image/color"
	"strings"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageWallet = "Wallet"

type moreItemText struct {
	signMessage,
	verifyMessage,
	viewProperty,
	rename,
	settings string
}

type walletPage struct {
	walletInfo    *wallet.MultiWalletInfo
	wallet        *wallet.Wallet
	walletAccount **wallet.Account
	theme         *decredmaterial.Theme
	current       wallet.InfoShort

	walletIcon                                 *widget.Image
	accountIcon                                *widget.Image
	addAcct, backupButton                      []decredmaterial.IconButton
	container, accountsList, walletsList, list layout.List
	line                                       *decredmaterial.Line
	toAddWalletPage                            *widget.Clickable
	walletCollapsible                          []*decredmaterial.CollapsibleWithOption
	toAcctDetails                              []*gesture.Click
	text                                       moreItemText
	errChann                                   chan error
	backButton                                 decredmaterial.IconButton
	toPrivacyPage                              decredmaterial.Button
	walletLayoutVisibility                     bool
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
		theme:           common.theme,
		wallet:          common.wallet,
		line:            common.theme.Line(),
		walletAccount:   &win.walletAccount,
		toAddWalletPage: new(widget.Clickable),
		toPrivacyPage:   common.theme.Button(new(widget.Clickable), "Set up mixer for this wallet"),
		backButton:      common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack),
		errChann:        common.errorChannels[PageWallet],
	}

	pg.line.Height = 1

	pg.text = moreItemText{
		signMessage:   "Sign message",
		verifyMessage: "Verify message",
		viewProperty:  "View property",
		rename:        "Rename",
		settings:      "Settings",
	}

	pg.walletCollapsible = make([]*decredmaterial.CollapsibleWithOption, 0)

	pg.addAcct = nil
	pg.backupButton = nil
	pg.toAcctDetails = make([]*gesture.Click, 0)
	pg.toPrivacyPage.Background = pg.theme.Color.Primary
	pg.backButton.Color = common.theme.Color.Text
	pg.backButton.Inset = layout.UniformInset(values.MarginPadding0)

	return func(gtx C) D {
		pg.Handle(common)
		return pg.Layout(gtx, common)
	}
}

// Layout lays out the widgets for the main wallets pg.
func (pg *walletPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	if !pg.walletLayoutVisibility {
		return pg.privacyLayout(gtx, &common)
	}

	if common.info.LoadedWallets == 0 {
		return common.Layout(gtx, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return common.theme.H3("No wallets loaded").Layout(gtx)
			})
		})
	}

	for i := 0; i < common.info.LoadedWallets; i++ {
		pg.walletCollapsible = append(pg.walletCollapsible, pg.theme.CollapsibleWithOption([]decredmaterial.MoreItem{
			{
				Text: pg.text.signMessage,
			},
			{
				Text: pg.text.verifyMessage,
			},
			{
				Text: pg.text.viewProperty,
			},
			{
				Text: pg.text.rename,
			},
			{
				Text: pg.text.settings,
			},
		}))

		addAcctBtn := common.theme.IconButton(new(widget.Clickable), common.icons.contentAdd)
		addAcctBtn.Inset = layout.UniformInset(values.MarginPadding0)
		addAcctBtn.Size = values.MarginPadding25
		addAcctBtn.Background = color.RGBA{}
		addAcctBtn.Color = common.theme.Color.Text
		pg.addAcct = append(pg.addAcct, addAcctBtn)

		backupBtn := common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowForward)
		backupBtn.Color = common.theme.Color.Surface
		backupBtn.Inset = layout.UniformInset(values.MarginPadding0)
		backupBtn.Size = values.MarginPadding20
		pg.backupButton = append(pg.backupButton, backupBtn)
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
		wIndex := i

		pg.updateAcctDetailsButtons(&accounts)

		collapsibleHeader := func(gtx C) D {
			walName := strings.Title(strings.ToLower(wn))
			walletNameLabel := pg.theme.Body1(walName)
			walletBalLabel := pg.theme.Body1(wb)
			walletBalLabel.Color = pg.theme.Color.Gray
			return pg.tableLayout(gtx, walletNameLabel, walletBalLabel, true, len(seed))
		}

		collapsibleFooter := func(gtx C) D {
			return pg.backupSeedNotification(gtx, common, i)
		}

		collapsibleBody := func(gtx C) D {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				gtx.Constraints.Min.Y = 100

				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.accountsList.Layout(gtx, len(accounts), func(gtx C, x int) D {
							accountsName := accounts[x].Name
							totalBalance := accounts[x].TotalBalance
							spendable := dcrutil.Amount(accounts[x].SpendableBalance).String()

							click := pg.toAcctDetails[x]
							pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
							click.Add(gtx.Ops)
							pg.goToAcctDetails(gtx, common, &accounts[x], wIndex, click)
							return pg.walletAccountsLayout(gtx, accountsName, totalBalance, spendable, common)
						})
					}),
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
									return pg.addAcct[i].Layout(gtx)
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

		return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
			if len(seed) > 0 {
				return pg.walletCollapsible[i].Layout(gtx, collapsibleHeader, collapsibleBody, collapsibleFooter)
			}
			return pg.walletCollapsible[i].Layout(gtx, collapsibleHeader, collapsibleBody, nil)
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
		m = values.MarginPaddingMinus5
	}

	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			if isIcon {
				inset := layout.Inset{
					Right: values.MarginPadding10,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return pg.walletIcon.Layout(gtx)
				})
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx C) D {
			inset := layout.Inset{
				Top: m,
			}
			return inset.Layout(gtx, func(gtx C) D {
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
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return rightLabel.Layout(gtx)
			})
		}),
	)
}

func (pg *walletPage) walletAccountsLayout(gtx layout.Context, name, totalBal, spendableBal string, common pageCommon) layout.Dimensions {
	pg.accountIcon = common.icons.accountIcon
	if name == "imported" {
		pg.accountIcon = common.icons.importedAccountIcon
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
								Right: values.MarginPadding20,
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

func (pg *walletPage) backupSeedNotification(gtx layout.Context, common pageCommon, i int) layout.Dimensions {
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
						return pg.backupButton[i].Layout(gtx)
					})
				})
			}),
		)
	})
}

func (pg *walletPage) sectionLayout(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return pg.theme.Card().Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding20).Layout(gtx, body)
	})
}

func (pg *walletPage) updateAcctDetailsButtons(walAcct *[]wallet.Account) {
	if len(*walAcct) != len(pg.toAcctDetails) {
		for i := 0; i < len(*walAcct); i++ {
			pg.toAcctDetails = append(pg.toAcctDetails, &gesture.Click{})
		}
	}
}

func (pg *walletPage) goToAcctDetails(gtx layout.Context, common pageCommon, acct *wallet.Account, index int, click *gesture.Click) {
	for _, e := range click.Events(gtx) {
		if e.Type == gesture.TypeClick {
			*pg.walletAccount = acct
			*common.page = PageAccountDetails
			*common.selectedWallet = index
		}
	}
}

func (pg *walletPage) Handle(common pageCommon) {
	for index, b := range pg.walletCollapsible {
		for b.Collapsible.Button.Clicked() {
			b.Collapsible.IsExpanded = !b.Collapsible.IsExpanded
		}

		for i, t := range b.Items {
			if i > 0 {
				for t.Button.Clicked() {
					*common.selectedWallet = index
					switch b.Items[i].Text {
					case pg.text.signMessage:
						*common.page = PageSignMessage
					case pg.text.verifyMessage:
						*common.page = PageVerifyMessage
					case pg.text.settings:
						*common.page = PageHelp
					case pg.text.rename:
						*common.page = PageAbout
					case pg.text.viewProperty:
						*common.page = PageHelp
					}
					b.Hide()
				}
			}
		}
	}

	for i, b := range pg.addAcct {
		for b.Button.Clicked() {
			walletID := pg.walletInfo.Wallets[i].ID
			go func() {
				common.modalReceiver <- &modalLoad{
					template: CreateAccountTemplate,
					title:    "Create new account",
					confirm: func(name string, passphrase string) {
						pg.wallet.AddAccount(walletID, name, []byte(passphrase), pg.errChann)
					},
					confirmText: "Create",
					cancel:      common.closeModal,
					cancelText:  "Cancel",
				}
			}()
			break
		}
	}

	for i, b := range pg.backupButton {
		for b.Button.Clicked() {
			*common.selectedWallet = i
			pg.current = pg.walletInfo.Wallets[i]
			*common.page = PageSeedBackup
		}
	}

	for pg.toAddWalletPage.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: CreateWalletTemplate,
				title:    "Create new wallet",
				confirm: func(name string, passphrase string) {
					pg.wallet.CreateWallet(name, passphrase, pg.errChann)
				},
				confirmText: "Create",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
		break
	}

	if pg.backButton.Button.Clicked() {
		pg.walletLayoutVisibility = true
	}

	if pg.toPrivacyPage.Button.Clicked() {
		*common.page = PagePrivacy
	}
}

func (pg *walletPage) privacyLayout(gtx layout.Context, c *pageCommon) layout.Dimensions {
	body := c.Layout(gtx, func(gtx C) D {
		return c.theme.Card().Layout(gtx, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.W.Layout(gtx, func(gtx C) D {
									return layout.Inset{Right: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
										return pg.backButton.Layout(gtx)
									})
								})
							}),
							layout.Rigid(func(gtx C) D {
								return pg.theme.H6("Privacy").Layout(gtx)
							}),
						)
					})
				}),
				layout.Flexed(1, func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Flex{Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
											c.icons.transactionFingerPrintIcon.Scale = 0.09
											return c.icons.transactionFingerPrintIcon.Layout(gtx)
										})
									}),
									layout.Rigid(func(gtx C) D {
										c.icons.arrowFowardIcon.Scale = 0.18
										return c.icons.arrowFowardIcon.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Left: values.MarginPadding5, Right: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
											c.icons.mixer.Scale = 0.25
											return c.icons.mixer.Layout(gtx)
										})
									}),
									layout.Rigid(func(gtx C) D {
										c.icons.arrowFowardIcon.Scale = 0.18
										return c.icons.arrowFowardIcon.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
											c.icons.transactionIcon.Scale = 0.09
											return c.icons.transactionIcon.Layout(gtx)
										})
									}),
								)
							}),
							layout.Rigid(func(gtx C) D {
								txt := pg.theme.H6("How does CoinShuffle++ mixer enhance your privacy?")
								txt2 := pg.theme.Body1("CoinShuffle++ mixer can mix your DCRs through CoinJoin transactions.")
								txt3 := pg.theme.Body1("Using mixed DCRs protects you from exposing your financial activities to the public (e.g. how much you own, who pays you).")
								txt.Alignment, txt2.Alignment, txt3.Alignment = text.Middle, text.Middle, text.Middle

								return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return txt.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
											return txt2.Layout(gtx)
										})
									}),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
											return txt3.Layout(gtx)
										})
									}),
								)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
						return pg.toPrivacyPage.Layout(gtx)
					})
				}),
			)
		})
	})
	return body
}

package ui

import (
	"image"
	"image/color"
	"strings"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageWallet = "Wallet"

type action struct {
	button   decredmaterial.IconButton
	topInset float32
}

type optionMenuItem struct {
	text   string
	button *widget.Clickable
	page   string
}

type walletPage struct {
	walletInfo                                 *wallet.MultiWalletInfo
	wallet                                     *wallet.Wallet
	walletAccount                              **wallet.Account
	theme                                      *decredmaterial.Theme
	current                                    wallet.InfoShort
	walletIcon                                 *widget.Image
	accountIcon                                *widget.Image
	walletAlertIcon                            *widget.Image
	addAcct, backupButton                      []decredmaterial.IconButton
	container, accountsList, walletsList, list layout.List
	line                                       *decredmaterial.Line
	toAddWalletPage                            *widget.Clickable
	collapsibles                               []*decredmaterial.Collapsible
	actions                                    map[int]action
	toAcctDetails                              []*gesture.Click
	iconButton                                 decredmaterial.IconButton
	errChann                                   chan error
	card                                       decredmaterial.Card
	popupTopInset                              float32
	backdrop                                   *widget.Clickable
	optionsMenuCard                            decredmaterial.Card
	optionsMenuItems                           []optionMenuItem
}

func (win *Window) WalletPage(common pageCommon) layout.Widget {
	pg := &walletPage{
		walletInfo:      win.walletInfo,
		container:       layout.List{Axis: layout.Vertical},
		accountsList:    layout.List{Axis: layout.Vertical},
		walletsList:     layout.List{Axis: layout.Vertical},
		list:            layout.List{Axis: layout.Vertical},
		theme:           common.theme,
		wallet:          common.wallet,
		line:            common.theme.Line(),
		card:            common.theme.Card(),
		walletAccount:   &win.walletAccount,
		toAddWalletPage: new(widget.Clickable),
		backdrop:        new(widget.Clickable),
		errChann:        common.errorChannels[PageWallet],
		popupTopInset:   -1,
	}

	pg.line.Height = 1
	pg.iconButton = decredmaterial.IconButton{
		IconButtonStyle: material.IconButtonStyle{
			Icon:       pg.theme.NavMoreIcon,
			Size:       unit.Dp(25),
			Background: color.NRGBA{},
			Color:      pg.theme.Color.Text,
			Inset:      layout.UniformInset(unit.Dp(0)),
		},
	}

	pg.optionsMenuCard = decredmaterial.Card{Color: pg.theme.Color.Surface}
	pg.optionsMenuCard.Radius = decredmaterial.CornerRadius{5, 5, 5, 5}

	pg.walletIcon = &widget.Image{Src: paint.NewImageOp(common.icons.walletIcon)}
	pg.walletIcon.Scale = 1

	pg.walletAlertIcon = common.icons.walletAlertIcon
	pg.walletAlertIcon.Scale = 1

	pg.optionsMenuItems = []optionMenuItem{
		{
			text:   "Sign message",
			button: new(widget.Clickable),
			page:   PageSignMessage,
		},
		{
			text:   "Verify message",
			button: new(widget.Clickable),
			page:   PageVerifyMessage,
		},
		{
			text:   "Settings",
			button: new(widget.Clickable),
			page:   PageHelp,
		},
		{
			text:   "Rename",
			button: new(widget.Clickable),
			page:   PageAbout,
		},
		{
			text:   "View property",
			button: new(widget.Clickable),
			page:   PageHelp,
		},
		{
			text:   "Privacy",
			button: new(widget.Clickable),
			page:   PagePrivacy,
		},
	}

	pg.collapsibles = make([]*decredmaterial.Collapsible, 0)
	pg.actions = make(map[int]action)

	pg.addAcct = nil
	pg.backupButton = nil
	pg.toAcctDetails = make([]*gesture.Click, 0)

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
				return common.theme.H3("No wallets loaded").Layout(gtx)
			})
		})
	}

	for index := 0; index < common.info.LoadedWallets; index++ {
		i := index
		collapsible := pg.theme.Collapsible()
		pg.collapsibles = append(pg.collapsibles, collapsible)

		if _, ok := pg.actions[i]; !ok {
			iconButton := pg.iconButton
			iconButton.Button = new(widget.Clickable)
			pg.actions[i] = action{
				button: iconButton,
			}
		}

		addAcctBtn := common.theme.IconButton(new(widget.Clickable), common.icons.contentAdd)
		addAcctBtn.Inset = layout.UniformInset(values.MarginPadding0)
		addAcctBtn.Size = values.MarginPadding25
		addAcctBtn.Background = color.NRGBA{}
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
			layout.Expanded(func(gtx C) D {
				if pg.popupTopInset > -1 {
					return pg.backdrop.Layout(gtx)
				}
				return D{}
			}),
			layout.Expanded(func(gtx C) D {
				if pg.popupTopInset > -1 {
					return layout.NE.Layout(gtx, func(gtx C) D {
						return layout.Inset{
							Top:   unit.Dp(pg.popupTopInset),
							Right: unit.Dp(10),
						}.Layout(gtx, pg.layoutOptionsMenu)
					})
				}
				return D{}
			}),
		)
	}
	return common.Layout(gtx, body)
}

func (pg *walletPage) layoutOptionsMenu(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			border := widget.Border{Color: pg.theme.Color.Background, CornerRadius: unit.Dp(5), Width: unit.Dp(2)}
			return border.Layout(gtx, func(gtx C) D {
				return pg.optionsMenuCard.Layout(gtx, func(gtx C) D {
					return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx C) D {
						return (&layout.List{Axis: layout.Vertical}).Layout(gtx, len(pg.optionsMenuItems), func(gtx C, i int) D {
							return material.Clickable(gtx, pg.optionsMenuItems[i].button, func(gtx C) D {
								return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
									return pg.theme.Body2(pg.optionsMenuItems[i].text).Layout(gtx)
								})
							})
						})
					})
				})
			})
		}),
	)
}

func (pg *walletPage) walletSection(gtx layout.Context, common pageCommon) layout.Dimensions {
	nextTopInset := float32(0)
	return pg.walletsList.Layout(gtx, len(common.info.Wallets), func(gtx C, i int) D {
		var headerDims, bodyDims, footerDims layout.Dimensions

		accounts := common.info.Wallets[i].Accounts
		pg.updateAcctDetailsButtons(&accounts)

		collapsibleMore := func(gtx C) D {
			return pg.actions[i].button.Layout(gtx)
		}

		collapsibleHeader := func(gtx C) D {
			headerDims = pg.layoutCollapsibleHeader(gtx, common.info.Wallets[i])
			return headerDims
		}

		collapsibleBody := func(gtx C) D {
			bodyDims = layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				gtx.Constraints.Min.Y = 100

				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.accountsList.Layout(gtx, len(accounts), func(gtx C, x int) D {
							click := pg.toAcctDetails[x]
							pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
							click.Add(gtx.Ops)
							pg.goToAcctDetails(gtx, common, &accounts[x], i, click)
							return pg.walletAccountsLayout(gtx, accounts[x].Name, accounts[x].TotalBalance, dcrutil.Amount(accounts[x].SpendableBalance).String(), common)
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

			return bodyDims
		}

		footerDims = layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
			var children []layout.FlexChild
			children = append(children, layout.Rigid(func(gtx C) D {
				return pg.collapsibles[i].Layout(gtx, collapsibleHeader, collapsibleBody, collapsibleMore)
			}))

			if len(common.info.Wallets[i].Seed) > 0 {
				children = append(children, layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: unit.Dp(-10)}.Layout(gtx, func(gtx C) D {
						pg.card.Color = pg.theme.Color.Orange
						pg.card.Radius = decredmaterial.CornerRadius{SW: 10, SE: 10}
						return pg.card.Layout(gtx, func(gtx C) D {
							return pg.backupSeedNotification(gtx, common, i)
						})
					})
				}))
			}
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		})

		currentAction := pg.actions[i]
		if i == 0 {
			currentAction.topInset = float32(headerDims.Size.Y)
			nextTopInset = float32(headerDims.Size.Y + footerDims.Size.Y)
		} else {
			currentAction.topInset = nextTopInset
			nextTopInset += float32(headerDims.Size.Y + bodyDims.Size.Y + footerDims.Size.Y)
		}
		pg.actions[i] = currentAction
		return footerDims
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

func (pg *walletPage) layoutCollapsibleHeader(gtx layout.Context, walletInfo wallet.InfoShort) D {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Right: values.MarginPadding10,
			}.Layout(gtx, func(gtx C) D {
				return pg.walletIcon.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.theme.Body1(strings.Title(strings.ToLower(walletInfo.Name))).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if len(walletInfo.Seed) > 0 {
						txt := pg.theme.Caption("Not backed up")
						txt.Color = pg.theme.Color.Orange
						return txt.Layout(gtx)
					}
					return D{}
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				balanceLabel := pg.theme.Body1(walletInfo.Balance)
				balanceLabel.Color = pg.theme.Color.Gray
				return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, balanceLabel.Layout)
			})
		}),
	)
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
				return pg.walletAlertIcon.Layout(gtx)
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
	pg.card.Color = pg.theme.Color.Surface
	pg.card.Radius = decredmaterial.CornerRadius{}

	return pg.card.Layout(gtx, func(gtx C) D {
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
	for pg.backdrop.Clicked() {
		pg.popupTopInset = -1
	}

	for index, action := range pg.actions {
		for action.button.Button.Clicked() {
			pg.popupTopInset = action.topInset
			*common.selectedWallet = index
		}
	}

	for index := range pg.optionsMenuItems {
		for pg.optionsMenuItems[index].button.Clicked() {
			*common.page = pg.optionsMenuItems[index].page
			pg.popupTopInset = -1
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

	select {
	case err := <-pg.errChann:
		if err.Error() == "invalid_passphrase" {
			e := "Password is incorrect"
			common.Notify(e, false)
			return
		}
		common.Notify(err.Error(), false)
	default:
	}
}

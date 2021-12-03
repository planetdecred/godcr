package wallets

import (
	"image/color"
	"sync"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/page/seedbackup"
	"github.com/planetdecred/godcr/ui/values"
)

const WalletPageID = components.WalletsPageID

type (
	C = layout.Context
	D = layout.Dimensions
)

type walletListItem struct {
	wal      *dcrlibwallet.Wallet
	accounts []*dcrlibwallet.Account

	totalBalance string
	optionsMenu  []menuItem
	accountsList *decredmaterial.ClickableList

	// normal wallets
	collapsible         *decredmaterial.CollapsibleWithOption
	addAcctClickable    *decredmaterial.Clickable
	backupAcctClickable *decredmaterial.Clickable

	// watch only
	moreButton decredmaterial.IconButton
}

type menuItem struct {
	text     string
	id       string
	button   *widget.Clickable
	action   func(*load.Load)
	separate bool
}

type WalletPage struct {
	*load.Load

	multiWallet *dcrlibwallet.MultiWallet
	listItems   []*walletListItem

	walletIcon               *decredmaterial.Image
	walletAlertIcon          *decredmaterial.Image
	container                *widget.List
	walletsList              layout.List
	watchWalletsList         *decredmaterial.ClickableList
	iconButton               decredmaterial.IconButton
	card                     decredmaterial.Card
	backdrop                 *widget.Clickable
	optionsMenuCard          decredmaterial.Card
	addWalletMenu            []menuItem
	openPopupIndex           int
	openAddWalletPopupButton *decredmaterial.Clickable
	isAddWalletMenuOpen      bool
	watchOnlyWalletLabel     decredmaterial.Label
	watchOnlyWalletIcon      *decredmaterial.Image
	shadowBox                *decredmaterial.Shadow
	separator                decredmaterial.Line
	addAcctIcon              *decredmaterial.Icon
	backupAcctIcon           *decredmaterial.Icon

	listLock sync.Mutex
}

func NewWalletPage(l *load.Load) *WalletPage {
	pg := &WalletPage{
		Load:        l,
		multiWallet: l.WL.MultiWallet,
		container: &widget.List{
			List: layout.List{Axis: layout.Vertical},
		},
		walletsList:              layout.List{Axis: layout.Vertical},
		watchWalletsList:         l.Theme.NewClickableList(layout.Vertical),
		card:                     l.Theme.Card(),
		backdrop:                 new(widget.Clickable),
		openAddWalletPopupButton: l.Theme.NewClickable(false),
		openPopupIndex:           -1,
		shadowBox:                l.Theme.Shadow(),
		separator:                l.Theme.Separator(),
		addAcctIcon:              decredmaterial.NewIcon(l.Icons.ContentAdd),
		backupAcctIcon:           decredmaterial.NewIcon(l.Icons.NavigationArrowForward),
	}

	pg.openAddWalletPopupButton.Radius = decredmaterial.Radius(24)

	pg.separator.Color = l.Theme.Color.Gray2

	pg.watchOnlyWalletLabel = pg.Theme.Body1(values.String(values.StrWatchOnlyWallets))
	pg.watchOnlyWalletLabel.Color = pg.Theme.Color.GrayText2

	pg.iconButton = decredmaterial.IconButton{
		IconButtonStyle: decredmaterial.IconButtonStyle{
			Size:  unit.Dp(25),
			Inset: layout.UniformInset(unit.Dp(0)),
		},
	}
	pg.iconButton.ChangeColorStyle(pg.Theme.Styles.IconButtonColorStyle)

	pg.optionsMenuCard = decredmaterial.Card{Color: pg.Theme.Color.Surface}
	pg.optionsMenuCard.Radius = decredmaterial.Radius(5)

	pg.walletIcon = pg.Icons.WalletIcon

	pg.walletAlertIcon = pg.Icons.WalletAlertIcon

	pg.initializeFloatingMenu()
	pg.watchOnlyWalletIcon = pg.Icons.WatchOnlyWalletIcon

	return pg
}

func (pg *WalletPage) ID() string {
	return WalletPageID
}

func (pg *WalletPage) OnResume() {
	pg.loadWalletAndAccounts()
}

func (pg *WalletPage) loadWalletAndAccounts() {
	wallets := pg.WL.SortedWalletList()

	listItems := make([]*walletListItem, 0)

	for _, wal := range wallets {
		accountsResult, err := wal.GetAccountsRaw()
		if err != nil {
			continue
		}

		var totalBalance int64
		for _, acc := range accountsResult.Acc {
			totalBalance += acc.TotalBalance
		}

		listItem := &walletListItem{
			wal:      wal,
			accounts: accountsResult.Acc,

			totalBalance: dcrutil.Amount(totalBalance).String(),
			optionsMenu:  pg.getWalletMenu(wal),
			accountsList: pg.Theme.NewClickableList(layout.Vertical),
		}

		if wal.IsWatchingOnlyWallet() {
			moreBtn := decredmaterial.IconButton{
				IconButtonStyle: decredmaterial.IconButtonStyle{
					Button: new(widget.Clickable),
					Icon:   pg.Icons.NavigationMore,
					Size:   values.MarginPadding25,
					Inset:  layout.UniformInset(values.MarginPadding0),
				},
			}
			moreBtn.ChangeColorStyle(&values.ColorStyle{
				Background: color.NRGBA{},
				Foreground: pg.Theme.Color.Text,
			})
			listItem.moreButton = moreBtn
		} else {
			listItem.addAcctClickable = pg.Theme.NewClickable(true)

			backupClickable := pg.Theme.NewClickable(false)
			backupClickable.ChangeStyle(&values.ClickableStyle{Color: pg.Theme.Color.OrangeRipple})
			backupClickable.Radius = decredmaterial.CornerRadius{BottomRight: 14, BottomLeft: 14}
			listItem.backupAcctClickable = backupClickable

			listItem.collapsible = pg.Theme.CollapsibleWithOption()
		}
		listItems = append(listItems, listItem)
	}

	pg.listLock.Lock()
	pg.listItems = listItems
	pg.listLock.Unlock()
}

func (pg *WalletPage) initializeFloatingMenu() {
	pg.addWalletMenu = []menuItem{
		{
			text:   values.String(values.StrCreateANewWallet),
			button: new(widget.Clickable),
			action: pg.showAddWalletModal,
		},
		{
			text:   values.String(values.StrImportExistingWallet),
			button: new(widget.Clickable),
			action: func(l *load.Load) {
				l.ChangeWindowPage(NewRestorePage(pg.Load), true)
			},
		},
		{
			text:   values.String(values.StrImportWatchingOnlyWallet),
			button: new(widget.Clickable),
			action: pg.showImportWatchOnlyWalletModal,
		},
	}
}

func (pg *WalletPage) getWalletMenu(wal *dcrlibwallet.Wallet) []menuItem {
	if wal.IsWatchingOnlyWallet() {
		return pg.getWatchOnlyWalletMenu(wal)
	}

	return []menuItem{
		{
			text:   values.String(values.StrSignMessage),
			button: new(widget.Clickable),
			id:     SignMessagePageID,
		},
		{
			text:     values.String(values.StrViewProperty),
			button:   new(widget.Clickable),
			separate: true,
			action:   func(load *load.Load) {},
		},
		{
			text:     values.String(values.StrStakeShuffle),
			button:   new(widget.Clickable),
			separate: true,
			id:       PrivacyPageID,
		},
		{
			text:   values.String(values.StrRename),
			button: new(widget.Clickable),
			action: func(l *load.Load) {
				textModal := modal.NewTextInputModal(l).
					Hint("Wallet name").
					PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
					PositiveButton(values.String(values.StrRename), func(newName string, tim *modal.TextInputModal) bool {
						err := pg.multiWallet.RenameWallet(wal.ID, newName)
						if err != nil {
							pg.Toast.NotifyError(err.Error())
							return false
						}
						return true
					})

				textModal.Title(values.String(values.StrRenameWalletSheetTitle)).
					NegativeButton(values.String(values.StrCancel), func() {})
				textModal.Show()
			},
		},
		{
			text:   values.String(values.StrSettings),
			button: new(widget.Clickable),
			id:     WalletSettingsPageID,
		},
	}
}

func (pg *WalletPage) getWatchOnlyWalletMenu(wal *dcrlibwallet.Wallet) []menuItem {
	return []menuItem{
		{
			text:   values.String(values.StrSettings),
			button: new(widget.Clickable),
			id:     WalletSettingsPageID,
		},
		{
			text:   values.String(values.StrRename),
			button: new(widget.Clickable),
			action: func(l *load.Load) {
				textModal := modal.NewTextInputModal(l).
					Hint("Wallet name").
					PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
					PositiveButton(values.String(values.StrRename), func(newName string, tim *modal.TextInputModal) bool {
						//TODO
						err := pg.multiWallet.RenameWallet(wal.ID, newName)
						if err != nil {
							pg.Toast.NotifyError(err.Error())
						} else {
							pg.Toast.Notify("Wallet renamed")
						}
						return true
					})

				textModal.Title(values.String(values.StrRenameWalletSheetTitle)).
					NegativeButton(values.String(values.StrCancel), func() {})
				textModal.Show()
			},
		},
	}
}

func (pg *WalletPage) showAddWalletModal(l *load.Load) {
	modal.NewCreatePasswordModal(l).
		Title("Create new wallet").
		EnableName(true).
		ShowWalletInfoTip(true).
		PasswordCreated(func(walletName, password string, m *modal.CreatePasswordModal) bool {
			go func() {
				_, err := pg.multiWallet.CreateNewWallet(walletName, password, dcrlibwallet.PassphraseTypePass)
				if err != nil {
					m.SetError(err.Error())
					m.SetLoading(false)
					return
				}
				pg.loadWalletAndAccounts()
				pg.Toast.Notify("Wallet created")
				m.Dismiss()
			}()
			return false
		}).Show()
}

func (pg *WalletPage) showImportWatchOnlyWalletModal(l *load.Load) {
	modal.NewCreateWatchOnlyModal(l).
		WatchOnlyCreated(func(walletName, extPubKey string, m *modal.CreateWatchOnlyModal) bool {
			go func() {
				_, err := pg.multiWallet.CreateWatchOnlyWallet(walletName, extPubKey)
				if err != nil {
					pg.Toast.NotifyError(err.Error())
					m.SetError(err.Error())
					m.SetLoading(false)
				} else {
					// pg.wallet.GetMultiWalletInfo() TODO
					pg.Toast.Notify(values.String(values.StrWatchOnlyWalletImported))
					m.Dismiss()
				}
			}()
			return false
		}).Show()
}

// Layout lays out the widgets for the main wallets pg.
func (pg *WalletPage) Layout(gtx layout.Context) layout.Dimensions {
	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.walletSection(gtx)
		},
		func(gtx C) D {
			return pg.watchOnlyWalletSection(gtx)
		},
	}

	body := func(gtx C) D {
		return layout.Stack{Alignment: layout.SE}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return pg.Theme.List(pg.container).Layout(gtx, 1, func(gtx C, i int) D {
					return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
						dims := pageContent[i](gtx)
						if pg.isAddWalletMenuOpen || pg.openPopupIndex != -1 {
							dims.Size.Y += 60
						}
						return dims
					})
				})
			}),
			layout.Stacked(func(gtx C) D {
				return pg.layoutAddWalletSection(gtx)
			}),
		)
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return components.UniformPadding(gtx, body)
		}),
		layout.Expanded(func(gtx C) D {
			if pg.isAddWalletMenuOpen || pg.openPopupIndex != -1 {
				return pg.backdrop.Layout(gtx)
			}
			return D{}
		}),
	)
}

func (pg *WalletPage) layoutOptionsMenu(gtx layout.Context, optionsMenuIndex int, listItem *walletListItem) {
	if pg.openPopupIndex != optionsMenuIndex {
		return
	}

	var leftInset float32
	if listItem.wal.IsWatchingOnlyWallet() {
		leftInset = -35
	} else {
		leftInset = -120
	}

	inset := layout.Inset{
		Top:  unit.Dp(30),
		Left: unit.Dp(leftInset),
	}
	menu := listItem.optionsMenu

	m := op.Record(gtx.Ops)
	inset.Layout(gtx, func(gtx C) D {
		width := unit.Value{U: unit.UnitDp, V: 150}
		gtx.Constraints.Max.X = gtx.Px(width)
		return pg.shadowBox.Layout(gtx, func(gtx C) D {
			return pg.optionsMenuCard.Layout(gtx, func(gtx C) D {
				return (&layout.List{Axis: layout.Vertical}).Layout(gtx, len(menu), func(gtx C, i int) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							card := pg.Theme.Card()
							card.Radius = decredmaterial.Radius(0)
							return card.HoverableLayout(gtx, menu[i].button, func(gtx C) D {
								return material.Clickable(gtx, menu[i].button, func(gtx C) D {
									m10 := values.MarginPadding10
									return layout.Inset{Top: m10, Bottom: m10, Left: m10, Right: m10}.Layout(gtx, func(gtx C) D {
										gtx.Constraints.Min.X = gtx.Constraints.Max.X
										return pg.Theme.Body1(menu[i].text).Layout(gtx)
									})
								})
							})
						}),
						layout.Rigid(func(gtx C) D {
							if menu[i].separate {
								m := values.MarginPadding5
								return layout.Inset{Top: m, Bottom: m}.Layout(gtx, pg.separator.Layout)
							}
							return D{}
						}),
					)
				})
			})
		})
	})
	op.Defer(gtx.Ops, m.Stop())
}

func (pg *WalletPage) walletSection(gtx layout.Context) layout.Dimensions {
	pg.listLock.Lock()
	listItems := pg.listItems
	pg.listLock.Unlock()

	return pg.walletsList.Layout(gtx, len(listItems), func(gtx C, i int) D {

		listItem := listItems[i]

		if listItem.wal.IsWatchingOnlyWallet() {
			return D{}
		}

		collapsibleMore := func(gtx C) {
			pg.layoutOptionsMenu(gtx, i, listItem)
		}

		collapsibleHeader := func(gtx C) D {
			return pg.layoutCollapsibleHeader(gtx, listItem)
		}

		collapsibleBody := func(gtx C) D {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				gtx.Constraints.Min.Y = 100

				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Left:  values.MarginPadding38,
							Right: values.MarginPadding10,
						}.Layout(gtx, pg.Theme.Separator().Layout)
					}),
					layout.Rigid(func(gtx C) D {
						return listItem.accountsList.Layout(gtx, len(listItem.accounts), func(gtx C, x int) D {
							return pg.walletAccountsLayout(gtx, listItem.accounts[x])
						})
					}),
					layout.Rigid(func(gtx C) D {
						return listItem.addAcctClickable.Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.Inset{
											Right: values.MarginPadding10,
											Left:  values.MarginPadding38,
										}.Layout(gtx, func(gtx C) D {
											pg.addAcctIcon.Color = pg.Theme.Color.Gray1
											return pg.addAcctIcon.Layout(gtx, values.MarginPadding25)
										})
									}),
									layout.Rigid(func(gtx C) D {
										txt := pg.Theme.Label(values.TextSize16, values.String(values.StrAddNewAccount))
										txt.Color = pg.Theme.Color.GrayText2
										return txt.Layout(gtx)
									}),
								)
							})
						})

					}),
				)
			})
		}

		return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
			var children []layout.FlexChild
			children = append(children, layout.Rigid(func(gtx C) D {
				return listItem.collapsible.Layout(gtx, collapsibleHeader, collapsibleBody, collapsibleMore)
			}))

			if len(listItem.wal.EncryptedSeed) > 0 {
				children = append(children, layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: unit.Dp(-10)}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								blankLine := pg.Theme.Line(10, gtx.Constraints.Max.X)
								blankLine.Color = pg.Theme.Color.Surface
								return blankLine.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								pg.card.Color = pg.Theme.Color.Danger
								pg.card.Radius = decredmaterial.CornerRadius{BottomLeft: 10, BottomRight: 10}
								return pg.card.Layout(gtx, func(gtx C) D {
									return pg.backupSeedNotification(gtx, listItem)
								})
							}),
						)
					})
				}))
			}
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		})
	})
}

func (pg *WalletPage) watchOnlyWalletSection(gtx layout.Context) layout.Dimensions {
	hasWatchOnly := false

	pg.listLock.Lock()
	listItems := pg.listItems
	pg.listLock.Unlock()

	for _, listItem := range listItems {
		if listItem.wal.IsWatchingOnlyWallet() {
			hasWatchOnly = true
			break
		}
	}
	if !hasWatchOnly {
		return D{}
	}
	card := pg.card
	card.Color = pg.Theme.Color.Surface
	card.Radius = decredmaterial.Radius(10)

	return card.Layout(gtx, func(gtx C) D {
		m := values.MarginPadding20
		return layout.Inset{Top: m, Left: m}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.watchOnlyWalletLabel.Layout),
				layout.Rigid(func(gtx C) D {
					m := values.MarginPadding10
					return layout.Inset{Top: m, Bottom: m}.Layout(gtx, pg.Theme.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return pg.layoutWatchOnlyWallets(gtx)
					})
				}),
			)
		})
	})
}

func (pg *WalletPage) layoutWatchOnlyWallets(gtx layout.Context) D {
	pg.listLock.Lock()
	listItems := pg.listItems
	pg.listLock.Unlock()

	return pg.watchWalletsList.Layout(gtx, len(listItems), func(gtx C, i int) D {

		listItem := listItems[i]

		if !listItem.wal.IsWatchingOnlyWallet() {
			return D{}
		}

		m := values.MarginPadding10
		return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Right: values.MarginPadding10,
							}
							return inset.Layout(gtx, pg.watchOnlyWalletIcon.Layout24dp)
						}),
						layout.Rigid(pg.Theme.Body2(listItem.wal.Name).Layout),
						layout.Flexed(1, func(gtx C) D {
							return layout.E.Layout(gtx, func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return pg.Theme.Body2(listItem.totalBalance).Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										pg.layoutOptionsMenu(gtx, i, listItem)
										return layout.Inset{Top: unit.Dp(-3)}.Layout(gtx, listItem.moreButton.Layout)
									}),
								)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10, Left: values.MarginPadding38, Right: values.MarginPaddingMinus10}.Layout(gtx, func(gtx C) D {
						if i == len(listItems)-1 {
							return D{}
						}
						return pg.Theme.Separator().Layout(gtx)
					})
				}),
			)
		})
	})
}

func (pg *WalletPage) layoutCollapsibleHeader(gtx layout.Context, listItem *walletListItem) D {
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Right: values.MarginPadding10,
			}.Layout(gtx, pg.walletIcon.Layout24dp)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.Theme.Body1(listItem.wal.Name).Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if len(listItem.wal.EncryptedSeed) > 0 {
						txt := pg.Theme.Caption(values.String(values.StrNotBackedUp))
						txt.Color = pg.Theme.Color.Danger
						return txt.Layout(gtx)
					}
					return D{}
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				balanceLabel := pg.Theme.Body1(listItem.totalBalance)
				balanceLabel.Color = pg.Theme.Color.GrayText2
				return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, balanceLabel.Layout)
			})
		}),
	)
}

func (pg *WalletPage) tableLayout(gtx layout.Context, leftLabel, rightLabel decredmaterial.Label) layout.Dimensions {
	m := values.MarginPadding0

	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			inset := layout.Inset{
				Top: m,
			}
			return inset.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return leftLabel.Layout(gtx)
					}),
				)
			})
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, rightLabel.Layout)
		}),
	)
}

func (pg *WalletPage) walletAccountsLayout(gtx layout.Context, account *dcrlibwallet.Account) layout.Dimensions {
	accountIcon := pg.Icons.AccountIcon
	if account.Number == load.MaxInt32 {
		accountIcon = pg.Icons.ImportedAccountIcon
		if account.TotalBalance == 0 {
			return D{}
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			inset := layout.Inset{
				Top:    values.MarginPadding10,
				Left:   values.MarginPadding38,
				Bottom: values.MarginPadding20,
			}
			return inset.Layout(gtx, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						inset := layout.Inset{
							Right: values.MarginPadding10,
							Top:   values.MarginPadding13,
						}
						return inset.Layout(gtx, accountIcon.Layout24dp)
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
											return pg.Theme.Label(values.TextSize18, account.Name).Layout(gtx)
										}),
										layout.Flexed(1, func(gtx C) D {
											return layout.E.Layout(gtx, func(gtx C) D {
												totalBal := dcrutil.Amount(account.TotalBalance).String()
												return components.LayoutBalance(gtx, pg.Load, totalBal)
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
									spendableLabel := pg.Theme.Body2(values.String(values.StrLabelSpendable))
									spendableLabel.Color = pg.Theme.Color.GrayText2

									spendableBal := dcrutil.Amount(account.Balance.Spendable).String()
									spendableBalLabel := pg.Theme.Body2(spendableBal)
									spendableBalLabel.Color = pg.Theme.Color.GrayText2
									return pg.tableLayout(gtx, spendableLabel, spendableBalLabel)
								})
							}),
						)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Left:  values.MarginPadding70,
				Right: values.MarginPadding10,
			}.Layout(gtx, pg.Theme.Separator().Layout)
		}),
	)
}

func (pg *WalletPage) backupSeedNotification(gtx layout.Context, listItem *walletListItem) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	textColor := pg.Theme.Color.InvText
	return listItem.backupAcctClickable.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.walletAlertIcon.Layout24dp(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					inset := layout.Inset{
						Left: values.MarginPadding10,
					}
					return inset.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								txt := pg.Theme.Body2(values.String(values.StrBackupSeedPhrase))
								txt.Color = textColor
								return txt.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								txt := pg.Theme.Caption(values.String(values.StrVerifySeedInfo))
								txt.Color = textColor
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
							return pg.backupAcctIcon.Layout(gtx, values.MarginPadding20)
						})
					})
				}),
			)
		})
	})
}

func (pg *WalletPage) layoutAddWalletMenu(gtx layout.Context) layout.Dimensions {
	inset := layout.Inset{
		Top:  unit.Dp(-100),
		Left: unit.Dp(-130),
	}

	return inset.Layout(gtx, func(gtx C) D {
		return pg.Theme.Shadow().Layout(gtx, func(gtx C) D {
			return pg.optionsMenuCard.Layout(gtx, func(gtx C) D {
				return (&layout.List{Axis: layout.Vertical}).Layout(gtx, len(pg.addWalletMenu), func(gtx C, i int) D {
					return material.Clickable(gtx, pg.addWalletMenu[i].button, func(gtx C) D {
						return layout.UniformInset(unit.Dp(10)).Layout(gtx, pg.Theme.Body2(pg.addWalletMenu[i].text).Layout)
					})
				})
			})
		})
	})
}

func (pg *WalletPage) layoutAddWalletSection(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	return layout.SE.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				if pg.isAddWalletMenuOpen {
					m := op.Record(gtx.Ops)
					pg.layoutAddWalletMenu(gtx)
					op.Defer(gtx.Ops, m.Stop())
				}
				return D{}
			}),
			layout.Rigid(func(gtx C) D {
				icon := pg.Icons.NewWalletIcon

				return decredmaterial.LinearLayout{
					Width:      decredmaterial.WrapContent,
					Height:     decredmaterial.WrapContent,
					Padding:    layout.UniformInset(values.MarginPadding12),
					Background: pg.Theme.Color.Surface,
					Clickable:  pg.openAddWalletPopupButton,
					Border:     decredmaterial.Border{Radius: pg.openAddWalletPopupButton.Radius},
				}.Layout2(gtx, icon.Layout24dp)
			}),
		)
	})
}

func (pg *WalletPage) closePopups() {
	pg.openPopupIndex = -1
	pg.isAddWalletMenuOpen = false
}

func (pg *WalletPage) openPopup(index int) {
	if pg.openPopupIndex >= 0 {
		if pg.openPopupIndex == index {
			pg.closePopups()
			return
		}
		pg.closePopups()
	}

	pg.openPopupIndex = index
}

func (pg *WalletPage) Handle() {
	for pg.backdrop.Clicked() {
		pg.closePopups()
	}

	if ok, selectedItem := pg.watchWalletsList.ItemClicked(); ok {
		pg.listLock.Lock()
		listItem := pg.listItems[selectedItem]
		pg.listLock.Unlock()

		// TODO: find default account using number
		pg.ChangeFragment(NewAcctDetailsPage(pg.Load, listItem.accounts[0]))
	}

	pg.listLock.Lock()
	listItems := pg.listItems
	pg.listLock.Unlock()

	for index, listItem := range listItems {
		*pg.SelectedWallet = index

		if ok, selectedItem := listItem.accountsList.ItemClicked(); ok {
			pg.ChangeFragment(NewAcctDetailsPage(pg.Load, listItem.accounts[selectedItem]))
		}

		if listItem.wal.IsWatchingOnlyWallet() {
			for listItem.moreButton.Button.Clicked() {
				pg.openPopup(index)
			}
		} else {
			for listItem.collapsible.MoreTriggered() {
				pg.openPopup(index)
			}

			for listItem.addAcctClickable.Clicked() {
				walletID := listItem.wal.ID
				textModal := modal.NewTextInputModal(pg.Load).
					Hint("Account name").
					ShowAccountInfoTip(true).
					PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
					PositiveButton(values.String(values.StrCreate), func(accountName string, tim *modal.TextInputModal) bool {
						if accountName != "" {
							modal.NewPasswordModal(pg.Load).
								Title(values.String(values.StrCreateNewAccount)).
								Hint("Spending password").
								NegativeButton(values.String(values.StrCancel), func() {}).
								PositiveButton(values.String(values.StrConfirm), func(password string, pm *modal.PasswordModal) bool {
									go func() {
										wal := pg.multiWallet.WalletWithID(walletID)
										_, err := wal.CreateNewAccount(accountName, []byte(password)) // TODO
										if err != nil {
											pg.Toast.NotifyError(err.Error())
											tim.SetError(err.Error())
										} else {
											pg.Toast.Notify("Account created")
											tim.Dismiss()
										}
										pg.loadWalletAndAccounts()
										pm.Dismiss()
									}()
									return false
								}).Show()
						}
						return true
					})
				textModal.Title(values.String(values.StrCreateNewAccount)).
					NegativeButton(values.String(values.StrCancel), func() {})
				textModal.Show()
				break
			}

			for listItem.backupAcctClickable.Clicked() {
				pg.ChangeFragment(seedbackup.NewBackupInstructionsPage(pg.Load, listItem.wal))
			}
		}

		for _, menu := range listItem.optionsMenu {
			if menu.button.Clicked() {
				switch menu.id {
				case SignMessagePageID:
					pg.ChangeFragment(NewSignMessagePage(pg.Load, listItem.wal))
				case PrivacyPageID:
					pg.ChangeFragment(NewPrivacyPage(pg.Load, listItem.wal))
				case WalletSettingsPageID:
					pg.ChangeFragment(NewWalletSettingsPage(pg.Load, listItem.wal))
				default:
					menu.action(pg.Load)
				}

				pg.openPopupIndex = -1
			}
		}
	}

	for index := range pg.addWalletMenu {
		for pg.addWalletMenu[index].button.Clicked() {
			pg.isAddWalletMenuOpen = false
			pg.addWalletMenu[index].action(pg.Load)
		}
	}

	for pg.openAddWalletPopupButton.Clicked() {
		pg.isAddWalletMenuOpen = !pg.isAddWalletMenuOpen
	}
}

func (pg *WalletPage) OnClose() {
	pg.closePopups()
}

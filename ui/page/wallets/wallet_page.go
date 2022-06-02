package wallets

import (
	"context"
	"image"
	"image/color"
	"sync"

	"gioui.org/gesture"
	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/page/privacy"
	"github.com/planetdecred/godcr/ui/page/seedbackup"
	"github.com/planetdecred/godcr/ui/values"
)

const WalletPageID = components.WalletsPageID

type (
	C = layout.Context
	D = layout.Dimensions
)

type badWalletListItem struct {
	*dcrlibwallet.Wallet
	deleteBtn decredmaterial.Button
}

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
	checkMixerClickable *decredmaterial.Clickable

	// watch only
	moreButton decredmaterial.IconButton
}

type menuItem struct {
	text     string
	id       string
	button   *decredmaterial.Clickable
	action   func(*load.Load)
	separate bool
}

type WalletPage struct {
	*load.Load
	*listeners.TxAndBlockNotificationListener
	ctx       context.Context // page context
	ctxCancel context.CancelFunc
	listLock  sync.Mutex

	multiWallet *dcrlibwallet.MultiWallet

	listItems      []*walletListItem
	badWalletsList []*badWalletListItem
	addWalletMenu  []menuItem

	container   *widget.List
	backdrop    *widget.Clickable
	walletsList layout.List

	walletIcon               *decredmaterial.Image
	walletAlertIcon          *decredmaterial.Image
	watchOnlyWalletIcon      *decredmaterial.Image
	watchWalletsList         *decredmaterial.ClickableList
	openAddWalletPopupButton *decredmaterial.Clickable
	shadowBox                *decredmaterial.Shadow
	addAcctIcon              *decredmaterial.Icon
	backupAcctIcon           *decredmaterial.Icon
	nextIcon                 *decredmaterial.Icon

	card                 decredmaterial.Card
	optionsMenuCard      decredmaterial.Card
	watchOnlyWalletLabel decredmaterial.Label
	separator            decredmaterial.Line

	mt    unit.Value    // option menu top margin padding
	click gesture.Click // page click event

	openPopupIndex      int
	isAddWalletMenuOpen bool
	hasWatchOnly        bool
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
		addAcctIcon:              decredmaterial.NewIcon(l.Theme.Icons.ContentAdd),
		backupAcctIcon:           decredmaterial.NewIcon(l.Theme.Icons.NavigationArrowForward),
	}

	pg.openAddWalletPopupButton.Radius = decredmaterial.Radius(10)

	pg.separator.Color = l.Theme.Color.Gray2

	pg.watchOnlyWalletLabel = pg.Theme.Body1(values.String(values.StrWatchOnlyWallets))
	pg.watchOnlyWalletLabel.Color = pg.Theme.Color.GrayText2

	pg.optionsMenuCard = decredmaterial.Card{Color: pg.Theme.Color.Surface}
	pg.optionsMenuCard.Radius = decredmaterial.Radius(5)

	pg.walletIcon = pg.Theme.Icons.WalletIcon

	pg.walletAlertIcon = pg.Theme.Icons.WalletAlertIcon

	pg.nextIcon = decredmaterial.NewIcon(pg.Theme.Icons.NavigationArrowForward)
	pg.nextIcon.Color = pg.Theme.Color.Primary

	pg.initializeFloatingMenu()
	pg.watchOnlyWalletIcon = pg.Theme.Icons.WatchOnlyWalletIcon

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *WalletPage) ID() string {
	return WalletPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *WalletPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())

	pg.listenForTxNotifications()
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
			pg.hasWatchOnly = true
			listItem.moreButton = pg.Theme.IconButtonWithStyle(
				decredmaterial.IconButtonStyle{
					Button: new(widget.Clickable),
					Icon:   pg.Theme.Icons.NavigationMore,
					Size:   values.MarginPadding25,
					Inset:  layout.UniformInset(values.MarginPadding0),
				},
				&values.ColorStyle{
					Background: color.NRGBA{},
					Foreground: pg.Theme.Color.Text,
				},
			)
		} else {
			listItem.addAcctClickable = pg.Theme.NewClickable(true)

			backupClickable := pg.Theme.NewClickable(false)
			backupClickable.ChangeStyle(&values.ClickableStyle{Color: pg.Theme.Color.OrangeRipple})
			backupClickable.Radius = decredmaterial.CornerRadius{BottomRight: 14, BottomLeft: 14}
			listItem.backupAcctClickable = backupClickable

			checkMixerClickable := pg.Theme.NewClickable(false)
			checkMixerClickable.Radius = decredmaterial.CornerRadius{BottomRight: 14, BottomLeft: 14}
			listItem.checkMixerClickable = checkMixerClickable

			listItem.collapsible = pg.Theme.CollapsibleWithOption()
		}
		listItems = append(listItems, listItem)
	}

	pg.listLock.Lock()
	pg.listItems = listItems
	pg.listLock.Unlock()

	pg.loadBadWallets()
}

func (pg *WalletPage) loadBadWallets() {
	badWallets := pg.WL.MultiWallet.BadWallets()
	pg.badWalletsList = make([]*badWalletListItem, 0, len(badWallets))
	for _, badWallet := range badWallets {
		listItem := &badWalletListItem{
			Wallet:    badWallet,
			deleteBtn: pg.Theme.OutlineButton(values.String(values.StrDeleted)),
		}
		listItem.deleteBtn.Color = pg.Theme.Color.Danger
		listItem.deleteBtn.Inset = layout.Inset{}
		pg.badWalletsList = append(pg.badWalletsList, listItem)
	}
}

func (pg *WalletPage) initializeFloatingMenu() {
	pg.addWalletMenu = []menuItem{
		{
			text:   values.String(values.StrCreateANewWallet),
			button: pg.Theme.NewClickable(true),
			action: pg.showAddWalletModal,
		},
		{
			text:   values.String(values.StrImportExistingWallet),
			button: pg.Theme.NewClickable(true),
			action: func(l *load.Load) {
				// The second nil parameter to NewRestorePage will cause the
				// restore page to pop back to this one after restore completes.
				l.ChangeWindowPage(NewRestorePage(pg.Load, nil), true)
			},
		},
		{
			text:   values.String(values.StrImportWatchingOnlyWallet),
			button: pg.Theme.NewClickable(true),
			action: pg.showImportWatchOnlyWalletModal,
		},
	}
}

func (pg *WalletPage) getWalletMenu(wal *dcrlibwallet.Wallet) []menuItem {
	if wal.IsWatchingOnlyWallet() {
		return pg.getWatchOnlyWalletMenu(wal)
	}
	privacyPageID := privacy.SetupPrivacyPageID
	if wal.AccountMixerConfigIsSet() {
		privacyPageID = privacy.AccountMixerPageID
	}
	return []menuItem{
		{
			text:   values.String(values.StrSignMessage),
			button: pg.Theme.NewClickable(true),
			id:     SignMessagePageID,
		},
		{
			text:     values.String(values.StrViewProperty),
			button:   pg.Theme.NewClickable(true),
			separate: true,
			action:   func(load *load.Load) {},
		},
		{
			text:     values.String(values.StrStakeShuffle),
			button:   pg.Theme.NewClickable(true),
			separate: true,
			id:       privacyPageID,
		},
		{
			text:   values.String(values.StrRename),
			button: pg.Theme.NewClickable(true),
			action: func(l *load.Load) {
				textModal := modal.NewTextInputModal(l).
					Hint(values.String(values.StrWalletName)).
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
			button: pg.Theme.NewClickable(true),
			id:     WalletSettingsPageID,
		},
	}
}

func (pg *WalletPage) getWatchOnlyWalletMenu(wal *dcrlibwallet.Wallet) []menuItem {
	return []menuItem{
		{
			text:   values.String(values.StrRename),
			button: pg.Theme.NewClickable(true),
			action: func(l *load.Load) {
				textModal := modal.NewTextInputModal(l).
					Hint(values.String(values.StrWalletName)).
					PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
					PositiveButton(values.String(values.StrRename), func(newName string, tim *modal.TextInputModal) bool {
						//TODO
						err := pg.multiWallet.RenameWallet(wal.ID, newName)
						if err != nil {
							pg.Toast.NotifyError(err.Error())
						} else {
							pg.Toast.Notify(values.String(values.StrWalletRenamed))
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
			button: pg.Theme.NewClickable(true),
			id:     WalletSettingsPageID,
		},
	}
}

func (pg *WalletPage) showAddWalletModal(l *load.Load) {
	modal.NewCreatePasswordModal(l).
		Title(values.String(values.StrCreateANewWallet)).
		EnableName(true).
		ShowWalletInfoTip(true).
		PasswordCreated(func(walletName, password string, m *modal.CreatePasswordModal) bool {
			go func() {
				wal, err := pg.multiWallet.CreateNewWallet(walletName, password, dcrlibwallet.PassphraseTypePass)
				if err != nil {
					m.SetError(err.Error())
					m.SetLoading(false)
					return
				}
				err = wal.CreateMixerAccounts("mixed", "unmixed", password)
				if err != nil {
					m.SetError(err.Error())
					m.SetLoading(false)
					return
				}
				pg.WL.MultiWallet.SetBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, true)
				pg.loadWalletAndAccounts()
				pg.Toast.Notify(values.String(values.StrWalletCreated))
				m.Dismiss()
			}()
			return false
		}).Show()
}

func (pg *WalletPage) showImportWatchOnlyWalletModal(l *load.Load) {
	modal.NewCreateWatchOnlyModal(l).
		EnableName(true).
		WatchOnlyCreated(func(walletName, extPubKey string, m *modal.CreateWatchOnlyModal) bool {
			go func() {
				_, err := pg.multiWallet.CreateWatchOnlyWallet(walletName, extPubKey)
				if err != nil {
					pg.Toast.NotifyError(err.Error())
					m.SetError(err.Error())
					m.SetLoading(false)
				} else {
					pg.loadWalletAndAccounts()
					pg.Toast.Notify(values.String(values.StrWatchOnlyWalletImported))
					m.Dismiss()
				}
			}()
			return false
		}).Show()
}

// moreOptionPositionEvent tracks the position of the click event on the page
func (pg *WalletPage) moreOptionPositionEvent(gtx layout.Context) {
	setUnitValue := func() {
		pg.mt = unit.Dp(-220)
	}

	for _, e := range pg.click.Events(gtx) {
		switch e.Type {
		case gesture.TypeClick:

			// calculate the click position making reference to list length
			pos := (e.Position.Y / float32(pg.container.Position.Length)) * 100

			switch {
			case pg.container.Position.Count > 1 && pos > -20 && pos < 1:
				setUnitValue()
			case pg.container.Position.Count > 1 && pos > 10:
				pg.mt = unit.Dp(-80) // set watchonly wallet top padding
			case pg.container.Position.Count == 1 && pos > 58 && !pg.hasWatchOnly:
				setUnitValue()
			case pg.container.Position.Count == 1 && pos > 30 && pg.hasWatchOnly:
				setUnitValue()
			default:
				pg.mt = values.MarginPadding30
			}
		}
	}
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
// Layout lays out the widgets for the main wallets pg.
func (pg *WalletPage) Layout(gtx layout.Context) layout.Dimensions {
	pg.moreOptionPositionEvent(gtx)
	pageContent := []func(gtx C) D{
		pg.walletSection,
	}

	if pg.hasWatchOnly {
		pageContent = append(pageContent, pg.watchOnlyWalletSection)
	}

	if len(pg.badWalletsList) != 0 {
		pageContent = append(pageContent, pg.badWalletsSection)
	}

	body := func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Flexed(0.85, func(gtx C) D {
				return pg.Theme.List(pg.container).Layout(gtx, len(pageContent), func(gtx C, i int) D {
					return layout.Inset{Right: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
						defer clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops).Pop()
						pg.click.Add(gtx.Ops)
						return pageContent[i](gtx)
					})
				})
			}),
			layout.Flexed(0.15, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
								return pg.layoutAddWalletSection(gtx)
							})
						})
					}),
				)
			}),
		)
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return components.UniformPadding(gtx, body)
		}),
		layout.Expanded(func(gtx C) D {
			if pg.isAddWalletMenuOpen || pg.openPopupIndex != -1 {
				return pg.backdrop.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					semantic.Button.Add(gtx.Ops)
					return layout.Dimensions{Size: gtx.Constraints.Min}
				})
			}
			return D{}
		}),
	)
}

func (pg *WalletPage) layoutOptionsMenu(gtx layout.Context, optionsMenuIndex int, listItem *walletListItem) {
	if pg.openPopupIndex != optionsMenuIndex {
		return
	}

	inset := layout.Inset{
		Top:  pg.mt,
		Left: unit.Dp(-120),
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
								return menu[i].button.Layout(gtx, func(gtx C) D {
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
				return listItem.collapsible.Layout(gtx, collapsibleHeader, collapsibleBody, collapsibleMore, listItem.wal.ID)
			}))

			if listItem.wal.IsAccountMixerActive() {
				children = append(children, layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: unit.Dp(-8)}.Layout(gtx, func(gtx C) D {
						pg.card.Color = pg.Theme.Color.Surface
						pg.card.Radius = decredmaterial.CornerRadius{BottomLeft: 10, BottomRight: 10}
						return pg.card.Layout(gtx, func(gtx C) D {
							return pg.checkMixerSection(gtx, listItem)
						})
					})
				}))
			}

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
	if !pg.hasWatchOnly {
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
					mp10 := values.MarginPadding10
					return layout.Inset{Right: mp10, Bottom: mp10}.Layout(gtx, pg.layoutWatchOnlyWallets)
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
		return layout.Inset{Top: m}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Left:  values.MarginPadding10,
								Right: values.MarginPadding10,
							}
							return inset.Layout(gtx, pg.watchOnlyWalletIcon.Layout24dp)
						}),
						layout.Rigid(pg.Theme.Body2(listItem.wal.Name).Layout),
						layout.Flexed(1, func(gtx C) D {
							return layout.E.Layout(gtx, func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										balanceLabel := pg.Theme.Body2(listItem.totalBalance)
										return layout.Inset{Right: values.MarginPadding5}.Layout(gtx, balanceLabel.Layout)
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

func (pg *WalletPage) badWalletsSection(gtx layout.Context) layout.Dimensions {
	m20 := values.MarginPadding20
	m10 := values.MarginPadding10

	layoutBadWallet := func(gtx C, badWallet *badWalletListItem, lastItem bool) D {
		return layout.Inset{Top: m10, Bottom: m10}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(pg.Theme.Body2(badWallet.Name).Layout),
						layout.Flexed(1, func(gtx C) D {
							return layout.E.Layout(gtx, func(gtx C) D {
								return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, badWallet.deleteBtn.Layout)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					if lastItem {
						return D{}
					}
					return layout.Inset{Top: values.MarginPadding10, Left: values.MarginPadding38, Right: values.MarginPaddingMinus10}.Layout(gtx, func(gtx C) D {
						return pg.Theme.Separator().Layout(gtx)
					})
				}),
			)
		})
	}

	card := pg.card
	card.Color = pg.Theme.Color.Surface
	card.Radius = decredmaterial.Radius(10)

	sectionTitleLabel := pg.Theme.Body1("Bad Wallets") // TODO: localize string
	sectionTitleLabel.Color = pg.Theme.Color.GrayText2

	return card.Layout(gtx, func(gtx C) D {
		return layout.Inset{Top: m20, Left: m20}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(sectionTitleLabel.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: m10, Bottom: m10}.Layout(gtx, pg.Theme.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return pg.Theme.NewClickableList(layout.Vertical).Layout(gtx, len(pg.badWalletsList), func(gtx C, i int) D {
							return layoutBadWallet(gtx, pg.badWalletsList[i], i == len(pg.badWalletsList)-1)
						})
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
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							var txt decredmaterial.Label
							if len(listItem.wal.EncryptedSeed) > 0 {
								txt = pg.Theme.Caption(values.String(values.StrNotBackedUp))
								txt.Color = pg.Theme.Color.Danger
								return txt.Layout(gtx)
							}
							return D{}
						}),
						layout.Rigid(func(gtx C) D {
							if listItem.wal.IsAccountMixerActive() {
								return layout.Inset{
									Left: values.MarginPadding4,
								}.Layout(gtx, func(gtx C) D {
									return decredmaterial.Card{
										Color: pg.Theme.Color.Gray4,
									}.Layout(gtx, func(gtx C) D {
										return layout.Inset{
											Left:  values.MarginPadding4,
											Right: values.MarginPadding4,
										}.Layout(gtx, func(gtx C) D {
											name := pg.Theme.Label(values.TextSize12, "Mixing...")
											name.Color = pg.Theme.Color.GrayText2
											return name.Layout(gtx)
										})
									})
								})
							}
							return D{}
						}),
					)
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
	accountIcon := pg.Theme.Icons.AccountIcon
	if account.Number == load.MaxInt32 {
		accountIcon = pg.Theme.Icons.ImportedAccountIcon
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
							pg.backupAcctIcon.Color = pg.Theme.Color.White
							return pg.backupAcctIcon.Layout(gtx, values.MarginPadding20)
						})
					})
				}),
			)
		})
	})
}

func (pg *WalletPage) checkMixerSection(gtx layout.Context, listItem *walletListItem) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return listItem.checkMixerClickable.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding8).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Top:   values.MarginPaddingMinus8,
						Left:  values.MarginPadding36,
						Right: values.MarginPadding10,
					}.Layout(gtx, pg.Theme.Separator().Layout)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							inset := layout.Inset{
								Top: values.MarginPadding5,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return layout.E.Layout(gtx, func(gtx C) D {
									txt := pg.Theme.Body2(values.String(values.StrCheckMixerStatus))
									txt.Color = pg.Theme.Color.Primary

									return layout.Flex{}.Layout(gtx,
										layout.Rigid(txt.Layout),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{
												Top:  values.MarginPadding2,
												Left: values.MarginPadding8,
											}.Layout(gtx, func(gtx C) D {
												return pg.nextIcon.Layout(gtx, values.MarginPadding16)
											})
										}),
									)
								})
							})
						}),
					)
				}),
			)
		})
	})
}

func (pg *WalletPage) layoutAddWalletMenu(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.X = gtx.Px(values.MarginPadding110)
	inset := layout.Inset{
		Top:  unit.Dp(-130),
		Left: unit.Dp(-80),
	}

	return inset.Layout(gtx, func(gtx C) D {
		return pg.Theme.Shadow().Layout(gtx, func(gtx C) D {
			return pg.optionsMenuCard.Layout(gtx, func(gtx C) D {
				return (&layout.List{Axis: layout.Vertical}).Layout(gtx, len(pg.addWalletMenu), func(gtx C, i int) D {
					return pg.addWalletMenu[i].button.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return pg.Theme.Body2(pg.addWalletMenu[i].text).Layout(gtx)
						})
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
				return decredmaterial.LinearLayout{
					Width:      decredmaterial.WrapContent,
					Height:     decredmaterial.WrapContent,
					Padding:    layout.UniformInset(values.MarginPadding12),
					Background: pg.Theme.Color.Surface,
					Clickable:  pg.openAddWalletPopupButton,
					Shadow:     pg.shadowBox,
					Border:     decredmaterial.Border{Radius: pg.openAddWalletPopupButton.Radius},
				}.Layout(gtx,
					layout.Rigid(pg.Theme.Icons.NewWalletIcon.Layout24dp),
					layout.Rigid(func(gtx C) D {
						return layout.Inset{
							Left: values.MarginPadding4,
							Top:  values.MarginPadding2,
						}.Layout(gtx, pg.Theme.Body2(values.String(values.StrAddWallet)).Layout)
					}),
				)
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

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *WalletPage) HandleUserInteractions() {
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
		if ok, selectedItem := listItem.accountsList.ItemClicked(); ok {
			pg.ChangeFragment(NewAcctDetailsPage(pg.Load, listItem.accounts[selectedItem]))
		}

		if listItem.wal.IsWatchingOnlyWallet() {
			for listItem.moreButton.Button.Clicked() {
				if pg.openPopupIndex != -1 {
					pg.closePopups()
				} else {
					pg.openPopup(index)
				}
			}
		} else {
			for listItem.collapsible.MoreTriggered() {
				pg.isAddWalletMenuOpen = false
				pg.openPopup(index)
			}

			for listItem.addAcctClickable.Clicked() {
				walletID := listItem.wal.ID
				textModal := modal.NewTextInputModal(pg.Load).
					Hint(values.String(values.StrAcctName)).
					ShowAccountInfoTip(true).
					PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
					PositiveButton(values.String(values.StrCreate), func(accountName string, tim *modal.TextInputModal) bool {
						if accountName != "" {
							modal.NewPasswordModal(pg.Load).
								Title(values.String(values.StrCreateNewAccount)).
								Hint(values.String(values.StrSpendingPassword)).
								NegativeButton(values.String(values.StrCancel), func() {}).
								PositiveButton(values.String(values.StrConfirm), func(password string, pm *modal.PasswordModal) bool {
									go func() {
										wal := pg.multiWallet.WalletWithID(walletID)
										_, err := wal.CreateNewAccount(accountName, []byte(password)) // TODO
										if err != nil {
											pg.Toast.NotifyError(err.Error())
											tim.SetError(err.Error())
										} else {
											pg.Toast.Notify(values.String(values.StrAcctCreated))
											tim.Dismiss()
										}
										pg.updateAccountBalance()
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

			for listItem.checkMixerClickable.Clicked() {
				pg.ChangeFragment(privacy.NewAccountMixerPage(pg.Load, listItem.wal))
			}
		}

		for _, menu := range listItem.optionsMenu {
			if menu.button.Clicked() {
				switch menu.id {
				case SignMessagePageID:
					pg.ChangeFragment(NewSignMessagePage(pg.Load, listItem.wal))
				case privacy.SetupPrivacyPageID:
					pg.ChangeFragment(privacy.NewSetupPrivacyPage(pg.Load, listItem.wal))
				case privacy.AccountMixerPageID:
					pg.ChangeFragment(privacy.NewAccountMixerPage(pg.Load, listItem.wal))
				case WalletSettingsPageID:
					pg.ChangeFragment(NewWalletSettingsPage(pg.Load, listItem.wal))
				default:
					menu.action(pg.Load)
				}

				pg.openPopupIndex = -1
			}
		}
	}

	for _, badWallet := range pg.badWalletsList {
		if badWallet.deleteBtn.Clicked() {
			pg.deleteBadWallet(badWallet.ID)
		}
	}

	for index := range pg.addWalletMenu {
		for pg.addWalletMenu[index].button.Clicked() {
			pg.isAddWalletMenuOpen = false
			pg.addWalletMenu[index].action(pg.Load)
		}
	}

	for pg.openAddWalletPopupButton.Clicked() {
		if pg.openPopupIndex != -1 {
			pg.closePopups()
		}
		pg.isAddWalletMenuOpen = !pg.isAddWalletMenuOpen
	}
}

func (pg *WalletPage) deleteBadWallet(badWalletID int) {
	modal.NewInfoModal(pg.Load).
		Title(values.String(values.StrRemoveWallet)).
		Body(values.String(values.StrWalletRestoreMsg)).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButtonStyle(pg.Load.Theme.Color.Surface, pg.Load.Theme.Color.Danger).
		PositiveButton(values.String(values.StrRemove), func(isChecked bool) {
			go func() {
				err := pg.WL.MultiWallet.DeleteBadWallet(badWalletID)
				if err != nil {
					pg.Toast.NotifyError(err.Error())
					return
				}
				pg.Toast.Notify(values.String(values.StrWalletRemoved))
				pg.loadBadWallets() // refresh bad wallets list
				pg.RefreshWindow()
			}()
		}).Show()
}

func (pg *WalletPage) listenForTxNotifications() {
	if pg.TxAndBlockNotificationListener != nil {
		return
	}
	pg.TxAndBlockNotificationListener = listeners.NewTxAndBlockNotificationListener()
	err := pg.WL.MultiWallet.AddTxAndBlockNotificationListener(pg.TxAndBlockNotificationListener, true, WalletPageID)
	if err != nil {
		log.Errorf("Error adding tx and block notification listener: %v", err)
		return
	}

	go func() {
		for {
			select {
			case n := <-pg.TxAndBlockNotifChan:
				switch n.Type {
				case listeners.BlockAttached:
					// refresh wallet account and balance on every new block
					// only if sync is completed.
					if pg.WL.MultiWallet.IsSynced() {
						pg.updateAccountBalance()
						pg.RefreshWindow()
					}
				case listeners.NewTransaction:
					// refresh wallets when new transaction is received
					pg.updateAccountBalance()
					pg.RefreshWindow()
				}
			case <-pg.ctx.Done():
				pg.WL.MultiWallet.RemoveTxAndBlockNotificationListener(WalletPageID)
				close(pg.TxAndBlockNotifChan)
				pg.TxAndBlockNotificationListener = nil

				return
			}
		}
	}()
}

func (pg *WalletPage) updateAccountBalance() {
	pg.listLock.Lock()
	defer pg.listLock.Unlock()

	for _, item := range pg.listItems {
		wal := pg.WL.MultiWallet.WalletWithID(item.wal.ID)
		if wal != nil {
			accountsResult, err := wal.GetAccountsRaw()
			if err != nil {
				continue
			}

			var totalBalance int64
			for _, acc := range accountsResult.Acc {
				totalBalance += acc.TotalBalance
			}

			item.totalBalance = dcrutil.Amount(totalBalance).String()
			item.accounts = accountsResult.Acc
		}
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *WalletPage) OnNavigatedFrom() {
	pg.ctxCancel()
	pg.closePopups()
}

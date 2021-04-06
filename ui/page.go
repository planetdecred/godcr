package ui

import (
	"image"
	"image/color"
	"strings"
	"time"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type pageIcons struct {
	contentAdd, navigationCheck, actionCheckCircle, actionInfo, navigationArrowBack,
	navigationArrowForward, actionCheck, chevronRight, navigationCancel, navMoreIcon,
	imageBrightness1, contentClear, dropDownIcon, cached *widget.Icon

	overviewIcon, overviewIconInactive, walletIconInactive, receiveIcon,
	transactionIcon, transactionIconInactive, sendIcon, moreIcon, moreIconInactive,
	pendingIcon, logo, redirectIcon, confirmIcon, newWalletIcon, walletAlertIcon,
	importedAccountIcon, accountIcon, editIcon, expandIcon, copyIcon, mixer,
	arrowForwardIcon, transactionFingerPrintIcon, settingsIcon, securityIcon, helpIcon,
	aboutIcon, debugIcon, verifyMessageIcon, locationPinIcon, alertGray, arrowDownIcon,
	watchOnlyWalletIcon, currencySwapIcon, syncingIcon, documentationIcon *widget.Image

	walletIcon image.Image
}

type navHandler struct {
	clickable     *widget.Clickable
	image         *widget.Image
	imageInactive *widget.Image
	page          string
}

type pageCommon struct {
	wallet          *wallet.Wallet
	info            *wallet.MultiWalletInfo
	selectedWallet  *int
	selectedAccount *int
	theme           *decredmaterial.Theme
	icons           pageIcons
	page            *string
	returnPage      *string
	navTab          *decredmaterial.Tabs
	walletTabs      *decredmaterial.Tabs
	accountTabs     *decredmaterial.Tabs
	keyEvents       chan *key.Event
	clipboard       chan interface{}
	toast           chan *toast
	toastLoad       *toast
	states          *states
	modal           *decredmaterial.Modal
	modalReceiver   chan *modalLoad
	modalLoad       *modalLoad
	modalTemplate   *ModalTemplate

	appBarNavItems          []navHandler
	drawerNavItems          []navHandler
	isNavDrawerMinimized    *bool
	minimizeNavDrawerButton decredmaterial.IconButton
	maximizeNavDrawerButton decredmaterial.IconButton
	testButton              decredmaterial.Button

	selectedUTXO map[int]map[int32]map[string]*wallet.UnspentOutput

	subPageBackButton decredmaterial.IconButton
	subPageInfoButton decredmaterial.IconButton

	changePage    func(string)
	setReturnPage func(string)
	refreshWindow func()
}

type (
	C = layout.Context
	D = layout.Dimensions
)

const (
	navDrawerWidth          = 160
	navDrawerMinimizedWidth = 80
)

func (win *Window) addPages(decredIcons map[string]image.Image) {
	ic := pageIcons{
		contentAdd:             mustIcon(widget.NewIcon(icons.ContentAdd)),
		navigationCheck:        mustIcon(widget.NewIcon(icons.NavigationCheck)),
		actionCheckCircle:      mustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		navigationArrowBack:    mustIcon(widget.NewIcon(icons.NavigationArrowBack)),
		navigationArrowForward: mustIcon(widget.NewIcon(icons.NavigationArrowForward)),
		actionInfo:             mustIcon(widget.NewIcon(icons.ActionInfo)),
		actionCheck:            mustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		navigationCancel:       mustIcon(widget.NewIcon(icons.NavigationCancel)),
		imageBrightness1:       mustIcon(widget.NewIcon(icons.ImageBrightness1)),
		chevronRight:           mustIcon(widget.NewIcon(icons.NavigationChevronRight)),
		contentClear:           mustIcon(widget.NewIcon(icons.ContentClear)),
		navMoreIcon:            mustIcon(widget.NewIcon(icons.NavigationMoreHoriz)),
		dropDownIcon:           mustIcon(widget.NewIcon(icons.NavigationArrowDropDown)),
		cached:                 mustIcon(widget.NewIcon(icons.ActionCached)),

		overviewIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["overview"])},
		overviewIconInactive:       &widget.Image{Src: paint.NewImageOp(decredIcons["overview_inactive"])},
		walletIconInactive:         &widget.Image{Src: paint.NewImageOp(decredIcons["wallet_inactive"])},
		receiveIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["receive"])},
		transactionIcon:            &widget.Image{Src: paint.NewImageOp(decredIcons["transaction"])},
		transactionIconInactive:    &widget.Image{Src: paint.NewImageOp(decredIcons["transaction_inactive"])},
		sendIcon:                   &widget.Image{Src: paint.NewImageOp(decredIcons["send"])},
		moreIcon:                   &widget.Image{Src: paint.NewImageOp(decredIcons["more"])},
		moreIconInactive:           &widget.Image{Src: paint.NewImageOp(decredIcons["more_inactive"])},
		logo:                       &widget.Image{Src: paint.NewImageOp(decredIcons["logo"])},
		confirmIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["confirmed"])},
		pendingIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["pending"])},
		redirectIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["redirect"])},
		newWalletIcon:              &widget.Image{Src: paint.NewImageOp(decredIcons["addNewWallet"])},
		walletAlertIcon:            &widget.Image{Src: paint.NewImageOp(decredIcons["walletAlert"])},
		accountIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["account"])},
		importedAccountIcon:        &widget.Image{Src: paint.NewImageOp(decredIcons["imported_account"])},
		editIcon:                   &widget.Image{Src: paint.NewImageOp(decredIcons["editIcon"])},
		expandIcon:                 &widget.Image{Src: paint.NewImageOp(decredIcons["expand_icon"])},
		copyIcon:                   &widget.Image{Src: paint.NewImageOp(decredIcons["copy_icon"])},
		mixer:                      &widget.Image{Src: paint.NewImageOp(decredIcons["mixer"])},
		transactionFingerPrintIcon: &widget.Image{Src: paint.NewImageOp(decredIcons["transaction_fingerprint"])},
		arrowForwardIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["arrow_forward"])},
		settingsIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["settings"])},
		securityIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["security"])},
		helpIcon:                   &widget.Image{Src: paint.NewImageOp(decredIcons["help_icon"])},
		aboutIcon:                  &widget.Image{Src: paint.NewImageOp(decredIcons["about_icon"])},
		debugIcon:                  &widget.Image{Src: paint.NewImageOp(decredIcons["debug"])},
		verifyMessageIcon:          &widget.Image{Src: paint.NewImageOp(decredIcons["verify_message"])},
		locationPinIcon:            &widget.Image{Src: paint.NewImageOp(decredIcons["location_pin"])},
		alertGray:                  &widget.Image{Src: paint.NewImageOp(decredIcons["alert-gray"])},
		arrowDownIcon:              &widget.Image{Src: paint.NewImageOp(decredIcons["arrow_down"])},
		watchOnlyWalletIcon:        &widget.Image{Src: paint.NewImageOp(decredIcons["watch_only_wallet"])},
		currencySwapIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["swap"])},
		syncingIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["syncing"])},
		documentationIcon:          &widget.Image{Src: paint.NewImageOp(decredIcons["documentation"])},

		walletIcon: decredIcons["wallet"],
	}
	win.theme.NavigationCheckIcon = ic.navigationCheck

	appBarNavItems := []navHandler{
		{
			clickable: new(widget.Clickable),
			image:     ic.sendIcon,
			page:      PageSend,
		},
		{
			clickable: new(widget.Clickable),
			image:     ic.receiveIcon,
			page:      PageReceive,
		},
	}

	drawerNavItems := []navHandler{
		{
			clickable:     new(widget.Clickable),
			image:         ic.overviewIcon,
			imageInactive: ic.overviewIconInactive,
			page:          PageOverview,
		},
		{
			clickable:     new(widget.Clickable),
			image:         ic.transactionIcon,
			imageInactive: ic.transactionIconInactive,
			page:          PageTransactions,
		},
		{
			clickable:     new(widget.Clickable),
			image:         &widget.Image{Src: paint.NewImageOp(ic.walletIcon)},
			imageInactive: ic.walletIconInactive,
			page:          PageWallet,
		},
		{
			clickable:     new(widget.Clickable),
			image:         &widget.Image{Src: paint.NewImageOp(ic.walletIcon)},
			imageInactive: ic.walletIconInactive,
			page:          PageTickets,
		},
		{
			clickable:     new(widget.Clickable),
			image:         ic.moreIcon,
			imageInactive: ic.moreIconInactive,
			page:          PageMore,
		},
	}

	common := pageCommon{
		wallet:                  win.wallet,
		info:                    win.walletInfo,
		selectedWallet:          &win.selected,
		selectedAccount:         &win.selectedAccount,
		theme:                   win.theme,
		icons:                   ic,
		returnPage:              &win.previous,
		page:                    &win.current,
		walletTabs:              win.walletTabs,
		accountTabs:             win.accountTabs,
		keyEvents:               win.keyEvents,
		clipboard:               win.clipboard,
		toast:                   win.toast,
		toastLoad:               &toast{},
		states:                  &win.states,
		appBarNavItems:          appBarNavItems,
		drawerNavItems:          drawerNavItems,
		minimizeNavDrawerButton: win.theme.PlainIconButton(new(widget.Clickable), ic.navigationArrowBack),
		maximizeNavDrawerButton: win.theme.PlainIconButton(new(widget.Clickable), ic.navigationArrowForward),
		selectedUTXO:            make(map[int]map[int32]map[string]*wallet.UnspentOutput),
		modal:                   win.theme.Modal(),
		modalReceiver:           win.modal,
		modalLoad:               &modalLoad{},
		subPageBackButton:       win.theme.PlainIconButton(new(widget.Clickable), ic.navigationArrowBack),
		subPageInfoButton:       win.theme.PlainIconButton(new(widget.Clickable), ic.actionInfo),
		changePage:              win.changePage,
		setReturnPage:           win.setReturnPage,
		refreshWindow:           win.refresh,
	}

	common.testButton = win.theme.Button(new(widget.Clickable), "test button")
	isNavDrawerMinimized := false
	common.isNavDrawerMinimized = &isNavDrawerMinimized

	iconColor := common.theme.Color.IconColor
	common.minimizeNavDrawerButton.Color, common.maximizeNavDrawerButton.Color = iconColor, iconColor

	zeroInset := layout.UniformInset(values.MarginPadding0)
	common.subPageBackButton.Color, common.subPageInfoButton.Color = iconColor, iconColor

	m25 := values.MarginPadding25
	common.subPageBackButton.Size, common.subPageInfoButton.Size = m25, m25
	common.subPageBackButton.Inset, common.subPageInfoButton.Inset = zeroInset, zeroInset

	common.modalTemplate = win.LoadModalTemplates()

	win.pages = make(map[string]layout.Widget)
	win.pages[PageWallet] = win.WalletPage(common)
	win.pages[PageOverview] = win.OverviewPage(common)
	win.pages[PageTransactions] = win.TransactionsPage(common)
	win.pages[PageMore] = win.MorePage(common)
	win.pages[PageCreateRestore] = win.CreateRestorePage(common)
	win.pages[PageReceive] = win.ReceivePage(common)
	win.pages[PageSend] = win.SendPage(common)
	win.pages[PageTransactionDetails] = win.TransactionDetailsPage(common)
	win.pages[PageSignMessage] = win.SignMessagePage(common)
	win.pages[PageVerifyMessage] = win.VerifyMessagePage(common)
	win.pages[PageSeedBackup] = win.BackupPage(common)
	win.pages[PageSettings] = win.SettingsPage(common)
	win.pages[PageWalletSettings] = win.WalletSettingsPage(common)
	win.pages[PageSecurityTools] = win.SecurityToolsPage(common)
	win.pages[PagePoliteia] = win.PoliteiaPage(common)
	win.pages[PageDebug] = win.DebugPage(common)
	win.pages[PageLog] = win.LogPage(common)
	win.pages[PageAbout] = win.AboutPage(common)
	win.pages[PageHelp] = win.HelpPage(common)
	win.pages[PageUTXO] = win.UTXOPage(common)
	win.pages[PageAccountDetails] = win.AcctDetailsPage(common)
	win.pages[PagePrivacy] = win.PrivacyPage(common)
	win.pages[PageTickets] = win.TicketPage(common)
	win.pages[ValidateAddress] = win.ValidateAddressPage(common)
}

func (page pageCommon) ChangePage(pg string) {
	page.changePage(pg)
}

func (page pageCommon) refreshPage() {
	page.refreshWindow()
}

func (page pageCommon) notify(text string, success bool) {
	go func() {
		page.toast <- &toast{
			text:    text,
			success: success,
		}
	}()
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
			page.changePage(page.appBarNavItems[i].page)
		}
	}

	for i := range page.drawerNavItems {
		for page.drawerNavItems[i].clickable.Clicked() {
			page.changePage(page.drawerNavItems[i].page)
		}
	}
}

func (page pageCommon) Layout(gtx layout.Context, body layout.Widget) layout.Dimensions {
	page.handleNavEvents()

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			// fill the entire window with a color if a user has no wallet created
			if *page.page == PageCreateRestore {
				return fill(gtx, page.theme.Color.Surface)
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(page.layoutAppBar),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							card := page.theme.Card()
							card.Radius = decredmaterial.CornerRadius{
								NE: 0,
								NW: 0,
								SE: 0,
								SW: 0,
							}
							return card.Layout(gtx, page.layoutNavDrawer)
						}),
						layout.Rigid(func(gtx C) D {
							return body(gtx)
						}),
					)
				}),
			)
		}),
		layout.Stacked(func(gtx C) D {
			// stack the page content on the entire window if a user has no wallet
			if *page.page == PageCreateRestore {
				return body(gtx)
			}
			return layout.Dimensions{}
		}),
		layout.Stacked(func(gtx C) D {
		outer:
			for {
				select {
				case load := <-page.modalReceiver:
					page.modalLoad.template = load.template
					page.modalLoad.title = load.title
					page.modalLoad.confirm = load.confirm
					page.modalLoad.confirmText = load.confirmText
					page.modalLoad.cancel = load.cancel
					page.modalLoad.cancelText = load.cancelText
					page.modalLoad.isReset = false
				default:
					break outer
				}
			}

			if page.modalLoad.cancel != nil {
				return page.modal.Layout(gtx, page.modalTemplate.Layout(page.theme, page.modalLoad),
					900)
			}

			return layout.Dimensions{}
		}),
		layout.Stacked(func(gtx C) D {
			toast := func(n *toast) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx C) D {
						return layout.Center.Layout(gtx, func(gtx C) D {
							return layout.Inset{Top: values.MarginPadding65}.Layout(gtx, func(gtx C) D {
								return displayToast(page.theme, gtx, n)
							})
						})
					}),
				)
			}

		outer:
			for {
				select {
				case n := <-page.toast:
					page.toastLoad.success = n.success
					page.toastLoad.text = n.text
					page.toastLoad.ResetTimer()
				default:
					break outer
				}
			}

			if page.toastLoad.text != "" {
				page.toastLoad.Timer(time.Second*3, func() {
					page.toastLoad.text = ""
				})
				return toast(page.toastLoad)
			}
			return layout.Dimensions{}
		}),
	)
}

func (page pageCommon) closeModal() {
	go func() {
		page.modalReceiver <- &modalLoad{
			title:   "",
			confirm: nil,
			cancel:  nil,
		}
	}()
}

func fill(gtx layout.Context, col color.NRGBA) layout.Dimensions {
	return decredmaterial.Fill(gtx, col)
}

func (page pageCommon) layoutAppBar(gtx layout.Context) layout.Dimensions {
	card := page.theme.Card()
	card.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
	return card.Layout(gtx, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								img := page.icons.logo
								img.Scale = 1.0
								return layout.Inset{
									Top:   values.MarginPadding15,
									Left:  values.MarginPadding25,
									Right: values.MarginPadding16,
								}.Layout(gtx, func(gtx C) D {
									return img.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								m := values.MarginPadding10
								return layout.Inset{
									Top: m,
								}.Layout(gtx, func(gtx C) D {
									return page.layoutBalance(gtx, page.info.TotalBalance)
								})
							}),
						)
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.E.Layout(gtx, func(gtx C) D {
							list := layout.List{Axis: layout.Horizontal}
							return list.Layout(gtx, len(page.appBarNavItems), func(gtx C, i int) D {
								background := page.theme.Color.Surface
								if page.appBarNavItems[i].page == *page.page {
									background = page.theme.Color.Background
								}
								card := page.theme.Card()
								card.Color = background
								card.Radius = decredmaterial.CornerRadius{
									NE: 0,
									NW: 0,
									SE: 0,
									SW: 0,
								}
								return card.Layout(gtx, func(gtx C) D {
									return decredmaterial.Clickable(gtx, page.appBarNavItems[i].clickable, func(gtx C) D {
										return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
											return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
												layout.Rigid(func(gtx C) D {
													page.appBarNavItems[i].image.Scale = 0.05
													return layout.Center.Layout(gtx, func(gtx C) D {
														img := page.appBarNavItems[i].image
														img.Scale = 1.0
														return page.appBarNavItems[i].image.Layout(gtx)
													})
												}),
												layout.Rigid(func(gtx C) D {
													return layout.Inset{
														Left: values.MarginPadding10,
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
			layout.Rigid(func(gtx C) D {
				l := page.theme.Line()
				l.Color = page.theme.Color.Background
				l.Width = gtx.Constraints.Min.X
				l.Height = 2
				return l.Layout(gtx)
			}),
		)
	})
}

func (page pageCommon) layoutNavDrawer(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			list := layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(page.drawerNavItems), func(gtx C, i int) D {
				background := page.theme.Color.Surface
				if page.drawerNavItems[i].page == *page.page {
					background = page.theme.Color.Background
				}
				txt := page.theme.Label(values.TextSize16, page.drawerNavItems[i].page)
				return decredmaterial.Clickable(gtx, page.drawerNavItems[i].clickable, func(gtx C) D {
					card := page.theme.Card()
					card.Color = background
					card.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
					return card.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
							axis := layout.Horizontal
							leftInset := float32(15)
							width := navDrawerWidth
							if *page.isNavDrawerMinimized {
								axis = layout.Vertical
								txt.TextSize = values.TextSize10
								leftInset = 0
								width = navDrawerMinimizedWidth
							}

							gtx.Constraints.Min.X = int(gtx.Metric.PxPerDp) * width
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
										Left: unit.Dp(leftInset),
										Top:  unit.Dp(4),
									}.Layout(gtx, func(gtx C) D {
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

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func (page pageCommon) layoutBalance(gtx layout.Context, amount string) layout.Dimensions {
	mainText, subText := page.breakBalance(amount)
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			txt := page.theme.H5(mainText)
			txt.Color = page.theme.Color.DeepBlue
			return txt.Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			txt := page.theme.Body1(subText)
			txt.Color = page.theme.Color.DeepBlue
			return txt.Layout(gtx)
		}),
	)
}

// breakBalance takes the balance string and returns it in two slices
func (page pageCommon) breakBalance(balance string) (b1, b2 string) {
	balanceParts := strings.Split(balance, ".")
	if len(balanceParts) == 1 {
		return balanceParts[0], ""
	}
	b1 = balanceParts[0]
	b2 = balanceParts[1]
	b1 = b1 + "." + b2[:2]
	b2 = b2[2:]
	return
}

func (page pageCommon) Modal(gtx layout.Context, body layout.Dimensions, modal layout.Dimensions) layout.Dimensions {
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

func (page pageCommon) LayoutWithWallets(gtx layout.Context, body layout.Widget) layout.Dimensions {
	bd := func(gtx C) D {
		if page.walletTabs.ChangeEvent() {
			*page.selectedWallet = page.walletTabs.Selected
			*page.selectedAccount = 0
			page.accountTabs.Selected = 0

			accounts := make([]decredmaterial.TabItem, len(page.info.Wallets[*page.selectedWallet].Accounts))
			for i, account := range page.info.Wallets[*page.selectedWallet].Accounts {
				if account.Name == "imported" {
					continue
				}
				accounts[i] = decredmaterial.TabItem{
					Title: page.info.Wallets[*page.selectedWallet].Accounts[i].Name,
				}
			}
			page.accountTabs.SetTabs(accounts)
		}
		return page.walletTabs.Layout(gtx, body)
	}
	return page.Layout(gtx, bd)
}

func (page pageCommon) LayoutWithAccounts(gtx layout.Context, body layout.Widget) layout.Dimensions {
	if page.accountTabs.ChangeEvent() {
		*page.selectedAccount = page.accountTabs.Selected
	}

	if page.selectedUTXO[page.info.Wallets[*page.selectedWallet].ID] == nil {
		current := page.info.Wallets[*page.selectedWallet]
		account := page.info.Wallets[*page.selectedWallet].Accounts[*page.selectedAccount]
		page.selectedUTXO[current.ID] = make(map[int32]map[string]*wallet.UnspentOutput)
		page.selectedUTXO[current.ID][account.Number] = make(map[string]*wallet.UnspentOutput)
	}

	return page.LayoutWithWallets(gtx, func(gtx C) D {
		return page.accountTabs.Layout(gtx, body)
	})
}

func (page pageCommon) SelectedAccountLayout(gtx layout.Context) layout.Dimensions {
	current := page.info.Wallets[*page.selectedWallet]
	account := page.info.Wallets[*page.selectedWallet].Accounts[*page.selectedAccount]

	selectedDetails := func(gtx C) D {
		return layout.UniformInset(values.MarginPadding10).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
								return page.theme.H6(account.Name).Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
								return page.theme.H6(dcrutil.Amount(account.SpendableBalance).String()).Layout(gtx)
							})
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
								return page.theme.Body2(current.Name).Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
								return page.theme.Body2(current.Balance).Layout(gtx)
							})
						}),
					)
				}),
			)
		})
	}

	card := page.theme.Card()
	card.Radius = decredmaterial.CornerRadius{
		NE: 0,
		NW: 0,
		SE: 0,
		SW: 0,
	}
	return card.Layout(gtx, selectedDetails)
}

func (page pageCommon) UniformPadding(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding24).Layout(gtx, body)
}

type SubPage struct {
	title        string
	subTitle     string
	walletName   string
	back         func()
	body         layout.Widget
	infoTemplate string
	extraItem    *widget.Clickable
	extra        layout.Widget
	handleExtra  func()
}

func (page pageCommon) SubPageLayout(gtx layout.Context, sp SubPage) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: values.MarginPadding15, Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
				return page.subpageHeader(gtx, sp)
			})
		}),
		layout.Rigid(sp.body),
	)
}

func (page pageCommon) subpageHeader(gtx layout.Context, sp SubPage) layout.Dimensions {
	page.subpageEventHandler(sp)

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Right: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
				return page.subPageBackButton.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if sp.subTitle == "" {
				return page.theme.H6(sp.title).Layout(gtx)
			}
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return page.theme.H6(sp.title).Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return page.theme.Body1(sp.subTitle).Layout(gtx)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if sp.walletName != "" {
				return layout.Inset{Left: values.MarginPadding5, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return decredmaterial.Card{
						Color: page.theme.Color.Surface,
					}.Layout(gtx, func(gtx C) D {
						return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
							walletText := page.theme.Caption(sp.walletName)
							walletText.Color = page.theme.Color.Gray
							return walletText.Layout(gtx)
						})
					})
				})
			}
			return layout.Dimensions{}
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Right: values.MarginPadding9}.Layout(gtx, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx C) D {
					if sp.infoTemplate != "" {
						return page.subPageInfoButton.Layout(gtx)
					} else if sp.extraItem != nil {
						return decredmaterial.Clickable(gtx, sp.extraItem, sp.extra)
					}
					return layout.Dimensions{}
				})
			})
		}),
	)
}

func (page pageCommon) SubpageSplitLayout(gtx layout.Context, sp SubPage) layout.Dimensions {
	card := page.theme.Card()
	card.Color = color.NRGBA{}
	return card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D { return page.subpageHeader(gtx, sp) }),
			layout.Rigid(sp.body),
		)
	})
}

func (page pageCommon) subpageEventHandler(sp SubPage) {
	if page.subPageInfoButton.Button.Clicked() {
		go func() {
			page.modalReceiver <- &modalLoad{
				template:   sp.infoTemplate,
				title:      sp.title,
				cancel:     page.closeModal,
				cancelText: "Got it",
			}
		}()
	}

	if page.subPageBackButton.Button.Clicked() {
		sp.back()
	}

	if sp.extraItem != nil && sp.extraItem.Clicked() {
		sp.handleExtra()
	}
}

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

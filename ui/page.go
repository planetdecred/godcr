package ui

import (
	"image"
	"image/color"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type pageIcons struct {
	contentAdd, navigationCheck, navigationMore, actionCheckCircle, actionInfo, navigationArrowBack,
	navigationArrowForward, actionCheck, chevronRight, navigationCancel, navMoreIcon,
	imageBrightness1, contentClear, dropDownIcon, cached *widget.Icon

	overviewIcon, overviewIconInactive, walletIconInactive, receiveIcon,
	transactionIcon, transactionIconInactive, sendIcon, moreIcon, moreIconInactive,
	pendingIcon, logo, redirectIcon, confirmIcon, newWalletIcon, walletAlertIcon,
	importedAccountIcon, accountIcon, editIcon, expandIcon, copyIcon, mixer,
	arrowForwardIcon, transactionFingerPrintIcon, settingsIcon, securityIcon, helpIcon,
	aboutIcon, debugIcon, verifyMessageIcon, locationPinIcon, alertGray, arrowDownIcon,
	watchOnlyWalletIcon, currencySwapIcon, syncingIcon *widget.Image

	walletIcon image.Image
}

type navHandler struct {
	clickable     *widget.Clickable
	image         *widget.Image
	imageInactive *widget.Image
	page          string
}

type pageCommon struct {
	printer         *message.Printer
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
		navigationMore:         mustIcon(widget.NewIcon(icons.NavigationMoreHoriz)),
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

		walletIcon: decredIcons["wallet"],
	}

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
		printer:                 message.NewPrinter(language.English),
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

	iconColor := common.theme.Color.Gray3
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

func (page pageCommon) closeModal() {
	go func() {
		page.modalReceiver <- &modalLoad{
			title:   "",
			confirm: nil,
			cancel:  nil,
		}
	}()
}

func (page pageCommon) Layout(gtx layout.Context, body layout.Widget) layout.Dimensions {
	page.handleNavEvents()

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			// fill the entire window with a color if a user has no wallet created
			if *page.page == PageCreateRestore {
				return decredmaterial.Fill(gtx, page.theme.Color.Surface)
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(page.layoutTopBar),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							card := page.theme.Card()
							card.Radius = decredmaterial.CornerRadius{}
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
			// global modal. Stack modal on all pages and contents
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
			// global toasts. Stack toast on all pages and contents
			t := func(n *toast) layout.Dimensions {
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
				return t(page.toastLoad)
			}
			return layout.Dimensions{}
		}),
	)
}

// Container is simply a wrapper for the Inset type. Its purpose is to differentiate the use of an inset as a padding or
// margin, making it easier to visualize the structure of a layout when reading UI code.
type Container struct {
	padding layout.Inset
}

func (c Container) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return c.padding.Layout(gtx, w)
}

func (page pageCommon) UniformPadding(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding24).Layout(gtx, body)
}

type SubPage struct {
	title        string
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
			return page.theme.H6(sp.title).Layout(gtx)
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

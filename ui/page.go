package ui

import (
	"encoding/json"
	"image"
	"image/color"
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type pageIcons struct {
	contentAdd, navigationCheck, navigationMore, actionCheckCircle, actionInfo, navigationArrowBack,
	navigationArrowForward, actionCheck, chevronRight, navigationCancel, navMoreIcon,
	imageBrightness1, contentClear, dropDownIcon, cached *widget.Icon

	overviewIcon, overviewIconInactive, walletIcon, walletIconInactive,
	receiveIcon, transactionIcon, transactionIconInactive, sendIcon, moreIcon, moreIconInactive,
	pendingIcon, logo, redirectIcon, confirmIcon, newWalletIcon, walletAlertIcon,
	importedAccountIcon, accountIcon, editIcon, expandIcon, copyIcon, mixer, mixerSmall,
	arrowForwardIcon, transactionFingerPrintIcon, settingsIcon, securityIcon, helpIcon,
	aboutIcon, debugIcon, verifyMessageIcon, locationPinIcon, alertGray, arrowDownIcon,
	watchOnlyWalletIcon, currencySwapIcon, syncingIcon, proposalIconActive, proposalIconInactive,
	restore, documentationIcon, downloadIcon, timerIcon, ticketIcon, ticketIconInactive, stakeyIcon,
	list, listGridIcon *widget.Image

	ticketPurchasedIcon,
	ticketImmatureIcon,
	ticketLiveIcon,
	ticketVotedIcon,
	ticketMissedIcon,
	ticketExpiredIcon,
	ticketRevokedIcon,
	ticketUnminedIcon *widget.Image
}

type Page interface {
	pageID() string
	Layout(layout.Context) layout.Dimensions
	handle()
	onClose()
}

type navHandler struct {
	clickable     *widget.Clickable
	image         *widget.Image
	imageInactive *widget.Image
	pageID        string
	label         string
}

type walletAccount struct {
	evt          *gesture.Click
	walletIndex  int
	accountIndex int
	accountName  string
	totalBalance string
	spendable    string
	number       int32
}

type wallectAccountOption struct {
	selectSendAccount           map[int][]walletAccount
	selectReceiveAccount        map[int][]walletAccount
	selectPurchaseTicketAccount map[int][]walletAccount
}

type DCRUSDTBittrex struct {
	LastTradeRate string
}

type pageCommon struct {
	printer          *message.Printer
	multiWallet      *dcrlibwallet.MultiWallet
	syncStatusUpdate *chan wallet.SyncStatusUpdate
	wallet           *wallet.Wallet
	selectedWallet   *int
	selectedAccount  *int
	theme            *decredmaterial.Theme
	icons            pageIcons
	returnPage       Page
	amountDCRtoUSD   float64
	usdExchangeRate  float64
	usdExchangeSet   bool
	dcrUsdtBittrex   DCRUSDTBittrex
	navTab           *decredmaterial.Tabs
	walletTabs       *decredmaterial.Tabs
	accountTabs      *decredmaterial.Tabs
	keyEvents        chan *key.Event
	toast            **toast
	states           *states
	modal            *decredmaterial.Modal
	modalReceiver    chan *modalLoad
	modalLoad        *modalLoad
	modalTemplate    *ModalTemplate

	testButton decredmaterial.Button

	selectedUTXO map[int]map[int32]map[string]*wallet.UnspentOutput

	subPageBackButton decredmaterial.IconButton
	subPageInfoButton decredmaterial.IconButton

	popPage    func()
	popToPage  func(string) error
	changePage func(Page)

	showModal    func(Modal)
	dismissModal func(Modal)

	setReturnPage func(Page) //TODO
	refreshWindow func()
}

type (
	C = layout.Context
	D = layout.Dimensions
)

func (win *Window) loadIcons(decredIcons map[string]image.Image) pageIcons {
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
		mixerSmall:                 &widget.Image{Src: paint.NewImageOp(decredIcons["mixer_small"])},
		transactionFingerPrintIcon: &widget.Image{Src: paint.NewImageOp(decredIcons["transaction_fingerprint"])},
		arrowForwardIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["arrow_forward"])},
		settingsIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["settings"])},
		securityIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["security"])},
		helpIcon:                   &widget.Image{Src: paint.NewImageOp(decredIcons["help_icon"])},
		aboutIcon:                  &widget.Image{Src: paint.NewImageOp(decredIcons["about_icon"])},
		debugIcon:                  &widget.Image{Src: paint.NewImageOp(decredIcons["debug"])},
		verifyMessageIcon:          &widget.Image{Src: paint.NewImageOp(decredIcons["verify_message"])},
		locationPinIcon:            &widget.Image{Src: paint.NewImageOp(decredIcons["location_pin"])},
		alertGray:                  &widget.Image{Src: paint.NewImageOp(decredIcons["alert_gray"])},
		arrowDownIcon:              &widget.Image{Src: paint.NewImageOp(decredIcons["arrow_down"])},
		watchOnlyWalletIcon:        &widget.Image{Src: paint.NewImageOp(decredIcons["watch_only_wallet"])},
		currencySwapIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["swap"])},
		syncingIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["syncing"])},
		documentationIcon:          &widget.Image{Src: paint.NewImageOp(decredIcons["documentation"])},
		proposalIconActive:         &widget.Image{Src: paint.NewImageOp(decredIcons["politeiaActive"])},
		proposalIconInactive:       &widget.Image{Src: paint.NewImageOp(decredIcons["politeiaInactive"])},
		restore:                    &widget.Image{Src: paint.NewImageOp(decredIcons["restore"])},
		downloadIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["downloadIcon"])},
		timerIcon:                  &widget.Image{Src: paint.NewImageOp(decredIcons["timerIcon"])},
		walletIcon:                 &widget.Image{Src: paint.NewImageOp(decredIcons["wallet"])},
		ticketIcon:                 &widget.Image{Src: paint.NewImageOp(decredIcons["ticket"])},
		ticketIconInactive:         &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_inactive"])},
		stakeyIcon:                 &widget.Image{Src: paint.NewImageOp(decredIcons["stakey"])},
		ticketPurchasedIcon:        &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_purchased"])},
		ticketImmatureIcon:         &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_immature"])},
		ticketUnminedIcon:          &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_unmined"])},
		ticketLiveIcon:             &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_live"])},
		ticketVotedIcon:            &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_voted"])},
		ticketMissedIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_missed"])},
		ticketExpiredIcon:          &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_expired"])},
		ticketRevokedIcon:          &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_revoked"])},
		list:                       &widget.Image{Src: paint.NewImageOp(decredIcons["list"])},
		listGridIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["list_grid"])},
	}

	return ic
}

func (win *Window) loadPageCommon(decredIcons map[string]image.Image, multiWallet *dcrlibwallet.MultiWallet) pageCommon {

	ic := win.loadIcons(decredIcons)

	common := pageCommon{
		printer:           message.NewPrinter(language.English),
		multiWallet:       multiWallet,
		syncStatusUpdate:  &win.syncStatusUpdate,
		wallet:            win.wallet,
		selectedWallet:    &win.selected,
		selectedAccount:   &win.selectedAccount,
		theme:             win.theme,
		icons:             ic,
		returnPage:        win.previousPage,
		walletTabs:        win.walletTabs,
		accountTabs:       win.accountTabs,
		keyEvents:         win.keyEvents,
		states:            &win.states,
		selectedUTXO:      make(map[int]map[int32]map[string]*wallet.UnspentOutput),
		modal:             win.theme.Modal(),
		modalReceiver:     win.modal,
		modalLoad:         &modalLoad{},
		subPageBackButton: win.theme.PlainIconButton(new(widget.Clickable), ic.navigationArrowBack),
		subPageInfoButton: win.theme.PlainIconButton(new(widget.Clickable), ic.actionInfo),
		popPage:           win.popPage,
		popToPage:         win.popToPage,
		changePage:        win.changePage,
		setReturnPage:     win.setReturnPage,
		refreshWindow:     win.refresh,
		toast:             &win.toast,
	}

	if common.fetchExchangeValue(&common.dcrUsdtBittrex) != nil {
		return
	}

	iconColor := common.theme.Color.Gray3

	common.testButton = win.theme.Button(new(widget.Clickable), "test button")

	zeroInset := layout.UniformInset(values.MarginPadding0)
	common.subPageBackButton.Color, common.subPageInfoButton.Color = iconColor, iconColor

	m25 := values.MarginPadding25
	common.subPageBackButton.Size, common.subPageInfoButton.Size = m25, m25
	common.subPageBackButton.Inset, common.subPageInfoButton.Inset = zeroInset, zeroInset

	common.modalTemplate = win.LoadModalTemplates()

	return common
}

func (page *pageCommon) fetchExchangeValue(target interface{}) error {
	url := "https://api.bittrex.com/v3/markets/DCR-USDT/ticker"
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(target)
	if err != nil {
		log.Error(err)
	}

	return nil
}

func (page pageCommon) refreshPage() {
	page.refreshWindow()
}

func (page pageCommon) notify(text string, success bool) {
	*page.toast = &toast{
		text:    text,
		success: success,
	}
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
	page.handler()

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			// fill the entire window with a color if a user has no wallet created
			if "page.page.pageID()" == PageCreateRestore { //TODO
				return decredmaterial.Fill(gtx, page.theme.Color.Surface)
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return body(gtx)
						}),
					)
				}),
			)
		}),
		layout.Stacked(func(gtx C) D {
			// stack the page content on the entire window if a user has no wallet
			if "page.page.pageID()" == PageCreateRestore { //TODO
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
			if *page.toast == nil {
				return layout.Dimensions{}
			}
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding65}.Layout(gtx, func(gtx C) D {
					return displayToast(page.theme, gtx, *page.toast)
				})
			})
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
	subTitle     string
	walletName   string
	back         func()
	body         layout.Widget
	infoTemplate string
	extraItem    *widget.Clickable
	extra        layout.Widget
	extraText    string
	handleExtra  func()
}

func (page pageCommon) SubPageLayout(gtx layout.Context, sp SubPage) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
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
			return layout.Inset{Right: values.MarginPadding20}.Layout(gtx, page.subPageBackButton.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if sp.subTitle == "" {
				return page.theme.H6(sp.title).Layout(gtx)
			}
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(page.theme.H6(sp.title).Layout),
				layout.Rigid(page.theme.Body1(sp.subTitle).Layout),
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
			return layout.E.Layout(gtx, func(gtx C) D {
				if sp.infoTemplate != "" {
					return page.subPageInfoButton.Layout(gtx)
				} else if sp.extraItem != nil {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if sp.extraText != "" {
								return layout.Inset{Right: values.MarginPadding10, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
									text := page.theme.Caption(sp.extraText)
									text.Color = page.theme.Color.DeepBlue
									return text.Layout(gtx)
								})
							}
							return layout.Dimensions{}
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return decredmaterial.Clickable(gtx, sp.extraItem, sp.extra)
						}),
					)
				}
				return layout.Dimensions{}
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

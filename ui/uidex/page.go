package uidex

import (
	"encoding/json"
	"image"
	"image/color"
	"net/http"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/utils"
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
	restore, documentationIcon, downloadIcon, timerIcon, ticketIcon, ticketIconInactive, stakeyIcon *widget.Image

	ticketPurchasedIcon,
	ticketImmatureIcon,
	ticketLiveIcon,
	ticketVotedIcon,
	ticketMissedIcon,
	ticketExpiredIcon,
	ticketRevokedIcon,
	ticketUnminedIcon *widget.Image
}

type navHandler struct {
	clickable     *widget.Clickable
	image         *widget.Image
	imageInactive *widget.Image
	page          string
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
	printer         *message.Printer
	wallet          *wallet.Wallet
	info            *wallet.MultiWalletInfo
	selectedWallet  *int
	selectedAccount *int
	theme           *decredmaterial.Theme
	icons           pageIcons
	page            *string
	returnPage      *string
	amountDCRtoUSD  float64
	usdExchangeRate float64
	usdExchangeSet  bool
	dcrUsdtBittrex  DCRUSDTBittrex
	navTab          *decredmaterial.Tabs
	walletTabs      *decredmaterial.Tabs
	accountTabs     *decredmaterial.Tabs
	keyEvents       chan *key.Event
	clipboard       chan interface{}
	toast           chan *toast
	toastLoad       *toast
	// states          *states
	modal *decredmaterial.Modal
	// modalReceiver   chan *modalLoad
	// modalLoad       *modalLoad
	// modalTemplate   *ModalTemplate

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
	switchView    *int
}

type (
	C = layout.Context
	D = layout.Dimensions
)

func (d *Dex) addPages(decredIcons map[string]image.Image) {
	ic := pageIcons{
		contentAdd:             utils.MustIcon(widget.NewIcon(icons.ContentAdd)),
		navigationCheck:        utils.MustIcon(widget.NewIcon(icons.NavigationCheck)),
		navigationMore:         utils.MustIcon(widget.NewIcon(icons.NavigationMoreHoriz)),
		actionCheckCircle:      utils.MustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		navigationArrowBack:    utils.MustIcon(widget.NewIcon(icons.NavigationArrowBack)),
		navigationArrowForward: utils.MustIcon(widget.NewIcon(icons.NavigationArrowForward)),
		actionInfo:             utils.MustIcon(widget.NewIcon(icons.ActionInfo)),
		actionCheck:            utils.MustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		navigationCancel:       utils.MustIcon(widget.NewIcon(icons.NavigationCancel)),
		imageBrightness1:       utils.MustIcon(widget.NewIcon(icons.ImageBrightness1)),
		chevronRight:           utils.MustIcon(widget.NewIcon(icons.NavigationChevronRight)),
		contentClear:           utils.MustIcon(widget.NewIcon(icons.ContentClear)),
		navMoreIcon:            utils.MustIcon(widget.NewIcon(icons.NavigationMoreHoriz)),
		dropDownIcon:           utils.MustIcon(widget.NewIcon(icons.NavigationArrowDropDown)),
		cached:                 utils.MustIcon(widget.NewIcon(icons.ActionCached)),
	}

	common := pageCommon{
		printer:                 message.NewPrinter(language.English),
		theme:                   d.theme,
		icons:                   ic,
		returnPage:              &d.previous,
		page:                    &d.current,
		toastLoad:               &toast{},
		minimizeNavDrawerButton: d.theme.PlainIconButton(new(widget.Clickable), ic.navigationArrowBack),
		maximizeNavDrawerButton: d.theme.PlainIconButton(new(widget.Clickable), ic.navigationArrowForward),
		selectedUTXO:            make(map[int]map[int32]map[string]*wallet.UnspentOutput),
		modal:                   d.theme.Modal(),
		subPageBackButton:       d.theme.PlainIconButton(new(widget.Clickable), ic.navigationArrowBack),
		subPageInfoButton:       d.theme.PlainIconButton(new(widget.Clickable), ic.actionInfo),
		changePage:              d.changePage,
		setReturnPage:           d.setReturnPage,
		refreshWindow:           d.refresh,

		switchView: d.switchView,
	}

	common.fetchExchangeValue(&common.dcrUsdtBittrex)

	iconColor := common.theme.Color.Gray3

	common.testButton = d.theme.Button(new(widget.Clickable), "test button")
	isNavDrawerMinimized := false
	common.isNavDrawerMinimized = &isNavDrawerMinimized

	common.minimizeNavDrawerButton.Color, common.maximizeNavDrawerButton.Color = iconColor, iconColor

	zeroInset := layout.UniformInset(values.MarginPadding0)
	common.subPageBackButton.Color, common.subPageInfoButton.Color = iconColor, iconColor

	m25 := values.MarginPadding25
	common.subPageBackButton.Size, common.subPageInfoButton.Size = m25, m25
	common.subPageBackButton.Inset, common.subPageInfoButton.Inset = zeroInset, zeroInset

	d.pages = make(map[string]layout.Widget)
	d.pages[PageDex] = d.DexPage(common)
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
	go func() {
		page.toast <- &toast{
			text:    text,
			success: success,
		}
	}()
}

func (page pageCommon) Layout(gtx layout.Context, body layout.Widget) layout.Dimensions {

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			// fill the entire window with a color if a user has no wallet created
			return decredmaterial.Fill(gtx, page.theme.Color.Surface)
		}),
		layout.Stacked(func(gtx C) D {
			// stack the page content on the entire window if a user has no wallet
			return body(gtx)
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

	if page.subPageBackButton.Button.Clicked() {
		sp.back()
	}

	if sp.extraItem != nil && sp.extraItem.Clicked() {
		sp.handleExtra()
	}
}

package ui

import (
	"encoding/json"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/page"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"image"
	"net/http"
	"time"
)

type DCRUSDTBittrex struct {
	LastTradeRate string
}

type Modal interface {
	ModalID() string
	OnResume()
	Layout(gtx layout.Context) layout.Dimensions
	OnDismiss()
	Dismiss()
	Show()
	Handle()
}

type Page interface {
	Layout(layout.Context) layout.Dimensions
	handle()
	onClose()
}

type toast struct {
	text    string
	success bool
	timer   *time.Timer
}

type Common struct {
	Theme              *decredmaterial.Theme

	printer          *message.Printer
	multiWallet      *dcrlibwallet.MultiWallet
	SyncStatusUpdate chan wallet.SyncStatusUpdate
	Wallet           *wallet.Wallet
	WalletAccount    **wallet.Account
	info             *wallet.MultiWalletInfo
	SelectedWallet   *int
	SelectedAccount  *int
	Icons            Icons
	page             *string
	returnPage       *string
	dcrUsdtBittrex   DCRUSDTBittrex
	navTab           *decredmaterial.Tabs
	keyEvents        chan *key.Event
	toast            *toast
	states           *states
	internalLog      *chan string
	walletSyncStatus *wallet.SyncStatus
	walletTransactions **wallet.Transactions
	walletTransaction  **wallet.Transaction
	acctMixerStatus    *chan *wallet.AccountMixer
	selectedProposal   **dcrlibwallet.Proposal
	proposals          **wallet.Proposals
	syncedProposal     chan *wallet.Proposal
	txAuthor           *dcrlibwallet.TxAuthor
	broadcastResult    *wallet.Broadcast
	signatureResult    **wallet.Signature
	walletTickets      **wallet.Tickets
	vspInfo            **wallet.VSP
	unspentOutputs     **wallet.UnspentOutputs
	ShowModal          func()
	DismissModal       func(Modal)

	testButton decredmaterial.Button

	selectedUTXO map[int]map[int32]map[string]*wallet.UnspentOutput

	SubPageBackButton decredmaterial.IconButton
	SubPageInfoButton decredmaterial.IconButton

	changeWindowPage func(Page)
	changePage       func(string)
	setReturnPage    func(string)

	WallAcctSelector *page.WalletAccountSelector
}

type Icons struct {
	contentAdd, navigationCheck, navigationMore, actionCheckCircle, actionInfo, NavigationArrowBack,
	NavigationArrowForward, actionCheck, chevronRight, navigationCancel, navMoreIcon,
	imageBrightness1, contentClear, dropDownIcon, cached *widget.Icon

	overviewIcon, overviewIconInactive, walletIcon, walletIconInactive,
	receiveIcon, transactionIcon, transactionIconInactive, sendIcon, moreIcon, moreIconInactive,
	pendingIcon, logo, redirectIcon, confirmIcon, newWalletIcon, walletAlertIcon,
	importedAccountIcon, accountIcon, editIcon, expandIcon, copyIcon, mixer, mixerSmall,
	arrowForwardIcon, transactionFingerPrintIcon, settingsIcon, securityIcon, helpIcon,
	aboutIcon, debugIcon, verifyMessageIcon, locationPinIcon, alertGray, arrowDownIcon,
	watchOnlyWalletIcon, currencySwapIcon, syncingIcon, proposalIconActive, proposalIconInactive,
	restore, documentationIcon, downloadIcon, timerIcon, ticketIcon, ticketIconInactive, stakeyIcon,
	list, listGridIcon, decredSymbolIcon *widget.Image

	ticketPurchasedIcon,
	ticketImmatureIcon,
	ticketLiveIcon,
	ticketVotedIcon,
	ticketMissedIcon,
	ticketExpiredIcon,
	ticketRevokedIcon,
	ticketUnminedIcon *widget.Image
}

func (win *Window) newCommon(decredIcons map[string]image.Image) *Common {
	ic := Icons{
		contentAdd:             MustIcon(widget.NewIcon(icons.ContentAdd)),
		navigationCheck:        MustIcon(widget.NewIcon(icons.NavigationCheck)),
		navigationMore:         MustIcon(widget.NewIcon(icons.NavigationMoreHoriz)),
		actionCheckCircle:      MustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		NavigationArrowBack:    MustIcon(widget.NewIcon(icons.NavigationArrowBack)),
		NavigationArrowForward: MustIcon(widget.NewIcon(icons.NavigationArrowForward)),
		actionInfo:             MustIcon(widget.NewIcon(icons.ActionInfo)),
		actionCheck:            MustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		navigationCancel:       MustIcon(widget.NewIcon(icons.NavigationCancel)),
		imageBrightness1:       MustIcon(widget.NewIcon(icons.ImageBrightness1)),
		chevronRight:           MustIcon(widget.NewIcon(icons.NavigationChevronRight)),
		contentClear:           MustIcon(widget.NewIcon(icons.ContentClear)),
		navMoreIcon:            MustIcon(widget.NewIcon(icons.NavigationMoreHoriz)),
		dropDownIcon:           MustIcon(widget.NewIcon(icons.NavigationArrowDropDown)),
		cached:                 MustIcon(widget.NewIcon(icons.ActionCached)),

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
		decredSymbolIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["decred_symbol"])},
	}

	common := &Common{
		printer:          message.NewPrinter(language.English),
		multiWallet:      win.wallet.GetMultiWallet(),
		SyncStatusUpdate: make(chan wallet.SyncStatusUpdate, 10),
		wallet:           win.wallet,
		walletAccount:    &win.walletAccount,
		info:             win.walletInfo,
		SelectedWallet:   &win.selected,
		SelectedAccount:  &win.selectedAccount,
		Theme:            win.theme,
		keyEvents:        win.keyEvents,
		states:           &win.states,
		Icons:            ic,
		walletSyncStatus: win.walletSyncStatus,
		walletTransactions: &win.walletTransactions,
		walletTransaction:  &win.walletTransaction,
		acctMixerStatus:    &win.walletAcctMixerStatus,
		selectedProposal:   &win.selectedProposal,
		proposals:          &win.proposals,
		syncedProposal:     win.proposal,
		txAuthor:           &win.txAuthor,
		broadcastResult:    &win.broadcastResult,
		signatureResult:    &win.signatureResult,
		walletTickets:      &win.walletTickets,
		vspInfo:            &win.vspInfo,
		unspentOutputs:     &win.walletUnspentOutputs,
		ShowModal:          win.showModal,
		DismissModal:       win.dismissModal,
		changeWindowPage:   win.changePage,

		selectedUTXO: make(map[int]map[int32]map[string]*wallet.UnspentOutput),
		toast:        win.toast,
		internalLog:  &win.internalLog,
	}

	common.subPageBackButton = win.theme.PlainIconButton(new(widget.Clickable), common.Icons.NavigationArrowBack)
	common.subPageInfoButton = win.theme.PlainIconButton(new(widget.Clickable), common.Icons.actionInfo)

	if common.fetchExchangeValue(&common.dcrUsdtBittrex) != nil {
		log.Info("Error fetching exchange value")
	}

	return common
}

func (c *Common) fetchExchangeValue(target interface{}) error {
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

func (c *Common) refreshTheme() {
	isDarkModeOn := c.wallet.ReadBoolConfigValueForKey("isDarkModeOn")
	if isDarkModeOn != c.Theme.DarkMode {
		c.Theme.SwitchDarkMode(isDarkModeOn)
	}
}

func (c *Common) notify(text string, success bool) {
	c.toast = &toast{
		text:    text,
		success: success,
	}
}

func (common *Common) UniformPadding(gtx layout.Context, body layout.Widget) layout.Dimensions {
	width := gtx.Constraints.Max.X

	padding := values.MarginPadding24

	if (width - 2*gtx.Px(padding)) > gtx.Px(MaxWidth) {
		paddingValue := float32(width-gtx.Px(MaxWidth)) / 2
		padding = unit.Px(paddingValue)
	}

	return layout.Inset{
		Top:    values.MarginPadding24,
		Right:  padding,
		Bottom: values.MarginPadding24,
		Left:   padding,
	}.Layout(gtx, body)
}

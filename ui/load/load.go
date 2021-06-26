package load

import (
	"encoding/json"
	"golang.org/x/text/language"
	"image"
	"net/http"
	"sort"
	"time"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/wallet"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/text/message"
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
	Handle()
	OnClose()
	OnResume()
}

type Toast struct {
	text    string
	success bool
	Timer   *time.Timer
}

type walletLoad struct {
	MultiWallet      *dcrlibwallet.MultiWallet
	TxAuthor         *dcrlibwallet.TxAuthor
	SelectedProposal *dcrlibwallet.Proposal

	Proposals       *wallet.Proposals
	SyncStatus      *wallet.SyncStatus
	Transactions    *wallet.Transactions
	Transaction     *wallet.Transaction
	BroadcastResult *wallet.Broadcast
	Tickets         *wallet.Tickets
	VspInfo         *wallet.VSP
	UnspentOutputs  *wallet.UnspentOutputs
	Wallet          *wallet.Wallet
	Account         *wallet.Account
	Info            *wallet.MultiWalletInfo

	SelectedWallet  *int
	SelectedAccount *int
}

type Receiver struct {
	InternalLog         chan string
	NotificationsUpdate chan interface{}
	KeyEvents           chan *key.Event
	AcctMixerStatus     chan *wallet.AccountMixer
	SyncedProposal      chan *wallet.Proposal
}

type Icons struct {
	ContentAdd, NavigationCheck, NavigationMore, ActionCheckCircle, ActionInfo, NavigationArrowBack,
	NavigationArrowForward, ActionCheck, ChevronRight, NavigationCancel, NavMoreIcon,
	ImageBrightness1, ContentClear, DropDownIcon, Cached *widget.Icon

	OverviewIcon, OverviewIconInactive, WalletIcon, WalletIconInactive,
	ReceiveIcon, TransactionIcon, TransactionIconInactive, SendIcon, MoreIcon, MoreIconInactive,
	pendingIcon, Logo, redirectIcon, confirmIcon, newWalletIcon, walletAlertIcon,
	importedAccountIcon, accountIcon, editIcon, expandIcon, copyIcon, mixer, mixerSmall,
	arrowForwardIcon, transactionFingerPrintIcon, settingsIcon, securityIcon, helpIcon,
	aboutIcon, debugIcon, verifyMessageIcon, locationPinIcon, alertGray, arrowDownIcon,
	watchOnlyWalletIcon, currencySwapIcon, syncingIcon, ProposalIconActive, ProposalIconInactive,
	restore, documentationIcon, downloadIcon, timerIcon, TicketIcon, TicketIconInactive, stakeyIcon,
	list, listGridIcon, decredSymbolIcon *widget.Image

	ticketPurchasedIcon,
	TicketImmatureIcon,
	TicketLiveIcon,
	TicketVotedIcon,
	TicketMissedIcon,
	TicketExpiredIcon,
	TicketRevokedIcon,
	TicketUnminedIcon *widget.Image
}

type Load struct {
	Theme *decredmaterial.Theme

	WL       *walletLoad
	Receiver *Receiver
	Printer  *message.Printer

	Icons          Icons
	Page           *string
	ReturnPage     *string
	DcrUsdtBittrex *DCRUSDTBittrex

	Toast *Toast

	SelectedUTXO map[int]map[int32]map[string]*wallet.UnspentOutput

	SubPageBackButton decredmaterial.IconButton
	SubPageInfoButton decredmaterial.IconButton

	ToggleSync       func()
	ShowModal        func(Modal)
	DismissModal     func(Modal)
	ChangeWindowPage func(Page)
	ChangePage       func(string)
	ChangeFragment   func(Page, string)
	SetReturnPage    func(string)
}

func NewLoad(th *decredmaterial.Theme, decredIcons map[string]image.Image) *Load {
	ic := Icons{
		ContentAdd:             mustIcon(widget.NewIcon(icons.ContentAdd)),
		NavigationCheck:        mustIcon(widget.NewIcon(icons.NavigationCheck)),
		NavigationMore:         mustIcon(widget.NewIcon(icons.NavigationMoreHoriz)),
		ActionCheckCircle:      mustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		NavigationArrowBack:    mustIcon(widget.NewIcon(icons.NavigationArrowBack)),
		NavigationArrowForward: mustIcon(widget.NewIcon(icons.NavigationArrowForward)),
		ActionInfo:             mustIcon(widget.NewIcon(icons.ActionInfo)),
		ActionCheck:            mustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		NavigationCancel:       mustIcon(widget.NewIcon(icons.NavigationCancel)),
		ImageBrightness1:       mustIcon(widget.NewIcon(icons.ImageBrightness1)),
		ChevronRight:           mustIcon(widget.NewIcon(icons.NavigationChevronRight)),
		ContentClear:           mustIcon(widget.NewIcon(icons.ContentClear)),
		NavMoreIcon:            mustIcon(widget.NewIcon(icons.NavigationMoreHoriz)),
		DropDownIcon:           mustIcon(widget.NewIcon(icons.NavigationArrowDropDown)),
		Cached:                 mustIcon(widget.NewIcon(icons.ActionCached)),

		OverviewIcon:            &widget.Image{Src: paint.NewImageOp(decredIcons["overview"])},
		OverviewIconInactive:    &widget.Image{Src: paint.NewImageOp(decredIcons["overview_inactive"])},
		WalletIconInactive:      &widget.Image{Src: paint.NewImageOp(decredIcons["wallet_inactive"])},
		ReceiveIcon:             &widget.Image{Src: paint.NewImageOp(decredIcons["receive"])},
		TransactionIcon:         &widget.Image{Src: paint.NewImageOp(decredIcons["transaction"])},
		TransactionIconInactive: &widget.Image{Src: paint.NewImageOp(decredIcons["transaction_inactive"])},
		SendIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["send"])},
		MoreIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["more"])},
		MoreIconInactive:        &widget.Image{Src: paint.NewImageOp(decredIcons["more_inactive"])},
		Logo:                    &widget.Image{Src: paint.NewImageOp(decredIcons["logo"])},
		confirmIcon:             &widget.Image{Src: paint.NewImageOp(decredIcons["confirmed"])},
		pendingIcon:             &widget.Image{Src: paint.NewImageOp(decredIcons["pending"])},
		redirectIcon:            &widget.Image{Src: paint.NewImageOp(decredIcons["redirect"])},
		newWalletIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["addNewWallet"])},
		walletAlertIcon:         &widget.Image{Src: paint.NewImageOp(decredIcons["walletAlert"])},
		accountIcon:             &widget.Image{Src: paint.NewImageOp(decredIcons["account"])},
		importedAccountIcon:     &widget.Image{Src: paint.NewImageOp(decredIcons["imported_account"])},
		editIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["editIcon"])},
		expandIcon:              &widget.Image{Src: paint.NewImageOp(decredIcons["expand_icon"])},
		copyIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["copy_icon"])},
		mixer:                      &widget.Image{Src: paint.NewImageOp(decredIcons["mixer"])},
		mixerSmall:                 &widget.Image{Src: paint.NewImageOp(decredIcons["mixer_small"])},
		transactionFingerPrintIcon: &widget.Image{Src: paint.NewImageOp(decredIcons["transaction_fingerprint"])},
		arrowForwardIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["arrow_forward"])},
		settingsIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["settings"])},
		securityIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["security"])},
		helpIcon:                   &widget.Image{Src: paint.NewImageOp(decredIcons["help_icon"])},
		aboutIcon:            &widget.Image{Src: paint.NewImageOp(decredIcons["about_icon"])},
		debugIcon:            &widget.Image{Src: paint.NewImageOp(decredIcons["debug"])},
		verifyMessageIcon:    &widget.Image{Src: paint.NewImageOp(decredIcons["verify_message"])},
		locationPinIcon:      &widget.Image{Src: paint.NewImageOp(decredIcons["location_pin"])},
		alertGray:            &widget.Image{Src: paint.NewImageOp(decredIcons["alert_gray"])},
		arrowDownIcon:        &widget.Image{Src: paint.NewImageOp(decredIcons["arrow_down"])},
		watchOnlyWalletIcon:  &widget.Image{Src: paint.NewImageOp(decredIcons["watch_only_wallet"])},
		currencySwapIcon:     &widget.Image{Src: paint.NewImageOp(decredIcons["swap"])},
		syncingIcon:          &widget.Image{Src: paint.NewImageOp(decredIcons["syncing"])},
		documentationIcon:    &widget.Image{Src: paint.NewImageOp(decredIcons["documentation"])},
		ProposalIconActive:   &widget.Image{Src: paint.NewImageOp(decredIcons["politeiaActive"])},
		ProposalIconInactive: &widget.Image{Src: paint.NewImageOp(decredIcons["politeiaInactive"])},
		restore:              &widget.Image{Src: paint.NewImageOp(decredIcons["restore"])},
		downloadIcon:         &widget.Image{Src: paint.NewImageOp(decredIcons["downloadIcon"])},
		timerIcon:            &widget.Image{Src: paint.NewImageOp(decredIcons["timerIcon"])},
		WalletIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["wallet"])},
		TicketIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["ticket"])},
		TicketIconInactive:   &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_inactive"])},
		stakeyIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["stakey"])},
		ticketPurchasedIcon:  &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_purchased"])},
		TicketImmatureIcon:   &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_immature"])},
		TicketUnminedIcon:    &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_unmined"])},
		TicketLiveIcon:       &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_live"])},
		TicketVotedIcon:      &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_voted"])},
		TicketMissedIcon:     &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_missed"])},
		TicketExpiredIcon:    &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_expired"])},
		TicketRevokedIcon:    &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_revoked"])},
		list:                 &widget.Image{Src: paint.NewImageOp(decredIcons["list"])},
		listGridIcon:        &widget.Image{Src: paint.NewImageOp(decredIcons["list_grid"])},
		decredSymbolIcon:    &widget.Image{Src: paint.NewImageOp(decredIcons["decred_symbol"])},
	}

	wl := &walletLoad{
		Wallet:          new(wallet.Wallet),
		Account:         new(wallet.Account),
		Info:            new(wallet.MultiWalletInfo),
		SyncStatus:      new(wallet.SyncStatus),
		Transactions:    new(wallet.Transactions),
		UnspentOutputs:  new(wallet.UnspentOutputs),
		Tickets:         new(wallet.Tickets),
		VspInfo:         new(wallet.VSP),
		BroadcastResult: new(wallet.Broadcast),
		Proposals:       new(wallet.Proposals),

		SelectedProposal: new(dcrlibwallet.Proposal),
		TxAuthor:         new(dcrlibwallet.TxAuthor),
	}

	r := &Receiver{
		AcctMixerStatus: make(chan *wallet.AccountMixer),
		SyncedProposal:  make(chan *wallet.Proposal),
	}

	l :=  &Load{
		Theme:    th,
		Icons:    ic,
		WL:       wl,
		Receiver: r,
		Toast: &Toast{},

		Printer: message.NewPrinter(language.English),

		SubPageBackButton: th.PlainIconButton(new(widget.Clickable), ic.NavigationArrowBack),
		SubPageInfoButton: th.PlainIconButton(new(widget.Clickable), ic.ActionInfo),
	}
	fetchExchangeValue(l.DcrUsdtBittrex)

	return l
}

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func (l *Load) CreateToast(text string, success bool) {
	   l.Toast = &Toast{
	   	  text: text,
	   	  success: success,
	   }
}

func (wl *walletLoad) SortedWalletList() []*dcrlibwallet.Wallet {
	wallets := wl.MultiWallet.AllWallets()

	sort.Slice(wallets, func(i, j int) bool {
		return wallets[i].ID < wallets[j].ID
	})

	return wallets
}

func fetchExchangeValue(target interface{}) {
	url := "https://api.bittrex.com/v3/markets/DCR-USDT/ticker"
	res, err := http.Get(url)
	if err != nil {
		log.Error(err)
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(target)
	if err != nil {
		log.Error(err)
	}
}

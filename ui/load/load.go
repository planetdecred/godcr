// The load package contains data structures that are shared by components in the ui package. It is not a dumping ground
// for code you feel might be shared with other components in the future. Before adding code here, ask yourself, can
// the code be isolated in the package you're calling it from? Is it really needed by other packages in the ui package?
// or you're just planning for a use case that might never used.

package load

import (
	"errors"

	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"gioui.org/io/key"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/assets"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/notification"
	"github.com/planetdecred/godcr/wallet"
)

type DCRUSDTBittrex struct {
	LastTradeRate string
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
	ImageBrightness1, ContentClear, DropDownIcon, Cached, ContentRemove *widget.Icon

	OverviewIcon, OverviewIconInactive, WalletIcon, WalletIconInactive,
	ReceiveIcon, Transferred, TransactionIcon, TransactionIconInactive, SendIcon, MoreIcon, MoreIconInactive,
	PendingIcon, Logo, RedirectIcon, ConfirmIcon, NewWalletIcon, WalletAlertIcon,
	ImportedAccountIcon, AccountIcon, EditIcon, expandIcon, CopyIcon, MixedTx, mixer, MixerSmall,
	ArrowForwardIcon, Next, TransactionFingerPrintIcon, SettingsIcon, SecurityIcon, HelpIcon,
	AboutIcon, DebugIcon, VerifyMessageIcon, LocationPinIcon, AlertGray, ArrowDownIcon,
	WatchOnlyWalletIcon, CurrencySwapIcon, SyncingIcon, ProposalIconActive, ProposalIconInactive,
	Restore, DocumentationIcon, DownloadIcon, TimerIcon, TicketIcon, TicketIconInactive, StakeyIcon,
	List, ListGridIcon, DecredSymbolIcon, DecredSymbol2 *widget.Image

	TicketPurchasedIcon,
	TicketImmatureIcon,
	TicketLiveIcon,
	TicketVotedIcon,
	TicketMissedIcon,
	TicketExpiredIcon,
	TicketRevokedIcon,
	TicketUnminedIcon *decredmaterial.Image
}

type Load struct {
	Theme *decredmaterial.Theme

	WL       *WalletLoad
	Receiver *Receiver
	Printer  *message.Printer
	Network  string

	Icons Icons

	Toast *notification.Toast

	SelectedWallet  *int
	SelectedAccount *int
	SelectedUTXO    map[int]map[int32]map[string]*wallet.UnspentOutput

	ToggleSync       func()
	RefreshWindow    func()
	ShowModal        func(Modal)
	DismissModal     func(Modal)
	ChangeWindowPage func(page Page, keepBackStack bool)
	PopWindowPage    func() bool
	ChangeFragment   func(page Page)
	PopFragment      func()
	PopToFragment    func(pageID string)
}

func NewLoad() (*Load, error) {

	wl := &WalletLoad{
		Wallet:         new(wallet.Wallet),
		Account:        new(wallet.Account),
		Info:           new(wallet.MultiWalletInfo),
		SyncStatus:     new(wallet.SyncStatus),
		Transactions:   new(wallet.Transactions),
		UnspentOutputs: new(wallet.UnspentOutputs),
		Tickets:        new(*wallet.Tickets),
		VspInfo:        new(wallet.VSP),
		Proposals:      new(wallet.Proposals),

		SelectedProposal: new(dcrlibwallet.Proposal),
	}

	r := &Receiver{
		AcctMixerStatus: make(chan *wallet.AccountMixer),
		SyncedProposal:  make(chan *wallet.Proposal),
	}

	icons := loadIcons()
	th := decredmaterial.NewTheme(assets.FontCollection(), assets.DecredIcons, false)
	if th == nil {
		return nil, errors.New("unexpected error while loading theme")
	}
	l := &Load{
		Theme:    th,
		Icons:    icons,
		WL:       wl,
		Receiver: r,
		Toast:    notification.NewToast(th),

		Printer: message.NewPrinter(language.English),
	}

	return l, nil
}

func (l *Load) RefreshTheme() {
	isDarkModeOn := l.WL.MultiWallet.ReadBoolConfigValueForKey("isDarkModeOn", false)
	if isDarkModeOn != l.Theme.DarkMode {
		l.Theme.SwitchDarkMode(isDarkModeOn)
	}
}

func loadIcons() Icons {
	decredIcons := assets.DecredIcons

	ic := Icons{
		ContentAdd:             decredmaterial.MustIcon(widget.NewIcon(icons.ContentAdd)),
		NavigationCheck:        decredmaterial.MustIcon(widget.NewIcon(icons.NavigationCheck)),
		NavigationMore:         decredmaterial.MustIcon(widget.NewIcon(icons.NavigationMoreHoriz)),
		ActionCheckCircle:      decredmaterial.MustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		NavigationArrowBack:    decredmaterial.MustIcon(widget.NewIcon(icons.NavigationArrowBack)),
		NavigationArrowForward: decredmaterial.MustIcon(widget.NewIcon(icons.NavigationArrowForward)),
		ActionInfo:             decredmaterial.MustIcon(widget.NewIcon(icons.ActionInfo)),
		ActionCheck:            decredmaterial.MustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		NavigationCancel:       decredmaterial.MustIcon(widget.NewIcon(icons.NavigationCancel)),
		ImageBrightness1:       decredmaterial.MustIcon(widget.NewIcon(icons.ImageBrightness1)),
		ChevronRight:           decredmaterial.MustIcon(widget.NewIcon(icons.NavigationChevronRight)),
		ContentClear:           decredmaterial.MustIcon(widget.NewIcon(icons.ContentClear)),
		NavMoreIcon:            decredmaterial.MustIcon(widget.NewIcon(icons.NavigationMoreHoriz)),
		DropDownIcon:           decredmaterial.MustIcon(widget.NewIcon(icons.NavigationArrowDropDown)),
		Cached:                 decredmaterial.MustIcon(widget.NewIcon(icons.ActionCached)),
		ContentRemove:          decredmaterial.MustIcon(widget.NewIcon(icons.ContentRemove)),

		OverviewIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["overview"])},
		OverviewIconInactive:       &widget.Image{Src: paint.NewImageOp(decredIcons["overview_inactive"])},
		WalletIconInactive:         &widget.Image{Src: paint.NewImageOp(decredIcons["wallet_inactive"])},
		ReceiveIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["receive"])},
		Transferred:                &widget.Image{Src: paint.NewImageOp(decredIcons["transferred"])},
		TransactionIcon:            &widget.Image{Src: paint.NewImageOp(decredIcons["transaction"])},
		TransactionIconInactive:    &widget.Image{Src: paint.NewImageOp(decredIcons["transaction_inactive"])},
		SendIcon:                   &widget.Image{Src: paint.NewImageOp(decredIcons["send"])},
		MoreIcon:                   &widget.Image{Src: paint.NewImageOp(decredIcons["more"])},
		MoreIconInactive:           &widget.Image{Src: paint.NewImageOp(decredIcons["more_inactive"])},
		Logo:                       &widget.Image{Src: paint.NewImageOp(decredIcons["logo"])},
		ConfirmIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["confirmed"])},
		PendingIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["pending"])},
		RedirectIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["redirect"])},
		NewWalletIcon:              &widget.Image{Src: paint.NewImageOp(decredIcons["addNewWallet"])},
		WalletAlertIcon:            &widget.Image{Src: paint.NewImageOp(decredIcons["walletAlert"])},
		AccountIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["account"])},
		ImportedAccountIcon:        &widget.Image{Src: paint.NewImageOp(decredIcons["imported_account"])},
		EditIcon:                   &widget.Image{Src: paint.NewImageOp(decredIcons["editIcon"])},
		expandIcon:                 &widget.Image{Src: paint.NewImageOp(decredIcons["expand_icon"])},
		CopyIcon:                   &widget.Image{Src: paint.NewImageOp(decredIcons["copy_icon"])},
		MixedTx:                    &widget.Image{Src: paint.NewImageOp(decredIcons["mixed_tx"])},
		mixer:                      &widget.Image{Src: paint.NewImageOp(decredIcons["mixer"])},
		MixerSmall:                 &widget.Image{Src: paint.NewImageOp(decredIcons["mixer_small"])},
		TransactionFingerPrintIcon: &widget.Image{Src: paint.NewImageOp(decredIcons["transaction_fingerprint"])},
		ArrowForwardIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["arrow_forward"])},
		Next:                       &widget.Image{Src: paint.NewImageOp(decredIcons["ic_next"])},
		SettingsIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["settings"])},
		SecurityIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["security"])},
		HelpIcon:                   &widget.Image{Src: paint.NewImageOp(decredIcons["help_icon"])},
		AboutIcon:                  &widget.Image{Src: paint.NewImageOp(decredIcons["about_icon"])},
		DebugIcon:                  &widget.Image{Src: paint.NewImageOp(decredIcons["debug"])},
		VerifyMessageIcon:          &widget.Image{Src: paint.NewImageOp(decredIcons["verify_message"])},
		LocationPinIcon:            &widget.Image{Src: paint.NewImageOp(decredIcons["location_pin"])},
		AlertGray:                  &widget.Image{Src: paint.NewImageOp(decredIcons["alert_gray"])},
		ArrowDownIcon:              &widget.Image{Src: paint.NewImageOp(decredIcons["arrow_down"])},
		WatchOnlyWalletIcon:        &widget.Image{Src: paint.NewImageOp(decredIcons["watch_only_wallet"])},
		CurrencySwapIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["swap"])},
		SyncingIcon:                &widget.Image{Src: paint.NewImageOp(decredIcons["syncing"])},
		DocumentationIcon:          &widget.Image{Src: paint.NewImageOp(decredIcons["documentation"])},
		ProposalIconActive:         &widget.Image{Src: paint.NewImageOp(decredIcons["politeiaActive"])},
		ProposalIconInactive:       &widget.Image{Src: paint.NewImageOp(decredIcons["politeiaInactive"])},
		Restore:                    &widget.Image{Src: paint.NewImageOp(decredIcons["restore"])},
		DownloadIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["downloadIcon"])},
		TimerIcon:                  &widget.Image{Src: paint.NewImageOp(decredIcons["timerIcon"])},
		WalletIcon:                 &widget.Image{Src: paint.NewImageOp(decredIcons["wallet"])},
		TicketIcon:                 &widget.Image{Src: paint.NewImageOp(decredIcons["ticket"])},
		TicketIconInactive:         &widget.Image{Src: paint.NewImageOp(decredIcons["ticket_inactive"])},
		StakeyIcon:                 &widget.Image{Src: paint.NewImageOp(decredIcons["stakey"])},
		TicketPurchasedIcon:        decredmaterial.NewImage(decredIcons["ticket_purchased"]),
		TicketImmatureIcon:         decredmaterial.NewImage(decredIcons["ticket_immature"]),
		TicketUnminedIcon:          decredmaterial.NewImage(decredIcons["ticket_unmined"]),
		TicketLiveIcon:             decredmaterial.NewImage(decredIcons["ticket_live"]),
		TicketVotedIcon:            decredmaterial.NewImage(decredIcons["ticket_voted"]),
		TicketMissedIcon:           decredmaterial.NewImage(decredIcons["ticket_missed"]),
		TicketExpiredIcon:          decredmaterial.NewImage(decredIcons["ticket_expired"]),
		TicketRevokedIcon:          decredmaterial.NewImage(decredIcons["ticket_revoked"]),
		List:                       &widget.Image{Src: paint.NewImageOp(decredIcons["list"])},
		ListGridIcon:               &widget.Image{Src: paint.NewImageOp(decredIcons["list_grid"])},
		DecredSymbolIcon:           &widget.Image{Src: paint.NewImageOp(decredIcons["decred_symbol"])},
		DecredSymbol2:              &widget.Image{Src: paint.NewImageOp(decredIcons["ic_decred02"])},
	}
	return ic
}

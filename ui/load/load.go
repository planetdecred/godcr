// The load package contains data structures that are shared by components in the ui package. It is not a dumping ground
// for code you feel might be shared with other components in the future. Before adding code here, ask yourself, can
// the code be isolated in the package you're calling it from? Is it really needed by other packages in the ui package?
// or you're just planning for a use case that might never used.

package load

import (
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/text/message"

	"gioui.org/io/key"
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
	KeyEvents map[string]chan *key.Event
}

type Icons struct {
	ContentAdd, NavigationCheck, NavigationMore, ActionCheckCircle, ActionInfo, NavigationArrowBack,
	NavigationArrowForward, ActionCheck, ChevronRight, NavigationCancel, NavMoreIcon,
	ImageBrightness1, ContentClear, DropDownIcon, Cached, ContentRemove, ConcealIcon, RevealIcon,
	SearchIcon, PlayIcon *widget.Icon

	OverviewIcon, OverviewIconInactive, WalletIcon, WalletIconInactive,
	ReceiveIcon, Transferred, TransactionsIcon, TransactionsIconInactive, SendIcon, MoreIcon, MoreIconInactive,
	PendingIcon, Logo, RedirectIcon, ConfirmIcon, NewWalletIcon, WalletAlertIcon,
	ImportedAccountIcon, AccountIcon, EditIcon, expandIcon, CopyIcon, MixedTx, Mixer, PrivacySetup,
	Next, SettingsIcon, SecurityIcon, HelpIcon,
	AboutIcon, DebugIcon, VerifyMessageIcon, LocationPinIcon, AlertGray, ArrowDownIcon,
	WatchOnlyWalletIcon, CurrencySwapIcon, SyncingIcon, ProposalIconActive, ProposalIconInactive,
	Restore, DocumentationIcon, DownloadIcon, TimerIcon, StakeIcon, StakeIconInactive, StakeyIcon,
	List, ListGridIcon, DecredSymbolIcon, DecredSymbol2, GovernanceActiveIcon, GovernanceInactiveIcon,
	LogoDarkMode, TimerDarkMode, Rebroadcast, SettingsActiveIcon, SettingsInactiveIcon, ActivatedActiveIcon, ActivatedInactiveIcon, LockinActiveIcon, LockinInactiveIcon *decredmaterial.Image

	NewStakeIcon,
	TicketImmatureIcon,
	TicketLiveIcon,
	TicketVotedIcon,
	TicketMissedIcon,
	TicketExpiredIcon,
	TicketRevokedIcon,
	TicketUnminedIcon *decredmaterial.Image

	DexIcon, DexIconInactive, BTC, DCR *decredmaterial.Image
}

type Load struct {
	Theme *decredmaterial.Theme

	WL       *WalletLoad
	Receiver *Receiver
	Printer  *message.Printer
	Network  string

	Icons Icons

	Toast *notification.Toast

	SelectedUTXO map[int]map[int32]map[string]*wallet.UnspentOutput

	ToggleSync          func()
	RefreshWindow       func()
	ShowModal           func(Modal)
	DismissModal        func(Modal)
	ChangeWindowPage    func(page Page, keepBackStack bool)
	PopWindowPage       func() bool
	ChangeFragment      func(page Page)
	PopFragment         func()
	PopToFragment       func(pageID string)
	SubscribeKeyEvent   func(eventChan chan *key.Event, pageID string) // Widgets call this function to recieve key events.
	UnsubscribeKeyEvent func(pageID string) error
	ReloadApp           func()
}

func (l *Load) RefreshTheme() {
	isDarkModeOn := l.WL.MultiWallet.ReadBoolConfigValueForKey(DarkModeConfigKey, false)
	l.Theme.SwitchDarkMode(isDarkModeOn, assets.DecredIcons)
	l.RefreshWindow()
}
func (l *Load) Dexc() *dcrlibwallet.DexClient {
	return l.WL.MultiWallet.DexClient()
}

func IconSet() Icons {
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
		ConcealIcon:            decredmaterial.MustIcon(widget.NewIcon(icons.ActionVisibility)),
		RevealIcon:             decredmaterial.MustIcon(widget.NewIcon(icons.ActionVisibilityOff)),
		SearchIcon:             decredmaterial.MustIcon(widget.NewIcon(icons.ActionSearch)),
		PlayIcon:               decredmaterial.MustIcon(widget.NewIcon(icons.AVPlayArrow)),

		OverviewIcon:             decredmaterial.NewImage(decredIcons["overview"]),
		OverviewIconInactive:     decredmaterial.NewImage(decredIcons["overview_inactive"]),
		WalletIconInactive:       decredmaterial.NewImage(decredIcons["wallet_inactive"]),
		ReceiveIcon:              decredmaterial.NewImage(decredIcons["receive"]),
		Transferred:              decredmaterial.NewImage(decredIcons["transferred"]),
		TransactionsIcon:         decredmaterial.NewImage(decredIcons["transactions"]),
		TransactionsIconInactive: decredmaterial.NewImage(decredIcons["transactions_inactive"]),
		SendIcon:                 decredmaterial.NewImage(decredIcons["send"]),
		MoreIcon:                 decredmaterial.NewImage(decredIcons["more"]),
		MoreIconInactive:         decredmaterial.NewImage(decredIcons["more_inactive"]),
		Logo:                     decredmaterial.NewImage(decredIcons["logo"]),
		ConfirmIcon:              decredmaterial.NewImage(decredIcons["confirmed"]),
		PendingIcon:              decredmaterial.NewImage(decredIcons["pending"]),
		RedirectIcon:             decredmaterial.NewImage(decredIcons["redirect"]),
		NewWalletIcon:            decredmaterial.NewImage(decredIcons["addNewWallet"]),
		WalletAlertIcon:          decredmaterial.NewImage(decredIcons["walletAlert"]),
		AccountIcon:              decredmaterial.NewImage(decredIcons["account"]),
		ImportedAccountIcon:      decredmaterial.NewImage(decredIcons["imported_account"]),
		EditIcon:                 decredmaterial.NewImage(decredIcons["editIcon"]),
		expandIcon:               decredmaterial.NewImage(decredIcons["expand_icon"]),
		CopyIcon:                 decredmaterial.NewImage(decredIcons["copy_icon"]),
		MixedTx:                  decredmaterial.NewImage(decredIcons["mixed_tx"]),
		Mixer:                    decredmaterial.NewImage(decredIcons["mixer"]),
		PrivacySetup:             decredmaterial.NewImage(decredIcons["privacy_setup"]),
		Next:                     decredmaterial.NewImage(decredIcons["ic_next"]),
		SettingsIcon:             decredmaterial.NewImage(decredIcons["settings"]),
		SecurityIcon:             decredmaterial.NewImage(decredIcons["security"]),
		HelpIcon:                 decredmaterial.NewImage(decredIcons["help_icon"]),
		AboutIcon:                decredmaterial.NewImage(decredIcons["about_icon"]),
		DebugIcon:                decredmaterial.NewImage(decredIcons["debug"]),
		VerifyMessageIcon:        decredmaterial.NewImage(decredIcons["verify_message"]),
		LocationPinIcon:          decredmaterial.NewImage(decredIcons["location_pin"]),
		AlertGray:                decredmaterial.NewImage(decredIcons["alert_gray"]),
		ArrowDownIcon:            decredmaterial.NewImage(decredIcons["arrow_down"]),
		WatchOnlyWalletIcon:      decredmaterial.NewImage(decredIcons["watch_only_wallet"]),
		CurrencySwapIcon:         decredmaterial.NewImage(decredIcons["swap"]),
		SyncingIcon:              decredmaterial.NewImage(decredIcons["syncing"]),
		DocumentationIcon:        decredmaterial.NewImage(decredIcons["documentation"]),
		ProposalIconActive:       decredmaterial.NewImage(decredIcons["politeiaActive"]),
		ProposalIconInactive:     decredmaterial.NewImage(decredIcons["politeiaInactive"]),
		Restore:                  decredmaterial.NewImage(decredIcons["restore"]),
		DownloadIcon:             decredmaterial.NewImage(decredIcons["downloadIcon"]),
		TimerIcon:                decredmaterial.NewImage(decredIcons["timerIcon"]),
		WalletIcon:               decredmaterial.NewImage(decredIcons["wallet"]),
		StakeIcon:                decredmaterial.NewImage(decredIcons["stake"]),
		StakeIconInactive:        decredmaterial.NewImage(decredIcons["stake_inactive"]),
		StakeyIcon:               decredmaterial.NewImage(decredIcons["stakey"]),
		NewStakeIcon:             decredmaterial.NewImage(decredIcons["stake_purchased"]),
		TicketImmatureIcon:       decredmaterial.NewImage(decredIcons["ticket_immature"]),
		TicketUnminedIcon:        decredmaterial.NewImage(decredIcons["ticket_unmined"]),
		TicketLiveIcon:           decredmaterial.NewImage(decredIcons["ticket_live"]),
		TicketVotedIcon:          decredmaterial.NewImage(decredIcons["ticket_voted"]),
		TicketMissedIcon:         decredmaterial.NewImage(decredIcons["ticket_missed"]),
		TicketExpiredIcon:        decredmaterial.NewImage(decredIcons["ticket_expired"]),
		TicketRevokedIcon:        decredmaterial.NewImage(decredIcons["ticket_revoked"]),
		List:                     decredmaterial.NewImage(decredIcons["list"]),
		ListGridIcon:             decredmaterial.NewImage(decredIcons["list_grid"]),
		DecredSymbolIcon:         decredmaterial.NewImage(decredIcons["decred_symbol"]),
		DecredSymbol2:            decredmaterial.NewImage(decredIcons["ic_decred02"]),
		GovernanceActiveIcon:     decredmaterial.NewImage(decredIcons["governance_active"]),
		GovernanceInactiveIcon:   decredmaterial.NewImage(decredIcons["governance_inactive"]),
		LogoDarkMode:             decredmaterial.NewImage(decredIcons["logo_darkmode"]),
		TimerDarkMode:            decredmaterial.NewImage(decredIcons["timer_dm"]),
		Rebroadcast:              decredmaterial.NewImage(decredIcons["rebroadcast"]),

		SettingsActiveIcon:       decredmaterial.NewImage(decredIcons["settings_active"]),
		SettingsInactiveIcon:     decredmaterial.NewImage(decredIcons["settings_inactive"]),
		ActivatedActiveIcon:      decredmaterial.NewImage(decredIcons["activated_active"]),
		ActivatedInactiveIcon:    decredmaterial.NewImage(decredIcons["activated_inactive"]),
		LockinActiveIcon:         decredmaterial.NewImage(decredIcons["lockin_active"]),
		LockinInactiveIcon:       decredmaterial.NewImage(decredIcons["lockin_inactive"]),

		DexIcon:         decredmaterial.NewImage(decredIcons["dex_icon"]),
		DexIconInactive: decredmaterial.NewImage(decredIcons["dex_icon_inactive"]),
		BTC:             decredmaterial.NewImage(decredIcons["dex_btc"]),
		DCR:             decredmaterial.NewImage(decredIcons["dex_dcr"]),
	}
	return ic
}

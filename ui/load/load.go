F// The load package contains data structures that are shared by components in the ui package. It is not a dumping ground
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
	LogoDarkMode, TimerDarkMode, Rebroadcast, SettingsActiveIcon, SettingsInactiveIcon, ActivatedActiveIcon,
	ActivatedInactiveIcon, LockinActiveIcon, LockinInactiveIcon *decredmaterial.Image

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

	DarkModeSettingChanged func(bool)
}

func (l *Load) RefreshTheme() {
	isDarkModeOn := l.WL.MultiWallet.ReadBoolConfigValueForKey(DarkModeConfigKey, false)
	l.Theme.SwitchDarkMode(isDarkModeOn, assets.DecredIcons)
	l.DarkModeSettingChanged(isDarkModeOn)
	l.RefreshWindow()
}

func (l *Load) Dexc() *dcrlibwallet.DexClient {
	return l.WL.MultiWallet.DexClient()
}

func IconSet(isDarkModeOn bool) Icons {
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
	}

	defaultIcons(&ic)
	// if isDarkModeOn{
	// 	darkModeIcons(&ic)
	// }

	return ic
}

func defaultIcons(ic *Icons) {
	decredIcons := assets.DecredIcons

	ic.OverviewIcon = decredmaterial.NewImage(decredIcons["overview"])
	ic.OverviewIconInactive = decredmaterial.NewImage(decredIcons["overview_inactive"])
	ic.WalletIconInactive = decredmaterial.NewImage(decredIcons["wallet_inactive"])
	ic.ReceiveIcon = decredmaterial.NewImage(decredIcons["receive"])
	ic.Transferred = decredmaterial.NewImage(decredIcons["transferred"])
	ic.TransactionsIcon = decredmaterial.NewImage(decredIcons["transactions"])
	ic.TransactionsIconInactive = decredmaterial.NewImage(decredIcons["transactions_inactive"])
	ic.SendIcon = decredmaterial.NewImage(decredIcons["send"])
	ic.MoreIcon = decredmaterial.NewImage(decredIcons["more"])
	ic.MoreIconInactive = decredmaterial.NewImage(decredIcons["more_inactive"])
	ic.Logo = decredmaterial.NewImage(decredIcons["logo"])
	ic.ConfirmIcon = decredmaterial.NewImage(decredIcons["confirmed"])
	ic.PendingIcon = decredmaterial.NewImage(decredIcons["pending"])
	ic.RedirectIcon = decredmaterial.NewImage(decredIcons["redirect"])
	ic.NewWalletIcon = decredmaterial.NewImage(decredIcons["addNewWallet"])
	ic.WalletAlertIcon = decredmaterial.NewImage(decredIcons["walletAlert"])
	ic.AccountIcon = decredmaterial.NewImage(decredIcons["account"])
	ic.ImportedAccountIcon = decredmaterial.NewImage(decredIcons["imported_account"])
	ic.EditIcon = decredmaterial.NewImage(decredIcons["editIcon"])
	ic.expandIcon = decredmaterial.NewImage(decredIcons["expand_icon"])
	ic.CopyIcon = decredmaterial.NewImage(decredIcons["copy_icon"])
	ic.MixedTx = decredmaterial.NewImage(decredIcons["mixed_tx"])
	ic.Mixer = decredmaterial.NewImage(decredIcons["mixer"])
	ic.PrivacySetup = decredmaterial.NewImage(decredIcons["privacy_setup"])
	ic.Next = decredmaterial.NewImage(decredIcons["ic_next"])
	ic.SettingsIcon = decredmaterial.NewImage(decredIcons["settings"])
	ic.SecurityIcon = decredmaterial.NewImage(decredIcons["security"])
	ic.HelpIcon = decredmaterial.NewImage(decredIcons["help_icon"])
	ic.AboutIcon = decredmaterial.NewImage(decredIcons["about_icon"])
	ic.DebugIcon = decredmaterial.NewImage(decredIcons["debug"])
	ic.VerifyMessageIcon = decredmaterial.NewImage(decredIcons["verify_message"])
	ic.LocationPinIcon = decredmaterial.NewImage(decredIcons["location_pin"])
	ic.AlertGray = decredmaterial.NewImage(decredIcons["alert_gray"])
	ic.ArrowDownIcon = decredmaterial.NewImage(decredIcons["arrow_down"])
	ic.WatchOnlyWalletIcon = decredmaterial.NewImage(decredIcons["watch_only_wallet"])
	ic.CurrencySwapIcon = decredmaterial.NewImage(decredIcons["swap"])
	ic.SyncingIcon = decredmaterial.NewImage(decredIcons["syncing"])
	ic.DocumentationIcon = decredmaterial.NewImage(decredIcons["documentation"])
	ic.ProposalIconActive = decredmaterial.NewImage(decredIcons["politeiaActive"])
	ic.ProposalIconInactive = decredmaterial.NewImage(decredIcons["politeiaInactive"])
	ic.Restore = decredmaterial.NewImage(decredIcons["restore"])
	ic.DownloadIcon = decredmaterial.NewImage(decredIcons["downloadIcon"])
	ic.TimerIcon = decredmaterial.NewImage(decredIcons["timerIcon"])
	ic.WalletIcon = decredmaterial.NewImage(decredIcons["wallet"])
	ic.StakeIcon = decredmaterial.NewImage(decredIcons["stake"])
	ic.StakeIconInactive = decredmaterial.NewImage(decredIcons["stake_inactive"])
	ic.StakeyIcon = decredmaterial.NewImage(decredIcons["stakey"])
	ic.NewStakeIcon = decredmaterial.NewImage(decredIcons["stake_purchased"])
	ic.TicketImmatureIcon = decredmaterial.NewImage(decredIcons["ticket_immature"])
	ic.TicketUnminedIcon = decredmaterial.NewImage(decredIcons["ticket_unmined"])
	ic.TicketLiveIcon = decredmaterial.NewImage(decredIcons["ticket_live"])
	ic.TicketVotedIcon = decredmaterial.NewImage(decredIcons["ticket_voted"])
	ic.TicketMissedIcon = decredmaterial.NewImage(decredIcons["ticket_missed"])
	ic.TicketExpiredIcon = decredmaterial.NewImage(decredIcons["ticket_expired"])
	ic.TicketRevokedIcon = decredmaterial.NewImage(decredIcons["ticket_revoked"])
	ic.List = decredmaterial.NewImage(decredIcons["list"])
	ic.ListGridIcon = decredmaterial.NewImage(decredIcons["list_grid"])
	ic.DecredSymbolIcon = decredmaterial.NewImage(decredIcons["decred_symbol"])
	ic.DecredSymbol2 = decredmaterial.NewImage(decredIcons["ic_decred02"])
	ic.GovernanceActiveIcon = decredmaterial.NewImage(decredIcons["governance_active"])
	ic.GovernanceInactiveIcon = decredmaterial.NewImage(decredIcons["governance_inactive"])
	ic.LogoDarkMode = decredmaterial.NewImage(decredIcons["logo_darkmode"])
	ic.TimerDarkMode = decredmaterial.NewImage(decredIcons["timer_dm"])
	ic.Rebroadcast = decredmaterial.NewImage(decredIcons["rebroadcast"])

	ic.SettingsActiveIcon = decredmaterial.NewImage(decredIcons["settings_active"])
	ic.SettingsInactiveIcon = decredmaterial.NewImage(decredIcons["settings_inactive"])
	ic.ActivatedActiveIcon = decredmaterial.NewImage(decredIcons["activated_active"])
	ic.ActivatedInactiveIcon = decredmaterial.NewImage(decredIcons["activated_inactive"])
	ic.LockinActiveIcon = decredmaterial.NewImage(decredIcons["lockin_active"])
	ic.LockinInactiveIcon = decredmaterial.NewImage(decredIcons["lockin_inactive"])

	ic.DexIcon = decredmaterial.NewImage(decredIcons["dex_icon"])
	ic.DexIconInactive = decredmaterial.NewImage(decredIcons["dex_icon_inactive"])
	ic.BTC = decredmaterial.NewImage(decredIcons["dex_btc"])
	ic.DCR = decredmaterial.NewImage(decredIcons["dex_dcr"])
}

func darkModeIcons(ic *Icons) {
	decredIcons := assets.DecredIcons

	ic.OverviewIcon = decredmaterial.NewImage(decredIcons["overview"])
	ic.OverviewIconInactive = decredmaterial.NewImage(decredIcons["overview_inactive"])
}

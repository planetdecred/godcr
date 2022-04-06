package decredmaterial

import (
	"gioui.org/widget"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/planetdecred/godcr/ui/assets"
)

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
	ActivatedInactiveIcon, LockinActiveIcon, LockinInactiveIcon *Image

	NewStakeIcon,
	TicketImmatureIcon,
	TicketLiveIcon,
	TicketVotedIcon,
	TicketMissedIcon,
	TicketExpiredIcon,
	TicketRevokedIcon,
	TicketUnminedIcon *Image

	DexIcon, DexIconInactive, BTC, DCR *Image
}

func (i *Icons) StandardMaterialIcons() *Icons {
	i.ContentAdd = MustIcon(widget.NewIcon(icons.ContentAdd))
	i.NavigationCheck = MustIcon(widget.NewIcon(icons.NavigationCheck))
	i.NavigationMore = MustIcon(widget.NewIcon(icons.NavigationMoreHoriz))
	i.ActionCheckCircle = MustIcon(widget.NewIcon(icons.ActionCheckCircle))
	i.NavigationArrowBack = MustIcon(widget.NewIcon(icons.NavigationArrowBack))
	i.NavigationArrowForward = MustIcon(widget.NewIcon(icons.NavigationArrowForward))
	i.ActionInfo = MustIcon(widget.NewIcon(icons.ActionInfo))
	i.ActionCheck = MustIcon(widget.NewIcon(icons.ActionCheckCircle))
	i.NavigationCancel = MustIcon(widget.NewIcon(icons.NavigationCancel))
	i.ImageBrightness1 = MustIcon(widget.NewIcon(icons.ImageBrightness1))
	i.ChevronRight = MustIcon(widget.NewIcon(icons.NavigationChevronRight))
	i.ContentClear = MustIcon(widget.NewIcon(icons.ContentClear))
	i.NavMoreIcon = MustIcon(widget.NewIcon(icons.NavigationMoreHoriz))
	i.DropDownIcon = MustIcon(widget.NewIcon(icons.NavigationArrowDropDown))
	i.Cached = MustIcon(widget.NewIcon(icons.ActionCached))
	i.ContentRemove = MustIcon(widget.NewIcon(icons.ContentRemove))
	i.ConcealIcon = MustIcon(widget.NewIcon(icons.ActionVisibility))
	i.RevealIcon = MustIcon(widget.NewIcon(icons.ActionVisibilityOff))
	i.SearchIcon = MustIcon(widget.NewIcon(icons.ActionSearch))
	i.PlayIcon = MustIcon(widget.NewIcon(icons.AVPlayArrow))

	return i
}

func (i *Icons) DefaultIcons() *Icons {
	decredIcons := assets.DecredIcons

	i.StandardMaterialIcons()
	i.OverviewIcon = NewImage(decredIcons["overview"])
	i.OverviewIconInactive = NewImage(decredIcons["overview_inactive"])
	i.WalletIconInactive = NewImage(decredIcons["wallet_inactive"])
	i.ReceiveIcon = NewImage(decredIcons["receive"])
	i.Transferred = NewImage(decredIcons["transferred"])
	i.TransactionsIcon = NewImage(decredIcons["transactions"])
	i.TransactionsIconInactive = NewImage(decredIcons["transactions_inactive"])
	i.SendIcon = NewImage(decredIcons["send"])
	i.MoreIcon = NewImage(decredIcons["more"])
	i.MoreIconInactive = NewImage(decredIcons["more_inactive"])
	i.Logo = NewImage(decredIcons["logo"])
	i.ConfirmIcon = NewImage(decredIcons["confirmed"])
	i.PendingIcon = NewImage(decredIcons["pending"])
	i.RedirectIcon = NewImage(decredIcons["redirect"])
	i.NewWalletIcon = NewImage(decredIcons["addNewWallet"])
	i.WalletAlertIcon = NewImage(decredIcons["walletAlert"])
	i.AccountIcon = NewImage(decredIcons["account"])
	i.ImportedAccountIcon = NewImage(decredIcons["imported_account"])
	i.EditIcon = NewImage(decredIcons["editIcon"])
	i.expandIcon = NewImage(decredIcons["expand_icon"])
	i.CopyIcon = NewImage(decredIcons["copy_icon"])
	i.MixedTx = NewImage(decredIcons["mixed_tx"])
	i.Mixer = NewImage(decredIcons["mixer"])
	i.PrivacySetup = NewImage(decredIcons["privacy_setup"])
	i.Next = NewImage(decredIcons["ic_next"])
	i.SettingsIcon = NewImage(decredIcons["settings"])
	i.SecurityIcon = NewImage(decredIcons["security"])
	i.HelpIcon = NewImage(decredIcons["help_icon"])
	i.AboutIcon = NewImage(decredIcons["about_icon"])
	i.DebugIcon = NewImage(decredIcons["debug"])
	i.VerifyMessageIcon = NewImage(decredIcons["verify_message"])
	i.LocationPinIcon = NewImage(decredIcons["location_pin"])
	i.AlertGray = NewImage(decredIcons["alert_gray"])
	i.ArrowDownIcon = NewImage(decredIcons["arrow_down"])
	i.WatchOnlyWalletIcon = NewImage(decredIcons["watch_only_wallet"])
	i.CurrencySwapIcon = NewImage(decredIcons["swap"])
	i.SyncingIcon = NewImage(decredIcons["syncing"])
	i.DocumentationIcon = NewImage(decredIcons["documentation"])
	i.ProposalIconActive = NewImage(decredIcons["politeiaActive"])
	i.ProposalIconInactive = NewImage(decredIcons["politeiaInactive"])
	i.Restore = NewImage(decredIcons["restore"])
	i.DownloadIcon = NewImage(decredIcons["downloadIcon"])
	i.TimerIcon = NewImage(decredIcons["timerIcon"])
	i.WalletIcon = NewImage(decredIcons["wallet"])
	i.StakeIcon = NewImage(decredIcons["stake"])
	i.StakeIconInactive = NewImage(decredIcons["stake_inactive"])
	i.StakeyIcon = NewImage(decredIcons["stakey"])
	i.NewStakeIcon = NewImage(decredIcons["stake_purchased"])
	i.TicketImmatureIcon = NewImage(decredIcons["ticket_immature"])
	i.TicketUnminedIcon = NewImage(decredIcons["ticket_unmined"])
	i.TicketLiveIcon = NewImage(decredIcons["ticket_live"])
	i.TicketVotedIcon = NewImage(decredIcons["ticket_voted"])
	i.TicketMissedIcon = NewImage(decredIcons["ticket_missed"])
	i.TicketExpiredIcon = NewImage(decredIcons["ticket_expired"])
	i.TicketRevokedIcon = NewImage(decredIcons["ticket_revoked"])
	i.List = NewImage(decredIcons["list"])
	i.ListGridIcon = NewImage(decredIcons["list_grid"])
	i.DecredSymbolIcon = NewImage(decredIcons["decred_symbol"])
	i.DecredSymbol2 = NewImage(decredIcons["ic_decred02"])
	i.GovernanceActiveIcon = NewImage(decredIcons["governance_active"])
	i.GovernanceInactiveIcon = NewImage(decredIcons["governance_inactive"])
	i.TimerDarkMode = NewImage(decredIcons["timer_dm"])
	i.Rebroadcast = NewImage(decredIcons["rebroadcast"])

	i.SettingsActiveIcon = NewImage(decredIcons["settings_active"])
	i.SettingsInactiveIcon = NewImage(decredIcons["settings_inactive"])
	i.ActivatedActiveIcon = NewImage(decredIcons["activated_active"])
	i.ActivatedInactiveIcon = NewImage(decredIcons["activated_inactive"])
	i.LockinActiveIcon = NewImage(decredIcons["lockin_active"])
	i.LockinInactiveIcon = NewImage(decredIcons["lockin_inactive"])

	i.DexIcon = NewImage(decredIcons["dex_icon"])
	i.DexIconInactive = NewImage(decredIcons["dex_icon_inactive"])
	i.BTC = NewImage(decredIcons["dex_btc"])
	i.DCR = NewImage(decredIcons["dex_dcr"])

	return i
}

func (i *Icons) DarkModeIcons() *Icons {
	decredIcons := assets.DecredIcons

	i.OverviewIcon = NewImage(decredIcons["overview"])
	i.OverviewIconInactive = NewImage(decredIcons["overview_inactive"])
	i.WalletIconInactive = NewImage(decredIcons["wallet_inactive"])
	i.ReceiveIcon = NewImage(decredIcons["receive"])
	i.Transferred = NewImage(decredIcons["transferred"])
	i.TransactionsIcon = NewImage(decredIcons["transactions"])
	i.TransactionsIconInactive = NewImage(decredIcons["transactions_inactive"])
	i.SendIcon = NewImage(decredIcons["send"])
	i.MoreIcon = NewImage(decredIcons["more"])
	i.MoreIconInactive = NewImage(decredIcons["more_inactive"])
	i.Logo = NewImage(decredIcons["logo_darkmode"])
	i.ConfirmIcon = NewImage(decredIcons["confirmed"])
	i.PendingIcon = NewImage(decredIcons["pending"])
	i.RedirectIcon = NewImage(decredIcons["redirect"])
	i.NewWalletIcon = NewImage(decredIcons["addNewWallet"])
	i.WalletAlertIcon = NewImage(decredIcons["walletAlert"])
	i.AccountIcon = NewImage(decredIcons["account"])
	i.ImportedAccountIcon = NewImage(decredIcons["imported_account"])
	i.EditIcon = NewImage(decredIcons["editIcon"])
	i.expandIcon = NewImage(decredIcons["expand_icon"])
	i.CopyIcon = NewImage(decredIcons["copy_icon"])
	i.MixedTx = NewImage(decredIcons["mixed_tx"])
	i.Mixer = NewImage(decredIcons["mixer"])
	i.PrivacySetup = NewImage(decredIcons["privacy_setup"])
	i.Next = NewImage(decredIcons["ic_next"])
	i.SettingsIcon = NewImage(decredIcons["settings"])
	i.SecurityIcon = NewImage(decredIcons["security"])
	i.HelpIcon = NewImage(decredIcons["help_icon"])
	i.AboutIcon = NewImage(decredIcons["about_icon"])
	i.DebugIcon = NewImage(decredIcons["debug"])
	i.VerifyMessageIcon = NewImage(decredIcons["verify_message"])
	i.LocationPinIcon = NewImage(decredIcons["location_pin"])
	i.AlertGray = NewImage(decredIcons["alert_gray"])
	i.ArrowDownIcon = NewImage(decredIcons["arrow_down"])
	i.WatchOnlyWalletIcon = NewImage(decredIcons["watch_only_wallet"])
	i.CurrencySwapIcon = NewImage(decredIcons["swap"])
	i.SyncingIcon = NewImage(decredIcons["syncing"])
	i.DocumentationIcon = NewImage(decredIcons["documentation"])
	i.ProposalIconActive = NewImage(decredIcons["politeiaActive"])
	i.ProposalIconInactive = NewImage(decredIcons["politeiaInactive"])
	i.Restore = NewImage(decredIcons["restore"])
	i.DownloadIcon = NewImage(decredIcons["downloadIcon"])
	i.TimerIcon = NewImage(decredIcons["timerIcon"])
	i.WalletIcon = NewImage(decredIcons["wallet"])
	i.StakeIcon = NewImage(decredIcons["stake"])
	i.StakeIconInactive = NewImage(decredIcons["stake_inactive"])
	i.StakeyIcon = NewImage(decredIcons["stakey"])
	i.NewStakeIcon = NewImage(decredIcons["stake_purchased"])
	i.TicketImmatureIcon = NewImage(decredIcons["ticket_immature"])
	i.TicketUnminedIcon = NewImage(decredIcons["ticket_unmined"])
	i.TicketLiveIcon = NewImage(decredIcons["ticket_live"])
	i.TicketVotedIcon = NewImage(decredIcons["ticket_voted"])
	i.TicketMissedIcon = NewImage(decredIcons["ticket_missed"])
	i.TicketExpiredIcon = NewImage(decredIcons["ticket_expired"])
	i.TicketRevokedIcon = NewImage(decredIcons["ticket_revoked"])
	i.List = NewImage(decredIcons["list"])
	i.ListGridIcon = NewImage(decredIcons["list_grid"])
	i.DecredSymbolIcon = NewImage(decredIcons["decred_symbol"])
	i.DecredSymbol2 = NewImage(decredIcons["ic_decred02"])
	i.GovernanceActiveIcon = NewImage(decredIcons["governance_active"])
	i.GovernanceInactiveIcon = NewImage(decredIcons["governance_inactive"])
	i.TimerDarkMode = NewImage(decredIcons["timer_dm"])
	i.Rebroadcast = NewImage(decredIcons["rebroadcast"])

	i.SettingsActiveIcon = NewImage(decredIcons["settings_active"])
	i.SettingsInactiveIcon = NewImage(decredIcons["settings_inactive"])
	i.ActivatedActiveIcon = NewImage(decredIcons["activated_active"])
	i.ActivatedInactiveIcon = NewImage(decredIcons["activated_inactive"])
	i.LockinActiveIcon = NewImage(decredIcons["lockin_active"])
	i.LockinInactiveIcon = NewImage(decredIcons["lockin_inactive"])

	i.DexIcon = NewImage(decredIcons["dex_icon"])
	i.DexIconInactive = NewImage(decredIcons["dex_icon_inactive"])
	i.BTC = NewImage(decredIcons["dex_btc"])
	i.DCR = NewImage(decredIcons["dex_dcr"])

	return i
}

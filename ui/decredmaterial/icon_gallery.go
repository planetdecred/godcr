package decredmaterial

import (
	"gioui.org/widget"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/planetdecred/godcr/ui/assets"
)

type Icons struct {
	ContentAdd, NavigationCheck, NavigationMore, ActionCheckCircle, ActionInfo, NavigationArrowBack,
	NavigationArrowForward, ActionCheck, ChevronRight, NavigationCancel, NavMoreIcon,
	ImageBrightness1, ContentClear, DropDownIcon, Cached, ContentRemove, SearchIcon, PlayIcon *widget.Icon

	OverviewIcon, OverviewIconInactive, WalletIcon, WalletIconInactive, MixerInactive, RedAlert,
	ReceiveIcon, Transferred, TransactionsIcon, TransactionsIconInactive, SendIcon, MoreIcon, MoreIconInactive,
	PendingIcon, Logo, RedirectIcon, ConfirmIcon, NewWalletIcon, WalletAlertIcon, ArrowForward,
	ImportedAccountIcon, AccountIcon, EditIcon, expandIcon, CopyIcon, MixedTx, Mixer, DcrWatchOnly,
	Next, SettingsIcon, SecurityIcon, HelpIcon, AboutIcon, DebugIcon, VerifyMessageIcon, LocationPinIcon, SignMessageIcon,
	HeaderSettingsIcon, AlertGray, ArrowDownIcon, WatchOnlyWalletIcon, CurrencySwapIcon, SyncingIcon, TransactionFingerprint,
	Restore, DocumentationIcon, TimerIcon, StakeIcon, StakeIconInactive, StakeyIcon, DecredLogo,
	DecredSymbol2, GovernanceActiveIcon, GovernanceInactiveIcon, LogoDarkMode, TimerDarkMode, Rebroadcast,
	SettingsActiveIcon, SettingsInactiveIcon, ActivatedActiveIcon, ActivatedInactiveIcon, LockinActiveIcon,
	LockinInactiveIcon, SuccessIcon, FailedIcon, ReceiveInactiveIcon, SendInactiveIcon, DarkmodeIcon,
	ChevronExpand, ChevronCollapse, ChevronLeft, MixedTxIcon, UnmixedTxIcon, MixerIcon, NotSynced *Image

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
	i.Next = NewImage(decredIcons["ic_next"])
	i.HeaderSettingsIcon = NewImage(decredIcons["header_settings"])
	i.SettingsIcon = NewImage(decredIcons["settings"])
	i.SecurityIcon = NewImage(decredIcons["security"])
	i.HelpIcon = NewImage(decredIcons["help_icon"])
	i.AboutIcon = NewImage(decredIcons["about_icon"])
	i.DebugIcon = NewImage(decredIcons["debug"])
	i.VerifyMessageIcon = NewImage(decredIcons["verify_message"])
	i.LocationPinIcon = NewImage(decredIcons["location_pin"])
	i.SignMessageIcon = NewImage(decredIcons["signMessage"])
	i.AlertGray = NewImage(decredIcons["alert_gray"])
	i.ArrowDownIcon = NewImage(decredIcons["arrow_down"])
	i.WatchOnlyWalletIcon = NewImage(decredIcons["watch_only_wallet"])
	i.CurrencySwapIcon = NewImage(decredIcons["swap"])
	i.SyncingIcon = NewImage(decredIcons["syncing"])
	i.DocumentationIcon = NewImage(decredIcons["documentation"])
	i.Restore = NewImage(decredIcons["restore"])
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
	i.DecredLogo = NewImage(decredIcons["decred_symbol"])
	i.DecredSymbol2 = NewImage(decredIcons["ic_decred02"])
	i.GovernanceActiveIcon = NewImage(decredIcons["governance_active"])
	i.GovernanceInactiveIcon = NewImage(decredIcons["governance_inactive"])
	i.Rebroadcast = NewImage(decredIcons["rebroadcast"])
	i.ConcealIcon = NewImage(decredIcons["reveal"])
	i.RevealIcon = NewImage(decredIcons["hide"])

	i.SettingsActiveIcon = NewImage(decredIcons["settings_active"])
	i.SettingsInactiveIcon = NewImage(decredIcons["settings_inactive"])
	i.ActivatedActiveIcon = NewImage(decredIcons["activated_active"])
	i.ActivatedInactiveIcon = NewImage(decredIcons["activated_inactive"])
	i.LockinActiveIcon = NewImage(decredIcons["lockin_active"])
	i.LockinInactiveIcon = NewImage(decredIcons["lockin_inactive"])
	i.TransactionFingerprint = NewImage(decredIcons["transaction_fingerprint"])
	i.ArrowForward = NewImage(decredIcons["arrow_fwd"])

	i.DexIcon = NewImage(decredIcons["dex_icon"])
	i.DexIconInactive = NewImage(decredIcons["dex_icon_inactive"])
	i.BTC = NewImage(decredIcons["dex_btc"])
	i.DCR = NewImage(decredIcons["dex_dcr"])
	i.SuccessIcon = NewImage(decredIcons["success_check"])
	i.FailedIcon = NewImage(decredIcons["crossmark_red"])
	i.ReceiveInactiveIcon = NewImage(decredIcons["receive_inactive"])
	i.SendInactiveIcon = NewImage(decredIcons["send_inactive"])
	i.DarkmodeIcon = NewImage(decredIcons["darkmodeIcon"])
	i.MixerInactive = NewImage(decredIcons["mixer_inactive"])
	i.DcrWatchOnly = NewImage(decredIcons["dcr_watch_only"])
	i.RedAlert = NewImage(decredIcons["red_alert"])
	i.ChevronExpand = NewImage(decredIcons["chevron_coll"])
	i.ChevronCollapse = NewImage(decredIcons["chevron_expand"])
	i.ChevronLeft = NewImage(decredIcons["chevron_left"])
	i.NotSynced = NewImage(decredIcons["notSynced"])
	i.UnmixedTxIcon = NewImage(decredIcons["unmixed_icon"])
	i.MixedTxIcon = NewImage(decredIcons["mixed_icon"])
	i.MixerIcon = NewImage(decredIcons["mixer_icon"])
	i.InfoAction = NewImage(decredIcons["info_icon"])
	i.DarkMode = NewImage(decredIcons["ic_moon"])
	i.LightMode = NewImage(decredIcons["ic_sun"])
	i.AddIcon = NewImage(decredIcons["addIcon"])

	return i
}

func (i *Icons) DarkModeIcons() *Icons {
	decredIcons := assets.DecredIcons

	i.OverviewIcon = NewImage(decredIcons["dm_overview"])
	i.OverviewIconInactive = NewImage(decredIcons["dm_overview_inactive"])
	i.WalletIconInactive = NewImage(decredIcons["dm_wallet_inactive"])
	i.TransactionsIcon = NewImage(decredIcons["dm_transactions"])
	i.TransactionsIconInactive = NewImage(decredIcons["dm_transactions_inactive"])
	i.MoreIcon = NewImage(decredIcons["dm_more"])
	i.MoreIconInactive = NewImage(decredIcons["dm_more_inactive"])
	i.Logo = NewImage(decredIcons["logo_darkmode"])
	i.RedirectIcon = NewImage(decredIcons["dm_redirect"])
	i.NewWalletIcon = NewImage(decredIcons["dm_addNewWallet"])
	i.WalletAlertIcon = NewImage(decredIcons["dm_walletAlert"])
	i.AccountIcon = NewImage(decredIcons["dm_account"])
	i.ImportedAccountIcon = NewImage(decredIcons["dm_imported_account"])
	i.EditIcon = NewImage(decredIcons["dm_editIcon"])
	i.CopyIcon = NewImage(decredIcons["dm_copy_icon"])
	i.Mixer = NewImage(decredIcons["dm_mixer"])
	i.Next = NewImage(decredIcons["dm_ic_next"])
	i.SettingsIcon = NewImage(decredIcons["dm_settings"])
	i.SecurityIcon = NewImage(decredIcons["dm_security"])
	i.HelpIcon = NewImage(decredIcons["dm_help_icon"])
	i.AboutIcon = NewImage(decredIcons["dm_info_icon"])
	i.DebugIcon = NewImage(decredIcons["dm_debug"])
	i.VerifyMessageIcon = NewImage(decredIcons["dm_verify_message"])
	i.LocationPinIcon = NewImage(decredIcons["dm_location_pin"])
	i.ArrowDownIcon = NewImage(decredIcons["dm_arrow_down"])
	i.WatchOnlyWalletIcon = NewImage(decredIcons["dm_watch_only_wallet"])
	i.CurrencySwapIcon = NewImage(decredIcons["dm_swap"])
	i.Restore = NewImage(decredIcons["dm_restore"])
	i.TimerIcon = NewImage(decredIcons["dm_timerIcon"])
	i.WalletIcon = NewImage(decredIcons["dm_wallet"])
	i.StakeIcon = NewImage(decredIcons["dm_stake"])
	i.TicketRevokedIcon = NewImage(decredIcons["dm_ticket_revoked"])
	i.DecredLogo = NewImage(decredIcons["dm_decred_symbol"])
	i.GovernanceActiveIcon = NewImage(decredIcons["dm_governance_active"])
	i.GovernanceInactiveIcon = NewImage(decredIcons["dm_governance_inactive"])
	i.Rebroadcast = NewImage(decredIcons["dm_rebroadcast"])
	i.ActivatedActiveIcon = NewImage(decredIcons["dm_activated_active"])
	i.LockinActiveIcon = NewImage(decredIcons["dm_lockin_active"])
	i.DexIcon = NewImage(decredIcons["dm_dex_icon"])
	i.TransactionFingerprint = NewImage(decredIcons["dm_transaction_fingerprint"])
	i.ArrowForward = NewImage(decredIcons["dm_arrow_fwd"])
	i.ChevronLeft = NewImage(decredIcons["chevron_left"])
	return i
}

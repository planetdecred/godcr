package load

const Uint32Size = 32 << (^uint32(0) >> 32 & 1) // 32 or 64
const MaxInt32 = 1<<(Uint32Size-1) - 1

const (
	// godcr config keys
<<<<<<< HEAD
	HideBalanceConfigKey             = "hide_balance"
	AutoSyncConfigKey                = "autoSync"
	LanguagePreferenceKey            = "app_language"
	DarkModeConfigKey                = "dark_mode"
	FetchProposalConfigKey           = "fetch_proposals"
	SeedBackupNotificationConfigKey  = "seed_backup_notification"
	ProposalNotificationConfigKey    = "proposal_notification_key"
	TransactionNotificationConfigKey = "transaction_notification_key"
=======
	HideBalanceConfigKey            = "hide_balance"
	AutoSyncConfigKey               = "autoSync"
	LanguagePreferenceKey           = "app_language"
	DarkModeConfigKey               = "dark_mode"
	FetchProposalConfigKey          = "fetch_proposals"
	SeedBackupNotificationConfigKey = "seed_backup_notification"
	ProposalNotificationConfigKey   = "proposal_notification_key"
<<<<<<< HEAD
	TicketBuyerConfigKey= "auto_ticekt_buyer"
>>>>>>> 5a26d0f (temp)
=======
	TicketBuyerConfigKey            = "auto_ticekt_buyer"
>>>>>>> be2c6bd (fix balance to mentain bug)
)

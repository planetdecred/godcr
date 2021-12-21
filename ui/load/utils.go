package load

const Uint32Size = 32 << (^uint32(0) >> 32 & 1) // 32 or 64
const MaxInt32 = 1<<(Uint32Size-1) - 1

const (
	// godcr config keys
	HideBalanceConfigKey            = "hide_balance"
	AutoSyncConfigKey               = "autoSync"
	LanguagePreferenceKey           = "app_language"
	DarkModeConfigKey               = "dark_mode"
	FetchProposalConfigKey          = "fetch_proposals"
	SeedBackupNotificationConfigKey = "seed_backup_notification"
	ProposalNotificationConfigKey   = "proposal_notification_key"
)

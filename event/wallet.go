package event

// WalletCmd represents commands sent to the wallet
type WalletCmd struct {
	Cmd       string
	Arguments *ArgumentQueue
}

// WalletResponse represents responses sent from the wallet
type WalletResponse struct {
	Resp    string
	Results *ArgumentQueue
}

// Sync represents sync events
type Sync struct {
	Event   string
	Payload interface{}
}

const (
	// ShutdownCmd tells the back end to clean up any operations then shutdown
	ShutdownCmd = "shutdown"
	// CreateCmd tells the wallet to create a new wallet given the provided passphrase and pass type
	CreateCmd = "create"
	// RestoreCmd tells the back end to restore the a wallet from the Payload string
	RestoreCmd = "restore"
	// InfoCmd tells the wallet to send back bulk info about the wallets
	InfoCmd = "infocmd"
	// StartSyncCmd tells the wallet to start the sync process
	StartSyncCmd = "startsync"
	// CreateTxCmd tells the wallet to create a new transaction
	CreateTxCmd = "createtx"

	// LoadedWalletsResp is the response for LoadedWalletsCmd
	LoadedWalletsResp = "loaded"
	// CreatedResp is the response returned when a new wallet has been created successfully
	CreatedResp = "created"
	// RestoredResp is the response returned when a wallet has been restored successfully
	RestoredResp = "restored"
	// CreatedTxResp is the response returned when a new transaction has been created
	CreatedTxResp = "createdTx"

	// SyncStart is the sync event sent when sync starts
	SyncStart = "syncstart"
	// SyncEnd is the sync event sent when sync ends
	SyncEnd = "syncend"
	// SyncProgress is the sync event sent when the wallet makes any sync progress
	SyncProgress = "syncprogress"
	// SyncCanceled is the sync event sent when the sync is canceled
	SyncCanceled = "synccanceled"
	// SyncError is the sync event sent when the sync ends with an error
	SyncError = "syncerror"
	// SyncPairsChanged is the sync event sent when peers are connected or disconnected
	SyncPairsChanged = "syncpairs"
)

// WalletInfo represents bulk information about the wallets returned by the wallet backend
type WalletInfo struct {
	LoadedWallets   int
	TotalBalance    int64
	BestBlockHeight int32
	BestBlockTime   int64
	Synced          bool
}

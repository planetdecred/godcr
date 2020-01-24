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

const (
	// ShutdownCmd tells the back end to clean up any operations then shutdown
	ShutdownCmd = "shutdown"
	// CreateCmd tells the wallet to create a new wallet given the provided passphrase and pass type
	CreateCmd = "create"
	// RestoreCmd tells the back end to restore the a wallet from the Payload string
	RestoreCmd = "restore"
	// LoadedWalletsCmd tells the wallet to send back the amount of loaded wallets
	LoadedWalletsCmd = "load"

	// LoadedWalletsResp is the response for LoadedWalletsCmd
	LoadedWalletsResp = "loaded"
	// CreatedResp is the response returned when a new wallet has been created successfully
	CreatedResp = "created"
	// RestoredResp is the response returned when a wallet has been restored successfully
	RestoredResp = "restored"

	// SyncStart is sync event send when sync starts
	SyncStart = "syncstart"
)

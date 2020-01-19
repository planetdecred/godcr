package event

// Loaded represents the event for the wallet finishing its loading sequence
type Loaded struct {
	WalletsLoadedCount int32
}

// WalletCmd represents commands sent to the wallet
type WalletCmd struct {
	Cmd string
}

const (
	ShutdownCmd = "shutdown"
)

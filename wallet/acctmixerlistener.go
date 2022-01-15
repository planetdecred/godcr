package wallet

type RunStatus int

const (
	MixerEnded RunStatus = iota
	MixerStarted
)

// AccountMixer is sent when account mixer started or ended.
type AccountMixer struct {
	WalletID  int
	RunStatus RunStatus
}

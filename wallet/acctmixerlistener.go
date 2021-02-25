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

func (l *listener) OnAccountMixerStarted(walletID int) {
	l.Send <- SyncStatusUpdate{
		Stage: AccountMixerStarted,
		AcctMixerInfo: AccountMixer{
			WalletID:  walletID,
			RunStatus: MixerStarted,
		},
	}
}

func (l *listener) OnAccountMixerEnded(walletID int) {
	l.Send <- SyncStatusUpdate{
		Stage: AccountMixerEnded,
		AcctMixerInfo: AccountMixer{
			WalletID:  walletID,
			RunStatus: MixerEnded,
		},
	}
}

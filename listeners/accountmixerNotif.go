package listeners

import "github.com/planetdecred/godcr/wallet"

// AccountMixerNotificationListener satisfies dcrlibwallet
// AccountMixerNotificationListener interface. Consumers interested in mixer
// notification instantiates this type.
type AccountMixerNotificationListener struct {
	MixerChan chan wallet.AccountMixer
}

func NewAccountMixerNotificationListener() *AccountMixerNotificationListener {
	return &AccountMixerNotificationListener{
		MixerChan: make(chan wallet.AccountMixer, 4),
	}
}

// OnAccountMixerStarted is a callback func called when the account mixer is
// started.
func (am *AccountMixerNotificationListener) OnAccountMixerStarted(walletID int) {
	am.UpdateNotification(wallet.AccountMixer{
		WalletID:  walletID,
		RunStatus: wallet.MixerStarted,
	})
}

// OnAccountMixerEnded is a callback func called when mixing ends.
func (am *AccountMixerNotificationListener) OnAccountMixerEnded(walletID int) {
	am.UpdateNotification(wallet.AccountMixer{
		WalletID:  walletID,
		RunStatus: wallet.MixerEnded,
	})
}

func (am *AccountMixerNotificationListener) UpdateNotification(signal wallet.AccountMixer) {
	am.MixerChan <- signal
}

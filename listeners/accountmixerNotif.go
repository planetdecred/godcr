package listeners

import (
	"github.com/planetdecred/godcr/wallet"
)

type AccountMixerNotificationListener struct {
	MixerChan chan wallet.AccountMixer
}

func NewAccountMixerNotificationListener(mixerCh chan wallet.AccountMixer) *AccountMixerNotificationListener {
	return &AccountMixerNotificationListener{
		MixerChan: mixerCh,
	}
}

func (am *AccountMixerNotificationListener) OnAccountMixerStarted(walletID int) {
	am.UpdateNotification(wallet.AccountMixer{
		WalletID:  walletID,
		RunStatus: wallet.MixerStarted,
	})
}

func (am *AccountMixerNotificationListener) OnAccountMixerEnded(walletID int) {
	am.UpdateNotification(wallet.AccountMixer{
		WalletID:  walletID,
		RunStatus: wallet.MixerEnded,
	})
}

func (am *AccountMixerNotificationListener) UpdateNotification(signal wallet.AccountMixer) {
	am.MixerChan <- signal
}

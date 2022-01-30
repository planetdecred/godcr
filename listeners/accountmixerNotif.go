package listeners

import (
	"github.com/planetdecred/godcr/wallet"
)

type AccountMixerNotif struct {
	MixerCh chan wallet.AccountMixer
}

func NewAccountMixerNotif(mixerCh chan wallet.AccountMixer) *AccountMixerNotif {
	return &AccountMixerNotif{
		MixerCh: mixerCh,
	}
}

func (am *AccountMixerNotif) OnAccountMixerStarted(walletID int) {
	am.UpdateNotification(wallet.AccountMixer{
		WalletID:  walletID,
		RunStatus: wallet.MixerStarted,
	})
}

func (am *AccountMixerNotif) OnAccountMixerEnded(walletID int) {
	am.UpdateNotification(wallet.AccountMixer{
		WalletID:  walletID,
		RunStatus: wallet.MixerEnded,
	})
}

func (am *AccountMixerNotif) UpdateNotification(signal wallet.AccountMixer) {
	am.MixerCh <- signal
}

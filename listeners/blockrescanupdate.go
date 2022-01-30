package listeners

import (
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/wallet"
)

type BlockRescanUpdate struct {
	BlockRescanCh chan wallet.RescanUpdate
}

func NewBlockRescanUpdate(blockRescanCh chan wallet.RescanUpdate) *BlockRescanUpdate {
	return &BlockRescanUpdate{
		BlockRescanCh: blockRescanCh,
	}
}

func (br *BlockRescanUpdate) OnBlocksRescanStarted(walletID int) {
	br.UpdateNotification(wallet.RescanUpdate{
		Stage:    wallet.RescanStarted,
		WalletID: walletID,
	})
}

func (br *BlockRescanUpdate) OnBlocksRescanProgress(progress *dcrlibwallet.HeadersRescanProgressReport) {
	br.UpdateNotification(wallet.RescanUpdate{
		Stage:          wallet.RescanProgress,
		WalletID:       progress.WalletID,
		ProgressReport: progress,
	})
}

func (br *BlockRescanUpdate) OnBlocksRescanEnded(walletID int, err error) {
	br.UpdateNotification(wallet.RescanUpdate{
		Stage:    wallet.RescanEnded,
		WalletID: walletID,
	})
}

func (br *BlockRescanUpdate) UpdateNotification(signal wallet.RescanUpdate) {
	br.BlockRescanCh <- signal
}

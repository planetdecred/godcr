package listeners

import (
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/wallet"
)

type BlocksRescanProgressListener struct {
	BlockRescanChan chan wallet.RescanUpdate
}

func NewBlocksRescanProgressListener(blockRescanCh chan wallet.RescanUpdate) *BlocksRescanProgressListener {
	return &BlocksRescanProgressListener{
		BlockRescanChan: blockRescanCh,
	}
}

func (br *BlocksRescanProgressListener) OnBlocksRescanStarted(walletID int) {
	br.UpdateNotification(wallet.RescanUpdate{
		Stage:    wallet.RescanStarted,
		WalletID: walletID,
	})
}

func (br *BlocksRescanProgressListener) OnBlocksRescanProgress(progress *dcrlibwallet.HeadersRescanProgressReport) {
	br.UpdateNotification(wallet.RescanUpdate{
		Stage:          wallet.RescanProgress,
		WalletID:       progress.WalletID,
		ProgressReport: progress,
	})
}

func (br *BlocksRescanProgressListener) OnBlocksRescanEnded(walletID int, err error) {
	br.UpdateNotification(wallet.RescanUpdate{
		Stage:    wallet.RescanEnded,
		WalletID: walletID,
	})
}

func (br *BlocksRescanProgressListener) UpdateNotification(signal wallet.RescanUpdate) {
	br.BlockRescanChan <- signal
}

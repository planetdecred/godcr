package listeners

import (
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/wallet"
)

// BlocksRescanProgressListener satisfies dcrlibwallet
// BlocksRescanProgressListener interface.
type BlocksRescanProgressListener struct {
	BlockRescanChan chan wallet.RescanUpdate
}

func NewBlocksRescanProgressListener() *BlocksRescanProgressListener {
	return &BlocksRescanProgressListener{
		BlockRescanChan: make(chan wallet.RescanUpdate, 4),
	}
}

// OnBlocksRescanStarted is a callback func called when block rescan is started.
func (br *BlocksRescanProgressListener) OnBlocksRescanStarted(walletID int) {
	br.UpdateNotification(wallet.RescanUpdate{
		Stage:    wallet.RescanStarted,
		WalletID: walletID,
	})
}

// OnBlocksRescanProgress is a callback func for block rescan progress report.
func (br *BlocksRescanProgressListener) OnBlocksRescanProgress(progress *dcrlibwallet.HeadersRescanProgressReport) {
	br.UpdateNotification(wallet.RescanUpdate{
		Stage:          wallet.RescanProgress,
		WalletID:       progress.WalletID,
		ProgressReport: progress,
	})
}

// OnBlocksRescanEnded is a callback func to notify the end of block rescan.
func (br *BlocksRescanProgressListener) OnBlocksRescanEnded(walletID int, err error) {
	br.UpdateNotification(wallet.RescanUpdate{
		Stage:    wallet.RescanEnded,
		WalletID: walletID,
	})
}

func (br *BlocksRescanProgressListener) UpdateNotification(signal wallet.RescanUpdate) {
	br.BlockRescanChan <- signal
}

package listeners

import (
	"encoding/json"

	"github.com/planetdecred/dcrlibwallet"
)

type TxAndBlockNotification struct {
	TxAndBlockNotifChan chan TxNotification
}

func NewTxAndBlockNotification(txAndBlockNotif chan TxNotification) *TxAndBlockNotification {
	return &TxAndBlockNotification{
		TxAndBlockNotifChan: txAndBlockNotif,
	}
}

func (txAndBlk *TxAndBlockNotification) OnTransaction(transaction string) {
	var tx dcrlibwallet.Transaction
	err := json.Unmarshal([]byte(transaction), &tx)
	if err == nil {
		update := TxNotification{
			NotificationType: NewTx,
			Transaction:      &tx,
		}
		txAndBlk.UpdateNotification(update)
	}
}

func (txAndBlk *TxAndBlockNotification) OnBlockAttached(walletID int, blockHeight int32) {
	txAndBlk.UpdateNotification(TxNotification{
		NotificationType: BlkAttached,
		WalletID:         walletID,
		BlockHeight:      blockHeight,
	})
}

func (txAndBlk *TxAndBlockNotification) OnTransactionConfirmed(walletID int, hash string, blockHeight int32) {
	txAndBlk.UpdateNotification(TxNotification{
		NotificationType: TxConfirmed,
		WalletID:         walletID,
		BlockHeight:      blockHeight,
		Hash:             hash,
	})
}

func (txAndBlk *TxAndBlockNotification) UpdateNotification(signal TxNotification) {
	txAndBlk.TxAndBlockNotifChan <- signal
}

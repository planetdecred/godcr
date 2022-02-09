package listeners

import (
	"encoding/json"

	"github.com/planetdecred/dcrlibwallet"
)

type TxAndBlockNotificationListener struct {
	TxAndBlockNotifChan chan TxNotification
}

func NewTxAndBlockNotificationListener(txAndBlockNotif chan TxNotification) *TxAndBlockNotificationListener {
	return &TxAndBlockNotificationListener{
		TxAndBlockNotifChan: txAndBlockNotif,
	}
}

func (txAndBlk *TxAndBlockNotificationListener) OnTransaction(transaction string) {
	var tx dcrlibwallet.Transaction
	err := json.Unmarshal([]byte(transaction), &tx)
	if err != nil {
		log.Errorf("Error unmarshalling transaction: %v", err)
		return
	}

	update := TxNotification{
		NotificationType: NewTx,
		Transaction:      &tx,
	}
	txAndBlk.UpdateNotification(update)
}

func (txAndBlk *TxAndBlockNotificationListener) OnBlockAttached(walletID int, blockHeight int32) {
	txAndBlk.UpdateNotification(TxNotification{
		NotificationType: BlkAttached,
		WalletID:         walletID,
		BlockHeight:      blockHeight,
	})
}

func (txAndBlk *TxAndBlockNotificationListener) OnTransactionConfirmed(walletID int, hash string, blockHeight int32) {
	txAndBlk.UpdateNotification(TxNotification{
		NotificationType: TxConfirmed,
		WalletID:         walletID,
		BlockHeight:      blockHeight,
		Hash:             hash,
	})
}

func (txAndBlk *TxAndBlockNotificationListener) UpdateNotification(signal TxNotification) {
	txAndBlk.TxAndBlockNotifChan <- signal
}

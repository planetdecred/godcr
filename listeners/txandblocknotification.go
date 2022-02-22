package listeners

import (
	"encoding/json"

	"github.com/planetdecred/dcrlibwallet"
)

// TxAndBlockNotificationListener satisfies dcrlibwallet
// TxAndBlockNotificationListener interface contract.
type TxAndBlockNotificationListener struct {
	TxAndBlockNotifChan chan TxNotification
}

func NewTxAndBlockNotificationListener() *TxAndBlockNotificationListener {
	return &TxAndBlockNotificationListener{
		TxAndBlockNotifChan: make(chan TxNotification, 4),
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
		NotificationType: NewTransaction,
		Transaction:      &tx,
	}
	txAndBlk.UpdateNotification(update)
}

func (txAndBlk *TxAndBlockNotificationListener) OnBlockAttached(walletID int, blockHeight int32) {
	txAndBlk.UpdateNotification(TxNotification{
		NotificationType: BlockAttached,
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

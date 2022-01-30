package listeners

import (
	"github.com/planetdecred/dcrlibwallet"
)

type TxNotifType int

const (
	NewTx TxNotifType = iota
	BlkAttached
	TxConfirmed
)

type TxNotification struct {
	NotificationType TxNotifType
	Transaction      *dcrlibwallet.Transaction
	WalletID         int
	BlockHeight      int32
	Hash             string
}

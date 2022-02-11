package listeners

import (
	"github.com/planetdecred/dcrlibwallet"
)

type TxNotifType int

const (
	// Transaction notification types
	NewTransaction TxNotifType = iota // 0 = New transaction.
	BlockAttached                     // 1 = block attached.
	TxConfirmed                       // 2 = Transaction confirmed.
)

// TxNotification models transaction notifications.
type TxNotification struct {
	NotificationType TxNotifType
	Transaction      *dcrlibwallet.Transaction
	WalletID         int
	BlockHeight      int32
	Hash             string
}

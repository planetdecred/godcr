package listeners

import (
	"github.com/planetdecred/dcrlibwallet"
)

type NotifType int

const (
	Tx NotifType = iota
	BlkAttached
	TxConfirmed
)

type TxNotification struct {
	NotificationType NotifType
	Transaction      *dcrlibwallet.Transaction
	WalletID         int
	BlockHeight      int32
	Hash             string
}

package wallet

import "github.com/planetdecred/dcrlibwallet"

// NewBlock is sent when a block is attached to the multiwallet.
type NewBlock struct {
	WalletID int
	Height   int32
}

// TxConfirmed is sent when a transaction is confirmed.
type TxConfirmed struct {
	WalletID int
	Height   int32
	Hash     string
}

type NewTransaction struct {
	Transaction *dcrlibwallet.Transaction
}

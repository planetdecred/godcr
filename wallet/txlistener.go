package wallet

// NewTransaction is sent when a new transaction is received.
type NewTransaction string

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

func (l *listener) OnTransaction(transaction string) {
	l.Send <- SyncStatusUpdate{}

}
func (l *listener) OnBlockAttached(walletID int, blockHeight int32) {
	l.Send <- SyncStatusUpdate{}
	// l.Send <- NewBlock{
	// 	WalletID: walletID,
	// 	Height:   blockHeight,
	// }
}
func (l *listener) OnTransactionConfirmed(walletID int, hash string, blockHeight int32) {
	l.Send <- SyncStatusUpdate{}
	// l.Send <- TxConfirmed{
	// 	WalletID: walletID,
	// 	Height:   blockHeight,
	// 	Hash:     hash,
	// }
}

package wallet

import "github.com/planetdecred/dcrlibwallet"

// TODO command.go file to be deprecated in subsiquent code clean up

// TODO move method to dcrlibwallet
// HaveAddress checks if the given address is valid for the wallet
func (wal *Wallet) HaveAddress(address string) (bool, string) {
	for _, wallet := range wal.multi.AllWallets() {
		result := wallet.HaveAddress(address)
		if result {
			return true, wallet.Name
		}
	}
	return false, ""
}

func (wal *Wallet) GetMultiWallet() *dcrlibwallet.MultiWallet {
	return wal.multi
}

package load

import "github.com/planetdecred/dcrlibwallet"

func (wl *WalletLoad) AllLiveTickets() ([]dcrlibwallet.Transaction, error) {
	var txs []dcrlibwallet.Transaction
	wallets := wl.MultiWallet.AllWallets()
	for _, w := range wallets {
		immatureTx, err := w.GetTransactionsRaw(0, 0, dcrlibwallet.TxFilterImmature, true)
		if err != nil {
			return txs, err
		}

		txs = append(txs, immatureTx...)

		liveTxs, err := w.GetTransactionsRaw(0, 0, dcrlibwallet.TxFilterLive, true)
		if err != nil {
			return txs, err
		}

		txs = append(txs, liveTxs...)
	}

	return txs, nil
}

func (wl *WalletLoad)


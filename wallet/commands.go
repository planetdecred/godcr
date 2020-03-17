package wallet

import (
	"fmt"
	"sort"
	"time"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
)

// CreateWallet creates a new wallet with the given parameters.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) CreateWallet(passphrase string) {
	go func() {
		var resp Response
		wall, err := wal.multi.CreateNewWallet(passphrase, dcrlibwallet.PassphraseTypePass)
		if err != nil {
			resp.Err = err
			wal.Send <- ResponseError(MultiWalletError{
				Message: "Could not create wallet",
				Err:     err,
			})
			return
		}
		resp.Resp = CreatedSeed{
			Seed: wall.Seed,
		}
		wal.Send <- resp
	}()
}

// RestoreWallet restores a wallet with the given parameters.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) RestoreWallet(seed, passphrase string) {
	go func() {
		var resp Response
		_, err := wal.multi.RestoreWallet(seed, passphrase, dcrlibwallet.PassphraseTypePass)
		if err != nil {
			resp.Err = err
			wal.Send <- ResponseError(MultiWalletError{
				Message: "Could not restore wallet",
				Err:     err,
			})
			return
		}
		resp.Resp = Restored{}
		wal.Send <- resp
	}()
}

// DeleteWallet deletes a wallet.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) DeleteWallet(walletID int, passphrase string) {
	log.Debug("Deleting Wallet")
	go func() {
		log.Debugf("Wallet %d: %+v", walletID, wal.multi.WalletWithID(walletID))
		err := wal.multi.DeleteWallet(walletID, []byte(passphrase))
		if err != nil {
			if err.Error() == dcrlibwallet.ErrInvalidPassphrase {
				err = ErrBadPass
			}
			wal.Send <- ResponseError(InternalWalletError{
				Message:  "Could not delete wallet",
				Affected: []int{walletID},
				Err:      err,
			})

		} else {
			wal.Send <- ResponseResp(DeletedWallet{ID: walletID})
		}
	}()
}

// CreateTransaction creates a TxAuthor with the given parameters.
// The created TxAuthor will have to have a destination added before broadcasting.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) CreateTransaction(walletID int, accountID int32) {
	go func() {
		var resp Response
		wallets, err := wal.wallets()
		if err != nil {
			resp.Err = err
			wal.Send <- resp
			return
		}

		if walletID > len(wallets) || walletID < 0 {
			resp.Err = err
			wal.Send <- resp
			return
		}

		if _, err := wallets[walletID].GetAccount(accountID, wal.confirms); err != nil {
			resp.Err = err
			wal.Send <- resp
			return
		}

		txAuthor := wallets[walletID].NewUnsignedTx(accountID, wal.confirms)
		if txAuthor == nil {
			resp.Err = err
			wal.Send <- resp
			return
		}

		resp.Resp = txAuthor
		wal.Send <- resp
	}()
}

// GetAllTransactions collects a per-wallet slice of transactions fitting the parameters.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) GetAllTransactions(offset, limit, txfilter int32) {
	go func() {
		var resp Response
		wallets, err := wal.wallets()
		if err != nil {
			resp.Err = err
			wal.Send <- resp
			return
		}
		alltxs := make([][]dcrlibwallet.Transaction, len(wallets))
		for i, wall := range wallets {
			txs, err := wall.GetTransactionsRaw(offset, limit, txfilter, true)
			if err != nil {
				resp.Err = err
				wal.Send <- resp
				return
			}
			alltxs[i] = txs
		}

		var recentTxs []dcrlibwallet.Transaction
		for _, tx := range alltxs {
			recentTxs = append(recentTxs, tx...)
		}
		sort.SliceStable(recentTxs, func(i, j int) bool {
			backTime := time.Unix(recentTxs[j].Timestamp, 0)
			frontTime := time.Unix(recentTxs[i].Timestamp, 0)
			return backTime.Before(frontTime)
		})
		recentTxsLimit := 5
		if len(recentTxs) > recentTxsLimit {
			recentTxs = recentTxs[:recentTxsLimit]
		}

		resp.Resp = &Transactions{
			Txs:    alltxs,
			Recent: recentTxs,
		}
		wal.Send <- resp
	}()
}

// GetTransactionsByWallet get list of transactions fitting the parameters.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) GetTransactionsByWallet(walletID int, offset, limit, txfilter int32) {
	go func() {
		var resp Response

		wallets, err := wal.wallets()

		if err != nil {
			resp.Err = err
			wal.Send <- resp

			return
		}

		var alltxs []dcrlibwallet.Transaction

		for _, wall := range wallets {
			if wall.ID == walletID {
				txs, err := wall.GetTransactionsRaw(offset, limit, txfilter, true)
				if err != nil {
					resp.Err = err
					wal.Send <- resp

					return
				}

				alltxs = txs
			}
		}

		resp.Resp = TransactionsWallet{
			Txs: alltxs,
		}
		wal.Send <- resp
	}()
}

// GetMultiWalletInfo gets bulk information about the loaded wallets.
// Information regarding transactions is collected with respect to wal.confirms as the
// number of required confirmations for said transactions.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) GetMultiWalletInfo() {
	go func() {
		log.Debug("Getting multiwallet info")
		var resp Response
		wallets, err := wal.wallets()
		if err != nil {
			resp.Err = err
			wal.Send <- resp
			return
		}

		var completeTotal int64
		infos := make([]InfoShort, len(wallets))
		i := 0
		for _, wall := range wallets {
			iter, err := wall.AccountsIterator(wal.confirms)
			if err != nil {
				resp.Err = err
				wal.Send <- resp
				return
			}

			walletAccounts := []Account{}
			var totalWalletBalance, spendableWalletBalance int64
			for acct := iter.Next(); acct != nil; acct = iter.Next() {
				totalWalletBalance += acct.TotalBalance
				spendableWalletBalance += acct.Balance.Spendable

				account := Account{
					Number:       fmt.Sprint(acct.Number),
					Name:         acct.Name,
					TotalBalance: fmt.Sprint(acct.Balance.Total),
				}
				walletAccounts = append(walletAccounts, account)
			}

			completeTotal += totalWalletBalance
			infos[i] = InfoShort{
				ID:              wall.ID,
				Name:            wall.Name,
				Accounts:        walletAccounts,
				BestBlockHeight: wall.GetBestBlock(),
				BlockTimestamp:  wall.GetBestBlockTimeStamp(),
				IsWaiting:       wall.IsWaiting(),
			}
			i++
		}
		best := wal.multi.GetBestBlock()

		if best == nil {
			if len(wallets) == 0 {
				wal.Send <- ResponseResp(MultiWalletInfo{})
				return
			}
			resp.Err = InternalWalletError{
				Message: "Could not get load best block",
			}
			wal.Send <- resp
			return
		}

		resp.Resp = MultiWalletInfo{
			LoadedWallets:   len(wallets),
			TotalBalance:    dcrutil.Amount(completeTotal).String(),
			BestBlockHeight: best.Height,
			BestBlockTime:   best.Timestamp,
			Wallets:         infos,
			Synced:          wal.multi.IsSynced(),
			Syncing:         wal.multi.IsSyncing(),
		}
		wal.Send <- resp
	}()
}

// RenameWallet renames the wallet identified by walletID.
func (wal *Wallet) RenameWallet(walletID int, name string) error {
	return wal.multi.RenameWallet(walletID, name)
}

// CurrentAddress returns the next address for the specified wallet account.
func (wal *Wallet) CurrentAddress(walletID int, accountID int32) (string, error) {
	wall := wal.multi.WalletWithID(walletID)
	if wall == nil {
		return "", ErrIDNotExist
	}
	return wall.CurrentAddress(accountID)
}

// NextAddress returns the next address for the specified wallet account.
func (wal *Wallet) NextAddress(walletID int, accountID int32) (string, error) {
	wall := wal.multi.WalletWithID(walletID)
	if wall == nil {
		return "", ErrIDNotExist
	}
	return wall.NextAddress(accountID)
}

// IsAddressValid checks if the given address is valid for the multiwallet network
func (wal *Wallet) IsAddressValid(address string) (bool, error) {
	wall := wal.multi.FirstOrDefaultWallet()
	if wall == nil {
		return false, InternalWalletError{
			Message: "No wallet loaded",
		}
	}
	return wall.IsAddressValid(address), nil
}

// StartSync starts the multiwallet SPV sync
func (wal *Wallet) StartSync() error {
	return wal.multi.SpvSync()
}

// CancelSync cancels the SPV sync
func (wal *Wallet) CancelSync() {
	go wal.multi.CancelSync()
}

package wallet

import (
	"encoding/base64"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/planetdecred/dcrlibwallet"
)

// CreateWallet creates a new wallet with the given parameters.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) CreateWallet(name, passphrase string, errChan chan error) {
	go func() {
		var resp Response
		wall, err := wal.multi.CreateNewWallet(name, passphrase, dcrlibwallet.PassphraseTypePass)
		sendErr := func(err error) {
			go func() {
				errChan <- err
			}()
			resp.Err = err
			wal.Send <- ResponseError(MultiWalletError{
				Message: "Could not create wallet",
				Err:     err,
			})
		}
		if err != nil {
			sendErr(err)
			return
		}
		seeds, err := wall.DecryptSeed([]byte(passphrase))
		if err != nil {
			sendErr(err)
			return
		}

		resp.Resp = CreatedSeed{
			Seed: seeds,
		}
		wal.Send <- resp
	}()
}

// RestoreWallet restores a wallet with the given parameters.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) RestoreWallet(seed, passphrase string, errChan chan error) {
	go func() {
		var resp Response
		_, err := wal.multi.RestoreWallet("wallet", seed, passphrase, dcrlibwallet.PassphraseTypePass)
		if err != nil {
			go func() {
				errChan <- err
			}()
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
func (wal *Wallet) DeleteWallet(walletID int, passphrase []byte, errChan chan error) {
	log.Debug("Deleting Wallet")
	go func() {
		var resp Response
		log.Debugf("Wallet %d: %+v", walletID, wal.multi.WalletWithID(walletID))
		err := wal.multi.DeleteWallet(walletID, passphrase)
		if err != nil {
			go func() {
				errChan <- err
			}()
			resp.Err = err
			wal.Send <- ResponseError(InternalWalletError{
				Message:  "Could not delete wallet",
				Affected: []int{walletID},
				Err:      err,
			})
			return
		}
		resp.Resp = DeletedWallet{
			ID: walletID,
		}
		wal.Send <- resp
	}()
}

// AddAccount adds an account to a wallet.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) AddAccount(walletID int, name string, pass []byte, errChan chan error) {
	go func() {
		var resp Response
		wall := wal.multi.WalletWithID(walletID)
		if wall == nil {
			go func() {
				errChan <- ErrIDNotExist
			}()
			resp.Err = ErrIDNotExist
			wal.Send <- Response{
				Resp: AddedAccount{},
				Err:  ErrIDNotExist,
			}
			return
		}

		id, err := wall.NextAccount(name)
		if err != nil {
			go func() {
				errChan <- err
			}()
			resp.Err = err
			wal.Send <- ResponseError(InternalWalletError{
				Message:  "Could not create account",
				Affected: []int{walletID},
				Err:      err,
			})
			return
		}
		resp.Resp = AddedAccount{
			ID: id,
		}
		wal.Send <- resp
	}()
}

// CreateTransaction creates a TxAuthor with the given parameters.
// The created TxAuthor will have to have a destination added before broadcasting.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) CreateTransaction(walletID int, accountID int32, errChan chan error) {
	go func() {
		var resp Response
		wall := wal.multi.WalletWithID(walletID)
		_, err := wall.GetAccount(accountID)
		if err != nil {
			errChan <- err
			return
		}

		txAuthor := wal.multi.NewUnsignedTx(wall, accountID)
		if txAuthor == nil {
			errChan <- err
			return
		}
		resp.Resp = txAuthor
		wal.Send <- resp
	}()
}

// transactionStatus accepts the bestBlockHeight, transactionBlockHeight returns a transaction status
// which could be confirmed/pending and confirmations count
func transactionStatus(bestBlockHeight, txnBlockHeight int32) (string, int32) {
	confirmations := bestBlockHeight - txnBlockHeight + 1
	if txnBlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
		return "confirmed", confirmations
	}
	return "pending", confirmations
}

// BroadcastTransaction broadcasts the transaction built with txAuthor to the network.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) BroadcastTransaction(txAuthor *dcrlibwallet.TxAuthor, passphrase []byte, errChan chan error) {
	go func() {
		var resp Response

		txHash, err := txAuthor.Broadcast(passphrase)
		if err != nil {
			errChan <- fmt.Errorf("error broadcasting transaction: %s", err.Error())
			return
		}

		hash, err := chainhash.NewHash(txHash)
		if err != nil {
			errChan <- fmt.Errorf("error parsing successful transaction hash: %s", err.Error())
			return
		}

		resp.Resp = &Broadcast{
			TxHash: hash.String(),
		}
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

		var recentTxs []Transaction

		transactions := make(map[int][]Transaction)
		ticketTxs := make(map[int][]Transaction)
		bestBestBlock := wal.multi.GetBestBlock()
		totalTxn := 0

		for _, wall := range wallets {
			txs, err := wall.GetTransactionsRaw(offset, limit, txfilter, true)
			if err != nil {
				resp.Err = err
				wal.Send <- resp
				return
			}
			for _, txnRaw := range txs {
				totalTxn++
				status, confirmations := transactionStatus(bestBestBlock.Height, txnRaw.BlockHeight)
				txn := Transaction{
					Txn:           txnRaw,
					Status:        status,
					Balance:       dcrutil.Amount(txnRaw.Amount).String(),
					WalletName:    wall.Name,
					Confirmations: confirmations,
					DateTime:      dcrlibwallet.ExtractDateOrTime(txnRaw.Timestamp),
				}
				recentTxs = append(recentTxs, txn)
				if txn.Txn.Type == dcrlibwallet.TxTypeTicketPurchase {
					ticketTxs[wall.ID] = append(ticketTxs[wall.ID], txn)
				}
				transactions[txnRaw.WalletID] = append(transactions[txnRaw.WalletID], txn)
			}
		}

		sort.SliceStable(recentTxs, func(i, j int) bool {
			backTime := time.Unix(recentTxs[j].Txn.Timestamp, 0)
			frontTime := time.Unix(recentTxs[i].Txn.Timestamp, 0)
			return backTime.Before(frontTime)
		})

		recentTxsLimit := 5
		if len(recentTxs) > recentTxsLimit {
			recentTxs = recentTxs[:recentTxsLimit]
		}

		resp.Resp = &Transactions{
			Total:   totalTxn,
			Txs:     transactions,
			Recent:  recentTxs,
			Tickets: ticketTxs,
		}
		wal.Send <- resp
	}()
}

// GetTransaction get transaction information by wallet ID and transaction hash
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) GetTransaction(walletID int, txnHash string) {
	go func() {
		var resp Response
		wall := wal.multi.WalletWithID(walletID)

		hash, err := chainhash.NewHashFromStr(txnHash)
		if err != nil {
			resp.Err = err
			wal.Send <- resp
			return
		}

		txn, err := wall.GetTransactionRaw(hash[:])
		if err != nil {
			resp.Err = err
			wal.Send <- resp
			return
		}
		bestBestBlock := wal.multi.GetBestBlock()
		status, confirmations := transactionStatus(bestBestBlock.Height, txn.BlockHeight)
		acct, err := wall.GetAccount(txn.Inputs[0].AccountNumber)
		if err != nil {
			resp.Err = err
			wal.Send <- resp
			return
		}
		resp.Resp = &Transaction{
			Txn:           *txn,
			Status:        status,
			Balance:       dcrutil.Amount(txn.Amount).String(),
			WalletName:    wall.Name,
			Confirmations: confirmations,
			DateTime:      dcrlibwallet.ExtractDateOrTime(txn.Timestamp),
			AccountName:   acct.Name,
		}
		wal.Send <- resp
	}()
}

// WalletSyncStatus returns the sync status of a single wallet
func walletSyncStatus(isWaiting bool, walletBestBlock, bestBlockHeight int32) string {
	if isWaiting {
		return "waiting for other wallets"
	}
	if walletBestBlock < bestBlockHeight {
		return "syncing..."
	}

	return "synced"
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
			iter, err := wall.AccountsIterator()
			if err != nil {
				resp.Err = err
				wal.Send <- resp
				return
			}
			var acctBalance, spendableBalance int64
			accts := make([]Account, 0)
			for acct := iter.Next(); acct != nil; acct = iter.Next() {
				var addr string
				if acct.Number != math.MaxInt32 {
					var er error
					addr, er = wall.CurrentAddress(acct.Number)
					if er != nil {
						log.Error("Could not get current address for wallet ", wall.ID, "account", acct.Number)
					}
				}
				accts = append(accts, Account{
					Number:           acct.Number,
					Name:             acct.Name,
					TotalBalance:     dcrutil.Amount(acct.TotalBalance).String(),
					SpendableBalance: acct.Balance.Spendable,
					Keys: struct {
						Internal, External, Imported string
					}{
						Internal: strconv.Itoa(int(acct.InternalKeyCount)),
						External: strconv.Itoa(int(acct.ExternalKeyCount)),
						Imported: strconv.Itoa(int(acct.ImportedKeyCount)),
					},
					HDPath:         wal.hdPrefix() + strconv.Itoa(int(acct.Number)) + "'",
					CurrentAddress: addr,
				})
				acctBalance += acct.TotalBalance
				spendableBalance += acct.Balance.Spendable
			}
			completeTotal += acctBalance

			infos[i] = InfoShort{
				ID:               wall.ID,
				Name:             wall.Name,
				Balance:          dcrutil.Amount(acctBalance).String(),
				SpendableBalance: spendableBalance,
				Accounts:         accts,
				BestBlockHeight:  wall.GetBestBlock(),
				BlockTimestamp:   wall.GetBestBlockTimeStamp(),
				DaysBehind:       fmt.Sprintf("%s behind", calculateDaysBehind(wall.GetBestBlockTimeStamp())),
				Status:           walletSyncStatus(wall.IsWaiting(), wall.GetBestBlock(), wal.OverallBlockHeight),
				Seed:             wall.EncryptedSeed,
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

		lastSyncTime := int64(time.Since(time.Unix(best.Timestamp, 0)).Seconds())
		resp.Resp = MultiWalletInfo{
			LoadedWallets:   len(wallets),
			TotalBalance:    dcrutil.Amount(completeTotal).String(),
			BestBlockHeight: best.Height,
			BestBlockTime:   best.Timestamp,
			LastSyncTime:    SecondsToDays(lastSyncTime),
			Wallets:         infos,
			Synced:          wal.multi.IsSynced(),
			Syncing:         wal.multi.IsSyncing(),
		}
		wal.Send <- resp
	}()
}

func (wal *Wallet) SignMessage(walletID int, passphrase []byte, address, message string, errChan chan error) {
	go func() {
		var resp Response

		wall := wal.multi.WalletWithID(walletID)
		if wall == nil {
			resp.Resp = &Signature{
				Err: InternalWalletError{
					Message: "No wallet found",
				},
			}
			wal.Send <- resp
			return
		}

		signedMessageBytes, err := wall.SignMessage(passphrase, address, message)
		if err != nil {
			go func() {
				errChan <- err
			}()
			resp.Resp = &Signature{
				Err: fmt.Errorf("error signing message: %s", err.Error()),
			}
			wal.Send <- resp
			return
		}

		resp.Resp = &Signature{
			Signature: base64.StdEncoding.EncodeToString(signedMessageBytes),
		}

		wal.Send <- resp
	}()
}

// RenameWallet renames the wallet identified by walletID.
func (wal *Wallet) RenameWallet(walletID int, name string, errChan chan error) {
	go func() {
		var resp Response
		err := wal.multi.RenameWallet(walletID, name)
		if err != nil {
			go func() {
				errChan <- err
			}()
			resp.Err = err
			wal.Send <- ResponseError(MultiWalletError{
				Message: "Could not rename wallet",
				Err:     err,
			})
			return
		}
		resp.Resp = Renamed{
			ID: walletID,
		}
		wal.Send <- resp
	}()
}

// ChangeWalletPassphrase changes the spending passphrase of the wallet identified by walletID.
func (wal *Wallet) ChangeWalletPassphrase(walletID int, oldPrivatePassphrase, newPrivatePassphrase string, errChan chan error) {
	go func() {
		var resp Response
		err := wal.multi.ChangePrivatePassphraseForWallet(walletID, []byte(oldPrivatePassphrase), []byte(newPrivatePassphrase), dcrlibwallet.PassphraseTypePass)
		if err != nil {
			go func() {
				errChan <- err
			}()
			resp.Err = err
			wal.Send <- ResponseError(InternalWalletError{
				Message:  "Could not change password",
				Affected: []int{walletID},
				Err:      err,
			})
			return
		}

		resp.Resp = &ChangePassword{
			ID: walletID,
		}
	}()
}

func (wal *Wallet) OpenWallets(passphrase string, errChan chan error) {
	go func() {
		var resp Response
		err := wal.multi.OpenWallets([]byte(passphrase))
		if err != nil {
			go func() {
				errChan <- err
			}()
			resp.Err = err
			wal.Send <- ResponseError(MultiWalletError{
				Message: "Could not open wallets",
				Err:     err,
			})
			return
		}

		resp.Resp = OpenWallet{}
		wal.Send <- resp
	}()
}

func (wal *Wallet) SetStartupPassphrase(passphrase string, errChan chan error) {
	go func() {
		var resp Response
		err := wal.multi.SetStartupPassphrase([]byte(passphrase), dcrlibwallet.PassphraseTypePass)
		if err != nil {
			go func() {
				errChan <- err
			}()
			resp.Err = err
			wal.Send <- ResponseError(MultiWalletError{
				Message: "Could not set up startup passphrase",
				Err:     err,
			})
			return
		}
		resp.Resp = &StartupPassphrase{
			Msg: "Startup password set",
		}
		wal.Send <- resp
	}()
}

func (wal *Wallet) ChangeStartupPassphrase(oldPrivatePassphrase, newPrivatePassphrase string, errChan chan error) {
	go func() {
		var resp Response
		err := wal.multi.ChangeStartupPassphrase([]byte(oldPrivatePassphrase), []byte(newPrivatePassphrase), dcrlibwallet.PassphraseTypePass)
		if err != nil {
			go func() {
				errChan <- err
			}()
			resp.Err = err
			wal.Send <- ResponseError(MultiWalletError{
				Message: "Could not change startup passphrase",
				Err:     err,
			})
			return
		}

		resp.Resp = &StartupPassphrase{
			Msg: "Startup password changed",
		}
		wal.Send <- resp
	}()
}

func (wal *Wallet) RemoveStartupPassphrase(passphrase string, errChan chan error) {
	go func() {
		var resp Response
		err := wal.multi.RemoveStartupPassphrase([]byte(passphrase))
		if err != nil {
			go func() {
				errChan <- err
			}()
			resp.Err = err
			wal.Send <- ResponseError(MultiWalletError{
				Message: "Could not remove startup passphrase",
				Err:     err,
			})
			return
		}
		resp.Resp = &StartupPassphrase{
			Msg: "Startup password removed",
		}
		wal.Send <- resp
	}()
}

// IsStartupSecuritySet checks if start up password is set
func (wal *Wallet) IsStartupSecuritySet() bool {
	return wal.multi.IsStartupSecuritySet()
}

// RenameAccount renames the acct of wallet with id walletID.
func (wal *Wallet) RenameAccount(walletID int, acct int32, name string, errChan chan<- error) {
	go func() {
		var resp Response
		wall := wal.multi.WalletWithID(walletID)
		if wall == nil {
			go func() {
				errChan <- ErrIDNotExist
			}()
			resp.Err = ErrIDNotExist
			wal.Send <- Response{
				Resp: UpdatedAccount{},
				Err:  ErrIDNotExist,
			}
			return
		}

		err := wall.RenameAccount(acct, name)
		if err != nil {
			go func() {
				errChan <- err
			}()
			resp.Err = err
			wal.Send <- ResponseError(InternalWalletError{
				Message:  "Could not rename account",
				Affected: []int{walletID},
				Err:      err,
			})
			return
		}
		resp.Resp = UpdatedAccount{
			ID: acct,
		}
		wal.Send <- resp
	}()
}

func (wal *Wallet) UnlockWallet(walletID int, passphrase []byte) error {
	return wal.multi.UnlockWallet(walletID, passphrase)
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
	return wal.multi.IsAddressValid(address), nil
}

// HaveAddress checks if the given address is valid for the wallet
func (wal *Wallet) HaveAddress(walletID int, address string) (bool, error) {
	wall := wal.multi.WalletWithID(walletID)
	if wall == nil {
		return false, ErrIDNotExist
	}
	return wall.HaveAddress(address), nil
}

// VerifyMessage checks if the given message matches the signature for the address.
func (wal *Wallet) VerifyMessage(address string, message string, signature string) (bool, error) {
	return wal.multi.VerifyMessage(address, message, signature)
}

// StartSync starts the multiwallet SPV sync
func (wal *Wallet) StartSync() error {
	return wal.multi.SpvSync()
}

// RescanBlocks rescans the multiwallet
func (wal *Wallet) RescanBlocks(walletID int) error {
	return wal.multi.RescanBlocks(walletID)
}

// CancelSync cancels the SPV sync
func (wal *Wallet) CancelSync() {
	go wal.multi.CancelSync()
}

func (wal *Wallet) GetWalletSeedPhrase(walletID int, password []byte) (string, error) {
	return wal.multi.WalletWithID(walletID).DecryptSeed(password)
}

func (wal *Wallet) VerifyWalletSeedPhrase(walletID int, seedPhrase string, privpass []byte) error {
	_, err := wal.multi.VerifySeedForWallet(walletID, seedPhrase, privpass)
	return err
}

func (wal *Wallet) SaveConfigValueForKey(key string, value interface{}) {
	wal.multi.SaveUserConfigValue(key, value)
}

func (wal *Wallet) ReadBoolConfigValueForKey(key string) bool {
	return wal.multi.ReadBoolConfigValueForKey(key, false)
}

func (wal *Wallet) ReadStringConfigValueForKey(key string) string {
	return wal.multi.ReadStringConfigValueForKey(key)
}

func (wal *Wallet) RemoveUserConfigValueForKey(key string) {
	wal.multi.DeleteUserConfigValueForKey(key)
}

func (wal *Wallet) LoadedWalletsCount() int32 {
	return wal.multi.LoadedWalletsCount()
}

func calculateDaysBehind(lastHeaderTime int64) string {
	diff := time.Since(time.Unix(lastHeaderTime, 0))
	daysBehind := int(math.Round(diff.Hours() / 24))
	if daysBehind < 1 {
		return "<1 day"
	} else if daysBehind == 1 {
		return "1 day"
	} else {
		return fmt.Sprintf("%d days", daysBehind)
	}
}

// SecondsToDays takes time in seconds and returns its string equivalent in the format ddhhmm.
func SecondsToDays(totalTimeLeft int64) string {
	q, r := divMod(totalTimeLeft, 24*60*60)
	timeLeft := time.Duration(r) * time.Second
	if q > 0 {
		return fmt.Sprintf("%dd%s", q, timeLeft.String())
	}
	return timeLeft.String()
}

// divMod divides a numerator by a denominator and returns its quotient and remainder.
func divMod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

// AllUnspentOutputs get all unspent outputs by walletID and acct
func (wal *Wallet) AllUnspentOutputs(walletID int, acct int32) {
	go func() {
		var resp Response
		wall := wal.multi.WalletWithID(walletID)
		if wall == nil {
			resp.Err = ErrIDNotExist
			wal.Send <- Response{
				Resp: UnspentOutputs{},
				Err:  ErrIDNotExist,
			}
			return
		}
		utxos, err := wall.UnspentOutputs(acct)
		if err != nil {
			resp.Err = err
			wal.Send <- ResponseError(InternalWalletError{
				Message:  "Could not get unspent outputs",
				Affected: []int{walletID, int(acct)},
				Err:      err,
			})
			return
		}

		var list []*UnspentOutput
		for _, utxo := range utxos {
			item := UnspentOutput{
				UTXO:     *utxo,
				Amount:   dcrutil.Amount(utxo.Amount).String(),
				DateTime: dcrlibwallet.ExtractDateOrTime(utxo.ReceiveTime),
			}
			list = append(list, &item)
		}
		resp.Resp = &UnspentOutputs{
			List: list,
		}
		wal.Send <- resp
	}()
}

// IsAccountMixerConfigSet check the wallet have account mixer config set
func (wal *Wallet) IsAccountMixerConfigSet(walletID int) bool {
	wall := wal.multi.WalletWithID(walletID)
	if wall == nil {
		return false
	}
	return wall.ReadBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, false)
}

// SetupAccountMixer setup account mixer with the given parameters.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) SetupAccountMixer(walletID int, walletPassphrase string, errChan chan error) {
	go func() {
		var resp Response

		wall := wal.multi.WalletWithID(walletID)
		if wall == nil {
			resp.Err = ErrIDNotExist
			wal.Send <- Response{
				Resp: SetupAccountMixer{},
				Err:  ErrIDNotExist,
			}
			return
		}

		err := wall.SetAccountMixerConfig("mixed", "unmixed", walletPassphrase)
		if err != nil {
			go func() {
				errChan <- err
			}()
			resp.Err = err
			wal.Send <- ResponseError(InternalWalletError{
				Message:  "Could not set account mixer",
				Err:      err,
				Affected: []int{walletID},
			})
			return
		}

		resp.Resp = SetupAccountMixer{}
		wal.Send <- resp
	}()
}

func (wal *Wallet) NewVSPD(walletID int, accountID int32) *dcrlibwallet.VSPD {
	return wal.multi.NewVSPD("http://dev.planetdecred.org:23125", walletID, accountID)
}

// TicketPrice get ticket price
func (wal *Wallet) TicketPrice(walletID int) string {
	wall := wal.multi.WalletWithID(walletID)
	pr, err := wall.TicketPrice()
	if err != nil {
		log.Error(err)
		return ""
	}
	return dcrutil.Amount(pr.TicketPrice).String()
}

// PurchaseTicket buy a ticket with given parameters
func (wal *Wallet) PurchaseTicket(walletID int, accountID int32, tickets uint32, passphrase []byte, expiry uint32) ([]string, error) {
	wall := wal.multi.WalletWithID(walletID)
	request := &dcrlibwallet.PurchaseTicketsRequest{
		Account:               uint32(accountID),
		Passphrase:            passphrase,
		NumTickets:            tickets,
		Expiry:                uint32(wal.multi.GetBestBlock().Height) + expiry,
		RequiredConfirmations: dcrlibwallet.DefaultRequiredConfirmations,
	}
	hashes, err := wall.PurchaseTickets(request, "")
	if err != nil {
		return []string{}, err
	}
	go func() {
		var resp Response
		resp.Resp = &TicketPurchase{}
		wal.Send <- resp
	}()
	return hashes, nil
}

// GetAllTickets collects a per-wallet slice of tickets fitting the parameters.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) GetAllTickets() {
	go func() {
		var resp Response
		wallets, err := wal.wallets()
		if err != nil {
			resp.Err = err
			wal.Send <- resp
			return
		}
		tickets := make(map[int][]Ticket)
		unconfirmedTickets := make(map[int][]UnconfirmedPurchase)
		totalTicket := 0

		for _, wall := range wallets {
			ticketsInfo, err := wall.GetTicketsForBlockHeightRange(0, wall.GetBestBlock(), math.MaxInt32)
			if err != nil {
				resp.Err = err
				wal.Send <- resp
				return
			}
			for _, tinfo := range ticketsInfo {
				var amount dcrutil.Amount
				for _, output := range tinfo.Ticket.MyOutputs {
					amount += output.Amount
				}
				info := Ticket{
					Info:     *tinfo,
					DateTime: dcrlibwallet.ExtractDateOrTime(tinfo.Ticket.Timestamp),
					Amount:   amount.String(),
					Fee:      tinfo.Ticket.Fee.String(),
				}
				tickets[wall.ID] = append(tickets[wall.ID], info)
			}

			unconfirmedTicketPurchases, err := getUnconfirmedPurchases(wall, tickets[wall.ID])
			if err != nil {
				resp.Err = err
				wal.Send <- resp
				return
			}
			unconfirmedTickets[wall.ID] = unconfirmedTicketPurchases
		}

		resp.Resp = &Tickets{
			Total:       totalTicket,
			Confirmed:   tickets,
			Unconfirmed: unconfirmedTickets,
		}
		wal.Send <- resp
	}()
}

func getUnconfirmedPurchases(wall dcrlibwallet.Wallet, tickets []Ticket) (unconfirmed []UnconfirmedPurchase, err error) {
	contains := func(slice []Ticket, item string) bool {
		set := make(map[string]struct{}, len(slice))
		for _, s := range slice {
			set[s.Info.Ticket.Hash.String()] = struct{}{}
		}

		_, ok := set[item]
		return ok
	}

	txs, err := wall.GetTransactionsRaw(0, 0, dcrlibwallet.TxFilterAll, true)
	if err != nil {
		return
	}

	ticketTxs := make(map[int][]dcrlibwallet.Transaction)
	for _, txn := range txs {
		if txn.Type == dcrlibwallet.TxTypeTicketPurchase {
			ticketTxs[wall.ID] = append(ticketTxs[wall.ID], txn)
		}
	}

	if len(tickets) == len(ticketTxs) {
		return
	}

	for _, txn := range ticketTxs[wall.ID] {
		var amount int64
		for _, output := range txn.Outputs {
			amount += output.Amount
		}

		if !contains(tickets, txn.Hash) {
			unconfirmed = append(unconfirmed, UnconfirmedPurchase{
				Hash:        txn.Hash,
				Status:      "UNCONFIRMED",
				DateTime:    dcrlibwallet.ExtractDateOrTime(txn.Timestamp),
				BlockHeight: txn.BlockHeight,
				Amount:      dcrutil.Amount(amount).String(),
			})
		}
	}

	return
}

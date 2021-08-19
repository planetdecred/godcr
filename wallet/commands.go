package wallet

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/planetdecred/dcrlibwallet"
)

// transactionStatus accepts the bestBlockHeight, transactionBlockHeight returns a transaction status
// which could be confirmed/pending and confirmations count
func transactionStatus(bestBlockHeight, txnBlockHeight int32) (string, int32) {
	confirmations := bestBlockHeight - txnBlockHeight + 1
	if txnBlockHeight != -1 && confirmations > dcrlibwallet.DefaultRequiredConfirmations {
		return "confirmed", confirmations
	}
	return "pending", confirmations
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

		txn, err := wall.GetTransactionRaw(txnHash)
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

// CancelSync cancels the SPV sync
func (wal *Wallet) CancelSync() {
	go wal.multi.CancelSync()
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
					Balance: Balance{
						Total:                   acct.Balance.Total,
						Spendable:               acct.Balance.Spendable,
						ImmatureReward:          acct.Balance.ImmatureReward,
						ImmatureStakeGeneration: acct.Balance.ImmatureStakeGeneration,
						LockedByTickets:         acct.Balance.LockedByTickets,
						VotingAuthority:         acct.Balance.VotingAuthority,
						UnConfirmed:             acct.Balance.UnConfirmed,
					},
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
				IsWatchingOnly:   wall.IsWatchingOnlyWallet(),
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
			TotalBalanceRaw: GetRawBalance(completeTotal, 0),
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

func (wal *Wallet) GetMultiWallet() *dcrlibwallet.MultiWallet {
	return wal.multi
}

func (wal *Wallet) GetAllProposals() {
	var resp Response
	go func() {
		proposals, err := wal.multi.Politeia.GetProposalsRaw(dcrlibwallet.ProposalCategoryAll, 0, 0, true)
		if err != nil {
			resp.Err = err
			wal.Send <- resp
			return
		}
		resp.Resp = &Proposals{
			Proposals: proposals,
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
func (wal *Wallet) HaveAddress(address string) (bool, string) {
	for _, wallet := range wal.multi.AllWallets() {
		result := wallet.HaveAddress(address)
		if result {
			return true, wallet.Name
		}
	}
	return false, ""
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

func (wal *Wallet) IsSyncingProposals() bool {
	return wal.multi.Politeia.IsSyncing()
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

// GetRawBalance gets the balance in int64, formats it and returns a string while also leaving out the "DCR" suffix
func GetRawBalance(balance int64, AmountUnit int) string {
	return strconv.FormatFloat(float64(balance)/math.Pow10(AmountUnit+8), 'f', -(AmountUnit + 8), 64)
}

// divMod divides a numerator by a denominator and returns its quotient and remainder.
func divMod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

// GetAllTickets collects a per-wallet slice of tickets fitting the parameters.
// It is non-blocking and sends its result or any error to wal.Send.
func (wal *Wallet) GetAllTickets() {
	go func() {
		log.Info("fetching all tickets")
		var resp Response
		wallets, err := wal.wallets()
		if err != nil {
			resp.Err = err
			wal.Send <- resp
			return
		}

		var liveRecentTickets []Ticket
		var recentActivity []Ticket

		tickets := make(map[int][]Ticket)
		unconfirmedTickets := make(map[int][]UnconfirmedPurchase)

		stackingRecordCounter := []struct {
			Status string
			Count  int
		}{
			{"UNMINED", 0},
			{"IMMATURE", 0},
			{"LIVE", 0},
			{"VOTED", 0},
			{"MISSED", 0},
			{"EXPIRED", 0},
			{"REVOKED", 0},
		}

		liveCounter := []struct {
			Status string
			Count  int
		}{
			{"UNMINED", 0},
			{"IMMATURE", 0},
			{"LIVE", 0},
		}

		for _, wall := range wallets {
			ticketsInfo, err := wall.GetTicketsForBlockHeightRange(0, wall.GetBestBlock(), math.MaxInt32)
			if err != nil {
				resp.Err = err
				wal.Send <- resp
				return
			}

			for _, tinfo := range ticketsInfo {
				if tinfo.Status == "UNKNOWN" {
					continue
				}

				var amount dcrutil.Amount
				for _, output := range tinfo.Ticket.MyOutputs {
					amount += output.Amount
				}
				info := Ticket{
					Info:       *tinfo,
					DateTime:   time.Unix(tinfo.Ticket.Timestamp, 0).Format("Jan 2, 2006 03:04:05 PM"),
					MonthDay:   time.Unix(tinfo.Ticket.Timestamp, 0).Format("Jan 2"),
					DaysBehind: calculateDaysBehind(tinfo.Ticket.Timestamp),
					Amount:     amount.String(),
					Fee:        tinfo.Ticket.Fee.String(),
					WalletName: wall.Name,
				}
				tickets[wall.ID] = append(tickets[wall.ID], info)

				for i := range liveCounter {
					if liveCounter[i].Status == tinfo.Status {
						liveCounter[i].Count++
					}
				}

				if tinfo.Status == "UNMINED" || tinfo.Status == "IMMATURE" || tinfo.Status == "LIVE" {
					liveRecentTickets = append(liveRecentTickets, info)
				}

				recentActivity = append(recentActivity, info)

				for i := range stackingRecordCounter {
					if stackingRecordCounter[i].Status == tinfo.Status {
						stackingRecordCounter[i].Count++
					}
				}
			}

			sort.SliceStable(tickets[wall.ID], func(i, j int) bool {
				backTime := time.Unix(tickets[wall.ID][j].Info.Ticket.Timestamp, 0)
				frontTime := time.Unix(tickets[wall.ID][i].Info.Ticket.Timestamp, 0)
				return backTime.Before(frontTime)
			})

			unconfirmedTicketPurchases, err := getUnconfirmedPurchases(wall, tickets[wall.ID])
			if err != nil {
				resp.Err = err
				wal.Send <- resp
				return
			}
			unconfirmedTickets[wall.ID] = unconfirmedTicketPurchases
		}

		sort.SliceStable(liveRecentTickets, func(i, j int) bool {
			backTime := time.Unix(liveRecentTickets[j].Info.Ticket.Timestamp, 0)
			frontTime := time.Unix(liveRecentTickets[i].Info.Ticket.Timestamp, 0)
			return backTime.Before(frontTime)
		})

		recentLimit := 5
		if len(liveRecentTickets) > recentLimit {
			liveRecentTickets = liveRecentTickets[:recentLimit]
		}

		sort.SliceStable(recentActivity, func(i, j int) bool {
			backTime := time.Unix(recentActivity[j].Info.Ticket.Timestamp, 0)
			frontTime := time.Unix(recentActivity[i].Info.Ticket.Timestamp, 0)
			return backTime.Before(frontTime)
		})

		if len(recentActivity) > recentLimit {
			recentActivity = recentActivity[:recentLimit]
		}

		resp.Resp = &Tickets{
			Confirmed:             tickets,
			Unconfirmed:           unconfirmedTickets,
			RecentActivity:        recentActivity,
			StackingRecordCounter: stackingRecordCounter,
			LiveRecent:            liveRecentTickets,
			LiveCounter:           liveCounter,
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

package load
//
//import (
//	"fmt"
//	"math"
//	"sort"
//	"time"
//
//	"github.com/decred/dcrd/dcrutil"
//	"github.com/planetdecred/dcrlibwallet"
//)
//
//type Ticket struct {
//	Status     string
//	Fee        string
//	Amount     string
//	DateTime   string
//	MonthDay   string
//	DaysBehind string
//	WalletName string
//
//	timestamp int64
//}
//
//const (
//	StakingLive     = "LIVE"
//	StakingImmature = "IMMATURE"
//	StakingExpired  = "EXPIRED"
//	StakingRevoked  = "REVOKED"
//	StakingVoted    = "VOTED"
//)
//
//func (wl *WalletLoad) StakingOverviewAllWallets() *dcrlibwallet.StakingOverview {
//	overview := new(dcrlibwallet.StakingOverview)
//	for _, w := range wl.MultiWallet.AllWallets() {
//		ov, _ := w.StakingOverview()
//		overview.All += ov.All
//		overview.Expired += ov.Expired
//		overview.Immature += ov.Immature
//		overview.Live += ov.Live
//		overview.Revoked += ov.Revoked
//		overview.Voted += ov.Voted
//	}
//
//	wl.MultiWallet.GetLowestBlockTimestamp()
//	return overview
//}
//
//func calculateDaysBehind(lastHeaderTime int64) string {
//	diff := time.Since(time.Unix(lastHeaderTime, 0))
//	daysBehind := int(math.Round(diff.Hours() / 24))
//	if daysBehind < 1 {
//		return "<1 day"
//	} else if daysBehind == 1 {
//		return "1 day"
//	} else {
//		return fmt.Sprintf("%d days", daysBehind)
//	}
//}
//
//func filterToStatus(txFilter int32) string {
//	switch txFilter {
//	case dcrlibwallet.TxFilterImmature:
//		return StakingImmature
//	case dcrlibwallet.TxFilterLive:
//		return StakingLive
//	case dcrlibwallet.TxFilterExpired:
//		return StakingExpired
//	case dcrlibwallet.TxFilterRevoked:
//		return StakingRevoked
//	case dcrlibwallet.TxFilterVoted:
//		return StakingVoted
//	}
//
//	return ""
//}
//
//func transactionToTicket(tx dcrlibwallet.Transaction, status, walletName string) Ticket {
//	return Ticket{
//		Status:     status,
//		Amount:     dcrutil.Amount(tx.Amount).String(),
//		DateTime:   time.Unix(tx.Timestamp, 0).Format("Jan 2, 2006 03:04:05 PM"),
//		MonthDay:   time.Unix(tx.Timestamp, 0).Format("Jan 2"),
//		DaysBehind: calculateDaysBehind(tx.Timestamp),
//		Fee:        dcrutil.Amount(tx.Fee).String(),
//		WalletName: walletName,
//		timestamp:  tx.Timestamp,
//	}
//}
//
//func transactionsToTickets(txs []dcrlibwallet.Transaction, status, walletName string) []Ticket {
//	var tickets []Ticket
//	for _, tx := range txs {
//		tickets = append(tickets, transactionToTicket(tx, status, walletName))
//	}
//
//	return tickets
//}
//
//func (wl *WalletLoad) getAllTickets(walletID int, newestFirst bool) ([]Ticket, error) {
//	w := wl.MultiWallet.WalletWithID(walletID)
//	var tickets []Ticket
//
//	addTickets := func(txFilter int32, newestFirst bool) error {
//		txs, err := w.GetTransactionsRaw(0, 0, txFilter, newestFirst)
//		if err != nil {
//			return err
//		}
//
//		tickets = append(tickets, transactionsToTickets(txs, filterToStatus(txFilter), w.Name)...)
//		return nil
//	}
//
//	filters := []int32{
//		dcrlibwallet.TxFilterImmature,
//		dcrlibwallet.TxFilterLive,
//		dcrlibwallet.TxFilterVoted,
//		dcrlibwallet.TxFilterExpired,
//		dcrlibwallet.TxFilterRevoked,
//	}
//
//	for _, filter := range filters {
//		err := addTickets(filter, newestFirst)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	sort.SliceStable(tickets, func(i, j int) bool {
//		if newestFirst {
//			return tickets[i].timestamp > tickets[j].timestamp
//		}
//		return tickets[i].timestamp < tickets[j].timestamp
//	})
//
//	return tickets, nil
//}
//
//func (wl *WalletLoad) GetTickets(walletID int, txFilter int32, newestFirst bool) ([]Ticket, error) {
//	if txFilter == dcrlibwallet.TxFilterStaking {
//		return wl.getAllTickets(walletID, newestFirst)
//	}
//
//	w := wl.MultiWallet.WalletWithID(walletID)
//	txs, err := w.GetTransactionsRaw(0, 0, txFilter, newestFirst)
//	if err != nil {
//		return nil, err
//	}
//
//	return transactionsToTickets(txs, filterToStatus(txFilter), w.Name), nil
//}
//
//func (wl *WalletLoad) AllLiveTickets() ([]Ticket, error) {
//	var tickets []Ticket
//	wallets := wl.MultiWallet.AllWallets()
//
//	liveTicketFilters := []int32{dcrlibwallet.TxFilterImmature, dcrlibwallet.TxFilterLive}
//	for _, w := range wallets {
//		for _, filter := range liveTicketFilters {
//			tx, err := w.GetTransactionsRaw(0, 0, filter, true)
//			if err != nil {
//				return tickets, err
//			}
//			tickets = append(tickets, transactionsToTickets(tx, filterToStatus(filter), w.Name)...)
//		}
//	}
//
//	return tickets, nil
//}

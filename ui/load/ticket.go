package load

import (
	"fmt"
	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"math"
	"time"
)

type Ticket struct {
	Status 	   string
	Fee        string
	Amount     string
	DateTime   string
	MonthDay   string
	DaysBehind string
	WalletName string
}

const(
	StakingLive     = "LIVE"
	StakingImmature = "IMMATURE"
)

func (wl *WalletLoad) StakingOverviewAllWallets() *dcrlibwallet.StakingOverview {
	overview := new(dcrlibwallet.StakingOverview)
	for _, w := range wl.MultiWallet.AllWallets() {
		ov, _ := w.StakingOverview()
		overview.All += ov.All
		overview.Expired += ov.Expired
		overview.Immature += ov.Immature
		overview.Live += ov.Live
		overview.Revoked += ov.Revoked
		overview.Voted += ov.Voted
	}

	return overview
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

func transactionToTicket(tx dcrlibwallet.Transaction, status, walletName string) Ticket {
	return Ticket{
		Status: status,
		Amount: dcrutil.Amount(tx.Amount).String(),
		DateTime: time.Unix(tx.Timestamp, 0).Format("Jan 2, 2006 03:04:05 PM"),
		MonthDay: time.Unix(tx.Timestamp, 0).Format("Jan 2"),
		DaysBehind: calculateDaysBehind(tx.Timestamp),
		Fee:       dcrutil.Amount(tx.Fee).String(),
		WalletName: walletName,
	}
}

func statusFromFilter(txFilter int32) string {
	switch txFilter {
	case dcrlibwallet.TxFilterImmature:
		return StakingImmature
	case dcrlibwallet.TxFilterLive:
		return StakingLive
	}

	return ""
}

func transactionsToTickets(txs []dcrlibwallet.Transaction, status, walletName string) []Ticket {
	var tickets []Ticket
	for _, tx := range txs {
		tickets = append(tickets, transactionToTicket(tx, status, walletName))
	}

	return tickets
}

func (wl *WalletLoad) GetTickets (walletID int, txFilter int32, newestFirst bool) ([]Ticket, error) {
	var tickets []Ticket

	w := wl.MultiWallet.WalletWithID(walletID)
	txs, err := w.GetTransactionsRaw(0, 0, txFilter, newestFirst)
	if err != nil {
		return tickets, err
	}

	return transactionsToTickets(txs, statusFromFilter(txFilter), w.Name), nil
}

func (wl *WalletLoad) AllLiveTickets() ([]Ticket, error) {
	var txs []Ticket
	wallets := wl.MultiWallet.AllWallets()

	for _, w := range wallets {
		immatureTx, err := w.GetTransactionsRaw(0, 0, dcrlibwallet.TxFilterImmature, true)
		if err != nil {
			return txs, err
		}

		txs = append(txs, transactionsToTickets(immatureTx, StakingImmature, w.Name)...)

		liveTxs, err := w.GetTransactionsRaw(0, 0, dcrlibwallet.TxFilterLive, true)
		if err != nil {
			return txs, err
		}

		txs = append(txs, transactionsToTickets(liveTxs, StakingLive, w.Name)...)
	}

	return txs, nil
}

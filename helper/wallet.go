package helper

import (
	"fmt"

	"github.com/decred/dcrd/hdkeychain"
	"github.com/decred/dcrwallet/walletseed"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet"
)

type (
	MultiWallet struct {
		*dcrlibwallet.MultiWallet
		WalletIDs []int
	}
)

func LoadWallet(appDataDir, netType string) (*MultiWallet, bool, bool, error) {
	multiWallet, err := dcrlibwallet.NewMultiWallet(appDataDir, "", netType)
	if err != nil {
		return nil, false, false, fmt.Errorf("Initialization error: %v", err)
	}

	mw := &MultiWallet{
		MultiWallet: multiWallet,
	}
	mw.WalletIDs = make([]int, 0)
	mw.WalletIDs = append(mw.WalletIDs, mw.OpenedWalletIDsRaw()...)

	if multiWallet.LoadedWalletsCount() == 0 {
		return mw, true, false, nil
	}

	err = multiWallet.OpenWallets(nil)
	if err != nil {
		return mw, false, false, fmt.Errorf("Error opening wallet db: %v", err)
	}

	for i := range mw.WalletIDs {
		fmt.Println(mw.WalletWithID(mw.WalletIDs[i]).WalletOpened())
	}

	err = multiWallet.SpvSync()
	if err != nil {
		return mw, false, false, fmt.Errorf("Spv sync attempt failed: %v", err)
	}

	return mw, false, false, nil
}

func (w *MultiWallet) RegisterWalletID(wID int) {
	for _, v := range w.WalletIDs {
		if v == wID {
			return
		}
	}

	w.WalletIDs = append(w.WalletIDs, wID)
	// TODO return and handle wallet is already registered error
}

func (w *MultiWallet) TotalBalance() (string, error) {
	var totalBalance int64

	for _, walletID := range w.WalletIDs {
		accounts, err := w.WalletWithID(walletID).GetAccountsRaw(dcrlibwallet.DefaultRequiredConfirmations)
		if err != nil {
			return "0", err
		}

		for _, account := range accounts.Acc {
			totalBalance += account.TotalBalance
		}
	}

	return dcrutil.Amount(totalBalance).String(), nil
}

func GenerateSeedWords() (string, error) {
	// generate seed
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		return "", fmt.Errorf("\nError generating seed for new wallet: %s.", err)
	}
	return walletseed.EncodeMnemonic(seed), nil
}

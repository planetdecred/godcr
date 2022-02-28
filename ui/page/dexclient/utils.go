package dexclient

import (
	"fmt"
	"sort"
	"strconv"

	"decred.org/dcrdex/client/asset"
	"decred.org/dcrdex/client/asset/btc"
	"decred.org/dcrdex/client/asset/dcr"
	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/dex"
	"github.com/planetdecred/godcr/ui/values"
)

// TODO: add localizable support for all these strings values
const (
	strLogin                    = "Login"
	strStartUseDex              = "Start using now"
	strStartSyncToUse           = "Start sync to continue"
	strAppPassword              = "App Password"
	strConfirmPassword          = "Confirm Password"
	strWalletPassword           = "Wallet Password"
	strSelectAccountForDex      = "Select DCR account to use with DEX"
	strWaitingConfirms          = "Waiting for confirmations..."
	strSetAppPassword           = "Set App Password"
	strWalletSetting            = "Wallets Settings"
	strSellectWallet            = "Sellect Wallet"
	strNoWallet                 = "No wallet"
	strReady                    = "Ready"
	strLocked                   = "Locked"
	strOff                      = "Off"
	strWithdraw                 = "Withdraw"
	strDeposit                  = "Deposit"
	strLock                     = "Lock"
	strUnLock                   = "Unlock"
	strSuccessful               = "Successfully!"
	strAddress                  = "Address"
	strAmount                   = "Amount"
	strDexAddr                  = "DEX Address"
	strSubmit                   = "Submit"
	strPickAServer              = "Pick a Server"
	strCustomServer             = "Custom Server"
	strAddADex                  = "Add a dex"
	strAddA                     = "Add a"
	strTLSCert                  = "TLS Certificate"
	strAddAFile                 = "Add a file"
	strChooseOtherFile          = "Choose other file"
	strRegister                 = "Register"
	strConfirmReg               = "Confirm Registration"
	strRequireWalletPayFee      = "Your wallet is required to pay registration fees."
	strConfirmSelectAssetPayFee = "How will you pay the registration fee?"
	strSetupNeeded              = "Setup Needed"
	strWalletReady              = "Wallet Ready"
	strMarket                   = "Market"
	strAllMarketAt              = "All markets at"
	strLotSize                  = "Lot Size"
	strInitDexPasswordDesc      = "Set your app password. This password will protect your DEX account keys and connected wallets."

	nStrNoWalletFound        = "no wallet found for %s"
	nStrWithdrawErr          = "withdraw error: %s"
	nStrNameWallet           = "%s Wallet"
	nStrAmountAvailable      = "%s available"
	nStrCreateAWallet        = "Create a %s Wallet"
	nStrUnlockWall           = "Unlock %s Wallet"
	nStrAlreadyConnectWallet = "Already connected a %s wallet"
	nStrConnHostError        = "Connection to dex server %s failed. You can close app and try again later or wait for it to reconnect"
	nStrConfirmationsStatus  = "In order to trade at %s, the registration fee payment needs %d confirmations."
	nStrNumberConfirmations  = "%d confirmations"
)

// supportedMarket check supported market for app depend on dcrlibwallet.
// TODO: update the logic or remove this when supported all markets.
func supportedMarket(mkt *core.Market) bool {
	// dcr_btc
	if mkt.BaseID == dcr.BipID && mkt.QuoteID == btc.BipID {
		return true
	}
	// btc_dcr
	if mkt.QuoteID == dcr.BipID && mkt.BaseID == btc.BipID {
		return true
	}
	return false
}

func formatAmount(amount uint64, unitInfo *dex.UnitInfo) string {
	convertedAmount := float64(amount) / float64(unitInfo.Conventional.ConversionFactor)
	return strconv.FormatFloat(convertedAmount, 'f', -1, 64)
}

func formatAmountUnit(assetID uint32, assetName string, amount uint64) string {
	assetInfo, err := asset.Info(assetID)
	if err != nil {
		return fmt.Sprintf("%d [%s units]", amount, assetName)
	}
	unitInfo := assetInfo.UnitInfo
	convertedLotSize := formatAmount(amount, &unitInfo)
	return fmt.Sprintf("%s %s", convertedLotSize, unitInfo.Conventional.Unit)
}

// sortExchanges convert map cert into a sorted slice
func sortExchanges(mapCert map[string][]byte) []string {
	servers := make([]string, 0, len(mapCert))
	for host := range mapCert {
		servers = append(servers, host)
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i] < servers[j]
	})
	return servers
}

// sortFeeAsset convert map FeeAsset into a sorted slice
func sortFeeAsset(mapFeeAsset map[string]*core.FeeAsset) []*core.FeeAsset {
	feeAssets := make([]*core.FeeAsset, 0, len(mapFeeAsset))
	for _, feeAsset := range mapFeeAsset {
		feeAssets = append(feeAssets, feeAsset)
	}
	sort.Slice(feeAssets, func(i, j int) bool {
		return feeAssets[i].ID < feeAssets[j].ID
	})
	return feeAssets
}

// sortServers convert mapExchanges into a sorted slice
func sortServers(mapExchanges map[string]*core.Exchange) []*core.Exchange {
	exchanges := make([]*core.Exchange, 0, len(mapExchanges))
	for _, dexServer := range mapExchanges {
		exchanges = append(exchanges, dexServer)
	}
	sort.Slice(exchanges, func(i, j int) bool {
		return exchanges[i].Host < exchanges[j].Host
	})
	return exchanges
}

func walletSyncPercentage(wall *core.WalletState) string {
	var percentage float32 = 0
	if wall.SyncProgress != 1 {
		percentage = wall.SyncProgress * 100
	}
	return fmt.Sprintf("%.2f%% %s", percentage, values.String(values.StrSynced))
}

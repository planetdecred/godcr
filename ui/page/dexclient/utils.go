package dexclient

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"decred.org/dcrdex/client/asset"
	"decred.org/dcrdex/client/asset/btc"
	"decred.org/dcrdex/client/asset/dcr"
	"decred.org/dcrdex/client/core"
	"decred.org/dcrdex/client/db"
	"decred.org/dcrdex/dex"
	"decred.org/dcrdex/dex/order"
	"gioui.org/text"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const (
	aYear   = 31536000000
	aMonth  = 2592000000
	aDay    = 86400000
	anHour  = 3600000
	aMinute = 60000
)

// TODO: add localizable support for all these strings values
const (
	strSetAppPassword               = "Set App Password"
	strConfirmPassword              = "Confirm Password"
	strWalletPassword               = "Wallet Password"
	strAppPassword                  = "App Password"
	strStartUseDex                  = "Start using now"
	strLogin                        = "Login"
	strAddADex                      = "Add a dex"
	strStartSyncToUse               = "Start sync to continue"
	strOrderHistory                 = "Order History"
	strDexSetting                   = "Dex Settings"
	strTrade                        = "Trade"
	strSide                         = "Side"
	strRate                         = "Rate"
	strQuantity                     = "Quantity"
	strFilled                       = "Filled"
	strSettled                      = "Settled"
	strStatus                       = "Status"
	strTime                         = "Time"
	strRegister                     = "Register"
	strConfirmReg                   = "Confirm Registration"
	strSubmit                       = "Submit"
	strDexAddr                      = "DEX Address"
	strPickAServer                  = "Pick a Server"
	strCustomServer                 = "Custom Server"
	strAddAFile                     = "Add a file"
	strChooseOtherFile              = "Choose other file"
	strNoWallet                     = "No wallet"
	strReady                        = "Ready"
	strLocked                       = "Locked"
	strOff                          = "Off"
	strWithdraw                     = "Withdraw"
	strDeposit                      = "Deposit"
	strLock                         = "Lock"
	strUnLock                       = "Unlock"
	strIHave                        = "I have"
	strIGet                         = "I get"
	strSuccessful                   = "Successfully!"
	strOk                           = "Ok"
	strAddress                      = "Address"
	strAmount                       = "Amount"
	strSellectWallet                = "Sellect Wallet"
	strSetupNeeded                  = "Setup Needed"
	strWalletReady                  = "Wallet Ready"
	strMarket                       = "Market"
	strAllMarketAt                  = "All markets at"
	strLotSize                      = "Lot Size"
	strTLSCert                      = "TLS Certificate"
	strAddA                         = "Add a"
	strImportAccount                = "Import Account"
	strExportAccount                = "Export Account"
	strDisableAccount               = "Disable Account"
	strChangeAppPassword            = "Change App Password"
	strWaitingConfirms              = "Waiting for confirmations..."
	strAuthOrderAppPassword         = "Authorize this order with your app password."
	strRequireWalletPayFee          = "Your wallet is required to pay registration fees."
	strConfirmSelectAssetPayFee     = "How will you pay the registration fee?"
	strErrInputAppPassword          = "Please input your application password"
	strInitDexPasswordDesc          = "Set your app password. This password will protect your DEX account keys and connected wallets."
	strDexClientSupportSimultaneous = "The Decred DEX Client supports simultaneous use of any number of DEX servers."
	strAuthorizeExport              = "Authorize Export"
	strAuthorizeImport              = "Authorize Import"
	strPasswordConfirmAcctImport    = "Enter your app password to confirm Account import."
	strCurrentPassword              = "Current Password"
	strNewPassword                  = "New Password"
	strConfirmNewPassword           = "Confirm new password"
	strErrRequireCertFile           = "Please choose a cert file"
	strNoteConfirmTradeMessage      = "IMPORTANT: Trades take time to settle, and you cannot turn off the DEX client software, or the BTC or DCR blockchain and/or wallet software, until settlement is complete. Settlement can complete as quickly as a few minutes or take as long as a few hours."
	strRecentActivity               = "Recent Activity"

	nStrNoWalletFound         = "no wallet found for %s"
	nStrCreateAWallet         = "Create a %s Wallet"
	nStrUnlockWall            = "Unlock %s Wallet"
	nStrWithdrawErr           = "withdraw error: %s"
	nStrAmountAvailable       = "%s available"
	nStrConfirmationsStatus   = "In order to trade at %s, the registration fee payment needs %d confirmations."
	nStrAlreadyConnectWallet  = "Already connected a %s wallet"
	nStrNameWallet            = "%s Wallet"
	nStrNumberConfirmations   = "%d confirmations"
	nStrConnHostError         = "Connection to dex server %s failed. You can close app and try again later or wait for it to reconnect"
	nStrAddressDex            = "Address DEX: %s"
	nStrAccountID             = "Account ID: %s"
	nStrConfirmDisableAccount = "Enter your app password to disable account: %s \n\nNote: This action is irreversible - once an account is disabled it can't be re-enabled."
	nStrConfirmExportAccount  = "Enter your app password to confirm Account export for: %s"
)

func sellString(ord *core.Order) string {
	if ord.Sell {
		return "sell"
	}

	return "buy"
}

func typeString(ord *core.Order) string {
	if ord.Type != order.LimitOrderType {
		return "market"
	}

	if ord.TimeInForce == order.ImmediateTiF {
		return "limit (i)"
	}

	return "limit"
}

func rateString(ord *core.Order) string {
	if ord.Type == order.MarketOrderType {
		return "market"
	}
	return formatCoinValue(ord.Rate)
}

// formatCoinValue formats the asset value to a string.
// TODO: should calculate base on Conventional.ConversionFactor
func formatCoinValue(val uint64) string {
	return fmt.Sprintf("%.6f", float64(val)/1e8)
}

// timeSince returns a string representation of the duration since the specified
// unix timestamp.
func timeSince(t uint64) string {
	var seconds = float64(time.Now().Unix()*1000 - int64(t))

	var result = ""
	var count = 0

	add := func(n float64, s string) bool {
		if n > 0 || count > 0 {
			count++
		}
		if n > 0 {
			result += fmt.Sprintf("%d %s ", int(n), s)
		}
		return count >= 2
	}

	var y, mo, d, h, m, s float64

	y, seconds = timeMod(seconds, aYear)
	if add(y, "y") {
		return result
	}
	mo, seconds = timeMod(seconds, aMonth)
	if add(mo, "mo") {
		return result
	}
	d, seconds = timeMod(seconds, aDay)
	if add(d, "d") {
		return result
	}
	h, seconds = timeMod(seconds, anHour)
	if add(h, "h") {
		return result
	}
	m, seconds = timeMod(seconds, aMinute)
	if add(m, "m") {
		return result
	}
	s, _ = timeMod(seconds, 1000)
	add(s, "s")
	if result == "" {
		return "0 s"
	}
	return result
}

// timeMod returns the quotient and remainder of t / dur.
func timeMod(t float64, dur float64) (float64, float64) {
	n := math.Floor(t / dur)
	return n, t - n*dur
}

// isMarketBuy will return true if the order is a market buy order.
func isMarketBuy(ord *core.Order) bool {
	return ord.Type == order.MarketOrderType && !ord.Sell
}

// settled sums the quantities of the matches that have completed.
func settled(ord *core.Order) float64 {
	if ord.Matches == nil {
		return 0
	}
	var qty func(m *core.Match) float64

	if isMarketBuy(ord) {
		qty = func(m *core.Match) float64 {
			return (float64(m.Qty*m.Rate) * 1e-8)
		}
	} else {
		qty = func(m *core.Match) float64 {
			return float64(m.Qty)
		}
	}

	var settle float64 = 0
	for _, match := range ord.Matches {
		if match.IsCancel {
			continue
		}
		redeemed := (match.Side == order.Maker && match.Status >= order.MakerRedeemed) ||
			(match.Side == order.Taker && match.Status >= order.MatchComplete)

		if redeemed {
			settle += settle + qty(match)
		} else {
			settle += settle
		}
	}

	return settle
}

// hasLiveMatches returns true if the order has matches that have not completed
// settlement yet.
func hasLiveMatches(ord *core.Order) bool {
	if ord.Matches == nil {
		return false
	}

	for _, m := range ord.Matches {
		if !m.Revoked && m.Status < order.MakerRedeemed {
			return true
		}
	}
	return false
}

// statusString converts the order status to a string
func statusString(ord *core.Order) string {
	isLive := hasLiveMatches(ord)

	switch ord.Status {
	case order.OrderStatusUnknown:
		return "unknown"
	case order.OrderStatusEpoch:
		return "epoch"
	case order.OrderStatusBooked:
		if ord.Cancelling {
			return "cancelling"
		}
		return "booked"
	case order.OrderStatusExecuted:
		if isLive {
			return "settling"
		}
		return "executed"
	case order.OrderStatusCanceled:
		if isLive {
			return "canceled/settling"
		}
		return "canceled"
	case order.OrderStatusRevoked:
		if isLive {
			return "revoked/settling"
		}
		return "revoked"
	}

	return "unknown"
}

func marketIDToAsset(marketID string) (baseInfo *asset.WalletInfo, quoteInfo *asset.WalletInfo, err error) {
	mktIDs := strings.Split(marketID, "_")
	baseID, ok := dex.BipSymbolID(mktIDs[0])
	if !ok {
		return nil, nil, errors.New("Invalid market")
	}
	b, err := asset.Info(baseID)
	if err != nil {
		return nil, nil, err
	}

	quoteID, ok := dex.BipSymbolID(mktIDs[1])
	if !ok {
		return nil, nil, errors.New("Invalid market")
	}
	q, err := asset.Info(quoteID)
	if err != nil {
		return nil, nil, err
	}

	return b, q, nil
}

func minMaxRateOrderBook(orders []*core.MiniOrder) (float64, float64) {
	if len(orders) == 0 {
		return 0, 0
	}
	max := orders[0].Rate
	min := orders[0].Rate
	for _, value := range orders {
		if max < value.Rate {
			max = value.Rate
		}
		if min > value.Rate {
			min = value.Rate
		}
	}
	return min, max
}

// sliceExchanges convert mapExchanges into a sorted slice
func sliceExchanges(mapExchanges map[string]*core.Exchange) []*core.Exchange {
	exchanges := make([]*core.Exchange, 0, len(mapExchanges))
	for _, dexServer := range mapExchanges {
		exchanges = append(exchanges, dexServer)
	}
	sort.Slice(exchanges, func(i, j int) bool {
		return exchanges[i].Host < exchanges[j].Host
	})
	return exchanges
}

// sliceMarkets convert mapMarkets into a sorted slice
func sliceMarkets(mapMarkets map[string]*core.Market) []*core.Market {
	markets := make([]*core.Market, 0, len(mapMarkets))
	for _, market := range mapMarkets {
		markets = append(markets, market)
	}
	sort.Slice(markets, func(i, j int) bool {
		return markets[i].Name < markets[j].Name
	})
	return markets
}

// sliceFeeAsset convert map FeeAsset into a sorted slice
func sliceFeeAsset(mapFeeAsset map[string]*core.FeeAsset) []*core.FeeAsset {
	feeAssets := make([]*core.FeeAsset, 0, len(mapFeeAsset))
	for _, feeAsset := range mapFeeAsset {
		feeAssets = append(feeAssets, feeAsset)
	}
	sort.Slice(feeAssets, func(i, j int) bool {
		return feeAssets[i].ID < feeAssets[j].ID
	})
	return feeAssets
}

// sliceSever convert map cert into a sorted slice
func sliceSever(mapCert map[string][]byte) []string {
	servers := make([]string, 0, len(mapCert))
	for host := range mapCert {
		servers = append(servers, host)
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i] < servers[j]
	})
	return servers
}

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

func formatAmountUnit(assetID uint32, assetName string, amount uint64) string {
	assetInfo, err := asset.Info(assetID)
	if err != nil {
		return fmt.Sprintf("%d [%s units]", amount, assetName)
	}
	unitInfo := assetInfo.UnitInfo
	convertedLotSize := formatAmount(amount, &unitInfo)
	return fmt.Sprintf("%s %s", convertedLotSize, unitInfo.Conventional.Unit)
}

func formatAmount(amount uint64, unitInfo *dex.UnitInfo) string {
	convertedAmount := float64(amount) / float64(unitInfo.Conventional.ConversionFactor)
	return strconv.FormatFloat(convertedAmount, 'f', -1, 64)
}

func walletSyncPercentage(wall *core.WalletState) string {
	var percentage float32 = 0
	if wall.SyncProgress != 1 {
		percentage = wall.SyncProgress * 100
	}
	return fmt.Sprintf("%.2f%% %s", percentage, values.String(values.StrSynced))
}

type notification struct {
	Acked     bool
	Details   string
	ID        string
	Severity  db.Severity
	TimeStamp uint64
	Subject   string
}

// saveNotification save one notifications to local db
func saveNotification(mw *dcrlibwallet.MultiWallet, n *core.Notification) {
	notifications := make(map[string]*notification)
	err := mw.ReadUserConfigValue(dexNotificationConfigKey, &notifications)
	if err == nil {
		notifications[(*n).ID().String()] = &notification{
			Acked:     (*n).Acked(),
			Details:   (*n).Details(),
			ID:        (*n).ID().String(),
			Severity:  (*n).Severity(),
			TimeStamp: (*n).Time(),
			Subject:   (*n).Subject(),
		}
		mw.SaveUserConfigValue(dexNotificationConfigKey, &notifications)
	}
}

// saveNotifications save all notifications to local db
func saveNotifications(mw *dcrlibwallet.MultiWallet, notifications []*db.Notification) {
	mapNotifications := make(map[string]*notification, len(notifications))
	for _, n := range notifications {
		mapNotifications[n.Id.String()] = &notification{
			Acked:     n.Acked(),
			Details:   n.Details(),
			ID:        n.ID().String(),
			Severity:  n.Severity(),
			TimeStamp: n.Time(),
			Subject:   n.Subject(),
		}
	}
	mw.SaveUserConfigValue(dexNotificationConfigKey, &mapNotifications)
}

// getNotifications get all notifications from local db with given parameters
func getNotifications(mw *dcrlibwallet.MultiWallet, onlyUnAcked bool) []*notification {
	mapNotifications := make(map[string]*notification)

	err := mw.ReadUserConfigValue(dexNotificationConfigKey, &mapNotifications)
	if err != nil {
		return nil
	}

	notifications := make([]*notification, 0)
	for _, ntfn := range mapNotifications {
		if onlyUnAcked {
			if !ntfn.Acked {
				notifications = append(notifications, ntfn)
			}
			continue
		}
		ntfn.Acked = true
		mapNotifications[ntfn.ID] = ntfn
		notifications = append(notifications, ntfn)
	}

	if !onlyUnAcked {
		mw.SaveUserConfigValue(dexNotificationConfigKey, &mapNotifications)
	}

	sort.Slice(notifications, func(i, j int) bool {
		return notifications[i].TimeStamp > notifications[j].TimeStamp
	})

	return notifications
}

func severityColor(severity db.Severity, colors *values.Color) color.NRGBA {
	switch severity {
	case db.Success:
		return colors.Success
	case db.WarningLevel:
		return colors.Yellow
	case db.ErrorLevel:
		return colors.Danger
	default:
		return colors.Background
	}
}

func severityIndicatorLabel(mw *dcrlibwallet.MultiWallet, theme *decredmaterial.Theme) (lb decredmaterial.Label) {
	notifications := getNotifications(mw, true)
	var severity db.Severity = 0
	if len(notifications) > 0 {
		severity = notifications[0].Severity
	}
	lb = theme.Label(values.TextSize12, "")
	lb.Font.Weight = text.Bold
	lb.Color = severityColor(severity, theme.Color)

	if len(notifications) > 0 {
		lb.Text = fmt.Sprintf("%d", len(notifications))
		return
	}

	if len(notifications) > 99 {
		lb.Text = "99+"
	}
	return
}

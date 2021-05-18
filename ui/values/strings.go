package values

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/planetdecred/godcr/ui/values/localizable"
)

var rex = regexp.MustCompile(`(?m)("(?:\\.|[^"\\])*")\s+=\s+("(?:\\.|[^"\\])*")`) // "key"="value"
var Languages = []string{"en"}

const DefaultLanguge = "en"

var en map[string]string

func init() {
	en = make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(localizable.EN))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		kv := rex.FindAllStringSubmatch(line, -1)[0]
		key := trimQuotes(kv[1])
		value := trimQuotes(kv[2])

		en[key] = value
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("localizable: error reading scanner input:", err)
	}
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func String(key string) string {
	str, ok := en[key]
	if !ok {
		return ""
	}

	return str
}

func StringF(key string, a ...interface{}) string {
	str := String(key)
	if str == "" {
		return str
	}

	return fmt.Sprintf(str, a...)
}

const (
	StrAppName                   = "app_name"
	StrAppTitle                  = "app_title"
	StrRecentTransactions        = "recentTransactions"
	StrSeeAll                    = "seeAll"
	StrSend                      = "send"
	StrReceive                   = "receive"
	StrOnline                    = "online"
	StrOffline                   = "offline"
	StrShowDetails               = "showDetails"
	StrHideDetails               = "hideDetails"
	StrConnectedPeersCount       = "connectedPeersCount"
	StrNoConnectedPeers          = "noConnectedPeer"
	StrDisconnect                = "disconnect"
	StrReconnect                 = "reconnect"
	StrCurrentTotalBalance       = "currentTotalBalance"
	StrWalletStatus              = "walletStatus"
	StrBlockHeadersFetched       = "blockHeadersFetched"
	StrNoTransactions            = "noTransactions"
	StrHeadersFetchProgress      = "headersFetchProgress"
	StrSyncSteps                 = "syncSteps"
	StrScanningTotalHeaders      = "scanningTotalHeaders"
	StrConnectedTo               = "connectedTo"
	StrWalletSynced              = "walletSynced"
	StrSynchronizing             = "synchronizing"
	StrWalletNotSynced           = "walletNotSynced"
	StrCancel                    = "cancel"
	StrUnlockToResumeRestoration = "unlockToResumeRestoration"
	StrUnlock                    = "unlock"
	StrSyncingProgress           = "syncingProgress"
	StrNoWalletLoaded            = "noWalletLoaded"
	StrLastBlockHeight           = "lastBlockHeight"
)

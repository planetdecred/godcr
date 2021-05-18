package values

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/planetdecred/godcr/ui/values/localizable"
)

var rex = regexp.MustCompile(`(?m)("(?:\\.|[^"\\])*")\s*=\s*("(?:\\.|[^"\\])*")`) // "key"="value"
var Languages = []string{"en", "zh"}
var UserLanguages = []string{"en"} // order of preference

const (
	DefaultLanguge = "en"
	commentPrefix  = "/"
)

var en map[string]string
var zh map[string]string

var languageStrings map[string]map[string]string

func init() {

	readIntoMap := func(m map[string]string, localizableStrings string) {
		scanner := bufio.NewScanner(strings.NewReader(localizableStrings))
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, commentPrefix) {
				continue
			}

			matches := rex.FindAllStringSubmatch(line, -1)
			if len(matches) == 0 {
				continue
			}

			kv := matches[0]
			key := trimQuotes(kv[1])
			value := trimQuotes(kv[2])

			m[key] = value
		}
	}

	en = make(map[string]string)
	zh = make(map[string]string)
	languageStrings = make(map[string]map[string]string)

	readIntoMap(en, localizable.EN)
	languageStrings["en"] = en
	readIntoMap(zh, localizable.ZH)
	languageStrings["zh"] = zh
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
	for _, lang := range UserLanguages {
		languageMap := languageStrings[lang]
		str, ok := languageMap[key]
		if ok {
			return str
		}
	}

	return ""
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

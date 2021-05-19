package values

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/planetdecred/godcr/ui/values/localizable"
)

const (
	DefaultLanguge = localizable.ENGLISH
	commentPrefix  = "/"
)

var rex = regexp.MustCompile(`(?m)("(?:\\.|[^"\\])*")\s*=\s*("(?:\\.|[^"\\])*")`) // "key"="value"
var Languages = []string{localizable.ENGLISH, localizable.CHINESE, localizable.FRENCH}
var UserLanguages = []string{DefaultLanguge} // order of preference

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

	en := make(map[string]string)
	zh := make(map[string]string)
	fr := make(map[string]string)
	languageStrings = make(map[string]map[string]string)

	readIntoMap(en, localizable.EN)
	languageStrings[localizable.ENGLISH] = en

	readIntoMap(zh, localizable.ZH)
	languageStrings[localizable.CHINESE] = zh

	readIntoMap(fr, localizable.FR)
	languageStrings[localizable.FRENCH] = fr
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
	StrAppName                     = "appName"
	StrSend                        = "send"
	StrReceive                     = "receive"
	StrUnlock                      = "unlock"
	StrWalletStatus                = "walletStatus"
	StrFetchingBlockHeaders        = "fetchingBlockHeaders"
	StrSyncingState                = "syncingState"
	StrResumeAccountDiscoveryTitle = "resumeAccountDiscoveryTitle"
	StrHideDetails                 = "hideDetails"
	StrSyncSteps                   = "syncSteps"
	StrBlockHeaderFetchedCount     = "blockHeaderFetchedCount"
	StrConnectedTo                 = "connectedTo"
	StrSynced                      = "synced"
	StrNoWalletLoaded              = "noWalletLoaded"
	StrReconnect                   = "reconnect"
	StrWalletNotSynced             = "walletNotSynced"
	StrDisconnect                  = "disconnect"
	StrSyncingProgress             = "syncingProgress"
	StrBlockHeaderFetched          = "blockHeaderFetched"
	StrNoTransactionsYet           = "noTransactionsYet"
	StrCancel                      = "cancel"
	StrAppTitle                    = "appTitle"
	StrSeeAll                      = "seeAll"
	StrOnline                      = "online"
	StrConnectedPeersCount         = "connectedPeersCount"
	StrNoConnectedPeer             = "noConnectedPeer"
	StrCurrentTotalBalance         = "currentTotalBalance"
	StrRecentTransactions          = "recentTransactions"
	StrOffline                     = "offline"
	StrShowDetails                 = "showDetails"
	StrLastBlockHeight             = "lastBlockHeight"
	StrAgo                         = "ago"
	StrNewest                      = "newest"
	StrOldest                      = "oldest"
	StrAll                         = "all"
	StrTransferred                 = "transferred"
	StrSent                        = "sent"
	StrReceived                    = "received"
	StrYourself                    = "yourself"
	StrStaking                     = "staking"
	StrNConfirmations              = "nConfirmations"
	StrFrom                        = "from"
	StrTo                          = "to"
	StrFee                         = "fee"
	StrIncludedInBlock             = "includedInBlock"
	StrType                        = "type"
	StrTransactionId               = "transactionId"
	StrXInputsConsumed             = "xInputsConsumed"
	StrXOutputCreated              = "xOutputCreated"
	StrViewOnDcrdata               = "viewOnDcrdata"
)

package values

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/planetdecred/godcr/ui/values/localizable"
)

const (
	DefaultLangauge = localizable.ENGLISH
	commentPrefix   = "/"
)

var rex = regexp.MustCompile(`(?m)("(?:\\.|[^"\\])*")\s*=\s*("(?:\\.|[^"\\])*")`) // "key"="value"
var Languages = []string{localizable.ENGLISH, localizable.CHINESE, localizable.FRENCH, localizable.SPANISH}
var UserLanguages = []string{DefaultLangauge} // order of preference

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
	es := make(map[string]string)
	languageStrings = make(map[string]map[string]string)

	readIntoMap(en, localizable.EN)
	languageStrings[localizable.ENGLISH] = en

	readIntoMap(zh, localizable.ZH)
	languageStrings[localizable.CHINESE] = zh

	readIntoMap(fr, localizable.FR)
	languageStrings[localizable.FRENCH] = fr

	readIntoMap(es, localizable.ES)
	languageStrings[localizable.SPANISH] = es
}

func hasLanguage(language string) bool {
	for _, lang := range Languages {
		if lang == language {
			return true
		}
	}

	return false
}

func SetUserLanguage(lang string) {
	if hasLanguage(lang) {
		languages := []string{lang}
		if lang != DefaultLangauge {
			languages = append(languages, DefaultLangauge)
		}

		UserLanguages = languages
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
	StrUnlockWithPassword          = "unlockWithPassword"
	StrWalletStatus                = "walletStatus"
	StrFetchingBlockHeaders        = "fetchingBlockHeaders"
	StrDiscoveringWalletAddress    = "discoveringWalletAddress"
	StrRescanningHeaders           = "rescanningHeaders"
	StrRescanningBlocks            = "rescanningBlocks"
	StrBlocksScanned               = "blocksScanned"
	StrBlocksLeft                  = "blocksLeft"
	StrSync                        = "sync"
	StrAutoSyncInfo                = "autoSyncInfo"
	StrSyncingState                = "syncingState"
	StrWaitingState                = "waitingState"
	StrResumeAccountDiscoveryTitle = "resumeAccountDiscoveryTitle"
	StrHideDetails                 = "hideDetails"
	StrSyncSteps                   = "syncSteps"
	StrBlockHeaderFetchedCount     = "blockHeaderFetchedCount"
	StrTimeLeft                    = "timeLeft"
	StrConnectedTo                 = "connectedTo"
	StrConnecting                  = "connecting"
	StrSynced                      = "synced"
	StrNoWalletLoaded              = "noWalletLoaded"
	StrReconnect                   = "reconnect"
	StrWalletNotSynced             = "walletNotSynced"
	StrDisconnect                  = "disconnect"
	StrSyncingProgress             = "syncingProgress"
	StrSyncingProgressStat         = "syncingProgressStat"
	StrBlockHeaderFetched          = "blockHeaderFetched"
	StrNoTransactions              = "noTransactions"
	StrCancel                      = "cancel"
	StrAppTitle                    = "appTitle"
	StrSeeAll                      = "seeAll"
	StrOnline                      = "online"
	StrPeers                       = "peers"
	StrConnectedPeersCount         = "connectedPeersCount"
	StrNoConnectedPeer             = "noConnectedPeer"
	StrCurrentTotalBalance         = "currentTotalBalance"
	StrTotalBalance                = "totalBalance"
	StrRecentTransactions          = "recentTransactions"
	StrRecentProposals             = "recentProposals"
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
	StrMixed                       = "mixed"
	StrUmined                      = "unmined"
	StrImmature                    = "immature"
	StrVoted                       = "voted"
	StrRevoked                     = "revoked"
	StrLive                        = "live"
	StrExpired                     = "expired"
	StrPurchased                   = "purchased"
	StrRevocation                  = "revocation"
	StrStaking                     = "staking"
	StrImmatureRewards             = "immatureRewards"
	StrLockedByTickets             = "lockedByTickets"
	StrVotingAuthority             = "votingAuthority"
	StrImmatureStakeGen            = "immatureStakeGen"
	StrUnknown                     = "unknown"
	StrNConfirmations              = "nConfirmations"
	StrFrom                        = "from"
	StrTo                          = "to"
	StrFee                         = "fee"
	StrIncludedInBlock             = "includedInBlock"
	StrType                        = "type"
	StrTransactionID               = "transactionId"
	StrRebroadcast                 = "rebroadcast"
	StrXInputsConsumed             = "xInputsConsumed"
	StrXOutputCreated              = "xOutputCreated"
	StrViewOnDcrdata               = "viewOnDcrdata"
	StrViewProperty                = "viewProperty"
	StrAddNewAccount               = "addNewAccount"
	StrBackupSeedPhrase            = "backupSeedPhrase"
	StrCreateNewAccount            = "createNewAccount"
	StrInvalidPassphrase           = "invalidPassphrase"
	StrPasswordNotMatch            = "passwordNotMatch"
	StrCreate                      = "create"
	StrNotBackedUp                 = "notBackedUp"
	StrLabelSpendable              = "labelSpendable"
	StrSignMessage                 = "signMessage"
	StrStakeShuffle                = "stakeShuffle"
	StrRenameWalletSheetTitle      = "renameWalletSheetTitle"
	StrSettings                    = "settings"
	StrImportWatchingOnlyWallet    = "importWatchingOnlyWallet"
	StrVerifySeedInfo              = "verifySeedInfo"
	StrVerifyMessage               = "verifyMessage"
	StrMessage                     = "message"
	StrRename                      = "rename"
	StrCreateANewWallet            = "createANewWallet"
	StrWalletCreated               = "walletCreated"
	StrImportExistingWallet        = "importExistingWallet"
	StrWatchOnlyWallets            = "watchOnlyWallets"
	StrWatchOnlyWalletImported     = "watchOnlyWalletImported"
	StrImport                      = "Import"
	StrRescanProgressNotification  = "rescanProgressNotification"
	StrRemove                      = "remove"
	StrConfirm                     = "confirm"
	StrSpendingPasswordInfo        = "spendingPasswordInfo"
	StrSpendingPasswordInfo2       = "spendingPasswordInfo2"
	StrSpendingPassword            = "spendingPassphrase"
	StrEnterSpendingPassword       = "enterSpendingPassword"
	StrConfirmSpendingPassword     = "confirmSpendingPassword"
	StrCurrentSpendingPassword     = "currentSpendingPassword"
	StrNewSpendingPassword         = "newSpendingPassword"
	StrConfirmNewSpendingPassword  = "confirmNewSpendingPassword"
	StrSpendingPasswordUpdated     = "spendingPasswordUpdated"
	StrDangerZone                  = "dangerZone"
	StrNotConnected                = "notConnected"
	StrConfirmToRemove             = "confirmToRemove"
	StrChangeSpendingPass          = "changeSpendingPass"
	StrDebug                       = "debug"
	StrBeepForNewBlocks            = "beepForNewBlocks"
	StrRemoveWallet                = "removeWallet"
	StrChange                      = "change"
	StrRescan                      = "rescan"
	StrRescanInfo                  = "rescanInfo"
	StrNotifications               = "notifications"
	StrRescanBlockchain            = "rescanBlockchain"
	StrNewStartupPass              = "newStartupPass"
	StrConfirmNewStartupPass       = "confirmNewStartupPass"
	StrStartupPassConfirm          = "startupPassConfirm"
	StrSetupStartupPassword        = "setupStartupPassword"
	StrStartupPasswordInfo         = "startupPasswordInfo"
	StrConfirmStartupPass          = "confirmStartupPass"
	StrCurrentStartupPass          = "currentStartupPass"
	StrStartupPassword             = "startupPassword"
	StrStartupPasswordEnabled      = "startupPasswordEnabled"
	StrChangeSpecificPeer          = "changeSpecificPeer"
	StrLanguage                    = "language"
	StrConnection                  = "connection"
	StrCustomUserAgent             = "CustomUserAgent"
	StrConfirmRemoveStartupPass    = "confirmRemoveStartupPass"
	StrUserAgentDialogTitle        = "userAgentDialogTitle"
	StrSecurity                    = "security"
	StrUnconfirmedFunds            = "unconfirmedFunds"
	StrConfirmed                   = "confirmed"
	StrChangeStartupPassword       = "changeStartupPassword"
	StrConnectToSpecificPeer       = "connectToSpecificPeer"
	StrUserAgentSummary            = "userAgentSummary"
	StrGeneral                     = "general"
	StrChangeUserAgent             = "changeUserAgent"
	StrCreateStartupPassword       = "createStartupPassword"
	StrCurrencyConversion          = "currencyConversion"
	StrTransactions                = "transactions"
	StrWallets                     = "wallets"
	StrMore                        = "more"
	StrOverview                    = "overview"
	StrEnglish                     = "english"
	StrFrench                      = "french"
	StrSpanish                     = "spanish"
	StrUsdBittrex                  = "usdBittrex"
	StrNone                        = "none"
	StrProposal                    = "proposals"
	StrDex                         = "dex"
	StrGovernance                  = "governance"
	StrPending                     = "pending"
	StrVote                        = "vote"
	StrRevoke                      = "revoke"
	StrMaturity                    = "maturity"
	StrYesterday                   = "yesterday"
	StrDays                        = "days"
	StrHours                       = "hours"
	StrMinutes                     = "mins"
	StrSeconds                     = "secs"
	StrYearAgo                     = "yearAgo"
	StrYearsAgo                    = "yearsAgo"
	StrMonthAgo                    = "monthAgo"
	StrMonthsAgo                   = "monthsAgo"
	StrWeekAgo                     = "weekAgo"
	StrWeeksAgo                    = "weeksAgo"
	StrDayAgo                      = "dayAgo"
	StrDaysAgo                     = "daysAgo"
	StrHourAgo                     = "hourAgo"
	StrHoursAgo                    = "hoursAgo"
	StrMinuteAgo                   = "minuteAgo"
	StrMinutesAgo                  = "minutesAgo"
	StrJustNow                     = "justNow"
	StrAwareOfRisk                 = "imawareOfRisk"
	StrUnmixedBalance              = "unmixedBalance"
	StrBackupLater                 = "backupLater"
	StrBackupNow                   = "backupNow"
	StrStatus                      = "status"
	StrDaysToVote                  = "daysToVote"
	StrReward                      = "reward"
	StrViewTicket                  = "viewTicket"
	StrExternal                    = "external"
	StrRepublished                 = "republished"
	StrCopied                      = "copied"
	StrTxHashCopied                = "txHashCopied"
	StrAddressCopied               = "addressCopied"
	StrAddress                     = "address"
	StrAcctDetailsKey              = "acctDetailsKey"
	StrAcctNum                     = "acctNum"
	StrAcctName                    = "acctName"
	StrAcctRenamed                 = "accRenamed"
	StrAcctCreated                 = "acctCreated"
	StrRenameAcct                  = "renameAcct"
	StrHDPath                      = "hdPath"
	StrKey                         = "key"
	StrValidateWalSeed             = "validateWalSeed"
	StrClearAll                    = "clearAll"
	StrRestoreWallet               = "restoreWallet"
	StrRestore                     = "restore"
	StrRestoreExistingWallet       = "restoreExistingWallet"
	StrEnterSeedPhrase             = "enterSeedPhrase"
	StrInvalidSeedPhrase           = "invalidSeedPhrase"
	StrSeedAlreadyExist            = "seedAlreadyExist"
	StrWalletRestored              = "walletRestored"
	StrWalletExist                 = "walletExist"
	StrWalletNotExist              = "walletNotExist"
	StrInvalidHex                  = "invalidHex"
	StrRestoreWithHex              = "restoreWithHex"
	StrEnterHex                    = "enterHex"
	StrSubmit                      = "submit"
	StrEnterWalDetails             = "enterWalletDetails"
	StrCopy                        = "copy"
	StrHowToCopy                   = "howToCopy"
	StrEnterAddressToSign          = "enterAddressToSign"
	StrSignCopied                  = "signCopied"
	StrSignature                   = "signature"
	StrConfirmToSign               = "confirmToSign"
	StrEnterValidAddress           = "enterValidAddress"
	StrEnterValidMsg               = "enterValidMsg"
	StrInvalidAddress              = "invalidAddress"
	StrValidAddress                = "validAddress"
	StrAddrNotOwned                = "addrNotOwned"
	StrDeleted                     = "delete"
	StrWalletName                  = "walletName"
	StrEnterWalletName             = "enterWalletName"
	StrAddWallet                   = "addWallet"
	StrSelectWallet                = "selectWallet"
	StrWalletRenamed               = "walletRenamed"
	StrCheckMixerStatus            = "checkMixerStatus"
	StrWalletRestoreMsg            = "walletRestoreMsg"
	StrWalletRemoved               = "walletRemoved"
	StrWalletRemoveInfo            = "walletRemoveInfo"
	StrWatchOnlyWalletRemoveInfo   = "watchOnlyWalletRemoveInfo"
	StrGotIt                       = "gotIt"
	StrNoValidAccountFound         = "noValidAccountFound"
	StrMixer                       = "mixer"
	StrReadyToMix                  = "readyToMix"
	StrMixerRunning                = "mixerRunning"
	StrUseMixer                    = "useMixer"
	StrKeepAppOpen                 = "keepAppOpen"
	StrMixerShutdown               = "mixerShutdown"
	StrVotingPreference            = "votingPreference"
	StrNoAgendaYet                 = "noAgendaYet"
	StrFetchingAgenda              = "fetchingAgenda"
	StrUpdatePreference            = "updatePreference"
	StrApproved                    = "approved"
	StrVoting                      = "voting"
	StrRejected                    = "rejected"
	StrAbandoned                   = "abandoned"
	StrInDiscussion                = "inDiscussion"
	StrFetchingProposals           = "fetchingProposals"
	StrFetchProposals              = "fetchProposals"
	StrUnderReview                 = "underReview"
	StrNoProposals                 = "noProposal"
	StrWaitingAuthor               = "waitingForAuthor"
	StrWaitingForAdmin             = "waitingForAdmin"
	StrVoteTooltip                 = "voteTooltip"
	StrYes                         = "yes"
	StrNo                          = "no"
	StrTotalVotes                  = "totalVotes"
	StrTotalVotesReverse           = "totalVotesReverse"
	StrQuorumRequirement           = "quorumRequirement"
	StrDiscussions                 = "discussions"
	StrPublished                   = "published"
	StrToken                       = "token"
	StrProposalVoteDetails         = "proposalVoteDetails"
	StrVotingServiceProvider       = "votingServiceProvider"
	StrSelectVSP                   = "selectVSP"
	StrAddVSP                      = "addVSP"
	StrSave                        = "save"
	StrNoVSPLoaded                 = "noVSPLoaded"
	StrExtendedPubKey              = "extendedPubKey"
	StrEnterExtendedPubKey         = "enterXpubKey"
	StrXpubKeyErr                  = "xpubKeyErr"
	StrXpubWalletExist             = "xpubWalletExist"
	StrHint                        = "hint"
	StrAddAcctWarn                 = "addAcctWarn"
	StrVerifyMessageInfo           = "verifyMessageInfo"
	StrSetupMixerInfo              = "setupMixerInfo"
	StrTxdetailsInfo               = "txDetailsInfo"
	StrBackupInfo                  = "backupInfo"
	StrSignMessageInfo             = "signMessageInfo"
	StrPrivacyInfo                 = "privacyInfo"
	StrAllowUnspendUnmixedAcct     = "allowUnspendUnmixedAcct"
	StrBalToMaintain               = "balToMaintain"
	StrAutoTicketPurchase          = "autoTicketPurchase"
	StrPurchasingAcct              = "purchasingAcct"
	StrLocked                      = "locked"
	StrBalance                     = "balance"
	StrAllTickets                  = "allTickets"
	StrNoTickets                   = "noTickets"
	StrNoActiveTickets             = "noActiveTickets"
	StrLiveTickets                 = "liveTickets"
	StrTicketRecord                = "ticketRecord"
	StrRewardsEarned               = "rewardsEarned"
	StrNoReward                    = "noReward"
	StrLoadingPrice                = "loadingPrice"
	StrNotAvailable                = "notAvailable"
	StrTicketConfirmed             = "ticketConfirmed"
	StrBackStaking                 = "backStaking"
	StrTicketSettingSaved          = "ticketSettingSaved"
	StrAutoTicketWarn              = "autoTicketWarn"
	StrAutoTicketInfo              = "autoTicketInfo"
	StrConfirmPurchase             = "confirmPurchase"
	StrTicketError                 = "ticketError"
	StrWalletToPurchaseFrom        = "walletToPurchaseFrom"
	StrSelectedAccount             = "selectedAcct"
	StrBalToMaintainValue          = "balToMaintainValue"
	StrStake                       = "stake"
	StrTotal                       = "total"
	StrInsufficentFund             = "insufficentFund"
	StrTicketPrice                 = "ticketPrice"
	StrUnminedInfo                 = "unminedInfo"
	StrImmatureInfo                = "immatureInfo"
	StrLiveInfo                    = "liveInfo"
	StrLiveInfoDisc                = "liveInfoDisc"
	StrLiveInfoDiscSub             = "liveInfoDiscSub"
	StrVotedInfo                   = "votedInfo"
	StrVotedInfoDisc               = "votedInfoDisc"
	StrRevokeInfo                  = "revokeInfo"
	StrRevokeInfoDisc              = "revokeInfoDisc"
	StrExpiredInfo                 = "expiredInfo"
	StrExpiredInfoDisc             = "expiredInfoDisc"
	StrExpiredInfoDiscSub          = "expiredInfoDiscSub"
	StrLiveIn                      = "liveIn"
	StrSpendableIn                 = "spendableIn"
	StrDuration                    = "duration"
	StrDaysToMiss                  = "daysToMiss"
	StrStakeAge                    = "stakeAge"
	StrSelectOption                = "selectOption"
	StrUpdatevotePref              = "updateVotePref"
	StrVoteUpdated                 = "voteUpdated"
	StrVotingWallet                = "votingWallet"
	StrVoteConfirm                 = "voteConfirm"
	StrVoteSent                    = "voteSent"
	StrNumberOfVotes               = "numberOfVotes"
	StrNotEnoughVotes              = "notEnoughVotes"
	StrSearch                      = "search"
	StrConsensusChange             = "consensusChange"
	StrOnChainVote                 = "onChainVote"
	StrOffChainVote                = "offChainVote"
	StrConsensusDashboard          = "consensusDashboard"
	StrCopyLink                    = "copyLink"
	StrWebURL                      = "webURL"
	StrVotingDashboard             = "votingDashboard"
	StrUpdated                     = "updated"
	StrViewOnPoliteia              = "viewOnPoliteia"
	StrVotingInProgress            = "votingInProgress"
	StrVersion                     = "version"
	StrPublished2                  = "published2"
	StrHowGovernanceWork           = "howGovernanceWork"
	StrGovernanceInfo              = "governanceInfo"
	StrProposalInfo                = "proposalInfo"
	StrSelectTicket                = "selectTicket"
	StrHash                        = "hash"
	StrMax                         = "max"
	StrnoValidWalletFound          = "noValidWalletFound"
	StrSecurityTools               = "securityTools"
	StrAbout                       = "about"
	StrHelp                        = "help"
	StrDarkMode                    = "darkMode"
	StrTxNotification              = "txNotification"
	StrPropNotification            = "propNotification"
	StrHTTPRequest                 = "httpReq"
	StrEnabled                     = "enabled"
	StrDisable                     = "disable"
	StrDisabled                    = "disabled"
	StrGovernanceSettingsInfo      = "governanceSettingsInfo"
	StrPropFetching                = "propFetching"
	StrCheckGovernace              = "checkGovernace"
	StrRemovePeer                  = "removePeer"
	StrRemovePeerWarn              = "removePeerWarn"
	StrRemoveUserAgent             = "removeUserAgent"
	StrRemoveUserAgentWarn         = "removeUserAgentWarn"
	StrIPAddress                   = "ipAddress"
	StrUserAgent                   = "userAgent"
	StrValidateMsg                 = "validateMsg"
	StrValidate                    = "validate"
	StrHelpInfo                    = "helpInfo"
	StrDocumentation               = "documentation"
	StrVerifyMsgNote               = "verifyMsgNote"
	StrVerifyMsgError              = "verifyMsgError"
	StrInvalidSignature            = "invalidSignature"
	StrValidSignature              = "validSignature"
	StrEmptyMsg                    = "emptyMsg"
	StrEmptySign                   = "emptySign"
	StrClear                       = "clear"
	StrValidateAddr                = "validateAddr"
	StrValidateNote                = "validateNote"
	StrNotOwned                    = "notOwned"
	StrOwned                       = "owned"
	StrBuildDate                   = "buildDate"
	StrNetwork                     = "network"
	StrLicense                     = "license"
	StrCheckWalletLog              = "checkWalletLog"
	StrCheckStatistics             = "checkStatistics"
	StrStatistics                  = "statistics"
	StrConfirmDexReset             = "confirmDexReset"
	StrDexResetInfo                = "dexResetInfo"
	StrDexStartupErr               = "dexStartupErr"
	StrResetDexClient              = "resetDexClient"
	StrWalletLog                   = "walletLog"
	StrBuild                       = "build"
	StrPeersConnected              = "peersConnected"
	StrUptime                      = "uptime"
	StrBestBlocks                  = "bestBlocks"
	StrBestBlockAge                = "bestBlockAge"
	StrBestBlockTimestamp          = "bestBlockTimestamp"
	StrWalletDirectory             = "walletDirectory"
	StrDateSize                    = "dateSize"
	StrExit                        = "exit"
	StrLoading                     = "loading"
	StrOpeningWallet               = "openingWallet"
	StrWelcomeNote                 = "welcomeNote"
	StrGenerateAddress             = "generateAddress"
	StrReceivingAddress            = "receivingAddress"
	StrYourAddress                 = "yourAddress"
	StrReceiveInfo                 = "receiveInfo"
	StrDcrReceived                 = "dcrReceived"
	StrTicektVoted                 = "ticektVoted"
	StrTicketRevoked               = "ticketRevoked"
	StrProposalAddedNotif          = "proposalAddedNotif"
	StrVoteStartedNotif            = "voteStartedNotif"
	StrVoteEndedNotif              = "voteEndedNotif"
	StrNewProposalUpdate           = "newProposalUpdate"
	StrWalletSyncing               = "walletSyncing"
	StrNext                        = "next"
	StrRetry                       = "retry"
	StrEstimatedTime               = "estimatedTime"
	StrEstimatedSize               = "estimatedSize"
	StrRate                        = "rate"
	StrTotalCost                   = "totalCost"
	StrBalanceAfter                = "balanceAfter"
	StrSendingAcct                 = "sendingAcct"
	StrTxEstimateErr               = "txEstimateErr"
	StrSendInfo                    = "sendInfo"
	StrAmount                      = "amount"
	StrTxSent                      = "txSent"
	StrConfirmSend                 = "confirmSend"
	StrSendingFrom                 = "sendingFrom"
	StrSendWarning                 = "sendWarning"
	StrDestAddr                    = "destAddr"
	StrMyAcct                      = "myAcct"
	StrSelectWalletToOpen          = "selectWalletToOpen"
	StrSecurityToolsInfo           = "securityToolsInfo"
	StrContinue                    = "continue"
	StrNewWallet                   = "newWallet"
	StrSelectWalletType            = "selectWalletType"
	StrWhatToCallWallet            = "whatToCallWallet"
	StrExistingWalletName          = "existingWalletName"
	StrTicketVotedTitle            = "ticketVotedTitle"
	StrTicketRevokedTitle          = "ticketRevokedTitle"
	StrSyncCompTime                = "syncCompTime"
	StrInfo                        = "info"
)

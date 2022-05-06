package localizable

const ENGLISH = "en"

// one string per line, no multiline
// semicolon is not compulsory
const EN = `
"appName" = "godcr";
"appTitle" = "godcr (%s)";
"recentTransactions" = "Recent Transactions";
"recentProposals" = "Recent Proposals";
"seeAll" = "See all";
"send" = "Send";
"receive" = "Receive";
"online" = "Online";
"offline" = "Offline";
"showDetails" = "Show details";
"hideDetails" = "Hide details";
"peers" = "peers";
"connectedPeersCount" = "Connected peers count";
"noConnectedPeer" = "No connected peers.";
"disconnect" = "Disconnect";
"reconnect" = "Reconnect";
"currentTotalBalance" = "Current Total Balance";
"totalBalance" = "Total Balance";
"walletStatus" = "Wallet Status";
"blockHeaderFetched" = "Block header fetched";
"noTransactions" = "No transactions";
"fetchingBlockHeaders" = "Fetching block headers · %v%%";
"discoveringWalletAddress" = "Discovering wallet address · %v%%";
"rescanningHeaders" = "Rescanning headers · %v%%";
"rescanningBlocks" = "Rescanning blocks";
"blocksScanned" = "Blocks scanned";
"blocksLeft" = "%d blocks left"
"syncSteps" = "Step %d/3";
"blockHeaderFetchedCount" = "%d of %d";
"timeLeft" = "%v left";
"connectedTo" = "Connected to";
"connecting" = "Connecting...";
"synced" = "Synced";
"autoSync" = "Auto sync";
"syncingState" = "Syncing...";
"waitingState" = "Waiting...";
"walletNotSynced" = "Not Synced";
"cancel" = "Cancel";
"resumeAccountDiscoveryTitle" = "Unlock to resume restoration";
"unlock" = "Unlock";
"syncingProgress" = "Syncing progress";
"syncingProgressStat" = "%s behind";
"noWalletLoaded" = "No wallet loaded";
"lastBlockHeight" = "Last Block Height";
"ago" = "ago";
"newest" = "Newest";
"oldest" = "Oldest";
"all" = "All";
"transferred" = "Transferred";
"sent" = "Sent";
"received" = "Received";
"yourself" = "Yourself";
"mixed" = "Mixed";
"unmined" = "Unmined";
"immature" = "Immature";
"voted" = "Voted";
"revoked" = "Revoked";
"live" = "Live";
"expired" = "Expired";
"purchased" = "Purchased";
"revocation" = "Revocation";
"rebroadcast" = "Rebroadcast";
"staking" = "Staking";
"immatureRewards" = "Immature Rewards";
"lockedByTickets" = "Locked By Tickets";
"immatureStakeGen" = "Immature Stake Gen";
"votingAuthority" = "Voting Authority";
"unknown" = "Unknown";
"nConfirmations" = "%d Confirmations";
"from" = "From";
"to" = "To";
"fee" = "Fee";
"includedInBlock" = "Included in block";
"type" = "Type";
"transactionId" = "Transaction ID";
"xInputsConsumed" = "%d Inputs consumed";
"xOutputCreated" = "%d Outputs created";
"viewOnDcrdata" = "View on dcrdata";
"watchOnlyWallets" = "Watch-only wallets";
"signMessage" = "Sign message";
"verifyMessage" = "Verify message";
"message" = "Message";
"viewProperty" = "View property";
"stakeShuffle" = "StakeShuffle";
"rename" = "Rename";
"renameWalletSheetTitle" = "Rename wallet";
"settings" = "Settings";
"createANewWallet" = "Create a new wallet"
"importExistingWallet" = "Import an existing wallet";
"importWatchingOnlyWallet" = "Import a watch-only wallet";
"create" = "Create";
"watchOnlyWalletImported" = "Watch only wallet imported";
"addNewAccount" = "Add new account";
"notBackedUp" = "Not backed up";
"labelSpendable" = "Spendable";
"backupSeedPhrase" = "Back up seed word";
"verifySeedInfo" = "Verify your seed word backup so you can recover your funds when needed.";
"createNewAccount" = "Create new account";
"invalidPassphrase" = "Password entered was not valid.";
"Import" = "Import";
"spendingPasswordInfo" = "A spending password helps secure your wallet transactions."
"spendingPassword" = "Spending password";
"changeSpendingPass" = "Change spending password";
"currentSpendingPassword" = "Current spending password";
"newSpendingPassword" = "New spending password";
"confirmNewSpendingPassword" = "Confirm new spending password";
"spendingPasswordUpdated" = "Spending password updated";
"notifications" = "Notifications";
"beepForNewBlocks" = "Beep for new blocks";
"debug" = "Debug";
"rescanBlockchain" = "Rescan blockchain";
"dangerZone" = "Danger zone";
"removeWallet" = "Remove wallet from device";
"change" = "Change";
"notConnected" = "Not connected to decred network";
"rescanProgressNotification" = "Check progress in overview.";
"rescan" = "Rescan";
"rescanInfo" = "Rescanning may help resolve some balance errors. This will take some time, as it scans the entire blockchain for transactions"
"confirmToRemove" = "Confirm to remove";
"remove" = "Remove";
"confirm" = "Confirm";
"general" = "General";
"unconfirmedFunds" = "Spend Unconfirmed Funds";
"confirmed" = "Confirmed";
"currencyConversion" = "Currency conversion";
"language" = "Language";
"security" = "Security";
"startupPassword" = "Startup password";
"changeStartupPassword" = "Change startup password";
"connection" = "Connection";
"connectToSpecificPeer" = "Connect to specific peer";
"changeSpecificPeer" = "Change specific peer";
"CustomUserAgent" = "Custom user agent";
"userAgentSummary" = "For exchange rate fetching";
"changeUserAgent" = "Change user agent";
"createStartupPassword" = "Create a startup password";
"confirmRemoveStartupPass" = "Confirm to turn off startup password";
"userAgentDialogTitle" = "Set up user agent";
"overview" = "Overview";
"transactions" = "Transactions";
"wallets" = "Wallets";
"tickets" = "Tickets";
"more" = "More";
"english" = "English";
"french" = "French";
"spanish" = "Spanish";
"usdBittrex" = "USD (Bittrex)";
"none" = "None";
"proposals" = "Proposals";
"dex" = "Dex";
"governance" = "Governance";
"pending" = "Pending";
"vote" = "Vote";
"revoke" = "Revoke";
"maturity" = "Maturity";
"yesterday" = "Yesterday";
"days" = "Days";
"hours" = "Hours";
"mins" = "Mins";
"secs" = "Secs";
"weekAgo" = "%d week ago";
"weeksAgo" = "%d weeks ago";
"yearAgo" = "%d year ago";
"yearsAgo" = "%d years ago";  
"monthAgo" = "%d month ago";
"monthsAgo" = "%d months ago";
"dayAgo" = "%d day ago";
"daysAgo" = "%d days ago";
"hourAgo" = "%d hour ago";
"hoursAgo" = " %d hours ago";
"minuteAgo" = "%d minute ago";
"minutesAgo" = "%d minutes ago";
"justNow" = "Just now";
"imawareOfRisk" = "I am aware of the risk";
"unmixedBalance" = "Unmixed balance";
"backupLater" = "Backup later";
"backupNow" = "Backup now";
"status" = "Status";
"daysToVote" = "Days to vote";
"reward" = "Reward";
"viewTicket" = "View associated ticket";
"external" = "External"
"republished" = "Republished unmined transactions to the decred network";
"copied" = "Copied";
"txHashCopied" = "Transaction Hash copied";
"addressCopied" = "Address copied";
"address" = "Address";
"acctDetailsKey" = "%d external, %d internal, %d imported";
"acctNum" = "Account Number";
"acctName" = "Account name";
"acctRenamed" = "Account renamed";
"acctCreated" = "Account created";
"renameAcct" = "Rename account";
"acctCreated" = "Account created"
"hdPath" = "HD Path";
"key" = "Key";
"validateWalSeed" = "Validate wallet seeds";
"clearAll" = "Clear all";
"restoreWallet" = "Restore wallet";
"enterSeedPhrase" = "Enter your seed phrase"
"invalidSeedPhrase" = "Invalid seed phrase"
"seedAlreadyExist" = "A wallet with an identical seed already exists."
"walletRestored" = "Wallet restored"
"enterWalletDetails" = "Enter wallet details"
"copy" = "Copy"
"enterAddressToSign" = "Enter an address and message to sign:"
"signCopied" = "Signature copied"
"signature" = "Signature"
"confirmToSign" = "Confirm to sign"
"enterValidAddress" = "Please enter a valid address"
"enterValidMsg" = "Please enter a valid message to sign"
"invalidAddress" = "Invalid address"
"addrNotOwned" = "Address not owned by any wallet"
"delete" = "Delete"
"walletName" = "Wallet name"
"walletRenamed" = "Wallet renamed"
"walletCreated" = "Wallet created"
"addWallet" = "Add wallet"
"checkMixerStatus" = "Check mixer status"
"walletRestoreMsg" = "You can restore this wallet from seed word after it is deleted."
"walletRemoved" = "Wallet removed"
"walletRemoveInfo" = "Make sure to have the seed word backed up before removing the wallet"
"watchOnlyWalletRemoveInfo" = "The watch-only wallet will be removed from your app"
"gotIt" = "Got it"
"noValidAccountFound" = "no valid account found"
"mixer" = "Mixer"
"readyToMix" = "Ready to mix"
"mixerRunning" = "Mixer is running..."
"keepAppOpen" = "Keep this app opened"
"mixerShutdown" = "The mixer will automatically stop when unmixed balance are fully mixed."
"votingPreference" = "Voting Preference:"
"noAgendaYet" = "No agendas yet"
"fetchingAgenda" = "Fetching agendas..."
"updatePreference" = "Update Preference"
"approved" = "Approved"
"voting" = "Votación"
"rejected" = "Rejected"
"abandoned" = "Abandoned"
"inDiscussion" = "In discussion"
"fetchingProposals" = "Fetching proposals..."
"underReview" = "Under review"
"noProposals" = "No proposals %s"
"waitingForAuthor" = "Waiting for author to authorize voting"
"waitingForAdmin" = "Waiting for admin to trigger the start of voting"
"voteTooltip" = "%d %% Yes votes required for approval"
"yes" = "Yes: "
"no" = "No: "
"totalVotes" = "Total votes:  %6.0f"
"totalVotesReverse" = "%d Total votes"
"quorumRequirement" = "Quorum requirement:  %6.0f"
"discussions" = "Discussions:   %d comments"
"published" = "Published:   %s"
"token" = "Token:   %s"
"proposalVoteDetails" = "Proposal vote details"
"votingServiceProvider" = "Voting service provider
"selectVSP" = "Select VSP..."
"addVSP" = "Add a new VSP..."
"save" = "Save"
"noVSPLoaded" = "No vsp loaded. Check internet connection and try again."
`

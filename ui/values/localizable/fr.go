package localizable

const FRENCH = "fr"

const FR = `
"watchOnlyWallets" = "Porte-monnaies en lecture seule";
"changeStartupPassword" = "Changer le mot de passe de démarrage";
"transactions" = "Transactions";
"unconfirmedFunds" = "Dépenser des fonds non-confirmés";
"confirmed" = "Confirmé";
"notBackedUp" = "Non sauvegardé";
"Import" = "Importer";
"overview" = "Vue d\'ensemble";
"from" = "De";
"addNewAccount" = "Ajouter un nouveau compte";
"changeSpendingPass" = "Changer le mot de passe de dépense";
"rescanBlockchain" = "Rescanner blockchain";
"xInputsConsumed" = "%d entrées consommées";
"userAgentDialogTitle" = "Configurer user agent";
"verifySeedInfo" = "Vérifier la sauvegarde de votre word de récupération afin de pouvoir récupérer vos fonds.";
"createNewAccount" = "Créer un nouveau compte";
"createANewWallet" = "Créer un nouveau porte-monnaie";
"seeAll" = "Voir tout";
"resumeAccountDiscoveryTitle" = "Déverouiller pour reprendre la restauration";
"offline" = "Hors ligne";
"showDetails" = "Montrer les détails";
"walletStatus" = "Statut de porte-monnaie";
"wallets" = "Porte-monnaies";
"settings" = "Paramètres";
"labelSpendable" = "Disponible";
"remove" = "Supprimer";
"synced" = "Synchronisé";
"more" = "Plus";
"sent" = "Envoyé";
"connection" = "Connexion";
"connecting" = "De liaison...";
"cancel" = "Annuler";
"backupSeedPhrase" = "Sauvegarder la word de récupération";
"fetchingBlockHeaders" = "Récupération des entêtes de blocs · %v%%";
"discoveringWalletAddress" = "Découvrir l'adresse du portefeuille · %v%%";
"rescanningHeaders" = "Réanalyser les en-têtes · %v%%";
"rescanningBlocks" = "Renumériser les blocs";
"rescanInfo" = "Une nouvelle analyse peut aider à résoudre certaines erreurs d'équilibrage. Cela prendra un certain temps, car il analyse l'intégralité de la blockchain pour les transactions"
"rescan" = "Renumériser"
"blocksScanned" = "Blocs scannés";
"blocksLeft" = "%d blocs restants";
"autoSync" = "Synchronisation automatique";
"blockHeaderFetchedCount" = "%d de %d";
"timeLeft" = "%v restant";
"currencyConversion" = "Conversion de devise";
"renameWalletSheetTitle" = "Renommer le porte-monnaie";
"send" = "Envoyer";
"transferred" = "Transféré";
"type" = "Type";
"to" = "À";
"importWatchingOnlyWallet" = "Importer un porte-monnaie en observation seule";
"create" = "Créer";
"hideDetails" = "Masquer les détails";
"rescanProgressNotification" = "Suivre l\'avancement dans la vue d\'ensemble!";
"general" = "Général";
"reconnect" = "Reconnecter";
"importExistingWallet" = "Importer un porte-monnaie existant";
"syncingProgress" = "Synchronisation en cours";
"syncingProgressStat" = "%s derrière";
"confirm" = "Confirmer";
"startupPassword" = "Mot de passe de démarrage";
"peers" = "pairs";
"connectedPeersCount" = "Compteur de pairs connectés";
"noConnectedPeer" = "Aucun pair connecté.";
"lastBlockHeight" = "Hauteur du dernier bloc";
"blockHeaderFetched" = "Entête de bloc récupéré";
"includedInBlock" = "Inclus dans le bloc";
"online" = "En ligne";
"unlock" = "Déverouiller";
"xOutputCreated" = "%d sorties créées";
"invalidPassphrase" = "La word secrète entrée n\'est pas valide.";
"passwordNotMatch" = "Les mots de passe ne correspondent pas"
"confirmToRemove" = "Confirmer la suppression";
"syncingState" = "Synchronisation...";
"waitingState" = "Attendre...";
"verifyMessage" = "Vérifier le message";
"message" = "Message";
"rebroadcast" = "Retransmettre";
"spendingPasswordInfo" = "Un mot de passe de dépenses permet de sécuriser les transactions de votre portefeuille."
"spendingPasswordInfo2" = "Ce mot de passe de dépenses est pour le nouveau portefeuille uniquement"
"spendingPassword" = "Mot de passe de dépense";
"enterSpendingPassword" = "Ingrese la contraseña de gastos"
"confirmSpendingPassword" = "Confirmer le mot de passe de dépenses";
"currentSpendingPassword" = "Mot de passe de dépenses actuel";
"newSpendingPassword" = "Nouveau mot de passe de dépenses";
"confirmNewSpendingPassword" = "Confirmer le nouveau mot de passe de dépenses";
"spendingPasswordUpdated" = "Spending password updated";
"change" = "Échanger";
"currentTotalBalance" = "Solde Total Actuel";
"totalBalance" = "Solde Total";
"signMessage" = "Signer le message";
"notifications" = "Notifications";
"disconnect" = "Déconnecter";
"receive" = "Recevoir";
"fee" = "Frais";
"userAgentSummary" = "Pour rafraîchir les taux de change";
"rename" = "Renommer";
"noTransactions" = "Aucune transaction pour le moment.";
"createStartupPassword" = "Créer un mot de passe de démarrage";
"transactionId" = "Transaction ID";
"debug" = "Debug";
"notConnected" = "Pas de connexion au réseau Decred";
"received" = "Reçu";
"mixed" = "Mezclado";
"unmined" = "Non miné";
"immature" = "Immature";
"voted" = "Voté";
"revoked" = "Révoqué";
"live" = "Vivre";
"expired" = "Expiré";
"maturity" = "Maturité";
"purchased" = "Acheté";
"revocation" ="Révocation";
"yourself" = "Toi même";
"immatureRewards" = "Récompenses immatures";
"lockedByTickets" = "Verrouillé par des billets;
"immatureStakeGen" = "Génération d'enjeux immatures";
"votingAuthority" = "Autorité de vote";
"unknown" = "Inconnue";
"removeWallet" = "Supprimer le porte-monnaie de cet appareil";
"recentTransactions" = "Transactions récentes";
"recentProposals" = "Propositions récentes";
"ago" = "avant";
"viewOnDcrdata" = "Voir sur dcrdata";
"beepForNewBlocks" = "Beep pour les nouveaux blocs";
"connectToSpecificPeer" = "Se connecter à un pair spécifique";
"staking" = "Jalonnement";
"english" = "Anglais";
"french" = "Français";
"spanish" = "Espanol";
"governance" = "Gouvernance";
"pending" = "En attendant";
"vote" = "Vote";
"revoke" = "Révoquer";
"yesterday" = "hier";
"days" = "jours";
"hours" = "Les heures";
"mins" = "minutes";
"secs" = "seconds";
"weekAgo" = "%d la semaine ago";
"weeksAgo" = "%d semaines ago";
"yearAgo" = "%d an ago";
"yearsAgo" = "%d ans ago";  
"monthAgo" = "%d mois ago";
"monthsAgo" = "%d mois ago";
"dayAgo" = "%d journée ago";
"daysAgo" = "%d jours ago";
"hourAgo" = "%d heure ago";
"hoursAgo" = "%d les heures ago";
"minuteAgo" = "%d minute ago";
"minutesAgo" = "%d minutes ago";
"justNow" = "Juste maintenant";
"imawareOfRisk" = "Je suis conscient du risque";
"unmixedBalance" = "Solde non mélangé";
"backupLater" = "Sauvegardez plus tard";
"backupNow" = "Sauvegarder maintenant";
"status" = "Statut";
"daysToVote" = "Jours pour voter";
"reward" = "Récompense";
"viewTicket" = "Voir le ticket associé";
"external" = "Externe";
"republished" = "Transactions non minées republiées sur le réseau décrété";
"copied" = "Copié";
"txHashCopied" = "Hachage de la transaction copié";
"addressCopied" = "Adresse copiée";
"address" = "Adresse";
"acctDetailsKey" = "%d externe, %d interne, %d importé";
"key" = "Clé"
"acctNum" = "Numéro de compte";
"acctName" = "Nom du compte";
"acctRenamed" = "Compte renommé";
"renameAcct" = "Renommer le compte";
"acctCreated" = "Compte créé"
"hdPath" = "Chemin HD";
"validateWalSeed" = "Valider les graines de portefeuille";
"clearAll" = "Tout effacer";
"restoreWallet" = "Restaurer le portefeuille";
"enterSeedPhrase" = "Entrez votre phrase de départ"
"invalidSeedPhrase" = "Phrase de départ non valide"
"walletExist" = "Le portefeuille avec le nom: %s existe déjà"
"walletNotExist" = "Portefeuille avec ID: %v n'existe pas"
"seedAlreadyExist" = "Un portefeuille avec une graine identique existe déjà."
"walletRestored" = "Portefeuille restauré"
"enterWalletDetails" = "Entrez les détails du portefeuille"
"copy" = "Copie"
"howToCopy" = "Comment copier"
"enterAddressToSign" = "Entrez une adresse et un message à signer:"
"signCopied" = "Signature copiée"
"signature" = "Signature"
"confirmToSign" = "Confirmer pour signer"
"enterValidAddress" = "s'il-vous-plaît entrez une adresse valide"
"enterValidMsg" = "Veuillez entrer un message valide pour signer"
"invalidAddress" = "Adresse invalide"
"addrNotOwned" = "Adresse n'appartenant à aucun portefeuille"
"delete" = "Supprimez"
"walletName" = "Nom du portefeuille"
"enterWalletName" = "Entrez le nom du portefeuille"
"walletRenamed" = "Portefeuille renommé"
"walletCreated" = "Portefeuille créé"
"addWallet" = "Ajouter un portefeuille"
"checkMixerStatus" = "Vérifier l'état du mélangeur"
"walletRestoreMsg" = "Vous pouvez restaurer ce portefeuille à partir du mot de départ après sa suppression."
"walletRemoved" = "Portefeuille supprimé"
"walletRemoveInfo" = "Assurez-vous d'avoir sauvegardé le mot de départ avant de retirer le portefeuille"
"watchOnlyWalletRemoveInfo" = "Le portefeuille réservé aux montres sera supprimé de votre application"
"gotIt" = "J'ai compris"
"noValidAccountFound" = "aucun compte valide trouvé"
"mixer" = "Mixer"
"readyToMix" = "Prêt à mélanger"
"mixerRunning" = "Le mélangeur est en marche..."
"useMixer" = "Comment utiliser le mixeur?"
"keepAppOpen" = "Gardez cette application ouverte"
"mixerShutdown" = "Le mélangeur s'arrêtera automatiquement lorsque la balance non mélangée sera entièrement mélangée."
"votingPreference" = "Préférence de vote:"
"noAgendaYet" = "Pas encore d'agendas"
"fetchingAgenda" = "Récupération des agendas"
"updatePreference" = "Préférence de mise à jour"
"approved" = "A approuvé"
"voting" = "Vote"
"rejected" = "Rejeté"
"abandoned" = "Abandonné"
"inDiscussion" = "En discussion"
"fetchingProposals" = "Récupération des propositions..."
"underReview" = "À l'étude"
"noProposals" = "Aucune proposition %s"
"waitingForAuthor" = "En attente de l'autorisation de vote de l'auteur"
"waitingForAdmin" = "En attente que l'admin déclenche le début du vote"
"voteTooltip" = "%d %% Oui votes requis pour approbation"
"yes" = "Oui: "
"no" = "Non: "
"totalVotes" = "Total des votes:  %6.0f"
"totalVotesReverse" = "%d Total des votes"
"quorumRequirement" = "Dont l'exigence:  %6.0f"
"discussions" = "Discussions:   %d commentaires"
"published" = "Publié:   %s"
"token" = "Jeton:   %s"
"proposalVoteDetails" = "Détails du vote de la proposition"
"votingServiceProvider" = "Fournisseur de services de vote"
"selectVSP" = "Sélectionnez PSV..."
"addVSP" = "Ajouter une nouvelle VSP..."
"save" = "Sauver"
"noVSPLoaded" = "Aucun vsp chargé. Vérifiez la connexion Internet et réessayez."
"extendedPubKey" = "Clé publique étendue"
"enterXpubKey" = "entrez une clé de pub étendue valide"
"xpubKeyErr" = "Erreur lors de la vérification de xpub: %v"
"xpubWalletExist" = "Un portefeuille avec une clé publique étendue identique existe déjà."
"hint" = "indice"
"addAcctWarn" = "Les comptes ne peuvent pas être supprimés une fois créés."
"balToMaintain" = "Équilibre à maintenir (DCR)"
"autoTicketPurchase" = "Achat automatique de billets"
"purchasingAcct" = "Compte d'achat"


`

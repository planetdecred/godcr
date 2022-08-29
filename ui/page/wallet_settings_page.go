package page

import (
	"gioui.org/layout"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/page/security"
	"github.com/planetdecred/godcr/ui/values"
)

const WalletSettingsPageID = "WalletSettings"

type clickableRowData struct {
	clickable *decredmaterial.Clickable
	labelText string
	title     string
}

type accountData struct {
	*dcrlibwallet.Account
	clickable *decredmaterial.Clickable
}

type WalletSettingsPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	wallet   *dcrlibwallet.Wallet
	accounts []*accountData

	pageContainer layout.List
	accountsList  *decredmaterial.ClickableList

	changePass, rescan, resetDexData           *decredmaterial.Clickable
	changeAccount, checklog, checkStats        *decredmaterial.Clickable
	changeWalletName, addAccount, deleteWallet *decredmaterial.Clickable
	verifyMessage, validateAddr, signMessage   *decredmaterial.Clickable
	updateConnectToPeer                        *decredmaterial.Clickable

	chevronRightIcon *decredmaterial.Icon
	backButton       decredmaterial.IconButton
	infoButton       decredmaterial.IconButton

	fetchProposal     *decredmaterial.Switch
	proposalNotif     *decredmaterial.Switch
	spendUnconfirmed  *decredmaterial.Switch
	spendUnmixedFunds *decredmaterial.Switch
	connectToPeer     *decredmaterial.Switch
	peerAddr          string
}

func NewWalletSettingsPage(l *load.Load) *WalletSettingsPage {
	pg := &WalletSettingsPage{
		Load:                l,
		GenericPageModal:    app.NewGenericPageModal(WalletSettingsPageID),
		wallet:              l.WL.SelectedWallet.Wallet,
		changePass:          l.Theme.NewClickable(false),
		rescan:              l.Theme.NewClickable(false),
		resetDexData:        l.Theme.NewClickable(false),
		changeAccount:       l.Theme.NewClickable(false),
		checklog:            l.Theme.NewClickable(false),
		checkStats:          l.Theme.NewClickable(false),
		changeWalletName:    l.Theme.NewClickable(false),
		addAccount:          l.Theme.NewClickable(false),
		deleteWallet:        l.Theme.NewClickable(false),
		verifyMessage:       l.Theme.NewClickable(false),
		validateAddr:        l.Theme.NewClickable(false),
		signMessage:         l.Theme.NewClickable(false),
		updateConnectToPeer: l.Theme.NewClickable(false),

		fetchProposal:     l.Theme.Switch(),
		proposalNotif:     l.Theme.Switch(),
		spendUnconfirmed:  l.Theme.Switch(),
		spendUnmixedFunds: l.Theme.Switch(),
		connectToPeer:     l.Theme.Switch(),

		pageContainer: layout.List{Axis: layout.Vertical},
		accountsList:  l.Theme.NewClickableList(layout.Vertical),
	}

	pg.chevronRightIcon = decredmaterial.NewIcon(l.Theme.Icons.ChevronRight)
	pg.chevronRightIcon.Color = pg.Theme.Color.Gray1

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *WalletSettingsPage) OnNavigatedTo() {
	// set switch button state on page load
	pg.fetchProposal.SetChecked(pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey, false))
	pg.proposalNotif.SetChecked(pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.ProposalNotificationConfigKey, false))
	pg.spendUnconfirmed.SetChecked(pg.WL.SelectedWallet.Wallet.ReadBoolConfigValueForKey(dcrlibwallet.SpendUnconfirmedConfigKey, false))
	pg.spendUnmixedFunds.SetChecked(pg.WL.SelectedWallet.Wallet.ReadBoolConfigValueForKey(load.SpendUnmixedFundsKey, false))

	pg.peerAddr = pg.WL.MultiWallet.ReadStringConfigValueForKey(dcrlibwallet.SpvPersistentPeerAddressesConfigKey)
	pg.connectToPeer.SetChecked(false)
	if pg.peerAddr != "" {
		pg.connectToPeer.SetChecked(true)
	}

	pg.loadWalletAccount()
}

func (pg *WalletSettingsPage) loadWalletAccount() {
	walletAccounts := make([]*accountData, 0)
	accounts, err := pg.wallet.GetAccountsRaw()
	if err != nil {
		log.Errorf("error retrieving wallet accounts: %v", err)
		return
	}

	for _, acct := range accounts.Acc {
		if acct.Number == dcrlibwallet.ImportedAccountNumber {
			continue
		}

		walletAccounts = append(walletAccounts, &accountData{
			Account:   acct,
			clickable: pg.Theme.NewClickable(false),
		})
	}

	pg.accounts = walletAccounts
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *WalletSettingsPage) Layout(gtx C) D {
	body := func(gtx C) D {
		w := []func(gtx C) D{
			func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding26,
				}.Layout(gtx, pg.Theme.Label(values.TextSize20, values.String(values.StrSettings)).Layout)
			},
			pg.generalSection(),
			pg.account(),
			pg.securityTools(),
			pg.debug(),
			pg.dangerZone(),
		}

		return pg.pageContainer.Layout(gtx, len(w), func(gtx C, i int) D {
			return layout.Inset{Left: values.MarginPadding50}.Layout(gtx, w[i])
		})
	}

	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx, body)
	}
	return pg.layoutDesktop(gtx, body)
}

func (pg *WalletSettingsPage) layoutDesktop(gtx C, body layout.Widget) D {
	return components.UniformPadding(gtx, body)
}

func (pg *WalletSettingsPage) layoutMobile(gtx C, body layout.Widget) D {
	return components.UniformMobile(gtx, false, false, body)
}

func (pg *WalletSettingsPage) generalSection() layout.Widget {
	dim := func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.sectionContent(pg.changePass, values.String(values.StrSpendingPassword))),
			layout.Rigid(pg.sectionContent(pg.changeWalletName, values.String(values.StrRenameWalletSheetTitle))),
			layout.Rigid(pg.subSectionSwitch(values.String(values.StrFetchProposals), pg.fetchProposal)),
			layout.Rigid(func(gtx C) D {
				if !pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.FetchProposalConfigKey, false) {
					return D{}
				}
				return pg.subSection(gtx, values.String(values.StrPropNotif), pg.proposalNotif.Layout)
			}),
			layout.Rigid(pg.subSectionSwitch(values.String(values.StrUnconfirmedFunds), pg.spendUnconfirmed)),
			layout.Rigid(pg.subSectionSwitch(values.String(values.StrAllowSpendingFromUnmixedAccount), pg.spendUnmixedFunds)),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(pg.subSectionSwitch(values.String(values.StrConnectToSpecificPeer), pg.connectToPeer)),
					layout.Rigid(func(gtx C) D {
						if pg.WL.MultiWallet.ReadStringConfigValueForKey(dcrlibwallet.SpvPersistentPeerAddressesConfigKey) == "" {
							return D{}
						}

						peerAddrRow := clickableRowData{
							title:     values.String(values.StrPeer),
							clickable: pg.updateConnectToPeer,
							labelText: pg.peerAddr,
						}
						return pg.clickableRow(gtx, peerAddrRow)
					}),
				)
			}),
		)
	}

	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrGeneral), dim)
	}
}

func (pg *WalletSettingsPage) account() layout.Widget {
	dim := func(gtx C) D {
		return pg.accountsList.Layout(gtx, len(pg.accounts), func(gtx C, a int) D {
			return pg.subSection(gtx, pg.accounts[a].Name, func(gtx C) D {
				return pg.chevronRightIcon.Layout(gtx, values.MarginPadding20)
			})
		})
	}
	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrAccount), dim)
	}
}

func (pg *WalletSettingsPage) debug() layout.Widget {
	dims := func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.sectionContent(pg.rescan, values.String(values.StrRescanBlockchain))),
			layout.Rigid(pg.sectionContent(pg.checklog, values.String(values.StrCheckWalletLog))),
			layout.Rigid(pg.sectionContent(pg.checkStats, values.String(values.StrCheckStatistics))),
			layout.Rigid(pg.sectionContent(pg.resetDexData, values.String(values.StrResetDexClient))),
		)
	}

	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrDebug), dims)
	}
}

func (pg *WalletSettingsPage) securityTools() layout.Widget {
	dims := func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.sectionContent(pg.verifyMessage, values.String(values.StrVerifyMessage))),
			layout.Rigid(pg.sectionContent(pg.validateAddr, values.String(values.StrValidateMsg))),
			layout.Rigid(pg.sectionContent(pg.signMessage, values.String(values.StrSignMessage))),
		)
	}

	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrSecurityTools), dims)
	}
}

func (pg *WalletSettingsPage) dangerZone() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrDangerZone),
			pg.sectionContent(pg.deleteWallet, values.String(values.StrRemoveWallet)),
		)
	}
}

func (pg *WalletSettingsPage) pageSections(gtx C, title string, body layout.Widget) D {
	dims := func(gtx C, title string, body layout.Widget) D {
		return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							txt := pg.Theme.Label(values.TextSize14, title)
							txt.Color = pg.Theme.Color.GrayText2
							return txt.Layout(gtx)
						}),
						layout.Flexed(1, func(gtx C) D {
							if title == values.String(values.StrSecurityTools) {
								pg.infoButton.Inset = layout.UniformInset(values.MarginPadding0)
								pg.infoButton.Size = values.MarginPadding16
								return layout.E.Layout(gtx, pg.infoButton.Layout)
							}
							if title == values.String(values.StrAccount) {
								return layout.E.Layout(gtx, func(gtx C) D {
									mGtx := gtx
									if pg.WL.SelectedWallet.Wallet.IsWatchingOnlyWallet() {
										mGtx = gtx.Disabled()
									}
									return pg.addAccount.Layout(mGtx, pg.Theme.Icons.AddIcon.Layout24dp)
								})
							}

							return D{}
						}),
					)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Bottom: values.MarginPadding10,
						Top:    values.MarginPadding7,
					}.Layout(gtx, pg.Theme.Separator().Layout)
				}),
				layout.Rigid(body),
			)
		})
	}

	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return dims(gtx, title, body)
	})
}

func (pg *WalletSettingsPage) sectionContent(clickable *decredmaterial.Clickable, title string) layout.Widget {
	return func(gtx C) D {
		return clickable.Layout(gtx, func(gtx C) D {
			textLabel := pg.Theme.Label(values.TextSize16, title)
			if title == values.String(values.StrRemoveWallet) {
				textLabel.Color = pg.Theme.Color.Danger
			}
			return layout.Inset{
				Bottom: values.MarginPadding20,
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(textLabel.Layout),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return pg.chevronRightIcon.Layout(gtx, values.MarginPadding20)
						})
					}),
				)
			})
		})
	}
}

func (pg *WalletSettingsPage) subSection(gtx C, title string, body layout.Widget) D {
	return layout.Inset{Top: values.MarginPadding5, Bottom: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(pg.Theme.Label(values.TextSize16, title).Layout),
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, body)
			}),
		)
	})
}

func (pg *WalletSettingsPage) subSectionSwitch(title string, option *decredmaterial.Switch) layout.Widget {
	return func(gtx C) D {
		return pg.subSection(gtx, title, option.Layout)
	}
}

func (pg *WalletSettingsPage) changeSpendingPasswordModal() {
	currentSpendingPasswordModal := modal.NewPasswordModal(pg.Load).
		Title(values.String(values.StrChangeSpendingPass)).
		Hint(values.String(values.StrCurrentSpendingPassword)).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton(values.String(values.StrConfirm), func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := pg.wallet.UnlockWallet([]byte(password))
				if err != nil {
					pm.SetError(err.Error())
					pm.SetLoading(false)
					return
				}
				pg.wallet.LockWallet()
				pm.Dismiss()

				// change password
				newSpendingPasswordModal := modal.NewCreatePasswordModal(pg.Load).
					Title(values.String(values.StrChangeSpendingPass)).
					EnableName(false).
					PasswordHint(values.String(values.StrNewSpendingPassword)).
					ConfirmPasswordHint(values.String(values.StrConfirmNewSpendingPassword)).
					PasswordCreated(func(walletName, newPassword string, m *modal.CreatePasswordModal) bool {
						go func() {
							err := pg.WL.MultiWallet.ChangePrivatePassphraseForWallet(pg.wallet.ID, []byte(password),
								[]byte(newPassword), dcrlibwallet.PassphraseTypePass)
							if err != nil {
								m.SetError(err.Error())
								m.SetLoading(false)
								return
							}
							pg.Toast.Notify(values.String(values.StrSpendingPasswordUpdated))
							m.Dismiss()
						}()
						return false
					})
				pg.ParentWindow().ShowModal(newSpendingPasswordModal)

			}()

			return false
		})
	pg.ParentWindow().ShowModal(currentSpendingPasswordModal)
}

func (pg *WalletSettingsPage) deleteWalletModal() {
	textModal := modal.NewTextInputModal(pg.Load)
	textModal.Hint(values.String(values.StrWalletName)).
		SetTextWithTemplate(modal.RemoveWalletInfoTemplate).
		PositiveButtonStyle(pg.Load.Theme.Color.Surface, pg.Load.Theme.Color.Danger).
		PositiveButton(values.String(values.StrRemove), func(walletName string, tim *modal.TextInputModal) bool {
			if walletName != pg.WL.SelectedWallet.Wallet.Name {
				pg.Toast.NotifyError("Wallet name entered does not match selected wallet.")
				textModal.SetLoading(false)
				return false
			}

			walletDeleted := func() {
				if pg.WL.MultiWallet.LoadedWalletsCount() > 0 {
					pg.Toast.Notify(values.String(values.StrWalletRemoved))
					textModal.Dismiss()
					pg.ParentNavigator().CloseCurrentPage()
					onWalSelected := func() {
						pg.ParentWindow().CloseCurrentPage()
					}
					onDexServerSelected := func(server string) {
						log.Info("Not implemented yet...", server)
					}
					pg.ParentWindow().Display(NewWalletDexServerSelector(pg.Load, onWalSelected, onDexServerSelected))
				} else {
					textModal.Dismiss()
					pg.ParentWindow().CloseAllPages()
				}
			}

			if pg.wallet.IsWatchingOnlyWallet() {
				textModal.SetLoading(true)
				go func() {
					// no password is required for watching only wallets.
					err := pg.WL.MultiWallet.DeleteWallet(pg.wallet.ID, nil)
					if err != nil {
						pg.Toast.NotifyError(err.Error())
						textModal.SetLoading(false)
					} else {
						walletDeleted()
					}
				}()
				return false
			}

			walletPasswordModal := modal.NewPasswordModal(pg.Load).
				Title(values.String(values.StrConfirmToRemove)).
				NegativeButton(values.String(values.StrCancel), func() {
					textModal.SetLoading(false)
				}).
				PositiveButton(values.String(values.StrConfirm), func(password string, pm *modal.PasswordModal) bool {
					go func() {
						err := pg.WL.MultiWallet.DeleteWallet(pg.wallet.ID, []byte(password))
						if err != nil {
							pm.SetError(err.Error())
							pm.SetLoading(false)
							return
						}

						walletDeleted()
						pm.Dismiss() // calls RefreshWindow.
					}()
					return false
				})
			pg.ParentWindow().ShowModal(walletPasswordModal)
			return false

		})

	textModal.Title(values.String(values.StrRemoveWallet)).
		NegativeButton(values.String(values.StrCancel), func() {})
	pg.ParentWindow().ShowModal(textModal)
}

func (pg *WalletSettingsPage) renameWalletModal() {
	textModal := modal.NewTextInputModal(pg.Load).
		Hint(values.String(values.StrWalletName)).
		PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
		PositiveButton(values.String(values.StrRename), func(newName string, tim *modal.TextInputModal) bool {
			err := pg.WL.MultiWallet.RenameWallet(pg.wallet.ID, newName)
			if err != nil {
				pg.Toast.NotifyError(err.Error())
				return false
			}
			return true
		})

	textModal.Title(values.String(values.StrRenameWalletSheetTitle)).
		NegativeButton(values.String(values.StrCancel), func() {})
	pg.ParentWindow().ShowModal(textModal)
}

func (pg *WalletSettingsPage) showSPVPeerDialog() {
	textModal := modal.NewTextInputModal(pg.Load).
		Hint(values.String(values.StrIPAddress)).
		PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
		PositiveButton(values.String(values.StrConfirm), func(ipAddress string, tim *modal.TextInputModal) bool {
			if ipAddress != "" {
				pg.WL.MultiWallet.SaveUserConfigValue(dcrlibwallet.SpvPersistentPeerAddressesConfigKey, ipAddress)
			}
			return true
		})

	textModal.Title(values.String(values.StrConnectToSpecificPeer)).
		NegativeButton(values.String(values.StrCancel), func() {})
	pg.ParentWindow().ShowModal(textModal)
}

func (pg *WalletSettingsPage) clickableRow(gtx C, row clickableRowData) D {
	return row.clickable.Layout(gtx, func(gtx C) D {
		return pg.subSection(gtx, row.title, func(gtx C) D {
			lbl := pg.Theme.Label(values.TextSize16, row.labelText)
			lbl.Color = pg.Theme.Color.GrayText2
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(lbl.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding2}.Layout(gtx, func(gtx C) D {
						return pg.chevronRightIcon.Layout(gtx, values.MarginPadding20)
					})
				}),
			)
		})
	})
}

func (pg *WalletSettingsPage) showWarningModalDialog(title, msg, key string) {
	info := modal.NewInfoModal(pg.Load).
		Title(title).
		Body(msg).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButtonStyle(pg.Theme.Color.Surface, pg.Theme.Color.Danger).
		PositiveButton(values.String(values.StrRemove), func(isChecked bool) bool {
			pg.WL.MultiWallet.DeleteUserConfigValueForKey(key)
			return true
		})
	pg.ParentWindow().ShowModal(info)
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *WalletSettingsPage) HandleUserInteractions() {
	for pg.changePass.Clicked() {
		pg.changeSpendingPasswordModal()
		break
	}

	for pg.rescan.Clicked() {
		go func() {
			info := modal.NewInfoModal(pg.Load).
				Title(values.String(values.StrRescanBlockchain)).
				Body(values.String(values.StrRescanInfo)).
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButtonStyle(pg.Theme.Color.Primary, pg.Theme.Color.Surface).
				PositiveButton(values.String(values.StrRescan), func(isChecked bool) bool {
					err := pg.WL.MultiWallet.RescanBlocks(pg.wallet.ID)
					if err != nil {
						if err.Error() == dcrlibwallet.ErrNotConnected {
							pg.Toast.NotifyError(values.String(values.StrNotConnected))
							return true
						}
						pg.Toast.NotifyError(err.Error())
						return true
					}
					msg := values.String(values.StrRescanProgressNotification)
					pg.Toast.Notify(msg)
					return true
				})

			pg.ParentWindow().ShowModal(info)
		}()
		break
	}

	for pg.deleteWallet.Clicked() {
		pg.deleteWalletModal()
		break
	}

	for pg.changeWalletName.Clicked() {
		pg.renameWalletModal()
		break
	}

	if pg.infoButton.Button.Clicked() {
		info := modal.NewInfoModal(pg.Load).
			PositiveButtonStyle(pg.Theme.Color.Primary, pg.Theme.Color.Surface).
			SetContentAlignment(layout.W, layout.Center).
			SetupWithTemplate(modal.SecurityToolsInfoTemplate).
			Title(values.String(values.StrSecurityTools)).
			PositiveButton(values.String(values.StrOK), func(isChecked bool) bool {
				return true
			})
		pg.ParentWindow().ShowModal(info)
	}

	if pg.fetchProposal.Changed() {
		if pg.fetchProposal.IsChecked() {
			go pg.WL.MultiWallet.Politeia.Sync()
			// set proposal notification config when proposal fetching is enabled
			pg.proposalNotif.SetChecked(pg.WL.MultiWallet.ReadBoolConfigValueForKey(load.ProposalNotificationConfigKey, false))
			pg.WL.MultiWallet.SaveUserConfigValue(load.FetchProposalConfigKey, true)
		} else {
			info := modal.NewInfoModal(pg.Load).
				Title(values.String(values.StrGovernance)).
				Body(values.String(values.StrGovernanceSettingsInfo)).
				NegativeButton(values.String(values.StrCancel), func() {
					pg.fetchProposal.SetChecked(true)
				}).
				PositiveButtonStyle(pg.Theme.Color.Surface, pg.Theme.Color.Danger).
				PositiveButton(values.String(values.StrDisable), func(isChecked bool) bool {
					if pg.WL.MultiWallet.Politeia.IsSyncing() {
						go pg.WL.MultiWallet.Politeia.StopSync()
					}

					pg.WL.MultiWallet.SaveUserConfigValue(load.FetchProposalConfigKey, false)
					pg.WL.MultiWallet.Politeia.ClearSavedProposals()
					// set proposal notification config when proposal fetching is disabled
					pg.WL.MultiWallet.SaveUserConfigValue(load.ProposalNotificationConfigKey, false)
					return true
				})
			pg.ParentWindow().ShowModal(info)
		}
	}

	if pg.spendUnconfirmed.Changed() {
		pg.WL.SelectedWallet.Wallet.SaveUserConfigValue(dcrlibwallet.SpendUnconfirmedConfigKey, pg.spendUnconfirmed.IsChecked())
	}

	if pg.spendUnconfirmed.Changed() {
		if pg.spendUnconfirmed.IsChecked() {
			textModal := modal.NewTextInputModal(pg.Load).
				SetTextWithTemplate(modal.AllowUnmixedSpendingTemplate).
				Hint("").
				PositiveButtonStyle(pg.Load.Theme.Color.Danger, pg.Load.Theme.Color.InvText).
				PositiveButton(values.String(values.StrConfirm), func(textInput string, tim *modal.TextInputModal) bool {
					if textInput != values.String(values.StrAwareOfRisk) {
						tim.SetError("confirmation text is incorrect")
						tim.SetLoading(false)
					} else {
						pg.WL.SelectedWallet.Wallet.SetBoolConfigValueForKey(load.SpendUnmixedFundsKey, true)
						tim.Dismiss()
					}
					return false
				})

			textModal.Title(values.String(values.StrConfirmUmixedSpending)).
				NegativeButton(values.String(values.StrCancel), func() {
					pg.spendUnconfirmed.SetChecked(false)
				})
			pg.ParentWindow().ShowModal(textModal)

		} else {
			pg.WL.SelectedWallet.Wallet.SetBoolConfigValueForKey(load.SpendUnmixedFundsKey, false)
		}
	}

	specificPeerKey := dcrlibwallet.SpvPersistentPeerAddressesConfigKey
	if pg.connectToPeer.Changed() {
		if pg.connectToPeer.IsChecked() {
			pg.showSPVPeerDialog()
			return
		}

		title := values.String(values.StrRemovePeer)
		msg := values.String(values.StrRemovePeerWarn)
		pg.showWarningModalDialog(title, msg, specificPeerKey)
	}

	for pg.updateConnectToPeer.Clicked() {
		pg.showSPVPeerDialog()
		break
	}

	if pg.verifyMessage.Clicked() {
		pg.ParentNavigator().Display(security.NewVerifyMessagePage(pg.Load))
	}

	if pg.validateAddr.Clicked() {
		pg.ParentNavigator().Display(security.NewValidateAddressPage(pg.Load))
	}

	if pg.signMessage.Clicked() {
		pg.ParentNavigator().Display(security.NewSignMessagePage(pg.Load))
	}

	if pg.checklog.Clicked() {
		pg.ParentNavigator().Display(NewLogPage(pg.Load))
	}

	if pg.checkStats.Clicked() {
		pg.ParentNavigator().Display(NewStatPage(pg.Load))
	}

	if pg.proposalNotif.Changed() {
		pg.WL.MultiWallet.SaveUserConfigValue(load.ProposalNotificationConfigKey, pg.proposalNotif.IsChecked())
	}

	if pg.resetDexData.Clicked() {
		pg.resetDexDataModal()
	}

	for pg.addAccount.Clicked() {
		newPasswordModal := modal.NewCreatePasswordModal(pg.Load).
			Title(values.String(values.StrCreateNewAccount)).
			EnableName(true).
			NameHint(values.String(values.StrAcctName)).
			EnableConfirmPassword(false).
			PasswordHint(values.String(values.StrSpendingPassword)).
			PasswordCreated(func(accountName, password string, m *modal.CreatePasswordModal) bool {
				go func() {
					_, err := pg.wallet.CreateNewAccount(accountName, []byte(password))
					if err != nil {
						m.SetError(err.Error())
						m.SetLoading(false)
						return
					}
					pg.Toast.Notify(values.String(values.StrAcctCreated))
					pg.loadWalletAccount()
					m.Dismiss()
				}()
				return false
			})
		pg.ParentWindow().ShowModal(newPasswordModal)
		break
	}

	if clicked, selectedItem := pg.accountsList.ItemClicked(); clicked {
		pg.ParentNavigator().Display(NewAcctDetailsPage(pg.Load, pg.accounts[selectedItem].Account))
	}
}

func (pg *WalletSettingsPage) resetDexDataModal() {
	// Show confirm modal before resetting dex client data.
	confirmModal := modal.NewInfoModal(pg.Load).
		Title(values.String(values.StrConfirmDexReset)).
		Body(values.String(values.StrDexResetInfo)).
		NegativeButton(values.String(values.StrCancel), func() {}).
		NegativeButtonStyle(pg.Theme.Color.Primary, pg.Theme.Color.Surface).
		PositiveButton(values.String(values.StrResetDexClient), func(isChecked bool) bool {
			if pg.Dexc().Reset() {
				pg.Toast.Notify("DEX client data reset complete.")
			} else {
				pg.Toast.NotifyError("DEX client data reset failed. Check the logs.")
			}
			return true
		})
	pg.ParentWindow().ShowModal(confirmModal)
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *WalletSettingsPage) OnNavigatedFrom() {}

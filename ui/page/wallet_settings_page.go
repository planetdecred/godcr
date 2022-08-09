package page

import (
	"gioui.org/layout"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const WalletSettingsPageID = "WalletSettings"

type WalletSettingsPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	wallet *dcrlibwallet.Wallet

	pageContainer layout.List

	changePass, rescan, resetDexData           *decredmaterial.Clickable
	changeAccount, checklog, checkStats        *decredmaterial.Clickable
	changeWalletName, addAccount, deleteWallet *decredmaterial.Clickable
	verifyMessage, validateMsg, signMessage    *decredmaterial.Clickable

	chevronRightIcon *decredmaterial.Icon
	backButton       decredmaterial.IconButton
	infoButton       decredmaterial.IconButton

	fetchProposal     *decredmaterial.Switch
	proposalNotif     *decredmaterial.Switch
	spendUnconfirmed  *decredmaterial.Switch
	SpendUnmixedFunds *decredmaterial.Switch
	connectToPeer     *decredmaterial.Switch
}

func NewWalletSettingsPage(l *load.Load) *WalletSettingsPage {
	pg := &WalletSettingsPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(WalletSettingsPageID),
		wallet:           l.WL.SelectedWallet.Wallet,
		changePass:       l.Theme.NewClickable(false),
		rescan:           l.Theme.NewClickable(false),
		resetDexData:     l.Theme.NewClickable(false),
		changeAccount:    l.Theme.NewClickable(false),
		checklog:         l.Theme.NewClickable(false),
		checkStats:       l.Theme.NewClickable(false),
		changeWalletName: l.Theme.NewClickable(false),
		addAccount:       l.Theme.NewClickable(false),
		deleteWallet:     l.Theme.NewClickable(false),
		verifyMessage:    l.Theme.NewClickable(false),
		validateMsg:      l.Theme.NewClickable(false),
		signMessage:      l.Theme.NewClickable(false),

		fetchProposal:     l.Theme.Switch(),
		proposalNotif:     l.Theme.Switch(),
		spendUnconfirmed:  l.Theme.Switch(),
		SpendUnmixedFunds: l.Theme.Switch(),
		connectToPeer:     l.Theme.Switch(),

		chevronRightIcon: decredmaterial.NewIcon(l.Theme.Icons.ChevronRight),
		pageContainer:    layout.List{Axis: layout.Vertical},
	}

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *WalletSettingsPage) OnNavigatedTo() {
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *WalletSettingsPage) Layout(gtx C) D {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrSettings),
			BackButton: pg.backButton,
			Back: func() {
				pg.ParentNavigator().CloseCurrentPage()
			},
			Body: func(gtx C) D {
				w := []func(gtx C) D{
					pg.generalSection(),
					pg.account(),
					pg.securityTools(),
					pg.debug(),
					pg.dangerZone(),
				}

				return pg.pageContainer.Layout(gtx, len(w), func(gtx C, i int) D {
					return w[i](gtx)
				})
			},
		}
		return sp.Layout(pg.ParentWindow(), gtx)
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
			layout.Rigid(pg.subSectionSwitch("Proposal notification", pg.proposalNotif)),
			layout.Rigid(pg.subSectionSwitch("Spend unconfirmed funds", pg.spendUnconfirmed)),
			layout.Rigid(pg.subSectionSwitch(values.String(values.StrAllowSpendingFromUnmixedAccount), pg.SpendUnmixedFunds)),
			layout.Rigid(pg.subSectionSwitch(values.String(values.StrConnectToSpecificPeer), pg.connectToPeer)),
		)
	}

	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrGeneral), dim)
	}
}

func (pg *WalletSettingsPage) account() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, values.String(values.StrAccount),
			pg.sectionContent(pg.addAccount, values.String(values.StrAddNewAccount)))
	}
}

func (pg *WalletSettingsPage) debug() layout.Widget {
	dims := func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.sectionContent(pg.rescan, values.String(values.StrRescanBlockchain))),
			layout.Rigid(pg.sectionContent(pg.checklog, "Check wallet log")),
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
			layout.Rigid(pg.sectionContent(pg.validateMsg, values.String(values.StrValidateMsg))),
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
							pg.chevronRightIcon.Color = pg.Theme.Color.Gray1
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

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *WalletSettingsPage) HandleUserInteractions() {
	for pg.changePass.Clicked() {
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
		break
	}

	for pg.rescan.Clicked() {
		go func() {
			info := modal.NewInfoModal(pg.Load).
				Title(values.String(values.StrRescanBlockchain)).
				Body(values.String(values.StrRescanInfo)).
				NegativeButton(values.String(values.StrCancel), func() {}).
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
		warningMsg := values.String(values.StrWalletRemoveInfo)
		if pg.wallet.IsWatchingOnlyWallet() {
			warningMsg = values.String(values.StrWatchOnlyWalletRemoveInfo)
		}
		confirmRemoveWalletModal := modal.NewInfoModal(pg.Load)
		confirmRemoveWalletModal.Title(values.String(values.StrRemoveWallet)).
			Body(warningMsg).
			NegativeButton(values.String(values.StrCancel), func() {}).
			PositiveButtonStyle(pg.Load.Theme.Color.Surface, pg.Load.Theme.Color.Danger).
			PositiveButton(values.String(values.StrRemove), func(isChecked bool) bool {
				walletDeleted := func() {
					if pg.WL.MultiWallet.LoadedWalletsCount() > 0 {
						pg.Toast.Notify(values.String(values.StrWalletRemoved))
						confirmRemoveWalletModal.Dismiss()
						pg.ParentNavigator().CloseCurrentPage()
						onWalSelected := func() {
							pg.ParentWindow().CloseCurrentPage()
						}
						onDexServerSelected := func(server string) {
							log.Info("Not implemented yet...", server)
						}
						pg.ParentWindow().Display(NewWalletDexServerSelector(pg.Load, onWalSelected, onDexServerSelected))
					} else {
						confirmRemoveWalletModal.Dismiss()
						pg.ParentWindow().CloseAllPages()
					}
				}

				if pg.wallet.IsWatchingOnlyWallet() {
					confirmRemoveWalletModal.SetLoading(true)
					go func() {
						// no password is required for watching only wallets.
						err := pg.WL.MultiWallet.DeleteWallet(pg.wallet.ID, nil)
						if err != nil {
							pg.Toast.NotifyError(err.Error())
							confirmRemoveWalletModal.SetLoading(false)
						} else {
							walletDeleted()
						}
					}()
					return false
				}

				walletPasswordModal := modal.NewPasswordModal(pg.Load).
					Title(values.String(values.StrConfirmToRemove)).
					NegativeButton(values.String(values.StrCancel), func() {}).
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
		pg.ParentWindow().ShowModal(confirmRemoveWalletModal)
	}

	if pg.infoButton.Button.Clicked() {
		info := modal.NewInfoModal(pg.Load).
			Title(values.String(values.StrSpendingPassword)).
			Body(values.String(values.StrSpendingPasswordInfo)).
			PositiveButton(values.String(values.StrGotIt), func(isChecked bool) bool {
				return true
			})
		pg.ParentWindow().ShowModal(info)
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *WalletSettingsPage) OnNavigatedFrom() {}

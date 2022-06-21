package privacy

import (
	"context"

	"gioui.org/layout"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/renderers"
	"github.com/planetdecred/godcr/ui/values"
)

const ManualMixerSetupPageID = "ManualMixerSetup"

type ManualMixerSetupPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	wallet                 *dcrlibwallet.Wallet
	mixedAccountSelector   *components.AccountSelector
	unmixedAccountSelector *components.AccountSelector

	backButton     decredmaterial.IconButton
	infoButton     decredmaterial.IconButton
	toPrivacySetup decredmaterial.Button
}

func NewManualMixerSetupPage(l *load.Load, wallet *dcrlibwallet.Wallet) *ManualMixerSetupPage {
	pg := &ManualMixerSetupPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(ManualMixerSetupPageID),
		wallet:           wallet,
		toPrivacySetup:   l.Theme.Button("Set up"),
	}

	// Mixed account picker
	pg.mixedAccountSelector = components.NewAccountSelector(l, wallet).
		Title("Mixed account").
		AccountSelected(func(selectedAccount *dcrlibwallet.Account) {}).
		AccountValidator(func(account *dcrlibwallet.Account) bool {
			wal := pg.Load.WL.MultiWallet.WalletWithID(account.WalletID)

			var unmixedAccNo int32 = -1
			if unmixedAcc := pg.unmixedAccountSelector.SelectedAccount(); unmixedAcc != nil {
				unmixedAccNo = unmixedAcc.Number
			}

			// Imported, watch only and default wallet accounts are invalid to use as a mixed account
			accountIsValid := account.Number != load.MaxInt32 && !wal.IsWatchingOnlyWallet() && account.Number != dcrlibwallet.DefaultAccountNum

			if !accountIsValid || account.Number == unmixedAccNo {
				return false
			}

			return true
		})

	// Unmixed account picker
	pg.unmixedAccountSelector = components.NewAccountSelector(l, wallet).
		Title("Unmixed account").
		AccountSelected(func(selectedAccount *dcrlibwallet.Account) {}).
		AccountValidator(func(account *dcrlibwallet.Account) bool {
			wal := pg.Load.WL.MultiWallet.WalletWithID(account.WalletID)

			var mixedAccNo int32 = -1
			if mixedAcc := pg.mixedAccountSelector.SelectedAccount(); mixedAcc != nil {
				mixedAccNo = mixedAcc.Number
			}

			// Imported, watch only and default wallet accounts are invalid to use as an unmixed account
			accountIsValid := account.Number != load.MaxInt32 && !wal.IsWatchingOnlyWallet() && account.Number != dcrlibwallet.DefaultAccountNum

			// Account is invalid if already selected by mixed account selector.
			if !accountIsValid || account.Number == mixedAccNo {
				return false
			}

			return true
		})

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *ManualMixerSetupPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())

	pg.mixedAccountSelector.SelectFirstWalletValidAccount(pg.wallet)
	pg.unmixedAccountSelector.SelectFirstWalletValidAccount(pg.wallet)
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *ManualMixerSetupPage) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		page := components.SubPage{
			Load:       pg.Load,
			Title:      "Manual setup",
			WalletName: pg.wallet.Name,
			BackButton: pg.backButton,
			Back: func() {
				pg.ParentNavigator().CloseCurrentPage()
			},
			Body: func(gtx C) D {
				return pg.Theme.Card().Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return pg.mixerAccountSections(gtx, "Mixed account", func(gtx layout.Context) layout.Dimensions {
										return pg.mixedAccountSelector.Layout(pg.ParentWindow(), gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Top: values.MarginPaddingMinus15}.Layout(gtx, func(gtx C) D {
										return pg.mixerAccountSections(gtx, "Unmixed account", func(gtx layout.Context) layout.Dimensions {
											return pg.unmixedAccountSelector.Layout(pg.ParentWindow(), gtx)
										})
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Top: values.MarginPadding10, Left: values.MarginPadding16, Right: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
										return layout.Flex{
											Axis: layout.Horizontal,
										}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												return pg.Theme.Icons.ActionInfo.Layout(gtx, pg.Theme.Color.Gray1)
											}),
											layout.Rigid(func(gtx C) D {
												txt := `<span style="text-color: grayText2">
											<b>Make sure to select the same accounts from the previous privacy setup. </b><br>Failing to do so could compromise wallet privacy.<br> You may not select the same account for mixed and unmixed.
										</span>`
												return layout.Inset{
													Left: values.MarginPadding8,
												}.Layout(gtx, renderers.RenderHTML(txt, pg.Theme).Layout)
											}),
										)
									})
								}),
							)
						}),
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.UniformInset(values.MarginPadding15).Layout(gtx, pg.toPrivacySetup.Layout)
						}),
					)
				})
			},
		}
		return page.Layout(pg.ParentWindow(), gtx)
	}

	return components.UniformPadding(gtx, body)
}

func (pg *ManualMixerSetupPage) mixerAccountSections(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
	return pg.Theme.Card().Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{
						Bottom: values.MarginPadding8,
					}.Layout(gtx, pg.Theme.Body1(title).Layout)
				}),
				layout.Rigid(body),
			)
		})
	})
}

func (pg *ManualMixerSetupPage) showModalSetupMixerAcct() {
	if pg.mixedAccountSelector.SelectedAccount().Number == pg.unmixedAccountSelector.SelectedAccount().Number {
		pg.Toast.NotifyError("Cannot use same account for mixed & unmixed")
		return
	}

	passwordModal := modal.NewPasswordModal(pg.Load).
		Title("Confirm to set mixer accounts").
		NegativeButton("Cancel", func() {}).
		PositiveButton("Confirm", func(password string, pm *modal.PasswordModal) bool {
			go func() {
				mixedAcctNumber := pg.mixedAccountSelector.SelectedAccount().Number
				unmixedAcctNumber := pg.unmixedAccountSelector.SelectedAccount().Number
				err := pg.wallet.SetAccountMixerConfig(mixedAcctNumber, unmixedAcctNumber, password)
				if err != nil {
					pm.SetError(err.Error())
					pm.SetLoading(false)
					return
				}
				pg.WL.MultiWallet.SetBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, true)

				// rename mixed account
				err = pg.wallet.RenameAccount(mixedAcctNumber, "mixed")
				if err != nil {
					log.Error(err)
				}

				// rename unmixed account
				err = pg.wallet.RenameAccount(unmixedAcctNumber, "unmixed")
				if err != nil {
					log.Error(err)
				}

				pm.Dismiss()

				pg.ParentNavigator().Display(NewAccountMixerPage(pg.Load, pg.wallet))
			}()

			return false
		})
	pg.ParentWindow().ShowModal(passwordModal)
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *ManualMixerSetupPage) HandleUserInteractions() {
	if pg.toPrivacySetup.Clicked() {
		go pg.showModalSetupMixerAcct()
	}

	if pg.mixedAccountSelector.SelectedAccount().Number == pg.unmixedAccountSelector.SelectedAccount().Number {
		pg.toPrivacySetup.SetEnabled(false)
	} else {
		pg.toPrivacySetup.SetEnabled(true)
	}

	// Disable set up button if either mixed or unmixed account is the default account.
	if pg.mixedAccountSelector.SelectedAccount().Number == dcrlibwallet.DefaultAccountNum ||
		pg.unmixedAccountSelector.SelectedAccount().Number == dcrlibwallet.DefaultAccountNum {
		pg.toPrivacySetup.SetEnabled(false)
	}

}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *ManualMixerSetupPage) OnNavigatedFrom() {
	pg.ctxCancel()
}

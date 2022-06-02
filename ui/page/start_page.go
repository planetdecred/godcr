package page

import (
	"os"

	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/wallets"
	"github.com/planetdecred/godcr/ui/values"
)

const StartPageID = "start_page"

type startPage struct {
	*load.Load

	createButton          decredmaterial.Button
	restoreButton         decredmaterial.Button
	watchOnlyWalletButton decredmaterial.Button

	loading bool
}

func NewStartPage(l *load.Load) load.Page {
	sp := &startPage{
		Load:    l,
		loading: true,

		watchOnlyWalletButton: l.Theme.Button(values.String(values.StrImportWatchingOnlyWallet)),
		createButton:          l.Theme.Button(values.String(values.StrCreateANewWallet)),
		restoreButton:         l.Theme.Button(values.String(values.StrRestoreExistingWallet)),
	}

	return sp
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (sp *startPage) ID() string {
	return StartPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (sp *startPage) OnNavigatedTo() {
	sp.WL.MultiWallet = sp.WL.Wallet.GetMultiWallet()

	// refresh theme now that config is available
	sp.RefreshTheme()

	if sp.WL.MultiWallet.LoadedWalletsCount() > 0 {
		if sp.WL.MultiWallet.IsStartupSecuritySet() {
			sp.unlock()
		} else {
			go sp.openWallets("")
		}

	} else {
		sp.loading = false
	}
}

func (sp *startPage) unlock() {
	modal.NewPasswordModal(sp.Load).
		Title(values.String(values.StrUnlockWithPassword)).
		Hint(values.String(values.StrStartupPassword)).
		NegativeButton(values.String(values.StrExit), func() {
			sp.WL.MultiWallet.Shutdown()
			os.Exit(0)
		}).
		PositiveButton(values.String(values.StrUnlock), func(password string, m *modal.PasswordModal) bool {
			go func() {
				err := sp.openWallets(password)
				if err != nil {
					m.SetError(translateErr(err))
					m.SetLoading(false)
					return
				}

				m.Dismiss()
			}()
			return false
		}).Show()
}

func (sp *startPage) openWallets(password string) error {
	err := sp.WL.MultiWallet.OpenWallets([]byte(password))
	if err != nil {
		log.Info("Error opening wallet:", err)
		// show err dialog
		return err
	}

	sp.ChangeWindowPage(NewMainPage(sp.Load), false)
	return nil
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (sp *startPage) HandleUserInteractions() {
	for sp.createButton.Clicked() {
		modal.NewCreatePasswordModal(sp.Load).
			Title(values.String(values.StrCreateANewWallet)).
			PasswordCreated(func(_, password string, m *modal.CreatePasswordModal) bool {
				go func() {
					wal, err := sp.WL.MultiWallet.CreateNewWallet("mywallet", password, dcrlibwallet.PassphraseTypePass)
					if err != nil {
						m.SetError(err.Error())
						m.SetLoading(false)
						return
					}
					err = wal.CreateMixerAccounts("mixed", "unmixed", password)
					if err != nil {
						m.SetError(err.Error())
						m.SetLoading(false)
						return
					}
					sp.WL.MultiWallet.SetBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, true)
					m.Dismiss()

					sp.ChangeWindowPage(NewMainPage(sp.Load), false)
				}()
				return false
			}).Show()
	}

	for sp.restoreButton.Clicked() {
		afterRestore := func() {
			sp.ChangeWindowPage(NewMainPage(sp.Load), false)
		}
		sp.ChangeWindowPage(wallets.NewRestorePage(sp.Load, afterRestore), true)
	}

	for sp.watchOnlyWalletButton.Clicked() {
		modal.NewCreateWatchOnlyModal(sp.Load).
			EnableName(false).
			WatchOnlyCreated(func(_, password string, m *modal.CreateWatchOnlyModal) bool {
				go func() {
					_, err := sp.WL.MultiWallet.CreateWatchOnlyWallet("mywallet", password)
					if err != nil {
						m.SetError(err.Error())
						m.SetLoading(false)
						return
					}
					m.Dismiss()

					sp.ChangeWindowPage(NewMainPage(sp.Load), false)
				}()
				return false
			}).Show()
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (sp *startPage) OnNavigatedFrom() {}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (sp *startPage) Layout(gtx C) D {
	gtx.Constraints.Min = gtx.Constraints.Max // use maximum height & width
	return layout.Stack{Alignment: layout.N}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return sp.loadingSection(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if sp.loading {
						return D{}
					}

					return sp.buttonSection(gtx)
				}),
			)
		}),
	)
}

func (sp *startPage) loadingSection(gtx C) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X // use maximum width
	if sp.loading {
		gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	} else {
		gtx.Constraints.Min.Y = (gtx.Constraints.Max.Y * 65) / 100 // use 65% of view height
	}

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle, Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return sp.Theme.Icons.DecredSymbolIcon.LayoutSize(gtx, values.MarginPadding150)
					})
				}),
				layout.Rigid(func(gtx C) D {
					netType := sp.WL.Wallet.Net
					if sp.WL.Wallet.Net == dcrlibwallet.Testnet3 {
						netType = "Testnet"
					}
					nType := sp.Theme.Label(values.TextSize20, netType)
					nType.Font.Weight = text.Medium
					return nType.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					if sp.loading {
						loadStatus := sp.Theme.Label(values.TextSize20, values.String(values.StrLoading))
						if sp.WL.MultiWallet.LoadedWalletsCount() > 0 {
							loadStatus.Text = values.String(values.StrOpeningWallet)
						}

						return layout.Inset{Top: values.MarginPadding24}.Layout(gtx, loadStatus.Layout)
					}

					welcomeText := sp.Theme.Label(values.TextSize24, values.String(values.StrWelcomeNote))
					return layout.Inset{Top: values.MarginPadding24}.Layout(gtx, welcomeText.Layout)
				}),
			)
		}),
	)
}

func (sp *startPage) buttonSection(gtx C) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X              // use maximum width
	gtx.Constraints.Min.Y = (gtx.Constraints.Max.Y * 35) / 100 // use 35% of view height
	return layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			return layout.Flex{Alignment: layout.Middle, Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					gtx.Constraints.Max.X = gtx.Px(values.AppWidth) // set button with to app width
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding24, Right: values.MarginPadding24}.Layout(gtx, func(gtx C) D {
								return sp.createButton.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Top: values.MarginPadding24, Left: values.MarginPadding24, Right: values.MarginPadding24}.Layout(gtx, func(gtx C) D {
								return sp.restoreButton.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Top: values.MarginPadding24, Bottom: values.MarginPadding24, Left: values.MarginPadding24, Right: values.MarginPadding24}.Layout(gtx, func(gtx C) D {
								return sp.watchOnlyWalletButton.Layout(gtx)
							})
						}),
					)
				}),
			)
		}),
	)
}

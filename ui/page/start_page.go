package page

import (
	"os"

	"gioui.org/layout"
	"gioui.org/text"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/values"
	// "github.com/planetdecred/godcr/ui/page/wallets"
)

const StartPageID = "start_page"

type startPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	addWalletButton decredmaterial.Button

	// to be removed after full layout migration
	newlayout    decredmaterial.Button
	legacyLayout decredmaterial.Button

	loading bool
}

func NewStartPage(l *load.Load) app.Page {
	sp := &startPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(StartPageID),
		loading:          true,

		addWalletButton: l.Theme.Button(values.String(values.StrAddWallet)),

		// to be removed after full layout migration
		newlayout:    l.Theme.Button("Continue v2 layout"),
		legacyLayout: l.Theme.Button("Continue legacy layout"),
	}

	return sp
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (sp *startPage) OnNavigatedTo() {
	sp.WL.MultiWallet = sp.WL.Wallet.GetMultiWallet()

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
	startupPasswordModal := modal.NewPasswordModal(sp.Load).
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
		})
	sp.ParentWindow().ShowModal(startupPasswordModal)
}

func (sp *startPage) openWallets(password string) error {
	err := sp.WL.MultiWallet.OpenWallets([]byte(password))
	if err != nil {
		log.Info("Error opening wallet:", err)
		// show err dialog
		return err
	}

	onWalSelected := func() {
		sp.ParentNavigator().ClearStackAndDisplay(NewMainPage(sp.Load))
	}
	sp.ParentNavigator().ClearStackAndDisplay(NewWalletList(sp.Load, onWalSelected))
	return nil
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (sp *startPage) HandleUserInteractions() {
	// <<<<<<< HEAD
	// 	for sp.createButton.Clicked() {
	// 		spendingPasswordModal := modal.NewCreatePasswordModal(sp.Load).
	// 			Title(values.String(values.StrCreateANewWallet)).
	// 			PasswordCreated(func(_, password string, m *modal.CreatePasswordModal) bool {
	// 				go func() {
	// 					wal, err := sp.WL.MultiWallet.CreateNewWallet("mywallet", password, dcrlibwallet.PassphraseTypePass)
	// 					if err != nil {
	// 						m.SetError(err.Error())
	// 						m.SetLoading(false)
	// 						return
	// 					}
	// 					err = wal.CreateMixerAccounts("mixed", "unmixed", password)
	// 					if err != nil {
	// 						m.SetError(err.Error())
	// 						m.SetLoading(false)
	// 						return
	// 					}
	// 					sp.WL.MultiWallet.SetBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, true)
	// 					m.Dismiss()

	// 					sp.ParentNavigator().ClearStackAndDisplay(NewMainPage(sp.Load))
	// 				}()
	// 				return false
	// 			})
	// 		sp.ParentWindow().ShowModal(spendingPasswordModal)
	// 	}

	// 	for sp.restoreButton.Clicked() {
	// =======
	// for sp.createButton.Clicked() {
	// 	modal.NewCreatePasswordModal(sp.Load).
	// 		Title(values.String(values.StrCreateANewWallet)).
	// 		PasswordCreated(func(_, password string, m *modal.CreatePasswordModal) bool {
	// 			go func() {
	// 				wal, err := sp.WL.MultiWallet.CreateNewWallet("mywallet", password, dcrlibwallet.PassphraseTypePass)
	// 				if err != nil {
	// 					m.SetError(err.Error())
	// 					m.SetLoading(false)
	// 					return
	// 				}
	// 				err = wal.CreateMixerAccounts("mixed", "unmixed", password)
	// 				if err != nil {
	// 					m.SetError(err.Error())
	// 					m.SetLoading(false)
	// 					return
	// 				}
	// 				sp.WL.MultiWallet.SetBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, true)
	// 				m.Dismiss()

	// 				sp.ChangeWindowPage(NewMainPage(sp.Load), false)
	// 			}()
	// 			return false
	// 		}).Show()
	// }

	// 	for sp.addWalletButton.Clicked() {
	// >>>>>>> temp
	// 		afterRestore := func() {
	// 			sp.ParentNavigator().ClearStackAndDisplay(NewMainPage(sp.Load))
	// 		}
	// 		sp.ParentNavigator().Display(wallets.NewRestorePage(sp.Load, afterRestore))
	// 	}

	// <<<<<<< HEAD
	// 	for sp.watchOnlyWalletButton.Clicked() {
	// 		createWatchOnlyModal := modal.NewCreateWatchOnlyModal(sp.Load).
	// 			EnableName(false).
	// 			WatchOnlyCreated(func(_, password string, m *modal.CreateWatchOnlyModal) bool {
	// 				go func() {
	// 					_, err := sp.WL.MultiWallet.CreateWatchOnlyWallet("mywallet", password)
	// 					if err != nil {
	// 						m.SetError(err.Error())
	// 						m.SetLoading(false)
	// 						return
	// 					}
	// 					m.Dismiss()

	// 					sp.ParentNavigator().ClearStackAndDisplay(NewMainPage(sp.Load))
	// 				}()
	// 				return false
	// 			})
	// 		sp.ParentWindow().ShowModal(createWatchOnlyModal)
	// 	}
	// =======
	// for sp.watchOnlyWalletButton.Clicked() {
	// 	modal.NewCreateWatchOnlyModal(sp.Load).
	// 		EnableName(false).
	// 		WatchOnlyCreated(func(_, password string, m *modal.CreateWatchOnlyModal) bool {
	// 			go func() {
	// 				_, err := sp.WL.MultiWallet.CreateWatchOnlyWallet("mywallet", password)
	// 				if err != nil {
	// 					m.SetError(err.Error())
	// 					m.SetLoading(false)
	// 					return
	// 				}
	// 				m.Dismiss()

	// 				sp.ChangeWindowPage(NewMainPage(sp.Load), false)
	// 			}()
	// 			return false
	// 		}).Show()
	// }
	// >>>>>>> temp
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
	return layout.Flex{
		Alignment: layout.Middle,
		Axis:      layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return sp.loadingSection(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			if sp.loading {
				return D{}
			}

			gtx.Constraints.Max.X = gtx.Dp(values.MarginPadding350)
			return layout.Inset{
				Left:  values.MarginPadding24,
				Right: values.MarginPadding24,
			}.Layout(gtx, sp.addWalletButton.Layout)
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
					return layout.Inset{Top: values.MarginPadding14}.Layout(gtx, nType.Layout)
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
					welcomeText.Color = sp.Theme.Color.GrayText1
					return layout.Inset{Top: values.MarginPadding24}.Layout(gtx, welcomeText.Layout)
				}),
			)
		}),
	)
}

// <<<<<<< HEAD

// func (sp *startPage) buttonSection(gtx C) D {
// 	gtx.Constraints.Min.X = gtx.Constraints.Max.X              // use maximum width
// 	gtx.Constraints.Min.Y = (gtx.Constraints.Max.Y * 35) / 100 // use 35% of view height
// 	return layout.Stack{Alignment: layout.S}.Layout(gtx,
// 		layout.Stacked(func(gtx C) D {
// 			return layout.Flex{Alignment: layout.Middle, Axis: layout.Vertical}.Layout(gtx,
// 				layout.Rigid(func(gtx layout.Context) layout.Dimensions {

// 					gtx.Constraints.Max.X = gtx.Dp(values.AppWidth) // set button with to app width
// 					gtx.Constraints.Min.X = gtx.Constraints.Max.X
// 					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
// 						layout.Rigid(func(gtx C) D {
// 							return layout.Inset{Left: values.MarginPadding24, Right: values.MarginPadding24}.Layout(gtx, func(gtx C) D {
// 								return sp.createButton.Layout(gtx)
// 							})
// 						}),
// 						layout.Rigid(func(gtx C) D {
// 							return layout.Inset{Top: values.MarginPadding24, Left: values.MarginPadding24, Right: values.MarginPadding24}.Layout(gtx, func(gtx C) D {
// 								return sp.restoreButton.Layout(gtx)
// 							})
// 						}),
// 						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
// 							return layout.Inset{Top: values.MarginPadding24, Bottom: values.MarginPadding24, Left: values.MarginPadding24, Right: values.MarginPadding24}.Layout(gtx, func(gtx C) D {
// 								return sp.watchOnlyWalletButton.Layout(gtx)
// 							})
// 						}),
// 					)
// 				}),
// 			)
// 		}),
// 	)
// }
// =======
// >>>>>>> temp

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
	"github.com/planetdecred/godcr/ui/page/wallets"
	"github.com/planetdecred/godcr/ui/values"
)

const StartPageID = "start_page"

type startPage struct {
	*app.App

	createButton  decredmaterial.Button
	restoreButton decredmaterial.Button

	loading bool
}

func NewStartPage(app *app.App) load.Page {
	sp := &startPage{
		App:     app,
		loading: true,

		createButton:  app.Theme.Button("Create a new wallet"),
		restoreButton: app.Theme.Button("Restore an existing wallet"),
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
	if sp.MultiWallet().LoadedWalletsCount() == 0 {
		sp.loading = false
	} else if sp.MultiWallet().IsStartupSecuritySet() {
		sp.unlock()
	} else {
		go sp.openWallets("")
	}
}

func (sp *startPage) unlock() {
	modal.NewPasswordModal(sp.Theme, nil).
		Title("Unlock with password").
		Hint("Startup password").
		NegativeButton("Exit", func() {
			sp.MultiWallet().Shutdown()
			os.Exit(0)
		}).
		PositiveButton("Unlock", func(password string, m *modal.PasswordModal) bool {
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
	err := sp.MultiWallet().OpenWallets([]byte(password))
	if err != nil {
		log.Info("Error opening wallet:", err)
		// show err dialog
		return err
	}

	sp.ChangePage(NewMainPage(sp.App), false)
	return nil
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (sp *startPage) HandleUserInteractions() {
	for sp.createButton.Clicked() {
		modal.NewCreatePasswordModal(nil).
			Title("Create new wallet").
			PasswordCreated(func(_, password string, m *modal.CreatePasswordModal) bool {
				go func() {
					_, err := sp.MultiWallet().CreateNewWallet("mywallet", password, dcrlibwallet.PassphraseTypePass)
					if err != nil {
						m.SetError(err.Error())
						m.SetLoading(false)
						return
					}
					m.Dismiss()

					sp.ChangePage(NewMainPage(sp.App), false)
				}()
				return false
			}).Show()
	}

	for sp.restoreButton.Clicked() {
		afterRestore := func() {
			sp.ChangePage(NewMainPage(sp.App), false)
		}
		sp.ChangePage(wallets.NewRestorePage(nil, afterRestore), true)
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

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (sp *startPage) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min = gtx.Constraints.Max // use maximum height & width
	return layout.Stack{Alignment: layout.N}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return sp.loadingSection(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if sp.loading {
						return layout.Dimensions{}
					}

					return sp.buttonSection(gtx)
				}),
			)
		}),
	)
}

func (sp *startPage) loadingSection(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X // use maximum width
	if sp.loading {
		gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
	} else {
		gtx.Constraints.Min.Y = (gtx.Constraints.Max.Y * 75) / 100 // use 75% of view height
	}

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle, Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return sp.Theme.Icons.DecredSymbolIcon.LayoutSize(gtx, values.MarginPadding150)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					netType := sp.MultiWallet().NetType()
					if netType == dcrlibwallet.Testnet3 {
						netType = "Testnet"
					}
					nType := sp.Theme.Label(values.TextSize20, netType)
					nType.Font.Weight = text.Medium
					return nType.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if sp.loading {
						loadStatus := sp.Theme.Label(values.TextSize20, "Loading")
						if sp.MultiWallet().LoadedWalletsCount() > 0 {
							loadStatus.Text = "Opening wallets"
						}

						return layout.Inset{Top: values.MarginPadding24}.Layout(gtx, loadStatus.Layout)
					}

					welcomeText := sp.Theme.Label(values.TextSize24, "Welcome to Decred Wallet, a secure & open-source mobile wallet.")
					return layout.Inset{Top: values.MarginPadding24}.Layout(gtx, welcomeText.Layout)
				}),
			)
		}),
	)
}

func (sp *startPage) buttonSection(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X              // use maximum width
	gtx.Constraints.Min.Y = (gtx.Constraints.Max.Y * 25) / 100 // use 25% of view height
	return layout.Stack{Alignment: layout.S}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Middle, Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Px(values.AppWidth) // set button with to app width
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Top: values.MarginPadding24, Left: values.MarginPadding24, Right: values.MarginPadding24}.Layout(gtx, func(gtx C) D {
								return sp.createButton.Layout(gtx)
							})
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Top: values.MarginPadding24, Bottom: values.MarginPadding24, Left: values.MarginPadding24, Right: values.MarginPadding24}.Layout(gtx, func(gtx C) D {
								return sp.restoreButton.Layout(gtx)
							})
						}),
					)
				}),
			)
		}),
	)
}

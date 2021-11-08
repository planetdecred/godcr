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

	loading bool

	decredSymbol  *decredmaterial.Image
	networkType   decredmaterial.Label
	loadStatus    decredmaterial.Label
	welcomeText   decredmaterial.Label
	createButton  decredmaterial.Button
	restoreButton decredmaterial.Button
}

func NewStartPage(l *load.Load) load.Page {
	sp := &startPage{
		Load: l,

		loading: true,

		decredSymbol: l.Icons.DecredSymbolIcon,
		networkType:  l.Theme.Label(values.TextSize20, "testnet"),
		loadStatus:   l.Theme.Label(values.TextSize20, "Loading"),
		welcomeText:  l.Theme.Label(values.TextSize24, "Welcome to Decred Wallet, a secure & open-source mobile wallet."),

		createButton:  l.Theme.Button("Create a new wallet"),
		restoreButton: l.Theme.Button("Restore an existing wallet"),
	}

	sp.networkType.Color = l.Theme.Color.DeepBlue
	sp.networkType.Font.Weight = text.Medium

	sp.loadStatus.Color = l.Theme.Color.DeepBlue
	sp.welcomeText.Color = l.Theme.Color.DeepBlue

	return sp
}

func (sp *startPage) ID() string {
	return StartPageID
}

func (sp *startPage) OnResume() {
	sp.WL.Wallet.InitMultiWallet()
	sp.WL.MultiWallet = sp.WL.Wallet.GetMultiWallet()

	// refresh theme now that config is available
	sp.RefreshTheme()

	if sp.WL.MultiWallet.LoadedWalletsCount() > 0 {
		sp.loadStatus.Text = "Opening wallets"

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
		Title("Unlock with password").
		Hint("Startup password").
		NegativeButton("Exit", func() {
			sp.WL.MultiWallet.Shutdown()
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
	err := sp.WL.MultiWallet.OpenWallets([]byte(password))
	if err != nil {
		log.Info("Error opening wallet:", err)
		// show err dialog
		return err
	}

	sp.proceedToMainPage()
	return nil
}

func (sp *startPage) proceedToMainPage() {
	sp.WL.Wallet.SetupListeners()
	sp.ChangeWindowPage(NewMainPage(sp.Load), false)
}

func (sp *startPage) Handle() {
	for sp.createButton.Clicked() {
		modal.NewCreatePasswordModal(sp.Load).
			Title("Create new wallet").
			PasswordCreated(func(_, password string, m *modal.CreatePasswordModal) bool {
				go func() {
					_, err := sp.WL.MultiWallet.CreateNewWallet("mywallet", password, dcrlibwallet.PassphraseTypePass)
					if err != nil {
						m.SetError(err.Error())
						m.SetLoading(false)
						return
					}
					m.Dismiss()

					sp.proceedToMainPage()
				}()
				return false
			}).Show()
	}

	for sp.restoreButton.Clicked() {
		sp.ChangeWindowPage(wallets.NewRestorePage(sp.Load), true)
	}
}

func (sp *startPage) OnClose() {}

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
						return sp.decredSymbol.LayoutSize(gtx, values.MarginPadding150)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return sp.networkType.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if sp.loading {
						return layout.Inset{Top: values.MarginPadding24}.Layout(gtx, sp.loadStatus.Layout)
					}

					return layout.Inset{Top: values.MarginPadding24}.Layout(gtx, sp.welcomeText.Layout)
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

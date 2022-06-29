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
		Hint(values.String(values.StrStartupPassphrase)).
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
	for sp.addWalletButton.Clicked() {
		sp.ParentNavigator().Display(NewCreateWallet(sp.Load))
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
	if sp.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return sp.layoutMobile(gtx)
	}
	return sp.layoutDesktop(gtx)
}

// Desktop layout
func (sp *startPage) layoutDesktop(gtx C) D {
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
						return sp.Theme.Icons.DecredLogo.LayoutSize(gtx, values.MarginPadding150)
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
					welcomeText.Alignment = text.Middle
					return layout.Inset{Top: values.MarginPadding24}.Layout(gtx, welcomeText.Layout)
				}),
			)
		}),
	)
}

// Mobile layout
func (sp *startPage) layoutMobile(gtx C) D {
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

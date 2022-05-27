package wallets

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const CreateRestorePageID = "Restore"

var tabTitles = []string{"Seed Words", "Hex"}

type Restore struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the ParentNavigator that displayed this page
	// and the root WindowNavigator. The ParentNavigator is also the root
	// WindowNavigator if this page is displayed from the StartPage, otherwise
	// the ParentNavigator is the MainPage.
	*app.GenericPageModal
	restoreComplete func()
	tabList         *decredmaterial.ClickableList
	tabIndex        int
	backButton      decredmaterial.IconButton
	seedRestorePage *SeedRestore
}

func NewRestorePage(l *load.Load, onRestoreComplete func()) *Restore {
	pg := &Restore{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(CreateRestorePageID),           l,
		seedRestorePage: NewSeedRestorePage(l, onRestoreComplete),
		tabIndex:        0,
		tabList:         l.Theme.NewClickableList(layout.Horizontal),
		restoreComplete: onRestoreComplete,
	}

	pg.backButton, _ = components.SubpageHeaderButtons(l)
	pg.backButton.Icon = pg.Theme.Icons.ContentClear
	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *Restore) OnNavigatedTo() {
	pg.seedRestorePage.OnNavigatedTo()
}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *Restore) Layout(gtx C) D {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrRestoreWallet),
			BackButton: pg.backButton,
			Back: func() {
				pg.ParentNavigator().CloseCurrentPage()
			},
			Body: func(gtx C) D {
				return pg.restoreLayout(gtx)
			},
		}
		return sp.Layout(pg.ParentWindow(), gtx)
	}
	return components.UniformPadding(gtx, body)
}

func (pg *Restore) restoreLayout(gtx layout.Context) layout.Dimensions {
	return components.UniformPadding(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.tabLayout),
			layout.Rigid(pg.Theme.Separator().Layout),
			layout.Flexed(1, func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, pg.indexLayout)
			}),
		)
	})
}

func (pg *Restore) indexLayout(gtx C) D {
	return pg.seedRestorePage.Layout(gtx)
}

func (pg *Restore) switchTab(tabIndex int) {
	if tabIndex == 0 {
		pg.seedRestorePage.OnNavigatedTo()
	} else {
		pg.showHexRestoreModal()
	}
}

func (pg *Restore) tabLayout(gtx C) D {
	var dims layout.Dimensions
	return layout.Inset{
		Top: values.MarginPaddingMinus30,
	}.Layout(gtx, func(gtx C) D {
		return pg.tabList.Layout(gtx, len(tabTitles), func(gtx C, i int) D {
			return layout.Stack{Alignment: layout.S}.Layout(gtx,
				layout.Stacked(func(gtx C) D {
					return layout.Inset{
						Right:  values.MarginPadding24,
						Bottom: values.MarginPadding8,
					}.Layout(gtx, func(gtx C) D {
						return layout.Center.Layout(gtx, func(gtx C) D {
							lbl := pg.Theme.Label(values.TextSize16, tabTitles[i])
							lbl.Color = pg.Theme.Color.GrayText1
							if pg.tabIndex == i {
								lbl.Color = pg.Theme.Color.Primary
								dims = lbl.Layout(gtx)
							}

							return lbl.Layout(gtx)
						})
					})
				}),
				layout.Stacked(func(gtx C) D {
					if pg.tabIndex != i {
						return D{}
					}

					tabHeight := gtx.Px(values.MarginPadding2)
					tabRect := image.Rect(0, 0, dims.Size.X, tabHeight)

					return layout.Inset{
						Left: values.MarginPaddingMinus22,
					}.Layout(gtx, func(gtx C) D {
						paint.FillShape(gtx.Ops, pg.Theme.Color.Primary, clip.Rect(tabRect).Op())
						return layout.Dimensions{
							Size: image.Point{X: dims.Size.X, Y: tabHeight},
						}
					})
				}),
			)
		})
	})
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *Restore) OnNavigatedFrom() {
	pg.seedRestorePage.OnNavigatedFrom()

	pg.PopWindowPage()
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *Restore) HandleUserInteractions() {
	if clicked, selectedItem := pg.tabList.ItemClicked(); clicked {
		if pg.tabIndex != selectedItem {
			pg.tabIndex = selectedItem
			pg.switchTab(pg.tabIndex)
		}
	}

	if pg.tabIndex == 0 {
		pg.seedRestorePage.HandleUserInteractions()
	}
}

func (pg *Restore) showHexRestoreModal() {
	hexModal := modal.NewTextInputModal(pg.Load).
		Hint(values.String(values.StrEnterHex)).
		PositiveButtonStyle(pg.Load.Theme.Color.Primary, pg.Load.Theme.Color.InvText).
		PositiveButton(values.String(values.StrSubmit), func(hex string, hm *modal.TextInputModal) bool {
			go func() {
				if !pg.verifyHex(hex) {
					hm.Dismiss()
					pg.tabIndex = 0
					pg.switchTab(pg.tabIndex)
					return
				}

				modal.NewCreatePasswordModal(pg.Load).
					Title(values.String(values.StrEnterWalDetails)).
					EnableName(true).
					ShowWalletInfoTip(true).
					SetParent(pg).
					PasswordCreated(func(walletName, password string, m *modal.CreatePasswordModal) bool {
						go func() {
							_, err := pg.WL.MultiWallet.RestoreWallet(walletName, hex, password, dcrlibwallet.PassphraseTypePass)
							if err != nil {
								m.SetError(components.TranslateErr(err))
								m.SetLoading(false)
								return
							}

							pg.Toast.Notify(values.String(values.StrWalletRestored))
							m.Dismiss()
							// Close this page and return to the previous page (most likely wallets page)
							// if there's no restoreComplete callback function.
							if pg.restoreComplete == nil {
								pg.PopWindowPage()
							} else {
								pg.restoreComplete()
							}
						}()
						return false
					}).Show()

				hm.Dismiss()
			}()
			return false
		})
	hexModal.Title(values.String(values.StrRestoreWithHex)).
		NegativeButton(values.String(values.StrCancel), func() {
			pg.tabIndex = 0
			pg.switchTab(pg.tabIndex)
		})
	hexModal.Show()
}

func (pg *Restore) verifyHex(hex string) bool {
	if !dcrlibwallet.VerifySeed(hex) {
		pg.Toast.NotifyError(values.String(values.StrInvalidHex))
		return false
	}

	// Compare with existing wallets seed. On positive match abort import
	// to prevent duplicate wallet. walletWithSameSeed >= 0 if there is a match.
	walletWithSameSeed, err := pg.WL.MultiWallet.WalletWithSeed(hex)
	if err != nil {
		log.Error(err)
		return false
	}

	if walletWithSameSeed != -1 {
		pg.Toast.NotifyError(values.String(values.StrSeedAlreadyExist))
		return false
	}

	return true
}

func (pg *Restore) OnNavigatedFrom() {}

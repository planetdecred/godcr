package page

import (
	"sync"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/page/wallets"
	"github.com/planetdecred/godcr/ui/values"
)

const CreateWalletID = "create_wallet"

type walletType struct {
	clickable *decredmaterial.Clickable
	logo      *decredmaterial.Image
	name      string
}

type decredAction struct {
	title     string
	clickable *decredmaterial.Clickable
	action    func()
	border    decredmaterial.Border
	width     unit.Dp
}

type CreateWallet struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	listLock        sync.Mutex
	scrollContainer *widget.List
	list            layout.List

	walletTypes   []*walletType
	decredActions []*decredAction

	walletTypeList *decredmaterial.ClickableList

	walletName         decredmaterial.Editor
	restoreWalletName  decredmaterial.Editor
	watchOnlyWalletHex decredmaterial.Editor
	watchOnlyCheckBox  decredmaterial.CheckBoxStyle

	continueBtn decredmaterial.Button
	restoreBtn  decredmaterial.Button
	importBtn   decredmaterial.Button

	selectedWalletType         int
	selectedDecredWalletAction int
}

func NewCreateWallet(l *load.Load) *CreateWallet {
	pg := &CreateWallet{
		GenericPageModal: app.NewGenericPageModal(CreateWalletID),
		scrollContainer: &widget.List{
			List: layout.List{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			},
		},
		list: layout.List{Axis: layout.Vertical},

		continueBtn:                l.Theme.Button("Continue"),
		restoreBtn:                 l.Theme.Button("Restore"),
		importBtn:                  l.Theme.Button("Import"),
		watchOnlyCheckBox:          l.Theme.CheckBox(new(widget.Bool), "Import watch only"),
		walletTypeList:             l.Theme.NewClickableList(layout.Horizontal),
		selectedWalletType:         -1,
		selectedDecredWalletAction: -1,

		Load: l,
	}

	pg.walletName = l.Theme.Editor(new(widget.Editor), values.String(values.StrWalletName))
	pg.walletName.Editor.SingleLine, pg.walletName.Editor.Submit, pg.walletName.IsTitleLabel = true, true, false

	pg.restoreWalletName = l.Theme.Editor(new(widget.Editor), values.String(values.StrWalletName))
	pg.restoreWalletName.Editor.SingleLine, pg.restoreWalletName.Editor.Submit, pg.restoreWalletName.IsTitleLabel = true, true, false

	pg.watchOnlyWalletHex = l.Theme.Editor(new(widget.Editor), "Extended public key")
	pg.restoreWalletName.Editor.SingleLine, pg.watchOnlyWalletHex.Editor.Submit, pg.watchOnlyWalletHex.IsTitleLabel = false, true, false

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *CreateWallet) OnNavigatedTo() {
	pg.initPageItems()
}

func (pg *CreateWallet) initPageItems() {
	walletTypes := []*walletType{
		{
			logo:      pg.Theme.Icons.DecredLogo,
			name:      "Decred",
			clickable: pg.Theme.NewClickable(true),
		},
		{
			logo:      pg.Theme.Icons.BTC,
			name:      "Bitcoin",
			clickable: pg.Theme.NewClickable(true),
		},
	}

	leftRadius := decredmaterial.CornerRadius{
		TopLeft:    8,
		BottomLeft: 8,
	}

	rightRadius := decredmaterial.CornerRadius{
		TopRight:    8,
		BottomRight: 8,
	}

	decredActions := []*decredAction{
		{
			title:     "New Wallet",
			clickable: pg.Theme.NewClickable(true),
			border: decredmaterial.Border{
				Radius: leftRadius,
				Color:  pg.Theme.Color.Gray1,
				Width:  values.MarginPadding2,
			},
			width: values.MarginPadding110,
		},
		{
			title:     "Restore existing wallet",
			clickable: pg.Theme.NewClickable(true),
			border: decredmaterial.Border{Radius: rightRadius,
				Color: pg.Theme.Color.Gray1,
				Width: values.MarginPadding2,
			},
			width: values.MarginPadding195,
		},
	}

	pg.walletTypes = walletTypes
	pg.decredActions = decredActions
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *CreateWallet) OnNavigatedFrom() {}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *CreateWallet) Layout(gtx C) D {
	pageContent := []func(gtx C) D{
		pg.Theme.H6("Select the type of wallet you want to create").Layout,
		pg.walletTypeSection,
		func(gtx C) D {
			switch pg.selectedWalletType {
			case 0:
				return pg.decredWalletOptions(gtx)
			case 1:
				return D{} // todo btc functionality
			default:
				return D{}
			}
		},
		func(gtx C) D {
			switch pg.selectedDecredWalletAction {
			case 0:
				return pg.createNewWallet(gtx)
			case 1:
				return pg.restoreWallet(gtx)
			default:
				return D{}
			}
		},
	}

	return decredmaterial.LinearLayout{
		Width:     decredmaterial.MatchParent,
		Height:    decredmaterial.MatchParent,
		Direction: layout.Center,
	}.Layout2(gtx, func(gtx C) D {
		return decredmaterial.LinearLayout{
			Width:     gtx.Dp(values.MarginPadding377),
			Height:    decredmaterial.MatchParent,
			Alignment: layout.Middle,
			Margin: layout.Inset{
				Top:    values.MarginPadding44,
				Bottom: values.MarginPadding30,
			},
		}.Layout2(gtx, func(gtx C) D {
			return pg.list.Layout(gtx, len(pageContent), func(gtx C, i int) D {
				return layout.Inset{
					Top:    values.MarginPadding26,
					Bottom: values.MarginPadding10,
				}.Layout(gtx, func(gtx C) D {
					return pageContent[i](gtx)
				})
			})
		})
	})
}

// todo bitcoin wallet creation
func (pg *CreateWallet) walletTypeSection(gtx C) D {
	list := layout.List{}
	return list.Layout(gtx, len(pg.walletTypes), func(gtx C, i int) D {
		item := pg.walletTypes[i]
		return decredmaterial.LinearLayout{
			Width:       gtx.Dp(values.MarginPadding172),
			Height:      gtx.Dp(values.MarginPadding174),
			Orientation: layout.Vertical,
			Alignment:   layout.Middle,
			Direction:   layout.Center,
			Background:  pg.Theme.Color.Surface,
			Clickable:   item.clickable,
			Margin: layout.Inset{
				Top:   values.MarginPadding10,
				Right: values.MarginPadding6,
			},
			Padding: layout.UniformInset(values.MarginPadding24),
		}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{
					Bottom: values.MarginPadding14,
				}.Layout(gtx, func(gtx C) D {
					return item.logo.LayoutSize(gtx, values.MarginPadding70)
				})
			}),
			layout.Rigid(pg.Theme.Label(values.TextSize16, item.name).Layout),
		)
	})
}

func (pg *CreateWallet) decredWalletOptions(gtx C) D {
	list := layout.List{}
	return list.Layout(gtx, len(pg.decredActions), func(gtx C, i int) D {
		item := pg.decredActions[i]
		return decredmaterial.LinearLayout{
			Width:       gtx.Dp(item.width),
			Height:      decredmaterial.WrapContent,
			Orientation: layout.Vertical,
			Alignment:   layout.Middle,
			Direction:   layout.Center,
			Background:  pg.Theme.Color.Surface,
			Clickable:   item.clickable,
			Border:      item.border,
			Padding:     layout.UniformInset(values.MarginPadding12),
		}.Layout2(gtx, pg.Theme.Label(values.TextSize16, item.title).Layout)
	})
}

func (pg *CreateWallet) createNewWallet(gtx C) D {
	return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(pg.Theme.Label(values.TextSize16, "What would you like to call your wallet?").Layout),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Top:    values.MarginPadding14,
				Bottom: values.MarginPadding20,
			}.Layout(gtx, func(gtx C) D {
				mGtx := gtx
				if pg.WL.MultiWallet.LoadedWalletsCount() == 0 {
					pg.walletName.Editor.SetText("myWallet")
					mGtx = gtx.Disabled()
				}

				return pg.walletName.Layout(mGtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						return pg.continueBtn.Layout(gtx)
					})
				}),
			)
		}),
	)
}

func (pg *CreateWallet) restoreWallet(gtx C) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(pg.Theme.Label(values.TextSize16, "What is your wallet existing wallet name?").Layout),
		layout.Rigid(pg.watchOnlyCheckBox.Layout),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Top:    values.MarginPadding14,
				Bottom: values.MarginPadding20,
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						mGtx := gtx
						if pg.WL.MultiWallet.LoadedWalletsCount() == 0 {
							pg.restoreWalletName.Editor.SetText("myWallet")
							mGtx = gtx.Disabled()
						}

						return pg.restoreWalletName.Layout(mGtx)
					}),
					layout.Rigid(func(gtx C) D {
						if pg.watchOnlyCheckBox.CheckBox.Value {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Inset{
										Top:    values.MarginPadding10,
										Bottom: values.MarginPadding8,
									}.Layout(gtx, pg.Theme.Label(values.TextSize16, "Extended public key").Layout)
								}),
								layout.Rigid(pg.watchOnlyWalletHex.Layout),
							)
						}
						return D{}
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Flexed(1, func(gtx C) D {
					return layout.E.Layout(gtx, func(gtx C) D {
						if pg.watchOnlyCheckBox.CheckBox.Value {
							return pg.importBtn.Layout(gtx)
						}
						return pg.restoreBtn.Layout(gtx)
					})
				}),
			)
		}),
	)
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *CreateWallet) HandleUserInteractions() {
	for i, item := range pg.walletTypes {
		for item.clickable.Clicked() {
			pg.selectedWalletType = i
		}
	}

	for i, item := range pg.decredActions {
		for item.clickable.Clicked() {
			pg.selectedDecredWalletAction = i
		}
	}

	isSubmit, isChanged := decredmaterial.HandleEditorEvents(pg.walletName.Editor, pg.watchOnlyWalletHex.Editor)
	if isChanged {
		// reset error when any editor is modified
		pg.walletName.SetError("")
	}

	if (pg.continueBtn.Clicked() || isSubmit) && pg.validInputsName() {
		spendingPasswordModal := modal.NewCreatePasswordModal(pg.Load).
			Title("Set spending password").
			NegativeButton(func() {}).
			PasswordCreated(func(_, password string, m *modal.CreatePasswordModal) bool {
				go func() {
					wal, err := pg.WL.MultiWallet.CreateNewWallet(pg.walletName.Editor.Text(), password, dcrlibwallet.PassphraseTypePass)
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
					pg.WL.MultiWallet.SetBoolConfigValueForKey(dcrlibwallet.AccountMixerConfigSet, true)
					m.Dismiss()

					// naviagate to homepage is wallet count is 1 else go back to the previous page
					if pg.WL.MultiWallet.LoadedWalletsCount() == 1 {
						onWalSelected := func() {
							pg.ParentNavigator().ClearStackAndDisplay(NewMainPage(pg.Load))
						}
						pg.ParentNavigator().ClearStackAndDisplay(NewWalletList(pg.Load, onWalSelected))
					} else {
						onWalSelected := func() {
							pg.ParentNavigator().CloseCurrentPage() // todo create new wallet from wallet page
						}
						pg.ParentNavigator().ClearStackAndDisplay(NewWalletList(pg.Load, onWalSelected))
					}
				}()
				return false
			})
		pg.ParentWindow().ShowModal(spendingPasswordModal)
	}

	if pg.restoreBtn.Clicked() {
		afterRestore := func() {
			// naviagate to homepage is wallet count is 1 else go back to the previous page
			if pg.WL.MultiWallet.LoadedWalletsCount() == 1 {
				onWalSelected := func() {
					pg.ParentNavigator().ClearStackAndDisplay(NewMainPage(pg.Load))
				}
				pg.ParentNavigator().ClearStackAndDisplay(NewWalletList(pg.Load, onWalSelected))
			} else {
				onWalSelected := func() {
					pg.ParentNavigator().CloseCurrentPage() // todo create new wallet from wallet page
				}
				pg.ParentNavigator().ClearStackAndDisplay(NewWalletList(pg.Load, onWalSelected))
			}
		}
		pg.ParentNavigator().Display(wallets.NewRestorePage(pg.Load, afterRestore))
	}

	if (pg.importBtn.Clicked() || isSubmit) && pg.validInputsName() {
		go func() {
			_, err := pg.WL.MultiWallet.CreateWatchOnlyWallet(pg.walletName.Editor.Text(), pg.watchOnlyWalletHex.Editor.Text())
			if err != nil {
				pg.watchOnlyWalletHex.SetError(err.Error())
				return
			}

			// naviagate to homepage is wallet count is 1 else go back to the previous page
			if pg.WL.MultiWallet.LoadedWalletsCount() == 1 {
				onWalSelected := func() {
					pg.ParentNavigator().ClearStackAndDisplay(NewMainPage(pg.Load))
				}
				pg.ParentNavigator().ClearStackAndDisplay(NewWalletList(pg.Load, onWalSelected))
			} else {
				onWalSelected := func() {
					pg.ParentNavigator().CloseCurrentPage() // todo create new wallet from wallet page
				}
				pg.ParentNavigator().ClearStackAndDisplay(NewWalletList(pg.Load, onWalSelected))
			}
		}()
	}
}

func (pg *CreateWallet) validInputsName() bool {
	pg.walletName.SetError("")
	pg.watchOnlyWalletHex.SetError("")
	if !components.StringNotEmpty(pg.walletName.Editor.Text()) {
		pg.walletName.SetError(values.String(values.StrEnterWalletName))
		return false
	}

	if pg.watchOnlyCheckBox.CheckBox.Value && !components.StringNotEmpty(pg.watchOnlyWalletHex.Editor.Text()) {
		pg.watchOnlyWalletHex.SetError(values.String(values.StrEnterExtendedPubKey))
		return false
	}

	return true
}

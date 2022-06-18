package page

import (
	// "context"
	"image/color"
	"sync"

	"gioui.org/layout"
	"gioui.org/widget"

	// "github.com/decred/dcrd/dcrutil/v4"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/app"
	// "github.com/planetdecred/godcr/listeners"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const CreateWalletID = "create_wallet"

type collapsibleItem struct {
	title       string
	collapsible *decredmaterial.Collapsible
	body        layout.Widget
}

type walletType struct {
	logo *decredmaterial.Image
	name string
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

	walletCreationOptions []*collapsibleItem
	walletTypes           []*walletType

	walletTypeList *decredmaterial.ClickableList
	walletName     decredmaterial.Editor
	continueBtn    decredmaterial.Button

	wallectSelected    func()
	selectedWalletType int
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

		continueBtn:        l.Theme.Button("Continue"),
		walletTypeList:     l.Theme.NewClickableList(layout.Horizontal),
		selectedWalletType: -1,

		Load: l,
	}

	pg.walletName = l.Theme.Editor(new(widget.Editor), values.String(values.StrWalletName))
	pg.walletName.Editor.SingleLine, pg.walletName.Editor.Submit, pg.walletName.IsTitleLabel = true, true, false

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
	walletCreationOptions := []*collapsibleItem{
		{
			title:       "Select the type of wallet you want to create",
			collapsible: pg.Theme.Collapsible(),
			body:        pg.walletTypeBody,
		},
		{
			title:       "Restore an existing wallet",
			collapsible: pg.Theme.Collapsible(),
			body:        nil,
		},
	}

	walletTypes := []*walletType{
		{
			logo: pg.Theme.Icons.DecredLogo,
			name: "Decred",
		},
		{
			logo: pg.Theme.Icons.BTC,
			name: "Bitcoin",
		},
	}

	pg.walletTypes = walletTypes
	pg.walletCreationOptions = walletCreationOptions
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
		pg.Theme.H6(values.String(values.StrCreateANewWallet)).Layout,
		pg.createWalletSection,
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
				Top:    values.MarginPadding70,
				Bottom: values.MarginPadding70,
			},
		}.Layout2(gtx, func(gtx C) D {
			return pg.list.Layout(gtx, len(pageContent), func(gtx C, i int) D {
				return layout.Inset{Top: values.MarginPadding26}.Layout(gtx, func(gtx C) D {
					return pageContent[i](gtx)
				})
			})
		})
	})
}

func (pg *CreateWallet) createWalletSection(gtx layout.Context) D {

	list := &layout.List{Axis: layout.Vertical}
	return list.Layout(gtx, len(pg.walletCreationOptions), func(gtx C, i int) D {
		collapsibleHeader := pg.Theme.Label(values.TextSize18, pg.walletCreationOptions[i].title).Layout

		collapsibleBody := func(gtx C) D {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, pg.walletCreationOptions[i].body)
		}

		return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
			pg.walletCreationOptions[i].collapsible.IconStyle = decredmaterial.Caret
			pg.walletCreationOptions[i].collapsible.IconPosition = decredmaterial.Before
			pg.Theme.Styles.CollapsibleStyle.Background = color.NRGBA{}
			return pg.walletCreationOptions[i].collapsible.Layout(gtx, collapsibleHeader, collapsibleBody)
		})
	})
}

// todo bitcoin wallet creation
func (pg *CreateWallet) walletTypeBody(gtx C) D {
	return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return pg.walletTypeList.Layout(gtx, len(pg.walletTypes), func(gtx C, i int) D {
				item := pg.walletTypes[i]
				return decredmaterial.LinearLayout{
					Width:       gtx.Dp(values.MarginPadding96),
					Height:      gtx.Dp(values.MarginPadding124),
					Orientation: layout.Vertical,
					Alignment:   layout.Middle,
					Direction:   layout.Center,
					Background:  pg.Theme.Color.Surface,
					Margin: layout.Inset{
						Top:    values.MarginPadding10,
						Bottom: values.MarginPadding30,
						Left:   values.MarginPadding30,
					},
					Padding: layout.UniformInset(values.MarginPadding24),
				}.Layout(gtx,
					layout.Rigid(item.logo.Layout48dp),
					layout.Rigid(pg.Theme.Label(values.TextSize16, item.name).Layout),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			if pg.selectedWalletType != -1 {
				if pg.walletTypes[pg.selectedWalletType].name == pg.walletTypes[0].name {
					return layout.Inset{Left: values.MarginPadding30}.Layout(gtx, func(gtx C) D {
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
					})
				} else if pg.walletTypes[pg.selectedWalletType].name == pg.walletTypes[1].name {
					return pg.Theme.Label(values.TextSize16, "BTC functionality is yet available").Layout(gtx)
				}
			}
			return D{}
		}),
	)
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *CreateWallet) HandleUserInteractions() {
	if ok, selectedItem := pg.walletTypeList.ItemClicked(); ok {
		pg.selectedWalletType = selectedItem
	}

	isSubmit, isChanged := decredmaterial.HandleEditorEvents(pg.walletName.Editor)
	if isChanged {
		// reset error when any editor is modified
		pg.walletName.SetError("")
	}

	if (pg.continueBtn.Clicked() || isSubmit) && pg.validWalletName() {
		spendingPasswordModal := modal.NewCreatePasswordModal(pg.Load).
			Title("Set spending password").
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

					// naviagate to homepage is wallet count is 1
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
}

func (pg *CreateWallet) validWalletName() bool {
	pg.walletName.SetError("")
	if !components.StringNotEmpty(pg.walletName.Editor.Text()) {
		pg.walletName.SetError(values.String(values.StrEnterWalletName))
		return false
	}

	return true
}

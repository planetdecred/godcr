package page

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// LoadingID is the id of the loading page.
const LandingID = "loading"

// Loading represents the loading page of the app.
type Landing struct {
	inset            layout.Inset
	container        layout.List
	heading          material.Label
	createErrorLabel material.Label

	restoreBtn material.Button
	restoreWdg *widget.Button
	createBtn  material.Button
	createWdg  *widget.Button
	walletsBtn material.Button
	walletsWdg *widget.Button

	isCreatingWallet           bool
	isShowingPasswordModal     bool
	walletCreationSuccessEvent interface{}
	passwordModal              *materialplus.Password

	states map[string]interface{}
	theme  *materialplus.Theme
	wal    *wallet.Wallet
}

// Init adds a heading and two buttons.
func (pg *Landing) Init(theme *materialplus.Theme, wal *wallet.Wallet, states map[string]interface{}) {
	pg.heading = theme.Label(units.Label, "Welcome to decred")
	pg.heading.Alignment = text.Middle

	pg.createErrorLabel = theme.Body2("")
	pg.createErrorLabel.Color = theme.Danger

	pg.createBtn = theme.Button("Create Wallet")
	pg.createWdg = new(widget.Button)

	pg.testnetLabel = theme.Label(unit.Dp(16), "testnet")
	pg.testnetLabel.Font.Weight = text.Bold
	pg.testnetLabel.Font.Size = unit.Px(50)
	pg.testnetLabel.Alignment = text.Middle
	pg.testnetLabel.Font.Weight = text.Bold

	pg.walletsBtn = theme.Button("Back to Wallets")
	pg.walletsWdg = new(widget.Button)

	pg.inset = layout.UniformInset(units.FlexInset)
	pg.container = layout.List{Axis: layout.Vertical}
	pg.isCreatingWallet = false
	pg.isShowingPasswordModal = false
	pg.walletCreationSuccessEvent = nil
	pg.passwordModal = theme.Password()
	pg.states = states
	pg.theme = theme
	pg.wal = wal
}

// Draw draws the page's to the given layout context.
// Does not react to any event but can return a Nav event.
func (pg *Landing) Draw(gtx *layout.Context) interface{} {
	ev := pg.walletCreationSuccessEvent
	if pg.walletCreationSuccessEvent != nil {
		pg.walletCreationSuccessEvent = nil
	}

	pg.checkForStatesUpdate()

	walletInfo := pg.states[StateWalletInfo].(*wallet.MultiWalletInfo)
	widgets := []func(){
		func() {
			gtx.Dimensions.Size.Y = 264
		},
		func() {
			topInset := float32(0)

			if pg.createErrorLabel.Text != "" {
				pg.createErrorLabel.Layout(gtx)
				topInset += 20
			}

			inset := layout.Inset{
				Top: unit.Dp(topInset),
			}
			inset.Layout(gtx, func() {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				for pg.createWdg.Clicked(gtx) {
					if !pg.isCreatingWallet {
						pg.isShowingPasswordModal = !pg.isShowingPasswordModal
					}
				}
				pg.createBtn.Layout(gtx, pg.createWdg)
			})
		},
		func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			if pg.restoreWdg.Clicked(gtx) {
				log.Debugf("{%s} Restore Btn clicked", LandingID)
				// res = EventNav {
				// 	Current: LandingID,
				// 	Next: CreateID,
				// }
			}
			pg.restoreBtn.Layout(gtx, pg.restoreWdg)
		},
		func() {
			if walletInfo.LoadedWallets != 0 {
				gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
				for pg.walletsWdg.Clicked(gtx) {
					if !pg.isCreatingWallet {
						pg.reset()
						ev = EventNav{
							Current: LandingID,
							Next:    WalletsID,
						}
					}
				}
				pg.walletsBtn.Layout(gtx, pg.walletsWdg)

			}
		},
	}

	pg.container.Layout(gtx, len(widgets),
		layout.ListElement(func(i int) {
			layout.UniformInset(units.FlexInset).Layout(gtx, widgets[i])
		}),
	)

	if pg.isShowingPasswordModal {
		pg.drawPasswordModal(gtx)
	}

	return ev
}

func (pg *Landing) checkForStatesUpdate() {
	err := pg.states[StateError]
	created := pg.states[StateWalletCreated]

	if err == nil && created == nil {
		return
	}

	pg.reset()

	if created != nil {
		pg.passwordModal.Reset()
		delete(pg.states, StateWalletCreated)

		pg.walletCreationSuccessEvent = EventNav{
			Current: LandingID,
			Next:    WalletsID,
		}
	} else if err != nil {
		pg.createErrorLabel.Text = err.(error).Error()
		delete(pg.states, StateError)
	}
}

func (pg *Landing) drawPasswordModal(gtx *layout.Context) {
	pg.theme.Modal(gtx, func() {
		pg.passwordModal.Draw(gtx, pg.confirm, pg.cancel)
	})
}

func (pg *Landing) confirm(password string) {
	pg.reset()

	pg.createBtn.Text = "Creating wallet..."
	pg.createBtn.Background = pg.theme.Disabled

	pg.isCreatingWallet = true
	pg.isShowingPasswordModal = false

	pg.wal.CreateWallet(password)
}

func (pg *Landing) cancel() {
	pg.isShowingPasswordModal = false
}

func (pg *Landing) reset() {
	pg.isCreatingWallet = false
	pg.createBtn.Text = "Create wallet"
	pg.createBtn.Background = pg.theme.Color.Primary
	pg.createErrorLabel.Text = ""
}

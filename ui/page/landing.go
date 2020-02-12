package page

import (
	"image/color"
	"image/png"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/markbates/pkger"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
<<<<<<< HEAD
	"github.com/raedahgroup/godcr-gio/ui/units"
=======
<<<<<<< HEAD
=======
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/ui/widgets"
>>>>>>> add back button
>>>>>>> add back button
	"github.com/raedahgroup/godcr-gio/wallet"
)

// LandingID is the id of the landing page.
const LandingID = "landing"

// Landing represents the landing page of the app.
// It should only be should shown if the app launches
// and cannot find any wallets.
type Landing struct {
<<<<<<< HEAD
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
=======
	container    *layout.List
	image        material.Image
	welcomeLabel material.Label

	createButton  *widget.Button
	restoreButton *widget.Button

	addIcon     material.Image
	restoreIcon material.Image
	backIcon    material.Image
	walletExits bool
}

// Init adds a heading and two buttons.
func (pg *Landing) Init(theme *materialplus.Theme, wallet *wallet.Wallet, states map[string]interface{}) {
	file, err := pkger.Open("/assets/icons/decred.png")
	if err != nil {
		log.Println(err)
	}
	image, err := png.Decode(file)
	if err != nil {
		log.Println(err)
	}
	pg.image = theme.Image(paint.NewImageOp(image))
>>>>>>> add back button

	pg.createErrorLabel = theme.Body2("")
	pg.createErrorLabel.Color = theme.Danger

	pg.createBtn = theme.Button("Create Wallet")
	pg.createWdg = new(widget.Button)

	file, err = pkger.Open("/assets/icons/add.png")
	if err != nil {
		log.Println(err)
	}
	image, err = png.Decode(file)
	if err != nil {
		log.Println(err)
	}
	pg.addIcon = theme.Image(paint.NewImageOp(image))

<<<<<<< HEAD
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
=======
	file, err = pkger.Open("/assets/icons/restore.png")
	if err != nil {
		log.Println(err)
	}
	image, err = png.Decode(file)
	if err != nil {
		log.Println(err)
	}
	pg.restoreIcon = theme.Image(paint.NewImageOp(image))

	file, err = pkger.Open("/assets/icons/back.png")
	if err != nil {
		log.Println(err)
	}
	image, err = png.Decode(file)
	if err != nil {
		log.Println(err)
	}
	pg.backIcon = theme.Image(paint.NewImageOp(image))
	pg.container = &layout.List{
		Axis: layout.Vertical,
	}
>>>>>>> add back button
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
			imgBt := widgets.NewImageButton(&pg.backIcon, "")
			imgBt.Background = color.RGBA{255, 255, 255, 255}
			imgBt.Src.Scale = 0.5
			imgBt.HPadding = unit.Dp(0)
			imgBt.Layout(gtx, pg.createButton, 20)
		},

		func() {
			layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func() {})
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

<<<<<<< HEAD
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
=======
	pg.container.Layout(gtx, len(widgets), func(i int) {
		layout.Inset{Tozp: unit.Dp(8), Left: unit.Dp(24), Right: unit.Dp(24), Bottom: unit.Dp(8)}.Layout(gtx, widgets[i])
>>>>>>> add back button
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

package page

import (
	"image/color"
	"image/png"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/markbates/pkger"
	"github.com/raedahgroup/godcr-gio/ui"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// LandingID is the id of the landing page.
const LandingID = "landing"

// Landing represents the landing page of the app.
// It's shown when users are to create or restore a wallet.
type Landing struct {
	container    *layout.List
	image        material.Image
	welcomeLabel material.Label

	createButton  *widget.Button
	restoreButton *widget.Button
	backButton    *widget.Button

	addIcon     material.Image
	restoreIcon material.Image
	backIcon    material.Image

	states map[string]interface{}
}

// Init adds a heading and two buttons.
func (pg *Landing) Init(theme *materialplus.Theme, wallet *wallet.Wallet, states map[string]interface{}) {
	file, err := pkger.Open("/assets/icons/decred.png")
	if err != nil {
		log.Error(err)
	}
	image, err := png.Decode(file)
	if err != nil {
		log.Error(err)
	}
	pg.image = theme.Image(paint.NewImageOp(image))

	pg.welcomeLabel = theme.Label(units.Label, "Welcome to\nDecred desktop wallet")

	pg.createButton = new(widget.Button)
	pg.restoreButton = new(widget.Button)
	pg.backButton = new(widget.Button)

	file, err = pkger.Open("/assets/icons/add.png")
	if err != nil {
		log.Error(err)
	}
	image, err = png.Decode(file)
	if err != nil {
		log.Error(err)
	}
	pg.addIcon = theme.Image(paint.NewImageOp(image))

	file, err = pkger.Open("/assets/icons/restore.png")
	if err != nil {
		log.Error(err)
	}
	image, err = png.Decode(file)
	if err != nil {
		log.Error(err)
	}
	pg.restoreIcon = theme.Image(paint.NewImageOp(image))

	file, err = pkger.Open("/assets/icons/back.png")
	if err != nil {
		log.Error(err)
	}
	image, err = png.Decode(file)
	if err != nil {
		log.Error(err)
	}
	pg.backIcon = theme.Image(paint.NewImageOp(image))
	pg.container = &layout.List{
		Axis: layout.Vertical,
	}

	pg.states = states
}

// Draw draws the page's to the given layout context.
// Does not react to any event but can return a Nav event.
func (pg *Landing) Draw(gtx *layout.Context) (res interface{}) {
	widgets := []func(){
		func() {
			imgBt := materialplus.NewImageButton(&pg.backIcon, "")
			imgBt.Background = color.RGBA{255, 255, 255, 255}
			imgBt.Src.Scale = 0.5
			imgBt.HPadding = unit.Dp(0)

			imgBt.Layout(gtx, pg.backButton, 20)
			for pg.backButton.Clicked(gtx) {
				ev := EventNav{
					Current: LandingID,
					Next:    WalletsID,
				}
				res = ev
			}
		},

		func() {
			layout.Inset{Bottom: unit.Dp(6)}.Layout(gtx, func() {})
		},

		func() {
			pg.image.Scale = 0.5
			pg.image.Layout(gtx)
		},

		func() {
			layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func() {})
		},

		func() {
			pg.welcomeLabel.Layout(gtx)
		},

		func() {
			layout.Inset{Bottom: unit.Dp(270)}.Layout(gtx, func() {})
		},

		func() {
			imgBt := materialplus.NewImageButton(&pg.addIcon, "Create a new wallet")
			imgBt.Background = ui.LightBlueColor
			imgBt.VPadding = unit.Dp(20)
			imgBt.Src.Scale = 0.3
			imgBt.Font.Size = units.SmallText

			imgBt.Layout(gtx, pg.createButton, 20)
		},

		func() {
			imgBt := materialplus.NewImageButton(&pg.restoreIcon, "Restore an existing wallet")
			imgBt.Background = ui.LighGreenColor
			imgBt.VPadding = unit.Dp(20)
			imgBt.Src.Scale = 0.3
			imgBt.Font.Size = units.SmallText

			imgBt.Layout(gtx, pg.restoreButton, 20)
		},

		func() {
			layout.Inset{Bottom: unit.Dp(20)}.Layout(gtx, func() {})
		},
	}

	walletInfo := pg.states[StateWalletInfo].(*wallet.MultiWalletInfo)
	if walletInfo.LoadedWallets == 0 {
		widgets[0] = func() {
			layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, func() {})
		}
		widgets[5] = func() {
			layout.Inset{Bottom: unit.Dp(310)}.Layout(gtx, func() {})
		}
	}

	pg.container.Layout(gtx, len(widgets), func(i int) {
		layout.Inset{Top: unit.Dp(8), Left: unit.Dp(24), Right: unit.Dp(24), Bottom: unit.Dp(8)}.Layout(gtx, widgets[i])
	})

	return res
}

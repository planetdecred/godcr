package page

import (
	"image/color"
	"image/png"
	"log"

	"github.com/markbates/pkger"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/widgets"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// LandingID is the id of the landing page.
const CreateAndRestoreLandingID = "Create"

// Landing represents the landing page of the app.
// It should only be should shown if the app launches
// and cannot find any wallets.
type CreateAndRestoreLanding struct {
	container    *layout.List
	image        material.Image
	welcomeLabel material.Label

	createButton  *widget.Button
	restoreButton *widget.Button

	addIcon     material.Image
	restoreIcon material.Image
}

// Init adds a heading and two buttons.
func (pg *CreateAndRestoreLanding) Init(theme *materialplus.Theme, _ *wallet.Wallet, states map[string]interface{}) {
	file, err := pkger.Open("/assets/icons/decred.png")
	if err != nil {
		log.Println(err)
	}
	image, err := png.Decode(file)
	pg.image = theme.Image(paint.NewImageOp(image))

	pg.welcomeLabel = theme.Label(unit.Sp(24), "Welcome to\nDecred desktop wallet")

	pg.createButton = new(widget.Button)
	pg.restoreButton = new(widget.Button)

	file, err = pkger.Open("/assets/icons/add.png")
	if err != nil {
		log.Println(err)
	}
	image, err = png.Decode(file)
	if err != nil {
		log.Println(err)
	}
	pg.addIcon = theme.Image(paint.NewImageOp(image))

	file, err = pkger.Open("/assets/icons/restore.png")
	if err != nil {
		log.Println(err)
	}
	image, err = png.Decode(file)
	if err != nil {
		log.Println(err)
	}
	pg.restoreIcon = theme.Image(paint.NewImageOp(image))

	pg.container = &layout.List{
		Axis: layout.Vertical,
	}
}

// Draw draws the page's to the given layout context.
// Does not react to any event but can return a Nav event.
func (pg *CreateAndRestoreLanding) Draw(gtx *layout.Context) (res interface{}) {
	widgets := []func(){
		func() {
			//bb.Layout(gtx, unit.Dp(32))
			//	th.IconButton(bb).Layout(gtx, button)
			gtx.Dimensions.Size.Y = 24
		},

		func() {
			pg.image.Scale = 0.5
			pg.image.Layout(gtx)
		},

		func() {
			gtx.Dimensions.Size.Y = 24
		},

		func() {
			pg.welcomeLabel.Layout(gtx)
		},

		func() {
			gtx.Dimensions.Size.Y = 550
		},

		func() {
			imgBt := widgets.NewImageButton(&pg.addIcon, "Create a new wallet")
			imgBt.Axis = layout.Horizontal
			imgBt.Background = color.RGBA{41, 112, 255, 255}
			imgBt.VPadding = unit.Dp(20)

			imgBt.Layout(gtx, pg.createButton, 20)
		},

		func() {
			imgBt := widgets.NewImageButton(&pg.restoreIcon, "Restore an existing wallet")
			imgBt.Axis = layout.Horizontal
			imgBt.Background = color.RGBA{45, 216, 163, 255}
			imgBt.VPadding = unit.Dp(20)

			imgBt.Layout(gtx, pg.restoreButton, 20)
		},
	}

	pg.container.Layout(gtx, len(widgets), func(i int) {
		layout.Inset{Top: unit.Dp(8), Left: unit.Dp(24), Right: unit.Dp(24), Bottom: unit.Dp(8)}.Layout(gtx, widgets[i])
	})

	return nil
}

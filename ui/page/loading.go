package page

import (
	"image/png"
	"log"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/markbates/pkger"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// LoadingID is the id of the loading page.
const LoadingID = "loading"

// Loading represents the loading page of the app.
type Loading struct {
	container    *layout.List
	loadingLabel material.Label
	icon         material.Image
	testnetLabel material.Label
}

// Init initializies the page with a label.
func (pg *Loading) Init(theme *materialplus.Theme, _ *wallet.Wallet, states map[string]interface{}) {
	pg.container = &layout.List{
		Axis: layout.Vertical,
	}

	pg.loadingLabel = theme.Label(unit.Dp(16), "Loading Wallets....")
	pg.loadingLabel.Font.Size = unit.Px(48)
	pg.loadingLabel.Alignment = text.Middle
	pg.loadingLabel.Font.Weight = text.Bold

	pg.testnetLabel = theme.Label(unit.Dp(16), "testnet")
	pg.testnetLabel.Font.Weight = text.Bold
	pg.testnetLabel.Font.Size = unit.Px(50)
	pg.testnetLabel.Alignment = text.Middle
	pg.testnetLabel.Font.Weight = text.Bold

	file, err := pkger.Open("/assets/icons/decred-loader.png")
	if err != nil {
		log.Println(err)
	}
	image, err := png.Decode(file)
	if err != nil {
		log.Println(err)
	}
	pg.icon = theme.Image(paint.NewImageOp(image))
}

// Draw renders the page widgets.
// It does not react to nor does it generate any event.
func (pg *Loading) Draw(gtx *layout.Context) (res interface{}) {
	widgets := []func(){
		func() {
			layout.Inset{Bottom: unit.Dp(180)}.Layout(gtx, func() {})
		},
		func() {
			pg.icon.Scale = 0.08
			pg.icon.Layout(gtx)
		},
		func() {
			layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, func() {})
		},
		func() {
			pg.testnetLabel.Layout(gtx)
		},
		func() {
			gtx.Dimensions.Size.Y = 20
		},
		func() {
			pg.loadingLabel.Layout(gtx)
		},
	}

	pg.container.Layout(gtx, len(widgets), func(i int) {
		layout.Align(layout.Center).Layout(gtx, widgets[i])
	})

	return nil
}

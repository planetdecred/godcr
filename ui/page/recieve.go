package page

import (
	"gioui.org/layout"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/event"
	"github.com/raedahgroup/godcr-gio/ui/units"
)

// ReceivingID is the id of the receiving page.
const ReceivingID = "receiving"

// Receiving represents the receiving page of the app.
type Receiving struct {
<<<<<<< HEAD
	label material.Label
=======
	titleLabel material.Label
	inset   layout.Inset
	stack   layout.Stack
	theme 	*material.Theme
	image 	*material.Image
	container layout.List
	pageTitle string
	walletName string
	walletIndex string
	walletAmount string
>>>>>>> cleaned up basic interface structure, working on acount selector
}

// Init initializies the page with a label.
func (pg *Receiving) Init(theme *material.Theme) {
<<<<<<< HEAD
	pg.label = theme.Label(units.Label, "Receive DCR")
	pg.label = theme.Label(units.Label, "Receive DCR")
	pg.label = theme.Label(units.Label, "Receive DCR")
	pg.label = theme.Label(units.Label, "Receive DCR")

=======
	pg.pageTitle = "Receive DCR"
	pg.inset = layout.UniformInset(unit.Dp(5))
	pg.walletName = "Default"
	pg.walletIndex = "Wallet-1"
	pg.walletAmount = "100.2345 DCR"
	pg.theme = theme
>>>>>>> cleaned up basic interface structure, working on acount selector
}

// Draw renders the page widgets.
// It does not react to nor does it generate any event.
<<<<<<< HEAD
func (pg *Receiving) Draw(gtx *layout.Context, _ event.Event) event.Event {
	pg.label.Layout(gtx)
	pg.label.Layout(gtx)
	return nil
=======
func (pg *Receiving) Draw(gtx *layout.Context, _ event.Event) (evt event.Event) {
	t := pg.theme
	layout.Stack{Alignment: layout.SE}.Layout(gtx,
		layout.Expanded(func() {
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func() {
					layout.UniformInset(unit.Dp(16)).Layout(gtx, func() {
						sz := gtx.Px(unit.Dp(32))
						cs := gtx.Constraints
						gtx.Constraints = layout.RigidConstraints(cs.Constrain(image.Point{X: sz, Y: sz}))
						t.Label(unit.Dp(22), pg.pageTitle).Layout(gtx)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Left: unit.Dp(60), Bottom: unit.Dp(30)}.Layout(gtx, func() {
						layout.Flex{}.Layout(gtx,
							layout.Rigid(func() {
								layout.Inset{Top: unit.Dp(6.5), Right: unit.Dp(15)}.Layout(gtx, func() {
									t.Label(unit.Dp(14), " *_+_* ").Layout(gtx)
								})
							}),
							layout.Rigid(func() {
								layout.Inset{Right: unit.Dp(50)}.Layout(gtx, func() {
									layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(func() {
											layout.Inset{Bottom: unit.Dp(5)}.Layout(gtx, func() {
												t.Label(unit.Dp(16), pg.walletName).Layout(gtx)
											})
										}),
										layout.Rigid(func() {
											layout.Inset{Left: unit.Dp(2)}.Layout(gtx, func() {
												t.Label(unit.Dp(10), pg.walletIndex).Layout(gtx)
											})
										}),
									)
								})
							}),
							layout.Rigid(func() {
								layout.Inset{Top: unit.Dp(6.5)}.Layout(gtx, func() {
									t.Label(unit.Dp(16), pg.walletAmount).Layout(gtx)
								})
							}),
						)
					})						
				}),
				layout.Rigid(func() {
					layout.Inset{Left: unit.Dp(70)}.Layout(gtx, func() {
						generateAddressAndQR(gtx, false, t)
					})
				}),
				layout.Flexed(1, func() {
					layout.UniformInset(unit.Dp(16)).Layout(gtx, func() {
						layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func() {
								label := t.Label(unit.Dp(16), "gvsdfvadfertt45656ynghdffgdfgb")
								label.Color = color.RGBA{
									R: 44,
									G: 114,
									B: 255,
									A: 255,
								}
								layout.Inset{Left: unit.Dp(80), Top: unit.Dp(20), Bottom: unit.Dp(10)}.Layout(gtx, func() {
									label.Layout(gtx)
								})
							}),
							layout.Rigid(func() {
								layout.Inset{Left: unit.Dp(150)}.Layout(gtx, func() {
									t.Label(unit.Dp(16), "(tap to copy)").Layout(gtx)
								})
							}),
						)
					})
				}),
			)
		}),
	)
	return 
}

func generateAddressAndQR(ctx *layout.Context, newAddress bool, t *material.Theme) {	
	addr := "Drnjjvsvnfkjjkvioer98493940940mb"
	
	qrCode, err := qrcode.New(addr, qrcode.Highest)
	qrCode.DisableBorder = true
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	img := t.Image(paint.NewImageOp(qrCode.Image(150)))

	material.Image(img).Layout(ctx)
}

func rgb(c uint32) color.RGBA {
	return argb((0xff << 24) | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

type fill struct {
	col color.RGBA
}

func (f fill) Layout(gtx *layout.Context) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: f.col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
>>>>>>> cleaned up basic interface structure, working on acount selector
}

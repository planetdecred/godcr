package page

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/ui/units"
	"github.com/raedahgroup/godcr-gio/wallet"
)

// RestoreID is the id of the restore page.
const RestoreID = "restore"

// Restore represents the restore page of the app.
type Restore struct {
	inset     layout.Inset
	container layout.Flex
	heading   material.Label

	restoreBtn material.Button
	restoreWdg *widget.Button

	inputGroupContainer layout.List
	inputGroupHeader    material.Label
	inputLabels         []material.Label
	inputs              []material.Editor
	editors             []*widget.Editor
}

// Init initializes restore page with heading, 33 inputs and button
func (pg *Restore) Init(theme *materialplus.Theme, _ *wallet.Wallet) {
	pg.heading = theme.Label(units.Label, "Restore from seed phrase")
	pg.heading.Alignment = text.Middle

	pg.inset = layout.UniformInset(unit.Dp(10))
	pg.container.Axis = layout.Vertical

	pg.inputGroupHeader = theme.Label(unit.Dp(16), "Enter your seed phrase in the correct order")
	pg.inputGroupHeader.Alignment = text.Middle
	pg.inputGroupHeader.Color = color.RGBA{102, 102, 102, 255}
	pg.inputGroupContainer.Axis = layout.Vertical

	pg.restoreBtn = theme.Button("Restore")
	pg.restoreBtn.Background = color.RGBA{196, 203, 210, 255}
	pg.restoreWdg = new(widget.Button)

	for i := 0; i <= 32; i++ {
		pg.inputs = append(pg.inputs, theme.Editor("Input phrase "+strconv.Itoa(i+1)+"..."))
		pg.inputs[i].Font.Size = unit.Sp(16)
		// pg.editors[i] = new(widget.Editor)
		pg.editors = append(pg.editors, &widget.Editor{SingleLine: true})
		pg.inputLabels = append(pg.inputLabels, theme.Label(unit.Dp(13), strconv.Itoa(i+1)))
	}
}

// Draw renders the page widgets
func (pg *Restore) Draw(gtx *layout.Context, _ ...interface{}) interface{} {
	pg.container.Layout(gtx,
		layout.Rigid(func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			pg.heading.Layout(gtx)
		}),
		layout.Rigid(func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			pg.inputGroupHeader.Layout(gtx)
		}),
		layout.Flexed(1, func() {
			pg.inputGroupContainer.Layout(gtx, len(pg.inputs),
				layout.ListElement(func(i int) {
					layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func() {
							in := layout.Inset{Left: unit.Dp(20)}
							in.Layout(gtx, func() {
								dim := gtx.Px(unit.Dp(22))
								sz := image.Point{X: dim, Y: dim}
								gtx.Constraints = layout.RigidConstraints(gtx.Constraints.Constrain(sz))
								pg.inputLabels[i].Layout(gtx)
							})
						}),
						layout.Flexed(1, func() {
							gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
							in := layout.Inset{Bottom: unit.Dp(20), Left: unit.Dp(20), Right: unit.Dp(20)}
							in.Layout(gtx, func() {
								pg.inputs[i].Layout(gtx, pg.editors[i])
							})
						}),
					)
				}),
			)
		}),
		layout.Rigid(func() {
			inset := layout.Inset{Bottom: unit.Dp(0)}
			gtx.Constraints.Height.Min = 44
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			inset.Layout(gtx, func() {
				pg.restoreBtn.Layout(gtx, pg.restoreWdg)
				if pg.restoreWdg.Clicked(gtx) {
					fmt.Println("ButtonClicked #15")
					txt := pg.editors[1].Text()
					fmt.Println(txt)
				}
			})
		}),
	)
	return nil
}

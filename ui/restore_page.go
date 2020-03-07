package ui

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
)

const (
	suggestionItems = 4   // Maximum of suggestionItems
	buttonWidth     = 210 // Width of the buttons
)

var (
	inputGroupContainerLeft  = *&layout.List{Axis: layout.Vertical}
	inputGroupContainerRight = *&layout.List{Axis: layout.Vertical}
)

// RestorePage lays out the main wallet page
func (win *Window) RestorePage() {
	body := func() {
		layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
			layout.Rigid(func() {
				txt := win.theme.H3("Restore from seed phrase")
				txt.Alignment = text.Middle
				txt.Layout(win.gtx)
			}),
			layout.Rigid(func() {
				txt := win.theme.H6("Enter your seed phrase in the correct order")
				txt.Alignment = text.Middle
				txt.Layout(win.gtx)
			}),
			layout.Rigid(func() {
				layout.Inset{Top: unit.Dp(20)}.Layout(win.gtx, func() {})
			}),
			layout.Flexed(1, func() {
				layout.Center.Layout(win.gtx, func() {
					layout.Flex{}.Layout(win.gtx,
						layout.Rigid(func() {
							drawInputGroup(win, &inputGroupContainerLeft, 16, 0)
						}),
						layout.Rigid(func() {
							drawInputGroup(win, &inputGroupContainerRight, 17, 16)
						}),
					)
				})
			}),
			layout.Rigid(func() {
				layout.Center.Layout(win.gtx, func() {
					win.gtx.Constraints.Width.Min = buttonWidth
					layout.Inset{Top: unit.Dp(15), Bottom: unit.Dp(15)}.Layout(win.gtx, func() {
						win.outputs.restoreDiag.Layout(win.gtx, &win.inputs.restoreDiag)
					})
				})
			}),
		)
	}

	win.Page(body)
}

// drawInputGroup lays out the list vertically with each row align input and label horizontally
func drawInputGroup(win *Window, l *layout.List, len int, startIndex int) {
	win.gtx.Constraints.Width.Min = win.gtx.Constraints.Width.Max / 2
	l.Layout(win.gtx, len, func(i int) {
		layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
			layout.Rigid(func() {
				layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(win.gtx,
					layout.Rigid(func() {
						win.theme.Label(unit.Dp(16), fmt.Sprintf("Word #%d", i+startIndex+1)).Layout(win.gtx)
					}),
					layout.Rigid(func() {
						layout.Inset{Left: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(win.gtx, func() {
							win.outputs.seeds[i+startIndex].Layout(win.gtx, &win.inputs.seeds[i+startIndex])
						})
						// pg.editorEventsHandler(gtx, i+startIndex)
					}),
				)
			}),
			// layout.Rigid(func() {
			// 	pg.drawAutoComplete(gtx, i+startIndex)
			// }),
		)
	})
}

package ui

// import (
// 	"gioui.org/layout"
// 	"gioui.org/text"
// 	"gioui.org/unit"
// )

// func (win *Window) verifyMessageDiag() {
// 	win.theme.Surface(win.gtx, func() {
// 		win.gtx.Constraints.Width.Min = win.gtx.Px(unit.Dp(550))
// 		win.gtx.Constraints.Width.Max = win.gtx.Constraints.Width.Min
// 		layout.UniformInset(unit.Dp(20)).Layout(win.gtx, func() {
// 			win.vFlex(
// 				rigid(func() {
// 					win.hFlex(
// 						rigid(func() {
// 							win.theme.H5("Verify Wallet Message").Layout(win.gtx)
// 						}),
// 						layout.Flexed(.7, func() {
// 							layout.E.Layout(win.gtx, func() {
// 								win.outputs.verifyInfo.Layout(win.gtx, &win.inputs.verifyInfo)
// 							})
// 						}),
// 						layout.Flexed(1, func() {
// 							layout.E.Layout(win.gtx, func() {
// 								win.outputs.cancelDiag.Layout(win.gtx, &win.cancelDialog)
// 							})
// 						}),
// 					)
// 				}),
// 				rigid(func() {
// 					win.Err()
// 				}),
// 				rigid(func() {
// 					win.vFlexSB(
// 						rigid(func() {
// 							inset := layout.Inset{
// 								Top:    unit.Dp(10),
// 								Bottom: unit.Dp(10),
// 							}
// 							inset.Layout(win.gtx, func() {
// 								win.outputs.addressInput.Layout(win.gtx)
// 							})
// 						}),
// 						rigid(func() {
// 							win.outputs.signInput.Layout(win.gtx)
// 						}),
// 						rigid(func() {
// 							inset := layout.Inset{
// 								Top:    unit.Dp(10),
// 								Bottom: unit.Dp(20),
// 							}
// 							inset.Layout(win.gtx, func() {
// 								win.outputs.messageInput.Layout(win.gtx)
// 							})
// 						}),
// 						rigid(func() {
// 							layout.Flex{}.Layout(win.gtx,
// 								layout.Flexed(.6, func() {
// 									win.outputs.verifyMessage.Layout(win.gtx)
// 								}),
// 								layout.Flexed(.4, func() {
// 									layout.Flex{}.Layout(win.gtx,
// 										layout.Flexed(.5, func() {
// 											layout.Inset{Left: unit.Dp(0), Right: unit.Dp(10)}.Layout(win.gtx, func() {
// 												win.outputs.clearBtn.Layout(win.gtx, &win.inputs.clearBtn)
// 											})
// 										}),
// 										layout.Flexed(.5, func() {
// 											win.outputs.verifyBtn.Layout(win.gtx, &win.inputs.verifyBtn)
// 										}),
// 									)
// 								}),
// 							)
// 						}),
// 					)
// 				}),
// 			)
// 		})
// 	})
// }

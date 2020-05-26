package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

func (win *Window) AddAccountDiag() {
	win.theme.Surface(win.gtx, func() {
		toMax(win.gtx)
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(win.gtx,
			layout.Rigid(func() {
				layout.E.Layout(win.gtx, func() {
					win.outputs.cancelDiag.Layout(win.gtx, &win.cancelDialog)
				})
			}),
			layout.Rigid(func() {
				d := win.theme.H3("Add account ")
				d.Layout(win.gtx)
			}),
			layout.Rigid(func() {
				win.outputs.dialog.Layout(win.gtx, &win.inputs.dialog)
			}),
			layout.Rigid(func() {
				win.outputs.spendingPassword.Layout(win.gtx, &win.inputs.spendingPassword)
			}),
			layout.Rigid(func() {
				win.Err()
			}),
			layout.Rigid(func() {
				win.outputs.addAccount.Layout(win.gtx, &win.inputs.addAccount)
			}),
		)
	})
}

func (win *Window) editPasswordDiag() {
	win.theme.Surface(win.gtx, func() {
		win.gtx.Constraints.Width.Min = win.gtx.Px(unit.Dp(450))
		win.gtx.Constraints.Width.Max = win.gtx.Constraints.Width.Min
		layout.UniformInset(unit.Dp(20)).Layout(win.gtx, func() {
			win.vFlex(
				rigid(func() {
					win.hFlex(
						rigid(func() {
							win.theme.H5("Change Wallet Password").Layout(win.gtx)
						}),
						layout.Flexed(1, func() {
							layout.E.Layout(win.gtx, func() {
								win.outputs.cancelDiag.Layout(win.gtx, &win.cancelDialog)
							})
						}),
					)
				}),
				rigid(func() {
					win.Err()
				}),
				rigid(func() {
					win.vFlexSB(
						rigid(func() {
							inset := layout.Inset{
								Top: unit.Dp(10),
							}
							inset.Layout(win.gtx, func() {
								win.vFlexSB(
									rigid(func() {
										win.theme.Body1("Old Password").Layout(win.gtx)
									}),
									rigid(func() {
										win.hFlexSB(
											layout.Flexed(1, func() {
												win.outputs.oldSpendingPassword.Layout(win.gtx, &win.inputs.oldSpendingPassword)
											}),
										)
									}),
								)
							})
						}),
						rigid(func() {
							win.vFlexSB(
								rigid(func() {
									win.passwordStrength()
								}),
								rigid(func() {
									win.theme.Body1("New Password").Layout(win.gtx)
								}),
								rigid(func() {
									win.hFlexSB(
										layout.Flexed(1, func() {
											win.outputs.spendingPassword.Layout(win.gtx, &win.inputs.spendingPassword)
										}),
									)
								}),
							)
						}),
						rigid(func() {
							win.theme.Body1("Confirm New Password").Layout(win.gtx)
						}),
						rigid(func() {
							win.hFlexSB(
								layout.Flexed(1, func() {
									win.outputs.matchSpending.Layout(win.gtx, &win.inputs.matchSpending)
								}),
							)
						}),
						rigid(func() {
							layout.Inset{Top: unit.Dp(20), Bottom: unit.Dp(15)}.Layout(win.gtx, func() {
								win.outputs.savePassword.Layout(win.gtx, &win.inputs.savePassword)
							})
						}),
					)
				}),
			)
		})
	})
}

func (win *Window) passwordStrength() {
	layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(8)}.Layout(win.gtx, func() {
		win.gtx.Constraints.Height.Max = 20
		win.outputs.passwordBar.Layout(win.gtx)
	})
}

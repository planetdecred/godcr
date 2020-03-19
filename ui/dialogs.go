package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
)

func (win *Window) CreateDiag() {
	win.theme.Surface(win.gtx, func() {
		toMax(win.gtx)
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(win.gtx,
			layout.Rigid(func() {
				layout.E.Layout(win.gtx, func() {
					win.outputs.cancelDiag.Layout(win.gtx, &win.cancelDialog)
				})
			}),
			layout.Rigid(func() {
				d := win.theme.H3("Create Wallet")
				d.Layout(win.gtx)
			}),
			layout.Rigid(func() {
				win.outputs.spendingPassword.Layout(win.gtx, &win.inputs.spendingPassword)
			}),
			layout.Rigid(func() {
				win.outputs.matchSpending.Layout(win.gtx, &win.inputs.matchSpending)
			}),
			layout.Rigid(func() {
				win.Err()
			}),
			layout.Rigid(func() {
				win.outputs.createWallet.Layout(win.gtx, &win.inputs.createWallet)
			}),
		)
	})
}

func (win *Window) DeleteDiag() {
	win.theme.Surface(win.gtx, func() {
		toMax(win.gtx)
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(win.gtx,
			layout.Rigid(func() {
				layout.E.Layout(win.gtx, func() {
					win.outputs.cancelDiag.Layout(win.gtx, &win.cancelDialog)
				})
			}),
			layout.Rigid(func() {
				d := win.theme.H3("Delete wallet")
				d.Layout(win.gtx)
			}),
			layout.Rigid(func() {
				win.outputs.spendingPassword.Layout(win.gtx, &win.inputs.spendingPassword)
			}),
			layout.Rigid(func() {
				win.Err()
			}),
			layout.Rigid(func() {
				win.outputs.deleteWallet.Layout(win.gtx, &win.inputs.deleteWallet)
			}),
		)
	})
}

func (win *Window) RestoreDiag() {
	win.theme.Surface(win.gtx, func() {
		toMax(win.gtx)
		layout.UniformInset(unit.Dp(30)).Layout(win.gtx, func() {
			layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(win.gtx,
				layout.Rigid(func() {
					layout.E.Layout(win.gtx, func() {
						win.outputs.cancelDiag.Layout(win.gtx, &win.cancelDialog)
					})
				}),
				layout.Rigid(func() {
					d := win.theme.H3("Restore Wallet")
					d.Layout(win.gtx)
				}),
				layout.Rigid(func() {
					win.outputs.spendingPassword.Layout(win.gtx, &win.inputs.spendingPassword)
				}),
				layout.Rigid(func() {
					win.outputs.matchSpending.Layout(win.gtx, &win.inputs.matchSpending)
				}),
				layout.Rigid(func() {
					win.Err()
				}),
				layout.Rigid(func() {
					win.outputs.restoreWallet.Layout(win.gtx, &win.inputs.restoreWallet)
				}),
			)
		})
	})
}

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

func (win *Window) infoDiag() {
	win.theme.Surface(win.gtx, func() {
		layout.Center.Layout(win.gtx, func() {
			selectedDetails := func() {
				layout.UniformInset(unit.Dp(10)).Layout(win.gtx, func() {
					layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(win.gtx,
						layout.Rigid(func() {
							layout.UniformInset(unit.Dp(10)).Layout(win.gtx, func() {
								win.outputs.pageInfo.Layout(win.gtx)
							})
						}),
						layout.Rigid(func() {
							inset := layout.Inset{
								Left: unit.Dp(190),
							}
							inset.Layout(win.gtx, func() {
								win.outputs.gotItDiag.Layout(win.gtx, &win.inputs.receiveIcons.gotItDiag)
							})
						}),
					)
				})
			}
			decredmaterial.Modal{}.Layout(win.gtx, selectedDetails)
		})
	})
}

func (win *Window) generateNewAddressDiag() {
	layout.Flex{}.Layout(win.gtx,
		layout.Flexed(0.80, func() {
		}),
		layout.Flexed(1, func() {
			inset := layout.Inset{
				Top: unit.Dp(150),
			}
			inset.Layout(win.gtx, func() {
				win.gtx.Constraints.Width.Min = 40
				win.gtx.Constraints.Height.Min = 40
				win.outputs.newAddress.Layout(win.gtx, &win.inputs.receiveIcons.newAddress)
			})
		}),
	)
}

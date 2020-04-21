package ui

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
)

var (
	inputGroupContainerLeft  = &layout.List{Axis: layout.Vertical}
	inputGroupContainerRight = &layout.List{Axis: layout.Vertical}
)

type createRestore struct {
	gtx          *layout.Context
	theme        *decredmaterial.Theme
	inputs       *inputs
	outputs      *outputs
	err          func()
	walletExists bool
}

// Loading lays out the loading widget with a faded background
func (win *Window) CreateRestorePage() {
	pg := createRestore{
		gtx:          win.gtx,
		theme:        win.theme,
		inputs:       &win.inputs,
		outputs:      &win.outputs,
		err:          win.Err,
		walletExists: win.walletInfo.LoadedWallets > 0,
	}
	win.theme.Surface(win.gtx, func() {
		toMax(win.gtx)
		pd := unit.Dp(15)
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(win.gtx,
			layout.Flexed(1, func() {
				layout.Inset{Top: pd, Left: pd, Right: pd}.Layout(win.gtx, func() {
					layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(win.gtx,
						layout.Rigid(func() {
							layout.W.Layout(win.gtx, func() {
								if pg.walletExists {
									win.outputs.icons.back.Layout(win.gtx, &win.backCreateRestore)
								}
							})
						}),
						layout.Flexed(1, func() {
							if win.states.restoreWallet {
								pg.Restore()()
							} else {
								pg.mainContent()()
							}
						}),
					)
				})
			}),
		)
	})
}

func (pg *createRestore) mainContent() layout.Widget {
	return func() {
		layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
			layout.Flexed(1, func() {
				layout.Center.Layout(pg.gtx, func() {
					title := pg.theme.H3("")
					title.Alignment = text.Middle
					if pg.walletExists {
						title.Text = "Create or Restore Wallet"
					} else {
						title.Text = "Welcome to Decred Wallet, a secure & open-source desktop wallet."
					}
					title.Layout(pg.gtx)
				})
			}),
			layout.Rigid(func() {
				btnPadding := unit.Dp(10)
				layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
					layout.Rigid(func() {
						layout.Inset{Top: btnPadding, Bottom: btnPadding}.Layout(pg.gtx, func() {
							pg.outputs.createDiag.Layout(pg.gtx, &pg.inputs.createDiag)
						})
					}),
					layout.Rigid(func() {
						layout.Inset{Top: btnPadding, Bottom: btnPadding}.Layout(pg.gtx, func() {
							pg.outputs.showRestoreWallet.Layout(pg.gtx, &pg.inputs.showRestoreWallet)
						})
					}),
				)
			}),
		)
	}
}

func (pg *createRestore) Restore() layout.Widget {
	return func() {
		layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
			layout.Rigid(func() {
				txt := pg.theme.H3("Restore from seed phrase")
				txt.Alignment = text.Middle
				txt.Layout(pg.gtx)
			}),
			layout.Rigid(func() {
				txt := pg.theme.H6("Enter your seed phrase in the correct order")
				txt.Alignment = text.Middle
				txt.Layout(pg.gtx)
			}),
			layout.Rigid(func() {
				layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10)}.Layout(pg.gtx, func() {
					layout.Center.Layout(pg.gtx, func() {
						pg.err()
					})
				})
			}),
			layout.Flexed(1, func() {
				layout.Center.Layout(pg.gtx, func() {
					layout.Flex{}.Layout(pg.gtx,
						layout.Rigid(func() {
							pg.inputsGroup(inputGroupContainerLeft, 16, 0)
						}),
						layout.Rigid(func() {
							pg.inputsGroup(inputGroupContainerRight, 17, 16)
						}),
					)
				})
			}),
			layout.Rigid(func() {
				layout.Center.Layout(pg.gtx, func() {
					layout.Inset{Top: unit.Dp(15), Bottom: unit.Dp(15)}.Layout(pg.gtx, func() {
						pg.outputs.restoreDiag.Layout(pg.gtx, &pg.inputs.restoreDiag)
					})
				})
			}),
		)
	}
}

func (pg *createRestore) inputsGroup(l *layout.List, len int, startIndex int) {
	pg.gtx.Constraints.Width.Min = pg.gtx.Constraints.Width.Max / 2
	l.Layout(pg.gtx, len, func(i int) {
		layout.Flex{Axis: layout.Vertical}.Layout(pg.gtx,
			layout.Rigid(func() {
				layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(pg.gtx,
					layout.Rigid(func() {
						pg.theme.Label(unit.Dp(16), fmt.Sprintf("Word #%d", i+startIndex+1)).Layout(pg.gtx)
					}),
					layout.Rigid(func() {
						layout.Inset{Left: unit.Dp(20), Bottom: unit.Dp(20)}.Layout(pg.gtx, func() {
							pg.outputs.seedEditors[i+startIndex].Layout(pg.gtx, &pg.inputs.seedEditors.editors[i+startIndex])
						})
					}),
				)
			}),
			layout.Rigid(func() {
				pg.autoComplete(i, startIndex)
			}),
		)
	})
}

func (pg *createRestore) autoComplete(index, startIndex int) {
	if !pg.inputs.seedEditors.editors[index+startIndex].Focused() {
		return
	}

	(&layout.List{Axis: layout.Horizontal}).Layout(pg.gtx, len(pg.inputs.seedsSuggestions), func(i int) {
		layout.Inset{Right: unit.Dp(4)}.Layout(pg.gtx, func() {
			pg.outputs.seedsSuggestions[i].Layout(pg.gtx, &pg.inputs.seedsSuggestions[i].button)
		})
	})
}

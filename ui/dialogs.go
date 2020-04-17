package ui

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
)

func (win *Window) CreateDiag() {
	win.theme.Surface(win.gtx, func() {
		toMax(win.gtx)
		pd := unit.Dp(15)
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(win.gtx,
			layout.Flexed(1, func() {
				layout.Inset{Top: pd, Left: pd, Right: pd}.Layout(win.gtx, func() {
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
					)
				})
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
		pd := unit.Dp(15)
		layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(win.gtx,
			layout.Flexed(1, func() {
				layout.Inset{Top: pd, Left: pd, Right: pd}.Layout(win.gtx, func() {
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
					)
				})
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
		layout.UniformInset(unit.Dp(10)).Layout(win.gtx, func() {
			layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(win.gtx,
				layout.Rigid(func() {
					layout.UniformInset(unit.Dp(10)).Layout(win.gtx, func() {
						win.outputs.pageInfo.Layout(win.gtx)
					})
				}),
				layout.Rigid(func() {
					inset := layout.Inset{
						Left: unit.Dp(10),
					}
					inset.Layout(win.gtx, func() {
						win.outputs.gotItDiag.Layout(win.gtx, &win.inputs.receiveIcons.gotItDiag)
					})
				}),
			)
		})
		//decredmaterial.Modal{}.Layout(win.gtx, selectedDetails)
	})
}

func (win *Window) transactionsFilters() {
	w := win.gtx.Constraints.Width.Max / 2
	win.theme.Surface(win.gtx, func() {
		layout.UniformInset(unit.Dp(20)).Layout(win.gtx, func() {
			win.gtx.Constraints.Width.Min = w
			layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
				layout.Rigid(func() {
					win.gtx.Constraints.Width.Min = w
					layout.E.Layout(win.gtx, func() {
						win.outputs.cancelDiag.Layout(win.gtx, &win.cancelDialog)
					})
				}),
				layout.Rigid(func() {
					layout.Stack{}.Layout(win.gtx,
						layout.Expanded(func() {
							win.gtx.Constraints.Width.Min = w
							headTxt := win.theme.H4("Transactions filters")
							headTxt.Alignment = text.Middle
							headTxt.Layout(win.gtx)
						}),
					)
				}),
				layout.Rigid(func() {
					win.gtx.Constraints.Width.Min = w
					layout.Flex{Axis: layout.Horizontal}.Layout(win.gtx,
						layout.Flexed(.25, func() {
							layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
								layout.Rigid(func() {
									win.theme.H5("Order").Layout(win.gtx)
								}),
								layout.Rigid(func() {
									(&layout.List{Axis: layout.Vertical}).
										Layout(win.gtx, len(win.outputs.transactionFilterSort), func(index int) {
											win.outputs.transactionFilterSort[index].Layout(win.gtx, win.inputs.transactionFilterSort)
										})
								}),
							)
						}),
						layout.Flexed(.25, func() {
							layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
								layout.Rigid(func() {
									win.theme.H5("Direction").Layout(win.gtx)
								}),
								layout.Rigid(func() {
									(&layout.List{Axis: layout.Vertical}).
										Layout(win.gtx, len(win.outputs.transactionFilterDirection), func(index int) {
											win.outputs.transactionFilterDirection[index].Layout(win.gtx, win.inputs.transactionFilterDirection)
										})
								}),
							)
						}),
					)
				}),
				layout.Rigid(func() {
					layout.Inset{Top: unit.Dp(20)}.Layout(win.gtx, func() {
						win.outputs.applyFiltersTransactions.Layout(win.gtx, &win.inputs.applyFiltersTransactions)
					})
				}),
			)
		})
	})
}

package ui

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
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

func (win *Window) verifyMessageDiag() {
	win.theme.Surface(win.gtx, func() {
		win.gtx.Constraints.Width.Min = win.gtx.Px(unit.Dp(550))
		win.gtx.Constraints.Width.Max = win.gtx.Constraints.Width.Min
		layout.UniformInset(unit.Dp(20)).Layout(win.gtx, func() {
			win.vFlex(
				rigid(func() {
					win.hFlex(
						rigid(func() {
							win.theme.H5("Verify Wallet Message").Layout(win.gtx)
						}),
						layout.Flexed(.7, func() {
							layout.E.Layout(win.gtx, func() {
								win.outputs.verifyInfo.Layout(win.gtx, &win.inputs.verifyInfo)
							})
						}),
						layout.Flexed(1, func() {
							layout.E.Layout(win.gtx, func() {
								win.outputs.cancelDiag.Padding = unit.Dp(5)
								win.outputs.cancelDiag.Size = unit.Dp(35)
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
										win.theme.H6("Enter Address").Layout(win.gtx)
									}),
									rigid(func() {
										win.hFlexSB(
											rigid(func() {
												decredmaterial.Card{}.Layout(win.gtx, func() {
													win.hFlexSB(
														layout.Flexed(0.9, func() {
															win.outputs.addressInput.Layout(win.gtx, &win.inputs.addressInput)
														}),
													)
												})
											}),
											rigid(func() {
												inset := layout.Inset{
													Left: unit.Dp(10),
												}
												inset.Layout(win.gtx, func() {
													if win.inputs.addressInput.Text() == "" {
														win.outputs.pasteAddr.Layout(win.gtx, &win.inputs.pasteAddr)
													} else {
														win.outputs.clearAddr.Layout(win.gtx, &win.inputs.clearAddr)
													}
												})
											}),
										)
									}),
								)
							})
						}),
						rigid(func() {
							inset := layout.Inset{
								Top:    unit.Dp(10),
								Bottom: unit.Dp(10),
							}
							inset.Layout(win.gtx, func() {
								win.vFlexSB(
									rigid(func() {
										win.theme.H6("Enter Signature").Layout(win.gtx)
									}),
									rigid(func() {
										win.hFlexSB(
											rigid(func() {
												decredmaterial.Card{}.Layout(win.gtx, func() {
													win.hFlexSB(
														layout.Flexed(0.9, func() {
															win.outputs.signInput.Layout(win.gtx, &win.inputs.signInput)
														}),
													)
												})
											}),
											rigid(func() {
												inset := layout.Inset{
													Left: unit.Dp(10),
												}
												inset.Layout(win.gtx, func() {
													if win.inputs.signInput.Text() == "" {
														win.outputs.pasteSign.Layout(win.gtx, &win.inputs.pasteSign)
													} else {
														win.outputs.clearSign.Layout(win.gtx, &win.inputs.clearSign)
													}
												})
											}),
										)
									}),
								)
							})
						}),
						rigid(func() {
							win.theme.H6("Enter Message").Layout(win.gtx)
						}),
						rigid(func() {
							win.hFlexSB(
								rigid(func() {
									decredmaterial.Card{}.Layout(win.gtx, func() {
										win.hFlexSB(
											layout.Flexed(0.9, func() {
												win.outputs.messageInput.Layout(win.gtx, &win.inputs.messageInput)
											}),
										)
									})
								}),
								rigid(func() {
									inset := layout.Inset{
										Left:   unit.Dp(10),
										Bottom: unit.Dp(10),
									}
									inset.Layout(win.gtx, func() {
										if win.inputs.messageInput.Text() == "" {
											win.outputs.pasteMsg.Layout(win.gtx, &win.inputs.pasteMsg)
										} else {
											win.outputs.clearMsg.Layout(win.gtx, &win.inputs.clearMsg)
										}
									})
								}),
							)
						}),
						rigid(func() {
							layout.Flex{}.Layout(win.gtx,
								layout.Flexed(.6, func() {
									layout.Inset{Bottom: unit.Dp(5), Top: unit.Dp(10)}.Layout(win.gtx, func() {
										win.outputs.verifyMessage.Layout(win.gtx)
									})
								}),
								layout.Flexed(.4, func() {
									layout.Flex{}.Layout(win.gtx,
										layout.Flexed(.5, func() {
											layout.Inset{Left: unit.Dp(0), Right: unit.Dp(10)}.Layout(win.gtx, func() {
												win.outputs.clearBtn.Layout(win.gtx, &win.inputs.clearBtn)
											})
										}),
										layout.Flexed(.5, func() {
											win.outputs.verifyBtn.Layout(win.gtx, &win.inputs.verifyBtn)
										}),
									)
								}),
							)
						}),
					)
				}),
			)
		})
	})
}

func (win *Window) msgInfoDiag() {
	var msg = "After you or your counterparty has genrated a signature, you can use this \nform to verify the signature." +
		"\n\nOnce you have entered the address, the message and the corresponding \nsignature, you will see VALID if the signature" +
		"appropriately matches \nthe address and message otherwise INVALID."
	win.theme.Surface(win.gtx, func() {
		layout.UniformInset(unit.Dp(10)).Layout(win.gtx, func() {
			layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(win.gtx,
				layout.Rigid(func() {
					win.theme.H5("Verify Message").Layout(win.gtx)
				}),
				layout.Rigid(func() {
					inset := layout.Inset{
						Top:    unit.Dp(10),
						Bottom: unit.Dp(10),
					}
					inset.Layout(win.gtx, func() {
						win.theme.Body1(msg).Layout(win.gtx)
					})
				}),
				layout.Rigid(func() {
					layout.Flex{}.Layout(win.gtx,
						layout.Rigid(func() {
							win.outputs.hideMsgInfo.Layout(win.gtx, &win.inputs.hideMsgInfo)
						}),
					)
				}),
			)
		})
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

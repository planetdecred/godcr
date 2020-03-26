package ui

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
	"github.com/skip2/go-qrcode"
)

var pageContainer = &layout.List{Axis: layout.Vertical}

func (win *Window) Receive() {
	if win.walletInfo.LoadedWallets == 0 {
		win.Page(func() {
			win.outputs.noWallet.Layout(win.gtx)
		})
		return
	}

	body := func() {
		layout.Stack{}.Layout(win.gtx,
			layout.Expanded(func() {
				layout.Flex{}.Layout(win.gtx,
					layout.Rigid(func() {
						win.combined.sel.Layout(win.gtx, func() {

						})
					}),
					layout.Rigid(func() {
						win.ReceivePageContents()
					}),
				)
			}),
		)
	}
	win.TabbedPage(body)
}

func (win *Window) ReceivePageContents() {
	ReceivePageContent := []func(){
		func() {
			win.pageHeaderColumn()
		},
		func() {
			win.selectedAcountColumn()
		},
		func() {
			win.qrCodeAddressColumn()
		},
		func() {
			layout.Flex{}.Layout(win.gtx,
				layout.Flexed(0.35, func() {
				}),
				layout.Flexed(1, func() {
					if win.addressCopiedLabel.Text != "" {
						win.addressCopiedLabel.Layout(win.gtx)
					}
				}),
			)
		},
		func() {
			layout.Flex{}.Layout(win.gtx,
				layout.Flexed(0.35, func() {
				}),
				layout.Flexed(1, func() {
					win.Err()
				}),
			)
		},
	}

	pageContainer.Layout(win.gtx, len(ReceivePageContent), func(i int) {
		layout.Inset{Left: unit.Dp(3)}.Layout(win.gtx, ReceivePageContent[i])
	})

	if newAddr {
		win.generateNewAddress()
	}
}

func (win *Window) pageHeaderColumn() {
	layout.Flex{}.Layout(win.gtx,
		layout.Flexed(.6, func() {
			win.outputs.pageTitle.Layout(win.gtx)
		}),
		layout.Flexed(.4, func() {
			layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(20)}.Layout(win.gtx, func() {
				layout.Flex{}.Layout(win.gtx,
					layout.Flexed(.5, func() {
						win.outputs.info.Layout(win.gtx, &win.inputs.receiveIcons.info)
					}),
					layout.Flexed(.5, func() {
						win.outputs.more.Layout(win.gtx, &win.inputs.receiveIcons.more)
					}),
				)
			})
		}),
	)
}

func (win *Window) qrCodeAddressColumn() {
	addrs := win.walletInfo.Wallets[win.selected].Accounts[win.selectedAccount].CurrentAddress
	qrCode, err := qrcode.New(addrs, qrcode.Highest)
	if err != nil {
		win.err = err.Error()
		return
	}
	win.err = ""
	qrCode.DisableBorder = true
	layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
		layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(16)}.Layout(win.gtx, func() {
				layout.Flex{}.Layout(win.gtx,
					layout.Flexed(0.2, func() {
					}),
					layout.Flexed(1, func() {
						img := win.theme.Image(paint.NewImageOp(qrCode.Image(520)))
						img.Src.Rect.Max.X = 521
						img.Src.Rect.Max.Y = 521
						img.Scale = 0.5
						img.Layout(win.gtx)
					}),
				)
			})
		}),
		layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(16)}.Layout(win.gtx, func() {
				layout.Flex{}.Layout(win.gtx,
					layout.Flexed(0.1, func() {
					}),
					layout.Flexed(1, func() {
						win.receiveAddressColumn(addrs)
					}),
				)
			})
		}),
	)
}

func (win *Window) receiveAddressColumn(addrs string) {
	layout.Flex{}.Layout(win.gtx,
		layout.Rigid(func() {
			win.outputs.receiveAddressLabel.Text = addrs
			win.outputs.receiveAddressLabel.Layout(win.gtx)
		}),
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(16)}.Layout(win.gtx, func() {
				win.outputs.copy.Layout(win.gtx, &win.inputs.receiveIcons.copy)
			})
		}),
	)
}

func (win *Window) selectedAcountColumn() {
	info := win.walletInfo.Wallets[win.selected]
	win.outputs.selectedWalletNameLabel.Text = info.Name
	win.outputs.selectedWalletBalLabel.Text = info.Balance

	account := win.walletInfo.Wallets[win.selected].Accounts[win.selectedAccount]
	win.outputs.selectedAccountNameLabel.Text = account.Name
	win.outputs.selectedAccountBalanceLabel.Text = account.SpendableBalance

	layout.Flex{}.Layout(win.gtx,
		layout.Flexed(0.22, func() {
		}),
		layout.Flexed(1, func() {
			layout.Stack{}.Layout(win.gtx,
				layout.Stacked(func() {
					selectedDetails := func() {
						layout.UniformInset(unit.Dp(10)).Layout(win.gtx, func() {
							layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
								layout.Rigid(func() {
									layout.Flex{}.Layout(win.gtx,
										layout.Rigid(func() {
											layout.Inset{Bottom: unit.Dp(5)}.Layout(win.gtx, func() {
												win.outputs.selectedAccountNameLabel.Layout(win.gtx)
											})
										}),
										layout.Rigid(func() {
											layout.Inset{Left: unit.Dp(20)}.Layout(win.gtx, func() {
												win.outputs.selectedAccountBalanceLabel.Layout(win.gtx)
											})
										}),
									)
								}),
								layout.Rigid(func() {
									layout.Inset{Left: unit.Dp(20)}.Layout(win.gtx, func() {
										layout.Flex{}.Layout(win.gtx,
											layout.Rigid(func() {
												layout.Inset{Bottom: unit.Dp(5)}.Layout(win.gtx, func() {
													win.outputs.selectedWalletNameLabel.Layout(win.gtx)
												})
											}),
											layout.Rigid(func() {
												layout.Inset{Left: unit.Dp(22)}.Layout(win.gtx, func() {
													win.outputs.selectedWalletBalLabel.Layout(win.gtx)
												})
											}),
										)
									})
								}),
							)
						})
					}
					decredmaterial.Card{}.Layout(win.gtx, selectedDetails)
				}),
			)
		}),
	)
}

func (win *Window) generateNewAddress() {
	layout.Flex{}.Layout(win.gtx,
		layout.Flexed(0.71, func() {
		}),
		layout.Flexed(1, func() {
			inset := layout.Inset{
				Top: unit.Dp(45),
			}
			inset.Layout(win.gtx, func() {
				win.gtx.Constraints.Width.Min = 40
				win.gtx.Constraints.Height.Min = 40
				win.outputs.newAddressDiag.Layout(win.gtx, &win.inputs.receiveIcons.newAddressDiag)
			})
		}),
	)
}

package ui

import (
	"time"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/atotto/clipboard"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr-gio/ui/decredmaterial"
	"github.com/raedahgroup/godcr-gio/wallet"
	"github.com/skip2/go-qrcode"
)

var (
	listContainer = &layout.List{Axis: layout.Vertical}

	generateNew bool
	addrs       string

	isInfoBtnModal = false
	isNewAddrModal = false

	selectedWallet  *wallet.InfoShort
	selectedAccount *wallet.Account

	pageTitle       = "Receiving DCR"
	ReceivePageInfo = "Each time you request a payment, a \nnew address is created to protect \nyour privacy."
)

func (win *Window) Receive() {
	if win.walletInfo.LoadedWallets != 0 {
		win.setDefaultPageValues()
	}

	body := func() {
		layout.Stack{}.Layout(win.gtx,
			layout.Expanded(func() {
				win.ReceivePageContents()
			}),
		)
	}
	win.TabbedPage(body)
}

func (win *Window) ReceivePageContents() {
	ReceivePageContent := []func(){
		func() {
			win.pageFirstColumn()
		},
		func() {
			win.selectedAccountColumn()
		},
		func() {
			win.qrCodeAddressColumn()
		},
		func() {
			layout.Flex{}.Layout(win.gtx,
				layout.Flexed(.35, func() {
				}),
				layout.Flexed(1, func() {
					if win.addressCopiedLabel.Text != "" {
						win.addressCopiedLabel.Layout(win.gtx)
					}
				}),
			)
		},
	}

	listContainer.Layout(win.gtx, len(ReceivePageContent), func(i int) {
		layout.Inset{Left: unit.Dp(20)}.Layout(win.gtx, ReceivePageContent[i])
	})
	if isNewAddrModal {
		win.drawMoreModal()
	}
	if isInfoBtnModal {
		win.drawInfoModal()
	}
}

func (win *Window) pageFirstColumn() {
	layout.Flex{}.Layout(win.gtx,
		layout.Flexed(.6, func() {
			win.theme.H4(pageTitle).Layout(win.gtx)
		}),
		layout.Flexed(.4, func() {
			layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(20)}.Layout(win.gtx, func() {
				layout.Flex{}.Layout(win.gtx,
					layout.Flexed(.5, func() {
						if win.inputs.info.Clicked(win.gtx) {
							isInfoBtnModal = true
							isNewAddrModal = false
						}
						win.outputs.info.Layout(win.gtx, &win.inputs.info)
					}),
					layout.Flexed(.5, func() {
						for win.inputs.more.Clicked(win.gtx) {
							if isNewAddrModal {
								isInfoBtnModal = false
								isNewAddrModal = false
							} else {
								isNewAddrModal = true
							}
						}
						win.outputs.more.Layout(win.gtx, &win.inputs.more)
					}),
				)
			})
		}),
	)
}

func (win *Window) selectedAccountColumn() {
	layout.Flex{}.Layout(win.gtx,
		layout.Flexed(.24, func() {
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

func (win *Window) qrCodeAddressColumn() {
	qrCode, err := qrcode.New(addrs, qrcode.Highest)
	if err != nil {
		win.outputs.err.Text = err.Error()
		return
	}

	qrCode.DisableBorder = true
	layout.Flex{Axis: layout.Vertical}.Layout(win.gtx,
		layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(16)}.Layout(win.gtx, func() {
				layout.Flex{}.Layout(win.gtx,
					layout.Flexed(.2, func() {
					}),
					layout.Flexed(1, func() {
						img := win.theme.Image(paint.NewImageOp(qrCode.Image(140)))
						img.Src.Rect.Max.X = 141
						img.Src.Rect.Max.Y = 141
						img.Layout(win.gtx)
					}),
				)
			})
		}),
		layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(16)}.Layout(win.gtx, func() {
				layout.Flex{}.Layout(win.gtx,
					layout.Flexed(.1, func() {
					}),
					layout.Flexed(1, func() {
						win.receiveAddressColumn()
					}),
				)
			})
		}),
	)
}

func (win *Window) receiveAddressColumn() {
	layout.Flex{}.Layout(win.gtx,
		layout.Rigid(func() {
			win.outputs.receiveAddressLabel.Text = addrs
			win.outputs.receiveAddressLabel.Layout(win.gtx)
		}),
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(16)}.Layout(win.gtx, func() {
				for win.inputs.copy.Clicked(win.gtx) {
					clipboard.WriteAll(addrs)
					win.addressCopiedLabel.Text = "Address Copied"
					time.AfterFunc(time.Second*9, func() {
						win.addressCopiedLabel.Text = ""
					})
				}
				win.outputs.copy.Layout(win.gtx, &win.inputs.copy)
			})
		}),
	)
}

func (win *Window) setDefaultPageValues() {
	info := win.walletInfo.Wallets[win.selected]

	for i := range info.Accounts {
		win.setSelectedAccount(info, info.Accounts[i], false)
		break
	}
}

func (win *Window) setSelectedAccount(wallet wallet.InfoShort, account wallet.Account, generateNew bool) {
	selectedWallet = &wallet
	selectedAccount = &account

	win.outputs.selectedAccountNameLabel.Text = account.Name
	win.outputs.selectedWalletNameLabel.Text = wallet.Name
	win.outputs.selectedWalletBalLabel.Text = dcrutil.Amount(account.SpendableBalance).String()
	win.outputs.selectedAccountBalanceLabel.Text = wallet.Balance

	var addr string
	var err error

	// create a new receive address everytime a new account is chosen
	if generateNew {
		addr, err = win.wallet.NextAddress(wallet.ID, account.Number)
		if err != nil {
			win.outputs.err.Text = err.Error()
			return
		}
	} else {
		addr, err = win.wallet.CurrentAddress(wallet.ID, account.Number)
		if err != nil {
			win.outputs.err.Text = err.Error()
			return
		}
	}
	addrs = addr
}

func (win *Window) drawInfoModal() {
	win.theme.Surface(win.gtx, func() {
		layout.Center.Layout(win.gtx, func() {
			selectedDetails := func() {
				layout.UniformInset(unit.Dp(10)).Layout(win.gtx, func() {
					layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEvenly}.Layout(win.gtx,
						layout.Rigid(func() {
							layout.UniformInset(unit.Dp(10)).Layout(win.gtx, func() {
								win.theme.Body1(ReceivePageInfo).Layout(win.gtx)
							})
						}),
						layout.Rigid(func() {
							inset := layout.Inset{
								Left: unit.Dp(190),
							}
							inset.Layout(win.gtx, func() {
								if win.inputs.gotIt.Clicked(win.gtx) {
									if isInfoBtnModal {
										isInfoBtnModal = false
									}
								}

								win.outputs.gotIt.Layout(win.gtx, &win.inputs.gotIt)
							})
						}),
					// })
					)
				})
			}
			decredmaterial.Modal{layout.SE, false}.Layout(win.gtx, selectedDetails)
		})
	})
}

func (win *Window) drawMoreModal() {
	layout.Flex{}.Layout(win.gtx,
		layout.Flexed(.73, func() {
		}),
		layout.Flexed(1, func() {
			inset := layout.Inset{
				Top: unit.Dp(50),
			}
			inset.Layout(win.gtx, func() {
				for win.inputs.newAddress.Clicked(win.gtx) {
					if isNewAddrModal {
						win.setSelectedAccount(*selectedWallet, *selectedAccount, true)
						isNewAddrModal = false
					}
				}

				win.gtx.Constraints.Width.Min = 40
				win.gtx.Constraints.Height.Min = 40
				win.outputs.newAddress.Layout(win.gtx, &win.inputs.newAddress)
			})
		}),
	)
	// }
	// decredmaterial.Modal{layout.SE, false}.Layout(win.gtx, selectedDetails)
	// })
}

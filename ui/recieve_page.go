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
			layout.Center.Layout(win.gtx, func() {
				if win.outputs.err.Text != "" {
					win.outputs.err.Layout(win.gtx)
				}
			})
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
	if win.isNewAddrModal {
		win.drawMoreModal()
	}
	if win.isInfoBtnModal {
		win.drawInfoModal()
	}
}

func (win *Window) pageFirstColumn() {
	layout.Flex{}.Layout(win.gtx,
		layout.Flexed(.6, func() {
			win.outputs.receivePageTitle.Layout(win.gtx)
		}),
		layout.Flexed(.4, func() {
			layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(20)}.Layout(win.gtx, func() {
				layout.Flex{}.Layout(win.gtx,
					layout.Flexed(.5, func() {
						if win.inputs.receiveIcons.info.Clicked(win.gtx) {
							win.isInfoBtnModal = true
							win.isNewAddrModal = false
						}
						win.outputs.info.Layout(win.gtx, &win.inputs.receiveIcons.info)
					}),
					layout.Flexed(.5, func() {
						for win.inputs.receiveIcons.more.Clicked(win.gtx) {
							if win.isNewAddrModal {
								win.isInfoBtnModal = false
								win.isNewAddrModal = false
							} else {
								win.isNewAddrModal = true
							}
						}
						win.outputs.more.Layout(win.gtx, &win.inputs.receiveIcons.more)
					}),
				)
			})
		}),
	)
}

func (win *Window) selectedAccountColumn() {
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

func (win *Window) qrCodeAddressColumn() {
	qrCode, err := qrcode.New(win.addrs, qrcode.Highest)
	if err != nil {
		win.outputs.err.Text = err.Error()
		return
	}

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
			win.outputs.receiveAddressLabel.Text = win.addrs
			win.outputs.receiveAddressLabel.Layout(win.gtx)
		}),
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(16)}.Layout(win.gtx, func() {
				for win.inputs.receiveIcons.copy.Clicked(win.gtx) {
					clipboard.WriteAll(win.addrs)
					win.addressCopiedLabel.Text = "Address Copied"
					time.AfterFunc(time.Second*9, func() {
						win.addressCopiedLabel.Text = ""
					})
				}
				win.outputs.copy.Layout(win.gtx, &win.inputs.receiveIcons.copy)
			})
		}),
	)
}

func (win *Window) setDefaultPageValues() {
	info := win.walletInfo.Wallets[win.selected]

	for i := range info.Accounts {
		if len(info.Accounts) == 0 {
			continue
		}

		for win.inputs.receiveIcons.newAddress.Clicked(win.gtx) {
			addr, err := win.wallet.NextAddress(info.ID, info.Accounts[i].Number)
			if err != nil {
				win.outputs.err.Text = err.Error()
				return
			}
			info.Accounts[i].CurrentAddress = addr
			win.isNewAddrModal = false
		}

		win.setSelectedAccount(info, info.Accounts[i])
		break
	}
}

func (win *Window) setSelectedAccount(wallet wallet.InfoShort, account wallet.Account) {
	win.selectedWallet = &wallet
	win.selectedAccount = &account

	win.outputs.selectedAccountNameLabel.Text = account.Name
	win.outputs.selectedWalletNameLabel.Text = wallet.Name
	win.outputs.selectedWalletBalLabel.Text = dcrutil.Amount(account.SpendableBalance).String()
	win.outputs.selectedAccountBalanceLabel.Text = wallet.Balance
	win.addrs = account.CurrentAddress
}

func (win *Window) drawInfoModal() {
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
								if win.inputs.receiveIcons.gotItDiag.Clicked(win.gtx) {
									if win.isInfoBtnModal {
										win.isInfoBtnModal = false
									}
								}

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

func (win *Window) drawMoreModal() {
	layout.Flex{}.Layout(win.gtx,
		layout.Flexed(0.73, func() {
		}),
		layout.Flexed(1, func() {
			inset := layout.Inset{
				Top: unit.Dp(50),
			}
			inset.Layout(win.gtx, func() {
				win.gtx.Constraints.Width.Min = 40
				win.gtx.Constraints.Height.Min = 40
				win.outputs.newAddress.Layout(win.gtx, &win.inputs.receiveIcons.newAddress)
			})
		}),
	)
}

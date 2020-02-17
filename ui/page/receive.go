package page

import (
	"image"
	"time"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	// "gioui.org/gesture"

	"github.com/decred/dcrd/dcrutil"
	"github.com/atotto/clipboard"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr-gio/ui"
	"github.com/raedahgroup/godcr-gio/ui/themes/materialplus"
	"github.com/raedahgroup/godcr-gio/wallet"
	"github.com/skip2/go-qrcode"
)

// ReceivingID is the id of the receiving page.
const ReceivingID = "receive"

var (
	ReceivePageInfo   = "Each time you request a payment, a \nnew address is created to protect \nyour privacy."
	accountModalTitle = "Select an Account"
	pageTitle         = "Receiving DCR"
)

type infoModalWidgets struct {
	gotItdBtnWdg *widget.Button
	infoLabel    material.Label
}

type moreModalWidgets struct {
	generateNewAddBtnWdg *widget.Button
}

// Receive represents the receiving page of the app.
type Receive struct {
	wallet  *wallet.Wallet
	wallets []wallet.InfoShort

	copyBtn     material.IconButton
	infoBtn     material.IconButton
	moreBtn     material.IconButton
	dropDownBtn material.IconButton
	minimizeBtn material.IconButton

	generateNewAddBtn material.Button
	gotItBtn          material.Button

	copyBtnWdg     *widget.Button
	infoBtnWdg     *widget.Button
	moreBtnWdg     *widget.Button
	dropDownBtnWdg *widget.Button
	minimizeBtnWdg *widget.Button

	infoModalWidgets *infoModalWidgets
	moreModalWidgets *moreModalWidgets

	isGenerateNewAddBtnModal bool
	isInfoBtnModal           bool
	isAccountModalOpen       bool

	selectedWallet  *wallet.InfoShort
	selectedAccount *wallet.Account

	pageTitleLabel              material.Label
	selectedAccountNameLabel    material.Label
	selectedWalletLabel         material.Label
	selectedAccountBalanceLabel material.Label
	receiveAddressLabel         material.Label
	accountModalTitleLabel      material.Label
	errorLabel                  material.Label
	addressCopiedLabel          material.Label

	accountModalLine *materialplus.Line

	accountSelectorButtons map[string]*widget.Button

	theme *materialplus.Theme
	listContainer  layout.List
	receiveAddress string
	states         map[string]interface{}
}

// Init initializies the page with a label.
func (pg *Receive) Init(theme *materialplus.Theme, wal *wallet.Wallet, states map[string]interface{}) {
	pg.theme = theme
	pg.wallet = wal
	pg.states = states

	pg.isGenerateNewAddBtnModal = false
	pg.isInfoBtnModal = false
	pg.isAccountModalOpen = false

	pg.pageTitleLabel = theme.H3(pageTitle)
	pg.accountModalTitleLabel = theme.H6(accountModalTitle)
	pg.selectedAccountNameLabel = pg.theme.Body1("")
	pg.selectedWalletLabel = pg.theme.Caption("")
	pg.selectedAccountBalanceLabel = pg.theme.Body2("")
	pg.receiveAddressLabel = pg.theme.H6("")
	pg.receiveAddressLabel.Color = ui.LightBlueColor
	pg.errorLabel = theme.Body1("")
	pg.errorLabel.Color = ui.DangerColor
	pg.addressCopiedLabel = theme.Caption("")
	pg.addressCopiedLabel.Color = ui.LightBlueColor

	pg.copyBtnWdg = new(widget.Button)
	pg.copyBtn = theme.IconButton(materialplus.ContentCopyIcon)
	pg.copyBtn.Background = ui.WhiteColor
	pg.copyBtn.Color = ui.GrayColor
	pg.copyBtn.Padding = unit.Dp(5)
	pg.copyBtn.Size = unit.Dp(30)

	pg.infoBtnWdg = new(widget.Button)
	pg.infoBtn = theme.IconButton(materialplus.ActionInfoIcon)
	pg.infoBtn.Background = ui.WhiteColor
	pg.infoBtn.Color = ui.GrayColor
	pg.infoBtn.Padding = unit.Dp(5)
	pg.infoBtn.Size = unit.Dp(40)

	pg.moreBtnWdg = new(widget.Button)
	pg.moreBtn = theme.IconButton(materialplus.NavigationMoreIcon)
	pg.moreBtn.Background = ui.WhiteColor
	pg.moreBtn.Color = ui.GrayColor
	pg.moreBtn.Padding = unit.Dp(5)
	pg.moreBtn.Size = unit.Dp(40)

	pg.dropDownBtnWdg = new(widget.Button)
	pg.dropDownBtn = theme.IconButton(materialplus.DropDownIcon)
	pg.dropDownBtn.Background = ui.LightGrayColor
	pg.dropDownBtn.Color = ui.BlackColor
	pg.dropDownBtn.Padding = unit.Dp(5)
	pg.dropDownBtn.Size = unit.Dp(40)

	pg.minimizeBtnWdg = new(widget.Button)
	pg.minimizeBtn = theme.IconButton(materialplus.CancelIcon)
	pg.minimizeBtn.Background = ui.LightGrayColor
	pg.minimizeBtn.Color = ui.BlackColor
	pg.minimizeBtn.Padding = unit.Dp(0)
	pg.minimizeBtn.Size = unit.Dp(20)

	pg.generateNewAddBtn = theme.Button("Generate New Address")
	pg.generateNewAddBtn.Background = ui.GrayColor

	pg.gotItBtn = theme.Button("Got It")

	pg.accountSelectorButtons = map[string]*widget.Button{}
	pg.listContainer = layout.List{Axis: layout.Vertical}

	pg.accountModalLine = pg.theme.Line()
	pg.accountModalLine.Width = 230

	pg.moreModalWidgets = &moreModalWidgets{
		generateNewAddBtnWdg: new(widget.Button),
	}

	pg.infoModalWidgets = &infoModalWidgets{
		infoLabel:    pg.theme.Body1(ReceivePageInfo),
		gotItdBtnWdg: new(widget.Button),
	}
}

// Draw renders the page materialplus.
// It does not react to nor does it generate any event.
func (pg *Receive) Draw(gtx *layout.Context) (res interface{}) {
	pg.checkForStatesUpdate()
	if pg.wallets == nil {
		pg.waitAndSetWalletInfo()
	}

	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			pg.ReceivePageContents(gtx)
		}),
	)

	return nil
}

func (pg *Receive) ReceivePageContents(gtx *layout.Context) {
	ReceivePageContent := []func(){
		func() {
			pg.pageFirstColumn(gtx)
		},
		func() {
			layout.Align(layout.Center).Layout(gtx, func() {
				if pg.errorLabel.Text != "" {
					pg.errorLabel.Layout(gtx)
				}
			})
		},
		func() {
			layout.Align(layout.Center).Layout(gtx, func() {
				pg.selectedAccountLabel(gtx)
			})
		},
		func() {
			pg.generateAddressQrCode(gtx)
		},
		func() {
			layout.Align(layout.Center).Layout(gtx, func() {
				if pg.addressCopiedLabel.Text != "" {
					pg.addressCopiedLabel.Layout(gtx)
				}
			})
		},
	}

	pg.listContainer.Layout(gtx, len(ReceivePageContent), func(i int) {
		layout.UniformInset(unit.Dp(10)).Layout(gtx, ReceivePageContent[i])
	})

	if pg.isGenerateNewAddBtnModal {
		pg.drawMoreModal(gtx)
	}
	if pg.isInfoBtnModal {
		pg.drawInfoModal(gtx)
	}
	if pg.isAccountModalOpen {
		pg.accountSelectedModal(gtx)
	}
}

func (pg *Receive) waitAndSetWalletInfo() {
	walletInfoState := pg.states[StateWalletInfo]
	if walletInfoState == nil {
		return
	}

	walletInfo := walletInfoState.(*wallet.MultiWalletInfo)
	pg.setDefaultAccount(walletInfo.Wallets)
}

func (pg *Receive) setDefaultAccount(wallets []wallet.InfoShort) {
	pg.wallets = wallets

	for i := range wallets {
		if len(wallets[i].Accounts) == 0 {
			continue
		}

		pg.setSelectedAccount(wallets[i], wallets[i].Accounts[0], false)
		break
	}
}

func (pg *Receive) pageFirstColumn(gtx *layout.Context) {
	layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
		layout.Rigid(func() {
			pg.pageTitleLabel.Layout(gtx)
		}),
		layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(10)}.Layout(gtx, func() {
				layout.Flex{}.Layout(gtx,
					layout.Rigid(func() {
						if pg.infoBtnWdg.Clicked(gtx) {
							pg.isInfoBtnModal = true
							pg.isGenerateNewAddBtnModal = false
							pg.isAccountModalOpen = false
						}
						pg.infoBtn.Layout(gtx, pg.infoBtnWdg)
					}),
					layout.Rigid(func() {
						if pg.moreBtnWdg.Clicked(gtx) {
							pg.isGenerateNewAddBtnModal = true
							pg.isInfoBtnModal = false
							pg.isAccountModalOpen = false
						}
						pg.moreBtn.Layout(gtx, pg.moreBtnWdg)
					}),
				)
			})
		}),
	)
}

func (pg *Receive) generateAddressQrCode(gtx *layout.Context) {
	qrCode, err := qrcode.New(pg.receiveAddress, qrcode.Highest)
	if err != nil {
		return
	}

	qrCode.DisableBorder = true
	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func() {
			layout.Align(layout.Center).Layout(gtx, func() {
				img := pg.theme.Image(paint.NewImageOp(qrCode.Image(140)))
				img.Layout(gtx)
			})
		}),
		layout.Rigid(func() {
			layout.Inset{Top: unit.Dp(16)}.Layout(gtx, func() {
				layout.Align(layout.Center).Layout(gtx, func() {
					pg.receiveAddressSection(gtx)
				})
			})
		}),
	)
}

func (pg *Receive) receiveAddressSection(gtx *layout.Context) {
	layout.Flex{}.Layout(gtx,
		layout.Rigid(func() {
			pg.receiveAddressLabel.Text = pg.receiveAddress
			pg.receiveAddressLabel.Layout(gtx)
		}),
		layout.Rigid(func() {
			layout.Inset{Left: unit.Dp(16)}.Layout(gtx, func() {
				if pg.copyBtnWdg.Clicked(gtx) {
					clipboard.WriteAll(pg.receiveAddress)
					pg.addressCopiedLabel.Text = "Address Copied"
					time.AfterFunc(time.Second*1, func() {
						pg.addressCopiedLabel.Text = ""
					})
				}
				pg.copyBtn.Layout(gtx, pg.copyBtnWdg)
			})
		}),
	)
}

func (pg *Receive) drawMoreModal(gtx *layout.Context) {
	modalWidgetFuncs := []func(){
		func() {
			inset := layout.Inset{
				Top:   unit.Dp(50),
				Right: unit.Dp(20),
			}
			inset.Layout(gtx, func() {
				if pg.moreModalWidgets.generateNewAddBtnWdg.Clicked(gtx) {
					if pg.isGenerateNewAddBtnModal {
						pg.setSelectedAccount(*pg.selectedWallet, *pg.selectedAccount, true)
						pg.isGenerateNewAddBtnModal = false
					}
				}
				gtx.Constraints.Width.Min = 40
				gtx.Constraints.Height.Min = 40
				pg.generateNewAddBtn.Layout(gtx, pg.moreModalWidgets.generateNewAddBtnWdg)
			})
		},
	}

	inset := layout.Inset{
		Top:  unit.Dp(0),
		Left: unit.Dp(0),
	}

	inset.Layout(gtx, func() {
		layout.Stack{Alignment: layout.NE}.Layout(gtx,
			layout.Expanded(func() {
				gtx.Constraints.Height.Min = 170
			}),
			layout.Stacked(func() {
				layout.Inset{}.Layout(gtx, func() {
					layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
						list := layout.List{}
						list.Layout(gtx, len(modalWidgetFuncs), func(i int) {
							layout.UniformInset(unit.Dp(0)).Layout(gtx, modalWidgetFuncs[i])
						})
					})
				})
			}),
		)
	})
}

func (pg *Receive) drawInfoModal(gtx *layout.Context) {
	modalWidgetFuncs := []func(){
		func() {
			layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
				pg.infoModalWidgets.infoLabel.Layout(gtx)
			})
		},
		func() {
			inset := layout.Inset{
				Left: unit.Dp(190),
			}
			inset.Layout(gtx, func() {
				if pg.infoModalWidgets.gotItdBtnWdg.Clicked(gtx) {
					if pg.isInfoBtnModal {
						pg.isInfoBtnModal = false
					}
				}
				gtx.Constraints.Width.Min = 20
				gtx.Constraints.Height.Min = 20
				pg.gotItBtn.Layout(gtx, pg.infoModalWidgets.gotItdBtnWdg)
			})
		},
	}

	inset := layout.Inset{
		Top:  unit.Dp(0),
		Left: unit.Dp(0),
	}
	inset.Layout(gtx, func() {
		pg.theme.ModalPopUp(gtx, 270, 150, func() {
			list := layout.List{Axis: layout.Vertical}
			list.Layout(gtx, len(modalWidgetFuncs), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, modalWidgetFuncs[i])
			})
		})
	})
}

func (pg *Receive) selectedAccountLabel(gtx *layout.Context) {
	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			layout.Inset{}.Layout(gtx, func() {
				materialplus.Fill(gtx, ui.LightGrayColor, 200, 60)
			})
		}),
		layout.Stacked(func() {
			layout.Inset{}.Layout(gtx, func() {
				layout.Flex{}.Layout(gtx,
					layout.Rigid(func() {
						layout.Inset{Right: unit.Dp(30)}.Layout(gtx, func() {
							layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func() {
									layout.Inset{Bottom: unit.Dp(5)}.Layout(gtx, func() {
										pg.selectedAccountNameLabel.Layout(gtx)
									})
								}),
								layout.Rigid(func() {
									layout.Inset{Left: unit.Dp(2)}.Layout(gtx, func() {
										pg.selectedWalletLabel.Layout(gtx)
									})
								}),
							)
						})
					}),
					layout.Rigid(func() {
						layout.Inset{Top: unit.Dp(6.5)}.Layout(gtx, func() {
							pg.selectedAccountBalanceLabel.Layout(gtx)
						})
					}),
					layout.Rigid(func() {
						layout.Inset{Left: unit.Dp(15)}.Layout(gtx, func() {
							if pg.dropDownBtnWdg.Clicked(gtx) {
								if pg.isAccountModalOpen {
									pg.isAccountModalOpen = false
								} else {
									pg.isAccountModalOpen = true
									pg.isInfoBtnModal = false
									pg.isGenerateNewAddBtnModal = false
								}
							}
							pg.dropDownBtn.Layout(gtx, pg.dropDownBtnWdg)
						})
					}),
					layout.Rigid(func() {
						layout.Inset{Left: unit.Dp(15)}.Layout(gtx, func() {
								if pg.dropDownBtnWdg.Clicked(gtx) {
									pg.isInfoBtnModal = true
								}
								pg.dropDownBtn.Layout(gtx, pg.dropDownBtnWdg)
							})
					}),
				)
			})
		}),
	)
}

func (pg *Receive) registerAccountSelectorButton(accountName string) {
	if _, ok := pg.accountSelectorButtons[accountName]; !ok {
		pg.accountSelectorButtons[accountName] = new(widget.Button)
	}
}

func (pg *Receive) accountSelectedModal(gtx *layout.Context) {
	modalWidgetFuncs := []func(){
		func() {
			layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
				layout.Flex{}.Layout(gtx,
					layout.Rigid(func() {
						pg.accountModalTitleLabel.Layout(gtx)
					}),
					layout.Rigid(func() {
						layout.Inset{Left: unit.Dp(35), Top: unit.Dp(5)}.Layout(gtx, func() {
							if pg.minimizeBtnWdg.Clicked(gtx) {
								pg.isAccountModalOpen = false
							}
							pg.minimizeBtn.Layout(gtx, pg.minimizeBtnWdg)
						})
					}),
				)
			})
		},
		func() {
			layout.UniformInset(unit.Dp(1)).Layout(gtx, func() {
				pg.accountModalLine.Layout(gtx)
			})
		},
		func() {
			list := layout.List{Axis: layout.Vertical}
			list.Layout(gtx, len(pg.wallets), func(i int) {
				layout.Inset{Left: unit.Dp(30), Bottom: unit.Dp(10)}.Layout(gtx, func() {
					wallet := pg.wallets[i]

					walletNameLabel := pg.theme.H6(wallet.Name + " \t" + dcrutil.Amount(wallet.TotalBalance).String())
					walletNameLabel.Layout(gtx)

					list := layout.List{Axis: layout.Vertical}
					list.Layout(gtx, len(wallet.Accounts), func(k int) {
						account := pg.wallets[i].Accounts[k]
						if account.Name != "imported" {
							buttonKey := wallet.Name + account.Name
							pg.registerAccountSelectorButton(buttonKey)

							for pg.accountSelectorButtons[buttonKey].Clicked(gtx) {
								pg.setSelectedAccount(wallet, account, false)
								pg.isAccountModalOpen = false
							}

							layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func() {
									layout.Inset{Left: unit.Dp(15)}.Layout(gtx, func() {
										layout.Align(layout.Center).Layout(gtx, func() {
											inset := layout.Inset{
												Top: unit.Dp(25),
											}
											inset.Layout(gtx, func() {
												sendAccountNameLabel := pg.theme.Body1(account.Name + " \t" + dcrutil.Amount(account.TotalBalance).String())
												sendAccountNameLabel.Layout(gtx)
											})

											inset = layout.Inset{
												Top: unit.Dp(50),
											}
											inset.Layout(gtx, func() {
												spendableBalanceLabel := pg.theme.Body2("Spendable: \t" + dcrutil.Amount(account.SpendableBalance).String())
												spendableBalanceLabel.Layout(gtx)
											})
										})
									})
								}),
							)
							pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
							pg.accountSelectorButtons[buttonKey].Layout(gtx)
						}
					})
				})
			})
		},
	}

	inset := layout.Inset{
		Top:   unit.Dp(160),
		Right: unit.Dp(10),
	}

	inset.Layout(gtx, func() {
		layout.Stack{Alignment: layout.N}.Layout(gtx,
			layout.Stacked(func() {
				layout.Inset{}.Layout(gtx, func() {
					materialplus.PaintArea(gtx, ui.LightGrayColor, 230, 210)
					list := layout.List{Axis: layout.Vertical}
					list.Layout(gtx, len(modalWidgetFuncs), func(i int) {
						layout.UniformInset(unit.Dp(0)).Layout(gtx, modalWidgetFuncs[i])
					})
				})
			}),
		)
	})
}

func (pg *Receive) setSelectedAccount(wallet wallet.InfoShort, account wallet.Account, generateNew bool) {
	pg.selectedWallet = &wallet
	pg.selectedAccount = &account

	pg.selectedWalletLabel.Text = wallet.Name
	pg.selectedAccountNameLabel.Text = account.Name
	pg.selectedAccountBalanceLabel.Text = dcrutil.Amount(account.SpendableBalance).String()

	var addr string
	var err error

	// create a new receive address everytime a new account is chosen
	if generateNew {
		addr, err = pg.wallet.NextAddress(wallet.ID, account.Number)
		if err != nil {
			pg.errorLabel.Text = err.Error()
			return
		}
	} else {
		addr, err = pg.wallet.CurrentAddress(wallet.ID, account.Number)
		if err != nil {
			pg.errorLabel.Text = err.Error()
			return
		}
	}
	pg.receiveAddress = addr
}

func (pg *Receive) checkForStatesUpdate() {
	err := pg.states[StateError]
	if err == nil {
		return
	}

	if err != nil {
		pg.errorLabel.Text = err.(error).Error()
		delete(pg.states, StateError)
		return
	}
}

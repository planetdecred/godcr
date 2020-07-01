package ui

import (
	"image/color"
	"time"

	"github.com/raedahgroup/godcr/ui/values"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/wallet"
)

const PageWallet = "wallet"

const (
	subWalletMain = iota
	subWalletRename
	subWalletDelete
	subWalletChangePass
	subWalletBackup
	subWalletAddAcct
)

type walletPage struct {
	walletInfo *wallet.MultiWalletInfo
	subPage    int
	current    wallet.InfoShort
	wallet     *wallet.Wallet
	result     **wallet.Signature
	icons      struct {
		main, delete, rename, sign, verify, addWallet, changePass,
		addAcct, backup widget.Button
		mainW, deleteW, signW, verifyW, addWalletW, renameW,
		changePassW, addAcctW, backupW decredmaterial.IconButton
	}
	container, accountsList         layout.List
	delete, addAcct, rename         widget.Button
	line                            *decredmaterial.Line
	renameW, deleteW, cancelDeleteW decredmaterial.Button
	errorLabel                      decredmaterial.Label
	editor                          widget.Editor
	editorW                         decredmaterial.Editor
	passwordModal                   *decredmaterial.Password
	isPasswordModalOpen             bool
	errChann                        chan error
	errorText                       string
}

func (win *Window) WalletPage(common pageCommon) layout.Widget {
	page := &walletPage{
		walletInfo: win.walletInfo,
		container: layout.List{
			Axis: layout.Vertical,
		},
		accountsList: layout.List{
			Axis: layout.Vertical,
		},
		wallet:        common.wallet,
		line:          common.theme.Line(),
		editorW:       common.theme.Editor("New wallet name"),
		renameW:       common.theme.Button("Rename Wallet"),
		errorLabel:    common.theme.Body2(""),
		result:        &win.signatureResult,
		deleteW:       common.theme.DangerButton("Confirm Delete Wallet"),
		cancelDeleteW: common.theme.Button("Cancel Wallet Delete"),
		passwordModal: common.theme.Password(),
		errChann:      common.errorChannels[PageWallet],
		errorText:     "",
	}
	page.line.Color = common.theme.Color.Gray
	page.line.Height = 1
	page.errorLabel.Color = common.theme.Color.Danger

	var iconPadding = values.MarginPadding5
	var iconSize = values.MarginPadding30

	page.icons.addAcctW = common.theme.IconButton(common.icons.contentAdd)
	page.icons.addAcctW.Padding = iconPadding
	page.icons.addAcctW.Size = iconSize
	page.icons.mainW = common.theme.IconButton(common.icons.navigationArrowBack)
	page.icons.mainW.Background = color.RGBA{}
	page.icons.mainW.Color = common.theme.Color.Hint
	page.icons.mainW.Padding = values.MarginPadding0
	page.icons.mainW.Size = iconSize
	page.icons.deleteW = common.theme.IconButton(common.icons.actionDelete)
	page.icons.deleteW.Size = iconSize
	page.icons.deleteW.Padding = iconPadding
	page.icons.deleteW.Background = common.theme.Color.Danger
	page.icons.signW = common.theme.IconButton(common.icons.communicationComment)
	page.icons.signW.Size = iconSize
	page.icons.signW.Padding = iconPadding
	page.icons.verifyW = common.theme.IconButton(common.icons.verifyAction)
	page.icons.verifyW.Size = iconSize
	page.icons.verifyW.Padding = iconPadding
	page.icons.addWalletW = common.theme.IconButton(common.icons.contentAdd)
	page.icons.addWalletW.Size = iconSize
	page.icons.addWalletW.Padding = iconPadding
	page.icons.renameW = common.theme.IconButton(common.icons.editorModeEdit)
	page.icons.renameW.Size = iconSize
	page.icons.renameW.Padding = iconPadding
	page.icons.changePassW = common.theme.IconButton(common.icons.actionLock)
	page.icons.changePassW.Size = iconSize
	page.icons.changePassW.Padding = iconPadding
	page.icons.backupW = common.theme.IconButton(common.icons.actionBackup)
	page.icons.backupW.Size = iconSize
	page.icons.backupW.Padding = iconPadding

	return func() {
		page.Layout(common)
		page.Handle(common)
	}
}

// Layout lays out the widgets for the main wallets page.
func (page *walletPage) Layout(common pageCommon) {
	switch page.subPage {
	case subWalletMain:
		page.subMain(common)
	case subWalletRename:
		page.subRename(common)
	case subWalletDelete:
		page.subDelete(common)
	}
	if common.states.deleted {
		page.subPage = subWalletMain
		common.states.deleted = false
	}
}

func (page *walletPage) subMain(common pageCommon) {
	gtx := common.gtx

	page.current = common.info.Wallets[*common.selectedWallet]

	body := func() {
		layout.Stack{}.Layout(gtx,
			layout.Expanded(func() {
				layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Flexed(0.88, func() {
							page.topRow(common)
						}),
						layout.Flexed(0.12, func() {
							page.bottomRow(common)
						}),
					)
				})
			}),
		)
	}

	common.LayoutWithWallets(gtx, body)
}

func (page *walletPage) topRow(common pageCommon) {
	gtx := common.gtx
	wdgs := []func(){
		func() {
			page.alert(common)
		},
		func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					common.theme.H5(page.current.Name).Layout(common.gtx)
				}),
			)
		},
		func() {
			common.theme.H6("Total Balance: " + page.current.Balance).Layout(gtx)
		},
		func() {
			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func() {
					common.theme.H6("Accounts").Layout(gtx)
				}),
				layout.Rigid(func() {
					layout.Inset{Left: values.MarginPadding5}.Layout(common.gtx, func() {
						page.icons.addAcctW.Layout(gtx, &page.icons.addAcct)
					})
				}),
			)
		},
		func() {
			page.accountsList.Layout(gtx, len(page.current.Accounts), func(i int) {
				acct := page.current.Accounts[i]
				a := func() {
					layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func() {
							common.theme.Body1(acct.Name).Layout(gtx)
						}),
						layout.Rigid(func() {
							common.theme.Body1(acct.TotalBalance).Layout(gtx)
						}),
						layout.Rigid(func() {
							common.theme.Body1("Keys: " + acct.Keys.External + " external, " + acct.Keys.Internal + " internal, " + acct.Keys.Imported + " imported").Layout(gtx)
						}),
						layout.Rigid(func() {
							common.theme.Body1("HD Path: " + acct.HDPath).Layout(gtx)
						}),
						layout.Rigid(func() {
							page.line.Width = gtx.Px(values.AccountLineWidth)
							page.line.Layout(gtx)
						}),
					)
				}
				layout.UniformInset(values.MarginPadding5).Layout(gtx, a)
			})
		},
	}

	page.container.Layout(gtx, len(wdgs), func(i int) {
		layout.Inset{Left: values.MarginPadding5}.Layout(gtx, wdgs[i])
	})
}

func (page *walletPage) bottomRow(common pageCommon) {
	gtx := common.gtx

	if page.walletInfo.Synced || page.walletInfo.Syncing {
		page.icons.addWalletW.Background = common.theme.Color.Hint
	} else {
		page.icons.addWalletW.Background = common.theme.Color.Primary
	}

	layout.UniformInset(values.MarginPadding5).Layout(gtx, func() {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(page.newRow(&common, page.icons.addWalletW, &page.icons.addWallet, "Add wallet")),
			layout.Rigid(page.newRow(&common, page.icons.renameW, &page.icons.rename, "Rename wallet")),
			layout.Rigid(page.newRow(&common, page.icons.signW, &page.icons.sign, "Sign message")),
			layout.Rigid(page.newRow(&common, page.icons.verifyW, &page.icons.verify, "Verify message")),
			layout.Rigid(page.newRow(&common, page.icons.changePassW, &page.icons.changePass, "Change passphrase")),
			layout.Rigid(page.newRow(&common, page.icons.deleteW, &page.icons.delete, "Delete wallet")),
			layout.Rigid(
				func() {
					if len(page.current.Seed) > 0 {
						page.newRow(&common, page.icons.backupW, &page.icons.backup, "Backup Seed")()
					}
				}),
		)
	})
}

func (page *walletPage) subRename(common pageCommon) {
	gtx := common.gtx
	list := layout.List{Axis: layout.Vertical}
	wdgs := []func(){
		func() {
			page.returnBtn(common)
			layout.Inset{Left: values.MarginPadding50}.Layout(gtx, func() {
				common.theme.H5("Rename Wallet").Layout(gtx)
			})
		},
		func() {
			layout.Flex{}.Layout(gtx,
				layout.Rigid(func() {
					layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func() {
						common.theme.Body1("Your are about to rename").Layout(gtx)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func() {
						txt := common.theme.H5(page.current.Name)
						txt.Color = common.theme.Color.Danger
						txt.Layout(gtx)
					})
				}),
			)
		},
		func() {
			inset := layout.Inset{
				Top:    values.MarginPadding20,
				Bottom: values.MarginPadding20,
			}
			inset.Layout(gtx, func() {
				page.editorW.Layout(gtx, &page.editor)
			})
		},
		func() {
			page.renameW.Layout(gtx, &page.rename)
		},
		func() {
			layout.Center.Layout(common.gtx, func() {
				layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func() {
					page.errorLabel.Layout(gtx)
				})
			})
		},
	}
	common.Layout(gtx, func() {
		list.Layout(gtx, len(wdgs), func(i int) {
			wdgs[i]()
		})
	})
}

func (page *walletPage) subDelete(common pageCommon) {
	gtx := common.gtx
	list := layout.List{Axis: layout.Vertical}
	wdgs := []func(){
		func() {
			common.theme.H5("Delete Wallet").Layout(gtx)
		},
		func() {
			layout.Flex{}.Layout(gtx,
				layout.Rigid(func() {
					layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func() {
						common.theme.Body1("Are you sure you want to delete ").Layout(gtx)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func() {
						txt := common.theme.H5(page.current.Name)
						txt.Color = common.theme.Color.Danger
						txt.Layout(gtx)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func() {
						common.theme.H5("?").Layout(gtx)
					})
				}),
			)
		},
		func() {
			inset := layout.Inset{
				Top:    values.MarginPadding20,
				Bottom: values.MarginPadding5,
			}
			inset.Layout(gtx, func() {
				page.cancelDeleteW.Layout(gtx, &page.icons.main)
			})
		},
		func() {
			page.deleteW.Layout(gtx, &page.delete)
		},
		func() {
			layout.Center.Layout(common.gtx, func() {
				layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func() {
					page.errorLabel.Layout(gtx)
				})
			})
		},
	}
	common.Layout(gtx, func() {
		list.Layout(gtx, len(wdgs), func(i int) {
			wdgs[i]()
		})
	})
	if page.isPasswordModalOpen {
		page.passwordModal.Layout(gtx, page.confirm, page.cancel)
	}
}

// Handle handles all widget inputs on the main wallets page.
func (page *walletPage) Handle(common pageCommon) {
	gtx := common.gtx
	if common.walletsTab.Changed() || common.navTab.Changed() { // reset everything
		page.subPage = subWalletMain
	}

	// Subs
	if page.icons.main.Clicked(gtx) {
		page.errorLabel.Text = ""
		page.subPage = subWalletMain
		return
	}

	if page.icons.rename.Clicked(gtx) {
		page.subPage = subWalletRename
		return
	}

	if page.icons.addWallet.Clicked(gtx) {
		if !(page.walletInfo.Synced || page.walletInfo.Syncing) {
			*common.page = PageCreateRestore
			return
		}
		if page.errorText == "" {
			page.errorText = "You have to stop sync to create a new wallet"
			time.AfterFunc(time.Second*2, func() {
				page.errorText = ""
			})
		}
		return
	}

	if page.icons.delete.Clicked(gtx) {
		page.subPage = subWalletDelete
		return
	}

	if page.icons.changePass.Clicked(gtx) {
		*common.page = PageWalletPassphrase
		return
	}

	if page.icons.sign.Clicked(gtx) {
		*common.page = PageSignMessage
	}

	if page.icons.verify.Clicked(gtx) {
		*common.page = PageVerifyMessage
		return
	}

	if page.icons.addAcct.Clicked(gtx) {
		*common.page = PageWalletAccounts
		return
	}

	if page.rename.Clicked(gtx) {
		name := page.editor.Text()
		if name == "" {
			return
		}

		err := common.wallet.RenameWallet(page.current.ID, name)
		if err != nil {
			log.Error(err)
			page.errorLabel.Text = err.Error()
			return
		}

		common.info.Wallets[*common.selectedWallet].Name = name
		page.subPage = subWalletMain
	}

	if page.editor.Text() == "" {
		page.renameW.Background = common.theme.Color.Hint
	} else {
		page.renameW.Background = common.theme.Color.Primary
	}

	if page.delete.Clicked(gtx) {
		page.errorLabel.Text = ""
		page.isPasswordModalOpen = true
	}

	if page.icons.backup.Clicked(gtx) {
		*common.page = PageSeedBackup
	}

	select {
	case err := <-page.errChann:
		if err.Error() == "invalid_passphrase" {
			page.errorLabel.Text = "Wallet passphrase is incorrect."
		} else {
			page.errorLabel.Text = err.Error()
		}
	default:
	}
}

func (page *walletPage) returnBtn(common pageCommon) {
	layout.W.Layout(common.gtx, func() {
		page.icons.mainW.Layout(common.gtx, &page.icons.main)
	})
}

func (page *walletPage) newRow(common *pageCommon, out decredmaterial.IconButton, in *widget.Button, label string) layout.Widget {
	return func() {
		layout.Inset{Right: values.MarginPadding15, Top: values.MarginPadding5}.Layout(common.gtx, func() {
			layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(common.gtx,
				layout.Rigid(func() {
					out.Layout(common.gtx, in)
				}),
				layout.Rigid(func() {
					layout.Inset{Top: values.MarginPadding5}.Layout(common.gtx, func() {
						common.theme.Caption(label).Layout(common.gtx)
					})
				}),
			)
		})
	}
}

func (page *walletPage) confirm(password []byte) {
	page.isPasswordModalOpen = false
	page.wallet.DeleteWallet(page.current.ID, password, page.errChann)
}

func (page *walletPage) cancel() {
	page.isPasswordModalOpen = false
}

func (page *walletPage) alert(common pageCommon) {
	if page.errorText != "" {
		common.theme.ErrorAlert(common.gtx, page.errorText)
	}
}

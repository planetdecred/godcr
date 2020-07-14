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
)

type walletPage struct {
	walletInfo *wallet.MultiWalletInfo
	subPage    int
	current    wallet.InfoShort
	wallet     *wallet.Wallet
	result     **wallet.Signature
	icons      struct {
		main, delete, sign, verify, addWallet, rename,
		changePass, addAcct, backup decredmaterial.IconButton
	}
	container, accountsList         layout.List
	line                            *decredmaterial.Line
	rename, delete, cancelDelete decredmaterial.Button
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
		editorW:       common.theme.Editor(new(widget.Editor), "New wallet name"),
		rename:       common.theme.Button(new(widget.Clickable), "Rename Wallet"),
		errorLabel:    common.theme.Body2(""),
		result:        &win.signatureResult,
		delete:       common.theme.DangerButton(new(widget.Clickable), "Confirm Delete Wallet"),
		cancelDelete: common.theme.Button(new(widget.Clickable), "Cancel Wallet Delete"),
		passwordModal: common.theme.Password(),
		errChann:      common.errorChannels[PageWallet],
		errorText:     "",
	}
	page.line.Color = common.theme.Color.Gray
	page.line.Height = 1
	page.errorLabel.Color = common.theme.Color.Danger

	var iconPadding = values.MarginPadding5
	var iconSize = values.MarginPadding30

	page.icons.addAcct = common.theme.IconButton(new(widget.Clickable), common.icons.contentAdd)
	page.icons.addAcct.Inset = layout.UniformInset(iconPadding)
	page.icons.addAcct.Size = iconSize
	page.icons.main = common.theme.IconButton(new(widget.Clickable), common.icons.navigationArrowBack)
	page.icons.main.Background = color.RGBA{}
	page.icons.main.Color = common.theme.Color.Hint
	page.icons.main.Inset = layout.UniformInset(iconPadding)
	page.icons.main.Size = iconSize
	page.icons.delete = common.theme.IconButton(new(widget.Clickable), common.icons.actionDelete)
	page.icons.delete.Size = iconSize
	page.icons.delete.Inset = layout.UniformInset(iconPadding)
	page.icons.delete.Background = common.theme.Color.Danger
	page.icons.sign = common.theme.IconButton(new(widget.Clickable), common.icons.communicationComment)
	page.icons.sign.Size = iconSize
	page.icons.sign.Inset = layout.UniformInset(iconPadding)
	page.icons.verify = common.theme.IconButton(new(widget.Clickable), common.icons.verifyAction)
	page.icons.verify.Size = iconSize
	page.icons.verify.Inset = layout.UniformInset(iconPadding)
	page.icons.addWallet = common.theme.IconButton(new(widget.Clickable), common.icons.contentAdd)
	page.icons.addWallet.Size = iconSize
	page.icons.addWallet.Inset = layout.UniformInset(iconPadding)
	page.icons.rename = common.theme.IconButton(new(widget.Clickable), common.icons.editorModeEdit)
	page.icons.rename.Size = iconSize
	page.icons.rename.Inset = layout.UniformInset(iconPadding)
	page.icons.changePass = common.theme.IconButton(new(widget.Clickable), common.icons.actionLock)
	page.icons.changePass.Size = iconSize
	page.icons.changePass.Inset = layout.UniformInset(iconPadding)
	page.icons.backup = common.theme.IconButton(new(widget.Clickable), common.icons.actionBackup)
	page.icons.backup.Size = iconSize
	page.icons.backup.Inset = layout.UniformInset(iconPadding)

	return func(gtx C) D {
		page.Handle(common)
		if page.isPasswordModalOpen {
			page.passwordModal.Layout(gtx, page.confirm, page.cancel)
		}
		return page.Layout(common)
	}
}

// Layout lays out the widgets for the main wallets page.
func (page *walletPage) Layout(common pageCommon) layout.Dimensions {
	if common.states.deleted {
		page.subPage = subWalletMain
		common.states.deleted = false
	}

	switch page.subPage {
	case subWalletMain:
		return page.subMain(common)
	case subWalletRename:
		return page.subRename(common)
	case subWalletDelete:
		return page.subDelete(common)
	}
	return page.subMain(common)
}

func (page *walletPage) subMain(common pageCommon) layout.Dimensions {
	gtx := common.gtx

	page.current = common.info.Wallets[*common.selectedWallet]

	body := func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Flexed(0.88, func(gtx C) D {
							return page.topRow(common)
						}),
						layout.Flexed(0.12, func(gtx C) D {
							return page.bottomRow(common)
						}),
					)
				})
			}),
		)
	}

	return common.LayoutWithWallets(gtx, body)
}

func (page *walletPage) topRow(common pageCommon) layout.Dimensions {
	gtx := common.gtx
	wdgs := []func(gtx C) D{
		func(gtx C) D {
			return page.alert(common)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return common.theme.H5(page.current.Name).Layout(common.gtx)
				}),
			)
		},
		func(gtx C) D {
			return common.theme.H6("Total Balance: " + page.current.Balance).Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return common.theme.H6("Accounts").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding5}.Layout(common.gtx, func(gtx C) D {
						return page.icons.addAcct.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return page.accountsList.Layout(gtx, len(page.current.Accounts), func(gtx C, i int) D {
				acct := page.current.Accounts[i]
				a := func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return common.theme.Body1(acct.Name).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return common.theme.Body1(acct.TotalBalance).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return common.theme.Body1("Keys: " + acct.Keys.External + " external, " + acct.Keys.Internal + " internal, " + acct.Keys.Imported + " imported").Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							return common.theme.Body1("HD Path: " + acct.HDPath).Layout(gtx)
						}),
						layout.Rigid(func(gtx C) D {
							page.line.Width = gtx.Px(values.AccountLineWidth)
							return page.line.Layout(gtx)
						}),
					)
				}
				return layout.UniformInset(values.MarginPadding5).Layout(gtx, a)
			})
		},
	}

	return page.container.Layout(gtx, len(wdgs), func(gtx C, i int) D {
		return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, wdgs[i])
	})
}

func (page *walletPage) bottomRow(common pageCommon) layout.Dimensions {
	gtx := common.gtx

	if page.walletInfo.Synced || page.walletInfo.Syncing {
		page.icons.addWallet.Background = common.theme.Color.Hint
	} else {
		page.icons.addWallet.Background = common.theme.Color.Primary
	}

	return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(page.newRow(&common, page.icons.addWallet, "Add wallet")),
			layout.Rigid(page.newRow(&common, page.icons.rename, "Rename wallet")),
			layout.Rigid(page.newRow(&common, page.icons.sign, "Sign message")),
			layout.Rigid(page.newRow(&common, page.icons.verify, "Verify message")),
			layout.Rigid(page.newRow(&common, page.icons.changePass, "Change passphrase")),
			layout.Rigid(page.newRow(&common, page.icons.delete, "Delete wallet")),
			layout.Rigid(
				func(gtx C) D {
					if len(page.current.Seed) > 0 {
						return page.newRow(&common, page.icons.backup, "Backup Seed")(gtx)
					}
					return layout.Dimensions{}
				}),
		)
	})
}

func (page *walletPage) subRename(common pageCommon) layout.Dimensions {
	gtx := common.gtx
	list := layout.List{Axis: layout.Vertical}
	wdgs := []func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return page.returnBtn(common)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding50}.Layout(gtx, func(gtx C) D {
						return common.theme.H5("Rename Wallet").Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return common.theme.Body1("Your are about to rename").Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						txt := common.theme.H5(page.current.Name)
						txt.Color = common.theme.Color.Danger
						return txt.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			inset := layout.Inset{
				Top:    values.MarginPadding20,
				Bottom: values.MarginPadding20,
			}
			return inset.Layout(gtx, func(gtx C) D {
				return page.editorW.Layout(gtx)
			})
		},
		func(gtx C) D {
			return page.rename.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Center.Layout(common.gtx, func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
					return page.errorLabel.Layout(gtx)
				})
			})
		},
	}
	return common.Layout(gtx, func(gtx C) D {
		return list.Layout(gtx, len(wdgs), func(gtx C, i int) D {
			return wdgs[i](gtx)
		})
	})
}

func (page *walletPage) subDelete(common pageCommon) layout.Dimensions {
	gtx := common.gtx
	list := layout.List{Axis: layout.Vertical}
	wdgs := []func(gtx C) D{
		func(gtx C) D {
			return common.theme.H5("Delete Wallet").Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
						return common.theme.Body1("Are you sure you want to delete ").Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						txt := common.theme.H5(page.current.Name)
						txt.Color = common.theme.Color.Danger
						return txt.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return common.theme.H5("?").Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			inset := layout.Inset{
				Top:    values.MarginPadding20,
				Bottom: values.MarginPadding5,
			}
			return inset.Layout(gtx, func(gtx C) D {
				return page.cancelDelete.Layout(gtx)
			})
		},
		func(gtx C) D {
			return page.delete.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Center.Layout(common.gtx, func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
					return page.errorLabel.Layout(gtx)
				})
			})
		},
	}
	return common.Layout(gtx, func(gtx C) D {
		return list.Layout(gtx, len(wdgs), func(gtx C, i int) D {
			return wdgs[i](gtx)
		})
	})
}

// Handle handles all widget inputs on the main wallets page.
func (page *walletPage) Handle(common pageCommon) {
	gtx := common.gtx

	for range common.walletsTab.ChangeEvent(gtx) {
		page.subPage = subWalletMain
	}

	for range common.navTab.ChangeEvent(gtx) {
		page.subPage = subWalletMain
	}

	// Subs
	if page.icons.main.Button.Clicked() || page.cancelDelete.Button.Clicked() {
		page.errorLabel.Text = ""
		page.subPage = subWalletMain
		return
	}

	if page.icons.rename.Button.Clicked() {
		page.subPage = subWalletRename
		return
	}

	if page.icons.addWallet.Button.Clicked() {
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

	if page.icons.delete.Button.Clicked() {
		page.subPage = subWalletDelete
		return
	}

	if page.icons.changePass.Button.Clicked() {
		*common.page = PageWalletPassphrase
		return
	}

	if page.icons.sign.Button.Clicked() {
		*common.page = PageSignMessage
	}

	if page.icons.verify.Button.Clicked() {
		*common.page = PageVerifyMessage
		return
	}

	if page.icons.addAcct.Button.Clicked() {
		*common.page = PageWalletAccounts
		return
	}

	if page.rename.Button.Clicked() {
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
		page.rename.Background = common.theme.Color.Hint
	} else {
		page.rename.Background = common.theme.Color.Primary
	}

	if page.delete.Button.Clicked() {
		page.errorLabel.Text = ""
		page.isPasswordModalOpen = true
	}

	if page.icons.backup.Button.Clicked() {
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

func (page *walletPage) returnBtn(common pageCommon) layout.Dimensions {
	return layout.W.Layout(common.gtx, func(gtx C) D {
		return page.icons.main.Layout(common.gtx)
	})
}

func (page *walletPage) newRow(common *pageCommon, button decredmaterial.IconButton, label string) layout.Widget {
	return func(gtx C) D {
		return layout.Inset{Right: values.MarginPadding15, Top: values.MarginPadding5}.Layout(common.gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(common.gtx,
				layout.Rigid(func(gtx C) D {
					return button.Layout(common.gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding5}.Layout(common.gtx, func(gtx C) D {
						return common.theme.Caption(label).Layout(common.gtx)
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

func (page *walletPage) alert(common pageCommon) layout.Dimensions {
	if page.errorText != "" {
		return common.theme.ErrorAlert(common.gtx, page.errorText)
	}
	return layout.Dimensions{}
}

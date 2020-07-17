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
	container, accountsList      layout.List
	line                         *decredmaterial.Line
	rename, delete, cancelDelete decredmaterial.Button
	errorLabel                   decredmaterial.Label
	editor                       widget.Editor
	editorW                      decredmaterial.Editor
	passwordModal                *decredmaterial.Password
	isPasswordModalOpen          bool
	errChann                     chan error
	errorText                    string
}

func (win *Window) WalletPage(common pageCommon) layout.Widget {
	pg := &walletPage{
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
		rename:        common.theme.Button(new(widget.Clickable), "Rename Wallet"),
		errorLabel:    common.theme.Body2(""),
		result:        &win.signatureResult,
		delete:        common.theme.DangerButton(new(widget.Clickable), "Confirm Delete Wallet"),
		cancelDelete:  common.theme.Button(new(widget.Clickable), "Cancel Wallet Delete"),
		passwordModal: common.theme.Password(),
		errChann:      common.errorChannels[PageWallet],
		errorText:     "",
	}
	pg.line.Color = common.theme.Color.Gray
	pg.line.Height = 1
	pg.errorLabel.Color = common.theme.Color.Danger

	var iconPadding = values.MarginPadding5
	var iconSize = values.MarginPadding30

	pg.icons.addAcct = common.theme.IconButton(new(widget.Clickable), common.icons.contentAdd)
	pg.icons.addAcct.Inset = layout.UniformInset(iconPadding)
	pg.icons.addAcct.Size = iconSize
	pg.icons.main = common.theme.IconButton(new(widget.Clickable), common.icons.navigationArrowBack)
	pg.icons.main.Background = color.RGBA{}
	pg.icons.main.Color = common.theme.Color.Hint
	pg.icons.main.Inset = layout.UniformInset(iconPadding)
	pg.icons.main.Size = iconSize
	pg.icons.delete = common.theme.IconButton(new(widget.Clickable), common.icons.actionDelete)
	pg.icons.delete.Size = iconSize
	pg.icons.delete.Inset = layout.UniformInset(iconPadding)
	pg.icons.delete.Background = common.theme.Color.Danger
	pg.icons.sign = common.theme.IconButton(new(widget.Clickable), common.icons.communicationComment)
	pg.icons.sign.Size = iconSize
	pg.icons.sign.Inset = layout.UniformInset(iconPadding)
	pg.icons.verify = common.theme.IconButton(new(widget.Clickable), common.icons.verifyAction)
	pg.icons.verify.Size = iconSize
	pg.icons.verify.Inset = layout.UniformInset(iconPadding)
	pg.icons.addWallet = common.theme.IconButton(new(widget.Clickable), common.icons.contentAdd)
	pg.icons.addWallet.Size = iconSize
	pg.icons.addWallet.Inset = layout.UniformInset(iconPadding)
	pg.icons.rename = common.theme.IconButton(new(widget.Clickable), common.icons.editorModeEdit)
	pg.icons.rename.Size = iconSize
	pg.icons.rename.Inset = layout.UniformInset(iconPadding)
	pg.icons.changePass = common.theme.IconButton(new(widget.Clickable), common.icons.actionLock)
	pg.icons.changePass.Size = iconSize
	pg.icons.changePass.Inset = layout.UniformInset(iconPadding)
	pg.icons.backup = common.theme.IconButton(new(widget.Clickable), common.icons.actionBackup)
	pg.icons.backup.Size = iconSize
	pg.icons.backup.Inset = layout.UniformInset(iconPadding)

	return func(gtx C) D {
		pg.Handle(common)
		if pg.isPasswordModalOpen {
			pg.passwordModal.Layout(gtx, pg.confirm, pg.cancel)
		}
		return pg.Layout(gtx, common)
	}
}

// Layout lays out the widgets for the main wallets pg.
func (pg *walletPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	if common.states.deleted {
		pg.subPage = subWalletMain
		common.states.deleted = false
	}

	switch pg.subPage {
	case subWalletMain:
		return pg.subMain(gtx, common)
	case subWalletRename:
		return pg.subRename(gtx, common)
	case subWalletDelete:
		return pg.subDelete(gtx, common)
	}
	return pg.subMain(gtx, common)
}

func (pg *walletPage) subMain(gtx layout.Context, common pageCommon) layout.Dimensions {
	pg.current = common.info.Wallets[*common.selectedWallet]

	body := func(gtx C) D {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Flexed(0.88, func(gtx C) D {
							return pg.topRow(gtx, common)
						}),
						layout.Flexed(0.12, func(gtx C) D {
							return pg.bottomRow(gtx, common)
						}),
					)
				})
			}),
		)
	}

	return common.LayoutWithWallets(gtx, body)
}

func (pg *walletPage) topRow(gtx layout.Context, common pageCommon) layout.Dimensions {
	wdgs := []func(gtx C) D{
		func(gtx C) D {
			return pg.alert(gtx, common)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return common.theme.H5(pg.current.Name).Layout(gtx)
				}),
			)
		},
		func(gtx C) D {
			return common.theme.H6("Total Balance: " + pg.current.Balance).Layout(gtx)
		},
		func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return common.theme.H6("Accounts").Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return pg.icons.addAcct.Layout(gtx)
					})
				}),
			)
		},
		func(gtx C) D {
			return pg.accountsList.Layout(gtx, len(pg.current.Accounts), func(gtx C, i int) D {
				acct := pg.current.Accounts[i]
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
							pg.line.Width = gtx.Px(values.AccountLineWidth)
							return pg.line.Layout(gtx)
						}),
					)
				}
				return layout.UniformInset(values.MarginPadding5).Layout(gtx, a)
			})
		},
	}

	return pg.container.Layout(gtx, len(wdgs), func(gtx C, i int) D {
		return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, wdgs[i])
	})
}

func (pg *walletPage) bottomRow(gtx layout.Context, common pageCommon) layout.Dimensions {
	if pg.walletInfo.Synced || pg.walletInfo.Syncing {
		pg.icons.addWallet.Background = common.theme.Color.Hint
	} else {
		pg.icons.addWallet.Background = common.theme.Color.Primary
	}

	return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(pg.newRow(&common, pg.icons.addWallet, "Add wallet")),
			layout.Rigid(pg.newRow(&common, pg.icons.rename, "Rename wallet")),
			layout.Rigid(pg.newRow(&common, pg.icons.sign, "Sign message")),
			layout.Rigid(pg.newRow(&common, pg.icons.verify, "Verify message")),
			layout.Rigid(pg.newRow(&common, pg.icons.changePass, "Change passphrase")),
			layout.Rigid(pg.newRow(&common, pg.icons.delete, "Delete wallet")),
			layout.Rigid(
				func(gtx C) D {
					if len(pg.current.Seed) > 0 {
						return pg.newRow(&common, pg.icons.backup, "Backup Seed")(gtx)
					}
					return layout.Dimensions{}
				}),
		)
	})
}

func (pg *walletPage) subRename(gtx layout.Context, common pageCommon) layout.Dimensions {
	list := layout.List{Axis: layout.Vertical}
	wdgs := []func(gtx C) D{
		func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return pg.returnBtn(gtx)
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
						txt := common.theme.H5(pg.current.Name)
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
				return pg.editorW.Layout(gtx)
			})
		},
		func(gtx C) D {
			return pg.rename.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
					return pg.errorLabel.Layout(gtx)
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

func (pg *walletPage) subDelete(gtx layout.Context, common pageCommon) layout.Dimensions {
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
						txt := common.theme.H5(pg.current.Name)
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
				return pg.cancelDelete.Layout(gtx)
			})
		},
		func(gtx C) D {
			return pg.delete.Layout(gtx)
		},
		func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
					return pg.errorLabel.Layout(gtx)
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

// Handle handles all widget inputs on the main wallets pg.
func (pg *walletPage) Handle(common pageCommon) {
	//if common.walletsTab.Selected != 1 {
	//	pg.subPage = subWalletMain
	//}

	if common.navTab.Selected != 1 {
		pg.subPage = subWalletMain
	}

	// Subs
	if pg.icons.main.Button.Clicked() || pg.cancelDelete.Button.Clicked() {
		pg.errorLabel.Text = ""
		pg.subPage = subWalletMain
		return
	}

	if pg.icons.rename.Button.Clicked() {
		pg.subPage = subWalletRename
		return
	}

	if pg.icons.addWallet.Button.Clicked() {
		if !(pg.walletInfo.Synced || pg.walletInfo.Syncing) {
			*common.page = PageCreateRestore
			return
		}
		if pg.errorText == "" {
			pg.errorText = "You have to stop sync to create a new wallet"
			time.AfterFunc(time.Second*2, func() {
				pg.errorText = ""
			})
		}
		return
	}

	if pg.icons.delete.Button.Clicked() {
		pg.subPage = subWalletDelete
		return
	}

	if pg.icons.changePass.Button.Clicked() {
		*common.page = PageWalletPassphrase
		return
	}

	if pg.icons.sign.Button.Clicked() {
		*common.page = PageSignMessage
	}

	if pg.icons.verify.Button.Clicked() {
		*common.page = PageVerifyMessage
		return
	}

	if pg.icons.addAcct.Button.Clicked() {
		*common.page = PageWalletAccounts
		return
	}

	if pg.rename.Button.Clicked() {
		name := pg.editor.Text()
		if name == "" {
			return
		}

		err := common.wallet.RenameWallet(pg.current.ID, name)
		if err != nil {
			log.Error(err)
			pg.errorLabel.Text = err.Error()
			return
		}

		common.info.Wallets[*common.selectedWallet].Name = name
		pg.subPage = subWalletMain
	}

	if pg.editor.Text() == "" {
		pg.rename.Background = common.theme.Color.Hint
	} else {
		pg.rename.Background = common.theme.Color.Primary
	}

	if pg.delete.Button.Clicked() {
		pg.errorLabel.Text = ""
		pg.isPasswordModalOpen = true
	}

	if pg.icons.backup.Button.Clicked() {
		*common.page = PageSeedBackup
	}

	select {
	case err := <-pg.errChann:
		if err.Error() == "invalid_passphrase" {
			pg.errorLabel.Text = "Wallet passphrase is incorrect."
		} else {
			pg.errorLabel.Text = err.Error()
		}
	default:
	}
}

func (pg *walletPage) returnBtn(gtx layout.Context) layout.Dimensions {
	return layout.W.Layout(gtx, func(gtx C) D {
		return pg.icons.main.Layout(gtx)
	})
}

func (pg *walletPage) newRow(common *pageCommon, button decredmaterial.IconButton, label string) layout.Widget {
	return func(gtx C) D {
		return layout.Inset{Right: values.MarginPadding15, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return button.Layout(gtx)
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
						return common.theme.Caption(label).Layout(gtx)
					})
				}),
			)
		})
	}
}

func (pg *walletPage) confirm(password []byte) {
	pg.isPasswordModalOpen = false
	pg.wallet.DeleteWallet(pg.current.ID, password, pg.errChann)
}

func (pg *walletPage) cancel() {
	pg.isPasswordModalOpen = false
}

func (pg *walletPage) alert(gtx layout.Context, common pageCommon) layout.Dimensions {
	if pg.errorText != "" {
		return common.theme.ErrorAlert(gtx, pg.errorText)
	}
	return layout.Dimensions{}
}

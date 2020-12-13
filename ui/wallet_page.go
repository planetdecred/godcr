package ui

import (
	"reflect"

	"github.com/planetdecred/godcr/ui/values"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/wallet"
)

const PageWallet = "Wallet"

const (
	subWalletMain = iota
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
	passwordModal                *decredmaterial.Password
	isPasswordModalOpen          bool
	errChann                     chan error

	// renameAcctIndex    int
	renameAcctButtons []decredmaterial.IconButton
	// renameAcctEditor   decredmaterial.Editor
	// renameAcctSubmit   decredmaterial.IconButton
	// renameAcctCancel   decredmaterial.IconButton
	renameAcctMinwidth int
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
		errorLabel:    common.theme.Body2(""),
		result:        &win.signatureResult,
		delete:        common.theme.DangerButton(new(widget.Clickable), "Confirm Delete Wallet"),
		cancelDelete:  common.theme.Button(new(widget.Clickable), "Cancel Wallet Delete"),
		passwordModal: common.theme.Password(),
		errChann:      common.errorChannels[PageWallet],
	}
	pg.line.Color = common.theme.Color.Gray
	pg.line.Height = 1
	pg.errorLabel.Color = common.theme.Color.Danger

	var iconPadding = values.MarginPadding5
	var iconSize = values.MarginPadding20

	pg.icons.addAcct = common.theme.IconButton(new(widget.Clickable), common.icons.contentAdd)
	pg.icons.addAcct.Inset = layout.UniformInset(iconPadding)
	pg.icons.addAcct.Size = iconSize
	pg.icons.main = common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack)
	pg.icons.main.Color = common.theme.Color.Hint
	pg.icons.main.Inset = layout.UniformInset(values.MarginPadding0)

	pg.icons.main.Size = values.MarginPadding30
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

	pg.renameAcctMinwidth = 250

	return func(gtx C) D {
		pg.Handle(common)
		return pg.Layout(gtx, common)
	}
}

// Layout lays out the widgets for the main wallets pg.
func (pg *walletPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	if common.states.deleted {
		pg.subPage = subWalletMain
		common.states.deleted = false
	}

	if pg.current.ID != common.info.Wallets[*common.selectedWallet].ID ||
		!reflect.DeepEqual(pg.current.Seed, common.info.Wallets[*common.selectedWallet].Seed) {
		pg.current = common.info.Wallets[*common.selectedWallet]
		pg.renameAcctButtons = nil

		for i := 0; i < len(pg.current.Accounts); i++ {
			btn := common.theme.IconButton(new(widget.Clickable), common.icons.editorModeEdit)
			btn.Size = values.TextSize12
			btn.Color = common.theme.Color.Primary
			btn.Background = common.theme.Color.Surface
			btn.Inset = layout.UniformInset(values.MarginPadding5)
			pg.renameAcctButtons = append(pg.renameAcctButtons, btn)
		}
	}

	var dims layout.Dimensions

	switch pg.subPage {
	case subWalletMain:
		dims = pg.subMain(gtx, common)
	case subWalletDelete:
		dims = pg.subDelete(gtx, common)
	default:
		dims = pg.subMain(gtx, common)
	}

	if pg.isPasswordModalOpen {
		return common.Modal(gtx, dims, pg.passwordModal.Layout(gtx, pg.confirm, pg.cancel))
	}
	return dims
}

func (pg *walletPage) subMain(gtx layout.Context, common pageCommon) layout.Dimensions {
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
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return common.theme.Body1(acct.Name).Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									if pg.renameAcctButtons != nil && acct.Name != "imported" {
										return pg.renameAcctButtons[i].Layout(gtx)
									}
									return layout.Dimensions{}
								}),
							)
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
							pg.line.Width = gtx.Px(values.MarginPadding350)
							return pg.line.Layout(gtx)
						}),
					)
				}
				return layout.UniformInset(values.MarginPadding5).Layout(gtx, a)
			})
		},
	}

	return pg.container.Layout(gtx, len(wdgs), func(gtx C, i int) D {
		return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, wdgs[i])
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
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
				return wdgs[i](gtx)
			})
		})
	})
}

// Handle handles all widget inputs on the main wallets pg.
func (pg *walletPage) Handle(common pageCommon) {
	// Subs
	if pg.icons.main.Button.Clicked() || pg.cancelDelete.Button.Clicked() {
		pg.errorLabel.Text = ""
		pg.subPage = subWalletMain
		return
	}

	if pg.icons.rename.Button.Clicked() {
		go func() {
			common.modalReceiver <- &modalLoad{
				template: RenameWalletTemplate,
				title:    "Rename wallet",
				confirm: func(name string) {
					common.wallet.RenameWallet(pg.current.ID, name, pg.errChann)
					// todo: handle success on page and not in state
				},
				confirmText: "Rename",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
	}

	if pg.icons.addWallet.Button.Clicked() {
		if !(pg.walletInfo.Synced || pg.walletInfo.Syncing) {
			*common.page = PageCreateRestore
			return
		}
		common.Notify("You have to stop sync to create a new wallet", false)
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

	for i := 0; i < len(pg.renameAcctButtons); i++ {
		if pg.renameAcctButtons[i].Button.Clicked() {
			go func() {
				common.modalReceiver <- &modalLoad{
					template: RenameAccountTemplate,
					title:    "Rename account",
					confirm: func(name string) {
						pg.wallet.RenameAccount(pg.current.ID, pg.current.Accounts[i].Number, name, pg.errChann)
						common.info.Wallets[*common.selectedWallet].Accounts[i].Name = name
						pg.current = common.info.Wallets[*common.selectedWallet]
						// todo: handle success on page and not in state
					},
					confirmText: "Rename",
					cancel:      common.closeModal,
					cancelText:  "Cancel",
				}
			}()
			break
		}
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
		common.Notify(pg.errorLabel.Text, false)
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

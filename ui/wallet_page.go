package ui

import (
	"strings"

	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"image/color"

	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageWallet = "Wallet"

type walletPage struct {
	walletInfo    *wallet.MultiWalletInfo
	subPage       int
	current       wallet.InfoShort
	wallet        *wallet.Wallet
	walletAccount **wallet.Account
	theme         *decredmaterial.Theme

	walletIcon                                 *widget.Image
	accountIcon                                *widget.Image
	addAcct                                    decredmaterial.IconButton
	container, accountsList, walletsList, list layout.List
	line                                       *decredmaterial.Line
	txFeeCollapsible                           *decredmaterial.Collapsible

	walletCollapsible []*decredmaterial.Collapsible
	toAcctDetails     []*gesture.Click

	// rename, delete, cancelDelete decredmaterial.Button
	// errorLabel                   decredmaterial.Label
	// walletNameEditor             decredmaterial.Editor
	// passwordModal                *decredmaterial.Password
	// isPasswordModalOpen          bool
	// errChann                     chan error
	// errorText                    string
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

		walletsList: layout.List{
			Axis: layout.Vertical,
		},
		list: layout.List{
			Axis: layout.Vertical,
		},
		theme:            common.theme,
		wallet:           common.wallet,
		txFeeCollapsible: common.theme.Collapsible(new(widget.Clickable)),
		line:             common.theme.Line(),
		walletAccount:    &win.walletAccount,
		// walletNameEditor: common.theme.Editor(new(widget.Editor), "New wallet name"),
		// renameAcctEditor: common.theme.Editor(new(widget.Editor), ""),
		// rename:           common.theme.Button(new(widget.Clickable), "Rename Wallet"),
		// errorLabel:       common.theme.Body2(""),
		// result:           &win.signatureResult,
		// delete:           common.theme.DangerButton(new(widget.Clickable), "Confirm Delete Wallet"),
		// cancelDelete:     common.theme.Button(new(widget.Clickable), "Cancel Wallet Delete"),
		// passwordModal:    common.theme.Password(),
		// errChann:         common.errorChannels[PageWallet],
		// errorText:        "",
		// renameAcctIndex:  -1,
	}
	pg.line.Height = 1

	// init wallet collapse
	for i := 0; i < 20; i++ {
		pg.walletCollapsible = append(pg.walletCollapsible, win.theme.Collapsible(new(widget.Clickable)))
	}
	// pg.walletCollapsible.focusIndex = -1
	// pg.errorLabel.Color = common.theme.Color.Danger

	// var iconPadding = values.MarginPadding5
	// var iconSize = values.MarginPadding20

	// pg.walletNameEditor.Editor.SingleLine = true

	pg.addAcct = common.theme.IconButton(new(widget.Clickable), common.icons.contentAdd)
	pg.addAcct.Inset = layout.UniformInset(values.MarginPadding0)
	pg.addAcct.Size = values.MarginPadding25
	pg.addAcct.Background = color.RGBA{}
	pg.addAcct.Color = common.theme.Color.Text
	// pg.icons.main = common.theme.PlainIconButton(new(widget.Clickable), common.icons.navigationArrowBack)
	// pg.icons.main.Color = common.theme.Color.Hint
	// pg.icons.main.Inset = layout.UniformInset(values.MarginPadding0)

	// pg.icons.main.Size = values.MarginPadding30
	// pg.icons.delete = common.theme.IconButton(new(widget.Clickable), common.icons.actionDelete)
	// pg.icons.delete.Size = iconSize
	// pg.icons.delete.Inset = layout.UniformInset(iconPadding)
	// pg.icons.delete.Background = common.theme.Color.Danger
	// pg.icons.sign = common.theme.IconButton(new(widget.Clickable), common.icons.communicationComment)
	// pg.icons.sign.Size = iconSize
	// pg.icons.sign.Inset = layout.UniformInset(iconPadding)
	// pg.icons.verify = common.theme.IconButton(new(widget.Clickable), common.icons.verifyAction)
	// pg.icons.verify.Size = iconSize
	// pg.icons.verify.Inset = layout.UniformInset(iconPadding)
	// pg.icons.addWallet = common.theme.IconButton(new(widget.Clickable), common.icons.contentAdd)
	// pg.icons.addWallet.Size = iconSize
	// pg.icons.addWallet.Inset = layout.UniformInset(iconPadding)
	// pg.icons.rename = common.theme.IconButton(new(widget.Clickable), common.icons.editorModeEdit)
	// pg.icons.rename.Size = iconSize
	// pg.icons.rename.Inset = layout.UniformInset(iconPadding)
	// pg.icons.changePass = common.theme.IconButton(new(widget.Clickable), common.icons.actionLock)
	// pg.icons.changePass.Size = iconSize
	// pg.icons.changePass.Inset = layout.UniformInset(iconPadding)
	// pg.icons.backup = common.theme.IconButton(new(widget.Clickable), common.icons.actionBackup)
	// pg.icons.backup.Size = iconSize
	// pg.icons.backup.Inset = layout.UniformInset(iconPadding)

	// pg.renameAcctEditor.IsTitleLabel = false
	// pg.renameAcctEditor.Editor.SetText("")
	// pg.renameAcctEditor.Editor.SingleLine = true
	// pg.renameAcctEditor.TextSize = values.MarginPadding15
	// pg.renameAcctSubmit = common.theme.IconButton(new(widget.Clickable), common.icons.navigationCheck)
	// pg.renameAcctSubmit.Size = values.TextSize12
	// pg.renameAcctSubmit.Color = common.theme.Color.Success
	// pg.renameAcctSubmit.Background = common.theme.Color.Surface
	// pg.renameAcctSubmit.Inset = layout.UniformInset(values.MarginPadding5)
	// pg.renameAcctCancel = common.theme.IconButton(new(widget.Clickable), common.icons.contentClear)
	// pg.renameAcctCancel.Size = values.TextSize12
	// pg.renameAcctCancel.Color = common.theme.Color.Danger
	// pg.renameAcctCancel.Background = common.theme.Color.Surface
	// pg.renameAcctCancel.Inset = layout.UniformInset(values.MarginPadding5)
	// pg.renameAcctMinwidth = 250

	return func(gtx C) D {
		pg.Handle(common)
		return pg.Layout(gtx, common)
	}
}

// Layout lays out the widgets for the main wallets pg.
func (pg *walletPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	if common.info.LoadedWallets == 0 {
		return common.Layout(gtx, func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return common.theme.H3("pg.text.noWallet").Layout(gtx)
			})
		})
	}

	pageContent := []func(gtx C) D{
		func(gtx C) D {
			return pg.walletSection(gtx, common)
		},
		func(gtx C) D {
			return pg.watchOnlyWalletSection(gtx, common)
		},
	}

	return common.Layout(gtx, func(gtx C) D {
		return pg.container.Layout(gtx, len(pageContent), func(gtx C, i int) D {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, pageContent[i])
		})
	})
	// if common.states.deleted {
	// 	pg.subPage = subWalletMain
	// 	common.states.deleted = false
	// }

	// if pg.current.ID != common.info.Wallets[*common.selectedWallet].ID ||
	// 	!reflect.DeepEqual(pg.current.Seed, common.info.Wallets[*common.selectedWallet].Seed) {
	// 	pg.current = common.info.Wallets[*common.selectedWallet]
	// 	pg.renameAcctButtons = nil
	// 	pg.renameAcctIndex = -1

	// 	for i := 0; i < len(pg.current.Accounts); i++ {
	// 		btn := common.theme.IconButton(new(widget.Clickable), common.icons.editorModeEdit)
	// 		btn.Size = values.TextSize12
	// 		btn.Color = common.theme.Color.Primary
	// 		btn.Background = common.theme.Color.Surface
	// 		btn.Inset = layout.UniformInset(values.MarginPadding5)
	// 		pg.renameAcctButtons = append(pg.renameAcctButtons, btn)
	// 	}
	// }

	// var dims layout.Dimensions

	// switch pg.subPage {
	// case subWalletMain:
	// 	dims = pg.subMain(gtx, common)
	// case subWalletRename:
	// 	dims = pg.subRename(gtx, common)
	// case subWalletDelete:
	// 	dims = pg.subDelete(gtx, common)
	// default:
	// 	dims = pg.subMain(gtx, common)
	// }

	// if pg.isPasswordModalOpen {
	// 	return common.Modal(gtx, dims, pg.passwordModal.Layout(gtx, pg.confirm, pg.cancel))
	// }
	// return dims
}

func (pg *walletPage) walletSection(gtx layout.Context, common pageCommon) layout.Dimensions {
	pg.walletIcon = &widget.Image{Src: paint.NewImageOp(common.icons.walletIcon)}
	pg.walletIcon.Scale = 0.05

	return pg.walletsList.Layout(gtx, len(common.info.Wallets), func(gtx C, i int) D {
		wn := common.info.Wallets[i].Name
		wb := common.info.Wallets[i].Balance
		accounts := common.info.Wallets[i].Accounts

		collapsibleHeader := func(gtx C) D {
			walName := strings.Title(strings.ToLower(wn))
			walletNameLabel := pg.theme.Body1(walName)
			walletBalLabel := pg.theme.Body1(wb)
			walletBalLabel.Color = pg.theme.Color.Gray
			return pg.tableLayout(gtx, walletNameLabel, walletBalLabel, true)
		}

		collapsibleBody := func(gtx C) D {
			return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				gtx.Constraints.Min.Y = 100

				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.accountsList.Layout(gtx, len(accounts), func(gtx C, i int) D {
							accountsName := accounts[i].Name
							totalBalance := accounts[i].TotalBalance
							spendable := dcrutil.Amount(accounts[i].SpendableBalance).String()

							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return pg.walletAccountsLayout(gtx, accountsName, totalBalance, spendable, common)
								}),
							)
						})
					}),
					// add to line()
					layout.Rigid(func(gtx C) D {
						pg.line.Width = gtx.Constraints.Max.X
						pg.line.Color = common.theme.Color.Background
						m := values.MarginPadding10
						return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
							return pg.line.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Bottom: values.MarginPadding5,
									Right:  values.MarginPadding10,
								}.Layout(gtx, func(gtx C) D {
									return pg.addAcct.Layout(gtx)
								})
							}),
							layout.Rigid(func(gtx C) D {
								txt := pg.theme.H6("Add new account")
								txt.Color = common.theme.Color.Gray
								return txt.Layout(gtx)
							}),
						)
					}),
				)
			})
		}

		return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
			return pg.sectionLayout(gtx, func(gtx C) D {
				inset := layout.Inset{
					Right: values.MarginPadding20,
				}
				return inset.Layout(gtx, func(gtx C) D {
				return pg.walletCollapsible[i].Layout(gtx, collapsibleHeader, collapsibleBody)
				})
			})
		})
	})
}

func (pg *walletPage) watchOnlyWalletSection(gtx layout.Context, common pageCommon) layout.Dimensions {
	return pg.sectionLayout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				txt := pg.theme.Body1("Watch-only Wallets")
				txt.Color = pg.theme.Color.Gray
				return txt.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				pg.line.Width = gtx.Constraints.Max.X
				pg.line.Color = common.theme.Color.Hint
				m := values.MarginPadding10
				inset := layout.Inset{
					Top:    m,
					Bottom: m,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return pg.line.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return pg.theme.H6("coming soon").Layout(gtx)
			}),
		)
	})
}

// func (pg *walletPage) walletAccountLayout(gtx layout.Context, i int, common pageCommon) layout.Dimensions {
// 	wal := common.info.Wallets[i]
// 	return pg.accountsList.Layout(gtx, len(wal.Accounts), func(gtx C, i int) D {
// 		accounts := pg.current.Accounts[i]

// 		spendable := dcrutil.Amount(accounts.SpendableBalance).String()
// 		return pg.accountsTableLayout(gtx, accounts.Name, accounts.TotalBalance, spendable)
// 		// a := func(gtx C) D {
// 		// 	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
// 		// 		layout.Rigid(func(gtx C) D {
// 		// 			if pg.renameAcctIndex == i {
// 		// 				return pg.renameAccountRow(gtx)
// 		// 			}
// 		// 			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
// 		// 				layout.Rigid(func(gtx C) D {
// 		// 					return pg.theme.Body1(accounts.Name).Layout(gtx)
// 		// 				}),
// 		// 				layout.Rigid(func(gtx C) D {
// 		// 					if pg.renameAcctButtons != nil && accounts.Name != "imported" {
// 		// 						return pg.renameAcctButtons[i].Layout(gtx)
// 		// 					}
// 		// 					return layout.Dimensions{}
// 		// 				}),
// 		// 			)
// 		// 		}),
// 		// 		layout.Rigid(func(gtx C) D {
// 		// 			return pg.theme.Body1(accounts.TotalBalance).Layout(gtx)
// 		// 		}),
// 		// 		layout.Rigid(func(gtx C) D {
// 		// 			return pg.theme.Body1("Keys: " + accounts.Keys.External + " external, " + accounts.Keys.Internal + " internal, " + accounts.Keys.Imported + " imported").Layout(gtx)
// 		// 		}),
// 		// 		layout.Rigid(func(gtx C) D {
// 		// 			return pg.theme.Body1("HD Path: " + accounts.HDPath).Layout(gtx)
// 		// 		}),
// 		// 		layout.Rigid(func(gtx C) D {
// 		// 			pg.line.Width = gtx.Px(values.MarginPadding350)
// 		// 			return pg.line.Layout(gtx)
// 		// 		}),
// 		// 	)
// 		// }
// 		// return layout.UniformInset(values.MarginPadding5).Layout(gtx, a)
// 	})
// }

func (pg *walletPage) tableLayout(gtx layout.Context, leftLabel, rightLabel decredmaterial.Label, isIcon bool) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					if isIcon {
						inset := layout.Inset{
							Right: values.MarginPadding10,
						}
						return inset.Layout(gtx, func(gtx C) D {
							return pg.walletIcon.Layout(gtx)
						})
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(func(gtx C) D {
					return leftLabel.Layout(gtx)
				}),
			)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				inset := layout.Inset{
					Right: values.MarginPadding10,
				}
				return inset.Layout(gtx, func(gtx C) D {
					return rightLabel.Layout(gtx)
				})
			})
		}),
	)
}

func (pg *walletPage) walletAccountsLayout(gtx layout.Context, name, totalBal, spendableBal string, common pageCommon) layout.Dimensions {
	pg.accountIcon = &widget.Image{Src: paint.NewImageOp(common.icons.accountIcon)}
	if name == "imported" {
		pg.accountIcon = &widget.Image{Src: paint.NewImageOp(common.icons.importedAccountIcon)}
	}
	pg.accountIcon.Scale = 0.8

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			pg.line.Width = gtx.Constraints.Max.X
			pg.line.Color = common.theme.Color.Background
			m := values.MarginPadding10
			return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
				return pg.line.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					inset := layout.Inset{
						Right: values.MarginPadding10,
						Top:   values.MarginPadding15,
					}
					return inset.Layout(gtx, func(gtx C) D {
						return pg.accountIcon.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Right: values.MarginPadding10,
							}
							return inset.Layout(gtx, func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										acctName := strings.Title(strings.ToLower(name))
										return pg.theme.H6(acctName).Layout(gtx)
									}),
									layout.Flexed(1, func(gtx C) D {
										return layout.E.Layout(gtx, func(gtx C) D {
											inset := layout.Inset{
												Right: values.MarginPadding10,
											}
											return inset.Layout(gtx, func(gtx C) D {
												return layoutBalance(gtx, totalBal, common)
											})
										})
									}),
								)
							})
						}),
						layout.Rigid(func(gtx C) D {
							inset := layout.Inset{
								Right: values.MarginPadding10,
							}
							return inset.Layout(gtx, func(gtx C) D {
								spendibleLabel := pg.theme.Body2("Spendable")
								spendibleLabel.Color = pg.theme.Color.Gray
								spendibleBalLabel := pg.theme.Body2(spendableBal)
								spendibleBalLabel.Color = pg.theme.Color.Gray
								return pg.tableLayout(gtx, spendibleLabel, spendibleBalLabel, false)
							})
						}),
					)
				}),
			)
		}),
	)
}

// drawlayout wraps the page tx and sync section in a card layout
func (pg *walletPage) sectionLayout(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return decredmaterial.Card{Color: pg.theme.Color.Surface, Rounded: true}.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding20).Layout(gtx, body)
	})
}

// func (pg *transactionsPage) updateToWalletAccountDetailsButtons(walAccts *[]wallet.Account) {
// 	if len(*walAccts) != len(pg.toAcctDetails) {
// 		pg.toAcctDetails = make([]*gesture.Click, len(*walAccts))
// 		for index := range *walTxs {
// 			pg.toAcctDetails[index] = &gesture.Click{}
// 		}
// 	}
// }

// func (pg *transactionsPage) goToAcctDetails(gtx layout.Context, c *pageCommon, txn *wallet.Account, click *gesture.Click) {
// 	for _, e := range click.Events(gtx) {
// 		if e.Type == gesture.TypeClick {
// 			*pg.walA = txn
// 			*c.page = PageTransactionDetails
// 		}
// 	}
// }

// func (pg *walletPage) subMain(gtx layout.Context, common pageCommon) layout.Dimensions {
// 	body := func(gtx C) D {
// 		return layout.Stack{}.Layout(gtx,
// 			layout.Expanded(func(gtx C) D {
// 				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
// 					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
// 						layout.Flexed(0.88, func(gtx C) D {
// 							return pg.topRow(gtx, common)
// 						}),
// 						layout.Flexed(0.12, func(gtx C) D {
// 							return pg.bottomRow(gtx, common)
// 						}),
// 					)
// 				})
// 			}),
// 		)
// 	}

// 	return common.LayoutWithWallets(gtx, body)
// }

// func (pg *walletPage) topRow(gtx layout.Context, common pageCommon) layout.Dimensions {
// 	wdgs := []func(gtx C) D{
// 		func(gtx C) D {
// 			return pg.alert(gtx, common)
// 		},
// 		func(gtx C) D {
// 			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
// 				layout.Rigid(func(gtx C) D {
// 					return common.theme.H5(pg.current.Name).Layout(gtx)
// 				}),
// 			)
// 		},
// 		func(gtx C) D {
// 			return common.theme.H6("Total Balance: " + pg.current.Balance).Layout(gtx)
// 		},
// 		func(gtx C) D {
// 			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
// 				layout.Rigid(func(gtx C) D {
// 					return common.theme.H6("Accounts").Layout(gtx)
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
// 						return pg.icons.addAcct.Layout(gtx)
// 					})
// 				}),
// 			)
// 		},
// 		func(gtx C) D {
// 			return pg.accountsList.Layout(gtx, len(pg.current.Accounts), func(gtx C, i int) D {
// 				accounts := pg.current.Accounts[i]
// 				a := func(gtx C) D {
// 					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
// 						layout.Rigid(func(gtx C) D {
// 							if pg.renameAcctIndex == i {
// 								return pg.renameAccountRow(gtx)
// 							}
// 							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
// 								layout.Rigid(func(gtx C) D {
// 									return common.theme.Body1(accounts.Name).Layout(gtx)
// 								}),
// 								layout.Rigid(func(gtx C) D {
// 									if pg.renameAcctButtons != nil && accounts.Name != "imported" {
// 										return pg.renameAcctButtons[i].Layout(gtx)
// 									}
// 									return layout.Dimensions{}
// 								}),
// 							)
// 						}),
// 						layout.Rigid(func(gtx C) D {
// 							return common.theme.Body1(accounts.TotalBalance).Layout(gtx)
// 						}),
// 						layout.Rigid(func(gtx C) D {
// 							return common.theme.Body1("Keys: " + accounts.Keys.External + " external, " + accounts.Keys.Internal + " internal, " + accounts.Keys.Imported + " imported").Layout(gtx)
// 						}),
// 						layout.Rigid(func(gtx C) D {
// 							return common.theme.Body1("HD Path: " + accounts.HDPath).Layout(gtx)
// 						}),
// 						layout.Rigid(func(gtx C) D {
// 							pg.line.Width = gtx.Px(values.MarginPadding350)
// 							return pg.line.Layout(gtx)
// 						}),
// 					)
// 				}
// 				return layout.UniformInset(values.MarginPadding5).Layout(gtx, a)
// 			})
// 		},
// 	}

// 	return pg.container.Layout(gtx, len(wdgs), func(gtx C, i int) D {
// 		return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, wdgs[i])
// 	})
// }

// func (pg *walletPage) bottomRow(gtx layout.Context, common pageCommon) layout.Dimensions {
// 	if pg.walletInfo.Synced || pg.walletInfo.Syncing {
// 		pg.icons.addWallet.Background = common.theme.Color.Hint
// 	} else {
// 		pg.icons.addWallet.Background = common.theme.Color.Primary
// 	}

// 	return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
// 		return layout.Flex{}.Layout(gtx,
// 			layout.Rigid(pg.newRow(&common, pg.icons.addWallet, "Add wallet")),
// 			layout.Rigid(pg.newRow(&common, pg.icons.rename, "Rename wallet")),
// 			layout.Rigid(pg.newRow(&common, pg.icons.sign, "Sign message")),
// 			layout.Rigid(pg.newRow(&common, pg.icons.verify, "Verify message")),
// 			layout.Rigid(pg.newRow(&common, pg.icons.changePass, "Change passphrase")),
// 			layout.Rigid(pg.newRow(&common, pg.icons.delete, "Delete wallet")),
// 			layout.Rigid(
// 				func(gtx C) D {
// 					if len(pg.current.Seed) > 0 {
// 						return pg.newRow(&common, pg.icons.backup, "Backup Seed")(gtx)
// 					}
// 					return layout.Dimensions{}
// 				}),
// 		)
// 	})
// }

// func (pg *walletPage) subRename(gtx layout.Context, common pageCommon) layout.Dimensions {
// 	list := layout.List{Axis: layout.Vertical}
// 	wdgs := []func(gtx C) D{
// 		func(gtx C) D {
// 			return layout.Flex{}.Layout(gtx,
// 				layout.Rigid(func(gtx C) D {
// 					return pg.returnBtn(gtx)
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return layout.Inset{Left: values.MarginPadding20}.Layout(gtx, func(gtx C) D {
// 						return common.theme.H5("Rename Wallet").Layout(gtx)
// 					})
// 				}),
// 			)
// 		},
// 		func(gtx C) D {
// 			return layout.Flex{}.Layout(gtx,
// 				layout.Rigid(func(gtx C) D {
// 					return layout.Inset{Top: values.TextSize12}.Layout(gtx, func(gtx C) D {
// 						return common.theme.Body1("Your are about to rename").Layout(gtx)
// 					})
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return layout.Inset{Left: values.MarginPadding5, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
// 						txt := common.theme.H5(pg.current.Name)
// 						txt.Color = common.theme.Color.Danger
// 						return txt.Layout(gtx)
// 					})
// 				}),
// 			)
// 		},
// 		func(gtx C) D {
// 			m := values.MarginPadding20
// 			inset := layout.Inset{
// 				Top:    m,
// 				Bottom: m,
// 			}
// 			return inset.Layout(gtx, func(gtx C) D {
// 				return pg.walletNameEditor.Layout(gtx)
// 			})
// 		},
// 		func(gtx C) D {
// 			return pg.rename.Layout(gtx)
// 		},
// 		func(gtx C) D {
// 			return layout.Center.Layout(gtx, func(gtx C) D {
// 				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
// 					return pg.errorLabel.Layout(gtx)
// 				})
// 			})
// 		},
// 	}
// 	return common.Layout(gtx, func(gtx C) D {
// 		return list.Layout(gtx, len(wdgs), func(gtx C, i int) D {
// 			return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
// 				return wdgs[i](gtx)
// 			})
// 		})
// 	})
// }

// func (pg *walletPage) subDelete(gtx layout.Context, common pageCommon) layout.Dimensions {
// 	list := layout.List{Axis: layout.Vertical}
// 	wdgs := []func(gtx C) D{
// 		func(gtx C) D {
// 			return common.theme.H5("Delete Wallet").Layout(gtx)
// 		},
// 		func(gtx C) D {
// 			return layout.Flex{}.Layout(gtx,
// 				layout.Rigid(func(gtx C) D {
// 					return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
// 						return common.theme.Body1("Are you sure you want to delete ").Layout(gtx)
// 					})
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
// 						txt := common.theme.H5(pg.current.Name)
// 						txt.Color = common.theme.Color.Danger
// 						return txt.Layout(gtx)
// 					})
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
// 						return common.theme.H5("?").Layout(gtx)
// 					})
// 				}),
// 			)
// 		},
// 		func(gtx C) D {
// 			inset := layout.Inset{
// 				Top:    values.MarginPadding20,
// 				Bottom: values.MarginPadding5,
// 			}
// 			return inset.Layout(gtx, func(gtx C) D {
// 				return pg.cancelDelete.Layout(gtx)
// 			})
// 		},
// 		func(gtx C) D {
// 			return pg.delete.Layout(gtx)
// 		},
// 		func(gtx C) D {
// 			return layout.Center.Layout(gtx, func(gtx C) D {
// 				return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
// 					return pg.errorLabel.Layout(gtx)
// 				})
// 			})
// 		},
// 	}
// 	return common.Layout(gtx, func(gtx C) D {
// 		return list.Layout(gtx, len(wdgs), func(gtx C, i int) D {
// 			return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
// 				return wdgs[i](gtx)
// 			})
// 		})
// 	})
// }

// Handle handles all widget inputs on the main wallets pg.
func (pg *walletPage) Handle(common pageCommon) {
	for _, b := range pg.walletCollapsible {
		for b.Button.Clicked() {
			b.IsExpanded = !b.IsExpanded
		}
	}

	// for pg.walletCollapsible[0].Button.Clicked() {
	// 	fmt.Println("clicked")
	// }

	// for i := 0; i < len(pg.walletCollapsible.collapsibles); i++ {
	// 	c := pg.walletCollapsible.collapsibles[i]

	// 	for c.ButtonWidget.Clicked(){
	// 		c.IsExpanded = !c.IsExpanded
	// 	}
	// }

	//if common.walletsTab.Selected != 1 {
	//	pg.subPage = subWalletMain
	//}

	// if common.navTab.Selected != 1 {
	// 	pg.subPage = subWalletMain
	// }

	// Subs
	// if pg.icons.main.Button.Clicked() || pg.cancelDelete.Button.Clicked() {
	// 	pg.errorLabel.Text = ""
	// 	pg.subPage = subWalletMain
	// 	return
	// }

	// if pg.icons.rename.Button.Clicked() {
	// 	pg.subPage = subWalletRename
	// 	return
	// }

	// if pg.icons.addWallet.Button.Clicked() {
	// 	if !(pg.walletInfo.Synced || pg.walletInfo.Syncing) {
	// 		*common.page = PageCreateRestore
	// 		return
	// 	}
	// 	if pg.errorText == "" {
	// 		pg.errorText = "You have to stop sync to create a new wallet"
	// 		time.AfterFunc(time.Second*2, func() {
	// 			pg.errorText = ""
	// 		})
	// 	}
	// 	return
	// }

	// if pg.icons.delete.Button.Clicked() {
	// 	pg.subPage = subWalletDelete
	// 	return
	// }

	// if pg.icons.changePass.Button.Clicked() {
	// 	*common.page = PageWalletPassphrase
	// 	return
	// }

	// if pg.icons.sign.Button.Clicked() {
	// 	*common.page = PageSignMessage
	// }

	// if pg.icons.verify.Button.Clicked() {
	// 	*common.page = PageVerifyMessage
	// 	return
	// }

	// if pg.icons.addAcct.Button.Clicked() {
	// 	*common.page = PageWalletAccounts
	// 	return
	// }

	// for i := 0; i < len(pg.renameAcctButtons); i++ {
	// 	if pg.renameAcctButtons[i].Button.Clicked() {
	// 		pg.renameAcctIndex = i
	// 		pg.renameAcctEditor.Editor.SetText(pg.current.Accounts[i].Name)
	// 		pg.renameAcctEditor.Editor.Move(len(pg.current.Accounts[i].Name))
	// 		break
	// 	}
	// }

	// if pg.renameAcctSubmit.Button.Clicked() {
	// 	pg.errorText = ""
	// 	pg.wallet.RenameAccount(pg.current.ID, pg.current.Accounts[pg.renameAcctIndex].Number, pg.renameAcctEditor.Editor.Text(), pg.errChann)
	// 	common.info.Wallets[*common.selectedWallet].Accounts[pg.renameAcctIndex].Name = pg.renameAcctEditor.Editor.Text()
	// 	pg.current = common.info.Wallets[*common.selectedWallet]
	// 	pg.renameAcctEditor.Editor.SetText("")
	// 	pg.renameAcctIndex = -1
	// }

	// if pg.renameAcctCancel.Button.Clicked() {
	// 	pg.renameAcctIndex = -1
	// 	pg.renameAcctEditor.Editor.SetText("")
	// }

	// if pg.rename.Button.Clicked() {
	// 	name := pg.walletNameEditor.Editor.Text()
	// 	if name == "" {
	// 		return
	// 	}

	// 	common.wallet.RenameWallet(pg.current.ID, name, pg.errChann)

	// 	common.info.Wallets[*common.selectedWallet].Name = name
	// 	pg.subPage = subWalletMain
	// 	pg.walletNameEditor.Editor.SetText("")
	// }

	// if pg.walletNameEditor.Editor.Text() == "" {
	// 	pg.rename.Background = common.theme.Color.Hint
	// } else {
	// 	pg.rename.Background = common.theme.Color.Primary
	// }

	// if pg.delete.Button.Clicked() {
	// 	pg.errorLabel.Text = ""
	// 	pg.isPasswordModalOpen = true
	// }

	// if pg.icons.backup.Button.Clicked() {
	// 	*common.page = PageSeedBackup
	// }

	// select {
	// case err := <-pg.errChann:
	// 	if err.Error() == "invalid_passphrase" {
	// 		pg.errorLabel.Text = "Wallet passphrase is incorrect."
	// 	} else {
	// 		pg.errorLabel.Text = err.Error()
	// 	}
	// 	if pg.subPage == subWalletMain {
	// 		pg.errorText = err.Error()
	// 		time.AfterFunc(time.Millisecond*3500, func() {
	// 			pg.errorText = ""
	// 		})
	// 	}
	// default:
	// }
}

// func (pg *walletPage) returnBtn(gtx layout.Context) layout.Dimensions {
// 	return layout.W.Layout(gtx, func(gtx C) D {
// 		return pg.icons.main.Layout(gtx)
// 	})
// }

// func (pg *walletPage) newRow(common *pageCommon, button decredmaterial.IconButton, label string) layout.Widget {
// 	return func(gtx C) D {
// 		return layout.Inset{Right: values.MarginPadding15, Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
// 			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
// 				layout.Rigid(func(gtx C) D {
// 					return button.Layout(gtx)
// 				}),
// 				layout.Rigid(func(gtx C) D {
// 					return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
// 						return common.theme.Caption(label).Layout(gtx)
// 					})
// 				}),
// 			)
// 		})
// 	}
// }

// func (pg *walletPage) renameAccountRow(gtx layout.Context) layout.Dimensions {
// 	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
// 		layout.Rigid(func(gtx C) D {
// 			{ // Calculate editor width using and estimated input width.
// 				editorAutoWidth := 17 * pg.renameAcctEditor.Editor.Len()
// 				pageWidth := gtx.Constraints.Max.X
// 				maxWidth := pageWidth - pg.renameAcctMinwidth
// 				gtx.Constraints.Max.X = pg.renameAcctMinwidth
// 				if editorAutoWidth >= pg.renameAcctMinwidth {
// 					if editorAutoWidth > maxWidth {
// 						gtx.Constraints.Max.X = maxWidth
// 					} else {
// 						gtx.Constraints.Max.X = editorAutoWidth
// 					}
// 				}
// 			}
// 			return pg.renameAcctEditor.Layout(gtx)
// 		}),
// 		layout.Rigid(func(gtx C) D {
// 			return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
// 				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
// 					layout.Rigid(func(gtx C) D {
// 						return pg.renameAcctSubmit.Layout(gtx)
// 					}),
// 					layout.Rigid(func(gtx C) D {
// 						return pg.renameAcctCancel.Layout(gtx)
// 					}),
// 				)
// 			})
// 		}),
// 	)
// }

// func (pg *walletPage) confirm(password []byte) {
// 	pg.isPasswordModalOpen = false
// 	pg.wallet.DeleteWallet(pg.current.ID, password, pg.errChann)
// }

// func (pg *walletPage) cancel() {
// 	pg.isPasswordModalOpen = false
// }

// func (pg *walletPage) alert(gtx layout.Context, common pageCommon) layout.Dimensions {
// 	if pg.errorText != "" {
// 		return common.theme.ErrorAlert(gtx, pg.errorText)
// 	}
// 	return layout.Dimensions{}
// }

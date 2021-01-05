package ui

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image/color"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

const PageSettings = "Settings"

type settingsPage struct {
	theme      *decredmaterial.Theme
	walletInfo *wallet.MultiWalletInfo
	wal        *wallet.Wallet

	changePass   decredmaterial.IconButton
	setupCoin    decredmaterial.IconButton
	rescan       decredmaterial.IconButton
	deleteWallet decredmaterial.IconButton

	notificationW *widget.Bool
	chevronIcon   *widget.Icon
	line          *decredmaterial.Line
	errChann      chan error
	result        **wallet.UnlockWallet
}

func (win *Window) SettingsPage(common pageCommon) layout.Widget {
	pg := &settingsPage{
		theme:         common.theme,
		walletInfo:    win.walletInfo,
		wal:           common.wallet,
		notificationW: new(widget.Bool),
		line:          common.theme.Line(),
		result:        &win.unlockWalletResult,
		errChann:      common.errorChannels[PageSettings],

		changePass: decredmaterial.IconButton{
			material.IconButtonStyle{
				Icon:       common.icons.chevronRight,
				Size:       values.MarginPadding25,
				Background: color.RGBA{},
				Color:      common.theme.Color.LightGray,
				Inset:      layout.UniformInset(values.MarginPadding0),
				Button:     new(widget.Clickable),
			},
		},
		setupCoin: decredmaterial.IconButton{
			material.IconButtonStyle{
				Icon:       common.icons.chevronRight,
				Size:       values.MarginPadding25,
				Background: color.RGBA{},
				Color:      common.theme.Color.LightGray,
				Inset:      layout.UniformInset(values.MarginPadding0),
				Button:     new(widget.Clickable),
			},
		},
		rescan: decredmaterial.IconButton{
			material.IconButtonStyle{
				Icon:       common.icons.chevronRight,
				Size:       values.MarginPadding25,
				Background: color.RGBA{},
				Color:      common.theme.Color.LightGray,
				Inset:      layout.UniformInset(values.MarginPadding25),
				Button:     new(widget.Clickable),
			},
		},
		deleteWallet: decredmaterial.IconButton{
			material.IconButtonStyle{
				Icon:       common.icons.chevronRight,
				Size:       values.MarginPadding25,
				Background: color.RGBA{},
				Color:      common.theme.Color.LightGray,
				Inset:      layout.UniformInset(values.MarginPadding25),
				Button:     new(widget.Clickable),
			},
		},
	}
	pg.line.Height = 2
	pg.line.Color = common.theme.Color.Background

	pg.chevronIcon = common.icons.chevronRight
	pg.chevronIcon.Color = common.theme.Color.Background

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

// main settings layout
func (pg *settingsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	body := func(gtx C) D {
		page := SubPage{
			title:      "Settings",
			walletName: common.info.Wallets[*common.selectedWallet].Name,
			back: func() {
				*common.page = PageWallet
			},
			body: func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(values.MarginPadding5).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(pg.changePassphrase()),
						layout.Rigid(pg.notification()),
						// layout.Rigid(pg.editors(pg.messageEditor)),
						// layout.Rigid(pg.drawButtonsRow()),
						// layout.Rigid(pg.drawResult()),
					)
				})
			},
			infoTemplate: "",
		}
		return common.SubPageLayout(gtx, page)
	}

	return common.Layout(gtx, body)
}

func (pg *settingsPage) changePassphrase() layout.Widget {
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				txt := pg.theme.Body2("Spending password")
				txt.Color = pg.theme.Color.Gray
				return txt.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.theme.H6("Change spending password").Layout(gtx)
					}),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return pg.changePass.Layout(gtx)
						})
					}),
				)
			}),
		)
	}
}

func (pg *settingsPage) notification() layout.Widget {
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return pg.drawLine(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				txt := pg.theme.Body2("Notification")
				txt.Color = pg.theme.Color.Gray
				return txt.Layout(gtx)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return pg.theme.H6("Incoming transactions").Layout(gtx)
					}),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return pg.theme.Switch(pg.notificationW).Layout(gtx)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return pg.drawLine(gtx)
			}),
		)
	}
}

func (pg *settingsPage) drawLine(gtx layout.Context) layout.Dimensions {
	m := values.MarginPadding10
	pg.line.Width = gtx.Constraints.Max.X
	return layout.Inset{Top: m, Bottom: m}.Layout(gtx, func(gtx C) D {
		return pg.line.Layout(gtx)
	})
}

// func (pg *settingsPage) drawRows(HeaderText, labelText string, isIconBtn bool) layout.Dimensions {
// 	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
// 		layout.Rigid(func(gtx C) D {
// 			return common.theme.Body2(HeaderText).Layout(gtx)
// 		}),
// 		layout.Rigid(func(gtx C) D {
// 			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
// 				layout.Rigid(func(gtx C) D {
// 					// inset := layout.Inset{
// 					// 	Right: values.MarginPadding5,
// 					// }
// 					// return inset.Layout(gtx, func(gtx C) D {
// 						return common.theme.Body2(HeaderText).Layout(gtx)
// 					// })
// 				}),
// 				layout.Flex(1,func(gtx C) D {
// 					return layout.E.Layout(gtx, func(gtx C) D {
// 						return pg.signButton.Layout(gtx)
// 					})
// 				}),
// 			)
// 		}),
// 	)
// }

func (pg *settingsPage) handle(common pageCommon) {
	for pg.changePass.Button.Clicked() {
		go func() {
			walletID := pg.walletInfo.Wallets[*common.selectedWallet].ID
			common.modalReceiver <- &modalLoad{
				template: PasswordTemplate,
				title:    "Confirm to change",
				confirm: func(password string) {
					pg.wal.UnlockWallet(walletID, []byte(password), pg.errChann)
				},
				confirmText: "Confirm",
				cancel:      common.closeModal,
				cancelText:  "Cancel",
			}
		}()
		break
	}

	if *pg.result != nil {
		if (*pg.result).Err == nil {
			fmt.Println(string((*pg.result).Pass))
			walletID := pg.walletInfo.Wallets[*common.selectedWallet].ID
			oldPassword := (*pg.result).Pass
			go func() {
				common.modalReceiver <- &modalLoad{
					template: ChangePasswordTemplate,
					title:    "Change spending password",
					confirm: func(newPassword string) {
						fmt.Println(string(oldPassword))
						pg.wal.ChangeWalletPassphrase(walletID, string(oldPassword), newPassword, pg.errChann)
					},
					confirmText: "Change",
					cancel:      common.closeModal,
					cancelText:  "Cancel",
				}
			}()
		}
		*pg.result = nil
	}
}

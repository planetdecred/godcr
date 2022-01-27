package dexclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const DexSettingsPageID = "DexSettings"

type DexSettingsPage struct {
	*load.Load
	pageContainer        layout.List
	backButton           decredmaterial.IconButton
	accountHandlerBtns   map[string]*accountHandlerButton
	addDexBtn            decredmaterial.Button
	importAccountBtn     decredmaterial.Button
	changeAppPasswordBtn decredmaterial.Button
}

type accountHandlerButton struct {
	exportAccountBtn  *decredmaterial.Clickable
	disableAccountBtn *decredmaterial.Clickable
}

func NewDexSettingsPage(l *load.Load) *DexSettingsPage {
	pg := &DexSettingsPage{
		Load:                 l,
		pageContainer:        layout.List{Axis: layout.Vertical},
		addDexBtn:            l.Theme.Button(strAddADex),
		importAccountBtn:     l.Theme.Button(strImportAccount),
		changeAppPasswordBtn: l.Theme.OutlineButton(strChangeAppPassword),
	}
	pg.backButton, _ = components.SubpageHeaderButtons(pg.Load)
	inset := layout.Inset{
		Top:    values.MarginPadding5,
		Bottom: values.MarginPadding5,
		Left:   values.MarginPadding9,
		Right:  values.MarginPadding9,
	}
	pg.addDexBtn.Background = l.Theme.Color.Success
	pg.importAccountBtn.Background = l.Theme.Color.Success
	pg.importAccountBtn.Inset, pg.addDexBtn.Inset = inset, inset
	pg.importAccountBtn.TextSize, pg.addDexBtn.TextSize = values.TextSize14, values.TextSize14

	return pg
}

func (pg *DexSettingsPage) ID() string {
	return DexSettingsPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *DexSettingsPage) OnNavigatedTo() {
	pg.initHandlerAccountBtns()
}

func (pg *DexSettingsPage) Layout(gtx layout.Context) D {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      strDexSetting,
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx layout.Context) D {
				return pg.Theme.Card().Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
						wdgs := []func(gtx C) D{
							pg.exchangesInfoLayout,
							pg.addDexAndImportAccountLayout,
							pg.changeAppPasswordLayout,
						}
						return pg.pageContainer.Layout(gtx, len(wdgs), func(gtx C, i int) D {
							return wdgs[i](gtx)
						})
					})
				})
			},
		}

		return sp.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

func (pg *DexSettingsPage) exchangesInfoLayout(gtx C) D {
	var wdgs []layout.FlexChild
	b := func(btn *decredmaterial.Clickable, labelBtn string) layout.Widget {
		return func(gtx C) D {
			return widget.Border{
				Color:        pg.Theme.Color.Gray2,
				CornerRadius: values.MarginPadding0,
				Width:        values.MarginPadding1,
			}.Layout(gtx, func(gtx C) D {
				return btn.Layout(gtx, func(gtx C) D {
					return layout.Inset{
						Top:    values.MarginPadding4,
						Bottom: values.MarginPadding4,
						Left:   values.MarginPadding8,
						Right:  values.MarginPadding8,
					}.Layout(gtx, pg.Theme.Label(values.MarginPadding12, labelBtn).Layout)
				})
			})
		}
	}

	exchanges := sliceExchanges(pg.Dexc().DEXServers())

	for _, e := range exchanges {
		if pg.accountHandlerBtns[e.Host] == nil {
			continue
		}

		dexServer := e
		exportAccountBtn := pg.accountHandlerBtns[dexServer.Host].exportAccountBtn
		disableAccountBtn := pg.accountHandlerBtns[dexServer.Host].disableAccountBtn
		card := pg.Theme.Card()
		card.Inset = layout.UniformInset(values.MarginPadding16)
		wdgs = append(wdgs, layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Bottom: values.MarginPadding10,
			}.Layout(gtx, func(gtx C) D {
				return widget.Border{
					Color:        pg.Theme.Color.Gray2,
					CornerRadius: values.MarginPadding4,
					Width:        values.MarginPadding1,
				}.Layout(gtx, func(gtx C) D {
					return card.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								dexAddress := fmt.Sprintf(nStrAddressDex, dexServer.Host)
								account := fmt.Sprintf(nStrAccountID, dexServer.AcctID)
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(pg.Theme.Label(values.TextSize12, dexAddress).Layout),
									layout.Rigid(pg.Theme.Label(values.TextSize12, account).Layout),
								)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Top: values.MarginPadding10,
								}.Layout(gtx, func(gtx C) D {
									return layout.Flex{}.Layout(gtx,
										layout.Rigid(b(exportAccountBtn, strExportAccount)),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{
												Left: values.MarginPadding10,
											}.Layout(gtx, b(disableAccountBtn, strDisableAccount))
										}),
									)
								})
							}),
						)
					})
				})
			})
		}))
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, wdgs...)
}

func (pg *DexSettingsPage) addDexAndImportAccountLayout(gtx C) D {
	return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.Theme.Label(values.TextSize14, strDexClientSupportSimultaneous).Layout),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(pg.addDexBtn.Layout),
						layout.Rigid(func(gtx C) D {
							return layout.Inset{Left: values.MarginPadding10}.Layout(gtx, pg.importAccountBtn.Layout)
						}),
					)
				})
			}),
		)
	})
}

func (pg *DexSettingsPage) changeAppPasswordLayout(gtx C) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding10, Bottom: values.MarginPadding10}.Layout(gtx, pg.Theme.Separator().Layout)
		}),
		layout.Rigid(pg.changeAppPasswordBtn.Layout),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{Top: values.MarginPadding10}.Layout(gtx, pg.Theme.Separator().Layout)
		}),
	)
}

func (pg *DexSettingsPage) initHandlerAccountBtns() {
	exchanges := sliceExchanges(pg.Dexc().DEXServers())
	pg.accountHandlerBtns = make(map[string]*accountHandlerButton, len(exchanges))
	clickable := func() *decredmaterial.Clickable {
		cl := pg.Theme.NewClickable(true)
		cl.Radius = decredmaterial.Radius(0)
		return cl
	}

	for _, ex := range exchanges {
		btn := &accountHandlerButton{
			exportAccountBtn:  clickable(),
			disableAccountBtn: clickable(),
		}
		pg.accountHandlerBtns[ex.Host] = btn
	}
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *DexSettingsPage) HandleUserInteractions() {
	for h, btn := range pg.accountHandlerBtns {
		dexHost := h
		if btn.disableAccountBtn.Clicked() {
			modal.NewPasswordModal(pg.Load).
				Title(strDisableAccount).
				Hint(strAppPassword).
				Description(fmt.Sprintf(nStrConfirmDisableAccount, dexHost)).
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButton(strDisableAccount, func(password string, pm *modal.PasswordModal) bool {
					go func() {
						err := pg.Dexc().Core().AccountDisable([]byte(password), dexHost)
						if err != nil {
							pm.SetError(err.Error())
							pm.SetLoading(false)
							return
						}
						pg.initHandlerAccountBtns()
						pm.Dismiss()
						pg.RefreshWindow()
					}()
					return false
				}).Show()
		}

		if btn.exportAccountBtn.Clicked() {
			modal.NewPasswordModal(pg.Load).
				Title(strAuthorizeExport).
				Hint(strAppPassword).
				Description(fmt.Sprintf(nStrConfirmExportAccount, dexHost)).
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButton(strAuthorizeExport, func(password string, pm *modal.PasswordModal) bool {
					go func() {
						errHandler := func(err error) {
							pm.SetError(err.Error())
							pm.SetLoading(false)
						}

						account, err := pg.Dexc().Core().AccountExport([]byte(password), dexHost)
						if err != nil {
							errHandler(err)
							return
						}

						b, err := json.Marshal(account)
						if err != nil {
							errHandler(err)
							return
						}

						fileName := fmt.Sprintf("dcrAccount-%s.json", dexHost)
						file, err := pg.Load.Expl.CreateFile(fileName)
						if err != nil {
							errHandler(err)
							return
						}
						defer func() {
							err := file.Close()
							if err != nil {
								errHandler(err)
								return
							}
						}()

						_, err = file.Write(b)
						if err != nil {
							errHandler(err)
							return
						}

						pm.Dismiss()
						pg.RefreshWindow()
					}()
					return false
				}).Show()
		}
	}

	if pg.addDexBtn.Button.Clicked() {
		newAddDexModal(pg.Load).OnDexAdded(func(_ *core.Exchange) {
			pg.initHandlerAccountBtns()
			pg.RefreshWindow()
		}).Show()
	}

	if pg.importAccountBtn.Button.Clicked() {
		go func() {
			file, err := pg.Load.Expl.ChooseFile("json")
			if err != nil {
				pg.Toast.NotifyError(err.Error())
				return
			}

			defer func() {
				err := file.Close()
				if err != nil {
					return
				}
			}()

			byteValue, err := ioutil.ReadAll(file)
			if err != nil {
				pg.Toast.NotifyError(err.Error())
				return
			}

			var account core.Account
			err = json.Unmarshal(byteValue, &account)

			if err != nil {
				pg.Toast.NotifyError(err.Error())
				return
			}

			modal.NewPasswordModal(pg.Load).
				Title(strAuthorizeImport).
				Hint(strAppPassword).
				Description(strPasswordConfirmAcctImport).
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButton(strAuthorizeImport, func(password string, pm *modal.PasswordModal) bool {
					go func() {
						err = pg.Dexc().Core().AccountImport([]byte(password), account)
						if err != nil {
							pm.SetError(err.Error())
							pm.SetLoading(false)
							return
						}

						pg.initHandlerAccountBtns()
						pm.Dismiss()
						pg.RefreshWindow()
					}()
					return false
				}).Show()
		}()
	}

	if pg.changeAppPasswordBtn.Button.Clicked() {
		modal.NewPasswordModal(pg.Load).
			Title(strCurrentPassword).
			Hint(strAppPassword).
			NegativeButton(values.String(values.StrCancel), func() {}).
			PositiveButton(values.String(values.StrConfirm), func(oldPassword string, pm *modal.PasswordModal) bool {
				go func() {
					pm.SetLoading(false)
					pm.Dismiss()
					modal.NewCreatePasswordModal(pg.Load).
						Title(strNewPassword).
						EnableName(false).
						PasswordHint(strNewPassword).
						ConfirmPasswordHint(strConfirmNewPassword).
						PasswordCreated(func(_, newPassword string, m *modal.CreatePasswordModal) bool {
							go func() {
								err := pg.Dexc().Core().ChangeAppPass([]byte(oldPassword), []byte(newPassword))
								// check if old password error then show previous modal
								if err != nil {
									// TODO: dont know if return in different language
									// find out more
									if strings.Contains(err.Error(), "old password error") {
										m.Dismiss()
										pm.Show()
										pm.SetError(err.Error())
										return
									}
									m.SetError(err.Error())
									m.SetLoading(false)
									return
								}
								pg.Toast.Notify(strSuccessful)
								m.Dismiss()
								pg.RefreshWindow()
							}()
							return false
						}).Show()
				}()
				return false
			}).Show()
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *DexSettingsPage) OnNavigatedFrom() {}

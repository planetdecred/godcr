package dexclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"decred.org/dcrdex/client/core"
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/ncruces/zenity"
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
	exchangesWdg         []*settingExchangeWidget
	addDexBtn            decredmaterial.Button
	importAccountBtn     decredmaterial.Button
	changeAppPasswordBtn decredmaterial.Button
}

type settingExchangeWidget struct {
	dexServer         *core.Exchange
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

	pg.initExchangeWidget()
	return pg
}

func (pg *DexSettingsPage) ID() string {
	return DexSettingsPageID
}

func (pg *DexSettingsPage) OnResume() {
}

func (pg *DexSettingsPage) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      strDexSetting,
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Left:  values.MarginPadding10,
					Right: values.MarginPadding10,
				}.Layout(gtx, func(gtx C) D {
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

	for _, e := range pg.exchangesWdg {
		eWdg := e
		card := pg.Theme.Card()
		card.Border = true
		card.Inset = layout.UniformInset(values.MarginPadding16)
		wdgs = append(wdgs, layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Bottom: values.MarginPadding10,
			}.Layout(gtx, func(gtx C) D {
				return card.Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							dexAddress := fmt.Sprintf(nStrAddressDex, eWdg.dexServer.Host)
							account := fmt.Sprintf(nStrAccountID, eWdg.dexServer.AcctID)
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
									layout.Rigid(b(eWdg.exportAccountBtn, strExportAccount)),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{
											Left: values.MarginPadding10,
										}.Layout(gtx, b(eWdg.disableAccountBtn, strDisableAccount))
									}),
								)
							})
						}),
					)
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

func (pg *DexSettingsPage) initExchangeWidget() {
	pg.exchangesWdg = make([]*settingExchangeWidget, 0)
	exchanges := sliceExchanges(pg.Dexc().DEXServers())
	clickable := func() *decredmaterial.Clickable {
		cl := pg.Theme.NewClickable(true)
		cl.Radius = decredmaterial.Radius(0)
		return cl
	}
	for _, ex := range exchanges {
		ew := &settingExchangeWidget{
			dexServer:         ex,
			exportAccountBtn:  clickable(),
			disableAccountBtn: clickable(),
		}
		pg.exchangesWdg = append(pg.exchangesWdg, ew)
	}
}

func (pg *DexSettingsPage) Handle() {
	for _, eWdg := range pg.exchangesWdg {
		if eWdg.disableAccountBtn.Clicked() {
			dexServer := eWdg.dexServer
			modal.NewPasswordModal(pg.Load).
				Title(strDisableAccount).
				Hint(strAppPassword).
				Description(fmt.Sprintf(nStrConfirmDisableAccount, dexServer.Host)).
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButton(strDisableAccount, func(password string, pm *modal.PasswordModal) bool {
					go func() {
						err := pg.Dexc().Core().AccountDisable([]byte(password), dexServer.Host)
						if err != nil {
							pm.SetError(err.Error())
							pm.SetLoading(false)
							return
						}
						pg.initExchangeWidget()
						pm.Dismiss()
						pg.RefreshWindow()
					}()
					return false
				}).Show()
		}

		if eWdg.exportAccountBtn.Clicked() {
			dexServer := eWdg.dexServer
			modal.NewPasswordModal(pg.Load).
				Title(strAuthorizeExport).
				Hint(strAppPassword).
				Description(fmt.Sprintf(nStrConfirmExportAccount, dexServer.Host)).
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButton(strAuthorizeExport, func(password string, pm *modal.PasswordModal) bool {
					go func() {
						account, err := pg.Dexc().Core().AccountExport([]byte(password), dexServer.Host)
						if err != nil {
							pm.SetError(err.Error())
							pm.SetLoading(false)
							return
						}

						file, err := json.Marshal(account)
						if err != nil {
							pm.SetError(err.Error())
							pm.SetLoading(false)
							return
						}

						fileName := fmt.Sprintf("dcrAccount-%s.json", dexServer.Host)
						filePath, err := zenity.SelectFileSave(
							zenity.Title("Save Your Account"),
							zenity.ConfirmOverwrite(),
							zenity.Filename(fileName),
							zenity.FileFilters{
								zenity.FileFilter{
									Name:     "JSON files",
									Patterns: []string{"*.json"},
								},
							})

						if err != nil {
							pm.SetError(err.Error())
							pm.SetLoading(false)
							return
						}

						err = ioutil.WriteFile(filePath, file, 0644)
						if err != nil {
							pm.SetError(err.Error())
							pm.SetLoading(false)
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
		newAddDexModal(pg.Load).DexCreated(func(_ *core.Exchange) {
			pg.initExchangeWidget()
			pg.RefreshWindow()
		}).Show()
	}

	if pg.importAccountBtn.Button.Clicked() {
		go func() {
			filePath, err := zenity.SelectFile(
				zenity.Title("Select Your Account"),
			)
			if err != nil {
				pg.Toast.NotifyError(err.Error())
				return
			}

			jsonFile, err := os.Open(filePath)
			defer func() {
				err := jsonFile.Close()
				if err != nil {
					return
				}
			}()

			if err != nil {
				pg.Toast.NotifyError(err.Error())
				return
			}

			byteValue, err := ioutil.ReadAll(jsonFile)
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

						pg.initExchangeWidget()
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
								pg.Toast.Notify(strSuccessfully)
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

func (pg *DexSettingsPage) OnClose() {
}

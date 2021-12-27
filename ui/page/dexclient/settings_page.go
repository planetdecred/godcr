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
	exchange          *core.Exchange
	exportAccountBtn  *decredmaterial.Clickable
	disableAccountBtn *decredmaterial.Clickable
}

func NewDexSettingsPage(l *load.Load) *DexSettingsPage {
	pg := &DexSettingsPage{
		Load:                 l,
		pageContainer:        layout.List{Axis: layout.Vertical},
		addDexBtn:            l.Theme.Button("Add a dex"),
		importAccountBtn:     l.Theme.Button("Import Account"),
		changeAppPasswordBtn: l.Theme.OutlineButton("Change App Password"),
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
			Title:      "Settings",
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
							dexAddress := fmt.Sprintf("Address DEX: %s", eWdg.exchange.Host)
							account := fmt.Sprintf("Account ID: %s", eWdg.exchange.AcctID)
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
									layout.Rigid(b(eWdg.exportAccountBtn, "Export Account")),
									layout.Rigid(func(gtx C) D {
										return layout.Inset{
											Left: values.MarginPadding10,
										}.Layout(gtx, b(eWdg.disableAccountBtn, "Disable Account"))
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
			layout.Rigid(func(gtx C) D {
				t := "The Decred DEX Client supports simultaneous use of any number of DEX servers."
				return pg.Theme.Label(values.TextSize14, t).Layout(gtx)
			}),
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
			exchange:          ex,
			exportAccountBtn:  clickable(),
			disableAccountBtn: clickable(),
		}
		pg.exchangesWdg = append(pg.exchangesWdg, ew)
	}
}

func (pg *DexSettingsPage) Handle() {
	for _, eWdg := range pg.exchangesWdg {
		if eWdg.disableAccountBtn.Clicked() {
			exchange := eWdg.exchange
			modal.NewPasswordModal(pg.Load).
				Title("Disable Account").
				Hint("Password").
				Description(fmt.Sprintf("Enter your app password to disable account: %s \n\nNote: This action is irreversible - once an account is disabled it can't be re-enabled.", exchange.Host)).
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButton("Disable Account", func(password string, pm *modal.PasswordModal) bool {
					go func() {
						err := pg.Dexc().Core().AccountDisable([]byte(password), exchange.Host)
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
			exchange := eWdg.exchange
			modal.NewPasswordModal(pg.Load).
				Title("Authorize Export").
				Hint("Password").
				Description(fmt.Sprintf("Enter your app password to confirm Account export for: %s", exchange.Host)).
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButton("Authorize Export", func(password string, pm *modal.PasswordModal) bool {
					go func() {
						account, err := pg.Dexc().Core().AccountExport([]byte(password), exchange.Host)
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

						fileName := fmt.Sprintf("dcrAccount-%s.json", exchange.Host)
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
		newAddDexModal(pg.Load, "").DexCreated(func(dex *core.Exchange) {
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
				Title("Authorize Import").
				Hint("Password").
				Description("Enter your app password to confirm Account import.").
				NegativeButton(values.String(values.StrCancel), func() {}).
				PositiveButton("Authorize Import", func(password string, pm *modal.PasswordModal) bool {
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
			Title("Current Password").
			Hint("Current app password").
			NegativeButton(values.String(values.StrCancel), func() {}).
			PositiveButton(values.String(values.StrConfirm), func(oldPassword string, pm *modal.PasswordModal) bool {
				go func() {
					pm.SetLoading(false)
					pm.Dismiss()
					modal.NewCreatePasswordModal(pg.Load).
						Title("New Password").
						EnableName(false).
						PasswordHint("New password").
						ConfirmPasswordHint("Confirm new password").
						PasswordCreated(func(walletName, newPassword string, m *modal.CreatePasswordModal) bool {
							go func() {
								err := pg.Dexc().Core().ChangeAppPass([]byte(oldPassword), []byte(newPassword))
								// check if old password error then show previous modal
								if err != nil {
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
								pg.Toast.Notify("Change password successfully!")
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

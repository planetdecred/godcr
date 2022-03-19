package privacy

import (
	"context"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/renderers"

	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

const SetupMixerAccountsPageID = "SetupMixerAccounts"

type SetupMixerAccountsPage struct {
	*load.Load

	ctx       context.Context // page context
	ctxCancel context.CancelFunc

	wallet *dcrlibwallet.Wallet

	backButton              decredmaterial.IconButton
	infoButton              decredmaterial.IconButton
	autoSetupClickable      *decredmaterial.Clickable
	manualSetupClickable    *decredmaterial.Clickable
	autoSetupIcon, nextIcon *decredmaterial.Icon
}

func NewSetupMixerAccountsPage(l *load.Load, wallet *dcrlibwallet.Wallet) *SetupMixerAccountsPage {
	pg := &SetupMixerAccountsPage{
		Load:   l,
		wallet: wallet,
	}
	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	pg.autoSetupIcon = decredmaterial.NewIcon(pg.Icons.ActionCheckCircle)
	pg.autoSetupIcon.Color = pg.Theme.Color.Success

	pg.nextIcon = decredmaterial.NewIcon(pg.Icons.NavigationArrowForward)
	pg.nextIcon.Color = pg.Theme.Color.Gray1

	pg.autoSetupClickable = pg.Theme.NewClickable(true)
	pg.manualSetupClickable = pg.Theme.NewClickable(true)

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *SetupMixerAccountsPage) ID() string {
	return SetupMixerAccountsPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *SetupMixerAccountsPage) OnNavigatedTo() {
	pg.ctx, pg.ctxCancel = context.WithCancel(context.TODO())
}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *SetupMixerAccountsPage) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		page := components.SubPage{
			Load:       pg.Load,
			Title:      "Set up needed accounts",
			WalletName: pg.wallet.Name,
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return pg.Theme.Card().Layout(gtx, func(gtx C) D {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
						layout.Flexed(1, func(gtx C) D {
							return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									txt := pg.Theme.Body1("Two dedicated accounts will be set up to use the mixer:")
									txt.Alignment = text.Start
									ic := decredmaterial.NewIcon(pg.Icons.ImageBrightness1)
									ic.Color = pg.Theme.Color.Gray1
									return layout.Inset{Top: values.MarginPadding16, Left: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
										return layout.Flex{Axis: layout.Vertical, Alignment: layout.Start}.Layout(gtx,
											layout.Rigid(txt.Layout),
											layout.Rigid(func(gtx C) D {
												return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
													return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
														layout.Rigid(func(gtx C) D {
															return layout.Inset{Bottom: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
																return ic.Layout(gtx, values.MarginPadding8)
															})
														}),
														layout.Rigid(func(gtx C) D {
															txt2 := `<span style="text-color: grayText2">
														<b>Mixed </b> account will be the outbounding spending account.
													</span>`

															return layout.Inset{
																Left: values.MarginPadding8,
															}.Layout(gtx, renderers.RenderHTML(txt2, pg.Theme).Layout)
														}),
													)
												})
											}),
											layout.Rigid(func(gtx C) D {
												return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
													return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
														layout.Rigid(func(gtx C) D {
															return layout.Inset{Bottom: values.MarginPadding12}.Layout(gtx, func(gtx C) D {
																return ic.Layout(gtx, values.MarginPadding8)
															})
														}),
														layout.Rigid(func(gtx C) D {
															txt3 := `<span style="text-color: grayText2">
													<b>Unmixed </b> account will be the change handling account.
												</span>`

															return layout.Inset{
																Left: values.MarginPadding8,
															}.Layout(gtx, renderers.RenderHTML(txt3, pg.Theme).Layout)
														}),
													)
												})
											}),
										)
									})
								}),
							)
						}),
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return pg.autoSetupClickable.Layout(gtx, pg.autoSetupLayout)
						}),
						layout.Rigid(func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return pg.manualSetupClickable.Layout(gtx, pg.manualSetupLayout)
						}),
					)
				})
			},
		}
		return page.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

func (pg *SetupMixerAccountsPage) autoSetupLayout(gtx C) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return pg.autoSetupIcon.Layout(gtx, values.MarginPadding20)
							}),
							layout.Rigid(func(gtx C) D {
								autoSetupText := pg.Theme.H6("Auto setup")
								txt := pg.Theme.Body2("Create and set up the needed accounts for you.")
								return layout.Inset{
									Left: values.MarginPadding16,
								}.Layout(gtx, func(gtx C) D {
									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(autoSetupText.Layout),
										layout.Rigid(txt.Layout),
									)
								})
							}),
						)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{
									Right: values.MarginPadding4,
									Top:   values.MarginPadding10,
								}.Layout(gtx, func(gtx C) D {
									return pg.nextIcon.Layout(gtx, values.MarginPadding20)
								})
							}),
						)
					}),
				)
			}),
		)
	})
}

func (pg *SetupMixerAccountsPage) manualSetupLayout(gtx C) D {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return layout.UniformInset(values.MarginPadding16).Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(pg.Icons.EditIcon.Layout24dp),
							layout.Rigid(func(gtx C) D {
								autoSetupText := pg.Theme.H6("Manual setup")
								txt := pg.Theme.Body2("For wallets that have enabled privacy before.")
								return layout.Inset{
									Left: values.MarginPadding16,
								}.Layout(gtx, func(gtx C) D {
									return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
										layout.Rigid(autoSetupText.Layout),
										layout.Rigid(txt.Layout),
									)
								})
							}),
						)
					}),
					layout.Rigid(func(gtx C) D {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Flex{}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										return layout.Inset{
											Right: values.MarginPadding4,
											Top:   values.MarginPadding10,
										}.Layout(gtx, func(gtx C) D {
											return pg.nextIcon.Layout(gtx, values.MarginPadding20)
										})
									}),
								)
							}),
						)
					}),
				)
			}),
		)
	})
}

func (pg *SetupMixerAccountsPage) showModalSetupMixerInfo() {
	info := modal.NewInfoModal(pg.Load).
		Title("Set up mixer by creating two needed accounts").
		SetupWithTemplate(modal.SetupMixerInfoTemplate).
		NegativeButton(values.String(values.StrCancel), func() {}).
		PositiveButton("Begin setup", func() {
			pg.showModalSetupMixerAcct()
		})
	pg.ShowModal(info)
}

func (pg *SetupMixerAccountsPage) showModalSetupMixerAcct() {
	accounts, _ := pg.wallet.GetAccountsRaw()
	for _, acct := range accounts.Acc {
		if acct.Name == "mixed" || acct.Name == "unmixed" {
			alert := decredmaterial.NewIcon(decredmaterial.MustIcon(widget.NewIcon(icons.AlertError)))
			alert.Color = pg.Theme.Color.DeepBlue
			info := modal.NewInfoModal(pg.Load).
				Icon(alert).
				Title("Account name is taken").
				Body("There are existing accounts named mixed or unmixed. Please change the name to something else for now. You can change them back after the setup.").
				PositiveButton("Go back & rename", func() {
					pg.PopFragment()
				})
			pg.ShowModal(info)
			return
		}
	}

	modal.NewPasswordModal(pg.Load).
		Title("Confirm to create needed accounts").
		NegativeButton("Cancel", func() {}).
		PositiveButton("Confirm", func(password string, pm *modal.PasswordModal) bool {
			go func() {
				err := pg.wallet.CreateMixerAccounts("mixed", "unmixed", password)
				if err != nil {
					pm.SetError(err.Error())
					pm.SetLoading(false)
					return
				}
				pm.Dismiss()
			}()

			return false
		}).Show()
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *SetupMixerAccountsPage) HandleUserInteractions() {
	if pg.autoSetupClickable.Clicked() {
		go pg.showModalSetupMixerInfo()
	}

	if pg.manualSetupClickable.Clicked() {
		pg.ChangeFragment(NewManualMixerSetupPage(pg.Load, pg.wallet))
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *SetupMixerAccountsPage) OnNavigatedFrom() {
	pg.ctxCancel()
}

package page

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const ValidateAddressPageID = "ValidateAddress"

const (
	none = iota
	valid
	invalid
	notOwned
)

type ValidateAddressPage struct {
	*load.Load
	addressEditor         decredmaterial.Editor
	clearBtn, validateBtn decredmaterial.Button
	stateValidate         int
	walletName            string
	isEnabled             bool
	backButton            decredmaterial.IconButton
}

func NewValidateAddressPage(l *load.Load) *ValidateAddressPage {
	pg := &ValidateAddressPage{
		Load: l,
	}

	pg.backButton, _ = components.SubpageHeaderButtons(l)

	pg.addressEditor = l.Theme.Editor(new(widget.Editor), "Address")
	pg.addressEditor.Editor.SingleLine = true
	pg.addressEditor.Editor.Submit = true

	pg.validateBtn = l.Theme.Button("Validate")
	pg.validateBtn.Font.Weight = text.Medium

	pg.clearBtn = l.Theme.OutlineButton("Clear")
	pg.clearBtn.Font.Weight = text.Medium

	pg.stateValidate = none

	return pg
}

func (pg *ValidateAddressPage) ID() string {
	return ValidateAddressPageID
}

func (pg *ValidateAddressPage) OnResume() {
	pg.addressEditor.Editor.Focus()
}

func (pg *ValidateAddressPage) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      "Validate address",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(pg.addressSection()),
					)
				})
			},
		}
		return sp.Layout(gtx)
	}
	return components.UniformPadding(gtx, body)
}

func (pg *ValidateAddressPage) addressSection() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.description()),
				layout.Rigid(pg.addressEditor.Layout),
				layout.Rigid(pg.actionButtons()),
				layout.Rigid(pg.showDisplayResult()),
			)
		})
	}
}

func (pg *ValidateAddressPage) description() layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		desc := pg.Theme.Caption("Enter an address to validate:")
		desc.Color = pg.Theme.Color.Gray
		return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, desc.Layout)
	}
}

func (pg *ValidateAddressPage) actionButtons() layout.Widget {
	return func(gtx C) D {
		dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: values.MarginPadding15}.Layout(gtx, func(gtx C) D {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, pg.clearBtn.Layout)
							}),
							layout.Rigid(pg.validateBtn.Layout),
						)
					})
				})
			}),
		)
		return dims
	}
}

func (pg *ValidateAddressPage) lineSeparator(gtx layout.Context) layout.Dimensions {
	m := values.MarginPadding10
	return layout.Inset{Top: m, Bottom: m}.Layout(gtx, pg.Theme.Separator().Layout)
}

func (pg *ValidateAddressPage) showDisplayResult() layout.Widget {
	if pg.stateValidate == none {
		return func(gtx C) D {
			return layout.Dimensions{}
		}
	}
	return func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(pg.lineSeparator),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.Inset{Right: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
							if pg.stateValidate == invalid {
								ic := decredmaterial.NewIcon(pg.Icons.NavigationCancel)
								return ic.Layout(gtx, values.MarginPadding25)
							}
							ic := decredmaterial.NewIcon(pg.Icons.ActionCheckCircle)
							return ic.Layout(gtx, values.MarginPadding25)
						})
					}),
					layout.Rigid(func(gtx C) D {
						if pg.stateValidate == invalid {
							txt := pg.Theme.Body1("Invalid Address")
							txt.Color = pg.Theme.Color.Danger
							txt.TextSize = values.TextSize16
							return txt.Layout(gtx)
						}
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								txt := pg.Theme.Body1("Valid address")
								txt.Color = pg.Theme.Color.Success
								txt.TextSize = values.TextSize16
								return txt.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										var text string
										if pg.stateValidate == valid {
											text = "Owned by you in"
										} else {
											text = "Not owned by you"
										}
										txt := pg.Theme.Body1(text)
										txt.TextSize = values.TextSize14
										txt.Color = pg.Theme.Color.Gray
										return txt.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										if pg.stateValidate == valid {
											if components.StringNotEmpty(pg.walletName) {
												return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
													return decredmaterial.Card{
														Color: pg.Theme.Color.LightGray,
													}.Layout(gtx, func(gtx C) D {
														return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
															walletText := pg.Theme.Caption(pg.walletName)
															walletText.Color = pg.Theme.Color.Gray
															return walletText.Layout(gtx)
														})
													})
												})
											}
										}
										return layout.Dimensions{}
									}),
								)
							}),
						)
					}),
				)
			}),
		)
	}
}

func (pg *ValidateAddressPage) pageSections(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
					layout.Rigid(body),
				)
			})
		})
	})
}

func (pg *ValidateAddressPage) Handle() {
	pg.updateButtonColors()

	isSubmit, isChanged := decredmaterial.HandleEditorEvents(pg.addressEditor.Editor)
	if isChanged {
		pg.stateValidate = none
	}

	if (pg.validateBtn.Clicked() || isSubmit) && pg.isEnabled {
		pg.validateAddress()
	}

	if pg.clearBtn.Clicked() {
		pg.clearPage()
	}
}

func (pg *ValidateAddressPage) clearPage() {
	pg.stateValidate = none
	pg.addressEditor.Editor.SetText("")
}

func (pg *ValidateAddressPage) validateAddress() {
	address := pg.addressEditor.Editor.Text()
	pg.addressEditor.SetError("")

	if !components.StringNotEmpty(address) {
		pg.addressEditor.SetError("Please enter a valid address")
		return
	}

	if !pg.WL.MultiWallet.IsAddressValid(address) {
		pg.stateValidate = invalid
		return
	}

	exist, walletName := pg.WL.Wallet.HaveAddress(address)
	if !exist {
		pg.stateValidate = notOwned
		return
	}

	pg.stateValidate = valid
	pg.walletName = walletName
}

func (pg *ValidateAddressPage) updateButtonColors() {
	if !components.StringNotEmpty(pg.addressEditor.Editor.Text()) {
		pg.validateBtn.Background = pg.Theme.Color.Hint
		pg.clearBtn.Color = pg.Theme.Color.Hint
		pg.isEnabled = false
	} else {
		pg.validateBtn.Background = pg.Theme.Color.Primary
		pg.clearBtn.Color = pg.Theme.Color.Primary
		pg.isEnabled = true
	}
}

func (pg *ValidateAddressPage) clearInputs() {
	pg.validateBtn.Background = pg.Theme.Color.Hint
	pg.addressEditor.Editor.SetText("")
}

func (pg *ValidateAddressPage) OnClose() {}

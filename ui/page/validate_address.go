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

	pg.addressEditor = l.Theme.Editor(new(widget.Editor), values.String(values.StrAddress))
	pg.addressEditor.Editor.SingleLine = true
	pg.addressEditor.Editor.Submit = true

	pg.validateBtn = l.Theme.Button(values.String(values.StrValidate))
	pg.validateBtn.Font.Weight = text.Medium

	pg.clearBtn = l.Theme.OutlineButton(values.String(values.StrClear))
	pg.clearBtn.Font.Weight = text.Medium

	pg.stateValidate = none

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *ValidateAddressPage) ID() string {
	return ValidateAddressPageID
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *ValidateAddressPage) OnNavigatedTo() {
	pg.addressEditor.Editor.Focus()
}

// Layout draws the page UI components into the provided C
// to be eventually drawn on screen.
// Part of the load.Page interface.
func (pg *ValidateAddressPage) Layout(gtx C) D {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrValidateAddr),
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
	return func(gtx C) D {
		desc := pg.Theme.Caption(values.String(values.StrValidateNote))
		desc.Color = pg.Theme.Color.GrayText2
		return layout.Inset{Bottom: values.MarginPadding20}.Layout(gtx, desc.Layout)
	}
}

func (pg *ValidateAddressPage) actionButtons() layout.Widget {
	return func(gtx C) D {
		dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx C) D {
				return layout.E.Layout(gtx, func(gtx C) D {
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

func (pg *ValidateAddressPage) lineSeparator(gtx C) D {
	m := values.MarginPadding10
	return layout.Inset{Top: m, Bottom: m}.Layout(gtx, pg.Theme.Separator().Layout)
}

func (pg *ValidateAddressPage) showDisplayResult() layout.Widget {
	if pg.stateValidate == none {
		return func(gtx C) D {
			return D{}
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
								ic := decredmaterial.NewIcon(pg.Theme.Icons.NavigationCancel)
								return ic.Layout(gtx, values.MarginPadding25)
							}
							ic := decredmaterial.NewIcon(pg.Theme.Icons.ActionCheckCircle)
							return ic.Layout(gtx, values.MarginPadding25)
						})
					}),
					layout.Rigid(func(gtx C) D {
						if pg.stateValidate == invalid {
							txt := pg.Theme.Body1(values.String(values.StrInvalidAddress))
							txt.Color = pg.Theme.Color.Danger
							txt.TextSize = values.TextSize16
							return txt.Layout(gtx)
						}
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								txt := pg.Theme.Body1(values.String(values.StrValidAddress))
								txt.Color = pg.Theme.Color.Success
								txt.TextSize = values.TextSize16
								return txt.Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx C) D {
										var text string
										if pg.stateValidate == valid {
											text = values.String(values.StrOwned)
										} else {
											text = values.String(values.StrNotOwned)
										}
										txt := pg.Theme.Body1(text)
										txt.TextSize = values.TextSize14
										txt.Color = pg.Theme.Color.GrayText2
										return txt.Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										if pg.stateValidate == valid {
											if components.StringNotEmpty(pg.walletName) {
												return layout.Inset{Left: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
													return decredmaterial.Card{
														Color: pg.Theme.Color.Gray4,
													}.Layout(gtx, func(gtx C) D {
														return layout.UniformInset(values.MarginPadding2).Layout(gtx, func(gtx C) D {
															walletText := pg.Theme.Caption(pg.walletName)
															walletText.Color = pg.Theme.Color.GrayText2
															return walletText.Layout(gtx)
														})
													})
												})
											}
										}
										return D{}
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

func (pg *ValidateAddressPage) pageSections(gtx C, body layout.Widget) D {
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

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *ValidateAddressPage) HandleUserInteractions() {
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
		pg.addressEditor.SetError(values.String(values.StrEnterValidAddress))
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
		pg.validateBtn.Background = pg.Theme.Color.Gray2
		pg.clearBtn.Color = pg.Theme.Color.GrayText4
		pg.isEnabled = false
	} else {
		pg.validateBtn.Background = pg.Theme.Color.Primary
		pg.clearBtn.Color = pg.Theme.Color.Primary
		pg.isEnabled = true
	}
}

func (pg *ValidateAddressPage) clearInputs() {
	pg.validateBtn.Background = pg.Theme.Color.Gray2
	pg.addressEditor.Editor.SetText("")
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *ValidateAddressPage) OnNavigatedFrom() {}

package security

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
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
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	addressEditor         decredmaterial.Editor
	clearBtn, validateBtn decredmaterial.Button
	stateValidate         int
	_                     string
	backButton            decredmaterial.IconButton
}

func NewValidateAddressPage(l *load.Load) *ValidateAddressPage {
	pg := &ValidateAddressPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(ValidateAddressPageID),
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

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *ValidateAddressPage) OnNavigatedTo() {
	pg.addressEditor.Editor.Focus()

	pg.validateBtn.SetEnabled(components.StringNotEmpty(pg.addressEditor.Editor.Text()))
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
				pg.ParentNavigator().CloseCurrentPage()
			},
			Body: func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Rigid(pg.addressSection()),
					)
				})
			},
		}
		return sp.Layout(pg.ParentWindow(), gtx)
	}
	if pg.Load.GetCurrentAppWidth() <= gtx.Dp(values.StartMobileView) {
		return pg.layoutMobile(gtx, body)
	}
	return pg.layoutDesktop(gtx, body)
}

func (pg *ValidateAddressPage) layoutDesktop(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return components.UniformPadding(gtx, body)
}

func (pg *ValidateAddressPage) layoutMobile(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return components.UniformMobile(gtx, false, false, body)
}

func (pg *ValidateAddressPage) addressSection() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(pg.description()),
				layout.Rigid(pg.addressEditor.Layout),
				layout.Rigid(pg.actionButtons()),
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
	pg.validateBtn.SetEnabled(components.StringNotEmpty(pg.addressEditor.Editor.Text()))

	isSubmit, isChanged := decredmaterial.HandleEditorEvents(pg.addressEditor.Editor)
	if isChanged {
		pg.stateValidate = none
	}

	if pg.validateBtn.Clicked() || isSubmit {
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

	var verifyMessageStatus *decredmaterial.Icon
	var verifyMessageText string

	if !pg.WL.MultiWallet.IsAddressValid(address) {
		verifyMessageText = values.String(values.StrInvalidAddress)
		verifyMessageStatus = decredmaterial.NewIcon(pg.Theme.Icons.NavigationCancel)
		verifyMessageStatus.Color = pg.Theme.Color.Danger
	} else {
		verifyMessageStatus = decredmaterial.NewIcon(pg.Theme.Icons.ActionCheck)
		verifyMessageStatus.Color = pg.Theme.Color.Success

		if !pg.WL.SelectedWallet.Wallet.HaveAddress(address) {
			verifyMessageText = values.String(values.StrNotOwned)
		} else {
			verifyMessageText = values.String(values.StrOwned)
		}
	}

	info := modal.NewInfoModal(pg.Load).
		Icon(verifyMessageStatus).
		Title(verifyMessageText).
		SetContentAlignment(layout.Center, layout.Center).
		PositiveButtonStyle(pg.Theme.Color.Primary, pg.Theme.Color.Surface).
		PositiveButton(values.String(values.StrGotIt), func(isChecked bool) bool {
			return true
		})
	pg.ParentWindow().ShowModal(info)
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *ValidateAddressPage) OnNavigatedFrom() {}

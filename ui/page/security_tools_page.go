package page

import (
	"image"

	"gioui.org/layout"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const SecurityToolsPageID = "SecurityTools"

type SecurityToolsPage struct {
	*load.Load
	verifyMessage   *decredmaterial.Clickable
	validateAddress *decredmaterial.Clickable

	backButton decredmaterial.IconButton
}

func NewSecurityToolsPage(l *load.Load) *SecurityToolsPage {
	pg := &SecurityToolsPage{
		Load:            l,
		verifyMessage:   l.Theme.NewClickable(true),
		validateAddress: l.Theme.NewClickable(true),
	}

	pg.verifyMessage.Radius = decredmaterial.Radius(14)
	pg.validateAddress.Radius = decredmaterial.Radius(14)

	pg.backButton, _ = components.SubpageHeaderButtons(l)

	return pg
}

// ID is a unique string that identifies the page and may be used
// to differentiate this page from other pages.
// Part of the load.Page interface.
func (pg *SecurityToolsPage) ID() string {
	return SecurityToolsPageID
}

// WillAppear is called when the page is about to displayed and may
// be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *SecurityToolsPage) WillAppear() {

}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
// main settings layout
func (pg *SecurityToolsPage) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      "Security Tools",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Flexed(.5, pg.message()),
						layout.Rigid(func(gtx C) D {
							size := image.Point{X: 15, Y: gtx.Constraints.Min.Y}
							return layout.Dimensions{Size: size}
						}),
						layout.Flexed(.5, pg.address()),
					)
				})
			},
		}
		return sp.Layout(gtx)
	}
	return components.UniformPadding(gtx, body)
}

func (pg *SecurityToolsPage) message() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, pg.Icons.VerifyMessageIcon, pg.verifyMessage, pg.Theme.Body1("Verify Message").Layout)
	}
}

func (pg *SecurityToolsPage) address() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, pg.Icons.LocationPinIcon, pg.validateAddress, pg.Theme.Body1("Validate Address").Layout)
	}
}

func (pg *SecurityToolsPage) pageSections(gtx layout.Context, icon *decredmaterial.Image, action *decredmaterial.Clickable, body layout.Widget) layout.Dimensions {
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return pg.Theme.Card().Layout(gtx, func(gtx C) D {
			return action.Layout(gtx, func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return icon.Layout24dp(gtx)
						}),
						layout.Rigid(body),
						layout.Rigid(func(gtx C) D {
							size := image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Min.Y}
							return layout.Dimensions{Size: size}
						}),
					)
				})
			})
		})
	})
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *SecurityToolsPage) HandleUserInteractions() {
	if pg.verifyMessage.Clicked() {
		pg.ChangeFragment(NewVerifyMessagePage(pg.Load))
	}

	if pg.validateAddress.Clicked() {
		pg.ChangeFragment(NewValidateAddressPage(pg.Load))
	}
}

// WillDisappear is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// WillAppear() will be called again. This method should not destroy UI
// components unless they'll be recreated in the WillAppear() method.
// Part of the load.Page interface.
func (pg *SecurityToolsPage) WillDisappear() {}

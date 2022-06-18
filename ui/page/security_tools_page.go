package page

import (
	"image"

	"gioui.org/layout"

	"github.com/planetdecred/godcr/app"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/modal"
	"github.com/planetdecred/godcr/ui/page/components"
	"github.com/planetdecred/godcr/ui/values"
)

const SecurityToolsPageID = "SecurityTools"

type SecurityToolsPage struct {
	*load.Load
	// GenericPageModal defines methods such as ID() and OnAttachedToNavigator()
	// that helps this Page satisfy the app.Page interface. It also defines
	// helper methods for accessing the PageNavigator that displayed this page
	// and the root WindowNavigator.
	*app.GenericPageModal

	verifyMessage   *decredmaterial.Clickable
	validateAddress *decredmaterial.Clickable
	signMsg         *decredmaterial.Clickable
	shadowBox       *decredmaterial.Shadow
	infoButton      decredmaterial.IconButton

	backButton decredmaterial.IconButton
}

func NewSecurityToolsPage(l *load.Load) *SecurityToolsPage {
	pg := &SecurityToolsPage{
		Load:             l,
		GenericPageModal: app.NewGenericPageModal(SecurityToolsPageID),
		verifyMessage:    l.Theme.NewClickable(true),
		validateAddress:  l.Theme.NewClickable(true),
		signMsg:          l.Theme.NewClickable(true),
	}

	pg.shadowBox = l.Theme.Shadow()
	pg.shadowBox.SetShadowRadius(14)

	pg.verifyMessage.Radius = decredmaterial.Radius(14)
	pg.validateAddress.Radius = decredmaterial.Radius(14)
	pg.signMsg.Radius = decredmaterial.Radius(14)

	pg.backButton, pg.infoButton = components.SubpageHeaderButtons(l)

	return pg
}

// OnNavigatedTo is called when the page is about to be displayed and
// may be used to initialize page features that are only relevant when
// the page is displayed.
// Part of the load.Page interface.
func (pg *SecurityToolsPage) OnNavigatedTo() {

}

// Layout draws the page UI components into the provided layout context
// to be eventually drawn on screen.
// Part of the load.Page interface.
// main settings layout
func (pg *SecurityToolsPage) Layout(gtx C) D {
	body := func(gtx C) D {
		sp := components.SubPage{
			Load:       pg.Load,
			Title:      values.String(values.StrSecurityTools),
			InfoButton: pg.infoButton,
			BackButton: pg.backButton,
			Back: func() {
				pg.ParentNavigator().CloseCurrentPage()
			},
			Body: func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(pg.message()),
					layout.Rigid(pg.address()),
					layout.Rigid(pg.signMessage()),
				)
			},
			InfoTemplate: modal.SecurityToolsInfoTemplate,
		}
		return sp.Layout(pg.ParentWindow(), gtx)
	}
	return components.UniformPadding(gtx, body)
}

func (pg *SecurityToolsPage) message() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, pg.Theme.Icons.VerifyMessageIcon, pg.verifyMessage, values.String(values.StrVerifyMessage))
	}
}

func (pg *SecurityToolsPage) address() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, pg.Theme.Icons.LocationPinIcon, pg.validateAddress, values.String(values.StrValidateMsg))
	}
}

func (pg *SecurityToolsPage) signMessage() layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, pg.Theme.Icons.SignMessageIcon, pg.signMsg, values.String(values.StrSignMessage))
	}
}

func (pg *SecurityToolsPage) pageSections(gtx C, icon *decredmaterial.Image, action *decredmaterial.Clickable, title string) D {
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return decredmaterial.LinearLayout{
			Orientation: layout.Vertical,
			Width:       decredmaterial.MatchParent,
			Height:      decredmaterial.WrapContent,
			Background:  pg.Theme.Color.Surface,
			Clickable:   action,
			Direction:   layout.Center,
			Alignment:   layout.Middle,
			Shadow:      pg.shadowBox,
			Border:      decredmaterial.Border{Radius: decredmaterial.Radius(14)},
			Padding:     layout.UniformInset(values.MarginPadding15),
			Margin:      layout.Inset{Bottom: values.MarginPadding4, Top: values.MarginPadding4}}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return icon.Layout24dp(gtx)
			}),
			layout.Rigid(pg.Theme.Body1(title).Layout),
			layout.Rigid(func(gtx C) D {
				size := image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Min.Y}
				return D{Size: size}
			}),
		)
	})
}

// HandleUserInteractions is called just before Layout() to determine
// if any user interaction recently occurred on the page and may be
// used to update the page's UI components shortly before they are
// displayed.
// Part of the load.Page interface.
func (pg *SecurityToolsPage) HandleUserInteractions() {
	if pg.verifyMessage.Clicked() {
		pg.ParentNavigator().Display(NewVerifyMessagePage(pg.Load))
	}

	if pg.validateAddress.Clicked() {
		pg.ParentNavigator().Display(NewValidateAddressPage(pg.Load))
	}
}

// OnNavigatedFrom is called when the page is about to be removed from
// the displayed window. This method should ideally be used to disable
// features that are irrelevant when the page is NOT displayed.
// NOTE: The page may be re-displayed on the app's window, in which case
// OnNavigatedTo() will be called again. This method should not destroy UI
// components unless they'll be recreated in the OnNavigatedTo() method.
// Part of the load.Page interface.
func (pg *SecurityToolsPage) OnNavigatedFrom() {}

package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/values"

	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const PageSecurityTools = "Security Tools"

type securityToolsPage struct {
	theme           *decredmaterial.Theme
	verifyMessage   *widget.Clickable
	validateAddress *widget.Clickable
}

func (win *Window) SecurityToolsPage(common pageCommon) layout.Widget {
	pg := &securityToolsPage{
		theme:           common.theme,
		verifyMessage:   new(widget.Clickable),
		validateAddress: new(widget.Clickable),
	}

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

// main settings layout
func (pg *securityToolsPage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	body := func(gtx C) D {
		page := SubPage{
			title: "Security Tools",
			back: func() {
				*common.page = PageMore
			},
			body: func(gtx C) D {
				return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Flexed(.5, pg.message(common)),
						layout.Rigid(func(gtx C) D {
							size := image.Point{X: 15, Y: gtx.Constraints.Min.Y}
							return layout.Dimensions{Size: size}
						}),
						layout.Flexed(.5, pg.address(common)),
					)
				})
			},
			infoTemplate: SecurityToolsInfoTemplate,
		}
		return common.SubPageLayout(gtx, page)
	}
	return common.Layout(gtx, body)
}

func (pg *securityToolsPage) message(common pageCommon) layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, common.icons.verifyMessageIcon, pg.verifyMessage, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return common.theme.Body1("Verify Message").Layout(gtx)
				}),
			)
		})
	}
}

func (pg *securityToolsPage) address(common pageCommon) layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, common.icons.locationPinIcon, pg.validateAddress, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return common.theme.Body1("Validate Address").Layout(gtx)
				}),
			)
		})
	}
}

func (pg *securityToolsPage) pageSections(gtx layout.Context, icon *decredmaterial.Image, action *widget.Clickable, body layout.Widget) layout.Dimensions {
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return pg.theme.Card().Layout(gtx, func(gtx C) D {
			return decredmaterial.Clickable(gtx, action, func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							return icon.Layout(gtx)
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

func (pg *securityToolsPage) handle(common pageCommon) {

}

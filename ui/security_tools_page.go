package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/values"

	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const PageSecurityTools = "Security_Tools"

type securityToolsPage struct {
	theme *decredmaterial.Theme

	dummyText decredmaterial.Label
}

func (win *Window) SecurityToolsPage(common pageCommon) layout.Widget {
	pg := &securityToolsPage{
		theme:     common.theme,
		dummyText: common.theme.H5("Coming soon"),
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
				gtx.Constraints.Min = gtx.Constraints.Max
				return layout.Inset{Top: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return layout.Flex{Spacing: layout.SpaceBetween}.Layout(gtx,
						layout.Flexed(.5, pg.message(common)),
						layout.Flexed(.5, pg.validateAddress(common)),
					)
				})
			},
			infoTemplate: SecurityToolsInfoTemplate,
		}
		return common.SubPageLayoutWithoutInfo(gtx, page)
	}
	return common.Layout(gtx, body)
}

func (pg *securityToolsPage) message(common pageCommon) layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, common.icons.verifyMessageIcon, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return common.theme.Body1("Verify Message").Layout(gtx)
				}),
			)
		})
	}
}

func (pg *securityToolsPage) validateAddress(common pageCommon) layout.Widget {
	return func(gtx C) D {
		return pg.pageSections(gtx, common.icons.locationPinIcon, func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return common.theme.Body1("Validate Address").Layout(gtx)
				}),
			)
		})
	}
}

func (pg *securityToolsPage) pageSections(gtx layout.Context, icon *widget.Image, body layout.Widget) layout.Dimensions {
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return pg.theme.Card().Layout(gtx, func(gtx C) D {
			return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						icon := icon
						icon.Scale = 1
						return icon.Layout(gtx)
					}),
					layout.Rigid(body),
				)
			})
		})
	})
}

func (pg *securityToolsPage) handle(common pageCommon) {

}

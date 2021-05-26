package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/widget"
	"github.com/planetdecred/godcr/ui/values"

	"github.com/planetdecred/godcr/ui/decredmaterial"
)

const PageSecurityTools = "SecurityTools"

type securityToolsPage struct {
	theme           *decredmaterial.Theme
	verifyMessage   *widget.Clickable
	validateAddress *widget.Clickable
	common          pageCommon
}

func (win *Window) SecurityToolsPage(common pageCommon) Page {
	pg := &securityToolsPage{
		theme:           common.theme,
		verifyMessage:   new(widget.Clickable),
		validateAddress: new(widget.Clickable),
		common:          common,
	}

	return pg
}

// main settings layout
func (pg *securityToolsPage) Layout(gtx layout.Context) layout.Dimensions {
	common := pg.common
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
	return common.Layout(gtx, func(gtx C) D {
		return common.UniformPadding(gtx, body)
	})
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

func (pg *securityToolsPage) pageSections(gtx layout.Context, icon *widget.Image, action *widget.Clickable, body layout.Widget) layout.Dimensions {
	return layout.Inset{Bottom: values.MarginPadding10}.Layout(gtx, func(gtx C) D {
		return pg.theme.Card().Layout(gtx, func(gtx C) D {
			return decredmaterial.Clickable(gtx, action, func(gtx C) D {
				return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceAround}.Layout(gtx,
						layout.Rigid(func(gtx C) D {
							icon := icon
							icon.Scale = 1
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

func (pg *securityToolsPage) handle() {
	common := pg.common
	if pg.verifyMessage.Clicked() {
		*common.returnPage = PageSecurityTools
		common.setReturnPage(PageSecurityTools)
		common.changePage(PageVerifyMessage)
	}

	if pg.validateAddress.Clicked() {
		*common.returnPage = PageSecurityTools
		common.setReturnPage(PageSecurityTools)
		common.changePage(ValidateAddress)
	}
}

func (pg *securityToolsPage) onClose() {}

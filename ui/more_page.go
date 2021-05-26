package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PageMore = "More"

type morePageHandler struct {
	clickable *widget.Clickable
	image     *widget.Image
	page      string
}

type morePage struct {
	common            pageCommon
	container         layout.Flex
	morePageListItems []morePageHandler
	page              *string
}

func (win *Window) MorePage(common pageCommon) Page {
	morePageListItems := []morePageHandler{
		{
			clickable: new(widget.Clickable),
			image:     common.icons.settingsIcon,
			page:      PageSettings,
		},
		{
			clickable: new(widget.Clickable),
			image:     common.icons.securityIcon,
			page:      PageSecurityTools,
		},
		{
			clickable: new(widget.Clickable),
			image:     common.icons.helpIcon,
			page:      PageHelp,
		},
		{
			clickable: new(widget.Clickable),
			image:     common.icons.aboutIcon,
			page:      PageAbout,
		},
		{
			clickable: new(widget.Clickable),
			image:     common.icons.debugIcon,
			page:      PageDebug,
		},
	}

	for i := range morePageListItems {
		morePageListItems[i].image.Scale = 1
	}

	pg := &morePage{
		container:         layout.Flex{Axis: layout.Vertical},
		morePageListItems: morePageListItems,
		page:              &win.current,
		common:            common,
	}

	return pg
}

func (pg *morePage) handleClickEvents(common pageCommon) {
	for i := range pg.morePageListItems {
		for pg.morePageListItems[i].clickable.Clicked() {
			common.changePage(pg.morePageListItems[i].page)
		}
	}
}

func (pg *morePage) Layout(gtx layout.Context) layout.Dimensions {
	common := pg.common
	pg.handleClickEvents(common)

	container := func(gtx C) D {
		pg.layoutMoreItems(gtx, common)
		return layout.Dimensions{Size: gtx.Constraints.Max}
	}
	return common.Layout(gtx, func(gtx C) D {
		return common.UniformPadding(gtx, container)
	})
}

func (pg *morePage) layoutMoreItems(gtx layout.Context, common pageCommon) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			list := layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(pg.morePageListItems), func(gtx C, i int) D {
				return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
					return decredmaterial.Clickable(gtx, pg.morePageListItems[i].clickable, func(gtx C) D {
						background := common.theme.Color.Surface
						card := common.theme.Card()
						card.Color = background
						return card.Layout(gtx, func(gtx C) D {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							return layout.Stack{}.Layout(gtx,
								layout.Stacked(func(gtx C) D {
									return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
										gtx.Constraints.Min.X = gtx.Constraints.Max.X
										return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												return layout.Center.Layout(gtx, func(gtx C) D {
													return pg.morePageListItems[i].image.Layout(gtx)
												})
											}),
											layout.Rigid(func(gtx C) D {
												return layout.Inset{
													Left: values.MarginPadding15,
													Top:  values.MarginPadding2,
												}.Layout(gtx, func(gtx C) D {
													return layout.Center.Layout(gtx, func(gtx C) D {
														page := pg.morePageListItems[i].page
														if page == PageSecurityTools {
															page = "Security Tools"
														}
														return common.theme.Body1(page).Layout(gtx)
													})
												})
											}),
										)
									})
								}),
							)
						})
					})
				})
			})
		}),
	)
}

func (pg *morePage) handle()  {}
func (pg *morePage) onClose() {}

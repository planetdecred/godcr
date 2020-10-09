package ui

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/godcr/ui/decredmaterial"
	"github.com/raedahgroup/godcr/ui/values"
)

const PageMore = "more"

type morePageIcons struct {

	settingsIcon, securityIcon, politeiaIcon, helpIcon, aboutIcon, debugIcon, logo image.Image
}

type morePageHandler struct {
	clickable *widget.Clickable
	image     *widget.Image
	page      string
}

type morePage struct {
	container                                   layout.Flex
	morePageListItems                           []morePageHandler
	page                                        *string
	icons                                       morePageIcons
}

func (win *Window) MorePage(decredIcons map[string]image.Image, common pageCommon) layout.Widget {

	ic := morePageIcons{
		settingsIcon:               decredIcons["overview"],
		securityIcon:               decredIcons["wallet_inactive"],
		politeiaIcon:               decredIcons["receive"],
		helpIcon:                   decredIcons["transaction_inactive"],
		aboutIcon:                  decredIcons["send"],
		debugIcon:                  decredIcons["transaction"],
	}

    morePageListItems := []morePageHandler{
		{
			clickable: new(widget.Clickable),
			image:     &widget.Image{Src: paint.NewImageOp(ic.settingsIcon)},
			page:      PageOverview,
		},
		{
			clickable: new(widget.Clickable),
			image:     &widget.Image{Src: paint.NewImageOp(ic.securityIcon)},
			page:      PageTransactions,
		},
		{
			clickable: new(widget.Clickable),
			image:     &widget.Image{Src: paint.NewImageOp(ic.politeiaIcon)},
			page:      PageWallet,
		},
		{
			clickable: new(widget.Clickable),
			image:     &widget.Image{Src: paint.NewImageOp(ic.helpIcon)},
			page:      PageMore,
		},
		{
			clickable: new(widget.Clickable),
			image:     &widget.Image{Src: paint.NewImageOp(ic.aboutIcon)},
			page:      PageWallet,
		},
		{
			clickable: new(widget.Clickable),
			image:     &widget.Image{Src: paint.NewImageOp(ic.debugIcon)},
			page:      PageMore,
		},
	}

	pg := morePage{
		container:              layout.Flex{Axis: layout.Vertical},
		morePageListItems:      morePageListItems,
		page:                   &win.current,
	}

	return func(gtx C) D {
		pg.Handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *morePage) handleClickEvents() {
	for i := range pg.morePageListItems {
		for pg.morePageListItems[i].clickable.Clicked() {
			*pg.page = pg.morePageListItems[i].page
		}
	}
}

func (pg *morePage) Layout(gtx layout.Context, common pageCommon) layout.Dimensions {
	pg.handleClickEvents()

	container := func(gtx C) D {
		return decredmaterial.Card{Rounded: true}.Layout(gtx, func(gtx C) D {
			pg.layoutMoreItems(gtx, common)
			return layout.Dimensions{Size: gtx.Constraints.Max}
		})
	}
	return common.Layout(gtx, container)
}

func (pg *morePage) layoutMoreItems(gtx layout.Context, common pageCommon) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			list := layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(pg.morePageListItems), func(gtx C, i int) D {
				return layout.Inset{Bottom: values.MarginPadding5}.Layout(gtx, func(gtx C) D {
				return decredmaterial.Clickable(gtx, pg.morePageListItems[i].clickable, func(gtx C) D {
					background := common.theme.Color.Surface

					return decredmaterial.Card{Color: background, Rounded: true}.Layout(gtx, func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Stack{}.Layout(gtx,
							layout.Stacked(func(gtx C) D {
								return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx C) D {
									axis := layout.Horizontal
									leftInset := float32(15)

									gtx.Constraints.Min.X = gtx.Constraints.Max.X
									return layout.Flex{Axis: axis}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											pg.morePageListItems[i].image.Scale = 0.05

											return layout.Center.Layout(gtx, func(gtx C) D {
												return pg.morePageListItems[i].image.Layout(gtx)
											})
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{
												Left: unit.Dp(leftInset),
											}.Layout(gtx, func(gtx C) D {
												return layout.Center.Layout(gtx, func(gtx C) D {
													return common.theme.Body1(pg.morePageListItems[i].page).Layout(gtx)
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

func (pg *morePage) Handle(common pageCommon) {
	
}

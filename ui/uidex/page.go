package uidex

import (
	"image"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/utils"
	"github.com/planetdecred/godcr/ui/values"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type pageIcons struct {
	contentAdd, navigationCheck, navigationMore, actionCheckCircle, actionInfo, navigationArrowBack,
	navigationArrowForward, actionCheck, chevronRight, navigationCancel, navMoreIcon,
	imageBrightness1, contentClear, dropDownIcon, cached *widget.Icon

	logo, btc, dcr, ltc, bch *widget.Image
}

type navHandler struct {
	clickable     *widget.Clickable
	image         *widget.Image
	imageInactive *widget.Image
	page          string
}
type pageCommon struct {
	printer        *message.Printer
	theme          *decredmaterial.Theme
	icons          pageIcons
	page           *string
	returnPage     *string
	navTab         *decredmaterial.Tabs
	keyEvents      chan *key.Event
	modal          *decredmaterial.Modal
	appBarNavItems []navHandler

	changePage    func(string)
	setReturnPage func(string)
	refreshWindow func()
	switchView    *int
}

type (
	C = layout.Context
	D = layout.Dimensions
)

func (d *Dex) addPages(decredIcons map[string]image.Image) {
	ic := pageIcons{
		contentAdd:             utils.MustIcon(widget.NewIcon(icons.ContentAdd)),
		navigationCheck:        utils.MustIcon(widget.NewIcon(icons.NavigationCheck)),
		navigationMore:         utils.MustIcon(widget.NewIcon(icons.NavigationMoreHoriz)),
		actionCheckCircle:      utils.MustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		navigationArrowBack:    utils.MustIcon(widget.NewIcon(icons.NavigationArrowBack)),
		navigationArrowForward: utils.MustIcon(widget.NewIcon(icons.NavigationArrowForward)),
		actionInfo:             utils.MustIcon(widget.NewIcon(icons.ActionInfo)),
		actionCheck:            utils.MustIcon(widget.NewIcon(icons.ActionCheckCircle)),
		navigationCancel:       utils.MustIcon(widget.NewIcon(icons.NavigationCancel)),
		imageBrightness1:       utils.MustIcon(widget.NewIcon(icons.ImageBrightness1)),
		chevronRight:           utils.MustIcon(widget.NewIcon(icons.NavigationChevronRight)),
		contentClear:           utils.MustIcon(widget.NewIcon(icons.ContentClear)),
		navMoreIcon:            utils.MustIcon(widget.NewIcon(icons.NavigationMoreHoriz)),
		dropDownIcon:           utils.MustIcon(widget.NewIcon(icons.NavigationArrowDropDown)),
		cached:                 utils.MustIcon(widget.NewIcon(icons.ActionCached)),

		logo: &widget.Image{Src: paint.NewImageOp(decredIcons["favicon"])},
		btc:  &widget.Image{Src: paint.NewImageOp(decredIcons["btc"])},
		dcr:  &widget.Image{Src: paint.NewImageOp(decredIcons["dcr"])},
		bch:  &widget.Image{Src: paint.NewImageOp(decredIcons["bch"])},
		ltc:  &widget.Image{Src: paint.NewImageOp(decredIcons["ltc"])},
	}

	common := pageCommon{
		printer:       message.NewPrinter(language.English),
		theme:         d.theme,
		icons:         ic,
		returnPage:    &d.previous,
		page:          &d.current,
		modal:         d.theme.Modal(),
		changePage:    d.changePage,
		setReturnPage: d.setReturnPage,
		refreshWindow: d.refresh,

		switchView: d.switchView,
	}

	d.pages = make(map[string]layout.Widget)
	d.pages[PageMarkets] = d.MarketsPage(common)
}

func (page pageCommon) refreshPage() {
	page.refreshWindow()
}

func (page pageCommon) notify(text string, success bool) {

}

func (page pageCommon) Layout(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(page.layoutTopBar),
				layout.Rigid(body),
			)
		}),
		layout.Stacked(func(gtx C) D {
			return layout.Dimensions{}
		}),
	)
}

// Container is simply a wrapper for the Inset type. Its purpose is to differentiate the use of an inset as a padding or
// margin, making it easier to visualize the structure of a layout when reading UI code.
type Container struct {
	padding layout.Inset
}

func (c Container) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return c.padding.Layout(gtx, w)
}

func (page pageCommon) UniformPadding(gtx layout.Context, body layout.Widget) layout.Dimensions {
	return layout.UniformInset(values.MarginPadding24).Layout(gtx, body)
}

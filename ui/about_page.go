package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PageAbout = "About"

type aboutPage struct {
	common    *pageCommon
	theme     *decredmaterial.Theme
	card      decredmaterial.Card
	container *layout.List

	version        decredmaterial.Label
	versionValue   decredmaterial.Label
	buildDate      decredmaterial.Label
	buildDateValue decredmaterial.Label
	network        decredmaterial.Label
	networkValue   decredmaterial.Label
	license        decredmaterial.Label
	licenseRow     *widget.Clickable

	chevronRightIcon *widget.Icon

	backButton decredmaterial.IconButton
}

func AboutPage(common *pageCommon) Page {
	pg := &aboutPage{
		common:           common,
		theme:            common.theme,
		card:             common.theme.Card(),
		container:        &layout.List{Axis: layout.Vertical},
		version:          common.theme.Body1("Version"),
		versionValue:     common.theme.Body1("v1.5.2"),
		buildDate:        common.theme.Body1("Build date"),
		buildDateValue:   common.theme.Body1("2020-09-10"),
		network:          common.theme.Body1("Network"),
		networkValue:     common.theme.Body1(common.wallet.Net),
		license:          common.theme.Body1("License"),
		licenseRow:       new(widget.Clickable),
		chevronRightIcon: common.icons.chevronRight,
	}

	pg.backButton, _ = common.SubPageHeaderButtons()
	pg.versionValue.Color = pg.theme.Color.Gray
	pg.buildDateValue.Color = pg.theme.Color.Gray
	pg.networkValue.Color = pg.theme.Color.Gray
	pg.chevronRightIcon.Color = pg.theme.Color.Gray

	return pg
}

func (pg *aboutPage) OnResume() {

}

func (pg *aboutPage) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		page := SubPage{
			title:      "About",
			backButton: pg.backButton,
			back: func() {
				pg.common.changePage(PageMore)
			},
			body: func(gtx C) D {
				return pg.card.Layout(gtx, func(gtx C) D {
					return pg.layoutRows(gtx)
				})
			},
		}
		return pg.common.SubPageLayout(gtx, page)
	}

	return pg.common.UniformPadding(gtx, body)
}

func (pg *aboutPage) layoutRows(gtx layout.Context) layout.Dimensions {
	w := []func(gtx C) D{
		func(gtx C) D {
			return endToEndRow(gtx, pg.version.Layout, pg.versionValue.Layout)
		},
		func(gtx C) D {
			return endToEndRow(gtx, pg.buildDate.Layout, pg.buildDateValue.Layout)
		},
		func(gtx C) D {
			return endToEndRow(gtx, pg.network.Layout, pg.networkValue.Layout)
		},
		func(gtx C) D {
			return decredmaterial.Clickable(gtx, pg.licenseRow, func(gtx C) D {
				return layout.Flex{}.Layout(gtx,
					layout.Rigid(pg.license.Layout),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return pg.chevronRightIcon.Layout(gtx, values.MarginPadding20)
						})
					}),
				)
			})
		},
	}

	return pg.container.Layout(gtx, len(w), func(gtx C, i int) D {
		return layout.Inset{}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return Container{
						layout.Inset{
							Top:    values.MarginPadding20,
							Bottom: values.MarginPadding20,
							Left:   values.MarginPadding16,
							Right:  values.MarginPadding16,
						},
					}.Layout(gtx, w[i])
				}),
				layout.Rigid(func(gtx C) D {
					if i == len(w)-1 {
						return layout.Dimensions{}
					}
					return layout.Inset{
						Left: values.MarginPadding16,
					}.Layout(gtx, pg.theme.Separator().Layout)
				}),
			)
		})
	})
}

func (pg *aboutPage) handle() {
	if pg.licenseRow.Clicked() {
		pg.common.changeFragment(LicensePage(pg.common), PageLicense)
	}
}
func (pg *aboutPage) onClose() {}

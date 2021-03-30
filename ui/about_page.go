package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PageAbout = "About"

type aboutPageRow struct {
	left  *decredmaterial.Label
	right *decredmaterial.Label
	icon  *widget.Icon
}

type aboutPage struct {
	theme     *decredmaterial.Theme
	card      decredmaterial.Card
	container *layout.List
	line      *decredmaterial.Line

	version        decredmaterial.Label
	versionValue   decredmaterial.Label
	buildDate      decredmaterial.Label
	buildDateValue decredmaterial.Label
	network        decredmaterial.Label
	networkValue   decredmaterial.Label
	license        decredmaterial.Label

	chevronRightIcon *widget.Icon
}

func (win *Window) AboutPage(common pageCommon) layout.Widget {
	pg := &aboutPage{
		theme:            common.theme,
		card:             common.theme.Card(),
		line:             common.theme.Line(),
		container:        &layout.List{Axis: layout.Vertical},
		version:          common.theme.Body1("Version"),
		versionValue:     common.theme.Body1("v1.5.2"),
		buildDate:        common.theme.Body1("Build date"),
		buildDateValue:   common.theme.Body1("2020-09-10"),
		network:          common.theme.Body1("Network"),
		networkValue:     common.theme.Body1(win.wallet.Net),
		license:          common.theme.Body1("License"),
		chevronRightIcon: common.icons.chevronRight,
	}
	pg.line.Height = 1
	pg.line.Color = common.theme.Color.Background

	pg.version.Color = pg.theme.Color.Text
	pg.buildDate.Color = pg.theme.Color.Text
	pg.network.Color = pg.theme.Color.Text
	pg.license.Color = pg.theme.Color.Text
	pg.versionValue.Color = pg.theme.Color.Gray
	pg.buildDateValue.Color = pg.theme.Color.Gray
	pg.networkValue.Color = pg.theme.Color.Gray

	pg.chevronRightIcon.Color = pg.theme.Color.Gray

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

func (pg *aboutPage) Layout(gtx C, common pageCommon) D {
	body := func(gtx C) D {
		page := SubPage{
			title: "About",
			back: func() {
				common.changePage(PageMore)
			},
			body: func(gtx layout.Context) layout.Dimensions {
				return pg.card.Layout(gtx, func(gtx C) D {
					return pg.layoutRows(gtx)
				})
			},
		}
		return common.SubPageLayout(gtx, page)
	}

	return common.Layout(gtx, func(gtx C) D {
		return common.UniformPadding(gtx, body)
	})
}

func (pg *aboutPage) layoutRows(gtx C) D {
	w := []func(gtx C) D{
		func(gtx C) D {
			row := aboutPageRow{
				left:  &pg.version,
				right: &pg.versionValue,
			}
			return pg.layoutRow(gtx, row, true)
		},
		func(gtx C) D {
			row := aboutPageRow{
				left:  &pg.buildDate,
				right: &pg.buildDateValue,
			}
			return pg.layoutRow(gtx, row, true)
		},
		func(gtx C) D {
			row := aboutPageRow{
				left:  &pg.network,
				right: &pg.networkValue,
			}
			return pg.layoutRow(gtx, row, true)
		},
		func(gtx C) D {
			row := aboutPageRow{
				left: &pg.license,
				icon: pg.chevronRightIcon,
			}
			return pg.layoutRow(gtx, row, false)
		},
	}

	return pg.container.Layout(gtx, len(w), func(gtx C, i int) D {
		return layout.UniformInset(values.MarginPadding0).Layout(gtx, w[i])
	})
}

func (pg *aboutPage) layoutRow(gtx C, row aboutPageRow, drawSeparator bool) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(row.left.Layout),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							if row.icon != nil {
								return row.icon.Layout(gtx, values.MarginPadding20)
							}
							return row.right.Layout(gtx)
						})
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			if !drawSeparator {
				return D{}
			}
			pg.line.Width = gtx.Constraints.Max.X
			return layout.Inset{
				Left: values.MarginPadding15,
			}.Layout(gtx, pg.line.Layout)
		}),
	)
}

func (pg *aboutPage) handle(common pageCommon) {

}

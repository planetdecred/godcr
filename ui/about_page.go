package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
)

const PageAbout = "About"

type aboutPageRow struct {
	leftLabel  *decredmaterial.Label
	rightLabel *decredmaterial.Label
	icon       *widget.Icon
}

type aboutPage struct {
	theme     *decredmaterial.Theme
	card      decredmaterial.Card
	container *layout.List
	line      *decredmaterial.Line

	versionLabel        decredmaterial.Label
	versionValueLabel   decredmaterial.Label
	buildDateLabel      decredmaterial.Label
	buildDateValueLabel decredmaterial.Label
	networkLabel        decredmaterial.Label
	networkValueLabel   decredmaterial.Label
	licenseLabel        decredmaterial.Label

	chevronRightIcon *widget.Icon
}

func (win *Window) AboutPage(common pageCommon) layout.Widget {
	pg := &aboutPage{
		theme:               common.theme,
		card:                common.theme.Card(),
		line:                common.theme.Line(),
		container:           &layout.List{Axis: layout.Vertical},
		versionLabel:        common.theme.Body1("Version"),
		versionValueLabel:   common.theme.Body2("v1.5.2"),
		buildDateLabel:      common.theme.Body1("Build date"),
		buildDateValueLabel: common.theme.Body2("2020-09-10"),
		networkLabel:        common.theme.Body1("Network"),
		networkValueLabel:   common.theme.Body2("Testnet3"),
		licenseLabel:        common.theme.Body1("License"),
		chevronRightIcon:    common.icons.chevronRight,
	}
	pg.line.Height = 1
	pg.line.Color = common.theme.Color.Background

	pg.versionValueLabel.Color = pg.theme.Color.Gray
	pg.buildDateValueLabel.Color = pg.theme.Color.Gray
	pg.networkValueLabel.Color = pg.theme.Color.Gray

	pg.chevronRightIcon.Color = pg.theme.Color.Gray

	return func(gtx C) D {
		pg.handle(common)
		return pg.Layout(gtx, common)
	}
}

// main settings layout
func (pg *aboutPage) Layout(gtx C, common pageCommon) D {
	body := func(gtx C) D {
		page := SubPage{
			title: "About",
			back: func() {
				common.changePage(PageMore)
			},
			body: func(gtx layout.Context) layout.Dimensions {
				return pg.card.Layout(gtx, func(gtx C) D {
					return layout.UniformInset(values.MarginPadding15).Layout(gtx, func(gtx C) D {
						return pg.layoutRows(gtx)
					})
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
				leftLabel:  &pg.versionLabel,
				rightLabel: &pg.versionValueLabel,
			}
			return pg.layoutRow(gtx, row, true)
		},
		func(gtx C) D {
			row := aboutPageRow{
				leftLabel:  &pg.buildDateLabel,
				rightLabel: &pg.buildDateValueLabel,
			}
			return pg.layoutRow(gtx, row, true)
		},
		func(gtx C) D {
			row := aboutPageRow{
				leftLabel:  &pg.networkLabel,
				rightLabel: &pg.networkValueLabel,
			}
			return pg.layoutRow(gtx, row, true)
		},
		func(gtx C) D {
			row := aboutPageRow{
				leftLabel: &pg.licenseLabel,
				icon:      pg.chevronRightIcon,
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
			return layout.Inset{
				Top:    values.MarginPadding5,
				Bottom: values.MarginPadding5,
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(row.leftLabel.Layout),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							if row.icon != nil {
								return row.icon.Layout(gtx, values.MarginPadding30)
							}
							return row.rightLabel.Layout(gtx)
						})
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			if !drawSeparator {
				return D{}
			}
			return layout.Inset{
				Top:    values.MarginPadding5,
				Bottom: values.MarginPadding5,
			}.Layout(gtx, func(gtx C) D {
				pg.line.Width = gtx.Constraints.Max.X
				return pg.line.Layout(gtx)
			})
		}),
	)
}

func (pg *aboutPage) layoutRowD(gtx C, leftLabel, rightLabel decredmaterial.Label) D {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Top:    values.MarginPadding5,
				Bottom: values.MarginPadding5,
			}.Layout(gtx, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(leftLabel.Layout),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, rightLabel.Layout)
					}),
				)
			})
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Top:    values.MarginPadding5,
				Bottom: values.MarginPadding5,
			}.Layout(gtx, func(gtx C) D {
				pg.line.Width = gtx.Constraints.Max.X
				return pg.line.Layout(gtx)
			})
		}),
	)
}

func (pg *aboutPage) handle(common pageCommon) {

}

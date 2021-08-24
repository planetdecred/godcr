package page

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/load"
	"github.com/planetdecred/godcr/ui/page/components"
)

const AboutPageID = "About"

type unclickableRow struct {
	leftWidget  decredmaterial.Label
	rightWidget decredmaterial.Label
}

type clickableLicense struct {
	clickable   *widget.Clickable
	leftWidget  decredmaterial.Label
	rightWidget *widget.Icon
}

type AboutPage struct {
	*load.Load
	container    *layout.List
	card         decredmaterial.Card
	versionRow   unclickableRow
	buildDateRow unclickableRow
	networkRow   unclickableRow
	licenseRow   clickableLicense
	backButton   decredmaterial.IconButton
}

/*type AboutPage struct {
	*load.Load
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
*/

func NewAboutPage2(l *load.Load) *AboutPage {
	versionRow := unclickableRow{
		leftWidget:  l.Theme.Body1("Version"),
		rightWidget: l.Theme.Body1("v1.5.2"),
	}

	buildDateRow := unclickableRow{
		leftWidget:  l.Theme.Body1("Build date"),
		rightWidget: l.Theme.Body1("2020-09-10"),
	}

	networkRow := unclickableRow{
		leftWidget:  l.Theme.Body1("Network"),
		rightWidget: l.Theme.Body1(l.WL.Wallet.Net),
	}

	licenseRow := clickableLicense{
		clickable:   new(widget.Clickable),
		leftWidget:  l.Theme.Body1("License"),
		rightWidget: l.Icons.ChevronRight,
	}
	pg := &AboutPage{
		Load:         l,
		card:         l.Theme.Card(),
		container:    &layout.List{Axis: layout.Vertical},
		versionRow:   versionRow,
		buildDateRow: buildDateRow,
		networkRow:   networkRow,
		licenseRow:   licenseRow,
	}

	pg.backButton, _ = components.SubpageHeaderButtons(l)
	pg.versionRow.rightWidget.Color = pg.Theme.Color.Gray
	pg.buildDateRow.rightWidget.Color = pg.Theme.Color.Gray
	pg.networkRow.rightWidget.Color = pg.Theme.Color.Gray
	pg.licenseRow.rightWidget.Color = pg.Theme.Color.Gray

	return pg
}

func (pg *AboutPage) ID() string {
	return AboutPageID
}

func (pg *AboutPage) OnResume() {

}

func (pg *AboutPage) Layout(gtx layout.Context) layout.Dimensions {
	body := func(gtx C) D {
		page := components.SubPage{
			Load:       pg.Load,
			Title:      "About",
			BackButton: pg.backButton,
			Back: func() {
				pg.PopFragment()
			},
			Body: func(gtx C) D {
				return pg.card.Layout(gtx, func(gtx C) D {
					return pg.Layout(gtx)
				})
			},
		}
		return page.Layout(gtx)
	}

	return components.UniformPadding(gtx, body)
}

/*func (pg *AboutPage) layoutRows(gtx layout.Context) layout.Dimensions {
	w := []func(gtx C) D{
		func(gtx C) D {
			return components.EndToEndRow(gtx, pg.version.Layout, pg.versionValue.Layout)
		},
		func(gtx C) D {
			return components.EndToEndRow(gtx, pg.buildDate.Layout, pg.buildDateValue.Layout)
		},
		func(gtx C) D {
			return components.EndToEndRow(gtx, pg.network.Layout, pg.networkValue.Layout)
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
					return components.Container{
						Padding: layout.Inset{
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
					}.Layout(gtx, pg.Theme.Separator().Layout)
				}),
			)
		})
	})
}
*/

func (pg *AboutPage) Handle() {
	/*if pg.licenseRow.Clicked() {
		pg.ChangeFragment(NewLicensePage(pg.Load))
	}*/
}

func (pg *AboutPage) OnClose() {}
